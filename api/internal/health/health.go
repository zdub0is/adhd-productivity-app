// Package health implements the service's liveness check.
package health

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Status string `json:"status"`
}

// Handler responds 200 OK with {"status":"ok"} to confirm the process is up.
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response{Status: "ok"})
}
