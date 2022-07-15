package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/bitcoin/proto"
	"github.com/patrickmn/go-cache"
)

type Bitcoin struct {
	Api   string
	Key   string
	Cache *cache.Cache
}

func New() *Bitcoin {
	// TODO: look for "bitcoin.provider" to determine the handler
	v, err := config.Get("finage.api")
	if err != nil {
		logger.Fatalf("finage.api config not found: %v", err)
	}
	api := v.String("")
	if len(api) == 0 {
		logger.Fatal("finage.api config not found")
	}
	v, err = config.Get("finage.key")
	if err != nil {
		logger.Fatalf("finage.key config not found: %v", err)
	}
	key := v.String("")
	if len(key) == 0 {
		logger.Fatal("finage.key config not found")
	}

	return &Bitcoin{
		Api:   api,
		Key:   key,
		Cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (b *Bitcoin) Balance(ctx context.Context, req *pb.BalanceRequest, rsp *pb.BalanceResponse) error {
	if len(req.Address) == 0 {
		return errors.BadRequest("bitcoin.balance", "missing address")
	}

	uri := fmt.Sprintf("https://blockchain.info/balance?active=%s", req.Address)

	resp, err := http.Get(uri)
	if err != nil {
		logger.Errorf("Failed to get balance: %v\n", err)
		return errors.InternalServerError("bitcoin.balance", "failed to get price")
	}
	defer resp.Body.Close()

	buf, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to get price (non 200): %d %v\n", resp.StatusCode, string(buf))
		return errors.InternalServerError("bitcoin.balance", "failed to get price")
	}

	var respBody map[string]interface{}

	if err := json.Unmarshal(buf, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal balance: %v\n", err)
		return errors.InternalServerError("bitcoin.balance", "failed to get price")
	}

	info := respBody[req.Address].(map[string]interface{})
	rsp.Balance = int64(info["final_balance"].(float64))

	return nil
}

func (b *Bitcoin) Price(ctx context.Context, req *pb.PriceRequest, rsp *pb.PriceResponse) error {
	if len(req.Symbol) <= 0 {
		req.Symbol = "BTCUSD"
	}

	if !strings.HasPrefix(req.Symbol, "BTC") {
		return errors.BadRequest("bitcoin.price", "Must be of format BTCXXX e.g BTCUSD")
	}

	// try the cache first
	if price, ok := b.Cache.Get("price:" + req.Symbol); ok {
		rsp.Symbol = req.Symbol
		rsp.Price = price.(float64)
		return nil
	}

	// get the price
	uri := fmt.Sprintf("%slast/crypto/%s?apikey=%s", b.Api, req.Symbol, b.Key)

	resp, err := http.Get(uri)
	if err != nil {
		logger.Errorf("Failed to get price: %v\n", err)
		return errors.InternalServerError("bitcoin.price", "failed to get price")
	}
	defer resp.Body.Close()

	buf, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to get price (non 200): %d %v\n", resp.StatusCode, string(buf))
		return errors.InternalServerError("bitcoin.price", "failed to get price")
	}

	var respBody map[string]interface{}

	if err := json.Unmarshal(buf, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal price: %v\n", err)
		return errors.InternalServerError("bitcoin.price", "failed to get price")
	}

	rsp.Symbol = req.Symbol
	rsp.Price = respBody["price"].(float64)

	// cache the price
	b.Cache.Set("price:"+req.Symbol, rsp.Price, time.Minute*5)

	return nil
}

func (b *Bitcoin) Transaction(ctx context.Context, req *pb.TransactionRequest, rsp *pb.TransactionResponse) error {
	if len(req.Hash) == 0 {
		return errors.BadRequest("bitcoin.transaction", "missing hash")
	}

	uri := fmt.Sprintf("https://blockchain.info/rawtx/%s", req.Hash)

	resp, err := http.Get(uri)
	if err != nil {
		logger.Errorf("Failed to get transaction: %v\n", err)
		return errors.InternalServerError("bitcoin.transaction", "failed to get transaction")
	}
	defer resp.Body.Close()

	buf, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Errorf("Failed to get transaction (non 200): %d %v\n", resp.StatusCode, string(buf))
		return errors.InternalServerError("bitcoin.transaction", "failed to get transaction")
	}

	var respBody map[string]interface{}

	if err := json.Unmarshal(buf, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal transaction: %v\n", err)
		return errors.InternalServerError("bitcoin.transaction", "failed to get transaction")
	}

	rsp.Hash = req.Hash
	rsp.Version = int64(respBody["ver"].(float64))
	rsp.VinSz = int64(respBody["vin_sz"].(float64))
	rsp.VoutSz = int64(respBody["vout_sz"].(float64))
	rsp.Size = int64(respBody["size"].(float64))
	rsp.Weight = int64(respBody["weight"].(float64))
	rsp.Fee = int64(respBody["fee"].(float64))
	rsp.Relay = respBody["relayed_by"].(string)
	rsp.LockTime = int64(respBody["lock_time"].(float64))
	rsp.TxIndex = int64(respBody["tx_index"].(float64))
	rsp.DoubleSpend = respBody["double_spend"].(bool)
	rsp.BlockIndex = int64(respBody["block_index"].(float64))
	rsp.BlockHeight = int64(respBody["block_height"].(float64))

	inputs := respBody["inputs"].([]interface{})
	outputs := respBody["outputs"].([]interface{})

	for _, input := range inputs {
		in := input.(map[string]interface{})

		prev := in["prev_out"].(map[string]interface{})

		rsp.Inputs = append(rsp.Inputs, &pb.Input{
			Script: in["script"].(string),
			PrevOut: &pb.Prev{
				Hash:    prev["hash"].(string),
				Value:   int64(prev["value"].(float64)),
				Script:  prev["script"].(string),
				Address: prev["address"].(string),
				Spent:   prev["spent"].(bool),
				TxIndex: int64(prev["tx_index"].(float64)),
				N:       prev["n"].(string),
			},
		})
	}

	for _, output := range outputs {
		out := output.(map[string]interface{})

		rsp.Outputs = append(rsp.Outputs, &pb.Output{
			Value:   int64(out["value"].(float64)),
			Spent:   out["spent"].(bool),
			Script:  out["script"].(string),
			Address: out["addr"].(string),
			TxIndex: int64(out["tx_index"].(float64)),
		})
	}

	return nil
}
