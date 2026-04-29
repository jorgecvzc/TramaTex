# Estrategia Técnica: Caché y Rendimiento (Post-MVP)

Este documento define la arquitectura de optimización para el ERP TramaTex, centrada en reducir la latencia en las operaciones más críticas: el motor de precios y la visualización del catálogo de productos.

---

## 1. Alcance de la Optimización

La estrategia se centra exclusivamente en datos de alta lectura y alto coste computacional:
- **Pricing**: Resultados de cálculos de precios basados en reglas complejas (clientes, grupos, tarifas).
- **Catálogo de Productos**: Proyecciones hidratadas (Producto + Variantes + Atributos) necesarias para la navegación rápida.

---

## 2. Estrategia de Gestión de Caché (Redis)

Se adoptará el patrón **Cache-Aside** con un enfoque de **Invalidación por Borrado** para garantizar la consistencia absoluta de los datos.

### 2.1 Flujo de Lectura (Lazy Loading)
1. El sistema recibe una petición de precio o producto.
2. Consulta en **Redis** mediante una llave única (ej: `price:{variant_id}:{party_id}:{quantity}`).
3. **Cache Hit**: Devuelve el dato instantáneamente.
4. **Cache Miss**: Consulta a PostgreSQL, realiza los cálculos, guarda el resultado en Redis y lo devuelve al usuario.

### 2.2 Estrategia de Invalidación (Consistencia Total)
Para evitar datos obsoletos, se implementará la **Invalidación Activa basada en Eventos**:
- **Acción**: Ante cualquier cambio en los datos origen (ej: edición de un precio base, cambio de atributos de producto o actualización de una regla de tarifa), el sistema **borra inmediatamente** todas las llaves de Redis relacionadas.
- **Resultado**: El siguiente acceso provocará un *Cache Miss*, forzando la recarga de datos frescos desde la base de datos maestra.

---

## 3. Implementación Técnica

### 3.1 Estructuras de Datos en Redis
- **Hashes/Strings**: Para resultados de precios calculados.
- **JSON**: Para objetos de producto hidratados, permitiendo una recuperación rápida de toda la ficha técnica en una sola operación de red.

### 3.2 Tiempos de Vida (TTL)
Aunque la invalidación es activa, se establecerá un **TTL (Time To Live)** de cortesía (ej: 24 horas) para limpiar datos de productos o precios que dejen de consultarse habitualmente, optimizando el uso de memoria en Redis.

---

## 4. UX: Rendimiento Percibido

La disponibilidad de datos en caché permite mejorar la fluidez de la interfaz:
- **Pre-Hydration**: El frontend puede mostrar datos cacheados de forma instantánea mientras valida en segundo plano si ha habido cambios.
- **Filtros Instantáneos**: Las búsquedas y filtrados en el catálogo se realizan sobre los objetos JSON en memoria, eliminando el "lag" de carga entre páginas.

---

*Última actualización: 2026-04-27*
