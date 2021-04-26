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
		Annotations: osrm.AnnotationsTrue,
		Overview:    osrm.OverviewFalse,
		Geometries:  osrm.GeometriesPolyline6,
	})
	if err != nil {
		return errors.InternalServerError("routing.route", "failed to get route: %v", err.Error())
	}

	for _, waypoint := range resp.Waypoints {
		rsp.Waypoints = append(rsp.Waypoints, &pb.Waypoint{
			Name:     waypoint.Name,
			Distance: float64(waypoint.Distance),
			Location: &pb.Point{
				Latitude:  waypoint.Location.Lat(),
				Longitude: waypoint.Location.Lng(),
			},
		})
	}

	return nil
}
