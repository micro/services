package handler_test

import (
	"testing"
	"time"

	"github.com/micro/micro/v3/service/events"
	"github.com/micro/services/streams/handler"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func testHandler(t *testing.T) *handler.Streams {
	// use an in memory DB
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}

	// migrate the database
	if err := db.AutoMigrate(&handler.Token{}); err != nil {
		t.Fatalf("Error migrating database: %v", err)
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
