package persistence

import (
	"context"

	"github.com/google/uuid"
	mes_app "github.com/joran-cortez/tramatex/internal/mes/application"
)

// WorkOrderSuspenderAdapter implements the sales WorkOrderSuspender interface
// by delegating to the MES service's Suspend/Reactivate methods.
type WorkOrderSuspenderAdapter struct {
	mesService *mes_app.MESService
}

func NewWorkOrderSuspenderAdapter(mesService *mes_app.MESService) *WorkOrderSuspenderAdapter {
	return &WorkOrderSuspenderAdapter{mesService: mesService}
}

func (a *WorkOrderSuspenderAdapter) SuspendWorkOrders(ctx context.Context, workOrderIDs []uuid.UUID) error {
	return a.mesService.SuspendWorkOrders(ctx, workOrderIDs)
}

func (a *WorkOrderSuspenderAdapter) ReactivateWorkOrders(ctx context.Context, workOrderIDs []uuid.UUID) error {
	return a.mesService.ReactivateWorkOrders(ctx, workOrderIDs)
}
