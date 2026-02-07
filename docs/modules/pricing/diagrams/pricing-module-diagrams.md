# Diagramas del Módulo de Precios

Este documento presenta los diagramas de módulo para el Bounded Context de Precios, siguiendo la filosofía de Domain-Driven Design y la Arquitectura Limpia adoptada en el proyecto TramaTex.

## 1. Diagrama de Componentes (C3) - Módulo de Precios

Este diagrama ilustra la estructura interna del módulo de Precios, mostrando sus capas de Arquitectura Limpia y los componentes clave (Entidades, Value Objects, Servicios de Dominio, Interfaces de Repositorio e Implementaciones de Infraestructura) tal como se definen en ADR-016.

```mermaid
%%{init: {'flowchart': {'defaultRenderer': 'dagre'}}}%%
graph TD
    subgraph "Pricing Bounded Context (tramatex-api)"
        direction LR
        subgraph "Interfaces Layer (API)"
            InterfacesAPI[API Endpoints: <br/>Get Selling Price, <br/>Calculate Sales Discount, <br/>Get Final Price & Discounts]
        end

        subgraph "Application Layer (Use Cases)"
            ApplicationUC[Use Cases: <br/>GetSellingPriceUC, <br/>CalculateSalesDiscountUC, <br/>GetFinalSellingPriceAndDiscountsUC]
        end

        subgraph "Domain Layer (Core Business Logic)"
            direction TB
            subgraph "Entities"
                PricingRule(PricingRule)
                ClientPricing(ClientPricing)
                ProductBasePrice(ProductBasePrice)
                VariantModifier(VariantModifier)
                BrandProfitMargin(BrandProfitMargin)
                SalesDiscountRule(SalesDiscountRule)
                PriceCalculation(PriceCalculation)
            end

            subgraph "Value Objects"
                Money[Money]
                Percentage[Percentage]
                Brand[Brand]
                ProductCode[ProductCode/VariantCode]
            end

            subgraph "Domain Services"
                SellingPriceCalculatorService[SellingPriceCalculatorService]
                SalesDiscountCalculatorService[SalesDiscountCalculatorService]
            end

            subgraph "Domain Interfaces"
                IPRRepo(IPricingRuleRepository)
                ICPRepo(IClientPricingRepository)
                IBPRepo(IBrandProfitMarginRepository)
                ISDRRepo(ISalesDiscountRuleRepository)
                IPCCRepo(IPriceCalculationRepository)
            end
        end

        subgraph "Infrastructure Layer (Persistence, External Services)"
            direction TB
            InfraPersistence[Pricing Repositories Impl (GORM)]
            InfraCache[In-Memory NoSQL Cache (Selling Price)]
            InfraProductClient[Product Module Client (e.g., HTTP/GRPC)]
            InfraPartyClient[Party Module Client (e.g., HTTP/GRPC)]
        end

        InterfacesAPI --> ApplicationUC
        
        ApplicationUC --> SellingPriceCalculatorService
        ApplicationUC --> SalesDiscountCalculatorService
        ApplicationUC --> IPCCRepo
        ApplicationUC --> InfraCache
        ApplicationUC --> InfraProductClient
        ApplicationUC --> InfraPartyClient

        SellingPriceCalculatorService --> PricingRule
        SellingPriceCalculatorService --> ProductBasePrice
        SellingPriceCalculatorService --> VariantModifier
        SellingPriceCalculatorService --> BrandProfitMargin
        SellingPriceCalculatorService --> Money
        SellingPriceCalculatorService --> Percentage
        SellingPriceCalculatorService --> Brand
        SellingPriceCalculatorService --> ProductCode
        SellingPriceCalculatorService --> IPRRepo
        SellingPriceCalculatorService --> IBPRepo
        
        SalesDiscountCalculatorService --> SalesDiscountRule
        SalesDiscountCalculatorService --> ClientPricing
        SalesDiscountCalculatorService --> Money
        SalesDiscountCalculatorService --> Percentage
        SalesDiscountCalculatorService --> ISDRRepo
        SalesDiscountCalculatorService --> ICPRepo

        IPRRepo -.-> InfraPersistence
        ICPRepo -.-> InfraPersistence
        IBPRepo -.-> InfraPersistence
        ISDRRepo -.-> InfraPersistence
        IPCCRepo -.-> InfraPersistence
    end

    style InterfacesAPI fill:#f9f,stroke:#333,stroke-width:2px
    style ApplicationUC fill:#bbf,stroke:#333,stroke-width:2px
    style SellingPriceCalculatorService fill:#fcc,stroke:#333,stroke-width:2px
    style SalesDiscountCalculatorService fill:#fcc,stroke:#333,stroke-width:2px
    style PricingRule fill:#afa,stroke:#333,stroke-width:2px
    style ClientPricing fill:#afa,stroke:#333,stroke-width:2px
    style ProductBasePrice fill:#afa,stroke:#333,stroke-width:2px
    style VariantModifier fill:#afa,stroke:#333,stroke-width:2px
    style BrandProfitMargin fill:#afa,stroke:#333,stroke-width:2px
    style SalesDiscountRule fill:#afa,stroke:#333,stroke-width:2px
    style PriceCalculation fill:#afa,stroke:#333,stroke-width:2px
    style Money fill:#def,stroke:#333,stroke-width:2px
    style Percentage fill:#def,stroke:#333,stroke-width:2px
    style Brand fill:#def,stroke:#333,stroke-width:2px
    style ProductCode fill:#def,stroke:#333,stroke-width:2px
    style IPRRepo fill:#ffd,stroke:#333,stroke-width:2px
    style ICPRepo fill:#ffd,stroke:#333,stroke-width:2px
    style IBPRepo fill:#ffd,stroke:#333,stroke-width:2px
    style ISDRRepo fill:#ffd,stroke:#333,stroke-width:2px
    style IPCCRepo fill:#ffd,stroke:#333,stroke-width:2px
    style InfraPersistence fill:#dff,stroke:#333,stroke-width:2px
    style InfraCache fill:#dff,stroke:#333,stroke-width:2px
    style InfraProductClient fill:#dff,stroke:#333,stroke-width:2px
    style InfraPartyClient fill:#dff,stroke:#333,stroke-width:2px
```

