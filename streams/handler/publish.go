package handler

import (
	"context"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/streams/proto"
)

func (s *Streams) Publish(ctx context.Context, req *pb.Message, rsp *pb.PublishResponse) error {
	// validate the request
	if len(req.Topic) == 0 {
		return ErrMissingTopic
	}
	if err := validateTopicInput(req.Topic); err != nil {
		return err
	}
	if len(req.Message) == 0 {
		return ErrMissingMessage
	}
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("UNAUTHORIZED", "Unauthorized")
	}

	// publish the message
	logger.Infof("Publishing message to topic: %v", req.Topic)
	return s.Events.Publish(fmtTopic(acc, req.Topic), req.Message)
}
