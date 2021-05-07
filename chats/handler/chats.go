package handler

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	pb "github.com/micro/services/chats/proto"
	"github.com/micro/services/pkg/tenant"

	"github.com/micro/micro/v3/service/errors"
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
	Time func() time.Time
}

type Chat struct {
	ID        string
	UserIDs   []string
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
		SentAt:   m.SentAt.UnixNano(),
	}
}

func (c *Chat) Index(ctx context.Context) string {
	sort.Strings(c.UserIDs)
	users := strings.Join(c.UserIDs, "-")

	key := fmt.Sprintf("chatByUserIDs:%s", users)

	t, ok := tenant.FromContext(ctx)
	if !ok {
		return key
	}

	return fmt.Sprintf("%s/%s", t, key)
}

func (c *Chat) Key(ctx context.Context) string {
	key := fmt.Sprintf("chat:%s", c.ID)

	t, ok := tenant.FromContext(ctx)
	if !ok {
		return key
	}

	return fmt.Sprintf("%s/%s", t, key)
}

func (m *Message) Key(ctx context.Context) string {
	key := fmt.Sprintf("message:%s:%s", m.ChatID, m.ID)

	t, ok := tenant.FromContext(ctx)
	if !ok {
		return key
	}

	return fmt.Sprintf("%s/%s", t, key)
}

func (m *Message) Index(ctx context.Context) string {
	key := fmt.Sprintf("messagesByChatID:%s", m.ChatID)

	if !m.SentAt.IsZero() {
		key = fmt.Sprintf("%s:%d", key, m.SentAt.UnixNano())

		if len(m.ID) > 0 {
			key = fmt.Sprintf("%s:%s", key, m.ID)
		}
	}

	t, ok := tenant.FromContext(ctx)
	if !ok {
		return key
	}

	return fmt.Sprintf("%s/%s", t, key)
}

func (c *Chat) Serialize() *pb.Chat {
	return &pb.Chat{
		Id:        c.ID,
		UserIds:   c.UserIDs,
		CreatedAt: c.CreatedAt.UnixNano(),
	}
}
