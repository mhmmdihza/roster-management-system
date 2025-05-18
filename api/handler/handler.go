package handler

import "github.com/gin-gonic/gin"

type Handler struct {
	*gin.Engine
}

func NewHandler() *Handler {
	router := gin.Default()
	public := router.Group("/")
	public.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	return &Handler{router}
}
