package test

import (
	"at/internal/services/sponsor"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func TestDbSponsorCreate(db *pgxpool.Pool) error {
	sponsorRepo := sponsor.NewPgSponsorRepository(db)

	s := &sponsor.Sponsor{
		TID:     111111,
		Link:    "test",
		Price:   8,
		Comment: "TEST",
		Active:  false,
		Name:    "tEsT",
	}

	err := sponsorRepo.CreateSponsor(s)
	if err != nil {
		log.Warn().Msgf("TEST FAILED: %v", err)
		return fmt.Errorf("TEST FAILED: %v", err)
	}

	log.Info().Msg("TEST SUCCESS")
	return nil
}

func TestDbSponsorGet(db *pgxpool.Pool) error {
	sponsorRepo := sponsor.NewPgSponsorRepository(db)

	id := 111111

	u, err := sponsorRepo.GetSponsor(uint64(id))
	if err != nil {
		log.Warn().Msgf("TEST FAILED: %v", err)
		return fmt.Errorf("TEST FAILED: %v", err)
	}

	log.Info().Msgf("TEST SUCCESS: %v", u)
	return nil
}

func TestDbSponsorsGet(db *pgxpool.Pool) error {
	sponsorRepo := sponsor.NewPgSponsorRepository(db)

	s, err := sponsorRepo.GetSponsors()
	if err != nil {
		log.Warn().Msgf("TEST FAILED: %v", err)
		return fmt.Errorf("TEST FAILED: %v", err)
	}
	log.Info().Msgf("TEST SUCCESS: %v", s)
	return nil
}
