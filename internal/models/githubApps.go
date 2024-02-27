package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type GithubApps struct {
	Id             uuid.UUID `gorm:"primaryKey"`
	AccountId      string    `gorm:"column:account_id"`
	InstallationId string    `gorm:"column:installation_id"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
	gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (GithubApps) TableName() string { return "github_apps" }
