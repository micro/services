package main

import (
	"github.com/micro/services/test/kv/handler"

	"micro.dev/v4/service"
)

func main() {
	srv := service.New(service.Name("example"))

	srv.Handle(new(handler.Example))

	srv.Run()
}
