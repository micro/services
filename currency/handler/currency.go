package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/currency/proto"
	"github.com/patrickmn/go-cache"
)

type Currency struct {
	Api   string
	Cache *cache.Cache
}

func (c *Currency) Rates(ctx context.Context, req *pb.RatesRequest, rsp *pb.RatesResponse) error {
	if len(req.Code) == 0 {
		return errors.BadRequest("currency.rates", "missing code")
	}
	if len(req.Code) != 3 {
		return errors.BadRequest("currency.rates", "code is invalid")
	}

	// try the cache
	if rates, ok := c.Cache.Get("rates:" + req.Code); ok {
		rsp.Code = req.Code
		rsp.Rates = rates.(map[string]float64)
		return nil
	}

	resp, err := http.Get(c.Api + "/latest/" + req.Code)
	if err != nil {
		logger.Errorf("Failed to get rates: %v\n", err)
		return errors.InternalServerError("currency.rates", "failed to get rates")
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to get rates (non 200): %d %v\n", resp.StatusCode, string(b))
		return errors.InternalServerError("currency.rates", "failed to get rates")
	}

	var respBody map[string]interface{}

	if err := json.Unmarshal(b, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal rates: %v\n", err)
		return errors.InternalServerError("currency.rates", "failed to get rates")
	}

	rates, ok := respBody["conversion_rates"].(map[string]interface{})
	if !ok {
		logger.Errorf("Failed to convert rates to map[string]interface{}: %v\n", ok)
		return errors.InternalServerError("currency.rates", "failed to get rates")
	}

	rsp.Code = req.Code
	rsp.Rates = make(map[string]float64)

	for code, rate := range rates {
		rsp.Rates[code], _ = rate.(float64)
	}

	// set for a period of time
	c.Cache.Set("rates:"+req.Code, rsp.Rates, cache.DefaultExpiration)

	return nil
}

func (c *Currency) Convert(ctx context.Context, req *pb.ConvertRequest, rsp *pb.ConvertResponse) error {
	if len(req.From) != 3 {
		return errors.BadRequest("currency.convert", "invalid from code")
	}
	if len(req.To) != 3 {
		return errors.BadRequest("currency.convert", "invalid to code")
	}

	uri := fmt.Sprintf("%s/pair/%s/%s", c.Api, req.From, req.To)

	// try the cache
	if req.Amount == 0 {
		rate, ok := c.Cache.Get("pair:" + req.From + req.To)
		if ok {
			rsp.From = req.From
			rsp.To = req.To
			rsp.Rate = rate.(float64)
			return nil
		}
	}

	if req.Amount > 0.0 {
		uri = fmt.Sprintf("%s/%v", uri, req.Amount)
	}

	resp, err := http.Get(uri)
	if err != nil {
		logger.Errorf("Failed to convert: %v\n", err)
		return errors.InternalServerError("currency.convert", "failed to convert")
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to get convert (non 200): %d %v\n", resp.StatusCode, string(b))
		return errors.InternalServerError("currency.convert", "failed to convert")
	}

	var respBody map[string]interface{}

	if err := json.Unmarshal(b, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal conversion: %v\n", err)
		return errors.InternalServerError("currency.convet", "failed to convert")
	}

	rsp.From = req.From
	rsp.To = req.To
	rsp.Rate, _ = respBody["conversion_rate"].(float64)
	rsp.Amount, _ = respBody["conversion_result"].(float64)

	// save for a period of time
	c.Cache.Set("pair:"+req.From+req.To, rsp.Rate, cache.DefaultExpiration)

	return nil
}
