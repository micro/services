package handler_test

import (
	"context"
	"testing"
	"time"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/store"
	"github.com/micro/micro/v3/service/store/memory"
	"github.com/micro/services/threads/handler"
	pb "github.com/micro/services/threads/proto"
	"github.com/stretchr/testify/assert"
)

func testHandler(t *testing.T) *handler.Threads {
	store.DefaultStore = memory.NewStore()
	return &handler.Threads{Time: func() time.Time { return time.Unix(1611327673, 0) }}
}

func assertThreadsMatch(t *testing.T, exp, act *pb.Thread) {
	if act == nil {
		t.Errorf("Thread not returned")
		return
	}

	// adapt this check so we can reuse the func in testing create, where we don't know the exact id
	// which will be generated
	if len(exp.Id) > 0 {
		assert.Equal(t, exp.Id, act.Id)
	} else {
		assert.NotEmpty(t, act.Id)
	}

	assert.Equal(t, exp.Topic, act.Topic)
	assert.Equal(t, exp.GroupId, act.GroupId)

	if act.CreatedAt == "" {
		t.Errorf("CreatedAt not set")
		return
	}

	assert.True(t, microSecondTime(exp.CreatedAt).Equal(microSecondTime(act.CreatedAt)))
}

func assertMessagesMatch(t *testing.T, exp, act *pb.Message) {
	if act == nil {
		t.Errorf("Message not returned")
		return
	}

	// adapt this check so we can reuse the func in testing create, where we don't know the exact id
	// which will be generated
	if len(exp.Id) > 0 {
		assert.Equal(t, exp.Id, act.Id)
	} else {
		assert.NotEmpty(t, act.Id)
	}

	assert.Equal(t, exp.Text, act.Text)
	assert.Equal(t, exp.AuthorId, act.AuthorId)
	assert.Equal(t, exp.ThreadId, act.ThreadId)

	if act.SentAt == "" {
		t.Errorf("SentAt not set")
		return
	}

	assert.True(t, microSecondTime(exp.SentAt).Equal(microSecondTime(act.SentAt)))
}

// postgres has a resolution of 100microseconds so just test that it's accurate to the second
func microSecondTime(t string) time.Time {
	tt := handler.ParseTime(t)
	return time.Unix(tt.Unix(), int64(tt.Nanosecond()-tt.Nanosecond()%1000))
}

func microAccountCtx() context.Context {
	return auth.ContextWithAccount(context.TODO(), &auth.Account{
		Issuer: "micro",
	})
}
