package interfaces

import (
	"lumon-backend/internal/domain/models"

	"github.com/google/uuid"
)

type LoanRequestRepository interface {
	Create(loanRequest *models.LoanRequest) error
	GetByID(id uuid.UUID) (*models.LoanRequest, error)
	GetByBorrower(borrowerID uuid.UUID) ([]models.LoanRequest, error)
	Update(loanRequest *models.LoanRequest) error
	Delete(id uuid.UUID) error
	List(page, pageSize int) ([]models.LoanRequest, int64, error)
	UpdateStatus(id uuid.UUID, status string) error
}
