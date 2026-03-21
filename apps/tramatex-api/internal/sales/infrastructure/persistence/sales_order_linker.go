package persistence

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SalesOrderLinkerAdapter implements the MES SalesOrderLinker interface.
// It updates order_work_setups.work_order_id so Sales can track which
// MES WorkOrder was generated for a given work setup config.
type SalesOrderLinkerAdapter struct {
	db *gorm.DB
}

func NewSalesOrderLinkerAdapter(db *gorm.DB) *SalesOrderLinkerAdapter {
	return &SalesOrderLinkerAdapter{db: db}
}

func (a *SalesOrderLinkerAdapter) LinkWorkOrder(ctx context.Context, orderWorkSetupID uuid.UUID, workOrderID uuid.UUID) error {
	result := a.db.WithContext(ctx).
		Model(&OrderWorkRefModel{}).
		Where("id = ? AND deleted_at IS NULL", orderWorkSetupID).
		Update("work_order_id", workOrderID)
	if result.Error != nil {
		return fmt.Errorf("update order_work_setup: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("order work setup %s not found", orderWorkSetupID)
	}
	return nil
}
