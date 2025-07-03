package service

import (
	"time"

	"github.com/scmbr/test-task/internal/repository"
	"github.com/scmbr/test-task/pkg/auth"
	"github.com/scmbr/test-task/pkg/hasher"
)

type Token interface {
	GenerateAccessToken(GUID string) (string, error)
	GenerateAndSaveRefreshToken(guid, userAgent, ip string) (string, error)
}
type TokenInfo struct {
	HashedRefresh string
	UserAgent     string
	IP            string
}
type Service struct {
	Token
}
type Deps struct {
	Repos           *repository.Repository
	Hasher          hasher.Hasher
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	TokenManager    auth.TokenManager
	WebhookUrl      string
}

func NewServices(deps Deps) *Service {
	tokenService := NewTokenService(deps.Repos.RefreshToken, deps.Hasher, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.TokenManager)

	return &Service{
		Token: tokenService,
	}
}
