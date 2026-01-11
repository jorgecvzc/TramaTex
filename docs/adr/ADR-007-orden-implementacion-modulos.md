# ADR-007 – Orden de Implementación de Módulos e Infraestructura (Revisión Final)

**Fecha:** 10/01/2026  
**Estado:** Aceptado 
**Autores:** Jorge Cortés Villalba, Claude  

---

## Contexto

TramaTex es un ERP/MES para microempresas que opera bajo las siguientes decisiones estratégicas:

- Stack tecnológico definido: Go, PostgreSQL, Vue.js, Docker (ADR-001)
- Arquitectura basada en **DDD + Clean Architecture con rigor asimétrico** (ADR-002)
- Monolito modular local-first, con dominios y subdominios claramente delimitados (ADR-003)
- Ciclo de vida incremental hasta MVP, con entregas funcionales y testeables (ADR-004)
- Estrategia de desarrollo dirigida por dominio, priorizando núcleo económico (tarificación) y dominios base (Party, Producto) (ADR-005, ADR-006)

Se requiere un **orden de implementación** que garantice:

- Dominio protegido y testable desde el inicio
- Integración progresiva de infraestructura
- MVP funcional desde las primeras fases
- Preparación para iteraciones futuras de manera controlada

---

## Decisión

Se adopta la siguiente **secuencia de implementación de módulos y creación de infraestructura**, organizada por **fases alineadas con el ciclo de vida MVP**. El **frontend se desarrolla en paralelo**, de modo que al final de cada fase ya haya interfaz operativa para los módulos implementados.

---

## Conceptos Clave

### ¿Qué es Party?

**Party** es un **patrón de modelado de dominio** que representa a cualquier **persona u organización** que tiene una relación con el sistema, independientemente del rol que desempeñe.

#### Problema que resuelve

**Enfoque tradicional (sin Party):**

- Duplicación de entidades: Cliente y Proveedor como tablas separadas
- Si una empresa es ambos roles: duplicación de datos (nombre, dirección, email)
- Inconsistencias al actualizar información
- Rigidez al cambiar roles

**Enfoque Party:**

- **Un Party, múltiples roles**: Una entidad puede ser Cliente, Proveedor, o ambos simultáneamente
- **Datos únicos**: Información común (nombre, dirección, contacto) se almacena una sola vez
- **Flexibilidad**: Activar/desactivar roles sin duplicar entidades
- **Consistencia**: Cambiar datos actualiza para todos los roles

#### Estructura en TramaTex

```
Party (Entidad Raíz)
├── ID, Tipo (Persona/Organización), Nombre, NIF/CIF
├── Dirección, Email, Teléfono
└── Roles activos

PartyRole
├── Tipo de rol (Cliente, Proveedor)
├── Estado (Activo/Inactivo)
└── Fechas (inicio, fin)

Customer (Especialización de PartyRole: Cliente)
├── Descuento base
├── Empresa matriz (jerarquía)
├── Límite de crédito
└── Condiciones de pago

Supplier (Especialización de PartyRole: Proveedor)
├── Código de proveedor
├── Días de entrega
├── Pedido mínimo
└── Costes base por producto/variante
```

#### Reglas de negocio específicas

**Clientes:**

- Pueden tener jerarquía empresarial (empresa matriz → dependientes)
- Los descuentos pueden heredarse de la empresa matriz
- Ejemplo: "Construcciones ABC S.L." (matriz) con descuento 10%, sus obras heredan ese descuento

**Proveedores:**

- **NO tienen jerarquía** (simplificación para MVP)
- Proporcionan costes base para productos/variantes
- Estos costes son inputs para el motor de tarificación

#### Ejemplo práctico

**Empresa que es Cliente Y Proveedor:**

```
Party: "Bordados Levante S.L."
├── Rol Cliente: Descuento 5%, Crédito 10.000€
└── Rol Proveedor: Código BL-001, Entrega 7 días
```

- Aparece **una sola vez** en Party
- Datos comunes no se duplican
- Cada rol tiene su información específica

---

## Fases de Implementación

### Fase 0 – Fundaciones Técnicas

**Objetivo:** Preparar la infraestructura mínima y estructura base del proyecto **sin implementar lógica de dominio**.

**Acciones:**

