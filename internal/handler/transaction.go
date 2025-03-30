package handler

import (
	"net/http"
	"strconv"

	"lumon-backend/internal/config"
	"lumon-backend/internal/middleware"
	"lumon-backend/internal/ml"
	"lumon-backend/internal/service"
	"lumon-backend/pkg/common/logger"
	"lumon-backend/pkg/common/response"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService *service.TransactionService
	userService        *service.UserService
	cfg                *config.Config
}

func NewTransactionHandler(transactionService *service.TransactionService, userService *service.UserService, cfg *config.Config) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		userService:        userService,
		cfg:                cfg,
	}
}

func (h *TransactionHandler) RegisterRoutes(r *gin.RouterGroup) {
	transaction := r.Group("/transactions")
	transaction.Use(middleware.JWTMiddleware(h.cfg))

	transaction = transaction.Group("", middleware.RequireRoles("common"))
	{
		transaction.POST("", h.CreateTransactions)
		transaction.GET("/credit", h.CreateCreditScore)
		transaction.GET("/item/:id", h.GetTransaction)
		transaction.GET("/type", h.ListTransactionsByType)
		transaction.GET("/account", h.ListTransactionsByAccountNumber)
		transaction.GET("/", h.ListDocuments)
	}
}

func (h *TransactionHandler) CreateTransactions(c *gin.Context) {
	req := struct {
		Url string `json:"url"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse(err.Error()))
		return
	}

	transactions, batchSize, err := ml.GetMoMoTransactions(c, req.Url)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse(err.Error()))
		return
	}

	user, exists := c.Get("user")
	if !exists {
		logger.APILogger.Error("Unauthorized request in CreateLoanRequest")
		c.JSON(http.StatusUnauthorized, response.NewFailureResponse("Unauthorized Request"))
		return
	}

	claims, ok := user.(jwt.MapClaims)
	if !ok {
		logger.APILogger.Error("Failed to parse user claims in CreateLoanRequest")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user claims"})
		return
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		logger.APILogger.Error("User ID not found in token in CreateLoanRequest")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in token"})
		return
	}

	if err := h.transactionService.CreateTransactionBatch(transactions, batchSize, userIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewServerResponse(err.Error()))
		return
	}

	c.JSON(
		http.StatusOK,
		response.NewSuccessResponse(gin.H{
			"created_transactions": transactions,
			"batch_size":           batchSize,
		}),
	)
}

func (h *TransactionHandler) CreateCreditScore(c *gin.Context) {
	userClaims, exists := c.Get("user")
	if !exists {
		logger.APILogger.Error("Unauthorized request in CreateLoanRequest")
		c.JSON(http.StatusUnauthorized, response.NewFailureResponse("Unauthorized Request"))
		return
	}

	claims, ok := userClaims.(jwt.MapClaims)
	if !ok {
		logger.APILogger.Error("Failed to parse user claims in CreateLoanRequest")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user claims"})
		return
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		logger.APILogger.Error("User ID not found in token in CreateLoanRequest")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in token"})
		return
	}

	transactions, err := h.transactionService.ListAllTransactions(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	dbUser, err := h.userService.GetUserModel(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var calculator *service.CreditScoreCalculator = service.NewCreditScoreCalculator(transactions)
	creditScore := calculator.Calculate()

	dbUser.CreditScore = int64(creditScore)

	if err := h.userService.UpdateUser(dbUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		response.NewSuccessResponse(gin.H{
			"credit_score": creditScore,
		}),
	)
}

func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse("ID is invalid"))
		return
	}

	transaction, err := h.transactionService.GetTransaction(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(transaction))
}

func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse("ID is invalid"))
		return
	}

	if err := h.transactionService.DeleteTransaction(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(gin.H{"message": "Transaction deleted successfully"}))
}

func (h *TransactionHandler) ListDocuments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	userClaims, exists := c.Get("user")
	if !exists {
		logger.APILogger.Error("Unauthorized request in ListDocumnents")
		c.JSON(http.StatusUnauthorized, response.NewFailureResponse("Unauthorized Request"))
		return
	}

	claims, ok := userClaims.(jwt.MapClaims)
	if !ok {
		logger.APILogger.Error("Failed to parse user claims in ListDocumnents")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user claims"})
		return
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		logger.APILogger.Error("User ID not found in token in ListDocumnents")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in token"})
		return
	}

	transactions, total, err := h.transactionService.ListTransactions(userIDStr, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		response.NewSuccessResponse(gin.H{
			"transactions": transactions,
			"meta": gin.H{
				"total":     total,
				"page":      page,
				"page_size": pageSize,
			},
		}),
	)
}

func (h *TransactionHandler) ListTransactionsByAccountNumber(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	req := struct {
		AccountNumber string `json:"account_number"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse(err.Error()))
		return
	}

	transactions, total, err := h.transactionService.GetTransactionsByAccountNumber(req.AccountNumber, page, pageSize)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transactions not found"})
		return
	}

	c.JSON(
		http.StatusOK,
		response.NewSuccessResponse(gin.H{
			"transactions": transactions,
			"meta": gin.H{
				"total":     total,
				"page":      page,
				"page_size": pageSize,
			},
		}),
	)
}

func (h *TransactionHandler) ListTransactionsByType(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	req := struct {
		Type string `json:"type"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse(err.Error()))
		return
	}

	transactions, total, err := h.transactionService.GetTransactionsByType(req.Type, page, pageSize)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transactions not found"})
		return
	}

	c.JSON(
		http.StatusOK,
		response.NewSuccessResponse(gin.H{
			"transactions": transactions,
			"meta": gin.H{
				"total":     total,
				"page":      page,
				"page_size": pageSize,
			},
		}),
	)
}
