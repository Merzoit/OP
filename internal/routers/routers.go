package routers

import (
	"at/internal/services/code"
	"at/internal/services/reflink"
	"at/internal/services/role"
	"at/internal/services/sponsor"
	user "at/internal/services/user"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterUserRoutes(router *mux.Router, uc *user.UserController) {
	ur := router.PathPrefix("/api/user").Subrouter()
	ur.HandleFunc("", uc.CreateUser).Methods(http.MethodPost)
	ur.HandleFunc("/{id}", uc.GetUser).Methods(http.MethodGet)
}

func RegisterCodeRoutes(router *mux.Router, controller *code.CodeController) {
	cr := router.PathPrefix("/api/code").Subrouter()
	cr.HandleFunc("", controller.CreateCode).Methods("POST")
	cr.HandleFunc("/{accessCode}", controller.GetCode).Methods("GET")
	cr.HandleFunc("/{accessCode}", controller.DeleteCode).Methods("DELETE")
	cr.HandleFunc("/increment/{accessCode}", controller.AddRequestCount).Methods("PATCH")
	cr.HandleFunc("/worker/{workerID}", controller.GetCodesByWorker).Methods("GET")
}

func RegisterReflinkRoutes(router *mux.Router, controller *reflink.RefLinkController) {
	rlr := router.PathPrefix("/api/reflink").Subrouter()
	rlr.HandleFunc("", controller.CreateLink).Methods("POST")
	rlr.HandleFunc("/{workerID}", controller.GetLink).Methods("GET")
	rlr.HandleFunc("/{workerID}", controller.UpdateLink).Methods("PATCH")
	rlr.HandleFunc("/clicks/{workerID}", controller.ClickAdd).Methods("PATCH")
	rlr.HandleFunc("/registrations/{workerID}", controller.RegistrationAdd).Methods("PATCH")
}

func RegisterRoleRoutes(router *mux.Router, controller *role.RoleController) {
	rrr := router.PathPrefix("/api/role").Subrouter()
	rrr.HandleFunc("/{roleID}", controller.GetRole).Methods("GET")
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
