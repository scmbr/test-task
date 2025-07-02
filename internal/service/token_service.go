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

func (r *TokenService) GenerateAccessToken(GUID string) (string, error) {
	accessToken, err := r.tokenManager.NewJWT(GUID, r.accessTokenTTL)
	return accessToken, err
}
func (r *TokenService) GenerateAndSaveRefreshToken(guid, userAgent, ip string) (string, error) {
	rawToken, err := r.tokenManager.NewRefreshToken()
	if err != nil {
		return "", err
	}
	hashedToken, err := r.hasher.Hash(rawToken)
	if err != nil {
		return "", err
	}
	if err := r.repo.SaveRefreshToken(guid, hashedToken, userAgent, ip, r.refreshTokenTTL); err != nil {
		return "", err
	}
	return rawToken, nil
}
