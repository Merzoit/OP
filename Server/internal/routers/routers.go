package routers

import (
	"at/internal/services/code"
	"at/internal/services/reflink"
	"at/internal/services/role"
	"at/internal/services/sponsor"
	"at/internal/services/subscribe"
	user "at/internal/services/user"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterUserRoutes(router chi.Router, uc *user.UserController) {
	router.Route("/api/user", func(r chi.Router) {
		r.Method("POST", "/create", http.HandlerFunc(uc.CreateUser))
		r.Method("GET", "/{id}", http.HandlerFunc(uc.GetUser))
	})
}

func RegisterCodeRoutes(router chi.Router, cc *code.CodeController) {
	router.Route("/api/code", func(r chi.Router) {
		r.Method("POST", "/create", http.HandlerFunc(cc.CreateCode))
		r.Method("GET", "/{accessCode}", http.HandlerFunc(cc.GetCode))
		r.Method("DELETE", "/{accessCode}", http.HandlerFunc(cc.DeleteCode))
		r.Method("PATCH", "/increment/{accessCode}", http.HandlerFunc(cc.AddRequestCount))
		r.Method("GET", "/worker/{workerID}", http.HandlerFunc(cc.GetCodesByWorker))
	})
}

func RegisterReflinkRoutes(router chi.Router, rc *reflink.RefLinkController) {
	router.Route("/api/reflink", func(r chi.Router) {
		r.Method("POST", "/create", http.HandlerFunc(rc.CreateLink))
		r.Method("GET", "/{workerID}", http.HandlerFunc(rc.GetLink))
		r.Method("PATCH", "/{workerID}", http.HandlerFunc(rc.UpdateLink))
		r.Method("PATCH", "/clicks/{workerID}", http.HandlerFunc(rc.ClickAdd))
		r.Method("PATCH", "/registrations/{workerID}", http.HandlerFunc(rc.RegistrationAdd))
	})
}

func RegisterRoleRoutes(router chi.Router, rc *role.RoleController) {
	router.Route("/api/role", func(r chi.Router) {
		r.Method("GET", "/{roleID}", http.HandlerFunc(rc.GetRole))
	})
}

func RegisterSponsorRoutes(router chi.Router, sc *sponsor.SponsorController) {
	router.Route("/api/sponsor", func(r chi.Router) {
		r.Method("POST", "/create", http.HandlerFunc(sc.CreateSponsor))
		r.Method("GET", "/all", http.HandlerFunc(sc.GetSponsors))
		r.Method("GET", "/{id}", http.HandlerFunc(sc.GetSponsor))
		r.Method("DELETE", "/{id}", http.HandlerFunc(sc.DeleteSponsor))
	})
}

func RegisterSubscribeRoutes(router chi.Router, sc *subscribe.SubscribeController) {
	router.Route("/api/subscribe", func(r chi.Router) {
		r.Method("POST", "/create", http.HandlerFunc(sc.AddSubscribe))
		r.Method("GET", "/user/{id}", http.HandlerFunc(sc.GetSubscribesByUser))
		r.Method("GET", "/sponsor/{id}", http.HandlerFunc(sc.GetSubscribesBySponsor))
	})
}
