package prayer

import(
	"github.com/m3o/m3o-go/client"
)

func NewPrayerService(token string) *PrayerService {
	return &PrayerService{
		client: client.NewClient(&client.Options{
			Token: token,
		}),
	}
}

type PrayerService struct {
	client *client.Client
}


// Get the prayer (salah) times for a location on a given date
func (t *PrayerService) Times(request *TimesRequest) (*TimesResponse, error) {
	rsp := &TimesResponse{}
	return rsp, t.client.Call("prayer", "Times", request, rsp)
}




type PrayerTime struct {
  // asr time
  Asr string `json:"asr"`
  // date for prayer times in YYYY-MM-DD format
  Date string `json:"date"`
  // fajr time
  Fajr string `json:"fajr"`
  // isha time
  Isha string `json:"isha"`
  // maghrib time
  Maghrib string `json:"maghrib"`
  // time of sunrise
  Sunrise string `json:"sunrise"`
  // zuhr time
  Zuhr string `json:"zuhr"`
}

type TimesRequest struct {
  // optional date in YYYY-MM-DD format, otherwise uses today
  Date string `json:"date"`
  // number of days to request times for
  Days int32 `json:"days"`
  // optional latitude used in place of location
  Latitude float64 `json:"latitude"`
  // location to retrieve prayer times for.
  // this can be a specific address, city, etc
  Location string `json:"location"`
  // optional longitude used in place of location
  Longitude float64 `json:"longitude"`
}

type TimesResponse struct {
  // date of request
  Date string `json:"date"`
  // number of days
  Days int32 `json:"days"`
  // latitude of location
  Latitude float64 `json:"latitude"`
  // location for the request
  Location string `json:"location"`
  // longitude of location
  Longitude float64 `json:"longitude"`
  // prayer times for the given location
  Times []PrayerTime `json:"times"`
}

# { Prayer

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/Prayer/api](https://m3o.com/Prayer/api).

Endpoints:

#times

// Get the prayer (salah) times for a location on a given date


[https://m3o.com/prayer/api#times](https://m3o.com/prayer/api#times)
