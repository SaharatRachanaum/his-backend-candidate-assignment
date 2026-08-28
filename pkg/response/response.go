package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Count   *int        `json:"count,omitempty"`
}

// Success sends a success JSON response
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SuccessWithCount sends a success JSON response including item count
func SuccessWithCount(c *gin.Context, statusCode int, message string, data interface{}, count int) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
		Count:   &count,
	})
}

// Error sends an error JSON response
func Error(c *gin.Context, statusCode int, message string, err interface{}) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// BadRequest sends a 400 Bad Request error
func BadRequest(c *gin.Context, message string, err interface{}) {
	Error(c, http.StatusBadRequest, message, err)
}

// Unauthorized sends a 401 Unauthorized error
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message, nil)
}

// Forbidden sends a 403 Forbidden error
func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message, nil)
}

// NotFound sends a 404 Not Found error
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message, nil)
}

// Conflict sends a 409 Conflict error
func Conflict(c *gin.Context, message string) {
	Error(c, http.StatusConflict, message, nil)
}

// InternalServerError sends a 500 Internal Server Error
func InternalServerError(c *gin.Context, message string, err interface{}) {
	Error(c, http.StatusInternalServerError, message, err)
}
