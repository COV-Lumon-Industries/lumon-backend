package interfaces

import (
	"lumon-backend/internal/domain/models"

	"github.com/google/uuid"
)

type DocumentRepository interface {
	Create(document *models.Document) error
	GetByID(id uuid.UUID) (*models.Document, error)
	GetByType(tag string) (*models.Document, error)
	Update(document *models.Document) error
	Delete(id uuid.UUID) error
	List(page, pageSize int) ([]models.Document, int64, error)
}
