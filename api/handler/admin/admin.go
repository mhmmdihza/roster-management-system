package admin

import (
	"payd/services/auth"
	"payd/services/role"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Admin struct {
	auth      auth.AuthInterface
	role      role.RoleManagerInterface
	validator *validator.Validate
}

type Option func(*Admin) error

func NewAdminHandler(router *gin.RouterGroup, opts ...Option) error {
	public := &Admin{}
	for _, opt := range opts {
		if err := opt(public); err != nil {
			return err
		}
	}
	router.POST("/register", public.register)
	return nil
}

func WithAuthSvc(auth auth.AuthInterface) Option {
	return func(s *Admin) error {
		s.auth = auth
		return nil
	}
}

func WithRoleManager(role role.RoleManagerInterface) Option {
	return func(s *Admin) error {
		s.role = role
		return nil
	}
}

func WithValidator(validator *validator.Validate) Option {
	return func(s *Admin) error {
		s.validator = validator
		return nil
	}
}
