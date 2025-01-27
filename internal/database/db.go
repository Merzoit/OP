package internal

import (
	"at/tools/errors"
	"context"
	"fmt"

	"at/tools"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var DB *pgxpool.Pool

func Connect(c *tools.Config) (*pgxpool.Pool, error) {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		c.DatabaseUser,
		c.DatabasePass,
		c.DatabaseHost,
		c.DatabasePort,
		c.DatabaseName,
	)

	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrDbParseConfig)
		return nil, fmt.Errorf(errors.ErrDbParseConfig)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Warn().Err(err).Msg(errors.ErrDbConnect)
		return nil, fmt.Errorf(errors.ErrDbConnect)
	}

	if err := db.Ping(context.Background()); err != nil {
		log.Warn().Err(err).Msg(errors.ErrDbPing)
		return nil, fmt.Errorf(errors.ErrDbPing)
	}

	log.Info().Msg("Database connected successfully")
	DB = db
	return db, nil
}
