# Módulo de Product (Catálogo de Productos)

## 1. Propósito

*   **Visión del Módulo:** Gestionar el catálogo de productos con variantes y clasificaciones.
*   **Objetivos Clave:**
    *   Proporcionar un sistema para gestionar el catálogo de productos.
    *   Soportar variantes de productos como tallas, colores y modificaciones personalizadas.

## 2. Requisitos

### 2.1. Requisitos Funcionales

*   **RF-001:** Crear y mantener productos base.
*   **RF-002:** Gestionar variantes de productos.
*   **RF-003:** Asignar SKU único a cada producto y variante.
*   **RF-004:** Categorización de productos.
*   **RF-005:** Gestionar disponibilidad de stock (básico).

## 3. Casos de Uso

### 3.1. Actores

*   **Gerente de Producto:** Define y gestiona el catálogo de productos.
*   **Vendedor:** Consulta el catálogo para realizar ventas.

### 3.2. Casos de Uso Principales

*   **CU-001: CreateProduct**
    *   **Actor:** Gerente de Producto
    *   **Descripción:** Crear un producto base en el catálogo.
*   **CU-002: CreateVariant**
    *   **Actor:** Gerente de Producto
    *   **Descripción:** Agregar una nueva variante (talla, color) a un producto existente.
*   **CU-003: UpdateProduct**
    *   **Actor:** Gerente de Producto
    *   **Descripción:** Actualizar la información de un producto.
*   **CU-004: GetProduct**
    *   **Actor:** Vendedor / Gerente de Producto
    *   **Descripción:** Obtener los detalles de un producto, incluyendo todas sus variantes.
*   **CU-005: ListProducts**
    *   **Actor:** Vendedor / Gerente de Producto
    *   **Descripción:** Listar el catálogo de productos con filtros.
*   **CU-006: ChangeStatus**
    *   **Actor:** Gerente de Producto
    *   **Descripción:** Activar o descontinuar un producto.

## 4. Historias de Usuario

*   **HU-001:** Como Gerente de Producto, quiero poder crear un nuevo producto con su nombre, SKU y descripción para agregarlo al catálogo.
*   **HU-002:** Como Gerente de Producto, quiero agregar variantes de talla y color a un producto para ofrecer más opciones a los clientes.
*   **HU-003:** Como Vendedor, quiero ver todas las variantes de un producto en una sola pantalla para facilitar la cotización.

## 5. Criterios de Aceptación

*   **Para HU-002:**
    *   **Criterio 1:** Dado un producto existente, cuando agrego una variante con talla "M" y color "Azul", entonces la nueva variante aparece en la lista de variantes del producto.

## 6. Modelo de Dominio

### Product (Raíz de Agregación)
- **ID**: UUID
- **Nombre**: String
- **SKU**: String (único)
- **Descripción**: String
- **Categoría**: Enum (Textiles, Personalizables, etc.)
- **Estado**: Enum (Activo, Descontinuado)
- **Variantes**: List<ProductVariant>
- **Metadata**: CreatedAt, UpdatedAt

### ProductVariant
- **ID**: UUID
- **Tamaño**: String (S, M, L, XL, etc.)
- **Color**: String
- **Modificaciones**: String (personalizaciones)
- **SKU**: String (único, derivado)
- **CostoUnitario**: Decimal

## 7. Decisiones de Diseño

*   **Relaciones con Otros Módulos:**
    *   **Pricing**: Cada variante tiene su precio según la `Party` y el volumen.
    *   **Sales**: Las órdenes de venta contienen `ProductVariants`.
    *   **Inventory (futuro):** Se gestionará el stock de cada `ProductVariant`.
*   **Fases de Desarrollo:**
    *   [X] Fase 1 (MVP): Producto base + variantes básicas.
    *   [ ] Fase 2: Categorización y búsqueda avanzada.
    *   [ ] Fase 3: Historial de cambios del producto.