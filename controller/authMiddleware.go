package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/FlowRays/WhisperTrail/service"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.Set("isLoggedIn", false)
			c.Next()
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		claims, err := service.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("isLoggedIn", true)
		c.Set("userID", claims.UserID)

		c.Next()
	}
}
