package handler_test

import (
	"testing"
	"time"

	"github.com/micro/services/streams/handler"
	pb "github.com/micro/services/streams/proto"
	"github.com/stretchr/testify/assert"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func testHandler(t *testing.T) *handler.Streams {
	// connect to the database
	db, err := gorm.Open(postgres.Open("postgresql://postgres@localhost:5432/postgres?sslmode=disable"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}

	// clean any data from a previous run
	if err := db.Exec("DROP TABLE IF EXISTS conversations, messages CASCADE").Error; err != nil {
		t.Fatalf("Error cleaning database: %v", err)
	}

	// migrate the database
	if err := db.AutoMigrate(&handler.Conversation{}, &handler.Message{}); err != nil {
		t.Fatalf("Error migrating database: %v", err)
	}

	return &handler.Streams{DB: db, Time: func() time.Time { return time.Unix(1611327673, 0) }}
}

func assertConversationsMatch(t *testing.T, exp, act *pb.Conversation) {
	if act == nil {
		t.Errorf("Conversation not returned")
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

	if act.CreatedAt == nil {
		t.Errorf("CreatedAt not set")
		return
	}

	assert.True(t, exp.CreatedAt.AsTime().Equal(act.CreatedAt.AsTime()))
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
	assert.Equal(t, exp.ConversationId, act.ConversationId)

	if act.SentAt == nil {
		t.Errorf("SentAt not set")
		return
	}

	assert.True(t, exp.SentAt.AsTime().Equal(act.SentAt.AsTime()))
}
