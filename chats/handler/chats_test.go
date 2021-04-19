package handler_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/services/chats/handler"
	pb "github.com/micro/services/chats/proto"
	"github.com/stretchr/testify/assert"
)

func testHandler(t *testing.T) *handler.Chats {
	// connect to the database
	addr := os.Getenv("POSTGRES_URL")
	if len(addr) == 0 {
		addr = "postgresql://postgres@localhost:5432/postgres?sslmode=disable"
	}
	sqlDB, err := sql.Open("pgx", addr)
	if err != nil {
		t.Fatalf("Failed to open connection to DB %s", err)
	}

	// clean any data from a previous run
	if _, err := sqlDB.Exec("DROP TABLE IF EXISTS micro_chats, micro_messages CASCADE"); err != nil {
		t.Fatalf("Error cleaning database: %v", err)
	}

	h := &handler.Chats{Time: func() time.Time { return time.Unix(1611327673, 0) }}
	h.DBConn(sqlDB).Migrations(&handler.Chat{}, &handler.Message{})
	return h
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

	if act.CreatedAt == 0 {
		t.Errorf("CreatedAt not set")
		return
	}

	assert.True(t, microSecondTime(exp.CreatedAt).Equal(microSecondTime(act.CreatedAt)))
}

// postgres has a resolution of 100microseconds so just test that it's accurate to the second
func microSecondTime(t int64) time.Time {
	tt := time.Unix(t, 0)
	return time.Unix(t, int64(tt.Nanosecond()-tt.Nanosecond()%1000))
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

	if act.SentAt == 0 {
		t.Errorf("SentAt not set")
		return
	}
	assert.True(t, microSecondTime(exp.SentAt).Equal(microSecondTime(act.SentAt)))
}

func microAccountCtx() context.Context {
	return auth.ContextWithAccount(context.TODO(), &auth.Account{
		Issuer: "micro",
	})
}
