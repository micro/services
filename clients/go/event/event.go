package event

import (
	"github.com/m3o/m3o-go/client"
)

func NewEventService(token string) *EventService {
	return &EventService{
		client: client.NewClient(&client.Options{
			Token: token,
		}),
	}
}

type EventService struct {
	client *client.Client
}

// Publish a message to the event stream.
func (t *EventService) Publish(request *PublishRequest) (*PublishResponse, error) {
	rsp := &PublishResponse{}
	return rsp, t.client.Call("event", "Publish", request, rsp)
}

// Subscribe to messages for a given topic.
func (t *EventService) Subscribe(request *SubscribeRequest) (*SubscribeResponse, error) {
	rsp := &SubscribeResponse{}
	return rsp, t.client.Call("event", "Subscribe", request, rsp)
}

type PublishRequest struct {
	// The json message to publish
	Message map[string]interface{} `json:"message"`
	// The topic to publish to
	Topic string `json:"topic"`
}

type PublishResponse struct {
}

type SubscribeRequest struct {
	// Optional group for the subscription
	Group string `json:"group"`
	// The topic to subscribe to
	Topic string `json:"topic"`
}

type SubscribeResponse struct {
	// The next json message on the topic
	Message map[string]interface{} `json:"message"`
	// The topic subscribed to
	Topic string `json:"topic"`
}
