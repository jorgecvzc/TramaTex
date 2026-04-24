package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/joran-cortez/tramatex/internal/party/domain"
	"gorm.io/gorm"
)

// GORMPartyRepository implements PartyRepository using GORM.
type GORMPartyRepository struct {
	db *gorm.DB
}

func NewGORMPartyRepository(db *gorm.DB) *GORMPartyRepository {
	return &GORMPartyRepository{db: db}
}

func (r *GORMPartyRepository) Save(ctx context.Context, party *domain.Party, createdBy string, modifiedBy string) error {
	if party == nil {
		return domain.NewValidationError("party cannot be nil")
	}

	now := time.Now()
	partyModel := partyDataModelFromDomain(party)

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing PartyDataModel
		result := tx.Select("id", "created_at", "created_by").First(&existing, "id = ?", partyModel.ID)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domain.WrapPersistence("failed to load party", result.Error)
		}

		partyModel.ModifiedAt = now
		partyModel.ModifiedBy = modifiedBy

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			partyModel.CreatedAt = now
			partyModel.CreatedBy = createdBy
			if err := tx.Create(partyModel).Error; err != nil {
				return domain.WrapPersistence("failed to create party", err)
			}
		} else {
			// Update using map to ensure 0 values (like discount) are not ignored by GORM
			updates := map[string]interface{}{
				"status":                      partyModel.Status,
				"default_discount_percentage": partyModel.DefaultDiscountPercentage,
				"modified_at":                 partyModel.ModifiedAt,
				"modified_by":                 partyModel.ModifiedBy,
			}
			if err := tx.Model(&PartyDataModel{}).Where("id = ?", partyModel.ID).Updates(updates).Error; err != nil {
				return domain.WrapPersistence("failed to update party", err)
			}
		}

		if party.PersonProfile() != nil {
			profile := party.PersonProfile()
			profileModel := &PersonProfileDataModel{
				PartyID:   partyModel.ID,
				FirstName: profile.FirstName(),
				LastName:  profile.LastName(),
			}
			if profile.Phone() != nil {
				phone := profile.Phone().Value()
				profileModel.Phone = &phone
			}
			if profile.Email() != nil {
				email := profile.Email().Value()
				profileModel.Email = &email
			}
			if err := tx.Save(profileModel).Error; err != nil {
				return domain.WrapPersistence("failed to save person profile", err)
			}
		} else {
			if err := tx.Where("party_id = ?", partyModel.ID).Delete(&PersonProfileDataModel{}).Error; err != nil {
				return domain.WrapPersistence("failed to delete person profile", err)
			}
		}

		if party.OrganizationProfile() != nil {
			org := party.OrganizationProfile()
			var taxID *string
			var taxType *string
			if org.TaxID() != nil {
				value := org.TaxID().Value()
				typeValue := org.TaxID().Type()
				taxID = &value
				taxType = &typeValue
			}
			orgModel := &OrganizationProfileDataModel{
				PartyID:   partyModel.ID,
				Name:      org.Name(),
				TaxID:     taxID,
				TaxIDType: taxType,
				Website:   org.Website(),
				Notes:     org.Notes(),
			}
			if org.Phone() != nil {
				phone := org.Phone().Value()
				orgModel.Phone = &phone
			}
			if org.Email() != nil {
				email := org.Email().Value()
				orgModel.Email = &email
			}
			if err := tx.Save(orgModel).Error; err != nil {
				return domain.WrapPersistence("failed to save organization profile", err)
			}

			if err := tx.Where("organization_party_id = ?", partyModel.ID).Delete(&ContactDetailsDataModel{}).Error; err != nil {
				return domain.WrapPersistence("failed to delete contact details", err)
			}

			for _, contact := range org.Contacts() {
				contactModel := &ContactDetailsDataModel{
					ID:                contact.ID().Value(),
					OrganizationParty: partyModel.ID,
					TypeDescription:   contact.TypeDescription(),
				}
				if contact.Phone() != nil {
					phone := contact.Phone().Value()
					contactModel.Phone = &phone
				}
				if contact.Email() != nil {
					email := contact.Email().Value()
					contactModel.Email = &email
				}
				if contact.RelatedPartyID() != nil {
					related := contact.RelatedPartyID().Value()
					contactModel.RelatedPartyID = &related
				}
				if err := tx.Save(contactModel).Error; err != nil {
					return domain.WrapPersistence("failed to save contact details", err)
				}
			}
		} else {
			if err := tx.Where("organization_party_id = ?", partyModel.ID).Delete(&ContactDetailsDataModel{}).Error; err != nil {
				return domain.WrapPersistence("failed to delete contact details", err)
			}
			if err := tx.Where("party_id = ?", partyModel.ID).Delete(&OrganizationProfileDataModel{}).Error; err != nil {
				return domain.WrapPersistence("failed to delete organization profile", err)
			}
		}

		if err := tx.Where("party_id = ?", partyModel.ID).Delete(&PartyRoleDataModel{}).Error; err != nil {
			return domain.WrapPersistence("failed to delete party roles", err)
		}
		for _, role := range party.Roles() {
			roleModel := &PartyRoleDataModel{
				PartyID: partyModel.ID,
				Role:    string(role.Type()),
			}
			if role.CreationIdentifier() != nil {
				roleModel.CreationIdentifier = role.CreationIdentifier()
			}
			if err := tx.Create(roleModel).Error; err != nil {
				return domain.WrapPersistence("failed to save party role", err)
			}
		}

		return nil
	})
}

