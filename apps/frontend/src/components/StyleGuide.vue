<template>
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
</template>

<style scoped>
.style-guide {
  padding: var(--spacing-xl);
  background-color: var(--color-surface);
  margin: var(--spacing-xl);
  border-radius: var(--border-radius-lg);
  box-shadow: var(--box-shadow-lg);
}

section {
  margin-bottom: var(--spacing-xxl);
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
  height: 100px;
  border-radius: var(--border-radius-md);
  margin-bottom: var(--spacing-sm);
  box-shadow: var(--box-shadow-sm);
}

.color-swatch-wrapper p {
  font-weight: var(--font-weight-medium);
  margin-bottom: var(--spacing-xs);
}

.color-swatch-wrapper code {
  font-size: var(--font-size-sm);
  background-color: var(--color-background);
  padding: 2px 6px;
  border-radius: var(--border-radius-sm);
}
</style>
