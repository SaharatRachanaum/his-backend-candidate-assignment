package postgres

import (
	"context"

	"github.com/agnos/hospital-middleware/internal/domain"
	"github.com/agnos/hospital-middleware/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type staffRepository struct {
	db *gorm.DB
}

// NewStaffRepository creates a new StaffRepository instance
func NewStaffRepository(db *gorm.DB) repository.StaffRepository {
	return &staffRepository{db: db}
}

func (r *staffRepository) Create(ctx context.Context, staff *domain.Staff) error {
	return r.db.WithContext(ctx).Create(staff).Error
}

func (r *staffRepository) GetByUsername(ctx context.Context, username string) (*domain.Staff, error) {
	var staff domain.Staff
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&staff).Error
	if err != nil {
		return nil, err
	}
	return &staff, nil
}

func (r *staffRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Staff, error) {
	var staff domain.Staff
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&staff).Error
	if err != nil {
		return nil, err
	}
	return &staff, nil
}
