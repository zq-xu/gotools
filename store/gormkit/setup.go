package gormkit

import (
	"fmt"

	"github.com/rotisserie/eris"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitGorm
func InitGorm() error {
	db, err := newGormDB()
	if err != nil {
		return eris.Wrap(err, "failed to new gorm db")
	}

	err = autoMigrate(db)
	if err != nil {
		return eris.Wrap(err, "failed to auto migrate")
	}

	GlobalDB = &gormDB{db: db}
	return nil
}

func newGormDB() (*gorm.DB, error) {
	return gorm.Open(newMysqlDialector(), newGormConfig())
}

func newMysqlDialector() gorm.Dialector {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=5s&readTimeout=6s",
		GormConfig.Username, GormConfig.Password, GormConfig.Address, GormConfig.Port, GormConfig.DatabaseName)

	return mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,
		DefaultDatetimePrecision:  nil,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	})
}

func newGormConfig() *gorm.Config {
	return &gorm.Config{
		Logger: logger.Default.LogMode(GormConfig.LogLevel),
	}
}
