package handler

import (
	"fmt"
	"time"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/events"

	"github.com/nats-io/nats-streaming-server/util"
	"gorm.io/gorm"
)

var (
	TokenTTL          = time.Minute
	ErrMissingTopic   = errors.BadRequest("MISSING_TOPIC", "Missing topic")
	ErrInvalidTopic   = errors.BadRequest("MISSING_TOPIC", "Invalid topic")
	ErrMissingToken   = errors.BadRequest("MISSING_TOKEN", "Missing token")
	ErrMissingMessage = errors.BadRequest("MISSING_MESSAGE", "Missing message")
	ErrInvalidToken   = errors.Forbidden("INVALID_TOKEN", "Invalid token")
	ErrExpiredToken   = errors.Forbidden("EXPIRED_TOKEN", "Token expired")
	ErrForbiddenTopic = errors.Forbidden("FORBIDDEN_TOPIC", "Token has not have permission to subscribe to this topic")
)

type Token struct {
	Token     string `gorm:"primaryKey"`
	Topic     string
	ExpiresAt time.Time
	Namespace string
}

type Streams struct {
	DB     *gorm.DB
	Events events.Stream
	Time   func() time.Time
}

// fmtTopic returns a topic string with namespace prefix
func fmtTopic(ns, topic string) string {
	return fmt.Sprintf("%s.%s", ns, topic)
}

// validateTopicInput validates that topic is alphanumeric
func validateTopicInput(topic string) error {
	if !util.IsChannelNameValid(topic, false) {
		return ErrInvalidTopic
	}
	return nil
}
