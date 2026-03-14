# Implementación - Módulo Pricing

## Estructura de Directorios

La estructura del módulo `pricing` sigue Clean Architecture y los estándares del proyecto:

```
apps/tramatex-api/internal/pricing/
├── application/             # Servicios de Aplicación y Motores
│   ├── pricing_engine_service.go # Motor de cálculo síncrono
│   ├── pricing_service.go        # Gestión de reglas y márgenes
│   └── dtos.go
├── domain/                  # Entidades, Reglas y Value Objects
│   ├── pricing_rule.go
│   ├── price_calculation.go
│   ├── money.go
│   ├── percentage.go
│   └── errors.go
├── infrastructure/          # Persistencia y Clientes Externos
│   └── persistence/
│       ├── gorm_pricing_repository.go
│       └── pricing_data_model.go
└── interfaces/              # Capa de Entrada (HTTP Handlers)
    └── http/handler/
        ├── pricing_engine_handler.go
        └── pricing_handler.go
```

## Estándares de Datos y Negocio

### 1. Mandatos Globales
- **Campos de Auditoría:** Los campos `CreatedAt`, `UpdatedAt`, `CreatedBy` y `UpdatedBy` **deben excluirse** de las entidades de dominio (`domain/`). Solo deben existir en los modelos de datos de persistencia (`persistence/pricing_data_model.go`).
- **Moneda Única:** Todas las operaciones monetarias y cálculos del motor de precios deben gestionarse exclusivamente en **Euros (€)**. El tipo `Money` debe validar esta restricción.

### 2. Mapeo de Errores
El módulo debe mapear sus errores de dominio a códigos HTTP estándar:
- `ErrCodeValidation` -> `400 Bad Request`
- `ErrCodeNotFound`   -> `404 Not Found`
- `ErrCodeConflict`   -> `409 Conflict`

## Dependencias Clave

* **GORM:** Para la capa de persistencia.
* **Gin:** Para los handlers HTTP.
* **Testify:** Para las pruebas unitarias exhaustivas del motor de cálculo.

## Flujo de Implementación Sugerido

1. **Definir Modelo de Dominio:** Entidades de reglas, cálculos y VOs monetarios.
2. **Definir Interfaces de Repositorio:** En la capa de dominio.
3. **Implementar Motor de Cálculo:** En `application/pricing_engine_service.go`.
4. **Implementar Persistencia:** Modelos de datos de GORM con campos de auditoría.
5. **Exponer Handlers HTTP:** Handlers diferenciados para el motor (`/calculate`) y las reglas (`/rules`).
6. **Añadir Pruebas (Objetivo >90%):** Implementar pruebas unitarias para todas las combinaciones de reglas de precio.