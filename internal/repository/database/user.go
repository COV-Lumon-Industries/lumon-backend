package database

import (
	"lumon-backend/internal/domain/models"
	"lumon-backend/internal/repository/interfaces"
	"lumon-backend/pkg/common/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	var wallet = &models.Wallet{}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			logger.APILogger.Error(err)
			return err
		}
		wallet.UserID = user.ID
		if err := tx.Create(wallet).Error; err != nil {
			logger.APILogger.Error(err)
			return err
		}
		return nil
	})
}

func (r *userRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Accounts").Preload("LoanRequests").First(&user, id).Error
	if err != nil {
		logger.APILogger.Error(err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByName(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).Preload("Accounts").Preload("LoanRequests").First(&user).Error
	if err != nil {
		logger.APILogger.Error(err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).Preload("Accounts").Preload("LoanRequests").First(&user).Error
	if err != nil {
		logger.APILogger.Error(err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *userRepository) List(page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	err := r.db.Model(&models.User{}).Count(&total).Error
	if err != nil {
		logger.APILogger.Error(err)
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = r.db.Offset(offset).Limit(pageSize).Preload("Accounts").Preload("LoanRequests").Find(&users).Error
	if err != nil {
		logger.APILogger.Error(err)
		return nil, 0, err
	}

	return users, total, nil
}

// func (r *userRepository) GetByResetToken(token string, currentTime time.Time) (*models.User, error) {
// 	var user models.User
// 	err := r.db.Where("password_reset_token = ? AND reset_token_expiry > ?",
// 		token, currentTime).First(&user).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }

// func (r *userRepository) UpdatePassword(userID uuid.UUID, hashedPassword string) error {
// 	return r.db.Model(&models.User{}).Where("id = ?", userID).
// 		Update("password", hashedPassword).Error
// }

// func (r *userRepository) ClearResetToken(userID uuid.UUID) error {
// 	return r.db.Model(&models.User{}).Where("id = ?", userID).
// 		Updates(map[string]interface{}{
// 			"password_reset_token": "",
// 			"reset_token_expiry":   nil,
// 		}).Error
// }

// func (r *userRepository) SetResetToken(userID uuid.UUID, token string, expiry time.Time) error {
// 	return r.db.Model(&models.User{}).Where("id = ?", userID).
// 		Updates(map[string]interface{}{
// 			"password_reset_token": token,
// 			"reset_token_expiry":   expiry,
// 		}).Error
// }

// func (r *userRepository) GetByVerificationToken(token string) (*models.User, error) {
// 	if token == "" {
// 		return nil, errors.New("verification token is required")
// 	}

// 	var user models.User
// 	err := r.db.Where("verification_token = ? AND token_expiry > ?",
// 		token, time.Now()).First(&user).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, errors.New("invalid or expired verification token")
// 		}
// 		return nil, err
// 	}

// 	return &user, nil
// }

// func (r *userRepository) SetEmailVerified(userID uuid.UUID) error {
// 	if userID.String() == "" {
// 		return errors.New("user ID is required")
// 	}

// 	result := r.db.Model(&models.User{}).Where("id = ?", userID).
// 		Updates(map[string]interface{}{
// 			"email_verified":     true,
// 			"verification_token": "",
// 			"token_expiry":       nil,
// 		})

// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	if result.RowsAffected == 0 {
// 		return errors.New("user not found")
// 	}

// 	return nil
// }

// func (r *userRepository) SetVerificationToken(userID uuid.UUID, token string, expiry time.Time) error {
// 	if userID.String() == "" {
// 		return errors.New("user ID is required")
// 	}

// 	if token == "" {
// 		return errors.New("verification token is required")
// 	}

// 	result := r.db.Model(&models.User{}).Where("id = ?", userID).
// 		Updates(map[string]interface{}{
// 			"verification_token": token,
// 			"token_expiry":       expiry,
// 			"email_verified":     false,
// 		})

// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	if result.RowsAffected == 0 {
// 		return errors.New("user not found")
// 	}

// 	return nil
// }
