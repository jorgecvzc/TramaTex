package persistence

import (
	"time"

	"github.com/joran-cortez/tramatex/internal/party/domain"
)

type PartyDataModel struct {
	ID         string    `gorm:"primaryKey;column:id"`
	Status     string    `gorm:"column:status"`
	CreatedBy  string    `gorm:"column:created_by"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	ModifiedBy string    `gorm:"column:modified_by"`
	ModifiedAt time.Time `gorm:"column:modified_at"`
}

func (PartyDataModel) TableName() string {
	return "parties"
}

type PersonProfileDataModel struct {
	PartyID   string  `gorm:"primaryKey;column:party_id"`
	FirstName string  `gorm:"column:first_name"`
	LastName  string  `gorm:"column:last_name"`
	Phone     *string `gorm:"column:phone"`
	Email     *string `gorm:"column:email"`
}

func (PersonProfileDataModel) TableName() string {
	return "person_profiles"
}

type OrganizationProfileDataModel struct {
	PartyID   string  `gorm:"primaryKey;column:party_id"`
	Name      string  `gorm:"column:name"`
	TaxID     *string `gorm:"column:tax_id"`
	TaxIDType *string `gorm:"column:tax_id_type"`
	Website   string  `gorm:"column:website"`
	Phone     *string `gorm:"column:phone"`
	Email     *string `gorm:"column:email"`
}

func (OrganizationProfileDataModel) TableName() string {
	return "organization_profiles"
}

type PartyRoleDataModel struct {
	PartyID            string  `gorm:"primaryKey;column:party_id"`
	Role               string  `gorm:"primaryKey;column:role"`
	CreationIdentifier *string `gorm:"column:creation_identifier"`
}

func (PartyRoleDataModel) TableName() string {
	return "party_roles"
}

type PartyRelationshipDataModel struct {
	ID         string    `gorm:"primaryKey;column:id"`
	FromParty  string    `gorm:"column:from_party_id"`
	ToParty    string    `gorm:"column:to_party_id"`
	Type       string    `gorm:"column:type"`
	CreatedBy  string    `gorm:"column:created_by"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	ModifiedBy string    `gorm:"column:modified_by"`
	ModifiedAt time.Time `gorm:"column:modified_at"`
}

func (PartyRelationshipDataModel) TableName() string {
	return "party_relationships"
}

type ContactDetailsDataModel struct {
	ID                string  `gorm:"primaryKey;column:id"`
	OrganizationParty string  `gorm:"column:organization_party_id"`
	TypeDescription   string  `gorm:"column:type_description"`
	Phone             *string `gorm:"column:phone"`
	Email             *string `gorm:"column:email"`
	RelatedPartyID    *string `gorm:"column:related_party_id"`
}

func (ContactDetailsDataModel) TableName() string {
	return "contact_details"
}

type PartyAddressDataModel struct {
	ID         string    `gorm:"primaryKey;column:id"`
	PartyID    string    `gorm:"column:party_id"`
	Street     string    `gorm:"column:street"`
	City       string    `gorm:"column:city"`
	Province   string    `gorm:"column:province"`
	PostalCode string    `gorm:"column:postal_code"`
	Country    string    `gorm:"column:country"`
	IsPrimary  bool      `gorm:"column:is_primary"`
	CreatedBy  string    `gorm:"column:created_by"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	ModifiedBy string    `gorm:"column:modified_by"`
	ModifiedAt time.Time `gorm:"column:modified_at"`
}

func (PartyAddressDataModel) TableName() string {
	return "party_addresses"
}

func partyDataModelFromDomain(party *domain.Party) *PartyDataModel {
	return &PartyDataModel{
		ID:     party.ID().Value(),
		Status: string(party.Status()),
	}
}
