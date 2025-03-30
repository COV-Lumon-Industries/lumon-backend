package service

import (
	"lumon-backend/internal/domain/models"
	"lumon-backend/internal/domain/schemas"
	"lumon-backend/internal/repository/interfaces"
	"lumon-backend/pkg/common/logger"

	"github.com/google/uuid"
)

type AccountService struct {
	repo interfaces.AccountRepository
}

func NewAccountService(repo interfaces.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) CreateAccount(account *models.Account) error {
	return s.repo.Create(account)
}

func (s *AccountService) GetAccount(id string) (*schemas.AccountResponse, error) {
	dbAccount, err := s.repo.GetByID(uuid.MustParse(id))
	if err != nil {
		logger.APILogger.Error(err)
		return nil, err
	}

	return &schemas.AccountResponse{
		ID:   dbAccount.ID.String(),
		Type: dbAccount.Type,
		Name: dbAccount.Name,
	}, nil
}

func (s *AccountService) GetAccountByType(accountType string) (*schemas.AccountResponse, error) {
	dbAccount, err := s.repo.GetByType(accountType)
	if err != nil {
		logger.APILogger.Error(err)
		return nil, err
	}

	return &schemas.AccountResponse{
		ID:   dbAccount.ID.String(),
		Type: dbAccount.Type,
		Name: dbAccount.Name,
	}, nil
}

func (s *AccountService) UpdateAccount(account *models.Account) error {
	return s.repo.Update(account)
}

func (s *AccountService) DeleteAccount(id string) error {
	return s.repo.Delete(uuid.MustParse(id))
}

func (s *AccountService) ListAccounts(page, pageSize int) ([]schemas.AccountResponse, int64, error) {
	dbAccounts, total, err := s.repo.List(page, pageSize)
	if err != nil {
		logger.APILogger.Error(err)
		return nil, 0, err
	}

	accounts := make([]schemas.AccountResponse, len(dbAccounts))

	for i, dbAccount := range dbAccounts {
		accounts[i] = schemas.AccountResponse{
			ID:   dbAccount.ID.String(),
			Type: dbAccount.Type,
			Name: dbAccount.Name,
		}
	}

	return accounts, total, nil
}