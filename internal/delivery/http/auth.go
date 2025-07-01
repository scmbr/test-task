package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) generateTokens(c *gin.Context) {
	c.JSON(http.StatusOK, "todo")
}
func (h *Handler) getGUID(c *gin.Context) {
	c.JSON(http.StatusOK, "todo")
}
func (h *Handler) refreshTokens(c *gin.Context) {
	c.JSON(http.StatusOK, "todo")
}
func (h *Handler) logOut(c *gin.Context) {
	c.JSON(http.StatusOK, "todo")
}
