package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) authMiddleware(c *gin.Context) {
	if c.Request.Method == "OPTIONS" {
		c.Next()
		return
	}
	header := c.GetHeader("Authorization")
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "empty auth header"})
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
		return
	}
	claims, err := h.tokenManager.Parse(headerParts[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
		return
	}

	c.Set("userGuid", claims.UserGUID)
	c.Next()
}
