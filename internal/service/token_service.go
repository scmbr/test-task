package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/scmbr/test-task/internal/dto"
	"github.com/scmbr/test-task/internal/models"
	"github.com/scmbr/test-task/internal/notifier"
	"github.com/scmbr/test-task/internal/repository"
	"github.com/scmbr/test-task/pkg/auth"
	"github.com/scmbr/test-task/pkg/hasher"
	"github.com/sirupsen/logrus"
)

type TokenService struct {
	repo             repository.RefreshToken
	hasher           hasher.Hasher
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
	tokenManager     auth.TokenManager
	ipChangeNotifier *notifier.IPNotifier
}

func NewTokenService(repo repository.RefreshToken, hasher hasher.Hasher, accessTTL, refreshTTL time.Duration, tokenManager auth.TokenManager, ipChangeNotifier *notifier.IPNotifier) *TokenService {
	return &TokenService{
		repo:             repo,
		hasher:           hasher,
		tokenManager:     tokenManager,
		accessTokenTTL:   accessTTL,
		refreshTokenTTL:  refreshTTL,
		ipChangeNotifier: ipChangeNotifier,
	}
}

func (s *TokenService) GenerateAccessToken(GUID string) (string, error) {
	accessToken, err := s.tokenManager.NewJWT(GUID, s.accessTokenTTL)
	return accessToken, err
}
func (s *TokenService) GenerateAndSaveRefreshToken(guid, userAgent, ip string) (string, error) {
	rawBytes, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	hashedToken, err := s.hasher.Hash(rawBytes)
	if err != nil {
		return "", fmt.Errorf("failed to hash token: %w", err)
	}

	if err := s.repo.SaveRefreshToken(guid, hashedToken, userAgent, ip, s.refreshTokenTTL); err != nil {
		return "", fmt.Errorf("failed to save token: %w", err)
	}
	encodedToken := base64.StdEncoding.EncodeToString(rawBytes)
	return encodedToken, nil
}
func (s *TokenService) RefreshTokenPair(refreshTokenBase64, accessToken, userAgent, clientIP string) (*dto.TokensResponse, error) {
	decodedRefresh, err := base64.StdEncoding.DecodeString(refreshTokenBase64)
	if err != nil {
		return nil, fmt.Errorf("invalid base64 token: %w", err)
	}

	accessClaims, err := s.tokenManager.Parse(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}

	userTokens, err := s.repo.GetUserRefreshTokens(accessClaims.UserGUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tokens: %w", err)
	}

	var matchedToken *models.RefreshToken
	for _, token := range userTokens {
		if s.hasher.Verify(decodedRefresh, token.TokenHash) {
			matchedToken = token
			break
		}
	}
	if matchedToken == nil {
		return nil, errors.New("invalid refresh token")
	}

	if accessClaims.UserGUID != matchedToken.UserGUID {
		return nil, errors.New("token pair mismatch")
	}

	if matchedToken.UserAgent != userAgent {
		if err := s.repo.DeleteAllUserRefreshTokens(matchedToken.UserGUID); err != nil {
			logrus.Errorf("failed to delete tokens: %v", err)
		}
		return nil, errors.New("user agent changed - all user's tokens deleted")
	}

	if matchedToken.IP != clientIP {
		go s.ipChangeNotifier.NotifyChange(matchedToken.UserGUID, matchedToken.IP, clientIP)
	}

	newAccess, err := s.tokenManager.NewJWT(matchedToken.UserGUID, s.accessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshRaw, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	hashedToken, err := s.hasher.Hash(newRefreshRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to hash token: %w", err)
	}

	if err := s.repo.SaveRefreshToken(accessClaims.UserGUID, hashedToken, userAgent, clientIP, s.refreshTokenTTL); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	if err := s.repo.DeleteRefreshToken(matchedToken.UserGUID, matchedToken.TokenHash); err != nil {
		logrus.Errorf("failed to delete old refresh token: %v", err)
	}

	newRefreshBase64 := base64.StdEncoding.EncodeToString(newRefreshRaw)

	return &dto.TokensResponse{
		AccessToken:  newAccess,
		RefreshToken: newRefreshBase64,
	}, nil
}
func (s *TokenService) Logout(accessToken string) error {
	if accessToken == "" {
		return errors.New("empty access token")
	}
	accessClaims, err := s.tokenManager.Parse(accessToken)
	if err != nil {
		return fmt.Errorf("invalid access token: %w", err)
	}
	err = s.repo.DeleteAllUserRefreshTokens(accessClaims.UserGUID)
	if err != nil {
		return fmt.Errorf("failed to delete user's refresh tokens: %w", err)
	}
	return nil
}
