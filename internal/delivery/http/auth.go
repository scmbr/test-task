package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scmbr/test-task/internal/dto"
)

// @Summary Генерация пары токенов
// @Description Создает access и refresh токены для указанного пользователя
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.GenerateTokensRequest true "User GUID"
// @Success 200 {object} dto.TokensResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/auth/token [post]
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

// @Summary Обновление пары токенов
// @Description Обновляет access и refresh токены парой токенов, которая была выдана вместе
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Tokens"
// @Success 200 {object} dto.TokensResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/auth/refresh [post]
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

// @Summary Получение GUID текущего пользователя
// @Description Возвращает GUID авторизованного пользователя
// @Tags User
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} dto.UserGUIDResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/auth/user [get]
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

// @Summary Выход из системы
// @Description Удаляет все refresh-токены пользователя
// @Tags Auth
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param access_token body dto.LogoutRequest true "Access Token"
// @Success 200 {object} dto.MessageResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/auth/logout [post]
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
