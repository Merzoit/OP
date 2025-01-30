package worker

import (
	"at/internal/services/worker"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func TestDbWorkerCreate(db *pgxpool.Pool) error {
	workerRepo := worker.NewPgWorkerRepository(db)

	w := &worker.Worker{
		PaymentRate: 8,
		UserID:      3,
		Description: "[T] Test",
	}

	err := workerRepo.CreateWorker(w)
	if err != nil {
		log.Warn().Msgf("TEST FAILED: %v", err)
		return fmt.Errorf("TEST FAILED: %v", err)
	}

	log.Info().Msg("TEST PASSED")
	return nil
}

func TestDbWorkerGet(db *pgxpool.Pool) error {
	workerRepo := worker.NewPgWorkerRepository(db)

	id := 1

	w, err := workerRepo.GetWorker(id)
	if err != nil {
		log.Warn().Msgf("TEST FAILED: %v", err)
		return fmt.Errorf("TEST FAILED: %v", err)
	}

	log.Info().Msgf("TEST SUCCESS: %v", w)
	return nil
}
