# Casos de Uso - Módulo Party

**Estado:** Definido y Documentado

Esta sección documenta los casos de uso vigentes del módulo Party, alineados con ADR-012 y los endpoints actuales.

## 1. Casos de Uso de Party

- **Crear Party:** Registrar una nueva Party con perfil de persona, organización o ambos.
- **Listar Parties:** Lista paginada con filtros por rol, tipo, estado, nombre y tax_id.
- **Obtener Party:** Recuperar una Party por su ID.
- **Actualizar Party:** Actualizar perfiles (persona/organización) y datos básicos.
- **Cambiar Estado:** Activar o desactivar una Party.

## 2. Casos de Uso de Roles

- **Añadir Rol a Party:** Agregar un rol (`CLIENT`, `SUPPLIER`, `EMPLOYEE`).
- **Eliminar Rol de Party:** Quitar un rol existente de la Party.

## 3. Casos de Uso de Relaciones

- **Crear Relación:** Enlazar dos Parties con un tipo (`IS_EMPLOYEE_OF`, `IS_SUBSIDIARY_OF`).
- **Listar Relaciones:** Obtener relaciones asociadas a una Party.
- **Eliminar Relación:** Eliminar una relación existente.

## 4. Casos de Uso de ContactDetails (Perfil Organización)

- **Añadir Contacto:** Agregar detalle de contacto a la organización.
- **Listar Contactos:** Recuperar contactos asociados a la organización.
- **Actualizar Contacto:** Modificar un contacto existente.
- **Eliminar Contacto:** Eliminar un contacto existente.

---
**Última Actualización:** 2026-02-03
