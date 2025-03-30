package middleware

import (
	"net/http"
	"strings"

	"lumon-backend/internal/config"
	"lumon-backend/pkg/common/logger"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWTMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			logger.APILogger.Error("No Authorization header provided")
			respondWithError(c, http.StatusUnauthorized, "User Not Authorized to perform action")
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(tokenString, bearerPrefix) {
			logger.APILogger.Error("Invalid token format: missing Bearer prefix")
			respondWithError(c, http.StatusUnauthorized, "Invalid token format")
			return
		}

		tokenString = strings.TrimPrefix(tokenString, bearerPrefix)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.SecretKey), nil
		})

		if err != nil || !token.Valid {
			logger.APILogger.Errorf("Invalid token: %v", err)
			respondWithError(c, http.StatusUnauthorized, "User Not Authorized to perform action")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			logger.APILogger.Error("Invalid token claims")
			respondWithError(c, http.StatusUnauthorized, "User Not Authorized to perform action")
			return
		}

		userRole, ok := claims["user_role"].(string)
		if !ok {
			logger.APILogger.Error("Missing or invalid user_role in token")
			respondWithError(c, http.StatusUnauthorized, "User Not Authorized to perform action")
			return
		}

		requiredRoles := c.GetStringSlice("requiredRoles")
		if len(requiredRoles) > 0 && !containsRole(requiredRoles, userRole) {
			logger.APILogger.Errorf("User role %s not in required roles: %v", userRole, requiredRoles)
			respondWithError(c, http.StatusForbidden, "User is forbidden to perform action")
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}

func respondWithError(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, gin.H{
		"message": message,
	})
}

func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("requiredRoles", roles)
		c.Next()
	}
}

func containsRole(roles []string, role string) bool {
	for _, r := range roles {
		if strings.EqualFold(r, role) {
			return true
		}
	}
	return false
}
