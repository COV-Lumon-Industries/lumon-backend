package interfaces

import (
	"lumon-backend/internal/domain/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	GetByName(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
	List(page, pageSize int) ([]models.User, int64, error)

	// GetByVerificationToken(token string) (*models.User, error)
	// SetEmailVerified(userID uuid.UUID) error
	// SetVerificationToken(userID uuid.UUID, token string, expiry time.Time) error

	// GetByResetToken(token string, currentTime time.Time) (*models.User, error)
	// UpdatePassword(userID uuid.UUID, hashedPassword string) error
	// ClearResetToken(userID uuid.UUID) error
	// SetResetToken(userID uuid.UUID, token string, expiry time.Time) error
}
