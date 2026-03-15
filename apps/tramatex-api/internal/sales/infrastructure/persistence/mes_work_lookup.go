package persistence

import (
	"context"

	"github.com/google/uuid"
	mes_app "github.com/joran-cortez/tramatex/internal/mes/application"
	sales_app "github.com/joran-cortez/tramatex/internal/sales/application"
)

// MESWorkLookupAdapter implements the sales MESWorkLookup interface
// by delegating to the MES module's WorkOrderQueryService.
type MESWorkLookupAdapter struct {
	mesQueryService *mes_app.WorkOrderQueryService
}

// NewMESWorkLookupAdapter creates a new adapter bridging Sales → MES.
func NewMESWorkLookupAdapter(mesQueryService *mes_app.WorkOrderQueryService) *MESWorkLookupAdapter {
	return &MESWorkLookupAdapter{mesQueryService: mesQueryService}
}

func (a *MESWorkLookupAdapter) GetWorkOrderProgress(ctx context.Context, workOrderID uuid.UUID) (*sales_app.WorkOrderProgress, error) {
	dto, err := a.mesQueryService.GetWorkOrderProgress(ctx, workOrderID)
	if err != nil {
		return nil, err
	}
	return toSalesProgress(dto), nil
}

func (a *MESWorkLookupAdapter) GetWorkOrdersProgress(ctx context.Context, workOrderIDs []uuid.UUID) ([]sales_app.WorkOrderProgress, error) {
	dtos, err := a.mesQueryService.GetWorkOrdersProgress(ctx, workOrderIDs)
	if err != nil {
		return nil, err
	}
	results := make([]sales_app.WorkOrderProgress, 0, len(dtos))
	for i := range dtos {
		results = append(results, *toSalesProgress(&dtos[i]))
	}
	return results, nil
}

// toSalesProgress converts a MES DTO to a Sales-local DTO,
// acting as the anti-corruption layer between modules.
func toSalesProgress(dto *mes_app.WorkOrderProgressDTO) *sales_app.WorkOrderProgress {
	lines := make([]sales_app.WorkOrderLineProgress, 0, len(dto.Lines))
	for _, l := range dto.Lines {
		lines = append(lines, sales_app.WorkOrderLineProgress{
			WorkTypeID:     l.WorkTypeID,
			PositionID:     l.PositionID,
			TotalTasks:     l.TotalTasks,
			CompletedTasks: l.CompletedTasks,
		})
	}
	return &sales_app.WorkOrderProgress{
		WorkOrderID:    dto.WorkOrderID,
		OrderNumber:    dto.OrderNumber,
		OrderName:      dto.OrderName,
		Status:         dto.Status,
		TotalTasks:     dto.TotalTasks,
		CompletedTasks: dto.CompletedTasks,
		Lines:          lines,
	}
}
