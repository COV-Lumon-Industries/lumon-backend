package handler

import (
	"net/http"
	"strings"
	"time"

	"lumon-backend/internal/config"
	"lumon-backend/internal/domain/models"
	"lumon-backend/internal/domain/schemas"
	"lumon-backend/internal/service"
	"lumon-backend/pkg/common/response"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userService *service.UserService
	cfg         *config.Config
}

func NewAuthHandler(userService *service.UserService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		cfg:         cfg,
	}
}

func (h *AuthHandler) RegisterRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)

		// auth.POST("/forgot-password", h.ForgotPassword)
		// auth.POST("/reset-password", h.ResetPassword)

		// auth.GET("/verify-email", h.VerifyEmail)
		// auth.POST("/resend-verification", h.ResendVerification)
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var request schemas.CreateAccountDetails

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{
		Username:    request.Username,
		Password:    request.Password,
		Email:       request.Email,
		PhoneNumber: request.PhoneNumber,
		UserRole:    "common",
	}

	if err := h.userService.CreateUser(user); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			c.JSON(http.StatusConflict, response.NewFailureResponse("Username or email already exists"))
			return
		}
		c.JSON(http.StatusInternalServerError, response.NewServerResponse(err.Error()))
		return
	}

	userSchema := schemas.UserResponse{
		ID:       user.ID,
		UserRole: user.UserRole,
	}
	token, err := h.generateJWTToken(&userSchema)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewServerResponse(err.Error()))
		return
	}

	c.SetCookie("Authorization", token, time.Now().Add(h.cfg.TokenDuration).Hour(), "/", "localhost", false, true)

	// _, err := h.emailService.SendVerificationEmail(user.Email, user.Username, user.VerificationToken)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, response.NewFailureResponse("Failed to send verification email"))
	// 	return
	// }

	user.Password = ""
	c.JSON(http.StatusCreated, response.NewSuccessResponse(gin.H{"token": token, "user": user}))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var credentials schemas.Credentials

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.GetUserByUsername(credentials.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewFailureResponse("Invalid Username or Password"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewFailureResponse("Invalid Username or Password"))
		return
	}

	// if !user.EmailVerified {
	// 	c.JSON(http.StatusUnauthorized, response.NewFailureResponse("User account not verified. Please verify your Account"))
	// 	return
	// }

	token, err := h.generateJWTToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewServerResponse(err.Error()))
		return
	}

	c.SetCookie("Authorization", token, time.Now().Add(h.cfg.TokenDuration).Hour(), "/", "localhost", false, true)

	c.JSON(
		http.StatusOK,
		response.NewSuccessResponse(gin.H{
			"token": token,
			"user":  user,
		}),
	)
}

func (h *AuthHandler) generateJWTToken(user *schemas.UserResponse) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"user_role": user.UserRole,
		"exp":       time.Now().Add(h.cfg.TokenDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.cfg.SecretKey)
}

// func (h *AuthHandler) ForgotPassword(c *gin.Context) {
// 	var req schemas.ForgotPasswordRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, response.NewFailureResponse(err.Error()))
// 		return
// 	}

// 	user, err := h.userService.GetUserByEmail(req.Email)
// 	if err != nil {
// 		c.JSON(
// 			http.StatusOK,
// 			response.NewSuccessResponse(
// 				gin.H{"message": "If your email is registered, you will receive reset instructions"},
// 			),
// 		)
// 		return
// 	}

// 	resetToken, err := h.userService.RequestPasswordReset(req.Email)
// 	if err != nil {
// 		c.JSON(
// 			http.StatusInternalServerError,
// 			response.NewServerResponse(
// 				gin.H{"error": "Failed to process request"},
// 			),
// 		)
// 		return
// 	}

// 	res, err := h.emailService.SendPasswordResetEmail(
// 		req.Email,
// 		resetToken,
// 		user.Username,
// 	)
// 	if err != nil {
// 		c.JSON(
// 			http.StatusInternalServerError,
// 			response.NewServerResponse(
// 				gin.H{"error": "Failed to send reset email"},
// 			),
// 		)
// 		return
// 	}

// 	c.JSON(
// 		http.StatusOK,
// 		response.NewSuccessResponse(
// 			gin.H{
// 				"message":  "If your email is registered, you will receive reset instructions",
// 				"email_id": res,
// 			},
// 		),
// 	)
// }

// func (h *AuthHandler) ResetPassword(c *gin.Context) {
// 	var req schemas.ResetPasswordRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, response.NewFailureResponse("Request parameter with likely missing Field"))
// 		return
// 	}

// 	if err := h.userService.ResetPassword(req.Token, req.NewPassword); err != nil {
// 		c.JSON(http.StatusBadRequest, response.NewFailureResponse("Invalid or expired reset token"))
// 		return
// 	}

// 	c.JSON(http.StatusOK, response.NewSuccessResponse(gin.H{"message": "Password successfully reset"}))
// }

// func (h *AuthHandler) VerifyEmail(c *gin.Context) {
// 	token := c.Query("token")
// 	if token == "" {
// 		c.JSON(http.StatusBadRequest, response.NewFailureResponse("verification token is required"))
// 		return
// 	}

// 	if err := h.userService.VerifyEmail(token); err != nil {
// 		c.JSON(http.StatusBadRequest, response.NewFailureResponse("Invalid or expired verification token"))
// 		return
// 	}

// 	c.JSON(http.StatusOK, response.NewSuccessResponse(gin.H{"message": "Email successfully verified"}))
// }

// func (h *AuthHandler) ResendVerification(c *gin.Context) {
// 	var req schemas.VerificationDetails

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, response.NewFailureResponse("Request parameter with likely missing Field"))
// 		return
// 	}

// 	token, err := h.userService.ResendVerification(req.Email)
// 	if err != nil {
// 		fmt.Println(err)
// 		c.JSON(http.StatusOK,
// 			response.NewSuccessResponse(
// 				gin.H{
// 					"message": "If your email is registered, you will receive verification instructions",
// 				},
// 			),
// 		)
// 		return
// 	}

// 	user, _ := h.userService.GetUserByEmail(req.Email)
// 	res, err := h.emailService.SendVerificationEmail(req.Email, user.Username, token)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, response.NewServerResponse(gin.H{"error": "Failed to send verification email"}))
// 		return
// 	}

// 	c.JSON(http.StatusOK, response.NewSuccessResponse(gin.H{"message": "Verification email has been sent", "email_id": res}))
// }
