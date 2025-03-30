package interfaces

import (
	"lumon-backend/internal/domain/models"
)

type WalletRepository interface {
	Create(wallet *models.Wallet) error
	Update(wallet *models.Wallet) error
	FindByID(id string) (*models.Wallet, error)
}
