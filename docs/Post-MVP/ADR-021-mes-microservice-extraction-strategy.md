# ADR-021: Estrategia de Extracción del Microservicio MES (Post-MVP)

## Estado
Propuesto (Estudio Post-MVP)

## Contexto
El módulo MES (Manufacturing Execution System) de TramaTex está actualmente integrado dentro del monolito modular para simplificar el desarrollo del MVP. Sin embargo, el negocio requiere que este módulo esté arquitectónicamente preparado para ser extraído como un microservicio independiente. Esto permitirá un escalado independiente, ciclos de despliegue especializados e integración con hardware físico de fábrica que puede residir en un segmento de red diferente.

## Decisión
Adoptaremos una arquitectura de **"Shared-Nothing"** (Nada Compartido) para la extracción de MES, pasando de llamadas a funciones en memoria a patrones de comunicación distribuida.

### 1. Patrones de Comunicación

#### A. Síncrona (Consultas) - gRPC
Para necesidades de datos en tiempo real (por ejemplo, verificar si un producto existe antes de comenzar una tarea), MES se comunicará con el ERP Core a través de **gRPC**.
- **Por qué:** Alto rendimiento, contratos estrictamente tipados mediante Protocol Buffers y excelente soporte en Go.
- **Mecanismo:** Las interfaces de servicios de aplicación actuales en MES serán satisfechas por nuevos adaptadores de cliente remoto que implementarán la misma interfaz pero realizarán llamadas gRPC en lugar de búsquedas en repositorios locales.

#### B. Asíncrona (Comandos/Efectos Secundarios) - Mensajería Fiable
Para flujos de trabajo reactivos donde una solicitud debe procesarse eventualmente sin bloquear al emisor (por ejemplo, una orden de venta creando un trabajo MES), utilizaremos **Mensajería Asíncrona Fiable** con **NATS JetStream**.

- **Tecnología: NATS JetStream.**
  - **Por qué:** Nativo de la nube (optimizado para Kubernetes), escrito en Go, proporciona latencia ultra baja y soporta "Streams" persistentes para una entrega sin pérdidas.
- **Patrón Transactional Outbox:** Para asegurar que no se pierdan mensajes en el origen, el Core guardará el mensaje en una tabla local `outbox` dentro de la misma transacción de base de datos que la entidad de negocio. Un proceso "Relay" en segundo plano los enviará a NATS.
- **Grupos de Consumidores:** MES utilizará "Pull Consumers" con suscripciones duraderas. Esto permite escalar el servicio MES a múltiples pods en Kubernetes; NATS distribuirá el trabajo automáticamente y asegurará que cada mensaje se procese exactamente una vez.
- **Confirmaciones y Reintentos:** NATS solo considerará un mensaje como entregado una vez que MES envíe un `Ack`. Si MES falla, NATS reenviará el mensaje según una política de reintentos.

#### C. Consultas de Información (Lectura) - Regla Híbrida 90/10
Para equilibrar la frescura de los datos con la resiliencia del sistema, adoptamos una estrategia híbrida:
- **90% - Proyecciones Locales (Vistas Materializadas vía Eventos NATS):** Los listados comunes se servirán desde modelos de lectura locales actualizados asíncronamente. Esto garantiza una respuesta instantánea de la UI y aislamiento total de fallos.
- **10% - Consultas Remotas Directas (gRPC):** Solo para datos de alta volatilidad o validaciones críticas en tiempo real (ej. telemetría de máquinas).

### 2. Estrategia de Datos
- **Esquema Independiente:** MES mantendrá su propia base de datos (instancia de PostgreSQL).
- **Referencias Lógicas:** Sin claves foráneas (FK) físicas a las tablas del Core. MES almacenará referencias `uuid`.
- **Replicación de Datos:** MES podrá mantener un "caché" local de datos esenciales del Core para seguir funcionando durante particiones de red.

### 3. Identidad y Seguridad
- **Propagación de JWT:** El token JWT del usuario se pasará en los metadatos de gRPC para asegurar la autorización en MES.

## Consecuencias

### Positivas
- **Escalado Independiente:** MES puede escalarse según la carga de la fábrica.
- **Aislamiento de Fallos:** Una caída en Ventas no detiene la producción en el taller.
- **Flexibilidad Tecnológica:** MES podría reescribirse en otro lenguaje si la integración con hardware lo requiere.

### Negativas
- **Complejidad Operativa:** Requiere gestionar un broker de mensajería.
- **Trazabilidad Distribuida:** Depurar es más complejo (requiere IDs de correlación).
- **Consistencia Eventual:** Se debe gestionar la sincronización de datos entre servicios.

## Implementation Roadmap (Conceptual)
1. **Phase 1:** Define `.proto` files for Core services.
2. **Phase 2:** Refactor MES Adapters to use the gRPC client instead of local Repo injection.
3. **Phase 3:** Introduce a Message Broker for `Sales` -> `MES` triggers.
4. **Phase 4:** Physical separation of the database.
