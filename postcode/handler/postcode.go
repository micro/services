package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/postcode/proto"
)

var (
	PostcodeAPI = "https://api.postcodes.io/postcodes/"
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
