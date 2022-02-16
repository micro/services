package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	"github.com/micro/services/pkg/api"
	"github.com/micro/services/pkg/auth"
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
	rsp.LogoUrl = v.getLogo(rsp.Make)
	return nil
}

const (
	prefixLogo = "logo"
)

func (v *Vehicle) getLogo(make string) string {
	recs, err := store.Read(logoKey(make))
	if err != nil {
		if err == store.ErrNotFound {
			logger.Errorf("No logo found for make %s", make)
		} else {
			logger.Errorf("Error reading logo %s", err)
		}
		return ""
	}
	return string(recs[0].Value)
}

func (v *Vehicle) SetLogo(ctx context.Context, request *pb.SetLogoRequest, response *pb.SetLogoResponse) error {
	if _, err := auth.VerifyMicroAdmin(ctx, "vehicleAdmin.SetLogo"); err != nil {
		return err
	}
	rec := store.NewRecord(logoKey(request.Make), request.Url)
	if err := store.Write(rec); err != nil {
		return err
	}
	return nil
}

func logoKey(make string) string {
	return fmt.Sprintf("%s/%s", prefixLogo, strings.ToLower(make))
}
