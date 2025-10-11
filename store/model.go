package store

import (
	"fmt"
	"time"
	"zq-xu/gotools/utils"

	"gorm.io/gorm"
)

type Model struct {
	ID        string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Comment   string         `gorm:"size:512"`
	Status    int
}

type ModelResponse struct {
	ID        string         `json:"id"`
	CreatedAt utils.UnixTime `json:"createdAt"`
	UpdatedAt utils.UnixTime `json:"updatedAt"`

	Comment string `json:"comment"`
}

// GenerateModel
func GenerateModel() Model {
	return GenerateModelWithID(utils.GenerateStringUUID())
}

// GenerateModelWithID
func GenerateModelWithID(id string) Model {
	return Model{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (m *Model) GetID() string     { return m.ID }
func (m *Model) GetStatus() string { return fmt.Sprintf("%d", m.Status) }

func (m Model) SetComment(str string) Model {
	m.Comment = str
	return m
}
