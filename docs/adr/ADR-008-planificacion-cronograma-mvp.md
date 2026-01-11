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

### Distribución de Esfuerzo por Fase

|Fase|Esfuerzo Original|Esfuerzo Ajustado|Semanas (8h/sem)|Meses|
|---|---|---|---|---|
|**Fase 0: Fundaciones**|80h|104h|13 semanas|3 meses|
|**Fase 1: Dominio Base**|200h|260h|32.5 semanas|8 meses|
|**Fase 2: Pedidos**|140h|182h|22.75 semanas|6 meses|
|**Fase 3: MES**|120h|156h|19.5 semanas|5 meses|
|**Estabilización**|-|80h|10 semanas|2.5 meses|
|**TOTAL**|**540h**|**782h**|**97.75 semanas**|**≈24 meses**|

---

## 4. Cronograma Detallado

### FASE 0: Fundaciones Técnicas

**Duración:** 13 semanas (3 meses)  
**Período:** Enero 2026 - Marzo 2026  
**Esfuerzo total:** 104 horas

#### Semanas 1-4 (32h): Setup de Proyecto

- **Semana 1 (8h):**
    
    - [4h] Crear repositorio Git y estructura de carpetas
    - [4h] Configurar Docker Compose básico (PostgreSQL)
- **Semana 2 (8h):**
    
    - [6h] Estructura Clean Architecture en Go (skeleton)
    - [2h] Configuración inicial Vue.js 3 + Vite
- **Semana 3 (8h):**
    
    - [4h] Configurar GORM y migraciones base
    - [4h] Setup Tailwind CSS y estructura frontend
- **Semana 4 (8h):**
    
    - [4h] Makefile y scripts de utilidad
    - [4h] README.md y documentación setup

**Hito 1 (Mes 1):** ✅ Proyecto arranca con `docker-compose up`

#### Semanas 5-9 (40h): Autenticación y Seguridad

- **Semana 5 (8h):**
    
    - [6h] Implementar entidad User en dominio
    - [2h] Escribir tests unitarios User
- **Semana 6 (8h):**
    
    - [6h] Implementar JWT (generación y validación)
    - [2h] Tests JWT
- **Semana 7 (8h):**
    
    - [4h] Hash de passwords (bcrypt)
    - [4h] Repositorio User en PostgreSQL
- **Semana 8 (8h):**
    
    - [6h] Middleware de autenticación
    - [2h] RBAC básico (roles: Admin, Comercial, Diseño, Taller)
- **Semana 9 (8h):**
    
    - [6h] Endpoint login/logout
    - [2h] Tests de integración autenticación

**Hito 2 (Mes 2):** ✅ Autenticación JWT funcional

#### Semanas 10-13 (32h): Frontend Base y Testing

- **Semana 10 (8h):**
    
    - [6h] Componente Login.vue
    - [2h] Store Pinia para auth
- **Semana 11 (8h):**
    
    - [4h] Router con rutas protegidas
    - [4h] Layout principal y navegación
- **Semana 12 (8h):**
    
    - [4h] Configurar pipeline CI/CD básico
    - [4h] Linters (golangci-lint, ESLint)
- **Semana 13 (8h):**
    
    - [4h] Tests end-to-end login
    - [4h] Documentación Fase 0

**Criterios de Aceptación Fase 0:**

- ✅ Sistema arranca con Docker Compose
- ✅ Login funcional (backend + frontend)
- ✅ Pipeline de tests ejecutándose
- ✅ Cobertura de tests ≥60% en auth

---

### FASE 1: Dominio Base para Tarificación

**Duración:** 32.5 semanas (8 meses)  
**Período:** Abril 2026 - Noviembre 2026  
**Esfuerzo total:** 260 horas

#### Sprint 1: Party (Semanas 14-21, 64h - 2 meses)

**Semanas 14-15 (16h): Diseño y Modelado**

