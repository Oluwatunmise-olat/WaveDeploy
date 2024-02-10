package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Projects struct {
	Id             uuid.UUID `gorm:"primaryKey"`
	AccountId      string    `gorm:"column:account_id"`
	Name           string    `gorm:"column:name"`
	GithubRepoUrl  string    `gorm:"column:github_repo_url"`
	GithubBranch   string    `gorm:"column:github_branch"`
	GithubCommit   string    `gorm:"column:github_commit"`
	IsLive         bool      `gorm:"column:is_live"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
	gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (Projects) TableName() string { return "projects" }
