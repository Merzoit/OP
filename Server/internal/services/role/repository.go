package role

import (
	"at/tools/errors"
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type RoleRepository interface {
	GetRole(roleID uint) (*Role, error)
}

type PgRoleRepository struct {
	db *pgxpool.Pool
}

func NewPgRoleRepository(db *pgxpool.Pool) RoleRepository {
	return &PgRoleRepository{db: db}
}

func (repo *PgRoleRepository) GetRole(roleID uint) (*Role, error) {
	role := &Role{}

	query := `
	SELECT id, name, permissions
	FROM roles
	WHERE id = $1
	`

	err := repo.db.QueryRow(context.Background(), query, roleID).Scan(
		&role.ID,
		&role.Name,
		&role.Permissions,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Err(err).Msg(errors.ErrRoleNotFound)
			return nil, fmt.Errorf(errors.ErrRoleNotFound)
		}
		log.Warn().Err(err).Msg(errors.ErrRoleFetching)
		return nil, fmt.Errorf(errors.ErrRoleFetching)
	}

	return role, nil
}
