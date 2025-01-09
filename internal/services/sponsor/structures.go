package sponsor

import "time"

type Sponsor struct {
	ID           int       `json:"id"`
	TelegramID   uint64    `json:"telegram_id"`
	TelegramLink string    `json:"telegram_link"`
	PricePerSub  float64   `json:"price_per_sub"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
}

type SponsorPayment struct {
	ID          uint      `json:"id"`
	SponsorID   uint      `json:"sponsor_id"`
	AmountDue   float64   `json:"amount_due"`
	Subscribers int       `json:"subscribers"`
	Paid        bool      `json:"paid"`
	CreatedAt   time.Time `json:"created_at"`
	Sponsor     Sponsor   `json:"sponsor"`
}
