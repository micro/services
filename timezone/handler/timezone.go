package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/timezone/proto"
	"github.com/tkuchiki/go-timezone"
)

type Timezone struct {
	Api string
	Key string
	TZ  *timezone.Timezone
}

func New() *Timezone {
	// TODO: look for "weather.provider" to determine the handler
	v, err := config.Get("weatherapi.api")
	if err != nil {
		logger.Fatalf("weatherapi.api config not found: %v", err)
	}
	api := v.String("")
	if len(api) == 0 {
		logger.Fatal("weatherapi.api config not found")
	}
	v, err = config.Get("weatherapi.key")
	if err != nil {
		logger.Fatalf("weatherapi.key config not found: %v", err)
	}
	key := v.String("")
	if len(key) == 0 {
		logger.Fatal("weatherapi.key config not found")
	}

	return &Timezone{
		Api: api,
		Key: key,
		TZ:  timezone.New(),
	}
}

func (t *Timezone) Info(ctx context.Context, req *pb.InfoRequest, rsp *pb.InfoResponse) error {
	if len(req.Location) == 0 {
		return errors.BadRequest("timezone.info", "invalid location")
	}

	vals := url.Values{}
	vals.Set("key", t.Key)
	vals.Set("q", req.Location)

	resp, err := http.Get(t.Api + "timezone.json?" + vals.Encode())
	if err != nil {
		logger.Errorf("Failed to get timezone info: %v\n", err)
		return errors.InternalServerError("weather.current", "failed to get timezone info")
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to get timezone info (non 200): %d %v\n", resp.StatusCode, string(b))
		return errors.InternalServerError("weather.current", "failed to get timezone info")
	}

	var respBody map[string]interface{}

	if err := json.Unmarshal(b, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal current: %v\n", err)
		return errors.InternalServerError("weather.current", "failed to get current")
	}

	location := respBody["location"].(map[string]interface{})

	rsp.Location = location["name"].(string)
	rsp.Region = location["region"].(string)
	rsp.Country = location["country"].(string)
	rsp.Latitude = location["lat"].(float64)
	rsp.Longitude = location["lon"].(float64)
	rsp.Timezone = location["tz_id"].(string)
	rsp.LocalTime = location["localtime"].(string)

	loc, _ := time.LoadLocation(rsp.Timezone)
	ti := time.Now().In(loc)
	isDST := t.TZ.IsDST(ti)
	rsp.Abbreviation, _ = t.TZ.GetTimezoneAbbreviation(rsp.Timezone, isDST)
	rsp.DaylightSavings = isDST

	return nil
}