## 2. Descripción de Componentes y Flujos

### Interfaces Layer (API)
Define los puntos de entrada HTTP para interactuar con el módulo de Precios, exponiendo los casos de uso definidos en ADR-016.

### Application Layer (Use Cases)
Contiene la orquestación de la lógica de negocio. Los Casos de Uso (UC) invocan a los Servicios de Dominio y Repositorios para cumplir con la funcionalidad requerida, gestionan transacciones y mapean DTOs.

### Domain Layer (Core Business Logic)
El corazón del módulo. Contiene:
*   **Entidades:** Objetos con identidad que encapsulan reglas de negocio.
*   **Value Objects:** Objetos inmutables que representan conceptos de dominio.
*   **Domain Services:** Contienen lógica de negocio que involucra múltiples Entidades o VOs, orquestando su comportamiento.
*   **Domain Interfaces (Repositorios):** Abstracciones que definen contratos para la persistencia y recuperación de Entidades.

### Infrastructure Layer (Persistence, External Services)
Contiene las implementaciones concretas de las interfaces de dominio, así como clientes para servicios externos.
*   **Pricing Repositories Impl (GORM):** Implementaciones de los repositorios de dominio utilizando GORM para interactuar con las nuevas tablas dedicadas de `Pricing`.
*   **In-Memory NoSQL Cache (Selling Price):** Implementación de la caché para almacenar precios de venta calculados.
*   **Product Module Client:** Cliente para comunicarse con el módulo `Product` (ej. vía HTTP/gRPC) para obtener `ProductBasePrice`, `VariantModifier` y `BrandProfitMargin`.
*   **Party Module Client:** Cliente para comunicarse con el módulo `Party` para obtener información del cliente (ej. para `ClientPricing`).

### Flujos Clave
*   Los `API Endpoints` delegan en los `Use Cases`.
*   Los `Use Cases` orquestan `Domain Services` y acceden a `Domain Interfaces`.
*   Los `Domain Services` (`SellingPriceCalculatorService`, `SalesDiscountCalculatorService`) encapsulan la lógica de cálculo principal, utilizando Entidades, Value Objects y las interfaces de los repositorios.
*   Las `Domain Interfaces` son implementadas por `Pricing Repositories Impl` en la capa de `Infrastructure`.
*   El módulo `Pricing` interactúa con los módulos `Product` y `Party` a través de clientes específicos en la capa de `Infrastructure`, obteniendo los datos base necesarios sin acoplamiento directo a sus bases de datos internas.
*   Los precios de venta calculados son gestionados por la `In-Memory NoSQL Cache` para optimizar la latencia.
