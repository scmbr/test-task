package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/scmbr/test-task/internal/dto"
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
	rawToken, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return "", err
	}
	hashedToken, err := s.hasher.Hash(rawToken)
	if err != nil {
		return "", err
	}
	if err := s.repo.SaveRefreshToken(guid, hashedToken, userAgent, ip, s.refreshTokenTTL); err != nil {
		return "", err
	}
	return rawToken, nil
}
func (s *TokenService) RefreshTokenPair(refreshToken, accessToken, userAgent, clientIP string) (*dto.TokensResponse, error) {
	accessClaims, err := s.tokenManager.Parse(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %w", err)
	}
	refreshTokenHash, err := s.hasher.Hash(refreshToken)
	if err != nil {
		return &dto.TokensResponse{}, err
	}
	refreshTokenModel, err := s.repo.ValidateRefreshToken(refreshTokenHash)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if accessClaims.UserGUID != refreshTokenModel.UserGUID {
		return nil, errors.New("token pair mismatch")
	}

	if refreshTokenModel.UserAgent != userAgent {
		if err := s.repo.DeleteAllUserRefreshTokens(refreshTokenModel.UserGUID); err != nil {
			logrus.Errorf("failed to delete tokens: %v", err)
		}
		return nil, errors.New("user agent changed - all user's tokens deleted")
	}

	if refreshTokenModel.IP != clientIP {
		go s.ipChangeNotifier.NotifyChange(refreshTokenModel.UserGUID, refreshTokenModel.IP, clientIP)
	}

	newAccess, err := s.tokenManager.NewJWT(refreshTokenModel.UserGUID, s.accessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefresh, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	if err := s.repo.DeleteRefreshToken(refreshTokenHash); err != nil {
		logrus.Errorf("failed to delete old refresh token: %v", err)
	}

	return &dto.TokensResponse{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}
