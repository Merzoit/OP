package worker

import (
	"at/constants"
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type WorkerRepository interface {
	CreateWorker(worker *Worker) error
	GetWorker(workerID int) (*Worker, error)
	DeleteWorker(workerID int) error
	UpdatePaymentRate(workerID int, rate float64) error
	BalanceUp(workerID int, count float64) error
	BalanceReset(workerID int) error
}

type PgWorkerRepository struct {
	db *pgxpool.Pool
}

func NewPgWorkerRepository(db *pgxpool.Pool) WorkerRepository {
	return &PgWorkerRepository{db: db}
}

func (repo *PgWorkerRepository) CreateWorker(worker *Worker) error {

	checkUserQuery := `
	SELECT id
	FROM users
	WHERE id = $1
	`

	var existingUserId int64
	err := repo.db.QueryRow(context.Background(), checkUserQuery, worker.UserID).Scan(&existingUserId)
	if err != nil {
		if err.Error() == constants.ErrNoRows {
			log.Warn().Msg(constants.ErrUserAlreadyExist)
			return fmt.Errorf(constants.ErrUserAlreadyExist)
		}

		log.Error().Err(err).Msgf(constants.ErrUserFetching)
		return fmt.Errorf(constants.ErrUserFetching)
	}

	insertWorkerQuery := `
	INSERT INTO workers (user_id, payment_rate, description)
	VALUES ($1, $2, $3)
	RETURNING id
	ON CONFLICT (user_id) DO NOTHING
	`

	var workerId int64
	err = repo.db.QueryRow(
		context.Background(),
		insertWorkerQuery,
		worker.UserID,
		worker.PaymentRate,
		worker.Description,
	).Scan(&workerId)

	if err != nil {
		log.Error().Err(err).Msg(constants.ErrWorkerCreate)
		return fmt.Errorf(constants.ErrWorkerCreate)
	}

	return nil
}

func (repo *PgWorkerRepository) GetWorker(workerID int) (*Worker, error) {
	worker := &Worker{}

	query := `
	SELECT 
		id,
		payment_rate,
		description,
		balance,
		created_at,
		user_id
	FROM workers
	WHERE id = $1
	`

	err := repo.db.QueryRow(context.Background(), query, workerID).Scan(
		&worker.ID,
		&worker.PaymentRate,
		&worker.Description,
		&worker.Balance,
		&worker.CreatedAt,
		&worker.UserID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Err(err).Msg(constants.ErrWorkerNotFound)
			return nil, fmt.Errorf(constants.ErrWorkerNotFound)
		}
		log.Error().Err(err).Msg(constants.ErrWorkerFetching)
		return nil, fmt.Errorf(constants.ErrWorkerFetching)
	}

	return worker, nil
}

func (repo *PgWorkerRepository) DeleteWorker(workerID int) error {
	query := `
	DELETE FROM workers
	WHERE id = $1
	`

	worker, err := repo.db.Exec(context.Background(), query, workerID)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrWorkerDelete)
		return fmt.Errorf(constants.ErrWorkerDelete)
	}

	if worker.RowsAffected() == 0 {
		log.Error().Err(err).Msg(constants.ErrWorkerNotFound)
		return fmt.Errorf(constants.ErrWorkerNotFound)
	}

	return nil
}

func (repo *PgWorkerRepository) UpdatePaymentRate(workerID int, rate float64) error {
	query := `
	UPDATE workers
	SET payment_rate = $1
	WHERE id = $2
	`

	cmdTag, err := repo.db.Exec(context.Background(), query, rate, workerID)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrWorkerUpdatePaymentRate)
		return fmt.Errorf(constants.ErrWorkerUpdatePaymentRate)
	}

	if cmdTag.RowsAffected() == 0 {
		log.Error().Err(err).Msg(constants.ErrWorkerNotFound)
		return fmt.Errorf(constants.ErrWorkerNotFound)
	}

	return nil
}

func (repo *PgWorkerRepository) BalanceUp(workerID int, count float64) error {
	query := `
	UPDATE workers
	SET balance = balance + $1
	WHERE id = $2
	`

	worker, err := repo.db.Exec(context.Background(), query, count, workerID)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrWorkerBalanceUp)
		return fmt.Errorf(constants.ErrWorkerBalanceUp)
	}

	if worker.RowsAffected() == 0 {
		log.Error().Err(err).Msg(constants.ErrWorkerNotFound)
		return fmt.Errorf(constants.ErrWorkerNotFound)
	}

	return nil
}

func (repo *PgWorkerRepository) BalanceReset(workerID int) error {
	query := `
	UPDATE workers
	SET balance = 0
	WHERE id = $1
	`

	cmdTag, err := repo.db.Exec(context.Background(), query, workerID)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrWorkerBalanceReset)
		return fmt.Errorf(constants.ErrWorkerBalanceReset)
	}

	if cmdTag.RowsAffected() == 0 {
		log.Error().Err(err).Msg(constants.ErrWorkerNotFound)
		return fmt.Errorf(constants.ErrWorkerNotFound)
	}

	return nil
}