- **Semana 14 (8h):**
    
    - [4h] Diagrama entidad-relación Party
    - [4h] Definir Value Objects (NIF, Email, Address)
- **Semana 15 (8h):**
    
    - [4h] Documentar reglas de negocio jerarquías
    - [4h] Documentar reglas herencia descuentos

**Semanas 16-19 (32h): Implementación Backend**

- **Semana 16 (8h):**
    
    - [6h] TDD: Entidad Party + tests
    - [2h] TDD: PartyRole + tests
- **Semana 17 (8h):**
    
    - [6h] TDD: Customer (descuentos, jerarquía) + tests
    - [2h] TDD: Supplier (sin jerarquía) + tests
- **Semana 18 (8h):**
    
    - [6h] TDD: SupplierCost + tests
    - [2h] Refactoring y optimización tests
- **Semana 19 (8h):**
    
    - [6h] Repositorio PostgreSQL Party
    - [2h] Migraciones BD Party

**Semanas 20-21 (16h): Casos de Uso y API**

- **Semana 20 (8h):**
    
    - [4h] Casos de uso: CreateParty, UpdateParty
    - [4h] Casos de uso: AssignRole, ManageHierarchy
- **Semana 21 (8h):**
    
    - [4h] Endpoints REST Party
    - [4h] Tests de integración API Party

**Hito 3 (Mes 5):** ✅ Backend Party completo y testeado

#### Sprint 2: Frontend Party (Semanas 22-25, 32h - 1 mes)

- **Semana 22 (8h):**
    
    - [6h] Store Pinia para Party
    - [2h] Service API Party
- **Semana 23 (8h):**
    
    - [6h] Componente PartyList.vue
    - [2h] Componente PartyForm.vue (skeleton)
- **Semana 24 (8h):**
    
    - [6h] PartyForm.vue completo (CRUD)
    - [2h] Validaciones frontend
- **Semana 25 (8h):**
    
    - [4h] Componente PartyHierarchy.vue (árbol jerárquico)
    - [4h] Tests E2E Party

**Hito 4 (Mes 6):** ✅ CRUD Party completo (frontend + backend)

#### Sprint 3: Producto (Semanas 26-32, 56h - 1.75 meses)

**Semanas 26-27 (16h): Diseño y Backend**

- **Semana 26 (8h):**
    
    - [4h] Diagrama ER Producto/Variante/Categoría
    - [4h] TDD: Entidad Product + tests
- **Semana 27 (8h):**
    
    - [4h] TDD: Variant + modificadores precio + tests
    - [4h] TDD: Category + tests

**Semanas 28-29 (16h): Repositorio y API**

- **Semana 28 (8h):**
    
    - [6h] Repositorio PostgreSQL Product
    - [2h] Migraciones BD Product
- **Semana 29 (8h):**
    
    - [4h] Casos de uso Product
    - [4h] Endpoints REST Product

**Semanas 30-32 (24h): Frontend Producto**

- **Semana 30 (8h):**
    
    - [6h] Store Pinia Product
    - [2h] Service API Product
- **Semana 31 (8h):**
    
    - [6h] ProductList.vue + ProductForm.vue
    - [2h] Gestión de variantes en formulario
- **Semana 32 (8h):**
    
    - [4h] ProductCatalog.vue (vista catálogo)
    - [4h] Tests E2E Product

**Hito 5 (Mes 8):** ✅ CRUD Producto completo

#### Sprint 4: Tarificación (Semanas 33-45, 108h - 3.25 meses)

**Semanas 33-35 (24h): Diseño Motor Tarificación**

- **Semana 33 (8h):**
    
    - [8h] Especificación detallada fórmula de cálculo
- **Semana 34 (8h):**
    
    - [8h] Documentar reglas de márgenes y descuentos
- **Semana 35 (8h):**
    
    - [8h] Casos límite y validaciones

**Semanas 36-40 (40h): Implementación Backend (TDD estricto)**

