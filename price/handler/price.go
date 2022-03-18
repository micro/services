package handler

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	"github.com/micro/services/price/crawler"
	pb "github.com/micro/services/price/proto"
)

type Price struct {
	Crawler *crawler.Crawler
}

func New() *Price {
	// TODO: look for "crypto.provider" to determine the handler
	v, err := config.Get("commodities.api_key")
	if err != nil {
		logger.Fatalf("commodities.api_key config not found: %v", err)
	}
	key := v.String("")
	if len(key) == 0 {
		logger.Fatal("commodities.api config not found")
	}

	c := crawler.New(key)
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

	// create a value
	value := &pb.Value{
		Name:      req.Name,
		Price:     req.Price,
		Currency:  req.Currency,
		Timestamp: timestamp.Format(time.RFC3339Nano),
	}

	// set response value
	rsp.Value = value

	// define a key
	key := path.Join(
		"price",
		strings.ToLower(value.Symbol),
		strings.ToLower(value.Currency),
		fmt.Sprintf("%d", timestamp.Unix()),
	)

	// TODO: add to index to search by name

	// create a record and store it
	rec := store.NewRecord(key, value)
	return store.Write(rec)
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

	// add currency if necessary
	if len(req.Currency) > 0 {
		key = path.Join(key, strings.ToLower(req.Currency))
	}

	key = path.Join("price", key)

	recs, err := store.Read(key, store.ReadPrefix(), store.ReadLimit(1), store.ReadOrder(store.OrderDesc))
	if err != nil && err != store.ErrNotFound {
		return err
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
