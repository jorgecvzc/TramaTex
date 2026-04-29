# Estrategia Técnica: Comunicación Asíncrona y Eventos (Post-MVP)

Este documento define la columna vertebral de comunicación de TramaTex. La transición hacia una arquitectura dirigida por eventos (Event-Driven) permite desacoplar los módulos, mejorar la escalabilidad y garantizar la integridad de los procesos de negocio entre el ERP y el MES.

---

## 1. Infraestructura de Mensajería: NATS JetStream

Se ha seleccionado **NATS JetStream** como el motor de mensajería debido a su ligereza, alto rendimiento y capacidades de persistencia nativa.

### 1.1 Características Clave
- **Persistencia (Streams)**: Los mensajes no se pierden si un consumidor está offline; se almacenan en el "Stream" hasta que son procesados exitosamente.
- **Acuse de Recibo (ACK)**: Garantiza que un evento de venta llegue al MES. Si el MES no confirma el procesamiento, NATS reintenta el envío.
- **Despliegue Simple**: Integrado en el stack de Docker del proyecto sin dependencias pesadas (a diferencia de Kafka).

---

## 2. Integridad: Patrón Transactional Outbox

Para evitar el problema de "Base de Datos actualizada pero Mensaje no enviado", implementaremos el patrón **Transactional Outbox**.

### 2.1 Flujo de Ejecución
1.  **Transacción Única**: Dentro de la misma transacción de PostgreSQL donde se guarda la entidad (ej: Pedido de Venta), se inserta un registro en una tabla técnica llamada `outbox`.
2.  **Relay Service**: Un servicio ligero en segundo plano (Worker) lee de la tabla `outbox`, publica el mensaje en NATS y, tras recibir el OK de NATS, marca el registro como enviado en la DB.
3.  **Resultado**: Garantía absoluta de que si el Pedido existe en la base de datos, el evento llegará al sistema de mensajería.

---

## 3. Coreografía de Eventos de Dominio

En lugar de que el módulo A llame al módulo B, los módulos "gritan" lo que ha pasado y otros "escuchan" si les interesa.

### 3.1 Eventos Principales
- **`SalesOrder.Confirmed`**:
    - **Publicador**: Sales.
    - **Suscriptor**: MES (Crea Órdenes de Trabajo) e Inventario (Reserva Stock).
- **`Manufacturing.WorkOrderFinished`**:
    - **Publicador**: MES.
    - **Suscriptor**: Sales (Marca pedido para facturar) y Logística (Prepara envío).
- **`Product.Updated`**:
    - **Publicador**: Product.
    - **Suscriptor**: MES y Sales (Actualizan sus proyecciones locales).

---

## 4. Consistencia Eventual y Resiliencia

### 4.1 Manejo de Errores (Dead Letter Queues)
Si un suscriptor falla repetidamente al procesar un mensaje (ej: error de lógica o base de datos bloqueada):
- El mensaje se mueve a una cola de **"Cartas Muertas" (DLQ)**.
- El sistema lanza una alerta técnica para intervención manual o reintento programado.

### 4.2 Idempotencia
Todos los consumidores deben ser **idempotentes**: recibir el mismo mensaje dos veces (ej: por un reintento de red) no debe duplicar el resultado (ej: no crear dos órdenes de trabajo para el mismo pedido).

---

## 5. UX: Visibilidad del Sistema

- **Monitor de Eventos**: Panel técnico en el área de administración para visualizar el tráfico de mensajes y detectar cuellos de botella o fallos de sincronización entre el ERP y el MES.

---

*Última actualización: 2026-04-27*
