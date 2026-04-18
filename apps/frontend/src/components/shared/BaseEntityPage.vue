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
  top: 76px; /* Ajustado a la Navbar superior */
  z-index: 500;
  background-color: white;
  border-bottom: 1px solid var(--color-border);
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  padding: 0.5rem 0;
}


.header-max-width {
  max-width: 1300px;
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
  max-width: 1300px;
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

/* --- 5. HEREDITARY STYLES FOR SLOTS (ESTÁNDAR VISUAL ÚNICO) --- */

/* Contenedores de filas de etiquetas */
:deep(.overview-tags-row), :deep(.related-history-grid) {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

:deep(.related-history-grid) {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
}

/* Caja de Etiqueta (Base común para Resumen y Trazabilidad) */
:deep(.summary-tag), :deep(.related-tag-card) {
  flex: 1;
  min-width: 260px; /* Tamaño estándar único */
  padding: 0.75rem 1.25rem;
  background: white;
  border: 1px solid var(--color-border);
  border-radius: 12px;
  display: flex;
  align-items: center;
  gap: 1rem;
  box-shadow: var(--box-shadow-sm);
  transition: all 0.2s ease;
}

/* Variación específica para Trazabilidad (Enlaces navegables) */
:deep(.related-tag-card) {
  border-left: 4px solid var(--color-secondary);
  text-decoration: none;
  cursor: pointer;
}
:deep(.related-tag-card.highlight-info) { border-left-color: #2563eb; }
:deep(.related-tag-card:hover) {
  transform: translateX(2px) translateY(-1px);
  box-shadow: var(--box-shadow-md);
  border-color: var(--color-secondary);
}

/* Iconos dentro de las etiquetas */
:deep(.summary-tag .icon), :deep(.related-tag-card .tag-icon) {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: rgba(0,0,0,0.03);
}

:deep(.icon.blue) { background: rgba(59, 130, 246, 0.1); color: #2563eb; }
:deep(.icon.yellow), :deep(.tag-icon.yellow) { background: rgba(230, 184, 0, 0.1); color: #d97706; }
:deep(.icon.purple) { background: rgba(168, 85, 247, 0.1); color: #9333ea; }
:deep(.icon.green), :deep(.tag-icon.success) { background: rgba(34, 197, 94, 0.1); color: #16a34a; }

/* Contenido de texto */
:deep(.tag-content) {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
  line-height: 1.2;
}

:deep(.tag-content label) {
  font-size: 0.65rem;
  font-weight: 800;
  text-transform: uppercase;
  color: var(--color-text-secondary);
  letter-spacing: 0.05em;
}

:deep(.tag-content strong) {
  font-size: 1rem;
  color: var(--color-text-primary);
  font-weight: 700;
}

:deep(.amount) {
  color: var(--color-success) !important;
  font-size: 1.15rem !important;
  font-family: var(--font-family-mono);
}

:deep(.jump-icon) {
  font-size: 18px;
  color: var(--color-text-secondary);
  opacity: 0.5;
  margin-left: auto;
}

/* --- 6. HEREDITARY STYLES: ECONOMIC SUMMARY (TOTALS) --- */
:deep(.totals-summary-container) {
  display: flex;
  justify-content: flex-end;
  margin-top: 0.5rem;
}

:deep(.totals-summary-card) {
  width: 100%;
  max-width: 400px;
  background: white;
  border: 1px solid var(--color-border-strong);
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: var(--box-shadow-md);
  position: relative;
}

:deep(.summary-row) {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 0;
  font-size: 0.95rem;
  color: var(--color-text-primary);
}

:deep(.summary-row label) {
  color: var(--color-text-secondary);
  font-weight: 600;
  margin: 0;
}

:deep(.summary-row.grand-total) {
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 2px solid var(--color-border);
  font-weight: 800;
  font-size: 1.3rem;
}

:deep(.summary-row.grand-total span) {
  color: var(--color-secondary);
}

:deep(.base-page-header) {
  margin-bottom: 0;
  padding: 0.5rem 0;
}

@media (max-width: 768px) {
  .identity-header-sticky { top: 50px; }
  .header-max-width, .entity-body-layout { padding: 0 0.75rem; }
  .entity-context-header { padding: 1.5rem 0; margin-bottom: 1.5rem; }
}
</style>
