# ADR-008 – Planificación y Cronograma MVP Ajustado a Disponibilidad Real

**Fecha:** 11/01/2026  
**Estado:** Aceptado  
**Autores:** Jorge Cortés Villalba, Claude (Anthropic)

---

## 1. Contexto

TramaTex requiere una planificación realista que considere:

- **Disponibilidad real:** 8 horas/semana
- **Equipo:** 1 desarrollador full-stack con IA como copiloto
- **Experiencia:** Poca experiencia con el stack (Go, Vue.js, PostgreSQL)
- **Complejidad:** MVP completo según ADR-007 (4 fases)
- **Compromiso:** Disponibilidad continuada durante todo el proyecto

**Restricciones identificadas:**

- Curva de aprendizaje en tecnologías nuevas
- Desarrollo paralelo backend + frontend por una sola persona
- Testing obligatorio (TDD) que incrementa tiempo de desarrollo
- Sin infraestructura previa montada

**Estimación original ADR-007:** 11-16 semanas (equipo tiempo completo, ~40h/semana)

---

## 2. Cálculo de Esfuerzo Ajustado

### 2.1 Conversión de Estimaciones

**Estimación original:**

- Mínimo: 11 semanas × 40h = 440 horas
- Máximo: 16 semanas × 40h = 640 horas
- Promedio: **540 horas**

**Factor de ajuste por curva de aprendizaje:**

- Desarrollador experimentado: 1.0×
- Desarrollador con poca experiencia + IA copiloto: **1.3×**

**Esfuerzo ajustado total:**

- 540h × 1.3 = **702 horas**

### 2.2 Traducción a Calendario

Con disponibilidad de **8 horas/semana**:

- 702h ÷ 8h/semana = **87.75 semanas**
- **≈ 88 semanas = 22 meses**

**Ajuste por buffers y contingencias (+15%):**

- 88 semanas × 1.15 = **101 semanas**
- **≈ 24 meses (2 años)**

---

## 3. Decisión Adoptada

Se adopta un **cronograma de 24 meses** dividido en **4 fases principales** más una fase de estabilización, con hitos verificables cada 2-3 meses.

---

## 4. Hitos y Entregables

### Q1 2026 (Enero - Marzo)
- **Hito 1:** Sistema arranca con Docker Compose  
- **Hito 2:** Autenticación JWT funcional  
- **Entregable:** Fase 0 completa

### Q2 2026 (Abril - Junio)
- **Hito 3:** Backend Party completo  
- **Hito 4:** CRUD Party frontend + backend  
- **Entregable:** Módulo Party operativo

### Q3 2026 (Julio - Septiembre)
- **Hito 5:** CRUD Producto completo  
- **Entregable:** Módulo Producto operativo

### Q4 2026 (Octubre - Diciembre)
- **Hito 6:** Núcleo económico (tarificación) funcional  
- **Entregable:** Fase 1 completa - Dominio base sólido

### Q1 2027 (Enero - Marzo)
- **Hito 7:** Backend Pedidos completo  
- **Hito 8:** Creación de pedidos funcional  
- **Entregable:** Módulo Pedidos operativo

### Q2 2027 (Abril - Junio)
- **Hito 9:** Flujo completo ventas con documentos  
- **Entregable:** Fase 2 completa - Ventas operativas

### Q3 2027 (Julio - Septiembre)
- **Hito 10:** Backend MES completo  
- **Hito 11:** Gestión documental funcional  
- **Entregable:** Infraestructura MES lista

### Q4 2027 (Octubre - Diciembre)
- **Hito 12:** MVP completo (todas las funcionalidades)  
- **Entregable:** Fase 3 completa - MES operativo

### Q1 2028 (Enero)
- **Hito Final:** Proyecto en producción estable  
- **Entregable:** MVP validado y aceptado por cliente

---

## 5. Consecuencias

### Positivas

- **Cronograma realista:** Considera disponibilidad real y curva de aprendizaje
- **Hitos verificables:** Cada 2-3 meses hay entregable demostrable
- **Ritmo sostenible:** 8h/semana permite mantener calidad sin burnout
- **Flexibilidad controlada:** Buffers del 15% permiten absorber imprevistos
- **Aprendizaje gradual:** Fases permiten dominar tecnologías progresivamente

### Negativas

- **Duración extendida:** 24 meses es largo para un MVP
- **Riesgo de obsolescencia:** Tecnologías pueden evolucionar en 2 años
- **Riesgo de cambio de requisitos:** Cliente puede cambiar necesidades
- **Mantenimiento de motivación:** Proyecto largo requiere disciplina constante
- **Dependencia de disponibilidad:** Cualquier cambio en disponibilidad impacta fuertemente

---

## Referencias

- Documento Consolidado TramaTex v3.0
- ADR-001 a ADR-007

---

**Fin del ADR-008**
