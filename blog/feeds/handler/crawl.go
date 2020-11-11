package handler

import (
	"context"

	"net/url"

	"github.com/SlyMarbo/rss"
	log "github.com/micro/micro/v3/service/logger"
	feeds "github.com/micro/services/blog/feeds/proto"
	posts "github.com/micro/services/blog/posts/proto"
)

func (e *Feeds) fetchAll() {
	fs := []*feeds.Feed{}
	err := e.feeds.List(e.feedsNameIndex.ToQuery(nil), &fs)
	if err != nil {
		log.Errorf("Error listing feeds: %v", err)
		return
	}
	if len(fs) == 0 {
		log.Infof("No feeds to fetch")
		return
	}
	for _, feed := range fs {
		log.Infof("Fetching address %v", feed.Address)
		fd, err := rss.Fetch(feed.Address)
		if err != nil {
			log.Errorf("Error fetching address %v: %v", feed.Address, err)
			continue
		}
		domain := getDomain(feed.Address)
		for _, item := range fd.Items {
			err = e.entries.Save(feeds.Entry{
				Id:      item.ID,
				Url:     item.Link,
				Title:   item.Title,
				Domain:  domain,
				Content: item.Summary,
				Date:    item.Date.Unix(),
			})
			if err != nil {
				log.Errorf("Error saving item: %v", err)
			}
			// @todo make this optional
			_, err := e.postsService.Save(context.TODO(), &posts.SaveRequest{
				Id:        item.ID,
				Title:     item.Title,
				Content:   item.Content,
				Timestamp: item.Date.Unix(),
				Metadata: map[string]string{
					"domain": domain,
					"link":   item.Link,
				},
			})
			if err != nil {
				log.Errorf("Error saving post: %v", err)
			}
		}
	}
}

func getDomain(address string) string {
	uri, _ := url.Parse(address)
	return uri.Host
}
