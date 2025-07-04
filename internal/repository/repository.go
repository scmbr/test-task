package repository

import (
	"time"

	"github.com/scmbr/test-task/internal/models"
	"gorm.io/gorm"
)

type RefreshToken interface {
	SaveRefreshToken(guid, hashedToken, userAgent, ip string, refreshTokenTTL time.Duration) error
	ValidateRefreshToken(guid, refreshTokenHash string) (*models.RefreshToken, error)
	DeleteAllUserRefreshTokens(guid string) error
	DeleteRefreshToken(guid, refreshTokenHash string) error
	GetUserRefreshTokens(guid string) ([]*models.RefreshToken, error)
}

type Repository struct {
	RefreshToken
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		RefreshToken: NewRefreshTokenRepo(db),
	}
}
