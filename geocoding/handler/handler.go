package handler

import (
	"context"
	"strings"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"googlemaps.github.io/maps"

	pb "github.com/micro/services/geocoding/proto"
)

var (
	ErrDownstream       = errors.InternalServerError("MAP_ERROR", "Unable to connect to map provider")
	ErrNoResults        = errors.BadRequest("NO_RESULTS", "Unable to geocode address, no results found")
	ErrMissingLatitude  = errors.BadRequest("MISSING_LATITUDE", "Missing latitude")
	ErrMissingLongitude = errors.BadRequest("MISSING_LONGITUDE", "Missing longitude")
)

type Geocoding struct {
	Maps *maps.Client
}

// Geocode an address
func (g *Geocoding) Geocode(ctx context.Context, req *pb.Address, rsp *pb.Address) error {
	// query google maps
	results, err := g.Maps.Geocode(ctx, &maps.GeocodingRequest{Address: toString(req)})
	if err != nil {
		logger.Errorf("Error geocoding: %v", err)
		return ErrDownstream
	}
	if len(results) == 0 {
		return ErrNoResults
	}

	// return the result
	serializeResult(results[0], rsp)
	return nil
}

// Reverse geocode an address
func (g *Geocoding) Reverse(ctx context.Context, req *pb.Coordinates, rsp *pb.Address) error {
	// validate the request
	if req.Latitude == nil {
		return ErrMissingLatitude
	}
	if req.Longitude == nil {
		return ErrMissingLongitude
	}

	// query google maps
	results, err := g.Maps.ReverseGeocode(ctx, &maps.GeocodingRequest{
		LatLng: &maps.LatLng{Lat: req.Latitude.Value, Lng: req.Longitude.Value},
	})
	if err != nil {
		logger.Errorf("Error geocoding: %v", err)
		return ErrDownstream
	}
	if len(results) == 0 {
		return ErrNoResults
	}

	// return the result
	serializeResult(results[0], rsp)
	return nil
}

func toString(a *pb.Address) string {
	var comps []string
	for _, c := range []string{a.LineOne, a.LineTwo, a.City, a.Postcode, a.Country} {
		t := strings.TrimSpace(c)
		if len(t) > 0 {
			comps = append(comps, t)
		}
	}
	return strings.Join(comps, ", ")
}

func serializeResult(r maps.GeocodingResult, a *pb.Address) {
	var street, number string
	for _, c := range r.AddressComponents {
		for _, t := range c.Types {
			switch t {
			case "street_number":
				number = c.LongName
			case "route":
				street = c.LongName
			case "neighborhood":
				a.LineTwo = c.LongName
			case "country":
				a.Country = c.LongName
			case "postal_code":
				a.Postcode = c.LongName
			case "postal_town":
				a.City = c.LongName
			}
		}
	}

	a.LineOne = strings.Join([]string{number, street}, " ")
	a.Latitude = r.Geometry.Location.Lat
	a.Longitude = r.Geometry.Location.Lng
}