- **Semana 36 (8h):**
    
    - [6h] TDD: PricingEngine (cálculo base) + tests
    - [2h] TDD: Money Value Object + tests
- **Semana 37 (8h):**
    
    - [6h] TDD: Aplicación márgenes + tests
    - [2h] TDD: Percentage Value Object + tests
- **Semana 38 (8h):**
    
    - [6h] TDD: Aplicación descuentos simples + tests
    - [2h] Refactoring
- **Semana 39 (8h):**
    
    - [6h] TDD: Herencia de descuentos + tests
    - [2h] TDD: Modificadores de variante + tests
- **Semana 40 (8h):**
    
    - [6h] TDD: Integración completa motor + tests
    - [2h] Tests casos límite (márgenes negativos, etc.)

**Semanas 41-42 (16h): Repositorio y API**

- **Semana 41 (8h):**
    
    - [6h] Repositorio PricingRules
    - [2h] Migraciones BD Pricing
- **Semana 42 (8h):**
    
    - [4h] Casos de uso Pricing
    - [4h] Endpoint CalculatePrice

**Semanas 43-45 (28h): Frontend Tarificación**

- **Semana 43 (8h):**
    
    - [6h] Store Pinia Pricing
    - [2h] Service API Pricing
- **Semana 44 (8h):**
    
    - [8h] PricingCalculator.vue (interfaz calculadora)
- **Semana 45 (12h):**
    
    - [6h] Desglose visual (coste + margen - descuento)
    - [4h] Tests E2E Tarificación
    - [2h] Auditoría básica cambios en precios

**Criterios de Aceptación Fase 1:**

- ✅ CRUD Party funcional (frontend + backend)
- ✅ CRUD Producto funcional
- ✅ Cálculo de tarificación con datos reales
- ✅ Cobertura tests ≥80% en tarificación
- ✅ Jerarquías funcionan correctamente

**Hito 6 (Mes 11):** ✅ Núcleo económico funcional y testeado

---

### FASE 2: Pedidos y Ventas

**Duración:** 22.75 semanas (6 meses)  
**Período:** Diciembre 2026 - Mayo 2027  
**Esfuerzo total:** 182 horas

#### Sprint 5: Pedidos Backend (Semanas 46-53, 64h - 2 meses)

**Semanas 46-47 (16h): Diseño**

- **Semana 46 (8h):**
    
    - [4h] Diagrama ER Pedido/Línea
    - [4h] Máquina de estados Pedido
- **Semana 47 (8h):**
    
    - [8h] Documentar casos de uso ventas

**Semanas 48-51 (32h): Implementación TDD**

- **Semana 48 (8h):**
    
    - [6h] TDD: Order (agregado raíz) + tests
    - [2h] TDD: OrderLine + tests
- **Semana 49 (8h):**
    
    - [6h] TDD: OrderState (máquina estados) + tests
    - [2h] Tests transiciones válidas/inválidas
- **Semana 50 (8h):**
    
    - [6h] TDD: Integración con Tarificación + tests
    - [2h] TDD: Validación de crédito + tests
- **Semana 51 (8h):**
    
    - [6h] TDD: Aplicación descuentos cliente + tests
    - [2h] Refactoring y optimización

**Semanas 52-53 (16h): Repositorio y API**

- **Semana 52 (8h):**
    
    - [6h] Repositorio Order
    - [2h] Migraciones BD Order
- **Semana 53 (8h):**
    
    - [4h] Casos de uso Order
    - [4h] Endpoints REST Order

**Hito 7 (Mes 13):** ✅ Backend Pedidos completo

#### Sprint 6: Frontend Pedidos (Semanas 54-60, 56h - 1.75 meses)

**Semanas 54-55 (16h): Store y Services**

- **Semana 54 (8h):**
    
    - [6h] Store Pinia Sales
    - [2h] Service API Sales
- **Semana 55 (8h):**
    
    - [8h] Composables useOrder, useOrderCalculation

