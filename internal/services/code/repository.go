package code

import (
	"at/tools/errors"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type CodeRepository interface {
	CreateCode(code *Code) (*Code, error)
	GetCode(code int) (*Code, error)
	GetCodesByWorker(workerID uint) ([]*Code, error)
	DeleteCode(code int) error
	AddRequestCount(accessCode, i int) error
}

type PgCodeRepository struct {
	db *pgxpool.Pool
}

func NewPgCodeRepository(db *pgxpool.Pool) CodeRepository {
	return &PgCodeRepository{db: db}
}

func (repo *PgCodeRepository) CreateCode(code *Code) (*Code, error) {
	tx, err := repo.db.Begin(context.Background())
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrBeginTransaction)
		return nil, err
	}
	tx.Rollback(context.Background())

	checkQuery := `
	SELECT id
	FROM codes
	WHERE access_code = $1
	`

	var id int
	err = tx.QueryRow(context.Background(), checkQuery, code.AccessCode).Scan(&id)

	if err == nil {
		log.Warn().Err(err).Msg(errors.ErrCodeAlreadyExist)
		return nil, fmt.Errorf(errors.ErrCodeAlreadyExist)
	} else if err.Error() != errors.ErrNoRows {
		log.Warn().Err(err).Msg(errors.ErrCodeCheck)
		return nil, fmt.Errorf(errors.ErrCodeCheck)
	}

	insertQuery := `
	INSERT INTO codes (access_code, title, year, description, added_by_worker_id, request_count, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
	`

	err = tx.QueryRow(
		context.Background(), insertQuery, code.AccessCode, code.Title,
		code.Year, code.Description, code.AddedByWorkerID, code.RequestCount, time.Now(),
	).Scan(&id)

	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrCodeCreate)
		return nil, fmt.Errorf(errors.ErrCodeCreate)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrCommitTransaction)
		return nil, err
	}

	code.Id = id
	return code, nil
}

func (repo *PgCodeRepository) GetCode(code int) (*Code, error) {
	c := &Code{}
	query := `
	SELECT id, access_code, title, year, description, added_by_worker_id, request_count, created_at
	FROM codes
	WHERE access_code=$1
	`

	err := repo.db.QueryRow(context.Background(), query, code).Scan(
		&c.Id,
		&c.AccessCode,
		&c.Title,
		&c.Year,
		&c.Description,
		&c.AddedByWorkerID,
		&c.RequestCount,
		&c.CreatedAt,
	)

	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrCodeFetching)
		return nil, fmt.Errorf(errors.ErrCodeFetching)
	}

	return c, nil
}

func (repo *PgCodeRepository) GetCodesByWorker(workerID uint) ([]*Code, error) {
	query := `
	SELECT id, access_code, title, year, description, added_by_worker_id, request_count, created_at
	FROM codes
	WHERE added_by_worker_id = $1
	`

	rows, err := repo.db.Query(context.Background(), query, workerID)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrCodeFetchingByWorker)
		return nil, fmt.Errorf(errors.ErrCodeFetchingByWorker)
	}
	defer rows.Close()

	var codes []*Code
	for rows.Next() {
		code := &Code{}
		err := rows.Scan(
			&code.Id,
			&code.AccessCode,
			&code.Title,
			&code.Year,
			&code.Description,
			&code.AddedByWorkerID,
			&code.RequestCount,
			&code.CreatedAt,
		)
		if err != nil {
			log.Warn().Err(err).Msg(errors.ErrCodeScan)
			return nil, fmt.Errorf(errors.ErrCodeScan)
		}
		codes = append(codes, code)
	}

	if rows.Err() != nil {
		log.Warn().Err(rows.Err()).Msg(errors.ErrCodesIterate)
		return nil, fmt.Errorf(errors.ErrCodesIterate)
	}

	return codes, nil
}

func (repo *PgCodeRepository) DeleteCode(code int) error {
	query := `
	DELETE 
	FROM codes 
	WHERE access_code=$1
	`

	_, err := repo.db.Exec(context.Background(), query, code)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrCodeDelete)
		return fmt.Errorf("%v: %v", errors.ErrCodeDelete, err)
	}

	return err
}

func (repo *PgCodeRepository) AddRequestCount(accessCode, i int) error {
	query := `
	UPDATE codes
	SET request_count = request_count + $2
	WHERE access_code = $1
	`

	code, err := repo.db.Exec(context.Background(), query, accessCode, i)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrCodeRequestCounter)
		return fmt.Errorf(errors.ErrCodeRequestCounter)
	}

	if code.RowsAffected() == 0 {
		log.Warn().Err(err).Msg(errors.ErrCodeNotFound)
		return fmt.Errorf(errors.ErrCodeNotFound)
	}

	return nil
}
