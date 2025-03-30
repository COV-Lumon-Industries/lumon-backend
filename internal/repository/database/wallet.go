package database

import (
	"lumon-backend/internal/domain/models"
	"lumon-backend/internal/repository/interfaces"

	"gorm.io/gorm"
)

type walletRepository struct {
	db *gorm.DB
}

func (w *walletRepository) FindByID(id string) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := w.db.Preload("User").First(&wallet, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func NewWalletRepository(db *gorm.DB) interfaces.WalletRepository {
	return &walletRepository{db: db}
}

func (w *walletRepository) Create(wallet *models.Wallet) error {
	return w.db.Create(wallet).Error
}

func (w *walletRepository) Update(wallet *models.Wallet) error {
	return w.db.Save(&wallet).Error
}
