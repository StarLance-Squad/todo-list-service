package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model         // Embedding gorm.Model which includes fields ID, CreatedAt, UpdatedAt, DeletedAt
	Username    string `gorm:"unique"`
	Email       string `gorm:"unique"`
	Password    string
	LastLoginAt time.Time
}

// UpdateLastLogin updates the LastLoginAt field to the current time
func (u *User) UpdateLastLogin(db *gorm.DB) error {
	u.LastLoginAt = time.Now()
	return db.Save(u).Error
}
