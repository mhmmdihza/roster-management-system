package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"payd/services/auth"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) VerifySignatureJWT(tokenStr string) (*auth.Identity, error) {
	args := m.Called(tokenStr)
	if identity, ok := args.Get(0).(*auth.Identity); ok {
		return identity, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAuthService) Login(ctx context.Context, username, password string) (*auth.Identity, error) {
	return nil, nil
}

func (m *MockAuthService) RegisterNewUser(ctx context.Context, email string, primaryRole int, roleAdmin bool) (string, error) {
	return "", nil
}

func TestJWTAuthorizeRoles(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		token          string
		cookiePresent  bool
		mockIdentity   *auth.Identity
		mockError      error
		allowedRoles   []string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "missing token",
			cookiePresent:  false,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Missing token",
		},
		{
			name:           "invalid token",
			cookiePresent:  true,
			token:          "invalid-token",
			mockError:      errors.New("invalid"),
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid token",
		},
		{
			name:           "unauthorized role",
			cookiePresent:  true,
			token:          "valid-token",
			mockIdentity:   &auth.Identity{Role: "employee"},
			allowedRoles:   []string{"admin"},
			expectedStatus: http.StatusForbidden,
			expectedBody:   "Insufficient role",
		},
		{
			name:           "authorized access",
			cookiePresent:  true,
			token:          "valid-token",
			mockIdentity:   &auth.Identity{Role: "admin"},
			allowedRoles:   []string{"admin"},
			expectedStatus: http.StatusOK,
			expectedBody:   "access granted",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAuth := new(MockAuthService)
			if tc.cookiePresent {
				mockAuth.On("VerifySignatureJWT", tc.token).Return(tc.mockIdentity, tc.mockError)
			}

			router := gin.New()
			router.Use(JWTAuthorizeRoles(mockAuth, tc.allowedRoles...))
			router.GET("/protected", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "access granted"})
			})

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tc.cookiePresent {
				req.AddCookie(&http.Cookie{Name: "token", Value: tc.token})
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectedBody)

			mockAuth.AssertExpectations(t)
		})
	}
}
