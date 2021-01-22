package handler

import (
	"context"

	pb "github.com/micro/services/streams/proto"
)

// RecentMessages returns the most recent messages in a group of conversations. By default the
// most messages retrieved per conversation is 10, however this can be overriden using the
// limit_per_conversation option
func (s *Streams) RecentMessages(ctx context.Context, req *pb.RecentMessagesRequest, rsp *pb.RecentMessagesResponse) error {
	return nil
}
