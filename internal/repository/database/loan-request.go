package database

import (
	"errors"
	"fmt"

	"lumon-backend/internal/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoanRequestRepositoryImpl struct {
	db *gorm.DB
}

func NewLoanRequestRepository(db *gorm.DB) *LoanRequestRepositoryImpl {
	return &LoanRequestRepositoryImpl{db: db}
}

func (r *LoanRequestRepositoryImpl) Create(loanRequest *models.LoanRequest) error {
	if loanRequest == nil {
		return errors.New("loan request cannot be nil")
	}

	return r.db.Create(loanRequest).Error
}

func (r *LoanRequestRepositoryImpl) GetByID(id uuid.UUID) (*models.LoanRequest, error) {
	var loanRequest models.LoanRequest

	err := r.db.Preload("Borrower").First(&loanRequest, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("loan request with ID %s not found", id)
		}
		return nil, err
	}

	return &loanRequest, nil
}

func (r *LoanRequestRepositoryImpl) GetByBorrower(borrowerID uuid.UUID) ([]models.LoanRequest, error) {
	var loanRequests []models.LoanRequest

	err := r.db.Where("borrower_id = ?", borrowerID).Preload("Borrower").Find(&loanRequests).Error
	if err != nil {
		return nil, err
	}

	return loanRequests, nil
}

func (r *LoanRequestRepositoryImpl) Update(loanRequest *models.LoanRequest) error {
	if loanRequest == nil {
		return errors.New("loan request cannot be nil")
	}

	return r.db.Save(&loanRequest).Error
}

func (r *LoanRequestRepositoryImpl) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.LoanRequest{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("loan request with ID %s not found", id)
	}

	return nil
}

func (r *LoanRequestRepositoryImpl) List(page, pageSize int) ([]models.LoanRequest, int64, error) {
	var loanRequests []models.LoanRequest
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&models.LoanRequest{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(pageSize).Preload("Borrower").Find(&loanRequests).Error; err != nil {
		return nil, 0, err
	}

	return loanRequests, total, nil
}

func (r *LoanRequestRepositoryImpl) UpdateStatus(id uuid.UUID, status string) error {
	result := r.db.Model(&models.LoanRequest{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("loan request with ID %s not found", id)
	}
	return nil
}
