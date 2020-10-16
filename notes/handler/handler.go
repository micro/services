package handler

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	pb "github.com/micro/services/notes/proto"
)

const storePrefix = "notes/"

// New returns an initialized notes handler
func New() pb.NotesHandler {
	return new(handler)
}

type handler struct{}

// List all the notes
func (h *handler) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	// query the store
	recs, err := store.Read("", store.Prefix(storePrefix))
	if err != nil {
		logger.Errorf("Error reading notes from the store: %v", err)
		return errors.InternalServerError("notes.List.Unknown", "Error reading from the store")
	}

	// serialize the response
	rsp.Notes = make([]*pb.Note, len(recs))
	for i, r := range recs {
		var note pb.Note
		if err := json.Unmarshal(r.Value, &note); err != nil {
			logger.Errorf("Error unmarshaling note: %v", err)
			return errors.InternalServerError("notes.List.Unknown", "Error unmarshaling note")
		}
		rsp.Notes[i] = &note
	}

	return nil
}

// Create a note
func (h *handler) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// validate the request
	if len(req.Title) == 0 {
		return errors.BadRequest("notes.Create.MissingTitle", "Missing title")
	}

	// construct the note
	note := &pb.Note{
		Id:      uuid.New().String(),
		Created: time.Now().Unix(),
		Title:   req.Title,
		Text:    req.Text,
	}

	// marshal the note to bytes
	bytes, err := json.Marshal(note)
	if err != nil {
		logger.Errorf("Error marshaling note: %v", err)
		return errors.InternalServerError("notes.Create.Unknown", "Error marshaling note")
	}

	// write to the store
	key := storePrefix + note.Id
	if err := store.Write(&store.Record{Key: key, Value: bytes}); err != nil {
		logger.Errorf("Error writing to store: %v", err)
		return errors.InternalServerError("notes.Create.Unknown", "Error writing to store")
	}

	// return the id
	rsp.Id = note.Id
	return nil
}

// Delete a note
func (h *handler) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	// validate the request
	if len(req.Id) == 0 {
		return errors.BadRequest("notes.Delete.MissingID", "Missing id")
	}

	// delete the note from the store
	if err := store.Delete(storePrefix + req.Id); err == store.ErrNotFound {
		return errors.NotFound("notes.Delete.InvalidID", "Note not found with this ID")
	} else if err != nil {
		logger.Errorf("Error deleting from the store: %v", err)
		return errors.InternalServerError("notes.Delete.Unknown", "Error deleting from the store")
	}

	return nil
}

// Update a note
func (h *handler) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	// validate the request
	if len(req.Id) == 0 {
		return errors.BadRequest("notes.Update.MissingID", "Missing id")
	}
	if len(req.Title) == 0 {
		return errors.BadRequest("notes.Update.MissingTitle", "Missing title")
	}

	// read the note from the store
	recs, err := store.Read(storePrefix + req.Id)
	if err == store.ErrNotFound {
		return errors.NotFound("notes.Update.InvalidID", "Note not found with this ID")
	} else if err != nil {
		logger.Errorf("Error reading from the store: %v", err)
		return errors.InternalServerError("notes.Update.Unknown", "Error reading from the store")
	}

	// unmarshal the note
	var note pb.Note
	if err := json.Unmarshal(recs[0].Value, &note); err != nil {
		logger.Errorf("Error unmarshaling note: %v", err)
		return errors.InternalServerError("notes.Update.Unknown", "Error unmarshaling note")
	}

	// assign the new title and text
	note.Title = req.Title
	note.Text = req.Text

	// marshal the note to bytes
	bytes, err := json.Marshal(note)
	if err != nil {
		logger.Errorf("Error marshaling note: %v", err)
		return errors.InternalServerError("notes.Update.Unknown", "Error marshaling note")
	}

	// write to the store
	key := storePrefix + note.Id
	if err := store.Write(&store.Record{Key: key, Value: bytes}); err != nil {
		logger.Errorf("Error writing to store: %v", err)
		return errors.InternalServerError("notes.Update.Unknown", "Error writing to store")
	}

	return nil
}

// UpdateStream updates a note every time an update is sent on the stream
func (h *handler) UpdateStream(ctx context.Context, stream pb.Notes_UpdateStreamStream) error {
	for {
		uReq, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return errors.InternalServerError("notes.UpdateStream.Unknown", "Error reading from stream")
		}

		if err := h.Update(ctx, uReq, nil); err != nil {
			return err
		}
	}
}
