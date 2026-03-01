# TramaTex – Bounded Contexts

**Versión:** 2.0  
**Fecha:** 11/01/2026  
**Autor:** Jorge Cortés Villalba  
**Copilotos LLM:** Claude (Anthropic)  
**Propósito:** Definición de contextos acotados, módulos y dependencias

---

## 🏗️ Visión General

Los módulos representan **Bounded Contexts internos** dentro de un **monolito modular**, con límites claros y contratos explícitos entre ellos. Cada módulo implementa **Clean Architecture con rigor asimétrico**:

- Dominio aislado de infraestructura y frontend
- Casos de uso testables de forma independiente
- Adaptadores de persistencia e interfaz intercambiables

---

## 📊 Diagrama Conceptual: Bounded Contexts

```
┌─────────────────────────────────────────────────────────────────┐
│                        MONOLITO MODULAR                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────┐         ┌─────────────────┐              │
│  │  Party /        │────────▶│  Producto /     │              │
│  │  Organización   │         │  Variante /     │              │
│  │  (Fundacional)  │         │  Categoría      │              │
│  │                 │         │  (Fundacional)  │              │
│  │ - Clientes      │         │                 │              │
│  │ - Proveedores   │         │ - Productos     │              │
│  │ - Jerarquías    │         │ - Variantes     │              │
│  │ - Costes base   │         │ - Categorías    │              │
│  └─────────────────┘         └─────────────────┘              │
│         │                            │                         │
│         │                            │                         │
│         ▼                            ▼                         │
│  ┌──────────────────────────────────────────┐                 │
│  │        Tarificación                      │                 │
│  │        (Núcleo Económico)                │                 │
│  │                                          │                 │
│  │  Coste base + Margen - Descuento        │                 │
│  │  = Precio final de venta                │                 │
│  └──────────────────────────────────────────┘                 │
│                    │                                           │
│                    ▼                                           │
│  ┌──────────────────────────────────────────┐                 │
│  │        Ventas                            │                 │
│  │        (Pedidos Estándar)                │                 │
│  │                                          │                 │
│  │  Cliente → Productos → Precio → Pedido  │                 │
│  └──────────────────────────────────────────┘                 │
│                    │                                           │
│                    │                                           │
│                    ▼                                           │
│  ┌──────────────────────────────────────────┐                 │
│  │        MES                               │                 │
│  │        (Producción Personalizada)        │                 │
│  │                                          │                 │
│  │  Pedido → Diseño → Producción → Entrega │                 │
│  └──────────────────────────────────────────┘                 │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│                    MÓDULOS TRANSVERSALES                        │
│  Seguridad │ Auditoría │ Gestión Documental │ i18n (frontend) │
└─────────────────────────────────────────────────────────────────┘

POST-MVP (No implementado en proyecto actual):
┌──────────────────┐      ┌──────────────────┐
│    Compras       │      │  Inventario /    │
│    (Formal)      │─────▶│  Stock           │
└──────────────────┘      └──────────────────┘
```

---

## 📋 Dominio Principal – Party / Organización (Fase 1 - Fundacional)
**Concepto:** Unifica la gestión de clientes y proveedores, permitiendo múltiples roles y relaciones entre ellos.
**Más detalles:** Consulte [docs/modules/party/README.md](../../../modules/party/README.md)
**Implementación:** Fase 1

---

## 📋 Dominio Principal – Producto / Variante / Categoría (Fase 1 - Fundacional)
**Concepto:** Gestiona el catálogo de productos, sus variantes y categorías, sirviendo como base para tarificación y ventas.
**Más detalles:** Consulte [docs/modules/product/README.md](../../../modules/product/README.md)
**Implementación:** Fase 1

---

## 📋 Dominio Principal – Tarificación (Fase 1 - Núcleo Económico)
**Concepto:** Motor central para el cálculo de precios de venta, aplicando costes, márgenes y descuentos.
**Más detalles:** Consulte [docs/modules/pricing/README.md](../../../modules/pricing/README.md)
**Implementación:** Fase 1

---

## 📋 Dominio Principal – Ventas (Fase 2)
**Concepto:** Gestiona el ciclo de vida de la venta, desde pedidos estándar hasta la generación de documentos mercantiles.
**Más detalles:** Consulte [docs/modules/sales/README.md](../../../modules/sales/README.md)
**Implementación:** Fase 2

