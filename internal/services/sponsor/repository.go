package sponsor

import (
	"at/tools/errors"
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type SponsorRepository interface {
	CreateSponsor(sponsor *Sponsor) (*Sponsor, error)
	GetSponsor(tid uint64) (*Sponsor, error)
	DeleteSponsor(tid uint64) error
	GetSponsors() ([]*Sponsor, error)
}

type PgSponsorRepository struct {
	db *pgxpool.Pool
}

func NewPgSponsorRepository(db *pgxpool.Pool) SponsorRepository {
	return &PgSponsorRepository{db: db}
}

func (repo *PgSponsorRepository) CreateSponsor(sponsor *Sponsor) (*Sponsor, error) {
	query := `
	INSERT INTO sponsors 
	(telegram_link, price_per_sub, name)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	err := repo.db.QueryRow(
		context.Background(),
		query,
		sponsor.TelegramLink,
		sponsor.PricePerSub,
		sponsor.Name,
	).Scan(&sponsor.ID)

	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrSponsorCreate)
		return nil, fmt.Errorf(errors.ErrSponsorCreate)
	}

	return sponsor, nil
}

func (repo *PgSponsorRepository) GetSponsor(tid uint64) (*Sponsor, error) {
	sponsor := &Sponsor{}

	query := `
	SELECT id, telegram_link, price_per_sub, name, created_at
	FROM sponsors
	WHERE id = $1
	`

	err := repo.db.QueryRow(context.Background(), query, tid).Scan(
		&sponsor.ID,
		&sponsor.TelegramLink,
		&sponsor.PricePerSub,
		&sponsor.Name,
		&sponsor.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Err(err).Msg(errors.ErrSponsorNotFound)
			return nil, fmt.Errorf(errors.ErrSponsorNotFound)
		}

		log.Warn().Err(err).Msg(errors.ErrSponsorFetching)
		return nil, fmt.Errorf(errors.ErrSponsorFetching)
	}

	return sponsor, nil
}

func (repo *PgSponsorRepository) DeleteSponsor(tid uint64) error {
	query := `
	DELETE FROM sponsors
	WHERE id = $1
	`

	sponsor, err := repo.db.Exec(context.Background(), query, tid)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrSponsorDelete)
		return fmt.Errorf(errors.ErrSponsorDelete)
	}

	if sponsor.RowsAffected() == 0 {
		log.Warn().Err(err).Msg(errors.ErrSponsorNotFound)
		return fmt.Errorf(errors.ErrSponsorNotFound)
	}

	return nil
}

func (repo *PgSponsorRepository) GetSponsors() ([]*Sponsor, error) {
	query := `
	SELECT id, telegram_link, price_per_sub, name, created_at
	FROM sponsors
	`

	rows, err := repo.db.Query(context.Background(), query)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrSponsorsGet)
		return nil, fmt.Errorf(errors.ErrSponsorsGet)
	}
	defer rows.Close()

	var sponsors []*Sponsor
	for rows.Next() {
		sponsor := &Sponsor{}
		err := rows.Scan(
			&sponsor.ID,
			&sponsor.TelegramLink,
			&sponsor.PricePerSub,
			&sponsor.Name,
			&sponsor.CreatedAt,
		)
		if err != nil {
			log.Warn().Err(err).Msg(errors.ErrSponsorScan)
			return nil, fmt.Errorf(errors.ErrSponsorScan)
		}
		sponsors = append(sponsors, sponsor)
	}

	if rows.Err() != nil {
		log.Warn().Err(err).Msg(errors.ErrSponsorsIterate)
		return nil, fmt.Errorf(errors.ErrSponsorsIterate)
	}

	return sponsors, nil
}
