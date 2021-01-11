package model

import (
	"time"
)

type Location struct {
	ID        string
	PlaceID   string `gorm:"index"`
	Latitude  float64
	Longitude float64
	Timestamp time.Time
}

// use the place id for the geoindex so only one result is returned per place
func (l *Location) Id() string {
	return l.PlaceID
}

func (l *Location) Lat() float64 {
	return l.Latitude
}

func (l *Location) Lon() float64 {
	return l.Longitude
}
