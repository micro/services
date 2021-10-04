package handler

import (
	"context"

	"github.com/micro/services/pkg/api"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/vehicle/proto"
)

var (
	apiURL = "https://driver-vehicle-licensing.api.gov.uk/vehicle-enquiry/v1/vehicles"
)

type Vehicle struct{
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

	rsp.Registration = resp["registrationNumber"].(string)
	rsp.Make = resp["make"].(string)
	rsp.Co2Emissions = resp["co2Emissions"].(float64)
	rsp.Colour = resp["colour"].(string)
	rsp.YearOfManufacture = int32(resp["yearOfManufacture"].(float64))
	rsp.EngineCapacity = int32(resp["engineCapacity"].(float64))
	rsp.FuelType = resp["fuelType"].(string)
	rsp.MonthOfFirstRegistration = resp["monthOfFirstRegistration"].(string)
	rsp.MotStatus = resp["motStatus"].(string)

	if v := resp["motExpiryDate"]; v != nil {
		rsp.MotExpiry = v.(string)
	}

	rsp.TaxDueDate = resp["taxDueDate"].(string)
	rsp.TaxStatus = resp["taxStatus"].(string)
	rsp.TypeApproval = resp["typeApproval"].(string)
	rsp.Wheelplan = resp["wheelplan"].(string)
	rsp.LastV5Issued = resp["dateOfLastV5CIssued"].(string)

	return nil
}

