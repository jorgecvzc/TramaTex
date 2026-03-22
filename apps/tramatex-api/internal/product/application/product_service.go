package application

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/product/domain"
)

// ProductService is the application service for product-related use cases.
type ProductService struct {
	productRepo            domain.ProductRepository
	brandRepo              domain.BrandRepository
	groupRepo              domain.ProductGroupRepository
	attributeRepo          domain.AttributeRepository
	variantRepo            domain.ProductVariantRepository
	partyServiceConfigRepo domain.PartyServiceConfigurationRepository // New
}

// NewProductService creates a new ProductService.
func NewProductService(
	productRepo domain.ProductRepository,
	brandRepo domain.BrandRepository,
	groupRepo domain.ProductGroupRepository,
	attributeRepo domain.AttributeRepository,
	variantRepo domain.ProductVariantRepository,
	partyServiceConfigRepo domain.PartyServiceConfigurationRepository, // New
) *ProductService {
	return &ProductService{
		productRepo:            productRepo,
		brandRepo:              brandRepo,
		groupRepo:              groupRepo,
		attributeRepo:          attributeRepo,
		variantRepo:            variantRepo,
		partyServiceConfigRepo: partyServiceConfigRepo, // New
	}
}

func withActorID(ctx context.Context, actorID string) (context.Context, error) {
	actorID = strings.TrimSpace(actorID)
	if actorID == "" {
		if existing, ok := ctx.Value("actorID").(string); ok && strings.TrimSpace(existing) != "" {
			return ctx, nil
		}
		if existing, ok := ctx.Value("userID").(string); ok && strings.TrimSpace(existing) != "" {
			return context.WithValue(ctx, "actorID", strings.TrimSpace(existing)), nil
		}
		return ctx, domain.NewValidationError("actor ID is required")
	}
	if existing, ok := ctx.Value("actorID").(string); ok && strings.TrimSpace(existing) == actorID {
		return ctx, nil
	}
	return context.WithValue(ctx, "actorID", actorID), nil
}

