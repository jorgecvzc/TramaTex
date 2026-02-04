package application

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

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

// CreateProduct handles the UC-P-003: Create a Product.
func (s *ProductService) CreateProduct(ctx context.Context, cmd CreateProductCommand) (*ProductDTO, error) {
	// 1. Validate Brand and Groups exist (external aggregates)
	brand, err := s.brandRepo.FindByID(ctx, cmd.BrandID)
	if err != nil {
		return nil, fmt.Errorf("brand not found: %w", err)
	}
	if brand == nil {
		return nil, fmt.Errorf("brand with ID %s does not exist", cmd.BrandID)
	}

	for _, groupID := range cmd.GroupIDs {
		group, err := s.groupRepo.FindByID(ctx, groupID)
		if err != nil {
			return nil, fmt.Errorf("product group not found: %w", err)
		}
		if group == nil {
			return nil, fmt.Errorf("product group with ID %s does not exist", groupID)
		}
	}

	// 2. Create Product domain entity
	product, err := domain.NewProduct(
		cmd.SKU,
		cmd.Name,
		cmd.LongName,
		cmd.Description,
		cmd.ProductType,
		cmd.BrandID,
		cmd.Barcode,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create product domain entity: %w", err)
	}

	// Add groups to product
	for _, groupID := range cmd.GroupIDs {
		product.AddGroup(groupID)
	}

	// 3. Persist Product
	if err := s.productRepo.Save(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to save product: %w", err)
	}

	return NewProductDTOFromDomain(product), nil
}

