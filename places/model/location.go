package model

import (
	"time"
)

type Location struct {
	ID        string
	UserID    string `gorm:"index"`
	Latitude  float64
	Longitude float64
	Timestamp time.Time
}

// use the user id for the geoindex so only one result is returned per user
func (l *Location) Id() string {
	return l.UserID
}

func (l *Location) Lat() float64 {
	return l.Latitude
}

func (l *Location) Lon() float64 {
	return l.Longitude
}
