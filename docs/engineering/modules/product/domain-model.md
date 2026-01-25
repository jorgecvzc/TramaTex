# Implementación - Módulo Product

## Estructura de Directorios

La estructura de directorios para el módulo `product` en el tramatex-api es la siguiente:

```
apps/tramatex-api/internal/product/
├── application/
│   └── usecase/
│       ├── create_product.go
│       └── ...
├── domain/
│   ├── model/
│   │   ├── product.go
│   │   └── product_variant.go
│   └── repository/
│       └── product_repository.go
├── infrastructure/
│   └── persistence/
│       └── product_postgres_repository.go
└── interfaces/
    └── http/
        └── handler/
            └── product_handler.go
```

## Dependencias Clave

*   **GORM:** Para la capa de persistencia.
*   **Testify:** Para aserciones en las pruebas.

## Flujo de Implementación Sugerido

1.  **Definir Modelos de Dominio:** Empezar por `product.go` y `product_variant.go`.
2.  **Definir Interface de Repositorio:** Crear `product_repository.go` en el dominio.
3.  **Implementar Casos de Uso:** Desarrollar `create_product.go` y `create_variant.go`.
4.  **Implementar Repositorio:** Crear la implementación de GORM.
5.  **Exponer Handlers HTTP:** Crear los endpoints para el catálogo.
6.  **Añadir Pruebas:** Implementar pruebas unitarias y de integración.