// AddGroupToProduct handles the UC-P-XXX: Add a group to an existing Product.
func (s *ProductService) AddGroupToProduct(ctx context.Context, cmd AddGroupCommand) (*ProductDTO, error) {
	product, err := s.productRepo.FindByID(ctx, cmd.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	if product == nil {
		return nil, fmt.Errorf("product with ID %s does not exist", cmd.ProductID)
	}

	group, err := s.groupRepo.FindByID(ctx, cmd.GroupID)
	if err != nil {
		return nil, fmt.Errorf("product group not found: %w", err)
	}
	if group == nil {
		return nil, fmt.Errorf("product group with ID %s does not exist", cmd.GroupID)
	}

	product.AddGroup(cmd.GroupID)

	if err := s.productRepo.Save(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to add group to product: %w", err)
	}

	return NewProductDTOFromDomain(product), nil
}

// AddDirectAttributeToProduct handles the UC-P-XXX: Add a direct attribute to an existing Product.
func (s *ProductService) AddDirectAttributeToProduct(ctx context.Context, cmd AddDirectAttributeCommand) (*ProductDTO, error) {
	product, err := s.productRepo.FindByID(ctx, cmd.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	if product == nil {
		return nil, fmt.Errorf("product with ID %s does not exist", cmd.ProductID)
	}

	// Validate attribute exists using attributeRepo
	attribute, err := s.attributeRepo.FindByID(ctx, cmd.AttributeID)
	if err != nil {
		return nil, fmt.Errorf("attribute not found: %w", err)
	}
	if attribute == nil {
		return nil, fmt.Errorf("attribute with ID %s does not exist", cmd.AttributeID)
	}

	product.AddDirectAttribute(cmd.AttributeID)

	if err := s.productRepo.Save(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to add direct attribute to product: %w", err)
	}

	return NewProductDTOFromDomain(product), nil
}

// UpdateProductSKU handles the UC-P-006: Modificar SKU de Producto.
func (s *ProductService) UpdateProductSKU(ctx context.Context, cmd UpdateProductSKUCommand) (*ProductDTO, error) {
	product, err := s.productRepo.FindByID(ctx, cmd.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	if product == nil {
		return nil, fmt.Errorf("product with ID %s does not exist", cmd.ProductID)
	}

	// Update product's own SKU
	product.SKU = cmd.NewSKU

	// This method in the repository is expected to handle the cascade update of ProductVariant SKUs
	if err := s.productRepo.UpdateSKUs(ctx, product.ID, cmd.NewSKU); err != nil {
		return nil, fmt.Errorf("failed to update product SKU and cascade to variants: %w", err)
	}

	// The product returned from repo.FindByID will have the old SKU. We need to fetch the updated product
	// or ensure the Save method updates the entity in memory
	updatedProduct, err := s.productRepo.FindByID(ctx, cmd.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated product: %w", err)
	}
	if updatedProduct == nil {
		return nil, fmt.Errorf("updated product with ID %s does not exist", cmd.ProductID)
	}

	return NewProductDTOFromDomain(updatedProduct), nil
}

// CreateAttribute handles the UC-P-001: Create an Attribute.
func (s *ProductService) CreateAttribute(ctx context.Context, cmd CreateAttributeCommand) (*AttributeDTO, error) {
	// 1. Validate ScopeBrandID and ScopeGroupID if provided
	if cmd.ScopeBrandID != nil {
		brand, err := s.brandRepo.FindByID(ctx, *cmd.ScopeBrandID)
		if err != nil {
			return nil, fmt.Errorf("scope brand not found: %w", err)
		}
		if brand == nil {
			return nil, fmt.Errorf("brand with ID %s does not exist", *cmd.ScopeBrandID)
		}
	}

	if cmd.ScopeGroupID != nil {
		group, err := s.groupRepo.FindByID(ctx, *cmd.ScopeGroupID)
		if err != nil {
			return nil, fmt.Errorf("scope product group not found: %w", err)
		}
		if group == nil {
			return nil, fmt.Errorf("product group with ID %s does not exist", *cmd.ScopeGroupID)
		}
	}

	// 2. Create Attribute domain entity
	attribute, err := domain.NewAttribute(
		cmd.Name,
		cmd.Code,
		cmd.SortOrder,
		cmd.ScopeBrandID,
		cmd.ScopeGroupID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create attribute domain entity: %w", err)
	}

	// 3. Add values to the attribute
	for _, valCmd := range cmd.Values {
		_, err := attribute.AddValue(valCmd.Value, valCmd.Code)
		if err != nil {
			return nil, fmt.Errorf("failed to add attribute value: %w", err)
		}
	}

	// 4. Persist Attribute
	if err := s.attributeRepo.Save(ctx, attribute); err != nil {
		return nil, fmt.Errorf("failed to save attribute: %w", err)
	}

	return NewAttributeDTOFromDomain(attribute), nil
}

// UpdateAttribute handles the UC-P-002: Modify an Attribute.
func (s *ProductService) UpdateAttribute(ctx context.Context, cmd UpdateAttributeCommand) (*AttributeDTO, error) {
	attribute, err := s.attributeRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("attribute not found: %w", err)
	}
	if attribute == nil {
		return nil, fmt.Errorf("attribute with ID %s does not exist", cmd.ID)
	}

	// 1. Validate ScopeBrandID and ScopeGroupID if provided and changed
	if cmd.ScopeBrandID != nil && (attribute.ScopeBrandID == nil || *attribute.ScopeBrandID != *cmd.ScopeBrandID) {
		brand, err := s.brandRepo.FindByID(ctx, *cmd.ScopeBrandID)
		if err != nil {
			return nil, fmt.Errorf("scope brand not found: %w", err)
		}
		if brand == nil {
			return nil, fmt.Errorf("brand with ID %s does not exist", *cmd.ScopeBrandID)
		}
	}

	if cmd.ScopeGroupID != nil && (attribute.ScopeGroupID == nil || *attribute.ScopeGroupID != *cmd.ScopeGroupID) {
		group, err := s.groupRepo.FindByID(ctx, *cmd.ScopeGroupID)
		if err != nil {
			return nil, fmt.Errorf("scope product group not found: %w", err)
		}
		if group == nil {
			return nil, fmt.Errorf("product group with ID %s does not exist", *cmd.ScopeGroupID)
		}
	}

	// 2. Update top-level attribute fields
	if cmd.Name != nil {
		attribute.Name = *cmd.Name
	}
	if cmd.Code != nil {
		attribute.Code = *cmd.Code
	}
	if cmd.SortOrder != nil {
		attribute.SortOrder = *cmd.SortOrder
	}
	if cmd.ScopeBrandID != nil {
		attribute.ScopeBrandID = cmd.ScopeBrandID
	}
	if cmd.ScopeGroupID != nil {
		attribute.ScopeGroupID = cmd.ScopeGroupID
	}

	// 3. Handle AttributeValue updates (add, modify, delete)
	// Map existing values for efficient lookup
	existingValuesMap := make(map[uuid.UUID]domain.AttributeValue)
	for _, val := range attribute.Values {
		existingValuesMap[val.ID] = val
	}

	for _, valCmd := range cmd.Values {
		if valCmd.ID == nil { // New value
			_, err := attribute.AddValue(valCmd.Value, valCmd.Code)
			if err != nil {
				return nil, fmt.Errorf("failed to add new attribute value: %w", err)
			}
		} else { // Existing value, potentially updated
			if _, exists := existingValuesMap[*valCmd.ID]; exists {
				err := attribute.UpdateValue(*valCmd.ID, valCmd.Value, valCmd.Code)
				if err != nil {
					return nil, fmt.Errorf("failed to update attribute value %s: %w", valCmd.ID.String(), err)
				}
				delete(existingValuesMap, *valCmd.ID) // Mark as processed
			} else {
				return nil, fmt.Errorf("attribute value with ID %s not found in existing values", valCmd.ID.String())
			}
		}
	}

	// Any values remaining in existingValuesMap should be deleted
	for id := range existingValuesMap {
		err := attribute.RemoveValue(id)
		if err != nil {
			return nil, fmt.Errorf("failed to remove attribute value %s: %w", id.String(), err)
		}
	}

	// 4. Persist updated Attribute
	if err := s.attributeRepo.Save(ctx, attribute); err != nil {
		return nil, fmt.Errorf("failed to save updated attribute: %w", err)
	}

	return NewAttributeDTOFromDomain(attribute), nil
}

// GetApplicableAttributesForProduct handles UC-P-005: Consultar Atributos Aplicables de un Producto.
func (s *ProductService) GetApplicableAttributesForProduct(ctx context.Context, productID uuid.UUID) ([]*AttributeDTO, error) {
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	if product == nil {
		return nil, fmt.Errorf("product with ID %s does not exist", productID)
	}

	// Map to store the winning attribute for each code based on precedence
	// Key: Attribute.Code, Value: *domain.Attribute
	// This map will store the highest precedence attribute for each code encountered.
	finalAttributesMap := make(map[string]*domain.Attribute)

	// Helper function to get precedence score (lower is higher precedence)
	// Precedence order: Direct (0), Group+Brand (1), Group (2), Brand (3), Generic (4)
	getAttributePrecedence := func(attr *domain.Attribute) int {
		// Check for direct attribute - highest precedence
		for _, directAttrID := range product.DirectAttributeIDs {
			if attr.ID == directAttrID {
				return 0 // Direct
			}
		}

		isScopedToProductBrand := attr.ScopeBrandID != nil && *attr.ScopeBrandID == product.BrandID
		isScopedToProductGroup := false
		for _, productGroupID := range product.GroupIDs {
			if attr.ScopeGroupID != nil && *attr.ScopeGroupID == productGroupID {
				isScopedToProductGroup = true
				break
			}
		}

		if isScopedToProductBrand && isScopedToProductGroup {
			return 1 // Group + Brand
		}
		if isScopedToProductGroup {
			return 2 // Product Group
		}
		if isScopedToProductBrand {
			return 3 // Brand
		}

		return 4 // Generic (lowest precedence)
	}

	// Fetch all candidate attributes and score them
	var scoredAttributes []struct {
		Attribute  *domain.Attribute
		Precedence int
	}

	// 1. Fetch Generic attributes
	genericAttrs, err := s.attributeRepo.FindByScope(ctx, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch generic attributes: %w", err)
	}
	for _, attr := range genericAttrs {
		scoredAttributes = append(scoredAttributes, struct {
			Attribute  *domain.Attribute
			Precedence int
		}{Attribute: attr, Precedence: getAttributePrecedence(attr)})
	}

	// 2. Fetch Brand scoped attributes
	brandScopedAttrs, err := s.attributeRepo.FindByScope(ctx, &product.BrandID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch brand scoped attributes: %w", err)
	}
	for _, attr := range brandScopedAttrs {
		scoredAttributes = append(scoredAttributes, struct {
			Attribute  *domain.Attribute
			Precedence int
		}{Attribute: attr, Precedence: getAttributePrecedence(attr)})
	}

	// 3. Fetch Group scoped attributes
	for _, groupID := range product.GroupIDs {
		groupScopedAttrs, err := s.attributeRepo.FindByScope(ctx, nil, &groupID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch group %s scoped attributes: %w", groupID.String(), err)
		}
		for _, attr := range groupScopedAttrs {
			scoredAttributes = append(scoredAttributes, struct {
				Attribute  *domain.Attribute
				Precedence int
			}{Attribute: attr, Precedence: getAttributePrecedence(attr)})
		}
	}

	// 4. Fetch Group + Brand scoped attributes (if any, using direct calls to repo)
	for _, groupID := range product.GroupIDs {
		groupBrandScopedAttrs, err := s.attributeRepo.FindByScope(ctx, &product.BrandID, &groupID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch group %s and brand %s scoped attributes: %w", groupID.String(), product.BrandID.String(), err)
		}
		for _, attr := range groupBrandScopedAttrs {
			scoredAttributes = append(scoredAttributes, struct {
				Attribute  *domain.Attribute
				Precedence int
			}{Attribute: attr, Precedence: getAttributePrecedence(attr)})
		}
	}

	// 5. Fetch Direct attributes
	for _, attrID := range product.DirectAttributeIDs {
		directAttr, err := s.attributeRepo.FindByID(ctx, attrID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch direct attribute %s: %w", attrID.String(), err)
		}
		if directAttr != nil {
			scoredAttributes = append(scoredAttributes, struct {
				Attribute  *domain.Attribute
				Precedence int
			}{Attribute: directAttr, Precedence: getAttributePrecedence(directAttr)})
		}
	}

	// Sort attributes: first by precedence (lower is better), then by SortOrder
	sort.Slice(scoredAttributes, func(i, j int) bool {
		if scoredAttributes[i].Precedence != scoredAttributes[j].Precedence {
			return scoredAttributes[i].Precedence < scoredAttributes[j].Precedence
		}
		return scoredAttributes[i].Attribute.SortOrder < scoredAttributes[j].Attribute.SortOrder
	})

	// Apply precedence: populate final map, only adding if no attribute with that code exists
	// or if the new attribute has higher precedence.
	// Since we've sorted by precedence, simply iterating will ensure the first one for each code wins.
	for _, sa := range scoredAttributes {
		if existing, exists := finalAttributesMap[sa.Attribute.Code]; !exists || getAttributePrecedence(sa.Attribute) < getAttributePrecedence(existing) {
			finalAttributesMap[sa.Attribute.Code] = sa.Attribute
		}
	}

	// Convert map values to slice of DTOs, sorted by SortOrder
	var result []*AttributeDTO
	for _, attr := range finalAttributesMap {
		result = append(result, NewAttributeDTOFromDomain(attr))
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].SortOrder < result[j].SortOrder
	})

	return result, nil
}

