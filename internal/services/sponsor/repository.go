package sponsor

import (
	"at/constants"
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type SponsorRepository interface {
	CreateSponsor(sponsor *Sponsor) error
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

func (repo *PgSponsorRepository) CreateSponsor(sponsor *Sponsor) error {
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
		log.Error().Err(err).Msg(constants.ErrSponsorCreate)
		return fmt.Errorf(constants.ErrSponsorCreate)
	}

	return nil
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
			log.Error().Err(err).Msg(constants.ErrSponsorNotFound)
			return nil, fmt.Errorf(constants.ErrSponsorNotFound)
		}

		log.Error().Err(err).Msg(constants.ErrSponsorFetching)
		return nil, fmt.Errorf(constants.ErrSponsorFetching)
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
		log.Error().Err(err).Msg(constants.ErrSponsorDelete)
		return fmt.Errorf(constants.ErrSponsorDelete)
	}

	if sponsor.RowsAffected() == 0 {
		log.Error().Err(err).Msg(constants.ErrSponsorNotFound)
		return fmt.Errorf(constants.ErrSponsorNotFound)
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
		log.Error().Err(err).Msg(constants.ErrSponsorsGet)
		return nil, fmt.Errorf(constants.ErrSponsorsGet)
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
			log.Error().Err(err).Msg(constants.ErrSponsorScan)
			return nil, fmt.Errorf(constants.ErrSponsorScan)
		}
		sponsors = append(sponsors, sponsor)
	}

	if rows.Err() != nil {
		log.Error().Err(err).Msg(constants.ErrSponsorsIterate)
		return nil, fmt.Errorf(constants.ErrSponsorsIterate)
	}

	return sponsors, nil
}
