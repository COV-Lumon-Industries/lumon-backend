package handler

import (
	"net/http"
	"strconv"

	"lumon-backend/internal/config"
	"lumon-backend/internal/domain/models"
	"lumon-backend/internal/middleware"
	"lumon-backend/internal/service"
	"lumon-backend/pkg/common/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AccountHandler struct {
	accountService *service.AccountService
	cfg            *config.Config
}

func NewAccountHandler(accountService *service.AccountService, cfg *config.Config) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
		cfg:            cfg,
	}
}

func (h *AccountHandler) RegisterRoutes(r *gin.RouterGroup) {
	account := r.Group("/accounts")
	account.Use(middleware.JWTMiddleware(h.cfg))

	account = account.Group("", middleware.RequireRoles("common"))
	{
		account.POST("/", h.CreateAccount)
		account.GET("/:id", h.GetAccount)
		account.GET("/type/:type", h.GetAccountByType)
		account.PUT("/:id", h.UpdateAccount)
		account.DELETE("/:id", h.DeleteAccount)
		account.GET("/", h.ListAccounts)
	}
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	req := struct {
		Name       string `json:"name"`
		Type       string `json:"type"`
		DocumentID string `json:"document_id"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse(err.Error()))
		return
	}

	account := models.Account{
		Name:       req.Name,
		Type:       req.Type,
		DocumentID: uuid.MustParse(req.DocumentID),
	}

	if err := h.accountService.CreateAccount(&account); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewServerResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse(account))
}

func (h *AccountHandler) GetAccount(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse("ID is invalid"))
		return
	}

	account, err := h.accountService.GetAccount(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(account))
}

func (h *AccountHandler) GetAccountByType(c *gin.Context) {
	accountType := c.Param("type")
	if accountType == "" {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse("Type is invalid"))
		return
	}

	account, err := h.accountService.GetAccountByType(accountType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(account))
}

func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse("ID is invalid"))
		return
	}

	var account models.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse(err.Error()))
		return
	}

	account.ID = uuid.MustParse(id)

	if err := h.accountService.UpdateAccount(&account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(account))
}

func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse("ID is invalid"))
		return
	}

	if err := h.accountService.DeleteAccount(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(gin.H{"message": "Account deleted successfully"}))
}

func (h *AccountHandler) ListAccounts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	accounts, total, err := h.accountService.ListAccounts(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		response.NewSuccessResponse(gin.H{
			"accounts": accounts,
			"meta": gin.H{
				"total":     total,
				"page":      page,
				"page_size": pageSize,
			},
		}),
	)
}
