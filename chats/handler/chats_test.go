package handler_test

import (
	"os"
	"testing"
	"time"

	"github.com/micro/services/chats/handler"
	pb "github.com/micro/services/chats/proto"
	"github.com/stretchr/testify/assert"

	"github.com/golang/protobuf/ptypes/timestamp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func testHandler(t *testing.T) *handler.Chats {
	// connect to the database
	addr := os.Getenv("POSTGRES_URL")
	if len(addr) == 0 {
		addr = "postgresql://postgres@localhost:5432/postgres?sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(addr), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}

	// clean any data from a previous run
	if err := db.Exec("DROP TABLE IF EXISTS chats, messages CASCADE").Error; err != nil {
		t.Fatalf("Error cleaning database: %v", err)
	}

	// migrate the database
	if err := db.AutoMigrate(&handler.Chat{}, &handler.Message{}); err != nil {
		t.Fatalf("Error migrating database: %v", err)
	}

	return &handler.Chats{DB: db, Time: func() time.Time { return time.Unix(1611327673, 0) }}
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

	if act.CreatedAt == nil {
		t.Errorf("CreatedAt not set")
		return
	}

	assert.True(t, microSecondTime(exp.CreatedAt).Equal(microSecondTime(act.CreatedAt)))
}

// postgres has a resolution of 100microseconds so just test that it's accurate to the second
func microSecondTime(t *timestamp.Timestamp) time.Time {
	tt := t.AsTime()
	return time.Unix(tt.Unix(), int64(tt.Nanosecond()-tt.Nanosecond()%1000))
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

	if act.SentAt == nil {
		t.Errorf("SentAt not set")
		return
	}
	assert.True(t, microSecondTime(exp.SentAt).Equal(microSecondTime(act.SentAt)))
}
