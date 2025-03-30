package service

import (
	"strings"
	"time"

	"lumon-backend/internal/domain/models"
	"lumon-backend/internal/domain/schemas"
	"lumon-backend/internal/repository/interfaces"
	"lumon-backend/pkg/common/logger"

	"github.com/google/uuid"
)

type TransactionService struct {
	repo interfaces.TransactionRepository
}

func NewTransactionService(repo interfaces.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) CreateTransaction(transaction *models.Transaction) error {
	return s.repo.Create(transaction)
}

func (s *TransactionService) CreateTransactionBatch(
	transactions []*schemas.MTNMoMoTransactionScrape, batchSize int, userID string,
) error {
	newTransactions := make([]*models.Transaction, len(transactions))

	for i, trans := range transactions {
		timeStr := strings.Trim(trans.TransactionDate, "\"")
		layout := "02-Jan-2006 03:04:05 PM"

		parsedTime, err := time.Parse(layout, timeStr)
		if err != nil {
			logger.APILogger.Error(err)
			return err
		}

		newTransactions[i] = &models.Transaction{
			TransactionDate: parsedTime,
			FromName:        trans.FromName,
			FromNumber:      trans.FromNumber,
			TransactionType: trans.TransactionType,
			Amount:          trans.Amount,
			ToNumber:        trans.ToNumber,
			ToName:          trans.ToName,
			Reference:       trans.Reference,
			FromAccount:     trans.FromAccount,
			Fees:            trans.Fees,
			BalanceBefore:   trans.BalanceBefore,
			BalanceAfter:    trans.BalanceAfter,
			ToAccount:       trans.ToAccount,
			UserID:          uuid.MustParse(userID),
		}
	}

	return s.repo.CreateBatch(newTransactions, batchSize)
}

func (s *TransactionService) UpdateTransaction(transaction *models.Transaction) error {
	return s.repo.Update(transaction)
}

func (s *TransactionService) DeleteTransaction(id string) error {
	return s.repo.Delete(uuid.MustParse(id))
}

func (s *TransactionService) GetTransaction(id string) (*schemas.TransactionResponse, error) {
	dbTransaction, err := s.repo.GetByID(uuid.MustParse(id))
	if err != nil {
		logger.APILogger.Error(err)
		return nil, err
	}

	return &schemas.TransactionResponse{
		UserID:          dbTransaction.UserID.String(),
		TransactionDate: dbTransaction.TransactionDate,
		FromName:        dbTransaction.FromName,
		FromNumber:      dbTransaction.FromNumber,
		TransactionType: dbTransaction.TransactionType,
		BalanceBefore:   dbTransaction.BalanceBefore,
		BalanceAfter:    dbTransaction.BalanceAfter,
		Amount:          dbTransaction.Amount,
		ToNumber:        dbTransaction.ToNumber,
		ToName:          dbTransaction.ToName,
		Reference:       dbTransaction.Reference,
	}, nil
}

func (s *TransactionService) GetTransactionsByType(
	transactionType string, page, pageSize int,
) ([]schemas.TransactionResponse, int64, error) {
	dbTransactions, total, err := s.repo.GetByType(transactionType, page, pageSize)
	if err != nil {
		logger.APILogger.Error(err)
		return nil, 0, err
	}

	transactions := make([]schemas.TransactionResponse, len(dbTransactions))
	for i, dbTransaction := range dbTransactions {
		transactions[i] = schemas.TransactionResponse{
			UserID:          dbTransaction.UserID.String(),
			TransactionDate: dbTransaction.TransactionDate,
			FromName:        dbTransaction.FromName,
			FromNumber:      dbTransaction.FromNumber,
			TransactionType: dbTransaction.TransactionType,
			BalanceBefore:   dbTransaction.BalanceBefore,
			BalanceAfter:    dbTransaction.BalanceAfter,
			Amount:          dbTransaction.Amount,
			ToNumber:        dbTransaction.ToNumber,
			ToName:          dbTransaction.ToName,
			Reference:       dbTransaction.Reference,
		}
	}

	return transactions, total, nil
}

