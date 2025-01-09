package test

import (
	user "at/internal/services/user"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func TestDbUserCreate(db *pgxpool.Pool) error {
	userRepo := user.NewPgUserRepository(db)

	u := &user.User{
		TelegramID: 111111,
		Username:   "TEST",
	}

	err := userRepo.CreateUser(u)
	if err != nil {
		log.Warn().Msgf("TEST FAILED: %v", err)
		return fmt.Errorf("TEST FAILED: %v", err)
	}

	log.Info().Msg("TEST SUCCESS")
	return nil
}

func TestDbUserGet(db *pgxpool.Pool) error {
	userRepo := user.NewPgUserRepository(db)

	id := 111111

	u, err := userRepo.GetUser(uint64(id))
	if err != nil {
		log.Warn().Msgf("TEST FAILED: %v", err)
		return fmt.Errorf("TEST FAILED: %v", err)
	}

	log.Info().Msgf("TEST SUCCESS: %v", u)
	return nil
}
