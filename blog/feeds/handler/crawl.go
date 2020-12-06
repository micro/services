package handler

import (
	"context"
	"crypto/md5"
	"fmt"

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
		err = e.fetch(feed.Url)
		if err != nil {
			log.Errorf("Error saving post: %v", err)
		}
	}
}

func (e *Feeds) fetch(url string) error {
	log.Infof("Fetching address %v", url)
	fd, err := rss.Fetch(url)
	if err != nil {
		return fmt.Errorf("Error fetching address %v: %v", url, err)
	}
	domain := getDomain(url)

	for _, item := range fd.Items {
		id := fmt.Sprintf("%x", md5.Sum([]byte(item.ID)))
		err = e.entries.Save(feeds.Entry{
			Id:      id,
			Url:     item.Link,
			Title:   item.Title,
			Domain:  domain,
			Content: item.Summary,
			Date:    item.Date.Unix(),
		})
		if err != nil {
			return fmt.Errorf("Error saving item: %v", err)
		}
		// @todo make this optional
		_, err := e.postsService.Save(context.TODO(), &posts.SaveRequest{
			Id:        id,
			Title:     item.Title,
			Content:   item.Content,
			Timestamp: item.Date.Unix(),
			Metadata: map[string]string{
				"domain": domain,
				"link":   item.Link,
			},
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
