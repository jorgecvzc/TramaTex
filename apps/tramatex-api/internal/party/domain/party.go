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
	createdBy           string
	createdAt           time.Time
	modifiedBy          string
	modifiedAt          time.Time
}

func NewParty(
	id PartyID,
	status PartyStatus,
	createdBy string,
	personProfile *PersonProfile,
	organizationProfile *OrganizationProfile,
) (*Party, error) {
	if id.String() == "" {
		return nil, fmt.Errorf("party ID cannot be empty")
	}
	if !status.IsValid() {
		return nil, fmt.Errorf("invalid party status: %s", status)
	}
	if createdBy == "" {
		return nil, fmt.Errorf("createdBy user ID cannot be empty")
	}
	if personProfile == nil && organizationProfile == nil {
		return nil, fmt.Errorf("party must have at least one profile")
	}

	now := time.Now()
	return &Party{
		id:                  id,
		status:              status,
		personProfile:       personProfile,
		organizationProfile: organizationProfile,
		roles:               make([]PartyRole, 0),
		relationships:       make([]PartyRelationship, 0),
		createdBy:           createdBy,
		createdAt:           now,
		modifiedBy:          createdBy,
		modifiedAt:          now,
	}, nil
}

// NewPartyFromPersistence hydrates a Party from stored data
func NewPartyFromPersistence(
	id PartyID,
	status PartyStatus,
	createdBy string,
	createdAt time.Time,
	modifiedBy string,
	modifiedAt time.Time,
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
	if createdBy == "" {
		return nil, fmt.Errorf("createdBy user ID cannot be empty")
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
		createdBy:           createdBy,
		createdAt:           createdAt,
		modifiedBy:          modifiedBy,
		modifiedAt:          modifiedAt,
	}, nil
}

func (p *Party) ID() PartyID {
	return p.id
}

func (p *Party) Status() PartyStatus {
	return p.status
}

func (p *Party) CreatedBy() string {
	return p.createdBy
}

func (p *Party) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Party) ModifiedBy() string {
	return p.modifiedBy
}

func (p *Party) ModifiedAt() time.Time {
	return p.modifiedAt
}

func (p *Party) Activate(modifiedBy string) error {
	if modifiedBy == "" {
		return fmt.Errorf("modifiedBy user ID cannot be empty")
	}
	if p.status == PartyStatusActive {
		return fmt.Errorf("party is already active")
	}
	p.status = PartyStatusActive
	p.modifiedBy = modifiedBy
	p.modifiedAt = time.Now()
	return nil
}

func (p *Party) Deactivate(modifiedBy string) error {
	if modifiedBy == "" {
		return fmt.Errorf("modifiedBy user ID cannot be empty")
	}
	if p.status == PartyStatusInactive {
		return fmt.Errorf("party is already inactive")
	}
	p.status = PartyStatusInactive
	p.modifiedBy = modifiedBy
	p.modifiedAt = time.Now()
	return nil
}

func (p *Party) PersonProfile() *PersonProfile {
	return p.personProfile
}

func (p *Party) OrganizationProfile() *OrganizationProfile {
	return p.organizationProfile
}

func (p *Party) SetPersonProfile(profile *PersonProfile, modifiedBy string) error {
	if modifiedBy == "" {
		return fmt.Errorf("modifiedBy user ID cannot be empty")
	}
	p.personProfile = profile
	p.modifiedBy = modifiedBy
	p.modifiedAt = time.Now()
	return nil
}

func (p *Party) SetOrganizationProfile(profile *OrganizationProfile, modifiedBy string) error {
	if modifiedBy == "" {
		return fmt.Errorf("modifiedBy user ID cannot be empty")
	}
	p.organizationProfile = profile
	p.modifiedBy = modifiedBy
	p.modifiedAt = time.Now()
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
