package session

import (
	"fmt"

	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	"github.com/tidwall/gjson"

	"github.com/micro/services/user/migrate/entity"
)

func generateStoreKey(tenantId string, id string) string {
	return fmt.Sprintf("%ssession/%s", entity.KeyPrefix(tenantId), id)
}

type session struct {
	to       store.Store
	tenantId string
}

func New(to store.Store, tenantId string) *session {
	return &session{
		to:       to,
		tenantId: tenantId,
	}
}

func (s *session) migrate(rows []*entity.Row) error {
	for _, rec := range rows {
		id := gjson.Get(rec.Data, "id").String()

		key := generateStoreKey(s.tenantId, id)
		err := s.to.Write(&store.Record{
			Key:   key,
			Value: []byte(rec.Data),
		})

		if err != nil {
			logger.Errorf("migrate session write error: %v, %+v", err, key)
			continue
		}
	}

	return nil
}

func (s *session) Migrate(rows []*entity.Row) error {
	return s.migrate(rows)
}
