# Modelo de Dominio - Módulo Party

### Party (Raíz de Agregación)
- **ID**: UUID
- **Tipo**: Enum (Cliente, Proveedor, Ambos)
- **Nombre**: String (nombre legal)
- **Categoría**: Enum (Mayorista, Minorista, Distribuidor, etc.)
- **Estado**: Enum (Activo, Inactivo, Suspendido)
- **Contactos**: List<Contact>
- **Direcciones**: List<Address>
- **Metadata**: CreatedAt, UpdatedAt
