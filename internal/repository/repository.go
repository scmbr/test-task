package repository

import "gorm.io/gorm"

type User interface {
}
type RefreshToken interface {
}
type WebhookLog interface {
}
type Repository struct {
	User
	RefreshToken
	WebhookLog
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:         NewUserRepo(db),
		RefreshToken: NewRefreshTokenRepo(db),
		WebhookLog:   NewWebhookLogRepo(db),
	}
}
