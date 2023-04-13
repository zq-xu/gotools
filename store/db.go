package store

import (
	"context"

	"gorm.io/gorm"
)

func InitDatabase() {
	InitGorm(DatabaseCfg)
}

func DB(ctx context.Context) *gorm.DB {
	return gormDB.WithContext(ctx)
}
