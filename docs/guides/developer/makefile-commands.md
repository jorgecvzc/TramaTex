# Comandos Make Comunes en el Proyecto TramaTex

Este documento lista los comandos `make` comunes utilizados en el proyecto TramaTex para automatizar tareas de desarrollo, testing y construcción, complementando la decisión arquitectónica definida en [ADR-009: Estructura de Carpetas y Organización del Proyecto](../../architecture/adrs/adr-009-project-structure.md).

---

## Comandos Make (Makefile)

### Backend (tramatex-api)

Estos comandos se ejecutan desde el directorio `apps/tramatex-api/`.

**Desarrollo:**
```makefile
make run              # Ejecutar aplicación
make test             # Ejecutar tests
make test-coverage    # Cobertura de tests
make lint             # Linter
make fmt              # Formatear código
```

**Base de datos:**
```makefile
make migrate-up       # Aplicar migraciones
make migrate-down     # Revertir última migración
make seed             # Seed de datos
```

**Build:**
```makefile
make build            # Compilar binario
make docker-build     # Build imagen Docker
make docker-up        # Levantar stack completo
make docker-down      # Detener stack
```

### Frontend (Vue.js)

Estos comandos se ejecutan desde el directorio `apps/frontend/`.

```makefile
make dev              # Servidor de desarrollo
make build            # Build para producción
make preview          # Preview del build
make lint             # ESLint
make format           # Prettier
```
