package models

import "time"

type WebhookLog struct {
	ID        uint   `gorm:"primaryKey"`
	UserGUID  string `gorm:"index;size:36"`
	EventType string
	OldIP     string
	NewIP     string
	Status    string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
