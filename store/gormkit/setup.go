package gormkit

import (
	"fmt"

	"github.com/rotisserie/eris"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"zq-xu/gotools/config"
)

// InitGorm
func InitGorm(cfg *config.Config) error {
	db, err := newGormDB(&cfg.DatabaseConfig)
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

func newGormDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	return gorm.Open(newMysqlDialector(cfg), newGormConfig(cfg))
}

func newMysqlDialector(cfg *config.DatabaseConfig) gorm.Dialector {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=5s&readTimeout=6s",
		cfg.Username, cfg.Password, cfg.Address, cfg.Port, cfg.DatabaseName)

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

func newGormConfig(cfg *config.DatabaseConfig) *gorm.Config {
	return &gorm.Config{
		Logger: logger.Default.LogMode(cfg.LogLevel),
	}
}
