package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/joran-cortez/tramatex/internal/product/domain"
	"github.com/stretchr/testify/assert"
)

func TestPartyServiceConfiguration_NewAndUpdate(t *testing.T) {
	partyID := uuid.New()
	payload := json.RawMessage(`{"key":"value"}`)

	config, err := domain.NewPartyServiceConfiguration(partyID, "service-1", "Config Name", payload)
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, partyID, config.PartyID)
	assert.Equal(t, "service-1", config.ServiceID)
	assert.Equal(t, "Config Name", config.Name)
	assert.Equal(t, payload, config.ConfigurationDetails)

	err = config.Update("service-2", "Updated", json.RawMessage(`{"key":"updated"}`))
	assert.NoError(t, err)
	assert.Equal(t, "service-2", config.ServiceID)
	assert.Equal(t, "Updated", config.Name)
}

func TestPartyServiceConfiguration_NewValidation(t *testing.T) {
	_, err := domain.NewPartyServiceConfiguration(uuid.Nil, "service", "name", nil)
	assert.Error(t, err)

	_, err = domain.NewPartyServiceConfiguration(uuid.New(), "", "name", nil)
	assert.Error(t, err)

	_, err = domain.NewPartyServiceConfiguration(uuid.New(), "service", "", nil)
	assert.Error(t, err)
}
