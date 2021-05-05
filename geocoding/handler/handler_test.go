package handler_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

			var rsp pb.LookupResponse
			address := tc.Address.LineOne
			if len(tc.Address.LineTwo) > 0 {
				address = fmt.Sprintf("%s, %s", address, tc.Address.LineTwo)
			}
			err = h.Lookup(context.TODO(), &pb.LookupRequest{
				Address:  address,
				Postcode: tc.Address.Postcode,
				Country:  tc.Address.Country,
			}, &rsp)
			assert.Equal(t, tc.MapQuery, query)
			assert.Equal(t, tc.Error, err)

			if tc.Result != nil {
				assert.Equal(t, tc.Result.LineOne, rsp.Address.LineOne)
				assert.Equal(t, tc.Result.LineTwo, rsp.Address.LineTwo)
				assert.Equal(t, tc.Result.City, rsp.Address.City)
				assert.Equal(t, tc.Result.Country, rsp.Address.Country)
				assert.Equal(t, tc.Result.Postcode, rsp.Address.Postcode)
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
		Latitude     float64
		Longitude    float64
		Result       *pb.Address
	}{
		{
			Name:     "Missing longitude",
			Latitude: 51.522214,
			Error:    handler.ErrMissingLongitude,
		},
		{
			Name:      "Missing latitude",
			Longitude: -0.113565,
			Error:     handler.ErrMissingLatitude,
		},
		{
			Name:         "Invalid address",
			ResponseBody: noResultsReponse,
			ResponseCode: http.StatusOK,
			Latitude:     999.999999,
			Longitude:    999.999999,
			Error:        handler.ErrNoResults,
		},
		{
			Name:         "Valid address",
			ResponseBody: validReponse,
			ResponseCode: http.StatusOK,
			Latitude:     51.522214,
			Longitude:    -0.113565,
			Result: &pb.Address{
				LineOne:  "160 Grays Inn Road",
				LineTwo:  "Holborn",
				Postcode: "WC1X 8ED",
				Country:  "United Kingdom",
			},
		},
		{
			Name:         "Maps error",
			Latitude:     51.522214,
			Longitude:    -0.113565,
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

			var rsp pb.ReverseResponse
			err = h.Reverse(context.TODO(), &pb.ReverseRequest{
				Latitude: tc.Latitude, Longitude: tc.Longitude,
			}, &rsp)
			assert.Equal(t, tc.Error, err)

			if tc.Latitude != 0.0 && tc.Longitude != 0.0 {
				assert.Equal(t, fmt.Sprintf("%f", tc.Latitude), lat)
				assert.Equal(t, fmt.Sprintf("%f", tc.Longitude), lng)
			}

			if tc.Result != nil {
				assert.Equal(t, tc.Result.LineOne, rsp.Address.LineOne)
				assert.Equal(t, tc.Result.LineTwo, rsp.Address.LineTwo)
				assert.Equal(t, tc.Result.City, rsp.Address.City)
				assert.Equal(t, tc.Result.Country, rsp.Address.Country)
				assert.Equal(t, tc.Result.Postcode, rsp.Address.Postcode)
			}
		})
	}

}
