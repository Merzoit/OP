package services

import (
	"at/internal/services/role"
	"time"
)

type User struct {
	ID         uint      `json:"id"`
	TelegramID int64     `json:"telegram_id"`
	Username   string    `json:"username"`
	RoleID     uint      `json:"role_id"`
	IsBanned   bool      `json:"is_banned"`
	CreatedAt  time.Time `json:"created_at"`
	Role       role.Role `json:"role"`
}
