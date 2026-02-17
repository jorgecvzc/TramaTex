package persistence

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/joran-cortez/tramatex/internal/product/domain"
)

// PartyServiceConfigurationModel represents the GORM model for a PartyServiceConfiguration.
type PartyServiceConfigurationModel struct {
	ID                   uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	PartyID              uuid.UUID    `gorm:"type:uuid;not null" json:"partyId"`
	ServiceID            string       `gorm:"type:varchar(255);not null" json:"serviceId"`
	Name                 string       `gorm:"type:varchar(255);not null" json:"name"`
	ConfigurationDetails []byte       `gorm:"type:jsonb" json:"configurationDetails"` // Store as JSONB
	CreatedAt            time.Time    `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt            time.Time    `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt            sql.NullTime `gorm:"index" json:"-"`
}

// TableName specifies the table name for GORM.
func (PartyServiceConfigurationModel) TableName() string {
	return "party_service_configurations"
}

// ToDomain converts a PartyServiceConfigurationModel to a domain.PartyServiceConfiguration.
func (m *PartyServiceConfigurationModel) ToDomain() (*domain.PartyServiceConfiguration, error) {
	configDetails := json.RawMessage(m.ConfigurationDetails)

	return &domain.PartyServiceConfiguration{
		ID:                   m.ID,
		PartyID:              m.PartyID,
		ServiceID:            m.ServiceID,
		Name:                 m.Name,
		ConfigurationDetails: configDetails,
	}, nil
}

// PartyServiceConfigurationFromDomain converts a domain.PartyServiceConfiguration to a PartyServiceConfigurationModel.
func PartyServiceConfigurationFromDomain(d *domain.PartyServiceConfiguration) (*PartyServiceConfigurationModel, error) {
	detailsBytes, err := d.ConfigurationDetails.MarshalJSON()
	if err != nil {
		return nil, domain.WrapPersistence("failed to marshal ConfigurationDetails to JSON", err)
	}

	return &PartyServiceConfigurationModel{
		ID:                   d.ID,
		PartyID:              d.PartyID,
		ServiceID:            d.ServiceID,
		Name:                 d.Name,
		ConfigurationDetails: detailsBytes,
	}, nil
}
