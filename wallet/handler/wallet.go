package handler

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	"github.com/micro/services/pkg/tenant"
	pb "github.com/micro/services/wallet/proto"
)

const (
	prefixCounter         = "wallet-service/counter"
	prefixStoreByCustomer = "transactionByUser"
)

type counter struct {
	sync.RWMutex
	redisClient *redis.Client
}

func (c *counter) incr(ctx context.Context, userID, walletID, path string, delta int64) (int64, error) {
	return c.redisClient.IncrBy(ctx, fmt.Sprintf("%s:%s:%s:%s", prefixCounter, userID, walletID, path), delta).Result()
}

func (c *counter) decr(ctx context.Context, userID, walletID, path string, delta int64) (int64, error) {
	return c.redisClient.DecrBy(ctx, fmt.Sprintf("%s:%s:%s:%s", prefixCounter, userID, walletID, path), delta).Result()
}

func (c *counter) read(ctx context.Context, userID, walletID, path string) (int64, error) {
	ret, err := c.redisClient.Get(ctx, fmt.Sprintf("%s:%s:%s:%s", prefixCounter, userID, walletID, path)).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return ret, err
}

func (c *counter) reset(ctx context.Context, userID, walletID, path string) error {
	return c.redisClient.Set(ctx, fmt.Sprintf("%s:%s:%s:%s", prefixCounter, userID, walletID, path), 0, 0).Err()
}

func (c *counter) deleteWallet(ctx context.Context, userID, walletID string) error {
	keys, err := c.redisClient.Keys(ctx, fmt.Sprintf("%s:%s:%s:*", prefixCounter, userID, walletID)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	if err := c.redisClient.Del(ctx, keys...).Err(); err != nil && err != redis.Nil {
		return err
	}

	return nil
}

// Transaction represents a wallet adjustment (not including normal API usage). e.g. credit being added, promo codes, manual adjustment for customer service etc
type Transaction struct {
	ID         string
	Created    time.Time
	Amount     int64  // positive is credit, negative is debit
	Reference  string // reference description
	Visible    bool   // should this be visible to the customer? If false, it only displays to admins
	WalletID   string
	ActionedBy string // who made the adjustment
	Metadata   map[string]string
}

type Wallet struct {
	c *counter // counts the wallet. Wallet is expressed in 1/10,000ths of a cent which allows us to price in fractions e.g. a request costs 0.0001 cents or 10,000 requests for 1 cent
	// for wallet transfers
	mtx sync.Mutex
}

func NewHandler(svc *service.Service) *Wallet {
	redisConfig := struct {
		Address  string
		User     string
		Password string
	}{}
	val, err := config.Get("micro.redis")
	if err != nil {
		log.Fatalf("No redis config found %s", err)
	}
	if err := val.Scan(&redisConfig); err != nil {
		log.Fatalf("Error parsing redis config %s", err)
	}
	if len(redisConfig.Password) == 0 || len(redisConfig.User) == 0 || len(redisConfig.Password) == 0 {
		log.Fatalf("Missing redis config %s", err)
	}
	rc := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Address,
		Username: redisConfig.User,
		Password: redisConfig.Password,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	})
	b := &Wallet{
		c: &counter{redisClient: rc},
	}
	return b
}

func (b *Wallet) Transfer(ctx context.Context, req *pb.TransferRequest, rsp *pb.TransferResponse) error {
	b.mtx.Lock()
	defer b.mtx.Unlock()

	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.BadRequest("wallet.transfer", "unauthorized")
	}

	if len(req.FromId) == 0 || len(req.ToId) == 0 {
		return errors.BadRequest("wallet.transfer", "missing ids")
	}

	amount, err := b.c.read(ctx, tnt, req.FromId, "$wallet$")
	if amount > req.Amount {
		return errors.BadRequest("wallet.transfer", "insufficient wallet")
	}

	_, err = b.c.decr(ctx, tnt, req.FromId, "$wallet$", req.Amount)
	if err != nil {
		return err
	}

	_, err = storeTransaction(tnt, -req.Amount, req.FromId, req.Reference, req.Visible, nil)
	if err != nil {
		return err
	}

	_, err = b.c.incr(ctx, tnt, req.ToId, "$wallet$", req.Amount)
	if err != nil {
		return err
	}

	_, err = storeTransaction(tnt, req.Amount, req.ToId, req.Reference, req.Visible, nil)
	if err != nil {
		return err
	}

	return nil
}

func (b Wallet) Credit(ctx context.Context, request *pb.CreditRequest, response *pb.CreditResponse) error {
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.BadRequest("wallet.credit", "unauthorized")
	}

	if len(request.Reference) == 0 {
		return errors.BadRequest("wallet.credit", "Missing reference")
	}

	// TODO idempotency
	// increment the wallet
	currBal, err := b.c.incr(ctx, tnt, request.Id, "$wallet$", request.Delta)
	if err != nil {
		return err
	}

	response.NewBalance = currBal
	_, err = storeTransaction(tnt, request.Delta, request.Id, request.Reference, request.Visible, nil)
	if err != nil {
		return err
	}

	return nil
}

