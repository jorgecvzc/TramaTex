# Guía de Ubicación de Tests en el Proyecto TramaTex

Este documento detalla dónde se deben ubicar los diferentes tipos de tests en el proyecto TramaTex, complementando la decisión arquitectónica definida en [ADR-009: Estructura de Carpetas y Organización del Proyecto](../../architecture/adrs/adr-009-project-structure.md) y la estrategia de testing de [ADR-011: Estrategia de Testing y Cobertura](../../architecture/adrs/adr-011-testing-coverage-strategy.md).

---

## Ubicación de Tests

### Tests Unitarios

**Ubicación:** Mismo paquete o directorio que el código que prueban.

**Ejemplo:**
```
domain/party/
├── party.go
└── party_test.go
```

### Tests de Integración

**Ubicación:** Mismo paquete o directorio que el adaptador o la implementación de infraestructura.

**Ejemplo:**
```
infrastructure/persistence/postgres/
├── party_repository.go
└── party_repository_test.go
```

### Tests End-to-End (E2E)

**Ubicación:** Carpeta separada, generalmente en el nivel superior de la aplicación o módulo (`apps/tramatex-api/test/e2e/`). Son opcionales y suelen implementarse en fases más avanzadas.

**Ejemplo:**
```
apps/tramatex-api/
└── test/
    └── e2e/
        └── party_e2e_test.go
```
