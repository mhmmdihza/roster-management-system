package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"payd/services/auth"
	"payd/services/role"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuth struct {
	mock.Mock
}

func (m *MockAuth) VerifySignatureJWT(tokenStr string) (*auth.Identity, error) {
	return nil, nil
}

func (m *MockAuth) Login(ctx context.Context, username, password string) (*auth.Identity, error) {
	return nil, nil
}

func (m *MockAuth) RegisterNewUser(ctx context.Context, email string, primaryRole int, roleAdmin bool) (string, error) {
	args := m.Called(ctx, email, primaryRole, roleAdmin)
	return args.String(0), args.Error(1)
}

type MockRoleService struct {
	mock.Mock
}

func (m *MockRoleService) GetRoles() []role.Role {
	args := m.Called()
	return args.Get(0).([]role.Role)
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		reqBody        interface{}
		mockRoles      []role.Role
		mockReturnID   string
		mockReturnErr  error
		wantStatusCode int
		wantRespBody   string
	}{
		{
			name: "success admin",
			reqBody: CreateUserRequest{
				Email:     "admin@example.com",
				RoleAdmin: true,
			},
			mockReturnID:   "new-admin-id",
			mockReturnErr:  nil,
			wantStatusCode: http.StatusOK,
			wantRespBody:   "user created successfully",
		},
		{
			name: "success employee",
			reqBody: CreateUserRequest{
				Email:       "user@example.com",
				PrimaryRole: ptrInt(2),
				RoleAdmin:   false,
			},
			mockRoles:      []role.Role{{ID: 2}},
			mockReturnID:   "new-user-id",
			mockReturnErr:  nil,
			wantStatusCode: http.StatusOK,
			wantRespBody:   "user created successfully",
		},
		{
			name: "validation error: email required",
			reqBody: map[string]interface{}{
				"roleAdmin": true,
			},
			wantStatusCode: http.StatusBadRequest,
			wantRespBody:   "Key: 'CreateUserRequest.Email'",
		},
		{
			name: "roleAdmin true but primaryRole set",
			reqBody: CreateUserRequest{
				Email:       "admin@example.com",
				RoleAdmin:   true,
				PrimaryRole: ptrInt(1),
			},
			wantStatusCode: http.StatusBadRequest,
			wantRespBody:   "primaryRole must not be defined when roleAdmin is true",
		},
		{
			name: "roleAdmin false but no primaryRole",
			reqBody: CreateUserRequest{
				Email:     "user@example.com",
				RoleAdmin: false,
			},
			wantStatusCode: http.StatusBadRequest,
			wantRespBody:   "primaryRole is required when roleAdmin is false",
		},
		{
			name: "invalid primaryRole",
			reqBody: CreateUserRequest{
				Email:       "user@example.com",
				PrimaryRole: ptrInt(999),
				RoleAdmin:   false,
			},
			mockRoles:      []role.Role{{ID: 1}, {ID: 2}},
			wantStatusCode: http.StatusBadRequest,
			wantRespBody:   "primaryRole is not a valid role ID",
		},
		{
			name: "conflict user already exists",
			reqBody: CreateUserRequest{
				Email:       "exists@example.com",
				PrimaryRole: ptrInt(2),
				RoleAdmin:   false,
			},
			mockRoles:      []role.Role{{ID: 2}},
			mockReturnErr:  auth.ErrAlreadyExists,
			wantStatusCode: http.StatusConflict,
			wantRespBody:   "already exists",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockAuth := new(MockAuth)
			mockRoleService := new(MockRoleService)

			if tc.mockRoles != nil {
				mockRoleService.On("GetRoles").Return(tc.mockRoles)
			}

			if tc.wantStatusCode == http.StatusOK || tc.mockReturnErr != nil {
				req := tc.reqBody.(CreateUserRequest)
				role := 0
				if req.PrimaryRole != nil {
					role = *req.PrimaryRole
				}
				mockAuth.On("RegisterNewUser", mock.Anything, req.Email, role, req.RoleAdmin).
					Return(tc.mockReturnID, tc.mockReturnErr)
			}

			a := &Admin{
				auth: mockAuth,
				role: mockRoleService,
			}

			router := gin.New()
			router.POST("/register", a.register)

			reqBodyJSON, err := json.Marshal(tc.reqBody)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(reqBodyJSON))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatusCode, w.Code)
			assert.Contains(t, w.Body.String(), tc.wantRespBody)

			mockAuth.AssertExpectations(t)
			mockRoleService.AssertExpectations(t)
		})
	}
}

func ptrInt(i int) *int {
	return &i
}
