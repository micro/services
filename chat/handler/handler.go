package handler

import (
	"context"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/events"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	pb "github.com/micro/services/chat/proto"
	"github.com/micro/services/pkg/tenant"
)

const (
	chatStoreKeyPrefix    = "chats/"
	chatEventKeyPrefix    = "chats/"
	messageStoreKeyPrefix = "messages/"
)

type Chat struct{}

func (c *Chat) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// get the tenant
	tenantId := tenant.Id(ctx)

	// generate a unique id for the chat
	roomId := uuid.New().String()

	// create a new room
	room := &pb.Room{
		Id:          roomId,
		Name:        req.Name,
		Description: req.Description,
		UserIds:     req.UserIds,
		Private:     req.Private,
		CreatedAt:   time.Now().Format(time.RFC3339Nano),
	}

	// key to lookup the chat in the store using, e.g. "chat/usera-userb-userc"
	key := path.Join(chatStoreKeyPrefix, tenantId, roomId)

	// create a new record for the room
	rec := store.NewRecord(key, room)

	// write a record for the new room
	if err := store.Write(rec); err != nil {
		logger.Errorf("Error writing to the store. Key: %v. Error: %v", key, err)
		return errors.InternalServerError("chat.new", "error creating chat room")
	}

	// return the room
	rsp.Room = room

	return nil
}

func (c *Chat) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	tenantId := tenant.Id(ctx)

	// validate the request
	if len(req.RoomId) == 0 {
		return errors.BadRequest("chat.delete", "missing room id")
	}

	key := path.Join(chatStoreKeyPrefix, tenantId, req.RoomId)

	// lookup the chat from the store to ensure it's valid
	recs, err := store.Read(key, store.ReadLimit(1))
	if err == store.ErrNotFound {
		return errors.BadRequest("chat.delete", "room not found")
	} else if err != nil {
		logger.Errorf("Error reading from the store. Room ID: %v. Error: %v", req.RoomId, err)
		return errors.InternalServerError("chat.delete", "error reading chat room")
	}

	room := new(pb.Room)
	err = recs[0].Decode(room)
	if err != nil {
		return errors.InternalServerError("chat.delete", "error reading chat room")
	}
	// set response
	rsp.Room = room

	// delete the room
	if err := store.Delete(key); err != nil {
		return errors.InternalServerError("chat.delete", "error deleting chat room")
	}

	// get all messages
	// TODO: paginate the list
	key = path.Join(messageStoreKeyPrefix, tenantId, req.RoomId)
	srecs, err := store.List(store.ListPrefix(key))
	if err != nil {
		return errors.InternalServerError("chat.delete", "failed to list messages")
	}

	// delete all the messages
	for _, rec := range srecs {
		if err := store.Delete(rec); err != nil {
			return errors.InternalServerError("chat.delete", "failed to list messages")
		}
	}

	// TODO: notify users of the event that the room is deleted

	return nil
}

func (c *Chat) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	// get the tenant
	tenantId := tenant.Id(ctx)

	key := path.Join(chatStoreKeyPrefix, tenantId) + "/"

	// read all the rooms from the store for the user
	recs, err := store.Read(key, store.ReadPrefix())
	if err != nil {
		return errors.InternalServerError("chat.list", "error listing chat rooms")
	}

	// list all the rooms
	for _, rec := range recs {
		room := new(pb.Room)
		err := rec.Decode(room)
		if err != nil {
			continue
		}

		if len(req.UserId) == 0 {
			rsp.Rooms = append(rsp.Rooms, room)
			continue
		}

		// check if there's a user id match
		for _, user := range room.UserIds {
			if user == req.UserId {
				rsp.Rooms = append(rsp.Rooms, room)
				break
			}
		}
	}

	return nil
}

