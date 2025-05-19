package handler

import (
	"payd/handler/admin"
	"payd/handler/public"
	"payd/services/auth"
	"payd/services/role"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	*gin.Engine
	auth      auth.AuthInterface
	validator *validator.Validate
	role      role.RoleManagerInterface
}

type Option func(*Handler) error

func NewHandler(opts ...Option) (*Handler, error) {
	router := gin.Default()
	handler := &Handler{Engine: router}

	for _, opt := range opts {
		if err := opt(handler); err != nil {
			return nil, err
		}
	}

	if err := public.PublicHandler(router.Group("/"),
		public.WithAuthSvc(handler.auth), public.WithValidator(handler.validator)); err != nil {
		return nil, err
	}
	if err := admin.NewAdminHandler(router.Group("/admin"),
		admin.WithAuthSvc(handler.auth), admin.WithValidator(handler.validator),
		admin.WithRoleManager(handler.role)); err != nil {
		return nil, err
	}
	return handler, nil
}

func WithAuthSvc(auth auth.AuthInterface) Option {
	return func(s *Handler) error {
		s.auth = auth
		return nil
	}
}
func WithRoleManager(role role.RoleManagerInterface) Option {
	return func(s *Handler) error {
		s.role = role
		return nil
	}
}

func WithValidator(validator *validator.Validate) Option {
	return func(s *Handler) error {
		s.validator = validator
		return nil
	}
}

func WithAllowedOrigins(allowedOrigins []string) Option {
	return func(s *Handler) error {
		s.Use(cors.New(cors.Config{
			AllowOrigins:     allowedOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
		return nil
	}
}
