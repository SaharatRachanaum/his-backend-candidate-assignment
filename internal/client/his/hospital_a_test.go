package his

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHospitalAClient_SearchPatientByID_Success(t *testing.T) {
	mockServer := NewMockHospitalAServer()
	defer mockServer.Close()

	client := NewHospitalAClient(mockServer.URL, 2*time.Second)
	assert.Equal(t, "hospital-a", client.HospitalID())

	// Positive: Search by National ID
	patient, err := client.SearchPatientByID(context.Background(), "1103701234567")
	assert.NoError(t, err)
	assert.NotNil(t, patient)
	assert.Equal(t, "HN-A-1001", patient.PatientHN)
	assert.Equal(t, "Somchai", patient.FirstNameEN)
	assert.Equal(t, "Jaidee", patient.LastNameEN)
	assert.Equal(t, "1103701234567", patient.NationalID)
	assert.Equal(t, "hospital-a", patient.HospitalID)

	// Positive: Search by Passport ID
	patientPassport, err := client.SearchPatientByID(context.Background(), "AA123456")
	assert.NoError(t, err)
	assert.NotNil(t, patientPassport)
	assert.Equal(t, "HN-A-1001", patientPassport.PatientHN)
}

func TestHospitalAClient_SearchPatientByID_NotFound(t *testing.T) {
	mockServer := NewMockHospitalAServer()
	defer mockServer.Close()

	client := NewHospitalAClient(mockServer.URL, 2*time.Second)

	// Negative: Non-existent patient ID
	patient, err := client.SearchPatientByID(context.Background(), "NON_EXISTENT_ID_999")
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrPatientNotFound)
	assert.Nil(t, patient)
}

func TestHospitalAClient_SearchPatientByID_EmptyID(t *testing.T) {
	client := NewHospitalAClient("http://localhost", 2*time.Second)

	// Negative: Empty ID
	patient, err := client.SearchPatientByID(context.Background(), "")
	assert.Error(t, err)
	assert.Nil(t, patient)
}

func TestHospitalAClient_SearchPatientByID_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewHospitalAClient(server.URL, 2*time.Second)

	// Negative: Server error 500
	patient, err := client.SearchPatientByID(context.Background(), "12345")
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrHISUnavailable)
	assert.Nil(t, patient)
}

func TestDefaultHISFactory(t *testing.T) {
	factory := NewHISFactory("https://hospital-a.api.co.th")

	// Positive: hospital-a client exists
	clientA, ok := factory.GetClient("hospital-a")
	assert.True(t, ok)
	assert.NotNil(t, clientA)
	assert.Equal(t, "hospital-a", clientA.HospitalID())

	// Negative: unconfigured hospital
	clientUnknown, ok := factory.GetClient("hospital-x")
	assert.False(t, ok)
	assert.Nil(t, clientUnknown)
}
