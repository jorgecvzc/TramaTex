package security

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
)

type RevokedTokenModel struct {
	TokenID   string    `gorm:"column:token_id;primaryKey;type:varchar(255)"`
	UserID    string    `gorm:"column:user_id;type:uuid;not null"`
	ExpiresAt time.Time `gorm:"type:timestamp with time zone"`
	CreatedAt time.Time `gorm:"column:revoked_at;autoCreateTime:milli;type:timestamp with time zone"`
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
func (b *PostgresTokenBlacklist) Revoke(token string, userID string, expiresAt time.Time) {
	if token == "" {
		return
	}
	hash := sha256.Sum256([]byte(token))
	model := RevokedTokenModel{
		TokenID:   hex.EncodeToString(hash[:]),
		UserID:    userID,
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
	err := b.db.First(&model, "token_id = ?", hex.EncodeToString(hash[:])).Error
	if err != nil {
		return false
	}

	if time.Now().After(model.ExpiresAt) {
		_ = b.db.Delete(&model).Error
		return false
	}
	return true
}
