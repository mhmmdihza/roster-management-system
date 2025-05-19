package admin

import (
	"net/http"
	"payd/util"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateNewShiftScheduleRequest struct {
	RoleID    int       `json:"roleId" binding:"required"`
	StartTime time.Time `json:"startTime" binding:"required"`
	EndTime   time.Time `json:"endTime" binding:"required"`
}

func (a *Admin) createNewShiftSchedule(c *gin.Context) {
	ctx := c.Request.Context()
	log := util.Log().WithContext(ctx)

	var req CreateNewShiftScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !a.isValidRoleID(req.RoleID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "roleId is not a valid role ID"})
		return
	}

	if !req.StartTime.Before(req.EndTime) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "startTime must be before endTime"})
		return
	}

	// Pretend to create and return a schedule ID
	scheduleID, err := a.shift.CreateNewShiftSchedule(ctx, req.RoleID, req.StartTime, req.EndTime)
	if err != nil {
		log.WithError(err).Error("create schedule")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "schedule created successfully",
		"id":      scheduleID,
	})
}
