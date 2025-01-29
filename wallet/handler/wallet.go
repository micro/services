package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/micro/micro/v5/service"
	"github.com/micro/micro/v5/service/errors"
	log "github.com/micro/micro/v5/service/logger"
	"github.com/micro/micro/v5/service/store"
	"github.com/micro/services/pkg/redis"
	"github.com/micro/services/pkg/tenant"
	pb "github.com/micro/services/wallet/proto"
)

const (
	accountPrefix     = "account"
	counterPrefix     = "wallet"
	transactionPrefix = "transaction"
)

// Transaction represents a wallet transaction
type Transaction struct {
	ID         string
	Created    time.Time
	Amount     int64  // positive is credit, negative is debit
	Reference  string // reference description
	Visible    bool   // should this be visible to the customer? If false, it only displays to admins
	WalletID   string
	ActionedBy string // who made the transaction
	Metadata   map[string]string
}

type Wallet struct {
	c *redis.Counter
	// for wallet transfers
	mtx sync.Mutex
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
		Key:   fmt.Sprintf("%s/%s/%s/%s", transactionPrefix, userID, walletID, rec.ID),
		Value: trx,
	}); err != nil {
		return nil, err
	}
	return rec, nil
}

func NewHandler(svc *service.Service) *Wallet {
	return &Wallet{
		c: redis.NewCounter(counterPrefix),
	}
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

	// check the wallets exist
	for _, id := range []string{req.FromId, req.ToId} {
		_, err := store.Read(fmt.Sprintf("%s/%s/%s", accountPrefix, tnt, id), store.ReadLimit(1))
		if err != nil {
			return errors.BadRequest("wallet.transfer", "invalid account")
		}
	}

	amount, err := b.c.Read(ctx, redis.Key(tnt, req.FromId), "$balance$")
	if amount < req.Amount {
		return errors.BadRequest("wallet.transfer", "insufficient credit")
	}

	_, err = b.c.Decr(ctx, redis.Key(tnt, req.FromId), "$balance$", req.Amount)
	if err != nil {
		return err
	}

	_, err = storeTransaction(tnt, -req.Amount, req.FromId, req.Reference, req.Visible, nil)
	if err != nil {
		return err
	}

	_, err = b.c.Incr(ctx, redis.Key(tnt, req.ToId), "$balance$", req.Amount)
	if err != nil {
		return err
	}

	_, err = storeTransaction(tnt, req.Amount, req.ToId, req.Reference, req.Visible, nil)
	if err != nil {
		return err
	}

	return nil
}

func (b *Wallet) Credit(ctx context.Context, request *pb.CreditRequest, response *pb.CreditResponse) error {
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.BadRequest("wallet.credit", "unauthorized")
	}

	if len(request.Id) == 0 {
		request.Id = "default"
	}

	if len(request.Reference) == 0 {
		return errors.BadRequest("wallet.credit", "Missing reference")
	}

	// TODO idempotency
	// increment the wallet
	currBal, err := b.c.Incr(ctx, redis.Key(tnt, request.Id), "$balance$", request.Amount)
	if err != nil {
		return err
	}

	response.Balance = currBal
	_, err = storeTransaction(tnt, request.Amount, request.Id, request.Reference, request.Visible, nil)
	if err != nil {
		return err
	}

	return nil
}

func (b *Wallet) Debit(ctx context.Context, request *pb.DebitRequest, response *pb.DebitResponse) error {
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.BadRequest("wallet.debit", "unauthorized")
	}

	if len(request.Reference) == 0 {
		return errors.BadRequest("wallet.debit", "Missing reference")
	}

	if len(request.Id) == 0 {
		request.Id = "default"
	}

	// TODO idempotency
	// decrement the wallet
	currBal, err := b.c.Decr(ctx, redis.Key(tnt, request.Id), "$balance$", request.Amount)
	if err != nil {
		return err
	}

	response.Balance = currBal

	_, err = storeTransaction(tnt, -request.Amount, request.Id, request.Reference, request.Visible, nil)
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

	if len(request.Id) == 0 {
		request.Id = "default"
	}

	currBal, err := b.c.Read(ctx, redis.Key(tnt, request.Id), "$balance$")
	if err != nil && err != redis.Nil {
		log.Errorf("Error reading from counter %s", err)
		return errors.InternalServerError("wallet.Balance", "Error retrieving balance")
	}

	response.Balance = currBal
	return nil
}

