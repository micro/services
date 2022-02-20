package handler

import (
	"context"
	"encoding/json"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	pb "github.com/micro/services/lists/proto"
	streamPb "github.com/micro/services/mq/proto"
	pauth "github.com/micro/services/pkg/auth"
	adminpb "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/pkg/tenant"
	"google.golang.org/protobuf/types/known/structpb"
)

// New returns an initialized Lists
func New(c client.Client) *Lists {
	return &Lists{
		Stream: streamPb.NewMqService("mq", c),
	}
}

// Lists implements the lists proto definition
type Lists struct {
	Stream streamPb.MqService
}

func newMessage(ev map[string]interface{}) *structpb.Struct {
	st := new(structpb.Struct)
	b, _ := json.Marshal(ev)
	json.Unmarshal(b, st)
	return st
}

// Create inserts a new list in the store
func (h *Lists) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	if len(req.Name) == 0 && len(req.Items) == 0 {
		return errors.BadRequest("lists.create", "missing name and text")
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
	// set the generated fields on the list
	list := &pb.List{
		Id:      id.String(),
		Created: t,
		Updated: t,
		Name:    req.Name,
		Items:   req.Items,
	}

	key := path.Join("list", tnt, id.String())
	rec := store.NewRecord(key, list)

	if err = store.Write(rec); err != nil {
		return errors.InternalServerError("lists.created", "failed to create list")
	}

	// return the list in the response
	rsp.List = list

	h.Stream.Publish(ctx, &streamPb.PublishRequest{
		Topic: "lists",
		Message: newMessage(map[string]interface{}{
			"event": "create",
			"list":  list,
		}),
	})

	return nil
}

func (h *Lists) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	if len(req.Id) == 0 {
		return errors.BadRequest("lists.read", "Missing List ID")
	}

	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "default"
	}

	key := path.Join("list", tnt, req.Id)

	// read the specific list
	recs, err := store.Read(key)
	if err == store.ErrNotFound {
		return errors.NotFound("lists.read", "List not found")
	} else if err != nil {
		return errors.InternalServerError("lists.read", "Error reading from store: %v", err.Error())
	}

	// Decode the list
	var list *pb.List
	if err := recs[0].Decode(&list); err != nil {
		return errors.InternalServerError("lists.update", "Error unmarshaling JSON: %v", err.Error())
	}

	// return the list
	rsp.List = list

	return nil
}

// Update is a unary API which updates a list in the store
func (h *Lists) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	// Validate the request
	if req.List == nil {
		return errors.BadRequest("lists.update", "Missing List")
	}
	if len(req.List.Id) == 0 {
		return errors.BadRequest("lists.update", "Missing List ID")
	}

	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "default"
	}

	key := path.Join("list", tnt, req.List.Id)

	// read the specific list
	recs, err := store.Read(key)
	if err == store.ErrNotFound {
		return errors.NotFound("lists.update", "List not found")
	} else if err != nil {
		return errors.InternalServerError("lists.update", "Error reading from store: %v", err.Error())
	}

	// Decode the list
	var list *pb.List
	if err := recs[0].Decode(&list); err != nil {
		return errors.InternalServerError("lists.update", "Error unmarshaling JSON: %v", err.Error())
	}

	// Update the lists name and text
	list.Name = req.List.Name
	list.Items = req.List.Items
	list.Updated = time.Now().Format(time.RFC3339)

	rec := store.NewRecord(key, list)

	// Write the updated list to the store
	if err = store.Write(rec); err != nil {
		return errors.InternalServerError("lists.update", "Error writing to store: %v", err.Error())
	}

	h.Stream.Publish(ctx, &streamPb.PublishRequest{
		Topic: "lists",
		Message: newMessage(map[string]interface{}{
			"event": "update",
			"list":  list,
		}),
	})

	rsp.List = list

	return nil
}

func (h *Lists) Events(ctx context.Context, req *pb.EventsRequest, stream pb.Lists_EventsStream) error {
	backendStream, err := h.Stream.Subscribe(ctx, &streamPb.SubscribeRequest{
		Topic: "lists",
	})
	if err != nil {
		return errors.InternalServerError("lists.subscribe", "Failed to subscribe to lists")
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

		list := rsp.List

		// filter if necessary by id
		if len(req.Id) > 0 && list.Id != req.Id {
			continue
		}

		// send back the event to the client
		if err := stream.Send(rsp); err != nil {
			return nil
		}
	}

	return nil
}

// Delete removes the list from the store, looking up using ID
func (h *Lists) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	// Validate the request
	if len(req.Id) == 0 {
		return errors.BadRequest("lists.delete", "Missing List ID")
	}

	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "default"
	}

	key := path.Join("list", tnt, req.Id)

	// read the specific list
	recs, err := store.Read(key)
	if err == store.ErrNotFound {
		return nil
	} else if err != nil {
		return errors.InternalServerError("lists.delete", "Error reading from store: %v", err.Error())
	}

	// Decode the list
	var list *pb.List
	if err := recs[0].Decode(&list); err != nil {
		return errors.InternalServerError("lists.delete", "Error unmarshaling JSON: %v", err.Error())
	}

	// now delete it
	if err := store.Delete(key); err != nil && err != store.ErrNotFound {
		return errors.InternalServerError("lists.delete", "Failed to delete list")
	}

	h.Stream.Publish(ctx, &streamPb.PublishRequest{
		Topic: "lists",
		Message: newMessage(map[string]interface{}{
			"event": "delete",
			"list":  list,
		}),
	})

	rsp.List = list

	return nil
}

// List returns all of the lists in the store
func (h *Lists) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		tnt = "default"
	}

	key := path.Join("list", tnt) + "/"

	// Retrieve all of the records in the store
	recs, err := store.Read(key, store.ReadPrefix())
	if err != nil {
		return errors.InternalServerError("lists.list", "Error reading from store: %v", err.Error())
	}

	// Initialize the response lists slice
	rsp.Lists = make([]*pb.List, len(recs))

	// Unmarshal the lists into the response
	for i, r := range recs {
		if err := r.Decode(&rsp.Lists[i]); err != nil {
			return errors.InternalServerError("lists.list", "Error decoding list: %v", err.Error())
		}
	}

	return nil
}

func (h *Lists) DeleteData(ctx context.Context, request *adminpb.DeleteDataRequest, response *adminpb.DeleteDataResponse) error {
	method := "admin.DeleteData"
	_, err := pauth.VerifyMicroAdmin(ctx, method)
	if err != nil {
		return err
	}

	if len(request.TenantId) < 10 { // deliberate length check so we don't delete all the things
		return errors.BadRequest(method, "Missing tenant ID")
	}

	keys, err := store.List(store.ListPrefix(path.Join("list", request.TenantId) + "/"))
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
