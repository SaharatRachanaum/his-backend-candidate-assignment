package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJWTManager_GenerateAndVerify(t *testing.T) {
	manager := NewJWTManager("test_secret_key_12345", 1*time.Hour)

	staffID := "a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d"
	username := "doctor_somchai"
	hospital := "hospital-a"

	// 1. Positive: Generate token
	tokenStr, expiresIn, err := manager.Generate(staffID, username, hospital)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)
	assert.Equal(t, int64(3600), expiresIn)

	// 2. Positive: Verify token
	claims, err := manager.Verify(tokenStr)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, staffID, claims.StaffID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, hospital, claims.Hospital)
}

func TestJWTManager_ExpiredToken(t *testing.T) {
	// Negative: Expired token
	manager := NewJWTManager("test_secret_key_12345", -1*time.Hour)

	tokenStr, _, err := manager.Generate("123", "expired_user", "hospital-a")
	assert.NoError(t, err)

	claims, err := manager.Verify(tokenStr)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTManager_InvalidSecret(t *testing.T) {
	manager1 := NewJWTManager("secret_key_one", 1*time.Hour)
	manager2 := NewJWTManager("secret_key_two", 1*time.Hour)

	tokenStr, _, err := manager1.Generate("123", "user", "hospital-a")
	assert.NoError(t, err)

	// Negative: Verify with wrong secret key
	claims, err := manager2.Verify(tokenStr)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTManager_MalformedToken(t *testing.T) {
	manager := NewJWTManager("test_secret_key_12345", 1*time.Hour)

	claims, err := manager.Verify("this.is.not.a.valid.jwt.token")
	assert.Error(t, err)
	assert.Nil(t, claims)
}
