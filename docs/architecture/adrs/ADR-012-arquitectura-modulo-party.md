# ADR-012: Arquitectura del Módulo Party

- **Estado:** En Propuesta
- **Fecha:** 2026-02-01
- **Proponentes:** Gemini, Usuario

## Contexto y Planteamiento del Problema

El módulo `Party` es fundamental para la gestión unificada de clientes y proveedores. Durante la revisión del dominio, ha surgido un requisito clave: **un cliente puede ser una persona individual, no necesariamente parte de una organización.**

La implementación actual asume que toda `Party` es una `Organization`, y que una `Person` siempre pertenece a una `Organization`. Este modelo es insuficiente para representar a clientes individuales.

Este ADR tiene como objetivo definir un modelo de dominio para el módulo `Party` que sea lo suficientemente flexible para representar tanto a organizaciones como a personas individuales como clientes o proveedores, y documentar formalmente la decisión para guiar la refactorización y el desarrollo futuro.

## Alternativas Consideradas

### Alternativa 1: Mantener la Estructura de Dominio Actual (Rechazada)

- **Descripción:** Continuar con el modelo donde `Organization` es el único tipo de `Party`. Para representar a una persona individual, se crearía una `Organization` con el mismo nombre que la persona.
- **Ventajas:**
  - Requiere el menor cambio inmediato en el código.
- **Desventajas:**
  - **Conceptualmente incorrecto:** Modela a una persona como una organización, lo que lleva a confusiones.
  - **Problemas de datos:** Campos como `TaxID` o `Website` no tendrían sentido para una persona.
  - **Complejidad futura:** Dificulta la gestión de atributos específicos de personas (ej. fecha de nacimiento).

### Alternativa 2: Modelo de Party Abstracto con Perfil Único (Recomendada)

- **Descripción:** Introducir una nueva entidad `Party` que actúe como raíz del agregado. Esta entidad `Party` tendría un ID y un tipo (ej. `PERSON` u `ORGANIZATION`), y estaría asociada a **un único perfil** (`Person` u `Organization`). En este modelo, una `Party` es *o* una Persona *o* una Organización, pero no ambas a la vez.

  ```mermaid
  classDiagram
      class Party {
          <<Aggregate Root>>
          +PartyID id
          +PartyType type
          +PartyProfile profile
      }
      class PartyProfile {
          <<Interface>>
      }
      class Person {
          <<Entity>>
          +string firstName
          +string lastName
      }
      class Organization {
          <<Entity>>
          +string name
          +string taxId
      }
      Party "1" -- "1" PartyProfile : has one
      PartyProfile <|.. Person
      PartyProfile <|.. Organization
  ```

- **Ventajas:**
  - **Modelo conceptualmente correcto:** Representa fielmente la realidad del negocio.
  - **Flexibilidad:** Permite añadir nuevos tipos de `Party` en el futuro.
  - **Datos limpios:** Cada tipo de `Party` tiene solo los atributos que le corresponden.
- **Desventajas:**
  - **Requiere una refactorización significativa.**
  - **Menos flexible si una Party puede tener múltiples roles/perfiles.**

### Alternativa 3: Modelo de Party con Múltiples Perfiles

- **Descripción:** Similar a la Alternativa 2, pero una `Party` puede tener **múltiples perfiles** a la vez. Por ejemplo, una `Party` podría ser simultáneamente una `Person` (el individuo) y una `Organization` (su empresa unipersonal).

- **Ventajas:**
  - **Máxima flexibilidad:** Permite modelar relaciones complejas.
- **Desventajas:**
  - **Mayor complejidad:** Tanto en el modelo de dominio como en la persistencia y la lógica de negocio.
  - **Potencial sobre-ingeniería:** Podría ser más complejo de lo necesario para los requisitos actuales.

## Decisión Tomada

Se propone adoptar la **Alternativa 2: Modelo de Party Abstracto con Perfil Único**.

**Justificación:**
Aunque la Alternativa 3 ofrece más flexibilidad, la Alternativa 2 satisface el requisito principal ("un cliente puede ser una persona") de una manera más simple y directa. La complejidad de gestionar múltiples perfiles por `Party` no parece necesaria para el MVP y puede añadirse en el futuro si el negocio lo requiere.

## Consecuencias y Ramificaciones

- **Impacto en el Código:**
  - **Refactorización Mayor:** Se debe crear la nueva entidad `Party`. `Organization` y `Person` se convertirán en perfiles.
- **Impacto en la Base de Datos:**
  - Se necesitará una nueva tabla `parties`. Se usará una estrategia de "tabla por clase" o "tabla única" para los perfiles.
- **Impacto en la API:**
  - Los endpoints (`/organizations`) se renombrarán a `/parties`.
  - Los DTOs de respuesta deberán ser capaces de representar tanto a personas como a organizaciones.

## Plan de Implementación

1.  **Aprobación del ADR:** Obtener la aprobación del usuario para esta decisión.
2.  **Refactorización del Dominio:** Implementar la nueva estructura.
3.  **Actualización de la Persistencia:** Modificar repositorios y migraciones.
4.  **Actualización de la Aplicación y la API:** Ajustar casos de uso, DTOs y handlers.

## Votos y Aprobación

- Usuario: ⬜ Pendiente de votación
- Gemini: ✅ Aprobado
