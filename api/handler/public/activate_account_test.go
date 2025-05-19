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

func (m *MockAuth) ActivateNewUser(ctx context.Context, userId string, name string, password string) error {
	args := m.Called(ctx, userId, name, password)
	return args.Error(0)
}

func TestPublicActivateAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type activateRequest struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	tests := []struct {
		name           string
		reqBody        interface{}
		mockError      error
		wantStatusCode int
		wantRespBody   string
	}{
		{
			name: "success activation",
			reqBody: activateRequest{
				ID:       "f47ac10b-58cc-4372-a567-0e02b2c3d479",
				Name:     "name",
				Password: "securepassword",
			},
			mockError:      nil,
			wantStatusCode: http.StatusOK,
			wantRespBody:   "Account activated successfully",
		},
		{
			name: "validation error - missing name",
			reqBody: map[string]interface{}{
				"id":       "f47ac10b-58cc-4372-a567-0e02b2c3d479",
				"password": "securepassword",
			},
			wantStatusCode: http.StatusBadRequest,
			wantRespBody:   "Key: 'ActivateAccountRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag",
		},
		{
			name: "invalid password error",
			reqBody: activateRequest{
				ID:       "f47ac10b-58cc-4372-a567-0e02b2c3d479",
				Name:     "name",
				Password: "securepassword",
			},
			mockError:      auth.ErrInvalidPassword,
			wantStatusCode: http.StatusBadRequest,
			wantRespBody:   auth.ErrInvalidPassword.Error(),
		},
		{
			name: "already exists error",
			reqBody: activateRequest{
				ID:       "f47ac10b-58cc-4372-a567-0e02b2c3d479",
				Name:     "name",
				Password: "securepassword",
			},
			mockError:      auth.ErrAlreadyExists,
			wantStatusCode: http.StatusConflict,
			wantRespBody:   auth.ErrAlreadyExists.Error(),
		},
		{
			name: "not found error",
			reqBody: activateRequest{
				ID:       "f47ac10b-58cc-4372-a567-0e02b2c3d479",
				Name:     "name",
				Password: "securepassword",
			},
			mockError:      auth.ErrNotFound,
			wantStatusCode: http.StatusNotFound,
			wantRespBody:   auth.ErrNotFound.Error(),
		},
		{
			name: "unexpected internal error",
			reqBody: activateRequest{
				ID:       "f47ac10b-58cc-4372-a567-0e02b2c3d479",
				Name:     "name",
				Password: "securepassword",
			},
			mockError:      assert.AnError,
			wantStatusCode: http.StatusInternalServerError,
			wantRespBody:   "internal error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAuth := new(MockAuth)

			if tc.mockError != nil || tc.wantStatusCode == http.StatusOK {
				body := tc.reqBody.(activateRequest)
				mockAuth.On("ActivateNewUser", mock.Anything, body.ID, body.Name, body.Password).
					Return(tc.mockError)
			}

			validate := validator.New()
			p := &Public{
				auth:      mockAuth,
				validator: validate,
			}

			router := gin.New()
			router.POST("/activate", p.activateAccount)

			reqBodyJSON, err := json.Marshal(tc.reqBody)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/activate", bytes.NewBuffer(reqBodyJSON))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatusCode, w.Code)
			assert.Contains(t, w.Body.String(), tc.wantRespBody)

			mockAuth.AssertExpectations(t)
		})
	}
}