func (r *GORMPartyRepository) FindByID(ctx context.Context, id domain.PartyID) (*domain.Party, error) {
	var partyModel PartyDataModel
	if err := r.db.WithContext(ctx).First(&partyModel, "id = ?", id.Value()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.NewNotFoundError("party not found")
		}
		return nil, domain.WrapPersistence("failed to query party", err)
	}

	personProfile, err := r.loadPersonProfile(ctx, id.Value())
	if err != nil {
		return nil, err
	}
	organizationProfile, err := r.loadOrganizationProfile(ctx, id.Value())
	if err != nil {
		return nil, err
	}
	roles, err := r.loadPartyRoles(ctx, id.Value())
	if err != nil {
		return nil, err
	}

	status := domain.PartyStatus(partyModel.Status)
	if !status.IsValid() {
		return nil, domain.NewValidationErrorf("invalid party status in storage: %s", partyModel.Status)
	}

	party, err := domain.NewPartyFromPersistence(id, status, personProfile, organizationProfile, roles)
	if err != nil {
		return nil, err
	}

	_ = party.SetDefaultDiscountPercentage(partyModel.DefaultDiscountPercentage)
	party.SetTimestamps(partyModel.CreatedAt, partyModel.ModifiedAt)
	return party, nil
}

func (r *GORMPartyRepository) FindAll(ctx context.Context, filters *PartyFilters) ([]*domain.Party, error) {
	query := r.db.WithContext(ctx).Model(&PartyDataModel{}).Select("parties.id").
		Joins("LEFT JOIN organization_profiles op ON op.party_id = parties.id").
		Joins("LEFT JOIN person_profiles pp ON pp.party_id = parties.id").
		Joins("LEFT JOIN party_roles pr ON pr.party_id = parties.id")

	if filters != nil {
		if filters.Status != nil {
			query = query.Where("parties.status = ?", string(*filters.Status))
		}
		if filters.Role != nil {
			if *filters.Role == domain.PartyRoleContact {
				query = query.Where("pr.role IN ?", []string{string(domain.PartyRoleContact), string(domain.PartyRoleEmployee)})
			} else {
				query = query.Where("pr.role = ?", string(*filters.Role))
			}
		}
		if filters.Type != "" {
			switch filters.Type {
			case "PERSON":
				query = query.Where("pp.party_id IS NOT NULL")
			case "ORGANIZATION":
				query = query.Where("op.party_id IS NOT NULL")
			case "BOTH":
				query = query.Where("pp.party_id IS NOT NULL AND op.party_id IS NOT NULL")
			}
		}
		if filters.Name != "" {
			query = query.Where("op.name ILIKE ? OR (pp.first_name || ' ' || pp.last_name) ILIKE ?", "%"+filters.Name+"%", "%"+filters.Name+"%")
		}
		if filters.TaxID != "" {
			query = query.Where("op.tax_id ILIKE ?", "%"+filters.TaxID+"%")
		}
		if filters.PageSize > 0 {
			query = query.Limit(filters.PageSize)
			offset := (filters.PageNumber - 1) * filters.PageSize
			query = query.Offset(offset)
		}
	}

	var ids []string
	if err := query.Order("parties.id").Distinct().Pluck("parties.id", &ids).Error; err != nil {
		return nil, domain.WrapPersistence("failed to query parties", err)
	}

	parties := make([]*domain.Party, 0, len(ids))
	for _, id := range ids {
		partyID, err := domain.NewPartyID(id)
		if err != nil {
			return nil, domain.WrapValidation("invalid party ID", err)
		}
		party, err := r.FindByID(ctx, partyID)
		if err != nil {
			return nil, err
		}
		parties = append(parties, party)
	}

	return parties, nil
}

