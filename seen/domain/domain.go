package domain

import (
	"time"

	"github.com/micro/micro/v3/service/store"

	model "github.com/micro/dev/model"
)

// Seen is the object which represents a user seeing a resource
type Seen struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	ResourceID   string    `json:"resource_id"`
	ResourceType string    `json:"resource_type"`
	Timestamp    time.Time `json:"timestamp"`
}

var (
	userIDIndex       = model.ByEquality("user_id")
	resourceIDIndex   = model.ByEquality("resource_id")
	resourceTypeIndex = model.ByEquality("resource_type")

	db = model.New(store.DefaultStore, Seen{}, []model.Index{
		userIDIndex, resourceIDIndex, resourceTypeIndex,
	}, &model.ModelOptions{})
)

// Create a seen object in the store
func Create(s *Seen) error {
	return db.Save(s)
}

// Delete a seen object from the store
func Delete(s *Seen) error {
	var result []Seen
	db.List(&model.Query{}, &result)
	// db.Where(s).Delete()
	return nil
}

// Read the timestamps from the store
func Read(userID, resourceType string, resourceIDs []string) (map[string]time.Time, error) {
	return nil, nil
}
