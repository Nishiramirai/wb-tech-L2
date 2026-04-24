package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *Handler) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(h.LoggingMiddleware)

	router.Post("/create_event", h.CreateEvent)
	router.Post("/update_event", h.UpdateEvent)
	router.Post("/delete_event", h.DeleteEvent)

	router.Get("/events_for_day", h.EventsForDay)
	router.Get("/events_for_week", h.EventsForWeek)
	router.Get("/events_for_month", h.EventsForMonth)

	return router
}
