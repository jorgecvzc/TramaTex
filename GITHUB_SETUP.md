# 🔧 Configuración de GitHub - Guía Rápida

## 📌 Contexto

Repositorio local listo para publicar. Necesitas configurar autenticación para:
- Push de commits al repositorio remoto
- Crear Pull Requests
- Configurar GitHub Actions (CI/CD)

---

## 🚀 OPCIÓN 1: Autenticación SSH (RECOMENDADA)

### 1. Verificar si ya tienes claves SSH

```powershell
ls ~/.ssh/
# Busca: id_rsa.pub, id_ed25519.pub, etc.
```

### 2. Si NO tienes, generar nueva clave

```powershell
# Generar clave ED25519 (más moderna y segura)
ssh-keygen -t ed25519 -C "tu-email@example.com"

# O RSA si tu sistema no soporta ED25519
ssh-keygen -t rsa -b 4096 -C "tu-email@example.com"

# Presiona Enter para ubicación por defecto
# Opcionalmente añade passphrase para mayor seguridad
```

### 3. Copiar clave pública

```powershell
# ED25519
cat ~/.ssh/id_ed25519.pub | clip

# O RSA
cat ~/.ssh/id_rsa.pub | clip
```

### 4. Añadir clave a GitHub

1. Ve a GitHub.com → Settings (tu perfil)
2. Sidebar: **SSH and GPG keys**
3. Click **New SSH key**
4. Title: "TramaTex Dev - [Nombre PC]"
5. Key: Pega lo que copiaste (Ctrl+V)
6. Click **Add SSH key**

### 5. Probar conexión

```powershell
ssh -T git@github.com
# Debe decir: "Hi [username]! You've successfully authenticated..."
```

### 6. Configurar remote SSH

```powershell
# Verificar remote actual
git remote -v

# Si es HTTPS, cambiarlo a SSH
git remote set-url origin git@github.com:TU-USUARIO/TramaTex.git

# Verificar
git remote -v
```

---

## 🔑 OPCIÓN 2: Personal Access Token (PAT)

### 1. Crear Token en GitHub

1. Ve a GitHub.com → Settings
2. Sidebar: **Developer settings** → **Personal access tokens** → **Tokens (classic)**
3. Click **Generate new token (classic)**
4. Note: "TramaTex Development"
5. Expiration: 90 días (o sin expiración si confías en tu máquina)
6. Scopes a marcar:
   - ✅ `repo` (Full control of private repositories)
   - ✅ `workflow` (Update GitHub Actions workflows)
   - ✅ `write:packages` (si usas GitHub Packages)
7. Click **Generate token**
8. **¡COPIA EL TOKEN AHORA!** (No podrás verlo de nuevo)

### 2. Usar Token para Push

```powershell
# Método 1: Cachear credenciales (recomendado)
git config --global credential.helper wincred

# Primer push te pedirá credenciales
git push origin main
# Username: tu-usuario-github
# Password: pega-tu-PAT-aquí

# Windows guardará el PAT para futuros push
```

```powershell
# Método 2: URL con token (menos seguro, visible en historial)
git remote set-url origin https://TU-PAT@github.com/TU-USUARIO/TramaTex.git
```

---

## 📤 Primer Push al Repositorio

```powershell
# 1. Verificar rama actual
git branch
# Debe mostrar: * main (o master)

# 2. Si no existe remote, añadirlo
git remote add origin git@github.com:TU-USUARIO/TramaTex.git
# O con HTTPS:
# git remote add origin https://github.com/TU-USUARIO/TramaTex.git

# 3. Verificar remote
git remote -v

# 4. Push inicial
git push -u origin main

# 5. Verificar en GitHub que aparece el código
```

---

## 🔄 Configurar GitHub Actions (Después del Push)

Una vez el código esté en GitHub:

### 1. Habilitar GitHub Actions

1. Ve a tu repositorio en GitHub
2. Tab: **Actions**
3. Si pregunta, click **I understand, enable Actions**

### 2. Crear Workflow para CI/CD (Sprint 04, Tarea 02)

Esto lo harás en el Sprint 04-02, pero el archivo base será:

```yaml
# .github/workflows/backend.yml
name: Backend CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.23'
      - name: Run tests
        run: cd apps/tramatex-api && go test ./...
```

### 3. Configurar Branch Protection (OPCIONAL pero recomendado)

1. Repositorio → Settings → Branches
2. Click **Add rule**
3. Branch name pattern: `main`
4. Marcar:
   - ✅ Require a pull request before merging
   - ✅ Require status checks to pass before merging
   - ✅ Require branches to be up to date before merging

---

## ✅ Checklist Final

Antes de empezar Sprint 04:

- [ ] Autenticación GitHub configurada (SSH o PAT)
- [ ] Remote origin configurado correctamente
- [ ] `git push` funciona sin errores
- [ ] Código visible en GitHub.com
- [ ] (Opcional) GitHub Actions habilitado
- [ ] (Opcional) Branch protection configurado

---

## 🆘 Troubleshooting

### Error: "Permission denied (publickey)"
**Causa:** Clave SSH no configurada o no reconocida
**Solución:**
```powershell
# Iniciar ssh-agent
ssh-agent bash
ssh-add ~/.ssh/id_ed25519

# Probar de nuevo
ssh -T git@github.com
```

### Error: "Authentication failed" con PAT
**Causa:** Token inválido o sin permisos
**Solución:**
- Verifica que copiaste el token completo
- Verifica que marcaste scope `repo`
- Genera un nuevo token si es necesario

### Error: "Repository not found"
**Causa:** URL del remote incorrecta o repositorio no existe
**Solución:**
```powershell
# Verificar
git remote -v

# Corregir
git remote set-url origin git@github.com:USUARIO-CORRECTO/TramaTex.git
```

---

## 📞 Recursos Adiciales

- [GitHub SSH Docs](https://docs.github.com/en/authentication/connecting-to-github-with-ssh)
- [GitHub PAT Docs](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)
- [GitHub Actions Docs](https://docs.github.com/en/actions)

---

**Una vez configurado, estarás listo para Sprint 04** 🚀
