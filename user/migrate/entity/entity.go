package entity

import (
	"fmt"
	"time"
)

type Row struct {
	Id        string
	Data      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func KeyPrefix(tenantId string) string {
	return fmt.Sprintf("user/%s/", tenantId)
}
