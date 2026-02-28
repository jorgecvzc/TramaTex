# ADR-012 – Arquitectura del Módulo Party

**Fecha:** 2026-02-01  
**Estado:** Aceptado  
**Autores:** Gemini CLI, Usuario  
**Última actualización:** 2026-02-08  

---

## 1. Contexto

El módulo `Party` debe gestionar clientes y proveedores, que pueden ser organizaciones o personas individuales. Durante la revisión, han surgido requisitos de relaciones complejas:
1.  Una `Party` puede ser una persona individual.
2.  Las organizaciones pueden tener relaciones jerárquicas (matrices/filiales).
3.  Una persona (empleado de una organización) puede ser también un cliente a título personal.
4.  Los puntos de contacto de una organización pueden ser personas o departamentos, y se busca una solución pragmática.

El modelo inicial era insuficiente para estos requisitos. Este ADR define un modelo de dominio robusto y flexible para el módulo `Party`.

---

## 2. Alternativas Consideradas

**Alternativa 1 – Mantener la Estructura de Dominio Actual**
- Ventajas: Coherencia con la implementación preexistente.
- Desventajas: Incapaz de representar a personas individuales como clientes/proveedores. Obliga a modelar a personas como organizaciones.

**Alternativa 2 – Modelo de Party Abstracto con Perfil Único**
- Ventajas: Más simple que un modelo de roles múltiples. Satisface el requisito de "persona individual".
- Desventajas: Para el caso del "empleado que es cliente", se necesitarían dos `Party`s para la misma persona, lo que es propenso a inconsistencias.

**Alternativa 3 – Modelo de Party con Roles y Relaciones (Decisión Final)**
- Ventajas: Única opción que resuelve la complejidad de relaciones entre `Party`s de forma nativa.
- Desventajas: Mayor complejidad inicial.

---

## 3. Criterios de Decisión

- **Simplicidad para MVP:** Priorizar soluciones que resuelvan la necesidad principal sin sobre-ingeniería.
- **Flexibilidad:** Capacidad de manejar diversos tipos de contactos y relaciones jerárquicas.
- **Adecuación al negocio:** Ajustarse a las necesidades de microempresas (gestión no excesivamente granular).
- **Consistencia DDD:** Mantener una única entidad `Party` como agregación central.

---

## 4. Decisión Adoptada

Se adopta el **Modelo de Party con Roles y Relaciones** con un **Manejo de Contactos Simplificado**.

**Justificación:**
El modelo de Roles y Relaciones resuelve la complejidad central de las relaciones entre `Party`s. Se elige la sub-alternativa de contactos simplificados (Value Object `ContactDetails` en `OrganizationProfile`) para priorizar la simplicidad y el menor coste de implementación para el MVP, permitiendo flexibilidad suficiente sin una gestión interna de contactos excesivamente granular.

---

## 5. Consecuencias

### Positivas
- Se elimina el riesgo de inconsistencias en relaciones empleado-cliente.
- Soporte nativo para jerarquías empresariales (matrices/filiales).
- Modelo extensible que puede evolucionar a una gestión de contactos más granular.

### Negativas
- Refactorización completa del módulo `Party`.
- Rediseño de DTOs y endpoints para reflejar el nuevo modelo centrado en `Party`.
- Consultas más complejas para tipos de contacto específicos.

---

## 6. Alcance

Aplica al diseño y desarrollo del Bounded Context de `Party`. Incluye el modelado de persistencia (Profiles, Roles, Relationships) y la integración de este módulo con `Product`, `Pricing` y `Sales`.

---

## 7. Integración con otros ADRs

- **ADR-005: Gestión Unificada de Clientes y Proveedores:** Evolución del patrón Party base definido originalmente.
- [ADR-011: Estrategia de Testing y Cobertura](adr-011-testing-coverage-strategy.md): El módulo Party debe cumplir el umbral de MVP (>= 75%).

---

## 8. Notas Adicionales / Consideraciones Especiales

### Estado de implementación (2026-02-08)
- Repositorios implementados con GORM en la capa `persistence`.
- Errores tipados en dominio y aplicación.
- Cobertura de pruebas actual en cumplimiento con el umbral de MVP.

### Plan de Implementación
1. Diseño detallado del nuevo esquema de base de datos.
2. Refactorización del Dominio.
3. Actualización de la Persistencia, Aplicación y API.

### Votos y Aprobación
- Usuario: ✅ Aprobado
- Gemini: ✅ Aprobado

---

## 9. Referencias

- [docs/modules/party/implementation-guide.md](../../modules/party/implementation-guide.md)
- [docs/archive/log/party-improvements.md](../../archive/log/party-improvements.md)