package models

import "time"

type Chat struct {
	ID           uint              `gorm:"primaryKey" json:"id"`
	IsGroup      bool              `json:"is_group"`
	CreatedAt    time.Time         `gorm:"autoCreateTime" json:"created_at"`
	Participants []ChatParticipant `gorm:"foreignKey:ChatID" json:"participants"`
}