// History returns the historical messages in a chat
func (c *Chat) History(ctx context.Context, req *pb.HistoryRequest, rsp *pb.HistoryResponse) error {
	// get the tenant
	tenantId := tenant.Id(ctx)

	// validate the request
	if len(req.RoomId) == 0 {
		return errors.BadRequest("chat.history", "missing room id")
	}

	key := path.Join(chatStoreKeyPrefix, tenantId, req.RoomId)

	// lookup the chat from the store to ensure it's valid
	if _, err := store.Read(key); err == store.ErrNotFound {
		return errors.BadRequest("chat.history", "room not found")
	} else if err != nil {
		logger.Errorf("Error reading from the store. Room ID: %v. Error: %v", req.RoomId, err)
		return errors.InternalServerError("chat.history", "error reading chat room")
	}

	// lookup the messages
	key = path.Join(messageStoreKeyPrefix, tenantId, req.RoomId)
	recs, err := store.Read(key+"/", store.ReadPrefix())
	if err != nil {
		logger.Errorf("Error reading messages the store. Room ID: %v. Error: %v", req.RoomId, err)
		return errors.InternalServerError("chat.history", "failed to read messages")
	}

	for _, rec := range recs {
		msg := new(pb.Message)
		err := rec.Decode(msg)
		if err != nil {
			return errors.InternalServerError("chat.history", "failed to decode message")
		}
		rsp.Messages = append(rsp.Messages, msg)
	}

	return nil
}

func (c *Chat) Invite(ctx context.Context, req *pb.InviteRequest, rsp *pb.InviteResponse) error {
	// get the tenant
	tenantId := tenant.Id(ctx)

	// validate the request
	if len(req.RoomId) == 0 {
		return errors.BadRequest("chat.invite", "missing room id")
	}

	if len(req.UserId) == 0 {
		return errors.BadRequest("chat.invite", "missing user id")
	}

	key := path.Join(chatStoreKeyPrefix, tenantId, req.RoomId)

	// lookup the chat from the store to ensure it's valid
	recs, err := store.Read(key)
	if err == store.ErrNotFound {
		return errors.BadRequest("chat.invite", "room not found")
	} else if err != nil {
		logger.Errorf("Error reading from the store. Room ID: %v. Error: %v", req.RoomId, err)
		return errors.InternalServerError("chat.invite", "error reading chat room")
	}

	// check the user is in the room
	room := new(pb.Room)
	err = recs[0].Decode(room)
	if err != nil {
		return errors.InternalServerError("chat.invite", "Error reading room")
	}

	var exists bool

	// check the user is in the room
	for _, user := range room.UserIds {
		if user == req.UserId {
			exists = true
			break
		}
	}

	// TODO: send join message
	if !exists {
		room.UserIds = append(room.UserIds, req.UserId)
		// write the record
		rec := store.NewRecord(key, room)
		if err := store.Write(rec); err != nil {
			return errors.InternalServerError("chat.invite", "Error adding user to room")
		}
	}

	rsp.Room = room

	return nil
}

// Send a single message to the chat, designed for ease of use via the API / CLI
func (c *Chat) Send(ctx context.Context, req *pb.SendRequest, rsp *pb.SendResponse) error {
	// get the tenant
	tenantId := tenant.Id(ctx)

	// validate the request
	if len(req.RoomId) == 0 {
		return errors.BadRequest("chat.send", "missing room id")
	}
	if len(req.UserId) == 0 {
		return errors.BadRequest("chat.send", "missing user id")
	}
	if len(req.Text) == 0 {
		return errors.BadRequest("chat.send", "missing text")
	}

	// check the room exists
	key := path.Join(chatStoreKeyPrefix, tenantId, req.RoomId)

	// lookup the chat room from the store to ensure it's valid
	recs, err := store.Read(key, store.ReadLimit(1))
	if err == store.ErrNotFound {
		return errors.BadRequest("chat.send", "room not found")
	} else if err != nil {
		logger.Errorf("Error reading from the store. Room ID: %v. Error: %v", req.RoomId, err)
		return errors.InternalServerError("chat.send", "error reading chat room")
	}

	// decode the room
	room := new(pb.Room)
	err = recs[0].Decode(room)
	if err != nil {
		return errors.InternalServerError("chat.send", "error reading chat room")
	}

	var exists bool

	// check the user is in the room
	for _, user := range room.UserIds {
		if user == req.UserId {
			exists = true
			break
		}
	}

	if !exists {
		return errors.BadRequest("chat.send", "user is not in the room")
	}

	// construct the message
	msg := &pb.Message{
		Id:      uuid.New().String(),
		Client:  req.Client,
		RoomId:  req.RoomId,
		UserId:  req.UserId,
		Subject: req.Subject,
		Text:    req.Text,
		SentAt:  time.Now().Format(time.RFC3339Nano),
	}

	// default the client id if not provided
	if len(msg.Client) == 0 {
		msg.Client = uuid.New().String()
	}

	// create the message
	if err := c.createMessage(tenantId, msg); err != nil {
		return err
	}

	// return the response
	rsp.Message = msg

	return nil
}

