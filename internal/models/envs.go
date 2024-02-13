package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Envs struct {
	Id             uuid.UUID `gorm:"primaryKey"`
	ProjectId      uuid.UUID `gorm:"column:project_id"`
	AccountId      uuid.UUID `gorm:"column:account_id"`
	Key            string    `gorm:"key"`
	Value          string    `gorm:"key"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
	gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (Envs) TableName() string { return "envs" }
