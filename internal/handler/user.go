package handler

import (
	"net/http"
	"strconv"

	"lumon-backend/internal/config"
	"lumon-backend/internal/domain/models"
	"lumon-backend/internal/domain/schemas"
	"lumon-backend/internal/middleware"
	"lumon-backend/internal/service"
	"lumon-backend/pkg/common/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userService *service.UserService
	cfg         *config.Config
}

func NewUserHandler(userService *service.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{
		userService: userService,
		cfg:         cfg,
	}
}

func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	users.Use(middleware.JWTMiddleware(h.cfg))

	admins := users.Group("", middleware.RequireRoles("admin"))
	{
		admins.DELETE("/:id", h.DeleteUser)
	}

	user := users.Group("", middleware.RequireRoles("common"))
	{
		user.GET("/:id", h.GetUser)
		user.PUT("/:id", h.UpdateUser)
		user.GET("/", h.ListUsers)
	}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse("ID is invalid"))
		return
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(user))
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse("ID is invalid"))
		return
	}

	var user_update schemas.UserUpdate
	if err := c.ShouldBindJSON(&user_update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user_update.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewServerResponse("Could not hash password."))
		return
	}

	user := models.User{
		ID:       uuid.MustParse(id),
		Username: user_update.Username,
		Password: string(hashedPassword),
	}

	if err := h.userService.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(user))
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, response.NewFailureResponse("ID is invalid"))
		return
	}

	if err := h.userService.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(gin.H{"message": "User deleted successfully"}))
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, err := h.userService.ListUsers(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(
		http.StatusOK,
		response.NewSuccessResponse(gin.H{
			"users": users,
			"meta": gin.H{
				"total":     total,
				"page":      page,
				"page_size": pageSize,
			},
		}),
	)
}
