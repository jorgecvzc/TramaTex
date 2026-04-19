# 🏛️ ADR-001: Selección del Stack y Estrategia Tecnológica Base

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.1 |
| **Estado** | ✅ Aceptado |
| **Fecha** | 06-01-2026 |
| **Autores** | Jorge Cortés Villalba, Gemini |

---

## 🎯 Contexto
TramaTex es un sistema ERP/MES dirigido a microempresas del sector textil y EPIs (5–15 empleados), con las siguientes características clave:
*   Operativa mayoritariamente bajo pedido.
*   Flujo productivo con estados críticos (Diseño, Marcaje, Taller).
*   Alta dependencia de una tarificación precisa.
*   Infraestructura **Local-First** con hardware limitado.

El objetivo es optimizar la **eficiencia, mantenibilidad y longevidad**, evitando dependencias de la nube o stacks sobredimensionados.

---

## 🔍 Alternativas Consideradas

### Backend (API)
*   Python (FastAPI)
*   .NET 8
*   **Go (Golang)**

### Generación de Documentos
*   Renderizado de PDFs en el Backend.
*   **Delegación al Cliente (Web-to-Print)**.

---

## ✅ Decisión Adoptada

### 1. Backend: Go (Golang)
Se selecciona Go como lenguaje principal por su compilación a binario único, bajo consumo de memoria, arranque instantáneo y concurrencia eficiente. Es ideal para entornos con recursos limitados y garantiza una larga vida útil al software.
*   **Framework Web:** Gin Gonic (minimalista y de alto rendimiento).

### 2. Persistencia: PostgreSQL 15+
Motor robusto y probado, con excelente soporte para integridad referencial y cargas transaccionales. Se utiliza **GORM** exclusivamente en la capa de infraestructura para operaciones CRUD simples.

### 3. Frontend: Vue.js 3
Se adopta Vue.js 3 con la Composition API y Pinia para la gestión de estado.
*   **Nota de Evolución:** Aunque inicialmente se consideró Tailwind CSS, el proyecto ha evolucionado hacia un **Sistema de Diseño propio en Vanilla CSS** para maximizar la flexibilidad y el control sobre la UI industrial.

### 4. Generación de Documentos (Web-to-Print)
Se delega la generación de PDFs e impresión al frontend. Esto reduce la carga de CPU en el servidor y permite utilizar estándares web (HTML/CSS) para el diseño de documentos mercantiles.

### 5. Despliegue: Docker & Docker Compose
Para garantizar instalaciones reproducibles y facilitar el mantenimiento en infraestructuras locales sin dependencia de orquestadores complejos.

---

## 📈 Consecuencias

### Positivas
*   Sistema altamente eficiente en hardware modesto.
*   Dominio protegido frente a decisiones tecnológicas (Clean Architecture).
*   Reducción drástica de costes operativos.

### Negativas
*   Curva de aprendizaje inicial ligeramente superior.
*   Necesidad de mayor disciplina arquitectónica.

---
[Volver al Índice de ADRs](./README.md)
