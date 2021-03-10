package handler_test

import (
	"os"
	"testing"
	"time"

	"github.com/micro/micro/v3/service/events"
	"github.com/micro/services/streams/handler"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func testHandler(t *testing.T) *handler.Streams {
	// connect to the database
	addr := os.Getenv("POSTGRES_URL")
	if len(addr) == 0 {
		addr = "postgresql://postgres@localhost:5432/postgres?sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(addr), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}

	// migrate the database
	if err := db.AutoMigrate(&handler.Token{}); err != nil {
		t.Fatalf("Error migrating database: %v", err)
	}

	// clean any data from a previous run
	if err := db.Exec("TRUNCATE TABLE tokens CASCADE").Error; err != nil {
		t.Fatalf("Error cleaning database: %v", err)
	}

	return &handler.Streams{
		DB:     db,
		Events: new(eventsMock),
		Time: func() time.Time {
			return time.Unix(1612787045, 0)
		},
	}
}

type eventsMock struct {
	PublishCount   int
	PublishTopic   string
	PublishMessage interface{}

	ConsumeTopic string
	ConsumeChan  <-chan events.Event
}

func (e *eventsMock) Publish(topic string, msg interface{}, opts ...events.PublishOption) error {
	e.PublishCount++
	e.PublishTopic = topic
	e.PublishMessage = msg
	return nil
}

func (e *eventsMock) Consume(topic string, opts ...events.ConsumeOption) (<-chan events.Event, error) {
	e.ConsumeTopic = topic
	return e.ConsumeChan, nil
}
