package application

import (
	"context"

	"github.com/google/uuid"
	pricing_app "github.com/joran-cortez/tramatex/internal/pricing/application"
	"github.com/joran-cortez/tramatex/internal/sales/domain"
)

type PricingEngine interface {
	CalculateFinalSalePrice(ctx context.Context, req pricing_app.CalculateFinalSalePriceRequest) (*pricing_app.CalculateFinalSalePriceResponse, error)
}

type PartyLookup interface {
	ExistsParty(ctx context.Context, partyID uuid.UUID) (bool, error)
	HasPartyRole(ctx context.Context, partyID uuid.UUID, role string) (bool, error)
}

type VariantInfo struct {
	ProductName         string
	VariantSKU          string
	OptionConfiguration map[string]string
}

type ProductVariantLookup interface {
	GetVariantInfo(ctx context.Context, variantID uuid.UUID) (*VariantInfo, error)
}

// WorkOrderProgress is a Sales-local DTO representing the execution state
// of a MES WorkOrder. Sales does not know MES internals â€” it receives
// a pre-computed progress snapshot.
type WorkOrderProgress struct {
	WorkOrderID    uuid.UUID
	OrderNumber    string
	OrderName      string
	Status         string // MES production status (opaque to Sales)
	TotalTasks     int
	CompletedTasks int
	Lines          []WorkOrderLineProgress
}

// WorkOrderLineProgress represents progress of one line (work type at a position)
// within a WorkOrder. Sales treats this as read-only information from MES.
type WorkOrderLineProgress struct {
	WorkTypeID     uuid.UUID
	PositionID     uuid.UUID
	TotalTasks     int
	CompletedTasks int
}

// MESWorkLookup provides read-only access to MES WorkOrder execution state.
// Implementation lives in infrastructure as an adapter that calls MES service.
type MESWorkLookup interface {
	GetWorkOrderProgress(ctx context.Context, workOrderID uuid.UUID) (*WorkOrderProgress, error)
	GetWorkOrdersProgress(ctx context.Context, workOrderIDs []uuid.UUID) ([]WorkOrderProgress, error)
}

// WorkOrderCreator creates a MES WorkOrder from Sales when an order is confirmed.
// Implementation lives in infrastructure as an adapter that calls MES service.
type WorkOrderCreator interface {
	CreateWorkOrder(ctx context.Context, workName, partyID, notes string, workSetupID *uuid.UUID, orderWorkSetupID uuid.UUID) (uuid.UUID, error)
}

// WorkOrderSuspender manages MES WorkOrder lifecycle in response to
// Sales order cancellation/reactivation. Implementation lives in
// infrastructure as an adapter that calls MES service.
type WorkOrderSuspender interface {
	SuspendWorkOrders(ctx context.Context, workOrderIDs []uuid.UUID) error
	ReactivateWorkOrders(ctx context.Context, workOrderIDs []uuid.UUID) error
}

// WorkOrderStatusProvider returns the current MES production status for a
// batch of WorkOrder IDs. Implementation lives in infrastructure so Sales
// never queries MES tables directly.
type WorkOrderStatusProvider interface {
	GetWorkOrderStatuses(ctx context.Context, workOrderIDs []uuid.UUID) (map[uuid.UUID]string, error)
}

type DocumentNumberGenerator interface {
	NextQuoteNumber(ctx context.Context) (domain.QuoteNumber, error)
	NextOrderNumber(ctx context.Context) (domain.OrderNumber, error)
	NextDeliveryNoteNumber(ctx context.Context) (domain.DeliveryNoteNumber, error)
	NextInvoiceNumber(ctx context.Context, series domain.InvoiceSeries) (domain.InvoiceNumber, error)
}

type SalesService struct {
	quoteRepo               domain.QuoteRepository
	orderRepo               domain.SalesOrderRepository
	deliveryRepo            domain.DeliveryNoteRepository
	invoiceRepo             domain.InvoiceRepository
	numberGen               DocumentNumberGenerator
	pricingEngine           PricingEngine
	partyLookup             PartyLookup
	productLookup           ProductVariantLookup
	mesLookup               MESWorkLookup
	workOrderCreator        WorkOrderCreator
	workOrderSuspender      WorkOrderSuspender
	workOrderStatusProvider WorkOrderStatusProvider
	txManager               TransactionManager
}

