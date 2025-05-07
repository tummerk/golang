package rest

import "github.com/go-chi/chi/v5"

func (s Server) RegisterRoutes(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Route("/schedules", func(r chi.Router) {
			r.Get("/", s.GetUserSchedules)
		})
		r.Route("/schedule", func(r chi.Router) {
			r.Get("/", s.GetUserSchedule)
			r.Post("/", s.CreateUserSchedule)
		})
		r.Route("/next_takings", func(r chi.Router) {
			r.Get("/", s.NextTakings)
		})
	})

}
