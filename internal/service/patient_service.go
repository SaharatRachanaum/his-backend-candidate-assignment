package service

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/agnos/hospital-middleware/internal/client/his"
	"github.com/agnos/hospital-middleware/internal/domain"
	"github.com/agnos/hospital-middleware/internal/repository"
)

type patientService struct {
	patientRepo repository.PatientRepository
	hisFactory  his.HISFactory
}

// NewPatientService creates a new PatientService instance
func NewPatientService(
	patientRepo repository.PatientRepository,
	hisFactory his.HISFactory,
) PatientService {
	return &patientService{
		patientRepo: patientRepo,
		hisFactory:  hisFactory,
	}
}

// Search retrieves patients scoped strictly to the staff's hospital
func (s *patientService) Search(ctx context.Context, hospitalID string, criteria domain.PatientSearchCriteria) ([]domain.PatientResponse, error) {
	if hospitalID == "" {
		return nil, errors.New("hospital identifier is required")
	}

	// 1. Search local PostgreSQL database scoped by hospitalID
	patients, err := s.patientRepo.Search(ctx, hospitalID, criteria)
	if err != nil {
		return nil, err
	}

	// 2. Middleware HIS integration:
	// If search criteria includes national_id or passport_id and we have an HIS client for this hospital,
	// query the hospital's HIS API to sync new/updated records seamlessly.
	identifier := strings.TrimSpace(criteria.NationalID)
	if identifier == "" {
		identifier = strings.TrimSpace(criteria.PassportID)
	}

	if identifier != "" && s.hisFactory != nil {
		if hisClient, ok := s.hisFactory.GetClient(hospitalID); ok && hisClient != nil {
			hisPatient, err := hisClient.SearchPatientByID(ctx, identifier)
			if err == nil && hisPatient != nil {
				// Ensure hospital isolation
				hisPatient.HospitalID = hospitalID

				// Upsert to local cache/database
				if err := s.patientRepo.Upsert(ctx, hisPatient); err != nil {
					log.Printf("Warning: failed to upsert HIS patient into DB: %v\n", err)
				}

				// Check if already in returned slice
				alreadyInList := false
				for _, p := range patients {
					if p.PatientHN == hisPatient.PatientHN ||
						(p.NationalID != "" && p.NationalID == hisPatient.NationalID) ||
						(p.PassportID != "" && strings.EqualFold(p.PassportID, hisPatient.PassportID)) {
						alreadyInList = true
						break
					}
				}

				if !alreadyInList {
					patients = append(patients, *hisPatient)
				}
			}
		}
	}

	return domain.ToResponseList(patients), nil
}
