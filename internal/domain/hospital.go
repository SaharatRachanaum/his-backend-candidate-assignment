package domain

import (
	"time"
)

// Hospital represents a hospital entity in the middleware system
type Hospital struct {
	ID        string    `gorm:"primaryKey;size:50" json:"id"` // e.g. "hospital-a"
	Name      string    `gorm:"size:100;not null" json:"name"`
	HISAPIURL string    `gorm:"size:255" json:"his_api_url,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
