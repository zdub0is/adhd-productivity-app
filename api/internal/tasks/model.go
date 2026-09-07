// Package tasks implements the /api/v1/tasks resource.
package tasks

import "time"

// Task is a single to-do item.
type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	DueDate     *string   `json:"due_date,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

var validStatuses = map[string]bool{
	"todo":        true,
	"in_progress": true,
	"done":        true,
	"archived":    true,
}

var validPriorities = map[string]bool{
	"low":    true,
	"normal": true,
	"high":   true,
}
