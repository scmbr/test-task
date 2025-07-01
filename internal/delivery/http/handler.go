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
			auth.GET("/user", h.getGUID)
			auth.POST("/token", h.generateTokens)
			auth.POST("/refresh", h.refreshTokens)
			auth.POST("/logout", h.logOut)
		}
	}
	return router
}
