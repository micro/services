package main

import (
	"log"

	"github.com/micro/micro/v3/service"
	"github.com/micro/services/location/handler"
	pb "github.com/micro/services/location/proto"
	"github.com/micro/services/pkg/tracing"
)

func main() {
	location := service.New(
		service.Name("location"),
	)

	pb.RegisterLocationHandler(location.Server(), new(handler.Location))

	// TODO reinstate me
	//service.Subscribe(subscriber.Topic, new(subscriber.Location))
	traceCloser := tracing.SetupOpentracing("location")
	defer traceCloser.Close()

	if err := location.Run(); err != nil {
		log.Fatal(err)
	}
}
