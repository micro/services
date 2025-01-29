package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/micro/micro/v5/service/config"
	"github.com/micro/micro/v5/service/errors"
	"github.com/micro/micro/v5/service/logger"
	pb "github.com/micro/services/forex/proto"
	"github.com/patrickmn/go-cache"
)

var (
	re = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
)

type Forex struct {
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

func New() *Forex {
	// TODO: look for "forex.provider" to determine the handler
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

	return &Forex{
		Api:   api,
		Key:   key,
		Cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (s *Forex) History(ctx context.Context, req *pb.HistoryRequest, rsp *pb.HistoryResponse) error {
	if len(req.Symbol) <= 0 {
		return errors.BadRequest("forex.history", "invalid symbol")
	}

	uri := fmt.Sprintf("%sagg/forex/prev-close/%s?apikey=%s", s.Api, req.Symbol, s.Key)

	resp, err := http.Get(uri)
	if err != nil {
		logger.Errorf("Failed to get history: %v\n", err)
		return errors.InternalServerError("forex.history", "failed to get history")
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to get history (non 200): %d %v\n", resp.StatusCode, string(b))
		return errors.InternalServerError("forex.history", "failed to get history")
	}

	var respBody Previous

	if err := json.Unmarshal(b, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal history: %v\n", err)
		return errors.InternalServerError("forex.history", "failed to get history")
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
func (s *Forex) Quote(ctx context.Context, req *pb.QuoteRequest, rsp *pb.QuoteResponse) error {
	if len(req.Symbol) <= 0 {
		return errors.BadRequest("forex.quote", "invalid symbol")
	}

	uri := fmt.Sprintf("%slast/forex/%s?apikey=%s", s.Api, req.Symbol, s.Key)

	resp, err := http.Get(uri)
	if err != nil {
		logger.Errorf("Failed to get quote: %v\n", err)
		return errors.InternalServerError("forex.quote", "failed to get quote")
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to get quote (non 200): %d %v\n", resp.StatusCode, string(b))
		return errors.InternalServerError("forex.quote", "failed to get quote")
	}

	var respBody Quote

	if err := json.Unmarshal(b, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal quote: %v\n", err)
		return errors.InternalServerError("forex.quote", "failed to get quote")
	}

	rsp.Symbol = respBody.Symbol
	rsp.AskPrice = respBody.Ask
	rsp.BidPrice = respBody.Bid
	rsp.Timestamp = time.Unix(0, int64(respBody.Timestamp)*int64(time.Millisecond)).UTC().Format(time.RFC3339Nano)

	return nil
}

func (s *Forex) Price(ctx context.Context, req *pb.PriceRequest, rsp *pb.PriceResponse) error {
	if len(req.Symbol) <= 0 {
		return errors.BadRequest("forex.price", "invalid symbol")
	}

	uri := fmt.Sprintf("%slast/trade/forex/%s?apikey=%s", s.Api, req.Symbol, s.Key)

	resp, err := http.Get(uri)
	if err != nil {
		logger.Errorf("Failed to get price: %v\n", err)
		return errors.InternalServerError("forex.trade", "failed to get price")
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to get price (non 200): %d %v\n", resp.StatusCode, string(b))
		return errors.InternalServerError("forex.quote", "failed to get price")
	}

	var respBody map[string]interface{}

	if err := json.Unmarshal(b, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal price: %v\n", err)
		return errors.InternalServerError("forex.price", "failed to get price")
	}

	rsp.Symbol = req.Symbol
	rsp.Price = respBody["price"].(float64)

	return nil
}
