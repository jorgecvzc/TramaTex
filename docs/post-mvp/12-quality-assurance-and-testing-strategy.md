# Estrategia de Calidad: QA y Testing Avanzado (Post-MVP)

Este documento define el plan para alcanzar una robustez total en TramaTex, asegurando que cada cambio sea verificado automáticamente antes de llegar a producción.

---

## 1. Objetivos de Cobertura (Coverage)

Se elevarán los estándares de calidad actuales:
- **Capa de Dominio (Domain)**: **100% de cobertura**. No puede existir lógica de negocio (reglas de precios, estados de factura, cálculos MES) sin test unitario.
- **Capa de Aplicación (Use Cases)**: **≥95% de cobertura**. Validación de todos los flujos de orquestación.
- **Frontend**: Test de componentes críticos y lógica de Pinia al 80%.

---

## 2. Tests de Integración y E2E (Playwright)

Implementación de una suite de tests de extremo a extremo que simule la operativa real del usuario:
- **Flujo Crítico 1 (Comercial)**: Presupuesto -> Pedido -> Albarán -> Factura -> Registro de Cobro.
- **Flujo Crítico 2 (Manufactura)**: Creación OT -> Inicio Fase -> Registro Tiempos -> Finalización OT.
- **Validación Legal**: Test automático que verifique la generación de XML firmados y su encadenamiento Hash.

---

## 3. Automatización en el Pipeline (CI/CD)

- **Calidad de Código (Linting & Static Analysis)**: Bloqueo de PRs que no cumplan con las reglas de estilo o presenten vulnerabilidades de seguridad (SonarQube/Snyk).
- **Entornos de Preview**: Despliegue automático de un entorno efímero por cada rama de Git para permitir la validación manual antes del merge a `develop`.

---

## 4. Tests de Rendimiento y Carga

- **K6 / JMeter**: Ejecución periódica de pruebas de carga para asegurar que el sistema soporta 100 operarios concurrentes y el motor de precios responde en <200ms bajo estrés.

---

*Última actualización: 2026-04-27*
