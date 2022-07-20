package crawler

import (
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	"github.com/micro/services/pkg/api"
	pb "github.com/micro/services/price/proto"
)

type Crawler struct {
	Url string
	Key string
}

type Response struct {
	Data Data `json:"data"`
}

type Data struct {
	Success   bool               `json:"success"`
	Timestamp int64              `json:"timestamp"`
	Date      string             `json:"date"`
	Base      string             `json:"base"`
	Rates     map[string]float64 `json:"rates"`
	Unit      string             `json:"unit,omitempty"`
}

func (c *Crawler) GetPrices(base string) {
	uri := c.Url + "/latest"
	vals := url.Values{}
	vals.Set("access_key", c.Key)
	vals.Set("base", base)

	var symbols []string

	for symbol, _ := range Index {
		symbols = append(symbols, symbol)
	}

	vals.Set("symbols", strings.Join(symbols, ","))

	q := vals.Encode()
	uri += "?" + q

	var rsp Response

	if err := api.Get(uri, &rsp); err != nil {
		logger.Errorf("Failed to get symbols: %v", err)
		return
	}

	for symbol, value := range rsp.Data.Rates {
		name := Index[symbol]

		val := &pb.Value{
			Name:      name,
			Price:     float64(1) / value,
			Symbol:    symbol,
			Currency:  rsp.Data.Base,
			Timestamp: time.Unix(rsp.Data.Timestamp, 0).Format(time.RFC3339Nano),
			Source:    c.Url,
			Author:    "Micro",
		}

		for _, suffix := range []string{"latest", fmt.Sprintf("%d", rsp.Data.Timestamp)} {
			// store it
			key := path.Join(
				"price",
				strings.ToLower(symbol),
				strings.ToLower(rsp.Data.Base),
				suffix,
			)

			rec := store.NewRecord(key, val)

			// save the record
			if err := store.Write(rec); err != nil {
				logger.Error("Failed to write symbol: %v error: %v", key, err)
			}
		}

		// index the item for the future
		key := path.Join(
			"index",
			strings.ToLower(symbol),
			strings.ToLower(rsp.Data.Base),
		)

		if err := store.Write(store.NewRecord(key, &pb.Index{
			Name:     val.Name,
			Symbol:   val.Symbol,
			Currency: val.Currency,
		})); err != nil {
			logger.Error("Failed to write index: %v error: %v", key, err)
		}
	}
}

func (c *Crawler) Run() {
	t := time.NewTicker(time.Minute * 10)
	defer t.Stop()

	// build the index
	var index map[string]interface{}

	vals := url.Values{}
	vals.Set("access_key", c.Key)
	q := vals.Encode()
	uri := c.Url + "/symbols?" + q

	if err := api.Get(uri, &index); err != nil {
		logger.Errorf("Failed to get index symbols: %v", err)
	}

	// update our built in index
	for k, v := range index {
		Index[k] = v.(string)
	}

	symbols := make(map[string]bool)

	// load the indexed items
	recs, err := store.List(store.ListPrefix("index/"))
	if err != nil {
		logger.Errorf("Failed to read index: %v", err)
	}

	// pull out the currencies for them
	for _, rec := range recs {
		parts := strings.Split(rec, "/")
		symbols[strings.ToUpper(parts[2])] = true
	}

	// if no currencies then use USD
	if len(symbols) == 0 {
		symbols["USD"] = true
	}

	// get prices now
	for symbol, _ := range symbols {
		c.GetPrices(symbol)
	}

	for {
		select {
		case <-t.C:
			// get prices now
			for symbol, _ := range symbols {
				c.GetPrices(symbol)
			}
		}
	}
}

func (c *Crawler) Get(symbol, currency string) (*pb.Value, error) {
	uri := c.Url + "/latest"
	vals := url.Values{}
	vals.Set("access_key", c.Key)
	vals.Set("base", currency)
	vals.Set("symbols", symbol)

	q := vals.Encode()
	uri += "?" + q

	var rsp Response

	if err := api.Get(uri, &rsp); err != nil {
		logger.Errorf("Failed to get symbols for %v:%v: %v", symbol, currency, err)
		return nil, err
	}

	if _, ok := rsp.Data.Rates[symbol]; !ok {
		return nil, errors.NotFound("crawler.get", "rate not found")
	}

	val := &pb.Value{
		Name:      Index[symbol],
		Price:     float64(1) / rsp.Data.Rates[symbol],
		Symbol:    symbol,
		Currency:  rsp.Data.Base,
		Timestamp: time.Unix(rsp.Data.Timestamp, 0).Format(time.RFC3339Nano),
		Source:    c.Url,
		Author:    "Micro",
	}

	// write historic record and latest
	for _, suffix := range []string{"latest", fmt.Sprintf("%d", rsp.Data.Timestamp)} {
		key := path.Join(
			"price",
			strings.ToLower(symbol),
			strings.ToLower(rsp.Data.Base),
			suffix,
		)

		rec := store.NewRecord(key, val)

		// save the record
		if err := store.Write(rec); err != nil {
			logger.Error("Failed to write symbol: %v error: %v", symbol, err)
		}
	}

	// index the item for the future
	key := path.Join(
		"index",
		strings.ToLower(symbol),
		strings.ToLower(rsp.Data.Base),
	)

	if err := store.Write(store.NewRecord(key, &pb.Index{
		Name:     val.Name,
		Symbol:   val.Symbol,
		Currency: val.Currency,
	})); err != nil {
		logger.Error("Failed to write index: %v error: %v", key, err)
	}

	// return value
	return val, nil
}

func (c *Crawler) List(currency string) ([]*pb.Value, error) {
	uri := c.Url + "/latest"
	vals := url.Values{}
	vals.Set("access_key", c.Key)
	vals.Set("base", currency)

	q := vals.Encode()
	uri += "?" + q

	var rsp Response

	if err := api.Get(uri, &rsp); err != nil {
		logger.Errorf("Failed to get list for currency %v: %v", currency, err)
		return nil, err
	}

	var values []*pb.Value

	for symbol, rate := range rsp.Data.Rates {
		val := &pb.Value{
			Name:      Index[symbol],
			Price:     float64(1) / rate,
			Symbol:    symbol,
			Currency:  rsp.Data.Base,
			Timestamp: time.Unix(rsp.Data.Timestamp, 0).Format(time.RFC3339Nano),
			Source:    c.Url,
			Author:    "Micro",
		}

		values = append(values, val)

		// write historic record and latest
		for _, suffix := range []string{"latest", fmt.Sprintf("%d", rsp.Data.Timestamp)} {
			key := path.Join(
				"price",
				strings.ToLower(symbol),
				strings.ToLower(rsp.Data.Base),
				suffix,
			)

			rec := store.NewRecord(key, val)

			// save the record
			if err := store.Write(rec); err != nil {
				logger.Error("Failed to write symbol: %v error: %v", symbol, err)
			}
		}

		// index the item for the future
		key := path.Join(
			"index",
			strings.ToLower(symbol),
			strings.ToLower(rsp.Data.Base),
		)

		if err := store.Write(store.NewRecord(key, &pb.Index{
			Name:     val.Name,
			Symbol:   val.Symbol,
			Currency: val.Currency,
		})); err != nil {
			logger.Error("Failed to write index: %v error: %v", key, err)
		}
	}

	// return value
	return values, nil
}

func New(url, key string) *Crawler {
	return &Crawler{Url: url, Key: key}
}
