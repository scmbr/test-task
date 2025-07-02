package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scmbr/test-task/internal/dto"
)

func (h *Handler) generateTokens(c *gin.Context) {
	var req dto.GenerateTokensRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := h.service.Token.GenerateAccessToken(req.GUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	refreshToken, err := h.service.Token.GenerateAndSaveRefreshToken(req.GUID, c.GetHeader("User-Agent"), c.ClientIP())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate and save refresh token"})
		return
	}

	c.JSON(http.StatusOK, dto.TokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
func (h *Handler) getCurrentUserGUID(c *gin.Context) {
	c.JSON(http.StatusOK, "todo")
}
func (h *Handler) refreshTokens(c *gin.Context) {
	c.JSON(http.StatusOK, "todo")
}
func (h *Handler) logOut(c *gin.Context) {
	c.JSON(http.StatusOK, "todo")
}
