package handler

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/SlyMarbo/rss"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/rss/parser"
	pb "github.com/micro/services/rss/proto"
)

var (
	rssSync  sync.RWMutex
	rssFeeds = map[string]*rss.Feed{}
)

func (e *Rss) fetchAll() {
	fs := []*pb.Feed{}
	err := e.feeds.Read(e.feedsNameIndex.ToQuery(nil), &fs)
	if err != nil {
		log.Errorf("Error listing pb: %v", err)
		return
	}
	if len(fs) == 0 {
		log.Infof("No pb to fetch")
		return
	}
	for _, feed := range fs {
		err = e.fetch(feed)
		if err != nil {
			log.Errorf("Error saving post: %v", err)
		}
	}
}

func (e *Rss) fetch(f *pb.Feed) error {
	url := f.Url
	log.Infof("Fetching address %v", url)

	// see if there's an existing rss feed
	rssSync.RLock()
	fd, ok := rssFeeds[f.Url]
	rssSync.RUnlock()

	if !ok {
		// create a new one if it doesn't exist
		var err error
		fd, err = rss.Fetch(f.Url)
		if err != nil {
			return fmt.Errorf("Error fetching address %v: %v", url, err)
		}
		// save the feed
		rssSync.Lock()
		rssFeeds[f.Url] = fd
		rssSync.Unlock()
	} else {
		// otherwise update the existing feed
		fd.Items = []*rss.Item{}
		fd.Unread = 0
		if err := fd.Update(); err != nil {
			return fmt.Errorf("Error updating address %v: %v", url, err)
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
		// e.g cnbc
		c, err := parser.Parse(item.Link)
		if err == nil && len(c) > 0 {
			content = c
		}

		err = e.entries.Create(pb.Entry{
			Id:       id,
			Title:    item.Title,
			Summary:  item.Summary,
			Feed:      f.Url,
			Link: item.Link,
			Content:  content,
			Date:     item.Date.Format(time.RFC3339Nano),
		})
		if err != nil {
			return fmt.Errorf("Error saving item: %v", err)
		}

	}

	return nil
}

func getDomain(address string) string {
	uri, _ := url.Parse(address)
	return uri.Host
}
