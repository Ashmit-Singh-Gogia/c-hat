package models

import "time"

type ChatParticipant struct {
	ID       uint      `json:"id"`
	ChatID   uint      `json:"chat_id"`
	UserID   uint      `json:"user_id"`
	JoinedAt time.Time `gorm:"autoCreateTime" json:"joined_at"`
	User     User      `gorm:"foreignKey:UserID" json:"user"`
}
