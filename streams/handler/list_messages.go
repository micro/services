package handler

import (
	"context"

	pb "github.com/micro/services/streams/proto"
)

// List the messages within a conversation in reverse chronological order, using sent_before to
// offset as older messages need to be loaded
func (s *Streams) ListMessages(ctx context.Context, req *pb.ListMessagesRequest, rsp *pb.ListMessagesResponse) error {
	return nil
}
