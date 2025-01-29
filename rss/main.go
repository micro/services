package main

import (
	"time"

	"github.com/micro/micro/v5/service"
	"github.com/micro/micro/v5/service/logger"
	"github.com/micro/micro/v5/service/store"
	admin "github.com/micro/services/pkg/service/proto"

	"github.com/micro/services/rss/handler"
	pb "github.com/micro/services/rss/proto"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("rss"),
	)

	st := store.DefaultStore
	crawl := handler.NewCrawl(st)
	rss := handler.NewRss(st, crawl)

	// crawl
	go func() {
		crawl.FetchAll()
		tick := time.NewTicker(1 * time.Minute)
		for _ = range tick.C {
			crawl.FetchAll()
		}
	}()

	// Register handler
	pb.RegisterRssHandler(srv.Server(), rss)
	admin.RegisterAdminHandler(srv.Server(), rss)

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
