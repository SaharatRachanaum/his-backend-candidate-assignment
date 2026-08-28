package service

import (
	"context"
	"testing"

	"github.com/agnos/hospital-middleware/internal/client/his"
	"github.com/agnos/hospital-middleware/internal/domain"
	"github.com/agnos/hospital-middleware/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupPatientService(mockHISURL string) (PatientService, *testutil.MockPatientRepository) {
	patientRepo := testutil.NewMockPatientRepository()

	// Seed patient data
	// Hospital A patients
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
		HospitalID:   "hospital-a",
		PatientHN:    "HN-A-1002",
		FirstNameTH:  "สมหญิง",
		LastNameTH:   "รักสงบ",
		FirstNameEN:  "Somying",
		LastNameEN:   "Raksangob",
		DateOfBirth:  "1995-10-20",
		NationalID:   "1103709876543",
		PassportID:   "AB987654",
		PhoneNumber:  "0898765432",
		Email:        "somying@example.com",
		Gender:       "F",
	})

	// Hospital B patients (isolated tenant)
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

	var hisFactory his.HISFactory
	if mockHISURL != "" {
		factory := his.NewHISFactory(mockHISURL)
		hisFactory = factory
	}

	svc := NewPatientService(patientRepo, hisFactory)
	return svc, patientRepo
}

func TestPatientService_Search_HospitalIsolation(t *testing.T) {
	svc, _ := setupPatientService("")

	// Hospital A staff searches all patients in hospital-a
	patientsA, err := svc.Search(context.Background(), "hospital-a", domain.PatientSearchCriteria{})
	assert.NoError(t, err)
	assert.Len(t, patientsA, 2)
	for _, p := range patientsA {
		assert.Equal(t, "hospital-a", p.HospitalID)
	}

	// Hospital B staff searches all patients in hospital-b
	patientsB, err := svc.Search(context.Background(), "hospital-b", domain.PatientSearchCriteria{})
	assert.NoError(t, err)
	assert.Len(t, patientsB, 1)
	assert.Equal(t, "HN-B-2001", patientsB[0].PatientHN)
	assert.Equal(t, "hospital-b", patientsB[0].HospitalID)

	// Hospital A staff tries searching for Hospital B patient's National ID
	crossSearch, err := svc.Search(context.Background(), "hospital-a", domain.PatientSearchCriteria{
		NationalID: "1209900112233", // belongs to hospital-b
	})
	assert.NoError(t, err)
	assert.Empty(t, crossSearch, "Staff A must NOT be able to view patient belonging to Hospital B")
}

func TestPatientService_Search_Filters(t *testing.T) {
	svc, _ := setupPatientService("")

	// Filter by National ID
	res, err := svc.Search(context.Background(), "hospital-a", domain.PatientSearchCriteria{
		NationalID: "1103701234567",
	})
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "HN-A-1001", res[0].PatientHN)

	// Filter by Passport ID
	res, err = svc.Search(context.Background(), "hospital-a", domain.PatientSearchCriteria{
		PassportID: "AB987654",
	})
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "HN-A-1002", res[0].PatientHN)

	// Filter by Thai first name
	res, err = svc.Search(context.Background(), "hospital-a", domain.PatientSearchCriteria{
		FirstName: "สมชาย",
	})
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "HN-A-1001", res[0].PatientHN)

	// Filter by English first name
	res, err = svc.Search(context.Background(), "hospital-a", domain.PatientSearchCriteria{
		FirstName: "Somying",
	})
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "HN-A-1002", res[0].PatientHN)

	// Filter by Phone
	res, err = svc.Search(context.Background(), "hospital-a", domain.PatientSearchCriteria{
		PhoneNumber: "0812345678",
	})
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "HN-A-1001", res[0].PatientHN)

	// Filter by Email
	res, err = svc.Search(context.Background(), "hospital-a", domain.PatientSearchCriteria{
		Email: "somying@example.com",
	})
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "HN-A-1002", res[0].PatientHN)

	// Filter by Date of birth
	res, err = svc.Search(context.Background(), "hospital-a", domain.PatientSearchCriteria{
		DateOfBirth: "1990-05-15",
	})
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "HN-A-1001", res[0].PatientHN)
}

func TestPatientService_Search_ExternalHISSync(t *testing.T) {
	mockServer := his.NewMockHospitalAServer()
	defer mockServer.Close()

	svc, repo := setupPatientService(mockServer.URL)

	// Patient "EXTERNAL_NEW_PATIENT_1" only exists in Hospital A external HIS mock, not in local DB yet
	res, err := svc.Search(context.Background(), "hospital-a", domain.PatientSearchCriteria{
		NationalID: "EXTERNAL_NEW_PATIENT_1",
	})
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "HN-A-9999", res[0].PatientHN)
	assert.Equal(t, "Somsak", res[0].FirstNameEN)
	assert.Equal(t, "hospital-a", res[0].HospitalID)

	// Verify that the record was automatically cached/persisted in the local repository
	cached, err := repo.GetByHN(context.Background(), "hospital-a", "HN-A-9999")
	assert.NoError(t, err)
	assert.NotNil(t, cached)
	assert.Equal(t, "Somsak", cached.FirstNameEN)
}

func TestPatientService_Search_EmptyHospital(t *testing.T) {
	svc, _ := setupPatientService("")

	// Negative: Empty hospital ID
	_, err := svc.Search(context.Background(), "", domain.PatientSearchCriteria{})
	assert.Error(t, err)
}
