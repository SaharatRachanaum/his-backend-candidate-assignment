package service

import (
	"context"
	"testing"
	"time"

	"github.com/agnos/hospital-middleware/internal/domain"
	"github.com/agnos/hospital-middleware/internal/testutil"
	"github.com/agnos/hospital-middleware/pkg/token"
	"github.com/stretchr/testify/assert"
)

func setupStaffService() (StaffService, *testutil.MockStaffRepository, *testutil.MockHospitalRepository) {
	staffRepo := testutil.NewMockStaffRepository()
	hospitalRepo := testutil.NewMockHospitalRepository()
	jwtManager := token.NewJWTManager("test_secret_for_staff_service_12345", 24*time.Hour)

	svc := NewStaffService(staffRepo, hospitalRepo, jwtManager)
	return svc, staffRepo, hospitalRepo
}

func TestStaffService_Create_Success(t *testing.T) {
	svc, _, _ := setupStaffService()

	req := domain.CreateStaffRequest{
		Username: "doctor_john",
		Password: "Password123!",
		Hospital: "hospital-a",
	}

	// Positive: create staff
	res, err := svc.Create(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "doctor_john", res.Username)
	assert.Equal(t, "hospital-a", res.Hospital)
	assert.NotEmpty(t, res.ID)
}

func TestStaffService_Create_DuplicateUsername(t *testing.T) {
	svc, _, _ := setupStaffService()

	req := domain.CreateStaffRequest{
		Username: "nurse_mary",
		Password: "Password123!",
		Hospital: "hospital-a",
	}

	_, err := svc.Create(context.Background(), req)
	assert.NoError(t, err)

	// Negative: Duplicate username
	_, err = svc.Create(context.Background(), req)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrUsernameAlreadyExists)
}

func TestStaffService_Create_MissingFields(t *testing.T) {
	svc, _, _ := setupStaffService()

	// Negative: Missing username
	_, err := svc.Create(context.Background(), domain.CreateStaffRequest{
		Username: "",
		Password: "Password123!",
		Hospital: "hospital-a",
	})
	assert.Error(t, err)

	// Negative: Missing password
	_, err = svc.Create(context.Background(), domain.CreateStaffRequest{
		Username: "doctor_tom",
		Password: "",
		Hospital: "hospital-a",
	})
	assert.Error(t, err)

	// Negative: Missing hospital
	_, err = svc.Create(context.Background(), domain.CreateStaffRequest{
		Username: "doctor_tom",
		Password: "Password123!",
		Hospital: "",
	})
	assert.Error(t, err)
}

func TestStaffService_Login_Success(t *testing.T) {
	svc, _, _ := setupStaffService()

	// Create staff first
	_, err := svc.Create(context.Background(), domain.CreateStaffRequest{
		Username: "doctor_somchai",
		Password: "SecurePassword123!",
		Hospital: "hospital-a",
	})
	assert.NoError(t, err)

	// Positive: Login with correct credentials
	loginRes, err := svc.Login(context.Background(), domain.LoginStaffRequest{
		Username: "doctor_somchai",
		Password: "SecurePassword123!",
		Hospital: "hospital-a",
	})
	assert.NoError(t, err)
	assert.NotNil(t, loginRes)
	assert.NotEmpty(t, loginRes.Token)
	assert.Equal(t, "Bearer", loginRes.TokenType)
	assert.Equal(t, "doctor_somchai", loginRes.Staff.Username)
	assert.Equal(t, "hospital-a", loginRes.Staff.Hospital)
}

func TestStaffService_Login_WrongPassword(t *testing.T) {
	svc, _, _ := setupStaffService()

	_, err := svc.Create(context.Background(), domain.CreateStaffRequest{
		Username: "doctor_somchai",
		Password: "SecurePassword123!",
		Hospital: "hospital-a",
	})
	assert.NoError(t, err)

	// Negative: Wrong password
	loginRes, err := svc.Login(context.Background(), domain.LoginStaffRequest{
		Username: "doctor_somchai",
		Password: "WrongPassword!",
		Hospital: "hospital-a",
	})
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidCredentials)
	assert.Nil(t, loginRes)
}

func TestStaffService_Login_WrongHospital(t *testing.T) {
	svc, _, _ := setupStaffService()

	_, err := svc.Create(context.Background(), domain.CreateStaffRequest{
		Username: "doctor_somchai",
		Password: "SecurePassword123!",
		Hospital: "hospital-a",
	})
	assert.NoError(t, err)

	// Negative: Staff belongs to hospital-a, tries logging in into hospital-b
	loginRes, err := svc.Login(context.Background(), domain.LoginStaffRequest{
		Username: "doctor_somchai",
		Password: "SecurePassword123!",
		Hospital: "hospital-b",
	})
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrHospitalMismatch)
	assert.Nil(t, loginRes)
}

func TestStaffService_Login_NonExistentUser(t *testing.T) {
	svc, _, _ := setupStaffService()

	// Negative: User does not exist
	loginRes, err := svc.Login(context.Background(), domain.LoginStaffRequest{
		Username: "unknown_user",
		Password: "Password123!",
		Hospital: "hospital-a",
	})
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidCredentials)
	assert.Nil(t, loginRes)
}
