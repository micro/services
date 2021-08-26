package handler

import (
	"context"

	"github.com/micro/micro/v3/service/logger"
	helloworld "github.com/micro/services/helloworld/proto"
)

type Helloworld struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Helloworld) Call(ctx context.Context, req *helloworld.CallRequest, rsp *helloworld.CallResponse) error {
	logger.Info("Received Helloworld.Call request")
	rsp.Message = "Hello " + req.Name
	return nil
}

func (e *Helloworld) Stream(ctx context.Context, req *helloworld.StreamRequest, rsp helloworld.Helloworld_StreamStream) error {
	// send one if none
	if req.Messages == 0 {
		req.Messages = 1
	}

	for i := 0; i < int(req.Messages); i++ {
		rsp.Send(&helloworld.StreamResponse{
			Message: "Hello " + req.Name,
		})
	}

	return nil
}
