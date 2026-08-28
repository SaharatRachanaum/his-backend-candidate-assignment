-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ==========================================
-- Table: hospitals
-- ==========================================
CREATE TABLE IF NOT EXISTS hospitals (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    his_api_url VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ==========================================
-- Table: staffs
-- ==========================================
CREATE TABLE IF NOT EXISTS staffs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    hospital_id VARCHAR(50) NOT NULL REFERENCES hospitals(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_staffs_hospital_id ON staffs(hospital_id);
CREATE INDEX IF NOT EXISTS idx_staffs_username ON staffs(username);

-- ==========================================
-- Table: patients
-- ==========================================
CREATE TABLE IF NOT EXISTS patients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hospital_id VARCHAR(50) NOT NULL REFERENCES hospitals(id) ON DELETE CASCADE,
    patient_hn VARCHAR(50) NOT NULL,
    first_name_th VARCHAR(100),
    middle_name_th VARCHAR(100),
    last_name_th VARCHAR(100),
    first_name_en VARCHAR(100),
    middle_name_en VARCHAR(100),
    last_name_en VARCHAR(100),
    date_of_birth VARCHAR(20),
    national_id VARCHAR(50),
    passport_id VARCHAR(50),
    phone_number VARCHAR(50),
    email VARCHAR(100),
    gender VARCHAR(10),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Multi-tenant and search optimization indexes
CREATE INDEX IF NOT EXISTS idx_patients_hospital_id ON patients(hospital_id);
CREATE INDEX IF NOT EXISTS idx_patients_hospital_hn ON patients(hospital_id, patient_hn);
CREATE INDEX IF NOT EXISTS idx_patients_hospital_national_id ON patients(hospital_id, national_id);
CREATE INDEX IF NOT EXISTS idx_patients_hospital_passport_id ON patients(hospital_id, passport_id);
CREATE INDEX IF NOT EXISTS idx_patients_hospital_phone ON patients(hospital_id, phone_number);
CREATE INDEX IF NOT EXISTS idx_patients_hospital_email ON patients(hospital_id, email);
CREATE INDEX IF NOT EXISTS idx_patients_hospital_dob ON patients(hospital_id, date_of_birth);
CREATE INDEX IF NOT EXISTS idx_patients_hospital_name_th ON patients(hospital_id, first_name_th, last_name_th);
CREATE INDEX IF NOT EXISTS idx_patients_hospital_name_en ON patients(hospital_id, first_name_en, last_name_en);
