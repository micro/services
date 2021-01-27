package handler

import (
	"context"
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/chats/proto"
)

// Create a chat between two or more users, if a chat already exists for these users, the existing
// chat will be returned
func (c *Chats) CreateChat(ctx context.Context, req *pb.CreateChatRequest, rsp *pb.CreateChatResponse) error {
	// validate the request
	if len(req.UserIds) < 2 {
		return ErrMissingUserIDs
	}

	// sort the user ids and then marshal to json
	sort.Strings(req.UserIds)
	bytes, err := json.Marshal(req.UserIds)
	if err != nil {
		logger.Errorf("Error mashaling user ids: %v", err)
		return errors.InternalServerError("ENCODING_ERROR", "Error encoding user ids")
	}

	// construct the chat
	chat := Chat{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		UserIDs:   string(bytes),
	}

	// write to the database, if we get a unique key error, the chat already exists
	err = c.DB.Create(&chat).Error
	if err == nil {
		rsp.Chat = chat.Serialize()
		return nil
	}
	if !strings.Contains(err.Error(), "idx_chats_user_ids") {
		logger.Errorf("Error creating chat: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}

	var existing Chat
	if err := c.DB.Where(&Chat{UserIDs: chat.UserIDs}).First(&existing).Error; err != nil {
		logger.Errorf("Error reading chat: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error connecting to database")
	}
	rsp.Chat = existing.Serialize()
	return nil
}