**Semanas 56-58 (24h): Interfaz de Pedidos**

- **Semana 56 (8h):**
    
    - [6h] OrderList.vue (listado + filtros)
    - [2h] OrderDetail.vue (visualización)
- **Semana 57 (8h):**
    
    - [8h] OrderForm.vue (creación pedido - parte 1)
- **Semana 58 (8h):**
    
    - [6h] OrderForm.vue (líneas de pedido - parte 2)
    - [2h] Cálculo automático totales

**Semanas 59-60 (16h): Gestión Estados y Tests**

- **Semana 59 (8h):**
    
    - [4h] Cambio de estados en interfaz
    - [4h] Validaciones frontend
- **Semana 60 (8h):**
    
    - [6h] Tests E2E Pedidos
    - [2h] Ajustes UX

**Hito 8 (Mes 15):** ✅ Creación de pedidos funcional

#### Sprint 7: Generación Documentos (Semanas 61-67, 62h - 2 meses)

**Semanas 61-63 (24h): Web-to-Print**

- **Semana 61 (8h):**
    
    - [8h] Investigación y setup Web-to-Print
- **Semana 62 (8h):**
    
    - [6h] Plantilla CSS Presupuesto
    - [2h] Generación PDF presupuesto
- **Semana 63 (8h):**
    
    - [6h] Plantilla CSS Albarán
    - [2h] Generación PDF albarán

**Semanas 64-66 (24h): Facturación**

- **Semana 64 (8h):**
    
    - [6h] Plantilla CSS Factura proforma
    - [2h] Generación PDF factura
- **Semana 65 (8h):**
    
    - [4h] Botones imprimir/descargar documentos
    - [4h] Integración con estados de pedido
- **Semana 66 (8h):**
    
    - [6h] Preview documentos antes de generar
    - [2h] Tests generación documentos

**Semana 67 (14h): Finalización Fase 2**

- [6h] Tests de integración completos
- [4h] Documentación Fase 2
- [4h] Demo y validación con cliente

**Criterios de Aceptación Fase 2:**

- ✅ Crear pedido completo desde frontend
- ✅ Precio calculado automáticamente
- ✅ Estados de pedido funcionan
- ✅ Documentos PDF funcionales
- ✅ Cobertura tests ≥70% en casos de uso

**Hito 9 (Mes 17):** ✅ Flujo completo de ventas operativo

---

### FASE 3: MES - Producción Personalizada

**Duración:** 19.5 semanas (5 meses)  
**Período:** Junio 2027 - Octubre 2027  
**Esfuerzo total:** 156 horas

#### Sprint 8: MES Backend (Semanas 68-74, 56h - 1.75 meses)

**Semanas 68-69 (16h): Diseño**

- **Semana 68 (8h):**
    
    - [4h] Diagrama ER Pedido Personalizado
    - [4h] Máquina estados Producción
- **Semana 69 (8h):**
    
    - [4h] Documentar flujo MES completo
    - [4h] Casos de uso MES

**Semanas 70-72 (24h): Implementación TDD**

- **Semana 70 (8h):**
    
    - [6h] TDD: CustomOrder (extensión Order) + tests
    - [2h] TDD: ProductionState + tests
- **Semana 71 (8h):**
    
    - [6h] TDD: WorkshopJob + tests
    - [2h] TDD: Transiciones estados MES + tests
- **Semana 72 (8h):**
    
    - [6h] TDD: Trazabilidad procesos + tests
    - [2h] Refactoring

**Semanas 73-74 (16h): Repositorio y API**

- **Semana 73 (8h):**
    
    - [6h] Repositorio MES
    - [2h] Migraciones BD MES
- **Semana 74 (8h):**
    
    - [4h] Casos de uso MES
    - [4h] Endpoints REST MES

**Hito 10 (Mes 19):** ✅ Backend MES completo

#### Sprint 9: Gestión Documental (Semanas 75-78, 32h - 1 mes)

