package main

import (
	"time"

	"micro.dev/v4/service"
	"micro.dev/v4/service/logger"
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
