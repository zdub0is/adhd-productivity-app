package apikey

import (
	"context"
	"database/sql"
	"fmt"
)

// Store persists API keys in Postgres.
type Store struct {
	db *sql.DB
}

// NewStore returns a Store backed by db.
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// Create generates a new API key named name, stores its hash, and returns
// the key including its one-time plaintext value.
func (s *Store) Create(ctx context.Context, name string) (Key, error) {
	plaintext, hash, err := generate()
	if err != nil {
		return Key{}, err
	}

	var key Key
	err = s.db.QueryRowContext(ctx, `
		INSERT INTO api_keys (name, key_hash)
		VALUES ($1, $2)
		RETURNING id, name, created_at
	`, name, hash).Scan(&key.ID, &key.Name, &key.CreatedAt)
	if err != nil {
		return Key{}, fmt.Errorf("insert api key: %w", err)
	}

	key.Plaintext = plaintext
	return key, nil
}

// Authenticate looks up the key whose hash matches plaintext and, if found,
// records that it was just used. It returns false if no active key matches.
func (s *Store) Authenticate(ctx context.Context, plaintext string) (bool, error) {
	hash := hashKey(plaintext)

	res, err := s.db.ExecContext(ctx, `
		UPDATE api_keys SET last_used_at = now() WHERE key_hash = $1
	`, hash)
	if err != nil {
		return false, fmt.Errorf("authenticate api key: %w", err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("authenticate api key: %w", err)
	}

	return n > 0, nil
}
