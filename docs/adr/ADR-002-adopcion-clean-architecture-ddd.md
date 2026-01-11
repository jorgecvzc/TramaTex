# ADR-002 – Adopción de Clean Architecture y DDD con Rigor Asimétrico

**Fecha:** 06/01/2026  
**Estado:** Aceptado  
**Autores:** Jorge Cortés Villalba, ChatGPT  
**LLM utilizado:** Claude (Anthropic)  

---

## 1. Contexto

TramaTex es un ERP/MES para microempresas del sector textil y EPIs, caracterizado por:

- Gestión compleja de productos con variantes (talla, color)  
- Alto volumen de pedidos personalizados  
- Flujo productivo con fases diferenciadas (diseño, marcaje, taller)  
- Dependencia crítica de una **tarificación correcta y consistente**  
- Interacción constante entre áreas comerciales, técnicas y operativas  

Adicionalmente, el sistema debe:

- Operar en infraestructura **local-first** con recursos limitados  
- Ser mantenible por terceros externos  
- Tener una vida útil larga  
- Permitir evolución funcional progresiva  
- Facilitar transición futura a arquitecturas distribuidas  

Riesgo principal: **degradación del modelo de dominio** por acoplamientos indebidos o dispersión de lógica de negocio.

---

## 2. Alternativas Consideradas

### Aplicación de Arquitectura

1. **No aplicar Clean Architecture / DDD:**  
   - Ventajas: rapidez inicial, menos complejidad.  
   - Desventajas: riesgo de deuda técnica, modelo de dominio vulnerable, difícil escalabilidad y testeo.  

2. **Aplicar Clean Architecture + DDD con rigor uniforme:**  
   - Ventajas: dominio protegido, testable, separación de responsabilidades.  
   - Desventajas: esfuerzo inicial elevado, curva de aprendizaje.  

3. **Aplicar Clean Architecture + DDD con rigor asimétrico (decisión adoptada):**  
   - Ventajas: protege núcleo crítico, aplica esfuerzo arquitectónico donde aporta valor, permite iteración ágil en áreas menos críticas.  
   - Desventajas: requiere disciplina, puede inducir confusión si se relaja el rigor sin control.

---

## 3. Criterios de Decisión

- Protección del dominio crítico (tarificación)  
- Facilidad de testeo y refactor controlado  
- Mantenibilidad por terceros  
- Evolución progresiva del sistema  
- Control de deuda técnica  
- Aplicación de esfuerzo arquitectónico **alineado con valor real**  

---

## 4. Decisiones

Se adopta **Domain-Driven Design (DDD)** junto con **Clean Architecture** aplicando **rigor asimétrico**.

### Dominio

- Implementado con Clean Architecture estricta  
- Contiene: Entidades, Value Objects, Servicios de dominio  
- No depende de: Frameworks, Infraestructura, ORM, Serialización, Transporte  
- Completamente testeable en aislamiento  
- Considerado activo estratégico principal, especialmente motor de tarificación  

### Capa de Aplicación

- Orquesta casos de uso y flujos de negocio  
- Se mantiene separación conceptual respecto al dominio  
- Permite abstracción menos estricta para CRUD simples y casos sin valor real adicional  
- Relajaciones explícitas, localizadas y justificadas  

### Infraestructura

- Sustituible, aislada del dominio  
- Incluye: Framework web, Persistencia, Adaptadores externos, Despliegue, Integraciones técnicas  
- **No contiene lógica de negocio**  

---

## 5. Consecuencias

### Positivas

- Dominio estable, expresivo y protegido  
- Alta mantenibilidad a medio y largo plazo  
- Refactors controlados y localizados  
- Facilita evolución funcional  
- Base sólida para transición futura a microservicios  

### Negativas

- Mayor coste inicial de diseño  
- Incremento del boilerplate estructural  
- Curva de entrada más elevada para perfiles junior  

---

## 6. Alcance

- Modelado de dominio  
- Diseño de entidades y Value Objects  
- Definición de casos de uso  
- Organización del código fuente  
- Estrategia de testing  
- Evolución futura hacia microservicios  

No prescribe tecnologías concretas, solo **principios estructurales obligatorios**.  

---

## 7. Referencias

- ADR-001: Selección del Stack Tecnológico  
- ADR-003: Tipo y Distribución de la Aplicación (Monolito Modular)  

---

## 8. Notas Finales

- Este ADR consolida y sustituye decisiones previas relacionadas con:  
  - Uso de Clean Architecture  
  - Aplicación de DDD  
  - Nivel de rigor arquitectónico  

- Cualquier desviación relevante deberá justificarse mediante **nuevo ADR**.
