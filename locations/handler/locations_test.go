package handler_test

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	geo "github.com/hailocab/go-geoindex"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/micro/services/locations/handler"
	"github.com/micro/services/locations/model"
	pb "github.com/micro/services/locations/proto"
)

func testHandler(t *testing.T) pb.LocationsHandler {
	// connect to the database
	db, err := gorm.Open(postgres.Open("postgresql://postgres@localhost:5432/locations?sslmode=disable"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}

	// migrate the database
	if err := db.AutoMigrate(&model.Location{}); err != nil {
		t.Fatalf("Error migrating database: %v", err)
	}

	// clean any data from a previous run
	if err := db.Exec("TRUNCATE TABLE locations CASCADE").Error; err != nil {
		t.Fatalf("Error cleaning database: %v", err)
	}

	return &handler.Locations{DB: db, Geoindex: geo.NewPointsIndex(geo.Km(0.1))}
}

func TestSave(t *testing.T) {
	tt := []struct {
		Name      string
		Locations []*pb.Location
		Error     error
	}{
		{
			Name:  "NoLocations",
			Error: handler.ErrMissingLocations,
		},
		{
			Name: "NoLatitude",
			Locations: []*pb.Location{
				{
					Longitude: &wrapperspb.DoubleValue{Value: -0.1246},
					UserId:    uuid.New().String(),
				},
			},
			Error: handler.ErrMissingLatitude,
		},
		{
			Name: "NoLongitude",
			Locations: []*pb.Location{
				{
					Latitude: &wrapperspb.DoubleValue{Value: -0.1246},
					UserId:   uuid.New().String(),
				},
			},
			Error: handler.ErrMissingLongitude,
		},
		{
			Name: "OneLocation",
			Locations: []*pb.Location{
				{
					Latitude:  &wrapperspb.DoubleValue{Value: 51.5007},
					Longitude: &wrapperspb.DoubleValue{Value: 0.1246},
					Timestamp: timestamppb.New(time.Now()),
					UserId:    uuid.New().String(),
				},
			},
		},
		{
			Name: "ManyLocations",
			Locations: []*pb.Location{
				{
					Latitude:  &wrapperspb.DoubleValue{Value: 51.5007},
					Longitude: &wrapperspb.DoubleValue{Value: 0.1246},
					Timestamp: timestamppb.New(time.Now()),
					UserId:    uuid.New().String(),
				},
				{
					Latitude:  &wrapperspb.DoubleValue{Value: 51.003},
					Longitude: &wrapperspb.DoubleValue{Value: -0.1246},
					UserId:    uuid.New().String(),
				},
			},
		},
	}

	h := testHandler(t)

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			err := h.Save(context.Background(), &pb.SaveRequest{
				Locations: tc.Locations,
			}, &pb.SaveResponse{})
			assert.Equal(t, tc.Error, err)
		})
	}
}

