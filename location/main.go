package main

import (
	"log"

	"github.com/micro/micro/v5/service"
	"github.com/micro/services/location/handler"
	pb "github.com/micro/services/location/proto"
	admin "github.com/micro/services/pkg/service/proto"
)

func main() {
	location := service.New(
		service.Name("location"),
	)

	h := new(handler.Location)
	pb.RegisterLocationHandler(location.Server(), h)
	admin.RegisterAdminHandler(location.Server(), h)

	// TODO reinstate me
	//service.Subscribe(subscriber.Topic, new(subscriber.Location))
	if err := location.Run(); err != nil {
		log.Fatal(err)
	}
}
