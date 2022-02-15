package handler

import (
	"context"
	"strings"
	"strconv"

	"googlemaps.github.io/maps"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"

	pb "github.com/micro/services/place/proto"
)

var (
	Key    string
	Url    = "https://maps.googleapis.com/maps/api/"
	Format = "json"
)

type Google struct {
	Maps *maps.Client
}

func NewGoogle() *Google {
        // Setup google maps
        c, err := config.Get("google.apikey")
        if err != nil {
                logger.Fatalf("Error loading config: %v", err)
        }
        apiKey := c.String("")
        if len(apiKey) == 0 {
                logger.Fatalf("Missing required config: google.apikey")
        }
        m, err := maps.NewClient(maps.WithAPIKey(apiKey))
        if err != nil {
                logger.Fatalf("Error configuring google maps client: %v", err)
        }

	return &Google{
		Maps: m,
	}
}

func (g *Google) Nearby(ctx context.Context, req *pb.NearbyRequest, rsp *pb.NearbyResponse) error {
	greq := &maps.NearbySearchRequest{}

	if len(req.Location) > 0 {
		parts := strings.Split(req.Location, ",")
		if len(parts) > 2 {
			return errors.BadRequest("place.nearby", "invalid location")
		}
		lat, _ := strconv.ParseFloat(parts[0], 64)
		lng, _ := strconv.ParseFloat(parts[1], 64)

		greq.Location = &maps.LatLng{
			Lat: lat,
			Lng: lng,
		}
	}

	if req.Radius == 0 {
		req.Radius = uint32(1000)
	}

	greq.Radius = uint(req.Radius)
	greq.Keyword = req.Keyword
	greq.Name = req.Name
	greq.OpenNow = req.OpenNow

	// set the place type
	// https://developers.google.com/maps/documentation/places/web-service/supported_types
	if len(req.Type) > 0 {
		pt, err := maps.ParsePlaceType(req.Type)
		if err != nil {
			return err
		}
		greq.Type = pt
	}

	resp, err := g.Maps.NearbySearch(ctx, greq)
	if err != nil {
		return err
	}

	for _, res := range resp.Results {
		var hours []string
		var openNow bool
		if res.OpeningHours != nil {
			hours = res.OpeningHours.WeekdayText
			if res.OpeningHours.OpenNow != nil {
				openNow = *res.OpeningHours.OpenNow
			}
		}

		rsp.Results = append(rsp.Results, &pb.Result{
			Address: res.FormattedAddress,
			Location: res.Geometry.Location.String(),
			Name: res.Name,
			IconUrl: res.Icon,
			Rating: float64(res.Rating),
			OpenNow: openNow,
			OpeningHours: hours,
			Vicinity: res.Vicinity,
			Types: res.Types,
		})
	}

	return nil
}

func (g *Google) Search(ctx context.Context, req *pb.SearchRequest, rsp *pb.SearchResponse) error {
	if len(req.Query) == 0 {
		return errors.BadRequest("place.search", "missing query")
	}

	greq := &maps.TextSearchRequest{}

	if len(req.Location) > 0 {
		parts := strings.Split(req.Location, ",")
		if len(parts) > 2 {
			return errors.BadRequest("place.search", "invalid location")
		}
		lat, _ := strconv.ParseFloat(parts[0], 64)
		lng, _ := strconv.ParseFloat(parts[1], 64)

		greq.Location = &maps.LatLng{
			Lat: lat,
			Lng: lng,
		}
	}

	if req.Radius == 0 {
		req.Radius = uint32(1000)
	}

	greq.Radius = uint(req.Radius)
	greq.OpenNow = req.OpenNow

	// set the place type
	// https://developers.google.com/maps/documentation/places/web-service/supported_types
	if len(req.Type) > 0 {
		pt, err := maps.ParsePlaceType(req.Type)
		if err != nil {
			return err
		}
		greq.Type = pt
	}

	resp, err := g.Maps.TextSearch(ctx, greq)
	if err != nil {
		return err
	}

	for _, res := range resp.Results {
		var hours []string
		var openNow bool
		if res.OpeningHours != nil {
			hours = res.OpeningHours.WeekdayText
			if res.OpeningHours.OpenNow != nil {
				openNow = *res.OpeningHours.OpenNow
			}
		}
		rsp.Results = append(rsp.Results, &pb.Result{
			Address: res.FormattedAddress,
			Location: res.Geometry.Location.String(),
			Name: res.Name,
			IconUrl: res.Icon,
			Rating: float64(res.Rating),
			OpenNow: openNow,
			OpeningHours: hours,
			Vicinity: res.Vicinity,
			Types: res.Types,
		})
	}

	return nil
}
