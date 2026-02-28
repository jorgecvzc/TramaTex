# 📖 Portal de Documentación del Proyecto

Bienvenido al centro de documentación técnica de **TramaTex**. Este directorio centraliza todo el conocimiento acumulado sobre la arquitectura, el diseño y la evolución del sistema.

---

## 🗺️ Mapa de Navegación

### 🏛️ [Arquitectura](./architecture/README.md)
Decisiones estratégicas y visión técnica del sistema.
- **[Visión y Alcance](./architecture/project-vision-and-scope.md):** El "por qué" y el "qué" del proyecto.
- **[Registro de Decisiones (ADRs)](./architecture/adrs/README.md):** Historial de elecciones técnicas justificadas.
- **[Glosario de Términos](./architecture/glossary.md):** Lenguaje ubicuo utilizado en todo el proyecto.

### 🧩 [Módulos de Negocio](./modules/README.md)
Documentación detallada de cada Bounded Context (Domain Models, Use Cases, Specs).
- **[IAM](./modules/iam/README.md):** Identidad y Acceso.
- **[Party](./modules/party/README.md):** Gestión de Terceros.
- **[Product](./modules/product/README.md):** Catálogo y Variantes.
- **[Pricing](./modules/pricing/README.md):** Motor de Precios.
- **[Sales](./modules/sales/README.md):** Ciclo Comercial.
- **[MES](./modules/mes/README.md):** Ejecución de Manufactura.

### 📚 [Guías y Estándares](./guides/README.md)
Manuales para desarrolladores y estándares de calidad.
- **[Inicio Rápido](./guides/quick-start.md):** Cómo empezar a desarrollar en 5 minutos.
- **[Estándares de Código](./guides/code-and-style-standards.md):** Clean Code, Testing y Go/Vue patterns.
- **[Estándares de Documentación](./guides/documentation-standards.md):** Reglas para mantener este portal.

### 📝 [Bitácora y Estado](./log/README.md)
Trazabilidad del desarrollo y situación actual.
- **[Estado del Proyecto](./log/project-status.md):** Hitos completados y foco actual.
- **[Registro de Sesiones](./log/session-log.md):** Historial de trabajo diario de los agentes.
- **[Resúmenes de Sprints](./log/sprints/README.md):** Evolución temporal del proyecto.

---

## 🎯 Estándares de este Directorio

Para mantener la coherencia de este portal, seguimos estas reglas críticas:

1.  **Nombres de Archivos/Carpetas:** Siempre en **Inglés** y `kebab-case`.
2.  **Contenido:** Siempre en **Castellano** (salvo términos técnicos).
3.  **Formato:** Markdown (`.md`) con diagramas en **Mermaid**.
4.  **Trazabilidad:** Cada cambio en el código debe reflejarse en su documentación correspondiente.

Para más detalles, consulta la **[Guía de Estándares de Documentación](./guides/documentation-standards.md)**.