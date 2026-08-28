package postgres

import (
	"context"
	"strings"

	"github.com/agnos/hospital-middleware/internal/domain"
	"github.com/agnos/hospital-middleware/internal/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type patientRepository struct {
	db *gorm.DB
}

// NewPatientRepository creates a new PatientRepository instance
func NewPatientRepository(db *gorm.DB) repository.PatientRepository {
	return &patientRepository{db: db}
}

// Search queries patients scoped strictly by hospital_id and optional criteria
func (r *patientRepository) Search(ctx context.Context, hospitalID string, criteria domain.PatientSearchCriteria) ([]domain.Patient, error) {
	query := r.db.WithContext(ctx).Where("hospital_id = ?", hospitalID)

	if criteria.NationalID != "" {
		query = query.Where("national_id = ?", strings.TrimSpace(criteria.NationalID))
	}

	if criteria.PassportID != "" {
		query = query.Where("LOWER(passport_id) = LOWER(?)", strings.TrimSpace(criteria.PassportID))
	}

	if criteria.FirstName != "" {
		fn := "%" + strings.TrimSpace(criteria.FirstName) + "%"
		query = query.Where("(first_name_th ILIKE ? OR first_name_en ILIKE ?)", fn, fn)
	}

	if criteria.MiddleName != "" {
		mn := "%" + strings.TrimSpace(criteria.MiddleName) + "%"
		query = query.Where("(middle_name_th ILIKE ? OR middle_name_en ILIKE ?)", mn, mn)
	}

	if criteria.LastName != "" {
		ln := "%" + strings.TrimSpace(criteria.LastName) + "%"
		query = query.Where("(last_name_th ILIKE ? OR last_name_en ILIKE ?)", ln, ln)
	}

	if criteria.DateOfBirth != "" {
		query = query.Where("date_of_birth = ?", strings.TrimSpace(criteria.DateOfBirth))
	}

	if criteria.PhoneNumber != "" {
		query = query.Where("phone_number = ?", strings.TrimSpace(criteria.PhoneNumber))
	}

	if criteria.Email != "" {
		query = query.Where("LOWER(email) = LOWER(?)", strings.TrimSpace(criteria.Email))
	}

	var patients []domain.Patient
	err := query.Order("created_at DESC").Find(&patients).Error
	if err != nil {
		return nil, err
	}

	return patients, nil
}

func (r *patientRepository) Create(ctx context.Context, patient *domain.Patient) error {
	return r.db.WithContext(ctx).Create(patient).Error
}

func (r *patientRepository) Upsert(ctx context.Context, patient *domain.Patient) error {
	var existing domain.Patient
	err := r.db.WithContext(ctx).
		Where("hospital_id = ? AND (patient_hn = ? OR (national_id != '' AND national_id = ?) OR (passport_id != '' AND passport_id = ?))",
			patient.HospitalID, patient.PatientHN, patient.NationalID, patient.PassportID).
		First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		return r.db.WithContext(ctx).Create(patient).Error
	} else if err != nil {
		return err
	}

	// Update existing record
	patient.ID = existing.ID
	return r.db.WithContext(ctx).
		Clauses(clause.Returning{}).
		Save(patient).Error
}

func (r *patientRepository) GetByHN(ctx context.Context, hospitalID, hn string) (*domain.Patient, error) {
	var patient domain.Patient
	err := r.db.WithContext(ctx).
		Where("hospital_id = ? AND patient_hn = ?", hospitalID, hn).
		First(&patient).Error
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

func (r *patientRepository) GetByIdentifier(ctx context.Context, hospitalID, id string) (*domain.Patient, error) {
	var patient domain.Patient
	err := r.db.WithContext(ctx).
		Where("hospital_id = ? AND (national_id = ? OR LOWER(passport_id) = LOWER(?))", hospitalID, id, id).
		First(&patient).Error
	if err != nil {
		return nil, err
	}
	return &patient, nil
}
