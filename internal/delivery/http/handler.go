package http

import "github.com/gin-gonic/gin"

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.GET("/token")
			auth.GET("/user")
			auth.POST("/refresh")
			auth.POST("/logout")
		}
	}
	return router
}
