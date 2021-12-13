package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"

	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"

	"github.com/micro/services/pkg/tenant"
	pb "github.com/micro/services/rss/proto"
)

type Rss struct {
	store store.Store
	crawl Crawler
}

// feedIdFromName generates md5 id by feed's name
func feedIdFromName(name string) string {
	hash := fnv.New64a()
	hash.Write([]byte(name))
	return fmt.Sprintf("%d", hash.Sum64())
}

// generateFeedKey returns feed key in store
func generateFeedKey(ctx context.Context, name string) string {
	tenantID, ok := tenant.FromContext(ctx)
	if !ok {
		tenantID = "micro"
	}

	return fmt.Sprintf("rss/feed/%s/%s", tenantID, feedIdFromName(name))
}

func NewRss(st store.Store, cr Crawler) *Rss {
	f := &Rss{
		store: st,
		crawl: cr,
	}

	return f
}

func (e *Rss) Add(ctx context.Context, req *pb.AddRequest, rsp *pb.AddResponse) error {
	log.Info("Received Rss.Add request")

	if len(req.Name) == 0 {
		return errors.BadRequest("rss.add", "require name")
	}

	key := generateFeedKey(ctx, req.Name)
	f := &pb.Feed{
		Id:       key,
		Name:     req.Name,
		Url:      req.Url,
		Category: req.Category,
	}

	// create the feed
	val, err := json.Marshal(f)
	if err != nil {
		return err
	}

	if err := e.store.Write(&store.Record{Key: key, Value: val}); err != nil {
		return err
	}

	// schedule immediate fetch
	go func() {
		_ = e.crawl.Fetch(f)
	}()

	return nil
}

func (e *Rss) Feed(ctx context.Context, req *pb.FeedRequest, rsp *pb.FeedResponse) error {
	log.Info("Received Rss.Entries request")

	prefix := generateFeedKey(ctx, req.Name)
	var records []*store.Record
	var err error

	// get records with prefix
	if len(req.Name) > 0 {
		records, err = e.store.Read(prefix)
	} else {
		records, err = e.store.Read(prefix, store.ReadPrefix())
	}

	if err != nil {
		return err
	}

	if len(records) == 0 {
		return nil
	}

	if req.Limit == 0 {
		req.Limit = int64(25)
	}

	for _, v := range records {
		// decode feed
		feed := pb.Feed{}
		if err := json.Unmarshal(v.Value, &feed); err != nil {
			log.Errorf("json unmarshal feed error: %v", err)
			continue
		}

		// read entries with prefix
		entryPrefix := generateEntryKey(feed.Url, "")
		entries, err := e.store.Read(entryPrefix, store.ReadPrefix())
		if err != nil {
			log.Errorf("read feed entry from store error: %v", err)
			continue
		}

		for _, val := range entries {
			var entry pb.Entry
			if err := json.Unmarshal(val.Value, &entry); err != nil {
				log.Errorf("json unmarshal entry error: %v", err)
				continue
			}

			rsp.Entries = append(rsp.Entries, &entry)
		}

		if len(rsp.Entries) >= int(req.Limit) {
			break
		}
	}

	return nil
}

func (e *Rss) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	prefix := generateFeedKey(ctx, "")
	records, err := e.store.Read(prefix, store.ReadPrefix())

	if err != nil {
		return err
	}

	for _, val := range records {
		var feed = pb.Feed{}
		if err := json.Unmarshal(val.Value, &feed); err != nil {
			continue
		}
		rsp.Feeds = append(rsp.Feeds, &feed)
	}

	return nil
}

func (e *Rss) Remove(ctx context.Context, req *pb.RemoveRequest, rsp *pb.RemoveResponse) error {
	if len(req.Name) == 0 {
		return errors.BadRequest("rss.remove", "name is required")
	}

	return e.store.Delete(generateFeedKey(ctx, req.Name))
}
