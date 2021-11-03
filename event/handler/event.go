package handler

import (
	"context"
	"fmt"
	"path"
	"time"

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

func (s *Event) Consume(ctx context.Context, req *pb.ConsumeRequest, stream pb.Event_ConsumeStream) error {
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
	offset := time.Now()
	if len(req.Group) > 0 {
		opts = append(opts, events.WithGroup(req.Group))
	}
	if len(req.Offset) > 0 {
		t, err := time.Parse(time.RFC3339Nano, req.Offset)
		if err == nil {
			offset = t
		}
	}
	opts = append(opts, events.WithOffset(offset))

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

		if err := stream.Send(&pb.ConsumeResponse{
			Topic:     req.Topic,
			Id:        msg.ID,
			Timestamp: msg.Timestamp.Format(time.RFC3339Nano),
			Message:   d,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *Event) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	if len(req.Topic) == 0 {
		return errors.BadRequest("event.read", "topic is blank")
	}

	id, ok := tenant.FromContext(ctx)
	if !ok {
		id = "default"
	}

	// create tenant based topics
	topic := path.Join("event", id, req.Topic)

	log.Infof("Tenant %v reading %v\n", id, req.Topic)
	limit := uint(25)
	offset := uint(0)

	if req.Limit > 0 {
		limit = uint(req.Limit)
	}

	if req.Offset > 0 {
		offset = uint(req.Offset)
	}

	events, err := events.Read(topic, events.ReadLimit(limit), events.ReadOffset(offset))
	if err != nil {
		return err
	}

	for _, ev := range events {
		// unmarshal the message into a struct
		d := &structpb.Struct{}
		d.UnmarshalJSON(ev.Payload)

		rsp.Events = append(rsp.Events, &pb.Ev{
			Id:        ev.ID,
			Timestamp: ev.Timestamp.Format(time.RFC3339Nano),
			Message:   d,
		})
	}

	return nil
}
