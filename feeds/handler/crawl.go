package handler

import (
	"context"
	"crypto/md5"
	"fmt"
	"net/url"
	"sync"

	"github.com/SlyMarbo/rss"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/feeds/parser"
	feeds "github.com/micro/services/feeds/proto"
	posts "github.com/micro/services/posts/proto"
)

var (
	rssSync  sync.RWMutex
	rssFeeds = map[string]*rss.Feed{}
)

func (e *Feeds) fetchAll() {
	fs := []*feeds.Feed{}
	err := e.feeds.Read(e.feedsNameIndex.ToQuery(nil), &fs)
	if err != nil {
		log.Errorf("Error listing feeds: %v", err)
		return
	}
	if len(fs) == 0 {
		log.Infof("No feeds to fetch")
		return
	}
	for _, feed := range fs {
		err = e.fetch(feed)
		if err != nil {
			log.Errorf("Error saving post: %v", err)
		}
	}
}

func (e *Feeds) fetch(f *feeds.Feed) error {
	url := f.Url
	log.Infof("Fetching address %v", url)

	// see if there's an existing rss feed
	rssSync.RLock()
	fd, ok := rssFeeds[f.Name]
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
		rssFeeds[f.Name] = fd
		rssSync.Unlock()
	} else {
		// otherwise update the existing feed
		fd.Items = []*rss.Item{}
		fd.Unread = 0
		if err := fd.Update(); err != nil {
			return fmt.Errorf("Error updating address %v: %v", url, err)
		}
	}

	domain := getDomain(url)

	// range over the feed and save the items
	for _, item := range fd.Items {
		id := fmt.Sprintf("%x", md5.Sum([]byte(item.ID)))

		err := e.entries.Create(feeds.Entry{
			Id:       id,
			Url:      item.Link,
			Title:    item.Title,
			Domain:   domain,
			Content:  item.Summary,
			Date:     item.Date.Unix(),
			Category: f.Category,
		})
		if err != nil {
			return fmt.Errorf("Error saving item: %v", err)
		}

		var tags []string
		if len(f.Category) > 0 {
			tags = append(tags, f.Category)
		}

		// check if content exists
		content := item.Content
		if len(content) == 0 && len(item.Summary) > 0 {
			content = item.Summary
		}

		// if we have a parser which returns content use it
		// e.g cnbc
		c, err := parser.Parse(item.Link)
		if err == nil && len(c) > 0 {
			content = c
		}

		// @todo make this optional
		_, err = e.postsService.Save(context.TODO(), &posts.SaveRequest{
			Id:        id,
			Title:     item.Title,
			Content:   content,
			Timestamp: item.Date.Unix(),
			Metadata: map[string]string{
				"domain": domain,
				"link":   item.Link,
			},
			Tags: tags,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func getDomain(address string) string {
	uri, _ := url.Parse(address)
	return uri.Host
}
