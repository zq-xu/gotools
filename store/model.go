package store

import (
	"time"

	"gorm.io/gorm"

	"zq-xu/helper/utils"
)

var (
	tableSet = make([]interface{}, 0)
)

type Model struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func GenerateModel() Model {
	return GenerateModelWithID(utils.GenerateUUID())
}

func GenerateModelWithID(id int64) Model {
	return Model{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func RegisterTable(m interface{}) {
	if utils.IsInterfaceValueNil(m) {
		return
	}

	tableSet = append(tableSet, m)
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(tableSet...)
}
