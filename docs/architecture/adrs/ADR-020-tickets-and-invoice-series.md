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

**Problemas a Resolver:**
- El módulo Sales actual (ADR-017) solo contempla flujo B2B: Quote → Order → DeliveryNote → Invoice (completa)
- No existe capacidad para emitir tickets/facturas simplificadas
- No hay sistema de series de numeración diferenciadas
- UI TPV (Punto de Venta) no está contemplada en el diseño

**Restricciones:**
- **Deadline crítico:** 23 de febrero de 2026 (9 días restantes)
- Módulos Sales + MES deben estar operativos para MVP
- No se puede comprometer la entrega con funcionalidades complejas

---

## 2. Alternativas Consideradas

### **Alternativa A – Entidad Separada `Ticket`**

**Descripción:** Crear una entidad de dominio `Ticket` completamente independiente de `Invoice`, con su propio agregado, repositorio y flujo.

**Ventajas:**
- Separación total de concerns (B2B vs retail)
- Modelo de dominio más puro
- UI TPV completamente desacoplada

**Desventajas:**
- ❌ Mayor complejidad arquitectónica
- ❌ Duplicación de lógica (LineItems, cálculos, impuestos)
- ❌ +5-7 días de desarrollo
- ❌ Riesgo para deadline del MVP

---

### **Alternativa B – Extensión de `Invoice` con `InvoiceType` (Adoptada para MVP)**

**Descripción:** Extender la entidad `Invoice` existente añadiendo:
- `InvoiceType` enum: `COMPLETA` | `SIMPLIFICADA`
- `InvoiceSeries` Value Object para gestionar series de numeración
- Validaciones específicas según el tipo de factura

**Ventajas:**
- ✅ Reutiliza lógica existente de `Invoice` y `InvoiceLineItem`
- ✅ Mínimo impacto arquitectónico (extensión, no nueva entidad)
- ✅ Cumple legislación española
- ✅ Series de numeración fácilmente configurables
- ✅ Tiempo razonable: +1 día desarrollo

**Desventajas:**
- ⚠️ `Invoice` tiene más responsabilidad (2 tipos con lógica diferente)
- ⚠️ Modelo no tan "puro" como entidad separada

---

### **Alternativa C – Posponer TODO a Post-MVP**

**Descripción:** No implementar tickets ni TPV en MVP. Solo flujo B2B tradicional.

**Ventajas:**
- ✅ 0 días extra de desarrollo
- ✅ Mayor margen para cumplir deadline

**Desventajas:**
- ❌ Sistema incompleto para negocio real (TramaTex tiene ventas retail)
- ❌ Refactoring complejo en el futuro si se implementa mal

---

## 3. Criterios de Decisión

- **Cumplimiento Legal:** Debe cumplir con normativa AEAT española
- **Time-to-Market:** No comprometer deadline del 23 de febrero
- **Extensibilidad:** Facilitar implementación completa de TPV en Post-MVP
- **Reutilización de Código:** Aprovechar lógica existente de Sales
- **Mantenibilidad:** Evitar duplicación innecesaria de código

---

## 4. Decisión Adoptada

Se adopta la **Alternativa B: Extensión de `Invoice` con `InvoiceType`** con implementación **híbrida MVP + Post-MVP**.

### **Fase 1: MVP (AHORA - antes del 23 de febrero)**

#### **Backend:**
1. **Nuevo Value Object: `InvoiceType`**
   ```go
   type InvoiceType string
   
   const (
       InvoiceTypeFull       InvoiceType = "COMPLETA"
       InvoiceTypeSimplified InvoiceType = "SIMPLIFICADA"
   )
   ```

2. **Nuevo Value Object: `InvoiceSeries`**
   ```go
   type InvoiceSeries struct {
       Code    string    // "A", "B", "TKT"
       Year    int       // 2026
       Prefix  string    // opcional: "FAC-", "TIK-"
   }
   
   // InvoiceNumber format: {Series}/{Number:05d}/{Year}
   // Examples: "A/00001/2026", "TKT/00123/2026"
   ```

