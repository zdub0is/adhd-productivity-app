// Package planning implements the Sunsama-style daily planning ritual:
// /api/v1/planning/today, .../start, and .../review.
package planning

import "time"

// Plan is a single day's plan, distinct from the individual tasks selected
// into it.
type Plan struct {
	Date        string     `json:"date"`
	Started     bool       `json:"started"`
	Intention   string     `json:"intention"`
	TaskIDs     []string   `json:"task_ids"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	Reviewed    bool       `json:"reviewed"`
	ReviewNotes string     `json:"review_notes"`
	ReviewedAt  *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

const dateLayout = "2006-01-02"
