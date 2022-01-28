package password

import (
	"fmt"

	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	"github.com/tidwall/gjson"

	"github.com/micro/services/user/migrate/entity"
)

func generatePasswordStoreKey(tenantId string, id string) string {
	return fmt.Sprintf("%spassword/%s", entity.KeyPrefix(tenantId), id)
}

type password struct {
	to       store.Store
	tenantId string
}

func New(to store.Store, tenantId string) *password {
	return &password{
		to:       to,
		tenantId: tenantId,
	}
}

func (u *password) migrate(rows []*entity.Row) error {
	for _, rec := range rows {
		id := gjson.Get(rec.Data, "id").String()

		key := generatePasswordStoreKey(u.tenantId, id)
		err := u.to.Write(&store.Record{
			Key:   key,
			Value: []byte(rec.Data),
		})

		if err != nil {
			logger.Errorf("migrate password write error: %v, %+v", err, key)
			continue
		}
	}

	return nil
}

func (u *password) Migrate(rows []*entity.Row) error {
	return u.migrate(rows)
}
