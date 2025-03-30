package handler

import (
	"net/http"
	"strconv"

	"lumon-backend/internal/config"
	"lumon-backend/internal/domain/models"
	"lumon-backend/internal/domain/schemas"
	"lumon-backend/internal/middleware"
	"lumon-backend/internal/service"
	"lumon-backend/pkg/common/logger"
	"lumon-backend/pkg/common/response"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LoanRequestHandler struct {
	loanRequestService *service.LoanRequestService
	userService        *service.UserService
	cfg                *config.Config
}

func NewLoanRequestHandler(loanRequestService *service.LoanRequestService, userService *service.UserService, cfg *config.Config) *LoanRequestHandler {
	return &LoanRequestHandler{
		loanRequestService: loanRequestService,
		userService:        userService,
		cfg:                cfg,
	}
}

func (h *LoanRequestHandler) RegisterRoutes(r *gin.RouterGroup) {
	loanRequests := r.Group("/loan-requests")
	loanRequests.Use(middleware.JWTMiddleware(h.cfg))

	loanRequests = loanRequests.Group("", middleware.RequireRoles("common"))
	{
		loanRequests.POST("", h.CreateLoanRequest)
		loanRequests.GET("/:id", h.GetLoanRequest)
		loanRequests.GET("/borrower/:borrower_id", h.GetLoanRequestsByBorrower)
		loanRequests.PUT("/:id", h.UpdateLoanRequest)
		loanRequests.DELETE("/:id", h.DeleteLoanRequest)
		loanRequests.GET("", h.ListLoanRequests)
		loanRequests.PATCH("/:id/status", h.UpdateLoanRequestStatus)
	}
}

func (h *LoanRequestHandler) CreateLoanRequest(c *gin.Context) {
	var request schemas.CreateLoanRequestDetails

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.APILogger.Error("Failed to bind JSON in CreateLoanRequest:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.APILogger.Error("Invalid user ID format in CreateLoanRequest:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	loanRequest := &models.LoanRequest{
		Amount:       request.Amount,
		BorrowerID:   userID,
		InterestRate: request.InterestRate,
		LoanDuration: request.LoanDuration,
		Purpose:      request.Purpose,
	}

	if err := h.loanRequestService.CreateLoanRequest(loanRequest); err != nil {
		logger.APILogger.Error("Failed to create loan request in CreateLoanRequest:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create loan request"})
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse("Loan Request Created"))
}

func (h *LoanRequestHandler) GetLoanRequest(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		logger.APILogger.Error("Invalid loan request ID in GetLoanRequest:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan request ID"})
		return
	}

	loanRequest, err := h.loanRequestService.GetLoanRequest(id)
	if err != nil {
		logger.APILogger.Error("Loan request not found in GetLoanRequest:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan request not found"})
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(gin.H{"loan_request": loanRequest}))
}

func (h *LoanRequestHandler) GetLoanRequestsByBorrower(c *gin.Context) {
	borrowerID := c.Param("borrower_id")
	if _, err := uuid.Parse(borrowerID); err != nil {
		logger.APILogger.Error("Invalid borrower ID in GetLoanRequestsByBorrower:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid borrower ID"})
		return
	}

	loanRequests, err := h.loanRequestService.GetLoanRequestsByBorrower(borrowerID)
	if err != nil {
		logger.APILogger.Error("Failed to retrieve loan requests in GetLoanRequestsByBorrower:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve loan requests"})
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(gin.H{"loan_requests": loanRequests}))
}

func (h *LoanRequestHandler) UpdateLoanRequest(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		logger.APILogger.Error("Invalid loan request ID in UpdateLoanRequest:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan request ID"})
		return
	}

	var request schemas.CreateLoanRequestDetails
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.APILogger.Error("Failed to bind JSON in UpdateLoanRequest:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lq, err := h.loanRequestService.GetLoanRequest(id)
	if err != nil {
		logger.APILogger.Error("Loan Request Not Found:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan Request not found"})
		return
	}

	if request.Amount != lq.Amount {
		lq.Amount = request.Amount
	}

	if request.InterestRate != lq.InterestRate {
		lq.InterestRate = request.InterestRate
	}

	if request.LoanDuration != lq.LoanDuration {
		lq.LoanDuration = request.LoanDuration
	}

	if request.Purpose != lq.Purpose {
		lq.Purpose = request.Purpose
	}

	if err := h.loanRequestService.UpdateLoanRequest(lq); err != nil {
		logger.APILogger.Error("Failed to update loan request in UpdateLoanRequest:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update loan request"})
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse("Loan request updated successfully"))
}

func (h *LoanRequestHandler) DeleteLoanRequest(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		logger.APILogger.Error("Invalid loan request ID in DeleteLoanRequest:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan request ID"})
		return
	}

	if err := h.loanRequestService.DeleteLoanRequest(id); err != nil {
		logger.APILogger.Error("Failed to delete loan request in DeleteLoanRequest:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete loan request"})
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse("Loan request deleted successfully"))
}

func (h *LoanRequestHandler) ListLoanRequests(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		logger.APILogger.Error("Invalid page number in ListLoanRequests:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		logger.APILogger.Error("Invalid page size in ListLoanRequests:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}

	loanRequests, total, err := h.loanRequestService.ListLoanRequests(page, pageSize)
	if err != nil {
		logger.APILogger.Error("Failed to list loan requests in ListLoanRequests:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list loan requests"})
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(gin.H{
		"loan_requests": loanRequests,
		"total":         total,
		"page":          page,
		"page_size":     pageSize,
	}))
}

func (h *LoanRequestHandler) UpdateLoanRequestStatus(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		logger.APILogger.Error("Invalid loan request ID in UpdateLoanRequestStatus:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan request ID"})
		return
	}

	var statusReq struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&statusReq); err != nil {
		logger.APILogger.Error("Failed to bind JSON in UpdateLoanRequestStatus:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.loanRequestService.UpdateLoanRequestStatus(id, statusReq.Status); err != nil {
		logger.APILogger.Error("Failed to update loan request status in UpdateLoanRequestStatus:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update loan request status"})
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse("Loan request status updated successfully"))
}
