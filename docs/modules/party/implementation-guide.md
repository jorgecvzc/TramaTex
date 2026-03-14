# Guía de Integración y Estándares: Módulo Party

Este documento establece las directrices para interactuar con el módulo Party desde otros contextos delimitados y los estándares de integridad de datos que deben respetarse.

---

## 1. Principios de Interacción

Cualquier módulo que necesite referenciar a un tercero (cliente, proveedor u operario) debe seguir estas reglas:

- **Referencia por ID:** Almacenar únicamente el `PartyID` (UUID). Nunca duplicar nombres o datos fiscales en las tablas de otros módulos.
- **Consulta de Metadatos:** Utilizar el servicio de aplicación de Party para resolver la información necesaria para visualización o facturación.
- **Respeto al Estado Operativo:** Antes de permitir una operación (ej. crear un pedido), el módulo llamador debe verificar que la Party no esté en estado `BLOCKED`.

---

## 2. Estándares de Integridad de Datos

### Soberanía del Dato
El módulo Party es el único propietario de la información de contacto y fiscal. Ningún otro módulo debe modificar estos datos directamente. Cualquier actualización debe fluir a través de los casos de uso definidos en Party.

### Gestión de Auditoría y Persistencia
- **Exclusión de Campos Técnicos:** Las entidades de dominio no deben contener campos de infraestructura (`CreatedAt`, `UpdatedAt`). Estos concerns se resuelven exclusivamente en la capa de persistencia para mantener el modelo de negocio puro y desacoplado.
- **Inmutabilidad de Identidad:** El `PartyID` es inmutable una vez creado. Cualquier cambio de naturaleza legal (ej. transformación de persona física a sociedad) debe gestionarse mediante la adición de un nuevo perfil a la misma Party o, en casos extremos, la creación de una nueva identidad con el traspaso de relaciones correspondiente.

---

## 3. Estándares Económicos y de Comunicación

### Moneda Única
En cumplimiento con el mandato del proyecto, todas las operaciones económicas vinculadas a la Party (ej. límites de crédito, descuentos base) deben gestionarse exclusivamente en **Euros (€)**.

### Mapeo de Errores de Negocio
Las integraciones deben estar preparadas para manejar los errores semánticos del módulo:
- **Validación Fallida:** Formatos de TaxID incorrectos o perfiles incompletos.
- **Conflicto de Identidad:** Intento de duplicar un TaxID ya existente.
- **Restricción de Borrado:** Intento de eliminar una entidad con actividad histórica activa.

---
**Última Actualización:** 2026-03-07
