package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/agnos/hospital-middleware/internal/domain"
	"github.com/agnos/hospital-middleware/internal/middleware"
	"github.com/agnos/hospital-middleware/internal/service"
	"github.com/agnos/hospital-middleware/internal/testutil"
	"github.com/agnos/hospital-middleware/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupPatientHandler() (*gin.Engine, *token.JWTManager, *testutil.MockPatientRepository) {
	gin.SetMode(gin.TestMode)
	patientRepo := testutil.NewMockPatientRepository()
	jwtManager := token.NewJWTManager("test_secret_for_patient_handler_12345", 24*time.Hour)

	// Seed patients
	patientRepo.Create(context.Background(), &domain.Patient{
		ID:           uuid.New(),
		HospitalID:   "hospital-a",
		PatientHN:    "HN-A-1001",
		FirstNameTH:  "สมชาย",
		LastNameTH:   "ใจดี",
		FirstNameEN:  "Somchai",
		LastNameEN:   "Jaidee",
		DateOfBirth:  "1990-05-15",
		NationalID:   "1103701234567",
		PassportID:   "AA123456",
		PhoneNumber:  "0812345678",
		Email:        "somchai@example.com",
		Gender:       "M",
	})
	patientRepo.Create(context.Background(), &domain.Patient{
		ID:           uuid.New(),
		HospitalID:   "hospital-b",
		PatientHN:    "HN-B-2001",
		FirstNameTH:  "ประเสริฐ",
		LastNameTH:   "มีสุข",
		FirstNameEN:  "Prasert",
		LastNameEN:   "Meesuk",
		DateOfBirth:  "1985-03-12",
		NationalID:   "1209900112233",
		PassportID:   "BA556677",
		PhoneNumber:  "0865551234",
		Email:        "prasert@example.com",
		Gender:       "M",
	})

	patientService := service.NewPatientService(patientRepo, nil)
	patientHandler := NewPatientHandler(patientService)

	r := gin.New()
	protected := r.Group("")
	protected.Use(middleware.AuthMiddleware(jwtManager))
	{
		protected.GET("/patient/search", patientHandler.SearchPatients)
		protected.POST("/patient/search", patientHandler.SearchPatients)
	}

	return r, jwtManager, patientRepo
}

func TestPatientHandler_SearchPatients_Success(t *testing.T) {
	r, jwtManager, _ := setupPatientHandler()

	tokenA, _, _ := jwtManager.Generate("staff-1", "doc_a", "hospital-a")

	// Positive: Authenticated GET search
	req, _ := http.NewRequest(http.MethodGet, "/patient/search?first_name=Somchai", nil)
	req.Header.Set("Authorization", "Bearer "+tokenA)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "HN-A-1001")
	assert.Contains(t, w.Body.String(), "Somchai")
	assert.Contains(t, w.Body.String(), `"count":1`)
}

func TestPatientHandler_SearchPatients_POST_Body(t *testing.T) {
	r, jwtManager, _ := setupPatientHandler()

	tokenA, _, _ := jwtManager.Generate("staff-1", "doc_a", "hospital-a")

	// Positive: Authenticated POST search with JSON body
	bodyJSON := []byte(`{"national_id":"1103701234567"}`)
	req, _ := http.NewRequest(http.MethodPost, "/patient/search", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Authorization", "Bearer "+tokenA)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "HN-A-1001")
	assert.Contains(t, w.Body.String(), "1103701234567")
}

func TestPatientHandler_SearchPatients_Unauthorized(t *testing.T) {
	r, _, _ := setupPatientHandler()

	// Negative: No Authorization header
	req, _ := http.NewRequest(http.MethodGet, "/patient/search", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPatientHandler_SearchPatients_HospitalIsolation(t *testing.T) {
	r, jwtManager, _ := setupPatientHandler()

	tokenA, _, _ := jwtManager.Generate("staff-a", "doc_a", "hospital-a")
	tokenB, _, _ := jwtManager.Generate("staff-b", "doc_b", "hospital-b")

	// Staff A queries all patients
	reqA, _ := http.NewRequest(http.MethodGet, "/patient/search", nil)
	reqA.Header.Set("Authorization", "Bearer "+tokenA)
	wA := httptest.NewRecorder()
	r.ServeHTTP(wA, reqA)
	assert.Equal(t, http.StatusOK, wA.Code)
	assert.Contains(t, wA.Body.String(), "HN-A-1001")
	assert.NotContains(t, wA.Body.String(), "HN-B-2001")

	// Staff B queries all patients
	reqB, _ := http.NewRequest(http.MethodGet, "/patient/search", nil)
	reqB.Header.Set("Authorization", "Bearer "+tokenB)
	wB := httptest.NewRecorder()
	r.ServeHTTP(wB, reqB)
	assert.Equal(t, http.StatusOK, wB.Code)
	assert.Contains(t, wB.Body.String(), "HN-B-2001")
	assert.NotContains(t, wB.Body.String(), "HN-A-1001")
}
