package respository

import (
	"github.com/Oluwatunmise-olat/WaveDeploy/db"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/models"
	"gorm.io/gorm"
	"sync"
)

var (
	accountsRepository        *AccountsRepository
	accountRepositoryInitOnce sync.Once
)

type AccountsRepository struct {
	DB *gorm.DB
}

func InitializeAccountsRepository() *AccountsRepository {
	accountRepositoryInitOnce.Do(func() {
		accountsRepository = &AccountsRepository{
			DB: db.DB,
		}
	})

	return accountsRepository
}

func (ar *AccountsRepository) GetAccountByEmail(email string) (*models.Accounts, error) {
	var account models.Accounts
	err := ar.DB.First(&account, "email = ?", email).Error
	return &account, err
}

func (ar *AccountsRepository) CreateAccount(account models.Accounts) error {
	err := ar.DB.Create(account).Error
	return err
}
