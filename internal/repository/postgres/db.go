package postgres

import (
	"fmt"
	"log"
	"time"

	"github.com/agnos/hospital-middleware/internal/config"
	"github.com/agnos/hospital-middleware/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDatabase initializes PostgreSQL connection using GORM
func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Bangkok",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
	)

	gormLogLevel := logger.Silent
	if cfg.AppEnv == "development" {
		gormLogLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB handle: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Successfully connected to PostgreSQL database")
	return db, nil
}

// AutoMigrate runs GORM schema migrations
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")
	err := db.AutoMigrate(
		&domain.Hospital{},
		&domain.Staff{},
		&domain.Patient{},
	)
	if err != nil {
		return fmt.Errorf("auto migration failed: %w", err)
	}

	// Create composite indexes if needed for optimal performance
	db.Exec("CREATE INDEX IF NOT EXISTS idx_patients_hospital_national_id ON patients(hospital_id, national_id);")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_patients_hospital_passport_id ON patients(hospital_id, passport_id);")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_patients_hospital_phone ON patients(hospital_id, phone_number);")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_patients_hospital_email ON patients(hospital_id, email);")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_patients_hospital_dob ON patients(hospital_id, date_of_birth);")

	log.Println("Database migration completed successfully")
	return nil
}

// SeedData inserts initial sample hospital and patient data if empty
func SeedData(db *gorm.DB) error {
	var hospitalCount int64
	db.Model(&domain.Hospital{}).Count(&hospitalCount)
	if hospitalCount > 0 {
		return nil
	}

	log.Println("Seeding initial hospital & sample patient data...")

	hospitals := []domain.Hospital{
		{
			ID:        "hospital-a",
			Name:      "Hospital A (General)",
			HISAPIURL: "https://hospital-a.api.co.th",
		},
		{
			ID:        "hospital-b",
			Name:      "Hospital B (Specialist)",
			HISAPIURL: "https://hospital-b.api.co.th",
		},
	}

	for _, h := range hospitals {
		if err := db.FirstOrCreate(&h, domain.Hospital{ID: h.ID}).Error; err != nil {
			return err
		}
	}

	patients := []domain.Patient{
		{
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
		},
		{
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
		},
		{
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
		},
	}

	for _, p := range patients {
		var existing domain.Patient
		if err := db.Where("hospital_id = ? AND patient_hn = ?", p.HospitalID, p.PatientHN).First(&existing).Error; err != nil {
			if err := db.Create(&p).Error; err != nil {
				return err
			}
		}
	}

	log.Println("Seed data loaded successfully")
	return nil
}
