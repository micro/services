package stream

import (
	"github.com/m3o/m3o-go/client"
)

func NewStreamService(token string) *StreamService {
	return &StreamService{
		client: client.NewClient(&client.Options{
			Token: token,
		}),
	}
}

type StreamService struct {
	client *client.Client
}

// List all the active channels
func (t *StreamService) ListChannels(request *ListChannelsRequest) (*ListChannelsResponse, error) {
	rsp := &ListChannelsResponse{}
	return rsp, t.client.Call("stream", "ListChannels", request, rsp)
}

// List messages for a given channel
func (t *StreamService) ListMessages(request *ListMessagesRequest) (*ListMessagesResponse, error) {
	rsp := &ListMessagesResponse{}
	return rsp, t.client.Call("stream", "ListMessages", request, rsp)
}

// Send a message to the stream.
func (t *StreamService) SendMessage(request *SendMessageRequest) (*SendMessageResponse, error) {
	rsp := &SendMessageResponse{}
	return rsp, t.client.Call("stream", "SendMessage", request, rsp)
}

type Channel struct {
	// last activity time
	LastActive string `json:"lastActive"`
	// name of the channel
	Name string `json:"name"`
}

type ListChannelsRequest struct {
}

type ListChannelsResponse struct {
	Channels []Channel `json:"channels"`
}

type ListMessagesRequest struct {
	// The channel to subscribe to
	Channel string `json:"channel"`
	// number of message to return
	Limit int32 `json:"limit"`
}

type ListMessagesResponse struct {
	// The channel subscribed to
	Channel string `json:"channel"`
	// Messages are chronological order
	Messages []Message `json:"messages"`
}

type Message struct {
	// the channel name
	Channel string `json:"channel"`
	// id of the message
	Id string `json:"id"`
	// the associated metadata
	Metadata map[string]string `json:"metadata"`
	// text of the message
	Text string `json:"text"`
	// time of message creation
	Timestamp string `json:"timestamp"`
}

type SendMessageRequest struct {
	// The channel to send to
	Channel string `json:"channel"`
	// The message text to send
	Text string `json:"text"`
}

type SendMessageResponse struct {
}
