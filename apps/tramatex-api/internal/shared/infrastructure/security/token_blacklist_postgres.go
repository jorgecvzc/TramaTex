package security

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
)

type RevokedTokenModel struct {
	TokenHash string    `gorm:"primaryKey;type:varchar(64)"`
	ExpiresAt time.Time `gorm:"type:timestamp with time zone"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli;type:timestamp with time zone"`
}

func (RevokedTokenModel) TableName() string {
	return "revoked_tokens"
}

// PostgresTokenBlacklist persists revoked tokens in PostgreSQL.
type PostgresTokenBlacklist struct {
	db *gorm.DB
}

// NewPostgresTokenBlacklist creates a new persistent blacklist.
func NewPostgresTokenBlacklist(db *gorm.DB) *PostgresTokenBlacklist {
	return &PostgresTokenBlacklist{db: db}
}

// Revoke stores a token hash until its expiration.
func (b *PostgresTokenBlacklist) Revoke(token string, expiresAt time.Time) {
	if token == "" {
		return
	}
	hash := sha256.Sum256([]byte(token))
	model := RevokedTokenModel{
		TokenHash: hex.EncodeToString(hash[:]),
		ExpiresAt: expiresAt,
	}

	_ = b.db.Create(&model).Error
}

// IsRevoked checks if a token is revoked and not yet expired.
func (b *PostgresTokenBlacklist) IsRevoked(token string) bool {
	if token == "" {
		return false
	}
	if b.db == nil {
		return false
	}

	hash := sha256.Sum256([]byte(token))
	var model RevokedTokenModel
	err := b.db.First(&model, "token_hash = ?", hex.EncodeToString(hash[:])).Error
	if err != nil {
		return false
	}

	if time.Now().After(model.ExpiresAt) {
		_ = b.db.Delete(&model).Error
		return false
	}
	return true
}
