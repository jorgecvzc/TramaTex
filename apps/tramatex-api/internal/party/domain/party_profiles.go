package domain

// PersonProfile holds person-specific data for a Party
// Optional: a Party may have a person profile, an organization profile, or both.
type PersonProfile struct {
	firstName string
	lastName  string
}

func NewPersonProfile(firstName, lastName string) (*PersonProfile, error) {
	if firstName == "" || lastName == "" {
		return nil, NewValidationError("first name and last name are required")
	}
	return &PersonProfile{firstName: firstName, lastName: lastName}, nil
}

func (p *PersonProfile) FirstName() string {
	return p.firstName
}

func (p *PersonProfile) LastName() string {
	return p.lastName
}

// OrganizationProfile holds organization-specific data for a Party
type OrganizationProfile struct {
	name     string
	taxID    *TaxID
	website  string
	contacts []*ContactDetails
}

func NewOrganizationProfile(name string, taxID *TaxID, website string) (*OrganizationProfile, error) {
	if name == "" {
		return nil, NewValidationError("organization name is required")
	}
	return &OrganizationProfile{
		name:     name,
		taxID:    taxID,
		website:  website,
		contacts: make([]*ContactDetails, 0),
	}, nil
}

func (o *OrganizationProfile) Name() string {
	return o.name
}

func (o *OrganizationProfile) TaxID() *TaxID {
	return o.taxID
}

func (o *OrganizationProfile) Website() string {
	return o.website
}

func (o *OrganizationProfile) Contacts() []*ContactDetails {
	return o.contacts
}

func (o *OrganizationProfile) AddContact(details *ContactDetails) error {
	if details == nil {
		return NewValidationError("contact details cannot be nil")
	}
	for _, c := range o.contacts {
		if c.ID() == details.ID() {
			return NewConflictError("contact details already exist")
		}
	}
	o.contacts = append(o.contacts, details)
	return nil
}

func (o *OrganizationProfile) UpdateContact(details *ContactDetails) error {
	if details == nil {
		return NewValidationError("contact details cannot be nil")
	}
	for i, c := range o.contacts {
		if c.ID() == details.ID() {
			o.contacts[i] = details
			return nil
		}
	}
	return NewNotFoundError("contact details not found")
}

func (o *OrganizationProfile) RemoveContact(id ContactDetailsID) error {
	for i, c := range o.contacts {
		if c.ID() == id {
			o.contacts = append(o.contacts[:i], o.contacts[i+1:]...)
			return nil
		}
	}
	return NewNotFoundError("contact details not found")
}

// ContactDetails is a value object for organization contacts
// relatedPartyID is optional and can link to another Party
// (e.g., a person contact stored as a Party)
type ContactDetails struct {
	id              ContactDetailsID
	typeDescription string
	phone           *Phone
	email           *Email
	relatedPartyID  *PartyID
}

func NewContactDetails(
	id ContactDetailsID,
	typeDescription string,
	phone *Phone,
	email *Email,
	relatedPartyID *PartyID,
) (*ContactDetails, error) {
	if id.String() == "" {
		return nil, NewValidationError("contact details ID cannot be empty")
	}
	if typeDescription == "" {
		return nil, NewValidationError("contact type description cannot be empty")
	}
	return &ContactDetails{
		id:              id,
		typeDescription: typeDescription,
		phone:           phone,
		email:           email,
		relatedPartyID:  relatedPartyID,
	}, nil
}

func (c *ContactDetails) ID() ContactDetailsID {
	return c.id
}

func (c *ContactDetails) TypeDescription() string {
	return c.typeDescription
}

func (c *ContactDetails) Phone() *Phone {
	return c.phone
}

func (c *ContactDetails) Email() *Email {
	return c.email
}

func (c *ContactDetails) RelatedPartyID() *PartyID {
	return c.relatedPartyID
}
