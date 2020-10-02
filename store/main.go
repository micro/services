package main

import (
	"github.com/micro/services/store/handler"

	"github.com/micro/micro/v3/service"
)

func main() {
	srv := service.New(service.Name("example"))

	srv.Handle(new(handler.Example))

	srv.Run()
}
