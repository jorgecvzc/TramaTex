# ADR-007 – Module Implementation Order (Orden de Implementación de Módulos)

**Fecha:** 2026-02-01 (Reconstruido)  
**Estado:** Aceptado  
**Autores:** Equipo de Arquitectura de TramaTex  

---

## 1. Contexto

El proyecto TramaTex se desarrolla como un monolito modular, siguiendo los principios de Arquitectura Limpia y Diseño Dirigido por Dominio (DDD) según se define en `ADR-002`. En este tipo de arquitectura, los diferentes dominios de negocio (Contextos Delimitados) tienen dependencias explícitas entre sí.

Para evitar dependencias circulares, asegurar que los módulos base estén disponibles antes que los módulos dependientes y alinear el desarrollo con las prioridades del negocio, es fundamental establecer un orden de implementación estricto. Esta decisión documenta la secuencia en la que se construirán los módulos del sistema.

---

## 2. Alternativas Consideradas

**Alternativa A – Desarrollo Ad-hoc**  
- **Descripción:** Implementar módulos según las necesidades inmediatas del sprint sin un plan de dependencias a largo plazo.
- **Ventajas:** Máxima flexibilidad a corto plazo.
- **Desventajas:** Alto riesgo de dependencias circulares, necesidad constante de refactorización, y bloqueos frecuentes cuando un módulo requiere otro que no está listo.

**Alternativa B – Desarrollo por Capas de Dependencia**  
- **Descripción:** Analizar el grafo de dependencias de los contextos delimitados y establecer un orden de implementación desde los más fundamentales (sin dependencias) hacia los más específicos (con muchas dependencias).
- **Ventajas:** Proceso predecible y ordenado, minimiza el retrabajo y los bloqueos, asegura una base estable para el desarrollo incremental.
- **Desventajas:** Menor flexibilidad para cambiar prioridades de negocio que contravengan el orden de dependencias.

---

## 3. Criterios de Decisión

- **Mantenibilidad:** La solución debe reducir la complejidad y facilitar la evolución del sistema.
- **Previsibilidad:** El flujo de trabajo debe ser claro y predecible para el equipo de desarrollo.
- **Alineación con el Negocio:** El orden debe reflejar la criticidad de los módulos para la operación del negocio.
- **Gestión de Dependencias:** La decisión debe resolver explícitamente el grafo de dependencias del dominio.

---

## 4. Decisión Adoptada

Se adopta la **Alternativa B: Desarrollo por Capas de Dependencia**. La implementación de los módulos seguirá un orden estricto basado en su criticidad y sus interdependencias.

El orden de implementación de los módulos será el siguiente:

1.  **Fase 0: Fundaciones**
    -   **Módulo `IAM` (Identity & Access Management):** Es la base de la seguridad y la gestión de usuarios. Todos los demás módulos dependen de él para el seguimiento de la autoría y la autorización.

2.  **Fase 1: Dominio Base**
    -   **Módulo `Party` (Gestión de Entidades):** Define a los clientes y proveedores, un concepto central requerido por `Product`, `Pricing` y `Sales`.
    -   **Módulo `Product` (Catálogo de Productos):** Define qué se vende. Es una dependencia directa para `Pricing` y `Sales`.

3.  **Fase 2: Lógica de Negocio Principal**
    -   **Módulo `Pricing` (Motor de Precios):** Calcula los precios de los productos. Depende de `Product` (costes) y `Party` (precios específicos del cliente). Es económicamente el módulo más crítico.
    -   **Módulo `Sales` (Ventas y Pedidos):** Orquesta el flujo de negocio principal, utilizando todos los módulos anteriores (`IAM`, `Party`, `Product`, `Pricing`).

4.  **Fase 3: Operaciones (MVP)**
    -   **Módulo `MES` (Manufacturing Execution System):** Gestiona la producción en el taller. Depende de `Sales` (para saber qué producir) y `Product`.

**Justificación:**
Esta secuencia respeta estrictamente las dependencias del dominio, garantizando que ninguna implementación comience sin que sus prerrequisitos estén completos y estables. Comienza con la seguridad (`IAM`), establece las entidades de negocio (`Party`, `Product`), construye la lógica económica sobre ellas (`Pricing`), y finalmente orquesta el flujo completo (`Sales` y `MES`).

---

## 5. Consecuencias

### Positivas
- Se elimina el riesgo de bloqueos por dependencias no resueltas.
- El proceso de desarrollo se vuelve más claro y planificable.
- Refuerza la disciplina de la Arquitectura Limpia al hacer explícitas las dependencias entre contextos.
- Permite la validación incremental y estable de cada capa del dominio.

### Negativas
- Reduce la capacidad de implementar una funcionalidad del módulo `Sales` si, por ejemplo, los requerimientos de `Pricing` no se han completado. La planificación de sprints debe adherirse a este orden.

---

## 6. Integración con otros ADRs

- **ADR-002 (Arquitectura Limpia y DDD):** Esta decisión es una implementación directa de los principios de DDD, donde se mapea el "Context Map" a un plan de acción.
- **ADR-006 (Domain-Driven Development Strategy):** Proporciona la justificación estratégica para la separación de contextos que este ADR ordena.
- **ADR-009 (Project Structure):** La estructura de directorios (`apps/internal/<module>`) está diseñada para soportar esta implementación modular.

---
*Nota: Este documento ha sido reconstruido el 2026-02-01 a partir de las referencias encontradas en el código y otros documentos del proyecto.*
