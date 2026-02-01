package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/joran-cortez/tramatex/internal/party/domain"
)

// PostgreSQLOrganizationRepository implements OrganizationRepository using PostgreSQL
type PostgreSQLOrganizationRepository struct {
	db *sql.DB
}

// NewPostgreSQLOrganizationRepository creates a new PostgreSQL organization repository
func NewPostgreSQLOrganizationRepository(db *sql.DB) *PostgreSQLOrganizationRepository {
	return &PostgreSQLOrganizationRepository{db: db}
}

// Save saves an organization to the database
func (r *PostgreSQLOrganizationRepository) Save(ctx context.Context, org *domain.Organization) error {
	query := `
		INSERT INTO organizations (id, name, role, status, tax_id, website, notes, created_by, created_at, modified_by, modified_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO UPDATE SET
			name = $2,
			role = $3,
			status = $4,
			tax_id = $5,
			website = $6,
			notes = $7,
			modified_by = $10,
			modified_at = $11
	`

	now := time.Now()
	taxIDValue := ""
	if org.TaxID() != nil {
		taxIDValue = org.TaxID().Value()
	}
	_, err := r.db.ExecContext(ctx, query,
		org.ID().Value(),
		org.Name(),
		org.Role().String(),
		org.Status().String(),
		taxIDValue,
		org.Website(),
		org.Notes(),
		org.CreatedBy(),
		org.CreatedAt(),
		org.ModifiedBy(),
		now,
	)

	return err
}

// FindByID finds an organization by ID
func (r *PostgreSQLOrganizationRepository) FindByID(ctx context.Context, id domain.OrganizationID) (*domain.Organization, error) {
	query := `
		SELECT id, name, role, status, tax_id, website, notes, created_by, created_at, modified_by, modified_at
		FROM organizations
		WHERE id = $1
	`

	var orgID, name, role, status, taxID, website, notes, createdBy, modifiedBy string
	var createdAt, modifiedAt time.Time

	err := r.db.QueryRowContext(ctx, query, id.Value()).Scan(
		&orgID, &name, &role, &status, &taxID, &website, &notes, &createdBy, &createdAt, &modifiedBy, &modifiedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("organization not found")
	}
	if err != nil {
		return nil, err
	}

	newID, _ := domain.NewOrganizationID(orgID)
	roleEnum := domain.ParseOrganizationRole(role)
	statusEnum := domain.ParseOrganizationStatus(status)

	var taxIDObj *domain.TaxID
	if taxID != "" {
		taxIDObj, _ = domain.NewTaxID(taxID, "NIF") // Default type
	}

	org, err := domain.NewOrganization(newID, name, roleEnum, taxIDObj, createdBy)
	if err != nil {
		return nil, err
	}

	// Set the website and notes if they exist
	if website != "" {
		if err := org.UpdateWebsite(website, modifiedBy); err != nil {
			return nil, err
		}
	}
	if notes != "" {
		if err := org.UpdateNotes(notes, modifiedBy); err != nil {
			return nil, err
		}
	}

	// Set the status if it's inactive
	if statusEnum == domain.OrganizationStatusInactive {
		if err := org.Deactivate(modifiedBy); err != nil {
			return nil, err
		}
	}

	return org, nil
}

