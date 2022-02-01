package handler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/cache/proto"
	pauth "github.com/micro/services/pkg/auth"
	"github.com/micro/services/pkg/cache"
	adminpb "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/pkg/tenant"
)

type Cache struct{}

func (c *Cache) Get(ctx context.Context, req *pb.GetRequest, rsp *pb.GetResponse) error {
	if len(req.Key) == 0 {
		return errors.BadRequest("cache.get", "missing key")
	}

	var value interface{}

	expires, err := cache.Context(ctx).Get(req.Key, &value)

	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			log.Errorf("Error querying cache %s", err)
			return errors.InternalServerError("cache.get", "Error querying cache")
		}
		value = ""
	}

	rsp.Key = req.Key
	// set the value
	rsp.Value = fmt.Sprintf("%v", value)
	// set the ttl
	rsp.Ttl = int64(expires.Sub(time.Now()).Seconds())

	if rsp.Ttl < 0 {
		rsp.Ttl = 0
	}

	return nil
}

func (c *Cache) Set(ctx context.Context, req *pb.SetRequest, rsp *pb.SetResponse) error {
	if len(req.Key) == 0 {
		return errors.BadRequest("cache.set", "missing key")
	}

	// max size 1mb e.g byte * 1024 * 1024
	if len(req.Value) > 1e6 {
		return errors.BadRequest("cache.set", "value is too big")
	}

	ttl := time.Time{}

	if req.Ttl > 0 {
		ttl = time.Now().Add(time.Duration(req.Ttl) * time.Second)
	}

	if err := cache.Context(ctx).Set(req.Key, req.Value, ttl); err != nil {
		log.Errorf("Error writing to cache %s", err)
		return errors.InternalServerError("cache.set", "Error writing to cache")
	}

	rsp.Status = "ok"

	return nil
}

func (c *Cache) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	if len(req.Key) == 0 {
		return errors.BadRequest("cache.delete", "missing key")
	}
	if err := cache.Context(ctx).Delete(req.Key); err != nil {
		log.Errorf("Error deleting from cache %s", err)
		return errors.InternalServerError("cache.delete", "Error deleting from cache")
	}

	rsp.Status = "ok"

	return nil
}

func (c *Cache) Increment(ctx context.Context, req *pb.IncrementRequest, rsp *pb.IncrementResponse) error {
	if len(req.Key) == 0 {
		return errors.BadRequest("cache.increment", "missing key")
	}

	// increment the value
	v, err := cache.Context(ctx).Increment(req.Key, req.Value)
	if err != nil {
		log.Errorf("Error incrementing cache %s", err)
		return errors.InternalServerError("cache.increment", "Error incrementing cache")
	}

	// set the response value
	rsp.Key = req.Key
	rsp.Value = v

	return nil
}

func (c *Cache) Decrement(ctx context.Context, req *pb.DecrementRequest, rsp *pb.DecrementResponse) error {
	if len(req.Key) == 0 {
		return errors.BadRequest("cache.decrement", "missing key")
	}

	v, err := cache.Context(ctx).Decrement(req.Key, req.Value)
	if err != nil {
		log.Errorf("Error decrementing cache %s", err)
		return errors.InternalServerError("cache.decrement", "Error decrementing cache")
	}

	// set the response value
	rsp.Key = req.Key
	rsp.Value = v

	return nil
}

func (c *Cache) ListKeys(ctx context.Context, req *pb.ListKeysRequest, rsp *pb.ListKeysResponse) error {
	keys, err := cache.Context(ctx).ListKeys()

	if err != nil {
		log.Errorf("Error listing keys in cache %s", err)
		return errors.InternalServerError("cache.listkeys", "Error listing keys in cache")
	}

	rsp.Keys = keys

	return nil
}

func (c *Cache) DeleteData(ctx context.Context, request *adminpb.DeleteDataRequest, response *adminpb.DeleteDataResponse) error {
	method := "admin.DeleteData"
	_, err := pauth.VerifyMicroAdmin(ctx, method)
	if err != nil {
		return err
	}

	if len(request.TenantId) == 0 {
		return errors.BadRequest(method, "Missing tenant ID")
	}

	split := strings.Split(request.TenantId, "/")
	tenant.NewContext(split[1], split[0], split[1])
	keys, err := cache.Context(ctx).ListKeys()
	if err != nil {
		return err
	}
	for _, k := range keys {
		if err := cache.Context(ctx).Delete(k); err != nil {
			return err
		}
	}
	return nil
}
