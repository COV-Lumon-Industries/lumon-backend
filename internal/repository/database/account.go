package database

import (
	"errors"
	"fmt"

	"lumon-backend/internal/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountRepositoryImpl struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepositoryImpl {
	return &AccountRepositoryImpl{db: db}
}

func (r *AccountRepositoryImpl) Create(account *models.Account) error {
	if account == nil {
		return errors.New("account cannot be nil")
	}

	return r.db.Create(account).Error
}

func (r *AccountRepositoryImpl) GetByID(id uuid.UUID) (*models.Account, error) {
	var account models.Account

	err := r.db.Preload("Document").First(&account, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("account with ID %s not found", id)
		}
		return nil, err
	}

	return &account, nil
}

func (r *AccountRepositoryImpl) GetByType(accountType string) (*models.Account, error) {
	var account models.Account

	err := r.db.Preload("Document").First(&account, "type = ?", accountType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("account with type %s not found", accountType)
		}
		return nil, err
	}

	return &account, nil
}

func (r *AccountRepositoryImpl) Update(account *models.Account) error {
	if account == nil {
		return errors.New("account cannot be nil")
	}

	return r.db.Save(account).Error
}

func (r *AccountRepositoryImpl) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.Account{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("account with ID %s not found", id)
	}

	return nil
}

func (r *AccountRepositoryImpl) List(page, pageSize int) ([]models.Account, int64, error) {
	var accounts []models.Account
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&models.Account{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(pageSize).Preload("Document").Find(&accounts).Error; err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}
