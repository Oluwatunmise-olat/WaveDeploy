package respository

import (
	"gorm.io/gorm"
)

type GithubInitAuthTokensRepository struct {
	DB *gorm.DB
}
