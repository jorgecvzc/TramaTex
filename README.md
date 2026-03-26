# TramaTex - ERP Textil para Microempresas

<div align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white" alt="Go 1.21+">
  <img src="https://img.shields.io/badge/Vue.js-3-4FC08D?logo=vue.js&logoColor=white" alt="Vue.js 3">
  <img src="https://img.shields.io/badge/PostgreSQL-15-4169E1?logo=postgresql&logoColor=white" alt="PostgreSQL 15">
  <img src="https://img.shields.io/badge/Docker-20.10+-2496ED?logo=docker&logoColor=white" alt="Docker">
</div>

> Un sistema ERP modular de código abierto especializado en la gestión integral de empresas de vestuario laboral y EPIs.

## 🌐 Demo Pública

Puedes explorar TramaTex en vivo sin instalar nada:

| | |
|---|---|
| **URL** | http://46.101.188.130 |
| **Email** | `admin@tramatex.local` |
| **Contraseña** | `admin123` |

> Los datos de la demo se restauran automáticamente cada domingo a las 3:00 AM UTC mediante un [workflow de GitHub Actions](.github/workflows/demo-reset.yml). Para más detalles sobre el mantenimiento de la demo, consulta la sección [Mantenimiento de la Demo](#-mantenimiento-de-la-demo).

## 💡 Sobre TramaTex

TramaTex nace con la misión de digitalizar y profesionalizar a las microempresas dedicadas a la venta de vestuario laboral y EPIs. Tradicionalmente, este sector ha carecido de herramientas adaptadas a su escala y a sus necesidades específicas de personalización y marcaje.

Este proyecto cierra esa brecha ofreciendo una solución integral para:
- Centralizar la administración de **clientes y proveedores**.
- Gestionar **catálogos complejos** con variantes JIT (tallas/colores).
- Controlar el **ciclo de ventas** (presupuestos, facturas, tickets).
- Supervisar la **ejecución en taller (MES)** de personalizaciones y arreglos.

## 🏛️ Arquitectura y Calidad Técnica

El sistema se ha desarrollado bajo los principios de **Clean Architecture** y **Domain-Driven Design (DDD)**, estructurado como un **monolito modular**. Esta elección asegura un software profesional, mantenible y escalable, diseñado para durar y evolucionar junto al negocio.

Para profundizar en los objetivos estratégicos y el modelo de negocio, consulta:
👉 **[Visión y Alcance del Proyecto](docs/architecture/project-vision-and-scope.md)**

## ✨ Funcionalidades

- **[Gestión de Identidad (IAM)](docs/modules/iam/README.md):** Sistema de autenticación y autorización basado en roles (JWT).
- **[Gestión de Terceros (Party)](docs/modules/party/README.md):** Unifica la administración de clientes y proveedores, permitiendo que una misma entidad se comporte como ambos sin duplicar información.
- **[Catálogo de Productos (Product)](docs/modules/product/README.md):** Permite la creación de productos con un sistema avanzado de variantes y atributos, tales como talla y color. Destaca la creación de variantes **"Just-In-Time" (JIT)**, que evita el engorroso trabajo de dar de alta manualmente miles de combinaciones.
- **[Motor de Precios (Pricing)](docs/modules/pricing/README.md):** Sistema para definir reglas de precios y calcular costes y precios de venta de forma dinámica, permitiendo el cálculo del precio de variantes a partir del producto base y atributos.
- **[Ciclo de Ventas (Sales)](docs/modules/sales/README.md):** Gestión completa del ciclo comercial, desde presupuestos hasta la emisión de facturas y tickets.
- **[Ejecución de Manufactura (MES)](docs/modules/mes/README.md):** Módulo para la gestión y seguimiento de las órdenes de producción en taller, referidas a las manipulaciones sobre los productos, tales como el marcado de logotipos o los arreglos.

## 🛠️ Tecnologías Utilizadas

Este proyecto se ha construido utilizando un stack de tecnologías modernas y robustas:

- **Backend:** Go (v1.21+)
- **Frontend:** Vue.js (v3) con TypeScript
- **Base de Datos:** PostgreSQL (v15+)
- **UI Framework:** Tailwind CSS
- **Contenerización:** Docker & Docker Compose

## 🚀 Cómo Empezar

Para poner en marcha una copia local del proyecto, sigue estos pasos.

### Prerrequisitos

La forma más sencilla y recomendada de ejecutar el proyecto es a través de Docker.
- **Docker**: Asegúrate de tenerlo instalado y en ejecución. Puede ser Docker Desktop (para entornos de desarrollo en Windows/macOS) o Docker Engine en un servidor Linux.

### Instalación en Desarrollo Local

1. **Clona el repositorio:**
   ```sh
   git clone git@github.com:jorgecvzc/TramaTex.git
   cd TramaTex
   ```
2. **Crea el archivo de entorno:**
   ```sh
   cp docker/.env.example docker/.env
   # Edita docker/.env con tus valores de desarrollo (o deja los que vienen por defecto)
   ```

3. **Levanta la base de datos y la API (modo desarrollo habitual):**
   ```powershell
   # Windows
   .\start-dev.ps1

   # Linux/macOS
   docker compose -f docker/docker-compose.local.yml --env-file docker/.env up -d --build
   ```
   El frontend se ejecuta aparte con hot-reload:
   ```sh
   cd apps/frontend && npm install && npm run dev
   # → http://localhost:5173
   ```

4. **O levanta el stack completo con Nginx (modo producción local):**
   ```powershell
   # Windows
   .\start-dev.ps1 -Full
   # → http://localhost:3000

   # Docker Compose directo
   docker compose -f docker/docker-compose.local.yml --env-file docker/.env --profile full up -d --build
   ```

### Instalación en Servidor de Producción

> **Arquitectura:** Los builds se realizan en GitHub Actions y las imágenes se publican en GHCR. El servidor solo descarga (`docker pull`) y ejecuta — no compila nada localmente.

#### Prerrequisitos del servidor

- Ubuntu 22.04 / 24.04
- Docker Engine + Docker Compose v2
- Git
- Usuario no-root con acceso al grupo `docker`
- Mínimo 1 GB RAM (no es necesario RAM para compilación)

#### 1. Preparar el servidor

```bash
# Crear usuario de despliegue y directorio
sudo useradd -m -s /bin/bash tramatex
sudo usermod -aG docker tramatex
sudo mkdir -p /opt/tramatex
sudo chown tramatex:tramatex /opt/tramatex

# Clonar el repositorio
sudo -u tramatex git clone https://github.com/jorgecvzc/TramaTex.git /opt/tramatex

# Generar clave SSH para GitHub Actions (sin passphrase)
ssh-keygen -t ed25519 -C "github-actions-deploy" -f ~/.ssh/deploy_key -N ""
cat ~/.ssh/deploy_key.pub >> /home/tramatex/.ssh/authorized_keys
```

#### 2. Configurar GitHub Secrets

En GitHub → **Settings → Secrets and variables → Actions**, crear:

| Secret            | Valor                                              |
|-------------------|----------------------------------------------------|
| `PROD_IP`         | IP pública del servidor (ej. `46.101.188.130`)     |
| `SSH_USER`        | Usuario SSH (ej. `tramatex`)                       |
| `SSH_PRIVATE_KEY` | Contenido del archivo `deploy_key` (clave privada) |
| `ENV_PROD`        | Contenido completo de `docker/.env` de producción  |

> **Plantilla para `ENV_PROD`:** copia `docker/.env.production.example` y rellena con tus valores reales.

#### 3. Primer despliegue

```bash
# En local: haz push a master para activar el workflow
git push origin master

# GitHub Actions ejecuta automáticamente:
# 1. Build de la imagen API    → ghcr.io/jorgecvzc/tramatex-api:latest
# 2. Build de la imagen Nginx  → ghcr.io/jorgecvzc/tramatex-frontend:latest
# 3. SSH al servidor → git pull + docker pull + docker compose up
```

Una vez completado, el sistema estará disponible en:

- **Aplicación Frontend:** `http://<IP-servidor>` (HTTP) o `https://<tu-dominio>` (tras configurar SSL)
- **API Health:** `http://<IP-servidor>/api/health`

#### 4. SSL / HTTPS (opcional pero recomendado)

```bash
# En el servidor, con certbot:
sudo apt install certbot
docker compose -f docker/docker-compose.remote.yml down
sudo certbot certonly --standalone -d tramatex.tudominio.com
# Editar docker/docker-compose.remote.yml: descomentar port 443 y volumes SSL
docker compose -f docker/docker-compose.remote.yml up -d
```

Para más detalles sobre cada entorno (local, staging, producción), consulta la:
👉 **[Guía de Despliegue Completa](docs/guides/developer/deployment-guide.md)**

## � Mantenimiento de la Demo

TramaTex sigue una filosofía **local-first**: el MVP está diseñado para instalación y operación local, sin depender de conexiones SSH ni infraestructura remota.

### Reset a datos de fábrica (local)

Para restaurar la base de datos al estado inicial con los datos de demostración, ejecuta:

```powershell
# Windows
docker compose -f docker/docker-compose.local.yml --env-file docker/.env down -v
.\start-dev.ps1

# Linux/macOS
docker compose -f docker/docker-compose.local.yml --env-file docker/.env down -v
docker compose -f docker/docker-compose.local.yml --env-file docker/.env up -d --build
```

Esto elimina el volumen de PostgreSQL y arranca desde cero. La API ejecuta automáticamente todas las migraciones y los datos semilla (`migrations/007_seed_data.sql`), restaurando el usuario admin y los datos de demostración.

### Reset automático de la demo pública

La instancia de demostración desplegada en producción se resetea automáticamente cada domingo a las 3:00 AM UTC mediante el workflow [`.github/workflows/demo-reset.yml`](.github/workflows/demo-reset.yml). También puede ejecutarse manualmente desde la pestaña **Actions** del repositorio en GitHub.

## �📂 Estructura del Proyecto

El proyecto sigue una estructura de monolito modular para separar las responsabilidades:

```
TramaTex/
├── apps/                 # Contiene las aplicaciones principales
│   ├── tramatex-api/     # El backend en Go (API REST)
│   └── frontend/         # El frontend en Vue.js (SPA)
│
├── docs/                 # Toda la documentación del proyecto
│   ├── architecture/     # Decisiones de arquitectura (ADRs), diagramas C4
│   ├── guides/           # Guías para desarrolladores y usuarios
│   ├── modules/          # Documentación funcional de cada módulo
│   └── log/              # Registro de trabajo (sprints, etc.)
│
├── agents/               # Definiciones para asistentes de IA
└── docker/               # Archivos de configuración de Docker Compose
```


## 📚 Documentación

El proyecto incluye documentación exhaustiva que cubre varios aspectos. Este `README.md` sirve como punto de partida para navegarla.

-   **Decisiones de Arquitectura (ADRs):** Razón y justificación detrás de las decisiones arquitectónicas clave del proyecto.
    -   [Visión y Alcance del Proyecto](docs/architecture/project-vision-and-scope.md)
    -   [Índice de ADRs](docs/architecture/adrs/README.md)
-   **Guías de Desarrollo:** Instrucciones y buenas prácticas para contribuir al código.
    -   [Guía de Inicio Rápido](docs/guides/quick-start.md)
    -   [Estándares de Código y Estilo](docs/guides/code-and-style-standards.md)
-   **Documentación de Módulos:** Especificaciones detalladas para cada Módulo de Negocio (Bounded Context).
    -   [Resumen de Módulos](docs/modules/README.md)
    -   [Módulo Party](docs/modules/party/README.md)
    -   [Módulo Product](docs/modules/product/README.md)
    -   [Módulo Pricing](docs/modules/pricing/README.md)
    -   [Módulo Sales](docs/modules/sales/README.md)
    -   [Módulo IAM](docs/modules/iam/README.md)
    -   [Módulo MES](docs/modules/mes/README.md)
-   **Registros y Estado del Proyecto:** Historial de trabajo, resúmenes de sprints y el estado actual del proyecto.
    -   [Estado del Proyecto](docs/log/project-status.md)
    -   [Registros de Sprints](docs/log/sprints/README.md)
    -   [Registro de Sesiones](docs/log/session-log.md)

## 📄 Licencia

Este proyecto está bajo la [Licencia MIT](LICENSE.md).

## 🌱 Proyecto de Scaffolding

Como resultado del desarrollo de TramaTex, se ha creado una tecnología de *scaffolding* para estandarizar la creación de futuros proyectos. Este sistema, ubicado en la carpeta `project-scaffolding/`, es una guía apoyada en Inteligencia Artificial para la generación de nuevos proyectos desde una base sólida y estandarizada.

El objetivo principal es permitir a los equipos de desarrollo enfocarse en la lógica de negocio desde el primer día, al proporcionar automáticamente:

-   Una estructura de directorios coherente.
-   Configuraciones iniciales para CI/CD.
-   Estrategias de testing.
-   Templates para documentación, ADRs y guías de desarrollo.

Al automatizar la configuración inicial, este proyecto de scaffolding fomenta la consistencia, las mejores prácticas y acelera significativamente el arranque de nuevos desarrollos.

## 👥 Autores

- **Jorge Cortés Villalba** - *Product Owner, Arquitectura y Desarrollo*
- **AI Assistant (Gemini, Claude, Copilot & Perplexity)** - *Copilotos Técnicos y Asistentes de Desarrollo*
