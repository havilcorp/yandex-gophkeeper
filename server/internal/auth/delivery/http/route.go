package http

import "github.com/go-chi/chi/v5"

func (h *handler) Register(mux *chi.Mux) {
	mux.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.login)
		r.Post("/registration", h.registration)
	})
}
