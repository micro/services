package entity

import "time"

type Row struct {
	Id        string
	Data      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
