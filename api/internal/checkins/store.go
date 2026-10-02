package checkins

import (
	"context"
	"database/sql"
	"fmt"
)

// Store persists check-ins in Postgres.
type Store struct {
	db *sql.DB
}

// NewStore returns a Store backed by db.
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// List returns all check-ins, most recent first.
func (s *Store) List(ctx context.Context) ([]Checkin, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, mood, energy, note, created_at
		FROM checkins
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("list checkins: %w", err)
	}
	defer rows.Close()

	result := []Checkin{}
	for rows.Next() {
		var c Checkin
		if err := rows.Scan(&c.ID, &c.Mood, &c.Energy, &c.Note, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("list checkins: %w", err)
		}
		result = append(result, c)
	}
	return result, rows.Err()
}

// Create inserts a new check-in.
func (s *Store) Create(ctx context.Context, mood, energy int, note string) (Checkin, error) {
	var c Checkin
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO checkins (mood, energy, note)
		VALUES ($1, $2, $3)
		RETURNING id, mood, energy, note, created_at
	`, mood, energy, note).Scan(&c.ID, &c.Mood, &c.Energy, &c.Note, &c.CreatedAt)
	if err != nil {
		return Checkin{}, fmt.Errorf("create checkin: %w", err)
	}
	return c, nil
}
