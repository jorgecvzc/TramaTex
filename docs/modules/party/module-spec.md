# Módulo de Party

## 1. Propósito

*   **Visión del Módulo:** Gestionar Parties (personas, organizaciones o ambas) con roles, relaciones y contactos.
*   **Objetivos Clave:**
    *   Proporcionar un sistema centralizado para gestionar todos los clientes y proveedores de la empresa.
    *   Gestionar información de contacto, direcciones, categorías, y relaciones comerciales.

## 2. Requisitos



### 2.1. Requisitos Funcionales



*   **RF-001:** Registro de nuevas Parties (persona, organización o ambas).

*   **RF-002:** Gestión de perfiles (persona/organización).

*   **RF-003:** Gestión de contactos de organización.

*   **RF-004:** Gestión de roles y relaciones entre Parties.

*   **RF-005:** Clasificación por categorías. **(Post-MVP)**

*   **RF-006:** Historial de interacciones. **(Post-MVP)**



## 3. Casos de Uso



### 3.1. Actores



*   **Vendedor:** Empleado que gestiona clientes.

*   **Comprador:** Empleado que gestiona proveedores.



### 3.2. Casos de Uso Principales (MVP)



#### Party

- **Crear Party:** Registrar una Party con perfil persona/organización.

- **Actualizar Party:** Actualizar perfiles y datos básicos.

- **Obtener Party:** Obtener detalles por ID.

- **Listar Parties:** Listar con filtros por rol, tipo y estado.

- **Cambiar Estado de Party:** Activar/desactivar.

#### Roles

- **Añadir Rol a Party**
- **Eliminar Rol de Party**

#### Relaciones

- **Crear Relación entre Parties**
- **Listar Relaciones**
- **Eliminar Relación**

#### ContactDetails (Organización)

- **Añadir Contacto**
- **Listar Contactos**
- **Actualizar Contacto**
- **Eliminar Contacto**



## 4. Historias de Usuario



*   **HU-001:** Como vendedor, quiero registrar una Party cliente para gestionar sus pedidos.

*   **HU-002:** Como comprador, quiero registrar una Party proveedor para gestionar órdenes de compra.

*   **HU-003:** Como vendedor, quiero agregar múltiples contactos de organización para gestionar comunicación.



## 5. Criterios de Aceptación



*   **Para HU-001:**

    *   **Criterio 1:** Dado que ingreso un perfil válido, cuando guardo el formulario, entonces se crea una nueva Party en el sistema.





## 7. Decisiones de Diseño



*   **Relaciones con Otros Módulos:**

    *   **IAM**: Cada Party puede estar vinculada a un usuario (opcional).

    *   **Product**: Party compra/vende productos.

    *   **Pricing**: Diferentes precios por categoría de Party.

    *   **Sales**: Órdenes se relacionan con Party.

*   **Fases de Desarrollo:**

    *   [X] Fase 1 (MVP): Casos de uso básicos (Crear, Leer, Actualizar, Listar).

    *   [ ] Fase 2: Categorías y clasificación **(Post-MVP)**.

    *   [ ] Fase 3: Historial de transacciones **(Post-MVP)**.