func TestLast(t *testing.T) {
	h := testHandler(t)

	t.Run("MissingUserIDs", func(t *testing.T) {
		err := h.Last(context.Background(), &pb.LastRequest{}, &pb.ListResponse{})
		assert.Equal(t, handler.ErrMissingUserIDs, err)
	})

	t.Run("NoMatches", func(t *testing.T) {
		var rsp pb.ListResponse
		err := h.Last(context.Background(), &pb.LastRequest{
			UserIds: []string{uuid.New().String()},
		}, &rsp)
		assert.NoError(t, err)
		assert.Empty(t, rsp.Locations)
	})

	// generate some example data to work with
	loc1 := &pb.Location{
		Latitude:  &wrapperspb.DoubleValue{Value: 51.5007},
		Longitude: &wrapperspb.DoubleValue{Value: 0.1246},
		Timestamp: timestamppb.New(time.Now()),
		UserId:    "a",
	}
	loc2 := &pb.Location{
		Latitude:  &wrapperspb.DoubleValue{Value: 51.6007},
		Longitude: &wrapperspb.DoubleValue{Value: 0.1546},
		Timestamp: timestamppb.New(time.Now()),
		UserId:    "b",
	}
	loc3 := &pb.Location{
		Latitude:  &wrapperspb.DoubleValue{Value: 52.6007},
		Longitude: &wrapperspb.DoubleValue{Value: 0.2546},
		Timestamp: timestamppb.New(time.Now()),
		UserId:    loc2.UserId,
	}
	err := h.Save(context.TODO(), &pb.SaveRequest{
		Locations: []*pb.Location{loc1, loc2, loc3},
	}, &pb.SaveResponse{})
	assert.NoError(t, err)

	t.Run("OneUser", func(t *testing.T) {
		var rsp pb.ListResponse
		err := h.Last(context.Background(), &pb.LastRequest{
			UserIds: []string{loc2.UserId},
		}, &rsp)
		assert.NoError(t, err)

		if len(rsp.Locations) != 1 {
			t.Fatalf("One location should be returned")
		}
		assert.Equal(t, loc3.UserId, rsp.Locations[0].UserId)
		assert.Equal(t, loc3.Latitude.Value, rsp.Locations[0].Latitude.Value)
		assert.Equal(t, loc3.Longitude.Value, rsp.Locations[0].Longitude.Value)
		assert.Equal(t, loc3.Timestamp.AsTime(), rsp.Locations[0].Timestamp.AsTime())
	})
	t.Run("ManyUser", func(t *testing.T) {
		var rsp pb.ListResponse
		err := h.Last(context.Background(), &pb.LastRequest{
			UserIds: []string{loc1.UserId, loc2.UserId},
		}, &rsp)
		assert.NoError(t, err)

		if len(rsp.Locations) != 2 {
			t.Fatalf("Two locations should be returned")
		}

		// sort using user_id so we can hardcode the index
		sort.Slice(rsp.Locations, func(i, j int) bool {
			return rsp.Locations[i].UserId > rsp.Locations[j].UserId
		})

		assert.Equal(t, loc1.UserId, rsp.Locations[1].UserId)
		assert.Equal(t, loc1.Latitude.Value, rsp.Locations[1].Latitude.Value)
		assert.Equal(t, loc1.Longitude.Value, rsp.Locations[1].Longitude.Value)
		assert.Equal(t, loc1.Timestamp.AsTime(), rsp.Locations[1].Timestamp.AsTime())

		assert.Equal(t, loc3.UserId, rsp.Locations[0].UserId)
		assert.Equal(t, loc3.Latitude.Value, rsp.Locations[0].Latitude.Value)
		assert.Equal(t, loc3.Longitude.Value, rsp.Locations[0].Longitude.Value)
		assert.Equal(t, loc3.Timestamp.AsTime(), rsp.Locations[0].Timestamp.AsTime())
	})
}

