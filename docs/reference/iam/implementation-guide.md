# Implementación - Módulo IAM

## Estructura de Directorios

La estructura de directorios para el módulo `iam` en el tramatex-api es la siguiente:

```
apps/tramatex-api/internal/iam/
├── application/
│   └── usecase/
│       ├── register_user.go
│       ├── login_user.go
│       └── ... (otros casos de uso)
├── domain/
│   ├── model/
│   │   ├── user.go
│   │   ├── role.go
│   │   └── permission.go
│   ├── repository/
│   │   ├── user_repository.go
│   │   └── role_repository.go
│   └── service/
│       └── auth_service.go
├── infrastructure/
│   └── persistence/
│       ├── user_postgres_repository.go
│       └── role_postgres_repository.go
└── interfaces/
    └── http/
        └── handler/
            └── iam_handler.go
```

## Dependencias Clave

*   **GORM:** Para la capa de persistencia.
*   **Go-JWT:** Para la generación y validación de JSON Web Tokens.
*   **Testify:** Para aserciones en las pruebas unitarias y de integración.

## Flujo de Implementación Sugerido

1.  **Definir Modelos de Dominio:** Empezar por `user.go`, `role.go`, y `permission.go`.
2.  **Definir Interfaces de Repositorio:** Crear `user_repository.go` y `role_repository.go` en el dominio.
3.  **Implementar Casos de Uso:** Desarrollar la lógica en la capa de aplicación (e.g., `register_user.go`).
4.  **Implementar Repositorios:** Crear las implementaciones de GORM en la capa de infraestructura.
5.  **Exponer Handlers HTTP:** Crear los endpoints en la capa de interfaces.
6.  **Añadir Pruebas:** Implementar pruebas unitarias para casos de uso y de integración para los repositorios.
