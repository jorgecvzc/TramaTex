# Flujo de Creación de Carpetas por Fase en el Proyecto TramaTex

Este documento detalla el flujo incremental de creación de carpetas por fase de desarrollo en el proyecto TramaTex, complementando la decisión arquitectónica definida en [ADR-009: Estructura de Carpetas y Organización del Proyecto](../../architecture/adrs/adr-009-project-structure.md). Este enfoque modular y por fases se alinea con el [ADR-007: Orden de Implementación de Módulos](../../architecture/adrs/adr-007-module-implementation-order.md).

---

## Flujo de Creación de Carpetas por Fase

### Fase 0 (Fundaciones)

En esta fase, se establecen los cimientos del proyecto, incluyendo la estructura básica, la seguridad y la configuración inicial.

```
tramatex/
├── docs/
│   └── adr/
├── apps/tramatex-api/
│   ├── cmd/api/
│   ├── internal/
│   │   ├── infrastructure/security/  # JWT, RBAC
│   │   └── interfaces/http/
│   │       ├── middleware/
│   │       └── handlers/
│   └── config/
├── apps/frontend/
│   └── src/
│       ├── views/auth/
│       ├── stores/
│       └── services/
└── docker/
```

### Fase 1 (Party + Producto + Tarificación)

Esta fase se enfoca en los módulos fundamentales del negocio: gestión de clientes/proveedores, catálogo de productos y el motor de precios.

```
# Se añaden:
apps/tramatex-api/internal/domain/
├── party/
├── product/
└── pricing/

apps/tramatex-api/internal/application/
├── party/
├── product/
└── pricing/

apps/tramatex-api/internal/infrastructure/persistence/postgres/
├── party_repository.go
├── product_repository.go
└── pricing_repository.go

apps/tramatex-api/internal/infrastructure/persistence/migrations/
├── 000002_create_parties.up.sql
├── 000003_create_products.up.sql
└── 000004_create_pricing.up.sql

apps/frontend/src/
├── views/party/
├── views/product/
├── views/pricing/
├── components/party/
├── components/product/
└── components/pricing/

docs/modules/
├── party/
├── product/
└── pricing/
```

### Fase 2 (Pedidos)

La segunda fase de implementación del dominio se centra en la gestión de ventas y pedidos.

```
# Se añaden:
apps/tramatex-api/internal/domain/sales/
apps/tramatex-api/internal/application/sales/
apps/tramatex-api/internal/infrastructure/persistence/postgres/order_repository.go
apps/frontend/src/views/sales/
apps/frontend/src/components/sales/
docs/modules/sales/
```

### Fase 3 (MES)

La fase final del MVP integra el sistema de ejecución de manufactura.

```
# Se añaden:
apps/tramatex-api/internal/domain/mes/
apps/tramatex-api/internal/application/mes/
apps/tramatex-api/internal/infrastructure/
├── persistence/postgres/mes_repository.go
└── storage/nas/
apps/frontend/src/views/mes/
apps/frontend/src/components/mes/
docs/modules/mes/
```
