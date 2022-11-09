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
	pb "github.com/micro/services/time/proto"
	"github.com/tkuchiki/go-timezone"
)

type Time struct {
	Api string
	Key string
	TZ  *timezone.Timezone
}

func New() *Time {
	// TODO: look for "weather.provider" to determine the handler
	v, err := config.Get("time.api")
	if err != nil {
		logger.Fatalf("time.api config not found: %v", err)
	}
	api := v.String("")
	if len(api) == 0 {
		logger.Fatal("time.api config not found")
	}
	v, err = config.Get("time.key")
	if err != nil {
		logger.Fatalf("time.key config not found: %v", err)
	}
	key := v.String("")
	if len(key) == 0 {
		logger.Fatal("time.key config not found")
	}

	return &Time{
		Api: api,
		Key: key,
		TZ:  timezone.New(),
	}
}

func (t *Time) Now(ctx context.Context, req *pb.NowRequest, rsp *pb.NowResponse) error {
	if len(req.Location) == 0 {
		ti := time.Now().In(time.UTC)
		rsp.Localtime = ti.Format("15:04:05")
		rsp.Timestamp = ti.Format(time.RFC3339Nano)
		rsp.Location = "Prime Meridian"
		rsp.Timezone = "UTC"
		rsp.Unix = ti.Unix()
		return nil
	}

	vals := url.Values{}
	vals.Set("key", t.Key)
	vals.Set("q", req.Location)

	resp, err := http.Get(t.Api + "timezone.json?" + vals.Encode())
	if err != nil {
		logger.Errorf("Failed to get time zone: %v\n", err)
		return errors.InternalServerError("time.now", "failed to get time")
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to get time zone (non 200): %d %v\n", resp.StatusCode, string(b))
		return errors.InternalServerError("time.now", "failed to get time")
	}

	var respBody map[string]interface{}

	if err := json.Unmarshal(b, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal current: %v\n", err)
		return errors.InternalServerError("time.now", "failed to get time")
	}

	location := respBody["location"].(map[string]interface{})

	loc, _ := time.LoadLocation(location["tz_id"].(string))
	ti := time.Now().In(loc)
	isDST := t.TZ.IsDST(ti)

	rsp.Localtime = ti.Format("15:04:05")
	rsp.Timezone, _ = t.TZ.GetTimezoneAbbreviation(location["tz_id"].(string), isDST)
	rsp.Location = ti.Location().String()
	rsp.Timestamp = ti.Format(time.RFC3339Nano)
	rsp.Unix = ti.Unix()

	return nil
}

func (t *Time) Zone(ctx context.Context, req *pb.ZoneRequest, rsp *pb.ZoneResponse) error {
	if len(req.Location) == 0 {
		return errors.BadRequest("time.zone", "invalid location")
	}

	vals := url.Values{}
	vals.Set("key", t.Key)
	vals.Set("q", req.Location)

	resp, err := http.Get(t.Api + "timezone.json?" + vals.Encode())
	if err != nil {
		logger.Errorf("Failed to get time zone: %v\n", err)
		return errors.InternalServerError("time.zone", "failed to get time zone")
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to get time zone (non 200): %d %v\n", resp.StatusCode, string(b))
		return errors.InternalServerError("time.zone", "failed to get time zone")
	}

	var respBody map[string]interface{}

	if err := json.Unmarshal(b, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal current: %v\n", err)
		return errors.InternalServerError("time.zone", "failed to get current")
	}

	location := respBody["location"].(map[string]interface{})

	rsp.Location = location["name"].(string)
	rsp.Region = location["region"].(string)
	rsp.Country = location["country"].(string)
	rsp.Latitude = location["lat"].(float64)
	rsp.Longitude = location["lon"].(float64)
	rsp.Timezone = location["tz_id"].(string)
	rsp.Localtime = location["localtime"].(string)
	loc, _ := time.LoadLocation(rsp.Timezone)
	ti := time.Now().In(loc)
	isDST := t.TZ.IsDST(ti)
	rsp.Abbreviation, _ = t.TZ.GetTimezoneAbbreviation(rsp.Timezone, isDST)
	rsp.Dst = isDST
	_, offset := ti.Zone()
	rsp.Offset = int32(offset / 60 / 60)
	return nil
}
