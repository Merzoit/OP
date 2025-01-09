package role

import (
	"at/internal/services/role"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func TestDbRoleGet(db *pgxpool.Pool) error {
	roleRepo := role.NewPgRoleRepository(db)

	roleID := uint(1)

	role, err := roleRepo.GetRole(roleID)
	if err != nil {
		log.Warn().Msgf("TEST FAILED: %v", err)
		return fmt.Errorf("TEST FAILED: %v", err)
	}

	log.Info().Msgf("TEST SUCCESS: %v", role)
	return nil
}
