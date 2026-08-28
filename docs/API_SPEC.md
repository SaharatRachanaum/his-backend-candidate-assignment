# รายละเอียดข้อกำหนด API (API Specification) - Hospital Middleware System

**Base URL**: `http://localhost` (ผ่าน Nginx Reverse Proxy บนพอร์ต 80) หรือ `http://localhost:8080` (ยิงตรงไปยัง Go Server)

ทุก Request และ Response มีการรับส่งข้อมูลในรูปแบบ `application/json`

---

## 1. ตรวจสอบสถานะระบบ (Health Check)

### `GET /health`
ตรวจสอบสถานะการทำงานของ Web Server และการเชื่อมต่อฐานข้อมูล PostgreSQL

#### ผลลัพธ์ (Response):
- **HTTP Status**: `200 OK`
```json
{
  "success": true,
  "message": "Service is healthy",
  "data": {
    "database": "up",
    "service": "hospital-middleware-api",
    "status": "healthy"
  }
}
```

---

## 2. ระบบจัดการเจ้าหน้าที่ (Staff Authentication)

### 2.1 สร้างบัญชีเจ้าหน้าที่ (`POST /staff/create`)
ลงทะเบียนเจ้าหน้าที่โรงพยาบาลใหม่พร้อมกำหนดรหัสผ่านและโรงพยาบาลที่สังกัด

#### Headers:
`Content-Type: application/json`

#### Request Body:
| ฟิลด์ (Field) | ชนิดข้อมูล | บังคับ | คำอธิบาย |
|---|---|---|---|
| `username` | string | ใช่ | ชื่อผู้ใช้งาน (ขั้นต่ำ 3 ตัวอักษร) |
| `password` | string | ใช่ | รหัสผ่าน (ขั้นต่ำ 6 ตัวอักษร ระบบจะนำไป Hash ด้วย Bcrypt) |
| `hospital` | string | ใช่ | รหัสโรงพยาบาล (เช่น `hospital-a`, `hospital-b`) |

```json
{
  "username": "doctor_somchai",
  "password": "Password123!",
  "hospital": "hospital-a"
}
```

#### ผลลัพธ์ (Response):
- **HTTP Status**: `201 Created`
```json
{
  "success": true,
  "message": "Staff created successfully",
  "data": {
    "id": "e81d432e-b4a8-4261-9f44-88ba42f5ee99",
    "username": "doctor_somchai",
    "hospital": "hospital-a",
    "created_at": "2026-08-28T09:00:00Z"
  }
}
```
- **HTTP Status**: `400 Bad Request` (ข้อมูลไม่ครบหรือรูปแบบไม่ถูกต้อง)
- **HTTP Status**: `409 Conflict` (มี Username นี้อยู่ในระบบแล้ว)

---

### 2.2 เข้าสู่ระบบเจ้าหน้าที่ (`POST /staff/login`)
ยืนยันตัวตนเจ้าหน้าที่เพื่อรับ JWT Bearer Token

#### Request Body:
```json
{
  "username": "doctor_somchai",
  "password": "Password123!",
  "hospital": "hospital-a"
}
```

#### ผลลัพธ์ (Response):
- **HTTP Status**: `200 OK`
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 86400,
    "staff": {
      "id": "e81d432e-b4a8-4261-9f44-88ba42f5ee99",
      "username": "doctor_somchai",
      "hospital": "hospital-a",
      "created_at": "2026-08-28T09:00:00Z"
    }
  }
}
```
- **HTTP Status**: `401 Unauthorized` (รหัสผ่านไม่ถูกต้อง หรือเข้าสู่ระบบด้วยโรงพยาบาลที่ไม่ตรงกับที่สังกัด)

---

## 3. ระบบค้นหาข้อมูลผู้ป่วย (Patient Search)

### 3.1 ค้นหาข้อมูลผู้ป่วย (`GET /patient/search` หรือ `POST /patient/search`)
ค้นหาข้อมูลผู้ป่วยตามเงื่อนไขที่ระบุ **ผลลัพธ์จะถูกจำกัดเฉพาะผู้ป่วยที่อยู่ในโรงพยาบาลเดียวกับเจ้าหน้าที่ที่เข้าสู่ระบบเท่านั้น** (Multi-Tenancy Isolation)

#### Headers:
`Authorization: Bearer <JWT_TOKEN>`

#### Parameter ค้นหา (Query Parameters - ทุกฟิลด์เป็นทางเลือก/Optional):
| Parameter | ชนิดข้อมูล | คำอธิบาย |
|---|---|---|
| `national_id` | string | เลขประจำตัวประชาชน |
| `passport_id` | string | เลขหนังสือเดินทาง (Passport) |
| `first_name` | string | ชื่อ (ค้นหาได้ทั้งภาษาไทยและภาษาอังกฤษแบบไม่เจาะจงตัวพิมพ์ใหญ่-เล็ก) |
| `middle_name` | string | ชื่อกลาง (ภาษาไทยหรืออังกฤษ) |
| `last_name` | string | นามสกุล (ภาษาไทยหรืออังกฤษ) |
| `date_of_birth` | string | วันเดือนปีเกิด (รูปแบบ: `YYYY-MM-DD` เช่น `1990-05-15`) |
| `phone_number` | string | เบอร์โทรศัพท์ |
| `email` | string | อีเมล |

#### ตัวอย่าง Request:
```bash
GET /patient/search?first_name=Somchai HTTP/1.1
Host: localhost
Authorization: Bearer eyJhbGciOiJIUzI1Ni...
```

#### ผลลัพธ์ (Response):
- **HTTP Status**: `200 OK`
```json
{
  "success": true,
  "message": "Patients retrieved successfully",
  "count": 1,
  "data": [
    {
      "id": "c3333333-3333-3333-3333-333333333333",
      "hospital_id": "hospital-a",
      "patient_hn": "HN-A-1001",
      "first_name_th": "สมชาย",
      "middle_name_th": "",
      "last_name_th": "ใจดี",
      "first_name_en": "Somchai",
      "middle_name_en": "",
      "last_name_en": "Jaidee",
      "date_of_birth": "1990-05-15",
      "national_id": "1103701234567",
      "passport_id": "AA123456",
      "phone_number": "0812345678",
      "email": "somchai@example.com",
      "gender": "M",
      "created_at": "2026-08-28T09:00:00Z"
    }
  ]
}
```
- **HTTP Status**: `401 Unauthorized` (ไม่ได้แนบ JWT Token หรือ Token ไม่ถูกต้อง/หมดอายุ)
