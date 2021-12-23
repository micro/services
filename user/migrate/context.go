package migrate

import (
	"github.com/micro/micro/v3/service/logger"
	"gorm.io/gorm"

	"github.com/micro/services/user/migrate/entity"
)

type Migration interface {
	Migrate([]*entity.Row) error
}

type context struct {
	db       *gorm.DB
	strategy Migration
}

func NewContext(db *gorm.DB, strg Migration) *context {
	return &context{
		db:       db,
		strategy: strg,
	}
}

func (c *context) Migrate(tableName string) error {
	var count int64

	db := c.db.Table(tableName)

	if err := db.Count(&count).Error; err != nil {
		return err
	}

	var offset, limit = 0, 1000

	for offset = 0; offset < int(count); offset = offset + limit {
		rows := make([]*entity.Row, 0)

		if err := db.Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
			logger.Errorf("migrate error, table:%v offset:%v limit:%v error:%v", tableName, offset, limit, err)
			continue
		}

		if err := c.strategy.Migrate(rows); err != nil {
			logger.Errorf("migrate error, table:%v offset:%v limit:%v error:%v", tableName, offset, limit, err)
			continue
		}
	}

	logger.Infof("migrate done, table: %v, rows count: %v", tableName, count)

	return nil
}
