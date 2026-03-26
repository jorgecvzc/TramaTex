---
marp: true
theme: default
paginate: true
header: 'TramaTex'
footer: 'TramaTex - La solución inteligente para la microempresa textil'
style: |
  section {
    font-family: 'Inter', sans-serif;
    background-color: #f1f5f9;
    color: #1e293b;
    font-size: 24px;
  }
  h1, h2 {
    color: #002395;
    font-family: 'Inter', sans-serif;
    font-weight: 800;
  }
  h1 { font-size: 1.8em; border-bottom: 4px solid #E6B800; padding-bottom: 10px; }
  h2 { font-size: 1.4em; }
  strong { color: #002395; }
  .brand { font-family: 'Calibri', sans-serif; font-style: italic; font-weight: bold; color: #002395; }
  .bg-primary { background-color: #002395; color: white; }
  .bg-primary h1 { color: #E6B800; border-color: #E6B800; }
  .card {
    background: white;
    border-radius: 8px;
    padding: 20px;
    box-shadow: 0 4px 6px rgba(0,0,0,0.05);
    border-left: 5px solid #E6B800;
  }
  .accent { color: #E6B800; font-weight: bold; }
  footer, header { color: #64748b; font-size: 12px; }
---

<!-- _class: bg-primary -->
<!-- _header: "" -->
<!-- _footer: "" -->

# <span class="brand">TramaTex</span>
### La solución inteligente de gestión para la microempresa textil

- **Cerrar la Brecha Digital:** Herramienta adaptada a recursos limitados.
- **Organización Integral:** Gestión centralizada de clientes y producción.
- **Impulso al Crecimiento:** Optimización de procesos y trazabilidad total.

---

# 🚀 Visión y Alcance

## Innovación en el Sector Textil

- **Visión:** Sistema líder para PYMES de vestuario laboral y EPIs.
- **Especialización:** Personalización y marcación de prendas.

<div class="card">

**El Desafío:** Superar la ineficiencia de procesos manuales y sistemas genéricos costosos mediante una arquitectura `local-first` eficiente.
</div>

**Objetivo:** Integrar ventas, tarificación inteligente y producción (MES) en un flujo único.

---

# 📊 Dashboard Principal

## Gestión Centralizada

<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 20px;">
<div>

- Interfaz moderna y minimalista.
- Acceso instantáneo a métricas críticas.
- Diseñado para la agilidad operativa.
</div>
<div class="card">

**Valor Añadido:**
Módulos de Clientes, Catálogo, Ventas y Producción a un solo click de distancia.
</div>
</div>

---

# 👥 Módulo Party

## Gestión de Terceros: Enfoque Unificado

- **Entidad Única:** Clientes y proveedores bajo el concepto de `Party`.
- **Roles Dinámicos:** Sin duplicidad de datos; una entidad puede ser ambos.
- **Perfiles:** Soporte para empresas (Organizaciones) y Personas.
- **Relaciones:** Gestión de matrices y múltiples contactos.

<div class="accent">Beneficio: Consistencia total en el histórico de transacciones.</div>

---

# 👕 Módulo Product

## Catálogo Inteligente y Variantes Dinámicas

- **Configurabilidad:** Atributos flexibles (tallas, colores, materiales).
- **Variantes JIT:** Creación automática según demanda (Just-In-Time).
- **Organización:** Clasificación por marcas, familias y grupos.
- **Precisión:** Base sólida para cálculo de costes y producción.

<div class="card">
Manejo sencillo de catálogos complejos de alta rotación sin carga manual masiva.
</div>

---

# 💰 Módulo Sales

## Del Presupuesto a la Factura: Ciclo Total

- **Ciclo de Vida:** Cotizaciones → Pedidos → Albaranes → Facturas.
- **Trazabilidad:** Transiciones automáticas manteniendo el rastro.
- **Flexibilidad:** Soporte para tickets retail y facturación simplificada.
- **Control de Precios:** Motor inteligente con ajustes supervisados.

<div class="accent">Resultado: Cero errores en facturación y visión clara de cada pedido.</div>

---

# 🏭 Módulo MES

## Taller y Producción: El Corazón Operativo

- **Órdenes de Producción:** Transformación automática desde pedidos de venta.
- **Monitorización:** Seguimiento en tiempo real en centros de trabajo.
- **Calidad:** Registro de verificaciones en cada etapa.
- **Sincronización:** Conexión directa oficina-taller.

<div class="card">
Optimización de tiempos de entrega y reducción de errores en personalizaciones.
</div>

---

# 🗓️ Hitos y Futuro

## Realidad Actual y Evolución

<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 20px;">
<div class="card">

### ✅ MVP (Completado)
- Core: Clientes, Productos, Precios.
- Ventas y Gestión de Taller (MES).
- Seguridad y accesos robustos.
</div>
<div>

### 🚀 Hoja de Ruta
- **BI:** Cuadros de mando avanzados.
- **Notificaciones:** WebSockets en tiempo real.
- **Expansión:** Compras, Logística y Contabilidad.
</div>
</div>

---

# 🛠️ Ingeniería de Clase Mundial

## Robustez Tecnológica

- **Arquitectura Modular:** Dominios independientes para estabilidad y crecimiento.
- **Clean Architecture:** Lógica de negocio protegida de cambios tecnológicos.
- **Domain-Driven Design:** El software habla el lenguaje del negocio textil.

<div class="accent">Garantía de longevidad y mantenibilidad del software.</div>

---

# ⚡ El Corazón Tecnológico (Go)

## Backend: Potencia y Eficiencia

<div style="display: grid; grid-template-columns: 1fr 2fr; gap: 20px;">
<div class="card" style="text-align: center; border-left: none; border-bottom: 5px solid #E6B800;">

# Go
PostgreSQL
</div>
<div>

- **Alto Rendimiento:** Procesamiento ultra-rápido.
- **Estabilidad:** Binario único, arranque instantáneo.
- **Concurrencia:** Gestión nativa de procesos simultáneos.
- **Integridad:** Seguridad total de la información financiera.
</div>
</div>

---

# 🎨 Experiencia de Usuario (Vue.js)

## Frontend: Agilidad y Modernidad

- **Tecnología:** Vue.js 3 + TypeScript + Tailwind CSS.
- **Interactividad:** Interfaz reactiva para mayor productividad.
- **Consistencia:** Sistema de diseño unificado (Design System).
- **Adaptabilidad:** Estaciones de trabajo y terminales de taller.

<div class="accent">Calidad: Tipado estricto que elimina errores de interfaz.</div>

---

# 📦 Simplicidad y Despliegue

## Filosofía Local-First

- **Independencia:** Operación 100% local (privacidad y sin internet).
- **Eficiencia:** Optimizado para hardware modesto.

### Facilidad de Gestión:
- **Docker:** Despliegue estandarizado y actualizaciones sin fricción.
- **Mantenimiento:** Diseñado para mínima intervención técnica.

---

# 🏗️ Estandarización de Proyectos

## Scaffolding: Innovación Acelerada

- **Estandarización:** Estructura profesional lista desde el día 1.
- **Calidad:** Pipeline de pruebas y documentación automática.
- **Velocidad:** Reducción drástica del "Time-to-Market" de nuevos módulos.
- **Consistencia:** Alineación total con los estándares de TramaTex.

---

<!-- _class: bg-primary -->
<!-- _header: "" -->
<!-- _footer: "" -->

# <span class="brand">TramaTex</span>
## Preparados para el Siguiente Nivel

- Solución ERP/MES específica para el sector textil.
- Arquitectura moderna, estable y escalable.
- Disciplina de ingeniería para un producto de alto rendimiento.

### **¡Gracias por su atención!**

<div style="margin-top: 40px; font-size: 0.8em; color: #E6B800;">
TramaTex: El aliado tecnológico definitivo.
</div>
