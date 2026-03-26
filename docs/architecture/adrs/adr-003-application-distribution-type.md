# ADR-003 – Tipo y Distribución de la Aplicación (Monolito Modular Local-First)

**Fecha:** 07/01/2026  
**Estado:** Aceptado  
**Autores:** Jorge Cortés Villalba, ChatGPT  
**LLM utilizado:** Claude (Anthropic)  

---

## 1. Contexto

TramaTex es un ERP/MES para microempresas (5–15 empleados) con restricciones y objetivos:

- Infraestructura **local-first**, sin dependencia obligatoria de cloud  
- Hardware limitado y heterogéneo (servidor i3, clientes modestos, tablets)  
- Necesidad de **mantenibilidad por terceros externos**  
- Dominio crítico con lógica de tarificación y pedidos personalizados  
- Ciclo de vida largo del software  
- Proyección futura **no inmediata**: escalado funcional, integraciones técnicas, posible SaaS  

Alternativas arquitectónicas evaluadas:

- **Microservicios:** excesiva complejidad, sobrecarga operativa, requiere madurez organizativa no presente  
- **Monolito tradicional:** riesgo de acoplamiento, degradación de dominio, difícil extracción futura de servicios  

Se requiere una solución que **maximice simplicidad operativa hoy** sin bloquear evolución futura.

---

## 2. Alternativas Consideradas

1. **Microservicios desde inicio:**  
   - Ventajas: escalabilidad futura, modularidad física  
   - Desventajas: complejidad inicial, sobrecoste operativo, dificultad de mantenimiento con equipo reducido  

2. **Monolito clásico:**  
   - Ventajas: simplicidad inicial  
   - Desventajas: acoplamiento del dominio, riesgo de "monolito de barro", difícil evolución  

3. **Monolito modular (decisión adoptada):**  
   - Ventajas: simplicidad operativa, modularidad lógica, dominio protegido, escalable hacia microservicios en el futuro  
   - Desventajas: menor aislamiento en tiempo de ejecución frente a microservicios, requiere disciplina arquitectónica  

---

## 3. Criterios de Decisión

- Protección del dominio frente a degradación  
- Simplicidad de despliegue y operación en infraestructura local  
- Modularidad lógica para evolución futura  
- Capacidad de mantenimiento por terceros  
- Preparación para refactor o extracción futura  

---

## 4. Decisión Adoptada

Se adopta un **Monolito Modular Local-First** con las siguientes características:

### tramatex-api

- **Un único tramatex-api** (binario Go)  
- **Una única base de datos PostgreSQL**  
- Dominios claramente separados:  
  - Dominio principal: Ventas, Tarificación, Clientes/Proveedores  
  - Subdominio especializado: MES – Producción personalizada  
  - Dominios Post-MVP: Compras, Inventario  
  - Módulos transversales: Seguridad, Auditoría, Gestión documental, i18n  

### Frontend

- SPA única (Vue.js 3 + Pinia + Tailwind CSS)  
- Cliente de la API, adaptador de presentación, motor Web-to-Print  
- No contiene lógica de negocio crítica  

### Persistencia

- Base de datos única, separación lógica por esquemas o convenciones  
- Integridad principalmente desde dominio  
- Evita joins transversales no justificados  

### Principios del Monolito Modular

1. Un proceso, múltiples dominios  
2. Límites primero: cada módulo define su modelo, reglas y casos de uso  
3. Comunicación explícita mediante interfaces/contratos  
4. Extraíble por diseño: cada módulo puede convertirse en servicio independiente  
5. Infraestructura compartida, dominio no  

### Evolución futura

- Extracción progresiva de motor de tarificación, MES o inventario  
- Introducción de comunicación asíncrona (eventos de dominio)  
- Transición eventual a modelo SaaS con multitenencia  

---

## 5. Consecuencias

### Positivas

- Simplicidad de despliegue y operación  
- Coste operativo mínimo  
- Dominio protegido frente a degradación  
- Base sólida para futuras extracciones de módulos  
- Compatible con hardware limitado y entornos locales  

### Negativas

- Requiere disciplina estricta para evitar acoplamientos indebidos  
- Menor aislamiento en tiempo de ejecución frente a microservicios  
- Necesidad de gobernanza arquitectónica continua  

---

## 6. Alcance

- Tipo de aplicación (monolito vs distribuido)  
- Distribución lógica de módulos  
- Estrategia de evolución arquitectónica  
- Forma de despliegue del MVP  

Cualquier cambio hacia microservicios, arquitectura distribuida o separación física de dominios requiere **nuevo ADR**.

---

## 7. Integración con otros ADRs

- ADR-001: Selección del Stack Tecnológico  
- ADR-002: Adopción de Clean Architecture y DDD con Rigor Asimétrico  

---

## 8. Notas Adicionales / Consideraciones Especiales

Este ADR, junto con ADR-001 y ADR-002, define el **marco arquitectónico completo** de TramaTex para el MVP.  
El desarrollo debe enfocarse en **modelar correctamente el dominio**, respetando los límites del monolito modular.

---

## 9. Referencias
