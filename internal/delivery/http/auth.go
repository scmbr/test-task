package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scmbr/test-task/internal/dto"
)

func (h *Handler) generateTokens(c *gin.Context) {
	var req dto.GenerateTokensRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	accessToken, err := h.service.Token.GenerateAccessToken(req.UserGUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	refreshToken, err := h.service.Token.GenerateAndSaveRefreshToken(req.UserGUID, c.GetHeader("User-Agent"), c.ClientIP())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
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
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	userAgent := c.GetHeader("User-Agent")
	clientIP := c.ClientIP()

	newTokens, err := h.service.Token.RefreshTokenPair(req.RefreshToken, req.AccessToken, userAgent, clientIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.TokensResponse{AccessToken: newTokens.AccessToken, RefreshToken: newTokens.RefreshToken})
}
func (h *Handler) getCurrentUserGUID(c *gin.Context) {
	userGUID, exists := c.Get("userGuid")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "user not authorized"})
		return
	}

	userGuid, ok := userGUID.(string)
	if !ok || userGuid == "" {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "invalid user data"})
		return
	}

	c.JSON(http.StatusOK, dto.UserGUIDResponse{UserGUID: userGuid})
}
func (h *Handler) logOut(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	if err := h.service.Logout(req.AccessToken); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "succesfully logged out"})

}
