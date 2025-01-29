package handler

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/SlyMarbo/rss"
	log "github.com/micro/micro/v5/service/logger"
	"github.com/micro/micro/v5/service/store"

	"github.com/micro/services/rss/parser"
	pb "github.com/micro/services/rss/proto"
)

var (
	rssSync  sync.RWMutex
	rssFeeds = map[string]*rss.Feed{}
)

type Crawler interface {
	Fetch(f *pb.Feed) error
	FetchAll()
}

type crawl struct {
	store store.Store
}

func NewCrawl(st store.Store) *crawl {
	return &crawl{store: st}
}

func generateEntryKey(feedUrl, id string) string {
	return fmt.Sprintf("rss/entry/%s/%s", feedUrl, id)
}

func (e *crawl) FetchAll() {
	prefix := "rss/feed/"
	records, err := e.store.Read(prefix, store.ReadPrefix())

	if err != nil {
		log.Errorf("get feeds list error: %v", err)
		return
	}

	if len(records) == 0 {
		log.Infof("No pb to fetch")
		return
	}

	currList := map[string]bool{}
	for _, v := range records {
		feed := pb.Feed{}
		if err := json.Unmarshal(v.Value, &feed); err != nil {
			log.Errorf("crawl.fetchAll json unmarshal feed error: %v", err)
			continue
		}

		err = e.Fetch(&feed)
		if err != nil {
			log.Errorf("Error saving post: %v", err)
		}
		currList[feed.Url] = true
	}

	// prune anything that has been deleted
	rssSync.Lock()
	defer rssSync.Unlock()
	for url, _ := range rssFeeds {
		if currList[url] {
			continue
		}
		// this isn't in the current list. delete from store any entries
		keys, _ := store.List(store.ListPrefix(generateEntryKey(url, "")))
		for _, k := range keys {
			store.Delete(k)
		}
		delete(rssFeeds, url)
	}
}

func (e *crawl) Fetch(f *pb.Feed) error {
	log.Infof("Fetching address %v", f.Url)

	// see if there's an existing rss feed
	rssSync.RLock()
	fd, ok := rssFeeds[f.Url]
	rssSync.RUnlock()

	if !ok {
		// create a new one if it doesn't exist
		var err error
		fd, err = rss.Fetch(f.Url)
		if err != nil {
			return fmt.Errorf("error fetching address %v: %v", f.Url, err)
		}
		// save the feed
		rssSync.Lock()
		rssFeeds[f.Url] = fd
		rssSync.Unlock()
	} else {
		// otherwise, update the existing feed
		fd.Items = []*rss.Item{}
		fd.Unread = 0
		if err := fd.Update(); err != nil {
			return fmt.Errorf("error updating address %v: %v", f.Url, err)
		}
	}

	// set the refresh time
	fd.Refresh = time.Now()

	// range over the feed and save the items
	for _, item := range fd.Items {
		id := fmt.Sprintf("%x", md5.Sum([]byte(item.ID)))

		// check if content exists
		content := item.Content

		// if we have a parser which returns content use it
		// e.g. cnbc
		c, err := parser.Parse(item.Link)
		if err == nil && len(c) > 0 {
			content = c
		}

		val, err := json.Marshal(&pb.Entry{
			Id:      id,
			Title:   item.Title,
			Summary: item.Summary,
			Feed:    f.Url,
			Link:    item.Link,
			Content: content,
			Date:    item.Date.Format(time.RFC3339Nano),
		})

		if err != nil {
			log.Errorf("json marshal entry error: %v", err)
			continue
		}

		// save
		err = e.store.Write(&store.Record{
			Key:   generateEntryKey(f.Url, id),
			Value: val,
		})
		if err != nil {
			return fmt.Errorf("error saving item: %v", err)
		}

	}

	return nil
}
