package public

import (
	"payd/services/auth"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Public struct {
	auth      auth.AuthInterface
	validator *validator.Validate
	cookie    cookie
}

type cookie struct {
	name     string
	path     string
	domain   string
	secure   bool
	httpOnly bool
}

func cookieConfig() cookie {
	return cookie{
		name:     "token",
		path:     "/",
		domain:   "",
		secure:   true,
		httpOnly: true,
	}
}

type Option func(*Public) error

func PublicHandler(router *gin.RouterGroup, opts ...Option) error {
	public := &Public{
		cookie: cookieConfig(),
	}
	for _, opt := range opts {
		if err := opt(public); err != nil {
			return err
		}
	}
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	router.POST("/login", public.login)
	router.POST("/logout", public.logout)
	router.POST("/activate", public.activateAccount)
	return nil
}

func WithAuthSvc(auth auth.AuthInterface) Option {
	return func(s *Public) error {
		s.auth = auth
		return nil
	}
}

func WithValidator(validator *validator.Validate) Option {
	return func(s *Public) error {
		s.validator = validator
		return nil
	}
}
