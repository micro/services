package handler

import (
	"fmt"
	"time"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/events"
	gorm2 "github.com/micro/services/pkg/gorm"

	"github.com/nats-io/nats-streaming-server/util"
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
}

type Streams struct {
	gorm2.Helper
	Events events.Stream
	Time   func() time.Time
}

// fmtTopic returns a topic string with namespace prefix
func fmtTopic(acc *auth.Account, topic string) string {
	owner := acc.Metadata["apikey_owner"]
	if len(owner) == 0 {
		owner = acc.ID
	}
	return fmt.Sprintf("%s.%s.%s", acc.Issuer, owner, topic)
}

// validateTopicInput validates that topic is alphanumeric
func validateTopicInput(topic string) error {
	if !util.IsChannelNameValid(topic, false) {
		return ErrInvalidTopic
	}
	return nil
}
