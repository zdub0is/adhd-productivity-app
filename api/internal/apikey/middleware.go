package apikey

import (
	"crypto/subtle"
	"net/http"

	"github.com/zdub0is/adhd-productivity-app/api/internal/httpx"
)

// RequireAPIKey wraps next, rejecting requests that don't carry a valid
// X-API-Key header with 401.
func RequireAPIKey(store *Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("X-API-Key")
			if key == "" {
				httpx.WriteError(w, http.StatusUnauthorized, "missing X-API-Key header")
				return
			}

			ok, err := store.Authenticate(r.Context(), key)
			if err != nil {
				httpx.WriteError(w, http.StatusInternalServerError, "authentication failed")
				return
			}
			if !ok {
				httpx.WriteError(w, http.StatusUnauthorized, "invalid API key")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAdminToken wraps next, rejecting requests that don't carry a valid
// X-Admin-Token header matching the configured bootstrap token. It protects
// the endpoint used to mint the very first API keys.
func RequireAdminToken(adminToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if subtle.ConstantTimeCompare([]byte(r.Header.Get("X-Admin-Token")), []byte(adminToken)) != 1 {
				httpx.WriteError(w, http.StatusUnauthorized, "missing or invalid X-Admin-Token header")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
