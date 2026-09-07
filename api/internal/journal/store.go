package journal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ErrNotFound is returned when an entry lookup by ID finds nothing.
var ErrNotFound = errors.New("journal entry not found")

// Store persists journal entries in Postgres.
type Store struct {
	db *sql.DB
}

// NewStore returns a Store backed by db.
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// List returns all journal entries, most recent first.
func (s *Store) List(ctx context.Context) ([]Entry, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, body, created_at, updated_at
		FROM journal_entries
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("list journal entries: %w", err)
	}
	defer rows.Close()

	result := []Entry{}
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.ID, &e.Title, &e.Body, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, fmt.Errorf("list journal entries: %w", err)
		}
		result = append(result, e)
	}
	return result, rows.Err()
}

// Get returns the entry with the given ID, or ErrNotFound.
func (s *Store) Get(ctx context.Context, id string) (Entry, error) {
	var e Entry
	err := s.db.QueryRowContext(ctx, `
		SELECT id, title, body, created_at, updated_at
		FROM journal_entries WHERE id = $1
	`, id).Scan(&e.ID, &e.Title, &e.Body, &e.CreatedAt, &e.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Entry{}, ErrNotFound
	}
	if err != nil {
		return Entry{}, fmt.Errorf("get journal entry: %w", err)
	}
	return e, nil
}

// Create inserts a new journal entry.
func (s *Store) Create(ctx context.Context, title, body string) (Entry, error) {
	var e Entry
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO journal_entries (title, body)
		VALUES ($1, $2)
		RETURNING id, title, body, created_at, updated_at
	`, title, body).Scan(&e.ID, &e.Title, &e.Body, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return Entry{}, fmt.Errorf("create journal entry: %w", err)
	}
	return e, nil
}
