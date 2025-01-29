package main

import (
	"time"

	"github.com/micro/micro/v5/service"
	"github.com/micro/micro/v5/service/logger"
)

func main() {
	srv := service.New()

	go func() {
		for {
			logger.Infof("This is a log line %s", time.Now())
			time.Sleep(1 * time.Second)
		}
	}()

	srv.Run()
}
