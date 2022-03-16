package handler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/store"
	"github.com/micro/services/pkg/tenant"

	pb "analytics/proto"
)

// Analytics implements the notes proto definition
type Analytics struct {
	lock sync.Mutex
}

// New returns an initialized Analytics
func New() *Analytics {
	return &Analytics{}
}

// Track inserts a new Event in the store
func (a *Analytics) Track(ctx context.Context, req *pb.TrackRequest, rsp *pb.TrackResponse) error {
	// Validate the request
	if len(req.Name) == 0 {
		return errors.BadRequest("analytics.track", "missing name")
	}

	defer func() {
		a.lock.Lock()
		defer a.lock.Unlock()
		tnt, ok := tenant.FromContext(ctx)
		if !ok {
			tnt = "default"
		}

		key := fmt.Sprintf("%s:%s", tnt, req.Name)

		var event *pb.Event

		// Create new Event if it doesn't exist or increment the value if it exists
		recs, err := store.Read(key)
		if err == store.ErrNotFound {
			t := time.Now().Format(time.RFC3339)
			event = &pb.Event{
				Name:    req.Name,
				Created: t,
				Value:   1,
			}
		} else if err == nil {
			if err := recs[0].Decode(&event); err != nil {
				return
			}
			event.Value = event.Value + 1
		} else {
			return
		}

		// write Event data to store
		rec := store.NewRecord(key, event)

		if err = store.Write(rec); err != nil {
			return
		}
	}()

	return nil
}

// Get returns a single Event
func (a *Analytics) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	// Validate the request
	if len(req.Name) == 0 {
		return errors.BadRequest("analytics.get", "missing name")
	}

	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "default"
	}

	key := fmt.Sprintf("%s:%s", tnt, req.Name)

	// Get the Event from the store
	recs, err := store.Read(key)
	if err == store.ErrNotFound {
		return errors.NotFound("analytics.get", "Event not found")
	} else if err != nil {
		return errors.InternalServerError("analytics.get", "Error reading from store: %v", err.Error())
	}

	// Decode the Event
	var event *pb.Event
	if err := recs[0].Decode(&event); err != nil {
		return errors.InternalServerError("analytics.get", "Error unmarshaling JSON: %v", err.Error())
	}

	rsp.Event = event

	return nil
}

// Delete removes the Event from the store
func (a *Analytics) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	// Validate the request
	if len(req.Name) == 0 {
		return errors.BadRequest("analytics.delete", "missing name")
	}

	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "default"
	}

	key := fmt.Sprintf("%s:%s", tnt, req.Name)

	// Get the Event from the store
	recs, err := store.Read(key)
	if err == store.ErrNotFound {
		return errors.NotFound("analytics.delete", "Event not found")
	} else if err != nil {
		return errors.InternalServerError("analytics.delete", "Error reading from store: %v", err.Error())
	}

	// Decode the Event
	var event *pb.Event
	if err := recs[0].Decode(&event); err != nil {
		return errors.InternalServerError("analytics.delete", "Error unmarshaling JSON: %v", err.Error())
	}

	// now delete it
	if err := store.Delete(key); err != nil && err != store.ErrNotFound {
		return errors.InternalServerError("analytics.delete", "Failed to delete event")
	}

	rsp.Event = event

	return nil
}

// List returns all of the Events in the store
func (a *Analytics) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "default"
	}

	// Read all events from the store
	recs, err := store.Read(tnt, store.ReadPrefix())
	if err != nil {
		return errors.InternalServerError("analytics.list", "Error reading from store: %v", err.Error())
	}

	// Initialize the response events slice
	rsp.Events = make([]*pb.Event, len(recs))

	// Retrieve all of the records in the store
	for i, rec := range recs {

		// Unmarshal the events into the response
		if err := rec.Decode(&rsp.Events[i]); err != nil {
			return errors.InternalServerError("analytics.list", "Error decoding event: %v", err.Error())
		}
	}

	return nil
}
