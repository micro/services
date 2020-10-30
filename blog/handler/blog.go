package handler

import (
	"context"

	log "github.com/micro/micro/v3/service/logger"

	blog "blog/proto"
)

type Blog struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Blog) Call(ctx context.Context, req *blog.Request, rsp *blog.Response) error {
	log.Info("Received Blog.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Blog) Stream(ctx context.Context, req *blog.StreamingRequest, stream blog.Blog_StreamStream) error {
	log.Infof("Received Blog.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&blog.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Blog) PingPong(ctx context.Context, stream blog.Blog_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&blog.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
