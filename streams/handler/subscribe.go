package handler

import (
	"context"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/streams/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func (s *Streams) Subscribe(ctx context.Context, req *pb.SubscribeRequest, stream pb.Streams_SubscribeStream) error {
	logger.Infof("Recieved subscribe request. Topic: '%v', Token: '%v'", req.Topic, req.Token)

	// validate the request
	if len(req.Token) == 0 {
		return ErrMissingToken
	}
	if len(req.Topic) == 0 {
		return ErrMissingTopic
	}

	// find the token and check to see if it has expired
	var token Token
	if err := s.DB.Where(&Token{Token: req.Token}).First(&token).Error; err == gorm.ErrRecordNotFound {
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

	// start the subscription
	evChan, err := s.Events.Consume(req.Topic, events.WithGroup(token.Token))
	logger.Infof("Subscribing to %v via queue %v", req.Topic, token.Topic)
	if err != nil {
		logger.Errorf("Error connecting to events stream: %v", err)
		return errors.InternalServerError("EVENTS_ERROR", "Error connecting to events stream")
	}
	go func() {
		defer stream.Close()
		for {
			msg, ok := <-evChan
			if !ok {
				return
			}
			logger.Infof("Sending message to subscriber %v", token.Topic)
			if err := stream.Send(&pb.Message{
				Topic:   msg.Topic,
				Message: string(msg.Payload),
				SentAt:  timestamppb.New(msg.Timestamp),
			}); err != nil {
				return
			}
		}
	}()

	return nil
}
