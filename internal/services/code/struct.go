package code

import (
	"at/internal/services/worker"
	"time"
)

type Code struct {
	Id              int           `json:"id"`
	AccessCode      int           `json:"access_code"`
	Title           string        `json:"title"`
	Year            int           `json:"year"`
	Description     string        `json:"description"`
	AddedByWorkerID uint          `json:"added_by_worker_id"`
	RequestCount    int           `json:"request_count"`
	CreatedAt       time.Time     `json:"created_at"`
	AddedByWorker   worker.Worker `json:"worker"`
}