// CreateProduct handles the UC-P-003: Create a Product.
func (s *ProductService) CreateProduct(ctx context.Context, cmd CreateProductCommand) (*ProductDTO, error) {
	ctx, err := withActorID(ctx, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	// 1. Validate Brand and Groups exist (external aggregates)
	brand, err := s.brandRepo.FindByID(ctx, cmd.BrandID)
	if err != nil {
		return nil, domain.WrapPersistence("brand not found", err)
	}
	if brand == nil {
		return nil, domain.NewNotFoundErrorf("brand with ID %s does not exist", cmd.BrandID)
	}

	for _, groupID := range cmd.GroupIDs {
		group, err := s.groupRepo.FindByID(ctx, groupID)
		if err != nil {
			return nil, domain.WrapPersistence("product group not found", err)
		}
		if group == nil {
			return nil, domain.NewNotFoundErrorf("product group with ID %s does not exist", groupID)
		}
	}

	// 2. Check if SKU already exists
	existingProduct, err := s.productRepo.FindBySKU(ctx, cmd.SKU)
	if err != nil {
		return nil, domain.WrapPersistence("error checking SKU uniqueness", err)
	}
	if existingProduct != nil {
		return nil, domain.NewConflictErrorf("producto con SKU '%s' ya existe", cmd.SKU)
	}

	// 3. Create Product domain entity
	product, err := domain.NewProduct(
		cmd.SKU,
		cmd.Name,
		cmd.LongName,
		cmd.Description,
		cmd.ProductType,
		cmd.BrandID,
		cmd.Barcode,
		cmd.BasePrice,
		cmd.TaxRate,
	)
	if err != nil {
		return nil, domain.WrapValidation("failed to create product domain entity", err)
	}

	// Add groups to product
	for _, groupID := range cmd.GroupIDs {
		product.AddGroup(groupID)
	}

	// 4. Validate and add direct attributes
	for _, attrID := range cmd.DirectAttributeIDs {
		attr, err := s.attributeRepo.FindByID(ctx, attrID)
		if err != nil {
			return nil, domain.WrapPersistence("attribute not found", err)
		}
		if attr == nil {
			return nil, domain.NewNotFoundErrorf("attribute with ID %s does not exist", attrID)
		}
		product.AddDirectAttribute(attrID)
	}

	// 5. Persist Product
	if err := s.productRepo.Save(ctx, product); err != nil {
		return nil, domain.WrapPersistence("failed to save product", err)
	}

	return NewProductDTOFromDomain(product), nil
}

// AddGroupToProduct handles the UC-P-XXX: Add a group to an existing Product.
func (s *ProductService) AddGroupToProduct(ctx context.Context, cmd AddGroupCommand) (*ProductDTO, error) {
	ctx, err := withActorID(ctx, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	product, err := s.productRepo.FindByID(ctx, cmd.ProductID)
	if err != nil {
		return nil, domain.WrapPersistence("product not found", err)
	}
	if product == nil {
		return nil, domain.NewNotFoundErrorf("product with ID %s does not exist", cmd.ProductID)
	}

	group, err := s.groupRepo.FindByID(ctx, cmd.GroupID)
	if err != nil {
		return nil, domain.WrapPersistence("product group not found", err)
	}
	if group == nil {
		return nil, domain.NewNotFoundErrorf("product group with ID %s does not exist", cmd.GroupID)
	}

	product.AddGroup(cmd.GroupID)

	if err := s.productRepo.Save(ctx, product); err != nil {
		return nil, domain.WrapPersistence("failed to add group to product", err)
	}

	return NewProductDTOFromDomain(product), nil
}

// AddDirectAttributeToProduct handles the UC-P-XXX: Add a direct attribute to an existing Product.
func (s *ProductService) AddDirectAttributeToProduct(ctx context.Context, cmd AddDirectAttributeCommand) (*ProductDTO, error) {
	ctx, err := withActorID(ctx, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	product, err := s.productRepo.FindByID(ctx, cmd.ProductID)
	if err != nil {
		return nil, domain.WrapPersistence("product not found", err)
	}
	if product == nil {
		return nil, domain.NewNotFoundErrorf("product with ID %s does not exist", cmd.ProductID)
	}

	// Validate attribute exists using attributeRepo
	attribute, err := s.attributeRepo.FindByID(ctx, cmd.AttributeID)
	if err != nil {
		return nil, domain.WrapPersistence("attribute not found", err)
	}
	if attribute == nil {
		return nil, domain.NewNotFoundErrorf("attribute with ID %s does not exist", cmd.AttributeID)
	}

	product.AddDirectAttribute(cmd.AttributeID)

	if err := s.productRepo.Save(ctx, product); err != nil {
		return nil, domain.WrapPersistence("failed to add direct attribute to product", err)
	}

	return NewProductDTOFromDomain(product), nil
}

// UpdateProductSKU handles the UC-P-006: Modificar SKU de Producto.
func (s *ProductService) UpdateProductSKU(ctx context.Context, cmd UpdateProductSKUCommand) (*ProductDTO, error) {
	ctx, err := withActorID(ctx, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	product, err := s.productRepo.FindByID(ctx, cmd.ProductID)
	if err != nil {
		return nil, domain.WrapPersistence("product not found", err)
	}
	if product == nil {
		return nil, domain.NewNotFoundErrorf("product with ID %s does not exist", cmd.ProductID)
	}

	// Update product's own SKU
	product.SKU = cmd.NewSKU

	// This method in the repository is expected to handle the cascade update of ProductVariant SKUs
	if err := s.productRepo.UpdateSKUs(ctx, product.ID, cmd.NewSKU); err != nil {
		return nil, domain.WrapPersistence("failed to update product SKU and cascade to variants", err)
	}

	// The product returned from repo.FindByID will have the old SKU. We need to fetch the updated product
	// or ensure the Save method updates the entity in memory
	updatedProduct, err := s.productRepo.FindByID(ctx, cmd.ProductID)
	if err != nil {
		return nil, domain.WrapPersistence("failed to retrieve updated product", err)
	}
	if updatedProduct == nil {
		return nil, domain.NewNotFoundErrorf("updated product with ID %s does not exist", cmd.ProductID)
	}

	return NewProductDTOFromDomain(updatedProduct), nil
}

// UpdateProduct updates a product's general information (name, description, brand, groups, attributes).
func (s *ProductService) UpdateProduct(ctx context.Context, cmd UpdateProductCommand) (*ProductDTO, error) {
	ctx, err := withActorID(ctx, cmd.ActorID)
	if err != nil {
		return nil, err
	}

	product, err := s.productRepo.FindByID(ctx, cmd.ProductID)
	if err != nil {
		return nil, domain.WrapPersistence("product not found", err)
	}
	if product == nil {
		return nil, domain.NewNotFoundErrorf("product with ID %s does not exist", cmd.ProductID)
	}

	// Update fields only if provided
	if cmd.Name != nil {
		product.Name = *cmd.Name
	}
	if cmd.LongName != nil {
		product.LongName = *cmd.LongName
	}
	if cmd.SKU != nil {
		product.SKU = *cmd.SKU
	}
	if cmd.Barcode != nil {
		product.Barcode = cmd.Barcode
	}
	if cmd.BasePrice != nil {
		if *cmd.BasePrice < 0 {
			return nil, domain.NewValidationError("base price cannot be negative")
		}
		product.BasePrice = *cmd.BasePrice
	}
	if cmd.TaxRate != nil {
		if *cmd.TaxRate < 0 || *cmd.TaxRate > 100 {
			return nil, domain.NewValidationError("tax rate must be between 0 and 100")
		}
		product.TaxRate = *cmd.TaxRate
	}
	if cmd.IsActive != nil {
		product.IsActive = *cmd.IsActive
	}
	if cmd.ProductType != nil {
		product.ProductType = *cmd.ProductType
	}
	if cmd.Description != nil {
		product.Description = *cmd.Description
	}
	if cmd.BrandID != nil {
		product.BrandID = *cmd.BrandID
	}
	if cmd.GroupIDs != nil {
		product.GroupIDs = cmd.GroupIDs
	}
	if cmd.DirectAttributeIDs != nil {
		product.DirectAttributeIDs = cmd.DirectAttributeIDs
	}

	if err := s.productRepo.Save(ctx, product); err != nil {
		return nil, domain.WrapPersistence("failed to update product", err)
	}

	// Fetch updated product to ensure we return the latest state
	updatedProduct, err := s.productRepo.FindByID(ctx, cmd.ProductID)
	if err != nil {
		return nil, domain.WrapPersistence("failed to retrieve updated product", err)
	}
	if updatedProduct == nil {
		return nil, domain.NewNotFoundErrorf("updated product with ID %s does not exist", cmd.ProductID)
	}

	return NewProductDTOFromDomain(updatedProduct), nil
}

// CreateAttribute handles the UC-P-001: Create an Attribute.
// Note: Scope validations removed for MVP simplicity.
func (s *ProductService) CreateAttribute(ctx context.Context, cmd CreateAttributeCommand) (*AttributeDTO, error) {
	ctx, err := withActorID(ctx, cmd.ActorID)
	if err != nil {
		return nil, err
	}

	// 1. Check if an attribute with this code already exists
	existing, err := s.attributeRepo.FindByCode(ctx, cmd.Code)
	if err != nil {
		return nil, domain.WrapPersistence("failed to check for existing attribute", err)
	}
	if existing != nil {
		return nil, domain.NewValidationErrorf("Ya existe un atributo con el código '%s': %s", cmd.Code, existing.Name)
	}

	// 2. Create Attribute domain entity
	attribute, err := domain.NewAttribute(
		cmd.Name,
		cmd.Code,
	)
	if err != nil {
		return nil, domain.WrapValidation("failed to create attribute domain entity", err)
	}

	// 3. Add values to the attribute
	for _, valCmd := range cmd.Values {
		modifierType := domain.PriceModifierType(valCmd.ModifierType)
		_, err := attribute.AddValueWithModifier(
			valCmd.Value,
			valCmd.Code,
			valCmd.HasPriceModifier,
			modifierType,
			valCmd.ModifierAmount,
		)
		if err != nil {
			return nil, domain.WrapValidation("failed to add attribute value", err)
		}
	}

	// 4. Persist Attribute
	if err := s.attributeRepo.Save(ctx, attribute); err != nil {
		return nil, domain.WrapPersistence("failed to save attribute", err)
	}

	return NewAttributeDTOFromDomain(attribute), nil
}

// UpdateAttribute handles the UC-P-002: Modify an Attribute.
// Note: Scope validations removed for MVP simplicity.
func (s *ProductService) UpdateAttribute(ctx context.Context, cmd UpdateAttributeCommand) (*AttributeDTO, error) {
	ctx, err := withActorID(ctx, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	attribute, err := s.attributeRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, domain.WrapPersistence("attribute not found", err)
	}
	if attribute == nil {
		return nil, domain.NewNotFoundErrorf("attribute with ID %s does not exist", cmd.ID)
	}

	// 1. Update top-level attribute fields
	if cmd.Name != nil {
		attribute.Name = *cmd.Name
	}
	if cmd.Code != nil {
		attribute.Code = *cmd.Code
	}

	// 2. Handle AttributeValue updates (add, modify, delete)
	// Map existing values for efficient lookup
	existingValuesMap := make(map[uuid.UUID]domain.AttributeValue)
	for _, val := range attribute.Values {
		existingValuesMap[val.ID] = val
	}

	for _, valCmd := range cmd.Values {
		if valCmd.ID == nil { // New value
			modifierType := domain.PriceModifierType(valCmd.ModifierType)
			_, err := attribute.AddValueWithModifier(
				valCmd.Value,
				valCmd.Code,
				valCmd.HasPriceModifier,
				modifierType,
				valCmd.ModifierAmount,
			)
			if err != nil {
				return nil, domain.WrapValidation("failed to add new attribute value", err)
			}
		} else { // Existing value, potentially updated
			if _, exists := existingValuesMap[*valCmd.ID]; exists {
				err := attribute.UpdateValue(*valCmd.ID, valCmd.Value, valCmd.Code)
				if err != nil {
					return nil, domain.WrapValidationf(err, "failed to update attribute value %s", valCmd.ID.String())
				}
				// Update price modifier fields directly
				for i := range attribute.Values {
					if attribute.Values[i].ID == *valCmd.ID {
						attribute.Values[i].HasPriceModifier = valCmd.HasPriceModifier
						attribute.Values[i].ModifierType = domain.PriceModifierType(valCmd.ModifierType)
						attribute.Values[i].ModifierAmount = valCmd.ModifierAmount
						break
					}
				}
				delete(existingValuesMap, *valCmd.ID) // Mark as processed
			} else {
				return nil, domain.NewNotFoundErrorf("attribute value with ID %s not found in existing values", valCmd.ID.String())
			}
		}
	}

	// Any values remaining in existingValuesMap should be deleted
	for id := range existingValuesMap {
		err := attribute.RemoveValue(id)
		if err != nil {
			return nil, domain.WrapNotFoundf(err, "failed to remove attribute value %s", id.String())
		}
	}

	// 4. Persist updated Attribute
	if err := s.attributeRepo.Save(ctx, attribute); err != nil {
		return nil, domain.WrapPersistence("failed to save updated attribute", err)
	}

	return NewAttributeDTOFromDomain(attribute), nil
}

// GetApplicableAttributesForProduct handles UC-P-005: Consultar Atributos Aplicables de un Producto.
// Note: Simplified for MVP - only returns directly assigned attributes.
func (s *ProductService) GetApplicableAttributesForProduct(ctx context.Context, productID uuid.UUID) ([]*AttributeDTO, error) {
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, domain.WrapPersistence("product not found", err)
	}
	if product == nil {
		return nil, domain.NewNotFoundErrorf("product with ID %s does not exist", productID)
	}

	// Fetch only directly assigned attributes
	var result []*AttributeDTO
	if len(product.DirectAttributeIDs) == 0 {
		return result, nil
	}

	directAttrs, err := s.attributeRepo.FindByIDs(ctx, product.DirectAttributeIDs)
	if err != nil {
		return nil, domain.WrapPersistence("failed to fetch direct attributes", err)
	}

	// Convert to DTOs - order is preserved from product.DirectAttributeIDs
	// via FindByIDs which returns attributes in the same order as the input IDs
	for i := range directAttrs {
		attr := directAttrs[i]
		result = append(result, NewAttributeDTOFromDomain(&attr))
	}

	return result, nil
}

// GetApplicableAttributesForProductLegacy is the old complex implementation kept for reference.
// Can be removed in future versions.
func (s *ProductService) GetApplicableAttributesForProductLegacy(ctx context.Context, productID uuid.UUID) ([]*AttributeDTO, error) {
	// [Legacy code removed for brevity - see git history if needed]
	return nil, domain.NewNotFoundError("legacy method not implemented")
}

// GetAttributeByID handles fetching a single attribute by its ID.
func (s *ProductService) GetAttributeByID(ctx context.Context, query GetAttributeByIDQuery) (*AttributeDTO, error) {
	attribute, err := s.attributeRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, domain.WrapPersistence("attribute not found", err)
	}
	if attribute == nil {
		return nil, domain.NewNotFoundErrorf("attribute with ID %s does not exist", query.ID)
	}
	return NewAttributeDTOFromDomain(attribute), nil
}

// ListAttributes handles fetching a list of attributes.
// Note: Scope-based filtering removed for MVP simplicity - all attributes are generic.
func (s *ProductService) ListAttributes(ctx context.Context, query ListAttributesQuery) ([]*AttributeDTO, error) {
	attributes, err := s.attributeRepo.FindByScope(ctx, nil, nil)
	if err != nil {
		return nil, domain.WrapPersistence("failed to list attributes", err)
	}

	var dtos []*AttributeDTO
	for _, attr := range attributes {
		dtos = append(dtos, NewAttributeDTOFromDomain(attr))
	}
	return dtos, nil
}

// GetProductByID handles fetching a single product by its ID.
func (s *ProductService) GetProductByID(ctx context.Context, query GetProductByIDQuery) (*ProductDTO, error) {
	product, err := s.productRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, domain.WrapPersistence("product not found", err)
	}
	if product == nil {
		return nil, domain.NewNotFoundErrorf("product with ID %s does not exist", query.ID)
	}
	return NewProductDTOFromDomain(product), nil
}

// ListProducts handles fetching a list of products with optional filtering.
func (s *ProductService) ListProducts(ctx context.Context, query ListProductsQuery) ([]*ProductDTO, error) {
	products, err := s.productRepo.FindAll(ctx)
	if err != nil {
		return nil, domain.WrapPersistence("failed to list products", err)
	}
	if len(products) == 0 {
		return []*ProductDTO{}, nil
	}

	dtos := make([]*ProductDTO, 0, len(products))
	for _, product := range products {
		if !productMatchesQuery(product, query) {
			continue
		}
		dtos = append(dtos, NewProductDTOFromDomain(product))
	}
	return dtos, nil
}

// ListProductVariantsByProductID handles fetching a list of product variants for a given product ID.
func (s *ProductService) ListProductVariantsByProductID(ctx context.Context, query ListProductVariantsByProductIDQuery) ([]*ProductVariantDTO, error) {
	variants, err := s.variantRepo.FindByProductID(ctx, query.ProductID)
	if err != nil {
		return nil, domain.WrapPersistence("failed to list product variants", err)
	}
	if len(variants) == 0 {
		return []*ProductVariantDTO{}, nil
	}

	// Fetch product for baseCost calculation
	product, err := s.productRepo.FindByID(ctx, query.ProductID)
	if err != nil {
		return nil, domain.WrapPersistence("failed to fetch product for baseCost calculation", err)
	}
	if product == nil {
		return nil, domain.NewNotFoundErrorf("product with ID %s does not exist", query.ProductID)
	}

	allAttributes, err := s.attributeRepo.FindByScope(ctx, nil, nil)
	if err != nil {
		return nil, domain.WrapPersistence("failed to fetch attributes for variants", err)
	}

	dtos := make([]*ProductVariantDTO, len(variants))
	for i, variant := range variants {
		dtos[i] = NewProductVariantDTOFromDomain(variant, product, allAttributes)
	}
	return dtos, nil
}

func productMatchesQuery(product *domain.Product, query ListProductsQuery) bool {
	if query.Search != nil {
		s := strings.ToLower(*query.Search)
		nameMatch := strings.Contains(strings.ToLower(product.Name), s)
		skuMatch := strings.Contains(strings.ToLower(product.SKU), s)
		longNameMatch := strings.Contains(strings.ToLower(product.LongName), s)
		if !nameMatch && !skuMatch && !longNameMatch {
			return false
		}
	}
	if query.BrandID != nil && product.BrandID != *query.BrandID {
		return false
	}
	if query.GroupID != nil && !productHasGroup(product, *query.GroupID) {
		return false
	}
	if query.IsActive != nil && product.IsActive != *query.IsActive {
		return false
	}
	if query.ProductType != nil && string(product.ProductType) != *query.ProductType {
		return false
	}
	return true
}

func productHasGroup(product *domain.Product, groupID uuid.UUID) bool {
	for _, id := range product.GroupIDs {
		if id == groupID {
			return true
		}
	}
	return false
}

// GetProductVariantByID handles fetching a single product variant by its ID.
func (s *ProductService) GetProductVariantByID(ctx context.Context, query GetProductVariantByIDQuery) (*ProductVariantDTO, error) {
	variant, err := s.variantRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, domain.WrapPersistence("product variant not found", err)
	}
	if variant == nil {
		return nil, domain.NewNotFoundErrorf("product variant with ID %s does not exist", query.ID)
	}
	// Fetch product for baseCost calculation
	product, err := s.productRepo.FindByID(ctx, variant.ProductID)
	if err != nil {
		return nil, domain.WrapPersistence("failed to fetch product for baseCost calculation", err)
	}
	if product == nil {
		return nil, domain.NewNotFoundErrorf("product with ID %s does not exist", variant.ProductID)
	}
	// Need to fetch all attributes to populate OptionConfiguration in ProductVariantDTO
	allAttributes, err := s.attributeRepo.FindByScope(ctx, nil, nil) // Fetch all attributes for now
	if err != nil {
		return nil, domain.WrapPersistence("failed to fetch attributes for variant", err)
	}
	return NewProductVariantDTOFromDomain(variant, product, allAttributes), nil
}

