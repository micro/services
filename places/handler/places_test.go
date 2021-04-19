package handler_test

import (
	"context"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	geo "github.com/hailocab/go-geoindex"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/micro/services/places/handler"
	"github.com/micro/services/places/model"
	pb "github.com/micro/services/places/proto"
)

func testHandler(t *testing.T) pb.PlacesHandler {
	// connect to the database
	addr := os.Getenv("POSTGRES_URL")
	if len(addr) == 0 {
		addr = "postgresql://postgres@localhost:5432/postgres?sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(addr), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}

	// clean any data from a previous run
	if err := db.Exec("DROP TABLE IF EXISTS locations CASCADE").Error; err != nil {
		t.Fatalf("Error cleaning database: %v", err)
	}

	// migrate the database
	if err := db.AutoMigrate(&model.Location{}); err != nil {
		t.Fatalf("Error migrating database: %v", err)
	}

	return &handler.Places{DB: db, Geoindex: geo.NewPointsIndex(geo.Km(0.1))}
}

func TestSave(t *testing.T) {
	tt := []struct {
		Name   string
		Places []*pb.Location
		Error  error
	}{
		{
			Name:  "NoPlaces",
			Error: handler.ErrMissingPlaces,
		},
		{
			Name: "NoLatitude",
			Places: []*pb.Location{
				{
					Longitude: -0.1246,
					Id:        uuid.New().String(),
				},
			},
			Error: handler.ErrMissingLatitude,
		},
		{
			Name: "NoLongitude",
			Places: []*pb.Location{
				{
					Latitude: -0.1246,
					Id:       uuid.New().String(),
				},
			},
			Error: handler.ErrMissingLongitude,
		},
		{
			Name: "OneLocation",
			Places: []*pb.Location{
				{
					Latitude:  51.5007,
					Longitude: 0.1246,
					Timestamp: time.Now().Unix(),
					Id:        uuid.New().String(),
				},
			},
		},
		{
			Name: "ManyPlaces",
			Places: []*pb.Location{
				{
					Latitude:  51.5007,
					Longitude: 0.1246,
					Timestamp: time.Now().Unix(),
					Id:        uuid.New().String(),
				},
				{
					Latitude:  51.003,
					Longitude: -0.1246,
					Id:        uuid.New().String(),
				},
			},
		},
	}

	h := testHandler(t)

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			err := h.Save(context.Background(), &pb.SaveRequest{
				Places: tc.Places,
			}, &pb.SaveResponse{})
			assert.Equal(t, tc.Error, err)
		})
	}
}

