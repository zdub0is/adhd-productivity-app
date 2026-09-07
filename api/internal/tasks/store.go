package tasks

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// ErrNotFound is returned when a task lookup by ID finds nothing.
var ErrNotFound = errors.New("task not found")

const dateLayout = "2006-01-02"

// Store persists tasks in Postgres.
type Store struct {
	db *sql.DB
}

// NewStore returns a Store backed by db.
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// NewTask holds the fields accepted when creating a task.
type NewTask struct {
	Title       string
	Description string
	Status      string
	Priority    string
	DueDate     *string
}

// TaskUpdate holds the fields accepted when patching a task. A nil field is
// left unchanged.
type TaskUpdate struct {
	Title       *string
	Description *string
	Status      *string
	Priority    *string
	DueDate     **string
}

func scanTask(row interface {
	Scan(dest ...any) error
}) (Task, error) {
	var t Task
	var due sql.NullTime
	if err := row.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &due, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return Task{}, err
	}
	if due.Valid {
		s := due.Time.Format(dateLayout)
		t.DueDate = &s
	}
	return t, nil
}

func parseDueDate(s *string) (sql.NullTime, error) {
	if s == nil || *s == "" {
		return sql.NullTime{}, nil
	}
	parsed, err := time.Parse(dateLayout, *s)
	if err != nil {
		return sql.NullTime{}, fmt.Errorf("due_date must be YYYY-MM-DD: %w", err)
	}
	return sql.NullTime{Time: parsed, Valid: true}, nil
}

// List returns all tasks, most recently created first.
func (s *Store) List(ctx context.Context) ([]Task, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()

	result := []Task{}
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, fmt.Errorf("list tasks: %w", err)
		}
		result = append(result, t)
	}
	return result, rows.Err()
}

// Get returns the task with the given ID, or ErrNotFound.
func (s *Store) Get(ctx context.Context, id string) (Task, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks WHERE id = $1
	`, id)

	t, err := scanTask(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Task{}, ErrNotFound
	}
	if err != nil {
		return Task{}, fmt.Errorf("get task: %w", err)
	}
	return t, nil
}

// Create inserts a new task.
func (s *Store) Create(ctx context.Context, in NewTask) (Task, error) {
	due, err := parseDueDate(in.DueDate)
	if err != nil {
		return Task{}, err
	}

	row := s.db.QueryRowContext(ctx, `
		INSERT INTO tasks (title, description, status, priority, due_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, title, description, status, priority, due_date, created_at, updated_at
	`, in.Title, in.Description, in.Status, in.Priority, due)

	t, err := scanTask(row)
	if err != nil {
		return Task{}, fmt.Errorf("create task: %w", err)
	}
	return t, nil
}

// Update applies a partial update to the task with the given ID.
func (s *Store) Update(ctx context.Context, id string, in TaskUpdate) (Task, error) {
	current, err := s.Get(ctx, id)
	if err != nil {
		return Task{}, err
	}

	title, description, status, priority := current.Title, current.Description, current.Status, current.Priority
	dueDate := current.DueDate

	if in.Title != nil {
		title = *in.Title
	}
	if in.Description != nil {
		description = *in.Description
	}
	if in.Status != nil {
		status = *in.Status
	}
	if in.Priority != nil {
		priority = *in.Priority
	}
	if in.DueDate != nil {
		dueDate = *in.DueDate
	}

	due, err := parseDueDate(dueDate)
	if err != nil {
		return Task{}, err
	}

	row := s.db.QueryRowContext(ctx, `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, priority = $4, due_date = $5, updated_at = now()
		WHERE id = $6
		RETURNING id, title, description, status, priority, due_date, created_at, updated_at
	`, title, description, status, priority, due, id)

	t, err := scanTask(row)
	if errors.Is(err, sql.ErrNoRows) {
		return Task{}, ErrNotFound
	}
	if err != nil {
		return Task{}, fmt.Errorf("update task: %w", err)
	}
	return t, nil
}

// Delete removes the task with the given ID. It returns ErrNotFound if it
// doesn't exist.
func (s *Store) Delete(ctx context.Context, id string) error {
	res, err := s.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