func (s *TransactionService) GetTransactionsByAccountNumber(
	accountNumber string, page, pageSize int,
) ([]schemas.TransactionResponse, int64, error) {
	dbTransactions, total, err := s.repo.GetByAccountNumber(accountNumber, page, pageSize)
	if err != nil {
		logger.APILogger.Error(err)
		return nil, 0, err
	}

	transactions := make([]schemas.TransactionResponse, len(dbTransactions))
	for i, dbTransaction := range dbTransactions {
		transactions[i] = schemas.TransactionResponse{
			UserID:          dbTransaction.UserID.String(),
			TransactionDate: dbTransaction.TransactionDate,
			FromName:        dbTransaction.FromName,
			FromNumber:      dbTransaction.FromNumber,
			TransactionType: dbTransaction.TransactionType,
			BalanceBefore:   dbTransaction.BalanceBefore,
			BalanceAfter:    dbTransaction.BalanceAfter,
			Amount:          dbTransaction.Amount,
			ToNumber:        dbTransaction.ToNumber,
			ToName:          dbTransaction.ToName,
			Reference:       dbTransaction.Reference,
		}
	}

	return transactions, total, nil
}

func (s *TransactionService) GetTransactionsByDateRange(
	start, end string, page, pageSize int,
) ([]schemas.TransactionResponse, int64, error) {
	startTime, err := time.Parse(time.RFC3339, start)
	if err != nil {
		logger.APILogger.Error(err)
		return nil, 0, err
	}

	endTime, err := time.Parse(time.RFC3339, end)
	if err != nil {
		logger.APILogger.Error(err)
		return nil, 0, err
	}

	dbTransactions, total, err := s.repo.GetByDateRange(startTime, endTime, page, pageSize)
	if err != nil {
		logger.APILogger.Error(err)
		return nil, 0, err
	}

	transactions := make([]schemas.TransactionResponse, len(dbTransactions))
	for i, dbTransaction := range dbTransactions {
		transactions[i] = schemas.TransactionResponse{
			UserID:          dbTransaction.UserID.String(),
			TransactionDate: dbTransaction.TransactionDate,
			FromName:        dbTransaction.FromName,
			FromNumber:      dbTransaction.FromNumber,
			TransactionType: dbTransaction.TransactionType,
			BalanceBefore:   dbTransaction.BalanceBefore,
			BalanceAfter:    dbTransaction.BalanceAfter,
			Amount:          dbTransaction.Amount,
			ToNumber:        dbTransaction.ToNumber,
			ToName:          dbTransaction.ToName,
			Reference:       dbTransaction.Reference,
		}
	}

	return transactions, total, nil
}

func (s *TransactionService) ListTransactions(userID string,
	page, pageSize int,
) ([]schemas.TransactionResponse, int64, error) {
	dbTransactions, total, err := s.repo.List(uuid.MustParse(userID), page, pageSize)
	if err != nil {
		logger.APILogger.Error(err)
		return nil, 0, err
	}

	transactions := make([]schemas.TransactionResponse, len(dbTransactions))
	for i, dbTransaction := range dbTransactions {
		transactions[i] = schemas.TransactionResponse{
			UserID:          dbTransaction.UserID.String(),
			TransactionDate: dbTransaction.TransactionDate,
			FromName:        dbTransaction.FromName,
			FromNumber:      dbTransaction.FromNumber,
			TransactionType: dbTransaction.TransactionType,
			BalanceBefore:   dbTransaction.BalanceBefore,
			BalanceAfter:    dbTransaction.BalanceAfter,
			Amount:          dbTransaction.Amount,
			ToNumber:        dbTransaction.ToNumber,
			ToName:          dbTransaction.ToName,
			Reference:       dbTransaction.Reference,
		}
	}

	return transactions, total, nil
}

func (s *TransactionService) ListAllTransactions(userID string) ([]models.Transaction, error) {
	dbTransactions, err := s.repo.ListAll(uuid.MustParse(userID))
	if err != nil {
		logger.APILogger.Error(err)
		return nil, err
	}

	return dbTransactions, nil
}
