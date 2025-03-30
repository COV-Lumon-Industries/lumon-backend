package service

import (
	"lumon-backend/internal/domain/models"
	"lumon-backend/internal/domain/schemas"
	"lumon-backend/internal/repository/interfaces"
	"lumon-backend/pkg/common/logger"

	"github.com/google/uuid"
)

type DocumentService struct {
	repo interfaces.DocumentRepository
}

func NewDocumentService(repo interfaces.DocumentRepository) *DocumentService {
	return &DocumentService{repo: repo}
}

func (s *DocumentService) CreateDocument(document *models.Document) error {
	return s.repo.Create(document)
}

func (s *DocumentService) GetDocument(id string) (*schemas.DocumentResponse, error) {
	dbDocument, err := s.repo.GetByID(uuid.MustParse(id))
	if err != nil {
		logger.APILogger.Error(err)
		return nil, err
	}

	return &schemas.DocumentResponse{
		ID:             dbDocument.ID.String(),
		ContentSummary: dbDocument.ContentSummary,
		Type:           dbDocument.ContentSummary,
		UploadedAt:     dbDocument.UploadedAt,
	}, nil
}

func (s *DocumentService) GetDocumentByTypeTag(tag string) (*schemas.DocumentResponse, error) {
	dbDocument, err := s.repo.GetByType(tag)
	if err != nil {
		logger.APILogger.Error(err)
		return nil, err
	}

	return &schemas.DocumentResponse{
		ID:             dbDocument.ID.String(),
		ContentSummary: dbDocument.ContentSummary,
		Type:           dbDocument.ContentSummary,
		UploadedAt:     dbDocument.UploadedAt,
	}, nil
}

func (s *DocumentService) DeleteDocument(id string) error {
	return s.repo.Delete(uuid.MustParse(id))
}

func (s *DocumentService) ListDocuments(page, pageSize int) ([]schemas.DocumentResponse, int64, error) {
	dbDocuments, total, err := s.repo.List(page, pageSize)
	if err != nil {
		logger.APILogger.Error(err)
		return nil, 0, err
	}

	documents := make([]schemas.DocumentResponse, total)

	for i, dbDocument := range dbDocuments {
		documents[i] = schemas.DocumentResponse{
			ID:             dbDocument.ID.String(),
			ContentSummary: dbDocument.ContentSummary,
			Type:           dbDocument.ContentSummary,
			UploadedAt:     dbDocument.UploadedAt,
		}
	}

	return documents, total, nil
}
