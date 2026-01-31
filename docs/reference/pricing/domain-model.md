# Implementación - Módulo Pricing

## Estructura de Directorios

La estructura de directorios para el módulo `pricing` en el tramatex-api es la siguiente:

```
apps/tramatex-api/internal/pricing/
├── application/
│   └── usecase/
│       ├── calculate_price.go
│       └── ...
├── domain/
│   ├── model/
│   │   ├── price_rule.go
│   │   └── price_calculation.go
│   └── repository/
│       └── price_rule_repository.go
├── infrastructure/
│   └── persistence/
│       └── price_rule_postgres_repository.go
└── interfaces/
    └── http/
        └── handler/
            └── pricing_handler.go
```

## Dependencias Clave

*   **GORM:** Para la capa de persistencia.
*   **Testify:** Para aserciones en las pruebas.

## Flujo de Implementación Sugerido

1.  **Definir Modelos de Dominio:** Empezar por `price_rule.go` y `price_calculation.go`.
2.  **Definir Interface de Repositorio:** Crear `price_rule_repository.go` en el dominio.
3.  **Implementar Caso de Uso Principal:** Desarrollar `calculate_price.go`.
4.  **Implementar Repositorio:** Crear la implementación de GORM.
5.  **Exponer Handler HTTP:** Crear el endpoint para el cálculo de precios.
6.  **Añadir Pruebas (Cobertura >95%):** Implementar pruebas unitarias exhaustivas para el motor de cálculo.