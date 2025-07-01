package models

import "time"

type User struct {
	GUID      string    `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
