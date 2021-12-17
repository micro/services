package migrate

import (
	goctx "context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	"github.com/pkg/errors"

	authPb "github.com/micro/micro/v3/proto/auth"

	db "github.com/micro/services/db/proto"
	"github.com/micro/services/user/domain"
	user "github.com/micro/services/user/proto"
)

func generateAccountStoreKey(userId string) string {
	return fmt.Sprintf("user/account/id/%s", userId)
}

func generateAccountEmailStoreKey(email string) string {
	return fmt.Sprintf("user/acccount/email/%s", email)
}

func generateAccountUsernameStoreKey(username string) string {
	return fmt.Sprintf("user/account/username/%s", username)
}

type userMigration struct {
	from        db.DbService
	to          store.Store
	domain      *domain.Domain
	authAccount authPb.AccountsService
}

func NewUserMigration(from db.DbService, to store.Store, authAccount authPb.AccountsService) *userMigration {
	return &userMigration{
		from:        from,
		to:          to,
		domain:      domain.New(to),
		authAccount: authAccount,
	}
}

func (u *userMigration) List(o, l int32, account *auth.Account) ([]*user.Account, error) {
	var limit int32 = 25
	var offset int32 = 0

	if l > 0 {
		limit = l
	}
	if o > 0 {
		offset = o
	}

	ctx := auth.ContextWithAccount(goctx.Background(), account)
	rsp, err := u.from.Read(ctx, &db.ReadRequest{
		Table:  "users",
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		return nil, err
	}

	if len(rsp.Records) == 0 {
		return nil, errors.New("no user records found")
	}

	ret := make([]*user.Account, len(rsp.Records))
	for i, v := range rsp.Records {
		m, _ := v.MarshalJSON()
		var acc user.Account
		json.Unmarshal(m, &acc)
		ret[i] = &acc
	}

	return ret, nil
}

func (u *userMigration) batchWrite(keys []string, val []byte) error {
	errs := make([]string, 0)

	for _, key := range keys {
		err := u.to.Write(&store.Record{
			Key:   key,
			Value: val,
		})

		if err != nil {
			errs = append(errs, err.Error())
		}

	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ";"))
	}

	return nil
}

// migratePerAccountUser migrates per account's users
func (u *userMigration) migratePerAccountUser(authAccount *auth.Account) error {
	// read all old users from db
	list, err := u.List(0, 10, authAccount)
	if err != nil {
		return err
	}

	for _, rec := range list {
		val, err := json.Marshal(rec)
		if err != nil {
			logger.Errorf("json marshal rec error: %v, %+v", err, rec)
			continue
		}

		keys := []string{
			generateAccountStoreKey(rec.Id),
			generateAccountEmailStoreKey(rec.Email),
			generateAccountUsernameStoreKey(rec.Username),
		}

		if err := u.batchWrite(keys, val); err != nil {
			logger.Errorf("rec migrate batch write error: %v, %+v", err, keys)
			continue
		}
	}

	// write to store

	return nil
}

func (u *userMigration) Do() error {
	ctx := context.DefaultContext

	// get accounts list
	resp, err := u.authAccount.List(ctx, &authPb.ListAccountsRequest{}, client.WithAuthToken())
	if err != nil {
		return errors.Wrap(err, "get account list error")
	}

	// migrate every account's data
	for _, act := range resp.Accounts {
		err := u.migratePerAccountUser(&auth.Account{
			ID:       act.Id,
			Type:     act.Type,
			Issuer:   act.Issuer,
			Metadata: act.Metadata,
			Scopes:   act.Scopes,
			Secret:   act.Secret,
			Name:     act.Name,
		})

		if err != nil {
			logger.Errorf("migrate account's users data error:", err)
		}
	}

	return nil
}
