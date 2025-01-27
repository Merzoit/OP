package subscribe

import (
	"at/tools/errors"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type SubscribeRepository interface {
	AddSubscribe(sub *Subscribe) (*Subscribe, error)
	GetSubscribesByUser(userID uint) ([]*Subscribe, error)
	GetSubscribesBySponsor(sponsorID uint) ([]*Subscribe, error)
}

type PgSubscribeRepository struct {
	db *pgxpool.Pool
}

func NewPgSubscribeRepository(db *pgxpool.Pool) SubscribeRepository {
	return &PgSubscribeRepository{db: db}
}

func (repo *PgSubscribeRepository) AddSubscribe(sub *Subscribe) (*Subscribe, error) {
	tx, err := repo.db.Begin(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to begin transaction")
		return nil, err
	}
	defer tx.Rollback(context.Background())

	checkQuery := `
	SELECT id FROM subscriptions
	WHERE user_id = $1 AND sponsor_id = $2
	`
	var existingID int
	err = tx.QueryRow(context.Background(), checkQuery, sub.UserID, sub.SponsorID).Scan(&existingID)

	if err == nil {
		log.Warn().Msgf("Subscription already exists for user_id=%d and sponsor_id=%d", sub.UserID, sub.SponsorID)
		return nil, fmt.Errorf("subscription already exists")
	} else if err != pgx.ErrNoRows {
		log.Error().Err(err).Msg("Error checking existing subscription")
		return nil, err
	}

	insertQuery := `
	INSERT INTO subscriptions (user_id, sponsor_id)
	VALUES ($1, $2)
	RETURNING id, user_id, sponsor_id, subscribed_at
	`
	err = tx.QueryRow(context.Background(), insertQuery, sub.UserID, sub.SponsorID).Scan(
		&sub.ID,
		&sub.UserID,
		&sub.SponsorID,
		&sub.SubscribedAt,
	)
	if err != nil {
		log.Error().Err(err).Msg("Error inserting subscription")
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return nil, err
	}

	return sub, nil
}

func (repo *PgSubscribeRepository) GetSubscribesByUser(userID uint) ([]*Subscribe, error) {
	query := `
	SELECT id, user_id, sponsor_id, subscribed_at
	FROM subscriptions
	WHERE user_id = $1
	`

	rows, err := repo.db.Query(context.Background(), query, userID)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrSubscribesByUser)
		return nil, fmt.Errorf(errors.ErrSubscribesByUser)
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
			log.Warn().Err(err).Msg(errors.ErrSubscribeScan)
			return nil, fmt.Errorf(errors.ErrSubscribeScan)
		}
		subscribes = append(subscribes, sub)
	}

	if rows.Err() != nil {
		log.Warn().Err(rows.Err()).Msg(errors.ErrSubscribeIterate)
		return nil, fmt.Errorf(errors.ErrSubscribeIterate)
	}

	return subscribes, nil
}

func (repo *PgSubscribeRepository) GetSubscribesBySponsor(sponsorID uint) ([]*Subscribe, error) {
	query := `
	SELECT id, user_id, sponsor_id, subscribed_at
	FROM subscriptions
	WHERE sponsor_id = $1
	`

	rows, err := repo.db.Query(context.Background(), query, sponsorID)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrSubscribesBySponsor)
		return nil, fmt.Errorf(errors.ErrSubscribesBySponsor)
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
			log.Warn().Err(err).Msg(errors.ErrSubscribeScan)
			return nil, fmt.Errorf(errors.ErrSubscribeScan)
		}
		subscribes = append(subscribes, sub)
	}

	if rows.Err() != nil {
		log.Warn().Err(rows.Err()).Msg(errors.ErrSubscribeIterate)
		return nil, fmt.Errorf(errors.ErrSubscribeIterate)
	}

	return subscribes, nil
}
