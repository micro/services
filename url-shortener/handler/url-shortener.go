package handler

import (
	"context"
	"errors"
	"strings"

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/model"
	"github.com/teris-io/shortid"

	"github.com/micro/services/pkg/tenant"
	urlshortener "github.com/micro/services/url-shortener/proto"
)

const hostPrefix = "https://cdn.m3ocontent.com"

type UrlShortener struct {
	pairs      model.Model
	ownerIndex model.Index
	hostPrefix string
}

func NewUrlShortener() *UrlShortener {
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

	return &UrlShortener{
		pairs: model.NewModel(
			model.WithKey("ShortURL"),
			model.WithIndexes(ownerIndex),
		),
		ownerIndex: ownerIndex,
	}
}

func (e *UrlShortener) Shorten(ctx context.Context, req *urlshortener.ShortenRequest, rsp *urlshortener.ShortenResponse) error {
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
	return e.pairs.Create(&urlshortener.URLPair{
		DestinationURL: req.DestinationURL,
		ShortURL:       id,
		Owner:          tenantID,
	})
}

func (e *UrlShortener) List(ctx context.Context, req *urlshortener.ListRequest, rsp *urlshortener.ListResponse) error {
	tenantID, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.New("Not authorized")
	}

	rsp.UrlPairs = []*urlshortener.URLPair{}
	err := e.pairs.Read(e.ownerIndex.ToQuery(e.ownerIndex.ToQuery(tenantID)), &rsp.UrlPairs)
	if err != nil {
		return err
	}
	for _, v := range rsp.UrlPairs {
		v.ShortURL = e.hostPrefix + "/" + v.ShortURL
	}
	return nil
}

func (e *UrlShortener) Get(ctx context.Context, req *urlshortener.GetRequest, rsp *urlshortener.GetResponse) error {
	tenantID, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.New("Not authorized")
	}

	var pair urlshortener.URLPair
	err := e.pairs.Read(e.ownerIndex.ToQuery(model.QueryEquals("ShortURL", e.hostPrefix+"/"+req.ShortURL)), pair)
	if err != nil {
		return err
	}
	if pair.Owner != tenantID {
		return errors.New("not authorized")
	}

	rsp.DestinationURL = pair.DestinationURL
	return nil
}
