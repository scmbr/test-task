package models

import "time"

type RefreshToken struct {
	ID        uint   `gorm:"primaryKey"`
	UserGUID  string `gorm:"index;size:36"`
	TokenHash string `gorm:"size:255"`
	UserAgent string
	IP        string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	ExpiresAt time.Time
}