// GetProductVariantBySKU handles fetching a single product variant by its SKU.
func (s *ProductService) GetProductVariantBySKU(ctx context.Context, query GetProductVariantBySKUQuery) (*ProductVariantDTO, error) {
	variant, err := s.variantRepo.FindBySKU(ctx, query.SKU)
	if err != nil {
		return nil, domain.WrapPersistence("product variant not found", err)
	}
	if variant == nil {
		return nil, domain.NewNotFoundErrorf("product variant with SKU %s does not exist", query.SKU)
	}
	// Fetch product for baseCost calculation
	product, err := s.productRepo.FindByID(ctx, variant.ProductID)
	if err != nil {
		return nil, domain.WrapPersistence("failed to fetch product for baseCost calculation", err)
	}
	if product == nil {
		return nil, domain.NewNotFoundErrorf("product with ID %s does not exist", variant.ProductID)
	}
	allAttributes, err := s.attributeRepo.FindByScope(ctx, nil, nil)
	if err != nil {
		return nil, domain.WrapPersistence("failed to fetch attributes for variant", err)
	}
	return NewProductVariantDTOFromDomain(variant, product, allAttributes), nil
}

// GenerateProductVariants handles UC-P-007: Pre-generate Variants for a Product.
func (s *ProductService) GenerateProductVariants(ctx context.Context, cmd GenerateProductVariantsCommand) error {
	product, err := s.productRepo.FindByID(ctx, cmd.ProductID)
	if err != nil {
		return domain.WrapPersistence("product not found", err)
	}
	if product == nil {
		return domain.NewNotFoundErrorf("product with ID %s does not exist", cmd.ProductID)
	}

	// UC-P-005: Get Applicable Attributes
	applicableAttributes, err := s.GetApplicableAttributesForProduct(ctx, product.ID)
	if err != nil {
		return domain.WrapPersistencef(err, "failed to get applicable attributes for product %s", product.ID)
	}

	// Convert AttributeDTOs back to domain.Attribute for internal use
	var domainAttributes []*domain.Attribute
	for _, attrDTO := range applicableAttributes {
		// This conversion is lossy for values (only string values in DTO)
		// A better approach would be to have a GetApplicableAttributes that returns domain.Attribute directly
		// For now, reconstruct with placeholder attribute values
		var fullDomainAttribute *domain.Attribute
		fullDomainAttribute, err = s.attributeRepo.FindByID(ctx, attrDTO.ID)
		if err != nil || fullDomainAttribute == nil {
			return domain.WrapPersistencef(err, "failed to retrieve full domain attribute %s", attrDTO.ID)
		}
		domainAttributes = append(domainAttributes, fullDomainAttribute)
	}

	if len(domainAttributes) == 0 {
		// No attributes → create or update a single default variant
		return s.ensureDefaultVariant(ctx, product)
	}

	// Attribute order is already correct from GetApplicableAttributesForProduct
	// which preserves the order of product.DirectAttributeIDs

	for _, attr := range domainAttributes {
		if len(attr.Values) == 0 {
			return nil
		}
	}

	combinations := buildAttributeValueCombinations(domainAttributes)
	for _, combination := range combinations {
		attributeValueIDs := make([]uuid.UUID, len(combination.values))
		copy(attributeValueIDs, combination.values)

		generatedSKU, skuErr := domain.GenerateVariantSKU(product.SKU, combination.skuParts)
		if skuErr != nil {
			return domain.WrapValidation("failed to generate variant SKU", skuErr)
		}

		existingVariant, findErr := s.variantRepo.FindByProductIDAndAttributeValues(ctx, product.ID, attributeValueIDs)
		if findErr != nil {
			return domain.WrapPersistence("failed to check existing variant by attributes", findErr)
		}

		if existingVariant != nil && !sameUUIDSet(existingVariant.AttributeValues, attributeValueIDs) {
			existingVariant = nil
		}

		if existingVariant != nil {
			shouldSave := false
			if existingVariant.SKU != generatedSKU {
				existingVariant.SKU = generatedSKU
				shouldSave = true
			}
			if existingVariant.Status != domain.StatusConfirmed {
				existingVariant.Status = domain.StatusConfirmed
				shouldSave = true
			}
			if !existingVariant.IsActive {
				existingVariant.IsActive = true
				shouldSave = true
			}
			if shouldSave {
				if saveErr := s.variantRepo.Save(ctx, existingVariant); saveErr != nil {
					return domain.WrapPersistencef(saveErr, "failed to update variant %s", existingVariant.ID)
				}
			}
			continue
		}

		newVariant, createErr := domain.NewProductVariant(product.ID, generatedSKU, nil, domain.StatusConfirmed, attributeValueIDs)
		if createErr != nil {
			return domain.WrapValidation("failed to create generated variant", createErr)
		}

		if saveErr := s.variantRepo.Save(ctx, newVariant); saveErr != nil {
			return domain.WrapPersistence("failed to save generated variant", saveErr)
		}
	}

	return nil
}

