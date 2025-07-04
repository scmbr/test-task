package http

import (
	"github.com/gin-gonic/gin"
	"github.com/scmbr/test-task/internal/service"
	"github.com/scmbr/test-task/pkg/auth"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	service      *service.Service
	tokenManager auth.TokenManager
}

func NewHandler(service *service.Service, tokenManager auth.TokenManager) *Handler {
	return &Handler{service: service, tokenManager: tokenManager}
}
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.GET("/user", h.authMiddleware, h.getCurrentUserGUID)
			auth.POST("/token", h.generateTokens)
			auth.POST("/refresh", h.refreshTokens)
			auth.POST("/logout", h.logOut)
		}
	}
	return router
}
