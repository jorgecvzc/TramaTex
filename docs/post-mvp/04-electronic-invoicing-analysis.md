# Estrategia Técnica: Facturación Electrónica y Veri*factu (Post-MVP)

Este documento detalla la adaptación integral de TramaTex a la normativa legal española (Ley Crea y Crece, Ley Antifraude y Veri*factu). El objetivo es automatizar la generación, firma y comunicación de facturas garantizando el cumplimiento absoluto para evitar sanciones por "software de doble uso".

---

## 1. Integridad y Encadenamiento (Veri*factu)

### 1.1 Encadenamiento de Registros (Hash Chaining)
Cada factura generada genera un registro XML de alta o anulación inalterable:
- **Hash de Vínculo**: Cada registro incluye el hash del registro anterior, la fecha y hora exacta (sincronizada vía **NTP**) y el número de serie del software.
- **Anulaciones**: No se permite el borrado de facturas. Las anulaciones generan un nuevo registro que referencia al original, manteniendo la integridad de la serie.

### 1.2 Registro de Eventos y Auditoría
- **Log de Eventos (Event Log)**: Registro inalterable de accesos, intentos de envío, errores de comunicación y cambios de configuración técnica.
- **Declaración Responsable**: El sistema incluirá una sección de "Certificación Legal" accesible al usuario que contiene la declaración responsable del fabricante, obligatoria por ley para que el usuario pueda facturar.

---

## 2. Firma Electrónica y Resiliencia

### 2.1 Política de Firma (XAdES)
- **Modo Veri*factu (Remisión)**: Envío directo a la AEAT. La firma es opcional pero recomendada para mayor seguridad.
- **Modo No Veri*factu (Conservación)**: Firma electrónica **obligatoria** de cada registro local mediante certificado de sello de entidad.

### 2.2 Gestión de la Continuidad (Modo Offline)
- El software debe permitir seguir operando si los servicios de la AEAT no están disponibles.
- Los registros se encadenan localmente y se marcan para **envío automático diferido** una vez se restablezca la conexión.

---

## 3. Ley Crea y Crece (Ciclo B2B)

### 3.1 Formatos Estructurados e Interoperabilidad
- Generación nativa de **FacturaE (XML)** y **UBL** para el intercambio entre empresas.
- Integración del **Código QR** en todos los PDF, permitiendo la verificación instantánea por parte del receptor.

### 3.2 Reporte de Estados de Pago
- Obligatoriedad de informar sobre la **Fecha de Pago Efectivo** para el control de la morosidad comercial (periodos medios de pago).
- Recepción automática de facturas de proveedores para cierre de ciclo financiero.

---

## 4. UX y Operativa Industrial

### 4.1 Panel de Cumplimiento Fiscal
- **Monitor Veri*factu**: Estado en tiempo real de la comunicación con la AEAT.
- **Buzón de Facturas Electrónicas**: Portal donde el cliente puede descargar sus facturas firmadas durante los 4 años obligatorios de custodia.

### 4.2 Alertas Críticas
- Aviso inmediato de desincronización horaria (NTP) que impida el encadenamiento correcto.
- Notificación de facturas rechazadas por la plataforma receptora con sugerencia de corrección técnica.

---

## 5. Especificaciones de Infraestructura

- **Custodia**: Almacenamiento seguro (WORM - Write Once Read Many) de los XML durante 6 años.
- **Sincronización NTP**: Configuración del servidor para consulta periódica a servidores de hora oficiales.
- **Web Services**: Implementación de clientes SOAP/REST compatibles con la sede electrónica de la AEAT.

---

*Última actualización: 2026-04-27 (Blindaje Legal Avanzado)*