func TestLast(t *testing.T) {
	h := testHandler(t)

	t.Run("MissingIDs", func(t *testing.T) {
		err := h.Last(context.Background(), &pb.LastRequest{}, &pb.ListResponse{})
		assert.Equal(t, handler.ErrMissingIDs, err)
	})

	t.Run("NoMatches", func(t *testing.T) {
		var rsp pb.ListResponse
		err := h.Last(context.Background(), &pb.LastRequest{
			Ids: []string{uuid.New().String()},
		}, &rsp)
		assert.NoError(t, err)
		assert.Empty(t, rsp.Places)
	})
	tn := time.Now()

	// generate some example data to work with
	loc1 := &pb.Location{
		Latitude:  51.5007,
		Longitude: 0.1246,
		Timestamp: tn.Unix(),
		Id:        "a",
	}
	loc2 := &pb.Location{
		Latitude:  51.6007,
		Longitude: 0.1546,
		Timestamp: tn.Add(1 * time.Second).Unix(),
		Id:        "b",
	}
	loc3 := &pb.Location{
		Latitude:  52.6007,
		Longitude: 0.2546,
		Timestamp: tn.Add(2 * time.Second).Unix(),
		Id:        loc2.Id,
	}
	err := h.Save(context.TODO(), &pb.SaveRequest{
		Places: []*pb.Location{loc1, loc2, loc3},
	}, &pb.SaveResponse{})
	assert.NoError(t, err)

	t.Run("OneUser", func(t *testing.T) {
		var rsp pb.ListResponse
		err := h.Last(context.Background(), &pb.LastRequest{
			Ids: []string{loc3.Id},
		}, &rsp)
		assert.NoError(t, err)

		if len(rsp.Places) != 1 {
			t.Fatalf("One location should be returned")
		}
		assert.Equal(t, loc3.Id, rsp.Places[0].Id)
		assert.Equal(t, loc3.Latitude, rsp.Places[0].Latitude)
		assert.Equal(t, loc3.Longitude, rsp.Places[0].Longitude)
		assert.Equal(t, microSecondTime(loc3.Timestamp), microSecondTime(rsp.Places[0].Timestamp))
	})

	t.Run("ManyUser", func(t *testing.T) {
		var rsp pb.ListResponse
		err := h.Last(context.Background(), &pb.LastRequest{
			Ids: []string{loc1.Id, loc2.Id},
		}, &rsp)
		assert.NoError(t, err)

		if len(rsp.Places) != 2 {
			t.Fatalf("Two places should be returned")
		}

		// sort using user_id so we can hardcode the index
		sort.Slice(rsp.Places, func(i, j int) bool {
			return rsp.Places[i].Id > rsp.Places[j].Id
		})

		assert.Equal(t, loc1.Id, rsp.Places[1].Id)
		assert.Equal(t, loc1.Latitude, rsp.Places[1].Latitude)
		assert.Equal(t, loc1.Longitude, rsp.Places[1].Longitude)
		assert.Equal(t, microSecondTime(loc1.Timestamp), microSecondTime(rsp.Places[1].Timestamp))

		assert.Equal(t, loc3.Id, rsp.Places[0].Id)
		assert.Equal(t, loc3.Latitude, rsp.Places[0].Latitude)
		assert.Equal(t, loc3.Longitude, rsp.Places[0].Longitude)
		assert.Equal(t, microSecondTime(loc3.Timestamp), microSecondTime(rsp.Places[0].Timestamp))
	})
}