func (r *GORMPartyRepository) Delete(ctx context.Context, id domain.PartyID) error {
	if err := r.db.WithContext(ctx).Delete(&PartyDataModel{}, "id = ?", id.Value()).Error; err != nil {
		return domain.WrapPersistence("failed to delete party", err)
	}
	return nil
}

func (r *GORMPartyRepository) Exists(ctx context.Context, id domain.PartyID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&PartyDataModel{}).Where("id = ?", id.Value()).Count(&count).Error; err != nil {
		return false, domain.WrapPersistence("failed to check party existence", err)
	}
	return count > 0, nil
}

func (r *GORMPartyRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&PartyDataModel{}).Count(&count).Error; err != nil {
		return 0, domain.WrapPersistence("failed to count parties", err)
	}
	return count, nil
}

func (r *GORMPartyRepository) CountByFilters(ctx context.Context, filters *PartyFilters) (int64, error) {
	query := r.db.WithContext(ctx).Model(&PartyDataModel{}).
		Joins("LEFT JOIN organization_profiles op ON op.party_id = parties.id").
		Joins("LEFT JOIN person_profiles pp ON pp.party_id = parties.id").
		Joins("LEFT JOIN party_roles pr ON pr.party_id = parties.id")

	if filters != nil {
		if filters.Status != nil {
			query = query.Where("parties.status = ?", string(*filters.Status))
		}
		if filters.Role != nil {
			if *filters.Role == domain.PartyRoleContact {
				query = query.Where("pr.role IN ?", []string{string(domain.PartyRoleContact), string(domain.PartyRoleEmployee)})
			} else {
				query = query.Where("pr.role = ?", string(*filters.Role))
			}
		}
		if filters.Type != "" {
			switch filters.Type {
			case "PERSON":
				query = query.Where("pp.party_id IS NOT NULL")
			case "ORGANIZATION":
				query = query.Where("op.party_id IS NOT NULL")
			case "BOTH":
				query = query.Where("pp.party_id IS NOT NULL AND op.party_id IS NOT NULL")
			}
		}
		if filters.Name != "" {
			query = query.Where("op.name ILIKE ? OR (pp.first_name || ' ' || pp.last_name) ILIKE ?", "%"+filters.Name+"%", "%"+filters.Name+"%")
		}
		if filters.TaxID != "" {
			query = query.Where("op.tax_id ILIKE ?", "%"+filters.TaxID+"%")
		}
	}

	var count int64
	if err := query.Distinct("parties.id").Count(&count).Error; err != nil {
		return 0, domain.WrapPersistence("failed to count filtered parties", err)
	}

	return count, nil
}

// HasContactDetailsReferences checks if a party is referenced in contact_details.related_party_id
func (r *GORMPartyRepository) HasContactDetailsReferences(ctx context.Context, partyID domain.PartyID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&ContactDetailsDataModel{}).
		Where("related_party_id = ?", partyID.Value()).
		Count(&count).Error; err != nil {
		return false, domain.WrapPersistence("failed to check contact details references", err)
	}
	return count > 0, nil
}

