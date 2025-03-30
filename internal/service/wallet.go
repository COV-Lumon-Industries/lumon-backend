package service

import (
	"fmt"
	"lumon-backend/internal/domain/models"
	"lumon-backend/internal/repository/interfaces"
	"lumon-backend/pkg/common/logger"

	"github.com/google/uuid"
)

type WalletService struct {
	repo interfaces.WalletRepository
}

func NewWalletService(repo interfaces.WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) TopUpAccount(wallet *models.Wallet) error {
	if wallet.Balance <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	existingWallet, err := s.repo.FindByID(wallet.ID.String())
	if err != nil {
		if err.Error() == "record not found" { 
			err = s.repo.Create(wallet)
			if err != nil {
				logger.APILogger.Errorf("Failed to create wallet: %v", err)
				return fmt.Errorf("failed to create wallet: %v", err)
			}
			return nil
		}
		logger.APILogger.Errorf("Failed to get wallet: %v", err)
		return fmt.Errorf("failed to get wallet: %v", err)
	}

	newBalance := existingWallet.Balance + wallet.Balance
	existingWallet.Balance = newBalance

	err = s.repo.Update(existingWallet)
	if err != nil {
		logger.APILogger.Errorf("Failed to update wallet: %v", err)
		return fmt.Errorf("failed to update wallet: %v", err)
	}

	return nil
}

func (s *WalletService) WithdrawAccount(walletID string, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	id, err := uuid.Parse(walletID)
	if err != nil {
		return fmt.Errorf("invalid wallet ID: %v", err)
	}

	currentWallet, err := s.repo.FindByID(id.String())
	if err != nil {
		logger.APILogger.Errorf("Failed to get wallet: %v", err)
		return fmt.Errorf("failed to get wallet: %v", err)
	}

	if currentWallet.Balance < amount {
		return fmt.Errorf("insufficient funds")
	}

	newBalance := currentWallet.Balance - amount
	currentWallet.Balance = newBalance

	err = s.repo.Update(currentWallet)
	if err != nil {
		logger.APILogger.Errorf("Failed to update wallet: %v", err)
		return fmt.Errorf("failed to update wallet: %v", err)
	}

	return nil
}