func storeTransaction(userID string, delta int64, walletID, reference string, visible bool, meta map[string]string) (*Transaction, error) {
	// record it
	rec := &Transaction{
		ID:         uuid.New().String(),
		Created:    time.Now(),
		Amount:     delta,
		Reference:  reference,
		Visible:    visible,
		WalletID:   walletID,
		ActionedBy: userID,
		Metadata:   meta,
	}

	trx, err := json.Marshal(rec)
	if err != nil {
		return nil, err
	}

	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf("%s/%s/%s/%s", prefixStoreByCustomer, userID, walletID, rec.ID),
		Value: trx,
	}); err != nil {
		return nil, err
	}
	return rec, nil
}

func (b *Wallet) Debit(ctx context.Context, request *pb.DebitRequest, response *pb.DebitResponse) error {
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.BadRequest("wallet.debit", "unauthorized")
	}

	if len(request.Reference) == 0 {
		return errors.BadRequest("wallet.debit", "Missing reference")
	}

	// TODO idempotency
	// decrement the wallet
	currBal, err := b.c.decr(ctx, tnt, request.Id, "$wallet$", request.Delta)
	if err != nil {
		return err
	}

	response.NewBalance = currBal

	_, err = storeTransaction(tnt, -request.Delta, request.Id, request.Reference, request.Visible, nil)
	if err != nil {
		return err
	}

	return nil
}

func (b *Wallet) Balance(ctx context.Context, request *pb.BalanceRequest, response *pb.BalanceResponse) error {
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.BadRequest("wallet.balance", "unauthorized")
	}

	currBal, err := b.c.read(ctx, tnt, request.Id, "$wallet$")
	if err != nil && err != redis.Nil {
		log.Errorf("Error reading from counter %s", err)
		return errors.InternalServerError("wallet.Balance", "Error retrieving current wallet")
	}

	response.Balance = currBal
	return nil
}

func (b *Wallet) Transactions(ctx context.Context, request *pb.TransactionsRequest, response *pb.TransactionsResponse) error {
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.BadRequest("wallet.transactions", "unauthorized")
	}

	recs, err := store.Read(fmt.Sprintf("%s/%s/%s/", prefixStoreByCustomer, tnt, request.Id), store.ReadPrefix())
	if err != nil {
		return err
	}

	ret := []*pb.Transaction{}
	for _, rec := range recs {
		var trx Transaction
		if err := json.Unmarshal(rec.Value, &trx); err != nil {
			return err
		}
		if !trx.Visible {
			continue
		}
		ret = append(ret, &pb.Transaction{
			Id:        trx.ID,
			Created:   trx.Created.Unix(),
			Delta:     trx.Amount,
			Reference: trx.Reference,
			Metadata:  trx.Metadata,
		})
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Created < ret[j].Created
	})

	response.Transactions = ret

	return nil
}

func (b *Wallet) deleteAccount(ctx context.Context, userID, walletID string) error {
	// delete the account
	key := fmt.Sprintf("wallet/%s/%s", userID, walletID)
	if err := store.Delete(key); err != nil {
		return err
	}

	if err := b.c.deleteWallet(ctx, userID, walletID); err != nil {
		return err
	}

	recs, err := store.List(store.ListPrefix(fmt.Sprintf("%s/%s/%s/", prefixStoreByCustomer, userID, walletID)))
	if err != nil {
		return err
	}

	for _, rec := range recs {
		if err := store.Delete(rec); err != nil {
			return err
		}
	}

	return nil
}

func (b *Wallet) Create(ctx context.Context, request *pb.CreateRequest, response *pb.CreateResponse) error {
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.BadRequest("wallet.create", "unauthorized")
	}

	// generate a new id
	id := uuid.New().String()

	// create a composite key
	key := fmt.Sprintf("wallet/%s/%s", tnt, id)

	// create a new record
	rec := store.NewRecord(key, &pb.Account{
		Id:          id,
		Name:        request.Name,
		Description: request.Description,
	})

	// store it
	if err := store.Write(rec); err != nil {
		return err
	}

	return nil
}

func (b *Wallet) Delete(ctx context.Context, request *pb.DeleteRequest, response *pb.DeleteResponse) error {
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.BadRequest("wallet.delete", "unauthorized")
	}

	if len(request.Id) == 0 {
		return errors.BadRequest("wallet.delete", "Missing wallet id")
	}

	if err := b.deleteAccount(ctx, tnt, request.Id); err != nil {
		logger.Errorf("Error deleting customer %s", err)
		return err
	}
	return nil
}

func (w *Wallet) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.BadRequest("wallet.create", "unauthorized")
	}

	recs, err := store.Read(fmt.Sprintf("wallet/%s/", tnt), store.ReadPrefix())
	if err != nil {
		return err
	}

	for _, rec := range recs {
		acc := new(pb.Account)
		rec.Decode(&acc)
		rsp.Accounts = append(rsp.Accounts, acc)
	}

	return nil
}
