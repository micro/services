package handler

import (
	"context"
	"fmt"
	"time"

	pb "etas/proto"

	"github.com/micro/micro/v3/service/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"googlemaps.github.io/maps"
)

type ETAs struct {
	Maps *maps.Client
}

// Calculate the ETAs for a route
func (e *ETAs) Calculate(ctx context.Context, req *pb.Route, rsp *pb.Response) error {
	// validate the request
	if req.Pickup == nil {
		return errors.BadRequest("etas.Calculate", "Missing pickup")
	}
	if len(req.Waypoints) == 0 {
		return errors.BadRequest("etas.Calculate", "One more more waypoints required")
	}
	if err := validatePoint(req.Pickup, "Pickup"); err != nil {
		return err
	}
	for i, p := range req.Waypoints {
		if err := validatePoint(p, fmt.Sprintf("Waypoint %v", i)); err != nil {
			return err
		}
	}

	// construct the request
	destinations := make([]string, len(req.Waypoints))
	for i, p := range req.Waypoints {
		destinations[i] = pointToCoords(p)
	}
	departureTime := "now"
	if req.StartTime != nil {
		departureTime = req.StartTime.String()
	}
	resp, err := e.Maps.DistanceMatrix(ctx, &maps.DistanceMatrixRequest{
		Origins:       []string{pointToCoords(req.Pickup)},
		Destinations:  destinations,
		DepartureTime: departureTime,
		Units:         "UnitsMetric",
		Mode:          maps.TravelModeDriving,
	})
	if err != nil {
		return err
	}

	// check the correct number of elements (route segments) were returned
	// from the Google API
	if len(resp.Rows[0].Elements) != len(destinations) {
		return errors.InternalServerError("etas.Calculate", "Invalid downstream response. Expected %v segments but got %v", len(destinations), len(resp.Rows[0].Elements))
	}

	// calculate the response
	currentTime := time.Now()
	if req.StartTime != nil {
		currentTime = req.StartTime.AsTime()
	}
	rsp.Points = make(map[string]*pb.ETA, len(req.Waypoints)+1)
	for i, p := range append([]*pb.Point{req.Pickup}, req.Waypoints...) {
		at := currentTime
		if i > 0 {
			at = at.Add(resp.Rows[0].Elements[i-1].Duration)
		}
		et := at.Add(time.Minute * time.Duration(p.WaitTime))

		rsp.Points[p.Id] = &pb.ETA{
			EstimatedArrivalTime:   timestamppb.New(at),
			EstimatedDepartureTime: timestamppb.New(et),
		}

		currentTime = et
	}

	return nil
}

func validatePoint(p *pb.Point, desc string) error {
	if len(p.Id) == 0 {
		return errors.BadRequest("etas.Calculate", "%v missing ID", desc)
	}
	if p.Latitude == 0 {
		return errors.BadRequest("etas.Calculate", "%v missing Latitude", desc)
	}
	if p.Longitude == 0 {
		return errors.BadRequest("etas.Calculate", "%v missing Longitude", desc)
	}
	return nil
}

func pointToCoords(p *pb.Point) string {
	return fmt.Sprintf("%v,%v", p.Latitude, p.Longitude)
}