// ensureDefaultVariant creates or reactivates a default variant for a product with no attributes.
func (s *ProductService) ensureDefaultVariant(ctx context.Context, product *domain.Product) error {
	existing, err := s.variantRepo.FindByProductIDAndAttributeValues(ctx, product.ID, []uuid.UUID{})
	if err != nil {
		return domain.WrapPersistence("failed to check existing default variant", err)
	}

	if existing != nil {
		shouldSave := false
		if existing.SKU != product.SKU {
			existing.SKU = product.SKU
			shouldSave = true
		}
		if existing.Status != domain.StatusConfirmed {
			existing.Status = domain.StatusConfirmed
			shouldSave = true
		}
		if !existing.IsActive {
			existing.IsActive = true
			shouldSave = true
		}
		if shouldSave {
			if saveErr := s.variantRepo.Save(ctx, existing); saveErr != nil {
				return domain.WrapPersistencef(saveErr, "failed to update default variant %s", existing.ID)
			}
		}
		return nil
	}

	newVariant, createErr := domain.NewDefaultProductVariant(product.ID, product.SKU)
	if createErr != nil {
		return domain.WrapValidation("failed to create default variant", createErr)
	}

	if saveErr := s.variantRepo.Save(ctx, newVariant); saveErr != nil {
		return domain.WrapPersistence("failed to save default variant", saveErr)
	}

	return nil
}

func buildAttributeValueCombinations(attributes []*domain.Attribute) []struct {
	values          []uuid.UUID
	attributeValues []*domain.AttributeValue
	skuParts        []struct{ AttributeCode, ValueCode string }
} {
	type skuPart = struct{ AttributeCode, ValueCode string }
	type combination = struct {
		values          []uuid.UUID
		attributeValues []*domain.AttributeValue
		skuParts        []skuPart
	}

	combinations := []combination{{
		values:          []uuid.UUID{},
		attributeValues: []*domain.AttributeValue{},
		skuParts:        []skuPart{},
	}}

	for _, attr := range attributes {
		next := make([]combination, 0, len(combinations)*len(attr.Values))
		for _, base := range combinations {
			for i := range attr.Values {
				val := &attr.Values[i]
				newValues := append(append([]uuid.UUID{}, base.values...), val.ID)
				newAttrValues := append(append([]*domain.AttributeValue{}, base.attributeValues...), val)
				newSkuParts := append(append([]skuPart{}, base.skuParts...), skuPart{AttributeCode: attr.Code, ValueCode: val.Code})
				next = append(next, combination{
					values:          newValues,
					attributeValues: newAttrValues,
					skuParts:        newSkuParts,
				})
			}
		}
		combinations = next
	}

	result := make([]struct {
		values          []uuid.UUID
		attributeValues []*domain.AttributeValue
		skuParts        []struct{ AttributeCode, ValueCode string }
	}, len(combinations))

	for i := range combinations {
		result[i] = struct {
			values          []uuid.UUID
			attributeValues []*domain.AttributeValue
			skuParts        []struct{ AttributeCode, ValueCode string }
		}{
			values:          combinations[i].values,
			attributeValues: combinations[i].attributeValues,
			skuParts:        combinations[i].skuParts,
		}
	}

	return result
}

func sameUUIDSet(left, right []uuid.UUID) bool {
	if len(left) != len(right) {
		return false
	}

	seen := make(map[uuid.UUID]int, len(left))
	for _, id := range left {
		seen[id]++
	}
	for _, id := range right {
		count, ok := seen[id]
		if !ok || count == 0 {
			return false
		}
		seen[id]--
	}

	for _, count := range seen {
		if count != 0 {
			return false
		}
	}

	return true
}

