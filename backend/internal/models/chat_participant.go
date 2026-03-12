package models

import "time"

type ChatParticipant struct {
	ID       uint
	ChatID   uint
	UserID   uint
	JoinedAt time.Time `gorm:"autoCreateTime"`
}
