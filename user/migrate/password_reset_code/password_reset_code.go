package password_reset_code

import (
	"fmt"

	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	"github.com/tidwall/gjson"

	"github.com/micro/services/user/migrate/entity"
)

func generatePasswordStoreKey(tenantId string, id string) string {
	return fmt.Sprintf("%spassword-reset-codes/%s", entity.KeyPrefix(tenantId), id)
}

type resetCode struct {
	to       store.Store
	tenantId string
}

func New(to store.Store, tenantId string) *resetCode {
	return &resetCode{
		to:       to,
		tenantId: tenantId,
	}
}

func (u *resetCode) migrate(rows []*entity.Row) error {
	for _, rec := range rows {
		id := gjson.Get(rec.Data, "id").String()

		key := generatePasswordStoreKey(u.tenantId, id)
		err := u.to.Write(&store.Record{
			Key:   key,
			Value: []byte(rec.Data),
		})

		if err != nil {
			logger.Errorf("migrate password-reset-code write error: %v, %+v", err, key)
			continue
		}
	}

	return nil
}

func (u *resetCode) Migrate(rows []*entity.Row) error {
	return u.migrate(rows)
}
