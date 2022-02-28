package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/crypto/proto"
	"github.com/patrickmn/go-cache"
)

var (
	re = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
)

type Article struct {
	Title       string
	Description string
	Url         string
	Source      string
	Date        string
}

type News struct {
	Ticker string
	Limit  int32
	News   []*Article
}

type Crypto struct {
	Api   string
	Key   string
	Cache *cache.Cache
}

type Quote struct {
	Symbol    string
	Ask       float64
	Bid       float64
	Asize     float64
	Bsize     float64
	Timestamp float64
}

type History struct {
	Open      float64 `json:"o"`
	High      float64 `json:"h"`
	Low       float64 `json:"l"`
	Close     float64 `json:"c"`
	Volume    float64 `json:"v"`
	Timestamp float64 `json:"t"`
}

type Previous struct {
	Symbol       string
	TotalResults int32
	Results      []*History
}

func New() *Crypto {
	// TODO: look for "crypto.provider" to determine the handler
	v, err := config.Get("finage.api")
	if err != nil {
		logger.Fatalf("finage.api config not found: %v", err)
	}
	api := v.String("")
	if len(api) == 0 {
		logger.Fatal("finage.api config not found")
	}
	v, err = config.Get("finage.key")
	if err != nil {
		logger.Fatalf("finage.key config not found: %v", err)
	}
	key := v.String("")
	if len(key) == 0 {
		logger.Fatal("finage.key config not found")
	}

	return &Crypto{
		Api:   api,
		Key:   key,
		Cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (c *Crypto) News(ctx context.Context, req *pb.NewsRequest, rsp *pb.NewsResponse) error {
	if len(req.Symbol) <= 0 {
		return errors.BadRequest("crypto.news", "invalid symbol")
	}

	uri := fmt.Sprintf("%snews/cryptocurrency/%s?apikey=%s", c.Api, req.Symbol, c.Key)

	resp, err := http.Get(uri)
	if err != nil {
		logger.Errorf("Failed to get news: %v\n", err)
		return errors.InternalServerError("crypto.news", "failed to get news")
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to get news (non 200): %d %v\n", resp.StatusCode, string(b))
		return errors.InternalServerError("crypto.news", "failed to get news")
	}

	var respBody News

	if err := json.Unmarshal(b, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal news: %v\n", err)
		return errors.InternalServerError("crypto.news", "failed to get news")
	}

	for _, article := range respBody.News {
		rsp.Articles = append(rsp.Articles, &pb.Article{
			Title:       article.Title,
			Description: article.Description,
			Url:         article.Url,
			Source:      article.Source,
			Date:        article.Date,
		})
	}

	return nil
}

func (c *Crypto) History(ctx context.Context, req *pb.HistoryRequest, rsp *pb.HistoryResponse) error {
	if len(req.Symbol) <= 0 {
		return errors.BadRequest("crypto.history", "invalid symbol")
	}

	uri := fmt.Sprintf("%sagg/crypto/prev-close/%s?apikey=%s", c.Api, req.Symbol, c.Key)

	resp, err := http.Get(uri)
	if err != nil {
		logger.Errorf("Failed to get history: %v\n", err)
		return errors.InternalServerError("crypto.history", "failed to get history")
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to get history (non 200): %d %v\n", resp.StatusCode, string(b))
		return errors.InternalServerError("crypto.history", "failed to get history")
	}

	var respBody Previous

	if err := json.Unmarshal(b, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal history: %v\n", err)
		return errors.InternalServerError("crypto.history", "failed to get history")
	}

	if len(respBody.Results) != 1 {
		return nil
	}

	res := respBody.Results[0]
	rsp.Symbol = req.Symbol
	rsp.Open = res.Open
	rsp.Close = res.Close
	rsp.High = res.High
	rsp.Low = res.Low
	rsp.Date = time.Unix(0, int64(res.Timestamp)*int64(time.Millisecond)).UTC().Format("2006-01-02")
	rsp.Volume = res.Volume

	return nil
}
func (c *Crypto) Quote(ctx context.Context, req *pb.QuoteRequest, rsp *pb.QuoteResponse) error {
	if len(req.Symbol) <= 0 {
		return errors.BadRequest("crypto.quote", "invalid symbol")
	}

	uri := fmt.Sprintf("%slast/quote/crypto/%s?apikey=%s", c.Api, req.Symbol, c.Key)

	resp, err := http.Get(uri)
	if err != nil {
		logger.Errorf("Failed to get quote: %v\n", err)
		return errors.InternalServerError("crypto.quote", "failed to get quote")
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to get quote (non 200): %d %v\n", resp.StatusCode, string(b))
		return errors.InternalServerError("crypto.quote", "failed to get quote")
	}

	var respBody Quote

	if err := json.Unmarshal(b, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal quote: %v\n", err)
		return errors.InternalServerError("crypto.quote", "failed to get quote")
	}

	rsp.Symbol = respBody.Symbol
	rsp.AskPrice = respBody.Ask
	rsp.BidPrice = respBody.Bid
	rsp.AskSize = respBody.Asize
	rsp.BidSize = respBody.Bsize
	rsp.Timestamp = time.Unix(0, int64(respBody.Timestamp)*int64(time.Millisecond)).UTC().Format(time.RFC3339Nano)

	return nil
}

func (c *Crypto) Price(ctx context.Context, req *pb.PriceRequest, rsp *pb.PriceResponse) error {
	if len(req.Symbol) <= 0 {
		return errors.BadRequest("crypto.price", "invalid symbol")
	}

	uri := fmt.Sprintf("%slast/crypto/%s?apikey=%s", c.Api, req.Symbol, c.Key)

	resp, err := http.Get(uri)
	if err != nil {
		logger.Errorf("Failed to get price: %v\n", err)
		return errors.InternalServerError("crypto.trade", "failed to get price")
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to get price (non 200): %d %v\n", resp.StatusCode, string(b))
		return errors.InternalServerError("crypto.quote", "failed to get price")
	}

	var respBody map[string]interface{}

	if err := json.Unmarshal(b, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal price: %v\n", err)
		return errors.InternalServerError("crypto.price", "failed to get price")
	}

	rsp.Symbol = req.Symbol
	rsp.Price = respBody["price"].(float64)

	return nil
}

func (c *Crypto) Symbols(ctx context.Context, req *pb.SymbolsRequest, rsp *pb.SymbolsResponse) error {
	cached, ok := c.Cache.Get("symbolsCache")
	if ok {
		rsp.Symbols, _ = cached.([]*pb.Symbol)
		return nil
	}

	toCache := []*pb.Symbol{}
	page := 1
	for {
		var symbolsRsp struct {
			Page      int32 `json:"page"`
			TotalPage int32 `json:"totalPage"`
			Symbols   []struct {
				Symbol string `json:"symbol"`
				Name   string `json:"name"`
			} `json:"symbols"`
		}

		uri := fmt.Sprintf("%ssymbol-list/crypto?page=%d&apikey=%s", c.Api, page, c.Key)

		resp, err := http.Get(uri)
		if err != nil {
			logger.Errorf("Failed to get price: %v\n", err)
			return errors.InternalServerError("crypto.trade", "failed to get price")
		}
		defer resp.Body.Close()

		b, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			logger.Errorf("Failed to load symbol list (non 200): %d %v\n", resp.StatusCode, string(b))
			return errors.InternalServerError("crypto.symbols", "failed to get symbols")
		}

		if err := json.Unmarshal(b, &symbolsRsp); err != nil {
			logger.Errorf("Error unmarshalling cyrpto symbols: %v\n", err)
			return errors.InternalServerError("crypto.symbols", "failed to get symbols")
		}

		for _, v := range symbolsRsp.Symbols {
			toCache = append(toCache, &pb.Symbol{
				Symbol: v.Symbol,
				Name:   v.Name,
			})
		}

		if symbolsRsp.Page == symbolsRsp.TotalPage {
			break
		}
		page++
	}
	c.Cache.Set("symbolsCache", toCache, 24*time.Hour)
	rsp.Symbols = toCache
	return nil

}
