package subscribe

import (
	"at/constants"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type SubscribeRepository interface {
	AddSubscribe(sub *Subscribe) error
	GetSubscribesByUser(userID uint) ([]*Subscribe, error)
	GetSubscribesBySponsor(sponsorID uint) ([]*Subscribe, error)
}

type PgSubscribeRepository struct {
	db *pgxpool.Pool
}

func NewPgSubscribeRepository(db *pgxpool.Pool) SubscribeRepository {
	return &PgSubscribeRepository{db: db}
}

func (repo *PgSubscribeRepository) AddSubscribe(sub *Subscribe) error {
	query := `
	INSERT INTO subscribes (user_id, sponsor_id)
	VALUES ($1, $2)
	RETURNING id
	`

	err := repo.db.QueryRow(
		context.Background(),
		query,
		sub.UserID,
		sub.SponsorID,
	).Scan(&sub.ID)

	if err != nil {
		log.Error().Err(err).Msg(constants.ErrSubscribeCreate)
		return fmt.Errorf(constants.ErrSubscribeCreate)
	}

	return nil
}

func (repo *PgSubscribeRepository) GetSubscribesByUser(userID uint) ([]*Subscribe, error) {
	query := `
	SELECT id, user_id, sponsor_id, subscribed_at
	FROM subscribes
	WHERE user_id = $1
	`

	rows, err := repo.db.Query(context.Background(), query, userID)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrSubscribesByUser)
		return nil, fmt.Errorf(constants.ErrSubscribesByUser)
	}
	defer rows.Close()

	var subscribes []*Subscribe
	for rows.Next() {
		sub := &Subscribe{}
		err := rows.Scan(
			&sub.ID,
			&sub.UserID,
			&sub.SponsorID,
			&sub.SubscribedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg(constants.ErrSubscribeScan)
			return nil, fmt.Errorf(constants.ErrSubscribeScan)
		}
		subscribes = append(subscribes, sub)
	}

	if rows.Err() != nil {
		log.Error().Err(rows.Err()).Msg(constants.ErrSubscribeIterate)
		return nil, fmt.Errorf(constants.ErrSubscribeIterate)
	}

	return subscribes, nil
}

func (repo *PgSubscribeRepository) GetSubscribesBySponsor(sponsorID uint) ([]*Subscribe, error) {
	query := `
	SELECT id, user_id, sponsor_id, subscribed_at
	FROM subscribes
	WHERE sponsor_id = $1
	`

	rows, err := repo.db.Query(context.Background(), query, sponsorID)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrSubscribesBySponsor)
		return nil, fmt.Errorf(constants.ErrSubscribesBySponsor)
	}
	defer rows.Close()

	var subscribes []*Subscribe
	for rows.Next() {
		sub := &Subscribe{}
		err := rows.Scan(
			&sub.ID,
			&sub.UserID,
			&sub.SponsorID,
			&sub.SubscribedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg(constants.ErrSubscribeScan)
			return nil, fmt.Errorf(constants.ErrSubscribeScan)
		}
		subscribes = append(subscribes, sub)
	}

	if rows.Err() != nil {
		log.Error().Err(rows.Err()).Msg(constants.ErrSubscribeIterate)
		return nil, fmt.Errorf(constants.ErrSubscribeIterate)
	}

	return subscribes, nil
}
