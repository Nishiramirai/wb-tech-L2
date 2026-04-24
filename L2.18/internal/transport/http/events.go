package http

import (
	"calendar/internal/domain"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req createEventRequest
	if ok := h.decodeJSON(w, r, &req); !ok {
		return
	}
	defer func() { _ = r.Body.Close() }()

	event, err := req.ToDomain()
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.CreateEvent(event); err != nil {
		h.logger.Error("service failed to create event", slog.Any("error", err))
		h.sendError(w, http.StatusServiceUnavailable, "internal server error")
		return
	}

	h.sendOK(w, "event created successfully")
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var req updateEventRequest
	if ok := h.decodeJSON(w, r, &req); !ok {
		return
	}
	defer func() { _ = r.Body.Close() }()

	event, err := req.ToDomain()
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.UpdateEvent(event); err != nil {
		h.logger.Info("service failed to update event", slog.Any("error", err))
		h.sendError(w, http.StatusServiceUnavailable, "internal server error")
		return
	}

	h.sendOK(w, "event updated successfully")
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	var req deleteEventRequest
	if ok := h.decodeJSON(w, r, &req); !ok {
		return
	}
	defer func() { _ = r.Body.Close() }()

	eventID, err := uuid.Parse(req.ID)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "invalid id format")
		return
	}

	if err := h.service.DeleteEvent(eventID); err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			h.sendError(w, http.StatusServiceUnavailable, "event not found")
			return
		}

		h.logger.Error("Service error", slog.Any("error", err))
		return
	}

	h.sendOK(w, "event deleted")
}

func (h *Handler) EventsForDay(w http.ResponseWriter, r *http.Request) {
	uid, date, err := h.parseFilterParams(r)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	events, err := h.service.GetEventsForDay(uid, date)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.sendJSON(w, 200, map[string]any{"result": events})
}

func (h *Handler) EventsForWeek(w http.ResponseWriter, r *http.Request) {
	uid, date, err := h.parseFilterParams(r)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	events, err := h.service.GetEventsForWeek(uid, date)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.sendJSON(w, 200, map[string]any{"result": events})
}

func (h *Handler) EventsForMonth(w http.ResponseWriter, r *http.Request) {
	uid, date, err := h.parseFilterParams(r)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	events, err := h.service.GetEventsForMonth(uid, date)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.sendJSON(w, 200, map[string]any{"result": events})
}
