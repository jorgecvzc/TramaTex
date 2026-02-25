# ADR-020 – Tickets (Facturas Simplificadas) y Series de Numeración

**Fecha:** 14 de febrero de 2026  
**Estado:** Aceptado  
**Autores:** Claude AI (Anthropic) + Joran (Product Owner)

---

## 1. Contexto

El sistema TramaTex necesita soportar **ventas al por menor** (retail/mostrador) además del flujo B2B tradicional. La legislación española (Real Decreto 1619/2012) establece que:

- **Ticket = Factura simplificada** para operaciones < 3.000€
- **Series de numeración** diferenciadas son obligatorias por AEAT
- Los tickets deben integrarse en el sistema contable junto con facturas completas

---

## 2. Alternativas Consideradas

**Alternativa A – Entidad Separada `Ticket`**
- Ventajas: Separación total de concerns (B2B vs retail).
- Desventajas: Mayor complejidad arquitectónica, duplicación de lógica.

**Alternativa B – Extensión de `Invoice` con `InvoiceType` (Adoptada)**
- Ventajas: Reutiliza lógica existente, mínimo impacto arquitectónico, cumple legislación.
- Desventajas: `Invoice` asume más responsabilidades.

**Alternativa C – Posponer a Post-MVP**
- Desventajas: Sistema incompleto para negocio real.

---

## 3. Criterios de Decisión

- Cumplimiento Legal (AEAT).
- Time-to-Market (Deadline 23 febrero).
- Reutilización de código existente de Sales.
- Mantenibilidad y extensibilidad para TPV futuro.

---

## 4. Decisión Adoptada

Se adopta la **Alternativa B: Extensión de `Invoice` con `InvoiceType`** con implementación híbrida MVP + Post-MVP.

### Implementación MVP:
- Nuevos Value Objects: `InvoiceType` y `InvoiceSeries`.
- Modificación de entidad `Invoice` para incluir tipo y serie.
- Validaciones específicas por tipo (límite 3.000€ para simplificadas).
- Party genérico "CONSUMIDOR_FINAL".

---

## 5. Consecuencias

### Positivas
- Cumplimiento legal desde el primer día.
- Base arquitectónica preparada para TPV completo.
- Trazabilidad unificada de todas las facturas.

### Negativas
- Mayor complejidad en la lógica del agregado `Invoice`.
- UX inicial para retail requiere crear `SalesOrder` (resuelto en Post-MVP).

---

## 6. Alcance

Aplica al módulo `Sales` (entidad `Invoice`, repositorios, casos de uso), base de datos (tabla `invoice_series`), frontend y la integración con el módulo `Party`.

---

## 7. Integración con otros ADRs

- **ADR-017:** Arquitectura del Módulo de Sales (base sobre la que se construye esta extensión).
- **ADR-005:** Gestión Unificada de Clientes (uso de Party genérico).

---

## 8. Notas Adicionales / Consideraciones Especiales

### Métricas de Éxito (MVP)
- Emisión de tickets en menos de 5 clicks.
- 100% de cumplimiento legal en formato.
- Cero duplicados en series de numeración.

### Plan de Migración Post-MVP
Se implementará una UI de TPV dedicada con diseño tipo caja registradora, gestión de métodos de pago y lector de códigos de barras. Si el volumen retail supera el 60%, se evaluará extraer `Ticket` a una entidad separada.

### Aprobación y Revisión
- **Aprobado por:** Joran (Product Owner)
- **Fecha de aprobación:** 14 de febrero de 2026.
- **Revisión próxima:** 3 meses después del MVP.

---

## 9. Referencias

- Real Decreto 1619/2012 (España)
- Normativa AEAT sobre series de numeración
- ADR-017: Arquitectura del Módulo de Sales