**Semanas 75-76 (16h): Storage NAS**

- **Semana 75 (8h):**
    
    - [6h] Adaptador almacenamiento NAS
    - [2h] Tests almacenamiento
- **Semana 76 (8h):**
    
    - [4h] Indexación PostgreSQL (ruta, tipo, tamaño)
    - [4h] Entidad DesignFile + tests

**Semanas 77-78 (16h): API Documentos**

- **Semana 77 (8h):**
    
    - [4h] Endpoint upload archivo
    - [4h] Endpoint download archivo
- **Semana 78 (8h):**
    
    - [4h] Asociar archivos a pedidos
    - [4h] Tests gestión documental

**Hito 11 (Mes 20):** ✅ Gestión documental funcional

#### Sprint 10: Terminal Taller (Semanas 79-86, 68h - 2.25 meses)

**Semanas 79-81 (24h): Frontend MES Desktop**

- **Semana 79 (8h):**
    
    - [6h] Store Pinia MES
    - [2h] Service API MES
- **Semana 80 (8h):**
    
    - [6h] ProductionDashboard.vue
    - [2h] Visualización estados producción
- **Semana 81 (8h):**
    
    - [4h] Gestión pedidos personalizados
    - [4h] Adjuntar diseños a pedidos

**Semanas 82-84 (24h): Terminal Tablet**

- **Semana 82 (8h):**
    
    - [8h] WorkshopTerminal.vue (interfaz simplificada)
- **Semana 83 (8h):**
    
    - [4h] Listado trabajos pendientes (touch-optimized)
    - [4h] Cambio estado con un toque
- **Semana 84 (8h):**
    
    - [4h] Visualización diseños en terminal
    - [4h] Captura observaciones/fotos

**Semanas 85-86 (20h): Visor y Tests**

- **Semana 85 (8h):**
    
    - [6h] Visor de archivos (PDF, imágenes)
    - [2h] Descarga diseños para producción
- **Semana 86 (12h):**
    
    - [6h] Tests E2E MES completo
    - [4h] Tests terminal en tablet real
    - [2h] Optimización rendimiento tablet

**Criterios de Aceptación Fase 3:**

- ✅ Terminal taller operativa en tablet
- ✅ Pedidos personalizados end-to-end
- ✅ Almacenamiento NAS funcional
- ✅ Diseños visualizables y descargables
- ✅ Trazabilidad completa de producción

**Hito 12 (Mes 22):** ✅ MVP COMPLETO - Todas las funcionalidades implementadas

---

### FASE 4: Estabilización y Despliegue

**Duración:** 10 semanas (2.5 meses)  
**Período:** Noviembre 2027 - Enero 2028  
**Esfuerzo total:** 80 horas

#### Semanas 87-89 (24h): Preparación Producción

- **Semana 87 (8h):**
    
    - [4h] Docker Compose para producción
    - [4h] Scripts backup automático PostgreSQL
- **Semana 88 (8h):**
    
    - [4h] Procedimientos de recuperación
    - [4h] Documentación despliegue
- **Semana 89 (8h):**
    
    - [4h] Monitoreo básico (logs, errores)
    - [4h] Configuración alertas críticas

#### Semanas 90-92 (24h): Documentación y Capacitación

- **Semana 90 (8h):**
    
    - [8h] Manual de usuario (Comerciales)
- **Semana 91 (8h):**
    
    - [4h] Manual de usuario (Taller)
    - [4h] Manual de usuario (Admin)
- **Semana 92 (8h):**
    
    - [4h] Sesiones de formación con usuarios
    - [4h] Videos tutoriales básicos

#### Semanas 93-96 (32h): Estabilización

- **Semanas 93-94 (16h):**
    
    - [16h] Uso supervisado con cliente real
    - Identificación y corrección de bugs
- **Semanas 95-96 (16h):**
    
    - [8h] Ajustes de usabilidad según feedback
    - [4h] Optimizaciones de rendimiento
    - [4h] Documentación final y cierre

