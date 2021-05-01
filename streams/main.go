package main

import (
	"time"

	"github.com/micro/services/streams/handler"
	pb "github.com/micro/services/streams/proto"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("streams"),
	)

	h := &handler.Streams{
		Events: events.DefaultStream,
		Time:   time.Now,
	}

	// Register handler
	pb.RegisterStreamsHandler(srv.Server(), h)

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
