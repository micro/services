package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	pb "github.com/micro/services/chats/proto"
)

// Create a message within a chat
func (c *Chats) SendMessage(ctx context.Context, req *pb.SendMessageRequest, rsp *pb.SendMessageResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.AuthorId) == 0 {
		return ErrMissingAuthorID
	}
	if len(req.ChatId) == 0 {
		return ErrMissingChatID
	}
	if len(req.Text) == 0 {
		return ErrMissingText
	}

	chat := &Chat{
		ID: req.ChatId,
	}

	recs, err := store.Read(chat.Key(ctx), store.ReadLimit(1))
	if err == store.ErrNotFound {
		return ErrNotFound
	} else if err != nil {
		logger.Errorf("Error reading chat: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// create the message
	msg := &Message{
		ID:       req.Id,
		Text:     req.Text,
		AuthorID: req.AuthorId,
		ChatID:   req.ChatId,
		SentAt: c.Time(),
	}
	if len(msg.ID) == 0 {
		msg.ID = uuid.New().String()
	}

	// check if the message already exists
	recs, err = store.Read(msg.Index(ctx), store.ReadLimit(1))
	if err == nil && len(recs) == 1 {
		// return the existing message
		msg = &Message{}
		recs[0].Decode(&msg)
		rsp.Message = msg.Serialize()
		return nil
	}

	// if there's an error then return
	if err != nil && err != store.ErrNotFound {
		logger.Errorf("Error creating message: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// otherwise write the record
	if err := store.Write(store.NewRecord(msg.Key(ctx), msg)); err != nil {
		logger.Errorf("Error creating message: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// write the chat based index
	if err := store.Write(store.NewRecord(msg.Index(ctx), msg)); err == nil {
		rsp.Message = msg.Serialize()
		return nil
	} else if err != nil {
		logger.Errorf("Error creating message: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	return nil
}
