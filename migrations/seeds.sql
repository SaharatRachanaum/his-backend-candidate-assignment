-- Seed Hospitals
INSERT INTO hospitals (id, name, his_api_url) VALUES
    ('hospital-a', 'Hospital A (General)', 'https://hospital-a.api.co.th'),
    ('hospital-b', 'Hospital B (Specialist)', 'https://hospital-b.api.co.th')
ON CONFLICT (id) DO NOTHING;

-- Seed Staff (Password is "Password123!")
-- Generated with bcrypt cost 10
INSERT INTO staffs (id, username, password_hash, hospital_id) VALUES
    ('a1111111-1111-1111-1111-111111111111', 'doctor_a', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'hospital-a'),
    ('b2222222-2222-2222-2222-222222222222', 'doctor_b', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'hospital-b')
ON CONFLICT (username) DO NOTHING;

-- Seed Patients for Hospital A
INSERT INTO patients (id, hospital_id, patient_hn, first_name_th, middle_name_th, last_name_th, first_name_en, middle_name_en, last_name_en, date_of_birth, national_id, passport_id, phone_number, email, gender) VALUES
    ('c3333333-3333-3333-3333-333333333333', 'hospital-a', 'HN-A-1001', 'สมชาย', '', 'ใจดี', 'Somchai', '', 'Jaidee', '1990-05-15', '1103701234567', 'AA123456', '0812345678', 'somchai@example.com', 'M'),
    ('c4444444-4444-4444-4444-444444444444', 'hospital-a', 'HN-A-1002', 'สมหญิง', '', 'รักสงบ', 'Somying', '', 'Raksangob', '1995-10-20', '1103709876543', 'AB987654', '0898765432', 'somying@example.com', 'F')
ON CONFLICT DO NOTHING;

-- Seed Patients for Hospital B (Tenant Isolation Test)
INSERT INTO patients (id, hospital_id, patient_hn, first_name_th, middle_name_th, last_name_th, first_name_en, middle_name_en, last_name_en, date_of_birth, national_id, passport_id, phone_number, email, gender) VALUES
    ('c5555555-5555-5555-5555-555555555555', 'hospital-b', 'HN-B-2001', 'ประเสริฐ', '', 'มีสุข', 'Prasert', '', 'Meesuk', '1985-03-12', '1209900112233', 'BA556677', '0865551234', 'prasert@example.com', 'M')
ON CONFLICT DO NOTHING;
