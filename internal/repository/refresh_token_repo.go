package repository

import (
	"fmt"
	"time"

	"github.com/scmbr/test-task/internal/models"
	"gorm.io/gorm"
)

type RefreshTokenRepo struct {
	db *gorm.DB
}

func NewRefreshTokenRepo(db *gorm.DB) *RefreshTokenRepo {
	return &RefreshTokenRepo{db: db}
}
func (r *RefreshTokenRepo) SaveRefreshToken(guid, hashedToken, userAgent, ip string, refreshTokenTTL time.Duration) error {
	token := models.RefreshToken{
		UserGUID:  guid,
		TokenHash: hashedToken,
		UserAgent: userAgent,
		IP:        ip,
		ExpiresAt: time.Now().Add(refreshTokenTTL),
	}
	if err := r.db.Create(&token).Error; err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}

	return nil
}
