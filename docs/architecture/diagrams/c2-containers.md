# Diagrama de Contenedores (C2)

**Versión:** 1.0
**Fecha:** 2026-01-31
**Autor:** Gemini
**Propósito:** Describir los principales bloques ejecutables o "contenedores" del sistema TramaTex y cómo interactúan.

---

## Diagrama

```mermaid
graph TD
    subgraph "Sistema TramaTex"
        direction LR
        
        user(👤<br>Usuario<br>[Comercial / Taller])
        
        subgraph "Contenedores"
            direction TB
            frontend(💻<br><b>Frontend</b><br><div style="font-size: 0.8em">Vue.js SPA</div>)
            api(🚀<br><b>Backend API</b><br><div style="font-size: 0.8em">Go Monolito Modular</div>)
            db(🗃️<br><b>Base de Datos</b><br><div style="font-size: 0.8em">PostgreSQL</div>)
        end

        user -- "Usa (HTTPS)" --> frontend
        frontend -- "Realiza llamadas API (JSON/HTTPS)" --> api
        api -- "Lee/Escribe" --> db
    end

    style user fill:#fff,stroke:#333,stroke-width:2px
    style frontend fill:#D1E7DD,stroke:#28A745,stroke-width:2px
    style api fill:#CDE8F5,stroke:#007BFF,stroke-width:2px
    style db fill:#F8D7DA,stroke:#DC3545,stroke-width:2px
```

---

## Descripción

Este diagrama muestra las principales piezas de software que componen el sistema TramaTex.

1.  **Usuario:** Representa a los actores humanos que interactúan con el sistema, como el personal de ventas o los operarios del taller. Acceden al sistema a través de un navegador web.

2.  **Frontend (Single-Page Application):** Es la aplicación web con la que el usuario interactúa directamente. Está construida con Vue.js y se ejecuta en el navegador del usuario. Es responsable de toda la presentación visual y la interacción. No contiene lógica de negocio crítica.

3.  **Backend API (Monolito Modular):** Es el corazón del sistema, desarrollado en Go. Expone una API REST que el frontend consume. Contiene toda la lógica de negocio, reglas de dominio y orquestación de casos de uso. A pesar de ser un único proceso (monolito), está dividido internamente en módulos lógicos (Party, Product, Pricing, etc.).

4.  **Base de Datos (PostgreSQL):** Es el sistema de persistencia donde se almacenan todos los datos del negocio (clientes, productos, pedidos, etc.). Solo la API de Backend tiene acceso directo a la base de datos.
