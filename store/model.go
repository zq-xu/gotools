package store

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/zq-xu/gotools/utils"
)

type Model struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Comment   string         `gorm:"size:512"`
	Status    int
}

// GenerateModel
func GenerateModel() Model {
	return GenerateModelWithID(utils.GenerateUUID())
}

// GenerateModelWithID
func GenerateModelWithID(id int64) Model {
	return Model{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (m *Model) GetID() string     { return fmt.Sprintf("%d", m.ID) }
func (m *Model) GetStatus() string { return fmt.Sprintf("%d", m.Status) }

func (m Model) SetComment(str string) Model {
	m.Comment = str
	return m
}
