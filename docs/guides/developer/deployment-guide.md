# Guía de Despliegue — TramaTex

> Referencia completa para desplegar TramaTex en los 4 entornos del proyecto.

---

## Arquitectura de Entornos

```
LOCAL (Windows)  →  PCELE (LAN)  →  STAGING (DigitalOcean)  →  PRODUCTION (DigitalOcean)
Docker Desktop      SSH directo      GitHub Actions              GitHub Actions
develop branch      develop branch   staging branch              master branch
```

| Entorno    | Servidor            | Acceso               | Puerto | Branch    |
|------------|---------------------|----------------------|--------|-----------|
| Local      | Docker Desktop (Win)| localhost             | 3000   | develop   |
| pcele      | 192.168.0.20        | SSH (user: ele)       | 80     | develop   |
| Staging    | DigitalOcean        | GitHub Actions + SSH  | 80/443 | staging   |
| Producción | DigitalOcean        | GitHub Actions + SSH  | 80/443 | master    |

---

## Stack de Servicios Docker

Todos los entornos ejecutan los mismos 3 contenedores:

| Servicio   | Imagen                  | Puerto interno | Descripción                     |
|------------|-------------------------|----------------|---------------------------------|
| postgres   | postgres:15-alpine      | 5432           | Base de datos PostgreSQL        |
| api        | Dockerfile (Go)         | 8080           | API REST (no expuesto en remoto)|
| frontend   | Dockerfile.frontend     | 80             | Nginx + Vue SPA + reverse proxy |

**Nginx** es el punto de entrada único:
- `/api/*` y `/auth/*` → reverse proxy al contenedor `api:8080`
- Todo lo demás → Vue SPA con fallback `index.html` (Vue Router history mode)

---

## 1. Entorno Local (Windows)

### Prerequisitos
- Docker Desktop instalado y corriendo
- Git

### Configuración inicial

```powershell
# Clonar el repo (si no lo tienes)
git clone <repo-url>
cd TramaTex

# Crear docker/.env desde el ejemplo
Copy-Item docker/.env.example docker/.env
# Editar docker/.env con los valores de desarrollo
```

### Levantar solo DB + API (modo desarrollo normal)

```powershell
# Opción 1: Script PowerShell
.\start-dev.ps1

# Opción 2: Makefile
make docker-up

# El frontend se ejecuta aparte con hot-reload:
cd apps/frontend
npm run dev
# → http://localhost:5173
```

### Levantar stack completo (DB + API + Frontend/Nginx)

```powershell
# Opción 1: Script PowerShell con -Full
.\start-dev.ps1 -Full
# → Frontend Nginx en http://localhost:3000

# Opción 2: Docker Compose directo
docker compose -f docker/docker-compose.local.yml --env-file docker/.env --profile full up -d --build
```

### Detener

```powershell
.\stop-dev.ps1
# o: make docker-down
```

---

## 2. Entorno pcele (LAN Pre-Staging)

### Configuración inicial del servidor

```bash
# En pcele (192.168.0.20), como usuario ele:
sudo mkdir -p /opt/tramatex
sudo chown ele:ele /opt/tramatex
cd /opt/tramatex

# Clonar el repositorio
git clone <repo-url> .

# Crear el archivo de entorno
cp docker/.env.pcele.example docker/.env
# Editar docker/.env — cambiar contraseñas y JWT_SECRET
nano docker/.env
```

### Prerequisitos en pcele
- Docker Engine + Docker Compose v2 instalados
- Git
- Acceso SSH desde Windows: `ssh ele@pcele` o `ssh ele@192.168.0.20`

### Desplegar

```bash
# Desde Windows (vía Makefile):
make deploy ENV=pcele

# O manualmente vía SSH:
ssh ele@pcele "cd /opt/tramatex && git pull origin develop && \
  docker compose -f docker/docker-compose.remote.yml --env-file docker/.env build && \
  docker compose -f docker/docker-compose.remote.yml --env-file docker/.env up -d"
```

### Acceder
- **Frontend**: http://192.168.0.20
- **API Health**: http://192.168.0.20/api/health

### Ver logs

```bash
ssh ele@pcele "cd /opt/tramatex && docker compose -f docker/docker-compose.remote.yml logs -f"
```

---

## 3. Entorno Staging (DigitalOcean)

### Configuración inicial del Droplet

```bash
# Crear Droplet en DigitalOcean (Ubuntu 22.04, Docker pre-instalado)
# Conectar por SSH con la IP del Droplet

# En el Droplet:
sudo mkdir -p /opt/tramatex
sudo chown $USER:$USER /opt/tramatex
cd /opt/tramatex

# Clonar el repo
git clone <repo-url> .
git checkout staging

# El archivo docker/.env será escrito por GitHub Actions
# Para la primera vez, crear manualmente:
cp docker/.env.production.example docker/.env
nano docker/.env  # Configurar con valores de staging
```

### GitHub Secrets necesarios

Configurar en **Settings → Secrets → Actions** del repositorio:

