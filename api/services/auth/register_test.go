package auth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	kratos "github.com/ory/kratos-client-go"
	"github.com/stretchr/testify/assert"
)

type registerNewUserScenarioFuncParam struct {
	email       string
	primaryRole int
	roleAdmin   bool
}

func scsRegisterNewUserAPIReqBody() kratos.CreateIdentityBody {
	pass := "123456"
	inactiveState := "inactive"
	traits := map[string]interface{}{
		"email":        "abc@gmail.com",
		"role":         "admin",
		"primary_role": 1,
	}
	requestBody := kratos.CreateIdentityBody{
		Credentials: &kratos.IdentityWithCredentials{
			Password: &kratos.IdentityWithCredentialsPassword{
				Config: &kratos.IdentityWithCredentialsPasswordConfig{
					Password: &pass,
				},
			},
		},
		SchemaId: "default",
		Traits:   traits,
		State:    &inactiveState,
	}
	return requestBody
}

func scsRegisterNewUserAPIRespBody() *kratos.Identity {
	identityResponse := &kratos.Identity{
		Id:     "10",
		Traits: make(map[string]string),
	}
	return identityResponse
}

var registerNewUserScenario = []struct {
	name               string
	mockApiResponse    mockApiResponse
	expectedApiRequest expectedApiRequest
	expectedError      error
	expectedId         string
	funcParams         registerNewUserScenarioFuncParam
}{
	{
		name:       "success register admin",
		expectedId: "10",
		funcParams: registerNewUserScenarioFuncParam{
			"abc@gmail.com", 1, true,
		},
		mockApiResponse: mockApiResponse{
			body:       scsRegisterNewUserAPIRespBody(),
			statusCode: 200,
		},
		expectedApiRequest: expectedApiRequest{
			method: "POST",
			path:   "/admin/identities",
			body:   scsRegisterNewUserAPIReqBody(),
		},
	},
	{
		name:       "success register employee",
		expectedId: "10",
		funcParams: registerNewUserScenarioFuncParam{
			"abc@gmail.com", 1, false,
		},
		mockApiResponse: mockApiResponse{
			body: kratos.Identity{
				Id:     "10",
				Traits: make(map[string]string),
			},
			statusCode: 200,
		},
		expectedApiRequest: expectedApiRequest{
			method: "POST",
			path:   "/admin/identities",
			body: func() kratos.CreateIdentityBody {
				requestBody := scsRegisterNewUserAPIReqBody()
				requestBody.Traits["role"] = "employee"
				return requestBody
			}(),
		},
	},
	{
		name:          "already exists",
		expectedError: ErrAlreadyExists,
		funcParams: registerNewUserScenarioFuncParam{
			"abc@gmail.com", 1, true,
		},
		mockApiResponse: mockApiResponse{
			statusCode: 409,
		},
		expectedApiRequest: expectedApiRequest{
			method: "POST",
			path:   "/admin/identities",
		},
	},
	{
		name:          "invalid email",
		expectedError: ErrInvalidEmail,
		funcParams: registerNewUserScenarioFuncParam{
			"aa", 1, true,
		},
		mockApiResponse: mockApiResponse{
			statusCode: 400,
		},
		expectedApiRequest: expectedApiRequest{
			method: "POST",
			path:   "/admin/identities",
		},
	},
}

func TestRegisterNewUser(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	for _, sc := range registerNewUserScenario {
		t.Run(sc.name, func(t *testing.T) {
			server := mockAPIServer(t, []expectedApiRequest{sc.expectedApiRequest}, []mockApiResponse{sc.mockApiResponse})
			defer server.Close()
			auth, err := NewAuth(WithKratosAdminURL(server.URL))
			assert.NoError(t, err)
			id, err := auth.RegisterNewUser(ctx, sc.funcParams.email, sc.funcParams.primaryRole, sc.funcParams.roleAdmin)
			assert.Equal(t, sc.expectedError, err)
			assert.Equal(t, sc.expectedId, id)
		})
	}
}

type expectedApiRequest struct {
	method     string
	path       string
	body       interface{}
	headers    map[string]string
	queryParam map[string]string
}

type mockApiResponse struct {
	body       interface{}
	statusCode int
	headers    map[string]string
}

// mockAPIServer creates a mock HTTP server for testing client-side API interactions.
// it verifies that incoming requests match the expected method, path, query parameters,
// headers, and body (if provided), and responds with predefined mock responses.
//
// Parameters:
// - expectedRequests: a slice of expectedApiRequest structs defining how each incoming request should look.
// - mockResponses: a slice of mockApiResponse structs specifying how the server should respond to each request.
//
// Returns:
// - an *httptest.Server that can be used as a stand-in for an actual API during tests.
//
// note: requests are matched and validated in the order they are received.
func mockAPIServer(
	t *testing.T,
	expectedRequests []expectedApiRequest,
	mockResponses []mockApiResponse,
) *httptest.Server {
	t.Helper()
	counter := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			counter++
		}()
		expectedRequest := expectedRequests[counter]
		mockResponse := mockResponses[counter]

		// Verify method
		assert.Equal(t, expectedRequest.method, r.Method)

		// Verify path
		assert.Equal(t, expectedRequest.path, r.URL.Path)

		// Verify query param
		for key, expectedValue := range expectedRequest.queryParam {
			assert.Equal(t, expectedValue, r.URL.Query().Get(key))
		}

		// Verify headers
		for key, expectedValue := range expectedRequest.headers {
			assert.Equal(t, expectedValue, r.Header.Get(key))
		}

		// Verify body if expected
		if expectedRequest.body != nil {
			defer r.Body.Close()
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("Failed reading body: %v", err)
			}

			actual := strings.TrimSpace(string(bodyBytes))

			expectedBytes, err := json.Marshal(expectedRequest.body)
			if err != nil {
				t.Errorf("Failed to marshal expected body: %v", err)
			}

			expected := strings.TrimSpace(string(expectedBytes))
			assert.Equal(t, expected, actual)
		}

		for keyHeader, valueHeader := range mockResponse.headers {
			w.Header().Set(keyHeader, valueHeader)
		}
		if mockResponse.statusCode == 0 {
			return
		}

		// Write mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(mockResponse.statusCode)
		if mockResponse.body != nil {
			var responseBytes []byte
			var err error

			switch v := mockResponse.body.(type) {
			case string:
				// assume it's already a JSON string (or plain text)
				responseBytes = []byte(v)
			default:
				// marshal the struct/map/etc. to JSON
				responseBytes, err = json.Marshal(v)
				if err != nil {
					t.Errorf("Failed to marshal mock response body: %v", err)
					return
				}
			}

			_, _ = w.Write(responseBytes)
		}
	}))
	return server
}
