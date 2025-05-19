package public

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (p *Public) logout(c *gin.Context) {
	c.SetCookie(
		p.cookie.name,
		"",
		-1,
		p.cookie.path,
		p.cookie.domain,
		p.cookie.secure,
		p.cookie.httpOnly,
	)

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
