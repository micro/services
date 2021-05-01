package handler_test

import (
	"testing"
	"time"

	"github.com/micro/micro/v3/service/store"
	"github.com/micro/micro/v3/service/store/memory"
	"github.com/micro/micro/v3/service/events"
	"github.com/micro/services/streams/handler"
)

func testHandler(t *testing.T) *handler.Streams {
	store.DefaultStore = memory.NewStore()

	h := &handler.Streams{
		Events: new(eventsMock),
		Time: func() time.Time {
			return time.Unix(1612787045, 0)
		},
	}
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
