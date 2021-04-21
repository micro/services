package handler_test

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/micro/micro/v3/service/events"
	"github.com/micro/services/streams/handler"
)

func testHandler(t *testing.T) *handler.Streams {
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
	if _, err := sqlDB.Exec("DROP TABLE IF EXISTS micro_users, micro_tokens CASCADE"); err != nil {
		t.Fatalf("Error cleaning database: %v", err)
	}

	h := &handler.Streams{
		Events: new(eventsMock),
		Time: func() time.Time {
			return time.Unix(1612787045, 0)
		},
	}
	h.DBConn(sqlDB).Migrations(&handler.Token{})
	return h
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
