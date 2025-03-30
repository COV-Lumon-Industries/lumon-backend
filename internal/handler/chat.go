package handler

import (
	"context"
	"io"
	"iter"
	"net/http"

	"lumon-backend/internal/config"
	"lumon-backend/internal/middleware"
	"lumon-backend/internal/ml"
	"lumon-backend/internal/service"
	"lumon-backend/pkg/common/response"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	transactionService *service.TransactionService
	cfg                *config.Config
}

func NewChatHandler(transactionService *service.TransactionService, cfg *config.Config) *ChatHandler {
	return &ChatHandler{
		transactionService: transactionService,
		cfg:                cfg,
	}
}

func (h *ChatHandler) RegisterRoutes(r *gin.RouterGroup) {
	chat := r.Group("/chat")
	chat.Use(middleware.JWTMiddleware(h.cfg))

	chat = chat.Group("", middleware.RequireRoles("common"))
	{
		chat.POST("/static", h.CreateChatResponse)
		chat.POST("/stream", h.CreateChatStream)
		chat.POST("/search", h.CreateSearch)
	}
}

func (h *ChatHandler) CreateChatResponse(c *gin.Context) {
	req := struct {
		Prompt string `json:"prompt"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse(err.Error()))
		return
	}

	chat, err := ml.GetChatResponse(c.Request.Context(), req.Prompt)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse(ml.ResponseToPartString(chat)))
}

func (h *ChatHandler) CreateSearch(c *gin.Context) {
	req := struct {
		Prompt string `json:"prompt"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse(err.Error()))
		return
	}

	search, err := ml.GetSearchResponse(c.Request.Context(), req.Prompt)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.NewSuccessResponse(ml.ResponseToPartString(search)))
}

func (h *ChatHandler) CreateChatStream(c *gin.Context) {
	req := struct {
		Prompt string `json:"prompt"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse(err.Error()))
		return
	}

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	stream, err := ml.GetChatResponseStream(ctx, req.Prompt)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse(err.Error()))
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	pull, stop := iter.Pull2(stream)
	defer stop()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-ctx.Done():
			c.SSEvent("message", "Client Disconnected")
			return false
		default:
			resp, err, ok := pull()

			if !ok {
				c.SSEvent("message", "Closed By API")
				return false
			}

			if err != nil {
				c.SSEvent("message", err)
				return false
			}

			if resp != nil && len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
				c.SSEvent("message", ml.ResponseToPartString(resp))
				return true
			}

			return false
		}
	})
}
