# ADR-006 – Estrategia de Desarrollo Dirigida por Dominio (MVP)

**Fecha:** 09/01/2026  
**Estado:** Aceptado  
**Autores:** Jorge Cortés Villalba, ChatGPT  

---

## 1. Contexto

En TramaTex, tras la definición del stack tecnológico (ADR-001), la arquitectura DDD + Clean Architecture (ADR-002) y la determinación del tipo de aplicación monolito modular local-first (ADR-003), se identifica que:

- El **principal riesgo** no es tecnológico, sino de ejecución.
- La implementación desordenada de módulos puede **diluir el dominio core** y generar retrabajo.
- Se necesita garantizar **trazabilidad** entre RF/RNF, dominio, código y documentación.

Restricciones y riesgos:

- Introducción prematura de infraestructura o lógica técnica antes de que el dominio lo requiera.
- Desarrollo de la tarificación sin soporte de los dominios Party y Producto.
- Pérdida de coherencia entre módulos y subdominios del MVP.

---

## 2. Alternativas Consideradas

**Alternativa A – Desarrollo basado en capas técnicas**  
- Ventajas: Desarrollo tradicional, rápido para equipos familiarizados con arquitecturas por capas.  
- Desventajas: Riesgo de mezclar dominio y lógica técnica, retrabajo elevado, dominio económico no protegido.  

**Alternativa B – Desarrollo dirigido por el dominio (propuesta)**  
- Ventajas: Protección del dominio, trazabilidad total, implementación incremental validable, menor deuda técnica.  
- Desventajas: Mayor esfuerzo inicial en análisis, ritmo de UI más lento al inicio.  

---

## 3. Decisión Adoptada

Se adopta una **estrategia de desarrollo incremental dirigida por el dominio** para el MVP.  

**Justificación:**  
- Garantiza que los módulos críticos (tarificación, Party, Producto) se implementen de manera coherente y testeable.  
- Minimiza riesgo de retrabajo y pérdida de trazabilidad.  
- Permite introducir infraestructura solo cuando existen casos de uso validados.  
- Facilita el desarrollo guiado por pruebas en dominio y capa de aplicación.

---

## 4. Consecuencias

### Positivas
- Dominio económico protegido y estable.  
- Implementación incremental y testeable.  
- Mayor facilidad de mantenimiento y evolución futura.

### Negativas
- Mayor esfuerzo inicial en análisis y diseño.  
- Ritmo inicial de UI más lento.  
- Requiere disciplina estricta en el desarrollo.

---

## 5. Alcance

Aplica exclusivamente al **desarrollo del dominio MVP**, incluyendo:

- Party / Organización  
- Producto / Variante / Categoría  
- Tarificación  
- Flujo de ventas y MES básico  

No aplica a Post-MVP, ni define infraestructura futura avanzada.

---

## 6. Integración con otros ADRs

- ADR-001: Stack tecnológico  
- ADR-002: Clean Architecture  
- ADR-003: Monolito Modular Local-First  
- ADR-004: Ciclo de vida hasta MVP  
- ADR-005: Definición de estrategia de desarrollo dirigida por dominio

---

## 7. Notas Adicionales / Consideraciones Especiales

- Desarrollo guiado por pruebas se aplica en dominio y capa de aplicación.  
- Persistencia inicial permitida solo para Party y Producto.  
- Frontend se desarrolla paralelo para validar casos de uso.  

---

## 8. Referencias

- Documentos internos: Documento Consolidado TramaTex v2.1  
- Diagramas de flujo de dominio y Bounded Contexts  
- Buenas prácticas DDD y Clean Architecture
