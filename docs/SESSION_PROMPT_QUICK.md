# SESSION_PROMPT_QUICK.md – Resumen Rápido

**Versión:** 1.0 (Quick Reference)  
**Para:** Inicio rápido de sesiones sin perder contexto

---

## 🎯 EN UN VISTAZO

**TramaTex** = ERP/MES para microempresas textil (pedidos + producción personalizada)

| Aspecto | Valor |
|--------|-------|
| **Stack** | Go (backend) + Vue.js 3 (frontend) + PostgreSQL |
| **Arquitectura** | Clean Architecture + DDD con rigor asimétrico |
| **Cronograma** | 24 meses (Enero 2026 - Enero 2028), 8h/semana, 782h totales |
| **Estado** | Fase 0 (Fundaciones) – Setup completado, Auth JWT próximo |
| **Sesiones** | 8 completadas (hasta 11/01/2026) |
| **Líneas código** | 0 (puro setup), 1435+ de documentación |

---

## 🏗️ ESTRUCTURA (5 Bounded Contexts)

```
Party (Clientes/Proveedores) → Producto (Variantes) → Tarificación
                          ↓
                       Ventas (Pedidos) → MES (Producción)
```

**Orden implementación:** Party → Producto → Tarificación → Ventas → MES

---

## 📦 FASE ACTUAL

**Fase 0 – Fundaciones Técnicas (4 semanas, 32h)**

✅ Completado:
- Estructura carpetas (Clean Architecture)
- Documentación (9 ADRs + Documento Consolidado 3.0)
- Git setup
- Docker base
- Makefile

⏳ Próximo (Sesión 09):
- Skeleton Go con Clean Architecture
- Setup Vue.js 3 + Vite
- Entidad User (dominio)
- JWT (generación/validación)
- Componente Login
- Tests iniciales

---

## 🔑 PRINCIPIOS CLAVE

1. **Dominio es activo crítico** → TDD obligatorio, cero dependencias externas
2. **Rigor asimétrico** → Dominio crítico estricto, infraestructura flexible
3. **Local-first** → Sin cloud MVP, 100% operativo sin internet
4. **Modular preparado** → Listo para extracción futura a microservicios

---

## 📁 ESTRUCTURA MÍNIMA (GO)

```
backend/internal/
├── domain/              # Sin frameworks
│   ├── common/          # Shared
│   └── [module]/        # party, product, pricing, sales, mes
├── application/         # Use cases
├── infrastructure/      # GORM, Persistencia
└── interfaces/          # Controllers REST
```

---

## ✅ CHECKLIST FIN SESIÓN

- [ ] Tests verdes: `make backend-test`
- [ ] Lint sin warnings: `golangci-lint run ./...`
- [ ] Documentación sesión: `/docs/sessions/2026-MM-DD-session-NN.md`
- [ ] PROJECT_STATUS.md actualizado
- [ ] Commits descriptivos

---

## 🔗 REFERENCIAS RÁPIDAS

| Necesito | Ver |
|---------|-----|
| Stack completo | [ADR-001](docs/adr/ADR-001-seleccion-stack-tecnologico.md) |
| Arquitectura | [ADR-002](docs/adr/ADR-002-adopcion-clean-architecture-ddd.md) |
| Especificación MVP | [DOCUMENTO-CONSOLIDADO-3.0.md](docs/consolidated/DOCUMENTO-CONSOLIDADO-3.0.md) |
| Estado actual | [PROJECT_STATUS.md](PROJECT_STATUS.md) |
| Sesión anterior | [docs/sessions/](docs/sessions/) |

---

## 🚀 COMANDOS ÚTILES

```bash
# Ver todo
make help

# Docker
make docker-up && make docker-down

# Backend
cd backend && make test && make run

# Frontend
cd frontend && npm run dev

# Documentación
make docs-view
```

---

## 📋 FILL THIS FOR EACH SESSION

**Session:** [YYYY-MM-DD Session-NN]  
**LLM Facilitator:** [GitHub Copilot / Claude / ...]  
**Estimated Duration:** [X hours]

**Objectives:**
1. 
2. 
3. 

**Definition of Done:**
- [ ] Objectives completed
- [ ] Tests ≥75% (critical path)
- [ ] Session documented
- [ ] PROJECT_STATUS.md updated

