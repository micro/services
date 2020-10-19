package handler

import (
	"context"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/context/metadata"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"

	// it's standard to import the services own proto under the alias pb
	pb "github.com/micro/services/chat/proto"
)

const (
	chatStoreKeyPrefix    = "chats/"
	chatEventKeyPrefix    = "chats/"
	messageStoreKeyPrefix = "messages/"
)

// Chat satisfies the ChatHandler interface. You can see this inteface defined in chat.pb.micro.go
type Chat struct{}

// New creates a chat for a group of users. The RPC is idempotent so if it's called multiple times
// for the same users, the same response will be returned. It's good practice to design APIs as
// idempotent since this enables safe retries.
func (c *Chat) New(ctx context.Context, req *pb.NewRequest, rsp *pb.NewResponse) error {
	// in a real world application we would authorize the request to ensure the authenticated user
	// is part of the chat they're attempting to create. We could do this by getting the user id from
	// auth.AccountFromContext(ctx) and then validating the presence of their id in req.UserIds. If
	// the user is not part of the request then we'd return a Forbidden error, which the micro api
	// would transform to a 403 status code.

	// validate the request
	if len(req.UserIds) == 0 {
		// Return a bad request error to the client, the first argument is a unique id which the client
		// can check for. The second argument is a human readable description. Returning the correct type
		// of error is important as it's used by the network to know if a request should be retried. Only
		// 500 (InternalServerError) and 408 (Timeout) errors are retried.
		return errors.BadRequest("chat.New.MissingUserIDs", "One or more user IDs are required")
	}

	// construct a key to identify the chat, we'll do this by sorting the user ids alphabetically and
	// then joining them. When a service calls the store, the data returned will be automatically scoped
	// to the service however it's still advised to use a prefix when writing data since this allows
	// other types of keys to be written in the future. We'll make a copy of the req.UserIds object as
	// it's a good practice to not mutate the request object.
	sortedIDs := make([]string, len(req.UserIds))
	copy(sortedIDs, req.UserIds)
	sort.Strings(sortedIDs)

	// key to lookup the chat in the store using, e.g. "chat/usera-userb-userc"
	key := chatStoreKeyPrefix + strings.Join(sortedIDs, "-")

	// read from the store to check if a chat with these users already exists
	recs, err := store.Read(key)
	if err == nil {
		// if an error wasn't returned, at least one record was found. The value returned by the store
		// is the bytes representation of the chat id. We'll convert this back into a string and return
		// it to the client.
		rsp.ChatId = string(recs[0].Value)
		return nil
	} else if err != store.ErrNotFound {
		// if no records were found then we'd expect to get a store.ErrNotFound error returned. If this
		// wasn't the case, the service could've experienced an issue connecting to the store so we should
		// log the error and return an InternalServerError to the client, indicating the request should
		// be retried
		logger.Errorf("Error reading from the store. Key: %v. Error: %v", key, err)
		return errors.InternalServerError("chat.New.Unknown", "Error reading from the store")
	}

	// no chat id was returned so we'll generate one, write it to the store and then return it to the
	// client
	chatID := uuid.New().String()
	record := store.Record{Key: chatStoreKeyPrefix + chatID, Value: []byte(chatID)}
	if err := store.Write(&record); err != nil {
		logger.Errorf("Error writing to the store. Key: %v. Error: %v", record.Key, err)
		return errors.InternalServerError("chat.New.Unknown", "Error writing to the store")
	}

	// The chat was successfully created so we'll log the event and then return the id to the client.
	// Note that we'll use logger.Infof here vs the Errorf above.
	logger.Infof("New chat created with ID %v", chatID)
	rsp.ChatId = chatID
	return nil
}

// History returns the historical messages in a chat
func (c *Chat) History(ctx context.Context, req *pb.HistoryRequest, rsp *pb.HistoryResponse) error {
	// as per the New function, in a real world application we would authorize the request to ensure
	// the authenticated user is part of the chat they're attempting to read the history of

	// validate the request
	if len(req.ChatId) == 0 {
		return errors.BadRequest("chat.History.MissingChatID", "ChatID is missing")
	}

	// lookup the chat from the store to ensure it's valid
	if _, err := store.Read(chatStoreKeyPrefix + req.ChatId); err == store.ErrNotFound {
		return errors.BadRequest("chat.History.InvalidChatID", "Chat not found with this ID")
	} else if err != nil {
		logger.Errorf("Error reading from the store. Chat ID: %v. Error: %v", req.ChatId, err)
		return errors.InternalServerError("chat.History.Unknown", "Error reading from the store")
	}

	// lookup the historical messages for the chat using the event store. lots of packages in micro
	// support options, in this case we'll pass the ReadLimit option to restrict the number of messages
	// we'll load from the events store.
	messages, err := events.Read(chatEventKeyPrefix+req.ChatId, events.ReadLimit(50))
	if err != nil {
		logger.Errorf("Error reading from the event store. Chat ID: %v. Error: %v", req.ChatId, err)
		return errors.InternalServerError("chat.History.Unknown", "Error reading from the event store")
	}

	// we've loaded the messages from the event store. next we need to serialize them and return them
	// to the client. The message is stored in the event payload, to retrieve it we need to unmarshal
	// the event into a message struct.
	rsp.Messages = make([]*pb.Message, len(messages))
	for i, ev := range messages {
		var msg pb.Message
		if err := ev.Unmarshal(&msg); err != nil {
			logger.Errorf("Error unmarshaling event: %v", err)
			return errors.InternalServerError("chat.History.Unknown", "Error unmarshaling event")
		}
		rsp.Messages[i] = &msg
	}

	return nil
}

