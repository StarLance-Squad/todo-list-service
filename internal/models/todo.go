package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Todo struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;"`
	Title       string
	Description string
	Completed   bool
	UserID      float64
	User        User // Associating with the User model
	gorm.Model
}

func (todo *Todo) BeforeCreate(tx *gorm.DB) (err error) {
	todo.ID = uuid.New()
	return
}
