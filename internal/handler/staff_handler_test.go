package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/agnos/hospital-middleware/internal/domain"
	"github.com/agnos/hospital-middleware/internal/service"
	"github.com/agnos/hospital-middleware/internal/testutil"
	"github.com/agnos/hospital-middleware/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupStaffHandler() (*gin.Engine, *StaffHandler) {
	gin.SetMode(gin.TestMode)
	staffRepo := testutil.NewMockStaffRepository()
	hospitalRepo := testutil.NewMockHospitalRepository()
	jwtManager := token.NewJWTManager("test_secret_key_12345", 24*time.Hour)

	staffService := service.NewStaffService(staffRepo, hospitalRepo, jwtManager)
	staffHandler := NewStaffHandler(staffService)

	r := gin.New()
	r.POST("/staff/create", staffHandler.CreateStaff)
	r.POST("/staff/login", staffHandler.LoginStaff)

	return r, staffHandler
}

func TestStaffHandler_CreateStaff_Success(t *testing.T) {
	r, _ := setupStaffHandler()

	payload := domain.CreateStaffRequest{
		Username: "doctor_somchai",
		Password: "Password123!",
		Hospital: "hospital-a",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/staff/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// Positive: 201 Created
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Staff created successfully")
	assert.Contains(t, w.Body.String(), "doctor_somchai")
}

func TestStaffHandler_CreateStaff_InvalidInput(t *testing.T) {
	r, _ := setupStaffHandler()

	// Negative: Missing hospital and short password
	payload := map[string]string{
		"username": "doc",
		"password": "12",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/staff/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStaffHandler_CreateStaff_DuplicateUsername(t *testing.T) {
	r, _ := setupStaffHandler()

	payload := domain.CreateStaffRequest{
		Username: "nurse_joy",
		Password: "Password123!",
		Hospital: "hospital-a",
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest(http.MethodPost, "/staff/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Negative: Duplicate 409 Conflict
	reqDup, _ := http.NewRequest(http.MethodPost, "/staff/create", bytes.NewBuffer(body))
	reqDup.Header.Set("Content-Type", "application/json")
	wDup := httptest.NewRecorder()
	r.ServeHTTP(wDup, reqDup)

	assert.Equal(t, http.StatusConflict, wDup.Code)
}

func TestStaffHandler_LoginStaff_Success(t *testing.T) {
	r, _ := setupStaffHandler()

	// Register first
	createPayload := domain.CreateStaffRequest{
		Username: "doctor_somchai",
		Password: "Password123!",
		Hospital: "hospital-a",
	}
	createBody, _ := json.Marshal(createPayload)
	reqCreate, _ := http.NewRequest(http.MethodPost, "/staff/create", bytes.NewBuffer(createBody))
	reqCreate.Header.Set("Content-Type", "application/json")
	wCreate := httptest.NewRecorder()
	r.ServeHTTP(wCreate, reqCreate)
	assert.Equal(t, http.StatusCreated, wCreate.Code)

	// Positive: Login with valid credentials
	loginPayload := domain.LoginStaffRequest{
		Username: "doctor_somchai",
		Password: "Password123!",
		Hospital: "hospital-a",
	}
	loginBody, _ := json.Marshal(loginPayload)
	reqLogin, _ := http.NewRequest(http.MethodPost, "/staff/login", bytes.NewBuffer(loginBody))
	reqLogin.Header.Set("Content-Type", "application/json")
	wLogin := httptest.NewRecorder()
	r.ServeHTTP(wLogin, reqLogin)

	assert.Equal(t, http.StatusOK, wLogin.Code)
	assert.Contains(t, wLogin.Body.String(), "token")
	assert.Contains(t, wLogin.Body.String(), "Bearer")
}

func TestStaffHandler_LoginStaff_InvalidCredentials(t *testing.T) {
	r, _ := setupStaffHandler()

	// Negative: Wrong password
	loginPayload := domain.LoginStaffRequest{
		Username: "doctor_somchai",
		Password: "WrongPassword!",
		Hospital: "hospital-a",
	}
	loginBody, _ := json.Marshal(loginPayload)
	reqLogin, _ := http.NewRequest(http.MethodPost, "/staff/login", bytes.NewBuffer(loginBody))
	reqLogin.Header.Set("Content-Type", "application/json")
	wLogin := httptest.NewRecorder()
	r.ServeHTTP(wLogin, reqLogin)

	assert.Equal(t, http.StatusUnauthorized, wLogin.Code)
}
