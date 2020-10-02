package handler

import (
	"context"

	"github.com/micro/micro/v3/service/logger"
	helloworld "github.com/micro/services/helloworld/proto"
)

type Helloworld struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Helloworld) Call(ctx context.Context, req *helloworld.Request, rsp *helloworld.Response) error {
	logger.Info("Received Helloworld.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}
