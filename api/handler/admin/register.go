package admin

import (
	"net/http"
	"payd/services/auth"
	"payd/util"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Email       string `json:"email" binding:"required,email"`
	PrimaryRole *int   `json:"primaryRole"` // must be nil if RoleAdmin true
	RoleAdmin   bool   `json:"roleAdmin,omitempty"`
}

// only admins can register new users
func (a *Admin) register(c *gin.Context) {
	ctx := c.Request.Context()
	log := util.Log().WithContext(ctx)

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	primaryRole := 0
	// custom validations
	if req.RoleAdmin {
		if req.PrimaryRole != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "primaryRole must not be defined when roleAdmin is true"})
			return
		}
	} else {
		if req.PrimaryRole == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "primaryRole is required when roleAdmin is false"})
			return
		}
		if !a.isValidRoleID(*req.PrimaryRole) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "primaryRole is not a valid role ID"})
			return
		}
		primaryRole = *req.PrimaryRole
	}

	userid, err := a.auth.RegisterNewUser(ctx, req.Email, primaryRole, req.RoleAdmin)
	if err != nil {
		switch err {
		case auth.ErrAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			log.WithError(err).Error("register")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user created successfully",
		"id":      userid,
	})
}

func (a *Admin) isValidRoleID(id int) bool {
	for _, role := range a.role.GetRoles() {
		if role.ID == id {
			return true
		}
	}
	return false
}
