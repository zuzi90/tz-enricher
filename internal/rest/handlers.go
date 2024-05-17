package rest

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/zuzi90/tz-enricher/docs"
)

func (s *Server) InitRoutes() {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:5005/swagger/doc.json"),
	))
	s.router.Group(func(r chi.Router) {

		r.Get("/metrics", s.metrics)
	})
	s.router.Group(func(r chi.Router) {
		s.router.Route("/api", func(r chi.Router) {
			r.Route("/v1", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Get("/users/{id}", s.getUser)
					r.Get("/users/", s.getUsers)
					r.Post("/users", s.addUser)
					r.Patch("/users/{id}", s.updateUser)
					r.Delete("/users/{id}", s.deleteUser)
				})
			})
		})
	})
}
