package application

import "context"

// TransactionManager provides service-level database transaction support
// for operations that span multiple repositories.
type TransactionManager interface {
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
