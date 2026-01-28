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

**Concepto:** Party es un patrón de modelado que representa a cualquier persona u organización que tiene relación con el sistema, independientemente del rol que desempeñe.

**Problema que resuelve:**
- Evita duplicación de entidades (Cliente y Proveedor separados)
- Permite que una entidad tenga múltiples roles simultáneamente
- Garantiza consistencia de datos (una sola fuente de verdad)

**Entidades:**

- **Party** (Entidad raíz): ID, Tipo (Persona/Organización), Nombre, NIF/CIF, Dirección, Contacto
- **PartyRole**: Tipo de rol (Cliente, Proveedor), Estado (Activo/Inactivo), Fechas
- **Customer** (Especialización para rol Cliente):
  - Descuento base
  - Empresa matriz (para jerarquías)
  - Límite de crédito
  - Condiciones de pago
- **Supplier** (Especialización para rol Proveedor):
  - Código de proveedor
  - Días de entrega
  - Pedido mínimo
- **SupplierCost**: Costes base por producto/variante (histórico)

**Reglas de negocio:**

**Clientes:**
- Pueden tener jerarquía empresarial (empresa matriz → dependientes)
- Los descuentos pueden heredarse de la empresa matriz
- Pueden sobrescribir descuentos heredados con descuentos propios
- Ejemplo: "Construcciones ABC S.L." (matriz) con descuento 10%, sus obras heredan ese descuento

**Proveedores:**
- **NO tienen jerarquía** (simplificación MVP, según ADR-005)
- Proporcionan costes base para productos/variantes
- Estos costes son inputs para el motor de tarificación
- Un mismo Party puede ser Cliente Y Proveedor simultáneamente

**Dependencias:**
- Ninguna (es fundacional)
- Provee servicios a: Tarificación, Ventas, MES

**Implementación:** Fase 1

---

## 📋 Dominio Principal – Producto / Variante / Categoría (Fase 1 - Fundacional)

**Entidades:**

- **Producto** (Entidad raíz): ID, Código, Nombre, Descripción
- **Variante**: ID, Producto (FK), Tipo (Talla, Color, Arreglo, etc.), Valor, Modificador de precio
- **Categoría**: ID, Nombre, Descripción, Categoría padre (jerarquía)

**Reglas de negocio:**

- Toda variante **puede afectar al precio** (modificador de precio)
- Las categorías **NO influyen en tarificación**, solo estructuran el catálogo
- Un producto puede tener múltiples variantes
- Las variantes son parte activa del modelo económico

**Dependencias:**
- Ninguna (es fundacional)
- Provee servicios a: Tarificación, Ventas, MES

**Implementación:** Fase 1

---

## 📋 Dominio Principal – Tarificación (Fase 1 - Núcleo Económico)

**Entidades:**

- **Tarifa**: Reglas de cálculo de precio
- **Regla de tarificación**: Márgenes, descuentos, condiciones
- **Precio calculado**: Resultado del motor de tarificación (no persistido, calculado on-demand)

**Reglas de negocio (Motor de tarificación):**

```
Precio final = ((Coste base + Mod. variante) × (1 + Margen) + Mod. otros) × (1 - Descuento total)

Donde:
- Coste base: de Supplier (Party rol Proveedor)
- Mod. variante: de Variante
- Margen: regla de tarificación (por producto, categoría, o general)
- Mod. otros: no dependientes de proveedor (p.e. arreglo de prenda)
- Descuento total: base cliente, heredado, u override específico
```

**Dependencias:**
- **Requiere:** Party (costes de proveedor, descuentos de cliente), Producto/Variante
- **Provee servicios a:** Ventas

**Nota crítica:** El motor de tarificación es **lógica de dominio pura**, completamente testeable en aislamiento, sin dependencias de infraestructura.

**Implementación:** Fase 1

---

## 📋 Dominio Principal – Ventas (Fase 2)

**Entidades:**

- **Pedido**: ID, Cliente (FK a Party), Fecha, Estado, Total
- **Línea de pedido**: ID, Pedido (FK), Producto (FK), Variante (FK), Cantidad, Precio unitario, Subtotal
- **Estado de pedido**: Borrador, Confirmado, En preparación, Entregado, Cancelado
- **Documentación mercantil**: Presupuesto, Pedido, Albarán, Factura

**Funciones:**

- Crear/editar pedidos/documentos de venta
- Añadir/modificar/eliminar líneas de pedido
- Calcular precio usando motor de tarificación
- Cambiar estado de pedido
- Visualizar historial de pedidos
- Generación de documentos **delegada a frontend vía Web-to-Print**

