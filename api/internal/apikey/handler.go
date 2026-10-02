package apikey

import (
	"net/http"
	"strings"

	"github.com/zdub0is/adhd-productivity-app/api/internal/httpx"
)

// Handler serves the API key management endpoints.
type Handler struct {
	store *Store
}

// NewHandler returns a Handler backed by store.
func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

type createKeyRequest struct {
	Name string `json:"name"`
}

// Create handles POST /api/v1/auth/keys.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req createKeyRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		httpx.WriteError(w, http.StatusBadRequest, "name is required")
		return
	}

	key, err := h.store.Create(r.Context(), name)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to create api key")
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, key)
}
