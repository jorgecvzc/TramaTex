package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PartyServiceConfiguration represents a configuration for a specific service tied to a party.
type PartyServiceConfiguration struct {
	ID                   uuid.UUID
	PartyID              uuid.UUID
	ServiceID            string
	Name                 string
	ConfigurationDetails json.RawMessage // Flexible JSON object
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// NewPartyServiceConfiguration creates a new PartyServiceConfiguration.
func NewPartyServiceConfiguration(partyID uuid.UUID, serviceID, name string, configDetails json.RawMessage) (*PartyServiceConfiguration, error) {
	if partyID == uuid.Nil {
		return nil, fmt.Errorf("party ID cannot be empty")
	}
	if serviceID == "" {
		return nil, fmt.Errorf("service ID cannot be empty")
	}
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	return &PartyServiceConfiguration{
		ID:                   uuid.New(),
		PartyID:              partyID,
		ServiceID:            serviceID,
		Name:                 name,
		ConfigurationDetails: configDetails,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}, nil
}

// Update updates the PartyServiceConfiguration's fields.
func (psc *PartyServiceConfiguration) Update(serviceID, name string, configDetails json.RawMessage) error {
	if serviceID != "" {
		psc.ServiceID = serviceID
	}
	if name != "" {
		psc.Name = name
	}
	if configDetails != nil {
		psc.ConfigurationDetails = configDetails
	}
	psc.UpdatedAt = time.Now()
	return nil
}
