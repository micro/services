package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/weather/proto"
)

type Weather struct {
	Api string
	Key string
}

func New() *Weather {
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

	return &Weather{
		Api: api,
		Key: key,
	}
}

func (w *Weather) Forecast(ctx context.Context, req *pb.ForecastRequest, rsp *pb.ForecastResponse) error {
	if len(req.Location) <= 0 {
		return errors.BadRequest("weather.forecast", "invalid location")
	}
	if req.Days <= 0 || req.Days > 10 {
		req.Days = 1
	}

	vals := url.Values{}
	vals.Set("key", w.Key)
	vals.Set("aqi", "no")
	vals.Set("q", req.Location)
	vals.Set("days", fmt.Sprintf("%d", req.Days))

	resp, err := http.Get(w.Api + "forecast.json?" + vals.Encode())
	if err != nil {
		logger.Errorf("Failed to get weather forecast: %v\n", err)
		return errors.InternalServerError("weather.forecast", "failed to get weather forecast")
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to get weather forecast (non 200): %d %v\n", resp.StatusCode, string(b))
		return errors.InternalServerError("weather.forecast", "failed to get weather forecast")
	}

	var respBody map[string]interface{}

	if err := json.Unmarshal(b, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal forecast: %v\n", err)
		return errors.InternalServerError("weather.forecast", "failed to get forecast")
	}

	forecast := respBody["forecast"].(map[string]interface{})

	for _, v := range forecast["forecastday"].([]interface{}) {
		fc := v.(map[string]interface{})

		day := fc["day"].(map[string]interface{})

		willrain := false
		chancerain := int32(0)
		if v := day["daily_will_it_rain"].(float64); v == 1.0 {
			willrain = true
		}

		// is this a string or a float64? Try both
		if dcr, ok := day["daily_chance_of_rain"].(string); ok {
			if v, _ := strconv.Atoi(dcr); v > 0 {
				chancerain = int32(v)
			}
		} else if dcr, ok := day["daily_chance_of_rain"].(float64); ok {
			chancerain = int32(dcr)
		}

		// set the daily forecast
		rsp.Forecast = append(rsp.Forecast, &pb.Forecast{
			Date:         fc["date"].(string),
			MaxTempC:     day["maxtemp_c"].(float64),
			MinTempC:     day["mintemp_c"].(float64),
			AvgTempC:     day["avgtemp_c"].(float64),
			MaxTempF:     day["maxtemp_f"].(float64),
			MinTempF:     day["mintemp_f"].(float64),
			AvgTempF:     day["avgtemp_f"].(float64),
			WillItRain:   willrain,
			ChanceOfRain: chancerain,
			Condition:    day["condition"].(map[string]interface{})["text"].(string),
			IconUrl:      day["condition"].(map[string]interface{})["icon"].(string),
			Sunrise:      fc["astro"].(map[string]interface{})["sunrise"].(string),
			Sunset:       fc["astro"].(map[string]interface{})["sunset"].(string),
			MaxWindMph:   day["maxwind_mph"].(float64),
			MaxWindKph:   day["maxwind_kph"].(float64),
		})
	}
	location := respBody["location"].(map[string]interface{})

	// set the location
	rsp.Location = location["name"].(string)
	rsp.Region = location["region"].(string)
	rsp.Country = location["country"].(string)
	rsp.Latitude = location["lat"].(float64)
	rsp.Longitude = location["lon"].(float64)
	rsp.Timezone = location["tz_id"].(string)
	rsp.LocalTime = location["localtime"].(string)

	return nil
}

func (w *Weather) Now(ctx context.Context, req *pb.NowRequest, rsp *pb.NowResponse) error {
	if len(req.Location) <= 0 {
		return errors.BadRequest("weather.current", "invalid location")
	}

	vals := url.Values{}
	vals.Set("key", w.Key)
	vals.Set("aqi", "no")
	vals.Set("q", req.Location)

	resp, err := http.Get(w.Api + "current.json?" + vals.Encode())
	if err != nil {
		logger.Errorf("Failed to get current weather: %v\n", err)
		return errors.InternalServerError("weather.current", "failed to get current weather")
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to get current weather (non 200): %d %v\n", resp.StatusCode, string(b))
		return errors.InternalServerError("weather.current", "failed to get current weather")
	}

	var respBody map[string]interface{}

	if err := json.Unmarshal(b, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal current: %v\n", err)
		return errors.InternalServerError("weather.current", "failed to get current")
	}

	location := respBody["location"].(map[string]interface{})
	current := respBody["current"].(map[string]interface{})

	// set the location
	rsp.Location = location["name"].(string)
	rsp.Region = location["region"].(string)
	rsp.Country = location["country"].(string)
	rsp.Latitude = location["lat"].(float64)
	rsp.Longitude = location["lon"].(float64)
	rsp.Timezone = location["tz_id"].(string)
	rsp.LocalTime = location["localtime"].(string)

	// set the time of day
	if current["is_day"].(float64) == 1.0 {
		rsp.Daytime = true
	}

	rsp.TempC = current["temp_c"].(float64)
	rsp.TempF = current["temp_f"].(float64)
	rsp.FeelsLikeC = current["feelslike_c"].(float64)
	rsp.FeelsLikeF = current["feelslike_f"].(float64)
	rsp.Humidity = int32(current["humidity"].(float64))
	rsp.Cloud = int32(current["cloud"].(float64))
	rsp.Condition = current["condition"].(map[string]interface{})["text"].(string)
	rsp.IconUrl = current["condition"].(map[string]interface{})["icon"].(string)
	rsp.WindMph = current["wind_mph"].(float64)
	rsp.WindKph = current["wind_kph"].(float64)
	rsp.WindDirection = current["wind_dir"].(string)
	rsp.WindDegree = int32(current["wind_degree"].(float64))

	return nil
}
