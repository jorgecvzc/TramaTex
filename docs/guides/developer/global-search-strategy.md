# Estrategia de Búsqueda Global: TramaTex Productivity

Esta guía define el funcionamiento del motor de búsqueda unificado (Ctrl+K). Se considera una **ampliación estratégica de la UI/UX** para maximizar la velocidad operativa de los usuarios.

---

## 🎯 Visión y Objetivos
El objetivo es permitir que cualquier usuario salte a cualquier punto de gestión en menos de 2 segundos, sin necesidad de navegar por los menús.

### Principios:
1.  **Omnipresencia**: Accesible desde cualquier lugar mediante `Ctrl+K`.
2.  **Multidominio**: Busca en Ventas, Producción, Catálogo y Entidades simultáneamente.
3.  **Prioridad al Cliente**: La búsqueda por nombre de cliente debe filtrar todos los documentos asociados.

---

## 🔍 Alcance de Búsqueda (Criterios de Éxito)

| Módulo | Entidad | Criterios de Búsqueda | Subtítulo Informativo |
| :--- | :--- | :--- | :--- |
| **Ventas** | Pedidos | Nº Pedido / Nombre Cliente | `[ESTADO] · Cliente: Nombre` |
| **Ventas** | Presupuestos | Nº Presupuesto / Nombre Cliente | `[ESTADO] · Cliente: Nombre` |
| **Ventas** | Facturas | Nº Factura / Nombre Cliente | `[ESTADO] · Total: 0,00€` |
| **Ventas** | Albaranes | Nº Albarán / Nombre Cliente | `[ESTADO] · Fecha: dd/mm/aa` |
| **Producción** | Órdenes (MES) | Nº Orden / Nombre Cliente (Origen) | `Trabajo para: Cliente` |
| **Catálogo** | Productos | SKU / Nombre del Producto | `SKU: XXX · Precio: 0,00€` |
| **Entidades** | Clientes/Prov. | Nombre / NIF / CIF | `NIF: XXX · Localidad` |

---

## 🛠️ Arquitectura Técnica

### Backend (Go)
*   **Endpoint Unificado**: `GET /api/search?q=query`.
*   **Autenticación**: Endpoint protegido. Requiere sesión autenticada/JWT válido.
*   **Concurrencia**: Uso de `goroutines` para consultar cada tabla en paralelo.
*   **Ranking**: Coincidencias exactas en números de documento o SKU tienen prioridad máxima.
*   **Limitación**: Máximo 5 resultados por categoría para mantener la limpieza.

### Frontend (Vue)
*   **Single Request**: Una sola llamada a la API por término de búsqueda.
*   **Categorización**: Resultados agrupados visualmente por módulo.
*   **Acción**: Navegación inmediata al hacer clic o pulsar `Enter`.

---

## 📏 Estándar de Visualización
Cada resultado debe seguir este patrón en el buscador:
*   **Título**: Identificador único (ej: `Pedido V-2026-001`).
*   **Subtítulo**: Contexto rico (ej: `[CONFIRMADO] · Cliente: Textil S.A.`).
*   **Icono**: El icono oficial de la familia UI correspondiente.
