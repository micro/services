package handler

import (
	"context"

	pb "github.com/micro/services/function/proto"
)

type Function struct {
	key string
}

func New(apiKey string) *Function {
	return &Function{
		key: apiKey,
	}
}

func (f *Function) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	return nil
}
func (f *Function) Call(ctx context.Context, req *pb.CallRequest, rsp *pb.CallResponse) error {
	return nil
}
func (f *Function) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	return nil
}
func (f *Function) Describe(ctx context.Context, req *pb.DescribeRequest, rsp *pb.DescribeResponse) error {
	return nil
}
func (f *Function) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	return nil
}
