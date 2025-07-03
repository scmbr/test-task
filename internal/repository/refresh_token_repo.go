package repository

import (
	"errors"
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
func (r *RefreshTokenRepo) ValidateRefreshToken(refreshTokenHash string) (*models.RefreshToken, error) {
	var token models.RefreshToken

	err := r.db.
		Where("token_hash = ? AND expires_at > ?", refreshTokenHash, time.Now()).
		First(&token).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("refresh token not found or expired")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &token, nil
}
func (r *RefreshTokenRepo) DeleteAllUserRefreshTokens(guid string) error {
	if guid == "" {
		return errors.New("user GUID cannot be empty")
	}
	result := r.db.
		Where("user_guid = ?", guid).
		Delete(&models.RefreshToken{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete user's refresh tokens: %w", result.Error)
	}
	return nil
}
