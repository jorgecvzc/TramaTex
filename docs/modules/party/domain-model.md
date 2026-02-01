# Modelo de Dominio - Módulo Party

Este documento describe el modelo de dominio para el módulo `Party`, que gestiona clientes y proveedores. Se adhiere a los principios de Domain-Driven Design (DDD), definiendo la `Organization` como la raíz del agregado, `Person` como una entidad dentro de ese agregado y `Address` como un objeto de valor.

## 1. Agregado: Organization

La `Organization` es la raíz del agregado principal en el módulo `Party`. Es el punto de entrada para todas las operaciones que involucran a `Person` y `Address`, asegurando la consistencia transaccional.

### Atributos de Organization:
- **ID**: `OrganizationID` (identificador único del agregado)
- **Nombre**: `string` (nombre legal de la organización)
- **Rol**: `OrganizationRole` (enum: Cliente, Proveedor, Ambos)
- **Estado**: `OrganizationStatus` (enum: Activo, Inactivo, Suspendido)
- **TaxID**: `*TaxID` (Número de Identificación Fiscal, ej. CIF, NIF)
- **Website**: `string` (URL del sitio web de la organización)
- **Notes**: `string` (notas internas)
- **CreatedBy**: `string` (ID del usuario que creó la organización)
- **CreatedAt**: `time.Time` (marca de tiempo de creación)
- **ModifiedBy**: `string` (ID del último usuario que modificó la organización)
- **ModifiedAt**: `time.Time` (marca de tiempo de la última modificación)

### Comportamientos Clave:
- Crear una nueva `Organization` con validaciones iniciales.
- Actualizar nombre, website o notas.
- Activar o desactivar la organización.
- Añadir `Person` (persona de contacto).
- Añadir `Address` (dirección).

## 2. Entidad: Person

La `Person` es una entidad con su propia identidad, pero su ciclo de vida y consistencia son gestionados por la `Organization` a la que pertenece.

### Atributos de Person:
- **ID**: `PersonID` (identificador único de la persona)
- **OrganizationID**: `OrganizationID` (referencia a la organización a la que pertenece)
- **FirstName**: `string` (nombre de pila)
- **LastName**: `string` (apellido)
- **Email**: `*Email` (objeto de valor para el email)
- **Phone**: `*Phone` (objeto de valor para el teléfono)
- **JobTitle**: `string` (cargo o puesto de trabajo)
- **IsPrimaryContact**: `bool` (indica si es el contacto principal de la organización)
- **CreatedBy**: `string`
- **CreatedAt**: `time.Time`
- **ModifiedBy**: `string`
- **ModifiedAt**: `time.Time`

## 3. Value Object: Address

`Address` es un objeto de valor que describe una ubicación física. No tiene una identidad conceptual propia dentro del dominio y se define por sus atributos.

### Atributos de Address:
- **Street**: `string` (calle y número)
- **City**: `string` (ciudad)
- **Province**: `string` (provincia)
- **PostalCode**: `string` (código postal)
- **Country**: `string` (país)

## 4. Enumeraciones

### OrganizationRole
- `CLIENT`
- `SUPPLIER`
- `BOTH`

### OrganizationStatus
- `ACTIVE`
- `INACTIVE`
- `SUSPENDED`

---
**Última Actualización:** 2026-02-01