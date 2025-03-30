package handler

import (
	"net/http"

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

type WalletsHandler struct {
	walletService *service.WalletService
	cfg           *config.Config
}

func NewWalletsHandler(walletService *service.WalletService, cfg *config.Config) *WalletsHandler {
	return &WalletsHandler{
		walletService: walletService,
		cfg:           cfg,
	}
}

func (h *WalletsHandler) RegisterRoutes(r *gin.RouterGroup) {
	wallets := r.Group("/wallets")
	wallets.Use(middleware.JWTMiddleware(h.cfg))
	wallets.Use(middleware.RequireRoles("common"))

	wallets.POST("/topup", h.TopUpWallet)
	wallets.POST("/withdraw", h.WithdrawWallet)
}

func (h *WalletsHandler) TopUpWallet(c *gin.Context) {
	var request schemas.TopUpWalletDetails

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.APILogger.Errorf("Failed to bind JSON in TopUpWallet: %v", err)
		c.JSON(http.StatusBadRequest, response.NewFailureResponse("Invalid request format"))
		return
	}

	user, exists := c.Get("user")
	if !exists {
		logger.APILogger.Error("Unauthorized request in TopUpWallet")
		c.JSON(http.StatusUnauthorized, response.NewFailureResponse("Unauthorized request"))
		return
	}

	claims, ok := user.(jwt.MapClaims)
	if !ok {
		logger.APILogger.Error("Failed to parse user claims in TopUpWallet")
		c.JSON(http.StatusInternalServerError, response.NewFailureResponse("Invalid token format"))
		return
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		logger.APILogger.Error("User ID not found in token in TopUpWallet")
		c.JSON(http.StatusInternalServerError, response.NewFailureResponse("Invalid user ID in token"))
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.APILogger.Errorf("Invalid user ID format in TopUpWallet: %v", err)
		c.JSON(http.StatusBadRequest, response.NewFailureResponse("Invalid user ID format"))
		return
	}

	if request.Amount <= 0 {
		logger.APILogger.Error("Invalid amount in TopUpWallet")
		c.JSON(http.StatusBadRequest, response.NewFailureResponse("Amount must be positive"))
		return
	}

	wallet := &models.Wallet{
		Balance: request.Amount,
		UserID:  userID,
	}

	if err := h.walletService.TopUpAccount(wallet); err != nil {
		logger.APILogger.Errorf("Failed to top up wallet: %v", err)
		c.JSON(http.StatusInternalServerError, response.NewFailureResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse("Wallet topped up successfully"))
}

func (h *WalletsHandler) WithdrawWallet(c *gin.Context) {
	var request schemas.TopUpWalletDetails

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.APILogger.Errorf("Failed to bind JSON in WithdrawWallet: %v", err)
		c.JSON(http.StatusBadRequest, response.NewFailureResponse("Invalid request format"))
		return
	}

	user, exists := c.Get("user")
	if !exists {
		logger.APILogger.Error("Unauthorized request in WithdrawWallet")
		c.JSON(http.StatusUnauthorized, response.NewFailureResponse("Unauthorized request"))
		return
	}

	claims, ok := user.(jwt.MapClaims)
	if !ok {
		logger.APILogger.Error("Failed to parse user claims in WithdrawWallet")
		c.JSON(http.StatusInternalServerError, response.NewFailureResponse("Invalid token format"))
		return
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		logger.APILogger.Error("User ID not found in token in WithdrawWallet")
		c.JSON(http.StatusInternalServerError, response.NewFailureResponse("Invalid user ID in token"))
		return
	}

	if request.Amount <= 0 {
		logger.APILogger.Error("Invalid amount in WithdrawWallet")
		c.JSON(http.StatusBadRequest, response.NewFailureResponse("Amount must be positive"))
		return
	}

	if err := h.walletService.WithdrawAccount(userIDStr, request.Amount); err != nil {
		logger.APILogger.Errorf("Failed to withdraw from wallet: %v", err)
		c.JSON(http.StatusBadRequest, response.NewFailureResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse("Withdrawal processed successfully"))
}