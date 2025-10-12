package gormkit

import (
	"gorm.io/gorm"

	"github.com/zq-xu/gotools/utils"
)

var (
	tableSet = make([]interface{}, 0)
)

func RegisterTable(m interface{}) {
	if utils.IsInterfaceValueNil(m) {
		return
	}

	tableSet = append(tableSet, m)
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(tableSet...)
}
