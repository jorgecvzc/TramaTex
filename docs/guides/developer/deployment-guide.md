# Guía de Despliegue — TramaTex

> Referencia completa para desplegar TramaTex en los 4 entornos del proyecto.

---

## Arquitectura de Entornos

```
LOCAL (Windows)  →  STAGING / pcele (LAN)  →  PRODUCTION (DigitalOcean)
Docker Desktop      SSH directo (manual)     GitHub Actions (auto)
develop branch      staging branch           master branch
```

| Entorno     | Servidor            | Acceso                | Puerto | Branch    |
|-------------|---------------------|-----------------------|--------|-----------|
| Local       | Docker Desktop (Win)| localhost              | 3000   | develop   |
| Staging     | pcele (192.168.0.20)| SSH (user: ele, manual)| 80     | staging   |
| Producción  | DigitalOcean        | GitHub Actions + SSH   | 80/443 | master    |

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

## 2. Entorno Staging (pcele LAN)

### Configuración inicial del servidor

```bash
# En pcele (192.168.0.20), como usuario ele:
sudo mkdir -p /opt/tramatex
sudo chown ele:ele /opt/tramatex
cd /opt/tramatex

# Clonar el repositorio
git clone <repo-url> .
git checkout staging

# Crear el archivo de entorno
cp docker/.env.staging.example docker/.env
# Editar docker/.env — cambiar contraseñas y JWT_SECRET
nano docker/.env
```

### Prerequisitos en pcele
- Docker Engine + Docker Compose v2 instalados
- Git
- Acceso SSH desde Windows: `ssh ele@pcele` o `ssh ele@192.168.0.20`

### Desplegar

```bash
# Desde Windows (script remoto por SSH):
powershell -ExecutionPolicy Bypass -File .\scripts\rebuild-staging-remote.ps1

# Si necesitas desplegar una rama temporal (FORZANDO construcción desde código fuente):
# El flag -BuildSource es CRÍTICO para ver cambios de ramas de feature que no están en GHCR.
powershell -ExecutionPolicy Bypass -File .\scripts\rebuild-staging-remote.ps1 -CheckoutRef origin/feature/mi-rama -BuildSource

# Ejecutando directamente en pcele (Linux):
./scripts/rebuild-staging-remote.sh --checkout-ref origin/staging --build-source
```

Opciones útiles:

- En Windows: `-NoCheckout`, `-PreserveDatabase`, `-SkipImageRemove`, `-BuildSource`
- En Linux: `--no-checkout`, `--preserve-database`, `--skip-image-remove`, `--build-source`

Nota: por defecto, el compose remoto intenta usar imagenes publicadas (GHCR). Si estás en una rama de desarrollo (feature), **debes usar `--build-source`** para compilar las imágenes localmente en pcele usando los Dockerfiles del repo.

> **Nota:** pcele está en la LAN (192.168.0.20), no es accesible desde internet.
> GitHub Actions no puede desplegar ahí. El deploy es manual desde una máquina en la misma red.
> Al hacer push a `staging`, GitHub Actions ejecuta CI (build + tests) para validar.

### Acceder
- **Frontend**: http://192.168.0.20
- **API Health**: http://192.168.0.20/api/health

### Ver logs

```bash
ssh ele@pcele "cd /opt/tramatex && docker compose -f docker/docker-compose.remote.yml logs -f"
```

---

## 3. Entorno Producción (DigitalOcean)

> **Estrategia de build:** Los builds NO ocurren en el Droplet (1 GB RAM → OOM con `npm run build`).
> Las imágenes se construyen en el runner de GitHub Actions (7 GB RAM) y se publican en GHCR.
> El Droplet solo ejecuta `docker pull` + `docker compose up`.

### Configuración inicial del Droplet

```bash
# En el Droplet (Ubuntu 22.04/24.04):

# Crear usuario de despliegue
sudo useradd -m -s /bin/bash tramatex
sudo usermod -aG docker tramatex
sudo chown tramatex:tramatex /home/tramatex

# Crear directorio y clonar el repo
sudo mkdir -p /opt/tramatex
sudo chown tramatex:tramatex /opt/tramatex
sudo -u tramatex git clone https://github.com/jorgecvzc/TramaTex.git /opt/tramatex
cd /opt/tramatex && git checkout master

# El archivo docker/.env será escrito por GitHub Actions en cada deploy
# Para la primera vez, crearlo manualmente (opcional):
cp docker/.env.production.example docker/.env
# nano docker/.env  → Configurar con valores de producción reales
```

