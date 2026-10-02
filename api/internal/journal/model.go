// Package journal implements the /api/v1/journal endpoints.
package journal

import "time"

// Entry is a single journal entry.
type Entry struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
