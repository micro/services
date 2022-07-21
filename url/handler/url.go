package handler

import (
	"context"
	"fmt"
	u "net/url"
	"regexp"
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

var (
	idRegex = regexp.MustCompile("[a-zA-Z0-9_-]+")
)

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

func (e *Url) Delete(ctx context.Context, req *url.DeleteRequest, rsp *url.DeleteResponse) error {
	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized("url.shorten", "not authorized")
	}

	if len(req.Id) == 0 && len(req.ShortURL) == 0 {
		return errors.BadRequest("url.delete", "missing id or short url")
	}

	id := req.Id

	if len(id) == 0 {
		id = strings.TrimPrefix(req.ShortURL, e.hostPrefix)
	}

	// check if exists
	recs, err := store.Read("urlOwner/" + tenantId + "/" + id)
	if err != nil || len(recs) == 0 {
		// swallow error
		return nil
	}

	// delete the url
	store.Delete("url/" + id)

	// delete the owner
	store.Delete("urlOwner/" + tenantId + "/" + id)

	return nil
}

func (e *Url) Create(ctx context.Context, req *url.CreateRequest, rsp *url.CreateResponse) error {
	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized("url.create", "not authorized")
	}

	if len(req.Id) == 0 {
		sid, err := shortid.New(1, shortid.DefaultABC, 2342)
		if err != nil {
			return err
		}

		id, err := sid.Generate()
		if err != nil {
			return err
		}
		req.Id = id
	}

	if !idRegex.MatchString(req.Id) {
		return errors.BadRequest("url.create", "invalid id")
	}
	_, err := u.Parse(req.DestinationURL)
	if err != nil {
		return errors.BadRequest("url.create", err.Error())
	}

	// the url id
	id := req.Id

	records, err := store.Read("url/" + req.Id)
	if err != nil && err != store.ErrNotFound {
		return err
	}

	if len(records) > 0 {
		return errors.BadRequest("url.create", "id already exists")
	}

	val := &url.URLPair{
		Id:             id,
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

	rsp.Url = val

	return nil
}

func (e *Url) Update(ctx context.Context, req *url.UpdateRequest, rsp *url.UpdateResponse) error {
	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.Unauthorized("url.shorten", "not authorized")
	}

	if len(req.Id) == 0 && len(req.ShortURL) == 0 {
		return errors.BadRequest("url.update", "missing id or short url")
	}

	id := req.Id

	if len(id) == 0 {
		id = strings.Replace(req.ShortURL, e.hostPrefix, "", -1)
	}

	// check the owner has this short url
	records, err := store.Read("urlOwner/" + tenantId + "/" + id)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return errors.NotFound("url.update", "not found")
	}

	uri := new(url.URLPair)
	if err := records[0].Decode(uri); err != nil {
		return err
	}

	// set the destination url
	uri.DestinationURL = req.DestinationURL

	// write a global key
	key := "url/" + id
	if err := store.Write(store.NewRecord(key, uri)); err != nil {
		return err
	}

	// write per owner key
	key = "urlOwner/" + tenantId + "/" + id
	if err := store.Write(store.NewRecord(key, uri)); err != nil {
		return err
	}

	return nil
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
		Id:             id,
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

		// set the id if not exists
		if len(uri.Id) == 0 {
			id := strings.Replace(req.ShortURL, e.hostPrefix, "", -1)
			uri.Id = id
		}

		uri.HitCount, _ = strconv.ParseInt(crsp.Value, 10, 64)
		rsp.UrlPairs = append(rsp.UrlPairs, uri)

	}

	return nil
}

func (e *Url) Resolve(ctx context.Context, req *url.ResolveRequest, rsp *url.ResolveResponse) error {
	id := strings.Replace(req.ShortURL, e.hostPrefix, "", -1)

	records, err := store.Read("url/" + id)
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return errors.NotFound("url.resolve", "not found")
	}

	uri := new(url.URLPair)
	if err := records[0].Decode(uri); err != nil {
		return err
	}

	rsp.DestinationURL = uri.DestinationURL

	go func() {
		_, err := e.cache.Increment(context.Background(), &cachepb.IncrementRequest{
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
		e.cache.Delete(ctx, &cachepb.DeleteRequest{Key: cacheKey(id)}, client.WithAuthToken())
	}
	logger.Infof("Deleted %d objects from S3 for %s", len(keys), request.TenantId)

	return nil
}

func (e *Url) Usage(ctx context.Context, request *adminpb.UsageRequest, response *adminpb.UsageResponse) error {
	return nil
}
