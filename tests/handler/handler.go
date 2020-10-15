package handler

import (
	"context"

	pb "github.com/m3o/services/tests/proto"
)

type Tests struct{}

func (t *Tests) Register(ctx context.Context, req *pb.RegisterRequest, rsp *pb.RegisterResponse) error {
	// TODO register the test to be run periodically
	return nil
}
