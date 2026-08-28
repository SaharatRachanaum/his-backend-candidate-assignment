package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/agnos/hospital-middleware/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtManager := token.NewJWTManager("test_secret_for_middleware_12345", 1*time.Hour)

	validToken, _, _ := jwtManager.Generate("staff-uuid-1", "doctor_somchai", "hospital-a")

	setupRouter := func() *gin.Engine {
		r := gin.New()
		r.Use(AuthMiddleware(jwtManager))
		r.GET("/protected", func(c *gin.Context) {
			hospital := GetStaffHospital(c)
			staffID := GetStaffID(c)
			c.JSON(http.StatusOK, gin.H{
				"hospital": hospital,
				"staff_id": staffID,
			})
		})
		return r
	}

	r := setupRouter()

	// 1. Positive: Valid Bearer token
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "hospital-a")
	assert.Contains(t, w.Body.String(), "staff-uuid-1")

	// 2. Negative: Missing Authorization Header
	reqNoAuth, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	wNoAuth := httptest.NewRecorder()
	r.ServeHTTP(wNoAuth, reqNoAuth)
	assert.Equal(t, http.StatusUnauthorized, wNoAuth.Code)

	// 3. Negative: Invalid Token Format
	reqBadFormat, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	reqBadFormat.Header.Set("Authorization", "Basic 12345")
	wBadFormat := httptest.NewRecorder()
	r.ServeHTTP(wBadFormat, reqBadFormat)
	assert.Equal(t, http.StatusUnauthorized, wBadFormat.Code)

	// 4. Negative: Invalid / Tampered Token
	reqInvalidToken, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	reqInvalidToken.Header.Set("Authorization", "Bearer invalid.token.value")
	wInvalidToken := httptest.NewRecorder()
	r.ServeHTTP(wInvalidToken, reqInvalidToken)
	assert.Equal(t, http.StatusUnauthorized, wInvalidToken.Code)
}