3. **Modificar entidad `Invoice`:**
   - Añadir campo: `Type: InvoiceType`
   - Añadir campo: `Series: InvoiceSeries`
   - Modificar: `InvoiceNumber` ahora usa series

4. **Tabla de configuración de series:**
   ```sql
   CREATE TABLE invoice_series (
       id UUID PRIMARY KEY,
       code VARCHAR(10) NOT NULL,
       year INT NOT NULL,
       current_number INT NOT NULL DEFAULT 0,
       invoice_type VARCHAR(20) NOT NULL, -- COMPLETA | SIMPLIFICADA
       prefix VARCHAR(10),
       is_active BOOLEAN NOT NULL DEFAULT true,
       UNIQUE(code, year)
   );
   ```

5. **Validaciones por tipo:**
   - **COMPLETA:**
     - Requiere `PartyID` con datos fiscales completos
     - Sin límite de importe
   - **SIMPLIFICADA (Ticket):**
     - `PartyID` puede ser "CONSUMIDOR_FINAL" (Party genérico)
     - Límite legal: `Total < 3.000 EUR`
     - Datos mínimos: NIF emisor, fecha, número, total con IVA

6. **Nuevo Caso de Uso: `CreateSimplifiedInvoice`**
   - Input: `SalesOrderID`, `InvoiceDate`, `DueDate`, `Series` (default: "TKT")
   - Validación: `Total < 3.000 EUR`
   - Output: `Invoice` con `Type = SIMPLIFICADA`

7. **Party genérico "CONSUMIDOR_FINAL":**
   - ID fijo: `00000000-0000-0000-0000-000000000001`
   - NIF: `99999999R` (NIF genérico válido para consumidor final)
   - Nombre: "Cliente de mostrador"

#### **Frontend:**
1. **En vista de SalesOrder Detail:**
   - Botón existente: "Emitir Factura" (genera Invoice COMPLETA)
   - **Nuevo botón:** "Emitir Ticket" (genera Invoice SIMPLIFICADA)

2. **Modal de confirmación simple:**
   - Título: "Emitir Ticket de Venta"
   - Muestra: Total del pedido
   - Validación: Si total > 3.000€ → error + sugerencia de factura completa
   - Botón: "Confirmar Emisión"

3. **Vista de ticket (PDF básico):**
   - Formato simplificado según legislación
   - Datos mínimos: Emisor, fecha, número, líneas, total IVA incluido

#### **Flujo MVP:**
```
Usuario crea/selecciona SalesOrder 
  → Click "Emitir Ticket"
  → Sistema valida (total < 3.000€)
  → Genera Invoice (type=SIMPLIFICADA, series="TKT")
  → Muestra PDF del ticket
```

---

### **Fase 2: Post-MVP (DESPUÉS del 23 de febrero)**

#### **UI TPV Dedicada:**
1. **Interfaz tipo POS:**
   - Diseño tipo caja registradora
   - Búsqueda rápida de productos (nombre, SKU, código de barras)
   - Carrito en pantalla con edición inline
   - Cálculo automático de totales

2. **Gestión de métodos de pago:**
   - Efectivo (con cálculo de cambio)
   - Tarjeta (integración con TPV físico)
   - Pago mixto (efectivo + tarjeta)

3. **Lector de códigos de barras:**
   - Integración USB con lectores estándar
   - Añadir productos al carrito escaneando

4. **Gestión de caja:**
   - Apertura de caja (fondo inicial)
   - Arqueo de caja
   - Cierre de caja (Z report)
   - Registro de movimientos (retiradas, ingresos)

5. **Informes:**
   - X report (consulta sin cerrar)
   - Z report (cierre de día)
   - Ventas por vendedor
   - Ventas por método de pago

6. **Flujo TPV dedicado:**
```
UI TPV → Añadir productos al carrito
       → Seleccionar método de pago
       → Genera SalesOrder automático (CONFIRMADO)
       → Genera Invoice (SIMPLIFICADA) automático
       → Imprime ticket
       → Abre cajón de efectivo (si aplica)
```

