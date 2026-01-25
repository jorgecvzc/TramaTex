package domain

import "time"

// Person represents a contact person within an organization
type Person struct {
	id               PersonID
	organizationID   OrganizationID
	firstName        string
	lastName         string
	email            *Email
	phone            *Phone
	jobTitle         string
	isPrimaryContact bool
	createdBy        string
	createdAt        time.Time
	modifiedBy       string
	modifiedAt       time.Time
}

// NewPerson creates a new Person
func NewPerson(
	id PersonID,
	organizationID OrganizationID,
	firstName string,
	lastName string,
	email *Email,
	createdBy string,
) *Person {
	now := time.Now()
	return &Person{
		id:               id,
		organizationID:   organizationID,
		firstName:        firstName,
		lastName:         lastName,
		email:            email,
		createdBy:        createdBy,
		createdAt:        now,
		modifiedBy:       createdBy,
		modifiedAt:       now,
		isPrimaryContact: false,
	}
}

// ID returns the person ID
func (p *Person) ID() PersonID {
	return p.id
}

// OrganizationID returns the organization ID
func (p *Person) OrganizationID() OrganizationID {
	return p.organizationID
}

// FirstName returns the first name
func (p *Person) FirstName() string {
	return p.firstName
}

// LastName returns the last name
func (p *Person) LastName() string {
	return p.lastName
}

// FullName returns full name
func (p *Person) FullName() string {
	return p.firstName + " " + p.lastName
}

// Email returns the email
func (p *Person) Email() *Email {
	return p.email
}

// Phone returns the phone
func (p *Person) Phone() *Phone {
	return p.phone
}

// SetPhone sets the phone
func (p *Person) SetPhone(phone *Phone) {
	p.phone = phone
}

// JobTitle returns the job title
func (p *Person) JobTitle() string {
	return p.jobTitle
}

// SetJobTitle sets the job title
func (p *Person) SetJobTitle(jobTitle string) {
	p.jobTitle = jobTitle
}

// IsPrimaryContact returns if this is the primary contact
func (p *Person) IsPrimaryContact() bool {
	return p.isPrimaryContact
}

// SetPrimaryContact marks this person as primary contact
func (p *Person) SetPrimaryContact(isPrimary bool) {
	p.isPrimaryContact = isPrimary
}

// CreatedBy returns who created the person
func (p *Person) CreatedBy() string {
	return p.createdBy
}

// CreatedAt returns when the person was created
func (p *Person) CreatedAt() time.Time {
	return p.createdAt
}

// ModifiedBy returns who last modified the person
func (p *Person) ModifiedBy() string {
	return p.modifiedBy
}

// ModifiedAt returns when the person was last modified
func (p *Person) ModifiedAt() time.Time {
	return p.modifiedAt
}

// SetAuditTrail sets the audit trail information
func (p *Person) SetAuditTrail(createdAt time.Time, modifiedBy string, modifiedAt time.Time) {
	p.createdAt = createdAt
	p.modifiedBy = modifiedBy
	p.modifiedAt = modifiedAt
}

// Update updates the person's details
func (p *Person) Update(firstName, lastName, jobTitle, modifiedBy string) {
	p.firstName = firstName
	p.lastName = lastName
	p.jobTitle = jobTitle
	p.modifiedBy = modifiedBy
	p.modifiedAt = time.Now()
}
