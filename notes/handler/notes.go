package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/store"
	pb "github.com/micro/services/notes/proto"
	"github.com/micro/services/pkg/tenant"
	streamPb "github.com/micro/services/stream/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

// New returns an initialized Notes
func New(c client.Client) *Notes {
	return &Notes{
		Stream: streamPb.NewStreamService("stream", c),
	}
}

// Notes implements the notes proto definition
type Notes struct {
	Stream streamPb.StreamService
}

func newMessage(ev map[string]interface{}) *structpb.Struct {
	v, _ := structpb.NewStruct(ev)
	return v
}

// Create inserts a new note in the store
func (h *Notes) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	if len(req.Title) == 0 && len(req.Text) == 0 {
		return errors.BadRequest("notes.create", "missing title and text")
	}

	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "default"
	}

	// generate a key (uuid v4)
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	t := time.Now().Format(time.RFC3339)
	// set the generated fields on the note
	note := &pb.Note{
		Id:      id.String(),
		Created: t,
		Updated: t,
		Title:   req.Title,
		Text:    req.Text,
	}

	key := fmt.Sprintf("%s:%s", tnt, id)
	rec := store.NewRecord(key, note)

	if err = store.Write(rec); err != nil {
		return errors.InternalServerError("notes.created", "failed to create note")
	}

	// return the note in the response
	rsp.Note = note

	go h.Stream.Publish(ctx, &streamPb.PublishRequest{
		Topic: "notes",
		Message: newMessage(map[string]interface{}{
			"type": "create",
			"note": note,
		}),
	})

	return nil
}

func (h *Notes) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("notes.read", "Missing Note ID")
	}

	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "default"
	}

	key := fmt.Sprintf("%s:%s", tnt, req.Id)

	// read the specific note
	recs, err := store.Read(key)
	if err == store.ErrNotFound {
		return errors.NotFound("notes.read", "Note not found")
	} else if err != nil {
		return errors.InternalServerError("notes.read", "Error reading from store: %v", err.Error())
	}

	// Decode the note
	var note *pb.Note
	if err := recs[0].Decode(&note); err != nil {
		return errors.InternalServerError("notes.update", "Error unmarshaling JSON: %v", err.Error())
	}

	// return the note
	rsp.Note = note

	return nil
}

// Update is a unary API which updates a note in the store
func (h *Notes) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	// Validate the request
	if req.Note == nil {
		return errors.BadRequest("notes.update", "Missing Note")
	}
	if len(req.Note.Id) == 0 {
		return errors.BadRequest("notes.update", "Missing Note ID")
	}

	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "default"
	}

	key := fmt.Sprintf("%s:%s", tnt, req.Note.Id)

	// read the specific note
	recs, err := store.Read(key)
	if err == store.ErrNotFound {
		return errors.NotFound("notes.update", "Note not found")
	} else if err != nil {
		return errors.InternalServerError("notes.update", "Error reading from store: %v", err.Error())
	}

	// Decode the note
	var note *pb.Note
	if err := recs[0].Decode(&note); err != nil {
		return errors.InternalServerError("notes.update", "Error unmarshaling JSON: %v", err.Error())
	}

	// Update the notes title and text
	note.Title = req.Note.Title
	note.Text = req.Note.Text
	note.Updated = time.Now().Format(time.RFC3339)

	rec := store.NewRecord(key, note)

	// Write the updated note to the store
	if err = store.Write(rec); err != nil {
		return errors.InternalServerError("notes.update", "Error writing to store: %v", err.Error())
	}

	go h.Stream.Publish(ctx, &streamPb.PublishRequest{
		Topic: "notes",
		Message: newMessage(map[string]interface{}{
			"type": "update",
			"note": note,
		}),
	})

	rsp.Note = note

	return nil
}

func (h *Notes) Subscribe(ctx context.Context, req *pb.SubscribeRequest, stream pb.Notes_SubscribeStream) error {
	backendStream, err := h.Stream.Subscribe(ctx, &streamPb.SubscribeRequest{
		Topic: "notes",
	})
	if err != nil {
		return errors.InternalServerError("notes.subscribe", "Failed to subscribe to notes")
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		// receive messages from the stream
		msg, err := backendStream.Recv()
		if err != nil {
			return nil
		}

		ev := msg.Message.AsMap()
		note := ev["note"].(*pb.Note)

		// filter if necessary by id
		if len(req.Id) > 0 && note.Id != req.Id {
			continue
		}

		// send back the event to the client
		if err := stream.Send(&pb.Event{
			Type: ev["type"].(string),
			Note: note,
		}); err != nil {
			return nil
		}
	}

	return nil
}

// Delete removes the note from the store, looking up using ID
func (h *Notes) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	// Validate the request
	if len(req.Id) == 0 {
		return errors.BadRequest("notes.delete", "Missing Note ID")
	}

	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "default"
	}

	key := fmt.Sprintf("%s:%s", tnt, req.Id)

	// read the specific note
	recs, err := store.Read(key)
	if err == store.ErrNotFound {
		return nil
	} else if err != nil {
		return errors.InternalServerError("notes.delete", "Error reading from store: %v", err.Error())
	}

	// Decode the note
	var note *pb.Note
	if err := recs[0].Decode(&note); err != nil {
		return errors.InternalServerError("notes.delete", "Error unmarshaling JSON: %v", err.Error())
	}

	// now delete it
	if err := store.Delete(key); err != nil && err != store.ErrNotFound {
		return errors.InternalServerError("notes.delete", "Failed to delete note")
	}

	go h.Stream.Publish(ctx, &streamPb.PublishRequest{
		Topic: "notes",
		Message: newMessage(map[string]interface{}{
			"type": "delete",
			"note": note,
		}),
	})

	rsp.Note = note

	return nil
}

// List returns all of the notes in the store
func (h *Notes) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "default"
	}

	key := fmt.Sprintf("%s:", tnt)

	// Retrieve all of the records in the store
	recs, err := store.Read(key, store.ReadPrefix())
	if err != nil {
		return errors.InternalServerError("notes.list", "Error reading from store: %v", err.Error())
	}

	// Initialize the response notes slice
	rsp.Notes = make([]*pb.Note, len(recs))

	// Unmarshal the notes into the response
	for i, r := range recs {
		if err := r.Decode(&rsp.Notes[i]); err != nil {
			return errors.InternalServerError("notes.list", "Error decoding note: %v", err.Error())
		}
	}

	return nil
}
