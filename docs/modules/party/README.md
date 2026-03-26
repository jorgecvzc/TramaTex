# Módulo Party: El Núcleo de Identidad

El módulo Party es el sistema de gestión de terceros de TramaTex. No es simplemente una base de datos de clientes; es el motor que unifica la identidad legal, comercial y operativa de cualquier entidad que interactúa con el ERP.

---

## 🎯 Propósito Estratégico

La misión de este módulo es eliminar la fragmentación del conocimiento. Al centralizar a clientes, proveedores y operarios en una estructura única (**Party**), el sistema garantiza que la información fiscal y de contacto sea consistente en todo el ciclo de vida del producto, desde la compra de materia prima hasta la entrega al cliente final.

---

## 🧩 Guías de Referencia

Para comprender el funcionamiento y la integración de este módulo, consulta las siguientes guías conceptuales:

1. **[Especificación del Módulo](./module-spec.md):** Definición formal de entidades, estados y componentes.
2. **[Lógica de Identidad y Dominio](./domain-model.md):** Explica cómo conviven los perfiles de personas y organizaciones, y cómo fluyen los roles comerciales.
3. **[Estrategia de Interacción (API)](./api-contracts.md):** Describe los puntos de integración y la intención detrás de las operaciones de identidad.
4. **[Procesos Operativos](./use-cases.md):** Guía sobre el ciclo de vida de un tercero, desde el alta hasta el bloqueo preventivo.
5. **[Estándares de Integración](./implementation-guide.md):** Directrices para que otros módulos consuman la información de Party respetando la soberanía del dato.
6. **[Diagrama de Dominio](./diagrams/domain-model.md):** Representación visual del modelo de dominio.

---

## 🏗️ Relación con el Ecosistema TramaTex

Party es un **Servicio de Contexto**. Su valor reside en alimentar a los demás módulos:
- **A Sales:** Le entrega la "Cara" del cliente y sus condiciones fiscales.
- **A Pricing:** Le proporciona los segmentos de cliente para aplicar tarifas dinámicas.
- **A MES:** Le identifica a los actores físicos que operan las máquinas del taller.

---
**Última Versión:** 2026-03-07