// FindOrCreateProductVariant handles UC-P-009: Obtener o Crear Variante (Find or Create).
func (s *ProductService) FindOrCreateProductVariant(ctx context.Context, cmd FindOrCreateProductVariantCommand) (*ProductVariantDTO, error) {
	ctx, err := withActorID(ctx, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	product, err := s.productRepo.FindByID(ctx, cmd.ProductID)
	if err != nil {
		return nil, domain.WrapPersistence("product not found", err)
	}
	if product == nil {
		return nil, domain.NewNotFoundErrorf("product with ID %s does not exist", cmd.ProductID)
	}

	if len(cmd.OptionConfiguration) == 0 {
		// No attributes selected → find or create default variant
		if err := s.ensureDefaultVariant(ctx, product); err != nil {
			return nil, err
		}
		variant, findErr := s.variantRepo.FindByProductIDAndAttributeValues(ctx, product.ID, []uuid.UUID{})
		if findErr != nil || variant == nil {
			return nil, domain.NewNotFoundErrorf("could not find or create default variant for product %s", product.ID)
		}
		return NewProductVariantDTOFromDomain(variant, product, nil), nil
	}

	// Get applicable attributes for validation and SKU construction
	applicableAttributesDTOs, err := s.GetApplicableAttributesForProduct(ctx, product.ID)
	if err != nil {
		return nil, domain.WrapPersistencef(err, "failed to get applicable attributes for product %s", product.ID)
	}

	// Map AttributeCode to domain.Attribute and AttributeValue
	attrCodeToAttribute := make(map[string]*domain.Attribute)
	attrValueToDomainValue := make(map[string]domain.AttributeValue) // Value (string) -> domain.AttributeValue
	// Preserve attribute order from DirectAttributeIDs (via applicableAttributesDTOs)
	attrCodeOrder := make(map[string]int)
	for i, attrDTO := range applicableAttributesDTOs {
		fullDomainAttribute, err := s.attributeRepo.FindByID(ctx, attrDTO.ID)
		if err != nil || fullDomainAttribute == nil {
			return nil, domain.WrapPersistencef(err, "failed to retrieve full domain attribute %s", attrDTO.ID)
		}
		attrCodeToAttribute[fullDomainAttribute.Code] = fullDomainAttribute
		attrCodeOrder[fullDomainAttribute.Code] = i
		for _, val := range fullDomainAttribute.Values {
			attrValueToDomainValue[val.Value] = val
		}
	}

	applicableAttributes := make([]*domain.Attribute, 0, len(attrCodeToAttribute))
	for _, attr := range attrCodeToAttribute {
		applicableAttributes = append(applicableAttributes, attr)
	}

	// Validate OptionConfiguration against applicable attributes and values
	var attributeValueIDs []uuid.UUID
	var sortedAttributeCodes []string
	attrCodeToValueCode := make(map[string]string) // For SKU construction

	for attributeName, value := range cmd.OptionConfiguration {
		attr, attrExists := attrCodeToAttribute[attributeName]
		if !attrExists {
			return nil, domain.NewValidationErrorf("attribute '%s' is not applicable to product '%s'", attributeName, product.ID)
		}
		val, valExists := attrValueToDomainValue[value]
		if !valExists || val.AttributeID != attr.ID { // Ensure value belongs to this attribute
			return nil, domain.NewValidationErrorf("value '%s' is not valid for attribute '%s'", value, attributeName)
		}
		attributeValueIDs = append(attributeValueIDs, val.ID)
		sortedAttributeCodes = append(sortedAttributeCodes, attr.Code)
		attrCodeToValueCode[attr.Code] = val.Code
	}

	// Sort attribute codes by product.DirectAttributeIDs order for deterministic SKU construction
	sort.Slice(sortedAttributeCodes, func(i, j int) bool {
		return attrCodeOrder[sortedAttributeCodes[i]] < attrCodeOrder[sortedAttributeCodes[j]]
	})

	// Construct the deterministic SKU
	generatedSKU := product.SKU
	for _, attrCode := range sortedAttributeCodes {
		valueCode := attrCodeToValueCode[attrCode]
		generatedSKU += fmt.Sprintf("-%s.%s", attrCode, valueCode)
	}

	// Try to find existing variant
	variant, err := s.variantRepo.FindByProductIDAndAttributeValues(ctx, product.ID, attributeValueIDs)
	if err != nil {
		// Log error, but proceed to create if not found
		// For production, differentiate between "not found" error and other errors
	}
	if variant != nil {
		if !variant.IsActive {
			return nil, domain.NewValidationError("La variante está inactiva. Actívala para continuar.")
		}

		// Ensure SKU matches; if not, update. This handles potential SKU changes on product or attributes.
		if variant.SKU != generatedSKU {
			variant.SKU = generatedSKU
			if err := s.variantRepo.Save(ctx, variant); err != nil { // Save might update
				return nil, domain.WrapPersistencef(err, "failed to update SKU of existing variant %s", variant.ID)
			}
		}
		return NewProductVariantDTOFromDomain(variant, product, applicableAttributes), nil
	}

	// If not found, create new variant
	newVariant, err := domain.NewProductVariant(product.ID, generatedSKU, nil, domain.StatusProvisional, attributeValueIDs)
	if err != nil {
		return nil, domain.WrapValidation("failed to create new product variant domain entity", err)
	}

	if err := s.variantRepo.Save(ctx, newVariant); err != nil {
		return nil, domain.WrapPersistence("failed to save new product variant", err)
	}

	return NewProductVariantDTOFromDomain(newVariant, product, applicableAttributes), nil
}

// UpdateProductVariant handles UC-P-008: Modificar una Variante Específica.
func (s *ProductService) UpdateProductVariant(ctx context.Context, cmd UpdateProductVariantCommand) (*ProductVariantDTO, error) {
	ctx, err := withActorID(ctx, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	variant, err := s.variantRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, domain.WrapPersistence("product variant not found", err)
	}
	if variant == nil {
		return nil, domain.NewNotFoundErrorf("product variant with ID %s does not exist", cmd.ID)
	}

	if cmd.Barcode != nil {
		variant.Barcode = cmd.Barcode
	}
	if cmd.IsActive != nil {
		variant.IsActive = *cmd.IsActive
	}
	if cmd.Status != nil {
		variant.Status = *cmd.Status
	} else if variant.Status == domain.StatusProvisional {
		// If other fields are updated and status is PROVISIONAL, automatically confirm it.
		variant.Status = domain.StatusConfirmed
	}

	if err := s.variantRepo.Save(ctx, variant); err != nil {
		return nil, domain.WrapPersistence("failed to save updated product variant", err)
	}

	// Fetch product for baseCost calculation
	product, err := s.productRepo.FindByID(ctx, variant.ProductID)
	if err != nil {
		return nil, domain.WrapPersistence("failed to fetch product for baseCost calculation", err)
	}
	if product == nil {
		return nil, domain.NewNotFoundErrorf("product with ID %s does not exist", variant.ProductID)
	}

	// Need to fetch all attributes to populate OptionConfiguration in ProductVariantDTO
	allAttributes, err := s.attributeRepo.FindByScope(ctx, nil, nil)
	if err != nil {
		return nil, domain.WrapPersistence("failed to fetch attributes for variant", err)
	}

	return NewProductVariantDTOFromDomain(variant, product, allAttributes), nil
}

// GetPartyServiceConfigurationByID handles fetching a single PartyServiceConfiguration by its ID and PartyID.
func (s *ProductService) GetPartyServiceConfigurationByID(ctx context.Context, query GetPartyServiceConfigurationByIDQuery) (*PartyServiceConfigurationDTO, error) {
	config, err := s.partyServiceConfigRepo.FindByID(ctx, query.PartyID, query.ID)
	if err != nil {
		return nil, domain.WrapPersistence("party service configuration not found", err)
	}
	if config == nil {
		return nil, domain.NewNotFoundErrorf("party service configuration with ID %s for party %s does not exist", query.ID, query.PartyID)
	}
	// Assuming a NewPartyServiceConfigurationDTOFromDomain function exists or create inline
	return &PartyServiceConfigurationDTO{
		ID:                   config.ID,
		PartyID:              config.PartyID,
		ServiceID:            config.ServiceID,
		Name:                 config.Name,
		ConfigurationDetails: config.ConfigurationDetails,
	}, nil
}

// ListPartyServiceConfigurationsByPartyID handles fetching all PartyServiceConfigurations for a given PartyID.
func (s *ProductService) ListPartyServiceConfigurationsByPartyID(ctx context.Context, query ListPartyServiceConfigurationsByPartyIDQuery) ([]*PartyServiceConfigurationDTO, error) {
	configs, err := s.partyServiceConfigRepo.FindByPartyID(ctx, query.PartyID)
	if err != nil {
		return nil, domain.WrapPersistence("failed to list party service configurations", err)
	}

	var dtos []*PartyServiceConfigurationDTO
	for _, config := range configs {
		dtos = append(dtos, &PartyServiceConfigurationDTO{
			ID:                   config.ID,
			PartyID:              config.PartyID,
			ServiceID:            config.ServiceID,
			Name:                 config.Name,
			ConfigurationDetails: config.ConfigurationDetails,
		})
	}
	return dtos, nil
}

// CreatePartyServiceConfiguration handles creating a new party service configuration.
func (s *ProductService) CreatePartyServiceConfiguration(ctx context.Context, cmd CreatePartyServiceConfigurationCommand) (*PartyServiceConfigurationDTO, error) {
	ctx, err := withActorID(ctx, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	var configDetails json.RawMessage
	if cmd.ConfigurationDetails != nil {
		payload, err := json.Marshal(cmd.ConfigurationDetails)
		if err != nil {
			return nil, domain.WrapValidation("failed to serialize configuration details", err)
		}
		configDetails = payload
	}

	config, err := domain.NewPartyServiceConfiguration(cmd.PartyID, cmd.ServiceID, cmd.Name, configDetails)
	if err != nil {
		return nil, domain.WrapValidation("failed to create party service configuration domain entity", err)
	}

	if err := s.partyServiceConfigRepo.Save(ctx, config); err != nil {
		return nil, domain.WrapPersistence("failed to save party service configuration", err)
	}
	return &PartyServiceConfigurationDTO{
		ID:                   config.ID,
		PartyID:              config.PartyID,
		ServiceID:            config.ServiceID,
		Name:                 config.Name,
		ConfigurationDetails: config.ConfigurationDetails,
	}, nil
}

// UpdatePartyServiceConfiguration handles updating an existing party service configuration.
func (s *ProductService) UpdatePartyServiceConfiguration(ctx context.Context, cmd UpdatePartyServiceConfigurationCommand) (*PartyServiceConfigurationDTO, error) {
	ctx, err := withActorID(ctx, cmd.ActorID)
	if err != nil {
		return nil, err
	}
	config, err := s.partyServiceConfigRepo.FindByID(ctx, cmd.PartyID, cmd.ID)
	if err != nil {
		return nil, domain.WrapPersistence("party service configuration not found", err)
	}
	if config == nil {
		return nil, domain.NewNotFoundErrorf("party service configuration with ID %s for party %s does not exist", cmd.ID, cmd.PartyID)
	}

	// Apply updates
	serviceID := config.ServiceID
	if cmd.ServiceID != nil {
		serviceID = *cmd.ServiceID
	}
	name := config.Name
	if cmd.Name != nil {
		name = *cmd.Name
	}
	configDetails := config.ConfigurationDetails
	if cmd.ConfigurationDetails != nil {
		payload, err := json.Marshal(cmd.ConfigurationDetails)
		if err != nil {
			return nil, domain.WrapValidation("failed to serialize configuration details", err)
		}
		configDetails = payload
	}

	if err := config.Update(serviceID, name, configDetails); err != nil {
		return nil, domain.WrapValidation("failed to update party service configuration domain entity", err)
	}

	if err := s.partyServiceConfigRepo.Save(ctx, config); err != nil {
		return nil, domain.WrapPersistence("failed to save updated party service configuration", err)
	}
	return &PartyServiceConfigurationDTO{
		ID:                   config.ID,
		PartyID:              config.PartyID,
		ServiceID:            config.ServiceID,
		Name:                 config.Name,
		ConfigurationDetails: config.ConfigurationDetails,
	}, nil
}

// DeletePartyServiceConfiguration handles deleting a party service configuration.
func (s *ProductService) DeletePartyServiceConfiguration(ctx context.Context, cmd DeletePartyServiceConfigurationCommand) error {
	ctx, err := withActorID(ctx, cmd.ActorID)
	if err != nil {
		return err
	}
	err = s.partyServiceConfigRepo.Delete(ctx, cmd.PartyID, cmd.ID)
	if err != nil {
		return domain.WrapPersistence("failed to delete party service configuration", err)
	}
	return nil
}

// ============================================================================
// Brand Queries
// ============================================================================

// BrandDTO represents a brand for API responses
type BrandDTO struct {
	ID                      uuid.UUID `json:"id"`
	Name                    string    `json:"name"`
	DefaultMarkupPercentage float64   `json:"defaultMarkupPercentage"`
	IsActive                bool      `json:"isActive"`
}

// ListBrands retrieves all brands
func (s *ProductService) ListBrands(ctx context.Context) ([]BrandDTO, error) {
	brands, err := s.brandRepo.FindAll(ctx)
	if err != nil {
		return nil, domain.WrapPersistence("failed to list brands", err)
	}

	result := make([]BrandDTO, len(brands))
	for i, brand := range brands {
		result[i] = BrandDTO{
			ID:                      brand.ID,
			Name:                    brand.Name,
			DefaultMarkupPercentage: brand.DefaultMarkupPercentage,
			IsActive:                brand.IsActive,
		}
	}
	return result, nil
}

// GetBrandByID retrieves a brand by ID
func (s *ProductService) GetBrandByID(ctx context.Context, id uuid.UUID) (*BrandDTO, error) {
	brand, err := s.brandRepo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.WrapPersistence("brand not found", err)
	}
	if brand == nil {
		return nil, domain.NewNotFoundErrorf("brand with ID %s does not exist", id)
	}

	return &BrandDTO{
		ID:                      brand.ID,
		Name:                    brand.Name,
		DefaultMarkupPercentage: brand.DefaultMarkupPercentage,
		IsActive:                brand.IsActive,
	}, nil
}

// ============================================================================
// Product Group Queries
// ============================================================================

// ProductGroupDTO represents a product group for API responses
type ProductGroupDTO struct {
	ID            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	Type          string     `json:"type"` // TANGIBLE or SERVICE
	ParentGroupID *uuid.UUID `json:"parent_group_id,omitempty"`
	IsActive      bool       `json:"isActive"`
}

// ListProductGroups retrieves all product groups
func (s *ProductService) ListProductGroups(ctx context.Context) ([]ProductGroupDTO, error) {
	groups, err := s.groupRepo.FindAll(ctx)
	if err != nil {
		return nil, domain.WrapPersistence("failed to list product groups", err)
	}

	result := make([]ProductGroupDTO, len(groups))
	for i, group := range groups {
		result[i] = ProductGroupDTO{
			ID:            group.ID,
			Name:          group.Name,
			Type:          string(group.Type),
			ParentGroupID: group.ParentGroupID,
			IsActive:      group.IsActive,
		}
	}
	return result, nil
}

// GetProductGroupByID retrieves a product group by ID
func (s *ProductService) GetProductGroupByID(ctx context.Context, id uuid.UUID) (*ProductGroupDTO, error) {
	group, err := s.groupRepo.FindByID(ctx, id)
	if err != nil {
		return nil, domain.WrapPersistence("product group not found", err)
	}
	if group == nil {
		return nil, domain.NewNotFoundErrorf("product group with ID %s does not exist", id)
	}

	return &ProductGroupDTO{
		ID:            group.ID,
		Name:          group.Name,
		Type:          string(group.Type),
		ParentGroupID: group.ParentGroupID,
		IsActive:      group.IsActive,
	}, nil
}

// CreateBrand creates a new brand
func (s *ProductService) CreateBrand(ctx context.Context, cmd CreateBrandCommand) (*BrandDTO, error) {
	ctx, err := withActorID(ctx, cmd.ActorID)
	if err != nil {
		return nil, err
	}

	brand, err := domain.NewBrand(cmd.Name, cmd.IsActive, cmd.DefaultMarkupPercentage)
	if err != nil {
		return nil, domain.WrapValidation("failed to create brand domain entity", err)
	}

	if err := s.brandRepo.Save(ctx, brand); err != nil {
		return nil, domain.WrapPersistence("failed to save brand", err)
	}

	return &BrandDTO{
		ID:                      brand.ID,
		Name:                    brand.Name,
		DefaultMarkupPercentage: brand.DefaultMarkupPercentage,
		IsActive:                brand.IsActive,
	}, nil
}

// UpdateBrand updates an existing brand
func (s *ProductService) UpdateBrand(ctx context.Context, cmd UpdateBrandCommand) (*BrandDTO, error) {
	ctx, err := withActorID(ctx, cmd.ActorID)
	if err != nil {
		return nil, err
	}

	brand, err := s.brandRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, domain.WrapPersistence("brand not found", err)
	}
	if brand == nil {
		return nil, domain.NewNotFoundErrorf("brand with ID %s does not exist", cmd.ID)
	}

	if cmd.Name != nil {
		if err := brand.UpdateName(*cmd.Name); err != nil {
			return nil, domain.WrapValidation("failed to update brand name", err)
		}
	}

	if cmd.DefaultMarkupPercentage != nil {
		if *cmd.DefaultMarkupPercentage < 0 {
			return nil, domain.NewValidationError("brand default markup percentage cannot be negative")
		}
		brand.DefaultMarkupPercentage = *cmd.DefaultMarkupPercentage
	}

	if cmd.IsActive != nil {
		brand.IsActive = *cmd.IsActive
	}

	if err := s.brandRepo.Save(ctx, brand); err != nil {
		return nil, domain.WrapPersistence("failed to save brand", err)
	}

	return &BrandDTO{
		ID:                      brand.ID,
		Name:                    brand.Name,
		DefaultMarkupPercentage: brand.DefaultMarkupPercentage,
		IsActive:                brand.IsActive,
	}, nil
}

