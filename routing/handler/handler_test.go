package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/micro/services/routing/handler"
	pb "github.com/micro/services/routing/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"googlemaps.github.io/maps"
)

const response = `{
		"routes" : [
			 {
					"bounds" : {
						 "northeast" : {
								"lat" : -33.8150985,
								"lng" : 151.2070825
						 },
						 "southwest" : {
								"lat" : -33.8770049,
								"lng" : 151.0031658
						 }
					},
					"overview_polyline": {
						"points" : "xvumEgs{y[V@|AH|@DdABbC@@?^@N?zD@\\?F@"
					},
					"copyrights" : "Map data ©2015 Google",
					"legs" : [
						 {
								"distance" : {
									 "text" : "23.8 km",
									 "value" : 23846
								},
								"duration" : {
									 "text" : "37 mins",
									 "value" : 2214
								},
								"end_address" : "Parramatta NSW, Australia",
								"end_location" : {
									 "lat" : -33.8150985,
									 "lng" : 151.0031658
								},
								"start_address" : "Sydney NSW, Australia",
								"start_location" : {
									 "lat" : -33.8674944,
									 "lng" : 151.2070825
								},
								"steps" : [
									 {
											"distance" : {
												 "text" : "0.4 km",
												 "value" : 366
											},
											"duration" : {
												 "text" : "2 mins",
												 "value" : 103
											},
											"end_location" : {
												 "lat" : -33.8707786,
												 "lng" : 151.206934
											},
											"html_instructions" : "Head \u003cb\u003esouth\u003c/b\u003e on \u003cb\u003eGeorge St\u003c/b\u003e toward \u003cb\u003eBarrack St\u003c/b\u003e",
											"polyline" : {
												 "points" : "xvumEgs{y[V@|AH|@DdABbC@@?^@N?zD@\\?F@"
											},
											"start_location" : {
												 "lat" : -33.8674944,
												 "lng" : 151.2070825
											},
											"transit_details" : {
												 "trip_short_name": "7108"
											},
											"travel_mode" : "DRIVING"
									 }
								],
								"via_waypoint" : []
						 }
					],
					"summary" : "A4 and M4"
			 }
		],
		"status" : "OK"
 }`

