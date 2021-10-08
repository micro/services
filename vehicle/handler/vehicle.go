package handler

import (
	"context"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/pkg/api"
	pb "github.com/micro/services/vehicle/proto"
)

var (
	apiURL = "https://driver-vehicle-licensing.api.gov.uk/vehicle-enquiry/v1/vehicles"
)

type Vehicle struct {
	Key string
}

func New(key string) *Vehicle {
	api.SetKey("x-api-key", key)

	return &Vehicle{
		Key: key,
	}
}

func (v *Vehicle) Lookup(ctx context.Context, req *pb.LookupRequest, rsp *pb.LookupResponse) error {
	if len(req.Registration) == 0 {
		return errors.BadRequest("vehicle.lookup", "missing registration number")
	}

	var resp map[string]interface{}

	if err := api.Post(apiURL, map[string]interface{}{
		"registrationNumber": req.Registration,
	}, &resp); err != nil {
		logger.Errorf("Failed to lookup vehicle %v: %v", req.Registration, err)
		return errors.InternalServerError("vehicle.lookup", "Failed to lookup vehicle")
	}

	rsp.Registration, _ = resp["registrationNumber"].(string)
	rsp.Make, _ = resp["make"].(string)
	rsp.Co2Emissions, _ = resp["co2Emissions"].(float64)
	rsp.Colour, _ = resp["colour"].(string)
	yom, _ := resp["yearOfManufacture"].(float64)
	rsp.YearOfManufacture = int32(yom)
	ec, _ := resp["engineCapacity"].(float64)
	rsp.EngineCapacity = int32(ec)
	rsp.FuelType, _ = resp["fuelType"].(string)
	rsp.MonthOfFirstRegistration, _ = resp["monthOfFirstRegistration"].(string)
	rsp.MotStatus, _ = resp["motStatus"].(string)

	if v := resp["motExpiryDate"]; v != nil {
		rsp.MotExpiry, _ = v.(string)
	}

	rsp.TaxDueDate, _ = resp["taxDueDate"].(string)
	rsp.TaxStatus, _ = resp["taxStatus"].(string)
	rsp.TypeApproval, _ = resp["typeApproval"].(string)
	rsp.Wheelplan, _ = resp["wheelplan"].(string)
	rsp.LastV5Issued, _ = resp["dateOfLastV5CIssued"].(string)

	return nil
}
