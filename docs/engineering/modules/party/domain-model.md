# Implementación - Módulo Party

## Estructura de Directorios

La estructura de directorios para el módulo `party` en el tramatex-api es la siguiente:

```
apps/tramatex-api/internal/party/
├── application/
│   └── usecase/
│       ├── create_party.go
│       └── ...
├── domain/
│   ├── model/
│   │   ├── party.go
│   │   ├── contact.go
│   │   └── address.go
│   └── repository/
│       └── party_repository.go
├── infrastructure/
│   └── persistence/
│       └── party_postgres_repository.go
└── interfaces/
    └── http/
        └── handler/
            └── party_handler.go
```

## Dependencias Clave

*   **GORM:** Para la capa de persistencia.
*   **Testify:** Para aserciones en las pruebas.

## Flujo de Implementación Sugerido

1.  **Definir Modelos de Dominio:** Empezar por `party.go`, `contact.go`, y `address.go`.
2.  **Definir Interface de Repositorio:** Crear `party_repository.go` en el dominio.
3.  **Implementar Casos de Uso:** Desarrollar `create_party.go` y otros casos de uso.
4.  **Implementar Repositorio:** Crear la implementación de GORM.
5.  **Exponer Handlers HTTP:** Crear los endpoints en la capa de interfaces.
6.  **Añadir Pruebas:** Implementar pruebas unitarias y de integración.