package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"payd/services/role"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockShiftService struct {
	mock.Mock
}

func (m *MockShiftService) CreateNewShiftSchedule(ctx context.Context, roleID int, startTime, endTime time.Time) (int, error) {
	args := m.Called(mock.Anything, roleID, startTime, endTime)
	return args.Int(0), args.Error(1)
}

func TestCreateSchedule(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type request struct {
		RoleID    int       `json:"roleId"`
		StartTime time.Time `json:"startTime"`
		EndTime   time.Time `json:"endTime"`
	}

	now := time.Now().Round(0)
	tests := []struct {
		name           string
		body           request
		validRoles     []int
		mockReturnID   int
		mockReturnErr  error
		wantStatusCode int
		wantRespBody   string
	}{
		{
			name: "success",
			body: request{
				RoleID:    1,
				StartTime: now,
				EndTime:   now.Add(1 * time.Hour),
			},
			mockReturnID:   123,
			mockReturnErr:  nil,
			wantStatusCode: http.StatusOK,
			wantRespBody:   "schedule created successfully",
		},
		{
			name: "invalid role ID",
			body: request{
				RoleID:    999,
				StartTime: now,
				EndTime:   now.Add(1 * time.Hour),
			},
			wantStatusCode: http.StatusBadRequest,
			wantRespBody:   "roleId is not a valid role ID",
		},
		{
			name: "startTime after endTime",
			body: request{
				RoleID:    1,
				StartTime: now.Add(2 * time.Hour),
				EndTime:   now,
			},
			wantStatusCode: http.StatusBadRequest,
			wantRespBody:   "startTime must be before endTime",
		},
		{
			name: "internal error",
			body: request{
				RoleID:    1,
				StartTime: now,
				EndTime:   now.Add(1 * time.Hour),
			},
			mockReturnErr:  assert.AnError,
			wantStatusCode: http.StatusInternalServerError,
			wantRespBody:   "internal error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockShift := new(MockShiftService)
			mockRoleService := new(MockRoleService)
			mockRoleService.On("GetRoles").Return([]role.Role{{ID: 1}})

			a := &Admin{
				shift: mockShift,
				role:  mockRoleService,
			}

			if tc.wantStatusCode == http.StatusOK || tc.mockReturnErr != nil {
				mockShift.On("CreateNewShiftSchedule",
					mock.Anything,
					tc.body.RoleID,
					mock.MatchedBy(func(t time.Time) bool {
						return t.Equal(tc.body.StartTime)
					}),
					mock.MatchedBy(func(t time.Time) bool {
						return t.Equal(tc.body.EndTime)
					})).
					Return(tc.mockReturnID, tc.mockReturnErr)
			}

			router := gin.New()
			router.POST("/schedule", a.createNewShiftSchedule)

			bodyJSON, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/schedule", bytes.NewReader(bodyJSON))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantStatusCode, w.Code)
			assert.Contains(t, w.Body.String(), tc.wantRespBody)

			mockShift.AssertExpectations(t)
		})
	}
}
