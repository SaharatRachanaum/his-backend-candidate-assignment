package handler

import (
	"net/http"

	"github.com/agnos/hospital-middleware/internal/domain"
	"github.com/agnos/hospital-middleware/internal/middleware"
	"github.com/agnos/hospital-middleware/internal/service"
	"github.com/agnos/hospital-middleware/pkg/response"
	"github.com/gin-gonic/gin"
)

// PatientHandler handles patient HTTP requests
type PatientHandler struct {
	patientService service.PatientService
}

// NewPatientHandler creates a new PatientHandler instance
func NewPatientHandler(patientService service.PatientService) *PatientHandler {
	return &PatientHandler{patientService: patientService}
}

// SearchPatients handles patient search scoped by staff's hospital
func (h *PatientHandler) SearchPatients(c *gin.Context) {
	hospital := middleware.GetStaffHospital(c)
	if hospital == "" {
		response.Unauthorized(c, "Hospital context is missing from authenticated token")
		return
	}

	var criteria domain.PatientSearchCriteria

	// Support both GET query parameters and POST request body
	if c.Request.Method == http.MethodPost && c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&criteria); err != nil {
			response.BadRequest(c, "Invalid search request body", err.Error())
			return
		}
	} else {
		if err := c.ShouldBindQuery(&criteria); err != nil {
			response.BadRequest(c, "Invalid query parameters", err.Error())
			return
		}
	}

	patients, err := h.patientService.Search(c.Request.Context(), hospital, criteria)
	if err != nil {
		response.InternalServerError(c, "Failed to search patients", err.Error())
		return
	}

	response.SuccessWithCount(c, http.StatusOK, "Patients retrieved successfully", patients, len(patients))
}