// HasMESWorkReferences checks if a party is referenced in work_orders.party_id
func (r *GORMPartyRepository) HasMESWorkReferences(ctx context.Context, partyID domain.PartyID) (bool, error) {
	var count int64
	// Check if work_orders table exists (it may not in all environments)
	var tableExists bool
	err := r.db.WithContext(ctx).Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_orders')").Scan(&tableExists).Error
	if err != nil {
		return false, domain.WrapPersistence("failed to check work_orders table existence", err)
	}

	if !tableExists {
		return false, nil
	}

	if err := r.db.WithContext(ctx).Table("work_orders").
		Where("party_id = ?", partyID.Value()).
		Count(&count).Error; err != nil {
		return false, domain.WrapPersistence("failed to check MES work references", err)
	}
	return count > 0, nil
}

// HasSalesReferences checks if a party is referenced in sales tables (quotes, orders, invoices, delivery_notes)
func (r *GORMPartyRepository) HasSalesReferences(ctx context.Context, partyID domain.PartyID) (bool, error) {
	partyIDStr := partyID.Value()

	salesTables := []struct {
		name   string
		errMsg string
	}{
		{"quotes", "failed to check quotes references"},
		{"sales_orders", "failed to check sales orders references"},
		{"invoices", "failed to check invoices references"},
		{"delivery_notes", "failed to check delivery notes references"},
	}

	for _, t := range salesTables {
		var tableExists bool
		if err := r.db.WithContext(ctx).Raw(
			"SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = ?)",
			t.name,
		).Scan(&tableExists).Error; err != nil {
			return false, domain.WrapPersistence("failed to check "+t.name+" table existence", err)
		}
		if !tableExists {
			continue
		}

		var count int64
		if err := r.db.WithContext(ctx).Table(t.name).
			Where("party_id::text = ?", partyIDStr).
			Count(&count).Error; err != nil {
			return false, domain.WrapPersistence(t.errMsg, err)
		}
		if count > 0 {
			return true, nil
		}
	}

	return false, nil
}

func (r *GORMPartyRepository) loadPersonProfile(ctx context.Context, partyID string) (*domain.PersonProfile, error) {
	var profile PersonProfileDataModel
	if err := r.db.WithContext(ctx).First(&profile, "party_id = ?", partyID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, domain.WrapPersistence("failed to load person profile", err)
	}

	var phone *domain.Phone
	if profile.Phone != nil && *profile.Phone != "" {
		p, err := domain.NewPhone(*profile.Phone)
		if err != nil {
			return nil, err
		}
		phone = p
	}

	var email *domain.Email
	if profile.Email != nil && *profile.Email != "" {
		e, err := domain.NewEmail(*profile.Email)
		if err != nil {
			return nil, err
		}
		email = e
	}

	return domain.NewPersonProfile(profile.FirstName, profile.LastName, phone, email)
}

func (r *GORMPartyRepository) loadOrganizationProfile(ctx context.Context, partyID string) (*domain.OrganizationProfile, error) {
	var profile OrganizationProfileDataModel
	if err := r.db.WithContext(ctx).First(&profile, "party_id = ?", partyID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, domain.WrapPersistence("failed to load organization profile", err)
	}

	var tax *domain.TaxID
	if profile.TaxID != nil {
		typeValue := "NIF"
		if profile.TaxIDType != nil && *profile.TaxIDType != "" {
			typeValue = *profile.TaxIDType
		}
		created, err := domain.NewTaxID(*profile.TaxID, typeValue)
		if err != nil {
			return nil, err
		}
		tax = created
	}

	var phone *domain.Phone
	if profile.Phone != nil && *profile.Phone != "" {
		p, err := domain.NewPhone(*profile.Phone)
		if err != nil {
			return nil, err
		}
		phone = p
	}

	var email *domain.Email
	if profile.Email != nil && *profile.Email != "" {
		e, err := domain.NewEmail(*profile.Email)
		if err != nil {
			return nil, err
		}
		email = e
	}

	org, err := domain.NewOrganizationProfile(profile.Name, tax, profile.Website, phone, email, profile.Notes)
	if err != nil {
		return nil, err
	}

	contacts, err := r.loadContactDetails(ctx, partyID)
	if err != nil {
		return nil, err
	}
	for _, contact := range contacts {
		if err := org.AddContact(contact); err != nil {
			return nil, err
		}
	}

	return org, nil
}

