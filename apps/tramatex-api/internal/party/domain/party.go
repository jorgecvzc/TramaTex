package domain

import "time"

// Party is the aggregate root for the Party module
// It can contain a PersonProfile, OrganizationProfile, or both.
type Party struct {
	id                        PartyID
	status                    PartyStatus
	personProfile             *PersonProfile
	organizationProfile       *OrganizationProfile
	roles                     []PartyRole
	relationships             []PartyRelationship
	defaultDiscountPercentage float64
	createdAt                 time.Time
	modifiedAt                time.Time
}

func NewParty(
	id PartyID,
	status PartyStatus,
	personProfile *PersonProfile,
	organizationProfile *OrganizationProfile,
) (*Party, error) {
	if id.String() == "" {
		return nil, NewValidationError("party ID cannot be empty")
	}
	if !status.IsValid() {
		return nil, NewValidationErrorf("invalid party status: %s", status)
	}
	if personProfile == nil && organizationProfile == nil {
		return nil, NewValidationError("party must have at least one profile")
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
		return nil, NewValidationError("party ID cannot be empty")
	}
	if !status.IsValid() {
		return nil, NewValidationErrorf("invalid party status: %s", status)
	}
	if personProfile == nil && organizationProfile == nil {
		return nil, NewValidationError("party must have at least one profile")
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

func (p *Party) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Party) ModifiedAt() time.Time {
	return p.modifiedAt
}

func (p *Party) SetTimestamps(createdAt, modifiedAt time.Time) {
	p.createdAt = createdAt
	p.modifiedAt = modifiedAt
}

func (p *Party) Activate() error {
	if p.status == PartyStatusActive {
		return NewConflictError("party is already active")
	}
	p.status = PartyStatusActive
	return nil
}

func (p *Party) Deactivate() error {
	if p.status == PartyStatusInactive {
		return NewConflictError("party is already inactive")
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

func (p *Party) DefaultDiscountPercentage() float64 {
	return p.defaultDiscountPercentage
}

func (p *Party) SetDefaultDiscountPercentage(percentage float64) error {
	if percentage < 0 || percentage > 100 {
		return NewValidationError("discount percentage must be between 0 and 100")
	}
	p.defaultDiscountPercentage = percentage
	return nil
}

func (p *Party) Roles() []PartyRole {
	return p.roles
}

func (p *Party) AddRole(role PartyRole) error {
	if !role.Type().IsValid() {
		return NewValidationErrorf("invalid party role: %s", role.Type())
	}
	for _, r := range p.roles {
		if r.Type() == role.Type() {
			return NewConflictError("role already exists")
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
	return NewNotFoundError("role not found")
}

func (p *Party) Relationships() []PartyRelationship {
	return p.relationships
}

func (p *Party) AddRelationship(rel PartyRelationship) error {
	if !rel.Type().IsValid() {
		return NewValidationErrorf("invalid relationship type: %s", rel.Type())
	}
	p.relationships = append(p.relationships, rel)
	return nil
}

// PartyRole links a party to a role type
type PartyRole struct {
	typeValue          PartyRoleType
	creationIdentifier *string
}

func NewPartyRole(roleType PartyRoleType, creationIdentifier *string) (PartyRole, error) {
	if !roleType.IsValid() {
		return PartyRole{}, NewValidationErrorf("invalid party role: %s", roleType)
	}
	return PartyRole{
		typeValue:          roleType,
		creationIdentifier: creationIdentifier,
	}, nil
}

func (r PartyRole) Type() PartyRoleType {
	return r.typeValue
}

func (r PartyRole) CreationIdentifier() *string {
	return r.creationIdentifier
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
		return PartyRelationship{}, NewValidationError("relationship ID cannot be empty")
	}
	if fromID.String() == "" || toID.String() == "" {
		return PartyRelationship{}, NewValidationError("from/to party IDs cannot be empty")
	}
	if !relType.IsValid() {
		return PartyRelationship{}, NewValidationErrorf("invalid relationship type: %s", relType)
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
