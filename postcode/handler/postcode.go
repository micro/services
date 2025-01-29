package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/micro/micro/v5/service/errors"
	"github.com/micro/micro/v5/service/logger"
	pb "github.com/micro/services/postcode/proto"
)

var (
	PostcodeAPI       = "https://api.postcodes.io/postcodes/"
	RandomPostcodeURL = "https://api.postcodes.io/random/postcodes"
)

type Postcode struct{}

func (e *Postcode) Lookup(ctx context.Context, req *pb.LookupRequest, rsp *pb.LookupResponse) error {
	if len(req.Postcode) == 0 {
		return errors.BadRequest("postcode.lookup", "missing postcode")
	}

	resp, err := http.Get(PostcodeAPI + req.Postcode)
	if err != nil {
		logger.Errorf("Failed to http call %v: %v", PostcodeAPI+req.Postcode, err.Error())
		return errors.BadRequest("postcode.lookup", "failed to lookup postcode")
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to http call (%v) %v: %v", resp.StatusCode, PostcodeAPI+req.Postcode, string(b))
		return errors.BadRequest("postcode.lookup", "failed to lookup postcode")
	}

	var response map[string]interface{}
	if err := json.Unmarshal(b, &response); err != nil {
		logger.Error("Failed to unmarshal %v: %v", PostcodeAPI+req.Postcode, err.Error())
		return errors.BadRequest("postcode.lookup", "failed to lookup postcode")
	}
	result := response["result"].(map[string]interface{})

	rsp.Postcode = result["postcode"].(string)
	rsp.Country = result["country"].(string)
	rsp.Region = result["region"].(string)
	rsp.Latitude = result["latitude"].(float64)
	rsp.Longitude = result["longitude"].(float64)
	rsp.District = result["admin_district"].(string)
	rsp.Ward = result["admin_ward"].(string)
	return nil
}

func (e *Postcode) Validate(ctx context.Context, req *pb.ValidateRequest, rsp *pb.ValidateResponse) error {
	if len(req.Postcode) == 0 {
		return errors.BadRequest("postcode.validate", "missing postcode")
	}

	resp, err := http.Get(PostcodeAPI + req.Postcode + "/validate")
	if err != nil {
		logger.Errorf("Failed to http call %v: %v", PostcodeAPI+req.Postcode+"/validate", err.Error())
		return errors.BadRequest("postcode.validate", "failed to validate postcode")
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to http call (%v) %v: %v", resp.StatusCode, PostcodeAPI+req.Postcode+"/validate", string(b))
		return errors.BadRequest("postcode.validate", "failed to validate postcode")
	}

	var response map[string]interface{}
	if err := json.Unmarshal(b, &response); err != nil {
		logger.Error("Failed to unmarshal %v: %v", PostcodeAPI+req.Postcode+"/validate", err.Error())
		return errors.BadRequest("postcode.validate", "failed to validate postcode")
	}
	result := response["result"].(bool)
	rsp.Valid = result
	return nil
}

func (e *Postcode) Random(ctx context.Context, req *pb.RandomRequest, rsp *pb.RandomResponse) error {
	resp, err := http.Get(RandomPostcodeURL)
	if err != nil {
		logger.Errorf("Failed to http call %v: %v", RandomPostcodeURL, err.Error())
		return errors.BadRequest("postcode.random", "failed to lookup random postcode")
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to http call (%v) %v: %v", resp.StatusCode, RandomPostcodeURL, string(b))
		return errors.BadRequest("postcode.random", "failed to lookup rabdom postcode")
	}

	var response map[string]interface{}
	if err := json.Unmarshal(b, &response); err != nil {
		logger.Error("Failed to unmarshal %v: %v", RandomPostcodeURL, err.Error())
		return errors.BadRequest("postcode.random", "failed to lookup random postcode")
	}
	result := response["result"].(map[string]interface{})

	rsp.Postcode = result["postcode"].(string)
	rsp.Country = result["country"].(string)
	rsp.Region = result["region"].(string)
	rsp.Latitude = result["latitude"].(float64)
	rsp.Longitude = result["longitude"].(float64)
	rsp.District = result["admin_district"].(string)
	rsp.Ward = result["admin_ward"].(string)
	return nil
}
