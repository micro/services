package handler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/micro/micro/v5/service/errors"
	"github.com/micro/micro/v5/service/logger"
	pb "github.com/micro/services/bitcoin/proto"
	"github.com/micro/services/pkg/api"
	"github.com/patrickmn/go-cache"
)

type Bitcoin struct {
	Cache *cache.Cache
}

func New() *Bitcoin {
	api.SetCache(true, time.Minute*5)

	return &Bitcoin{
		Cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (b *Bitcoin) Balance(ctx context.Context, req *pb.BalanceRequest, rsp *pb.BalanceResponse) error {
	if len(req.Address) == 0 {
		return errors.BadRequest("bitcoin.balance", "missing address")
	}

	uri := fmt.Sprintf("https://blockchain.info/balance?active=%s", req.Address)

	var respBody map[string]interface{}

	if err := api.Get(uri, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal balance: %v\n", err)
		return errors.InternalServerError("bitcoin.balance", "failed to get price")
	}

	info := respBody[req.Address].(map[string]interface{})
	rsp.Balance = int64(info["final_balance"].(float64))

	return nil
}

func (b *Bitcoin) Lookup(ctx context.Context, req *pb.LookupRequest, rsp *pb.LookupResponse) error {
	if len(req.Address) == 0 {
		return errors.BadRequest("bitcoin.lookup", "missing address")
	}

	uri := fmt.Sprintf("https://blockchain.info/rawaddr/%s", req.Address)

	if req.Limit <= 0 || req.Limit > 50 {
		req.Limit = 50
	}

	uri += fmt.Sprintf("?limit=%d&offset=%d", req.Limit, req.Offset)

	var respBody map[string]interface{}

	if err := api.Get(uri, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal address: %v\n", err)
		return errors.InternalServerError("bitcoin.lookup", "failed to get address")
	}

	rsp.Address = req.Address
	rsp.Hash = respBody["hash160"].(string)
	rsp.NumTx = int64(respBody["n_tx"].(float64))
	rsp.NumUnredeemed = int64(respBody["n_unredeemed"].(float64))
	rsp.TotalReceived = int64(respBody["total_received"].(float64))
	rsp.TotalSent = int64(respBody["total_sent"].(float64))
	rsp.FinalBalance = int64(respBody["final_balance"].(float64))

	for _, tx := range respBody["txs"].([]interface{}) {
		transaction := tx.(map[string]interface{})

		// result of transaction
		rspTx := new(pb.Transaction)

		rspTx.Result = int64(transaction["result"].(float64))
		rspTx.Balance = int64(transaction["balance"].(float64))
		rspTx.Version = int64(transaction["ver"].(float64))
		rspTx.VinSz = int64(transaction["vin_sz"].(float64))
		rspTx.VoutSz = int64(transaction["vout_sz"].(float64))
		rspTx.Size = int64(transaction["size"].(float64))
		rspTx.Weight = int64(transaction["weight"].(float64))
		rspTx.Fee = int64(transaction["fee"].(float64))
		rspTx.Relay = transaction["relayed_by"].(string)
		rspTx.LockTime = int64(transaction["lock_time"].(float64))
		rspTx.TxIndex = int64(transaction["tx_index"].(float64))
		rspTx.DoubleSpend = transaction["double_spend"].(bool)
		rspTx.BlockIndex = int64(transaction["block_index"].(float64))
		rspTx.BlockHeight = int64(transaction["block_height"].(float64))

		inputs := transaction["inputs"].([]interface{})
		outputs := transaction["out"].([]interface{})

		for _, input := range inputs {
			in := input.(map[string]interface{})

			prev := in["prev_out"].(map[string]interface{})

			script, _ := prev["script"].(string)
			addr, _ := prev["addr"].(string)
			spent, _ := prev["spent"].(bool)
			n, _ := prev["n"].(float64)
			txIndex, _ := prev["tx_index"].(float64)
			value, _ := prev["value"].(float64)

			rspTx.Inputs = append(rspTx.Inputs, &pb.Input{
				Script: in["script"].(string),
				PrevOut: &pb.Prev{
					Value:   int64(value),
					Script:  script,
					Address: addr,
					Spent:   spent,
					TxIndex: int64(txIndex),
					N:       int64(n),
				},
			})
		}

		for _, output := range outputs {
			out := output.(map[string]interface{})

			rspTx.Outputs = append(rspTx.Outputs, &pb.Output{
				Value:   int64(out["value"].(float64)),
				Spent:   out["spent"].(bool),
				Script:  out["script"].(string),
				Address: out["addr"].(string),
				TxIndex: int64(out["tx_index"].(float64)),
			})
		}

		rsp.Transactions = append(rsp.Transactions, rspTx)
	}

	return nil
}

func (b *Bitcoin) Price(ctx context.Context, req *pb.PriceRequest, rsp *pb.PriceResponse) error {
	if len(req.Symbol) <= 0 {
		req.Symbol = "USD"
	}

	if strings.HasPrefix(req.Symbol, "BTC") {
		req.Symbol = strings.TrimPrefix(req.Symbol, "BTC")
	}

	// upper case
	req.Symbol = strings.ToUpper(req.Symbol)

	// try the cache first
	if price, ok := b.Cache.Get("price:" + req.Symbol); ok {
		rsp.Symbol = req.Symbol
		rsp.Price = price.(float64)
		return nil
	}

	// get the price

	uri := "https://blockchain.info/ticker"

	var respBody map[string]interface{}

	if err := api.Get(uri, &respBody); err != nil {
		logger.Errorf("Failed to unmarshal price: %v\n", err)
		return errors.InternalServerError("bitcoin.price", "failed to get price")
	}

	data, ok := respBody[req.Symbol]
	if !ok {
		return errors.InternalServerError("bitcoin.price", "unsupported symbol")
	}

	rsp.Symbol = req.Symbol
	rsp.Price = data.(map[string]interface{})["last"].(float64)

	// cache the price
	b.Cache.Set("price:"+req.Symbol, rsp.Price, time.Minute*5)

	return nil
}

func (b *Bitcoin) Transaction(ctx context.Context, req *pb.TransactionRequest, rsp *pb.TransactionResponse) error {
	if len(req.Hash) == 0 {
		return errors.BadRequest("bitcoin.transaction", "missing hash")
	}

	uri := fmt.Sprintf("https://blockchain.info/rawtx/%s", req.Hash)

	var respBody map[string]interface{}

	if err := api.Get(uri, &respBody); err != nil {
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
	outputs := respBody["out"].([]interface{})

	for _, input := range inputs {
		in := input.(map[string]interface{})

		prev := in["prev_out"].(map[string]interface{})

		rsp.Inputs = append(rsp.Inputs, &pb.Input{
			Script: in["script"].(string),
			PrevOut: &pb.Prev{
				Value:   int64(prev["value"].(float64)),
				Script:  prev["script"].(string),
				Address: prev["addr"].(string),
				Spent:   prev["spent"].(bool),
				TxIndex: int64(prev["tx_index"].(float64)),
				N:       int64(prev["n"].(float64)),
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
