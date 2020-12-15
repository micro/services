package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"etas/handler"
	pb "etas/proto"

	"googlemaps.github.io/maps"
)

func TestCalculate(t *testing.T) {
	// mock the API response from Google Maps
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		fmt.Fprintln(w, `{
			"rows": [
				{
					"elements": [
						{
							"duration": {
								"text": "10 mins",
								"value": 600
							},
							"status": "OK"
						},
						{
							"duration": {
								"text": "6 mins",
								"value": 360
							},
							"status": "OK"
						}
					]
				}
			],
			"status": "OK"
		}`)
	}))
	defer s.Close()
	m, err := maps.NewClient(maps.WithAPIKey("notrequired"), maps.WithBaseURL(s.URL))
	if err != nil {
		t.Fatal(err)
	}

	// construct the handler and test the response
	e := handler.ETAs{m}
	t.Run("MissingPickup", func(t *testing.T) {
		err := e.Calculate(context.TODO(), &pb.Route{
			Waypoints: []*pb.Point{
				&pb.Point{
					Id:        "shenfield-station",
					Latitude:  51.6308,
					Longitude: 0.3295,
				},
			},
		}, &pb.Response{})
		assert.Error(t, err)
	})

	t.Run("MissingWaypoints", func(t *testing.T) {
		err := e.Calculate(context.TODO(), &pb.Route{
			Pickup: &pb.Point{
				Id:        "shenfield-station",
				Latitude:  51.6308,
				Longitude: 0.3295,
			},
		}, &pb.Response{})
		assert.Error(t, err)
	})

	t.Run("Valid", func(t *testing.T) {
		st := time.Unix(1609459200, 0)

		var rsp pb.Response
		err := e.Calculate(context.TODO(), &pb.Route{
			StartTime: timestamppb.New(st),
			Pickup: &pb.Point{
				Id:        "shenfield-station",
				Latitude:  51.6308,
				Longitude: 0.3295,
				WaitTime:  5,
			},
			Waypoints: []*pb.Point{
				{
					Id:        "nandos",
					Latitude:  51.6199,
					Longitude: 0.2999,
					WaitTime:  10,
				},
				{
					Id:        "brentwood-station",
					Latitude:  51.6136,
					Longitude: 0.2996,
				},
			},
		}, &rsp)

		assert.NoError(t, err)
		assert.NotNilf(t, rsp.Points, "Points should be returned")

		p := rsp.Points["shenfield-station"]
		ea := st
		ed := ea.Add(time.Minute * 5)
		assert.True(t, p.EstimatedArrivalTime.AsTime().Equal(ea))
		assert.True(t, p.EstimatedDepartureTime.AsTime().Equal(ed))

		p = rsp.Points["nandos"]
		ea = ed.Add(time.Minute * 10) // drive time
		ed = ea.Add(time.Minute * 10) // wait time
		assert.True(t, p.EstimatedArrivalTime.AsTime().Equal(ea))
		assert.True(t, p.EstimatedDepartureTime.AsTime().Equal(ed))

		p = rsp.Points["brentwood-station"]
		ea = ed.Add(time.Minute * 6) // drive time
		ed = ea
		assert.True(t, p.EstimatedArrivalTime.AsTime().Equal(ea))
		assert.True(t, p.EstimatedDepartureTime.AsTime().Equal(ed))
	})
}