**Criterios de Finalización del Proyecto:**

- ✅ MVP operativo en producción
- ✅ 2 semanas de uso estable sin incidencias críticas
- ✅ Cobertura tests ≥75% global
- ✅ Cobertura tests ≥80% en tarificación
- ✅ Backups automáticos funcionando
- ✅ Usuarios operan sistema sin soporte

**Hito Final (Mes 24):** ✅ PROYECTO MVP COMPLETADO Y EN PRODUCCIÓN

---

## 5. Hitos y Entregables por Trimestre

### Q1 2026 (Enero - Marzo)

**Hito 1:** Sistema arranca con Docker Compose  
**Hito 2:** Autenticación JWT funcional  
**Entregable:** Fase 0 completa

### Q2 2026 (Abril - Junio)

**Hito 3:** Backend Party completo  
**Hito 4:** CRUD Party frontend + backend  
**Entregable:** Módulo Party operativo

### Q3 2026 (Julio - Septiembre)

**Hito 5:** CRUD Producto completo  
**Entregable:** Módulo Producto operativo

### Q4 2026 (Octubre - Diciembre)

**Hito 6:** Núcleo económico (tarificación) funcional  
**Entregable:** Fase 1 completa - Dominio base sólido

### Q1 2027 (Enero - Marzo)

**Hito 7:** Backend Pedidos completo  
**Hito 8:** Creación de pedidos funcional  
**Entregable:** Módulo Pedidos operativo

### Q2 2027 (Abril - Junio)

**Hito 9:** Flujo completo ventas con documentos  
**Entregable:** Fase 2 completa - Ventas operativas

### Q3 2027 (Julio - Septiembre)

**Hito 10:** Backend MES completo  
**Hito 11:** Gestión documental funcional  
**Entregable:** Infraestructura MES lista

### Q4 2027 (Octubre - Diciembre)

**Hito 12:** MVP completo (todas las funcionalidades)  
**Entregable:** Fase 3 completa - MES operativo

### Q1 2028 (Enero)

**Hito Final:** Proyecto en producción estable  
**Entregable:** MVP validado y aceptado por cliente

---

## 6. Estrategia de Gestión de Riesgos

### Riesgos Identificados y Mitigaciones

|Riesgo|Probabilidad|Impacto|Mitigación|
|---|---|---|---|
|**Curva de aprendizaje mayor de lo esperado**|Alta|Alto|- Buffer del 15% incluido<br>- IA copiloto para acelerar<br>- Documentación detallada por módulo<br>- Fases permiten aprendizaje gradual|
|**Disponibilidad reducida (< 8h/semana)**|Media|Alto|- Cronograma con 10% de margen adicional<br>- Hitos cada 2-3 meses permiten ajustes<br>- Priorización estricta de MVP|
|**Scope creep (añadir funcionalidades)**|Alta|Crítico|- ADR-004: MVP cerrado, Post-MVP = nuevo proyecto<br>- Revisión de alcance cada hito<br>- Stakeholder consciente de cronograma|
|**Bugs críticos en producción**|Media|Alto|- TDD obligatorio<br>- Cobertura ≥75%<br>- Fase de estabilización de 2.5 meses<br>- Uso supervisado antes de producción|
|**Problemas de integración entre módulos**|Media|Medio|- Clean Architecture separa concerns<br>- Tests de integración en cada sprint<br>- Contratos API claros desde diseño|
|**Cambios en requisitos del cliente**|Media|Medio|- Validación cada trimestre<br>- Hitos permiten pivotes controlados<br>- Documentación versionada|
|**Agotamiento (burnout)**|Media|Medio|- Ritmo sostenible (8h/semana)<br>- Hitos pequeños dan sensación de logro<br>- Descansos entre fases|

---

## 7. Métricas de Seguimiento

### KPIs por Fase

