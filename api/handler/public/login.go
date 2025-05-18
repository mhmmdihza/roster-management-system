package public

import (
	"net/http"
	"time"

	"payd/services/auth"
	"payd/util"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required,email"`
	Password string `json:"password" binding:"required,min=1"`
}

func (p *Public) login(c *gin.Context) {
	ctx := c.Request.Context()
	log := util.Log().WithContext(ctx)
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := p.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	identity, err := p.auth.Login(ctx, req.Username, req.Password)
	if err != nil {
		switch err {
		case auth.ErrInvalidCredential:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case auth.ErrNotYetActivatingAccount:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			log.WithContext(ctx).WithError(err).Error("login")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}
	token, err := identity.GenerateJWT(time.Minute * 15)
	if err != nil {
		log.WithContext(ctx).WithError(err).Error("generate jwt")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.SetCookie(
		"token", // cookie name
		token,   // value
		900,     // maxAge in seconds (15 minutes)
		"/",     // path
		"",      // domain (empty = current domain)
		true,    // secure (true = HTTPS only)
		true,    // httpOnly (not accessible via JS)
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
	})
}
