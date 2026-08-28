package testutil

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/agnos/hospital-middleware/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// MockHospitalRepository in-memory implementation for unit testing
type MockHospitalRepository struct {
	mu        sync.RWMutex
	hospitals map[string]*domain.Hospital
}

func NewMockHospitalRepository() *MockHospitalRepository {
	return &MockHospitalRepository{
		hospitals: make(map[string]*domain.Hospital),
	}
}

func (m *MockHospitalRepository) GetByID(ctx context.Context, id string) (*domain.Hospital, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	h, ok := m.hospitals[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return h, nil
}

func (m *MockHospitalRepository) Create(ctx context.Context, hospital *domain.Hospital) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.hospitals[hospital.ID]; exists {
		return errors.New("hospital already exists")
	}
	hospital.CreatedAt = time.Now()
	hospital.UpdatedAt = time.Now()
	m.hospitals[hospital.ID] = hospital
	return nil
}

func (m *MockHospitalRepository) List(ctx context.Context) ([]domain.Hospital, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make([]domain.Hospital, 0, len(m.hospitals))
	for _, h := range m.hospitals {
		res = append(res, *h)
	}
	return res, nil
}

// MockStaffRepository in-memory implementation for unit testing
type MockStaffRepository struct {
	mu     sync.RWMutex
	staffs map[string]*domain.Staff // keyed by username
	byID   map[uuid.UUID]*domain.Staff
}

func NewMockStaffRepository() *MockStaffRepository {
	return &MockStaffRepository{
		staffs: make(map[string]*domain.Staff),
		byID:   make(map[uuid.UUID]*domain.Staff),
	}
}

func (m *MockStaffRepository) Create(ctx context.Context, staff *domain.Staff) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.staffs[staff.Username]; exists {
		return errors.New("username already exists")
	}
	if staff.ID == uuid.Nil {
		staff.ID = uuid.New()
	}
	staff.CreatedAt = time.Now()
	staff.UpdatedAt = time.Now()
	m.staffs[staff.Username] = staff
	m.byID[staff.ID] = staff
	return nil
}

func (m *MockStaffRepository) GetByUsername(ctx context.Context, username string) (*domain.Staff, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.staffs[username]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return s, nil
}

func (m *MockStaffRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Staff, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.byID[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return s, nil
}

// MockPatientRepository in-memory implementation for unit testing
type MockPatientRepository struct {
	mu       sync.RWMutex
	patients []*domain.Patient
}

func NewMockPatientRepository() *MockPatientRepository {
	return &MockPatientRepository{
		patients: make([]*domain.Patient, 0),
	}
}

func (m *MockPatientRepository) Search(ctx context.Context, hospitalID string, criteria domain.PatientSearchCriteria) ([]domain.Patient, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []domain.Patient
	for _, p := range m.patients {
		// Strict hospital isolation check
		if p.HospitalID != hospitalID {
			continue
		}

		// Filter matching
		if criteria.NationalID != "" && p.NationalID != strings.TrimSpace(criteria.NationalID) {
			continue
		}
		if criteria.PassportID != "" && !strings.EqualFold(p.PassportID, strings.TrimSpace(criteria.PassportID)) {
			continue
		}
		if criteria.FirstName != "" {
			fn := strings.ToLower(strings.TrimSpace(criteria.FirstName))
			if !strings.Contains(strings.ToLower(p.FirstNameTH), fn) && !strings.Contains(strings.ToLower(p.FirstNameEN), fn) {
				continue
			}
		}
		if criteria.MiddleName != "" {
			mn := strings.ToLower(strings.TrimSpace(criteria.MiddleName))
			if !strings.Contains(strings.ToLower(p.MiddleNameTH), mn) && !strings.Contains(strings.ToLower(p.MiddleNameEN), mn) {
				continue
			}
		}
		if criteria.LastName != "" {
			ln := strings.ToLower(strings.TrimSpace(criteria.LastName))
			if !strings.Contains(strings.ToLower(p.LastNameTH), ln) && !strings.Contains(strings.ToLower(p.LastNameEN), ln) {
				continue
			}
		}
		if criteria.DateOfBirth != "" && p.DateOfBirth != strings.TrimSpace(criteria.DateOfBirth) {
			continue
		}
		if criteria.PhoneNumber != "" && p.PhoneNumber != strings.TrimSpace(criteria.PhoneNumber) {
			continue
		}
		if criteria.Email != "" && !strings.EqualFold(p.Email, strings.TrimSpace(criteria.Email)) {
			continue
		}

		results = append(results, *p)
	}

	return results, nil
}

func (m *MockPatientRepository) Create(ctx context.Context, patient *domain.Patient) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if patient.ID == uuid.Nil {
		patient.ID = uuid.New()
	}
	patient.CreatedAt = time.Now()
	patient.UpdatedAt = time.Now()
	m.patients = append(m.patients, patient)
	return nil
}

func (m *MockPatientRepository) Upsert(ctx context.Context, patient *domain.Patient) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, p := range m.patients {
		if p.HospitalID == patient.HospitalID && (p.PatientHN == patient.PatientHN ||
			(p.NationalID != "" && p.NationalID == patient.NationalID) ||
			(p.PassportID != "" && strings.EqualFold(p.PassportID, patient.PassportID))) {
			patient.ID = p.ID
			patient.CreatedAt = p.CreatedAt
			patient.UpdatedAt = time.Now()
			m.patients[i] = patient
			return nil
		}
	}

	if patient.ID == uuid.Nil {
		patient.ID = uuid.New()
	}
	patient.CreatedAt = time.Now()
	patient.UpdatedAt = time.Now()
	m.patients = append(m.patients, patient)
	return nil
}

func (m *MockPatientRepository) GetByHN(ctx context.Context, hospitalID, hn string) (*domain.Patient, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, p := range m.patients {
		if p.HospitalID == hospitalID && p.PatientHN == hn {
			return p, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockPatientRepository) GetByIdentifier(ctx context.Context, hospitalID, id string) (*domain.Patient, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, p := range m.patients {
		if p.HospitalID == hospitalID && (p.NationalID == id || strings.EqualFold(p.PassportID, id)) {
			return p, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}