func NewSalesService(
	quoteRepo domain.QuoteRepository,
	orderRepo domain.SalesOrderRepository,
	deliveryRepo domain.DeliveryNoteRepository,
	invoiceRepo domain.InvoiceRepository,
	numberGen DocumentNumberGenerator,
	pricingEngine PricingEngine,
	partyLookup PartyLookup,
	productLookup ProductVariantLookup,
	mesLookup MESWorkLookup,
) *SalesService {
	return &SalesService{
		quoteRepo:     quoteRepo,
		orderRepo:     orderRepo,
		deliveryRepo:  deliveryRepo,
		invoiceRepo:   invoiceRepo,
		numberGen:     numberGen,
		pricingEngine: pricingEngine,
		partyLookup:   partyLookup,
		productLookup: productLookup,
		mesLookup:     mesLookup,
	}
}

// SetWorkOrderCreator configures the optional cross-module WorkOrder creator.

func (s *SalesService) SetWorkOrderCreator(creator WorkOrderCreator) {
	s.workOrderCreator = creator
}

// SetWorkOrderSuspender configures the optional cross-module WorkOrder suspender.

func (s *SalesService) SetWorkOrderSuspender(suspender WorkOrderSuspender) {
	s.workOrderSuspender = suspender
}

// SetWorkOrderStatusProvider configures the optional cross-module status checker.

func (s *SalesService) SetWorkOrderStatusProvider(provider WorkOrderStatusProvider) {
	s.workOrderStatusProvider = provider
}

// SetTransactionManager configures service-level transaction support.

func (s *SalesService) SetTransactionManager(txManager TransactionManager) {
	s.txManager = txManager
}

// runInTransaction wraps fn in a DB transaction if a TransactionManager is configured.
// If no TransactionManager is set (e.g., in tests), fn runs directly.

func (s *SalesService) runInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if s.txManager != nil {
		return s.txManager.RunInTransaction(ctx, fn)
	}
	return fn(ctx)
}

// processMesWorkRefs converts MesWorkRefInputs to domain objects.
// WorkSetupID is optional â€” if provided, it links the reference to an existing MES setup.
// If ID is provided it is reused so existing WorkOrder links (WorkOrderID) are preserved.

func (s *SalesService) processMesWorkRefs(_ context.Context, _ uuid.UUID, inputs []MesWorkRefInput) ([]domain.WorkReference, error) {
	if len(inputs) == 0 {
		return nil, nil
	}
	refs := make([]domain.WorkReference, len(inputs))
	for i, d := range inputs {
		var wsID *uuid.UUID
		if d.WorkSetupID != nil && *d.WorkSetupID != uuid.Nil {
			wsID = d.WorkSetupID
		}
		var woID *uuid.UUID
		if d.WorkOrderID != nil && *d.WorkOrderID != uuid.Nil {
			woID = d.WorkOrderID
		}
		id := uuid.New()
		if d.ID != nil && *d.ID != uuid.Nil {
			id = *d.ID
		}
		refs[i] = domain.WorkReference{
			ID:          id,
			WorkSetupID: wsID,
			WorkOrderID: woID,
			Description: d.Description,
			Sequence:    i + 1,
		}
	}
	return refs, nil
}

func (s *SalesService) ensurePartyExists(ctx context.Context, partyID uuid.UUID) error {
	if s.partyLookup == nil {
		return nil
	}
	exists, err := s.partyLookup.ExistsParty(ctx, partyID)
	if err != nil {
		return err
	}
	if !exists {
		return domain.NewNotFoundError("party not found")
	}
	isClient, err := s.partyLookup.HasPartyRole(ctx, partyID, "CLIENT")
	if err != nil {
		return err
	}
	if !isClient {
		return domain.NewValidationError("party must have CLIENT role to sell")
	}
	return nil
}

type orderLineItemSeed struct {
	ID               *uuid.UUID
	ProductVariantID uuid.UUID
	Quantity         int
	UnitPrice        *MoneyDTO
	DiscountPercent  *float64
}

func toDomainMoney(dto pricing_app.MoneyDTO) (domain.Money, error) {
	return domain.NewMoney(dto.Amount, dto.Currency)
}

func zeroMoney() (domain.Money, error) {
	return domain.NewMoney(0, domain.DefaultCurrency)
}

func (s *SalesService) lookupVariant(ctx context.Context, variantID uuid.UUID) (string, string, map[string]string) {
	if s.productLookup == nil {
		return "", "", nil
	}
	info, err := s.productLookup.GetVariantInfo(ctx, variantID)
	if err != nil || info == nil {
		return "", "", nil
	}
	return info.ProductName, info.VariantSKU, info.OptionConfiguration
}
