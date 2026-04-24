package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func (h *Handler) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		h.logger.Info("HTTP Request",
			slog.String("method", r.Method),
			slog.String("url", r.URL.Path),
			slog.Duration("duration", time.Since(start)),
			slog.String("req_id", middleware.GetReqID(r.Context())),
		)
	})
}
