package handler

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	pb "github.com/micro/services/chats/proto"
)

// Create a chat between two or more users, if a chat already exists for these users, the existing
// chat will be returned
func (c *Chats) CreateChat(ctx context.Context, req *pb.CreateChatRequest, rsp *pb.CreateChatResponse) error {
	_, ok := auth.AccountFromContext(ctx)
	if !ok {
		errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}
	// validate the request
	if len(req.UserIds) < 2 {
		return ErrMissingUserIDs
	}

	// sort the user ids
	sort.Strings(req.UserIds)

	id := uuid.New().String()
	if len(req.Id) > 0 {
		id = req.Id
	}

	// construct the chat
	chat := &Chat{
		ID:        id,
		CreatedAt: time.Now(),
		UserIDs:   req.UserIds,
	}

	// read the chat by the unique composition of ids
	recs, err := store.Read(chat.Key(ctx), store.ReadLimit(1))
	if err == nil && len(recs) == 1 {
		// found an existing record
		recs[0].Decode(&chat)
		rsp.Chat = chat.Serialize()
		return nil
	}

	// if not found check it exists by user index key
	if err == store.ErrNotFound {
		recs, err = store.Read(chat.Index(ctx), store.ReadLimit(1))
		if err == nil && len(recs) > 0 {
			recs[0].Decode(&chat)
			rsp.Chat = chat.Serialize()
			return nil
		}
	}

	// ok otherwise we're creating an entirely new record
	newRec := store.NewRecord(chat.Key(ctx), chat)
	if err := store.Write(newRec); err != nil {
		logger.Errorf("Error creating chat: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// write the user composite key
	newRec = store.NewRecord(chat.Index(ctx), chat)
	if err := store.Write(newRec); err != nil {
		logger.Errorf("Error creating chat: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	// return the record
	rsp.Chat = chat.Serialize()

	return nil
}
