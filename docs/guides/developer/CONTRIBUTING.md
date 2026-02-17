# 🤝 Guía de Contribución - TramaTex

¡Gracias por tu interés en contribuir a TramaTex! Esta guía te ayudará a entender el proceso de desarrollo y los estándares del proyecto.

---

## Workflow de Desarrollo

### 1. Crear Branch

Usa prefijos para tus ramas para indicar el propósito del cambio:
- `feature/`: Nueva funcionalidad (ej. `feature/pricing-engine`)
- `fix/`: Corrección de un bug (ej. `fix/login-error-500`)
- `refactor/`: Mejoras en el código sin cambiar la funcionalidad (ej. `refactor/move-validators`)
- `docs/`: Cambios en la documentación (ej. `docs/update-adr-011`)
- `test/`: Añadir o mejorar tests (ej. `test/add-coverage-for-party-module`)

### 2. Desarrollo Guiado por Pruebas

Para asegurar la calidad y el correcto funcionamiento del software, se fomenta un enfoque de desarrollo guiado por pruebas. Se recomienda un ciclo iterativo donde la escritura de tests acompaña la implementación de la lógica de negocio (capas de dominio y aplicación):
1.  **Definir el comportamiento esperado**: Describe la funcionalidad a implementar mediante tests.
2.  **Implementar la funcionalidad**: Escribe el código necesario para cumplir con el comportamiento definido.
3.  **Verificar y refinar**: Ejecuta los tests para asegurar que la funcionalidad es correcta y refactoriza el código y los tests para mejorar la calidad.

### 3. Formato de Commits

Usa [Conventional Commits](https://www.conventionalcommits.org/). Esto nos ayuda a generar changelogs automáticamente y a entender la historia del proyecto.

**Formato:** `tipo(alcance): descripción breve`

- **Tipos:** `feat`, `fix`, `refactor`, `test`, `docs`, `chore`
- **Alcance (opcional):** El módulo afectado (ej. `iam`, `party`, `ci`)

**Ejemplos:**
```
feat(party): implementar CRUD de organizaciones
fix(iam): corregir validación de contraseña con caracteres especiales
refactor(domain): extraer lógica de validación a value objects
test(party): agregar tests de integración para el repositorio de organizaciones
docs(adr): documentar la nueva estrategia de testing
```

---

## Estándares de Código

### Backend (Go)
- Sigue las directrices de [Effective Go](https://go.dev/doc/effective_go).
- Usa `golangci-lint` para el análisis estático. La configuración está en `.golangci.yml`.
- Comenta el código en **Inglés**.

### Frontend (Vue 3)
- Usa la **Composition API** exclusivamente.
- Escribe nuevos componentes con `<script setup lang="ts">`.
- Sigue las reglas de `ESLint` y `Prettier` definidas en el proyecto.
- Comenta el código en **Inglés**.

---

## Herramientas de Calidad (Linting y Formateo)

Para asegurar la consistencia y calidad del código, usamos un conjunto de herramientas que deben ejecutarse antes de crear un Pull Request.

### Instalación

La primera vez, necesitas instalar las herramientas de desarrollo requeridas (incluyendo el linter):
```bash
make install-tools
```

### Backend (Go)

- **Linting:** Para ejecutar el linter y verificar si hay errores de estilo o posibles bugs:
  ```bash
  make tramatex-api-lint
  ```

- **Formateo:** Para formatear el código según las convenciones del proyecto:
  ```bash
  make tramatex-api-fmt
  ```

### Frontend (Vue)

- **Linting y Formateo:** Las herramientas de frontend (ESLint, Prettier) suelen estar integradas en el editor de código y se ejecutan automáticamente. También puedes ejecutarlas manualmente:
  ```bash
  # Desde la carpeta apps/frontend
  npm run lint
  ```

---

## Testing

La estrategia de testing es una parte fundamental de nuestra calidad. Consulta el [ADR-011: Estrategia de Testing y Coverage](../../architecture/adrs/ADR-011-testing-coverage-strategy.md) para los detalles completos.

**Comandos clave:**
```bash
# Backend
make test
make coverage

# Frontend
npm run test:unit
npm run test:unit -- --coverage
```

---

## Proceso de Pull Request (PR)

1.  **Antes de crear el PR:**
    - Asegúrate de que todos los tests pasan localmente.
    - Asegúrate de que los linters no reportan errores.
    - Actualiza tu rama con la última versión de `main` (`git pull origin main`).

2.  **Crear el PR:**
    - Usa un título claro que siga el formato de Conventional Commits.
    - En la descripción, explica **qué** hace el PR y **por qué**.
    - Enlaza cualquier `issue` o tarea relacionada.

3.  **Revisión de Código:**
    - Se requiere al menos **una aprobación** de otro miembro del equipo.
    - La pipeline de CI/CD en GitHub Actions debe pasar con éxito.
    - Responde a los comentarios y realiza los ajustes necesarios.

4.  **Merge:**
    - El merge se realizará usando **"Squash and Merge"** para mantener un historial de `main` limpio, donde cada commit corresponde a una feature o fix completo.