**Dependencias:**
- **Requiere:** Party (cliente), Producto/Variante, Tarificación
- **Provee servicios a:** MES (para pedidos personalizados)

**Implementación:** Fase 2

---

## 📋 Subdominio MES – Producción Personalizada (Fase 3)

**Entidades:**

- **Pedido personalizado** (extensión de Pedido con requisitos de producción)
- **Estado de producción**: Diseño, Aprobación de diseño, Impresión, Marcaje, Taller, Control de calidad, Listo para entrega. No todas las fases son siempre necesarias. Pueden crearse nuevas según las necesidades
- **Trabajo de taller**: ID, Pedido personalizado (FK), Operario, Observaciones
- **Diseño/Archivo asociado**: ID, Pedido personalizado (FK), Ruta NAS, Tipo archivo, Tamaño

**Funciones:**

- Gestionar pedidos personalizados
- Registrar estados de producción
- Terminal de taller (interfaz simplificada para tablets)
- Adjuntar/visualizar diseños
- Trazabilidad de trabajos y cambios de estado
- Integración con inventario para consumos de materiales (Post-MVP)

**Dependencias:**
- **Requiere:** Ventas (pedido base), Producto/Variante, Gestión documental
- **Provee servicios a:** Ninguno (es subdominio final)

**Implementación:** Fase 3

---

## 📋 Dominio Principal – Compras (Post-MVP)

**Rol en MVP:** NO existe como dominio independiente. Proveedores existen **solo como fuente de costes** dentro de Party.

**Rol Post-MVP:** Cierre del ciclo económico, entrada formal de mercancía, impacto en stock real y contabilidad.

**Entidades previstas (Post-MVP):**
- Pedido de compra
- Línea de pedido de compra
- Recepción de mercancía
- Albarán y Factura de proveedor

**Dependencias previstas:**
- **Requiere:** Party (proveedor), Producto/Variante, Inventario/Stock

---

## 📋 Dominio Principal – Inventario / Stock (Post-MVP)

**Rol en MVP:** NO existe. No hay control de stock en MVP.

**Entidades previstas (Post-MVP):**
- Producto (referencia a Producto/Variante)
- Stock físico, Lote, Ubicación
- Movimiento de stock
- Regularización

**Dependencias previstas:**
- **Requiere:** Producto/Variante, Compras (entradas), Ventas (salidas)

---

## 🔄 Módulos Transversales

### Seguridad (MVP: básica; Post-MVP: avanzada)

**Referencia Arquitectónica:** [ADR-010 - Estrategia de Seguridad](./adr/ADR-010-estrategia-seguridad-defensa-profundidad.md)

**MVP:**
- Autenticación JWT
- Roles básicos: Admin, Comercial, Diseño, Taller
- Control de acceso por rol (RBAC básico)
- Hash de passwords

**Post-MVP:**
- RBAC avanzado con permisos granulares
- Roles personalizados
- Autenticación multifactor

**Implementación MVP:** Fase 0

---

### Auditoría (MVP: mínima; Post-MVP: completa)

**MVP:**
- Log de cambios críticos en tarificación (precios, márgenes, descuentos)
- Log de cambios en estados de producción (MES)
- Registro de quién, cuándo, qué cambió

**Post-MVP:**
- Trazabilidad completa de todos los cambios
- Retención configurable de datos
- Informes de auditoría
- Cumplimiento normativo

**Implementación MVP:** Fase 1 - auditoría mínima

---

### Gestión Documental (MVP: mínima; Post-MVP: avanzada)

**MVP:**
- Almacenamiento de diseños en NAS
- Indexación en PostgreSQL (ID, ruta, tipo, tamaño, fecha)
- Adjuntar/visualizar archivos asociados a pedidos personalizados
- Trazabilidad básica de archivos

**Post-MVP:**
- Versionado de documentos
- Control de acceso granular
- Búsqueda avanzada
- OCR y extracción de metadatos

**Implementación MVP:** Fase 3

---

### i18n - Internacionalización (MVP: frontend; Post-MVP: tramatex-api)

**MVP:**
- Etiquetas estáticas en frontend (Vue-i18n)
- Idiomas: Español (por defecto), posible catalán/valenciano
- Mensajes de interfaz traducibles
- **NO incluye tramatex-api i18n**

**Post-MVP:**
- i18n completa tramatex-api (go-i18n)
- Mensajes de error, validaciones, notificaciones traducibles
- Soporte multiidioma completo

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

**Para especificación del MVP, ver:** [Project Charter & MVP Specification](../1_project/README.md)  
**Para charter del proyecto, ver:** [Project Charter & MVP Specification](../1_project/README.md)  
**Para glosario de términos, ver:** [Glosario Unificado](./glossary.md)
