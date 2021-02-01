package handler

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/store"
	pb "github.com/micro/services/mail/proto"
)

const (
	messagePrefix = "message"
	joinKey       = "/"
)

type Mail struct{}

// Send a message
func (m *Mail) Send(ctx context.Context, req *pb.SendRequest, rsp *pb.SendResponse) error {
	// validate the request
	if len(req.To) == 0 {
		return errors.BadRequest("mail.Send.MissingTo", "Missing to")
	}
	if len(req.From) == 0 {
		return errors.BadRequest("mail.Send.MissingFrom", "Missing from")
	}
	if len(req.Text) == 0 {
		return errors.BadRequest("mail.Send.MissingText", "Missing text")
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
		return errors.BadRequest("mail.Send.Unknown", "Error encoding message")
	}

	// write the message to the store under the recipients key
	key := strings.Join([]string{messagePrefix, req.To, msg.Id}, joinKey)
	if err := store.Write(&store.Record{Key: key, Value: bytes}); err != nil {
		return errors.BadRequest("mail.Send.Unknown", "Error writing to the store")
	}

	// write the message to the store under the id (so it can be looked up without needing to know
	// the users id)
	key = strings.Join([]string{messagePrefix, msg.Id}, joinKey)
	if err := store.Write(&store.Record{Key: key, Value: bytes}); err != nil {
		return errors.BadRequest("mail.Send.Unknown", "Error writing to the store")
	}

	return nil
}

// List mail for a user
func (m *Mail) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	// validate the request
	if len(req.User) == 0 {
		return errors.BadRequest("mail.List.MissingUser", "Missing user")
	}

	// query the store for any mail sent to the user
	prefix := strings.Join([]string{messagePrefix, req.User}, joinKey)
	recs, err := store.Read(prefix, store.ReadPrefix())
	if err != nil {
		return errors.BadRequest("mail.List.Unknown", "Error reading from the store")
	}

	// serialize the result
	rsp.Mail = make([]*pb.Message, len(recs))
	for i, r := range recs {
		var msg pb.Message
		if err := json.Unmarshal(r.Value, &msg); err != nil {
			return errors.BadRequest("mail.List.Unknown", "Error decoding message")
		}
		rsp.Mail[i] = &msg
	}

	return nil
}

// Read a message
func (m *Mail) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	// validate the request
	if len(req.Id) == 0 {
		return errors.BadRequest("mail.Read.MissingUser", "Missing user")
	}

	// query the store
	key := strings.Join([]string{messagePrefix, req.Id}, joinKey)
	recs, err := store.Read(key)
	if err == store.ErrNotFound {
		return errors.NotFound("message.Read.InvalidID", "Message not found with ID")
	} else if err != nil {
		return errors.BadRequest("mail.Read.Unknown", "Error reading from the store")
	}

	// serialize the response
	var msg pb.Message
	if err := json.Unmarshal(recs[0].Value, &msg); err != nil {
		return errors.BadRequest("mail.Read.Unknown", "Error decoding message")
	}
	rsp.Message = &msg

	return nil
}
