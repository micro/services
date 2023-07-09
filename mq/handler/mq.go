package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"path"

	"github.com/asim/mq/broker"
	pb "github.com/micro/services/mq/proto"
	"github.com/micro/services/pkg/tenant"
	"google.golang.org/protobuf/types/known/structpb"
	"micro.dev/v4/service/errors"
	log "micro.dev/v4/service/logger"
)

type Mq struct{}

func (mq *Mq) Publish(ctx context.Context, req *pb.PublishRequest, rsp *pb.PublishResponse) error {
	if len(req.Topic) == 0 {
		return errors.BadRequest("mq.publish", "topic is blank")
	}

	// get the tenant
	id, ok := tenant.FromContext(ctx)
	if !ok {
		id = "default"
	}

	// create tenant based topics
	topic := path.Join("event", id, req.Topic)

	// marshal the data
	b, _ := json.Marshal(req.Message.AsMap())

	log.Infof("Tenant %v publishing to %v\n", id, req.Topic)

	// publish the message
	broker.Publish(topic, b)

	return nil
}

func (mq *Mq) Subscribe(ctx context.Context, req *pb.SubscribeRequest, stream pb.Mq_SubscribeStream) error {
	if len(req.Topic) == 0 {
		return errors.BadRequest("mq.publish", "topic is blank")
	}

	id, ok := tenant.FromContext(ctx)
	if !ok {
		id = "default"
	}

	// create tenant based topics
	topic := path.Join("event", id, req.Topic)

	log.Infof("Tenant %v subscribing to %v\n", id, req.Topic)

	sub, err := broker.Subscribe(topic)
	if err != nil {
		return errors.InternalServerError("mq.subscribe", "failed to subscribe to mq")
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
