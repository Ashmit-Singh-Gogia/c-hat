package models

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Provider   string    `gorm:"uniqueIndex:idx_provider_id;not null" json:"provider"`
	ProviderID string    `gorm:"uniqueIndex:idx_provider_id;not null" json:"provider_id"`
	Username   string    `gorm:"not null" json:"username"`
	Email      string    `gorm:"not null" json:"email"`
	Avatar     string    `json:"avatar"`
	Password   string    `json:"-"` // Leave empty for OAuth, use later for Email/Pass
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
