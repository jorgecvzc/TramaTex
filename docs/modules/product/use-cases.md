# Product Module: Use Cases

- **Version:** 1.2.0
- **Status:** Implementado

Este documento describe los casos de uso para el módulo de `product`, alineado con la implementación final.

---

## 1. Gestión de Atributos (`Attribute`)

### UC-P-001: Crear un Atributo

- **Actor:** Administrador del Sistema
- **Descripción:** Crea un nuevo atributo configurable (ej: "Talla", "Color") y sus posibles valores, definiendo su código para la composición de SKU y su alcance.
- **Flujo Básico:**
  1. El administrador proporciona:
     - `Name` (ej: "Tallas Estándar")
     - `Code` (ej: "T")
     - `Order` (posición en la composición del SKU)
     - Lista de `AttributeValue`s, cada uno con `Value` (ej: "Small") y `Code` (ej: "S").
  2. El administrador define el **alcance** de este atributo:
     - **Genérico:** No se especifica `BrandID` ni `GroupIDs`.
     - **Por Marca:** Se selecciona una `BrandID`.
     - **Por Grupo de Producto:** Se selecciona un `GroupID`.
     - **Por Marca y Grupo:** Se seleccionan una `BrandID` y un `GroupID`.
- **Postcondiciones:** El atributo está disponible para ser asignado o heredado por productos.

### UC-P-002: Modificar un Atributo

- **Actor:** Administrador del Sistema
- **Descripción:** Permite ajustar los detalles de un atributo existente, incluyendo sus valores o su código.
- **Flujo Básico:**
  1. El administrador selecciona un `Attribute`.
  2. Modifica su `Name`, `Code`, `Order`, `AttributeValue`s (añadir, modificar, eliminar) o su alcance.
  3. El sistema guarda los cambios.
- **Regla Crítica (Impacto en SKU):**
  - **Si se modifica el `Code` de un `Attribute`:** El sistema deberá ejecutar un proceso síncrono y transaccional que identifique todas las `ProductVariant`s dependientes de este `Attribute`, recalcule sus SKUs y actualice sus registros correspondientes. Este proceso debe ser rápido y garantizar la consistencia inmediata.

---

## 2. Gestión de Productos (`Product`)

### UC-P-003: Crear un Producto

- **Actor:** Administrador del Sistema
- **Descripción:** Registra una nueva plantilla de producto, incluyendo sus nombres (corto y largo) y el SKU base.
- **Flujo Básico:**
  1. El administrador proporciona:
     - `Name` (nombre corto)
     - `LongName` (nombre largo)
     - `SKU` (SKU base para la composición de variantes)
     - `Barcode` (opcional)
     - `ProductType` (`TANGIBLE`/`SERVICE`)
     - `BrandID` y `GroupIDs`.
     - `BasePrice` (Precio de coste base para cálculo de variantes).
     - `TaxRate` (Tipo de IVA aplicable, ej: 21.0).
  2. El sistema crea la entidad `Product`.
- **Postcondiciones:** El producto hereda automáticamente todos los `Attribute`s que coincidan con su marca y grupos, además de los genéricos.

### UC-P-004: Asignar un Atributo Directo a un Producto

- **Actor:** Administrador del Sistema
- **Descripción:** Asigna un `Attribute` a un producto específico para anular o complementar la herencia de atributos.
- **Flujo Básico:**
  1. El administrador selecciona un `Product`.
  2. Busca y selecciona un `Attribute` existente.
  3. El sistema añade el `AttributeID` a la lista `DirectAttributeIDs` del producto.

### UC-P-005: Consultar Atributos Aplicables de un Producto

- **Actor:** Administrador del Sistema / Sistema Externo
- **Descripción:** Obtiene la lista completa y calculada de todos los atributos que se aplican a un producto, respetando el orden de precedencia.
- **Flujo Básico:**
  1. El actor solicita los atributos para un `ProductID`.
  2. El sistema ejecuta la lógica de herencia (descrita en el `domain-model.md`).
  3. El sistema devuelve la lista consolidada de `Attribute`s aplicables.

### UC-P-006: Modificar SKU de Producto

