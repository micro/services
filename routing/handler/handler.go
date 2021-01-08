package handler

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/wrapperspb"
	"googlemaps.github.io/maps"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/routing/proto"
)

var (
	ErrDownstream         = errors.InternalServerError("ROUTING_ERROR", "Unable to connect to routing provider")
	ErrMissingOrigin      = errors.BadRequest("MISSING_ORIGIN", "Missing origin")
	ErrMissingDestination = errors.BadRequest("MISSING_DESTINATION", "Missing destination")
	ErrMissingLatitude    = errors.BadRequest("MISSING_LATITUDE", "Missing latitude")
	ErrMissingLongitude   = errors.BadRequest("MISSING_LONGITUDE", "Missing longitude")
	ErrNoRoutes           = errors.BadRequest("NO_ROUTES", "No routes found")
)

type Routing struct {
	Maps *maps.Client
}

func (r *Routing) Route(ctx context.Context, req *pb.RouteRequest, rsp *pb.RouteResponse) error {
	// validate the request
	if req.Origin == nil {
		return ErrMissingOrigin
	}
	if req.Destination == nil {
		return ErrMissingDestination
	}
	if err := validatePoint(req.Origin); err != nil {
		return err
	}
	if err := validatePoint(req.Destination); err != nil {
		return err
	}

	// query google maps
	routes, _, err := r.Maps.Directions(ctx, &maps.DirectionsRequest{
		Origin: pointToString(req.Origin), Destination: pointToString(req.Destination),
	})
	if err != nil {
		logger.Errorf("Error geocoding: %v. Origin: '%v', Destination: '%v'", err, pointToString(req.Origin), pointToString(req.Destination))
		return ErrDownstream
	}
	if len(routes) == 0 {
		return ErrNoRoutes
	}

	// decode the points
	points, err := routes[0].OverviewPolyline.Decode()
	if err != nil {
		logger.Errorf("Error decoding polyline: %v", err)
		return ErrDownstream
	}

	// return the result
	rsp.Waypoints = make([]*pb.Point, len(points))
	for i, p := range points {
		rsp.Waypoints[i] = &pb.Point{
			Latitude:  &wrapperspb.DoubleValue{Value: p.Lat},
			Longitude: &wrapperspb.DoubleValue{Value: p.Lng},
		}
	}
	return nil
}

func validatePoint(p *pb.Point) error {
	if p.Latitude == nil {
		return ErrMissingLatitude
	}
	if p.Longitude == nil {
		return ErrMissingLongitude
	}
	return nil
}

func pointToString(p *pb.Point) string {
	return fmt.Sprintf("%v,%v", p.Latitude.Value, p.Longitude.Value)
}
