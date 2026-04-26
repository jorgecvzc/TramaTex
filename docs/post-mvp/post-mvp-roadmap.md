# Hoja de Ruta Post-MVP — TramaTex

> Enumeración y orden  de  las mejoras, funcionalidades y módulos planificados para después del MVP.
> Habrá que definir contexto, prioridad estimada y referencias a la documentación existente.

---

## Índice

1. [Unificación UI/UX y Sistema de Diseño](#1-unificación-uiux-y-sistema-de-diseño)
2. [Facturación Consolidada Multi-Albarán](#2-facturación-consolidada-multi-albarán)
3. [Adopción sistema de facturación electrónica](#3-facturacion-electronica)
4. [Extracción MES como Microservicio](#4-extracción-mes-como-microservicio)
5. [Caché de Productos/Variantes y Precios](#5-caché-de-productosvariantes-y-precios)
6. [Comunicación Asíncrona (Message Broker)](#6-comunicación-asíncrona-message-broker)
7. [Inteligencia de Negocio y Analítica](#7-inteligencia-de-negocio-y-analítica)
8. [Notificaciones en Tiempo Real](#8-notificaciones-en-tiempo-real)
9. [Búsqueda Global Avanzada](#9-búsqueda-global-avanzada)
10. [Nuevos Módulos de Negocio](#10-nuevos-módulos-de-negocio)
11. [Mejoras Técnicas de Infraestructura](#11-mejoras-técnicas-de-infrastructure)
12. [Cobertura de Tests](#12-cobertura-de-tests)
13. [Gestión Avanzada de Archivos de Diseño (MES)](#13-gestión-avanzada-de-archivos-de-diseño-mes)
14. [Integración Sales ↔ MES (Producto-Trabajo)](#14-integración-sales--mes-producto-trabajo)
15. [Asignación de Tareas MES a Operarios](#15-asignación-de-tareas-mes-a-operarios)

---

## 1. Unificación UI/UX y Sistema de Diseño

**Prioridad:** Máxima (Primera tarea Post-MVP)  
**Referencia Técnica Única:**
- [Plan Maestro de Unificación UI/UX](01-ui-ux-unification-master-plan.md)

**Contexto:**
Auditoría completa detectó inconsistencias significativas: patrones de navegación mixtos, ausencia de base CSS global, uso de emojis y paleta fragmentada. Esta fase transforma la UI en una herramienta industrial de alta eficiencia para operarios siguiendo el **Plan Maestro**.

**Tareas de Implementación:**

### 1.1 Sistema de Diseño Base y Estilos
- [ ] **Iconografía Profesional**: Instalar `lucide-vue-next` y sustituir todos los emojis y Material Symbols residuales (🗑️→Trash, 🔧→Package, ⚙️→Settings).
- [ ] **CSS Global**: Crear `_buttons.css` (estilos .btn uniformes) y `_dashboards.css` (estilos de KPIs y tarjetas comunes).
- [ ] **Identidad Visual**: Forzar uso estricto de variables en `_variables.css` (Amarillo Oro #E6B800, Azul Profundo #1B3A6B).
- [ ] **BasePageHeader**: Refactorizar para incluir indicadores de atajos `<kbd>` y migas de pan consistentes.

### 1.2 Navegación "Keyboard-First"
- [ ] **Listados Maestros (BaseCatalog)**: Implementar `selectedIndex` y captura de flechas `Up`/`Down` + `Enter` para selección sin ratón.
- [ ] **Dashboards**: Alinear todos los paneles usando `BaseDashboardPage` e implementar atajos numéricos (`1`-`4`) para KPIs y `Alt+R` para refrescar.
- [ ] **Líneas de Documentos (Ventas/Pedidos)**:
  - Navegación bidimensional con cursores entre celdas (Producto ↔ Cantidad).
  - Uso de <kbd>+</kbd> y <kbd>-</kbd> para ajustar cantidades rápidamente.
  - <kbd>Enter</kbd> para confirmar línea y crear una nueva automáticamente.
  - <kbd>Delete</kbd> para eliminación rápida de líneas.

### 1.3 Experiencia de Usuario y Ayuda
- [ ] **Ayuda en Línea**: Integrar la guía de atajos directamente en la interfaz (ej: modal de ayuda o panel lateral).
- [ ] **Refactor de Formularios**: Migrar de `fieldset/legend` a diseño de tarjetas (`.card`) y mejorar contraste de etiquetas.
- [ ] **Foco Visual**: Implementar un estilo de foco (`:focus-visible`) vibrante basado en el color primario para una clara indicación de posición.

### 1.4 Ergonomía Avanzada y Feedback (Nuevas Propuestas)
- [ ] **Sistema de Notificaciones (Toasts)**: Sustituir `alert()` por un store global de notificaciones no intrusivas (Éxito/Error/Aviso).
- [ ] **Indicadores de Campos Calculados**: Estilo visual distintivo (`.input-calculated`) para campos que el usuario no puede editar por ser lógica de backend.
- [ ] **Skeleton Loaders**: Implementar pantallas de carga progresiva que imiten la estructura de los datos para reducir la fatiga visual.
- [ ] **Chips de Filtros Activos**: Etiquetas eliminables en la parte superior de `BaseCatalog` para gestionar filtros aplicados de un vistazo.
- [ ] **Tooltips de Atajos**: Mostrar el atajo de teclado al pasar el ratón por botones o elementos interactivos para mejorar el aprendizaje.
- [ ] **Validación Inline**: Feedback visual inmediato (rojo/verde) al perder el foco en los campos de los formularios.
- [ ] **Modo Industrial (MES)**: Selector de alto contraste y fuentes aumentadas para terminales de taller con iluminación difícil.

---

## Secuencia de Implementación Sugerida (Fase UI/UX)

Para maximizar la eficiencia y el feedback temprano, se propone el siguiente orden de ejecución:

1.  **Fundamentos y Estilo**: Iconos Lucide + CSS Global (_buttons, _variables).
2.  **Claridad Operativa**: Indicadores de campos calculados + Validación Inline.
3.  **Feedback Maestro**: Sistema de Notificaciones (Toasts) global.
4.  **Navegación Eficiente**: Keyboard navigation en BaseCatalog + Chips de filtros.
5.  **Descubribilidad**: Tooltips de atajos + Guía de ayuda integrada.
6.  **Consolidación Visual**: Refactor de Dashboards + Skeleton Loaders.
7.  **Especialización**: Modo Industrial de alto contraste para el MES.

---

## 2. Facturación Consolidada Multi-Albarán

**Prioridad:** Media  
**Referencia:** [module-spec.md](../modules/sales/module-spec.md) — Fase 6

Permite agrupar múltiples albaranes de un mismo cliente en una única factura (ej: facturación mensual).

- [ ] Pantalla de selección de albaranes pendientes por cliente y rango de fechas.
- [ ] Consolidación de líneas: fusión de líneas con mismo producto y condiciones sumando cantidades.
- [ ] Vinculación N:1 de albaranes a la factura consolidada.

---

## 3. Extracción MES como Microservicio

**Prioridad:** Media  
**Referencia:** [ADR-022](adr-022-mes-microservice-extraction-strategy.md)

Extraer el módulo MES a un microservicio independiente para escalado y mantenimiento autónomo.

- [ ] Comunicación síncrona vía **gRPC**.
- [ ] Comunicación asíncrona vía **NATS JetStream**.
- [ ] Base de datos PostgreSQL independiente para MES.

---

## 4. Caché y Rendimiento

**Prioridad:** Media
- [ ] Implementar caché (Redis) de productos/variantes y reglas de precios.
- [ ] Estrategia de invalidación automática ante cambios de precios base.

---

## 5. Comunicación Asíncrona (Message Broker)

**Prioridad:** Media
- [ ] Introducir **NATS JetStream** para desacoplar módulos (Sales → MES).
- [ ] Implementar patrón **Transactional Outbox** para garantizar integridad.

---

## 6. Inteligencia de Negocio y Analítica

**Prioridad:** Baja
- [ ] Dashboards avanzados con tendencias, márgenes y KPIs históricos.
- [ ] Reportes exportables personalizados.

---

## 7. Búsqueda Global Avanzada

**Prioridad:** Baja
- [ ] Búsqueda full-text unificada (clientes, productos, pedidos) vía `Ctrl+K`.
- [ ] Autocompletado inteligente y sugerencias de búsqueda.

---

## 8. Nuevos Módulos de Negocio

**Prioridad:** Variable
- [ ] **Compras** (Purchases): Gestión de aprovisionamiento.
- [ ] **Inventario Avanzado**: Movimientos de almacén y alertas de stock.
- [ ] **Logística**: Gestión de envíos y transportistas.

---

## 9. Cobertura de Tests

**Prioridad:** Alta
- [ ] Elevar cobertura Domain al 100%.
- [ ] Elevar cobertura Application al ≥95%.
- [ ] Implementar tests E2E con Playwright para flujos críticos.

---

## 10. Gestión Avanzada de Archivos de Diseño (MES)

**Prioridad:** Media
- [ ] Vista previa (thumbnails) de archivos técnicos en paneles de producción.
- [ ] Integración con aplicaciones nativas para apertura de archivos desde el navegador.

---

*Última actualización: 2026-04-26*
