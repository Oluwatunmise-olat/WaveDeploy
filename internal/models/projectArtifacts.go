package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ProjectArtifacts struct {
	Id             uuid.UUID `gorm:"primaryKey"`
	ProjectId      string    `gorm:"column:project_id"`
	AccountId      string    `gorm:"column:account_id"`
	IsLive         bool      `gorm:"column:is_live"`
	Tag            string    `gorm:"column:tag"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
	gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (ProjectArtifacts) TableName() string { return "project_artifacts" }
