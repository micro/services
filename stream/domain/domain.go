package domain

import (
	"errors"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/golang/groupcache/lru"
	"github.com/google/uuid"
)

const (
	defaultStream  = "_"
	maxMessageSize = 512
	maxMessages    = 1000
	maxStreams     = 1000
	streamTTL      = 8.64e13
)

type Metadata struct {
	Created     int64
	Title       string
	Description string
	Type        string
	Image       string
	Url         string
	Site        string
}

type Stream struct {
	Id          string
	Description string
	Messages    []*Message
	Updated     int64
}

type Message struct {
	Id       string
	Text     string
	Created  int64 `json:",string"`
	Stream   string
	Metadata *Metadata
}

type Store struct {
	Created int64
	Updates chan *Message

	mtx       sync.RWMutex
	Streams   *lru.Cache
	streams   map[string]*Stream
	metadatas map[string]*Metadata
}

var (
	C = newStore()
)

func newStore() *Store {
	return &Store{
		Created:   time.Now().UnixNano(),
		Streams:   lru.New(maxStreams),
		Updates:   make(chan *Message, 100),
		streams:   make(map[string]*Stream),
		metadatas: make(map[string]*Metadata),
	}
}

func newStream(id, desc string) *Stream {
	return &Stream{
		Id:          id,
		Description: desc,
		Updated:     time.Now().UnixNano(),
	}
}

func newMessage(text, stream string) *Message {
	return &Message{
		Id:      uuid.New().String(),
		Text:    text,
		Created: time.Now().UnixNano(),
		Stream:  stream,
	}
}

func getMetadata(uri string) *Metadata {
	u, err := url.Parse(uri)
	if err != nil {
		return nil
	}

	d, err := goquery.NewDocument(u.String())
	if err != nil {
		return nil
	}

	g := &Metadata{
		Created: time.Now().UnixNano(),
	}

	for _, node := range d.Find("meta").Nodes {
		if len(node.Attr) < 2 {
			continue
		}

		p := strings.Split(node.Attr[0].Val, ":")
		if len(p) < 2 || (p[0] != "twitter" && p[0] != "og") {
			continue
		}

		switch p[1] {
		case "site_name":
			g.Site = node.Attr[1].Val
		case "site":
			if len(g.Site) == 0 {
				g.Site = node.Attr[1].Val
			}
		case "title":
			g.Title = node.Attr[1].Val
		case "description":
			g.Description = node.Attr[1].Val
		case "card", "type":
			g.Type = node.Attr[1].Val
		case "url":
			g.Url = node.Attr[1].Val
		case "image":
			if len(p) > 2 && p[2] == "src" {
				g.Image = node.Attr[1].Val
			} else if len(g.Image) == 0 {
				g.Image = node.Attr[1].Val
			}
		}
	}

	if len(g.Type) == 0 || len(g.Image) == 0 || len(g.Title) == 0 || len(g.Url) == 0 {
		return nil
	}

	return g
}

func (c *Store) CreateStream(name, description string) {
	c.mtx.Lock()
	ch, ok := c.streams[name]
	if ok {
		ch.Description = description
	} else {
		ch = newStream(name, description)
	}
	c.streams[name] = ch
	c.mtx.Unlock()
}

func (c *Store) Metadata(t *Message) {
	parts := strings.Split(t.Text, " ")
	for _, part := range parts {
		g := getMetadata(part)
		if g == nil {
			continue
		}
		c.mtx.Lock()
		c.metadatas[t.Id] = g
		c.mtx.Unlock()
		return
	}
}

func (c *Store) ListStreams() map[string]*Stream {
	c.mtx.RLock()
	streams := c.streams
	c.mtx.RUnlock()
	return streams
}

