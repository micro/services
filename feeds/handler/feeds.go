package handler

import (
	"context"
	"time"

	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/model"

	feeds "github.com/micro/services/feeds/proto"
	posts "github.com/micro/services/posts/proto"
)

type Feeds struct {
	feeds            model.Model
	entries          model.Model
	postsService     posts.PostsService
	feedsIdIndex     model.Index
	feedsNameIndex   model.Index
	entriesDateIndex model.Index
	entriesURLIndex  model.Index
}

func NewFeeds(postsService posts.PostsService) *Feeds {
	idIndex := model.ByEquality("url")
	idIndex.Order.Type = model.OrderTypeUnordered

	nameIndex := model.ByEquality("name")
	nameIndex.Unique = true
	nameIndex.Order.Type = model.OrderTypeUnordered

	dateIndex := model.ByEquality("date")
	dateIndex.Order.Type = model.OrderTypeDesc

	entriesURLIndex := model.ByEquality("url")
	entriesURLIndex.Order.Type = model.OrderTypeDesc
	entriesURLIndex.Order.FieldName = "date"

	f := &Feeds{
		feeds: model.NewModel(
			model.WithKey("Name"),
			model.WithNamespace("feeds"),
			model.WithIndexes(nameIndex),
		),
		entries: model.NewModel(
			model.WithNamespace("entries"),
			model.WithIndexes(dateIndex, entriesURLIndex),
		),
		postsService:     postsService,
		feedsIdIndex:     idIndex,
		feedsNameIndex:   nameIndex,
		entriesDateIndex: dateIndex,
		entriesURLIndex:  entriesURLIndex,
	}

	// register model instances
	f.feeds.Register(new(feeds.Feed))
	f.entries.Register(new(feeds.Entry))

	go f.crawl()
	return f
}

func (e *Feeds) crawl() {
	e.fetchAll()
	tick := time.NewTicker(1 * time.Minute)
	for _ = range tick.C {
		e.fetchAll()
	}
}

func (e *Feeds) Add(ctx context.Context, req *feeds.AddRequest, rsp *feeds.AddResponse) error {
	log.Info("Received Feeds.New request")
	e.feeds.Create(feeds.Feed{
		Name: req.Name,
		Url:  req.Url,
	})
	return nil
}

func (e *Feeds) Entries(ctx context.Context, req *feeds.EntriesRequest, rsp *feeds.EntriesResponse) error {
	log.Info("Received Feeds.New request")
	err := e.fetch(req.Url)
	if err != nil {
		return err
	}
	return e.entries.Read(e.entriesURLIndex.ToQuery(req.Url), &rsp.Entries)
}

func (e *Feeds) List(ctx context.Context, req *feeds.ListRequest, rsp *feeds.ListResponse) error {
	var feeds []*feeds.Feed

	err := e.feeds.Read(model.QueryAll(), &feeds)
	if err != nil {
		return errors.InternalServerError("feeds.list", "failed to read list of feeds: %v", err)
	}

	rsp.Feeds = feeds
	return nil
}

func (e *Feeds) Remove(ctx context.Context, req *feeds.RemoveRequest, rsp *feeds.RemoveResponse) error {
	if len(req.Name) == 0 {
		return errors.BadRequest("feeds.remove", "blank name provided")
	}

	e.feeds.Delete(model.QueryEquals("name", req.Name))
	return nil
}
