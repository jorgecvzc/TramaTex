# Estrategia Técnica: Evolución de Infraestructura (Post-MVP)

Este documento define la ruta de escalabilidad técnica para soportar el crecimiento de TramaTex hacia entornos de producción masiva y alta disponibilidad.

---

## 1. Orquestación y Despliegue (Kubernetes)

Migración del despliegue actual basado en Docker Compose hacia un clúster de **Kubernetes (K8s)**.

### 1.1 Ventajas Operativas
- **Auto-escalado**: Capacidad de aumentar el número de réplicas de los microservicios (especialmente el MES) ante picos de carga en el taller.
- **Self-healing**: Reinicio automático de contenedores en caso de fallo crítico.
- **Zero-Downtime Deployments**: Actualizaciones del sistema sin interrupción del servicio mediante *Rolling Updates*.

---

## 2. Observabilidad Centralizada

Implementación de un stack de monitorización profesional (Prometheus + Grafana + Loki).

### 2.1 Métricas y Alertas
- **Infraestructura**: Control de CPU, Memoria y Disco de los servidores.
- **Negocio**: Alertas en tiempo real si la tasa de errores en el motor de facturación o comunicación con la AEAT supera un umbral.
- **Taller**: Monitorización de la latencia en los terminales de operario para detectar problemas de red local.

---

## 3. Seguridad Avanzada e Infraestructura como Código (IaC)

### 3.1 Gestión de Secretos
- Migración hacia **HashiCorp Vault** o soluciones nativas de Cloud (AWS Secrets Manager / GCP Secret Manager) para la gestión segura de certificados de firma electrónica y claves de API.

### 3.2 Terraform / CloudFormation
- Definición de toda la infraestructura mediante código para asegurar entornos de Staging y Producción idénticos y reproducibles.

---

## 4. Gestión del Ciclo de Vida del Dato (Archiving)

Para mantener el rendimiento de PostgreSQL ante el crecimiento masivo de registros (especialmente en el MES):

### 4.1 Estrategia de Cold Storage
- **Datos Activos**: Documentos de los últimos 2 años y órdenes de trabajo abiertas permanecen en la base de datos de alto rendimiento (SSD).
- **Archivado Automático**: Los registros de más de 2 años se mueven a una instancia de base de datos de "solo lectura" o a archivos de almacenamiento de objetos (Parquet/S3) para análisis histórico.
- **Acceso Transparente**: Implementación de una capa de abstracción que permite consultar datos históricos desde el ERP si es necesario, aunque con una latencia mayor.

---

*Última actualización: 2026-04-27*
