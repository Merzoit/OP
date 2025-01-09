package worker

import (
	user "at/internal/services/user"

	"time"
)

type Worker struct {
	ID          uint       `json:"id"`
	UserID      uint       `json:"user_id"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	PaymentRate float64    `json:"payment_rate"`
	Balance     float64    `json:"balance"`
	User        *user.User `json:"user"`
}
