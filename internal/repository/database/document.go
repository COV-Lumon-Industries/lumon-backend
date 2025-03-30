package database

import (
	"errors"
	"fmt"

	"lumon-backend/internal/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentRepositoryImpl struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) *DocumentRepositoryImpl {
	return &DocumentRepositoryImpl{db: db}
}

func (r *DocumentRepositoryImpl) Create(document *models.Document) error {
	if document == nil {
		return errors.New("document cannot be nil")
	}

	return r.db.Create(document).Error
}

func (r *DocumentRepositoryImpl) GetByID(id uuid.UUID) (*models.Document, error) {
	var document models.Document

	err := r.db.Preload("User").First(&document, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("document with ID %s not found", id)
		}
		return nil, err
	}

	return &document, nil
}

func (r *DocumentRepositoryImpl) GetByType(docType string) (*models.Document, error) {
	var document models.Document

	err := r.db.Preload("User").First(&document, "type = ?", docType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("document with type %s not found", docType)
		}
		return nil, err
	}

	return &document, nil
}

func (r *DocumentRepositoryImpl) Update(document *models.Document) error {
	if document == nil {
		return errors.New("document cannot be nil")
	}

	return r.db.Save(document).Error
}

func (r *DocumentRepositoryImpl) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.Document{}, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("document with ID %s not found", id)
	}

	return nil
}

func (r *DocumentRepositoryImpl) List(page, pageSize int) ([]models.Document, int64, error) {
	var documents []models.Document
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&models.Document{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(pageSize).Preload("User").Find(&documents).Error; err != nil {
		return nil, 0, err
	}

	return documents, total, nil
}
