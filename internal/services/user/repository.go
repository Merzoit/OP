package services

import (
	"at/constants"
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type UserRepository interface {
	CreateUser(user *User) error
	GetUser(tid uint64) (*User, error)
	BanUser(tid uint64) error
}

type PgUserRepository struct {
	db *pgxpool.Pool
}

func NewPgUserRepository(db *pgxpool.Pool) UserRepository {
	return &PgUserRepository{db: db}
}

func (repo *PgUserRepository) GetUser(tid uint64) (*User, error) {
	user := &User{}
	query := `
	SELECT id, telegram_id, username, role_id, is_banned, created_at
	FROM users 
	WHERE telegram_id=$1
	`

	log.Info().Msgf(constants.CallRepoGetUserById, tid)
	err := repo.db.QueryRow(context.Background(), query, tid).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.RoleID,
		&user.IsBanned,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Err(err).Msg(constants.ErrUserNotFound)
			return nil, fmt.Errorf(constants.ErrUserNotFound)
		}
		log.Error().Err(err).Msg(constants.ErrUserFetching)
		return nil, fmt.Errorf(constants.ErrUserFetching)
	}

	return user, nil
}

func (repo *PgUserRepository) CreateUser(user *User) error {
	query := `
	INSERT INTO users (telegram_id, username)
	VALUES ($1, $2)
	ON CONFLICT (telegram_id) DO UPDATE
	SET username = EXCLUDED.username
	RETURNING telegram_id
	`

	log.Info().Msgf(constants.CallRepoCreateUser, user.TelegramID)

	var returnedID uint64
	err := repo.db.QueryRow(
		context.Background(), query,
		user.TelegramID, user.Username,
	).Scan(&returnedID)

	if err != nil {
		if err.Error() == constants.ErrNoRows {
			log.Warn().Msg(constants.ErrUserAlreadyExist)
			return fmt.Errorf(constants.ErrUserAlreadyExist)
		}

		log.Error().Err(err).Msgf(constants.ErrUserCreate)
		return fmt.Errorf(constants.ErrUserCreate)
	}

	log.Info().Msgf(constants.SuccessfullyUserCreate, returnedID)
	return nil
}

func (repo *PgUserRepository) BanUser(tid uint64) error {
	query := `
	UPDATE users
	SET is_banned = TRUE
	WHERE telegram_id = $1
	`

	user, err := repo.db.Exec(context.Background(), query, tid)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrUserBanned)
		return fmt.Errorf(constants.ErrUserBanned)
	}

	if user.RowsAffected() == 0 {
		log.Error().Err(err).Msg(constants.ErrUserNotFound)
		return fmt.Errorf(constants.ErrUserNotFound)
	}

	return nil
}
