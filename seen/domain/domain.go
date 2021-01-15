package domain

import (
	"time"

	model "github.com/micro/dev/model"
	"github.com/micro/micro/v3/service/store/file"
)

// Seen is the object which represents a user seeing a resource
type Seen struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	ResourceID   string    `json:"resource_id"`
	ResourceType string    `json:"resource_type"`
	Timestamp    time.Time `json:"timestamp"`
}

type Domain struct {
	db model.Model
}

func New() *Domain {
	userIDIndex := model.ByEquality("user_id")
	resourceIDIndex := model.ByEquality("resource_id")
	resourceTypeIndex := model.ByEquality("resource_type")

	db := model.New(file.NewStore(), Seen{}, []model.Index{
		userIDIndex, resourceIDIndex, resourceTypeIndex,
	}, &model.ModelOptions{})

	return &Domain{db: db}
}

// Create a seen object in the store
func (d *Domain) Create(s Seen) error {
	return d.db.Save(s)
}

// Delete a seen object from the store
func (d *Domain) Delete(s Seen) error {
	// var result []Seen
	// db.List(model.Equals("user_id", s.UserID), &result)
	// db.Where(s).Delete()
	return nil
}

// Read the timestamps from the store
func (d *Domain) Read(userID, resourceType string, resourceIDs []string) (map[string]time.Time, error) {
	return nil, nil
}
