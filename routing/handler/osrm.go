package handler

import (
	"context"
	"net/url"
	"time"

	"github.com/gojuno/go.osrm"
	"github.com/micro/micro/v3/service/errors"
	pb "github.com/micro/services/routing/proto"
	"github.com/paulmach/go.geo"
)

type OSRM struct {
	// api address
	Address string
	// osrm client
	Client *osrm.OSRM
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

	for _, routes := range resp.Routes {
		for _, leg := range routes.Legs {
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
	}

	return nil
}
