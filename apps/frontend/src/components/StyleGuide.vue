<template>
  <div class="style-guide-wrapper">
    <Navbar />
    
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
                <td><button class="btn btn-primary"><span class="material-symbols-outlined">save</span><span>Guardar</span></button></td>
                <td><button class="btn btn-primary" disabled>Guardar</button></td>
                <td><code>.btn-primary</code></td>
              </tr>
              <!-- Secondary -->
              <tr>
                <td><strong>Acción de Proceso</strong><br><small>Confirmar, Facturar, Albarán</small></td>
                <td><button class="btn btn-secondary">Confirmar</button></td>
                <td><button class="btn btn-secondary"><span class="material-symbols-outlined">check_circle</span><span>Confirmar</span></button></td>
                <td><button class="btn btn-secondary" disabled>Confirmar</button></td>
                <td><code>.btn-secondary</code></td>
              </tr>
              <!-- Outline -->
              <tr>
                <td><strong>Acción Neutra</strong><br><small>Cancelar, Volver, Filtros</small></td>
                <td><button class="btn btn-outline">Cancelar</button></td>
                <td><button class="btn btn-outline"><span class="material-symbols-outlined">close</span><span>Cancelar</span></button></td>
                <td><button class="btn btn-outline" disabled>Cancelar</button></td>
                <td><code>.btn-outline</code></td>
              </tr>
              <!-- Danger -->
              <tr>
                <td><strong>Acción Crítica</strong><br><small>Eliminar, Anular, Borrar</small></td>
                <td><button class="btn btn-danger">Eliminar</button></td>
                <td><button class="btn btn-danger"><span class="material-symbols-outlined">delete</span><span>Eliminar</span></button></td>
                <td><button class="btn btn-danger" disabled>Eliminar</button></td>
                <td><code>.btn-danger</code></td>
              </tr>
              <!-- Success -->
              <tr>
                <td><strong>Acción Positiva</strong><br><small>Aprobar, Completar, Validar</small></td>
                <td><button class="btn btn-success">Completar</button></td>
                <td><button class="btn btn-success"><span class="material-symbols-outlined">task_alt</span><span>Completar</span></button></td>
                <td><button class="btn btn-success" disabled>Completar</button></td>
                <td><code>.btn-success</code></td>
              </tr>
              <!-- Ghost -->
              <tr>
                <td><strong>Acción Minimalista</strong><br><small>Iconos en tablas, Cierres</small></td>
                <td><button class="btn btn-ghost">Editar</button></td>
                <td><button class="btn btn-ghost"><span class="material-symbols-outlined">edit</span></button></td>
                <td><button class="btn btn-ghost" disabled><span class="material-symbols-outlined">edit</span></button></td>
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

    </div>
  </div>
</template>

<script setup>
import Navbar from './layout/Navbar.vue';
</script>

<style scoped>
.style-guide-wrapper {
  background-color: var(--color-background);
  min-height: 100vh;
}

.style-guide {
  padding: 3rem var(--spacing-xl);
  max-width: 1200px;
  margin: 0 auto;
  background-color: var(--color-surface);
  border-radius: var(--border-radius-lg);
  box-shadow: var(--box-shadow-lg);
}

section {
  margin-bottom: 4rem;
}

h2 {
  border-bottom: 2px solid var(--color-primary);
  display: inline-block;
  margin-bottom: 1.5rem;
  padding-right: 2rem;
}

.section-intro {
  margin-bottom: 2rem;
  color: var(--color-text-secondary);
}

.color-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: var(--spacing-lg);
}

.color-swatch-wrapper {
  text-align: center;
}

.color-swatch {
  width: 100%;
  height: 80px;
  border-radius: var(--border-radius-md);
  margin-bottom: var(--spacing-sm);
  box-shadow: var(--box-shadow-sm);
  border: 1px solid rgba(0,0,0,0.05);
}

.color-swatch-wrapper p {
  font-weight: 600;
  font-size: 0.85rem;
  margin-bottom: 0.25rem;
}

.color-swatch-wrapper code {
  font-size: 0.75rem;
  background-color: var(--color-background);
  padding: 2px 6px;
  border-radius: var(--border-radius-sm);
}

/* Showcase de Botones */
.showcase-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 1rem;
}

.showcase-table th {
  text-align: left;
  padding: 1rem;
  background: var(--color-background);
  font-size: 0.8rem;
  text-transform: uppercase;
  color: var(--color-text-secondary);
}

.showcase-table td {
  padding: 1.25rem 1rem;
  border-bottom: 1px solid var(--color-border);
}

.showcase-table tr:last-child td {
  border-bottom: none;
}

.sizes-row {
  display: flex;
  align-items: flex-end;
  gap: 2rem;
  padding: 2rem;
  background: var(--color-background);
  border-radius: 12px;
}

.size-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
}

.size-item code {
  font-size: 0.8rem;
}

.mt-8 { margin-top: 2rem; }

.btn-success {
  background-color: #16a34a;
  color: white;
}
.btn-success:hover:not(:disabled) {
  background-color: #15803d;
}
</style>
