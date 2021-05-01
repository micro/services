package handler

import (
	"context"
	"fmt"
	"io"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	pb "github.com/micro/services/streams/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Streams) Subscribe(ctx context.Context, req *pb.SubscribeRequest, stream pb.Streams_SubscribeStream) error {
	logger.Infof("Received subscribe request. Topic: '%v', Token: '%v'", req.Topic, req.Token)

	// validate the request
	if len(req.Token) == 0 {
		return ErrMissingToken
	}
	if len(req.Topic) == 0 {
		return ErrMissingTopic
	}
	if err := validateTopicInput(req.Topic); err != nil {
		return err
	}

	// find the token and check to see if it has expired
	var token Token
	if err := s.Cache.Get(req.Token, &token); err == store.ErrNotFound {
		return ErrInvalidToken
	} else if err != nil {
		logger.Errorf("Error reading token from store: %v", err)
		return errors.InternalServerError("DATABASE_ERROR", "Error reading token from database")
	}
	if token.ExpiresAt.Before(s.Time()) {
		return ErrExpiredToken
	}

	// if the token was scoped to a channel, ensure the channel is the one being requested
	if len(token.Topic) > 0 && token.Topic != req.Topic {
		return ErrForbiddenTopic
	}

	var topic string

	// attempt to create a unique topic for the account
	acc, ok := auth.AccountFromContext(ctx)
	if ok {
		topic = fmtTopic(acc, req.Topic)
	} else if len(token.Account) > 0 {
		// use the account in the token if present
		topic = fmt.Sprintf("%s.%s", token.Account, token.Topic)
	} else {
		topic = req.Topic
	}

	// start the subscription
	logger.Infof("Subscribing to %v via queue %v", req.Topic, token.Token)

	evChan, err := s.Events.Consume(topic, events.WithGroup(token.Token))
	if err != nil {
		logger.Errorf("Error connecting to events stream: %v", err)
		return errors.InternalServerError("EVENTS_ERROR", "Error connecting to events stream")
	}

	for {
		msg, ok := <-evChan
		if !ok {
			return nil
		}

		logger.Infof("Sending message to subscriber %v", token.Topic)
		pbMsg := &pb.Message{
			Topic:   req.Topic, // use req.Topic not msg.Topic because topic is munged for multitenancy
			Message: string(msg.Payload),
			SentAt:  timestamppb.New(msg.Timestamp),
		}

		if err := stream.Send(pbMsg); err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
	}
}
