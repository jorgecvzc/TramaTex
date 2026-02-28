package security

import (
	"sync"
	"time"
)

// TokenBlacklist defines a blacklist for revoked tokens.
// Used to invalidate access/refresh tokens on logout.
type TokenBlacklist interface {
	Revoke(token string, userID string, expiresAt time.Time)
	IsRevoked(token string) bool
}

// InMemoryTokenBlacklist is a simple in-memory blacklist implementation.
// NOTE: This is process-local and will be cleared on restart.
type InMemoryTokenBlacklist struct {
	items sync.Map
}

// NewInMemoryTokenBlacklist creates a new in-memory blacklist.
func NewInMemoryTokenBlacklist() *InMemoryTokenBlacklist {
	return &InMemoryTokenBlacklist{}
}

// Revoke stores a token until its expiration.
func (b *InMemoryTokenBlacklist) Revoke(token string, userID string, expiresAt time.Time) {
	if token == "" {
		return
	}
	b.items.Store(token, expiresAt)
}

// IsRevoked checks if a token is revoked and not yet expired.
func (b *InMemoryTokenBlacklist) IsRevoked(token string) bool {
	if token == "" {
		return false
	}
	value, ok := b.items.Load(token)
	if !ok {
		return false
	}
	expiresAt, ok := value.(time.Time)
	if !ok {
		b.items.Delete(token)
		return false
	}
	if time.Now().After(expiresAt) {
		b.items.Delete(token)
		return false
	}
	return true
}
