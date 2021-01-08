package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"googlemaps.github.io/maps"

	"github.com/micro/services/geocoding/handler"
	pb "github.com/micro/services/geocoding/proto"
)

const (
	validReponse = `{
		"results":[
			 {
					"address_components":[
						 {
								"long_name":"160",
								"types":[
									 "street_number"
								]
						 },
						 {
								"long_name":"Grays Inn Road",
								"types":[
									 "route"
								]
						 },
						 {
								"long_name":"Holborn",
								"types":[
									 "neighborhood"
								]
						 },
						 {
								"long_name":"Santa Clara County",
								"types":[
									 "administrative_area_level_2",
									 "political"
								]
						 },
						 {
								"long_name":"London",
								"types":[
									 "political_town"
								]
						 },
						 {
								"long_name":"United Kingdom",
								"types":[
									 "country",
									 "political"
								]
						 },
						 {
								"long_name":"WC1X 8ED",
								"types":[
									 "postal_code"
								]
						 }
					],
					"geometry":{
						 "location":{
								"lat":51.522214,
								"lng":-0.113565
						 }
					},
					"partial_math":false,
					"place_id":"ChIJ2eUgeAK6j4ARbn5u_wAGqWA",
					"types":[
						 "street_address"
					]
			 }
		],
		"status":"OK"
 }`
	noResultsReponse = `{
	"results": [],
	"status": "OK"
}`
)

func TestGeocoding(t *testing.T) {
	tt := []struct {
		Name         string
		ResponseBody string
		ResponseCode int
		MapQuery     string
		Error        error
		Address      *pb.Address
		Result       *pb.Address
	}{
		{
			Name:         "Invalid address",
			ResponseBody: noResultsReponse,
			ResponseCode: http.StatusOK,
			Address: &pb.Address{
				LineOne: "Foobar Street",
			},
			Error:    handler.ErrNoResults,
			MapQuery: "Foobar Street",
		},
		{
			Name:         "Valid address",
			ResponseBody: validReponse,
			ResponseCode: http.StatusOK,
			Address: &pb.Address{
				LineOne:  "160 Grays Inn Road",
				LineTwo:  "Holborn",
				Postcode: "wc1x8ed",
				Country:  "United Kingdom",
			},
			MapQuery: "160 Grays Inn Road, Holborn, wc1x8ed, United Kingdom",
		},
		{
			Name:         "Maps error",
			ResponseCode: http.StatusInternalServerError,
			Address: &pb.Address{
				LineOne: "Foobar Street",
			},
			Error:        handler.ErrDownstream,
			ResponseBody: "{}",
			MapQuery:     "Foobar Street",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			var query string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				query = r.URL.Query().Get("address")
				w.WriteHeader(tc.ResponseCode)
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				fmt.Fprintln(w, tc.ResponseBody)
			}))
			defer server.Close()

			m, err := maps.NewClient(maps.WithBaseURL(server.URL), maps.WithAPIKey("shh"))
			if err != nil {
				t.Fatal(err)
			}
			h := &handler.Geocoding{Maps: m}

			var rsp pb.Address
			err = h.Geocode(context.TODO(), tc.Address, &rsp)
			assert.Equal(t, tc.MapQuery, query)
			assert.Equal(t, tc.Error, err)

			if tc.Result != nil {
				assert.Equal(t, tc.Result.LineOne, rsp.LineOne)
				assert.Equal(t, tc.Result.LineTwo, rsp.LineTwo)
				assert.Equal(t, tc.Result.City, rsp.City)
				assert.Equal(t, tc.Result.Country, rsp.Country)
				assert.Equal(t, tc.Result.Postcode, rsp.Postcode)
			}
		})
	}

}

func TestReverseGeocoding(t *testing.T) {
	tt := []struct {
		Name         string
		ResponseBody string
		ResponseCode int
		Error        error
		Latitude     *wrapperspb.DoubleValue
		Longitude    *wrapperspb.DoubleValue
		Result       *pb.Address
	}{
		{
			Name:     "Missing longitude",
			Latitude: &wrapperspb.DoubleValue{Value: 51.522214},
			Error:    handler.ErrMissingLongitude,
		},
		{
			Name:      "Missing latitude",
			Longitude: &wrapperspb.DoubleValue{Value: -0.113565},
			Error:     handler.ErrMissingLatitude,
		},
		{
			Name:         "Invalid address",
			ResponseBody: noResultsReponse,
			ResponseCode: http.StatusOK,
			Latitude:     &wrapperspb.DoubleValue{Value: 999.999999},
			Longitude:    &wrapperspb.DoubleValue{Value: 999.999999},
			Error:        handler.ErrNoResults,
		},
		{
			Name:         "Valid address",
			ResponseBody: validReponse,
			ResponseCode: http.StatusOK,
			Latitude:     &wrapperspb.DoubleValue{Value: 51.522214},
			Longitude:    &wrapperspb.DoubleValue{Value: -0.113565},
			Result: &pb.Address{
				LineOne:  "160 Grays Inn Road",
				LineTwo:  "Holborn",
				Postcode: "WC1X 8ED",
				Country:  "United Kingdom",
			},
		},
		{
			Name:         "Maps error",
			Latitude:     &wrapperspb.DoubleValue{Value: 51.522214},
			Longitude:    &wrapperspb.DoubleValue{Value: -0.113565},
			ResponseCode: http.StatusInternalServerError,
			Error:        handler.ErrDownstream,
			ResponseBody: "{}",
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			var lat, lng string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if coords := strings.Split(string(r.URL.Query().Get("latlng")), ","); len(coords) == 2 {
					lat = coords[0]
					lng = coords[1]
				}

				w.WriteHeader(tc.ResponseCode)
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				fmt.Fprintln(w, tc.ResponseBody)
			}))
			defer server.Close()

			m, err := maps.NewClient(maps.WithBaseURL(server.URL), maps.WithAPIKey("shh"))
			if err != nil {
				t.Fatal(err)
			}
			h := &handler.Geocoding{Maps: m}

			var rsp pb.Address
			err = h.Reverse(context.TODO(), &pb.Coordinates{
				Latitude: tc.Latitude, Longitude: tc.Longitude,
			}, &rsp)
			assert.Equal(t, tc.Error, err)

			if tc.Latitude != nil && tc.Longitude != nil {
				assert.Equal(t, fmt.Sprintf("%f", tc.Latitude.Value), lat)
				assert.Equal(t, fmt.Sprintf("%f", tc.Longitude.Value), lng)
			}

			if tc.Result != nil {
				assert.Equal(t, tc.Result.LineOne, rsp.LineOne)
				assert.Equal(t, tc.Result.LineTwo, rsp.LineTwo)
				assert.Equal(t, tc.Result.City, rsp.City)
				assert.Equal(t, tc.Result.Country, rsp.Country)
				assert.Equal(t, tc.Result.Postcode, rsp.Postcode)
			}
		})
	}

}
