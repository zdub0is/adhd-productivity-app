// Package apikey implements API key generation, storage, and the
// authentication middleware that gates the v1 API.
package apikey

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Key is an API key as returned to a caller. Plaintext is only ever
// populated once, at creation time, and is never stored or returned again.
type Key struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Plaintext  string     `json:"key,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
}

// generate returns a new random API key and the SHA-256 hex hash used to
// store and look it up. API keys carry enough entropy on their own, so a
// fast deterministic hash (rather than a slow password hash like bcrypt) is
// the right tool here: it allows an indexed equality lookup by hash.
func generate() (plaintext string, hash string, err error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", fmt.Errorf("generate key: %w", err)
	}
	plaintext = "ak_" + hex.EncodeToString(buf)
	hash = hashKey(plaintext)
	return plaintext, hash, nil
}

func hashKey(plaintext string) string {
	sum := sha256.Sum256([]byte(plaintext))
	return hex.EncodeToString(sum[:])
}
