package handler

import (
	"context"
	"net/url"
	"time"

	"github.com/gojuno/go.osrm"
	pb "github.com/micro/services/routing/proto"
	"github.com/paulmach/go.geo"
	"micro.dev/v4/service/errors"
)

type OSRM struct {
	// api address
	Address string
	// osrm client
	Client *osrm.OSRM
}

func (o *OSRM) Directions(ctx context.Context, req *pb.DirectionsRequest, rsp *pb.DirectionsResponse) error {
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

	orig := req.Origin
	dest := req.Destination

	if o.Client == nil {
		u, _ := url.Parse(o.Address)
		if u.Scheme == "" {
			u.Scheme = "http"
		}
		o.Client = osrm.NewFromURL(u.String())
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second)
	defer cancelFn()

	resp, err := o.Client.Route(ctx, osrm.RouteRequest{
		Profile: "car",
		Coordinates: osrm.NewGeometryFromPointSet(geo.PointSet{
			{orig.Longitude, orig.Latitude},
			{dest.Longitude, dest.Latitude},
		}),
		Steps:       osrm.StepsTrue,
		Annotations: osrm.AnnotationsFalse,
		Overview:    osrm.OverviewFalse,
	})
	if err != nil {
		return errors.InternalServerError("routing.eta", "failed to get route: %v", err.Error())
	}

	if len(resp.Routes) == 0 {
		return nil
	}

	// set the estimated duration and distance
	rsp.Duration = float64(resp.Routes[0].Duration)
	rsp.Distance = float64(resp.Routes[0].Distance)

	for _, leg := range resp.Routes[0].Legs {
		for _, step := range leg.Steps {
			// set the waypoints for the route
			for _, intersect := range step.Intersections {
				rsp.Waypoints = append(rsp.Waypoints, &pb.Waypoint{
					Name: step.Name,
					Location: &pb.Point{
						Latitude:  intersect.Location.Lat(),
						Longitude: intersect.Location.Lng(),
					},
				})
			}

			instruction := step.Maneuver.Modifier

			switch step.Maneuver.Type {
			case "new name":
				instruction = "go " + step.Maneuver.Modifier
			default:
				instruction = step.Maneuver.Type + " " + step.Maneuver.Modifier
			}

			var intersections []*pb.Intersection

			for _, is := range step.Intersections {
				var bearings []float64

				for _, bearing := range is.Bearings {
					bearings = append(bearings, float64(bearing))
				}

				intersections = append(intersections, &pb.Intersection{
					Bearings: bearings,
					Location: &pb.Point{
						Latitude:  is.Location.Lat(),
						Longitude: is.Location.Lng(),
					},
				})
			}

			action := step.Maneuver.Type
			if action == "new name" {
				action = "continue"
			}

			// set the directions for the route
			rsp.Directions = append(rsp.Directions, &pb.Direction{
				Name:        step.Name,
				Instruction: instruction,
				Distance:    float64(step.Distance),
				Duration:    float64(step.Duration),
				Maneuver: &pb.Maneuver{
					BearingBefore: float64(step.Maneuver.BearingBefore),
					BearingAfter:  float64(step.Maneuver.BearingAfter),
					Location: &pb.Point{
						Latitude:  float64(step.Maneuver.Location.Lat()),
						Longitude: float64(step.Maneuver.Location.Lng()),
					},
					Action:    action,
					Direction: step.Maneuver.Modifier,
				},
				Intersections: intersections,
			})
		}
	}

	return nil
}

func (o *OSRM) Eta(ctx context.Context, req *pb.EtaRequest, rsp *pb.EtaResponse) error {
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

	orig := req.Origin
	dest := req.Destination

	if o.Client == nil {
		u, _ := url.Parse(o.Address)
		if u.Scheme == "" {
			u.Scheme = "http"
		}
		o.Client = osrm.NewFromURL(u.String())
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second)
	defer cancelFn()

	resp, err := o.Client.Route(ctx, osrm.RouteRequest{
		Profile: "car",
		Coordinates: osrm.NewGeometryFromPointSet(geo.PointSet{
			{orig.Longitude, orig.Latitude},
			{dest.Longitude, dest.Latitude},
		}),
		Steps:       osrm.StepsFalse,
		Annotations: osrm.AnnotationsFalse,
		Overview:    osrm.OverviewFalse,
	})
	if err != nil {
		return errors.InternalServerError("routing.eta", "failed to get route: %v", err.Error())
	}

	if len(resp.Routes) == 0 {
		return nil
	}

	// distance is meters
	distance := resp.Routes[0].Distance
	// duration is seconds
	duration := resp.Routes[0].Duration

	// nothing to calculate
	if distance == 0.0 {
		return nil
	}

	// set the duration
	rsp.Duration = float64(duration)

	// TODO: calculate transport/speed

	return nil
}

func (o *OSRM) Route(ctx context.Context, req *pb.RouteRequest, rsp *pb.RouteResponse) error {
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

	if o.Client == nil {
		u, _ := url.Parse(o.Address)
		if u.Scheme == "" {
			u.Scheme = "http"
		}
		o.Client = osrm.NewFromURL(u.String())
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), time.Second)
	defer cancelFn()

	orig := req.Origin
	dest := req.Destination

	resp, err := o.Client.Route(ctx, osrm.RouteRequest{
		Profile: "car",
		Coordinates: osrm.NewGeometryFromPointSet(geo.PointSet{
			{orig.Longitude, orig.Latitude},
			{dest.Longitude, dest.Latitude},
		}),
		Steps:       osrm.StepsTrue,
		Annotations: osrm.AnnotationsFalse,
		Overview:    osrm.OverviewFalse,
	})
	if err != nil {
		return errors.InternalServerError("routing.route", "failed to get route: %v", err.Error())
	}

	if len(resp.Routes) == 0 {
		return nil
	}

	// set distance and duration
	rsp.Duration = float64(resp.Routes[0].Duration)
	rsp.Distance = float64(resp.Routes[0].Distance)

	for _, leg := range resp.Routes[0].Legs {
		for _, step := range leg.Steps {
			for _, intersect := range step.Intersections {
				rsp.Waypoints = append(rsp.Waypoints, &pb.Waypoint{
					Name: step.Name,
					Location: &pb.Point{
						Latitude:  intersect.Location.Lat(),
						Longitude: intersect.Location.Lng(),
					},
				})
			}
		}
	}

	return nil
}
