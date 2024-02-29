package respository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectsRepository struct {
	DB *gorm.DB
}

type UpdateProjectPayload struct {
	Project   map[string]interface{}
	ProjectId uuid.UUID
	AccountId uuid.UUID
	Trx       *gorm.DB
}

type EnvsRepository struct {
	DB *gorm.DB
}

type GithubInitAuthTokensRepository struct {
	DB *gorm.DB
}

type GithubAppsRepository struct {
	DB *gorm.DB
}

type AccountsRepository struct {
	DB *gorm.DB
}
