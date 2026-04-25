package models

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	GoogleID  string `gorm:"uniqueIndex;not null"`
	Username  string `gorm:"not null"`
	Email     string `gorm:"uniqueIndex;not null"`
	Avatar    string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
