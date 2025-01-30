package main

import (
	"at/internal/services/code"
	"at/internal/services/reflink"
	"at/internal/services/role"
	"at/internal/services/sponsor"
	"at/internal/services/subscribe"
	"at/tools/errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"net/http"
	"os"

	internal "at/internal/database"
	"at/internal/routers"
	user "at/internal/services/user"
	"at/tools"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	log.Info().Msg("Starting application...")

	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msgf(errors.ErrEnvLoading)
	}

	connectionString, err := tools.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg(errors.ErrConnectionString)
	}

	if err := connectionString.Validate(); err != nil {
		log.Fatal().Err(err).Msg(errors.ErrConnectionStringValidate)
	}

	log.Info().Msg(fmt.Sprintf("Application starting on port %d", connectionString.AppPort))

	db, err := internal.Connect(connectionString)
	if err != nil {
		log.Fatal().Err(err).Msg(errors.ErrDbConnect)
	}
	defer db.Close()

	userRepo := user.NewPgUserRepository(db)
	userController := user.NewUserController(userRepo)

	codeRepo := code.NewPgCodeRepository(db)
	codeController := code.NewCodeController(codeRepo)

	reflinkRepo := reflink.NewPgReflinkRepository(db)
	reflinkController := reflink.NewRefLinkController(reflinkRepo)

	roleRepo := role.NewPgRoleRepository(db)
	roleController := role.NewRoleController(roleRepo)

	sponsorRepo := sponsor.NewPgSponsorRepository(db)
	sponsorController := sponsor.NewSponsorController(sponsorRepo)

	subRepo := subscribe.NewPgSubscribeRepository(db)
	subController := subscribe.NewSubscribeController(subRepo)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	routers.RegisterUserRoutes(router, userController)
	routers.RegisterCodeRoutes(router, codeController)
	routers.RegisterReflinkRoutes(router, reflinkController)
	routers.RegisterRoleRoutes(router, roleController)
	routers.RegisterSponsorRoutes(router, sponsorController)
	routers.RegisterSubscribeRoutes(router, subController)

	log.Fatal().Err(http.ListenAndServe(fmt.Sprintf(":%d", connectionString.AppPort), router)).Msg("HTTP server stopped")
}
