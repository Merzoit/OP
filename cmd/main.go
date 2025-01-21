package main

import (
	"at/internal/services/code"
	"at/internal/services/reflink"
	"at/internal/services/role"
	"at/internal/services/sponsor"
	"at/internal/services/subscribe"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"

	"at/constants"
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

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "../configs/config.yaml"
	}

	config, err := tools.LoadConfig(configPath)
	if err != nil {
		log.Fatal().Err(err).Msg(constants.ErrConfigLoad)
	}

	if err := config.Validate(); err != nil {
		log.Fatal().Err(err).Msg(constants.ErrConfigValidate)
	}

	log.Info().Msg(fmt.Sprintf("Application starting on port %d", config.App.Port))

	db, err := internal.Connect(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the database")
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
	//ts.TestDbSponsorCreate(db)
	//ts.TestDbSponsorsGet(db)
	//tc.TestDbCodeDelete(db)
	//+tc.TestDbCodeCreate(db)
	//+tc.TestDbCodeGet(db)
	//test.TestDbUserCreate(db)
	//test.TestDbUserGet(db)
	//tw.TestDbWorkerCreate(db)
	//tw.TestDbWorkerGet(db)
	//tr.TestDbRoleGet(db)
	//router := routers.InitRouter(userController)
	router := mux.NewRouter()
	routers.RegisterUserRoutes(router, userController)
	routers.RegisterCodeRoutes(router, codeController)
	routers.RegisterReflinkRoutes(router, reflinkController)
	routers.RegisterRoleRoutes(router, roleController)
	routers.RegisterSponsorRoutes(router, sponsorController)
	routers.RegisterSubscribeRoutes(router, subController)

	log.Fatal().Err(http.ListenAndServe(fmt.Sprintf(":%d", config.App.Port), router)).Msg("HTTP server stopped")
}
