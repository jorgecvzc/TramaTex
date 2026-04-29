# Estrategia Técnica: Gestión Avanzada de Archivos de Diseño en MES (Post-MVP)

Este documento define cómo TramaTex manejará los archivos técnicos y de diseño (ej: planos, fichas técnicas de tejido) para que el operario tenga toda la información necesaria en el terminal de taller.

---

## 1. Visualización Industrial

### 1.1 Thumbnails e Imágenes de Referencia
- **Generación Automática**: El sistema generará miniaturas de los archivos de diseño (PDF, PNG, JPG) al subirlos al ERP.
- **Visualización en OT**: En el terminal MES, la Orden de Trabajo mostrará de forma prominente la imagen del diseño para evitar errores de producción.

---

## 2. Integración con Aplicaciones Nativas

Para archivos que requieren software especializado (ej: CAD, Editores de Diseño):
- **Protocolo Personalizado (`tramatex://`)**: Implementación de un protocolo que permite, desde el navegador, lanzar la apertura del archivo en la aplicación nativa instalada en el equipo del taller.
- **Sincronización de Cambios**: El operario puede editar un diseño y, al guardar, el sistema detecta el cambio y propone subir la nueva versión como revisión.

---

## 3. Almacenamiento y Seguridad

### 3.1 Gestión de Versiones (Versioning)
- Cada cambio en un archivo de diseño genera una nueva versión. El MES siempre mostrará la versión aprobada para producción, evitando el uso de planos obsoletos.

### 3.2 Acceso Offline
- Los archivos necesarios para las Órdenes de Trabajo activas se cachearán localmente en los terminales del taller para garantizar el acceso incluso durante caídas de red (alineado con la estrategia del Microservicio MES).

---

*Última actualización: 2026-04-27*
