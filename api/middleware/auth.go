package middleware

import (
	"net/http"
	"payd/services/auth"

	"github.com/gin-gonic/gin"
)

func JWTAuthorizeRoles(authService auth.AuthInterface, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}

		identity, err := authService.VerifySignatureJWT(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// check if the user's role is allowed
		roleAllowed := false
		for _, role := range allowedRoles {
			if identity.Role == role {
				roleAllowed = true
				break
			}
		}
		if !roleAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient role"})
			return
		}

		c.Set("identity", identity)
		c.Next()
	}
}
