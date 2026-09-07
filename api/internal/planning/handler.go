package planning

import (
	"errors"
	"net/http"
	"time"

	"github.com/zdub0is/adhd-productivity-app/api/internal/httpx"
)

// Handler serves the /api/v1/planning endpoints.
type Handler struct {
	store *Store
}

// NewHandler returns a Handler backed by store.
func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

// Register wires the planning routes onto mux.
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/planning/today", h.today)
	mux.HandleFunc("POST /api/v1/planning/{date}/start", h.start)
	mux.HandleFunc("POST /api/v1/planning/{date}/review", h.review)
}

func (h *Handler) today(w http.ResponseWriter, r *http.Request) {
	date := time.Now().UTC().Format(dateLayout)

	p, err := h.store.GetByDate(r.Context(), date)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to get today's plan")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, p)
}

func parseDate(w http.ResponseWriter, r *http.Request) (string, bool) {
	date := r.PathValue("date")
	if _, err := time.Parse(dateLayout, date); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "date must be YYYY-MM-DD")
		return "", false
	}
	return date, true
}

type startRequest struct {
	Intention string   `json:"intention"`
	TaskIDs   []string `json:"task_ids"`
}

func (h *Handler) start(w http.ResponseWriter, r *http.Request) {
	date, ok := parseDate(w, r)
	if !ok {
		return
	}

	var req startRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.TaskIDs == nil {
		req.TaskIDs = []string{}
	}

	p, err := h.store.Start(r.Context(), date, req.Intention, req.TaskIDs)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	httpx.WriteJSON(w, http.StatusOK, p)
}

type reviewRequest struct {
	Notes string `json:"notes"`
}

func (h *Handler) review(w http.ResponseWriter, r *http.Request) {
	date, ok := parseDate(w, r)
	if !ok {
		return
	}

	var req reviewRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	p, err := h.store.Review(r.Context(), date, req.Notes)
	if errors.Is(err, ErrNotFound) {
		httpx.WriteError(w, http.StatusNotFound, "no plan has been started for this date")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to review plan")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, p)
}
