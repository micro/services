// Package main is a client for the chat service to demonstrate how it would work for a client. To
// run the client, first launch the chat service by running `micro run ./chat` from the top level of
// this repo. Then run `micro run ./chat/client` and `micro logs -f client` to follow the logs of
// the client.
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/context/metadata"
	"github.com/micro/micro/v3/service/logger"
	chat "github.com/micro/services/chat/proto"
)

var (
	userOneID = "user-one-" + uuid.New().String()
	userTwoID = "user-two-" + uuid.New().String()
)

func main() {
	// create a chat service client
	srv := service.New()
	cli := chat.NewChatService("chat", srv.Client())

	// create a chat for our users
	userIDs := []string{userOneID, userTwoID}
	nRsp, err := cli.New(context.TODO(), &chat.NewRequest{UserIds: userIDs})
	if err != nil {
		logger.Fatalf("Error creating the chat: %v", err)
	}
	chatID := nRsp.GetChatId()
	logger.Infof("Chat Created. ID: %v", chatID)

	// list the number messages in the chat history
	hRsp, err := cli.History(context.TODO(), &chat.HistoryRequest{ChatId: chatID})
	if err != nil {
		logger.Fatalf("Error getting the chat history: %v", err)
	}
	logger.Infof("Chat has %v message(s)", len(hRsp.Messages))

	// create a channel to handle errors
	errChan := make(chan error)

	// run user one
	go func() {
		ctx := metadata.NewContext(context.TODO(), metadata.Metadata{
			"user-id": userOneID, "chat-id": chatID,
		})
		stream, err := cli.Connect(ctx)
		if err != nil {
			errChan <- err
			return
		}

		for i := 1; true; i++ {
			// send a message to the chat
			err = stream.Send(&chat.Message{
				ClientId: uuid.New().String(),
				SentAt:   time.Now().Unix(),
				Subject:  "Message from user one",
				Text:     fmt.Sprintf("Message #%v", i),
			})
			if err != nil {
				errChan <- err
				return
			}
			logger.Infof("User one sent message")

			// wait for user two to respond
			msg, err := stream.Recv()
			if err != nil {
				errChan <- err
				return
			}
			logger.Infof("User one recieved message %v from %v", msg.Text, msg.UserId)
			time.Sleep(time.Second)
		}
	}()

	// run user two
	go func() {
		ctx := metadata.NewContext(context.TODO(), metadata.Metadata{
			"user-id": userTwoID, "chat-id": chatID,
		})
		stream, err := cli.Connect(ctx)
		if err != nil {
			errChan <- err
			return
		}

		for i := 1; true; i++ {
			// send a response to the chat
			err = stream.Send(&chat.Message{
				ClientId: uuid.New().String(),
				SentAt:   time.Now().Unix(),
				Subject:  "Response from user two",
				Text:     fmt.Sprintf("Response #%v", i),
			})
			if err != nil {
				errChan <- err
				return
			}
			logger.Infof("User two sent message")

			// wait for a message from user one
			msg, err := stream.Recv()
			if err != nil {
				errChan <- err
				return
			}
			logger.Infof("User two recieved message %v from %v", msg.Text, msg.UserId)
			time.Sleep(time.Second)
		}
	}()

	logger.Fatal(<-errChan)
}