| Secret            | Descripción                                          | Ejemplo                    |
|-------------------|------------------------------------------------------|----------------------------|
| `SSH_PRIVATE_KEY` | Clave privada SSH para acceder al Droplet            | `-----BEGIN OPENSSH...`    |
| `SSH_USER`        | Usuario SSH en el Droplet                            | `root` o `deploy`          |
| `STAGING_IP`      | IP del Droplet de staging                            | `164.90.xxx.xxx`           |
| `ENV_STAGING`     | Contenido completo del archivo docker/.env de staging| (ver docker/.env.production.example) |

### GitHub Environments

Crear dos environments en **Settings → Environments**:
1. **staging** — sin protección adicional
2. **production** — con revisores requeridos (opcional)

### Flujo de despliegue

```bash
# Desde local, push de develop a staging:
make deploy ENV=staging
# Equivale a: git push origin develop:staging

# GitHub Actions detecta push a staging y:
# 1. Conecta por SSH al Droplet
# 2. Escribe docker/.env desde el secret ENV_STAGING
# 3. docker compose build + up
# 4. Verifica health check
# 5. Limpia imágenes antiguas
```

### Crear la rama staging

```bash
# Primera vez — crear staging desde develop:
git checkout develop
git push origin develop:staging
```

---

## 4. Entorno Producción (DigitalOcean)

Mismo proceso que staging, con sus propios secrets:

| Secret        | Descripción                                |
|---------------|--------------------------------------------|
| `PROD_IP`     | IP del Droplet de producción               |
| `ENV_PROD`    | Contenido del docker/.env de producción    |

### Flujo de despliegue

```bash
# Promover staging a producción:
make deploy ENV=prod
# Equivale a: git push origin staging:master
# (pide confirmación antes de ejecutar)
```

---

## 5. SSL/HTTPS (Producción)

### Opción A: Let's Encrypt con Certbot (recomendado)

```bash
# En el Droplet de producción:

# 1. Instalar certbot
sudo apt install certbot

# 2. Obtener certificado (con Nginx detenido temporalmente)
docker compose -f docker/docker-compose.remote.yml --env-file docker/.env down
sudo certbot certonly --standalone -d tramatex.tudominio.com

# 3. Habilitar SSL en docker-compose.remote.yml:
#    - Descomentar el port 443:443
#    - Descomentar los volumes de letsencrypt y nginx-ssl.conf
#    - Descomentar el volume certbot_webroot

# 4. Reiniciar
docker compose -f docker/docker-compose.remote.yml --env-file docker/.env up -d
```

### Opción B: Cloudflare Proxy (alternativa simple)

Si usas Cloudflare como proxy DNS:
- El SSL termina en Cloudflare → Nginx recibe HTTP en puerto 80
- No necesitas `nginx-ssl.conf` — deja la configuración HTTP actual
- Configura Cloudflare con SSL/TLS mode: "Full"

### Renovación automática de certificados

```bash
# Añadir cron job para renovación:
sudo crontab -e
# Añadir:
0 3 * * * certbot renew --quiet && docker exec tramatex_frontend nginx -s reload
```

---

## Estructura de Archivos de Despliegue

```
TramaTex/
├── Dockerfile                          # Go API (multi-stage build)
├── Dockerfile.frontend                 # Vue + Nginx (multi-stage build)
├── Makefile                            # Comandos: docker-up, deploy, etc.
├── start-dev.ps1                       # Script Windows (dev + full mode)
├── stop-dev.ps1                        # Script Windows (detener)
├── docker/
│   ├── docker-compose.yml              # Base (standalone, development)
│   ├── docker-compose.local.yml        # Local Docker Desktop
│   ├── docker-compose.remote.yml       # pcele + DigitalOcean
│   ├── nginx.conf                      # Nginx HTTP (default)
│   ├── nginx-ssl.conf                  # Nginx HTTPS (Let's Encrypt)
│   ├── .env                            # Variables actuales (NO committear)
│   ├── .env.example                    # Template genérico
│   ├── .env.pcele.example              # Template para pcele
│   └── .env.production.example         # Template para DO staging/prod
└── .github/workflows/
    ├── deploy-staging.yml              # Auto-deploy on push to staging
    └── deploy-production.yml           # Auto-deploy on push to master
```

---

## Troubleshooting

### Docker build falla con error de red
```
dialing docker-images-prod...cloudflarestorage.com:443: connectex: timeout
```
Problema de conectividad con Docker Hub. Verificar VPN/proxy o reintentar.

### API no responde después de deploy
```bash
# Ver logs del contenedor API:
docker logs tramatex_api --tail 50

# Verificar que las migraciones se ejecutaron:
docker exec tramatex_api ls /app/migrations

# Verificar conectividad DB desde el contenedor API:
docker exec tramatex_api wget -qO- http://localhost:8080/api/health
```

### Frontend muestra página en blanco
```bash
# Verificar que el build de Vue se completó:
docker exec tramatex_frontend ls /usr/share/nginx/html/

# Debe contener: index.html, assets/
# Verificar config de Nginx:
docker exec tramatex_frontend nginx -t
```

### Nginx 502 Bad Gateway en /api/
- El contenedor `api` no está corriendo o no está healthy
- Verificar: `docker ps` — el contenedor api debe mostrar "(healthy)"
- Verificar red: los contenedores deben estar en la misma red Docker (`tramatex_network`)
