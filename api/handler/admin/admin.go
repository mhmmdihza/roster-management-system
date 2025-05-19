package admin

import (
	"payd/middleware"
	"payd/services/auth"
	"payd/services/role"
	"payd/services/shift"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Admin struct {
	auth      auth.AuthInterface
	role      role.RoleManagerInterface
	shift     shift.ShiftInterface
	validator *validator.Validate
}

type Option func(*Admin) error

func NewAdminHandler(router *gin.RouterGroup, opts ...Option) error {
	admin := &Admin{}
	for _, opt := range opts {
		if err := opt(admin); err != nil {
			return err
		}
	}
	router.Use(middleware.JWTAuthorizeRoles(admin.auth, "admin"))
	router.POST("/register", admin.register)
	router.GET("/list-role", admin.listRole)
	router.POST("/schedules", admin.createNewShiftSchedule)

	return nil
}

func WithShiftSvc(shift shift.ShiftInterface) Option {
	return func(s *Admin) error {
		s.shift = shift
		return nil
	}
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
