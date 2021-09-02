package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	pb "github.com/micro/services/address/proto"
	"github.com/micro/services/pkg/api"
)

var (
	AddressURL = "https://api.ideal-postcodes.co.uk/v1/"
)

type Address struct {
	Url string
	Key string
}

func field(key string, vals map[string]interface{}) string {
	if v, ok := vals[key].(string); ok {
		return v
	}
	return ""
}

func (a *Address) lookupPostcode(q string, rsp interface{}) error {
	u := fmt.Sprintf("%saddresses?api_key=%s&postcode=%s", a.Url, a.Key, q)
	return api.Get(u, rsp)
}

func (a *Address) LookupPostcode(ctx context.Context, req *pb.LookupPostcodeRequest, rsp *pb.LookupPostcodeResponse) error {
	if len(req.Postcode) == 0 {
		return errors.BadRequest("address.lookup-postcode", "missing postcode")
	}

	// check the store to see if we have it already
	postcode := strings.ToLower(req.Postcode)
	postcode = strings.ReplaceAll(postcode, " ", "")

	// check the store
	rec, err := store.Read(postcode)
	if err == nil && len(rec) > 0 {
		for _, r := range rec {
			var records []*pb.Record
			r.Decode(&records)
			rsp.Addresses = records
		}
		return nil
	}

	var resp map[string]interface{}

	// lookup the address api for the given postcode
	if err := a.lookupPostcode(postcode, &resp); err != nil {
		logger.Errorf("Failed to lookup postcode %v: %v", req.Postcode, err)
		return errors.InternalServerError("address.lookup-postcode", "failed to lookup postcode")
	}

	result := resp["result"].(map[string]interface{})
	hits := result["hits"].([]interface{})

	for _, res := range hits {
		addr := res.(map[string]interface{})
		rsp.Addresses = append(rsp.Addresses, &pb.Record{
			LineOne:      field("line_1", addr),
			LineTwo:      field("line_2", addr),
			Organisation: field("organisation_name", addr),
			BuildingName: field("building_name", addr),
			Premise:      field("premise", addr),
			Street:       field("thoroughfare", addr),
			Locality:     field("dependent_locality", addr),
			Town:         field("post_town", addr),
			County:       field("county", addr),
			Postcode:     field("postcode", addr),
		})
	}

	// cache the record if we haven't done so already
	store.Write(store.NewRecord(postcode, rsp.Addresses))

	return nil
}
