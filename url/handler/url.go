package handler

import (
	"context"
	"errors"
	"strings"

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/model"
	"github.com/teris-io/shortid"

	"github.com/micro/services/pkg/tenant"
	url "github.com/micro/services/url/proto"
)

const hostPrefix = "https://m3o.one/u"

type Url struct {
	pairs      model.Model
	ownerIndex model.Index
	hostPrefix string
}

func NewUrl() *Url {
	var hp string
	cfg, err := config.Get("micro.url_shortener.host_prefix")
	if err != nil {
		hp = cfg.String(hostPrefix)
	}
	if len(strings.TrimSpace(hp)) == 0 {
		hp = hostPrefix
	}

	ownerIndex := model.ByEquality("Owner")
	ownerIndex.Order.Type = model.OrderTypeUnordered

	return &Url{
		pairs: model.NewModel(
			model.WithKey("ShortURL"),
			model.WithIndexes(ownerIndex),
		),
		ownerIndex: ownerIndex,
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
	return e.pairs.Create(&url.URLPair{
		DestinationURL: req.DestinationURL,
		ShortURL:       id,
		Owner:          tenantID,
	})
}

func (e *Url) List(ctx context.Context, req *url.ListRequest, rsp *url.ListResponse) error {
	tenantID, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.New("Not authorized")
	}

	rsp.UrlPairs = []*url.URLPair{}
	err := e.pairs.Read(e.ownerIndex.ToQuery(e.ownerIndex.ToQuery(tenantID)), &rsp.UrlPairs)
	if err != nil {
		return err
	}
	for _, v := range rsp.UrlPairs {
		v.ShortURL = e.hostPrefix + "/" + v.ShortURL
	}
	return nil
}

func (e *Url) Proxy(ctx context.Context, req *url.ProxyRequest, rsp *url.ProxyResponse) error {
	var pair url.URLPair
	err := e.pairs.Read(e.ownerIndex.ToQuery(model.QueryEquals("ShortURL", e.hostPrefix+"/"+req.ShortURL)), pair)
	if err != nil {
		return err
	}

	rsp.DestinationURL = pair.DestinationURL
	return nil
}