func TestNear(t *testing.T) {
	lat := &wrapperspb.DoubleValue{Value: 51.510357}
	lng := &wrapperspb.DoubleValue{Value: -0.116773}
	rad := &wrapperspb.DoubleValue{Value: 2.0}

	inBoundsLat := &wrapperspb.DoubleValue{Value: 51.5110}
	inBoundsLng := &wrapperspb.DoubleValue{Value: -0.1142}

	outOfBoundsLat := &wrapperspb.DoubleValue{Value: 51.5415}
	outOfBoundsLng := &wrapperspb.DoubleValue{Value: -0.0028}

	tt := []struct {
		Name           string
		Locations      []*pb.Location
		Results        []*pb.Location
		QueryLatitude  *wrapperspb.DoubleValue
		QueryLongitude *wrapperspb.DoubleValue
		QueryRadius    *wrapperspb.DoubleValue
		Error          error
	}{
		{
			Name:           "MissingLatitude",
			QueryLongitude: lng,
			QueryRadius:    rad,
			Error:          handler.ErrMissingLatitude,
		},
		{
			Name:          "MissingLongitude",
			QueryLatitude: lat,
			QueryRadius:   rad,
			Error:         handler.ErrMissingLongitude,
		},
		{
			Name:           "MissingRadius",
			QueryLatitude:  lat,
			QueryLongitude: lng,
			Error:          handler.ErrMissingRadius,
		},
		{
			Name:           "NoLocations",
			QueryLatitude:  lat,
			QueryLongitude: lng,
			QueryRadius:    rad,
		},
		{
			Name:           "OneWithinRadius",
			QueryLatitude:  lat,
			QueryLongitude: lng,
			QueryRadius:    rad,
			Locations: []*pb.Location{
				&pb.Location{
					Latitude:  inBoundsLat,
					Longitude: inBoundsLng,
					UserId:    "in",
				},
				&pb.Location{
					Latitude:  outOfBoundsLat,
					Longitude: outOfBoundsLng,
					UserId:    "out",
				},
			},
			Results: []*pb.Location{
				&pb.Location{
					Latitude:  inBoundsLat,
					Longitude: inBoundsLng,
					UserId:    "in",
				},
			},
		},
		{
			Name:           "NoneWithinRadius",
			QueryLatitude:  lat,
			QueryLongitude: lng,
			QueryRadius:    &wrapperspb.DoubleValue{Value: 0.01},
			Locations: []*pb.Location{
				&pb.Location{
					Latitude:  inBoundsLat,
					Longitude: inBoundsLng,
					UserId:    "in",
				},
				&pb.Location{
					Latitude:  outOfBoundsLat,
					Longitude: outOfBoundsLng,
					UserId:    "out",
				},
			},
		},
		{
			Name:           "TwoLocationsForUser",
			QueryLatitude:  lat,
			QueryLongitude: lng,
			QueryRadius:    rad,
			Locations: []*pb.Location{
				&pb.Location{
					Latitude:  inBoundsLat,
					Longitude: inBoundsLng,
					UserId:    "in",
				},
				&pb.Location{
					Latitude:  outOfBoundsLat,
					Longitude: outOfBoundsLng,
					UserId:    "out",
				},
				&pb.Location{
					Latitude:  inBoundsLat,
					Longitude: inBoundsLng,
					UserId:    "out",
				},
			},
			Results: []*pb.Location{
				&pb.Location{
					Latitude:  inBoundsLat,
					Longitude: inBoundsLng,
					UserId:    "in",
				},
				&pb.Location{
					Latitude:  inBoundsLat,
					Longitude: inBoundsLng,
					UserId:    "out",
				},
			},
		},
		{
			Name:           "ManyWithinRadius",
			QueryLatitude:  lat,
			QueryLongitude: lng,
			QueryRadius:    &wrapperspb.DoubleValue{Value: 20},
			Locations: []*pb.Location{
				&pb.Location{
					Latitude:  inBoundsLat,
					Longitude: inBoundsLng,
					UserId:    "in",
				},
				&pb.Location{
					Latitude:  outOfBoundsLat,
					Longitude: outOfBoundsLng,
					UserId:    "out",
				},
			},
			Results: []*pb.Location{
				&pb.Location{
					Latitude:  inBoundsLat,
					Longitude: inBoundsLng,
					UserId:    "in",
				},
				&pb.Location{
					Latitude:  outOfBoundsLat,
					Longitude: outOfBoundsLng,
					UserId:    "out",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			h := testHandler(t)

			// create the locations
			if len(tc.Locations) > 0 {
				err := h.Save(context.TODO(), &pb.SaveRequest{Locations: tc.Locations}, &pb.SaveResponse{})
				assert.NoError(t, err)
			}

			// find near locations
			var rsp pb.ListResponse
			err := h.Near(context.TODO(), &pb.NearRequest{
				Latitude:  tc.QueryLatitude,
				Longitude: tc.QueryLongitude,
				Radius:    tc.QueryRadius,
			}, &rsp)
			assert.Equal(t, tc.Error, err)

			// check the count of the results matches
			if len(tc.Results) != len(rsp.Locations) {
				t.Errorf("Incorrect number of results returned. Expected %v, got %v", len(tc.Results), len(rsp.Locations))
			}

			// validate the results match
			sort.Slice(rsp.Locations, func(i, j int) bool {
				return rsp.Locations[i].UserId > rsp.Locations[j].UserId
			})
			sort.Slice(tc.Results, func(i, j int) bool {
				return tc.Results[i].UserId > tc.Results[j].UserId
			})
			for i, r := range tc.Results {
				l := rsp.Locations[i]
				assert.Equal(t, r.UserId, l.UserId)
				assert.Equal(t, r.Latitude.Value, l.Latitude.Value)
				assert.Equal(t, r.Longitude.Value, l.Longitude.Value)
			}
		})
	}
}

func TestRead(t *testing.T) {
	h := testHandler(t)

	baseTime := time.Now().Add(time.Hour * -24)

	t.Run("MissingUserIDs", func(t *testing.T) {
		err := h.Read(context.Background(), &pb.ReadRequest{
			After:  timestamppb.New(baseTime),
			Before: timestamppb.New(baseTime),
		}, &pb.ListResponse{})
		assert.Equal(t, handler.ErrMissingUserIDs, err)
	})

	t.Run("MissingAfter", func(t *testing.T) {
		err := h.Read(context.Background(), &pb.ReadRequest{
			UserIds: []string{uuid.New().String()},
			Before:  timestamppb.New(baseTime),
		}, &pb.ListResponse{})
		assert.Equal(t, handler.ErrMissingAfter, err)
	})

	t.Run("MissingBefore", func(t *testing.T) {
		err := h.Read(context.Background(), &pb.ReadRequest{
			UserIds: []string{uuid.New().String()},
			After:   timestamppb.New(baseTime),
		}, &pb.ListResponse{})
		assert.Equal(t, handler.ErrMissingBefore, err)
	})

	// generate some example data to work with
	loc1 := &pb.Location{
		Latitude:  &wrapperspb.DoubleValue{Value: 51.5007},
		Longitude: &wrapperspb.DoubleValue{Value: 0.1246},
		Timestamp: timestamppb.New(baseTime.Add(time.Minute * 10)),
		UserId:    "a",
	}
	loc2 := &pb.Location{
		Latitude:  &wrapperspb.DoubleValue{Value: 51.6007},
		Longitude: &wrapperspb.DoubleValue{Value: 0.1546},
		Timestamp: timestamppb.New(baseTime.Add(time.Minute * 20)),
		UserId:    "b",
	}
	loc3 := &pb.Location{
		Latitude:  &wrapperspb.DoubleValue{Value: 52.6007},
		Longitude: &wrapperspb.DoubleValue{Value: 0.2546},
		Timestamp: timestamppb.New(baseTime.Add(time.Minute * 40)),
		UserId:    loc2.UserId,
	}
	err := h.Save(context.TODO(), &pb.SaveRequest{
		Locations: []*pb.Location{loc1, loc2, loc3},
	}, &pb.SaveResponse{})
	assert.NoError(t, err)

	t.Run("NoMatches", func(t *testing.T) {
		var rsp pb.ListResponse
		err := h.Read(context.Background(), &pb.ReadRequest{
			UserIds: []string{uuid.New().String()},
			After:   timestamppb.New(baseTime),
			Before:  timestamppb.New(baseTime.Add(time.Hour)),
		}, &rsp)
		assert.NoError(t, err)
		assert.Empty(t, rsp.Locations)
	})

	t.Run("OneUserID", func(t *testing.T) {
		var rsp pb.ListResponse
		err := h.Read(context.Background(), &pb.ReadRequest{
			UserIds: []string{loc2.UserId},
			After:   timestamppb.New(baseTime),
			Before:  timestamppb.New(baseTime.Add(time.Hour)),
		}, &rsp)
		assert.NoError(t, err)

		if len(rsp.Locations) != 2 {
			t.Fatalf("Two locations should be returned")
		}
		assert.Equal(t, loc2.UserId, rsp.Locations[0].UserId)
		assert.Equal(t, loc2.Latitude.Value, rsp.Locations[0].Latitude.Value)
		assert.Equal(t, loc2.Longitude.Value, rsp.Locations[0].Longitude.Value)
		assert.Equal(t, loc2.Timestamp.AsTime(), rsp.Locations[0].Timestamp.AsTime())

		assert.Equal(t, loc3.UserId, rsp.Locations[1].UserId)
		assert.Equal(t, loc3.Latitude.Value, rsp.Locations[1].Latitude.Value)
		assert.Equal(t, loc3.Longitude.Value, rsp.Locations[1].Longitude.Value)
		assert.Equal(t, loc3.Timestamp.AsTime(), rsp.Locations[1].Timestamp.AsTime())
	})

	t.Run("OneUserIDReducedTime", func(t *testing.T) {
		var rsp pb.ListResponse
		err := h.Read(context.Background(), &pb.ReadRequest{
			UserIds: []string{loc2.UserId},
			After:   timestamppb.New(baseTime),
			Before:  timestamppb.New(baseTime.Add(time.Minute * 30)),
		}, &rsp)
		assert.NoError(t, err)

		if len(rsp.Locations) != 1 {
			t.Fatalf("One location should be returned")
		}
		assert.Equal(t, loc2.UserId, rsp.Locations[0].UserId)
		assert.Equal(t, loc2.Latitude.Value, rsp.Locations[0].Latitude.Value)
		assert.Equal(t, loc2.Longitude.Value, rsp.Locations[0].Longitude.Value)
		assert.Equal(t, loc2.Timestamp.AsTime(), rsp.Locations[0].Timestamp.AsTime())
	})

	t.Run("TwoUserIDs", func(t *testing.T) {
		var rsp pb.ListResponse
		err := h.Read(context.Background(), &pb.ReadRequest{
			UserIds: []string{loc1.UserId, loc2.UserId},
			After:   timestamppb.New(baseTime),
			Before:  timestamppb.New(baseTime.Add(time.Minute * 30)),
		}, &rsp)
		assert.NoError(t, err)

		if len(rsp.Locations) != 2 {
			t.Fatalf("Two locations should be returned")
		}
		assert.Equal(t, loc1.UserId, rsp.Locations[0].UserId)
		assert.Equal(t, loc1.Latitude.Value, rsp.Locations[0].Latitude.Value)
		assert.Equal(t, loc1.Longitude.Value, rsp.Locations[0].Longitude.Value)
		assert.Equal(t, loc1.Timestamp.AsTime(), rsp.Locations[0].Timestamp.AsTime())

		assert.Equal(t, loc2.UserId, rsp.Locations[1].UserId)
		assert.Equal(t, loc2.Latitude.Value, rsp.Locations[1].Latitude.Value)
		assert.Equal(t, loc2.Longitude.Value, rsp.Locations[1].Longitude.Value)
		assert.Equal(t, loc2.Timestamp.AsTime(), rsp.Locations[1].Timestamp.AsTime())
	})
}
