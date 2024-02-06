package models

import (
	"gorm.io/gorm"
	"time"
)

type GithubApps struct {
	Id             string    `gorm:"primaryKey"`
	AccountId      string    `gorm:"column:account_id"`
	InstallationId string    `gorm:"column:installation_id"`
	Code           string    `gorm:"column:code"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
	gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (GithubApps) TableName() string { return "github_apps" }