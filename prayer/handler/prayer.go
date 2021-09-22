package handler

import (
	"context"

	"github.com/hablullah/go-prayer"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/errors"
	pb "github.com/micro/services/prayer/proto"
	geocode "github.com/micro/services/geocoding/proto"
	timepb "github.com/micro/services/time/proto"
)

type Prayer struct{}

func New(c client.Client) *Prayer {
	return &Prayer{
		Geocode: geocode.NewGeocodingService("geocoding", c),
		Time: timepb.NewTimeService("time", c),
	}
}

func (p *Prayer) Times(ctx context.Context, req *pb.TimesRequest, rsp *pb.TimesResponse) error {
	// geocode the location
	if len(req.Location) == 0 && req.Latitude == 0.0 && req.Longitude == 0.0 {
		return errors.BadRequest("prayer.times", "missing location")
	}

	latitude := req.Latitude
	longitude := req.Longitude

	// get lat/lng if location is specified
	if len(req.Location) > 0 {
		resp, err := p.Geocode.Lookup(ctx, &geocode.LookupRequest{
			Address: req.Location,
		})
		if err != nil {
			return errors.InternalServerError("prayer.times", "failed to lookup location")
		}
		latitude = resp.Location.Latitude
		longitude = resp.Location.Longitude
	}

	if latitude == 0.0 && longitude == 0.0 {
		return errors.BadRequest("prayer.times", "missing location")
	}

	// get the timezone
	resp, err := p.Time.Zone(ctx, &timepb.ZoneRequest{
		Location: fmt.Sprintf("%v,%v", latitude, longitude),
	})
	if err != nil {
		return errors.InternalServerError("prayer.times", "failed to lookup timezone")
	}

	cfg := prayer.Config{
		Latitude:          latitude,
		Longitude:         longitude,
		CalculationMethod: prayer.Kemenag,
		AsrConvention:     prayer.Shafii,
		PreciseToSeconds:  false,
	}

	// current date
	date := time.Now()

	// if date is specified then change it
	if len(req.Date) > 0 {
		d, err := time.Parse("20060102", req.Date)
		if err == nil {
			date = d
		}
	}

	// set time zone
	zone := time.FixedZone(resp.Abbreviation, 0)
	date = date.In(loc)

	rsp.Location = req.Location
	rsp.Latitude = latitude
	rsp.Longitude = longitude

	times, err := prayer.Calculate(cfg, date)
	if err != nil {
		return errors.InternalServerError("prayer.times", "failed to retrieve prayer times")
	}

	rsp.Times = append(rsp.Times, &pb.PrayerTime{
		Date: date.Format("2006-01-02"),
		Fajr: times.Fajr.Format("15:04"),
		Sunrise: times.Sunrise.Format("15:04"),
		Zuhr: times.Zuhr.Format("15:04"),
		Asr: times.Asr.Format("15:04"),
		Maghrib: times.Maghrib.Format("15:04"),
		Isha: times.Isha.Format("15:04"),
	})

	return nil
}