|Métrica|Objetivo|Frecuencia Medición|
|---|---|---|
|**Horas trabajadas/semana**|8h ± 1h|Semanal|
|**Cobertura de tests**|≥75% global, ≥80% tarificación|Al finalizar cada sprint|
|**Velocidad de desarrollo**|Ajustar estimaciones cada hito|Cada hito (2-3 meses)|
|**Bugs críticos abiertos**|0 antes de pasar a siguiente fase|Continuo|
|**Deuda técnica acumulada**|< 10% del tiempo por sprint|Al finalizar cada sprint|
|**Cumplimiento de hitos**|100% (con buffers)|Cada hito|

### Herramientas de Seguimiento Recomendadas

- **Time tracking:** Toggl / Clockify (registrar horas semanales)
- **Gestión tareas:** GitHub Projects (tablero Kanban por sprint)
- **Cobertura tests:** Codecov / SonarQube
- **Documentación sesiones:** Carpeta `/docs/sessions/` con registro semanal

---

## 8. Criterios de Ajuste del Cronograma

El cronograma se revisará cada **hito** (2-3 meses) y se ajustará si:

### Adelanto (mejor de lo esperado)

- Velocidad real > estimada en 20%+
- **Acción:** NO reducir cronograma, usar tiempo para:
    - Mejorar cobertura de tests
    - Refactoring y optimización
    - Documentación adicional
    - Exploración de tecnologías para Post-MVP

### Retraso (peor de lo esperado)

- Velocidad real < estimada en 20%+
- **Acción:**
    1. Analizar causa raíz (¿curva aprendizaje? ¿complejidad subestimada?)
    2. Ajustar estimaciones de sprints restantes
    3. Evaluar si es necesario simplificar alcance MVP (mantener núcleo económico intacto)
    4. Comunicar al stakeholder nuevo cronograma

### Cambios en Disponibilidad

- Si disponibilidad baja a < 6h/semana: recalcular cronograma completo
- Si disponibilidad sube a > 10h/semana: mantener cronograma (usar para calidad)

---

## 9. Estructura de Trabajo Semanal Recomendada

Con **8 horas/semana**, se recomienda distribución:

### Opción A: 2 sesiones de 4h (recomendada)

- **Sesión 1 (4h):** Backend (dominio + tests)
- **Sesión 2 (4h):** Frontend o infraestructura

**Ventaja:** Permite cambio de contexto, reduce fatiga mental

### Opción B: 1 sesión de 8h

- **Domingo/Sábado (8h):** Sprint completo (backend + frontend)

**Ventaja:** Contexto continuo, menos overhead de setup

### Distribución de Tiempo Típica (por sesión de 8h)

```
[1.5h] Desarrollo con TDD (tests primero)
[1.5h] Implementación funcionalidad
[1.0h] Refactoring y optimización
[1.0h] Documentación (código, ADRs, sesión)
[0.5h] Code review y ajustes
[0.5h] Commits y push
[2.0h] Frontend (si aplica en la fase)
```

---

## 10. Plantilla de Documentación de Sesión

Cada semana, documentar en `/docs/sessions/YYYY-MM-DD-session-NN.md`:

```markdown
# Sesión NN - [Fecha]

## Contexto
- **Fase actual:** X
- **Sprint actual:** Y
- **Semana del cronograma:** NN

## Objetivos de la sesión
- [ ] Objetivo 1
- [ ] Objetivo 2
- [ ] Objetivo 3

## Trabajo realizado
- [Xh] Descripción tarea 1
- [Xh] Descripción tarea 2

## Decisiones tomadas
- Decisión 1: ...
- Decisión 2: ...

## Problemas encontrados
- Problema 1: ... → Solución: ...
- Problema 2: ... → Pendiente resolver

## Aprendizajes
- Aprendizaje técnico: ...
- Aprendizaje de dominio: ...

## Próxima sesión
- [ ] Tarea 1
- [ ] Tarea 2

## Métricas
- Horas trabajadas: Xh
- Cobertura tests: X%
- Commits: N
```

