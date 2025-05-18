package handler

import "github.com/gin-gonic/gin"

type Auth interface{}

type Handler struct {
	*gin.Engine
	auth Auth
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

	public := router.Group("/")
	public.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	return handler, nil
}

func WithAuthSvc(auth Auth) Option {
	return func(s *Handler) error {
		s.auth = auth
		return nil
	}
}
