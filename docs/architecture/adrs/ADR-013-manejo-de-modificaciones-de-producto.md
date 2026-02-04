# ADR-013 – Manejo de Modificaciones de Producto (Arreglos/Marcajes)

**Fecha:** 2026-02-02
**Estado:** Aceptado
**Autores:** Gemini, Usuario

---

## 1. Contexto

El módulo de `product` debe gestionar no solo productos vendibles estándar, sino también servicios de modificación como arreglos (ej: ajustar un bajo) o marcajes (ej: serigrafiar un logo). Surge una pregunta fundamental sobre cómo modelar estos servicios, ya que pueden aplicarse tanto a productos nuevos comprados en la tienda como a productos que el cliente ya posee.

La decisión sobre el modelo de datos es crítica porque impactará directamente en la gestión de inventario, la fijación de precios (pricing), el procesamiento de órdenes (sales) y las operaciones de producción (mes).

---

## 2. Alternativas Consideradas

### Alternativa A – Modificaciones como Variantes de un Producto Base

En este modelo, una modificación es una "variante" de un producto.

- **Para un producto nuevo:** Se añadiría una variante "Arreglo" o "Marcaje" al producto principal en la orden de venta. Por ejemplo, una "Camiseta Blanca" tendría una variante "Estampado Logo Cliente".
- **Para un producto del cliente:** Se usaría un producto genérico "Servicio sobre artículo de cliente" con precio base cero, y se le añadiría la variante específica del servicio (ej: "Ajuste de bajo"), que sí tendría un precio.

**Ventajas:**
- **Atomicidad:** La modificación está intrínsecamente ligada al producto que modifica, manteniendo la orden de venta cohesiva.
- **Reutilización:** El concepto de "variante" ya es necesario para atributos como talla y color, por lo que se extiende un concepto existente.
- **Pricing:** El motor de precios podría tratar las variantes de modificación de forma similar a otras variantes que afectan al precio.

**Desventajas:**
- **Complejidad del Producto Genérico:** El caso de uso "producto del cliente" requiere un artificio (producto base con precio cero) que puede ser poco intuitivo y complicar la lógica de negocio.
- **Gestión de Órdenes:** Puede ser complejo diferenciar en una orden si una variante es un atributo (color) o un servicio añadido (arreglo).
- **Modelo de Datos:** El concepto de "variante" se sobrecarga con dos responsabilidades: describir un producto y añadir servicios.

### Alternativa B – Modificaciones como Productos Separados (Tipo "Servicio")

En este modelo, una modificación es un producto en sí mismo, de un tipo especial "Servicio".

- **Para un producto nuevo:** La orden de venta incluiría dos líneas de pedido: una para el producto ("Camiseta Blanca") y otra para el servicio ("Servicio de Estampado").
- **Para un producto del cliente:** La orden de venta incluiría una única línea de pedido para el servicio ("Servicio de Ajuste de Bajo"), posiblemente con una nota indicando el artículo del cliente.

**Ventajas:**
- **Claridad y Simplicidad:** El modelo es explícito y fácil de entender. Un producto es un bien tangible, un servicio es una acción. No hay ambigüedad.
- **Flexibilidad:** Los servicios se pueden vender y gestionar de forma independiente, sin necesidad de asociarlos a un "producto base". Esto simplifica enormemente el caso de "producto del cliente".
- **Gestión de Órdenes (MES):** Es más claro para el sistema de producción (MES) ver un "Servicio de Estampado" como una tarea discreta en una orden de trabajo.
- **Pricing:** Las reglas de precios para servicios pueden ser gestionadas de forma independiente a las de los productos tangibles.

**Desventajas:**
- **Asociación Implícita:** La relación entre el producto y el servicio que se le aplica en una orden de venta es implícita (dos líneas en la misma orden), no explícita en el modelo. Podría requerir lógica adicional para vincularlas si es necesario.
- **Fragmentación del Catálogo:** El catálogo de productos se mezcla con un catálogo de servicios, aunque esto puede mitigarse con una buena categorización.

---

## 3. Criterios de Decisión

- **Simplicidad del Modelo:** ¿Cuán fácil es de entender, implementar y mantener?
- **Cobertura de Casos de Uso:** ¿Cubre de forma natural el escenario de "producto del cliente"?
- **Cohesión vs. Acoplamiento:** ¿El modelo agrupa conceptos que cambian juntos y mantiene separados los que no?
- **Impacto en otros módulos (Pricing, MES):** ¿Qué modelo facilita la implementación en los módulos dependientes?

---

## 4. Decisión Adoptada

**{{PENDIENTE DE DECISIÓN}}**

**Justificación:**
- {{Se rellenará tras la discusión y análisis de las alternativas.}}

---

## 5. Consecuencias

{{Se rellenará tras la decisión.}}

---

## 6. Alcance

Esta decisión aplica al diseño del dominio del módulo `product` y su interacción con los módulos `sales`, `pricing` y `mes`.

---

## 7. Integración con otros ADRs

- N/A

---

## 8. Notas Adicionales / Consideraciones Especiales

Se recomienda evaluar cómo cada alternativa impacta en la experiencia de usuario al crear una orden de venta.

---

## 9. Referencias

- Tarea de Sprint: `docs/log/sprints/sprint-06/01-definicion-dominio-producto.md`

---

## 10. Extensión (Decidido): Configuraciones de Servicio por Cliente

- **Estado de esta sección:** En Análisis

### 10.1. Contexto Adicional

Al adoptar la **Alternativa B** (servicios como productos), surge un requisito avanzado para el `MES` y `Sales`: los clientes recurrentes necesitan guardar "plantillas" de configuración para los servicios.

**Problema:** Un cliente puede pedir frecuentemente el mismo tipo de marcaje (un producto `SERVICE`) pero con detalles específicos (un logo concreto, en una posición determinada). Se necesita una forma de guardar estas "preferencias" o "configuraciones" para no tener que especificarlas en cada nueva venta.

**Ejemplo:**
- El Cliente "Construcciones Acme" siempre pide el servicio "Marcaje de Logo" con su fichero `acme-logo.png` en la manga izquierda de las prendas.
- Puede tener varias configuraciones: una para ropa de trabajo y otra para ropa de eventos.

### 10.2. Propuesta de Solución

Se propone la creación de una nueva entidad, `PartyServiceConfiguration`, que actúe como un puente entre el módulo `party` y el módulo `product`.

- **Entidad Propuesta: `PartyServiceConfiguration`**
  - `ConfigurationID` (ID único)
  - `PartyID` (FK a `Party`: el cliente)
  - `ServiceID` (FK a `Product`: el producto de tipo `SERVICE`)
  - `Name` (Nombre descriptivo para la plantilla, ej: "Logo estándar pecho izquierdo")
  - `ConfigurationDetails` (JSON o similar): Un campo flexible para guardar los detalles específicos de la configuración (ej: `{ "logo_url": "...", "position": "left_chest", "thread_color": "blue" }`).

### 10.3. Ubicación en el Dominio

La ubicación más lógica para este nuevo agregado sería el módulo `party`, ya que representa una propiedad o configuración específica de un cliente (`Party`). De este modo, el módulo `product` se mantiene genérico.

### 10.4. Impacto

- **Módulo `sales`:** Al añadir un producto de tipo `SERVICE` a una orden para un cliente, la API debería permitir consultar y seleccionar una de estas configuraciones guardadas.
- **Módulo `mes`:** La orden de trabajo para el servicio debe incluir los `ConfigurationDetails` para que producción sepa exactamente qué hacer.

Esta extensión se estudiará en detalle antes de finalizar los contratos de API.
