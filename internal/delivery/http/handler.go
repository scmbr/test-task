package http

import (
	"github.com/gin-gonic/gin"
	"github.com/scmbr/test-task/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.GET("/user", h.getGUID)
			auth.POST("/token", h.generateTokens)
			auth.POST("/refresh", h.refreshTokens)
			auth.POST("/logout", h.logOut)
		}
	}
	return router
}
