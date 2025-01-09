package main

import (
	"at/internal/services/code"
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
	/*
		sponsorRepo := sponsor.NewPgSponsorRepository(db)
		sponsorController := sponsor.NewSponsorController(sponsorRepo)
	*/
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
	//routers.RegisterSponsorRoutes(router, sponsorController)

	log.Fatal().Err(http.ListenAndServe(fmt.Sprintf(":%d", config.App.Port), router)).Msg("HTTP server stopped")
}
