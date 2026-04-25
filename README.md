# 🏛️ TramaTex: Gestión Inteligente para la Microempresa Textil

<div align="center">
  <img src="https://img.shields.io/badge/Proyecto-TFM-1b3a6b?style=for-the-badge" alt="TFM">
  <img src="https://img.shields.io/badge/Backend-Go_1.23-1b3a6b?style=flat-square" alt="Go">
  <img src="https://img.shields.io/badge/Frontend-Vue.js_3-1b3a6b?style=flat-square" alt="Vue">
  <img src="https://img.shields.io/badge/Base_de_Datos-PostgreSQL_15-1b3a6b?style=flat-square" alt="PostgreSQL">
  <img src="https://img.shields.io/badge/Licencia-MIT-E6B800?style=flat-square" alt="Licencia">
</div>

<div align="center">
  <a href="https://github.com/jorgecvzc/TramaTex/actions/workflows/backend.yml">
    <img src="https://github.com/jorgecvzc/TramaTex/actions/workflows/backend.yml/badge.svg" alt="Backend CI">
  </a>
  <a href="https://github.com/jorgecvzc/TramaTex/actions/workflows/frontend.yml">
    <img src="https://github.com/jorgecvzc/TramaTex/actions/workflows/frontend.yml/badge.svg" alt="Frontend CI">
  </a>
</div>

---

## Descripción General
TramaTex es un sistema integral de gestión (ERP/MES) de código abierto, desarrollado como **Trabajo Fin de Máster (TFM)**. Su misión es digitalizar y profesionalizar las operaciones de microempresas dedicadas al vestuario laboral y EPIs, aportando soluciones a necesidades críticas de personalización (marcajes, arreglos y tallajes complejos).

El sistema apuesta por la **soberanía tecnológica** mediante una arquitectura **Local-First**, garantizando que la empresa mantenga el control total de sus datos y opere con total independencia de la nube.

