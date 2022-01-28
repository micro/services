package migrate

import (
	"strings"
	"sync"

	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	pwdMgr "github.com/micro/services/user/migrate/password"
	resetMgr "github.com/micro/services/user/migrate/password_reset_code"
	sessionMgr "github.com/micro/services/user/migrate/session"
	tokenMgr "github.com/micro/services/user/migrate/token"
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

	for _, t := range tables {

		wg.Add(1)
		concurrencyChan <- struct{}{}

		go func(tableName string) {
			defer func() {
				wg.Done()
				<-concurrencyChan
			}()
			var strg Migration

			if strings.HasSuffix(tableName, "_users") {
				strg = userMgr.New(m.store, strings.TrimSuffix(tableName, "_users"))
			} else if strings.HasSuffix(tableName, "_passwords") {
				strg = pwdMgr.New(m.store, strings.TrimSuffix(tableName, "_passwords"))
			} else if strings.HasSuffix(tableName, "_sessions") {
				strg = sessionMgr.New(m.store, strings.TrimSuffix(tableName, "_sessions"))
			} else if strings.HasSuffix(tableName, "_tokens") {
				strg = tokenMgr.New(m.store, strings.TrimSuffix(tableName, "_tokens"))
			} else if strings.HasSuffix(tableName, "_password_reset_codes") {
				strg = resetMgr.New(m.store, strings.TrimSuffix(tableName, "_password_reset_codes"))
			} else {
				logger.Infof("ignore table: %s", tableName)
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
