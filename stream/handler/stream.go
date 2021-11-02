package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"path"

	"github.com/asim/mq/broker"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/pkg/tenant"
	pb "github.com/micro/services/stream/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type Stream struct{}

func (s *Stream) Publish(ctx context.Context, req *pb.PublishRequest, rsp *pb.PublishResponse) error {
	if len(req.Topic) == 0 {
		return errors.BadRequest("stream.publish", "topic is blank")
	}

	// get the tenant
	id, ok := tenant.FromContext(ctx)
	if !ok {
		id = "default"
	}

	// create tenant based topics
	topic := path.Join("stream", id, req.Topic)

	// marshal the data
	b, _ := json.Marshal(req.Message.AsMap())

	log.Infof("Tenant %v publishing to %v\n", id, req.Topic)

	// publish the message
	broker.Publish(topic, b)

	return nil
}

func (s *Stream) Subscribe(ctx context.Context, req *pb.SubscribeRequest, stream pb.Stream_SubscribeStream) error {
	if len(req.Topic) == 0 {
		return errors.BadRequest("stream.publish", "topic is blank")
	}

	id, ok := tenant.FromContext(ctx)
	if !ok {
		id = "default"
	}

	// create tenant based topics
	topic := path.Join("stream", id, req.Topic)

	log.Infof("Tenant %v subscribing to %v\n", id, req.Topic)

	sub, err := broker.Subscribe(topic)
	if err != nil {
		return errors.InternalServerError("stream.subscribe", "failed to subscribe to stream")
	}
	defer broker.Unsubscribe(req.Topic, sub)

	// range over the messages until the subscriber is closed
	for msg := range sub {
		fmt.Println("got message, sending")
		// unmarshal the message into a struct
		d := &structpb.Struct{}
		d.UnmarshalJSON(msg)

		if err := stream.Send(&pb.SubscribeResponse{
			Topic:   req.Topic,
			Message: d,
		}); err != nil {
			return err
		}
	}

	return nil
}
