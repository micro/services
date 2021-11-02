package handler

import (
	"context"
	"fmt"
	"path"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/events"
	log "github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/event/proto"
	"github.com/micro/services/pkg/tenant"
	"google.golang.org/protobuf/types/known/structpb"
)

type Event struct{}

func (s *Event) Publish(ctx context.Context, req *pb.PublishRequest, rsp *pb.PublishResponse) error {
	if len(req.Topic) == 0 {
		return errors.BadRequest("event.publish", "topic is blank")
	}

	// get the tenant
	id, ok := tenant.FromContext(ctx)
	if !ok {
		id = "default"
	}

	// create tenant based topics
	topic := path.Join("event", id, req.Topic)

	log.Infof("Tenant %v publishing to %v\n", id, req.Topic)

	// publish the message
	return events.Publish(topic, req.Message.AsMap())
}

func (s *Event) Subscribe(ctx context.Context, req *pb.SubscribeRequest, stream pb.Event_SubscribeStream) error {
	if len(req.Topic) == 0 {
		return errors.BadRequest("event.publish", "topic is blank")
	}

	id, ok := tenant.FromContext(ctx)
	if !ok {
		id = "default"
	}

	// create tenant based topics
	topic := path.Join("event", id, req.Topic)

	log.Infof("Tenant %v subscribing to %v\n", id, req.Topic)

	// check if a group os provided
	opts := []events.ConsumeOption{}
	if len(req.Group) > 0 {
		opts = append(opts, events.WithGroup(req.Group))
	}

	sub, err := events.Consume(topic, opts...)
	if err != nil {
		return errors.InternalServerError("event.subscribe", "failed to subscribe to event")
	}

	// range over the messages until the subscriber is closed
	for msg := range sub {
		fmt.Println("got message, sending")
		// unmarshal the message into a struct
		d := &structpb.Struct{}
		d.UnmarshalJSON(msg.Payload)

		if err := stream.Send(&pb.SubscribeResponse{
			Topic:   req.Topic,
			Message: d,
		}); err != nil {
			return err
		}
	}

	return nil
}
