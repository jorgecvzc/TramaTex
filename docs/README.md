# 📖 Portal de Conocimiento TramaTex

Bienvenido al centro de conocimiento estratégico de TramaTex. Este portal no es un espejo del código, sino una guía conceptual diseñada para entender la lógica de negocio, la arquitectura y los flujos operativos del sistema.

---

## 🏛️ Gestión del Conocimiento (Plan Maestro)

Antes de navegar o contribuir a esta documentación, es obligatorio conocer nuestro **[Plan Maestro de Documentación](./guides/documentation-standards.md)**. 

**Nuestros Pilares:**
1. **Comportamiento Primero:** Explicamos el *porqué* y el *cómo* del negocio, no la lista de campos del código.
2. **Anti-Redundancia:** Si está en el código y es evidente, no se duplica aquí.
3. **Soberanía del Dominio:** El conocimiento está organizado por contextos de negocio (Bounded Contexts).

---

## 🗺️ Mapa de Navegación

### 🏛️ [Arquitectura y Estrategia](./architecture/README.md)
Visión de alto nivel, principios técnicos y registro de decisiones críticas.
- **[Visión del Sistema](./architecture/architecture-vision.md)** | **[Glosario Ubicuo](./architecture/glossary.md)** | **[Registro de Decisiones (ADRs)](./architecture/adrs/README.md)**

### 🧩 [Módulos de Dominio](./modules/README.md)
Guías de comportamiento y reglas de negocio por cada área del sistema.
- **[Party (Terceros)](./modules/party/README.md):** Lógica de identidad y relaciones.
- **[Product (Catálogo)](./modules/product/README.md):** Variantes dinámicas y herencia.
- **[Pricing (Precios)](./modules/pricing/README.md):** Motor de reglas económicas.
- **[Sales (Ventas)](./modules/sales/README.md):** Ciclo documental y comercial.
- **[MES (Producción)](./modules/mes/README.md):** Ejecución y trazabilidad en taller.

### 📚 [Guías Operativas](./guides/README.md)
Estándares de ingeniería, procesos de despliegue y manuales de usuario.
- **[Inicio Rápido](./guides/quick-start.md)** | **[Estándares de Ingeniería](./guides/code-and-style-standards.md)** | **[Estrategia de Búsqueda Global](./guides/developer/global-search-strategy.md)**

### 🚀 [Hoja de Ruta Post-MVP](./post-mvp/post-mvp-roadmap.md)
Mejoras, funcionalidades y módulos planificados tras el MVP.
- **Primera prioridad:** Unificación UI/UX y Sistema de Diseño.

### 🎓 [Presentaciones](./presentations/slides_spec.md)
Material de presentación y defensa del proyecto.
- **[Presentación TramaTex](./presentations/tramatex-presentation.md)**

### 📝 [Bitácora del Proyecto](./log/README.md)
Trazabilidad histórica del desarrollo y estado actual de los hitos.
- **[Estado del Proyecto](./log/project-status.md)** | **[Historial de Sprints](./log/sprints/README.md)**

---
**Directiva de Calidad:** Si encuentras documentación obsoleta o redundante con el código, por favor aplica el proceso de [Refinamiento](./guides/documentation-standards.md#4-flujo-de-mantenimiento-y-refinamiento).
