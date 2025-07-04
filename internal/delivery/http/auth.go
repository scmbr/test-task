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
func (h *Handler) refreshTokens(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userAgent := c.GetHeader("User-Agent")
	clientIP := c.ClientIP()

	newTokens, err := h.service.RefreshTokenPair(req.RefreshToken, req.AccessToken, userAgent, clientIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newTokens)
}
func (h *Handler) getCurrentUserGUID(c *gin.Context) {
	userGUID, exists := c.Get("userGuid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authorized"})
		return
	}

	guid, ok := userGUID.(string)
	if !ok || guid == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_guid": guid})
}

func (h *Handler) logOut(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.Logout(req.AccessToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "succesfully logged out"})
}