// DeleteBrand deletes a brand
func (s *ProductService) DeleteBrand(ctx context.Context, cmd DeleteBrandCommand) error {
	ctx, err := withActorID(ctx, cmd.ActorID)
	if err != nil {
		return err
	}

	brand, err := s.brandRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return domain.WrapPersistence("brand not found", err)
	}
	if brand == nil {
		return domain.NewNotFoundErrorf("brand with ID %s does not exist", cmd.ID)
	}

	products, err := s.productRepo.FindAll(ctx)
	if err != nil {
		return domain.WrapPersistence("failed to validate brand usage", err)
	}
	for _, product := range products {
		if product.BrandID == cmd.ID {
			return domain.NewValidationError("cannot delete brand because it is used by one or more products")
		}
	}

	if err := s.brandRepo.Delete(ctx, cmd.ID); err != nil {
		return domain.WrapPersistence("failed to delete brand", err)
	}

	return nil
}

// CreateProductGroup creates a new product group
func (s *ProductService) CreateProductGroup(ctx context.Context, cmd CreateProductGroupCommand) (*ProductGroupDTO, error) {
	ctx, err := withActorID(ctx, cmd.ActorID)
	if err != nil {
		return nil, err
	}

	// Validate parent exists if provided
	if cmd.ParentID != nil {
		parent, err := s.groupRepo.FindByID(ctx, *cmd.ParentID)
		if err != nil {
			return nil, domain.WrapPersistence("parent group not found", err)
		}
		if parent == nil {
			return nil, domain.NewNotFoundErrorf("parent group with ID %s does not exist", *cmd.ParentID)
		}
	}

	// Parse and validate group type
	groupType := domain.ProductGroupType(cmd.Type)
	if !groupType.IsValid() {
		return nil, domain.NewValidationError("invalid group type: must be TANGIBLE or SERVICE")
	}

	group, err := domain.NewProductGroup(cmd.Name, groupType, cmd.ParentID, cmd.IsActive)
	if err != nil {
		return nil, domain.WrapValidation("failed to create product group domain entity", err)
	}

	if err := s.groupRepo.Save(ctx, group); err != nil {
		return nil, domain.WrapPersistence("failed to save product group", err)
	}

	return &ProductGroupDTO{
		ID:            group.ID,
		Name:          group.Name,
		Type:          string(group.Type),
		ParentGroupID: group.ParentGroupID,
		IsActive:      group.IsActive,
	}, nil
}

