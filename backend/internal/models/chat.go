package models

import "time"

type Chat struct {
	ID           uint              `gorm:"primaryKey"`
	IsGroup      bool              `json:"is_group"`
	CreatedAt    time.Time         `gorm:"autoCreateTime"`
	Participants []ChatParticipant `gorm:"foreignKey:ChatID"`
}
