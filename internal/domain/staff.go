package domain

import (
	"time"

	"github.com/google/uuid"
)

// Staff represents a hospital staff member entity
type Staff struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Username     string    `gorm:"size:50;uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	HospitalID   string    `gorm:"size:50;index;not null" json:"hospital"`
	Hospital     *Hospital `gorm:"foreignKey:HospitalID;references:ID" json:"-"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// CreateStaffRequest DTO for staff registration
type CreateStaffRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=72"`
	Hospital string `json:"hospital" binding:"required,min=1,max=50"`
}

// StaffResponse DTO for staff public details
type StaffResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Hospital  string    `json:"hospital"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginStaffRequest DTO for staff login
type LoginStaffRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Hospital string `json:"hospital" binding:"required"`
}

// LoginStaffResponse DTO for authentication success
type LoginStaffResponse struct {
	Token     string        `json:"token"`
	TokenType string        `json:"token_type"`
	ExpiresIn int64         `json:"expires_in"`
	Staff     StaffResponse `json:"staff"`
}

// ToResponse maps Staff domain model to StaffResponse DTO
func (s *Staff) ToResponse() StaffResponse {
	return StaffResponse{
		ID:        s.ID,
		Username:  s.Username,
		Hospital:  s.HospitalID,
		CreatedAt: s.CreatedAt,
	}
}
