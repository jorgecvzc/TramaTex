<template>
  <div class="style-guide-wrapper">
    <div class="style-guide">
      <h1 class="font-brand" style="font-size: 2.5rem; color: var(--color-secondary);">Guía de Estilos de TramaTex</h1>
      <p>Este componente sirve como una guía visual de los elementos de UI definidos en el sistema de diseño, alineados con los mockups más recientes.</p>

      <!-- Sección de Colores -->
      <section>
        <h2>Paleta de Colores</h2>
        <div class="color-grid">
          <div class="color-swatch-wrapper">
            <div class="color-swatch" style="background-color: var(--color-primary);"></div>
            <p>Primary (Yellow)</p>
            <code>var(--color-primary)</code>
          </div>
          <div class="color-swatch-wrapper">
            <div class="color-swatch" style="background-color: var(--color-secondary);"></div>
            <p>Secondary (Blue)</p>
            <code>var(--color-secondary)</code>
          </div>
          <div class="color-swatch-wrapper">
            <div class="color-swatch" style="background-color: var(--color-surface); border: 1px solid var(--color-border);"></div>
            <p>Surface</p>
            <code>var(--color-surface)</code>
          </div>
          <div class="color-swatch-wrapper">
            <div class="color-swatch" style="background-color: var(--color-background); border: 1px solid var(--color-border);"></div>
            <p>Background</p>
            <code>var(--color-background)</code>
          </div>
          <div class="color-swatch-wrapper">
            <div class="color-swatch" style="background-color: var(--color-text-primary);"></div>
            <p>Text Primary</p>
            <code>var(--color-text-primary)</code>
          </div>
          <div class="color-swatch-wrapper">
            <div class="color-swatch" style="background-color: var(--color-text-on-primary); border: 1px solid var(--color-border);"></div>
            <p>Text on Primary</p>
            <code>var(--color-text-on-primary)</code>
          </div>
          <div class="color-swatch-wrapper">
            <div class="color-swatch" style="background-color: var(--color-text-on-secondary);"></div>
            <p>Text on Secondary</p>
            <code>var(--color-text-on-secondary)</code>
          </div>
        </div>
      </section>

      <!-- Sección de Botones (NUEVA) -->
      <section>
        <h2>Estandarización de Botones</h2>
        <p class="section-intro">Todos los botones deben usar la clase base <code>.btn</code> seguida de su tipo semántico.</p>
        
        <div class="button-showcase">
          <table class="showcase-table">
            <thead>
              <tr>
                <th>Tipo / Propósito</th>
                <th>Normal</th>
                <th>Con Icono</th>
                <th>Deshabilitado</th>
                <th>Clase CSS</th>
              </tr>
            </thead>
            <tbody>
              <!-- Primary -->
              <tr>
                <td><strong>Acción Principal</strong><br><small>Guardar, Crear, Aceptar</small></td>
                <td><button class="btn btn-primary">Guardar</button></td>
                <td><button class="btn btn-primary"><Save :size="18" /><span>Guardar</span></button></td>
                <td><button class="btn btn-primary" disabled>Guardar</button></td>
                <td><code>.btn-primary</code></td>
              </tr>
              <!-- Secondary -->
              <tr>
                <td><strong>Acción de Proceso</strong><br><small>Confirmar, Facturar, Albarán</small></td>
                <td><button class="btn btn-secondary">Confirmar</button></td>
                <td><button class="btn btn-secondary"><CheckCircle :size="18" /><span>Confirmar</span></button></td>
                <td><button class="btn btn-secondary" disabled>Confirmar</button></td>
                <td><code>.btn-secondary</code></td>
              </tr>
              <!-- Outline -->
              <tr>
                <td><strong>Acción Neutra</strong><br><small>Cancelar, Volver, Filtros</small></td>
                <td><button class="btn btn-outline">Cancelar</button></td>
                <td><button class="btn btn-outline"><X :size="18" /><span>Cancelar</span></button></td>
                <td><button class="btn btn-outline" disabled>Cancelar</button></td>
                <td><code>.btn-outline</code></td>
              </tr>
              <!-- Danger -->
              <tr>
                <td><strong>Acción Crítica</strong><br><small>Eliminar, Anular, Borrar</small></td>
                <td><button class="btn btn-danger">Eliminar</button></td>
                <td><button class="btn btn-danger"><Trash2 :size="18" /><span>Eliminar</span></button></td>
                <td><button class="btn btn-danger" disabled>Eliminar</button></td>
                <td><code>.btn-danger</code></td>
              </tr>
              <!-- Success -->
              <tr>
                <td><strong>Acción Positiva</strong><br><small>Aprobar, Completar, Validar</small></td>
                <td><button class="btn btn-success">Completar</button></td>
                <td><button class="btn btn-success"><CheckCircle2 :size="18" /><span>Completar</span></button></td>
                <td><button class="btn btn-success" disabled>Completar</button></td>
                <td><code>.btn-success</code></td>
              </tr>
              <!-- Ghost -->
              <tr>
                <td><strong>Acción Minimalista</strong><br><small>Iconos en tablas, Cierres</small></td>
                <td><button class="btn btn-ghost">Editar</button></td>
                <td><button class="btn btn-ghost"><Pencil :size="18" /></button></td>
                <td><button class="btn btn-ghost" disabled><Pencil :size="18" /></button></td>
                <td><code>.btn-ghost</code></td>
              </tr>
            </tbody>
          </table>
        </div>

        <h3 class="mt-8">Tamaños de Botón</h3>
        <div class="sizes-row">
          <div class="size-item">
            <button class="btn btn-primary btn-sm">Pequeño (sm)</button>
            <code>.btn-sm</code>
          </div>
          <div class="size-item">
            <button class="btn btn-primary">Normal (base)</button>
            <span>Por defecto</span>
          </div>
          <div class="size-item">
            <button class="btn btn-primary btn-lg">Grande (lg)</button>
            <code>.btn-lg</code>
          </div>
        </div>
      </section>

      <!-- Sección de Tipografía -->
      <section>
        <h2>Tipografía</h2>
        <h1 class="font-brand" style="color: var(--color-secondary)">TramaTex (Brand Font)</h1>
        <h1>Heading 1: El veloz murciélago hindú</h1>
        <h2>Heading 2: comía feliz cardillo y kiwi.</h2>
        <h3>Heading 3: La cigüeña tocaba el saxofón</h3>
        <h4>Heading 4: detrás del palenque de paja.</h4>
        <p>Párrafo: Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed non risus. Suspendisse lectus tortor, dignissim sit amet, adipiscing nec, ultricies sed, dolor. Cras elementum ultrices diam. Maecenas ligula massa, varius a, semper congue, euismod non, mi. </p>
        <p>Este es un <a href="#">enlace de ejemplo</a>.</p>
        <p>Y aquí tenemos <code>código inline</code>.</p>
      </section>

      <!-- Sección de Iconografía (SSOT) -->
      <section>
        <h2>Iconografía (Lucide SSOT)</h2>
        <p class="section-intro">Iconografía industrial unificada. El sistema utiliza <code>getIcon(name)</code> para mapear nombres antiguos de Material Symbols a los nuevos componentes Lucide.</p>
        
        <div class="icon-search-box mb-4">
          <input v-model="iconSearch" type="text" placeholder="Filtrar iconos por nombre..." class="form-input" />
        </div>

        <div class="icon-grid">
          <div v-for="name in filteredIconNames" :key="name" class="icon-card">
            <div class="icon-preview">
              <component :is="getIcon(name)" :size="24" />
            </div>
            <div class="icon-info">
              <span class="icon-key">{{ name }}</span>
            </div>
          </div>
        </div>
        <p class="mt-4 text-muted"><small>Total de iconos registrados: {{ allIconNames.length }}</small></p>
      </section>

    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { Save, CheckCircle, X, Trash2, CheckCircle2, Pencil } from 'lucide-vue-next';
import { getIcon, getAllIconNames } from '@/utils/icons';

const iconSearch = ref('');
const allIconNames = getAllIconNames();

const filteredIconNames = computed(() => {
  if (!iconSearch.value) return allIconNames;
  const q = iconSearch.value.toLowerCase();
  return allIconNames.filter(n => n.toLowerCase().includes(q));
});
</script>

<style scoped>
/* ... (rest of styles) ... */

/* Icon Grid */
.icon-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 1rem;
}

.icon-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 1rem;
  background: var(--color-background);
  border: 1px solid var(--color-border);
  border-radius: 8px;
  transition: 0.2s;
}

.icon-card:hover {
  border-color: var(--color-primary);
  transform: translateY(-2px);
  box-shadow: var(--box-shadow-sm);
}

.icon-preview {
  margin-bottom: 0.75rem;
  color: var(--color-secondary);
}

.icon-info {
  text-align: center;
  width: 100%;
}

.icon-key {
  font-size: 0.7rem;
  font-family: var(--font-family-mono);
  word-break: break-all;
  color: var(--color-text-secondary);
}

.mb-4 { margin-bottom: 1rem; }
.form-input {
  width: 100%;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  border: 1px solid var(--color-border);
  font-family: inherit;
}
</style>
