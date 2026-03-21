package persistence

import (
	"context"

	mes_app "github.com/joran-cortez/tramatex/internal/mes/application"
	"github.com/joran-cortez/tramatex/internal/sales/domain"
)

// PendingSetupProviderAdapter implements MES PendingWorkSetupProvider by
// querying Sales repositories. This lets MES retrieve pending setups
// without directly accessing Sales tables.
type PendingSetupProviderAdapter struct {
	orderRepo domain.SalesOrderRepository
}

func NewPendingSetupProviderAdapter(orderRepo domain.SalesOrderRepository) *PendingSetupProviderAdapter {
	return &PendingSetupProviderAdapter{orderRepo: orderRepo}
}

func (a *PendingSetupProviderAdapter) ListPending(ctx context.Context) ([]mes_app.PendingWorkSetupDTO, error) {
	status := domain.SalesOrderStatusInPreparation
	filter := domain.SalesOrderFilter{Status: &status}
	orders, err := a.orderRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	var result []mes_app.PendingWorkSetupDTO
	for _, order := range orders {
		for _, ws := range order.WorkReferences {
			if ws.WorkOrderID != nil {
				continue
			}
			result = append(result, mes_app.PendingWorkSetupDTO{
				ID:           ws.ID,
				WorkSetupID:  ws.WorkSetupID,
				Description:  ws.Description,
				OrderID:      order.ID,
				OrderNumber:  order.OrderNumber.String(),
				DeliveryDate: order.DeliveryDate,
				PartyID:      order.PartyID,
			})
		}
	}
	return result, nil
}
