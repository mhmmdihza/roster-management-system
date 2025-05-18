package auth

import (
	"context"
	"fmt"
	st "payd/storage"
	"testing"

	kratos "github.com/ory/kratos-client-go"
	"github.com/stretchr/testify/assert"
)

// SelectEmployeeByID implements storage.
func (m *mockStorage) SelectEmployeeByID(ctx context.Context, id int) (*st.Employee, error) {
	if m.selectEmployeeByIDFunc != nil {
		return m.selectEmployeeByIDFunc(ctx, id)
	}
	return nil, fmt.Errorf("not implemented")
}

type loginFuncParam struct {
	username string
	password string
}

var loginApiRequests = []expectedApiRequest{
	{
		method: "GET",
		path:   "/self-service/login/api",
	},
	{
		method: "POST",
		path:   "/self-service/login",
		body: kratos.UpdateLoginFlowBody{
			UpdateLoginFlowWithPasswordMethod: &kratos.UpdateLoginFlowWithPasswordMethod{
				Identifier: "abc@mail.com",
				Password:   "password123",
				Method:     "password",
			},
		},
	},
}

var loginApiResponses = []mockApiResponse{
	{
		statusCode: 200,
		body: kratos.LoginFlow{
			Id:    "flow-id-123",
			State: "password",
		},
	},
	{
		statusCode: 200,
		body: kratos.SuccessfulNativeLogin{
			Session: kratos.Session{
				Identity: &kratos.Identity{
					Id:    "userid",
					State: func() *string { s := "active"; return &s }(),
					Traits: map[string]interface{}{
						"employee_id": "4",
						"email":       "abc@mail.com",
						"role":        "employee",
					},
				},
			},
		},
	},
}

var loginScenarios = []struct {
	name                string
	mockApiResponses    []mockApiResponse
	expectedApiRequests []expectedApiRequest
	expectedError       error
	funcParams          loginFuncParam
	mockStorage         *mockStorage
}{
	{
		name: "success",
		funcParams: loginFuncParam{
			username: "abc@mail.com",
			password: "password123",
		},
		expectedApiRequests: loginApiRequests,
		mockApiResponses:    loginApiResponses,
		mockStorage: &mockStorage{
			selectEmployeeByIDFunc: func(ctx context.Context, id int) (*st.Employee, error) {
				return &st.Employee{
					ID:          4,
					Name:        "name",
					PrimaryRole: 7,
				}, nil
			},
		},
	},
	{
		name: "invalid credential",
		funcParams: loginFuncParam{
			username: "wrong@mail.com",
			password: "wrongpassword",
		},
		expectedApiRequests: func() []expectedApiRequest {
			expectedReq := make([]expectedApiRequest, len(loginApiRequests))
			copy(expectedReq, loginApiRequests)
			req := expectedReq[1]
			req.body = kratos.UpdateLoginFlowBody{
				UpdateLoginFlowWithPasswordMethod: &kratos.UpdateLoginFlowWithPasswordMethod{
					Identifier: "wrong@mail.com",
					Password:   "wrongpassword",
					Method:     "password",
				},
			}
			expectedReq[1] = req
			return expectedReq
		}(),
		mockApiResponses: func() []mockApiResponse {
			mockResp := make([]mockApiResponse, len(loginApiResponses))
			copy(mockResp, loginApiResponses)
			mockResp[1] = mockApiResponse{
				statusCode: 400,
				body: kratos.SuccessfulNativeLogin{
					Session: kratos.Session{
						Identity: &kratos.Identity{
							Id:    "userid",
							State: func() *string { s := "active"; return &s }(),
							Traits: map[string]interface{}{
								"employee_id": "4",
								"email":       "wrong@mail.com",
								"role":        "wrongpassword",
							},
						},
					},
				},
			}
			return mockResp
		}(),
		expectedError: ErrInvalidCredential,
		mockStorage:   &mockStorage{},
	},
	{
		name: "identity is nil",
		funcParams: loginFuncParam{
			username: "abc@mail.com",
			password: "password123",
		},
		expectedApiRequests: loginApiRequests,
		mockApiResponses: []mockApiResponse{
			{
				statusCode: 200,
				body: kratos.LoginFlow{
					Id:    "flow-id-123",
					State: "password",
				},
			},
			{
				statusCode: 200,
				body: kratos.SuccessfulNativeLogin{
					Session: kratos.Session{
						Identity: nil,
					},
				},
			},
		},
		expectedError: fmt.Errorf("unexpected identity nil"),
		mockStorage:   &mockStorage{},
	},
	{
		name: "not yet activating account",
		funcParams: loginFuncParam{
			username: "abc@mail.com",
			password: "password123",
		},
		expectedApiRequests: loginApiRequests,
		mockApiResponses: func() []mockApiResponse {
			mockResp := make([]mockApiResponse, len(loginApiResponses))
			copy(mockResp, loginApiResponses)
			mockResp[1] = mockApiResponse{
				statusCode: 200,
				body: kratos.SuccessfulNativeLogin{
					Session: kratos.Session{
						Identity: &kratos.Identity{
							Id:    "userid",
							State: func() *string { s := "inactive"; return &s }(),
							Traits: map[string]interface{}{
								"employee_id": "4",
								"email":       "abc@mail.com",
								"role":        "employee",
							},
						},
					},
				},
			}
			return mockResp
		}(),
		mockStorage: &mockStorage{
			selectEmployeeByIDFunc: func(ctx context.Context, id int) (*st.Employee, error) {
				return &st.Employee{
					ID:          4,
					Name:        "name",
					PrimaryRole: 7,
				}, nil
			},
		},
		expectedError: ErrNotYetActivatingAccount,
	},
	{
		name: "db error",
		funcParams: loginFuncParam{
			username: "abc@mail.com",
			password: "password123",
		},
		expectedApiRequests: loginApiRequests,
		mockApiResponses:    loginApiResponses,
		mockStorage: &mockStorage{
			selectEmployeeByIDFunc: func(ctx context.Context, id int) (*st.Employee, error) {
				return nil, fmt.Errorf("db error")
			},
		},
		expectedError: fmt.Errorf("db error"),
	},
}

func TestLogin(t *testing.T) {
	t.Parallel()
	for _, sc := range loginScenarios {
		t.Run(sc.name, func(t *testing.T) {
			ctx := context.Background()
			server := mockAPIServer(t, sc.expectedApiRequests, sc.mockApiResponses)
			defer server.Close()

			auth, err := NewAuth(sc.mockStorage, WithKratosPublicURL(server.URL))
			assert.NoError(t, err)

			identity, err := auth.Login(ctx, sc.funcParams.username, sc.funcParams.password)

			if sc.expectedError != nil {
				assert.Equal(t, sc.expectedError.Error(), err.Error())
				assert.Nil(t, identity)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, identity)
				assert.Equal(t, "userid", identity.ID)
				assert.Equal(t, "abc@mail.com", identity.Email)
				assert.Equal(t, "employee", identity.Role)
				assert.Equal(t, "4", identity.EmployeeId)
				assert.Equal(t, "name", identity.EmployeeName)
				assert.Equal(t, 7, identity.PrimaryRole)
			}
		})
	}
}