func (c *Store) Save(message *Message) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	var stream *Stream

	if obj, ok := c.Streams.Get(message.Stream); ok {
		stream = obj.(*Stream)
	} else {
		stream = newStream(message.Stream, "")
		c.Streams.Add(message.Stream, stream)
	}

	stream.Messages = append(stream.Messages, message)
	if len(stream.Messages) > maxMessages {
		stream.Messages = stream.Messages[1:]
	}
	stream.Updated = time.Now().UnixNano()
}

func (c *Store) Retrieve(message string, streem string, direction, last, limit int64) []*Message {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	var stream *Stream

	if message, ok := c.Streams.Get(streem); ok {
		stream = message.(*Stream)
	} else {
		return []*Message{}
	}

	if len(message) == 0 {
		var messages []*Message

		if limit <= 0 {
			return messages
		}

		li := int(limit)

		// go back in time
		if direction < 0 {
			for i := len(stream.Messages) - 1; i >= 0; i-- {
				if len(messages) >= li {
					return messages
				}

				message := stream.Messages[i]

				if message.Created < last {
					if g, ok := c.metadatas[message.Id]; ok {
						tc := *message
						tc.Metadata = g
						messages = append(messages, &tc)
					} else {
						messages = append(messages, message)
					}
				}
			}
			return messages
		}

		start := 0
		if len(stream.Messages) > li {
			start = len(stream.Messages) - li
		}

		for i := start; i < len(stream.Messages); i++ {
			if len(messages) >= li {
				return messages
			}

			message := stream.Messages[i]

			if message.Created > last {
				if g, ok := c.metadatas[message.Id]; ok {
					tc := *message
					tc.Metadata = g
					messages = append(messages, &tc)
				} else {
					messages = append(messages, message)
				}
			}
		}
		return messages
	}

	// retrieve one
	for _, t := range stream.Messages {
		var messages []*Message
		if message == t.Id {
			if g, ok := c.metadatas[t.Id]; ok {
				tc := *t
				tc.Metadata = g
				messages = append(messages, &tc)
			} else {
				messages = append(messages, t)
			}
			return messages
		}
	}

	return []*Message{}
}

func (c *Store) Run() {
	t1 := time.NewTicker(time.Hour)
	t2 := time.NewTicker(time.Minute)
	streams := make(map[string]*Stream)

	for {
		select {
		case message := <-c.Updates:
			c.Save(message)
			ch, ok := streams[message.Stream]
			if !ok {
				ch = newStream(message.Stream, "")
				streams[message.Stream] = ch
			}
			ch.Updated = time.Now().UnixNano()
			streams[message.Stream] = ch
			go c.Metadata(message)
		case <-t1.C:
			now := time.Now().UnixNano()
			for stream, ch := range streams {
				if d := now - ch.Updated; d > streamTTL {
					c.Streams.Remove(stream)
					delete(streams, stream)
				}
			}
			c.mtx.Lock()
			for metadata, g := range c.metadatas {
				if d := now - g.Created; d > streamTTL {
					delete(c.metadatas, metadata)
				}
			}
			c.mtx.Unlock()
		case <-t2.C:
			c.mtx.Lock()
			c.streams = streams
			c.mtx.Unlock()
		}
	}
}

func CreateChannel(name, description string) {
	C.CreateStream(name, description)
}

func ListChannels() map[string]*Stream {
	return C.ListStreams()
}

func ListMessages(channel string, limit int64) []*Message {
	message := ""
	last := int64(0)
	direction := int64(1)

	if limit <= 0 {
		limit = 25
	}

	// default stream
	if len(channel) == 0 {
		channel = defaultStream
	}

	return C.Retrieve(message, channel, direction, last, limit)
}

func SendMessage(channel, message string) error {
	// default stream
	if len(channel) == 0 {
		channel = defaultStream
	}

	// default length
	if len(message) > maxMessageSize {
		message = message[:maxMessageSize]
	}

	select {
	case C.Updates <- newMessage(message, channel):
	case <-time.After(time.Second):
		return errors.New("timed out creating message")
	}

	return nil
}

// TODO: streams per user
func Setup() {
	go C.Run()
}
