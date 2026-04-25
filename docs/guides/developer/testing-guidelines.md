# Guía de Testing en TramaTex

Este documento detalla la estrategia, ubicación y ejecución de los tests en el proyecto TramaTex, asegurando la calidad y estabilidad del sistema bajo los estándares de Clean Architecture y DDD.

---

## 🏗️ Estrategia de Testing

Seguimos una estrategia basada en la pirámide de tests, priorizando los tests unitarios para la lógica de negocio y los de integración para los adaptadores de infraestructura.

### 1. Tests Unitarios (Dominio y Aplicación)
Prueban la lógica pura de las entidades de dominio y los servicios de aplicación sin dependencias externas.
- **Ubicación:** Mismo directorio que el código que prueban.
- **Formato:** `nombre_test.go`.
- **Mocks:** Utilizamos `sqlmock` para simular la base de datos en los repositorios cuando sea necesario probar lógica de persistencia sin infraestructura real.

### 2. Tests de Integración (Infraestructura)
Verifican la correcta interacción con la base de datos PostgreSQL real.
- **Ubicación:** En las carpetas de `persistence` de cada módulo.
- **Infraestructura:** Requieren una instancia de PostgreSQL (vía Docker o local).
- **Helpers:** Utilizamos `test_helpers.go` en cada módulo para gestionar el ciclo de vida del esquema de pruebas.

---

## 🚀 Ejecución de Tests

### Backend (Go)

Desde el directorio `apps/tramatex-api`:

```bash
# Ejecutar todos los tests (Unitarios + Integración si hay DB)
go test -v ./...

# Ejecutar tests de un módulo específico
go test -v ./internal/party/...

# Ejecutar con detección de condiciones de carrera (Recomendado CI)
go test -v -race ./...
```

### Frontend (Vue.js)

Desde el directorio `apps/frontend`:

```bash
# Ejecutar tests unitarios (Vitest)
npm run test:unit

# Ejecutar linter
npm run lint
```

---

## ⚙️ Configuración del Entorno de Tests

Los `test_helpers` del backend están diseñados para ser flexibles y compatibles con entornos locales y de CI (GitHub Actions).

### Prioridad de Configuración:
El sistema carga las credenciales de la base de datos siguiendo este orden (de mayor a menor prioridad):
1. **Variables Específicas:** `TRAMATEX_TEST_DB_HOST`, `TRAMATEX_TEST_DB_USER`, etc.
2. **Archivos Locales:** `.env.local` o `.env.remote` en la raíz del proyecto.
3. **Variables Estándar:** `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`.
4. **Valores por Defecto:** `localhost`, `tramatex`.

### Autodeteción y Saltado (Safe Skipping):
Si un test de integración no logra conectar con la base de datos (por ejemplo, si Docker no está corriendo localmente), el test **no fallará**, sino que se marcará como **SKIP**. Esto permite que el flujo de desarrollo local sea ágil incluso sin infraestructura completa.

---

## 🧹 Limpieza de Salida (Logs)

Para facilitar la identificación de errores reales, hemos configurado GORM para ser silencioso durante los tests unitarios que utilizan `sqlmock`. Esto evita que los logs se llenen de mensajes de "db error" cuando forzamos fallos deliberados para probar el manejo de errores.

---

## 🔍 Mejores Prácticas

1.  **Independencia:** Cada test debe ser capaz de ejecutarse de forma aislada.
2.  **Limpieza:** Utilizar `t.Cleanup()` o `defer` para cerrar conexiones y limpiar datos de prueba.
3.  **Aserciones Claras:** Utilizar mensajes descriptivos en `t.Fatalf` o `t.Errorf` para facilitar el debug.
4.  **Esquemas Dinámicos:** Evitar depender de datos pre-existentes en la base de datos; el test debe crear su propio escenario.

---

**Relacionado:**
- [ADR-011: Estrategia de Testing y Cobertura](../../architecture/adrs/adr-011-testing-coverage-strategy.md)
- [Guía de CI/CD](ci-cd.md)
