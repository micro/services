package handler

import (
	"context"
	"fmt"

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
	ErrUnimplemented      = errors.InternalServerError("UNIMPLEMENTED", "endpoint is unimplemented")
)

type Google struct {
	Maps *maps.Client
}

func (r *Google) ETA(ctx context.Context, req *pb.ETARequest, rsp *pb.ETAResponse) error {
	// TODO: implement eta
	return ErrUnimplemented
}

func (r *Google) Route(ctx context.Context, req *pb.RouteRequest, rsp *pb.RouteResponse) error {
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
	rsp.Waypoints = make([]*pb.Waypoint, len(points))
	for i, p := range points {
		rsp.Waypoints[i] = &pb.Waypoint{
			Location: &pb.Point{
				Latitude:  p.Lat,
				Longitude: p.Lng,
			},
		}
	}
	return nil
}

func validatePoint(p *pb.Point) error {
	if p.Latitude == 0.0 {
		return ErrMissingLatitude
	}
	if p.Longitude == 0.0 {
		return ErrMissingLongitude
	}
	return nil
}

func pointToString(p *pb.Point) string {
	return fmt.Sprintf("%v,%v", p.Latitude, p.Longitude)
}
