package journal

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/zdub0is/adhd-productivity-app/api/internal/httpx"
)

var uuidRe = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// Handler serves the /api/v1/journal endpoints.
type Handler struct {
	store *Store
}

// NewHandler returns a Handler backed by store.
func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

// Register wires the journal routes onto mux.
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/journal", h.list)
	mux.HandleFunc("POST /api/v1/journal", h.create)
	mux.HandleFunc("GET /api/v1/journal/{id}", h.get)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	list, err := h.store.List(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list journal entries")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, list)
}

type createRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Body = strings.TrimSpace(req.Body)
	if req.Body == "" {
		httpx.WriteError(w, http.StatusBadRequest, "body is required")
		return
	}

	e, err := h.store.Create(r.Context(), req.Title, req.Body)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to create journal entry")
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, e)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !uuidRe.MatchString(id) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid journal entry id")
		return
	}

	e, err := h.store.Get(r.Context(), id)
	if errors.Is(err, ErrNotFound) {
		httpx.WriteError(w, http.StatusNotFound, "journal entry not found")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to get journal entry")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, e)
}
