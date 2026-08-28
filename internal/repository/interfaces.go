package repository

import (
	"context"

	"github.com/agnos/hospital-middleware/internal/domain"
	"github.com/google/uuid"
)

// HospitalRepository defines storage operations for Hospital
type HospitalRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Hospital, error)
	Create(ctx context.Context, hospital *domain.Hospital) error
	List(ctx context.Context) ([]domain.Hospital, error)
}

// StaffRepository defines storage operations for Staff
type StaffRepository interface {
	Create(ctx context.Context, staff *domain.Staff) error
	GetByUsername(ctx context.Context, username string) (*domain.Staff, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Staff, error)
}

// PatientRepository defines storage operations for Patient
type PatientRepository interface {
	Search(ctx context.Context, hospitalID string, criteria domain.PatientSearchCriteria) ([]domain.Patient, error)
	Create(ctx context.Context, patient *domain.Patient) error
	Upsert(ctx context.Context, patient *domain.Patient) error
	GetByHN(ctx context.Context, hospitalID, hn string) (*domain.Patient, error)
	GetByIdentifier(ctx context.Context, hospitalID, id string) (*domain.Patient, error)
}
