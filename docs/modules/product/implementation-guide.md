# Guia de Implementacion - Modulo Product

## Estructura de Directorios

La estructura del modulo `product` en tramatex-api sigue Clean Architecture y refleja el modelo Product:

```
apps/tramatex-api/internal/product/
├── application/
│   ├── commands.go
│   ├── queries.go
│   ├── dtos.go
│   └── product_service.go
├── domain/
│   ├── attribute.go
│   ├── product.go
│   ├── variant.go
│   ├── party_service_configuration.go
│   ├── enums.go
│   ├── errors.go
│   └── repository.go
├── infrastructure/
│   └── persistence/
│       ├── gorm_attribute_repository.go
│       ├── gorm_brand_repository.go
│       ├── gorm_product_group_repository.go
│       ├── gorm_product_repository.go
│       ├── gorm_variant_repository.go
│       ├── gorm_party_service_configuration_repository.go
│       ├── product_data_model.go
│       ├── attribute_data_model.go
│       ├── variant_data_model.go
│       ├── party_service_configuration_data_model.go
│       ├── audit_context.go
│       └── test_helpers.go
└── interfaces/
	└── http/handler/
		├── product_handler.go
		└── product_handler_test.go
```

## Dependencias Clave

* **GORM:** Para la capa de persistencia (`infrastructure/persistence`).
* **Gin:** Para los handlers HTTP (`interfaces/http/handler`).
* **Testify:** Para aserciones en las pruebas.

## Flujo de Implementacion Sugerido

1. **Definir Modelo de Dominio:** `attribute.go`, `product.go`, `variant.go`, `party_service_configuration.go`.
2. **Definir Interfaces de Repositorio:** `domain/repository.go`.
3. **Implementar Casos de Uso:** `commands.go`, `queries.go` y `product_service.go` en `application`.
4. **Implementar Repositorios GORM:** archivos `gorm_*.go` y `*_data_model.go`.
5. **Exponer Handlers HTTP:** `interfaces/http/handler/product_handler.go`.
6. **Añadir Pruebas:** unitarias y de integracion por capa.