// GetAttributeByID handles fetching a single attribute by its ID.
func (s *ProductService) GetAttributeByID(ctx context.Context, query GetAttributeByIDQuery) (*AttributeDTO, error) {
	attribute, err := s.attributeRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("attribute not found: %w", err)
	}
	if attribute == nil {
		return nil, fmt.Errorf("attribute with ID %s does not exist", query.ID)
	}
	return NewAttributeDTOFromDomain(attribute), nil
}

// ListAttributes handles fetching a list of attributes with optional filtering.
func (s *ProductService) ListAttributes(ctx context.Context, query ListAttributesQuery) ([]*AttributeDTO, error) {
	// Implement filtering logic based on query.ScopeType, query.BrandID, query.ProductGroupID
	// For simplicity, let's assume FindByScope can handle all combinations for now.
	// A more robust implementation would build the query dynamically.

	attributes, err := s.attributeRepo.FindByScope(ctx, query.BrandID, query.ProductGroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to list attributes: %w", err)
	}

	var dtos []*AttributeDTO
	for _, attr := range attributes {
		dtos = append(dtos, NewAttributeDTOFromDomain(attr))
	}
	// TODO: Further filtering by ScopeType if needed, as FindByScope might not cover all cases
	return dtos, nil
}

// GetProductByID handles fetching a single product by its ID.
func (s *ProductService) GetProductByID(ctx context.Context, query GetProductByIDQuery) (*ProductDTO, error) {
	product, err := s.productRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	if product == nil {
		return nil, fmt.Errorf("product with ID %s does not exist", query.ID)
	}
	return NewProductDTOFromDomain(product), nil
}

