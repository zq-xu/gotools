package store

import (
	"zq-xu/gotools/setup"
	"zq-xu/gotools/store/database"
	"zq-xu/gotools/store/gormkit"
)

var (
	RegisterTable = gormkit.RegisterTable
)

func init() {
	setup.RegisterSetup("Gorm", gormkit.InitGorm)
}

func DB() database.Database {
	return gormkit.GlobalDB
}

func SetDB(d database.Database) {
	gormkit.GlobalDB = d
}
