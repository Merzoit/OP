package routers

import (
	code "at/internal/services/code"
	"at/internal/services/sponsor"
	user "at/internal/services/user"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterUserRoutes(router *mux.Router, uc *user.UserController) {
	ur := router.PathPrefix("/api/user").Subrouter()
	ur.HandleFunc("", uc.CreateUser).Methods(http.MethodPost)
	ur.HandleFunc("/{id}", uc.GetUserById).Methods(http.MethodGet)
}

func RegisterCodeRoutes(router *mux.Router, cc *code.CodeController) {
	cr := router.PathPrefix("/api/code").Subrouter()
	cr.HandleFunc("", cc.CreateCode).Methods(http.MethodPost)
	cr.HandleFunc("/{id}", cc.GetCode).Methods(http.MethodGet)
}

func RegisterSponsorRoutes(router *mux.Router, sc *sponsor.SponsorController) {
	sr := router.PathPrefix("/api/sponsor").Subrouter()
	sr.HandleFunc("", sc.CreateSponsor).Methods(http.MethodPost)
	sr.HandleFunc("/all", sc.GetSponsors).Methods(http.MethodGet)
	sr.HandleFunc("/{id}", sc.GetSponsor).Methods(http.MethodGet)
	sr.HandleFunc("/{id}", sc.DeleteSponsor).Methods(http.MethodDelete)
}

func InitRouter(userController *user.UserController) *mux.Router {
	router := mux.NewRouter()
	RegisterUserRoutes(router, userController)
	return router
}