## 📺 Presentación Rápida
Explora de un vistazo la propuesta de valor y los pilares de ingeniería en nuestra presentación corporativa interactiva:
👉 **[Ver Presentación Visual del Proyecto](https://jorgecvzc.github.io/TramaTex/presentations/presentation.html)**

---

## 🌐 Demo Pública
Accede a una instancia funcional para evaluar el sistema de forma inmediata:

| Recurso | Detalle |
| :--- | :--- |
| **URL de acceso** | [http://46.101.188.130](http://46.101.188.130) |
| **Usuario Administrador** | `admin@tramatex.local` |
| **Contraseña de acceso** | `admin123` |

> 🔄 **Mantenimiento Automático:** Para garantizar una experiencia limpia a todos los evaluadores, la base de datos de la demo se restaura a su estado inicial **cada domingo a las 3:00 AM UTC**.

---

## Stack Tecnológico
Selección técnica orientada al rendimiento, la precisión y la longevidad del software:

*   **Backend:** [Go (Golang)](https://go.dev/) siguiendo los patrones de **Clean Architecture** y **Domain-Driven Design (DDD)**.
*   **Frontend:** [Vue.js 3](https://vuejs.org/) con TypeScript y un **Sistema de Diseño propio** basado en **Vanilla CSS**, optimizado para la agilidad en el día a día.
*   **Base de Datos:** [PostgreSQL](https://www.postgresql.org/) con manejo de precisión decimal para integridad financiera.
*   **DevOps:** [Docker](https://www.docker.com/) para una orquestación estandarizada y multiplataforma.

---

## Instalación y Ejecución
El proyecto fomenta el uso de sistemas de código abierto, proporcionando automatización tanto para entornos Windows como Linux.

### 🐧 Entornos Linux (Recomendado para Producción/OSS)
1.  **Instalación:** `chmod +x scripts/install.sh && ./scripts/install.sh`
2.  **Ejecución:** `docker compose -f docker/docker-compose.local.yml --env-file docker/.env up -d`

### 🪟 Entornos Windows (Desarrollo Ágil)
1.  **Configurar entorno:** `cp docker/.env.example docker/.env`
2.  **Lanzar servicios:** `.\start-dev.ps1` (o `.\start-dev.ps1 -Full` para el stack completo).

### 🛠️ Herramientas de Gestión
Consulta el índice de utilidades para tareas de mantenimiento, migraciones o auditoría:
👉 **[Guía Maestra de Scripts y Utilidades](docs/guides/developer/scripts-index.md)**
👉 **[Pipeline CI/CD (GitHub Actions)](docs/guides/developer/ci-cd.md)**

---

## Estructura del Proyecto
Organización modular por responsabilidades para facilitar la mantenibilidad:

*   `apps/` - Aplicaciones: Backend (`tramatex-api`) y Frontend (`frontend`).
*   `docs/` - **Árbol de Conocimiento**: El núcleo documental del proyecto.
*   `project-scaffolding/` - **Ecosistema Metodológico**: Herramientas de estandarización.
*   `docker/` - Definiciones de infraestructura y contenedores.
*   `scripts/` - Utilidades operativas para despliegue y seguridad.

---

## Funcionalidades Principales
Ciclo de negocio modular y completo:

1.  **[Gestión de Identidad (IAM)](docs/modules/iam/README.md):** Seguridad robusta y control de acceso por roles.
2.  **[Gestión de Terceros (Party)](docs/modules/party/README.md):** Clientes y proveedores unificados.
3.  **[Catálogo de Productos (Product)](docs/modules/product/README.md):** Variantes dinámicas **Just-In-Time (JIT)**.
4.  **[Motor de Precios (Pricing)](docs/modules/pricing/README.md):** Cálculos con precisión decimal y reglas de negocio.
5.  **[Ciclo de Ventas (Sales)](docs/modules/sales/README.md):** Trazabilidad desde el presupuesto a la factura.
6.  **[Control de Producción (MES)](docs/modules/mes/README.md):** Gestión y seguimiento de taller en tiempo real.

---

## 🌱 Metodología y Scaffolding
TramaTex se ha construido utilizando un sistema de **Scaffolding** desarrollado en paralelo (`project-scaffolding/`). No es una simple utilidad, sino el motor que genera el ecosistema de agentes, carpetas y guías necesario para que el desarrollo del proyecto sea metódico y posible bajo estándares de ingeniería.

Pieza fundamental de esta metodología es la **[Bitácora de Sesiones](docs/log/session-log.md)**, que permite una trazabilidad técnica exhaustiva de cada decisión y avance realizado en el proyecto.

Para conocer el sistema que hace posible esta estructura, visita:
👉 **[Documentación de Scaffolding y Metodología](project-scaffolding/README.md)**

---

## 📚 Navegación del Árbol de Documentación
Explora el conocimiento técnico y estratégico siguiendo estas ramas principales:

*   🏛️ **[Arquitectura y Decisiones](docs/architecture/README.md):** ADRs, visión técnica y glosario.
*   📦 **[Detalle de Módulos](docs/modules/README.md):** Especificaciones funcionales de cada dominio.
*   🛠️ **[Guías y Estándares](docs/guides/README.md):** Normas de código y operativa.
*   🚀 **[Evolución Post-MVP](docs/post-mvp/post-mvp-roadmap.md):** Hoja de ruta estratégica tras el TFM.

---

## 📄 Licencia
Este proyecto es software libre y se distribuye bajo la **[Licencia MIT](LICENSE.md)**, garantizando la libertad de uso, modificación y distribución para la comunidad y la industria textil.

## 👥 Autores
- **Jorge Cortés Villalba** - *Diseño, Arquitectura e Implementación.*
- **AI Collaborative Ecosystem** - *Gemini, Claude, GitHub Copilot, ChatGPT & Perplexity en roles de copilotos técnicos.*

---
© 2026 TramaTex - Software Engineering for Textile Industry.
