package public

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"payd/services/auth"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuth mocks auth.AuthInterface
type MockAuth struct {
	mock.Mock
}

func (m *MockAuth) Login(ctx context.Context, username, password string) (*auth.Identity, error) {
	args := m.Called(ctx, username, password)
	identity, _ := args.Get(0).(*auth.Identity)
	return identity, args.Error(1)
}

// Sample login request body struct matching your expected input
type loginRequest struct {
	Username string `json:"username" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func TestPublicLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		reqBody        interface{}
		mockLoginResp  *auth.Identity
		mockLoginErr   error
		wantStatusCode int
		wantRespBody   string
	}{
		{
			name: "success login",
			reqBody: loginRequest{
				Username: "abc@mail.com",
				Password: "password123",
			},
			mockLoginResp: &auth.Identity{
				ID:           "userid",
				Email:        "abc@mail.com",
				Role:         "employee",
				EmployeeId:   "4",
				EmployeeName: "name",
				PrimaryRole:  7,
			},
			mockLoginErr:   nil,
			wantStatusCode: http.StatusOK,
			wantRespBody:   "Login successful",
		},
		{
			name: "validation error - missing username",
			reqBody: map[string]string{
				"password": "password123",
			},
			wantStatusCode: http.StatusBadRequest,
			wantRespBody:   "Key: 'LoginRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag",
		},
		{
			name: "validation error - invalid email",
			reqBody: loginRequest{
				Username: "invalid-email",
				Password: "password123",
			},
			wantStatusCode: http.StatusBadRequest,
			wantRespBody:   "Key: 'LoginRequest.Username' Error:Field validation for 'Username' failed on the 'email' tag",
		},
		{
			name: "auth login error",
			reqBody: loginRequest{
				Username: "abc@mail.com",
				Password: "wrongpass",
			},
			mockLoginResp:  nil,
			mockLoginErr:   auth.ErrInvalidCredential,
			wantStatusCode: http.StatusBadRequest,
			wantRespBody:   "invalid credential",
		},

		{
			name: "account not yet active",
			reqBody: loginRequest{
				Username: "abc@mail.com",
				Password: "wrongpass",
			},
			mockLoginResp:  nil,
			mockLoginErr:   auth.ErrNotYetActivatingAccount,
			wantStatusCode: http.StatusBadRequest,
			wantRespBody:   "the user has not yet activated the account",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAuth := new(MockAuth)
			if tc.mockLoginResp != nil || tc.mockLoginErr != nil {
				loginReq := tc.reqBody.(loginRequest)
				mockAuth.On("Login", mock.Anything, loginReq.Username, loginReq.Password).
					Return(tc.mockLoginResp, tc.mockLoginErr)
			}

			validate := validator.New()

			p := &Public{
				auth:      mockAuth,
				validator: validate,
			}

			router := gin.New()
			router.POST("/login", p.login)

			reqBodyJSON, err := json.Marshal(tc.reqBody)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBodyJSON))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatusCode, w.Code)
			if tc.wantStatusCode == http.StatusOK {
				cookies := w.Result().Cookies()

				var tokenCookie *http.Cookie
				for _, c := range cookies {
					if c.Name == "token" {
						tokenCookie = c
						break
					}
				}
				assert.NotNil(t, tokenCookie, "expected token cookie to be set")
				assert.Equal(t, "/", tokenCookie.Path)
				assert.True(t, tokenCookie.Secure)
				assert.True(t, tokenCookie.HttpOnly)
				assert.Greater(t, tokenCookie.MaxAge, 0)
			}

			assert.Contains(t, w.Body.String(), tc.wantRespBody)

			mockAuth.AssertExpectations(t)
		})
	}
}
