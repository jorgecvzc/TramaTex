<template>
  <div class="brand-showcase">
    <div class="showcase-container">
      <!-- CABECERA DE LUJO -->
      <header class="showcase-header">
        <div class="brand-badge">TramaTex Design System v1.0</div>
        <h1>Identidad Visual y Experiencia</h1>
        <p class="lead">
          Nuestra interfaz combina la robustez industrial con una experiencia de usuario moderna y fluida. 
          Este sistema garantiza la coherencia en todos los puntos de contacto del ERP.
        </p>
      </header>

      <!-- SECCIÓN 1: PALETA DE COLORES -->
      <section class="showcase-section">
        <div class="section-intro">
          <h2>01. Paleta Cromática</h2>
          <p>Colores institucionales que definen nuestra marca y estados de negocio.</p>
        </div>
        <div class="color-grid">
          <div v-for="color in brandColors" :key="color.name" class="color-card">
            <div class="color-swatch" :style="{ backgroundColor: color.hex }"></div>
            <div class="color-info">
              <strong>{{ color.name }}</strong>
              <code>{{ color.hex }}</code>
              <small>{{ color.var }}</small>
            </div>
          </div>
        </div>
      </section>

      <!-- SECCIÓN 2: TIPOGRAFÍA -->
      <section class="showcase-section alt-bg">
        <div class="section-intro">
          <h2>02. Tipografía</h2>
          <p>Jerarquía visual diseñada para la lectura de datos complejos y gestión densa.</p>
        </div>
        <div class="typo-demo">
          <div class="typo-row">
            <span class="label">Brand Font</span>
            <h1 class="font-brand" style="color: var(--color-secondary)">TramaTex Industrial</h1>
          </div>
          <div class="typo-row">
            <span class="label">Encabezados</span>
            <div class="headers-stack">
              <h1>H1. Título de Sección</h1>
              <h2>H2. Título de Bloque</h2>
              <h3>H3. Título de Tarjeta</h3>
            </div>
          </div>
          <div class="typo-row">
            <span class="label">Cuerpo</span>
            <p>Texto estándar para descripciones y contenido general. Legible y equilibrado.</p>
          </div>
        </div>
      </section>

      <!-- SECCIÓN 3: COMPONENTES INTERACTIVOS -->
      <section class="showcase-section">
        <div class="section-intro">
          <h2>03. Elementos de Acción</h2>
          <p>Botones y controles con estados claros para una interactividad sin errores.</p>
        </div>
        <div class="components-grid">
          <div class="component-box">
            <h3>Botones Estándar</h3>
            <div class="actions-stack">
              <button class="btn btn-primary">Acción Principal</button>
              <button class="btn btn-secondary">Acción de Proceso</button>
              <button class="btn btn-outline">Acción Neutra</button>
              <button class="btn btn-danger">Acción Crítica</button>
            </div>
          </div>
          <div class="component-box">
            <h3>Variantes de Tamaño</h3>
            <div class="actions-stack">
              <button class="btn btn-primary btn-lg">Grande (LG)</button>
              <button class="btn btn-primary">Normal</button>
              <button class="btn btn-primary btn-sm">Compacto (SM)</button>
            </div>
          </div>
        </div>
      </section>

      <!-- SECCIÓN 4: ICONOGRAFÍA (SSOT) -->
      <section class="showcase-section alt-bg">
        <div class="section-intro">
          <h2>04. Iconografía (Lucide SSOT)</h2>
          <p>Iconografía industrial unificada. El sistema utiliza una <strong>Verdad Única</strong> en <code>src/utils/icons.ts</code>.</p>
        </div>
        
        <div class="icon-search-box mb-6">
          <input v-model="iconSearch" type="text" placeholder="Filtrar iconos por nombre o alias..." class="form-input-lux" />
        </div>

        <div class="icon-grid">
          <div v-for="name in filteredIconNames" :key="name" class="icon-card-lux">
            <div class="icon-preview">
              <component :is="getIcon(name)" :size="24" />
            </div>
            <div class="icon-info">
              <span class="icon-key">{{ name }}</span>
            </div>
          </div>
        </div>
        <p class="mt-4 text-muted"><small>Total de entradas registradas: {{ allIconNames.length }}</small></p>
      </section>

      <!-- SECCIÓN 5: FORMULARIOS (NUEVO ESTÁNDAR) -->
      <section class="showcase-section">
        <div class="section-intro">
          <h2>05. Entrada de Datos</h2>
          <p>Campos de edición diseñados para la precisión y la rapidez en la gestión.</p>
        </div>
        <div class="forms-demo">
          <div class="form-grid">
            <div class="form-group">
              <label>Campo de Texto</label>
              <input type="text" placeholder="Escribe algo aquí..." />
            </div>
            <div class="form-group">
              <label>Selección Desplegable</label>
              <select>
                <option>Opción por defecto</option>
                <option>Opción secundaria</option>
              </select>
            </div>
            <div class="form-group">
              <label>Selector de Fecha</label>
              <input type="date" />
            </div>
            <div class="form-group">
              <label>Estado Deshabilitado</label>
              <input type="text" value="No editable" disabled />
            </div>
          </div>
        </div>
      </section>

      <footer class="showcase-footer">
        <p>&copy; 2026 TramaTex Industrial ERP. Todos los derechos reservados.</p>
        <RouterLink to="/dashboard" class="btn btn-outline btn-sm mt-4">Cerrar Guía de Estilos</RouterLink>
      </footer>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { RouterLink } from 'vue-router'
