# 🏛️ ADR-007: Orden de Implementación de Módulos

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.0 |
| **Estado** | ✅ Aceptado |
| **Fecha** | 01-02-2026 |
| **Autores** | Equipo de Arquitectura de TramaTex |

---

## 🎯 Contexto
Al ser un monolito modular con principios de Clean Architecture y DDD, los dominios de negocio tienen dependencias explícitas. Es vital establecer un orden de construcción que evite dependencias circulares y asegure que los módulos base estén listos antes que los dependientes.

---

## 🔍 Alternativas Consideradas
1. **Desarrollo Ad-hoc:** Implementar según necesidad inmediata. Riesgo alto de deuda técnica y bloqueos constantes por prerrequisitos no listos.
2. **Desarrollo por Capas de Dependencia (Decisión Adoptada):** Analizar el grafo de dependencias y construir desde lo fundamental hacia lo específico.

---

## ✅ Decisión Adoptada
Se adopta la implementación por **Capas de Dependencia** en el siguiente orden:

### 1. Fase 0: Fundaciones (Seguridad)
*   **Módulo `IAM`:** Base de autenticación y roles. Requerido por todos los demás módulos para autoría y permisos.

### 2. Fase 1: Dominio Base (Maestros)
*   **Módulo `Party`:** Gestión de clientes y proveedores.
*   **Módulo `Product`:** Catálogo de productos y variantes.

### 3. Fase 2: Lógica Crítica (Económico)
*   **Módulo `Pricing`:** Motor de precios. Depende de `Product` (costes) y `Party` (tarifas).
*   **Módulo `Sales`:** Orquestación comercial. Utiliza todos los módulos anteriores.

### 4. Fase 3: Operaciones (Taller)
*   **Módulo `MES`:** Gestión de producción. Depende de `Sales` (demandas) y `Product`.

---

## 📈 Consecuencias
### Positivas
*   Eliminación de bloqueos por dependencias no resueltas.
*   Proceso de desarrollo predecible y ordenado.
*   Validación incremental estable del dominio de negocio.

### Negativas
*   Menor flexibilidad para cambiar prioridades de negocio que rompan el orden técnico de dependencias.

---
[Volver al Índice de ADRs](./README.md)
