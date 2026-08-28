package his

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/agnos/hospital-middleware/internal/domain"
)

var (
	ErrPatientNotFound = errors.New("patient not found in HIS")
	ErrHISUnavailable  = errors.New("HIS service unavailable")
)

// HospitalAPatientResponse represents the raw JSON schema returned by Hospital A API
type HospitalAPatientResponse struct {
	FirstNameTH  string `json:"first_name_th"`
	MiddleNameTH string `json:"middle_name_th"`
	LastNameTH   string `json:"last_name_th"`
	FirstNameEN  string `json:"first_name_en"`
	MiddleNameEN string `json:"middle_name_en"`
	LastNameEN   string `json:"last_name_en"`
	DateOfBirth  string `json:"date_of_birth"`
	PatientHN    string `json:"patient_hn"`
	NationalID   string `json:"national_id"`
	PassportID   string `json:"passport_id"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
	Gender       string `json:"gender"`
}

// HospitalAClient implements HISClient for Hospital A
type HospitalAClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewHospitalAClient creates a new Hospital A HIS client
func NewHospitalAClient(baseURL string, timeout time.Duration) *HospitalAClient {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	if baseURL == "" {
		baseURL = "https://hospital-a.api.co.th"
	}
	return &HospitalAClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// HospitalID returns the identifier for Hospital A
func (c *HospitalAClient) HospitalID() string {
	return "hospital-a"
}

// SearchPatientByID calls GET /patient/search/{id} on Hospital A API
func (c *HospitalAClient) SearchPatientByID(ctx context.Context, identifier string) (*domain.Patient, error) {
	if identifier == "" {
		return nil, errors.New("identifier is required")
	}

	escapedID := url.PathEscape(identifier)
	endpoint := fmt.Sprintf("%s/patient/search/%s", c.baseURL, escapedID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HIS request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrHISUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrPatientNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: received HTTP %d", ErrHISUnavailable, resp.StatusCode)
	}

	var hisPatient HospitalAPatientResponse
	if err := json.NewDecoder(resp.Body).Decode(&hisPatient); err != nil {
		return nil, fmt.Errorf("failed to decode HIS response: %w", err)
	}

	// Map to internal domain model
	patient := &domain.Patient{
		HospitalID:   c.HospitalID(),
		PatientHN:    hisPatient.PatientHN,
		FirstNameTH:  hisPatient.FirstNameTH,
		MiddleNameTH: hisPatient.MiddleNameTH,
		LastNameTH:   hisPatient.LastNameTH,
		FirstNameEN:  hisPatient.FirstNameEN,
		MiddleNameEN: hisPatient.MiddleNameEN,
		LastNameEN:   hisPatient.LastNameEN,
		DateOfBirth:  hisPatient.DateOfBirth,
		NationalID:   hisPatient.NationalID,
		PassportID:   hisPatient.PassportID,
		PhoneNumber:  hisPatient.PhoneNumber,
		Email:        hisPatient.Email,
		Gender:       hisPatient.Gender,
	}

	return patient, nil
}

// DefaultHISFactory manages registry of HIS clients
type DefaultHISFactory struct {
	clients map[string]HISClient
}

// NewHISFactory creates a new HIS factory with registered hospital adapters
func NewHISFactory(hospitalAURL string) *DefaultHISFactory {
	factory := &DefaultHISFactory{
		clients: make(map[string]HISClient),
	}

	factory.clients["hospital-a"] = NewHospitalAClient(hospitalAURL, 5*time.Second)
	return factory
}

// RegisterClient adds a client to factory (useful for testing or dynamic addition)
func (f *DefaultHISFactory) RegisterClient(hospitalID string, client HISClient) {
	f.clients[hospitalID] = client
}

// GetClient retrieves the HIS client for a given hospital ID
func (f *DefaultHISFactory) GetClient(hospitalID string) (HISClient, bool) {
	client, ok := f.clients[hospitalID]
	return client, ok
}
