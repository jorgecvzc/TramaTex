package persistence

import (
	"context"

	"github.com/google/uuid"
	mes_app "github.com/joran-cortez/tramatex/internal/mes/application"
	"gorm.io/gorm"
)

// WorkOrderSalesInfoAdapter implements MES WorkOrderSalesInfoProvider by
// querying the order_work_setups + sales_orders tables. This lets MES
// enrich WorkOrders with their Sales order references without directly
// accessing Sales tables.
type WorkOrderSalesInfoAdapter struct {
	db *gorm.DB
}

func NewWorkOrderSalesInfoAdapter(db *gorm.DB) *WorkOrderSalesInfoAdapter {
	return &WorkOrderSalesInfoAdapter{db: db}
}

type workOrderSalesRow struct {
	WorkOrderID uuid.UUID `gorm:"column:work_order_id"`
	OrderID     uuid.UUID `gorm:"column:order_id"`
	OrderNumber string    `gorm:"column:order_number"`
}

func (a *WorkOrderSalesInfoAdapter) GetSalesInfoByWorkOrderIDs(ctx context.Context, workOrderIDs []uuid.UUID) (map[uuid.UUID]mes_app.WorkOrderSalesInfo, error) {
	if len(workOrderIDs) == 0 {
		return map[uuid.UUID]mes_app.WorkOrderSalesInfo{}, nil
	}

	var rows []workOrderSalesRow
	err := a.db.WithContext(ctx).
		Table("order_work_setups AS ows").
		Select("ows.work_order_id, so.id AS order_id, so.order_number").
		Joins("JOIN sales_orders so ON so.id = ows.order_id AND so.deleted_at IS NULL").
		Where("ows.work_order_id IN ? AND ows.deleted_at IS NULL", workOrderIDs).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make(map[uuid.UUID]mes_app.WorkOrderSalesInfo, len(rows))
	for _, row := range rows {
		result[row.WorkOrderID] = mes_app.WorkOrderSalesInfo{
			SalesOrderID:     row.OrderID,
			SalesOrderNumber: row.OrderNumber,
		}
	}
	return result, nil
}
