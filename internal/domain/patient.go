package domain

import (
	"time"

	"github.com/google/uuid"
)

// Patient represents a patient record in the middleware database
type Patient struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	HospitalID   string    `gorm:"size:50;index;not null" json:"hospital_id"`
	Hospital     *Hospital `gorm:"foreignKey:HospitalID;references:ID" json:"-"`
	PatientHN    string    `gorm:"size:50;index;not null" json:"patient_hn"`
	FirstNameTH  string    `gorm:"size:100;index" json:"first_name_th,omitempty"`
	MiddleNameTH string    `gorm:"size:100" json:"middle_name_th,omitempty"`
	LastNameTH   string    `gorm:"size:100;index" json:"last_name_th,omitempty"`
	FirstNameEN  string    `gorm:"size:100;index" json:"first_name_en,omitempty"`
	MiddleNameEN string    `gorm:"size:100" json:"middle_name_en,omitempty"`
	LastNameEN   string    `gorm:"size:100;index" json:"last_name_en,omitempty"`
	DateOfBirth  string    `gorm:"size:20;index" json:"date_of_birth,omitempty"` // Format: YYYY-MM-DD
	NationalID   string    `gorm:"size:50;index" json:"national_id,omitempty"`
	PassportID   string    `gorm:"size:50;index" json:"passport_id,omitempty"`
	PhoneNumber  string    `gorm:"size:50;index" json:"phone_number,omitempty"`
	Email        string    `gorm:"size:100;index" json:"email,omitempty"`
	Gender       string    `gorm:"size:10" json:"gender,omitempty"` // "M", "F", "Other"
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// PatientSearchCriteria represents search parameters for querying patients
type PatientSearchCriteria struct {
	NationalID  string `form:"national_id" json:"national_id"`
	PassportID  string `form:"passport_id" json:"passport_id"`
	FirstName   string `form:"first_name" json:"first_name"`
	MiddleName  string `form:"middle_name" json:"middle_name"`
	LastName    string `form:"last_name" json:"last_name"`
	DateOfBirth string `form:"date_of_birth" json:"date_of_birth"`
	PhoneNumber string `form:"phone_number" json:"phone_number"`
	Email       string `form:"email" json:"email"`
}

// PatientResponse represents the public JSON output for patient search
type PatientResponse struct {
	ID           uuid.UUID `json:"id"`
	HospitalID   string    `json:"hospital_id"`
	PatientHN    string    `json:"patient_hn"`
	FirstNameTH  string    `json:"first_name_th,omitempty"`
	MiddleNameTH string    `json:"middle_name_th,omitempty"`
	LastNameTH   string    `json:"last_name_th,omitempty"`
	FirstNameEN  string    `json:"first_name_en,omitempty"`
	MiddleNameEN string    `json:"middle_name_en,omitempty"`
	LastNameEN   string    `json:"last_name_en,omitempty"`
	DateOfBirth  string    `json:"date_of_birth,omitempty"`
	NationalID   string    `json:"national_id,omitempty"`
	PassportID   string    `json:"passport_id,omitempty"`
	PhoneNumber  string    `json:"phone_number,omitempty"`
	Email        string    `json:"email,omitempty"`
	Gender       string    `json:"gender,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// ToResponse converts Patient entity to PatientResponse DTO
func (p *Patient) ToResponse() PatientResponse {
	return PatientResponse{
		ID:           p.ID,
		HospitalID:   p.HospitalID,
		PatientHN:    p.PatientHN,
		FirstNameTH:  p.FirstNameTH,
		MiddleNameTH: p.MiddleNameTH,
		LastNameTH:   p.LastNameTH,
		FirstNameEN:  p.FirstNameEN,
		MiddleNameEN: p.MiddleNameEN,
		LastNameEN:   p.LastNameEN,
		DateOfBirth:  p.DateOfBirth,
		NationalID:   p.NationalID,
		PassportID:   p.PassportID,
		PhoneNumber:  p.PhoneNumber,
		Email:        p.Email,
		Gender:       p.Gender,
		CreatedAt:    p.CreatedAt,
	}
}

// ToResponseList converts slice of Patient to slice of PatientResponse
func ToResponseList(patients []Patient) []PatientResponse {
	res := make([]PatientResponse, len(patients))
	for i, p := range patients {
		res[i] = p.ToResponse()
	}
	return res
}
