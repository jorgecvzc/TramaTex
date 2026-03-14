# Procesos y Casos de Uso: Módulo Party

Este documento describe los flujos operativos clave para la gestión de terceros, centrando la atención en la lógica de negocio y las transiciones de estado de las entidades.

---

## 1. Ciclo de Vida de la Identidad

### Alta y Configuración de Perfiles
El proceso de registro captura la esencia legal del tercero. El sistema exige la definición de al menos un perfil (**Persona** u **Organización**). 
- **Lógica de Validación:** No se permite una Party "vacía". Al crear una organización, el flujo sugiere inmediatamente la creación de al menos una dirección primaria para habilitar operaciones comerciales futuras.

### Evolución de Roles
A medida que la relación con el tercero progresa, el usuario puede asignar roles comerciales (CLIENT, SUPPLIER).
- **Impacto:** La asignación de un rol activa visibilidad en otros módulos. Un tercero sin el rol `CLIENT` no aparecerá en los selectores del módulo **Sales**, protegiendo la integridad del flujo de venta.

### Bloqueo y Suspensión
El cambio de estado a `BLOCKED` o `INACTIVE` es una medida de control de riesgo.
- **Comportamiento:** Una Party bloqueada conserva su historial (pedidos antiguos, facturas), pero el sistema impide su uso en cualquier documento nuevo. Es una "congelación" operativa completa.

---

## 2. Gestión de Estructuras Complejas

### Establecimiento de Relaciones
Permite conectar identidades para reflejar la realidad del cliente.
- **Uso Común:** Vincular múltiples "Personas" como empleados de una "Organización". Esto facilita que, al crear un pedido para la organización, el sistema pueda sugerir los contactos autorizados para la firma.

### Centralización de Contactos
Para organizaciones, el sistema permite gestionar una agenda de contactos internos.
- **Vínculo Opcional:** Un contacto puede ser simplemente un nombre y teléfono, o puede estar vinculado a una `Party` de tipo persona ya existente, permitiendo que esa persona tenga su propia ficha independiente.

---

## 3. Seguridad e Integridad Operativa

### Borrado Inteligente (Smart Deletion)
El sistema protege la integridad histórica mediante validaciones cruzadas. El proceso de eliminación fallará si la Party tiene:
1. Documentos de venta emitidos (Quote, Order, Invoice).
2. Tareas de producción asignadas o realizadas (MES).
3. Relaciones de jerarquía activas.
**Alternativa:** En estos casos, el sistema sugiere la **Desactivación** en lugar del borrado físico.

### Resolución de Identidad en Lote (Batch Resolution)
Diseñado para optimizar la interfaz de usuario en listados de alta densidad (ej. ver todos los pedidos del mes). 
- **Proceso:** En lugar de solicitar los datos de cada cliente uno a uno, el sistema recolecta todos los IDs necesarios y resuelve sus nombres y estados en una única operación de alto rendimiento.

---
**Última Actualización:** 2026-03-07
