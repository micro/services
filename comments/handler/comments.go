package handler

import (
	"context"
	"encoding/json"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v5/service/client"
	"github.com/micro/micro/v5/service/errors"
	"github.com/micro/micro/v5/service/logger"
	"github.com/micro/micro/v5/service/store"
	pb "github.com/micro/services/comments/proto"
	streamPb "github.com/micro/services/mq/proto"
	pauth "github.com/micro/services/pkg/auth"
	adminpb "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/pkg/tenant"
	"google.golang.org/protobuf/types/known/structpb"
)

// New returns an initialized Comments
func New(c client.Client) *Comments {
	return &Comments{
		Stream: streamPb.NewMqService("mq", c),
	}
}

// Comments implements the comments proto definition
type Comments struct {
	Stream streamPb.MqService
}

func newMessage(ev map[string]interface{}) *structpb.Struct {
	st := new(structpb.Struct)
	b, _ := json.Marshal(ev)
	json.Unmarshal(b, st)
	return st
}

// Create inserts a new comment in the store
func (h *Comments) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	if len(req.Subject) == 0 && len(req.Text) == 0 {
		return errors.BadRequest("comments.create", "missing name and text")
	}

	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "default"
	}

	// generate a key (uuid v4)
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	t := time.Now().Format(time.RFC3339)
	// set the generated fields on the comment
	comment := &pb.Comment{
		Id:      id.String(),
		Created: t,
		Updated: t,
		Subject: req.Subject,
		Text:    req.Text,
	}

	key := path.Join("comment", tnt, id.String())
	rec := store.NewRecord(key, comment)

	if err = store.Write(rec); err != nil {
		return errors.InternalServerError("comments.created", "failed to create comment")
	}

	// return the comment in the response
	rsp.Comment = comment

	h.Stream.Publish(ctx, &streamPb.PublishRequest{
		Topic: "comments",
		Message: newMessage(map[string]interface{}{
			"event":   "create",
			"comment": comment,
		}),
	})

	return nil
}

func (h *Comments) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("comments.read", "Missing Comment ID")
	}

	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "default"
	}

	key := path.Join("comment", tnt, req.Id)

	// read the specific comment
	recs, err := store.Read(key)
	if err == store.ErrNotFound {
		return errors.NotFound("comments.read", "Comment not found")
	} else if err != nil {
		return errors.InternalServerError("comments.read", "Error reading from store: %v", err.Error())
	}

	// Decode the comment
	var comment *pb.Comment
	if err := recs[0].Decode(&comment); err != nil {
		return errors.InternalServerError("comments.update", "Error unmarshaling JSON: %v", err.Error())
	}

	// return the comment
	rsp.Comment = comment

	return nil
}

// Update is a unary API which updates a comment in the store
func (h *Comments) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	// Validate the request
	if req.Comment == nil {
		return errors.BadRequest("comments.update", "Missing Comment")
	}
	if len(req.Comment.Id) == 0 {
		return errors.BadRequest("comments.update", "Missing Comment ID")
	}

	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "default"
	}

	key := path.Join("comment", tnt, req.Comment.Id)

	// read the specific comment
	recs, err := store.Read(key)
	if err == store.ErrNotFound {
		return errors.NotFound("comments.update", "Comment not found")
	} else if err != nil {
		return errors.InternalServerError("comments.update", "Error reading from store: %v", err.Error())
	}

	// Decode the comment
	var comment *pb.Comment
	if err := recs[0].Decode(&comment); err != nil {
		return errors.InternalServerError("comments.update", "Error unmarshaling JSON: %v", err.Error())
	}

	// Update the comments name and text
	comment.Subject = req.Comment.Subject
	comment.Text = req.Comment.Text
	comment.Updated = time.Now().Format(time.RFC3339)

	rec := store.NewRecord(key, comment)

	// Write the updated comment to the store
	if err = store.Write(rec); err != nil {
		return errors.InternalServerError("comments.update", "Error writing to store: %v", err.Error())
	}

	h.Stream.Publish(ctx, &streamPb.PublishRequest{
		Topic: "comments",
		Message: newMessage(map[string]interface{}{
			"event":   "update",
			"comment": comment,
		}),
	})

	rsp.Comment = comment

	return nil
}