---

## 11. Consecuencias

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

### Mitigaciones de Consecuencias Negativas

- **Obsolescencia:** Stack elegido es estable (Go, PostgreSQL, Vue 3)
- **Cambios requisitos:** Validación trimestral con stakeholder
- **Motivación:** Hitos pequeños y frecuentes generan sensación de progreso
- **Disponibilidad:** Cronograma incluye buffers para absorber variaciones

---

## 12. Alcance

Este ADR aplica a:

- Cronograma completo del MVP (24 meses)
- Distribución de esfuerzo por fase y sprint
- Hitos y criterios de aceptación por trimestre
- Estrategia de gestión de riesgos
- Métricas de seguimiento y ajuste

**No aplica a:**

- Post-MVP (será un proyecto independiente según ADR-004)
- Evolución futura a microservicios
- Funcionalidades no incluidas en MVP

---

## 13. Integración con otros ADRs

- **ADR-001:** Stack tecnológico (define qué aprender)
- **ADR-002:** Clean Architecture (define cómo estructurar)
- **ADR-003:** Monolito Modular (define arquitectura a implementar)
- **ADR-004:** Ciclo de vida MVP (define alcance a completar en 24 meses)
- **ADR-005:** Party/Organización (primer dominio a implementar)
- **ADR-006:** Desarrollo dirigido por dominio (estrategia de implementación)
- **ADR-007:** Orden de módulos (secuencia respetada en cronograma)

---

## 14. Notas Adicionales

### Checkpoints Críticos

**Mes 6 (Checkpoint 1):**

- **Evaluar:** ¿La curva de aprendizaje fue la esperada?
- **Decisión:** Ajustar factor de aprendizaje para fases restantes

**Mes 12 (Checkpoint 2):**

- **Evaluar:** ¿El núcleo económico es sólido?
- **Decisión:** ¿Continuar con Fase 2 o reforzar Fase 1?

**Mes 18 (Checkpoint 3):**

- **Evaluar:** ¿El flujo de ventas está validado con cliente?
- **Decisión:** ¿Ajustar UX antes de MES o continuar?

**Mes 22 (Checkpoint 4):**

- **Evaluar:** ¿MVP completo cumple expectativas?
- **Decisión:** ¿Cuánto tiempo de estabilización es necesario?

### Plan de Contingencia

Si en cualquier punto el proyecto debe detenerse o pausarse:

1. **Documentar estado actual:** Qué está completo, qué falta
2. **Generar informe de cierre parcial:** Entregables hasta el momento
3. **Crear roadmap de continuación:** Pasos para retomar
4. **Preservar código y documentación:** Repositorio accesible y documentado

---

## 15. Referencias

- Documento Consolidado TramaTex v3.0
- ADR-001 a ADR-007
- Estimaciones estándar de desarrollo con Clean Architecture
- Factores de productividad en desarrollo en solitario con IA copiloto

---

## 16. Resumen Ejecutivo

**TramaTex MVP será desarrollado en 24 meses (96 semanas) con:**

- **Esfuerzo total:** 782 horas
- **Disponibilidad:** 8 horas/semana
- **Equipo:** 1 desarrollador full-stack + IA copiloto
- **4 Fases principales + Estabilización**
- **12 Hitos verificables cada 2-3 meses**
- **Entrega final:** Enero 2028

**Fecha de inicio:** Enero 2026  
**Fecha de finalización estimada:** Enero 2028

El cronograma incluye:

- ✅ Curva de aprendizaje considerada (+30%)
- ✅ Buffers de contingencia (+15%)
- ✅ TDD obligatorio en dominio
- ✅ Frontend paralelo en todas las fases
- ✅ Validación con cliente cada trimestre
- ✅ Ritmo sostenible (sin burnout)

---

**Fin del ADR-008 – Planificación y Cronograma MVP Ajustado**