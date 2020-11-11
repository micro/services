package handler

import (
	"context"
	"time"

	"github.com/micro/dev/model"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"

	feeds "github.com/micro/services/blog/feeds/proto"
	posts "github.com/micro/services/blog/posts/proto"
)

type Feeds struct {
	feeds            model.Model
	entries          model.Model
	postsService     posts.PostsService
	feedsIdIndex     model.Index
	feedsNameIndex   model.Index
	entriesDateIndex model.Index
}

func NewFeeds(postsService posts.PostsService) *Feeds {
	idIndex := model.ByEquality("address")
	idIndex.Order.Type = model.OrderTypeUnordered

	nameIndex := model.ByEquality("name")
	nameIndex.Unique = true
	nameIndex.Order.Type = model.OrderTypeUnordered

	dateIndex := model.ByEquality("date")
	dateIndex.Order.Type = model.OrderTypeDesc

	f := &Feeds{
		feeds: model.New(
			store.DefaultStore,
			"feeds",
			model.Indexes(nameIndex),
			&model.ModelOptions{
				Debug:   false,
				IdIndex: idIndex,
			},
		),
		entries: model.New(
			store.DefaultStore,
			"entries",
			model.Indexes(dateIndex),
			&model.ModelOptions{
				Debug: false,
			},
		),
		postsService:     postsService,
		feedsIdIndex:     idIndex,
		feedsNameIndex:   nameIndex,
		entriesDateIndex: dateIndex,
	}

	go f.crawl()
	return f
}

func (e *Feeds) crawl() {
	e.fetchAll()
	tick := time.NewTicker(10 * time.Minute)
	for _ = range tick.C {
		e.fetchAll()
	}
}

func (e *Feeds) New(ctx context.Context, req *feeds.NewRequest, rsp *feeds.NewResponse) error {
	log.Info("Received Feeds.New request")
	e.feeds.Save(feeds.Feed{
		Name:    req.Name,
		Address: req.Address,
	})
	return nil
}
