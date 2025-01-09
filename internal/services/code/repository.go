package code

import (
	"at/constants"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type CodeRepository interface {
	CreateCode(code *Code) error
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

func (repo *PgCodeRepository) CreateCode(code *Code) error {

	checkQuery := `
	SELECT id
	FROM codes
	WHERE access_code = $1
	`

	var id int
	err := repo.db.QueryRow(context.Background(), checkQuery, code.AccessCode).Scan(&id)

	if err == nil {
		log.Warn().Msg(constants.ErrCodeAlreadyExist)
		return nil
	} else if err.Error() != constants.ErrNoRows {
		log.Error().Err(err).Msg(constants.ErrCodeCheck)
		return fmt.Errorf(constants.ErrCodeCheck)
	}

	insertQuery := `
	INSERT INTO codes (access_code, title, year, description, added_by_worker_id, request_count, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
	`

	err = repo.db.QueryRow(
		context.Background(), insertQuery, code.AccessCode, code.Title,
		code.Year, code.Description, code.AddedByWorkerID, code.RequestCount, time.Now(),
	).Scan(&id)

	if err != nil {
		log.Error().Err(err).Msg(constants.ErrCodeCreate)
		return fmt.Errorf(constants.ErrCodeCreate)
	}

	return nil
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
		log.Error().Err(err).Msg(constants.ErrCodeFetching)
		return nil, fmt.Errorf(constants.ErrCodeFetching)
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
		log.Error().Err(err).Msg(constants.ErrCodeFetchingByWorker)
		return nil, fmt.Errorf(constants.ErrCodeFetchingByWorker)
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
			log.Error().Err(err).Msg(constants.ErrCodeScan)
			return nil, fmt.Errorf(constants.ErrCodeScan)
		}
		codes = append(codes, code)
	}

	if rows.Err() != nil {
		log.Error().Err(rows.Err()).Msg(constants.ErrCodesIterate)
		return nil, fmt.Errorf(constants.ErrCodesIterate)
	}

	return codes, nil
}

func (repo *PgCodeRepository) DeleteCode(code int) error {
	query := `
	DELETE 
	FROM codes 
	WHERE code=$1
	`

	log.Info().Msgf(constants.CallRepoDeleteCode, code)

	_, err := repo.db.Exec(context.Background(), query, code)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrCodeDelete)
		return fmt.Errorf("%v: %v", constants.ErrCodeDelete, err)
	}

	return err
}

func (repo *PgCodeRepository) AddRequestCount(accessCode, i int) error {
	query := `
	UPDATE codes
	SET request_count = request_count + $2
	WHERE access_code = $1
	`

	cmdTag, err := repo.db.Exec(context.Background(), query, accessCode, i)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrCodeRequestCounter)
		return fmt.Errorf(constants.ErrCodeRequestCounter)
	}

	if cmdTag.RowsAffected() == 0 {
		log.Error().Err(err).Msg(constants.ErrCodeNotFound)
		return fmt.Errorf(constants.ErrCodeNotFound)
	}

	return nil
}
