package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	cachepb "github.com/micro/services/cache/proto"
	pauth "github.com/micro/services/pkg/auth"
	adminpb "github.com/micro/services/pkg/service/proto"
	"github.com/micro/services/pkg/tenant"
	url "github.com/micro/services/url/proto"
	"github.com/teris-io/shortid"
)

const hostPrefix = "https://m3o.one/u/"

type Url struct {
	cache      cachepb.CacheService
	hostPrefix string
}

func NewUrl(svc *service.Service) *Url {
	var hp string

	cfg, err := config.Get("micro.url.host_prefix")
	if err != nil {
		hp = cfg.String(hostPrefix)
	}

	if len(strings.TrimSpace(hp)) == 0 {
		hp = hostPrefix
	}

	return &Url{
		cache:      cachepb.NewCacheService("cache", svc.Client()),
		hostPrefix: hp,
	}
}

func (e *Url) Shorten(ctx context.Context, req *url.ShortenRequest, rsp *url.ShortenResponse) error {
	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized("url.shorten", "not authorized")
	}
	sid, err := shortid.New(1, shortid.DefaultABC, 2342)
	if err != nil {
		return err
	}

	id, err := sid.Generate()
	if err != nil {
		return err
	}

	val := &url.URLPair{
		DestinationURL: req.DestinationURL,
		ShortURL:       e.hostPrefix + id,
		Created:        time.Now().Format(time.RFC3339Nano),
	}

	// write a global key
	key := "url/" + id
	if err := store.Write(store.NewRecord(key, val)); err != nil {
		return err
	}

	// write per owner key
	key = "urlOwner/" + tenantId + "/" + id
	if err := store.Write(store.NewRecord(key, val)); err != nil {
		return err
	}

	rsp.ShortURL = val.ShortURL

	return nil
}

func (e *Url) List(ctx context.Context, req *url.ListRequest, rsp *url.ListResponse) error {
	method := "url.shorten"
	errInternal := errors.InternalServerError(method, "Error listing URLs")
	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "not authorized")
	}

	var err error

	prefix := "urlOwner/" + tenantId + "/"
	key := prefix

	var opts []store.ReadOption

	if len(req.ShortURL) > 0 {
		id := strings.Replace(req.ShortURL, e.hostPrefix, "", -1)
		key += id
	} else {
		opts = append(opts, store.ReadPrefix())
	}

	records, err := store.Read(key, opts...)
	if err != nil {
		logger.Errorf("Error reading record %s", err)
		return errInternal
	}

	for _, rec := range records {
		uri := new(url.URLPair)

		if err := rec.Decode(uri); err != nil {
			continue
		}
		crsp, err := e.cache.Get(ctx, &cachepb.GetRequest{Key: cacheKey(strings.TrimPrefix(rec.Key, prefix))}, client.WithAuthToken())
		if err != nil {
			logger.Errorf("Error reading cache %s", err)
			return errInternal
		}
		uri.HitCount, _ = strconv.ParseInt(crsp.Value, 10, 64)
		rsp.UrlPairs = append(rsp.UrlPairs, uri)

	}

	return nil
}

func (e *Url) Proxy(ctx context.Context, req *url.ProxyRequest, rsp *url.ProxyResponse) error {
	id := strings.Replace(req.ShortURL, e.hostPrefix, "", -1)

	records, err := store.Read("url/" + id)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return errors.NotFound("url.proxy", "not found")
	}

	uri := new(url.URLPair)
	if err := records[0].Decode(uri); err != nil {
		return err
	}

	rsp.DestinationURL = uri.DestinationURL
	go func() {
		_, err := e.cache.Increment(ctx, &cachepb.IncrementRequest{
			Key:   cacheKey(id),
			Value: 1,
		}, client.WithAuthToken())
		if err != nil {
			logger.Errorf("Error incrementing cache %s", err)
		}
	}()

	return nil
}

func cacheKey(id string) string {
	logger.Infof("Cache key %s", id)
	return fmt.Sprintf("url/HitCount/%s", id)
}

func (e *Url) DeleteData(ctx context.Context, request *adminpb.DeleteDataRequest, response *adminpb.DeleteDataResponse) error {
	method := "admin.DeleteData"
	_, err := pauth.VerifyMicroAdmin(ctx, method)
	if err != nil {
		return err
	}

	if len(request.TenantId) < 10 { // deliberate length check so we don't delete all the things
		return errors.BadRequest(method, "Missing tenant ID")
	}

	prefix := "urlOwner/" + request.TenantId + "/"

	keys, err := store.List(store.ListPrefix(prefix))
	if err != nil {
		return err
	}

	for _, key := range keys {
		id := strings.TrimPrefix(key, prefix)
		if err := store.Delete("url/" + id); err != nil {
			return err
		}
		if err := store.Delete(key); err != nil {
			return err
		}
		e.cache.Delete(ctx, &cachepb.DeleteRequest{Key: cacheKey(id)})
	}
	logger.Infof("Deleted %d objects from S3 for %s", len(keys), request.TenantId)

	return nil
}
