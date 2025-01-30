package reflink

import (
	"at/tools/errors"
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type ReflinkRepository interface {
	CreateLink(link *ReferralLink) (*ReferralLink, error)
	GetLink(workerID int) (*ReferralLink, error)
	UpdateLink(workerID int, tag string) error
	ClickAdd(workerID, count int) error
	RegistrationAdd(workerID, count int) error
}

type PgReflinkRepository struct {
	db *pgxpool.Pool
}

func NewPgReflinkRepository(db *pgxpool.Pool) ReflinkRepository {
	return &PgReflinkRepository{db: db}
}

func (repo *PgReflinkRepository) CreateLink(link *ReferralLink) (*ReferralLink, error) {
	globalLink := "https://t.me/pokeepsi_movie_bot?start="
	ref := fmt.Sprintf("%s%s", globalLink, link.Link)

	query := `
	INSERT INTO referral_links (worker_id, referral_link)
	VALUES ($1, $2)
	RETURNING id
	`

	err := repo.db.QueryRow(
		context.Background(),
		query,
		link.WorkerID,
		ref,
	).Scan(&link.ID)

	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrRefCreate)
		return nil, fmt.Errorf(errors.ErrRefCreate)
	}

	return link, nil
}

func (repo *PgReflinkRepository) GetLink(workerID int) (*ReferralLink, error) {
	link := &ReferralLink{}

	query := `
	SELECT id, worker_id, referral_link, clicks, registrations, created_at
	FROM referral_links
	WHERE worker_id = $1
	`

	err := repo.db.QueryRow(context.Background(), query, workerID).Scan(
		&link.ID,
		&link.WorkerID,
		&link.Link,
		&link.Clicks,
		&link.Registrations,
		&link.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Err(err).Msg(errors.ErrNoRows)
			return nil, fmt.Errorf(errors.ErrNoRows)
		}

		log.Warn().Err(err).Msg(errors.ErrRefGet)
		return nil, fmt.Errorf(errors.ErrRefGet)
	}

	return link, nil
}

func (repo *PgReflinkRepository) UpdateLink(workerID int, tag string) error {

	globalLink := "https://t.me/pokeepsi_movie_bot?start="
	ref := fmt.Sprintf("%s%s", globalLink, tag)

	query := `
	UPDATE referral_links
	SET referral_link = $1
	WHERE worker_id = $2
	`

	link, err := repo.db.Exec(context.Background(), query, ref, workerID)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrRefUpdate)
		return fmt.Errorf(errors.ErrRefUpdate)
	}

	if link.RowsAffected() == 0 {
		log.Warn().Err(err).Msg(errors.ErrRefNotFound)
		return fmt.Errorf(errors.ErrRefNotFound)
	}

	return nil
}

func (repo *PgReflinkRepository) ClickAdd(workerID int, count int) error {
	query := `
	UPDATE referral_links
	SET clicks = clicks + $1
	WHERE worker_id = $2
	`

	link, err := repo.db.Exec(context.Background(), query, count, workerID)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrRefClickAdd)
		return fmt.Errorf(errors.ErrRefClickAdd)
	}

	if link.RowsAffected() == 0 {
		log.Warn().Err(err).Msg(errors.ErrRefNotFound)
		return fmt.Errorf(errors.ErrRefNotFound)
	}

	return nil
}

func (repo *PgReflinkRepository) RegistrationAdd(workerID int, count int) error {
	query := `
	UPDATE referral_links
	SET registrations = registrations + $1
	WHERE worker_id = $2
	`

	link, err := repo.db.Exec(context.Background(), query, count, workerID)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrRefRegAdd)
		return fmt.Errorf(errors.ErrRefRegAdd)
	}

	if link.RowsAffected() == 0 {
		log.Warn().Err(err).Msg(errors.ErrRefNotFound)
		return fmt.Errorf(errors.ErrRefNotFound)
	}

	return nil
}
