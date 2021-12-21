package migrate

import (
	"strings"
	"sync"

	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	pwdMgr "github.com/micro/services/user/migrate/password"
	userMgr "github.com/micro/services/user/migrate/user"
)

type migration struct {
	db    *gorm.DB
	store store.Store
}

func NewMigration(db *gorm.DB) *migration {
	return &migration{
		db:    db,
		store: store.DefaultStore,
	}
}

func (m *migration) Do() error {
	// get all tables
	var tables []string
	if err := m.db.Table("information_schema.tables").
		Where("table_schema = ?", "public").
		Pluck("table_name", &tables).Error; err != nil {
		return errors.Wrap(err, "get pgx tables error")
	}

	// migrate all data in concurrency
	wg := sync.WaitGroup{}
	// max concurrency is 5
	concurrencyChan := make(chan struct{}, 5)

	var strg Migration

	for _, t := range tables {

		wg.Add(1)
		concurrencyChan <- struct{}{}

		go func(tableName string) {
			defer func() {
				wg.Done()
				<-concurrencyChan
			}()

			if strings.HasSuffix(tableName, "_users") {
				strg = userMgr.New(m.store, strings.TrimSuffix(tableName, "_users"))
			} else if strings.HasSuffix(tableName, "_passwords") {
				strg = pwdMgr.New(m.store, strings.TrimSuffix(tableName, "_passwords"))
			}

			if strg == nil {
				return
			}

			ctx := NewContext(m.db, strg)
			if err := ctx.Migrate(tableName); err != nil {
				logger.Errorf("migrate table:%s error", tableName)
			}

		}(t)
	}

	wg.Wait()

	return nil

}
