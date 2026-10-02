package checkins

import (
	"net/http"

	"github.com/zdub0is/adhd-productivity-app/api/internal/httpx"
)

// Handler serves the /api/v1/checkins endpoints.
type Handler struct {
	store *Store
}

// NewHandler returns a Handler backed by store.
func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

// Register wires the check-in routes onto mux.
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/checkins", h.list)
	mux.HandleFunc("POST /api/v1/checkins", h.create)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	list, err := h.store.List(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list checkins")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, list)
}

type createRequest struct {
	Mood   int    `json:"mood"`
	Energy int    `json:"energy"`
	Note   string `json:"note"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Mood < 1 || req.Mood > 5 {
		httpx.WriteError(w, http.StatusBadRequest, "mood must be between 1 and 5")
		return
	}
	if req.Energy < 1 || req.Energy > 5 {
		httpx.WriteError(w, http.StatusBadRequest, "energy must be between 1 and 5")
		return
	}

	c, err := h.store.Create(r.Context(), req.Mood, req.Energy, req.Note)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to create checkin")
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, c)
}
