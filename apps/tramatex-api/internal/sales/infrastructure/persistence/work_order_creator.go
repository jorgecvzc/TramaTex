package persistence

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	mes_app "github.com/joran-cortez/tramatex/internal/mes/application"
)

// WorkOrderCreatorAdapter implements the sales WorkOrderCreator interface
// by delegating to the MES service's CreateWorkOrder method.
type WorkOrderCreatorAdapter struct {
	mesService *mes_app.MESService
}

func NewWorkOrderCreatorAdapter(mesService *mes_app.MESService) *WorkOrderCreatorAdapter {
	return &WorkOrderCreatorAdapter{mesService: mesService}
}

// CreateWorkOrder creates a MES WorkOrder and returns its ID.
// workSetupID is optional — if nil, the WorkOrder is created without lines.
func (a *WorkOrderCreatorAdapter) CreateWorkOrder(
	ctx context.Context,
	workName, partyID, notes string,
	workSetupID *uuid.UUID,
	orderWorkSetupID uuid.UUID,
) (uuid.UUID, error) {
	status := "PENDING"
	cmd := mes_app.CreateWorkOrderCommand{
		WorkName:         workName,
		PartyID:          partyID,
		WorkSetupID:      workSetupID,
		Notes:            &notes,
		Status:           &status,
		OrderWorkSetupID: &orderWorkSetupID,
	}
	result, err := a.mesService.CreateWorkOrder(ctx, cmd)
	if err != nil {
		return uuid.Nil, fmt.Errorf("mes create work order: %w", err)
	}
	return result.ID, nil
}
