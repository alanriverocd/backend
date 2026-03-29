package middleware

import (
	"net/http"
	"strings"

	"github.com/finatiol/backend/config"
	"github.com/gin-gonic/gin"
)

const userUIDKey = "userUID"

// AuthRequired valida el Firebase ID token en el header Authorization.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Header Authorization requerido"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Formato de token invalido. Usar: Bearer <token>"})
			return
		}

		idToken := parts[1]
		token, err := config.AuthClient.VerifyIDToken(c.Request.Context(), idToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token invalido o expirado"})
			return
		}

		c.Set(userUIDKey, token.UID)
		c.Next()
	}
}

// GetUserUID obtiene el UID del usuario autenticado del contexto Gin.
func GetUserUID(c *gin.Context) (string, bool) {
	uid, exists := c.Get(userUIDKey)
	if !exists {
		return "", false
	}
	return uid.(string), true
}
