# Módulo de Party (Clientes y Proveedores)

## 1. Propósito

*   **Visión del Módulo:** Gestionar la información centralizada de clientes y proveedores.
*   **Objetivos Clave:**
    *   Proporcionar un sistema centralizado para gestionar todos los clientes y proveedores de la empresa.
    *   Gestionar información de contacto, direcciones, categorías, y relaciones comerciales.

## 2. Requisitos

### 2.1. Requisitos Funcionales

*   **RF-001:** Registro de nuevos clientes y proveedores.
*   **RF-002:** Gestión de información de contacto.
*   **RF-003:** Gestión de direcciones (envío, facturación).
*   **RF-004:** Clasificación de clientes/proveedores por categoría.
*   **RF-005:** Historial de interacciones.

## 3. Casos de Uso

### 3.1. Actores

*   **Vendedor:** Empleado que gestiona clientes.
*   **Comprador:** Empleado que gestiona proveedores.

### 3.2. Casos de Uso Principales

*   **CU-001: CreateParty**
    *   **Actor:** Vendedor/Comprador
    *   **Descripción:** Registrar un nuevo cliente o proveedor.
*   **CU-002: UpdateParty**
    *   **Actor:** Vendedor/Comprador
    *   **Descripción:** Actualizar la información de un cliente o proveedor.
*   **CU-003: GetParty**
    *   **Actor:** Vendedor/Comprador
    *   **Descripción:** Obtener los detalles de un cliente o proveedor.
*   **CU-004: ListParties**
    *   **Actor:** Vendedor/Comprador
    *   **Descripción:** Listar clientes o proveedores con filtros.
*   **CU-005: AddAddress**
    *   **Actor:** Vendedor/Comprador
    *   **Descripción:** Agregar una dirección a un cliente o proveedor.
*   **CU-006: ChangeStatus**
    *   **Actor:** Vendedor/Comprador
    *   **Descripción:** Cambiar el estado (activo/inactivo) de un cliente o proveedor.

## 4. Historias de Usuario

*   **HU-001:** Como vendedor, quiero registrar nuevos clientes para poder gestionar sus pedidos.
*   **HU-002:** Como comprador, quiero registrar nuevos proveedores para poder gestionar las órdenes de compra.
*   **HU-003:** Como vendedor, quiero poder agregar múltiples direcciones a un cliente para gestionar envíos a diferentes sucursales.

## 5. Criterios de Aceptación

*   **Para HU-001:**
    *   **Criterio 1:** Dado que ingreso el nombre legal y el tipo de un nuevo cliente, cuando guardo el formulario, entonces se crea una nueva "party" en el sistema.


## 7. Decisiones de Diseño

*   **Relaciones con Otros Módulos:**
    *   **IAM**: Cada Party puede estar vinculada a un usuario (opcional).
    *   **Product**: Party compra/vende productos.
    *   **Pricing**: Diferentes precios por categoría de Party.
    *   **Sales**: Órdenes se relacionan con Party.
*   **Fases de Desarrollo:**
    *   [X] Fase 1 (MVP): Casos de uso básicos (Crear, Leer, Actualizar, Listar).
    *   [ ] Fase 2: Categorías y clasificación.
    *   [ ] Fase 3: Historial de transacciones.