package handler

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/store"
	pb "github.com/micro/services/messages/proto"
)

const (
	messagePrefix = "message"
	joinKey       = "/"
)

type Messages struct{}

// Send a message
func (m *Messages) Send(ctx context.Context, req *pb.SendRequest, rsp *pb.SendResponse) error {
	// validate the request
	if len(req.To) == 0 {
		return errors.BadRequest("messages.Send.MissingTo", "Missing to")
	}
	if len(req.From) == 0 {
		return errors.BadRequest("messages.Send.MissingFrom", "Missing from")
	}
	if len(req.Text) == 0 {
		return errors.BadRequest("messages.Send.MissingText", "Missing text")
	}

	// construct the message and marshal it to json
	msg := &pb.Message{
		Id:      uuid.New().String(),
		To:      req.To,
		From:    req.From,
		Subject: req.Subject,
		Text:    req.Text,
		SentAt:  time.Now().Unix(),
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		return errors.BadRequest("messages.Send.Unknown", "Error encoding message")
	}

	// write the message to the store under the recipients key
	key := strings.Join([]string{messagePrefix, req.To, msg.Id}, joinKey)
	if err := store.Write(&store.Record{Key: key, Value: bytes}); err != nil {
		return errors.BadRequest("messages.Send.Unknown", "Error writing to the store")
	}

	// write the message to the store under the id (so it can be looked up without needing to know
	// the users id)
	key = strings.Join([]string{messagePrefix, msg.Id}, joinKey)
	if err := store.Write(&store.Record{Key: key, Value: bytes}); err != nil {
		return errors.BadRequest("messages.Send.Unknown", "Error writing to the store")
	}

	return nil
}

// List messages for a user
func (m *Messages) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	// validate the request
	if len(req.User) == 0 {
		return errors.BadRequest("messages.List.MissingUser", "Missing user")
	}

	// query the store for any messages sent to the user
	prefix := strings.Join([]string{messagePrefix, req.User}, joinKey)
	recs, err := store.Read("", store.Prefix(prefix))
	if err != nil {
		return errors.BadRequest("messages.List.Unknown", "Error reading from the store")
	}

	// serialize the result
	rsp.Messages = make([]*pb.Message, len(recs))
	for i, r := range recs {
		var msg pb.Message
		if err := json.Unmarshal(r.Value, &msg); err != nil {
			return errors.BadRequest("messages.List.Unknown", "Error decoding message")
		}
		rsp.Messages[i] = &msg
	}

	return nil
}

// Read a message
func (m *Messages) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	// validate the request
	if len(req.Id) == 0 {
		return errors.BadRequest("messages.Read.MissingUser", "Missing user")
	}

	// query the store
	key := strings.Join([]string{messagePrefix, req.Id}, joinKey)
	recs, err := store.Read(key)
	if err == store.ErrNotFound {
		return errors.NotFound("message.Read.InvalidID", "Message not found with ID")
	} else if err != nil {
		return errors.BadRequest("messages.Read.Unknown", "Error reading from the store")
	}

	// serialize the response
	var msg pb.Message
	if err := json.Unmarshal(recs[0].Value, &msg); err != nil {
		return errors.BadRequest("messages.Read.Unknown", "Error decoding message")
	}
	rsp.Message = &msg

	return nil
}
