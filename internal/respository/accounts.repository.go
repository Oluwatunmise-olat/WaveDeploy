package respository

import (
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/db"
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/models"
	"sync"
)

var (
	accountsRepository        *AccountsRepository
	accountRepositoryInitOnce sync.Once
)

func (ar *AccountsRepository) initializeAccountsRepository() *AccountsRepository {
	accountRepositoryInitOnce.Do(func() {
		accountsRepository = &AccountsRepository{
			DB: db.DB,
		}
	})

	return accountsRepository
}

func (ar *AccountsRepository) GetAccountByEmail(email string) (*models.Accounts, error) {
	var account models.Accounts
	err := ar.initializeAccountsRepository().DB.Model(&models.Accounts{}).First(&account, "email = ? and deleted_at is null", email).Error
	return &account, err
}

func (ar *AccountsRepository) GetAccountById(accountId string) (*models.Accounts, error) {
	var account models.Accounts
	err := ar.initializeAccountsRepository().DB.Model(&models.Accounts{}).First(&account, "id = ? and deleted_at is null", accountId).Error
	return &account, err
}

func (ar *AccountsRepository) CreateAccount(account models.Accounts) error {
	err := ar.initializeAccountsRepository().DB.Model(&models.Accounts{}).Create(account).Error
	return err
}