- **Actor:** Administrador del Sistema
- **Descripción:** Permite cambiar el SKU base de un producto y propaga el cambio a todos sus `ProductVariant`s.
- **Flujo Básico:**
  1. El administrador selecciona un `Product`.
  2. Modifica el campo `SKU`.
  3. El sistema ejecuta un proceso síncrono y transaccional que:
     - Identifica todas las `ProductVariant`s asociadas a este `Product`.
     - Para cada `ProductVariant`, recalcula su `SKU` utilizando el nuevo `Product.SKU` y los `Attribute.Code`/`AttributeValue.Code` de sus atributos.
     - Actualiza el `SKU` de la `ProductVariant` en la base de datos.
  4. El sistema guarda los cambios en la entidad `Product`.
- **Postcondiciones:** El `Product` tiene su nuevo `SKU` y todas sus `ProductVariant`s tienen sus SKUs actualizados.

---

## 3. Gestión de Variantes de Producto (`ProductVariant`)

El módulo de producto opera bajo el principio fundamental de "Creación Just-in-Time" para las `ProductVariant`s. Esto significa que las variantes se crean en la base de datos solo cuando son necesarias, ya sea explícitamente por un usuario o automáticamente por el sistema ante una demanda. Las variantes creadas automáticamente (Just-in-Time) inician con un `Status = PROVISIONAL` y requieren validación para pasar a `CONFIRMED`.


### UC-P-007: Pre-generar Variantes para un Producto (Opcional)

- **Actor:** Administrador del Sistema
- **Descripción:** Permite generar explícitamente un subconjunto o todas las posibles `ProductVariant` para un producto. Esto es útil para pre-cargar variantes con `Status = CONFIRMED`.
- **Flujo Básico:**
  1. El administrador selecciona un `Product`.
  2. Opcionalmente, puede filtrar o seleccionar subconjuntos de atributos y valores.
  3. El sistema invoca el UC-P-005 para obtener todos los atributos aplicables al producto.
  4. El sistema calcula las combinaciones de `AttributeValue`s.
  5. Para cada combinación:
     - Si la `ProductVariant` no existe, el sistema la crea con el `SKU` determinista y `Status = CONFIRMED`.
     - Si existe con `Status = PROVISIONAL`, actualiza su `Status` a `CONFIRMED`.
- **Postcondiciones:** Las variantes seleccionadas existen con `Status = CONFIRMED` y están listas para operaciones de inventario o venta.

### UC-P-008: Modificar una Variante Específica

- **Actor:** Administrador del Sistema
- **Descripción:** Permite ajustar detalles de una única variante y, si está `PROVISIONAL`, confirmarla.
- **Flujo Básico:**
  1. El administrador busca y selecciona una `ProductVariant` específica (ej: por su SKU).
  2. Modifica su `Barcode`, `IsActive`, o cualquier otro campo editable.
  3. Si la `ProductVariant` tenía `Status = PROVISIONAL`, su `Status` se actualiza a `CONFIRMED`.
  4. El sistema guarda los cambios en la entidad `ProductVariant`.

### UC-P-009: Obtener o Crear Variante (Find or Create)

- **Actor:** Sistema de Ventas (ej: TPV, E-commerce)
- **Descripción:** Proporciona un mecanismo "Just-in-Time" para obtener una `ProductVariant` existente o crearla dinámicamente si no existe, basándose en un `Product` y una combinación de `AttributeValueID`s.
- **Flujo Básico:**
  1. El sistema de ventas solicita una `ProductVariant` para un `ProductID` dado y una lista específica de `AttributeValueID`s.
  2. El sistema busca una `ProductVariant` que coincida con `ProductID` y la lista de `AttributeValueID`s.
  3. **Si se encuentra:** La variante existente se devuelve.
  4. **Si no se encuentra:**
     - El sistema valida que la combinación de `AttributeValueID`s sea válida para el `ProductID` (es decir, que los atributos y valores existan y sean aplicables al producto).
     - Se construye el `SKU` determinista para la nueva variante.
     - Se crea una nueva `ProductVariant` con `Status = PROVISIONAL` y `IsActive = TRUE`.
     - Se persiste en la base de datos.
     - La nueva variante se devuelve.
- **Postcondiciones:** Una `ProductVariant` válida y persistida está disponible para ser utilizada en una transacción. Su `Status` será `PROVISIONAL` si es creada aquí.