package reflink

import (
	"at/internal/services/worker"
	"time"
)

type ReferralLink struct {
	ID            uint          `json:"id"`
	WorkerID      uint          `json:"worker_id"`
	Link          string        `json:"referral_link"`
	Clicks        int           `json:"clicks"`
	Registrations int           `json:"registrations"`
	CreatedAt     time.Time     `json:"created_at"`
	Worker        worker.Worker `json:"worker"`
}
