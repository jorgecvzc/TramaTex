# ADR-004 – Ciclo de Vida de Desarrollo e Implementación hasta MVP

**Fecha:** 07/01/2026  
**Estado:** Aceptado  
**Autores:** Jorge Cortés Villalba, ChatGPT  
**LLM utilizado:** Claude (Anthropic)  

---

## 1. Contexto

TramaTex es un ERP/MES para microempresas del sector textil y EPIs que debe alcanzar un **MVP plenamente operativo en entorno real**, bajo las siguientes condiciones:

- Equipo de desarrollo reducido  
- Infraestructura local limitada (servidor i3)  
- Dominio crítico en tarificación y precios finales  
- Necesidad de calidad técnica elevada:
  - Clean Architecture  
  - TDD  
  - Alta cobertura de tests  
- El MVP **no es un prototipo**, sino un sistema productivo  

Se requiere un **ciclo de vida de implementación acotado exclusivamente al MVP**, minimizando riesgos y garantizando estabilidad.

---

## 2. Alternativas Consideradas

1. **Ciclo de vida completo (MVP + Post-MVP):**  
   - Rechazada: diluye foco, aumenta riesgo de no entregar MVP sólido  

2. **Prototipado rápido sin rigor técnico:**  
   - Rechazada: alto riesgo en tarificación y finanzas, deuda técnica inasumible  

3. **Ciclo incremental orientado a MVP (decisión adoptada):**  
   - Incrementos funcionales completos  
   - Priorización estricta del núcleo de negocio  
   - Validación temprana en entorno real  
   - Rigor técnico desde el inicio  

---

## 3. Criterios de Decisión

- MVP funcional desde primeras iteraciones  
- Dominio crítico protegido y testeable  
- Incrementos entregables y validables  
- Minimizar riesgo operativo y técnico  
- Preparación para iteraciones futuras controladas  

---

## 4. Decisiones

Se adopta un **Ciclo de Vida Incremental hasta MVP** con fases definidas:

### Fase 0 – Fundaciones Técnicas

- Preparar infraestructura mínima y Clean Architecture base  
- Configuración Docker, repositorios, pipeline de tests, linters  
- Esqueleto de autenticación y autorización  
- Sin funcionalidad de negocio  

### Fase 1 – Implementación del MVP

- Desarrollar módulos críticos:
  - Ventas y atención al cliente  
  - Tarificación completa  
  - Gestión básica de proveedores  
  - Flujo de producción esencial (MES)  
- Interfaces gráficas funcionales (taller, ventas, finanzas básicas)  

### Estrategia de Iteración

- Iteraciones cortas y controladas  
- Cada incremento debe ser funcional y testeable  
- Enfoque ágil pragmático, no Scrum completo  

### Criterio de Finalización

- MVP desplegado en servidor local  
- Sistema soporta flujo real de trabajo  
- Tarificación y ventas fiables  
- Sistema mantenible por terceros  

---

## 5. Consecuencias

### Positivas

- Foco en valor de negocio inmediato  
- Reducción de riesgo de sobreingeniería  
- MVP estable y productivo  
- Control sobre alcance y plazos  
- Facilita iniciar un nuevo proyecto post-MVP  

### Negativas

- Mayor disciplina requerida en planificación  
- Restricciones en funcionalidad avanzada durante MVP  

---

## 6. Alcance

- Ciclo de vida y fases de implementación hasta MVP  
- Entregables funcionales parciales y finales del MVP  
- Validación de dominio y casos de uso críticos  

---

## 7. Referencias

- ADR-001: Selección del Stack Tecnológico  
- ADR-002: Adopción de Clean Architecture y DDD con Rigor Asimétrico  
- ADR-003: Tipo y Distribución de la Aplicación (Monolito Modular)  

---

## 8. Notas Finales

Este ADR define la **estrategia de ejecución controlada del MVP**, asegurando:

- Dominio crítico protegido  
- Incrementos funcionales testeables  
- Preparación para iteraciones futuras sin comprometer estabilidad
