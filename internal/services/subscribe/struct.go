package subscribe

import "time"

type Subscribe struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	SponsorID    int       `json:"sponsor_id"`
	SubscribedAt time.Time `json:"subscribed_at"`
}
