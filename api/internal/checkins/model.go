// Package checkins implements the /api/v1/checkins mood/energy log.
package checkins

import "time"

// Checkin is a single mood/energy log entry.
type Checkin struct {
	ID        string    `json:"id"`
	Mood      int       `json:"mood"`
	Energy    int       `json:"energy"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}
