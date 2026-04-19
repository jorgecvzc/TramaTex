# 🏛️ ADR-005: Gestión Unificada de Terceros (Patrón Party)

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.0 |
| **Estado** | ✅ Aceptado |
| **Fecha** | 09-01-2026 |
| **Autores** | Jorge Cortés Villalba, ChatGPT |

---

## 🎯 Contexto
Los módulos de Compras y Ventas necesitan interactuar con entidades externas (clientes y proveedores). Mantener entidades separadas genera duplicidad de datos, inconsistencias en el histórico de transacciones y dificultades para gestionar empresas que actúan en ambos roles simultáneamente.

---

## 🔍 Alternativas Consideradas
1. **Entidades Separadas:** Simplicidad inicial, pero con alta redundancia de datos fiscales y problemas de sincronización.
2. **Entidad Base Unificada (Decisión Adoptada):** Uso del **Patrón Party**. Una única entidad base para Personas y Organizaciones que puede asumir múltiples roles dinámicos.

---

## ✅ Decisión Adoptada
Se adopta la creación de una entidad base unificada **`Party`**.

### Principios del Modelo:
*   **Entidad Única:** Evita registrar la misma empresa dos veces (ej: una como cliente y otra como proveedor).
*   **Roles Dinámicos:** Una `Party` puede ser Cliente, Proveedor o ambos sin duplicar datos de contacto o fiscales.
*   **Jerarquías:** Soporte nativo para relaciones entre matrices y filiales.
*   **Histórico Consistente:** Visión 360º de todas las transacciones vinculadas a un tercero.

---

## 📈 Consecuencias
### Positivas
*   Eliminación total de la redundancia de datos de terceros.
*   Máxima integridad referencial en transacciones y precios.
*   Facilidad de mantenimiento y escalabilidad del maestro de datos.

### Negativas
*   Requiere lógica adicional en la aplicación para el filtrado por rol.
*   La interfaz de usuario debe manejar la visualización condicional de campos específicos.

---
[Volver al Índice de ADRs](./README.md)
