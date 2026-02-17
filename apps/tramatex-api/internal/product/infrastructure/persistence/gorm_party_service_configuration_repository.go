package persistence

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/joran-cortez/tramatex/internal/product/domain"
)

// GORMPartyServiceConfigurationRepository implements domain.PartyServiceConfigurationRepository using GORM.
type GORMPartyServiceConfigurationRepository struct {
	db *gorm.DB
}

// NewGORMPartyServiceConfigurationRepository creates a new GORMPartyServiceConfigurationRepository.
func NewGORMPartyServiceConfigurationRepository(db *gorm.DB) *GORMPartyServiceConfigurationRepository {
	return &GORMPartyServiceConfigurationRepository{db: db}
}

// Save creates or updates a PartyServiceConfiguration.
func (r *GORMPartyServiceConfigurationRepository) Save(ctx context.Context, config *domain.PartyServiceConfiguration) error {
	model, err := PartyServiceConfigurationFromDomain(config)
	if err != nil {
		return domain.WrapPersistence("failed to convert domain to model", err)
	}
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return domain.WrapPersistence("failed to save party service configuration", err)
	}
	return nil
}

// FindByID finds a PartyServiceConfiguration by its ID and PartyID.
func (r *GORMPartyServiceConfigurationRepository) FindByID(ctx context.Context, partyID, id uuid.UUID) (*domain.PartyServiceConfiguration, error) {
	var model PartyServiceConfigurationModel
	if err := r.db.WithContext(ctx).Where("id = ? AND party_id = ?", id, partyID).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Not found
		}
		return nil, domain.WrapPersistence("failed to find party service configuration by ID", err)
	}
	return model.ToDomain()
}

// FindByPartyID finds all PartyServiceConfigurations for a given PartyID.
func (r *GORMPartyServiceConfigurationRepository) FindByPartyID(ctx context.Context, partyID uuid.UUID) ([]*domain.PartyServiceConfiguration, error) {
	var models []PartyServiceConfigurationModel
	if err := r.db.WithContext(ctx).Where("party_id = ?", partyID).Find(&models).Error; err != nil {
		return nil, domain.WrapPersistence("failed to find party service configurations by party ID", err)
	}

	var configs []*domain.PartyServiceConfiguration
	for _, model := range models {
		config, err := model.ToDomain()
		if err != nil {
			return nil, domain.WrapPersistence("failed to convert model to domain", err)
		}
		configs = append(configs, config)
	}
	return configs, nil
}

// Delete deletes a PartyServiceConfiguration by its ID and PartyID.
func (r *GORMPartyServiceConfigurationRepository) Delete(ctx context.Context, partyID, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ? AND party_id = ?", id, partyID).Delete(&PartyServiceConfigurationModel{}).Error; err != nil {
		return domain.WrapPersistence("failed to delete party service configuration", err)
	}
	return nil
}
