package gormkit

import (
	"gorm.io/gorm"

	"github.com/zq-xu/gotools/utilsx"
)

var (
	tableSet = make([]interface{}, 0)
)

func RegisterTable(m interface{}) {
	if utilsx.IsInterfaceValueNil(m) {
		return
	}

	tableSet = append(tableSet, m)
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(tableSet...)
}
