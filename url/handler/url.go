package handler

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/model"
	cache "github.com/patrickmn/go-cache"
	"github.com/teris-io/shortid"

	"github.com/micro/services/pkg/tenant"
	url "github.com/micro/services/url/proto"
)

const hostPrefix = "https://m3o.one/u/"

type Url struct {
	pairs      model.Model
	ownerIndex model.Index
	cache      *cache.Cache
	hostPrefix string
}

func NewUrl() *Url {
	var hp string
	cfg, err := config.Get("micro.url.host_prefix")
	if err != nil {
		hp = cfg.String(hostPrefix)
	}
	if len(strings.TrimSpace(hp)) == 0 {
		hp = hostPrefix
	}

	ownerIndex := model.ByEquality("owner")
	ownerIndex.Order.Type = model.OrderTypeUnordered

	m := model.NewModel(
		model.WithKey("shortURL"),
		model.WithIndexes(ownerIndex),
	)
	m.Register(&url.URLPair{})
	return &Url{
		pairs:      m,
		ownerIndex: ownerIndex,
		hostPrefix: hp,
		cache:      cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

func (e *Url) Shorten(ctx context.Context, req *url.ShortenRequest, rsp *url.ShortenResponse) error {
	tenantID, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.New("Not authorized")
	}
	sid, err := shortid.New(1, shortid.DefaultABC, 2342)
	if err != nil {
		return err
	}

	id, err := sid.Generate()
	if err != nil {
		return err
	}

	p := &url.URLPair{
		DestinationURL: req.DestinationURL,
		ShortURL:       id,
		Owner:          tenantID,
		Created:        time.Now().Unix(),
	}
	rsp.ShortURL = e.hostPrefix + id

	return e.pairs.Create(p)
}

func (e *Url) List(ctx context.Context, req *url.ListRequest, rsp *url.ListResponse) error {
	tenantID, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.New("Not authorized")
	}

	rsp.UrlPairs = []*url.URLPair{}
	err := e.pairs.Read(e.ownerIndex.ToQuery(tenantID), &rsp.UrlPairs)
	if err != nil {
		return err
	}
	for _, v := range rsp.UrlPairs {
		// get the counter and add it to db value to improve
		// accuracy
		count, ok := e.cache.Get(v.ShortURL)
		if ok {
			v.HitCount += count.(int64)
		}
		v.ShortURL = e.hostPrefix + v.ShortURL
	}
	return nil
}

func (e *Url) Proxy(ctx context.Context, req *url.ProxyRequest, rsp *url.ProxyResponse) error {
	var pair url.URLPair
	id := strings.Replace(req.ShortURL, e.hostPrefix, "", -1)
	err := e.pairs.Read(model.QueryEquals("shortURL", id), &pair)
	if err != nil {
		return err
	}
	v, found := e.cache.Get(id)
	if !found {
		e.cache.Set(id, int64(1), cache.NoExpiration)
	} else {
		// we null out the counter
		e.cache.Set(id, 0, cache.NoExpiration)
		if v.(int64)%10 == 0 {
			go func() {
				// We add instead of set in case the service runs in multiple
				// instances
				pair.HitCount += v.(int64) + int64(1)
				err = e.pairs.Update(pair)
				if err != nil {
					logger.Error(err)
				}
			}()
		}
	}

	rsp.DestinationURL = pair.DestinationURL
	return nil
}
