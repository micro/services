package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/services/pkg/tenant"
	pb "github.com/micro/services/threads/proto"
)

var (
	ErrMissingID        = errors.BadRequest("MISSING_ID", "Missing ID")
	ErrMissingGroupID   = errors.BadRequest("MISSING_GROUP_ID", "Missing GroupID")
	ErrMissingTopic     = errors.BadRequest("MISSING_TOPIC", "Missing Topic")
	ErrMissingAuthorID  = errors.BadRequest("MISSING_AUTHOR_ID", "Missing Author ID")
	ErrMissingText      = errors.BadRequest("MISSING_TEXT", "Missing text")
	ErrMissingThreadID  = errors.BadRequest("MISSING_CONVERSATION_ID", "Missing Thread ID")
	ErrMissingThreadIDs = errors.BadRequest("MISSING_CONVERSATION_IDS", "One or more Thread IDs are required")
	ErrNotFound         = errors.NotFound("NOT_FOUND", "Thread not found")
)

type Threads struct {
	Time func() time.Time
}

type Message struct {
	ID       string
	AuthorID string
	ThreadID string
	Text     string
	SentAt   time.Time
}

func (m *Message) Serialize() *pb.Message {
	return &pb.Message{
		Id:       m.ID,
		AuthorId: m.AuthorID,
		ThreadId: m.ThreadID,
		Text:     m.Text,
		SentAt:   m.SentAt.Format(time.RFC3339Nano),
	}
}

type Thread struct {
	ID        string
	GroupID   string
	Topic     string
	CreatedAt time.Time
}

func (c *Thread) Serialize() *pb.Thread {
	return &pb.Thread{
		Id:        c.ID,
		GroupId:   c.GroupID,
		Topic:     c.Topic,
		CreatedAt: c.CreatedAt.Format(time.RFC3339Nano),
	}
}

func ParseTime(v string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, v)
	if err == nil {
		return t
	}
	t, err = time.Parse(time.RFC3339, v)
	if err == nil {
		return t
	}
	return time.Time{}
}

func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339Nano)
}

func (t *Thread) Key(ctx context.Context) string {
	if len(t.ID) == 0 {
		return ""
	}

	key := fmt.Sprintf("thread:%s", t.ID)

	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return key
	}

	return fmt.Sprintf("%s/%s", tnt, key)
}

func (t *Thread) Index(ctx context.Context) string {
	key := fmt.Sprintf("threadsByGroupID:%s:%s", t.GroupID, t.ID)

	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return key
	}

	return fmt.Sprintf("%s/%s", tnt, key)
}

func (t *Thread) Value() interface{} {
	return t
}

func (m *Message) Key(ctx context.Context) string {
	if len(m.ID) == 0 {
		return ""
	}

	key := fmt.Sprintf("message:%s:%s", m.ID, m.ThreadID)

	t, ok := tenant.FromContext(ctx)
	if !ok {
		return key
	}

	return fmt.Sprintf("%s/%s", t, key)
}

func (m *Message) Index(ctx context.Context) string {
	key := fmt.Sprintf("messagesByThreadID:%s", m.ThreadID)

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

func (m *Message) Value() interface{} {
	return m
}
