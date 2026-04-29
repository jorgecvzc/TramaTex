# Estrategia Técnica: Búsqueda Global Avanzada (Post-MVP)

Este documento define la evolución del sistema de búsqueda de TramaTex hacia un motor unificado y de alta velocidad, centralizado en el atajo `Ctrl + K`.

---

## 1. El Concepto "Command Center" (`Ctrl + K`)

La búsqueda dejará de estar fragmentada por módulos. Un único punto de entrada permitirá localizar cualquier entidad del sistema mediante lenguaje natural o códigos.

### 1.1 Entidades Indexadas
- **Comercial**: Clientes (Nombre, CIF), Facturas (Nº Serie), Pedidos, Presupuestos.
- **Catálogo**: Productos (Nombre, SKU, Atributos técnicos).
- **Producción**: Órdenes de Trabajo (ID), Operarios, Máquinas.
- **Navegación**: "Ir a Ventas", "Nuevo Producto", "Ayuda".

---

## 2. Motor de Búsqueda y Rendimiento

### 2.1 Full-Text Search (PostgreSQL + Trigramas)
Para evitar la complejidad de un servidor ElasticSearch dedicado en el MVP/Post-MVP temprano, se potenciará la capacidad nativa de PostgreSQL:
- **Indices GIN/GIST**: Uso de extensiones `pg_trgm` para búsquedas "borrosas" (Fuzzy Search) que toleren errores tipográficos.
- **Rankeo de Relevancia**: Los resultados se ordenarán por importancia (ej: una coincidencia exacta en SKU tiene más peso que una coincidencia parcial en descripción).

### 2.2 Autocompletado Instantáneo
- Al escribir las primeras 3 letras, se mostrará un desplegable con previsualización de resultados.
- **Ejecución Local**: Los comandos de navegación se cachean en el frontend para respuesta inmediata (0ms de latencia).

---

## 3. UX y Ergonomía

### 3.1 Previsualización Rápida
- Al navegar por los resultados con las flechas del teclado, un panel lateral mostrará un resumen de la entidad (ej: si es un cliente, mostrar su deuda actual y teléfono de contacto).

### 3.2 Acciones Contextuales
- Cada resultado permitirá ejecutar acciones directas:
    - Cliente -> `Enter` para detalle, `Alt+N` para nuevo pedido.
    - Factura -> `Enter` para PDF, `Alt+P` para registrar cobro.

---

*Última actualización: 2026-04-27*
