package middleware

import (
	"strings"

	"github.com/agnos/hospital-middleware/pkg/response"
	"github.com/agnos/hospital-middleware/pkg/token"
	"github.com/gin-gonic/gin"
)

const (
	ContextKeyStaffID  = "staff_id"
	ContextKeyUsername = "username"
	ContextKeyHospital = "hospital"
)

// AuthMiddleware ensures requests have a valid Bearer JWT token
func AuthMiddleware(jwtManager *token.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Unauthorized(c, "Invalid authorization format. Format: Bearer <token>")
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(parts[1])
		if tokenString == "" {
			response.Unauthorized(c, "Token is missing")
			c.Abort()
			return
		}

		claims, err := jwtManager.Verify(tokenString)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Inject authenticated staff claims into gin Context
		c.Set(ContextKeyStaffID, claims.StaffID)
		c.Set(ContextKeyUsername, claims.Username)
		c.Set(ContextKeyHospital, claims.Hospital)

		c.Next()
	}
}

// GetStaffHospital retrieves the hospital ID of the authenticated staff from context
func GetStaffHospital(c *gin.Context) string {
	val, exists := c.Get(ContextKeyHospital)
	if !exists {
		return ""
	}
	hospital, ok := val.(string)
	if !ok {
		return ""
	}
	return hospital
}

// GetStaffID retrieves the staff ID from context
func GetStaffID(c *gin.Context) string {
	val, exists := c.Get(ContextKeyStaffID)
	if !exists {
		return ""
	}
	id, ok := val.(string)
	if !ok {
		return ""
	}
	return id
}
