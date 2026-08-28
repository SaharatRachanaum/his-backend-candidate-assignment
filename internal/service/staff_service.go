package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/agnos/hospital-middleware/internal/domain"
	"github.com/agnos/hospital-middleware/internal/repository"
	"github.com/agnos/hospital-middleware/pkg/token"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUsernameAlreadyExists = errors.New("username is already taken")
	ErrInvalidCredentials    = errors.New("invalid username, password, or hospital")
	ErrHospitalNotFound      = errors.New("hospital not found")
	ErrHospitalMismatch      = errors.New("staff does not belong to the specified hospital")
)

type staffService struct {
	staffRepo    repository.StaffRepository
	hospitalRepo repository.HospitalRepository
	jwtManager   *token.JWTManager
}

// NewStaffService creates a new StaffService instance
func NewStaffService(
	staffRepo repository.StaffRepository,
	hospitalRepo repository.HospitalRepository,
	jwtManager *token.JWTManager,
) StaffService {
	return &staffService{
		staffRepo:    staffRepo,
		hospitalRepo: hospitalRepo,
		jwtManager:   jwtManager,
	}
}

func (s *staffService) Create(ctx context.Context, req domain.CreateStaffRequest) (*domain.StaffResponse, error) {
	username := strings.TrimSpace(req.Username)
	hospitalID := strings.TrimSpace(req.Hospital)

	if username == "" || req.Password == "" || hospitalID == "" {
		return nil, errors.New("username, password, and hospital are required")
	}

	// Check if hospital exists in DB, if not auto-register hospital
	_, err := s.hospitalRepo.GetByID(ctx, hospitalID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Auto create hospital entry
			newHospital := &domain.Hospital{
				ID:   hospitalID,
				Name: fmt.Sprintf("Hospital %s", strings.ToUpper(hospitalID)),
			}
			if err := s.hospitalRepo.Create(ctx, newHospital); err != nil {
				return nil, fmt.Errorf("failed to register hospital: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error querying hospital: %w", err)
		}
	}

	// Check if username already exists
	existingStaff, err := s.staffRepo.GetByUsername(ctx, username)
	if err == nil && existingStaff != nil {
		return nil, ErrUsernameAlreadyExists
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error checking existing staff: %w", err)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	staff := &domain.Staff{
		Username:     username,
		PasswordHash: string(hashedPassword),
		HospitalID:   hospitalID,
	}

	if err := s.staffRepo.Create(ctx, staff); err != nil {
		return nil, fmt.Errorf("failed to create staff: %w", err)
	}

	res := staff.ToResponse()
	return &res, nil
}

func (s *staffService) Login(ctx context.Context, req domain.LoginStaffRequest) (*domain.LoginStaffResponse, error) {
	username := strings.TrimSpace(req.Username)
	hospitalID := strings.TrimSpace(req.Hospital)

	if username == "" || req.Password == "" || hospitalID == "" {
		return nil, ErrInvalidCredentials
	}

	staff, err := s.staffRepo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("error retrieving staff: %w", err)
	}

	// Check hospital matching (Task 3: Staff can only access same hospital)
	if staff.HospitalID != hospitalID {
		return nil, ErrHospitalMismatch
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(staff.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT Token
	jwtToken, expiresIn, err := s.jwtManager.Generate(staff.ID.String(), staff.Username, staff.HospitalID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	return &domain.LoginStaffResponse{
		Token:     jwtToken,
		TokenType: "Bearer",
		ExpiresIn: expiresIn,
		Staff:     staff.ToResponse(),
	}, nil
}
