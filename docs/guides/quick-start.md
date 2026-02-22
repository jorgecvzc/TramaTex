# 🚀 Quick Start Guide - TramaTex ERP

**Última actualización:** 2026-02-09

---

## ✅ Credenciales de Acceso VERIFICADAS

```
Email:    admin@tramatex.local
Password: admin123
```

**Estado:** ✅ Verificado y funcionando en la API
**Hash bcrypt:** `$2a$10$Gd8.JP/L3j.vNvap81EpjuX7G4u5KKLmf10TSxmu779Mq/HdC/B9e`

---

## 🚀 Iniciar la Aplicación

### Paso 1: Iniciar Backend (API + Base de Datos)

```bash
# Desde la raíz del proyecto TramaTex
docker-compose -f docker/docker-compose.yml up -d
```

**Verificar que estén corriendo:**
```bash
docker ps
```

Deberías ver:
- `tramatex_api` (puerto 8080)
- `tramatex_db` (puerto 5432)

### Paso 2: Iniciar Frontend

```bash
# Abrir otra terminal
cd apps/frontend
npm run dev
```

El frontend estará disponible en: **http://localhost:5173** (o 5174 si 5173 está ocupado)

---

## 🔐 Probar el Login

### Opción 1: Desde el Navegador (Recomendado)

1. Abre: **http://localhost:5173/login** (o **http://localhost:5174/login** si el puerto cambió)
2. Ingresa:
   - **Email:** `admin@tramatex.local`
   - **Password:** `admin123`
3. Click en **"Ingresar"**
4. Deberías ser redirigido a `/dashboard`

