package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	GoogleID  string    `gorm:"uniqueIndex;not null" json:"google_id"`
	Username  string    `gorm:"not null" json:"username"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