func (c *Chat) Join(ctx context.Context, req *pb.JoinRequest, stream pb.Chat_JoinStream) error {
	// get the tenant
	tenantId := tenant.Id(ctx)

	// validate the request
	if len(req.RoomId) == 0 {
		return errors.BadRequest("chat.send", "missing room id")
	}
	if len(req.UserId) == 0 {
		return errors.BadRequest("chat.send", "missing user id")
	}

	key := path.Join(chatStoreKeyPrefix, tenantId, req.RoomId)

	// lookup the chat from the store to ensure it's valid
	recs, err := store.Read(key, store.ReadLimit(1))
	if err == store.ErrNotFound {
		return errors.BadRequest("chat.join", "room not found")
	} else if err != nil {
		logger.Errorf("Error reading from the store. Room ID: %v. Error: %v", req.RoomId, err)
		return errors.InternalServerError("chat.join", "Error reading room")
	}

	// check the user is in the room
	room := new(pb.Room)
	err = recs[0].Decode(room)
	if err != nil {
		return errors.InternalServerError("chat.join", "Error reading room")
	}

	var exists bool

	// check the user is in the room
	for _, user := range room.UserIds {
		if user == req.UserId {
			exists = true
			break
		}
	}

	// TODO: send join message
	if !exists {
		room.UserIds = append(room.UserIds, req.UserId)
		// write the record
		rec := store.NewRecord(key, room)
		if err := store.Write(rec); err != nil {
			return errors.InternalServerError("chat.join", "Error adding user to room")
		}
	}

	// create a channel to send errors on, because the subscriber / publisher will run in seperate go-
	// routines, they need a way of returning errors to the client
	errChan := make(chan error)

	eventKey := path.Join(chatEventKeyPrefix, tenantId, req.RoomId)

	// create an event stream to consume messages posted by other users into the chat. we'll use the
	// user id as a queue to ensure each user recieves the message
	evStream, err := events.Consume(eventKey, events.WithGroup(req.UserId), events.WithContext(ctx))
	if err != nil {
		logger.Errorf("Error streaming events. Room ID: %v. Error: %v", req.RoomId, err)
		return errors.InternalServerError("chat.join", "Error joining the room")
	}

	for {
		select {
		case <-ctx.Done():
			// the context has been cancelled or timed out, stop subscribing to new messages
			return nil
		case ev := <-evStream:
			// recieved a message, unmarshal it into a message struct. if an error occurs log it and
			// cancel the context
			var msg pb.Message
			if err := ev.Unmarshal(&msg); err != nil {
				logger.Errorf("Error unmarshaling message. Room ID: %v. Error: %v", req.RoomId, err)
				errChan <- err
				return nil
			}

			// ignore any messages published by the current user
			if msg.UserId == req.UserId {
				continue
			}

			// publish the message to the stream
			if err := stream.Send(&pb.JoinResponse{Message: &msg}); err != nil {
				logger.Errorf("Error sending message to stream. ChatID: %v. Message ID: %v. Error: %v", msg.RoomId, msg.Id, err)
				errChan <- err
				return nil
			}
		}
	}

	return nil
}

