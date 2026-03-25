# Guía de Implementación - Módulo Product

## Estructura de Directorios

La estructura del módulo `product` en tramatex-api sigue Clean Architecture y los estándares del proyecto:

```
apps/tramatex-api/internal/product/
├── application/             # Casos de Uso y Servicios
│   ├── commands.go
│   ├── queries.go
│   ├── dtos.go
│   └── product_service.go
├── domain/                  # Lógica de Negocio (Entidades, VOs, Repo interfaces)
│   ├── attribute.go
│   ├── product.go
│   ├── variant.go
│   ├── enums.go
│   └── errors.go
├── persistence/             # Infraestructura GORM (Data Models y Repositorios)
│   ├── gorm_product_repository.go
│   ├── gorm_variant_repository.go
│   ├── product_data_model.go
│   ├── variant_data_model.go
│   └── repository.go        # Interfaces de dominio implementadas aquí
└── interfaces/              # Capa de Entrada (HTTP Handlers)
    └── http/handler/
        ├── product_handler.go
        └── ...
```

## Estándares de Datos y Negocio

### 1. Mandatos Globales
- **Campos de Auditoría:** Los campos `CreatedAt`, `UpdatedAt`, `CreatedBy` y `UpdatedBy` **deben excluirse** de las entidades de dominio (`domain/`). Solo deben existir en los modelos de datos de persistencia (`persistence/*_data_model.go`).
- **Moneda Única:** Todas las operaciones de precios y costes deben gestionarse exclusivamente en **Euros (€)**.

### 2. Gestión de Errores Estandarizada
El módulo `Product` delega la traducción de errores de dominio a respuestas HTTP en el `ErrorHandlerMiddleware` de la capa `shared`. 

Para que esto funcione:
1. **Definir Errores en Dominio**: Todos los errores de negocio deben definirse en `internal/product/domain/errors.go`.
2. **Implementar `HTTPStatuser`**: Los errores de dominio deben implementar la interfaz `shared/domain.HTTPStatuser` para indicar su código HTTP correspondiente (ej. `ErrProductNotFound` devuelve `404`).
3. **Delegación en Handlers**: Los controladores Gin NO deben formatear respuestas de error manualmente. Deben simplemente adjuntar el error al contexto: `c.Error(err)`. El middleware se encargará de sanitizar la respuesta y registrar el log con el ID de petición.

Mapeo estándar sugerido:
- `ErrCodeValidation` -> `400 Bad Request`
- `ErrCodeNotFound`   -> `404 Not Found`
- `ErrCodeConflict`   -> `409 Conflict`
- `ErrCodeInternal`   -> `500 Internal Server Error`

## Dependencias Clave

* **GORM:** Para la capa de persistencia e infraestructura.
* **Gin:** Para los handlers HTTP.
* **Testify:** Para aserciones en las pruebas unitarias e integración.

## Flujo de Implementación Sugerido

1. **Definir Modelo de Dominio:** Entidades, VOs y enums sin campos de infraestructura.
2. **Definir Interfaces de Repositorio:** En la capa de dominio.
3. **Implementar Casos de Uso:** En la capa de aplicación (servicios).
4. **Implementar Repositorios GORM:** Modelos de datos y mapeadores en la capa de persistencia.
5. **Exponer Handlers HTTP:** Handlers en la capa de interfaces.
6. **Añadir Pruebas:** unitarias por capa y de integración para los repositorios.
