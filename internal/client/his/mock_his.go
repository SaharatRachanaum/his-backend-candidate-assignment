package his

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
)

// NewMockHospitalAServer creates an httptest.Server that mocks the Hospital A API
func NewMockHospitalAServer() *httptest.Server {
	mockPatients := map[string]HospitalAPatientResponse{
		"1103701234567": {
			FirstNameTH:  "สมชาย",
			MiddleNameTH: "",
			LastNameTH:   "ใจดี",
			FirstNameEN:  "Somchai",
			MiddleNameEN: "",
			LastNameEN:   "Jaidee",
			DateOfBirth:  "1990-05-15",
			PatientHN:    "HN-A-1001",
			NationalID:   "1103701234567",
			PassportID:   "AA123456",
			PhoneNumber:  "0812345678",
			Email:        "somchai@example.com",
			Gender:       "M",
		},
		"AA123456": {
			FirstNameTH:  "สมชาย",
			MiddleNameTH: "",
			LastNameTH:   "ใจดี",
			FirstNameEN:  "Somchai",
			MiddleNameEN: "",
			LastNameEN:   "Jaidee",
			DateOfBirth:  "1990-05-15",
			PatientHN:    "HN-A-1001",
			NationalID:   "1103701234567",
			PassportID:   "AA123456",
			PhoneNumber:  "0812345678",
			Email:        "somchai@example.com",
			Gender:       "M",
		},
		"1103709876543": {
			FirstNameTH:  "สมหญิง",
			MiddleNameTH: "",
			LastNameTH:   "รักสงบ",
			FirstNameEN:  "Somying",
			MiddleNameEN: "",
			LastNameEN:   "Raksangob",
			DateOfBirth:  "1995-10-20",
			PatientHN:    "HN-A-1002",
			NationalID:   "1103709876543",
			PassportID:   "AB987654",
			PhoneNumber:  "0898765432",
			Email:        "somying@example.com",
			Gender:       "F",
		},
		"EXTERNAL_NEW_PATIENT_1": {
			FirstNameTH:  "สมศักดิ์",
			MiddleNameTH: "",
			LastNameTH:   "มั่งคั่ง",
			FirstNameEN:  "Somsak",
			MiddleNameEN: "",
			LastNameEN:   "Mangkang",
			DateOfBirth:  "1988-08-08",
			PatientHN:    "HN-A-9999",
			NationalID:   "EXTERNAL_NEW_PATIENT_1",
			PassportID:   "EXT999",
			PhoneNumber:  "0899999999",
			Email:        "somsak@example.com",
			Gender:       "M",
		},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		const prefix = "/patient/search/"
		if !strings.HasPrefix(path, prefix) {
			http.NotFound(w, r)
			return
		}

		id := strings.TrimPrefix(path, prefix)
		patient, found := mockPatients[id]
		if !found {
			// Case-insensitive lookup for passport
			for k, v := range mockPatients {
				if strings.EqualFold(k, id) {
					patient = v
					found = true
					break
				}
			}
		}

		if !found {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "patient not found"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(patient)
	})

	return httptest.NewServer(handler)
}
