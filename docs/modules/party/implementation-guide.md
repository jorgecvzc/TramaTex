# Implementación - Módulo Party

## Estructura de Directorios

La estructura de directorios para el módulo `party` en el tramatex-api sigue los principios de Clean Architecture:

```
apps/tramatex-api/internal/party/
├── application/
│   ├── commands.go
│   ├── queries.go
│   └── ...
├── domain/
│   ├── organization.go
│   ├── person.go
│   ├── enums.go
│   ├── value_objects.go
│   └── ...
├── persistence/
│   ├── repository.go
│   ├── postgresql.go
│   └── in_memory.go
└── interfaces/
    └── http/
        ├── handlers.go
        └── dto.go
```

## Dependencias Clave

*   **GORM:** Para la capa de persistencia (`persistence`).
*   **Gin:** Para el enrutamiento y los handlers HTTP (`interfaces/http`).
*   **Testify:** Para aserciones en las pruebas.

## Flujo de Implementación Sugerido

1.  **Definir Modelos de Dominio:** Empezar por `organization.go`, `person.go`, `enums.go` y `value_objects.go`.
2.  **Definir Interfaces de Repositorio:** Crear las interfaces en `persistence/repository.go`.
3.  **Implementar Casos de Uso:** Desarrollar los `commands.go` y `queries.go` en la capa de aplicación.
4.  **Implementar Repositorio:** Crear la implementación de GORM en `persistence/postgresql.go`.
5.  **Exponer Handlers HTTP:** Crear los endpoints en `interfaces/http/handlers.go`, usando los DTOs definidos en `dto.go`.
6.  **Añadir Pruebas:** Implementar pruebas unitarias y de integración para todas las capas.
