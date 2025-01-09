package subscribe

type Subscribe struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	SponsorID    int    `json:"channel_id"`
	SubscribedAt string `json:"subscribed_at"`
}
