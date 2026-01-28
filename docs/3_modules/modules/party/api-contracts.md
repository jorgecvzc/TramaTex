# Party Module API Contracts

This document specifies the API contracts for the Party module, which handles organizations, persons (contacts), and addresses.

---

## 1. Organizations

### 1.1 Create Organization
- **Endpoint:** `POST /organizations`
- **Description:** Creates a new organization (client or supplier).
- **Request Body:**
  ```json
  {
    "name": "New Company S.L.",
    "role": "CLIENT", // or "SUPPLIER" or "BOTH"
    "tax_id": "B12345678",
    "website": "https://newcompany.com",
    "notes": "Initial notes."
  }
  ```
- **Success Response (201 Created):** OrganizationDTO

### 1.2 List Organizations
- **Endpoint:** `GET /organizations`
- **Description:** Lists organizations with filtering and pagination.
- **Query Parameters:**
  - `role` (string, optional): "CLIENT", "SUPPLIER", "BOTH"
  - `status` (string, optional): "ACTIVE", "INACTIVE"
  - `page` (int, optional): default 1
  - `page_size` (int, optional): default 10
- **Success Response (200 OK):** `ListResponse<OrganizationDTO>`

### 1.3 Get Organization
- **Endpoint:** `GET /organizations/{id}`
- **Description:** Retrieves a single organization by its ID.
- **Success Response (200 OK):** `OrganizationDTO`

### 1.4 Update Organization
- **Endpoint:** `PUT /organizations/{id}`
- **Description:** Updates an existing organization's details.
- **Request Body:**
  ```json
  {
    "name": "Updated Company Name S.L.",
    "website": "https://updated-company.com",
    "notes": "Updated notes."
  }
  ```
- **Success Response (200 OK):** `OrganizationDTO`

### 1.5 Change Organization Status
- **Endpoint:** `PATCH /organizations/{id}/status`
- **Description:** Changes the status of an organization (e.g., from ACTIVE to INACTIVE).
- **Request Body:**
  ```json
  {
    "status": "INACTIVE"
  }
  ```
- **Success Response (200 OK):** `OrganizationDTO`

---

## 2. Persons (Contacts)

### 2.1 Add Person to Organization
- **Endpoint:** `POST /organizations/{org_id}/persons`
- **Description:** Adds a new contact person to an organization.
- **Request Body:**
  ```json
  {
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "phone": "+34600123456",
    "job_title": "Purchasing Manager",
    "is_primary_contact": true
  }
  ```
- **Success Response (201 Created):** `PersonDTO`

### 2.2 Get Person
- **Endpoint:** `GET /persons/{id}`
- **Description:** Retrieves a single person by their ID.
- **Success Response (200 OK):** `PersonDTO`

### 2.3 List Persons for Organization
- **Endpoint:** `GET /organizations/{org_id}/persons`
- **Description:** Lists all contact persons for a given organization.
- **Success Response (200 OK):** `ListResponse<PersonDTO>`

### 2.4 Get Primary Contact
- **Endpoint:** `GET /organizations/{org_id}/primary-contact`
- **Description:** Retrieves the designated primary contact for an organization.
- **Success Response (200 OK):** `PersonDTO`

---

## 3. Addresses

### 3.1 Add Address to Organization
- **Endpoint:** `POST /organizations/{org_id}/addresses`
- **Description:** Adds a new address to an organization.
- **Request Body:**
  ```json
  {
    "street": "123 Main St",
    "city": "Anytown",
    "province": "Anyprovince",
    "postal_code": "12345",
    "country": "Spain",
    "is_primary": true
  }
  ```
- **Success Response (201 Created):** `AddressDTO`

### 3.2 List Addresses for Organization
- **Endpoint:** `GET /organizations/{org_id}/addresses`
- **Description:** Lists all addresses for a given organization.
- **Success Response (200 OK):** `ListResponse<AddressDTO>`

### 3.3 Get Primary Address
- **Endpoint:** `GET /organizations/{org_id}/primary-address`
- **Description:** Retrieves the designated primary address for an organization.
- **Success Response (200 OK):** `AddressDTO`