// UpdateProductGroup updates an existing product group
func (s *ProductService) UpdateProductGroup(ctx context.Context, cmd UpdateProductGroupCommand) (*ProductGroupDTO, error) {
	ctx, err := withActorID(ctx, cmd.ActorID)
	if err != nil {
		return nil, err
	}

	group, err := s.groupRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, domain.WrapPersistence("product group not found", err)
	}
	if group == nil {
		return nil, domain.NewNotFoundErrorf("product group with ID %s does not exist", cmd.ID)
	}

	if cmd.Name != nil {
		if err := group.UpdateName(*cmd.Name); err != nil {
			return nil, domain.WrapValidation("failed to update product group name", err)
		}
	}

	if cmd.Type != nil {
		groupType := domain.ProductGroupType(*cmd.Type)
		if err := group.UpdateType(groupType); err != nil {
			return nil, domain.WrapValidation("failed to update product group type", err)
		}
	}

	if cmd.ParentID != nil {
		// Validate parent exists
		parent, err := s.groupRepo.FindByID(ctx, *cmd.ParentID)
		if err != nil {
			return nil, domain.WrapPersistence("parent group not found", err)
		}
		if parent == nil {
			return nil, domain.NewNotFoundErrorf("parent group with ID %s does not exist", *cmd.ParentID)
		}
		group.ParentGroupID = cmd.ParentID
	}

	if cmd.IsActive != nil {
		group.IsActive = *cmd.IsActive
	}

	if err := s.groupRepo.Save(ctx, group); err != nil {
		return nil, domain.WrapPersistence("failed to save product group", err)
	}

	return &ProductGroupDTO{
		ID:            group.ID,
		Name:          group.Name,
		Type:          string(group.Type),
		ParentGroupID: group.ParentGroupID,
		IsActive:      group.IsActive,
	}, nil
}

// DeleteProductGroup deletes a product group
func (s *ProductService) DeleteProductGroup(ctx context.Context, cmd DeleteProductGroupCommand) error {
	ctx, err := withActorID(ctx, cmd.ActorID)
	if err != nil {
		return err
	}

	group, err := s.groupRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return domain.WrapPersistence("product group not found", err)
	}
	if group == nil {
		return domain.NewNotFoundErrorf("product group with ID %s does not exist", cmd.ID)
	}

	allGroups, err := s.groupRepo.FindAll(ctx)
	if err != nil {
		return domain.WrapPersistence("failed to validate product group usage", err)
	}
	for _, g := range allGroups {
		if g.ParentGroupID != nil && *g.ParentGroupID == cmd.ID {
			return domain.NewValidationError("cannot delete product group because it is parent of another group")
		}
	}

	products, err := s.productRepo.FindAll(ctx)
	if err != nil {
		return domain.WrapPersistence("failed to validate product group usage", err)
	}
	for _, product := range products {
		if containsUUID(product.GroupIDs, cmd.ID) {
			return domain.NewValidationError("cannot delete product group because it is used by one or more products")
		}
	}

	if err := s.groupRepo.Delete(ctx, cmd.ID); err != nil {
		return domain.WrapPersistence("failed to delete product group", err)
	}

	return nil
}

// DeleteAttribute deletes an attribute
func (s *ProductService) DeleteAttribute(ctx context.Context, cmd DeleteAttributeCommand) error {
	ctx, err := withActorID(ctx, cmd.ActorID)
	if err != nil {
		return err
	}

	attribute, err := s.attributeRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return domain.WrapPersistence("attribute not found", err)
	}
	if attribute == nil {
		return domain.NewNotFoundErrorf("attribute with ID %s does not exist", cmd.ID)
	}

	products, err := s.productRepo.FindAll(ctx)
	if err != nil {
		return domain.WrapPersistence("failed to validate attribute usage", err)
	}
	for _, product := range products {
		if containsUUID(product.DirectAttributeIDs, cmd.ID) {
			return domain.NewValidationError("cannot delete attribute because it is assigned to one or more products")
		}
	}

	attributeValueIDs := make(map[uuid.UUID]struct{}, len(attribute.Values))
	for _, value := range attribute.Values {
		attributeValueIDs[value.ID] = struct{}{}
	}
	if len(attributeValueIDs) > 0 {
		for _, product := range products {
			variants, err := s.variantRepo.FindByProductID(ctx, product.ID)
			if err != nil {
				return domain.WrapPersistence("failed to validate attribute usage", err)
			}
			for _, variant := range variants {
				for _, variantValueID := range variant.AttributeValues {
					if _, used := attributeValueIDs[variantValueID]; used {
						return domain.NewValidationError("cannot delete attribute because one or more product variants depend on its values")
					}
				}
			}
		}
	}

	if err := s.attributeRepo.Delete(ctx, cmd.ID); err != nil {
		return domain.WrapPersistence("failed to delete attribute", err)
	}

	return nil
}

func containsUUID(ids []uuid.UUID, target uuid.UUID) bool {
	for _, id := range ids {
		if id == target {
			return true
		}
	}
	return false
}