### Opción 2: Desde la API directamente (Testing)

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@tramatex.local","password":"admin123"}'
```

**Respuesta esperada:**
```json
{
  "user": {
    "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
    "email": "admin@tramatex.local",
    "role": "admin"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 900
}
```

---

## 🎯 Rutas Disponibles

### Frontend (Vue.js)

- **Login:** http://localhost:5173/login (o 5174)
- **Dashboard:** http://localhost:5173/dashboard
- **Productos (Lista):** http://localhost:5173/products
- **Producto (Detalle):** http://localhost:5173/products/1
- **Entidades:** http://localhost:5173/parties
- **Style Guide:** http://localhost:5173/style-guide

> **Nota:** Si el puerto 5173 está ocupado, Vite usará automáticamente 5174. Verifica la URL en la terminal donde corriste `npm run dev`.

### API Backend (Go)

- **Health Check:** http://localhost:8080/api/health
- **Login:** http://localhost:8080/auth/login
- **Logout:** http://localhost:8080/auth/logout
- **Refresh Token:** http://localhost:8080/auth/refresh
- **Swagger/Docs:** (Por implementar)

---

## 🐛 Solución de Problemas

### "Credenciales inválidas" en el navegador

**Causa:** El frontend no está conectándose correctamente a la API.

**Solución:**
1. Verifica que el backend esté corriendo: `docker ps`
2. **Verifica en qué puerto está el frontend:**
   ```bash
   # Windows
   netstat -ano | findstr "5173 5174"
   ```
   Busca la línea que diga `LISTENING` y anota el puerto
3. **Abre el navegador en el puerto correcto:**
   - Si está en 5173: http://localhost:5173/login
   - Si está en 5174: http://localhost:5174/login
4. **Si aún no funciona, reinicia el frontend:**
   ```bash
   # Detener el proceso de Vite (Ctrl+C en la terminal)
   cd apps/frontend
   npm run dev
   ```
   Observa el mensaje que dice `Local: http://localhost:XXXX/` y usa ese puerto

### El usuario no existe en la base de datos

```bash
# Verificar usuario
docker exec tramatex_db psql -U tramatex -d tramatex -c \
  "SELECT id, email, role FROM users WHERE email = 'admin@tramatex.local';"
```

**Si no hay resultados, insertar manualmente:**
```bash
docker exec tramatex_db psql -U tramatex -d tramatex -c \
  "INSERT INTO users (id, email, password, role, is_active) VALUES \
  ('f47ac10b-58cc-4372-a567-0e02b2c3d479', \
   'admin@tramatex.local', \
   '\$2a\$10\$Gd8.JP/L3j.vNvap81EpjuX7G4u5KKLmf10TSxmu779Mq/HdC/B9e', \
   'admin', \
   true) \
  ON CONFLICT (email) DO UPDATE SET password = EXCLUDED.password;"
```

### El frontend no se conecta a la API

**Verificar configuración de proxy en `apps/frontend/vite.config.js`:**

```javascript
server: {
  port: 5173,
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, ''),
    },
  },
}
```

### Error: "time: unknown unit "d" in duration"

**Causa:** Formato incorrecto de JWT TTL en docker-compose.

**Solución:** Editar `docker/docker-compose.yml`:
```yaml
JWT_REFRESH_TOKEN_TTL: 168h  # 7 días = 168 horas (no usar "7d")
```

Luego reiniciar:
```bash
docker-compose -f docker/docker-compose.yml restart api
```

---

## 📊 Verificar Estado del Sistema

### Backend

```bash
# Logs de la API
docker logs tramatex_api --tail 50

# Logs de la base de datos
docker logs tramatex_db --tail 50

# Estado de contenedores
docker ps

# Health check
curl http://localhost:8080/api/health
```

### Frontend

```bash
# Ver en qué puerto está corriendo (5173 o 5174)
netstat -ano | findstr "5173 5174"

# Verificar que responde (ajusta el puerto según corresponda)
curl -I http://localhost:5173
# o
curl -I http://localhost:5174
```

---

## 🔧 Reiniciar Todo (Solución Definitiva)

Si nada funciona, reinicia todo desde cero:

```bash
# 1. Detener todo
docker-compose -f docker/docker-compose.yml down -v
# Ctrl+C en la terminal del frontend

# 2. Reiniciar backend
docker-compose -f docker/docker-compose.yml up -d

# 3. Esperar 10 segundos para que inicie
timeout /t 10

# 4. Verificar que funciona
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@tramatex.local","password":"admin123"}'

# 5. Iniciar frontend
cd apps/frontend
npm run dev

# 6. Observar en qué puerto inició (verás un mensaje como "Local: http://localhost:5173/")
# 7. Abrir navegador en la URL que indica (5173 o 5174)
```

---

## 📚 Documentación Adicional

- **Credenciales completas:** `TEST_CREDENTIALS.md`
- **Credenciales completas:** `docs/guides/developer/test-credentials.md`
- **Generar nuevos hashes:** `apps/tramatex-api/cmd/hashgen/`
- **Arquitectura del proyecto:** `docs/architecture/`
- **Configuración impresión Sales (perfil fiscal emisor):** `apps/frontend/README.md` (sección "Configuración de emisor para impresión (Sales)")
- **Sistema de agentes:** `AGENTS.md`
- **Sprint actual:** `docs/log/sprints/sprint-09/`

---

## ✅ Checklist de Inicio

- [ ] Backend corriendo (docker ps muestra tramatex_api y tramatex_db)
- [ ] Frontend corriendo (verificar puerto con `netstat -ano | findstr "5173 5174"`)
- [ ] Navegador abierto en el puerto correcto (5173 o 5174)
- [ ] Usuario admin existe en BD
- [ ] Login desde API funciona (curl http://localhost:8080/auth/login)
- [ ] Login desde navegador funciona con admin@tramatex.local / admin123
- [ ] Puedo navegar a /products

### 🎯 Truco Rápido
Si ves "Credenciales inválidas" en el navegador:
1. Abre la consola del navegador (F12 → Network)
2. Intenta hacer login
3. Mira la petición a `/api/auth/login`
4. Si ves un error 404 o no hay petición, el puerto está mal
5. Verifica la URL en la barra del navegador (debe ser 5173 o 5174)

---

## 🆘 Contacto de Emergencia

Si después de seguir todos los pasos aún no funciona:

1. **Revisa los logs:** `docker logs tramatex_api --tail 100`
2. **Verifica la base de datos:** Ejecuta las queries de verificación arriba
3. **Revisa la consola del navegador:** F12 → Console → Network
4. **Comparte el error específico** que ves

---

**Última verificación exitosa:** 2026-02-09 17:56 UTC
**Hash verificado con:** `bcrypt.CompareHashAndPassword()` ✅
**API probada con:** `curl` ✅
**Estado:** Funcionando correctamente