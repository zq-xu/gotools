package store

import (
	"github.com/zq-xu/gotools/config"
	"github.com/zq-xu/gotools/store/database"
	"github.com/zq-xu/gotools/store/gormkit"
)

var (
	RegisterTable = gormkit.RegisterTable
)

func init() {
	config.Register("database", &gormkit.GormConfig, gormkit.InitGorm)
}

func DB() database.Database {
	return gormkit.GlobalDB
}

func SetDB(d database.Database) {
	gormkit.GlobalDB = d
}