// SmartSearch performs an intelligent search for products/variants by SKU, barcode, or partial reference.
// Search priority:
// 1. Exact variant SKU match
// 2. Exact variant barcode match
// 3. Exact product SKU match (returns product + option sets)
// 4. Exact product barcode match
// 5. Prefix match on variant SKU (parses pre-selected attributes from SKU suffix)
// 6. Prefix match on product SKU
// 7. Text search on product name/SKU/longName
func (s *ProductService) SmartSearch(ctx context.Context, query SmartSearchQuery) (*SmartSearchResultDTO, error) {
	q := strings.TrimSpace(query.Query)
	if q == "" {
		return &SmartSearchResultDTO{Type: "no_match"}, nil
	}

	// 1. Exact variant SKU match
	variant, err := s.variantRepo.FindBySKU(ctx, q)
	if err != nil {
		return nil, domain.WrapPersistence("smart search: variant SKU lookup failed", err)
	}
	if variant != nil {
		product, err := s.productRepo.FindByID(ctx, variant.ProductID)
		if err != nil {
			return nil, domain.WrapPersistence("smart search: product lookup failed", err)
		}
		if product != nil {
			allAttributes, _ := s.attributeRepo.FindByScope(ctx, nil, nil)
			variantDTO := NewProductVariantDTOFromDomain(variant, product, allAttributes)
			variantDTO.ProductName = product.Name
			return &SmartSearchResultDTO{
				Type:    "exact_variant",
				Product: NewProductDTOFromDomain(product),
				Variant: variantDTO,
			}, nil
		}
	}

	// 2. Exact variant barcode match
	variant, err = s.variantRepo.FindByBarcode(ctx, q)
	if err != nil {
		return nil, domain.WrapPersistence("smart search: variant barcode lookup failed", err)
	}
	if variant != nil {
		product, err := s.productRepo.FindByID(ctx, variant.ProductID)
		if err != nil {
			return nil, domain.WrapPersistence("smart search: product lookup failed", err)
		}
		if product != nil {
			allAttributes, _ := s.attributeRepo.FindByScope(ctx, nil, nil)
			variantDTO := NewProductVariantDTOFromDomain(variant, product, allAttributes)
			variantDTO.ProductName = product.Name
			return &SmartSearchResultDTO{
				Type:    "exact_variant",
				Product: NewProductDTOFromDomain(product),
				Variant: variantDTO,
			}, nil
		}
	}

	// 3. Exact product SKU match → return product + its option sets
	product, err := s.productRepo.FindBySKU(ctx, q)
	if err != nil {
		return nil, domain.WrapPersistence("smart search: product SKU lookup failed", err)
	}
	if product != nil {
		return s.buildExactProductResult(ctx, product)
	}

	// 4. Exact product barcode match
	product, err = s.productRepo.FindByBarcode(ctx, q)
	if err != nil {
		return nil, domain.WrapPersistence("smart search: product barcode lookup failed", err)
	}
	if product != nil {
		return s.buildExactProductResult(ctx, product)
	}

	// 5. Prefix match on variant SKU → parse pre-selected attributes
	if strings.Contains(q, "-") {
		// Try to find variants with this prefix
		variants, err := s.variantRepo.FindBySKUPrefix(ctx, q)
		if err != nil {
			return nil, domain.WrapPersistence("smart search: variant prefix lookup failed", err)
		}
		if len(variants) > 0 {
			// All matching variants belong to the same product (they share the product SKU prefix)
			product, err := s.productRepo.FindByID(ctx, variants[0].ProductID)
			if err != nil {
				return nil, domain.WrapPersistence("smart search: product lookup failed", err)
			}
			if product != nil {
				return s.buildPartialMatchResult(ctx, product, q, variants)
			}
		}

		// Try to extract the product SKU part (before the first "-ATTR." segment)
		productSKUPart := s.extractProductSKUFromPartialRef(q)
		if productSKUPart != "" && productSKUPart != q {
			product, err := s.productRepo.FindBySKU(ctx, productSKUPart)
			if err == nil && product != nil {
				return s.buildPartialMatchResult(ctx, product, q, nil)
			}
		}
	}

	// 6. Prefix match on product SKU
	products, err := s.productRepo.FindBySKUPrefix(ctx, q)
	if err != nil {
		return nil, domain.WrapPersistence("smart search: product prefix lookup failed", err)
	}
	if len(products) == 1 {
		// Single match → treat as exact product
		return s.buildExactProductResult(ctx, products[0])
	}
	if len(products) > 1 {
		dtos := make([]*ProductDTO, len(products))
		for i, p := range products {
			dtos[i] = NewProductDTOFromDomain(p)
		}
		return &SmartSearchResultDTO{
			Type:     "product_list",
			Products: dtos,
		}, nil
	}

	// 7. Text search fallback (name, SKU, longName)
	allProducts, err := s.productRepo.FindAll(ctx)
	if err != nil {
		return nil, domain.WrapPersistence("smart search: product list failed", err)
	}
	qLower := strings.ToLower(q)
	var matchedProducts []*ProductDTO
	for _, p := range allProducts {
		if strings.Contains(strings.ToLower(p.Name), qLower) ||
			strings.Contains(strings.ToLower(p.SKU), qLower) ||
			strings.Contains(strings.ToLower(p.LongName), qLower) {
			matchedProducts = append(matchedProducts, NewProductDTOFromDomain(p))
		}
		if len(matchedProducts) >= 20 {
			break
		}
	}
	if len(matchedProducts) > 0 {
		return &SmartSearchResultDTO{
			Type:     "product_list",
			Products: matchedProducts,
		}, nil
	}

	return &SmartSearchResultDTO{Type: "no_match"}, nil
}

// buildExactProductResult builds a result for an exact product match, including option sets.
func (s *ProductService) buildExactProductResult(ctx context.Context, product *domain.Product) (*SmartSearchResultDTO, error) {
	optionSets, err := s.GetApplicableAttributesForProduct(ctx, product.ID)
	if err != nil {
		optionSets = []*AttributeDTO{} // Non-fatal
	}
	return &SmartSearchResultDTO{
		Type:       "exact_product",
		Product:    NewProductDTOFromDomain(product),
		OptionSets: optionSets,
	}, nil
}

// buildPartialMatchResult builds a result for a partial SKU match, extracting pre-selected attributes.
func (s *ProductService) buildPartialMatchResult(ctx context.Context, product *domain.Product, query string, matchingVariants []*domain.ProductVariant) (*SmartSearchResultDTO, error) {
	optionSets, err := s.GetApplicableAttributesForProduct(ctx, product.ID)
	if err != nil {
		optionSets = []*AttributeDTO{} // Non-fatal
	}

	// Parse the query suffix to extract pre-selected attributes
	// SKU format: PRODUCT_SKU-ATTR_CODE.VALUE_CODE-ATTR_CODE.VALUE_CODE
	// Use case-insensitive prefix removal since the user may type in different case
	selectedAttributes := make(map[string]string)
	suffix := ""
	if len(query) > len(product.SKU) && strings.EqualFold(query[:len(product.SKU)], product.SKU) {
		suffix = query[len(product.SKU):]
	}
	if strings.HasPrefix(suffix, "-") {
		suffix = suffix[1:]
	}
	if suffix != "" {
		parts := strings.Split(suffix, "-")
		for _, part := range parts {
			dotIdx := strings.Index(part, ".")
			if dotIdx > 0 && dotIdx < len(part)-1 {
				attrCode := strings.ToUpper(part[:dotIdx])
				valueCode := strings.ToUpper(part[dotIdx+1:])
				selectedAttributes[attrCode] = valueCode
			}
		}
	}

	// Build variant DTOs if provided
	var variantDTOs []*ProductVariantDTO
	if len(matchingVariants) > 0 {
		allAttributes, _ := s.attributeRepo.FindByScope(ctx, nil, nil)
		for _, v := range matchingVariants {
			dto := NewProductVariantDTOFromDomain(v, product, allAttributes)
			dto.ProductName = product.Name
			variantDTOs = append(variantDTOs, dto)
		}
	}

	return &SmartSearchResultDTO{
		Type:               "partial_match",
		Product:            NewProductDTOFromDomain(product),
		OptionSets:         optionSets,
		SelectedAttributes: selectedAttributes,
		MatchingVariants:   variantDTOs,
	}, nil
}

// extractProductSKUFromPartialRef tries to extract the product SKU from a partial variant reference.
// For example, "FY5678-SIZE.M" → "FY5678"
func (s *ProductService) extractProductSKUFromPartialRef(ref string) string {
	// Find the first segment that looks like ATTR_CODE.VALUE_CODE
	parts := strings.Split(ref, "-")
	var productSKUParts []string
	for _, part := range parts {
		if strings.Contains(part, ".") {
			break
		}
		productSKUParts = append(productSKUParts, part)
	}
	return strings.Join(productSKUParts, "-")
}
