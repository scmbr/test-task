package repository

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken interface {
	SaveRefreshToken(guid, hashedToken, userAgent, ip string, refreshTokenTTL time.Duration) error
}

type Repository struct {
	RefreshToken
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		RefreshToken: NewRefreshTokenRepo(db),
	}
}