### Generar clave SSH para el deploy

```bash
# En local (o en el Droplet):
ssh-keygen -t ed25519 -C "github-actions-tramatex" -f ./deploy_key -N ""

# Añadir la clave pública al servidor
cat deploy_key.pub >> /home/tramatex/.ssh/authorized_keys
chmod 600 /home/tramatex/.ssh/authorized_keys

# El contenido de 'deploy_key' (privada) va en el GitHub Secret SSH_PRIVATE_KEY
```

### GitHub Secrets necesarios

Configurar en **Settings → Secrets → Actions** del repositorio:

| Secret            | Descripción                                                     | Ejemplo                      |
|-------------------|-----------------------------------------------------------------|------------------------------|
| `SSH_PRIVATE_KEY` | Clave privada SSH para acceder al Droplet (ed25519, sin pass)   | `-----BEGIN OPENSSH...`      |
| `SSH_USER`        | Usuario SSH en el Droplet                                       | `tramatex`                   |
| `PROD_IP`         | IP del Droplet de producción                                    | `46.101.188.130`             |
| `ENV_PROD`        | Contenido completo del archivo `docker/.env` de producción      | (ver `docker/.env.production.example`) |

> `GITHUB_TOKEN` es automático — lo usa GitHub Actions para autenticarse en GHCR. No hay que crearlo.

### GitHub Environments

Crear en **Settings → Environments**:
1. **production** — con revisores requeridos (opcional, para control adicional)

### Flujo del workflow deploy-production.yml

```
push a master
    │
    ├── build-api (ubuntu-latest, 7GB RAM)
    │     docker build --no-cache -f Dockerfile .
    │     docker push ghcr.io/jorgecvzc/tramatex-api:latest
    │
    ├── build-frontend (ubuntu-latest, 7GB RAM)
    │     docker build --no-cache -f Dockerfile.frontend .
    │     docker push ghcr.io/jorgecvzc/tramatex-frontend:latest
    │
    └── deploy (needs: build-api, build-frontend)
          SSH al Droplet →
            git pull master
            printf '%s' "$ENV_PROD" > docker/.env
            docker login ghcr.io
            docker compose pull          ← solo descarga las imágenes ya construidas
            docker compose up -d --force-recreate
            docker image prune -f
```

### Promover a producción manualmente

```bash
# Haz push a master (cualquier commit desencadena el deploy):
git push origin master

# O desde una rama de feature:
git push origin mi-rama:master
```

### Verificar el deploy

```bash
# Desde local:
ssh -i tmp/do-setup/deploy_final tramatex@46.101.188.130 \
  "docker ps --format 'table {{.Names}}\t{{.Status}}'"

# Health check API:
curl http://46.101.188.130/api/health
# → {"status":"healthy"}
```

### Acceder al entorno de producción
- **Frontend**: http://46.101.188.130 (HTTP) o https://tu-dominio.com (con SSL)
- **API Health**: http://46.101.188.130/api/health

### Ver logs en producción

```bash
ssh -i tmp/do-setup/deploy_final tramatex@46.101.188.130 \
  "cd /opt/tramatex && docker compose -f docker/docker-compose.remote.yml logs -f"
```

---

## 4. Reseteo Semanal de la Demo (Mantenimiento)

Para garantizar que la demo pública siempre esté limpia y funcional para nuevos evaluadores, el proyecto incluye un flujo de mantenimiento automático.

- **Frecuencia:** Todos los domingos a las 03:00 AM UTC.
- **Flujo (`demo-reset.yml`):**
  1. Detención completa del stack.
  2. Borrado de volúmenes persistentes (`docker compose down -v`).
  3. Re-arranque y ejecución de migraciones.
  4. Verificación de carga de datos iniciales (Seed Data).
- **Verificación Manual:** Se puede disparar manualmente desde la pestaña "Actions" en GitHub seleccionando "Demo Weekly Reset".

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
│   ├── .env.staging.example              # Template para staging (pcele)
│   └── .env.production.example         # Template para DO producción
└── .github/workflows/
    ├── deploy-staging.yml              # CI on push to staging (build + tests)
    └── deploy-production.yml           # Auto-deploy on push to master (DO)
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
