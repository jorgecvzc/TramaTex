package domain

import (
	"fmt"
	"time"
)

// Party is the aggregate root for the Party module
// It can contain a PersonProfile, OrganizationProfile, or both.
type Party struct {
	id                  PartyID
	status              PartyStatus
	personProfile       *PersonProfile
	organizationProfile *OrganizationProfile
	roles               []PartyRole
	relationships       []PartyRelationship
}

func NewParty(
	id PartyID,
	status PartyStatus,
	personProfile *PersonProfile,
	organizationProfile *OrganizationProfile,
) (*Party, error) {
	if id.String() == "" {
		return nil, fmt.Errorf("party ID cannot be empty")
	}
	if !status.IsValid() {
		return nil, fmt.Errorf("invalid party status: %s", status)
	}
	if personProfile == nil && organizationProfile == nil {
		return nil, fmt.Errorf("party must have at least one profile")
	}

	return &Party{
		id:                  id,
		status:              status,
		personProfile:       personProfile,
		organizationProfile: organizationProfile,
		roles:               make([]PartyRole, 0),
		relationships:       make([]PartyRelationship, 0),
	}, nil
}

// NewPartyFromPersistence hydrates a Party from stored data
func NewPartyFromPersistence(
	id PartyID,
	status PartyStatus,
	personProfile *PersonProfile,
	organizationProfile *OrganizationProfile,
	roles []PartyRole,
) (*Party, error) {
	if id.String() == "" {
		return nil, fmt.Errorf("party ID cannot be empty")
	}
	if !status.IsValid() {
		return nil, fmt.Errorf("invalid party status: %s", status)
	}
	if personProfile == nil && organizationProfile == nil {
		return nil, fmt.Errorf("party must have at least one profile")
	}

	return &Party{
		id:                  id,
		status:              status,
		personProfile:       personProfile,
		organizationProfile: organizationProfile,
		roles:               roles,
		relationships:       make([]PartyRelationship, 0),
	}, nil
}

func (p *Party) ID() PartyID {
	return p.id
}

func (p *Party) Status() PartyStatus {
	return p.status
}



func (p *Party) Activate() error {
	if p.status == PartyStatusActive {
		return fmt.Errorf("party is already active")
	}
	p.status = PartyStatusActive
	return nil
}

func (p *Party) Deactivate() error {
	if p.status == PartyStatusInactive {
		return fmt.Errorf("party is already inactive")
	}
	p.status = PartyStatusInactive
	return nil
}

func (p *Party) PersonProfile() *PersonProfile {
	return p.personProfile
}

func (p *Party) OrganizationProfile() *OrganizationProfile {
	return p.organizationProfile
}

func (p *Party) SetPersonProfile(profile *PersonProfile) error {
	p.personProfile = profile
	return nil
}

func (p *Party) SetOrganizationProfile(profile *OrganizationProfile) error {
	p.organizationProfile = profile
	return nil
}

func (p *Party) Roles() []PartyRole {
	return p.roles
}

func (p *Party) AddRole(role PartyRole) error {
	if !role.Type().IsValid() {
		return fmt.Errorf("invalid party role: %s", role.Type())
	}
	for _, r := range p.roles {
		if r.Type() == role.Type() {
			return fmt.Errorf("role already exists")
		}
	}
	p.roles = append(p.roles, role)
	return nil
}

func (p *Party) RemoveRole(roleType PartyRoleType) error {
	for i, r := range p.roles {
		if r.Type() == roleType {
			p.roles = append(p.roles[:i], p.roles[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("role not found")
}

func (p *Party) Relationships() []PartyRelationship {
	return p.relationships
}

func (p *Party) AddRelationship(rel PartyRelationship) error {
	if !rel.Type().IsValid() {
		return fmt.Errorf("invalid relationship type: %s", rel.Type())
	}
	p.relationships = append(p.relationships, rel)
	return nil
}

// PartyRole links a party to a role type
type PartyRole struct {
	typeValue PartyRoleType
}

func NewPartyRole(roleType PartyRoleType) (PartyRole, error) {
	if !roleType.IsValid() {
		return PartyRole{}, fmt.Errorf("invalid party role: %s", roleType)
	}
	return PartyRole{typeValue: roleType}, nil
}

func (r PartyRole) Type() PartyRoleType {
	return r.typeValue
}

// PartyRelationship links two parties with a relationship type
type PartyRelationship struct {
	id        PartyRelationshipID
	fromID    PartyID
	toID      PartyID
	typeValue RelationshipType
}

func NewPartyRelationship(id PartyRelationshipID, fromID PartyID, toID PartyID, relType RelationshipType) (PartyRelationship, error) {
	if id.String() == "" {
		return PartyRelationship{}, fmt.Errorf("relationship ID cannot be empty")
	}
	if fromID.String() == "" || toID.String() == "" {
		return PartyRelationship{}, fmt.Errorf("from/to party IDs cannot be empty")
	}
	if !relType.IsValid() {
		return PartyRelationship{}, fmt.Errorf("invalid relationship type: %s", relType)
	}
	return PartyRelationship{
		id:        id,
		fromID:    fromID,
		toID:      toID,
		typeValue: relType,
	}, nil
}

func (r PartyRelationship) ID() PartyRelationshipID {
	return r.id
}

func (r PartyRelationship) FromID() PartyID {
	return r.fromID
}

func (r PartyRelationship) ToID() PartyID {
	return r.toID
}

func (r PartyRelationship) Type() RelationshipType {
	return r.typeValue
}
