package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	pb "github.com/micro/services/currency/proto"
	"github.com/patrickmn/go-cache"
	"micro.dev/v4/service/errors"
	"micro.dev/v4/service/logger"
)

var (
	re = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
)

type Currency struct {
	Api   string
	Cache *cache.Cache
}

func (c *Currency) Codes(ctx context.Context, req *pb.CodesRequest, rsp *pb.CodesResponse) error {
	// try the cache
	if codes, ok := c.Cache.Get("codes"); ok {
		rsp.Codes = codes.([]*pb.Code)
		return nil
	}

	resp, err := http.Get(c.Api + "/codes")
	if err != nil {
		logger.Errorf("Failed to get codes: %v\n", err)
		return errors.InternalServerError("currency.codes", "failed to get codes")
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to get codes (non 200): %d %v\n", resp.StatusCode, string(b))
		return errors.InternalServerError("currency.codes", "failed to get codes")
	}

	var respBody map[string]interface{}

	if err := json.Unmarshal(b, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal codes: %v\n", err)
		return errors.InternalServerError("currency.codes", "failed to get codes")
	}

	codes, ok := respBody["supported_codes"].([]interface{})
	if !ok {
		logger.Errorf("Failed to convert rates to map[string]interface{}: %v\n", ok)
		return errors.InternalServerError("currency.rates", "failed to get rates")
	}

	for _, code := range codes {
		c := code.([]interface{})
		rsp.Codes = append(rsp.Codes, &pb.Code{
			Name:     c[0].(string),
			Currency: c[1].(string),
		})
	}

	// set for a period of time
	c.Cache.Set("codes", rsp.Codes, time.Hour)

	return nil
}

func (c *Currency) History(ctx context.Context, req *pb.HistoryRequest, rsp *pb.HistoryResponse) error {
	if len(req.Code) == 0 {
		return errors.BadRequest("currency.rates", "missing code")
	}
	if len(req.Code) != 3 {
		return errors.BadRequest("currency.rates", "code is invalid")
	}

	if len(req.Date) == 0 {
		return errors.BadRequest("currency.history", "missing date")
	}

	if !re.MatchString(req.Date) {
		return errors.BadRequest("currency.history", "invalid date")
	}

	// try the cache
	if rates, ok := c.Cache.Get("history:" + req.Code + req.Date); ok {
		rsp.Code = req.Code
		rsp.Date = req.Date
		rsp.Rates = rates.(map[string]float64)
		return nil
	}

	parts := strings.Split(req.Date, "-")

	resp, err := http.Get(fmt.Sprintf("%s/history/%s/%s/%s/%s", c.Api, req.Code, parts[0], parts[1], parts[2]))
	if err != nil {
		logger.Errorf("Failed to get historic rates: %v\n", err)
		return errors.InternalServerError("currency.history", "failed to get history")
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to get historic rates (non 200): %d %v\n", resp.StatusCode, string(b))
		return errors.InternalServerError("currency.history", "failed to get history")
	}

	var respBody map[string]interface{}

	if err := json.Unmarshal(b, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal historic rates: %v\n", err)
		return errors.InternalServerError("currency.history", "failed to get history")
	}

	rates, ok := respBody["conversion_rates"].(map[string]interface{})
	if !ok {
		logger.Errorf("Failed to convert historic rates to map[string]interface{}: %v\n", ok)
		return errors.InternalServerError("currency.history", "failed to get history")
	}

	rsp.Code = req.Code
	rsp.Date = req.Date
	rsp.Rates = make(map[string]float64)

	for code, rate := range rates {
		rsp.Rates[code], _ = rate.(float64)
	}

	// set for a period of time
	c.Cache.Set("history:"+req.Code+req.Date, rsp.Rates, time.Hour*24)

	return nil
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