---

## 📋 Subdominio MES – Producción Personalizada (Fase 3)
**Concepto:** Gestiona el ciclo de vida de producción personalizada, desde el diseño hasta el control de calidad.
**Más detalles:** Consulte [docs/modules/mes/README.md](../../../modules/mes/README.md)
**Implementación:** Fase 3

---

## 📋 Dominio Principal – Compras (Post-MVP)
**Concepto:** Cierre del ciclo económico con la entrada formal de mercancía, impactando en stock real y contabilidad.
**Rol en MVP:** NO existe como dominio independiente. Proveedores existen **solo como fuente de costes** dentro de Party.

---

## 📋 Dominio Principal – Inventario / Stock (Post-MVP)
**Concepto:** Gestión de inventario físico, lotes, ubicación y movimientos de stock.
**Rol en MVP:** NO existe. No hay control de stock en MVP.

---

## 🔄 Módulos Transversales

### Seguridad (MVP: básica; Post-MVP: avanzada)
**Concepto:** Autenticación JWT, RBAC básico y hash de passwords para proteger el sistema.
**Más detalles:** Consulte [ADR-010 - Estrategia de Seguridad](../adrs/adr-010-defense-in-depth-security-strategy.md)
**Implementación MVP:** Fase 0

---

### Auditoría (MVP: mínima; Post-MVP: completa)
**Concepto:** Registro de cambios críticos en tarificación y estados de producción.
**Implementación MVP:** Fase 1 - auditoría mínima

---

### Gestión Documental (MVP: mínima; Post-MVP: avanzada)
**Concepto:** Almacenamiento y gestión de diseños y archivos asociados a pedidos personalizados.
**Implementación MVP:** Fase 3

---

### i18n - Internacionalización (MVP: frontend; Post-MVP: tramatex-api)
**Concepto:** Soporte multi-idioma con etiquetas estáticas en frontend.
**Implementación MVP:** Integrado en frontend desde Fase 1

---

## 📊 Diagrama Conceptual: Flujo de Negocio MVP

```
PEDIDO ESTÁNDAR:
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│ Cliente  │───▶│ Selección│───▶│ Cálculo  │───▶│  Pedido  │
│ (Party)  │    │ Producto/│    │  Precio  │    │Confirmado│
│          │    │ Variante │    │(Tarif.)  │    │          │
└──────────┘    └──────────┘    └──────────┘    └──────────┘
                                                       │
                                                       ▼
                                                ┌──────────┐
                                                │  Entrega │
                                                └──────────┘

PEDIDO PERSONALIZADO:
┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐
│ Cliente  │───▶│ Selección│───▶│ Cálculo  │───▶│  Pedido  │
│ (Party)  │    │ Producto/│    │  Precio  │    │Confirmado│
│          │    │ Variante │    │(Tarif.)  │    │+ Diseño  │
└──────────┘    └──────────┘    └──────────┘    └──────────┘
                                                       │
                                                       ▼
                                                ┌──────────┐
                                                │   MES    │
                                                │ (Fases)  │
                                                └──────────┘
                                                       │
                        ┌──────────┬──────────┬────────┴───┬──────────┐
                        ▼          ▼          ▼            ▼          ▼
                    ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐
                    │ Diseño │▶│Aprob.  │▶│Marcaje │▶│ Taller │▶│Control │
                    │        │ │Diseño  │ │        │ │        │ │Calidad │
                    └────────┘ └────────┘ └────────┘ └────────┘ └────────┘
                                                                       │
                                                                       ▼
                                                                ┌──────────┐
                                                                │  Entrega │
                                                                └──────────┘
```

---

## 📋 Prioridad de Implementación

**Orden obligatorio (ADR-007):**

1. **Fase 0:** Fundaciones (infraestructura)
2. **Fase 1:** Party + Producto + Tarificación
3. **Fase 2:** Pedidos
4. **Fase 3:** MES

**Frontend paralelo:** En todas las fases, garantizando operatividad progresiva.

---

**Para especificación del MVP, ver:** [Project Vision and Scope](../../architecture/project-vision-and-scope.md)
**Para charter del proyecto, ver:** [Project Vision and Scope](../../architecture/project-vision-and-scope.md)**Para glosario de términos, ver:** [Glosario Unificado](../../architecture/glossary.md)