// ListProducts handles fetching a list of products with optional filtering.
func (s *ProductService) ListProducts(ctx context.Context, query ListProductsQuery) ([]*ProductDTO, error) {
	// TODO: Implement actual listing with filters based on query parameters.
	// For now, returning an empty list as a placeholder.
	return []*ProductDTO{}, nil
}

// ListProductVariantsByProductID handles fetching a list of product variants for a given product ID.
func (s *ProductService) ListProductVariantsByProductID(ctx context.Context, query ListProductVariantsByProductIDQuery) ([]*ProductVariantDTO, error) {
	// TODO: Implement domain.ProductVariantRepository.FindByProductID
	// For now, returning an empty list as a placeholder.
	return []*ProductVariantDTO{}, nil
}

// GetProductVariantByID handles fetching a single product variant by its ID.
func (s *ProductService) GetProductVariantByID(ctx context.Context, query GetProductVariantByIDQuery) (*ProductVariantDTO, error) {
	variant, err := s.variantRepo.FindByID(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("product variant not found: %w", err)
	}
	if variant == nil {
		return nil, fmt.Errorf("product variant with ID %s does not exist", query.ID)
	}
	// Need to fetch all attributes to populate OptionConfiguration in ProductVariantDTO
	allAttributes, err := s.attributeRepo.FindByScope(ctx, nil, nil) // Fetch all attributes for now
	if err != nil {
		return nil, fmt.Errorf("failed to fetch attributes for variant: %w", err)
	}
	return NewProductVariantDTOFromDomain(variant, allAttributes), nil
}