func (r *GORMPartyRepository) loadPartyRoles(ctx context.Context, partyID string) ([]domain.PartyRole, error) {
	var roleModels []PartyRoleDataModel
	if err := r.db.WithContext(ctx).Where("party_id = ?", partyID).Find(&roleModels).Error; err != nil {
		return nil, domain.WrapPersistence("failed to load party roles", err)
	}

	roles := make([]domain.PartyRole, 0, len(roleModels))
	for _, role := range roleModels {
		roleType := domain.PartyRoleType(role.Role)
		if !roleType.IsValid() {
			return nil, domain.NewValidationErrorf("invalid role in storage: %s", role.Role)
		}
		partyRole, err := domain.NewPartyRole(roleType, role.CreationIdentifier)
		if err != nil {
			return nil, err
		}
		roles = append(roles, partyRole)
	}

	return roles, nil
}

func (r *GORMPartyRepository) loadContactDetails(ctx context.Context, partyID string) ([]*domain.ContactDetails, error) {
	var contactModels []ContactDetailsDataModel
	if err := r.db.WithContext(ctx).Where("organization_party_id = ?", partyID).Find(&contactModels).Error; err != nil {
		return nil, domain.WrapPersistence("failed to load contact details", err)
	}

	contacts := make([]*domain.ContactDetails, 0, len(contactModels))
	for _, model := range contactModels {
		contactID, err := domain.NewContactDetailsID(model.ID)
		if err != nil {
			return nil, err
		}

		var phone *domain.Phone
		if model.Phone != nil {
			phone, err = domain.NewPhone(*model.Phone)
			if err != nil {
				return nil, err
			}
		}

		var email *domain.Email
		if model.Email != nil {
			email, err = domain.NewEmail(*model.Email)
			if err != nil {
				return nil, err
			}
		}

		var related *domain.PartyID
		if model.RelatedPartyID != nil {
			pid, err := domain.NewPartyID(*model.RelatedPartyID)
			if err != nil {
				return nil, err
			}
			related = &pid
		}

		contact, err := domain.NewContactDetails(contactID, model.TypeDescription, phone, email, related)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	return contacts, nil
}

// GORMPartyRelationshipRepository implements PartyRelationshipRepository using GORM.
type GORMPartyRelationshipRepository struct {
	db *gorm.DB
}

func NewGORMPartyRelationshipRepository(db *gorm.DB) *GORMPartyRelationshipRepository {
	return &GORMPartyRelationshipRepository{db: db}
}

func (r *GORMPartyRelationshipRepository) Save(ctx context.Context, relationship domain.PartyRelationship, createdBy string, modifiedBy string) error {
	now := time.Now()
	model := &PartyRelationshipDataModel{
		ID:         relationship.ID().Value(),
		FromParty:  relationship.FromID().Value(),
		ToParty:    relationship.ToID().Value(),
		Type:       string(relationship.Type()),
		CreatedBy:  createdBy,
		CreatedAt:  now,
		ModifiedBy: modifiedBy,
		ModifiedAt: now,
	}

	var existing PartyRelationshipDataModel
	result := r.db.WithContext(ctx).Select("id", "created_at", "created_by").First(&existing, "id = ?", model.ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return domain.WrapPersistence("failed to load relationship", result.Error)
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		model.CreatedAt = now
		model.CreatedBy = createdBy
	} else {
		model.CreatedAt = existing.CreatedAt
		model.CreatedBy = existing.CreatedBy
	}

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return domain.WrapPersistence("failed to save relationship", err)
	}
	return nil
}

