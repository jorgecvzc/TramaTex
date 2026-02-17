# Implementación - Módulo Party

## Estructura de Directorios

La estructura del módulo `party` en tramatex-api sigue Clean Architecture y refleja el modelo Party:

```
apps/tramatex-api/internal/party/
├── application/
│   ├── party_commands.go
│   ├── party_queries.go
│   └── ...
├── domain/
│   ├── party.go
│   ├── party_profiles.go
│   ├── party_types.go
│   ├── value_objects.go
│   └── ...
├── persistence/
│   ├── repository.go
│   ├── gorm_party.go
│   ├── party_data_model.go
│   └── test_helpers.go
└── interfaces/
    ├── party_handlers.go
    ├── party_dto.go
    └── helpers.go
```

## Dependencias Clave

* **GORM:** Para la capa de persistencia (`persistence`).
* **Gin:** Para los handlers HTTP (`interfaces`).
* **Testify:** Para aserciones en las pruebas.

## Flujo de Implementación Sugerido

1. **Definir Modelo de Dominio:** `party.go`, `party_profiles.go`, `party_types.go`, `value_objects.go`.
2. **Definir Interfaces de Repositorio:** `persistence/repository.go`.
3. **Implementar Casos de Uso:** `party_commands.go` y `party_queries.go` en `application`.
4. **Implementar Repositorio:** `persistence/gorm_party.go` + `party_data_model.go`.
5. **Exponer Handlers HTTP:** `party_handlers.go` y mappers en `party_dto.go`.
6. **Añadir Pruebas:** unitarias y de integración por capa.
