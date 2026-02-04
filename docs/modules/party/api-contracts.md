# Party Module API Contracts (ADR-012)

This document specifies the API contracts for the Party module aligned with ADR-012: **Party + Profiles + Roles + Relationships + ContactDetails**.

---

## 1. Parties

### 1.1 Create Party
- **Endpoint:** `POST /parties`
- **Description:** Creates a new party (person, organization, or both profiles).
- **Request Body:**
  ```json
  {
    "status": "ACTIVE",
    "roles": ["CLIENT", "SUPPLIER"],
    "person_profile": {
      "first_name": "Ana",
      "last_name": "Pérez"
    },
    "organization_profile": {
      "name": "Textiles Pérez S.L.",
      "tax_id": "B12345678",
      "website": "https://textilesperez.com",
      "contacts": [
        {
          "type_description": "Ventas",
          "phone": "+34 600 123 456",
          "email": "ventas@textilesperez.com",
          "related_party_id": null
        }
      ]
    }
  }
  ```
- **Success Response (201 Created):** `PartyDTO`

### 1.2 List Parties
- **Endpoint:** `GET /parties`
- **Description:** Lists parties with filtering and pagination.
- **Query Parameters:**
  - `role` (string, optional): "CLIENT", "SUPPLIER", "EMPLOYEE", etc.
  - `type` (string, optional): "PERSON", "ORGANIZATION", "BOTH"
  - `status` (string, optional): "ACTIVE", "INACTIVE"
  - `name` (string, optional): matches person/organization name
  - `tax_id` (string, optional)
  - `page` (int, optional): default 1
  - `page_size` (int, optional): default 10
- **Success Response (200 OK):** `ListResponse<PartyDTO>`

### 1.3 Get Party
- **Endpoint:** `GET /parties/{id}`
- **Description:** Retrieves a single party by its ID.
- **Success Response (200 OK):** `PartyDTO`

### 1.4 Update Party
- **Endpoint:** `PUT /parties/{id}`
- **Description:** Updates party profiles and status.
- **Request Body:**
  ```json
  {
    "status": "ACTIVE",
    "person_profile": {
      "first_name": "Ana",
      "last_name": "Pérez"
    },
    "organization_profile": {
      "name": "Textiles Pérez S.L.",
      "tax_id": "B12345678",
      "website": "https://textilesperez.com"
    }
  }
  ```
- **Success Response (200 OK):** `PartyDTO`

### 1.5 Change Party Status
- **Endpoint:** `PATCH /parties/{id}/status`
- **Description:** Changes the party status (e.g., ACTIVE → INACTIVE).
- **Request Body:**
  ```json
  {
    "status": "INACTIVE"
  }
  ```
- **Success Response (200 OK):** `PartyDTO`

---

## 2. Party Roles

### 2.1 Add Role
- **Endpoint:** `POST /parties/{id}/roles`
- **Description:** Adds a role to a party.
- **Request Body:**
  ```json
  {
    "role": "CLIENT"
  }
  ```
- **Success Response (200 OK):** `PartyDTO`

### 2.2 Remove Role
- **Endpoint:** `DELETE /parties/{id}/roles/{role}`
- **Description:** Removes a role from a party.
- **Success Response (200 OK):** `PartyDTO`

---

## 3. Party Relationships

### 3.1 Add Relationship
- **Endpoint:** `POST /parties/{id}/relationships`
- **Description:** Creates a relationship between parties.
- **Request Body:**
  ```json
  {
    "to_party_id": "party_456",
    "type": "IS_EMPLOYEE_OF"
  }
  ```
- **Success Response (201 Created):** `PartyRelationshipDTO`

### 3.2 List Relationships
- **Endpoint:** `GET /parties/{id}/relationships`
- **Description:** Lists all relationships for a party.
- **Success Response (200 OK):** `ListResponse<PartyRelationshipDTO>`

### 3.3 Remove Relationship
- **Endpoint:** `DELETE /parties/{id}/relationships/{relationship_id}`
- **Description:** Deletes a relationship.
- **Success Response (204 No Content)**

---

## 4. Contact Details (Organization Profile)

### 4.1 Add Contact Detail
- **Endpoint:** `POST /parties/{id}/contact-details`
- **Description:** Adds a contact detail to the organization profile.
- **Request Body:**
  ```json
  {
    "type_description": "Ventas",
    "phone": "+34 600 123 456",
    "email": "ventas@textilesperez.com",
    "related_party_id": null
  }
  ```
- **Success Response (201 Created):** `ContactDetailsDTO`

### 4.2 List Contact Details
- **Endpoint:** `GET /parties/{id}/contact-details`
- **Description:** Lists contact details for the organization profile.
- **Success Response (200 OK):** `ListResponse<ContactDetailsDTO>`

### 4.3 Update Contact Detail
- **Endpoint:** `PUT /parties/{id}/contact-details/{contact_id}`
- **Success Response (200 OK):** `ContactDetailsDTO`

### 4.4 Remove Contact Detail
- **Endpoint:** `DELETE /parties/{id}/contact-details/{contact_id}`
- **Success Response (204 No Content)**

---

## DTOs (Resumen)

```json
PartyDTO {
  "id": "party_123",
  "status": "ACTIVE",
  "roles": ["CLIENT"],
  "person_profile": {
    "first_name": "Ana",
    "last_name": "Pérez"
  },
  "organization_profile": {
    "name": "Textiles Pérez S.L.",
    "tax_id": "B12345678",
    "website": "https://textilesperez.com",
    "contacts": [ContactDetailsDTO]
  },
  "created_at": "2026-02-02T00:00:00Z",
  "modified_at": "2026-02-02T00:00:00Z"
}
```