func TestNear(t *testing.T) {
	lat := 51.510357
	lng := -0.116773
	rad := 2.0

	inBoundsLat := 51.5110
	inBoundsLng := -0.1142

	outOfBoundsLat := 51.5415
	outOfBoundsLng := -0.0028

	tt := []struct {
		Name           string
		Places         []*pb.Location
		Results        []*pb.Location
		QueryLatitude  float64
		QueryLongitude float64
		QueryRadius    float64
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
			Name:           "NoPlaces",
			QueryLatitude:  lat,
			QueryLongitude: lng,
			QueryRadius:    rad,
		},
		{
			Name:           "OneWithinRadius",
			QueryLatitude:  lat,
			QueryLongitude: lng,
			QueryRadius:    rad,
			Places: []*pb.Location{
				&pb.Location{
					Latitude:  inBoundsLat,
					Longitude: inBoundsLng,
					Id:        "in",
				},
				&pb.Location{
					Latitude:  outOfBoundsLat,
					Longitude: outOfBoundsLng,
					Id:        "out",
				},
			},
			Results: []*pb.Location{
				&pb.Location{
					Latitude:  inBoundsLat,
					Longitude: inBoundsLng,
					Id:        "in",
				},
			},
		},
		{
			Name:           "NoneWithinRadius",
			QueryLatitude:  lat,
			QueryLongitude: lng,
			QueryRadius:    0.01,
			Places: []*pb.Location{
				&pb.Location{
					Latitude:  inBoundsLat,
					Longitude: inBoundsLng,
					Id:        "in",
				},
				&pb.Location{
					Latitude:  outOfBoundsLat,
					Longitude: outOfBoundsLng,
					Id:        "out",
				},
			},
		},
		{
			Name:           "TwoPlacesForUser",
			QueryLatitude:  lat,
			QueryLongitude: lng,
			QueryRadius:    rad,
			Places: []*pb.Location{
				&pb.Location{
					Latitude:  inBoundsLat,
					Longitude: inBoundsLng,
					Id:        "in",
				},
				&pb.Location{
					Latitude:  outOfBoundsLat,
					Longitude: outOfBoundsLng,
					Id:        "out",
				},
				&pb.Location{
					Latitude:  inBoundsLat,
					Longitude: inBoundsLng,
					Id:        "out",
				},
			},
			Results: []*pb.Location{
				&pb.Location{
					Latitude:  inBoundsLat,
					Longitude: inBoundsLng,
					Id:        "in",
				},
				&pb.Location{
					Latitude:  inBoundsLat,
					Longitude: inBoundsLng,
					Id:        "out",
				},
			},
		},
		{
			Name:           "ManyWithinRadius",
			QueryLatitude:  lat,
			QueryLongitude: lng,
			QueryRadius:    20,
			Places: []*pb.Location{
				&pb.Location{
					Latitude:  inBoundsLat,
					Longitude: inBoundsLng,
					Id:        "in",
				},
				&pb.Location{
					Latitude:  outOfBoundsLat,
					Longitude: outOfBoundsLng,
					Id:        "out",
				},
			},
			Results: []*pb.Location{
				&pb.Location{
					Latitude:  inBoundsLat,
					Longitude: inBoundsLng,
					Id:        "in",
				},
				&pb.Location{
					Latitude:  outOfBoundsLat,
					Longitude: outOfBoundsLng,
					Id:        "out",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			h := testHandler(t)

			// create the places
			if len(tc.Places) > 0 {
				err := h.Save(context.TODO(), &pb.SaveRequest{Places: tc.Places}, &pb.SaveResponse{})
				assert.NoError(t, err)
			}

			// find near places
			var rsp pb.ListResponse
			err := h.Near(context.TODO(), &pb.NearRequest{
				Latitude:  tc.QueryLatitude,
				Longitude: tc.QueryLongitude,
				Radius:    tc.QueryRadius,
			}, &rsp)
			assert.Equal(t, tc.Error, err)

			// check the count of the results matches
			if len(tc.Results) != len(rsp.Places) {
				t.Errorf("Incorrect number of results returned. Expected %v, got %v", len(tc.Results), len(rsp.Places))
			}

			// validate the results match
			sort.Slice(rsp.Places, func(i, j int) bool {
				return rsp.Places[i].Id > rsp.Places[j].Id
			})
			sort.Slice(tc.Results, func(i, j int) bool {
				return tc.Results[i].Id > tc.Results[j].Id
			})
			for i, r := range tc.Results {
				l := rsp.Places[i]
				assert.Equal(t, r.Id, l.Id)
				assert.Equal(t, r.Latitude, l.Latitude)
				assert.Equal(t, r.Longitude, l.Longitude)
			}
		})
	}
}

