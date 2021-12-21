package user

import (
	"fmt"
	"strings"

	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/micro/services/user/migrate/entity"
)

func getStoreKeyPrefix(tenantId string) string {
	return fmt.Sprintf("user/%s/", tenantId)
}

func generateAccountStoreKey(tenantId, userId string) string {
	return fmt.Sprintf("%saccount/id/%s", getStoreKeyPrefix(tenantId), userId)
}

func generateAccountEmailStoreKey(tenantId, email string) string {
	return fmt.Sprintf("%sacccount/email/%s", getStoreKeyPrefix(tenantId), email)
}

func generateAccountUsernameStoreKey(tenantId, username string) string {
	return fmt.Sprintf("%saccount/username/%s", getStoreKeyPrefix(tenantId), username)
}

type user struct {
	to       store.Store
	tenantId string
}

func New(to store.Store, tenantId string) *user {
	return &user{
		to:       to,
		tenantId: tenantId,
	}
}

func (u *user) batchWrite(keys []string, val []byte) error {
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

func (u *user) migrate(rows []*entity.Row) error {
	for _, rec := range rows {
		id := gjson.Get(rec.Data, "id").String()
		email := gjson.Get(rec.Data, "email").String()
		username := gjson.Get(rec.Data, "username").String()

		fmt.Println("--> username", id, email)

		keys := []string{
			generateAccountStoreKey(u.tenantId, id),
			generateAccountEmailStoreKey(u.tenantId, email),
			generateAccountUsernameStoreKey(u.tenantId, username),
		}

		if err := u.batchWrite(keys, []byte(rec.Data)); err != nil {
			logger.Errorf("migrate users batch write error: %v, %+v", err, keys)
			continue
		}
	}

	return nil
}

func (u *user) Migrate(rows []*entity.Row) error {
	return u.migrate(rows)
}
