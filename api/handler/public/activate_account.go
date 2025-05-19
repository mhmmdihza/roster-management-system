package public

import (
	"net/http"
	"payd/services/auth"
	"payd/util"

	"github.com/gin-gonic/gin"
)

type ActivateAccountRequest struct {
	ID       string `json:"id" binding:"required,uuid"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

func (p *Public) activateAccount(c *gin.Context) {
	ctx := c.Request.Context()
	log := util.Log().WithContext(ctx)

	var req ActivateAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := p.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := p.auth.ActivateNewUser(ctx, req.ID, req.Name, req.Password)
	if err != nil {
		switch err {
		case auth.ErrInvalidPassword:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case auth.ErrAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case auth.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			log.WithError(err).Error("login")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account activated successfully",
	})
}
