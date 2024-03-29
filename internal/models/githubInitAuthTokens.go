package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type GithubInitAuthTokens struct {
	Id             uuid.UUID `gorm:"primaryKey"`
	Token          string    `gorm:"column:token"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
	gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (GithubInitAuthTokens) TableName() string { return "github_init_auth_tokens" }