- Crear repositorios Git y estructura de carpetas según Clean Architecture
- Configurar Docker Compose y entorno local (PostgreSQL, contenedores mínimos)
- Implementar esqueleto de **autenticación y autorización básica** (JWT, roles: Admin, Comercial, Diseño, Taller)
- Configurar pipeline básico de calidad: TDD, cobertura de tests inicial, linters
- Preparar infraestructura de migraciones de base de datos (sin dominios de negocio)

**Módulos transversales incluidos:**

- **Seguridad básica**: Login, gestión de sesiones JWT, RBAC mínimo

**Frontend paralelo:**

- Pantalla de login funcional
- Estructura de navegación base (menús, rutas)
- Layout general de la aplicación
- **Sin módulos de negocio funcionales**

**Criterios de aceptación Fase 0:**

- ✅ Sistema arranca con `docker-compose up`
- ✅ Login funcional con usuario de prueba
- ✅ Pipeline de tests ejecutándose correctamente
- ✅ Estructura Clean Architecture verificable en código

---

### Fase 1 – Dominio Base para Tarificación

**Objetivo:** Construir los módulos necesarios para que el núcleo económico (tarificación) sea funcional.

### Fase 2 – Casos de Uso Core y Pedidos

**Objetivo:** Orquestar los flujos de negocio críticos usando los módulos base ya estables.

### Fase 3 – Subdominio Secundario MVP: MES

**Objetivo:** Completar el subdominio necesario para gestionar pedidos personalizados y alcanzar el **MVP completo**.

---

## Entregables Finales del MVP

Al completar la Fase 3, el sistema debe:

1. **Operar completamente el flujo de negocio real**:
    
    - Pedidos estándar: desde cotización hasta entrega
    - Pedidos personalizados: desde diseño hasta producción y entrega
2. **Ser estable y mantenible**:
    
    - Cobertura de tests global **≥75%**
    - Dominio crítico (tarificación) con cobertura **≥80%**
    - Documentación técnica actualizada
3. **Estar desplegado en infraestructura local**:
    
    - Docker Compose funcional
    - Backups automáticos configurados
    - Procedimientos de recuperación documentados
4. **Ser operado por usuarios reales**:
    
    - Comerciales gestionan pedidos
    - Taller usa terminal para producción
    - Administrador gestiona catálogo y tarifas

---

## Estrategia de Frontend Paralelo

En **todas las fases**, el desarrollo frontend se realiza simultáneamente con el backend, siguiendo estos principios:

1. **Contratos primero**: Backend expone APIs REST con contratos JSON claros
2. **Validación temprana**: Cada caso de uso se valida con interfaz funcional
3. **Iteración rápida**: Frontend consume APIs reales, no mocks
4. **Feedback inmediato**: Usuarios pueden probar funcionalidad al final de cada fase

**Ventaja clave:** Al finalizar cada fase, los módulos implementados ya son **operativos y utilizables**, no solo "técnicamente completos".

---

## Consecuencias

### Positivas

- MVP operativo y funcional desde fases tempranas
- Dominio crítico (tarificación) protegido y testeable
- Integración progresiva de infraestructura evita sobrecarga inicial
- Frontend disponible al final de cada fase para validar funcionalidad
- Modularidad preparada para evolución futura
- Reducción de riesgo de entrega: cada fase es un incremento estable

### Negativas

- Secuenciación estricta puede ralentizar desarrollo inicial de UI completo
- Dependencias del core (Party, Producto) retrasan implementación de dominios secundarios
- Requiere disciplina estricta y cumplimiento de TDD
- Mayor esfuerzo inicial en diseño y análisis

Estas consecuencias se aceptan explícitamente como parte de la estrategia de calidad del proyecto.

---

## Alcance

Este ADR aplica a:

- Orden de desarrollo de módulos y subdominios del MVP
- Secuencia de integración de infraestructura y frontend
- Criterios de validación de fases completadas
- Estrategia de entrega incremental

**Cualquier desviación en este orden requiere un nuevo ADR con justificación explícita.**

---

## Referencias

- ADR-001: Selección del Stack Tecnológico
- ADR-002: Adopción de Clean Architecture y DDD
- ADR-003: Tipo y Distribución de la Aplicación (Monolito Modular)
- ADR-004: Ciclo de Vida de Desarrollo e Implementación hasta MVP
- ADR-005: Gestión Unificada de Clientes y Proveedores (Party / Organización)
- ADR-006: Estrategia de Desarrollo Dirigida por Dominio

---

**Fin del ADR-007**
