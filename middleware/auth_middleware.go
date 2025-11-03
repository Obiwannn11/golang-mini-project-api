package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"rakamin-evermos/utils" 

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token from header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Token authentication not provided")
			c.Abort() 
			return
		}

		// format Bearer tokenString
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		if tokenString == "" {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Token authentication not valid")
			c.Abort()
			return
		}

		token, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.SendErrorResponse(c, http.StatusUnauthorized, fmt.Sprintf("Token not valid: %s", err.Error()))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Token authentication not valid")
			c.Abort()
			return
		}

		// Set user information to context gin
		userID := uint(claims["user_id"].(float64)) // JWT number is float64
		isAdmin := claims["is_admin"].(bool)

		c.Set("currentUserID", userID)
		c.Set("currentUserIsAdmin", isAdmin)

		c.Next()
	}
}

// AdminOnlyMiddleware verifies if user is admin
func AdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ensure AuthMiddleware has run before to set currentUserIsAdmin
		isAdmin, exists := c.Get("currentUserIsAdmin")
		if !exists || !isAdmin.(bool) {
			utils.SendErrorResponse(c, http.StatusForbidden, "Access denied: Only Admins are allowed")
			c.Abort()
			return
		}

		c.Next()
	}
}