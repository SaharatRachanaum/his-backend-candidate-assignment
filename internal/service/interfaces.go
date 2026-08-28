package service

import (
	"context"

	"github.com/agnos/hospital-middleware/internal/domain"
)

// StaffService defines business logic for Staff operations
type StaffService interface {
	Create(ctx context.Context, req domain.CreateStaffRequest) (*domain.StaffResponse, error)
	Login(ctx context.Context, req domain.LoginStaffRequest) (*domain.LoginStaffResponse, error)
}

// PatientService defines business logic for Patient operations
type PatientService interface {
	Search(ctx context.Context, hospitalID string, criteria domain.PatientSearchCriteria) ([]domain.PatientResponse, error)
}
