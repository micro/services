package handler

import (
	"time"

	"github.com/micro/micro/v3/service/errors"
	pb "github.com/micro/services/threads/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

var (
	ErrMissingID              = errors.BadRequest("MISSING_ID", "Missing ID")
	ErrMissingGroupID         = errors.BadRequest("MISSING_GROUP_ID", "Missing GroupID")
	ErrMissingTopic           = errors.BadRequest("MISSING_TOPIC", "Missing Topic")
	ErrMissingAuthorID        = errors.BadRequest("MISSING_AUTHOR_ID", "Missing Author ID")
	ErrMissingText            = errors.BadRequest("MISSING_TEXT", "Missing text")
	ErrMissingConversationID  = errors.BadRequest("MISSING_CONVERSATION_ID", "Missing Conversation ID")
	ErrMissingConversationIDs = errors.BadRequest("MISSING_CONVERSATION_IDS", "One or more Conversation IDs are required")
	ErrNotFound               = errors.NotFound("NOT_FOUND", "Conversation not found")
)

type Threads struct {
	DB   *gorm.DB
	Time func() time.Time
}

type Message struct {
	ID             string
	AuthorID       string
	ConversationID string
	Text           string
	SentAt         time.Time
}

func (m *Message) Serialize() *pb.Message {
	return &pb.Message{
		Id:             m.ID,
		AuthorId:       m.AuthorID,
		ConversationId: m.ConversationID,
		Text:           m.Text,
		SentAt:         timestamppb.New(m.SentAt),
	}
}

type Conversation struct {
	ID        string
	GroupID   string
	Topic     string
	CreatedAt time.Time
}

func (c *Conversation) Serialize() *pb.Conversation {
	return &pb.Conversation{
		Id:        c.ID,
		GroupId:   c.GroupID,
		Topic:     c.Topic,
		CreatedAt: timestamppb.New(c.CreatedAt),
	}
}
