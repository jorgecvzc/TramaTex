package persistence

import (
	"context"

	"gorm.io/gorm"
)

type txContextKey struct{}

// GORMTransactionManager implements application.TransactionManager using GORM.
type GORMTransactionManager struct {
	db *gorm.DB
}

func NewGORMTransactionManager(db *gorm.DB) *GORMTransactionManager {
	return &GORMTransactionManager{db: db}
}

func (m *GORMTransactionManager) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txContextKey{}, tx)
		return fn(txCtx)
	})
}

// getDB returns the transaction from context if present, otherwise the base db.
// This enables transparent transaction propagation across repositories.
func getDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txContextKey{}).(*gorm.DB); ok {
		return tx.WithContext(ctx)
	}
	return db.WithContext(ctx)
}