func (b *Wallet) Transactions(ctx context.Context, request *pb.TransactionsRequest, response *pb.TransactionsResponse) error {
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.BadRequest("wallet.transactions", "unauthorized")
	}

	if len(request.Id) == 0 {
		request.Id = "default"
	}

	recs, err := store.Read(fmt.Sprintf("%s/%s/%s/", transactionPrefix, tnt, request.Id), store.ReadPrefix())
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
			Created:   trx.Created.Format(time.RFC3339Nano),
			Amount:    trx.Amount,
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

func (b *Wallet) Create(ctx context.Context, request *pb.CreateRequest, response *pb.CreateResponse) error {
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.BadRequest("wallet.create", "unauthorized")
	}

	// generate a new id
	id := request.Id
	if len(id) == 0 {
		id = uuid.New().String()
	}

	// create a composite key
	key := fmt.Sprintf("%s/%s/%s", accountPrefix, tnt, id)

	acc := &pb.Account{
		Id:          id,
		Name:        request.Name,
		Description: request.Description,
	}

	// create a new record
	rec := store.NewRecord(key, acc)
	// store it
	if err := store.Write(rec); err != nil {
		return err
	}

	response.Account = acc

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

	userID := tnt
	walletID := request.Id

	// delete the account
	key := fmt.Sprintf("%s/%s/%s", accountPrefix, userID, walletID)
	if err := store.Delete(key); err != nil {
		return err
	}

	// delete the wallet
	if err := b.c.Delete(ctx, redis.Key(userID, walletID)); err != nil {
		return err
	}

	// delete all related transactions
	recs, err := store.List(store.ListPrefix(fmt.Sprintf("%s/%s/%s/", transactionPrefix, userID, walletID)))
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

func (w *Wallet) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.BadRequest("wallet.read", "unauthorized")
	}

	if len(req.Id) == 0 {
		req.Id = "default"
	}

	recs, err := store.Read(fmt.Sprintf("%s/%s/%s", accountPrefix, tnt, req.Id), store.ReadLimit(1))
	if err != nil {
		return err
	}
	if len(recs) == 0 {
		return nil
	}

	acc := new(pb.Account)
	recs[0].Decode(&acc)

	bal, err := w.c.Read(ctx, redis.Key(tnt, acc.Id), "$balance$")
	if err != nil && err != redis.Nil {
		log.Errorf("Error reading from counter %s", err)
		return errors.BadRequest("wallet.read", "error reading balance")
	}

	// set balance
	acc.Balance = bal

	rsp.Account = acc

	return nil
}

func (w *Wallet) List(ctx context.Context, req *pb.ListRequest, rsp *pb.ListResponse) error {
	tnt, ok := tenant.FromContext(ctx)
	if !ok {
		return errors.BadRequest("wallet.create", "unauthorized")
	}

	recs, err := store.Read(fmt.Sprintf("%s/%s/", accountPrefix, tnt), store.ReadPrefix())
	if err != nil {
		return err
	}

	var def bool

	for _, rec := range recs {
		acc := new(pb.Account)
		rec.Decode(&acc)

		bal, err := w.c.Read(ctx, redis.Key(tnt, acc.Id), "$balance$")
		if err != nil && err != redis.Nil {
			log.Errorf("Error reading from counter %s", err)
			continue
		}

		// set balance
		acc.Balance = bal
		rsp.Accounts = append(rsp.Accounts, acc)
	}

	// someone has an account with id default
	if def {
		return nil
	}

	bal, err := w.c.Read(ctx, redis.Key(tnt, "default"), "$balance$")
	if err != nil && err != redis.Nil {
		log.Errorf("Error reading from counter %s", err)
		return nil
	}

	// add default
	rsp.Accounts = append(rsp.Accounts, &pb.Account{
		Id:      "default",
		Name:    "Default account",
		Balance: bal,
	})

	return nil
}
