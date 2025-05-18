package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAndVerifyJWT(t *testing.T) {
	secret := []byte("supersecretkey")

	identity := &Identity{
		ID:           "user123",
		Email:        "user@example.com",
		EmployeeId:   "42",
		EmployeeName: "emp",
		Role:         "admin",
		PrimaryRole:  1,
		jwtSecret:    secret,
	}

	tokenStr, err := identity.GenerateJWT(time.Minute)
	assert.NoError(t, err)

	auth := &Auth{jwtSecret: secret}

	verifiedIdentity, err := auth.VerifySignatureJWT(tokenStr)
	assert.NoError(t, err)

	assert.Equal(t, verifiedIdentity.Email, identity.Email)

	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":           identity.ID,
		"email":         identity.Email,
		"employee_id":   identity.EmployeeId,
		"employee_name": identity.EmployeeName,
		"role":          identity.Role,
		"primary_role":  identity.PrimaryRole,
		"exp":           time.Now().Add(-time.Hour).Unix(), // expired
	})

	expiredTokenStr, err := expiredToken.SignedString(secret)
	assert.NoError(t, err)

	_, err = auth.VerifySignatureJWT(expiredTokenStr)
	if err == nil || err.Error() != "token has invalid claims: token is expired" {
		t.Errorf("expected token expired error, got %v", err)
	}

	tamperedTokenStr := tokenStr[:len(tokenStr)-1] + "x"

	_, err = auth.VerifySignatureJWT(tamperedTokenStr)
	assert.Error(t, err)

}