func (r *GORMPartyRelationshipRepository) FindByPartyID(ctx context.Context, partyID domain.PartyID) ([]domain.PartyRelationship, error) {
	var models []PartyRelationshipDataModel
	if err := r.db.WithContext(ctx).Where("from_party_id = ? OR to_party_id = ?", partyID.Value(), partyID.Value()).Find(&models).Error; err != nil {
		return nil, domain.WrapPersistence("failed to load relationships", err)
	}

	relationships := make([]domain.PartyRelationship, 0, len(models))
	for _, model := range models {
		relID, err := domain.NewPartyRelationshipID(model.ID)
		if err != nil {
			return nil, err
		}
		fromID, err := domain.NewPartyID(model.FromParty)
		if err != nil {
			return nil, err
		}
		toID, err := domain.NewPartyID(model.ToParty)
		if err != nil {
			return nil, err
		}

		relType := domain.RelationshipType(model.Type)
		if !relType.IsValid() {
			return nil, domain.NewValidationErrorf("invalid relationship type in storage: %s", model.Type)
		}

		rel, err := domain.NewPartyRelationship(relID, fromID, toID, relType)
		if err != nil {
			return nil, err
		}
		relationships = append(relationships, rel)
	}

	return relationships, nil
}

func (r *GORMPartyRelationshipRepository) Delete(ctx context.Context, id domain.PartyRelationshipID) error {
	if err := r.db.WithContext(ctx).Delete(&PartyRelationshipDataModel{}, "id = ?", id.Value()).Error; err != nil {
		return domain.WrapPersistence("failed to delete relationship", err)
	}
	return nil
}

// GORMPartyAddressRepository implements PartyAddressRepository using GORM.
type GORMPartyAddressRepository struct {
	db *gorm.DB
}

func NewGORMPartyAddressRepository(db *gorm.DB) *GORMPartyAddressRepository {
	return &GORMPartyAddressRepository{db: db}
}

func (r *GORMPartyAddressRepository) Save(ctx context.Context, address *domain.Address, addressID domain.AddressID, partyID domain.PartyID, createdBy string, modifiedBy string) error {
	if address == nil {
		return domain.NewValidationError("address cannot be nil")
	}

	now := time.Now()
	model := &PartyAddressDataModel{
		ID:         addressID.Value(),
		PartyID:    partyID.Value(),
		Street:     address.Street(),
		City:       address.City(),
		Province:   address.Province(),
		PostalCode: address.PostalCode(),
		Country:    address.Country(),
		IsPrimary:  false,
		CreatedBy:  createdBy,
		CreatedAt:  now,
		ModifiedBy: modifiedBy,
		ModifiedAt: now,
	}

	var existing PartyAddressDataModel
	result := r.db.WithContext(ctx).Select("id", "created_at", "created_by").First(&existing, "id = ?", model.ID)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return domain.WrapPersistence("failed to load address", result.Error)
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		model.CreatedAt = now
		model.CreatedBy = createdBy
	} else {
		model.CreatedAt = existing.CreatedAt
		model.CreatedBy = existing.CreatedBy
	}

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return domain.WrapPersistence("failed to save address", err)
	}
	return nil
}

func (r *GORMPartyAddressRepository) FindByPartyID(ctx context.Context, partyID domain.PartyID) ([]*AddressWithID, error) {
	var models []PartyAddressDataModel
	if err := r.db.WithContext(ctx).Where("party_id = ?", partyID.Value()).Order("created_at").Find(&models).Error; err != nil {
		return nil, domain.WrapPersistence("failed to load addresses", err)
	}

	addresses := make([]*AddressWithID, 0, len(models))
	for _, model := range models {
		address, err := domain.NewAddress(model.Street, model.City, model.Province, model.PostalCode, model.Country)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, &AddressWithID{
			ID:      model.ID,
			Address: address,
		})
	}

	return addresses, nil
}

func (r *GORMPartyAddressRepository) FindPrimary(ctx context.Context, partyID domain.PartyID) (*domain.Address, error) {
	var model PartyAddressDataModel
	if err := r.db.WithContext(ctx).Where("party_id = ? AND is_primary = true", partyID.Value()).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.NewNotFoundError("primary address not found")
		}
		return nil, domain.WrapPersistence("failed to load primary address", err)
	}

	return domain.NewAddress(model.Street, model.City, model.Province, model.PostalCode, model.Country)
}

func (r *GORMPartyAddressRepository) Delete(ctx context.Context, id domain.AddressID) error {
	if err := r.db.WithContext(ctx).Delete(&PartyAddressDataModel{}, "id = ?", id.Value()).Error; err != nil {
		return domain.WrapPersistence("failed to delete address", err)
	}
	return nil
}
