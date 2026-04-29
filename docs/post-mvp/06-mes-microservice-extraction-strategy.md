# Estrategia Técnica: Extracción del MES como Microservicio (Post-MVP)

Este documento detalla el plan maestro para desacoplar el Módulo de Ejecución de Manufactura (MES) del monolito ERP, transformándolo en un microservicio autónomo. Esta evolución garantiza la resiliencia de la planta de producción frente a caídas del sistema comercial y permite el escalado independiente de los terminales de taller.

---

## 1. Arquitectura de Soberanía de Datos

El principio fundamental de esta extracción es la **Soberanía Estricta de Datos**: ningún servicio accederá directamente a la base de datos de otro.

### 1.1 El MES como "Single Source of Truth" (SSOT)
- **Propiedad**: El microservicio MES es el único dueño de las Entidades de Fabricación: Órdenes de Trabajo (`WorkOrders`), Fases de Producción (`WorkPhases`), Registros de Tiempo (`TimeEntries`) y Estados de Máquina.
- **Acceso del ERP**: Si el módulo de Ventas (Sales) necesita conocer el progreso de fabricación de un Pedido, **deberá consultar al MES** a través de una API (gRPC) o consumir eventos de actualización de estado. No existirá ninguna tabla de "estado de fabricación" masterizada en el ERP.

### 1.2 Datos de Referencia (Proyección Local en MES)
El MES necesita conocer ciertos datos del ERP para contextualizar el trabajo en planta, pero bajo una política de **Mínima Información Necesaria**:
- **Persistencia**: Los datos se almacenarán en **tablas de proyección (solo lectura)** dentro de la propia base de datos PostgreSQL del microservicio MES.
- **Simplicidad**: Al ser datos estáticos (sin cálculos pesados), no se requiere una capa de caché adicional (Redis); la base de datos relacional local garantiza la integridad y la persistencia necesaria para operar offline.
- **Datos Sincronizados**: Identificadores (UUID), SKUs de productos, nombres de clientes (para etiquetas/identificación) y especificaciones técnicas.
- **Exclusión Financiera**: El MES **no almacena ni procesa precios, descuentos, impuestos ni importes**.


---

## 2. Patrones de Comunicación Híbrida

Para conectar el ERP central con el microservicio MES, se utilizará un enfoque dual:

### 2.1 Comunicación Síncrona (gRPC)
Uso reservado para consultas inmediatas donde el usuario (en el ERP) está esperando una respuesta bloqueante.
- **Ejemplo**: El comercial pulsa "Ver Detalle de Fabricación" en un Pedido. El ERP hace una llamada gRPC al MES: `GetManufacturingStatus(OrderID)`.

### 2.2 Comunicación Asíncrona (Message Broker - NATS JetStream)
El canal principal para el flujo de procesos y la coreografía de dominio, garantizando el desacoplamiento temporal.
- **Flujo de Inicio**: Sales publica el evento `SalesOrderConfirmed`. El MES está suscrito, recibe el evento y genera automáticamente las `WorkOrders` correspondientes.
- **Flujo de Fin**: El MES publica `WorkOrderCompleted`. Sales (o Inventario) lo consume para preparar el albarán de entrega o actualizar el stock.

---

## 3. Resiliencia Industrial (Local-First y Offline)

El taller no puede parar. La arquitectura del microservicio MES debe diseñarse asumiendo que la conexión con las oficinas centrales (ERP) es inestable o puede caer.

### 3.1 Operativa Aislada
- Las bases de datos locales del MES y las cachés de referencia permiten a los operarios seguir fichando tiempos, cambiando estados de máquinas y completando tareas **sin conexión al ERP**.
- Los eventos generados durante la caída (ej: `PhaseCompleted`) se almacenan en una cola transaccional local (Patrón *Outbox*).

### 3.2 Reconciliación Automática
- Al recuperar la conexión, el MES vuelca todos los eventos acumulados al Message Broker central (NATS), sincronizando el estado con el ERP de forma transparente para el usuario.

---

## 4. Hoja de Ruta de Migración (Patrón Strangler Fig)

La extracción no se hará con un "Big Bang", sino de forma progresiva y segura:

1.  **Proxy de API**: Redirigir todas las llamadas del frontend relacionadas con fabricación a un nuevo Gateway. Inicialmente, el Gateway enruta al monolito.
2.  **Duplicación de Escritura (Shadowing)**: El monolito empieza a escribir datos de fabricación tanto en su base de datos actual como en la nueva base de datos del futuro microservicio.
3.  **Despliegue del Microservicio MES**: Se levanta el servicio independiente, que lee de la nueva base de datos. Se realizan pruebas de carga y validación cruzada.
4.  **Corte (Cut-over)**: El Gateway cambia el enrutamiento. Ahora todas las peticiones MES van al nuevo microservicio.
5.  **Limpieza (Cleanup)**: Se elimina el código del módulo MES y sus tablas del monolito original.

---

*Última actualización: 2026-04-27*
