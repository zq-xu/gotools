package store

import (
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"zq-xu/helper/log"
)

const (
	dsnFmt = "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=5s&readTimeout=6s"
)

var (
	gormDB *gorm.DB

	once sync.Once
)

type DatabaseInfo struct {
	Address      string
	Port         string
	Username     string
	Password     string
	DatabaseName string
}

func (dbInfo *DatabaseInfo) GenerateMysqlDSN() string {
	return fmt.Sprintf(dsnFmt,
		dbInfo.Username,
		dbInfo.Password,
		dbInfo.Address,
		dbInfo.Port,
		dbInfo.DatabaseName)
}

func InitGorm(dbCfg *DatabaseConfig) {
	once.Do(func() {
		gormDB = NewGormDB(dbCfg)

		err := autoMigrate(gormDB)
		if err != nil {
			log.Logger.Errorf("Failed to auto migrate, %v", err)
			return
		}

		log.Logger.Infof("Succeed to init db！")
	})
}

func NewGormDB(dbCfg *DatabaseConfig) *gorm.DB {
	log.Logger.Debugf("init db connection with %+v", dbCfg)

	db, err := gorm.Open(newMysqlDialector(dbCfg), newGormConfig(dbCfg))
	if err != nil {
		log.Logger.Fatal(err)
	}

	return db
}

func newMysqlDialector(dbCfg *DatabaseConfig) gorm.Dialector {
	return mysql.New(mysql.Config{
		DSN:                       dbCfg.GenerateMysqlDSN(),
		DefaultStringSize:         256,
		DefaultDatetimePrecision:  nil,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	})
}

func newGormConfig(dbCfg *DatabaseConfig) *gorm.Config {
	return &gorm.Config{
		Logger: logger.Default.LogMode(dbCfg.LogLevel),
	}
}
