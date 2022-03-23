package handler

import (
	"context"
	"fmt"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	"github.com/micro/services/pkg/tenant"
	"github.com/micro/services/price/crawler"
	pb "github.com/micro/services/price/proto"
)

type Price struct {
	Crawler *crawler.Crawler
}

var (
	re = regexp.MustCompile("^[A-Z]{3}$")
)

func New() *Price {
	// TODO: look for "crypto.provider" to determine the handler
	v, err := config.Get("commodities.api_url")
	if err != nil {
		logger.Fatalf("commodities.api_url config not found: %v", err)
	}
	url := v.String("")
	// TODO: look for "crypto.provider" to determine the handler
	v, err = config.Get("commodities.api_key")
	if err != nil {
		logger.Fatalf("commodities.api_key config not found: %v", err)
	}
	key := v.String("")
	if len(key) == 0 {
		logger.Fatal("commodities.api config not found")
	}

	c := crawler.New(url, key)
	go c.Run()

	return &Price{c}
}

func (p *Price) Add(ctx context.Context, req *pb.AddRequest, rsp *pb.AddResponse) error {
	if len(req.Name) == 0 {
		return errors.BadRequest("price.add", "missing name")
	}
	if len(req.Symbol) == 0 {
		return errors.BadRequest("price.add", "missing symbol")
	}
	// price can technically be zero so no validation
	//if req.Price <= 0.0 {
	//	returns errors.BadRequest("price.add", "missing price")
	//}
	if len(req.Currency) == 0 {
		return errors.BadRequest("price.add", "missing currency")
	}

	timestamp := time.Now()

	symbol := strings.ToUpper(req.Symbol)
	if _, ok := crawler.Index[symbol]; ok {
		return errors.BadRequest("price.add", "already indexed")
	}

	// can't use our author name
	if len(req.Author) == 0 || req.Author == "Micro" {
		req.Author = "User"
	}

	if len(req.Source) == 0 {
		req.Source = "N/A"
	}

	// create a value
	value := &pb.Value{
		Name:      req.Name,
		Price:     req.Price,
		Currency:  req.Currency,
		Timestamp: timestamp.Format(time.RFC3339Nano),
		Source:    req.Source,
		Author:    req.Author,
	}

	// set response value
	rsp.Value = value

	for _, suffix := range []string{"latest", fmt.Sprintf("%d", timestamp.Unix())} {
		// define a key
		key := path.Join(
			"price",
			strings.ToLower(value.Symbol),
			strings.ToLower(value.Currency),
			suffix,
		)

		// TODO: add to index to search by name

		// create a record and store it
		rec := store.NewRecord(key, value)
		if err := store.Write(rec); err != nil {
			return err
		}
	}

	return nil
}

func (p *Price) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	// add currency if necessary
	if len(req.Currency) == 0 {
		req.Currency = "USD"
	}

	key := "/" + strings.ToLower(req.Currency) + "/latest"

	if req.Limit <= 0 {
		req.Limit = 100
	}

	if req.Offset < 0 {
		req.Offset = 0
	}

	offset := uint(req.Offset)
	limit := uint(req.Limit)

	recs, err := store.Read(key, store.ReadSuffix(), store.ReadOffset(offset), store.ReadLimit(limit))
	if err != nil && err != store.ErrNotFound {
		return err
	}

	if len(recs) == 0 {
		values, err := p.Crawler.List(
			strings.ToUpper(req.Currency),
		)
		if err != nil {
			return errors.NotFound("price.lis", "no values available")
		}
		rsp.Values = values
		return nil
	}

	for _, rec := range recs {
		value := new(pb.Value)
		if err := rec.Decode(value); err != nil {
			continue
		}
		if len(value.Name) == 0 {
			continue
		}

		rsp.Values = append(rsp.Values, value)
	}

	return nil
}

func (p *Price) Get(ctx context.Context, req *pb.GetRequest, rsp *pb.GetResponse) error {
	if len(req.Name) == 0 && len(req.Symbol) == 0 {
		return errors.BadRequest("price.get", "missing name or symbol")
	}

	var key string
	var symbol string

	if len(req.Name) > 0 {
		for k, v := range crawler.Index {
			name := strings.ToLower(v)
			if name == strings.ToLower(req.Name) {
				key = strings.ToLower(k)
				symbol = k
				break
			}

		}
	} else {
		// key is defined as the value name
		key = strings.ToLower(req.Symbol)
		symbol = strings.ToUpper(req.Symbol)
	}

	if len(key) == 0 {
		return errors.NotFound("price.get", "value not found")
	}

	var recs []*store.Record

	// hard code to USD if no currency
	// in future drop to allow listing
	if len(req.Currency) == 0 {
		req.Currency = "USD"
	}

	// add currency if necessary
	if len(req.Currency) > 0 {
		key = path.Join("price", key, strings.ToLower(req.Currency), "latest")

		r, err := store.Read(key, store.ReadPrefix(), store.ReadOrder(store.OrderDesc))
		if err != nil && err != store.ErrNotFound {
			return err
		}
		recs = r
	} else {
		// get a list of keys with prefix price/symbol and suffix /latest
		keys, err := store.List(
			store.ListPrefix(path.Join("price", key)),
			store.ListSuffix("/latest"),
		)
		if err != nil && err != store.ErrNotFound {
			return err
		}
		for _, key := range keys {
			r, err := store.Read(key, store.ReadLimit(1))
			if err != nil {
				continue
			}
			recs = append(recs, r...)
		}
	}

	// try get it directly
	if len(recs) == 0 {
		value, err := p.Crawler.Get(
			strings.ToUpper(symbol),
			strings.ToUpper(req.Currency),
		)
		if err != nil {
			return errors.NotFound("price.get", "value not found")
		}
		rsp.Values = append(rsp.Values, value)
		return nil
	}

	for _, rec := range recs {
		// only get the latest valeus
		if !strings.HasSuffix(rec.Key, "/latest") {
			continue
		}
		value := new(pb.Value)
		rec.Decode(value)
		rsp.Values = append(rsp.Values, value)
	}

	return nil
}

func (p *Price) Index(ctx context.Context, req *pb.IndexRequest, rsp *pb.IndexResponse) error {
	index, err := store.List(store.ListPrefix("index/"))
	if err != nil {
		logger.Errorf("Failed to retrieve index: %v", err)
		for k, v := range crawler.Index {
			rsp.Index = append(rsp.Index, &pb.Index{
				Name:     v,
				Symbol:   k,
				Currency: "USD",
			})
		}
		return nil
	}

	for _, idx := range index {
		parts := strings.Split(idx, "/")
		symbol := strings.ToUpper(parts[1])
		currency := strings.ToUpper(parts[2])
		name := crawler.Index[symbol]

		rsp.Index = append(rsp.Index, &pb.Index{
			Name:     name,
			Symbol:   symbol,
			Currency: currency,
		})
	}

	return nil
}

func (p *Price) Report(ctx context.Context, req *pb.ReportRequest, rsp *pb.ReportResponse) error {
	if len(req.Name) == 0 {
		return errors.BadRequest("price.report", "missing name")
	}
	if len(req.Symbol) == 0 {
		return errors.BadRequest("price.report", "missing symbol")
	}
	if len(req.Comment) == 0 {
		return errors.BadRequest("price.report", "missing comment")
	}

	id, _ := tenant.FromContext(ctx)

	rec := store.NewRecord(
		path.Join("report", req.Symbol, fmt.Sprintf("%d", time.Now().UnixNano())),
		&pb.Report{Name: req.Name, Symbol: req.Symbol, Comment: req.Comment, Author: id},
	)

	return store.Write(rec)
}
