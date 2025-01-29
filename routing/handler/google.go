package handler

import (
	"context"
	"fmt"
	"time"

	"googlemaps.github.io/maps"

	"github.com/micro/micro/v5/service/errors"
	"github.com/micro/micro/v5/service/logger"
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

func (r *Google) Directions(ctx context.Context, req *pb.DirectionsRequest, rsp *pb.DirectionsResponse) error {
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

	legs := routes[0].Legs

	var distance int
	var duration time.Duration

	for i, leg := range legs {
		var intersection []*pb.Intersection
		for _, waypoint := range leg.ViaWaypoint {
			intersection = append(intersection, &pb.Intersection{
				Location: &pb.Point{
					Latitude:  waypoint.Location.Lat,
					Longitude: waypoint.Location.Lng,
				},
			})
		}
		rsp.Directions = append(rsp.Directions, &pb.Direction{
			Name:          fmt.Sprintf("leg %d", i),
			Distance:      float64(leg.Distance.Meters),
			Duration:      float64(leg.Duration.Seconds()),
			Intersections: intersection,
		})
		distance += leg.Distance.Meters
		duration += leg.Duration
	}

	// total distance/duration

	// in meters
	rsp.Distance = float64(distance)
	// in seconds
	rsp.Duration = float64(duration.Seconds())

	return nil
}

// Calculate the ETAs for a route
func (r *Google) Eta(ctx context.Context, req *pb.EtaRequest, rsp *pb.EtaResponse) error {
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

	// construct the request
	resp, err := r.Maps.DistanceMatrix(ctx, &maps.DistanceMatrixRequest{
		Origins:       []string{pointToString(req.Origin)},
		Destinations:  []string{pointToString(req.Destination)},
		DepartureTime: "now",
		Units:         "UnitsMetric",
		Mode:          maps.TravelModeDriving,
	})
	if err != nil {
		return err
	}

	// check the correct number of elements (route segments) were returned
	// from the Google API
	if len(resp.Rows[0].Elements) != 1 {
		return errors.InternalServerError("routing.eta", "Invalid downstream response. Expected %v segments but got %v", 1, len(resp.Rows[0].Elements))
	}

	// set the response duration in seconds
	rsp.Duration = float64(resp.Rows[0].Elements[0].Duration.Seconds())

	return nil
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
