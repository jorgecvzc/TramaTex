package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/joran-cortez/tramatex/internal/party/domain"
)

// PostgreSQLPartyRepository implements PartyRepository using PostgreSQL

type PostgreSQLPartyRepository struct {
	db *sql.DB
}

func NewPostgreSQLPartyRepository(db *sql.DB) *PostgreSQLPartyRepository {
	return &PostgreSQLPartyRepository{db: db}
}

func (r *PostgreSQLPartyRepository) Save(ctx context.Context, party *domain.Party) error {
	if party == nil {
		return fmt.Errorf("party cannot be nil")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	now := time.Now()

	partyQuery := `
		INSERT INTO parties (id, status, created_by, created_at, modified_by, modified_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			status = $2,
			modified_by = $5,
			modified_at = $6
	`

	_, err = tx.ExecContext(ctx, partyQuery,
		party.ID().Value(),
		string(party.Status()),
		party.CreatedBy(),
		party.CreatedAt(),
		party.ModifiedBy(),
		now,
	)
	if err != nil {
		return err
	}

	// Person profile
	if party.PersonProfile() != nil {
		profileQuery := `
			INSERT INTO person_profiles (party_id, first_name, last_name)
			VALUES ($1, $2, $3)
			ON CONFLICT (party_id) DO UPDATE SET
				first_name = $2,
				last_name = $3
		`
		_, err = tx.ExecContext(ctx, profileQuery,
			party.ID().Value(),
			party.PersonProfile().FirstName(),
			party.PersonProfile().LastName(),
		)
		if err != nil {
			return err
		}
	} else {
		_, err = tx.ExecContext(ctx, "DELETE FROM person_profiles WHERE party_id = $1", party.ID().Value())
		if err != nil {
			return err
		}
	}

	// Organization profile
	if party.OrganizationProfile() != nil {
		org := party.OrganizationProfile()
		taxID := sql.NullString{Valid: false}
		taxIDType := sql.NullString{Valid: false}
		if org.TaxID() != nil {
			taxID = sql.NullString{String: org.TaxID().Value(), Valid: true}
			taxIDType = sql.NullString{String: org.TaxID().Type(), Valid: true}
		}

		orgQuery := `
			INSERT INTO organization_profiles (party_id, name, tax_id, tax_id_type, website)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (party_id) DO UPDATE SET
				name = $2,
				tax_id = $3,
				tax_id_type = $4,
				website = $5
		`
		_, err = tx.ExecContext(ctx, orgQuery,
			party.ID().Value(),
			org.Name(),
			taxID,
			taxIDType,
			org.Website(),
		)
		if err != nil {
			return err
		}

		// Replace contacts
		_, err = tx.ExecContext(ctx, "DELETE FROM contact_details WHERE organization_party_id = $1", party.ID().Value())
		if err != nil {
			return err
		}

		contactQuery := `
			INSERT INTO contact_details (id, organization_party_id, type_description, phone, email, related_party_id)
			VALUES ($1, $2, $3, $4, $5, $6)
		`
		for _, contact := range org.Contacts() {
			phone := sql.NullString{Valid: false}
			if contact.Phone() != nil {
				phone = sql.NullString{String: contact.Phone().Value(), Valid: true}
			}
			email := sql.NullString{Valid: false}
			if contact.Email() != nil {
				email = sql.NullString{String: contact.Email().Value(), Valid: true}
			}
			related := sql.NullString{Valid: false}
			if contact.RelatedPartyID() != nil {
				related = sql.NullString{String: contact.RelatedPartyID().String(), Valid: true}
			}

			_, err = tx.ExecContext(ctx, contactQuery,
				contact.ID().Value(),
				party.ID().Value(),
				contact.TypeDescription(),
				phone,
				email,
				related,
			)
			if err != nil {
				return err
			}
		}
	} else {
		_, err = tx.ExecContext(ctx, "DELETE FROM contact_details WHERE organization_party_id = $1", party.ID().Value())
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, "DELETE FROM organization_profiles WHERE party_id = $1", party.ID().Value())
		if err != nil {
			return err
		}
	}

	// Roles (replace)
	_, err = tx.ExecContext(ctx, "DELETE FROM party_roles WHERE party_id = $1", party.ID().Value())
	if err != nil {
		return err
	}
	roleQuery := `
		INSERT INTO party_roles (party_id, role) VALUES ($1, $2)
	`
	for _, role := range party.Roles() {
		_, err = tx.ExecContext(ctx, roleQuery, party.ID().Value(), string(role.Type()))
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *PostgreSQLPartyRepository) FindByID(ctx context.Context, id domain.PartyID) (*domain.Party, error) {
	query := `
		SELECT id, status, created_by, created_at, modified_by, modified_at
		FROM parties
		WHERE id = $1
	`

	var partyID, status, createdBy, modifiedBy string
	var createdAt, modifiedAt time.Time

	err := r.db.QueryRowContext(ctx, query, id.Value()).Scan(&partyID, &status, &createdBy, &createdAt, &modifiedBy, &modifiedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("party not found")
	}
	if err != nil {
		return nil, err
	}

	parsedID, _ := domain.NewPartyID(partyID)
	statusEnum := domain.PartyStatus(status)
	if !statusEnum.IsValid() {
		return nil, fmt.Errorf("invalid party status in storage: %s", status)
	}

	personProfile, err := r.loadPersonProfile(ctx, parsedID)
	if err != nil {
		return nil, err
	}
	organizationProfile, err := r.loadOrganizationProfile(ctx, parsedID)
	if err != nil {
		return nil, err
	}

	roles, err := r.loadPartyRoles(ctx, parsedID)
	if err != nil {
		return nil, err
	}

	return domain.NewPartyFromPersistence(
		parsedID,
		statusEnum,
		createdBy,
		createdAt,
		modifiedBy,
		modifiedAt,
		personProfile,
		organizationProfile,
		roles,
	)
}

func (r *PostgreSQLPartyRepository) FindAll(ctx context.Context, filters *PartyFilters) ([]*domain.Party, error) {
	query := `
		SELECT DISTINCT p.id
		FROM parties p
		LEFT JOIN organization_profiles op ON op.party_id = p.id
		LEFT JOIN person_profiles pp ON pp.party_id = p.id
		LEFT JOIN party_roles pr ON pr.party_id = p.id
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	if filters != nil {
		if filters.Status != nil {
			query += fmt.Sprintf(" AND p.status = $%d", argCount)
			args = append(args, string(*filters.Status))
			argCount++
		}
		if filters.Role != nil {
			query += fmt.Sprintf(" AND pr.role = $%d", argCount)
			args = append(args, string(*filters.Role))
			argCount++
		}
		if filters.Type != "" {
			switch filters.Type {
			case "PERSON":
				query += " AND pp.party_id IS NOT NULL"
			case "ORGANIZATION":
				query += " AND op.party_id IS NOT NULL"
			case "BOTH":
				query += " AND pp.party_id IS NOT NULL AND op.party_id IS NOT NULL"
			}
		}
		if filters.Name != "" {
			query += fmt.Sprintf(" AND (op.name ILIKE $%d OR (pp.first_name || ' ' || pp.last_name) ILIKE $%d)", argCount, argCount)
			args = append(args, "%"+filters.Name+"%")
			argCount++
		}
		if filters.TaxID != "" {
			query += fmt.Sprintf(" AND op.tax_id ILIKE $%d", argCount)
			args = append(args, "%"+filters.TaxID+"%")
			argCount++
		}
	}

	query += " ORDER BY p.id"

	if filters != nil && filters.PageSize > 0 {
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

	parties := make([]*domain.Party, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		partyID, _ := domain.NewPartyID(id)
		party, err := r.FindByID(ctx, partyID)
		if err != nil {
			return nil, err
		}
		parties = append(parties, party)
	}

	return parties, rows.Err()
}

func (r *PostgreSQLPartyRepository) Delete(ctx context.Context, id domain.PartyID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM parties WHERE id = $1", id.Value())
	return err
}

func (r *PostgreSQLPartyRepository) Exists(ctx context.Context, id domain.PartyID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM parties WHERE id = $1)"
	var exists bool
	if err := r.db.QueryRowContext(ctx, query, id.Value()).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *PostgreSQLPartyRepository) Count(ctx context.Context) (int64, error) {
	query := "SELECT COUNT(*) FROM parties"
	var count int64
	if err := r.db.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *PostgreSQLPartyRepository) loadPersonProfile(ctx context.Context, partyID domain.PartyID) (*domain.PersonProfile, error) {
	query := "SELECT first_name, last_name FROM person_profiles WHERE party_id = $1"
	var firstName, lastName string
	err := r.db.QueryRowContext(ctx, query, partyID.Value()).Scan(&firstName, &lastName)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return domain.NewPersonProfile(firstName, lastName)
}

func (r *PostgreSQLPartyRepository) loadOrganizationProfile(ctx context.Context, partyID domain.PartyID) (*domain.OrganizationProfile, error) {
	query := "SELECT name, tax_id, tax_id_type, website FROM organization_profiles WHERE party_id = $1"
	var name, taxID, taxIDType, website sql.NullString
	err := r.db.QueryRowContext(ctx, query, partyID.Value()).Scan(&name, &taxID, &taxIDType, &website)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var tax *domain.TaxID
	if taxID.Valid {
		typeValue := "NIF"
		if taxIDType.Valid && taxIDType.String != "" {
			typeValue = taxIDType.String
		}
		created, err := domain.NewTaxID(taxID.String, typeValue)
		if err != nil {
			return nil, err
		}
		tax = created
	}

	profile, err := domain.NewOrganizationProfile(name.String, tax, website.String)
	if err != nil {
		return nil, err
	}

	contacts, err := r.loadContactDetails(ctx, partyID)
	if err != nil {
		return nil, err
	}
	for _, contact := range contacts {
		if err := profile.AddContact(contact); err != nil {
			return nil, err
		}
	}

	return profile, nil
}

func (r *PostgreSQLPartyRepository) loadPartyRoles(ctx context.Context, partyID domain.PartyID) ([]domain.PartyRole, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT role FROM party_roles WHERE party_id = $1", partyID.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]domain.PartyRole, 0)
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roleType := domain.PartyRoleType(role)
		if !roleType.IsValid() {
			return nil, fmt.Errorf("invalid role in storage: %s", role)
		}
		partyRole, err := domain.NewPartyRole(roleType)
		if err != nil {
			return nil, err
		}
		roles = append(roles, partyRole)
	}

	return roles, rows.Err()
}

func (r *PostgreSQLPartyRepository) loadContactDetails(ctx context.Context, partyID domain.PartyID) ([]*domain.ContactDetails, error) {
	query := "SELECT id, type_description, phone, email, related_party_id FROM contact_details WHERE organization_party_id = $1"
	rows, err := r.db.QueryContext(ctx, query, partyID.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contacts := make([]*domain.ContactDetails, 0)
	for rows.Next() {
		var id, typeDesc string
		var phone, email, related sql.NullString
		if err := rows.Scan(&id, &typeDesc, &phone, &email, &related); err != nil {
			return nil, err
		}

		contactID, err := domain.NewContactDetailsID(id)
		if err != nil {
			return nil, err
		}

		var phoneVO *domain.Phone
		if phone.Valid {
			phoneVO, err = domain.NewPhone(phone.String)
			if err != nil {
				return nil, err
			}
		}

		var emailVO *domain.Email
		if email.Valid {
			emailVO, err = domain.NewEmail(email.String)
			if err != nil {
				return nil, err
			}
		}

		var relatedID *domain.PartyID
		if related.Valid {
			pid, err := domain.NewPartyID(related.String)
			if err != nil {
				return nil, err
			}
			relatedID = &pid
		}

		contact, err := domain.NewContactDetails(contactID, typeDesc, phoneVO, emailVO, relatedID)
		if err != nil {
			return nil, err
		}

		contacts = append(contacts, contact)
	}

	return contacts, rows.Err()
}

// PostgreSQLPartyRelationshipRepository implements PartyRelationshipRepository

type PostgreSQLPartyRelationshipRepository struct {
	db *sql.DB
}

func NewPostgreSQLPartyRelationshipRepository(db *sql.DB) *PostgreSQLPartyRelationshipRepository {
	return &PostgreSQLPartyRelationshipRepository{db: db}
}

func (r *PostgreSQLPartyRelationshipRepository) Save(ctx context.Context, relationship domain.PartyRelationship) error {
	query := `
		INSERT INTO party_relationships (id, from_party_id, to_party_id, type, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			from_party_id = $2,
			to_party_id = $3,
			type = $4
	`
	_, err := r.db.ExecContext(ctx, query,
		relationship.ID().Value(),
		relationship.FromID().Value(),
		relationship.ToID().Value(),
		string(relationship.Type()),
		time.Now(),
	)
	return err
}

func (r *PostgreSQLPartyRelationshipRepository) FindByPartyID(ctx context.Context, partyID domain.PartyID) ([]domain.PartyRelationship, error) {
	query := `
		SELECT id, from_party_id, to_party_id, type
		FROM party_relationships
		WHERE from_party_id = $1 OR to_party_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, partyID.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	relationships := make([]domain.PartyRelationship, 0)
	for rows.Next() {
		var id, fromID, toID, relType string
		if err := rows.Scan(&id, &fromID, &toID, &relType); err != nil {
			return nil, err
		}

		relID, err := domain.NewPartyRelationshipID(id)
		if err != nil {
			return nil, err
		}
		fromPartyID, err := domain.NewPartyID(fromID)
		if err != nil {
			return nil, err
		}
		toPartyID, err := domain.NewPartyID(toID)
		if err != nil {
			return nil, err
		}

		typeEnum := domain.RelationshipType(relType)
		if !typeEnum.IsValid() {
			return nil, fmt.Errorf("invalid relationship type in storage: %s", relType)
		}

		relationship, err := domain.NewPartyRelationship(relID, fromPartyID, toPartyID, typeEnum)
		if err != nil {
			return nil, err
		}
		relationships = append(relationships, relationship)
	}

	return relationships, rows.Err()
}

func (r *PostgreSQLPartyRelationshipRepository) Delete(ctx context.Context, id domain.PartyRelationshipID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM party_relationships WHERE id = $1", id.Value())
	return err
}

// PostgreSQLPartyAddressRepository implements PartyAddressRepository

type PostgreSQLPartyAddressRepository struct {
	db *sql.DB
}

func NewPostgreSQLPartyAddressRepository(db *sql.DB) *PostgreSQLPartyAddressRepository {
	return &PostgreSQLPartyAddressRepository{db: db}
}

func (r *PostgreSQLPartyAddressRepository) Save(ctx context.Context, address *domain.Address, addressID domain.AddressID, partyID domain.PartyID, createdBy string, modifiedBy string) error {
	query := `
		INSERT INTO party_addresses (id, party_id, street, city, province, postal_code, country, is_primary, created_by, created_at, modified_by, modified_at)
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
		addressID.Value(),
		partyID.Value(),
		address.Street(),
		address.City(),
		address.Province(),
		address.PostalCode(),
		address.Country(),
		false,
		createdBy,
		now,
		modifiedBy,
		now,
	)
	return err
}

func (r *PostgreSQLPartyAddressRepository) FindByPartyID(ctx context.Context, partyID domain.PartyID) ([]*domain.Address, error) {
	query := `
		SELECT street, city, province, postal_code, country
		FROM party_addresses
		WHERE party_id = $1
		ORDER BY created_at
	`

	rows, err := r.db.QueryContext(ctx, query, partyID.Value())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	addresses := make([]*domain.Address, 0)
	for rows.Next() {
		var street, city, province, postalCode, country string
		if err := rows.Scan(&street, &city, &province, &postalCode, &country); err != nil {
			return nil, err
		}
		address, err := domain.NewAddress(street, city, province, postalCode, country)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	return addresses, rows.Err()
}

func (r *PostgreSQLPartyAddressRepository) FindPrimary(ctx context.Context, partyID domain.PartyID) (*domain.Address, error) {
	query := `
		SELECT street, city, province, postal_code, country
		FROM party_addresses
		WHERE party_id = $1 AND is_primary = true
		LIMIT 1
	`

	var street, city, province, postalCode, country string
	err := r.db.QueryRowContext(ctx, query, partyID.Value()).Scan(&street, &city, &province, &postalCode, &country)
	if err == sql.ErrNoRows {
		return nil, errors.New("primary address not found")
	}
	if err != nil {
		return nil, err
	}

	return domain.NewAddress(street, city, province, postalCode, country)
}

func (r *PostgreSQLPartyAddressRepository) Delete(ctx context.Context, id domain.AddressID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM party_addresses WHERE id = $1", id.Value())
	return err
}