// FindAll finds all organizations with filters
func (r *PostgreSQLOrganizationRepository) FindAll(ctx context.Context, filters *OrganizationFilters) ([]*domain.Organization, error) {
	query := "SELECT id, name, role, status, tax_id, website, notes, created_by, created_at, modified_by, modified_at FROM organizations WHERE 1=1"

	args := []interface{}{}
	argCount := 1

	if filters.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, filters.Status.String())
		argCount++
	}
	if filters.Role != nil {
		query += fmt.Sprintf(" AND role = $%d", argCount)
		args = append(args, filters.Role.String())
		argCount++
	}
	if filters.Name != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argCount)
		args = append(args, "%"+filters.Name+"%")
		argCount++
	}

	query += " ORDER BY name"

	if filters.PageSize > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.PageSize)
		argCount++

		query += fmt.Sprintf(" OFFSET $%d", argCount)
		offset := (filters.PageNumber - 1) * filters.PageSize
		args = append(args, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var organizations []*domain.Organization
	for rows.Next() {
		var orgID, name, role, status, taxID, website, notes, createdBy, modifiedBy string
		var createdAt, modifiedAt time.Time

		err := rows.Scan(&orgID, &name, &role, &status, &taxID, &website, &notes, &createdBy, &createdAt, &modifiedBy, &modifiedAt)
		if err != nil {
			return nil, err
		}

		newID, _ := domain.NewOrganizationID(orgID)
		roleEnum := domain.ParseOrganizationRole(role)

		var taxIDObj *domain.TaxID
		if taxID != "" {
			taxIDObj, _ = domain.NewTaxID(taxID, "NIF")
		}

		org, err := domain.NewOrganization(newID, name, roleEnum, taxIDObj, createdBy)
		if err != nil {
			return nil, err
		}

		// Set website and notes if they exist
		if website != "" {
			if err := org.UpdateWebsite(website, modifiedBy); err != nil {
				return nil, err
			}
		}
		if notes != "" {
			if err := org.UpdateNotes(notes, modifiedBy); err != nil {
				return nil, err
			}
		}

		organizations = append(organizations, org)
	}

	return organizations, rows.Err()
}

