package handler_test

import (
	"context"
	"testing"
	"time"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/store"
	"github.com/micro/micro/v3/service/store/memory"
	"github.com/micro/services/chats/handler"
	pb "github.com/micro/services/chats/proto"
	"github.com/stretchr/testify/assert"
)

func testHandler(t *testing.T) *handler.Chats {
	store.DefaultStore = memory.NewStore()
	return &handler.Chats{Time: func() time.Time { return time.Unix(1611327673, 0) }}
}

func assertChatsMatch(t *testing.T, exp, act *pb.Chat) {
	if act == nil {
		t.Errorf("Chat not returned")
		return
	}

	// adapt this check so we can reuse the func in testing create, where we don't know the exact id
	// which will be generated
	if len(exp.Id) > 0 {
		assert.Equal(t, exp.Id, act.Id)
	} else {
		assert.NotEmpty(t, act.Id)
	}

	assert.Equal(t, exp.UserIds, act.UserIds)

	if len(act.CreatedAt) == 0 {
		t.Errorf("CreatedAt not set")
		return
	}

	assert.True(t, exp.CreatedAt == act.CreatedAt)
}

func assertMessagesMatch(t *testing.T, exp, act *pb.Message) {
	if act == nil {
		t.Errorf("Message not returned")
		return
	}

	// adapt these checks so we can reuse the func in testing create, where we don't know the exact id /
	// idempotent_id which will be generated
	if len(exp.Id) > 0 {
		assert.Equal(t, exp.Id, act.Id)
	} else {
		assert.NotEmpty(t, act.Id)
	}
	assert.Equal(t, exp.Text, act.Text)
	assert.Equal(t, exp.AuthorId, act.AuthorId)
	assert.Equal(t, exp.ChatId, act.ChatId)

	if len(act.SentAt) == 0 {
		t.Errorf("SentAt not set")
		return
	}

	assert.True(t, exp.SentAt == act.SentAt)
}

func microAccountCtx() context.Context {
	return auth.ContextWithAccount(context.TODO(), &auth.Account{
		Issuer: "micro",
	})
}