func TestRead(t *testing.T) {
	h := testHandler(t)

	baseTime := time.Now().Add(time.Hour * -24)

	t.Run("MissingIDs", func(t *testing.T) {
		err := h.Read(context.Background(), &pb.ReadRequest{
			After:  baseTime.Unix(),
			Before: baseTime.Unix(),
		}, &pb.ListResponse{})
		assert.Equal(t, handler.ErrMissingIDs, err)
	})

	t.Run("MissingAfter", func(t *testing.T) {
		err := h.Read(context.Background(), &pb.ReadRequest{
			Ids:    []string{uuid.New().String()},
			Before: baseTime.Unix(),
		}, &pb.ListResponse{})
		assert.Equal(t, handler.ErrMissingAfter, err)
	})

	t.Run("MissingBefore", func(t *testing.T) {
		err := h.Read(context.Background(), &pb.ReadRequest{
			Ids:   []string{uuid.New().String()},
			After: baseTime.Unix(),
		}, &pb.ListResponse{})
		assert.Equal(t, handler.ErrMissingBefore, err)
	})

	// generate some example data to work with
	loc1 := &pb.Location{
		Latitude:  51.5007,
		Longitude: 0.1246,
		Timestamp: baseTime.Add(time.Minute * 10).Unix(),
		Id:        "a",
	}
	loc2 := &pb.Location{
		Latitude:  51.6007,
		Longitude: 0.1546,
		Timestamp: baseTime.Add(time.Minute * 20).Unix(),
		Id:        "b",
	}
	loc3 := &pb.Location{
		Latitude:  52.6007,
		Longitude: 0.2546,
		Timestamp: baseTime.Add(time.Minute * 40).Unix(),
		Id:        loc2.Id,
	}
	err := h.Save(context.TODO(), &pb.SaveRequest{
		Places: []*pb.Location{loc1, loc2, loc3},
	}, &pb.SaveResponse{})
	assert.NoError(t, err)

	t.Run("NoMatches", func(t *testing.T) {
		var rsp pb.ListResponse
		err := h.Read(context.Background(), &pb.ReadRequest{
			Ids:    []string{uuid.New().String()},
			After:  baseTime.Unix(),
			Before: baseTime.Add(time.Hour).Unix(),
		}, &rsp)
		assert.NoError(t, err)
		assert.Empty(t, rsp.Places)
	})

	t.Run("OnePlaceID", func(t *testing.T) {
		var rsp pb.ListResponse
		err := h.Read(context.Background(), &pb.ReadRequest{
			Ids:    []string{loc2.Id},
			After:  baseTime.Unix(),
			Before: baseTime.Add(time.Hour).Unix(),
		}, &rsp)
		assert.NoError(t, err)

		if len(rsp.Places) != 2 {
			t.Fatalf("Two places should be returned")
		}
		assert.Equal(t, loc2.Id, rsp.Places[0].Id)
		assert.Equal(t, loc2.Latitude, rsp.Places[0].Latitude)
		assert.Equal(t, loc2.Longitude, rsp.Places[0].Longitude)
		assert.Equal(t, microSecondTime(loc2.Timestamp), microSecondTime(rsp.Places[0].Timestamp))

		assert.Equal(t, loc3.Id, rsp.Places[1].Id)
		assert.Equal(t, loc3.Latitude, rsp.Places[1].Latitude)
		assert.Equal(t, loc3.Longitude, rsp.Places[1].Longitude)
		assert.Equal(t, microSecondTime(loc3.Timestamp), microSecondTime(rsp.Places[1].Timestamp))
	})

	t.Run("OnePlaceIDReducedTime", func(t *testing.T) {
		var rsp pb.ListResponse
		err := h.Read(context.Background(), &pb.ReadRequest{
			Ids:    []string{loc2.Id},
			After:  baseTime.Unix(),
			Before: baseTime.Add(time.Minute * 30).Unix(),
		}, &rsp)
		assert.NoError(t, err)

		if len(rsp.Places) != 1 {
			t.Fatalf("One location should be returned")
		}
		assert.Equal(t, loc2.Id, rsp.Places[0].Id)
		assert.Equal(t, loc2.Latitude, rsp.Places[0].Latitude)
		assert.Equal(t, loc2.Longitude, rsp.Places[0].Longitude)
		assert.Equal(t, microSecondTime(loc2.Timestamp), microSecondTime(rsp.Places[0].Timestamp))
	})

	t.Run("TwoPlaceIDs", func(t *testing.T) {
		var rsp pb.ListResponse
		err := h.Read(context.Background(), &pb.ReadRequest{
			Ids:    []string{loc1.Id, loc2.Id},
			After:  baseTime.Unix(),
			Before: baseTime.Add(time.Minute * 30).Unix(),
		}, &rsp)
		assert.NoError(t, err)

		if len(rsp.Places) != 2 {
			t.Fatalf("Two places should be returned")
		}
		assert.Equal(t, loc1.Id, rsp.Places[0].Id)
		assert.Equal(t, loc1.Latitude, rsp.Places[0].Latitude)
		assert.Equal(t, loc1.Longitude, rsp.Places[0].Longitude)
		assert.Equal(t, microSecondTime(loc1.Timestamp), microSecondTime(rsp.Places[0].Timestamp))

		assert.Equal(t, loc2.Id, rsp.Places[1].Id)
		assert.Equal(t, loc2.Latitude, rsp.Places[1].Latitude)
		assert.Equal(t, loc2.Longitude, rsp.Places[1].Longitude)
		assert.Equal(t, microSecondTime(loc2.Timestamp), microSecondTime(rsp.Places[1].Timestamp))
	})
}

// postgres has a resolution of 100microseconds so just test that it's accurate to the second
func microSecondTime(t int64) time.Time {
	tt := time.Unix(t, 0)
	return time.Unix(t, int64(tt.Nanosecond()-tt.Nanosecond()%1000))
}
