package migrate

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/micro/micro/v3/service/context"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"

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
	from   db.DbService
	to     store.Store
	domain *domain.Domain
}

func NewUserMigration(from db.DbService, to store.Store) *userMigration {
	return &userMigration{
		from:   from,
		to:     to,
		domain: domain.New(to),
	}
}

func (u *userMigration) List(o, l int32) ([]*user.Account, error) {
	var limit int32 = 25
	var offset int32 = 0

	if l > 0 {
		limit = l
	}
	if o > 0 {
		offset = o
	}

	rsp, err := u.from.Read(context.DefaultContext, &db.ReadRequest{
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
		var account user.Account
		json.Unmarshal(m, &account)
		ret[i] = &account
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

func (u *userMigration) Do() error {
	// read all old users from db
	list, err := u.List(0, 100)
	if err != nil {
		return err
	}

	for _, account := range list {
		val, err := json.Marshal(account)
		if err != nil {
			logger.Errorf("json marshal user error: %v, %+v", err, account)
			continue
		}

		keys := []string{
			generateAccountStoreKey(account.Id),
			generateAccountEmailStoreKey(account.Email),
			generateAccountUsernameStoreKey(account.Username),
		}

		if err := u.batchWrite(keys, val); err != nil {
			logger.Errorf("user migrate batch write error: %v, %+v", err, keys)
			continue
		}
	}

	// write to store

	return nil
}
