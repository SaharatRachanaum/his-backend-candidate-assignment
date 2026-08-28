package his

import (
	"context"

	"github.com/agnos/hospital-middleware/internal/domain"
)

// HISClient defines interface for external Hospital Information Systems
type HISClient interface {
	SearchPatientByID(ctx context.Context, identifier string) (*domain.Patient, error)
	HospitalID() string
}

// HISFactory creates appropriate HIS client for a given hospital
type HISFactory interface {
	GetClient(hospitalID string) (HISClient, bool)
}
