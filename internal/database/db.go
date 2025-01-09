package internal

import (
	"at/constants"
	"at/tools"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var DB *pgxpool.Pool

func Connect(c *tools.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrDbParseConfig)
		return nil, fmt.Errorf(constants.ErrDbParseConfig)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Error().Err(err).Msg(constants.ErrDbConnect)
		return nil, fmt.Errorf(constants.ErrDbConnect)
	}

	if err := db.Ping(context.Background()); err != nil {
		log.Error().Err(err).Msg(constants.ErrDbPing)
		return nil, fmt.Errorf(constants.ErrDbPing)
	}

	log.Info().Msg("Database connect successfully")
	DB = db
	return db, nil
}
