package handler_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/events"
	"github.com/micro/services/streams/handler"
	pb "github.com/micro/services/streams/proto"
	"github.com/stretchr/testify/assert"
)

func TestSubscribe(t *testing.T) {
	t.Run("MissingToken", func(t *testing.T) {
		h := testHandler(t)
		s := new(streamMock)

		ctx := auth.ContextWithAccount(context.TODO(), &auth.Account{Issuer: "foo"})
		err := h.Subscribe(ctx, &pb.SubscribeRequest{
			Topic: "helloworld",
		}, s)

		assert.Equal(t, handler.ErrMissingToken, err)
		assert.Empty(t, s.Messages)
	})

	t.Run("MissingTopic", func(t *testing.T) {
		h := testHandler(t)
		s := new(streamMock)

		ctx := auth.ContextWithAccount(context.TODO(), &auth.Account{Issuer: "foo"})
		err := h.Subscribe(ctx, &pb.SubscribeRequest{
			Token: uuid.New().String(),
		}, s)

		assert.Equal(t, handler.ErrMissingTopic, err)
		assert.Empty(t, s.Messages)
	})

	t.Run("InvalidToken", func(t *testing.T) {
		h := testHandler(t)
		s := new(streamMock)

		ctx := auth.ContextWithAccount(context.TODO(), &auth.Account{Issuer: "foo"})
		err := h.Subscribe(ctx, &pb.SubscribeRequest{
			Topic: "helloworld",
			Token: uuid.New().String(),
		}, s)

		assert.Equal(t, handler.ErrInvalidToken, err)
		assert.Empty(t, s.Messages)
	})

	t.Run("ExpiredToken", func(t *testing.T) {
		h := testHandler(t)

		var tRsp pb.TokenResponse
		ctx := auth.ContextWithAccount(context.TODO(), &auth.Account{Issuer: "foo"})
		err := h.Token(ctx, &pb.TokenRequest{
			Topic: "helloworld",
		}, &tRsp)
		assert.NoError(t, err)

		ct := h.Time()
		h.Time = func() time.Time { return ct.Add(handler.TokenTTL * 2) }
		s := new(streamMock)
		err = h.Subscribe(ctx, &pb.SubscribeRequest{
			Topic: "helloworld",
			Token: tRsp.Token,
		}, s)

		assert.Equal(t, handler.ErrExpiredToken, err)
		assert.Empty(t, s.Messages)
	})

	t.Run("ForbiddenTopic", func(t *testing.T) {
		h := testHandler(t)

		var tRsp pb.TokenResponse
		ctx := auth.ContextWithAccount(context.TODO(), &auth.Account{Issuer: "foo"})
		err := h.Token(ctx, &pb.TokenRequest{
			Topic: "helloworldx",
		}, &tRsp)
		assert.NoError(t, err)

		s := new(streamMock)
		err = h.Subscribe(ctx, &pb.SubscribeRequest{
			Topic: "helloworld",
			Token: tRsp.Token,
		}, s)

		assert.Equal(t, handler.ErrForbiddenTopic, err)
		assert.Empty(t, s.Messages)
	})

	t.Run("Valid", func(t *testing.T) {
		defer func() {
			if i := recover(); i != nil {
				t.Logf("%+v", i)
			}
		}()
		h := testHandler(t)
		c := make(chan events.Event)
		h.Events.(*eventsMock).ConsumeChan = c

		var tRsp pb.TokenResponse
		ctx := auth.ContextWithAccount(context.TODO(), &auth.Account{Issuer: "foo"})
		err := h.Token(ctx, &pb.TokenRequest{
			Topic: "helloworld",
		}, &tRsp)
		assert.NoError(t, err)

		s := &streamMock{Messages: []*pb.Message{}}
		wg := sync.WaitGroup{}
		wg.Add(1)
		var subsErr error
		go func() {
			defer wg.Done()
			subsErr = h.Subscribe(ctx, &pb.SubscribeRequest{
				Topic: "helloworld",
				Token: tRsp.Token,
			}, s)
		}()

		e1 := events.Event{
			ID:        uuid.New().String(),
			Topic:     "helloworld",
			Timestamp: h.Time().Add(time.Second * -2),
			Payload:   []byte("abc"),
		}
		e2 := events.Event{
			ID:        uuid.New().String(),
			Topic:     "helloworld",
			Timestamp: h.Time().Add(time.Second * -1),
			Payload:   []byte("123"),
		}

		timeout := time.NewTimer(time.Millisecond * 100).C
		select {
		case <-timeout:
			t.Fatal("Events not consumed from stream")
			return
		case c <- e1:
			t.Log("Event1 consumed")
		}
		select {
		case <-timeout:
			t.Fatal("Events not consumed from stream")
			return
		case c <- e2:
			t.Log("Event2 consumed")
		}

		close(c)
		wg.Wait()
		assert.NoError(t, subsErr)
		assert.Equal(t, "foo.helloworld", h.Events.(*eventsMock).ConsumeTopic)

		// sleep to wait for the subscribe loop to push the message to the stream
		//time.Sleep(1 * time.Second)
		if len(s.Messages) != 2 {
			t.Fatalf("Expected 2 messages, got %v", len(s.Messages))
			return
		}

		assert.Equal(t, e1.Topic, s.Messages[0].Topic)
		assert.Equal(t, string(e1.Payload), s.Messages[0].Message)
		assert.True(t, e1.Timestamp.Equal(time.Unix(s.Messages[0].SentAt, 0)))

		assert.Equal(t, e2.Topic, s.Messages[1].Topic)
		assert.Equal(t, string(e2.Payload), s.Messages[1].Message)
		assert.True(t, e2.Timestamp.Equal(time.Unix(s.Messages[1].SentAt, 0)))
	})

	t.Run("TokenForDifferentIssuer", func(t *testing.T) {
		h := testHandler(t)

		var tRsp pb.TokenResponse
		ctx := auth.ContextWithAccount(context.TODO(), &auth.Account{Issuer: "foo"})
		err := h.Token(ctx, &pb.TokenRequest{
			Topic: "tokfordiff",
		}, &tRsp)
		assert.NoError(t, err)

		s := new(streamMock)
		ctx = auth.ContextWithAccount(context.TODO(), &auth.Account{Issuer: "bar"})
		err = h.Subscribe(ctx, &pb.SubscribeRequest{
			Topic: "tokfordiff",
			Token: tRsp.Token,
		}, s)
		assert.Equal(t, handler.ErrInvalidToken, err)
		assert.Empty(t, s.Messages)
	})

	t.Run("BadTopic", func(t *testing.T) {
		h := testHandler(t)

		var tRsp pb.TokenResponse
		ctx := auth.ContextWithAccount(context.TODO(), &auth.Account{Issuer: "foo"})
		err := h.Token(ctx, &pb.TokenRequest{}, &tRsp)
		assert.NoError(t, err)

		s := new(streamMock)
		ctx = auth.ContextWithAccount(context.TODO(), &auth.Account{Issuer: "bar"})
		err = h.Subscribe(ctx, &pb.SubscribeRequest{
			Topic: "tok-for/diff",
			Token: tRsp.Token,
		}, s)
		assert.Equal(t, handler.ErrInvalidTopic, err)
		assert.Empty(t, s.Messages)
	})

}

type streamMock struct {
	Messages []*pb.Message
	pb.Streams_SubscribeStream
}

func (x *streamMock) Send(m *pb.Message) error {
	x.Messages = append(x.Messages, m)
	return nil
}
