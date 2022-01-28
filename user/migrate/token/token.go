package token

import (
	"fmt"

	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	"github.com/tidwall/gjson"

	"github.com/micro/services/user/migrate/entity"
)

func generateStoreKey(tenantId string, id string) string {
	return fmt.Sprintf("%sverification-token/%s", entity.KeyPrefix(tenantId), id)
}

type token struct {
	to       store.Store
	tenantId string
}

func New(to store.Store, tenantId string) *token {
	return &token{
		to:       to,
		tenantId: tenantId,
	}
}

func (s *token) migrate(rows []*entity.Row) error {
	for _, rec := range rows {
		id := gjson.Get(rec.Data, "id").String()

		key := generateStoreKey(s.tenantId, id)
		err := s.to.Write(&store.Record{
			Key:   key,
			Value: []byte(rec.Data),
		})

		if err != nil {
			logger.Errorf("migrate token write error: %v, %+v", err, key)
			continue
		}
	}

	return nil
}

func (s *token) Migrate(rows []*entity.Row) error {
	return s.migrate(rows)
}
