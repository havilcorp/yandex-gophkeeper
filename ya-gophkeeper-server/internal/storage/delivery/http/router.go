package http

import (
	middleware "yandex-gophkeeper-server/internal/middlewares"

	"github.com/go-chi/chi/v5"
)

func (h *handler) Register(mux *chi.Mux) {
	mux.Route("/storage", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware3(h.conf.JWTKey))
		r.Post("/", h.Create)
		r.Get("/{id}", h.GetOne)
		r.Get("/", h.GetAll)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Remove)
	})
}
