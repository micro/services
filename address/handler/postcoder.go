package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/micro/micro/v5/service/errors"
	"github.com/micro/micro/v5/service/logger"
	"github.com/micro/micro/v5/service/store"
	pb "github.com/micro/services/address/proto"
	"github.com/micro/services/pkg/api"
)

var (
	PostcoderURL = "https://ws.postcoder.com/pcw/"
)

type Postcoder struct {
	Url string
	Key string
}

func (a *Postcoder) lookupPostcode(q string, rsp interface{}) error {
	u := fmt.Sprintf("%s/%s/address/uk/%s?format=json&lines=2", a.Url, a.Key, q)
	return api.Get(u, rsp)
}

func (a *Postcoder) LookupPostcode(ctx context.Context, req *pb.LookupPostcodeRequest, rsp *pb.LookupPostcodeResponse) error {
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

	var resp []interface{}

	// lookup the address api for the given postcode
	if err := a.lookupPostcode(req.Postcode, &resp); err != nil {
		logger.Errorf("Failed to lookup postcode %v: %v", req.Postcode, err)
		return errors.InternalServerError("address.lookup-postcode", "failed to lookup postcode")
	}

	for _, res := range resp {
		addr := res.(map[string]interface{})
		rsp.Addresses = append(rsp.Addresses, &pb.Record{
			LineOne:      field("addressline1", addr),
			LineTwo:      field("addressline2", addr),
			Summary:      field("summaryline", addr),
			Organisation: field("organisation", addr),
			BuildingName: field("buildingname", addr),
			Premise:      field("premise", addr),
			Street:       field("street", addr),
			Locality:     field("dependentlocality", addr),
			Town:         field("posttown", addr),
			County:       field("county", addr),
			Postcode:     field("postcode", addr),
		})
	}

	// cache the record if we haven't done so already
	store.Write(store.NewRecord(postcode, rsp.Addresses))

	return nil
}
