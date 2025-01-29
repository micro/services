package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/micro/micro/v5/service/config"
	"github.com/micro/micro/v5/service/errors"
	"github.com/micro/micro/v5/service/logger"
	log "github.com/micro/micro/v5/service/logger"
	"github.com/micro/micro/v5/service/store"

	pb "github.com/micro/services/holidays/proto"
)

type Holidays struct {
	conf conf
}

type conf struct {
	NagerHost string `json:"nager_host"`
}

func New() *Holidays {
	val, err := config.Get("holidays")
	if err != nil {
		log.Fatalf("Failed to load config %s", err)
	}
	var conf conf
	if err := val.Scan(&conf); err != nil {
		log.Fatalf("Failed to load config %s", err)
	}
	if len(conf.NagerHost) == 0 {
		log.Fatalf("Nager host not configured")
	}
	return &Holidays{conf: conf}
}

type nagerCountry struct {
	CountryCode string `json:"countryCode"`
	Name        string `json:"name"`
}

func (h *Holidays) Countries(ctx context.Context, request *pb.CountriesRequest, response *pb.CountriesResponse) error {
	key := "countries"
	recs, err := store.Read(key, store.ReadLimit(1))
	if err == nil && len(recs) == 1 {
		var countries []*pb.Country
		if err := recs[0].Decode(&countries); err == nil {
			response.Countries = countries
			return nil
		} else {
			logger.Errorf("Failed to get countries from store: %v", err)
		}
	}

	rsp, err := http.Get(h.conf.NagerHost + "/api/v3/AvailableCountries")
	if err != nil {
		log.Errorf("Error listing available countries %s", err)
		return errors.InternalServerError("holidays.countries", "Error retrieving country list")
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != 200 {
		log.Errorf("Error listing available countries %s", rsp.Status)
		return errors.InternalServerError("holidays.countries", "Error retrieving country list")
	}
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Errorf("Error processing available countries %s", err)
		return errors.InternalServerError("holidays.countries", "Error retrieving country list")
	}
	var rspArr []nagerCountry
	if err := json.Unmarshal(b, &rspArr); err != nil {
		log.Errorf("Error processing available countries %s", err)
		return errors.InternalServerError("holidays.countries", "Error retrieving country list")
	}

	var countries []*pb.Country

	for _, c := range rspArr {
		countries = append(countries, &pb.Country{
			Code: c.CountryCode,
			Name: c.Name,
		})
	}

	// save the countries
	store.Write(store.NewRecord(key, countries))

	// set response
	response.Countries = countries

	return nil
}

type nagerHoliday struct {
	Date        string   `json:"date"`
	LocalName   string   `json:"localName"`
	Name        string   `json:"name"`
	CountryCode string   `json:"countryCode"`
	Counties    []string `json:"counties"`
	Types       []string `json:"types"`
}

func (h Holidays) List(ctx context.Context, request *pb.ListRequest, response *pb.ListResponse) error {
	if request.Year == 0 {
		return errors.BadRequest("holidays.list", "Missing year argument")
	}
	if len(request.CountryCode) == 0 {
		return errors.BadRequest("holidays.list", "Missing country code argument")
	}

	key := fmt.Sprintf("holidays/%d/%s", request.Year, request.CountryCode)

	recs, err := store.Read(key, store.ReadLimit(1))
	if err == nil && len(recs) == 1 {
		var holidays []*pb.Holiday
		if err := recs[0].Decode(&holidays); err == nil {
			response.Holidays = holidays
			return nil
		} else {
			logger.Errorf("Failed to get holidays from store: %v", err)
		}
	}

	rsp, err := http.Get(fmt.Sprintf("%s/api/v3/PublicHolidays/%d/%s", h.conf.NagerHost, request.Year, request.CountryCode))
	if err != nil {
		log.Errorf("Error listing available countries %s", err)
		return errors.InternalServerError("holidays.list", "Error retrieving holidays list")
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != 200 {
		log.Errorf("Error listing holidays %s", rsp.Status)
		return errors.InternalServerError("holidays.list", "Error retrieving holidays list")
	}
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Errorf("Error processing holidays %s", err)
		return errors.InternalServerError("holidays.countries", "Error retrieving holidays list")
	}
	var rspArr []nagerHoliday
	if err := json.Unmarshal(b, &rspArr); err != nil {
		log.Errorf("Error processing holidays %s", err)
		return errors.InternalServerError("holidays.countries", "Error retrieving holidays list")
	}

	var holidays []*pb.Holiday

	for _, c := range rspArr {
		holidays = append(holidays, &pb.Holiday{
			Date:        c.Date,
			Name:        c.Name,
			LocalName:   c.LocalName,
			CountryCode: c.CountryCode,
			Regions:     c.Counties,
			Types:       c.Types,
		})
	}

	// save holidays
	store.Write(store.NewRecord(key, holidays))

	// write response
	response.Holidays = holidays

	return nil
}
