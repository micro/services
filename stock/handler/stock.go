package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/stock/proto"
	"github.com/patrickmn/go-cache"
)

var (
	re = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
)

type Stock struct {
	Api   string
	Key   string
	Cache *cache.Cache
}

type OrderBook struct {
	Symbol  string
	Results []*Order
}

type Order struct {
	Symbol string
	Ask    float64
	Bid    float64
	Asize  int32
	Bsize  int32
	T      int64
}

type Quote struct {
	Symbol    string
	Ask       float64
	Bid       float64
	Asize     int32
	Bsize     int32
	Timestamp int64
}

type History struct {
	Symbol string
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int32
	From   string
}

func (s *Stock) History(ctx context.Context, req *pb.HistoryRequest, rsp *pb.HistoryResponse) error {
	if len(req.Stock) <= 0 || len(req.Stock) > 5 {
		return errors.BadRequest("stock.history", "invalid symbol")
	}

	if !re.MatchString(req.Date) {
		return errors.BadRequest("stock.history", "invalid date")
	}

	uri := fmt.Sprintf("%shistory/stock/open-close?stock=%s&date=%s&apikey=%s", s.Api, req.Stock, req.Date, s.Key)

	resp, err := http.Get(uri)
	if err != nil {
		logger.Errorf("Failed to get history: %v\n", err)
		return errors.InternalServerError("stock.history", "failed to get history")
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to get history (non 200): %d %v\n", resp.StatusCode, string(b))
		return errors.InternalServerError("stock.history", "failed to get history")
	}

	var respBody History

	if err := json.Unmarshal(b, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal history: %v\n", err)
		return errors.InternalServerError("stock.history", "failed to get history")
	}

	rsp.Symbol = respBody.Symbol
	rsp.Open = respBody.Open
	rsp.Close = respBody.Close
	rsp.High = respBody.High
	rsp.Low = respBody.Low
	rsp.Date = respBody.From
	rsp.Volume = respBody.Volume
	return nil
}

func (s *Stock) Quote(ctx context.Context, req *pb.QuoteRequest, rsp *pb.QuoteResponse) error {
	if len(req.Symbol) <= 0 || len(req.Symbol) > 5 {
		return errors.BadRequest("stock.quote", "invalid symbol")
	}

	uri := fmt.Sprintf("%slast/stock/%s?apikey=%s", s.Api, req.Symbol, s.Key)

	resp, err := http.Get(uri)
	if err != nil {
		logger.Errorf("Failed to get quote: %v\n", err)
		return errors.InternalServerError("stock.quote", "failed to get quote")
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to get quote (non 200): %d %v\n", resp.StatusCode, string(b))
		return errors.InternalServerError("stock.quote", "failed to get quote")
	}

	var respBody Quote

	if err := json.Unmarshal(b, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal quote: %v\n", err)
		return errors.InternalServerError("stock.quote", "failed to get quote")
	}

	rsp.Symbol = respBody.Symbol
	rsp.AskPrice = respBody.Ask
	rsp.BidPrice = respBody.Bid
	rsp.AskSize = respBody.Asize
	rsp.BidSize = respBody.Bsize
	rsp.Timestamp = time.Unix(0, respBody.Timestamp*int64(time.Millisecond)).UTC().Format(time.RFC3339Nano)

	return nil
}

func (s *Stock) Price(ctx context.Context, req *pb.PriceRequest, rsp *pb.PriceResponse) error {
	if len(req.Symbol) <= 0 || len(req.Symbol) > 5 {
		return errors.BadRequest("stock.price", "invalid symbol")
	}

	uri := fmt.Sprintf("%slast/trade/stock/%s?apikey=%s", s.Api, req.Symbol, s.Key)

	resp, err := http.Get(uri)
	if err != nil {
		logger.Errorf("Failed to get price: %v\n", err)
		return errors.InternalServerError("stock.trade", "failed to get price")
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to get price (non 200): %d %v\n", resp.StatusCode, string(b))
		return errors.InternalServerError("stock.price", "failed to get price")
	}

	var respBody map[string]interface{}

	if err := json.Unmarshal(b, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal price: %v\n", err)
		return errors.InternalServerError("stock.price", "failed to get price")
	}

	rsp.Symbol = req.Symbol
	rsp.Price = respBody["price"].(float64)

	return nil
}
