package persistence

import (
	"context"

	"github.com/google/uuid"
	mes_app "github.com/joran-cortez/tramatex/internal/mes/application"
)

// WorkOrderStatusCheckerAdapter implements sales.WorkOrderStatusProvider by
// delegating to the MES service. Sales never queries MES tables directly.
type WorkOrderStatusCheckerAdapter struct {
	mesService *mes_app.MESService
}

func NewWorkOrderStatusCheckerAdapter(mesService *mes_app.MESService) *WorkOrderStatusCheckerAdapter {
	return &WorkOrderStatusCheckerAdapter{mesService: mesService}
}

func (a *WorkOrderStatusCheckerAdapter) GetWorkOrderStatuses(ctx context.Context, workOrderIDs []uuid.UUID) (map[uuid.UUID]string, error) {
	result := make(map[uuid.UUID]string, len(workOrderIDs))
	for _, id := range workOrderIDs {
		dto, err := a.mesService.GetWorkOrderByID(ctx, mes_app.GetWorkOrderByIDQuery{ID: id})
		if err != nil || dto == nil {
			continue
		}
		result[id] = dto.Status
	}
	return result, nil
}