func TestRoute(t *testing.T) {
	var oLat, oLng, dLat, dLng string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if comps := strings.Split(r.URL.Query().Get("origin"), ","); len(comps) == 2 {
			oLat = comps[0]
			oLng = comps[1]
		} else {
			oLat = ""
		}

		if comps := strings.Split(r.URL.Query().Get("destination"), ","); len(comps) == 2 {
			dLat = comps[0]
			dLng = comps[1]
		} else {
			dLat = ""
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, response)
	}))
	defer server.Close()

	m, err := maps.NewClient(maps.WithBaseURL(server.URL), maps.WithAPIKey("shh"))
	if err != nil {
		t.Fatal(err)
	}
	h := &handler.Routing{Maps: m}

	originLat := &wrapperspb.DoubleValue{Value: 33.8688}
	originLng := &wrapperspb.DoubleValue{Value: 151.2093}

	destinationLat := &wrapperspb.DoubleValue{Value: 33.8136}
	destinationLng := &wrapperspb.DoubleValue{Value: 151.0034}

	tt := []struct {
		Name        string
		Origin      *pb.Point
		Destination *pb.Point
		Error       error
		Result      []*pb.Point
	}{
		{
			Name:        "MissingOrigin",
			Destination: &pb.Point{Latitude: destinationLat, Longitude: destinationLng},
			Error:       handler.ErrMissingOrigin,
		},
		{
			Name:   "MissingDestination",
			Origin: &pb.Point{Latitude: originLat, Longitude: originLng},
			Error:  handler.ErrMissingDestination,
		},
		{
			Name:        "MissingLatitude",
			Origin:      &pb.Point{Longitude: originLng},
			Destination: &pb.Point{Latitude: destinationLat, Longitude: destinationLng},
			Error:       handler.ErrMissingLatitude,
		},
		{
			Name:        "MissingLongitude",
			Origin:      &pb.Point{Latitude: originLat},
			Destination: &pb.Point{Latitude: destinationLat, Longitude: destinationLng},
			Error:       handler.ErrMissingLongitude,
		},
		{
			Name:        "Valid",
			Origin:      &pb.Point{Latitude: originLat, Longitude: originLng},
			Destination: &pb.Point{Latitude: destinationLat, Longitude: destinationLng},
			Result: []*pb.Point{
				{
					Latitude:  &wrapperspb.DoubleValue{Value: -33.867490000000004},
					Longitude: &wrapperspb.DoubleValue{Value: 151.20708000000002},
				},
				{
					Latitude:  &wrapperspb.DoubleValue{Value: -33.867610000000006},
					Longitude: &wrapperspb.DoubleValue{Value: 151.20707000000002},
				},
				{
					Latitude:  &wrapperspb.DoubleValue{Value: -33.868080000000006},
					Longitude: &wrapperspb.DoubleValue{Value: 151.20702},
				},
				{
					Latitude:  &wrapperspb.DoubleValue{Value: -33.868390000000005},
					Longitude: &wrapperspb.DoubleValue{Value: 151.20699000000002},
				},
				{
					Latitude:  &wrapperspb.DoubleValue{Value: -33.86874},
					Longitude: &wrapperspb.DoubleValue{Value: 151.20697},
				},
				{
					Latitude:  &wrapperspb.DoubleValue{Value: -33.869400000000006},
					Longitude: &wrapperspb.DoubleValue{Value: 151.20696},
				},
				{
					Latitude:  &wrapperspb.DoubleValue{Value: -33.86941},
					Longitude: &wrapperspb.DoubleValue{Value: 151.20696},
				},
				{
					Latitude:  &wrapperspb.DoubleValue{Value: -33.86957},
					Longitude: &wrapperspb.DoubleValue{Value: 151.20695},
				},
				{
					Latitude:  &wrapperspb.DoubleValue{Value: -33.86965},
					Longitude: &wrapperspb.DoubleValue{Value: 151.20695},
				},
				{
					Latitude:  &wrapperspb.DoubleValue{Value: -33.87059},
					Longitude: &wrapperspb.DoubleValue{Value: 151.20694},
				},
				{
					Latitude:  &wrapperspb.DoubleValue{Value: -33.870740000000005},
					Longitude: &wrapperspb.DoubleValue{Value: 151.20694},
				},
				{
					Latitude:  &wrapperspb.DoubleValue{Value: -33.87078},
					Longitude: &wrapperspb.DoubleValue{Value: 151.20693},
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			var rsp pb.RouteResponse
			err := h.Route(context.Background(), &pb.RouteRequest{
				Origin: tc.Origin, Destination: tc.Destination,
			}, &rsp)

			assert.Equal(t, tc.Error, err)
			if err != nil {
				return
			}

			// check the right info was sent to google maps
			if tc.Origin != nil && tc.Origin.Latitude != nil {
				assert.Equal(t, fmt.Sprintf("%v", tc.Origin.Latitude.Value), oLat)
			}
			if tc.Origin != nil && tc.Origin.Longitude != nil {
				assert.Equal(t, fmt.Sprintf("%v", tc.Origin.Longitude.Value), oLng)
			}
			if tc.Destination != nil && tc.Destination.Latitude != nil {
				assert.Equal(t, fmt.Sprintf("%v", tc.Destination.Latitude.Value), dLat)
			}
			if tc.Destination != nil && tc.Destination.Longitude != nil {
				assert.Equal(t, fmt.Sprintf("%v", tc.Destination.Longitude.Value), dLng)
			}

			// check the response is correct
			if len(tc.Result) != len(rsp.Waypoints) {
				t.Errorf("Incorrect number of waypoints returned, expected %v got %v", len(tc.Result), len(rsp.Waypoints))
			}
			for i, p := range tc.Result {
				w := rsp.Waypoints[i]
				assert.Equal(t, p.Latitude.Value, w.Latitude.Value)
				assert.Equal(t, p.Longitude.Value, w.Longitude.Value)
			}
		})
	}

}
