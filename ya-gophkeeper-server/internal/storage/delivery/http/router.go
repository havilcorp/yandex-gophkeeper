package http

import (
	middleware "yandex-gophkeeper-server/internal/middlewares"

	"github.com/go-chi/chi/v5"
)

func (h *handler) Register(mux *chi.Mux) {
	mux.Route("/storage", func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware(h.conf.JWTKey))
		r.Post("/", h.save)
		r.Get("/", h.getAll)
	})
}
