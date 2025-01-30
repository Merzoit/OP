package clients

import (
	"PB/constants"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type Worker struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	PaymentRate float64   `json:"payment_rate"`
	Balance     float64   `json:"balance"`
	User        *User     `json:"user"`
}

type WorkerClient struct {
	baseURL string
	client  *http.Client
}

func NewWorkerClient(baseURL string) *WorkerClient {
	return &WorkerClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *WorkerClient) CreateWorker(worker *Worker) error {
	url := fmt.Sprintf("%s%s", c.baseURL, "worker")

	data, err := json.Marshal(worker)
	if err != nil {
		return err
	}
	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Warn().Msgf(constants.ErrWorkerCreate, resp.Status)
		return fmt.Errorf(constants.ErrWorkerCreate, resp.Status)
	}

	return nil
}
