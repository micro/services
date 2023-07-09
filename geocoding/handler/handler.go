package handler

import (
	"context"
	"strings"

	"googlemaps.github.io/maps"
	"micro.dev/v4/service/errors"
	"micro.dev/v4/service/logger"

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

func (g *Geocoding) Lookup(ctx context.Context, req *pb.LookupRequest, rsp *pb.LookupResponse) error {
	// query google maps
	results, err := g.Maps.Geocode(ctx, &maps.GeocodingRequest{Address: toString(req)})
	if err != nil {
		logger.Errorf("Error geocoding: %v", err)
		return ErrDownstream
	}
	if len(results) == 0 {
		return ErrNoResults
	}

	rsp.Address = new(pb.Address)
	rsp.Location = new(pb.Location)

	// return the result
	serializeResult(results[0], rsp.Address, rsp.Location)

	return nil
}

// Reverse geocode an address
func (g *Geocoding) Reverse(ctx context.Context, req *pb.ReverseRequest, rsp *pb.ReverseResponse) error {
	// validate the request
	if req.Latitude == 0.0 {
		return ErrMissingLatitude
	}
	if req.Longitude == 0.0 {
		return ErrMissingLongitude
	}

	// query google maps
	results, err := g.Maps.ReverseGeocode(ctx, &maps.GeocodingRequest{
		LatLng: &maps.LatLng{Lat: req.Latitude, Lng: req.Longitude},
	})
	if err != nil {
		logger.Errorf("Error geocoding: %v", err)
		return ErrDownstream
	}
	if len(results) == 0 {
		return ErrNoResults
	}

	rsp.Address = new(pb.Address)
	rsp.Location = new(pb.Location)

	// return the result
	serializeResult(results[0], rsp.Address, rsp.Location)
	return nil
}

func toString(l *pb.LookupRequest) string {
	var comps []string
	for _, c := range []string{l.Address, l.City, l.Postcode, l.Country} {
		t := strings.TrimSpace(c)
		if len(t) > 0 {
			comps = append(comps, t)
		}
	}
	return strings.Join(comps, ", ")
}

func serializeResult(r maps.GeocodingResult, a *pb.Address, l *pb.Location) {
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
	l.Latitude = r.Geometry.Location.Lat
	l.Longitude = r.Geometry.Location.Lng
}