import { getIcon, getAllIconNames } from '@/utils/icons'

const brandColors = [
  { name: 'Amarillo TramaTex', hex: '#E6B800', var: '--color-primary' },
  { name: 'Azul Corporativo', hex: '#1b3a6b', var: '--color-secondary' },
  { name: 'Azul Suave', hex: '#324e7a', var: '--color-secondary-light' },
  { name: 'Éxito', hex: '#16a34a', var: '--color-success' },
  { name: 'Error', hex: '#dc2626', var: '--color-error' },
  { name: 'Fondo', hex: '#f1f5f9', var: '--color-background' }
]

const iconSearch = ref('')
const allIconNames = getAllIconNames()

const filteredIconNames = computed(() => {
  if (!iconSearch.value) return allIconNames
  const q = iconSearch.value.toLowerCase()
  return allIconNames.filter(n => n.toLowerCase().includes(q))
})
</script>

<style scoped>
.brand-showcase {
  background-color: var(--color-surface);
  min-height: 100vh;
  color: var(--color-text-primary);
  font-family: 'Inter', sans-serif;
}

.showcase-container {
  max-width: 1100px;
  margin: 0 auto;
  padding: 4rem 2rem;
}

.showcase-header {
  text-align: center;
  margin-bottom: 5rem;
}

.brand-badge {
  display: inline-block;
  padding: 0.4rem 1rem;
  background: var(--color-primary-light);
  color: #856404;
  border-radius: 20px;
  font-size: 0.75rem;
  font-weight: 800;
  text-transform: uppercase;
  margin-bottom: 1.5rem;
}

.showcase-header h1 {
  font-size: 3.5rem;
  font-weight: 900;
  letter-spacing: -0.02em;
  margin-bottom: 1.5rem;
  color: var(--color-secondary);
}

.lead {
  font-size: 1.25rem;
  color: var(--color-text-secondary);
  max-width: 700px;
  margin: 0 auto;
  line-height: 1.6;
}

.showcase-section {
  margin-bottom: 4rem;
  padding: 3rem;
  border-radius: 24px;
  border: 1px solid var(--color-border);
}

.showcase-section.alt-bg {
  background-color: var(--color-background);
  border: none;
}

.section-intro {
  margin-bottom: 2.5rem;
}

.section-intro h2 {
  font-size: 1.75rem;
  font-weight: 800;
  color: var(--color-secondary);
  margin-bottom: 0.5rem;
}

.color-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 1.5rem;
}

.color-card {
  background: white;
  padding: 0.75rem;
  border-radius: 16px;
  box-shadow: var(--box-shadow-sm);
}

.color-swatch {
  height: 100px;
  border-radius: 10px;
  margin-bottom: 1rem;
  border: 1px solid rgba(0,0,0,0.05);
}

.color-info strong { display: block; font-size: 0.85rem; }
.color-info code { display: block; font-size: 0.75rem; color: var(--color-text-secondary); margin: 0.25rem 0; }
.color-info small { font-size: 0.65rem; color: var(--color-primary); font-weight: 700; }

.typo-row {
  display: flex;
  margin-bottom: 2.5rem;
  align-items: flex-start;
}

.typo-row .label {
  min-width: 150px;
  font-size: 0.75rem;
  font-weight: 800;
  text-transform: uppercase;
  color: var(--color-text-secondary);
  padding-top: 0.5rem;
}

.headers-stack h1, .headers-stack h2, .headers-stack h3 { margin-bottom: 1rem; }

.actions-stack {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.components-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 3rem;
}

.component-box h3 {
  font-size: 0.9rem;
  text-transform: uppercase;
  color: var(--color-text-secondary);
  margin-bottom: 1.5rem;
}

/* Icon Grid Lux */
.icon-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 1rem;
}

.icon-card-lux {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 1.25rem;
  background: white;
  border: 1px solid var(--color-border);
  border-radius: 12px;
  transition: 0.2s;
}

.icon-card-lux:hover {
  border-color: var(--color-primary);
  transform: translateY(-2px);
  box-shadow: var(--box-shadow-md);
}

.icon-preview {
  margin-bottom: 1rem;
  color: var(--color-secondary);
}

.icon-key {
  font-size: 0.7rem;
  font-family: var(--font-family-mono);
  word-break: break-all;
  color: var(--color-text-secondary);
  text-align: center;
}

.form-input-lux {
  width: 100%;
  padding: 1rem 1.5rem;
  border-radius: 12px;
  border: 2px solid var(--color-border);
  font-family: inherit;
  font-size: 1rem;
  outline: none;
  transition: 0.2s;
}

.form-input-lux:focus {
  border-color: var(--color-primary);
  background: white;
  box-shadow: 0 0 0 4px rgba(230, 184, 0, 0.1);
}

.mb-6 { margin-bottom: 1.5rem; }

.showcase-footer {
  text-align: center;
  padding-top: 4rem;
  border-top: 1px solid var(--color-border);
  color: var(--color-text-secondary);
  font-size: 0.85rem;
}

@media (max-width: 768px) {
  .showcase-header h1 { font-size: 2.5rem; }
  .showcase-section { padding: 1.5rem; }
  .typo-row { flex-direction: column; }
  .typo-row .label { margin-bottom: 1rem; }
}
</style>