func (c *Chat) Kick(ctx context.Context, req *pb.KickRequest, rsp *pb.KickResponse) error {
	// get the tenant
	tenantId := tenant.Id(ctx)

	// validate the request
	if len(req.RoomId) == 0 {
		return errors.BadRequest("chat.kick", "missing room id")
	}
	if len(req.UserId) == 0 {
		return errors.BadRequest("chat.kick", "missing user id")
	}

	key := path.Join(chatStoreKeyPrefix, tenantId, req.RoomId)

	// lookup the chat from the store to ensure it's valid
	recs, err := store.Read(key, store.ReadLimit(1))
	if err == store.ErrNotFound {
		return errors.BadRequest("chat.kick", "room not found")
	} else if err != nil {
		logger.Errorf("Error reading from the store. Chat ID: %v. Error: %v", req.RoomId, err)
		return errors.InternalServerError("chat.kick", "Error reading room")
	}

	// check the user is in the room
	room := new(pb.Room)
	err = recs[0].Decode(room)
	if err != nil {
		return errors.InternalServerError("chat.kick", "Error reading room")
	}

	var users []string

	// check the user is in the room
	for _, user := range room.UserIds {
		if user == req.UserId {
			continue
		}
		users = append(users, user)
	}

	room.UserIds = users

	rec := store.NewRecord(key, room)
	if err := store.Write(rec); err != nil {
		return errors.InternalServerError("chat.kick", "Error leaveing from room")
	}

	// TODO: send leave message
	// TODO: disconnect the actual event consumption
	rsp.Room = room

	return nil
}
func (c *Chat) Leave(ctx context.Context, req *pb.LeaveRequest, rsp *pb.LeaveResponse) error {
	// get the tenant
	tenantId := tenant.Id(ctx)

	// validate the request
	if len(req.RoomId) == 0 {
		return errors.BadRequest("chat.leave", "missing room id")
	}
	if len(req.UserId) == 0 {
		return errors.BadRequest("chat.leave", "missing user id")
	}

	key := path.Join(chatStoreKeyPrefix, tenantId, req.RoomId)

	// lookup the chat from the store to ensure it's valid
	recs, err := store.Read(key, store.ReadLimit(1))
	if err == store.ErrNotFound {
		return errors.BadRequest("chat.leave", "room not found")
	} else if err != nil {
		logger.Errorf("Error reading from the store. Chat ID: %v. Error: %v", req.RoomId, err)
		return errors.InternalServerError("chat.leave", "Error reading room")
	}

	// check the user is in the room
	room := new(pb.Room)
	err = recs[0].Decode(room)
	if err != nil {
		return errors.InternalServerError("chat.leave", "Error reading room")
	}

	var users []string

	// check the user is in the room
	for _, user := range room.UserIds {
		if user == req.UserId {
			continue
		}
		users = append(users, user)
	}

	room.UserIds = users

	rec := store.NewRecord(key, room)
	if err := store.Write(rec); err != nil {
		return errors.InternalServerError("chat.leave", "Error leaveing from room")
	}

	// TODO: send leave message
	// TODO: disconnect the actual event consumption
	rsp.Room = room

	return nil
}

// createMessage is a helper function which creates a message in the event stream. It handles the
// logic for ensuring client id is unique.
func (c *Chat) createMessage(tenantId string, msg *pb.Message) error {
	storekey := path.Join(messageStoreKeyPrefix, tenantId, msg.RoomId, msg.Id)
	eventKey := path.Join(chatEventKeyPrefix, tenantId, msg.RoomId)

	// send the message to the event stream
	if err := events.Publish(eventKey, msg); err != nil {
		return err
	}

	// create a new record
	rec := store.NewRecord(storekey, msg)

	// record the messages client id
	return store.Write(rec)
}
