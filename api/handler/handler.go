package handler

import (
	"payd/handler/public"
	"payd/services/auth"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	*gin.Engine
	auth      auth.AuthInterface
	validator *validator.Validate
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
	return handler, nil
}

func WithAuthSvc(auth auth.AuthInterface) Option {
	return func(s *Handler) error {
		s.auth = auth
		return nil
	}
}

func WithValidator(validator *validator.Validate) Option {
	return func(s *Handler) error {
		s.validator = validator
		return nil
	}
}