// FindByRole finds all organizations with a specific role
func (r *PostgreSQLOrganizationRepository) FindByRole(ctx context.Context, role domain.OrganizationRole) ([]*domain.Organization, error) {
	query := `
		SELECT id, name, role, status, tax_id, website, notes, created_by, created_at, modified_by, modified_at
		FROM organizations
		WHERE role = $1
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query, role.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var organizations []*domain.Organization
	for rows.Next() {
		var orgID, name, roleStr, status, taxID, website, notes, createdBy, modifiedBy string
		var createdAt, modifiedAt time.Time

		err := rows.Scan(&orgID, &name, &roleStr, &status, &taxID, &website, &notes, &createdBy, &createdAt, &modifiedBy, &modifiedAt)
		if err != nil {
			return nil, err
		}

		newID, _ := domain.NewOrganizationID(orgID)
		roleEnum := domain.ParseOrganizationRole(roleStr)

		var taxIDObj *domain.TaxID
		if taxID != "" {
			taxIDObj, _ = domain.NewTaxID(taxID, "NIF")
		}

		org, err := domain.NewOrganization(newID, name, roleEnum, taxIDObj, createdBy)
		if err != nil {
			return nil, err
		}

		if website != "" {
			if err := org.UpdateWebsite(website, modifiedBy); err != nil {
				return nil, err
			}
		}
		if notes != "" {
			if err := org.UpdateNotes(notes, modifiedBy); err != nil {
				return nil, err
			}
		}

		organizations = append(organizations, org)
	}

	return organizations, rows.Err()
}

// Delete deletes an organization
func (r *PostgreSQLOrganizationRepository) Delete(ctx context.Context, id domain.OrganizationID) error {
	query := "DELETE FROM organizations WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id.Value())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("organization not found")
	}

	return nil
}

// Exists checks if an organization exists
func (r *PostgreSQLOrganizationRepository) Exists(ctx context.Context, id domain.OrganizationID) (bool, error) {
	query := "SELECT 1 FROM organizations WHERE id = $1"
	var exists int
	err := r.db.QueryRowContext(ctx, query, id.Value()).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// Count counts total organizations
func (r *PostgreSQLOrganizationRepository) Count(ctx context.Context) (int64, error) {
	query := "SELECT COUNT(*) FROM organizations"
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

// PostgreSQLPersonRepository implements PersonRepository using PostgreSQL
type PostgreSQLPersonRepository struct {
	db *sql.DB
}

// NewPostgreSQLPersonRepository creates a new PostgreSQL person repository
func NewPostgreSQLPersonRepository(db *sql.DB) *PostgreSQLPersonRepository {
	return &PostgreSQLPersonRepository{db: db}
}

// Save saves a person to the database
func (r *PostgreSQLPersonRepository) Save(ctx context.Context, person *domain.Person) error {
	query := `
		INSERT INTO persons (id, organization_id, first_name, last_name, email, phone, job_title, is_primary_contact, created_by, created_at, modified_by, modified_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (id) DO UPDATE SET
			first_name = $3,
			last_name = $4,
			email = $5,
			phone = $6,
			job_title = $7,
			is_primary_contact = $8,
			modified_by = $11,
			modified_at = $12
	`

	now := time.Now()
	phone := ""
	if person.Phone() != nil {
		phone = person.Phone().Value()
	}

	_, err := r.db.ExecContext(ctx, query,
		person.ID().Value(),
		person.OrganizationID().Value(),
		person.FirstName(),
		person.LastName(),
		person.Email().Value(),
		phone,
		person.JobTitle(),
		person.IsPrimaryContact(),
		person.CreatedBy(),
		person.CreatedAt(),
		person.ModifiedBy(),
		now,
	)

	return err
}

// FindByID finds a person by ID
func (r *PostgreSQLPersonRepository) FindByID(ctx context.Context, id domain.PersonID) (*domain.Person, error) {
	query := `
		SELECT id, organization_id, first_name, last_name, email, phone, job_title, is_primary_contact, created_by, created_at, modified_by, modified_at
		FROM persons
		WHERE id = $1
	`

	var personID, orgID, firstName, lastName, email, phone, jobTitle, createdBy, modifiedBy string
	var isPrimaryContact bool
	var createdAt, modifiedAt time.Time

	err := r.db.QueryRowContext(ctx, query, id.Value()).Scan(
		&personID, &orgID, &firstName, &lastName, &email, &phone, &jobTitle, &isPrimaryContact, &createdBy, &createdAt, &modifiedBy, &modifiedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("person not found")
	}
	if err != nil {
		return nil, err
	}

	newID, _ := domain.NewPersonID(personID)
	newOrgID, _ := domain.NewOrganizationID(orgID)
	emailVO, _ := domain.NewEmail(email)
	var phoneVO *domain.Phone
	if phone != "" {
		phoneVO, _ = domain.NewPhone(phone)
	}

	person := domain.NewPerson(newID, newOrgID, firstName, lastName, emailVO, createdBy)
	if phoneVO != nil {
		person.SetPhone(phoneVO)
	}
	if jobTitle != "" {
		person.SetJobTitle(jobTitle)
	}
	if isPrimaryContact {
		person.SetPrimaryContact(true)
	}

	return person, nil
}

// FindByOrganization finds all persons in an organization
func (r *PostgreSQLPersonRepository) FindByOrganization(ctx context.Context, orgID domain.OrganizationID) ([]*domain.Person, error) {
	query := `
		SELECT id, organization_id, first_name, last_name, email, phone, job_title, is_primary_contact, created_by, created_at, modified_by, modified_at
		FROM persons
		WHERE organization_id = $1
		ORDER BY last_name, first_name
	`

	rows, err := r.db.QueryContext(ctx, query, orgID.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var persons []*domain.Person
	for rows.Next() {
		var personID, retrievedOrgID, firstName, lastName, email, phone, jobTitle, createdBy, modifiedBy string
		var isPrimaryContact bool
		var createdAt, modifiedAt time.Time

		err := rows.Scan(&personID, &retrievedOrgID, &firstName, &lastName, &email, &phone, &jobTitle, &isPrimaryContact, &createdBy, &createdAt, &modifiedBy, &modifiedAt)
		if err != nil {
			return nil, err
		}

		newID, _ := domain.NewPersonID(personID)
		emailVO, _ := domain.NewEmail(email)
		var phoneVO *domain.Phone
		if phone != "" {
			phoneVO, _ = domain.NewPhone(phone)
		}

		person := domain.NewPerson(newID, orgID, firstName, lastName, emailVO, createdBy)
		if phoneVO != nil {
			person.SetPhone(phoneVO)
		}
		if jobTitle != "" {
			person.SetJobTitle(jobTitle)
		}
		if isPrimaryContact {
			person.SetPrimaryContact(true)
		}

		persons = append(persons, person)
	}

	return persons, rows.Err()
}

// FindByEmail finds a person by email
func (r *PostgreSQLPersonRepository) FindByEmail(ctx context.Context, email string) (*domain.Person, error) {
	query := `
		SELECT id, organization_id, first_name, last_name, email, phone, job_title, is_primary_contact, created_by, created_at, modified_by, modified_at
		FROM persons
		WHERE LOWER(email) = LOWER($1)
	`

	var personID, orgID, firstName, lastName, emailStr, phone, jobTitle, createdBy, modifiedBy string
	var isPrimaryContact bool
	var createdAt, modifiedAt time.Time

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&personID, &orgID, &firstName, &lastName, &emailStr, &phone, &jobTitle, &isPrimaryContact, &createdBy, &createdAt, &modifiedBy, &modifiedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("person not found")
	}
	if err != nil {
		return nil, err
	}

	newID, _ := domain.NewPersonID(personID)
	newOrgID, _ := domain.NewOrganizationID(orgID)
	emailVO, _ := domain.NewEmail(emailStr)
	var phoneVO *domain.Phone
	if phone != "" {
		phoneVO, _ = domain.NewPhone(phone)
	}

	person := domain.NewPerson(newID, newOrgID, firstName, lastName, emailVO, createdBy)
	if phoneVO != nil {
		person.SetPhone(phoneVO)
	}
	if jobTitle != "" {
		person.SetJobTitle(jobTitle)
	}
	if isPrimaryContact {
		person.SetPrimaryContact(true)
	}

	return person, nil
}

// FindPrimaryContact finds the primary contact for an organization
func (r *PostgreSQLPersonRepository) FindPrimaryContact(ctx context.Context, orgID domain.OrganizationID) (*domain.Person, error) {
	query := `
		SELECT id, organization_id, first_name, last_name, email, phone, job_title, is_primary_contact, created_by, created_at, modified_by, modified_at
		FROM persons
		WHERE organization_id = $1 AND is_primary_contact = true
		LIMIT 1
	`

	var personID, retrievedOrgID, firstName, lastName, email, phone, jobTitle, createdBy, modifiedBy string
	var isPrimaryContact bool
	var createdAt, modifiedAt time.Time

	err := r.db.QueryRowContext(ctx, query, orgID.Value()).Scan(
		&personID, &retrievedOrgID, &firstName, &lastName, &email, &phone, &jobTitle, &isPrimaryContact, &createdBy, &createdAt, &modifiedBy, &modifiedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("primary contact not found")
	}
	if err != nil {
		return nil, err
	}

	newID, _ := domain.NewPersonID(personID)
	emailVO, _ := domain.NewEmail(email)
	var phoneVO *domain.Phone
	if phone != "" {
		phoneVO, _ = domain.NewPhone(phone)
	}

	person := domain.NewPerson(newID, orgID, firstName, lastName, emailVO, createdBy)
	if phoneVO != nil {
		person.SetPhone(phoneVO)
	}
	if jobTitle != "" {
		person.SetJobTitle(jobTitle)
	}

	return person, nil
}

// Delete deletes a person
func (r *PostgreSQLPersonRepository) Delete(ctx context.Context, id domain.PersonID) error {
	query := "DELETE FROM persons WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id.Value())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("person not found")
	}

	return nil
}

// Exists checks if a person exists
func (r *PostgreSQLPersonRepository) Exists(ctx context.Context, id domain.PersonID) (bool, error) {
	query := "SELECT 1 FROM persons WHERE id = $1"
	var exists int
	err := r.db.QueryRowContext(ctx, query, id.Value()).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// PostgreSQLAddressRepository implements AddressRepository using PostgreSQL
type PostgreSQLAddressRepository struct {
	db *sql.DB
}

// NewPostgreSQLAddressRepository creates a new PostgreSQL address repository
func NewPostgreSQLAddressRepository(db *sql.DB) *PostgreSQLAddressRepository {
	return &PostgreSQLAddressRepository{db: db}
}

// Save saves an address to the database
func (r *PostgreSQLAddressRepository) Save(ctx context.Context, addr *domain.Address, id domain.AddressID, orgID domain.OrganizationID) error {
	query := `
		INSERT INTO addresses (id, organization_id, street, city, province, postal_code, country, is_primary, created_by, created_at, modified_by, modified_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (id) DO UPDATE SET
			street = $3,
			city = $4,
			province = $5,
			postal_code = $6,
			country = $7,
			is_primary = $8,
			modified_by = $11,
			modified_at = $12
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query,
		id.Value(),
		orgID.Value(),
		addr.Street(),
		addr.City(),
		addr.Province(),
		addr.PostalCode(),
		addr.Country(),
		false, // is_primary
		"system",
		now,
		"system",
		now,
	)

	return err
}

// FindByID finds an address by ID
func (r *PostgreSQLAddressRepository) FindByID(ctx context.Context, id domain.AddressID) (*domain.Address, error) {
	query := `
		SELECT street, city, province, postal_code, country
		FROM addresses
		WHERE id = $1
	`

	var street, city, province, postalCode, country string

	err := r.db.QueryRowContext(ctx, query, id.Value()).Scan(&street, &city, &province, &postalCode, &country)

	if err == sql.ErrNoRows {
		return nil, errors.New("address not found")
	}
	if err != nil {
		return nil, err
	}

	addr, _ := domain.NewAddress(street, city, province, postalCode, country)
	return addr, nil
}

// FindByOrganization finds all addresses for an organization
func (r *PostgreSQLAddressRepository) FindByOrganization(ctx context.Context, orgID domain.OrganizationID) ([]*domain.Address, error) {
	query := `
		SELECT street, city, province, postal_code, country
		FROM addresses
		WHERE organization_id = $1
		ORDER BY is_primary DESC, city
	`

	rows, err := r.db.QueryContext(ctx, query, orgID.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []*domain.Address
	for rows.Next() {
		var street, city, province, postalCode, country string

		err := rows.Scan(&street, &city, &province, &postalCode, &country)
		if err != nil {
			return nil, err
		}

		addr, _ := domain.NewAddress(street, city, province, postalCode, country)
		addresses = append(addresses, addr)
	}

	return addresses, rows.Err()
}

// FindPrimary finds the primary address for an organization
func (r *PostgreSQLAddressRepository) FindPrimary(ctx context.Context, orgID domain.OrganizationID) (*domain.Address, error) {
	query := `
		SELECT street, city, province, postal_code, country
		FROM addresses
		WHERE organization_id = $1 AND is_primary = true
		LIMIT 1
	`

	var street, city, province, postalCode, country string

	err := r.db.QueryRowContext(ctx, query, orgID.Value()).Scan(&street, &city, &province, &postalCode, &country)

	if err == sql.ErrNoRows {
		return nil, errors.New("primary address not found")
	}
	if err != nil {
		return nil, err
	}

	addr, _ := domain.NewAddress(street, city, province, postalCode, country)
	return addr, nil
}

// Delete deletes an address
func (r *PostgreSQLAddressRepository) Delete(ctx context.Context, id domain.AddressID) error {
	query := "DELETE FROM addresses WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id.Value())
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("address not found")
	}

	return nil
}

// Exists checks if an address exists
func (r *PostgreSQLAddressRepository) Exists(ctx context.Context, id domain.AddressID) (bool, error) {
	query := "SELECT 1 FROM addresses WHERE id = $1"
	var exists int
	err := r.db.QueryRowContext(ctx, query, id.Value()).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
