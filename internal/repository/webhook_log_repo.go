package repository

import "gorm.io/gorm"

type WebhookLogRepo struct {
	db *gorm.DB
}

func NewWebhookLogRepo(db *gorm.DB) *WebhookLogRepo {
	return &WebhookLogRepo{db: db}
}
