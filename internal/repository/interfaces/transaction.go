package interfaces

import (
	"time"

	"lumon-backend/internal/domain/models"

	"github.com/google/uuid"
)

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	CreateBatch(transactions []*models.Transaction, batchSize int) error
	GetByID(id uuid.UUID) (*models.Transaction, error)
	GetByType(transactionType string, page, pageSize int) ([]models.Transaction, int64, error)
	GetByAccountNumber(accountNumber string, page, pageSize int) ([]models.Transaction, int64, error)
	GetByDateRange(start, end time.Time, page, pageSize int) ([]models.Transaction, int64, error)
	GetByAmountRange(min, max float64, page, pageSize int) ([]models.Transaction, int64, error)
	Update(transaction *models.Transaction) error
	Delete(id uuid.UUID) error
	List(userId uuid.UUID, page, pageSize int) ([]models.Transaction, int64, error)
	ListAll(userId uuid.UUID) ([]models.Transaction, error)
}
