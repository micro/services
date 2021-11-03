package handler

import (
	"context"
	"path"
	"strings"
	"time"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/services/pkg/tenant"
	"github.com/micro/services/stream/domain"
	pb "github.com/micro/services/stream/proto"
)

type Stream struct{}

func New() *Stream {
	domain.Setup()
	return &Stream{}
}

func (s *Stream) SendMessage(ctx context.Context, req *pb.SendMessageRequest, rsp *pb.SendMessageResponse) error {
	if len(req.Channel) == 0 {
		return errors.BadRequest("stream.sendmessage", "channel is blank")
	}
	if len(req.Text) == 0 {
		return errors.BadRequest("stream.sendmessage", "message is blank")
	}

	// get the tenant
	id, ok := tenant.FromContext(ctx)
	if !ok {
		id = "default"
	}

	// create tenant based channels
	channel := path.Join(id, req.Channel)

	// sendmessage the message
	if err := domain.SendMessage(channel, req.Text); err != nil {
		return errors.InternalServerError("stream.sendmessage", err.Error())
	}

	return nil
}

func (s *Stream) ListMessages(ctx context.Context, req *pb.ListMessagesRequest, rsp *pb.ListMessagesResponse) error {
	if len(req.Channel) == 0 {
		return errors.BadRequest("stream.sendmessage", "channel is blank")
	}
	if req.Limit <= 0 {
		req.Limit = 25
	}

	id, ok := tenant.FromContext(ctx)
	if !ok {
		id = "default"
	}

	// create tenant based channels
	channel := path.Join(id, req.Channel)
	rsp.Channel = req.Channel

	for _, message := range domain.ListMessages(channel, int64(req.Limit)) {
		metadata := map[string]string{}

		if message.Metadata != nil {
			metadata["created"] = time.Unix(0, message.Metadata.Created).Format(time.RFC3339Nano)
			metadata["title"] = message.Metadata.Title
			metadata["description"] = message.Metadata.Description
			metadata["type"] = message.Metadata.Type
			metadata["image"] = message.Metadata.Image
			metadata["url"] = message.Metadata.Url
			metadata["site"] = message.Metadata.Site
		}

		rsp.Messages = append(rsp.Messages, &pb.Message{
			Id:        message.Id,
			Text:      message.Text,
			Timestamp: time.Unix(0, message.Created).Format(time.RFC3339Nano),
			Channel:   req.Channel,
			Metadata:  metadata,
		})
	}

	return nil
}

func (s *Stream) ListChannels(ctx context.Context, req *pb.ListChannelsRequest, rsp *pb.ListChannelsResponse) error {
	// get the tenant
	id, ok := tenant.FromContext(ctx)
	if !ok {
		id = "default"
	}

	for channel, active := range domain.ListChannels() {
		if !strings.HasPrefix(channel, id+"/") {
			continue
		}

		channel = strings.TrimPrefix(channel, id+"/")

		rsp.Channels = append(rsp.Channels, &pb.Channel{
			Name:       channel,
			LastActive: time.Unix(0, active).Format(time.RFC3339Nano),
		})
	}

	return nil
}
