package handler

import (
	"context"
	"time"

	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/model"
	pb "github.com/micro/services/rss/proto"
)

type Rss struct {
	feeds            model.Model
	entries          model.Model
	feedsIdIndex     model.Index
	feedsNameIndex   model.Index
	entriesDateIndex model.Index
	entriesURLIndex  model.Index
}

func NewRss() *Rss {
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

	f := &Rss{
		feeds: model.NewModel(
			model.WithKey("Name"),
			model.WithNamespace("feeds"),
			model.WithIndexes(nameIndex),
		),
		entries: model.NewModel(
			model.WithNamespace("entries"),
			model.WithIndexes(dateIndex, entriesURLIndex),
		),
		feedsIdIndex:     idIndex,
		feedsNameIndex:   nameIndex,
		entriesDateIndex: dateIndex,
		entriesURLIndex:  entriesURLIndex,
	}

	// register model instances
	f.feeds.Register(new(pb.Feed))
	f.entries.Register(new(pb.Entry))

	go f.crawl()
	return f
}

func (e *Rss) crawl() {
	e.fetchAll()
	tick := time.NewTicker(1 * time.Minute)
	for _ = range tick.C {
		e.fetchAll()
	}
}

func (e *Rss) Add(ctx context.Context, req *pb.AddRequest, rsp *pb.AddResponse) error {
	log.Info("Received Rss.Add request")

	if len(req.Name) == 0 {
		return errors.BadRequest("rss.add", "require name")
	}

	rssSync.RLock()
	defer rssSync.RUnlock()

	// check if the feed already exists
	if _, ok := rssRss[req.Name]; ok {
		return errors.BadRequest("rss.add", "%s already exists", req.Name)
	}

	f := pb.Feed{
		Name:     req.Name,
		Url:      req.Url,
		Category: req.Category,
	}

	// create the feed
	e.feeds.Create(f)

	// schedule immediate fetch
	go e.fetch(&f)

	return nil
}

func (e *Rss) Feed(ctx context.Context, req *pb.FeedRequest, rsp *pb.FeedResponse) error {
	log.Info("Received Rss.Entries request")
	if len(req.Name) == 0 {
		return errors.BadRequest("rss.feed", "missing feed name")
	}

	feed := new(pb.Feed)
        q := model.QueryEquals("Name", req.Name)

	// get the feed
        if err := e.feeds.Read(q, feed); err != nil {
		return errors.InternalServerError("rss.feeds", "could not read feed")
	}

	// get the entries for each
	return e.entries.Read(e.entriesURLIndex.ToQuery(feed.Url), &rsp.Entries)
}

func (e *Rss) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	var feeds []*pb.Feed
	q := model.QueryAll()
	q.Index.FieldName = "Name"
	err := e.feeds.Read(q, &feeds)
	if err != nil {
		return errors.InternalServerError("rss.list", "failed to read list of feeds: %v", err)
	}

	rsp.Feeds = feeds
	return nil
}

func (e *Rss) Remove(ctx context.Context, req *pb.RemoveRequest, rsp *pb.RemoveResponse) error {
	if len(req.Name) == 0 {
		return errors.BadRequest("rss.remove", "blank name provided")
	}

	e.feeds.Delete(model.QueryEquals("name", req.Name))
	return nil
}
