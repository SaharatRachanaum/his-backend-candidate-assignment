package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid or expired token")
	ErrInvalidClaim = errors.New("invalid token claims")
)

// StaffClaims represents JWT claims payload for authenticated staff
type StaffClaims struct {
	StaffID  string `json:"staff_id"`
	Username string `json:"username"`
	Hospital string `json:"hospital"`
	jwt.RegisteredClaims
}

// JWTManager handles generation and validation of JWT tokens
type JWTManager struct {
	secretKey     []byte
	tokenDuration time.Duration
}

// NewJWTManager creates a new JWTManager instance
func NewJWTManager(secretKey string, duration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     []byte(secretKey),
		tokenDuration: duration,
	}
}

// Generate creates a signed JWT token for a staff member
func (m *JWTManager) Generate(staffID, username, hospital string) (string, int64, error) {
	expiresAt := time.Now().Add(m.tokenDuration)
	claims := StaffClaims{
		StaffID:  staffID,
		Username: username,
		Hospital: hospital,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   staffID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", 0, err
	}

	return signedToken, int64(m.tokenDuration.Seconds()), nil
}

// Verify parses and verifies the token, returning the staff claims
func (m *JWTManager) Verify(accessToken string) (*StaffClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&StaffClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}
			return m.secretKey, nil
		},
	)

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*StaffClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidClaim
	}

	return claims, nil
}
