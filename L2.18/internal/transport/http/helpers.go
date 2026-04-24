package http

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			h.sendError(w, http.StatusBadRequest, "request body is empty")
		} else {
			h.sendError(w, http.StatusBadRequest, "invalid json format")
		}
		return false
	}
	return true
}

func (h *Handler) sendJSON(w http.ResponseWriter, status int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		h.logger.Error("Failed to marshal response payload", slog.Any("error", err))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"internal server error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(data)
	if err != nil {
		h.logger.Warn("Failed to write response to client", slog.Any("error", err))
	}

}

func (h *Handler) sendOK(w http.ResponseWriter, msg string) {
	h.sendJSON(w, http.StatusOK, map[string]string{"result": msg})
}

func (h *Handler) sendError(w http.ResponseWriter, status int, msg string) {
	h.sendJSON(w, status, map[string]string{"error": msg})
}

func (h *Handler) parseFilterParams(r *http.Request) (int, time.Time, error) {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		return 0, time.Time{}, errors.New("invalid user_id")
	}

	dateStr := r.URL.Query().Get("date")
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return 0, time.Time{}, errors.New("invalid date format, use YYYY-MM-DD")
	}

	return userID, t, nil
}