// GetProductVariantBySKU handles fetching a single product variant by its SKU.
func (s *ProductService) GetProductVariantBySKU(ctx context.Context, query GetProductVariantBySKUQuery) (*ProductVariantDTO, error) {
	variant, err := s.variantRepo.FindBySKU(ctx, query.SKU)
	if err != nil {
		return nil, fmt.Errorf("product variant not found: %w", err)
	}
	if variant == nil {
		return nil, fmt.Errorf("product variant with SKU %s does not exist", query.SKU)
	}
	allAttributes, err := s.attributeRepo.FindByScope(ctx, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch attributes for variant: %w", err)
	}
	return NewProductVariantDTOFromDomain(variant, allAttributes), nil
}

// GenerateProductVariants handles UC-P-007: Pre-generate Variants for a Product.
func (s *ProductService) GenerateProductVariants(ctx context.Context, cmd GenerateProductVariantsCommand) error {
	product, err := s.productRepo.FindByID(ctx, cmd.ProductID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}
	if product == nil {
		return fmt.Errorf("product with ID %s does not exist", cmd.ProductID)
	}

	// UC-P-005: Get Applicable Attributes
	applicableAttributes, err := s.GetApplicableAttributesForProduct(ctx, product.ID)
	if err != nil {
		return fmt.Errorf("failed to get applicable attributes for product %s: %w", product.ID, err)
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
			return fmt.Errorf("failed to retrieve full domain attribute %s: %w", attrDTO.ID, err)
		}
		domainAttributes = append(domainAttributes, fullDomainAttribute)
	}

	// Logic to generate all possible combinations of variants based on applicable attributes
	// This is a complex combinatorial problem and requires careful implementation.
	// Placeholder: This part will iterate through all combinations and create/update variants.

	// For each combination:
	//   - Construct expected SKU
	//   - Check if variant already exists by SKU or by ProductID + AttributeValues
	//   - If not exists, create with Status = CONFIRMED (as it's pre-generated)
	//   - If exists and PROVISIONAL, update Status to CONFIRMED

	return nil // Placeholder for actual implementation
}

