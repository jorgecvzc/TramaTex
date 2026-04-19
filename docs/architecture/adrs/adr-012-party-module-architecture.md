# 🏛️ ADR-012: Arquitectura del Módulo Party

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.0 |
| **Estado** | ✅ Aceptado |
| **Fecha** | 01-02-2026 |
| **Autores** | Gemini CLI, Usuario |

---

## 🎯 Contexto
El módulo `Party` debe gestionar clientes y proveedores con relaciones complejas (matrices/filiales, personas físicas vs. jurídicas, y empleados que también son clientes). El modelo inicial era insuficiente para representar estas jerarquías sin generar redundancia.

---

## 🔍 Alternativas Consideradas
1. **Mantener Estructura Actual:** Incapaz de representar personas individuales de forma natural.
2. **Modelo Abstracto de Perfil Único:** Simple, pero propenso a inconsistencias en casos híbridos.
3. **Modelo de Roles y Relaciones (Decisión Adoptada):** Resuelve la complejidad de forma nativa permitiendo que una entidad asuma múltiples identidades de negocio.

---

## ✅ Decisión Adoptada
Se adopta el **Modelo de Party con Roles y Relaciones** y un manejo de contactos simplificado.

### Claves del Diseño:
*   **Separación de Identidad y Rol:** Una `Party` es la entidad (empresa o persona), mientras que el `Role` define su comportamiento (Cliente, Proveedor).
*   **Soporte de Jerarquías:** Las relaciones permiten vincular matrices con filiales de forma explícita.
*   **Pragmatismo en Contactos:** Se utiliza un Value Object para los detalles de contacto en lugar de una gestión excesivamente granular, priorizando la agilidad para el MVP.

---

## 📈 Consecuencias
### Positivas
*   Soporte nativo para la complejidad real de las relaciones industriales textiles.
*   Modelo extensible y preparado para evoluciones futuras de CRM.
*   Eliminación de inconsistencias en perfiles híbridos.

### Negativas
*   Refactorización completa del módulo existente.
*   Consultas a base de datos ligeramente más complejas.

---
[Volver al Índice de ADRs](./README.md)
