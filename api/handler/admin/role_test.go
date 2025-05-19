package admin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"payd/services/role"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestListRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockRoles      []role.Role
		wantStatusCode int
		wantRespBody   []RoleResponse
	}{
		{
			name: "list roles successfully",
			mockRoles: []role.Role{
				{ID: 1, Name: "Admin"},
				{ID: 2, Name: "Employee"},
			},
			wantStatusCode: http.StatusOK,
			wantRespBody: []RoleResponse{
				{ID: 1, RoleName: "Admin"},
				{ID: 2, RoleName: "Employee"},
			},
		},
		{
			name:           "no roles",
			mockRoles:      []role.Role{},
			wantStatusCode: http.StatusOK,
			wantRespBody:   []RoleResponse{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRoleService := new(MockRoleService)
			mockRoleService.On("GetRoles").Return(tc.mockRoles)

			a := &Admin{
				role: mockRoleService,
			}

			router := gin.New()
			router.GET("/roles", a.listRole)

			req := httptest.NewRequest(http.MethodGet, "/roles", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatusCode, w.Code)

			var got []RoleResponse
			err := json.Unmarshal(w.Body.Bytes(), &got)
			assert.NoError(t, err)
			if got == nil {
				got = []RoleResponse{}
			}
			assert.Equal(t, tc.wantRespBody, got)

			mockRoleService.AssertExpectations(t)
		})
	}
}