func (h *Comments) Events(ctx context.Context, req *pb.EventsRequest, stream pb.Comments_EventsStream) error {
	backendStream, err := h.Stream.Subscribe(ctx, &streamPb.SubscribeRequest{
		Topic: "comments",
	})
	if err != nil {
		return errors.InternalServerError("comments.subscribe", "Failed to subscribe to comments")
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		// receive messages from the stream
		msg, err := backendStream.Recv()
		if err != nil {
			return nil
		}

		v, err := msg.Message.MarshalJSON()
		if err != nil {
			continue
		}

		rsp := new(pb.EventsResponse)

		if err := json.Unmarshal(v, rsp); err != nil {
			continue
		}

		comment := rsp.Comment

		// filter if necessary by id
		if len(req.Id) > 0 && comment.Id != req.Id {
			continue
		}

		// send back the event to the client
		if err := stream.Send(rsp); err != nil {
			return nil
		}
	}

	return nil
}

// Delete removes the comment from the store, looking up using ID
func (h *Comments) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	// Validate the request
	if len(req.Id) == 0 {
		return errors.BadRequest("comments.delete", "Missing Comment ID")
	}

	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "default"
	}

	key := path.Join("comment", tnt, req.Id)

	// read the specific comment
	recs, err := store.Read(key)
	if err == store.ErrNotFound {
		return nil
	} else if err != nil {
		return errors.InternalServerError("comments.delete", "Error reading from store: %v", err.Error())
	}

	// Decode the comment
	var comment *pb.Comment
	if err := recs[0].Decode(&comment); err != nil {
		return errors.InternalServerError("comments.delete", "Error unmarshaling JSON: %v", err.Error())
	}

	// now delete it
	if err := store.Delete(key); err != nil && err != store.ErrNotFound {
		return errors.InternalServerError("comments.delete", "Failed to delete comment")
	}

	h.Stream.Publish(ctx, &streamPb.PublishRequest{
		Topic: "comments",
		Message: newMessage(map[string]interface{}{
			"event":   "delete",
			"comment": comment,
		}),
	})

	rsp.Comment = comment

	return nil
}

// Comment returns all of the comments in the store
func (h *Comments) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "default"
	}

	key := path.Join("comment", tnt) + "/"

	// Retrieve all of the records in the store
	recs, err := store.Read(key, store.ReadPrefix())
	if err != nil {
		return errors.InternalServerError("comments.list", "Error reading from store: %v", err.Error())
	}

	// Initialize the response comments slice
	rsp.Comments = make([]*pb.Comment, len(recs))

	// Unmarshal the comments into the response
	for i, r := range recs {
		if err := r.Decode(&rsp.Comments[i]); err != nil {
			return errors.InternalServerError("comments.list", "Error decoding comment: %v", err.Error())
		}
	}

	return nil
}

func (h *Comments) DeleteData(ctx context.Context, request *adminpb.DeleteDataRequest, response *adminpb.DeleteDataResponse) error {
	method := "admin.DeleteData"
	_, err := pauth.VerifyMicroAdmin(ctx, method)
	if err != nil {
		return err
	}

	if len(request.TenantId) < 10 { // deliberate length check so we don't delete all the things
		return errors.BadRequest(method, "Missing tenant ID")
	}

	keys, err := store.List(store.ListPrefix(path.Join("comment", request.TenantId) + "/"))
	if err != nil {
		return err
	}

	for _, k := range keys {
		if err := store.Delete(k); err != nil {
			return err
		}
	}

	logger.Infof("Deleted %d keys for %s", len(keys), request.TenantId)
	return nil
}

func (h *Comments) Usage(ctx context.Context, request *adminpb.UsageRequest, response *adminpb.UsageResponse) error {
	return nil
}
