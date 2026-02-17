package domain

import (
	"encoding/json"

	"github.com/google/uuid"
)

// PartyServiceConfiguration represents a configuration for a specific service tied to a party.
type PartyServiceConfiguration struct {
	ID                   uuid.UUID
	PartyID              uuid.UUID
	ServiceID            string
	Name                 string
	ConfigurationDetails json.RawMessage // Flexible JSON object
}

// NewPartyServiceConfiguration creates a new PartyServiceConfiguration.
func NewPartyServiceConfiguration(partyID uuid.UUID, serviceID, name string, configDetails json.RawMessage) (*PartyServiceConfiguration, error) {
	if partyID == uuid.Nil {
		return nil, NewValidationError("party ID cannot be empty")
	}
	if serviceID == "" {
		return nil, NewValidationError("service ID cannot be empty")
	}
	if name == "" {
		return nil, NewValidationError("name cannot be empty")
	}

	return &PartyServiceConfiguration{
		ID:                   uuid.New(),
		PartyID:              partyID,
		ServiceID:            serviceID,
		Name:                 name,
		ConfigurationDetails: configDetails,
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
	return nil
}
