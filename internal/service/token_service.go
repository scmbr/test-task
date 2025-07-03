package service

import (
	"time"

	"github.com/scmbr/test-task/internal/repository"
	"github.com/scmbr/test-task/pkg/auth"
	"github.com/scmbr/test-task/pkg/hasher"
)

type TokenService struct {
	repo            repository.RefreshToken
	hasher          hasher.Hasher
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	tokenManager    auth.TokenManager
}

func NewTokenService(repo repository.RefreshToken, hasher hasher.Hasher, accessTTL, refreshTTL time.Duration, tokenManager auth.TokenManager) *TokenService {
	return &TokenService{
		repo:            repo,
		hasher:          hasher,
		tokenManager:    tokenManager,
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
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
