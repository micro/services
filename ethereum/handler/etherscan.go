package handler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	pb "github.com/micro/services/ethereum/proto"
	"github.com/micro/services/pkg/api"
	"micro.dev/v4/service/config"
	"micro.dev/v4/service/errors"
	"micro.dev/v4/service/logger"
)

// Etherscan handler
type Etherscan struct {
	apiKey string
}

var (
	etherscanURL = "https://api.etherscan.io/api"
)

func New() *Etherscan {
	v, err := config.Get("etherscan.api_key")
	if err != nil {
		logger.Fatal("etherscan.api_key config not found: %v", err)
	}

	key := v.String("")
	if len(key) == 0 {
		logger.Fatal("etherscan.api_key config not found")
	}

	// set the cache
	api.SetCache(true, time.Minute*5)

	return &Etherscan{
		apiKey: key,
	}
}

func (e *Etherscan) url(module, action string, args map[string]interface{}) string {
	uri := fmt.Sprintf("%s?apikey=%s", etherscanURL, e.apiKey)
	uri += "&module=" + module
	uri += "&action=" + action

	for k, v := range args {
		uri += fmt.Sprintf("&%s=%v", k, v)
	}

	return uri
}

func (e *Etherscan) Balance(ctx context.Context, req *pb.BalanceRequest, rsp *pb.BalanceResponse) error {
	if len(req.Address) == 0 {
		return errors.BadRequest("ethereum.balance", "missing address")
	}

	uri := e.url("account", "balance", map[string]interface{}{
		"tag":     "latest",
		"address": req.Address,
	})

	var resp map[string]interface{}
	if err := api.Get(uri, &resp); err != nil {
		return errors.InternalServerError("ethereum.balance", err.Error())
	}

	if v, ok := resp["message"]; !ok || v.(string) != "OK" {
		logger.Error("Failed to get balance ", v.(string))
		return errors.InternalServerError("ethereum.balance", "failed to get balance")
	}

	bal := resp["result"].(string)

	rsp.Balance, _ = strconv.ParseInt(bal, 10, 64)

	return nil
}

func (e *Etherscan) Broadcast(ctx context.Context, req *pb.BroadcastRequest, rsp *pb.BroadcastResponse) error {
	if len(req.Hex) == 0 {
		return errors.BadRequest("ethereum.broadcast", "missing hex")
	}

	uri := e.url("proxy", "eth_sendRawTransaction", map[string]interface{}{
		"hex": req.Hex,
	})

	var resp map[string]interface{}
	if err := api.Get(uri, &resp); err != nil {
		return errors.InternalServerError("ethereum.broadcast", err.Error())
	}

	if v, ok := resp["id"]; !ok || v.(float64) != 1.00 {
		logger.Errorf("Failed to broadcast %v", resp)
		return errors.InternalServerError("ethereum.broadcast", "failed to broadcast transaction")
	}

	rsp.Hash, _ = resp["result"].(string)
	return nil
}

func (e *Etherscan) Transaction(ctx context.Context, req *pb.TransactionRequest, rsp *pb.TransactionResponse) error {
	if len(req.Hash) == 0 {
		return errors.BadRequest("ethereum.transaction", "missing hash")
	}

	uri := e.url("proxy", "eth_getTransactionByHash", map[string]interface{}{
		"txhash": req.Hash,
	})

	var resp map[string]interface{}
	if err := api.Get(uri, &resp); err != nil {
		return errors.InternalServerError("ethereum.transaction", err.Error())
	}

	if v, ok := resp["id"]; !ok || v.(float64) != 1.00 {
		logger.Errorf("Failed to get transaction %v", resp)
		return errors.InternalServerError("ethereum.transaction", "failed to get transaction")
	}

	tx := resp["result"].(map[string]interface{})

	rsp.Hash = tx["hash"].(string)
	rsp.BlockHash = tx["blockHash"].(string)
	rsp.BlockNumber = tx["blockNumber"].(string)
	rsp.FromAddress = tx["from"].(string)
	rsp.Gas = tx["gas"].(string)
	rsp.GasPrice = tx["gasPrice"].(string)
	rsp.MaxFeePerGas = tx["maxFeePerGas"].(string)
	rsp.MaxPriorityFeePerGas = tx["maxPriorityFeePerGas"].(string)
	rsp.Input = tx["input"].(string)
	rsp.Nonce = tx["nonce"].(string)
	rsp.ToAddress = tx["to"].(string)
	rsp.TxIndex = tx["transactionIndex"].(string)
	rsp.Value = tx["value"].(string)
	rsp.Type = tx["type"].(string)
	rsp.ChainId = tx["chainId"].(string)
	rsp.V = tx["v"].(string)
	rsp.R = tx["r"].(string)
	rsp.S = tx["s"].(string)

	return nil
}
