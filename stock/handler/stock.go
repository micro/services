package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/stock/proto"
)

type Stock struct{
	Api string
	Key string
	Cache *cache.Cache
}

type Quote struct {
	Symbol string
	Ask float64
	Bid float64
	Asize int32
	Bsize int32
	Timestamp int64
}

func (s *Stock) Quote(ctx context.Context, req *pb.QuoteRequest, rsp *pb.QuoteResponse) error {
	if len(req.Symbol) < 0 || len(req.Symbol) > 5 {
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
	rsp.Timestamp = time.Unix(0, respBody.Timestamp * int64(time.Millisecond)).UTC().Format(time.RFC3339Nano)

	return nil
}

func (s *Stock) Price(ctx context.Context, req *pb.PriceRequest, rsp *pb.PriceResponse) error {
	if len(req.Symbol) < 0 || len(req.Symbol) > 5 {
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
                return errors.InternalServerError("stock.quote", "failed to get price")
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
