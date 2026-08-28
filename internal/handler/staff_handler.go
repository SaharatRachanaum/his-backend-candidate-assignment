package handler

import (
	"errors"
	"net/http"

	"github.com/agnos/hospital-middleware/internal/domain"
	"github.com/agnos/hospital-middleware/internal/service"
	"github.com/agnos/hospital-middleware/pkg/response"
	"github.com/gin-gonic/gin"
)

// StaffHandler handles staff HTTP requests
type StaffHandler struct {
	staffService service.StaffService
}

// NewStaffHandler creates a new StaffHandler instance
func NewStaffHandler(staffService service.StaffService) *StaffHandler {
	return &StaffHandler{staffService: staffService}
}

// CreateStaff handles staff registration
func (h *StaffHandler) CreateStaff(c *gin.Context) {
	var req domain.CreateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	res, err := h.staffService.Create(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrUsernameAlreadyExists) {
			response.Conflict(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "Staff created successfully", res)
}

// LoginStaff handles staff authentication and token issuance
func (h *StaffHandler) LoginStaff(c *gin.Context) {
	var req domain.LoginStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	res, err := h.staffService.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) || errors.Is(err, service.ErrHospitalMismatch) {
			response.Unauthorized(c, err.Error())
			return
		}
		response.InternalServerError(c, "Login failed", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login successful", res)
}