// Send a single message to the chat, designed for ease of use via the API / CLI
func (c *Chat) Send(ctx context.Context, req *pb.SendRequest, rsp *pb.SendResponse) error {
	// validate the request
	if len(req.ChatId) == 0 {
		return errors.BadRequest("chat.Send.MissingChatID", "ChatID is missing")
	}
	if len(req.UserId) == 0 {
		return errors.BadRequest("chat.Send.MissingUserID", "UserID is missing")
	}
	if len(req.Text) == 0 {
		return errors.BadRequest("chat.Send.MissingText", "Text is missing")
	}

	// construct the message
	msg := &pb.Message{
		Id:       uuid.New().String(),
		ClientId: req.ClientId,
		ChatId:   req.ChatId,
		UserId:   req.UserId,
		Subject:  req.Subject,
		Text:     req.Text,
	}

	// default the client id if not provided
	if len(msg.ClientId) == 0 {
		msg.ClientId = uuid.New().String()
	}

	// create the message
	return c.createMessage(msg)
}

// Connect to a chat using a bidirectional stream enabling the client to send and recieve messages
// over a single RPC. When a message is sent on the stream, it will be added to the chat history
// and sent to the other connected users. When opening the connection, the client should provide
// the chat_id and user_id in the context so the server knows which messages to stream.
func (c *Chat) Connect(ctx context.Context, stream pb.Chat_ConnectStream) error {
	// the client passed the chat id and user id in the request context. we'll load that information
	// now and validate it. If any information is missing we'll return a BadRequest error to the client
	userID, ok := metadata.Get(ctx, "user-id")
	if !ok {
		return errors.BadRequest("chat.Connect.MissingUserID", "UserID missing in context")
	}
	chatID, ok := metadata.Get(ctx, "chat-id")
	if !ok {
		return errors.BadRequest("chat.Connect.MissingChatID", "ChatId missing in context")
	}

	// lookup the chat from the store to ensure it's valid
	if _, err := store.Read(chatStoreKeyPrefix + chatID); err == store.ErrNotFound {
		return errors.BadRequest("chat.Connect.InvalidChatID", "Chat not found with this ID")
	} else if err != nil {
		logger.Errorf("Error reading from the store. Chat ID: %v. Error: %v", chatID, err)
		return errors.InternalServerError("chat.Connect.Unknown", "Error reading from the store")
	}

	// as per the New and Connect functions, at this point in a real world application we would
	// authorize the request to ensure the authenticated user is part of the chat they're attempting
	// to read the history of

	// create a new context which can be cancelled, in the case either the consumer of publisher errors
	// we don't want one to keep running in the background
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// create a channel to send errors on, because the subscriber / publisher will run in seperate go-
	// routines, they need a way of returning errors to the client
	errChan := make(chan error)

	// create an event stream to consume messages posted by other users into the chat. we'll use the
	// user id as a queue to ensure each user recieves the message
	evStream, err := events.Consume(chatEventKeyPrefix+chatID, events.WithGroup(userID))
	if err != nil {
		logger.Errorf("Error streaming events. Chat ID: %v. Error: %v", chatID, err)
		return errors.InternalServerError("chat.Connect.Unknown", "Error connecting to the event stream")
	}
	go func() {
		for {
			select {
			case <-cancelCtx.Done():
				// the context has been cancelled or timed out, stop subscribing to new messages
				return
			case ev := <-evStream:
				// recieved a message, unmarshal it into a message struct. if an error occurs log it and
				// cancel the context
				var msg pb.Message
				if err := ev.Unmarshal(&msg); err != nil {
					logger.Errorf("Error unmarshaling message. ChatID: %v. Error: %v", chatID, err)
					errChan <- err
					return
				}

				// ignore any messages published by the current user
				if msg.UserId == userID {
					continue
				}

				// publish the message to the stream
				if err := stream.Send(&msg); err != nil {
					logger.Errorf("Error sending message to stream. ChatID: %v. Message ID: %v. Error: %v", chatID, msg.Id, err)
					errChan <- err
					return
				}
			}
		}
	}()

	// transform the stream.Recv into a channel which can be used in the select statement below
	msgChan := make(chan *pb.Message)
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				errChan <- err
				close(msgChan)
				return
			}
			msgChan <- msg
		}
	}()

	for {
		select {
		case <-cancelCtx.Done():
			// the context has been cancelled or timed out, stop subscribing to new messages
			return nil
		case err := <-errChan:
			// an error occured in another goroutine, terminate the stream
			return err
		case msg := <-msgChan:
			// set the defaults
			msg.UserId = userID
			msg.ChatId = chatID

			// create the message
			if err := c.createMessage(msg); err != nil {
				return err
			}
		}
	}
}

// createMessage is a helper function which creates a message in the event stream. It handles the
// logic for ensuring client id is unique.
func (c *Chat) createMessage(msg *pb.Message) error {
	// a message was recieved from the client. validate it hasn't been recieved before
	if _, err := store.Read(messageStoreKeyPrefix + msg.ClientId); err == nil {
		// the message has already been processed
		return nil
	} else if err != store.ErrNotFound {
		// an unexpected error occured
		return err
	}

	// send the message to the event stream
	if err := events.Publish(chatEventKeyPrefix+msg.ChatId, msg); err != nil {
		return err
	}

	// record the messages client id
	if err := store.Write(&store.Record{Key: messageStoreKeyPrefix + msg.ClientId}); err != nil {
		return err
	}

	return nil
}
