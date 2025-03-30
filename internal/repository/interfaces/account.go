package interfaces

import (
	"lumon-backend/internal/domain/models"

	"github.com/google/uuid"
)

type AccountRepository interface {
	Create(account *models.Account) error
	GetByID(id uuid.UUID) (*models.Account, error)
	GetByType(accountType string) (*models.Account, error)
	Update(account *models.Account) error
	Delete(id uuid.UUID) error
	List(page, pageSize int) ([]models.Account, int64, error)
}
