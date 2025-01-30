package worker

import (
	"at/tools/errors"
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
		if err.Error() == errors.ErrNoRows {
			log.Warn().Msg(errors.ErrUserAlreadyExist)
			return fmt.Errorf(errors.ErrUserAlreadyExist)
		}

		log.Warn().Err(err).Msgf(errors.ErrUserFetching)
		return fmt.Errorf(errors.ErrUserFetching)
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
		log.Warn().Err(err).Msg(errors.ErrWorkerCreate)
		return fmt.Errorf(errors.ErrWorkerCreate)
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
			log.Warn().Err(err).Msg(errors.ErrWorkerNotFound)
			return nil, fmt.Errorf(errors.ErrWorkerNotFound)
		}
		log.Warn().Err(err).Msg(errors.ErrWorkerFetching)
		return nil, fmt.Errorf(errors.ErrWorkerFetching)
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
		log.Warn().Err(err).Msg(errors.ErrWorkerDelete)
		return fmt.Errorf(errors.ErrWorkerDelete)
	}

	if worker.RowsAffected() == 0 {
		log.Warn().Err(err).Msg(errors.ErrWorkerNotFound)
		return fmt.Errorf(errors.ErrWorkerNotFound)
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
		log.Warn().Err(err).Msg(errors.ErrWorkerUpdatePaymentRate)
		return fmt.Errorf(errors.ErrWorkerUpdatePaymentRate)
	}

	if cmdTag.RowsAffected() == 0 {
		log.Warn().Err(err).Msg(errors.ErrWorkerNotFound)
		return fmt.Errorf(errors.ErrWorkerNotFound)
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
		log.Warn().Err(err).Msg(errors.ErrWorkerBalanceUp)
		return fmt.Errorf(errors.ErrWorkerBalanceUp)
	}

	if worker.RowsAffected() == 0 {
		log.Warn().Err(err).Msg(errors.ErrWorkerNotFound)
		return fmt.Errorf(errors.ErrWorkerNotFound)
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
		log.Warn().Err(err).Msg(errors.ErrWorkerBalanceReset)
		return fmt.Errorf(errors.ErrWorkerBalanceReset)
	}

	if cmdTag.RowsAffected() == 0 {
		log.Warn().Err(err).Msg(errors.ErrWorkerNotFound)
		return fmt.Errorf(errors.ErrWorkerNotFound)
	}

	return nil
}
