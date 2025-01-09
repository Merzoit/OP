package test

import (
	"at/internal/services/code"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func TestDbCodeCreate(db *pgxpool.Pool) error {
	codeRepo := code.NewPgCodeRepository(db)

	c := &code.Code{
		AccessCode:      111111,
		Title:           "TEST",
		Year:            0000,
		Description:     "TEST",
		AddedByWorkerID: 2,
	}

	err := codeRepo.CreateCode(c)
	if err != nil {
		log.Warn().Msgf("TEST FAILED: %v", err)
		return fmt.Errorf("TEST FAILED: %v", err)
	}

	log.Info().Msg("TEST SUCCESS")
	return nil
}

func TestDbCodeGet(db *pgxpool.Pool) error {
	codeRepo := code.NewPgCodeRepository(db)

	id := 111111

	u, err := codeRepo.GetCode(id)
	if err != nil {
		log.Warn().Msgf("TEST FAILED: %v", err)
		return fmt.Errorf("TEST FAILED: %v", err)
	}

	log.Info().Msgf("TEST SUCCESS: %v", u)
	return nil
}

func TestDbCodeDelete(db *pgxpool.Pool) error {
	codeRepo := code.NewPgCodeRepository(db)

	id := 111111

	err := codeRepo.DeleteCode(id)
	if err != nil {
		log.Warn().Msgf("TEST FAILED: %v", err)
		return fmt.Errorf("TEST FAILED: %v", err)
	}

	log.Info().Msg("TEST SUCCESS")
	return nil
}
