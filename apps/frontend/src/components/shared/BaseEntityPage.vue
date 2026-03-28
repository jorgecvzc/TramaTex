<script setup lang="ts">
/**
 * BaseEntityPage.vue
 * 
 * Plantilla maestra para la gestión integral de una Entidad.
 * Estética: Cabecera fija blanca + Zona de contexto a ancho completo (gris apagado) + Main Content.
 */
</script>

<template>
  <div class="entity-page-container">
    
    <!-- 1. IDENTITY HEADER (Fijo - Blanco Puro) -->
    <div v-if="$slots.header" class="identity-header-sticky">
      <div class="header-max-width">
        <slot name="header"></slot>
      </div>
    </div>

    <!-- 2. SCROLLING AREA -->
    <div class="entity-scroll-area">
      
      <!-- A. CONTEXT HEADER (Ancho completo - Gris Apagado) -->
      <header v-if="$slots.toolbar || $slots.summary || $slots.related" class="entity-context-header">
        <div class="header-max-width">
          <div class="metadata-grid">
            <div v-if="$slots.toolbar" class="metadata-item"><slot name="toolbar"></slot></div>
            <div v-if="$slots.summary" class="metadata-item"><slot name="summary"></slot></div>
            <div v-if="$slots.related" class="metadata-item"><slot name="related"></slot></div>
          </div>
        </div>
      </header>

      <!-- B. MAIN CONTENT (Área de Trabajo - Gris Base) -->
      <div class="entity-body-layout">
        <main class="entity-main-content">
          <slot></slot>
        </main>

        <!-- C. FOOTER -->
        <footer v-if="$slots.footer" class="entity-footer-area">
          <div class="footer-divider"></div>
          <slot name="footer"></slot>
        </footer>
      </div>

    </div>
  </div>
</template>

<style scoped>
.entity-page-container {
  width: 100%;
  min-height: 100vh;
  background-color: var(--color-background);
}

/* --- 1. IDENTITY HEADER --- */
.identity-header-sticky {
  position: sticky;
  top: 60px;
  z-index: 500; /* Reducido para permitir que el Navbar esté por encima */
  background-color: white;
  border-bottom: 1px solid var(--color-border);
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  padding: 0.5rem 0;
}


.header-max-width {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 1rem;
}

/* --- 2. SCROLLING AREA --- */
.entity-scroll-area {
  padding-bottom: 4rem;
}

/* Context Header: Fondo a todo lo ancho, más apagado que el header blanco */
.entity-context-header {
  background-color: #f9fafb; /* Gris muy suave (apagado respecto al blanco) */
  border-bottom: 1px solid var(--color-border);
  padding: 2rem 0;
  margin-bottom: 2.5rem;
}

.metadata-grid {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

/* --- 3. MAIN CONTENT --- */
.entity-body-layout {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 1rem;
}

.entity-main-content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  width: 100%;
}

/* Refuerzo de visibilidad para las tarjetas del Main */
.entity-main-content :deep(.card) {
  border: 1px solid #d1d5db;
  box-shadow: 
    0 4px 6px -1px rgba(0, 0, 0, 0.1), 
    0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

/* --- 4. FOOTER --- */
.footer-divider {
  width: 100%;
  height: 1px;
  background-color: var(--color-border);
  margin: 3rem 0 1rem 0;
  opacity: 0.6;
}

.entity-footer-area {
  padding: 1rem 0;
  color: var(--color-text-secondary);
  font-size: 0.85rem;
}

:deep(.page-header) {
  margin-bottom: 0;
  padding: 0.5rem 0;
}

@media (max-width: 768px) {
  .identity-header-sticky { top: 50px; }
  .header-max-width, .entity-body-layout { padding: 0 0.75rem; }
  .entity-context-header { padding: 1.5rem 0; margin-bottom: 1.5rem; }
}
</style>
