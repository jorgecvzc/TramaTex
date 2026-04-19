# 🛠️ Índice Maestro de Scripts y Utilidades

| Metadato | Valor |
| :--- | :--- |
| **Versión** | 1.0 |
| **Estado** | ✅ Vigente |
| **Última Actualización** | 19-04-2026 |

---

## 🎯 Propósito
Este documento centraliza y explica todas las herramientas de automatización, scripts de mantenimiento y utilidades de soporte del ecosistema TramaTex. El objetivo es facilitar la operativa diaria y garantizar que ninguna utilidad técnica quede como un archivo desconocido.

---

## 🏁 Gestión del Entorno de Desarrollo (PowerShell)
Scripts principales para controlar el ciclo de vida de los contenedores Docker en Windows.

*   **`start-dev.ps1`**: Orquestador de arranque.
    *   *Uso común:* `.\start-dev.ps1` (Arranca DB y API).
    *   *Uso completo:* `.\start-dev.ps1 -Full` (Incluye Frontend y Nginx).
*   **`stop-dev.ps1`**: Detiene los servicios.
    *   *Limpieza profunda:* `.\stop-dev.ps1 -FullCleanup` (Borra volúmenes y orphans).
*   **`rebuild-dev.ps1`**: Reconstrucción total desde cero (Clean Slate). Útil ante cambios estructurales o limpieza de caché de Docker.

---

## 🧪 Utilidades de Datos y Validación (Go)
Herramientas para interactuar con el sistema de forma directa mediante el lenguaje del backend.

*   **`check_db.go`**: Utilidad de diagnóstico para verificar la conectividad y el estado de las tablas.
*   **`apply_migration.go`**: Permite ejecutar manualmente archivos SQL de migración.
*   **`restore_sales_db.go`**: Restaura específicamente los datos del módulo Sales al estado inicial de demo.
*   **`validate_sales_flow.go`**: **Script de Validación Crítico**. Simula un flujo de venta completo para asegurar la integridad del sistema.
*   **`translate_statuses.go`**: Herramienta de migración para normalizar estados en la base de datos.

---

## 🚀 Infraestructura y Despliegue Remoto (`/scripts`)
Automatización enfocada en la producción y entornos de `staging`.

*   **`install.sh`**: Script de instalación para servidores basados en Linux (Debian/Ubuntu).
*   **`rebuild-staging-remote.ps1`**: Orquestador SSH para actualizar automáticamente el servidor de pruebas.
*   **`verify-connectivity.ps1`**: Comprueba la accesibilidad entre los componentes del stack.
*   **`generate_presentation.py`**: Utilidad Python para generar la documentación de presentación corporativa.

---

## 🛡️ Seguridad y Calidad
*   **`protect-trunk-branches.ps1`**: Gancho preventivo para proteger las ramas maestras de commits directos.
*   **`setup-ssh-keys.ps1`**: Configuración automatizada de identidades para acceso remoto seguro.

---
[Volver al README Principal](../../../README.md)
