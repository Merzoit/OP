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

type UserClient struct {
	baseURL string
	client  *http.Client
}

func NewUserClient(baseURL string) *UserClient {
	return &UserClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type Role struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Permissions string `json:"permissions"`
}

type User struct {
	ID         uint      `json:"id"`
	TelegramID int64     `json:"telegram_id"`
	Username   string    `json:"username"`
	RoleID     uint      `json:"role_id"`
	IsBanned   bool      `json:"is_banned"`
	CreatedAt  time.Time `json:"created_at"`
	Role       Role      `json:"role"`
}

func (c *UserClient) CreateUser(user *User) error {
	url := fmt.Sprintf("%s%s", c.baseURL, "user")

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Warn().Msgf(constants.ErrUserCreate, resp.Status)
		return fmt.Errorf(constants.ErrUserCreate, resp.Status)
	}

	return nil
}

func (c *UserClient) GetUser(id int64) (*User, error) {
	url := fmt.Sprintf("%suser/%v", c.baseURL, id)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user: status %s", resp.Status)
	}

	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &user, nil
}
