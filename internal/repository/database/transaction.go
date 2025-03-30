package database

import (
	"errors"
	"fmt"
	"time"

	"lumon-backend/internal/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepositoryImpl struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepositoryImpl {
	return &TransactionRepositoryImpl{db: db}
}

func (r *TransactionRepositoryImpl) Create(transaction *models.Transaction) error {
	if transaction == nil {
		return errors.New("transaction cannot be nil")
	}

	return r.db.Create(transaction).Error
}

func (r *TransactionRepositoryImpl) CreateBatch(transactions []*models.Transaction, batchSize int) error {
	if len(transactions) == 0 {
		return errors.New("no transactions provided for batch creation")
	}

	return r.db.CreateInBatches(transactions, batchSize).Error
}

func (r *TransactionRepositoryImpl) GetByID(id uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction

	err := r.db.Preload("User").First(&transaction, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("transaction with ID %s not found", id)
		}
		return nil, err
	}

	return &transaction, nil
}

func (r *TransactionRepositoryImpl) GetByType(
	transactionType string,
	page, pageSize int,
) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&models.Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where(
		"transaction_type = ?",
		transactionType,
	).Offset(offset).Limit(pageSize).Preload("User").Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *TransactionRepositoryImpl) GetByAccountNumber(
	accountNumber string,
	page, pageSize int,
) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&models.Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where(
		"to_number = ?",
		accountNumber,
	).Offset(offset).Limit(pageSize).Preload("User").Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *TransactionRepositoryImpl) GetByDateRange(
	start, end time.Time,
	page, pageSize int,
) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&models.Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where(
		"transaction_date BETWEEN ? AND ?",
		start, end,
	).Offset(offset).Limit(pageSize).Preload("User").Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *TransactionRepositoryImpl) GetByAmountRange(
	min, max float64,
	page, pageSize int,
) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&models.Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Where(
		"amount BETWEEN ? AND ?",
		min, max,
	).Offset(offset).Limit(pageSize).Preload("User").Find(&transactions).Error
	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *TransactionRepositoryImpl) Update(transaction *models.Transaction) error {
	if transaction == nil {
		return errors.New("transaction cannot be nil")
	}

	return r.db.Save(transaction).Error
}

func (r *TransactionRepositoryImpl) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.Transaction{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("transaction with ID %s not found", id)
	}

	return nil
}

func (r *TransactionRepositoryImpl) List(userId uuid.UUID, page, pageSize int) ([]models.Transaction, int64, error) {
	var transactions []models.Transaction
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&models.Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(pageSize).Preload("User").Find(&transactions, "user_id", userId).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *TransactionRepositoryImpl) ListAll(userId uuid.UUID) ([]models.Transaction, error) {
	var transactions []models.Transaction

	if err := r.db.Preload("User").Find(&transactions, "user_id", userId).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}
