# Modelo de Dominio - Módulo Party

Este documento describe el modelo de dominio vigente para el módulo `Party`, alineado con ADR-012 y la implementación actual. El agregado principal es `Party`, que puede representar una persona, una organización o ambas (perfiles compuestos).

## 1. Agregado: Party

`Party` es la raíz del agregado. Centraliza el estado, los perfiles y las relaciones, y garantiza la consistencia transaccional.

### Atributos principales
- **ID**: `PartyID`
- **Status**: `PartyStatus` (enum: `ACTIVE`, `INACTIVE`)
- **DefaultDiscountPercentage**: `float64` (0-100, descuento por defecto cuando actúa como cliente, usado por el módulo Pricing como fallback)
- **PersonProfile**: `*PersonProfile` (opcional)
- **OrganizationProfile**: `*OrganizationProfile` (opcional)
- **Roles**: `[]PartyRole`
- **Relationships**: `[]PartyRelationship`
- **CreatedBy / CreatedAt / ModifiedBy / ModifiedAt**

### Reglas clave
- Un `Party` debe tener **al menos un perfil** (persona u organización).
- Puede tener múltiples roles (cliente, proveedor, empleado).
- Puede relacionarse con otros `Party` mediante relaciones tipadas.

## 2. Perfil de Persona (`PersonProfile`)

Representa los datos personales de un `Party`.

- **FirstName**: `string`
- **LastName**: `string`

## 3. Perfil de Organización (`OrganizationProfile`)

Representa los datos corporativos de un `Party`.

- **Name**: `string`
- **TaxID**: `*TaxID`
- **Website**: `string`
- **Contacts**: `[]ContactDetails`

## 4. ContactDetails

Detalle de contacto asociado a una organización.

- **ID**: `ContactDetailsID`
- **TypeDescription**: `string` (ej. “Ventas”, “Soporte”)
- **Phone**: `*Phone`
- **Email**: `*Email`
- **RelatedPartyID**: `*PartyID` (opcional, para enlazar a otra Party)

## 5. Roles y Relaciones

### PartyRole
- **Type**: `PartyRoleType` (`CLIENT`, `SUPPLIER`, `EMPLOYEE`)

### PartyRelationship
- **FromID**: `PartyID`
- **ToID**: `PartyID`
- **Type**: `RelationshipType` (`IS_EMPLOYEE_OF`, `IS_SUBSIDIARY_OF`)

## 6. Value Objects y Enumeraciones

**Value Objects:** `PartyID`, `ContactDetailsID`, `Email`, `Phone`, `TaxID`.

**Enumeraciones:** `PartyStatus`, `PartyRoleType`, `RelationshipType`.

---
**Última Actualización:** 2026-03-01