# ADR-001 – Selección del Stack Tecnológico y Estrategia Tecnológica Base

**Fecha:** 06/01/2026  
**Estado:** Aceptado  
**Autores:** Jorge Cortés Villalba, Gemini  
**LLM utilizado:** Claude (Anthropic)  

---

## 1. Contexto

TramaTex es un sistema ERP/MES dirigido a microempresas del sector textil y EPIs (5–15 empleados), con las siguientes características clave:

- Operativa **mayoritariamente bajo pedido**
- Flujo productivo con estados críticos (Diseño, Marcaje, Taller)
- Alta dependencia de una **tarificación correcta**
- Infraestructura **local-first** con hardware limitado (servidor i3, 16GB RAM)
- Necesidad de mantenimiento por terceros externos
- Proyección futura hacia:
  - Escalado funcional
  - Integración con maquinaria
  - Posible modelo SaaS

El objetivo tecnológico no es maximizar sofisticación, sino **optimizar eficiencia, mantenibilidad y longevidad**, evitando dependencia temprana de cloud o stacks sobredimensionados.

---

## 2. Alternativas Consideradas

### Backend
- **Python (FastAPI)**  
- **.NET 8**  
- **Go (Golang)**  

### Arquitectura
- Arquitectura Hexagonal  
- Clean Architecture  
> Nota: La decisión arquitectónica se formaliza en ADR-002.

### Generación de documentos
- Renderizado de PDFs en backend  
- Delegación de renderizado al cliente (Web-to-Print)  

---

## 3. Criterios de Decisión

Las decisiones se tomaron considerando:

1. Consumo de recursos en hardware limitado  
2. Facilidad de mantenimiento por terceros  
3. Capacidad de testeo (TDD, alto coverage en dominio)  
4. Claridad y protección del modelo de dominio  
5. Escalabilidad futura (local → SaaS)  
6. Carga operativa del servidor  
7. Simplicidad de despliegue y operación  
8. Longevidad tecnológica (evitar obsolescencia prematura)  

---

## 4. Decisiones

### Lenguaje y Backend

Se selecciona **Go (Golang)** como lenguaje principal del backend.

**Justificación:**

- Compilación a binario único, sin runtime externo  
- Consumo de memoria bajo y predecible  
- Arranque instantáneo  
- Concurrencia eficiente mediante goroutines  
- Código explícito y legible, adecuado para mantenimiento externo  
- Ecosistema maduro y estable a largo plazo  
- Uso probado en infraestructuras críticas y sistemas de larga vida útil  

---

### Framework Web

Se selecciona **Gin Gonic** como framework HTTP.

**Justificación:**

- Minimalista y de alto rendimiento  
- Bajo overhead en comparación con frameworks más complejos  
- Fácil integración con Clean Architecture  
- Amplia adopción y documentación  
- Adecuado para APIs REST sencillas y estables  

> Gin se considera **infraestructura**, no parte del dominio.

---

### Persistencia y ORM

- **Base de datos:** PostgreSQL 15+  
- **ORM:** GORM (uso controlado)  

**Justificación PostgreSQL:**

- Motor robusto, ACID y probado  
- Excelente soporte para integridad referencial  
- Adecuado para cargas transaccionales  
- Compatible con despliegues locales y futuros entornos SaaS  

**Justificación GORM:**

- Reduce fricción en operaciones CRUD simples  
- Amplia adopción en el ecosistema Go  
- Se utiliza **exclusivamente en la capa de infraestructura**  
- No se permite dependencia del ORM en el dominio  

> En lógica crítica (tarificación, reglas económicas), se prioriza el control explícito sobre la abstracción.

---

### Frontend

- **Framework:** Vue.js 3  
- **Patrón:** Composition API  
- **Estado:** Pinia  
- **Estilos:** Tailwind CSS  

**Justificación:**

- Curva de aprendizaje moderada  
- Ecosistema maduro  
- Separación clara de responsabilidades  
- Buen soporte para aplicaciones administrativas y terminales de taller  
- Facilidad para delegar tareas de presentación complejas (documentos, impresión)  

---

### Generación de Documentos (Web-to-Print)

Se decide **delegar la generación de documentos (PDF / impresión)** al frontend.

**Justificación:**

- Reduce carga CPU y memoria en el servidor  
- Evita dependencias pesadas de renderizado en backend  
- Garantiza compatibilidad con hardware limitado  
- Permite usar estándares web (HTML + CSS Media Queries)  
- El backend se limita a exponer datos estructurados y contratos claros  

El frontend actúa como **adaptador de infraestructura de presentación**.

---

### Contenerización y Despliegue

- **Docker / Docker Compose** para despliegues controlados  
- Sin dependencia de orquestadores complejos (Kubernetes fuera de alcance MVP)  

**Justificación:**

- Simplifica instalación y mantenimiento  
- Facilita backups y restauraciones  
- Reduce dependencia del entorno host  
- Compatible con infraestructuras locales  

---

## 5. Consecuencias

### Positivas

- Sistema eficiente y estable en hardware limitado  
- Backend predecible y fácil de operar  
- Dominio protegido frente a decisiones tecnológicas  
- Reducción de costes operativos  
- Base sólida para evolución futura  

### Negativas

- Menor velocidad inicial frente a stacks más "rápidos de prototipar"  
- Necesidad de mayor disciplina arquitectónica  
- Menor disponibilidad de perfiles Go frente a stacks mainstream  

Estas consecuencias se aceptan explícitamente como parte de la estrategia del proyecto.

---

## 6. Alcance

Este ADR aplica a:

- Selección del stack tecnológico base  
- Estrategia de despliegue MVP  
- Relación backend ↔ frontend  
- Política de generación de documentos  
- Criterios de eficiencia y mantenibilidad  

Cualquier cambio relevante en estos puntos requerirá un **nuevo ADR**.

---

## 7. Referencias

- ADR-002: Adopción de Clean Architecture y DDD con Rigor Asimétrico  
- ADR-003: Tipo y Distribución de la Aplicación (Monolito Modular)  
- Documento de Visión del Proyecto  

---

## 8. Notas Finales

Este ADR define el **suelo tecnológico** de TramaTex.  
Las decisiones arquitectónicas estructurales se complementan con:

- ADR-002 (Clean Architecture + DDD)  
- ADR-003 (Monolito modular con proyección a microservicios)  

A partir de este punto, el foco del proyecto se desplaza del *qué tecnología usar* al *cómo modelar correctamente el dominio*.
