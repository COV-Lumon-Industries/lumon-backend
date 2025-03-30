package migrations

import (
	"lumon-backend/internal/domain/models"
)

func GetMigrationModels() []interface{} {
	mgrModel := []any{
		&models.User{},
		&models.Document{},
		&models.Transaction{},
		&models.LoanRequest{},
		&models.Account{},
		&models.Wallet{},
	}

	return mgrModel
}
