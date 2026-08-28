package postgres

import (
	"context"

	"github.com/agnos/hospital-middleware/internal/domain"
	"github.com/agnos/hospital-middleware/internal/repository"
	"gorm.io/gorm"
)

type hospitalRepository struct {
	db *gorm.DB
}

// NewHospitalRepository creates a new HospitalRepository instance
func NewHospitalRepository(db *gorm.DB) repository.HospitalRepository {
	return &hospitalRepository{db: db}
}

func (r *hospitalRepository) GetByID(ctx context.Context, id string) (*domain.Hospital, error) {
	var hospital domain.Hospital
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&hospital).Error
	if err != nil {
		return nil, err
	}
	return &hospital, nil
}

func (r *hospitalRepository) Create(ctx context.Context, hospital *domain.Hospital) error {
	return r.db.WithContext(ctx).Create(hospital).Error
}

func (r *hospitalRepository) List(ctx context.Context) ([]domain.Hospital, error) {
	var hospitals []domain.Hospital
	err := r.db.WithContext(ctx).Find(&hospitals).Error
	return hospitals, err
}
