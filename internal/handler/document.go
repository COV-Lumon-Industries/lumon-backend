package handler

import (
	"net/http"
	"strconv"
	"time"

	"lumon-backend/internal/config"
	"lumon-backend/internal/domain/models"
	"lumon-backend/internal/middleware"
	"lumon-backend/internal/ml"
	"lumon-backend/internal/service"
	"lumon-backend/pkg/common/response"

	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
	documentService *service.DocumentService
	cfg             *config.Config
}

func NewDocumentHandler(documentService *service.DocumentService, cfg *config.Config) *DocumentHandler {
	return &DocumentHandler{
		documentService: documentService,
		cfg:             cfg,
	}
}

func (h *DocumentHandler) RegisterRoutes(r *gin.RouterGroup) {
	document := r.Group("/documents")
	document.Use(middleware.JWTMiddleware(h.cfg))

	document = document.Group("", middleware.RequireRoles("common"))
	{
		document.POST("/", h.CreateDocument)
		document.GET("/item/:id", h.GetDocumentInformation)
		document.GET("/tag/:tag", h.GetDocumentInformationByTag)
		document.GET("/", h.ListDocuments)
		document.DELETE("/item/:id", h.DeleteDocument)
	}
}

func (h *DocumentHandler) CreateDocument(c *gin.Context) {
	req := struct {
		Url string `json:"url"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse(err.Error()))
		return
	}

	summary, err := ml.GetDocumentSummary(c, req.Url)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse(err.Error()))
		return
	}

	docType, err := ml.GetDocumentType(c, req.Url)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse(err.Error()))
		return
	}
	now := time.Now()

	document := &models.Document{
		ContentSummary: summary,
		Type:           docType,
		UploadedAt:     &now,
	}

	if err := h.documentService.CreateDocument(document); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewServerResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse(document))
}

func (h *DocumentHandler) GetDocumentInformation(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse("ID is invalid"))
		return
	}

	document, err := h.documentService.GetDocument(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(document))
}

func (h *DocumentHandler) GetDocumentInformationByTag(c *gin.Context) {
	tag := c.Param("tag")
	if tag == "" {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse("Tag is invalid"))
		return
	}

	document, err := h.documentService.GetDocumentByTypeTag(tag)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(document))
}

func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse("ID is invalid"))
		return
	}

	if err := h.documentService.DeleteDocument(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(gin.H{"message": "Document deleted successfully"}))
}

func (h *DocumentHandler) ListDocuments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	documents, total, err := h.documentService.ListDocuments(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		response.NewSuccessResponse(gin.H{
			"documents": documents,
			"meta": gin.H{
				"total":     total,
				"page":      page,
				"page_size": pageSize,
			},
		}),
	)
}
