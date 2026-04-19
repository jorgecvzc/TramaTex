# 🌿 Guía de Control de Versiones (Git)

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.1 |
| **Estado** | ✅ Vigente |
| **Última Actualización** | 19-04-2026 |
| **Referencia** | [ADR-021](../../architecture/adrs/adr-021-version-control-and-branching-strategy.md) |

---

## 🎯 Propósito
Establecer un flujo de trabajo predecible y estandarizado para la gestión del código fuente, asegurando la estabilidad de la producción y la agilidad en el desarrollo del ecosistema TramaTex.

---

## 📍 Estrategia de Ramas (Branching Model)

### Ramas de Larga Duración
*   **`master`**: Refleja el código estable y desplegado en producción. No se permite el desarrollo directo. Solo recibe merges de `develop` (vía `release/`) o de `hotfix/*`.
*   **`develop`**: Rama principal de integración. Contiene el último código de desarrollo aprobado para la próxima entrega.

### Ramas Temporales y Nomenclatura
Usamos el formato: `<tipo>/<descripcion-kebab-case>`. Las descripciones deben ser claras, en minúsculas y sin superar los 50 caracteres.

| Tipo | Propósito | Origen | Destino PR |
| :--- | :--- | :--- | :--- |
| **`feature/`** | Nuevas funcionalidades o mejoras funcionales. | `develop` | `develop` |
| **`bugfix/`** | Corrección de errores detectados en desarrollo. | `develop` | `develop` |
| **`hotfix/`** | Emergencias críticas en producción. | `master` | `master` |
| **`infra/`** | Cambios en Docker, CI/CD, Nginx o scripts. | `develop` | `develop` |
| **`release/`** | Estabilización y preparación de versión. | `develop` | `master` |
| **`docs/`** | Cambios exclusivos de documentación. | `develop` | `develop` |

---

## 🔄 Flujo de Trabajo Estándar

### 1. Iniciar una Tarea
```bash
git checkout develop
git pull origin develop
git checkout -b feature/nombre-tarea
```

### 2. Commits Convencionales
Todo commit debe seguir el estándar **Conventional Commits**: `tipo(ambito): descripción`.
*   `feat(sales): implement order cancellation logic`
*   `fix(party): resolve NIF validation edge case`
*   `refactor(domain): improve money object performance`
*   `docs(guides): update deployment steps`

### 3. Integración y Limpieza
Tras completar el desarrollo y asegurar que todos los tests pasan:
1. Subir rama: `git push origin feature/nombre-tarea`.
2. Abrir Pull Request (PR) hacia el destino correspondiente.
3. Tras la revisión y aprobación, realizar el merge (preferiblemente Squash o Merge Commit según política de repositorio).
4. Borrar la rama local y remota.

---

## 🔐 Gestión de Variables de Entorno (.env)
**Regla de Seguridad:** Los archivos `.env` reales contienen secretos y **NUNCA** deben subirse al repositorio.
*   Utiliza los archivos `.env.example` distribuidos en las carpetas `apps/` y `docker/` como base.
*   Cualquier nueva variable de entorno debe añadirse primero al archivo `.env.example` correspondiente.

---

## 🏷️ Versionado Semántico (SemVer)
TramaTex utiliza el estándar `MAJOR.MINOR.PATCH` para etiquetar sus entregas:
*   **MAJOR**: Cambios arquitectónicos o de API que rompen la compatibilidad.
*   **MINOR**: Nuevas funcionalidades o módulos añadidos de forma compatible.
*   **PATCH**: Corrección de errores, mejoras de rendimiento o actualizaciones menores.

---
[Volver al README Principal](../../../README.md)
