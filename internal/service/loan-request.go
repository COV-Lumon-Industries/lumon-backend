package service

import (
	"lumon-backend/internal/domain/models"
	"lumon-backend/internal/repository/interfaces"
	"lumon-backend/pkg/common/logger"

	"github.com/google/uuid"
)

type LoanRequestService struct {
	repo interfaces.LoanRequestRepository
}

func NewLoanRequestService(repo interfaces.LoanRequestRepository) *LoanRequestService {
	return &LoanRequestService{repo: repo}
}

func (s *LoanRequestService) CreateLoanRequest(loanRequest *models.LoanRequest) error {
	err := s.repo.Create(loanRequest)
	if err != nil {
		logger.APILogger.Error(err)
		return err
	}
	return nil
}

func (s *LoanRequestService) GetLoanRequest(id string) (*models.LoanRequest, error) {
	dbLoanRequest, err := s.repo.GetByID(uuid.MustParse(id))
	if err != nil {
		logger.APILogger.Error(err)
		return nil, err
	}

	return dbLoanRequest, nil
}

func (s *LoanRequestService) GetLoanRequestsByBorrower(borrowerID string) ([]models.LoanRequest, error) {
	requests, err := s.repo.GetByBorrower(uuid.MustParse(borrowerID))
	if err != nil {
		logger.APILogger.Error(err)
		return nil, err
	}

	return requests, nil
}

func (s *LoanRequestService) UpdateLoanRequest(loanRequest *models.LoanRequest) error {
	err := s.repo.Update(loanRequest)
	if err != nil {
		logger.APILogger.Error(err)
		return err
	}
	return nil
}

func (s *LoanRequestService) DeleteLoanRequest(id string) error {
	return s.repo.Delete(uuid.MustParse(id))
}

func (s *LoanRequestService) ListLoanRequests(page, pageSize int) ([]models.LoanRequest, int64, error) {
	dbLoanRequests, total, err := s.repo.List(page, pageSize)
	if err != nil {
		logger.APILogger.Error(err)
		return nil, 0, err
	}

	return dbLoanRequests, total, nil
}

func (s *LoanRequestService) UpdateLoanRequestStatus(id string, status string) error {
	err := s.repo.UpdateStatus(uuid.MustParse(id), status)
	if err != nil {
		logger.APILogger.Error(err)
		return err
	}
	return nil
}
