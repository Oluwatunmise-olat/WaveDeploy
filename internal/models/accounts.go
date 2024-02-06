package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Accounts struct {
	Id             uuid.UUID `gorm:"primaryKey"`
	UserName       string    `gorm:"column:username"`
	Email          string    `gorm:"column:email"`
	Password       string    `gorm:"column:password"`
	LastAuthAt     time.Time `gorm:"column:last_auth_at"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
	gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (Accounts) TableName() string { return "accounts" }
