package domain

import (
	"time"

	"github.com/micro/micro/v3/service/model"
	"github.com/micro/micro/v3/service/store/file"
)

// Seen is the object which represents a user seeing a resource
type Seen struct {
	ID           string
	UserID       string
	ResourceID   string
	ResourceType string
	Timestamp    time.Time
}

type Domain struct {
	db model.Model
}

var (
	userIDIndex       = model.ByEquality("UserID")
	resourceIDIndex   = model.ByEquality("ResourceID")
	resourceTypeIndex = model.ByEquality("ResourceType")
)

func New() *Domain {
	db := model.New(file.NewStore(), Seen{}, []model.Index{
		userIDIndex, resourceIDIndex, resourceTypeIndex,
	}, &model.ModelOptions{})

	return &Domain{db: db}
}

// Create a seen object in the store
func (d *Domain) Create(s Seen) error {
	return d.db.Create(s)
}

// Delete a seen object from the store
func (d *Domain) Delete(s Seen) error {
	// load all the users objects and then delete only the ones which match the resource, unfortunately
	// the model doesn't yet support querying by multiple columns
	var all []Seen
	if err := d.db.Read(model.Equals("UserID", s.UserID), &all); err != nil {
		return err
	}
	for _, a := range all {
		if s.ResourceID != a.ResourceID {
			continue
		}
		if s.ResourceType != s.ResourceType {
			continue
		}

		q := model.Equals("ID", s.ID)
		q.Order.Type = model.OrderTypeUnordered
		if err := d.db.Delete(q); err != nil {
			return err
		}
	}
	return nil
}

// Read the timestamps from the store
func (d *Domain) Read(userID, resourceType string, resourceIDs []string) (map[string]time.Time, error) {
	// load all the users objects and then return only the timestamps for the ones which match the
	// resource, unfortunately the model doesn't yet support querying by multiple columns
	var all []Seen
	if err := d.db.Read(model.Equals("UserID", userID), &all); err != nil {
		return nil, err
	}

	result := map[string]time.Time{}
	for _, a := range all {
		if a.ResourceType != resourceType {
			continue
		}

		for _, id := range resourceIDs {
			if id != a.ResourceID {
				continue
			}

			result[id] = a.Timestamp
			break
		}
	}

	return result, nil
}
