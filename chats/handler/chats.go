package handler

import (
	"encoding/json"
	"time"

	pb "github.com/micro/services/chats/proto"

	"github.com/micro/micro/v3/service/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

var (
	ErrMissingID       = errors.BadRequest("MISSING_ID", "Missing ID")
	ErrMissingAuthorID = errors.BadRequest("MISSING_AUTHOR_ID", "Missing Author ID")
	ErrMissingText     = errors.BadRequest("MISSING_TEXT", "Missing text")
	ErrMissingChatID   = errors.BadRequest("MISSING_CHAT_ID", "Missing Chat ID")
	ErrMissingUserIDs  = errors.BadRequest("MISSING_USER_IDs", "Two or more user IDs are required")
	ErrNotFound        = errors.NotFound("NOT_FOUND", "Chat not found")
)

type Chats struct {
	DB   *gorm.DB
	Time func() time.Time
}

type Chat struct {
	ID        string
	UserIDs   string `gorm:"uniqueIndex"` // sorted json array
	CreatedAt time.Time
}

type Message struct {
	ID       string
	AuthorID string
	ChatID   string
	Text     string
	SentAt   time.Time
}

func (m *Message) Serialize() *pb.Message {
	return &pb.Message{
		Id:       m.ID,
		AuthorId: m.AuthorID,
		ChatId:   m.ChatID,
		Text:     m.Text,
		SentAt:   timestamppb.New(m.SentAt),
	}
}

func (c *Chat) Serialize() *pb.Chat {
	var userIDs []string
	json.Unmarshal([]byte(c.UserIDs), &userIDs)

	return &pb.Chat{
		Id:        c.ID,
		UserIds:   userIDs,
		CreatedAt: timestamppb.New(c.CreatedAt),
	}
}