// FindOrCreateProductVariant handles UC-P-009: Obtener o Crear Variante (Find or Create).
func (s *ProductService) FindOrCreateProductVariant(ctx context.Context, cmd FindOrCreateProductVariantCommand) (*ProductVariantDTO, error) {
	product, err := s.productRepo.FindByID(ctx, cmd.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	if product == nil {
		return nil, fmt.Errorf("product with ID %s does not exist", cmd.ProductID)
	}

	// Get applicable attributes for validation and SKU construction
	applicableAttributesDTOs, err := s.GetApplicableAttributesForProduct(ctx, product.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get applicable attributes for product %s: %w", product.ID, err)
	}

	// Map AttributeCode to domain.Attribute and AttributeValue
	attrCodeToAttribute := make(map[string]*domain.Attribute)
	attrValueToDomainValue := make(map[string]domain.AttributeValue) // Value (string) -> domain.AttributeValue
	for _, attrDTO := range applicableAttributesDTOs {
		fullDomainAttribute, err := s.attributeRepo.FindByID(ctx, attrDTO.ID)
		if err != nil || fullDomainAttribute == nil {
			return nil, fmt.Errorf("failed to retrieve full domain attribute %s: %w", attrDTO.ID, err)
		}
		attrCodeToAttribute[fullDomainAttribute.Code] = fullDomainAttribute
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

	for _, item := range cmd.OptionConfiguration {
		attr, attrExists := attrCodeToAttribute[item.AttributeName]
		if !attrExists {
			return nil, fmt.Errorf("attribute '%s' is not applicable to product '%s'", item.AttributeName, product.ID)
		}
		val, valExists := attrValueToDomainValue[item.Value]
		if !valExists || val.AttributeID != attr.ID { // Ensure value belongs to this attribute
			return nil, fmt.Errorf("value '%s' is not valid for attribute '%s'", item.Value, item.AttributeName)
		}
		attributeValueIDs = append(attributeValueIDs, val.ID)
		sortedAttributeCodes = append(sortedAttributeCodes, attr.Code)
		attrCodeToValueCode[attr.Code] = val.Code
	}

	// Sort attribute codes to ensure deterministic SKU construction
	sort.Slice(sortedAttributeCodes, func(i, j int) bool {
		attrI := attrCodeToAttribute[sortedAttributeCodes[i]]
		attrJ := attrCodeToAttribute[sortedAttributeCodes[j]]
		return attrI.SortOrder < attrJ.SortOrder
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
		// Ensure SKU matches; if not, update. This handles potential SKU changes on product or attributes.
		if variant.SKU != generatedSKU {
			variant.SKU = generatedSKU
			if err := s.variantRepo.Save(ctx, variant); err != nil { // Save might update
				return nil, fmt.Errorf("failed to update SKU of existing variant %s: %w", variant.ID, err)
			}
		}
		return NewProductVariantDTOFromDomain(variant, applicableAttributes), nil
	}

	// If not found, create new variant
	newVariant, err := domain.NewProductVariant(product.ID, generatedSKU, nil, domain.StatusProvisional, attributeValueIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to create new product variant domain entity: %w", err)
	}

	if err := s.variantRepo.Save(ctx, newVariant); err != nil {
		return nil, fmt.Errorf("failed to save new product variant: %w", err)
	}

	return NewProductVariantDTOFromDomain(newVariant, applicableAttributes), nil
}

// UpdateProductVariant handles UC-P-008: Modificar una Variante Específica.
func (s *ProductService) UpdateProductVariant(ctx context.Context, cmd UpdateProductVariantCommand) (*ProductVariantDTO, error) {
	variant, err := s.variantRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("product variant not found: %w", err)
	}
	if variant == nil {
		return nil, fmt.Errorf("product variant with ID %s does not exist", cmd.ID)
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
		return nil, fmt.Errorf("failed to save updated product variant: %w", err)
	}

	// Need to fetch all attributes to populate OptionConfiguration in ProductVariantDTO
	allAttributes, err := s.attributeRepo.FindByScope(ctx, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch attributes for variant: %w", err)
	}

	return NewProductVariantDTOFromDomain(variant, allAttributes), nil
}

// GetPartyServiceConfigurationByID handles fetching a single PartyServiceConfiguration by its ID and PartyID.
func (s *ProductService) GetPartyServiceConfigurationByID(ctx context.Context, query GetPartyServiceConfigurationByIDQuery) (*PartyServiceConfigurationDTO, error) {
	config, err := s.partyServiceConfigRepo.FindByID(ctx, query.PartyID, query.ID)
	if err != nil {
		return nil, fmt.Errorf("party service configuration not found: %w", err)
	}
	if config == nil {
		return nil, fmt.Errorf("party service configuration with ID %s for party %s does not exist", query.ID, query.PartyID)
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
		return nil, fmt.Errorf("failed to list party service configurations: %w", err)
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
	var configDetails json.RawMessage
	if cmd.ConfigurationDetails != nil {
		payload, err := json.Marshal(cmd.ConfigurationDetails)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize configuration details: %w", err)
		}
		configDetails = payload
	}

	config, err := domain.NewPartyServiceConfiguration(cmd.PartyID, cmd.ServiceID, cmd.Name, configDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to create party service configuration domain entity: %w", err)
	}

	if err := s.partyServiceConfigRepo.Save(ctx, config); err != nil {
		return nil, fmt.Errorf("failed to save party service configuration: %w", err)
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
	config, err := s.partyServiceConfigRepo.FindByID(ctx, cmd.PartyID, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("party service configuration not found: %w", err)
	}
	if config == nil {
		return nil, fmt.Errorf("party service configuration with ID %s for party %s does not exist", cmd.ID, cmd.PartyID)
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
			return nil, fmt.Errorf("failed to serialize configuration details: %w", err)
		}
		configDetails = payload
	}

	if err := config.Update(serviceID, name, configDetails); err != nil {
		return nil, fmt.Errorf("failed to update party service configuration domain entity: %w", err)
	}

	if err := s.partyServiceConfigRepo.Save(ctx, config); err != nil {
		return nil, fmt.Errorf("failed to save updated party service configuration: %w", err)
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
	err := s.partyServiceConfigRepo.Delete(ctx, cmd.PartyID, cmd.ID)
	if err != nil {
		return fmt.Errorf("failed to delete party service configuration: %w", err)
	}
	return nil
}
