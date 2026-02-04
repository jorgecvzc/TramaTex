# ADR-012: Arquitectura del Módulo Party

- **Estado:** Finalizado
- **Fecha:** 2026-02-01
- **Proponentes:** Gemini, Usuario

## Contexto y Planteamiento del Problema

El módulo `Party` debe gestionar clientes y proveedores, que pueden ser organizaciones o personas individuales. Durante la revisión, han surgido requisitos de relaciones complejas:
1.  Una `Party` puede ser una persona individual.
2.  Las organizaciones pueden tener relaciones jerárquicas (matrices/filiales).
3.  Una persona (empleado de una organización) puede ser también un cliente a título personal.
4.  Los puntos de contacto de una organización pueden ser personas o departamentos, y se busca una solución pragmática.

El modelo actual es insuficiente para estos requisitos. Este ADR busca definir un modelo de dominio robusto y flexible para el módulo `Party`.

## Alternativas Consideradas

### Alternativa 1: Mantener la Estructura de Dominio Actual (Rechazada)

- **Descripción:** La implementación inicial del módulo `Party` modela a una `Party` como una `Organization` únicamente. Las `Person`s están siempre asociadas a una `Organization`. No existe el concepto de una `Person` individual que actúe como `Party`.
- **Ventajas:**
  - Coherencia con la implementación preexistente, lo que minimiza el esfuerzo de refactorización inmediato.
- **Desventajas:**
  - **No satisface requisitos clave:** Incapaz de representar a personas individuales como clientes/proveedores.
  - **Conceptualmente incorrecto:** Obliga a modelar a personas como organizaciones si actúan como clientes individuales.
  - **Inflexibilidad:** No puede manejar relaciones complejas (ej. empleado que también es cliente, jerarquías de empresas).

### Alternativa 2: Modelo de Party Abstracto con Perfil Único (Rechazada)

- **Descripción:** Introduce una entidad `Party` abstracta que es *o* una `Person` *o* una `Organization`, pero no ambas a la vez. Las relaciones (empleado-de, filial-de) se gestionarían con enlaces explícitos entre diferentes `Party`s.
- **Ventajas:**
  - Más simple que un modelo de roles múltiples y más flexible que la Alternativa 1.
  - Satisface el requisito de "persona individual como cliente".
- **Desventajas:**
  - **No satisface los requisitos complejos:** Para el caso del "empleado que es cliente", se necesitarían dos `Party`s para la misma persona (uno como cliente, y otro como `Person` dentro de la `Organization`), y un enlace entre ellos. Esto es complejo y propenso a inconsistencias.

### Alternativa 3: Modelo de Party con Roles y Relaciones (Decisión Final)

- **Descripción General:** Se establece una única entidad `Party` como agregación central. Esta `Party` puede tener asociados múltiples **roles** (ej. Cliente, Empleado), **perfiles** de datos (`PersonProfile`, `OrganizationProfile`), y **puntos de contacto**. Las relaciones entre `Party`s se gestionan a través de `PartyRelationship`.

#### Sub-alternativa 3.1: Modelo de Contactos Detallado (Rechazada)

- **Descripción:** Implementa `ContactPoint` como una interfaz con implementaciones específicas (`PersonContact`, `DepartmentContact`).
- **Ventajas:** Máxima flexibilidad, corrección conceptual.
- **Desventajas:** Máxima complejidad, sobre-ingeniería para un MVP dirigido a microempresas.

#### Sub-alternativa 3.2: Modelo de Contactos Simplificado (Decisión Final)

- **Descripción:** En lugar de interfaces `ContactPoint` complejas, una `Organization` (a través de su `PartyProfile`) tendrá una lista de `ContactDetails`. `ContactDetails` es un `Value Object` simple.

  ```mermaid
  classDiagram
      direction LR
      class Party {
          <<Aggregate Root>>
          +PartyID id
      }
      class PartyProfile {
          <<Interface>>
      }
      class PersonProfile {
          +string firstName
          +string lastName
      }
      class OrganizationProfile {
          +string name
          +string taxId
          +ContactDetails[] contacts
      }
      class PartyRelationship {
          +PartyID fromParty
          +PartyID toParty
          +RelationshipType type
      }
      class PartyRole {
        +RoleType type
      }
      class ContactDetails {
          <<Value Object>>
          +string typeDescription  // Ej: "Ventas", "Soporte", "Juan Pérez (Responsable)"
          +string phone
          +string email
          +PartyID? relatedPartyId // Opcional: enlace a otra Party (ej. una persona específica)
      }

      Party "1" -- "1..*" PartyProfile : has profiles
      Party "1" -- "1..*" PartyRole : has roles
      Party "1" -- "0..*" PartyRelationship : participates in
      
      PartyProfile <|.. PersonProfile
      PartyProfile <|.. OrganizationProfile
      OrganizationProfile o-- ContactDetails
  ```

- **Ventajas:**
  - **Menor complejidad inicial:** Más rápido de implementar y mantener para un MVP.
  - **Flexible:** Permite guardar diversos tipos de contactos con un esquema simple.
  - **Adecuado para microempresas:** Responde a la necesidad principal sin una gestión interna de contactos excesivamente granular.
- **Desventajas:**
  - Menos estructurado y con validación más débil que el modelo detallado.
  - Consultas más complejas para tipos de contacto específicos.

## Decisión Tomada

Se ha decidido adoptar la **Alternativa 3: Modelo de Party con Roles y Relaciones** como enfoque general para el módulo `Party`. Para el manejo de contactos, se ha elegido la **Sub-alternativa 3.2 (Modelo de Contactos Simplificado)**.

**Justificación:**
El modelo de Roles y Relaciones resuelve la complejidad central de las relaciones entre `Party`s, siendo la única opción viable para los requisitos expuestos. La elección del modelo de contactos simplificado (3.2) se basa en la priorización de la simplicidad y el menor coste de implementación para un MVP dirigido a microempresas. Este enfoque permite flexibilidad suficiente para los requisitos actuales y puede evolucionar a un modelo más detallado si la necesidad de una gestión de contactos más granular surge en el futuro.

## Comentarios y Requisitos Adicionales a la Decisión

### 1. Relaciones entre Empresas (Matrices y Filiales)
- **Solución:** Se gestionará a través de la entidad `PartyRelationship` con un tipo de relación como `IS_SUBSIDIARY_OF`.

### 2. Empleado que es Cliente
- **Solución:** Se crean dos `Party`s (uno para la organización, otro para la persona) y se enlazan con una `PartyRelationship` de tipo `IS_EMPLOYEE_OF`. La `Party` de la persona puede tener además un rol de `CLIENTE`.

## Consecuencias y Ramificaciones

- **Impacto en el Código:** Refactorización completa del módulo `Party`.
- **Impacto en la Base de Datos:** Diseño de un esquema de tablas que soporte `Party`, `Profiles` (Person y Organization), `Roles`, `Relationships` y `ContactDetails`.
- **Impacto en la API:** Rediseño de los DTOs y endpoints para reflejar el nuevo modelo centrado en `Party`.

## Plan de Implementación

1.  **Aprobación del ADR.** (Ya se ha discutido y aprobado con esta última iteración.)
2.  **Diseño detallado del nuevo esquema de base de datos.**
3.  **Refactorización del Dominio.**
4.  **Actualización de la Persistencia, Aplicación y API.**

## Votos y Aprobación

- Usuario: ✅ Aprobado
- Gemini: ✅ Aprobado