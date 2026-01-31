# Implementación - Módulo Sales

## Estructura de Directorios

La estructura de directorios para el módulo `sales` en el tramatex-api es la siguiente:

```
apps/tramatex-api/internal/sales/
├── application/
│   └── usecase/
│       ├── create_quote.go
│       ├── convert_to_order.go
│       └── ...
├── domain/
│   ├── model/
│   │   ├── sales_order.go
│   │   └── order_line_item.go
│   └── repository/
│       └── sales_order_repository.go
├── infrastructure/
│   └── persistence/
│       └── sales_order_postgres_repository.go
└── interfaces/
    └── http/
        └── handler/
            └── sales_handler.go
```

## Dependencias Clave

*   **GORM:** Para la capa de persistencia.
*   **Testify:** Para aserciones en las pruebas.
*   **Módulos `party`, `product`, `pricing`:** Para obtener datos y calcular precios.

## Flujo de Implementación Sugerido

1.  **Definir Modelos de Dominio:** Empezar por `sales_order.go` y `order_line_item.go`.
2.  **Definir Interface de Repositorio:** Crear `sales_order_repository.go` en el dominio.
3.  **Implementar Casos de Uso:** Desarrollar `create_quote.go` y `convert_to_order.go`.
4.  **Implementar Repositorio:** Crear la implementación de GORM.
5.  **Exponer Handlers HTTP:** Crear los endpoints para la gestión de órdenes.
6.  **Añadir Pruebas:** Implementar pruebas unitarias y de integración.