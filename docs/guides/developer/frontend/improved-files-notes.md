# Notas sobre archivos .IMPROVED.vue

**Fecha:** 2026-01-17  
**Ubicación:** `apps/apps/frontend/src/`

---

## 📋 Archivos Identificados

```
apps/apps/frontend/src/
├── components/auth/
│   ├── LoginForm.vue              (VERSIÓN ACTIVA)
│   └── LoginForm.IMPROVED.vue     ⚠️ (VERSIÓN ALTERNATIVA)
└── pages/auth/
    ├── LoginPage.vue              (VERSIÓN ACTIVA)
    └── LoginPage.IMPROVED.vue     ⚠️ (VERSIÓN ALTERNATIVA)
```

---

## ⚠️ Problema

- Los archivos `.IMPROVED.vue` contienen ~99% de contenido idéntico a las versiones activas.
- No está claro cuál es el propósito (backup, experimental, mejorada pero no adoptada).
- Riesgo de desincronización si uno se actualiza pero el otro no.

---

## 🔍 Diferencias Detectadas

El análisis preliminar muestra que los archivos son esencialmente idénticos con posibles cambios menores en:
- Nombres de variables
- Comentarios
- Estructura de CSS

---

## ⚡ Acciones Recomendadas

### OPCIÓN A: Eliminar (Recomendada)
```bash
rm apps/apps/frontend/src/components/auth/LoginForm.IMPROVED.vue
rm apps/apps/frontend/src/pages/auth/LoginPage.IMPROVED.vue
```

**Razón:** No hay evidencia de que sean mejoras activamente mantenidas.

### OPCIÓN B: Conservar con Propósito Claro
Si hay mejoras reales que deban ser evaluadas:
1. Renombrar a `LoginForm.experimental.vue`
2. Crear una rama `feature/improved-login`
3. Documentar las diferencias en la descripción del PR
4. Revisar antes de hacer merge

### OPCIÓN C: Consolidar
Si `.IMPROVED.vue` es realmente mejor:
1. Hacer merge de los cambios en `LoginForm.vue`
2. Eliminar `.IMPROVED.vue`
3. Hacer commit de los cambios

---

## ✅ Decisión Pendiente

**Requiere input del equipo:**
- ¿Son estas versiones experimentales para ser evaluadas?
- ¿Son intentos fallidos que deben ser descartados?
- ¿Son backups que pueden eliminarse?

**Recomendación:** Eliminar hasta que se justifique su mantenimiento.

---

**Estado:** PENDIENTE DE DECISIÓN  
**Propietario:** Equipo de frontend