---

## 5. Consecuencias

### **Positivas**
- ✅ **Cumplimiento Legal:** Sistema cumple normativa AEAT desde MVP
- ✅ **Time-to-Market:** No compromete deadline (+1 día vs +5-7 días alternativa A)
- ✅ **Extensibilidad:** Arquitectura preparada para UI TPV completa Post-MVP
- ✅ **Reutilización:** Aprovecha lógica de `Invoice` y `InvoiceLineItem` existente
- ✅ **Flexibilidad:** Series de numeración configurables por año y tipo
- ✅ **Trazabilidad:** Tickets quedan registrados como facturas en el sistema

### **Negativas**
- ⚠️ **Responsabilidad de Invoice:** Entidad maneja 2 tipos con lógicas diferentes
- ⚠️ **UX Inicial:** Para venta de mostrador en MVP, usuario debe crear SalesOrder primero (no es flujo natural de TPV). Se resuelve en Post-MVP con UI dedicada.
- ⚠️ **Validaciones Complejas:** Lógica de validación por tipo añade complejidad al agregado `Invoice`

### **Riesgos y Mitigación**

| Riesgo | Probabilidad | Impacto | Mitigación |
|--------|--------------|---------|------------|
| Validación de límite 3.000€ incorrecta | Baja | Alto | Tests unitarios exhaustivos, validación en backend + frontend |
| Duplicación de números en series | Media | Alto | Transacciones con locks en generación de números |
| Refactoring costoso a entidad `Ticket` separada | Baja | Medio | ADR documenta decisión, se puede mantener arquitectura actual |

---

## 6. Alcance

Esta decisión aplica a:
- **Módulo Sales:** Entidad `Invoice`, repositorios, casos de uso
- **Base de Datos:** Nueva tabla `invoice_series`, modificación de `invoices`
- **Frontend:** Vistas de SalesOrder, modal de tickets, PDF renderer
- **Integración:** Módulo Party (Party genérico CONSUMIDOR_FINAL)

---

## 7. Métricas de Éxito (MVP)

| Métrica | Objetivo |
|---------|----------|
| **Emisión de tickets** | Usuario puede emitir ticket desde SalesOrder en < 5 clicks |
| **Validación legal** | 100% de tickets cumplen formato legislación española |
| **Unicidad de números** | 0 duplicados en series de numeración |
| **Performance** | Generación de ticket < 500ms |
| **Tiempo de desarrollo** | ≤ 1 día adicional al módulo Sales base |

---

## 8. Plan de Migración Post-MVP (Evaluación Futura)

Si el volumen de ventas retail crece significativamente, se puede considerar:

**Opción de refactoring a entidad `Ticket` separada:**
1. Crear entidad `Ticket` con campos simplificados
2. Migrar registros `Invoice` con `Type=SIMPLIFICADA` a tabla `tickets`
3. Mantener `Invoice` solo para facturas completas B2B
4. Actualizar UI TPV para usar entidad `Ticket` directamente

**Umbral de decisión:** Si ventas retail > 60% del total O si lógica de tickets diverge significativamente de facturas completas.

---

## 9. Referencias

- **ADR-017:** Arquitectura del Módulo de Sales (base)
- **Real Decreto 1619/2012:** Obligaciones de facturación (España)
- **Normativa AEAT:** Series de numeración de facturas
- **Ley 58/2003 (LGT):** Obligaciones tributarias
- **Reglamento de Facturación:** Facturas simplificadas (tickets)

---

## 10. Aprobación y Revisión

- **Aprobado por:** Joran (Product Owner)
- **Fecha de aprobación:** 14 de febrero de 2026
- **Revisión:** Post-MVP (después del 23 de febrero de 2026)
- **Próxima revisión:** Al alcanzar 1.000 tickets emitidos o 3 meses después del MVP

---
