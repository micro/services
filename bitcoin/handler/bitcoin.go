package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/bitcoin/proto"
	"github.com/patrickmn/go-cache"
)

type Bitcoin struct {
	Api   string
	Key   string
	Cache *cache.Cache
}

func New() *Bitcoin {
	// TODO: look for "bitcoin.provider" to determine the handler
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

	return &Bitcoin{
		Api:   api,
		Key:   key,
		Cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (b *Bitcoin) Price(ctx context.Context, req *pb.PriceRequest, rsp *pb.PriceResponse) error {
	if len(req.Symbol) <= 0 {
		req.Symbol = "BTCUSD"
	}

	// try the cache first
	if price, ok := b.Cache.Get("price:" + req.Symbol); ok {
		rsp.Symbol = req.Symbol
		rsp.Price = price.(float64)
		return nil
	}

	// get the price
	uri := fmt.Sprintf("%slast/crypto/%s?apikey=%s", b.Api, req.Symbol, b.Key)

	resp, err := http.Get(uri)
	if err != nil {
		logger.Errorf("Failed to get price: %v\n", err)
		return errors.InternalServerError("bitcoin.price", "failed to get price")
	}
	defer resp.Body.Close()

	buf, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to get price (non 200): %d %v\n", resp.StatusCode, string(buf))
		return errors.InternalServerError("bitcoin.price", "failed to get price")
	}

	var respBody map[string]interface{}

	if err := json.Unmarshal(buf, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal price: %v\n", err)
		return errors.InternalServerError("bitcoin.price", "failed to get price")
	}

	rsp.Symbol = req.Symbol
	rsp.Price = respBody["price"].(float64)

	// cache the price
	b.Cache.Set("price:"+req.Symbol, rsp.Price, time.Minute*5)

	return nil
}
