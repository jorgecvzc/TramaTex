<script setup lang="ts">
/**
 * BaseDashboardPage.vue - Plantilla Maestro para Paneles Operativos
 * 
 * Diseñado para herramientas como Terminal de Tickets, Taller o Analytics.
 * Soporta un layout de dos columnas (Main + Sidebar) con scroll independiente.
 */
defineProps<{
  isLoading?: boolean
  sidebarPosition?: 'left' | 'right'
}>()
</script>

<template>
  <div class="dashboard-page-container">
    <!-- CAPA 1: IDENTIDAD -->
    <header class="dashboard-header-sticky">
      <div class="header-content-wrapper">
        <slot name="header"></slot>
      </div>
    </header>

    <!-- CARGA -->
    <div v-if="isLoading" class="dashboard-loading">
      <div class="spinner"></div>
      <p>Iniciando panel operativo...</p>
    </div>

    <!-- CAPA 2 y 3: DASHBOARD LAYOUT -->
    <div v-else :class="['dashboard-body-layout', sidebarPosition || 'right']">
      <!-- ÁREA PRINCIPAL -->
      <main class="dashboard-main-content">
        <slot></slot>
      </main>

      <!-- ÁREA LATERAL (SIDEBAR / CONTEXTO) -->
      <aside v-if="$slots.sidebar" class="dashboard-sidebar">
        <slot name="sidebar"></slot>
      </aside>
    </div>
  </div>
</template>

<style scoped>
.dashboard-page-container {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 64px); /* Restar altura de la Navbar */
  overflow: hidden;
  background-color: var(--color-background);
}

.dashboard-header-sticky {
  background: white;
  border-bottom: 1px solid var(--color-border);
  box-shadow: var(--box-shadow-sm);
  z-index: 100;
  display: flex;
  align-items: stretch;
  min-height: 88px;
}

.header-content-wrapper {
  display: flex;
  align-items: stretch;
  max-width: 1300px;
  margin: 0 auto;
  padding: 0 2rem;
  width: 100%;
}

.dashboard-body-layout {
  display: flex;
  flex: 1;
  overflow: hidden;
  max-width: 1300px;
  width: 100%;
  margin: 0 auto;
}

.dashboard-body-layout.left { flex-direction: row-reverse; }

.dashboard-main-content {
  flex: 1;
  overflow-y: auto;
  padding: 2rem;
  padding-right: 1rem;
}

.dashboard-sidebar {
  width: 320px;
  background: white;
  border-left: 1px solid var(--color-border);
  overflow-y: auto;
  padding: 2rem;
  box-shadow: -4px 0 15px rgba(0,0,0,0.02);
}

.dashboard-body-layout.left .dashboard-sidebar {
  border-left: none;
  border-right: 1px solid var(--color-border);
  box-shadow: 4px 0 15px rgba(0,0,0,0.02);
}

.dashboard-loading {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--color-text-secondary);
}

/* Custom Scrollbars para el modo Dashboard */
.dashboard-main-content::-webkit-scrollbar,
.dashboard-sidebar::-webkit-scrollbar {
  width: 6px;
}

.dashboard-main-content::-webkit-scrollbar-thumb,
.dashboard-sidebar::-webkit-scrollbar-thumb {
  background: var(--color-border);
  border-radius: 10px;
}

@media (max-width: 1200px) {
  .dashboard-body-layout { flex-direction: column; overflow-y: auto; }
  .dashboard-main-content, .dashboard-sidebar { width: 100%; overflow: visible; padding: 1rem; }
  .dashboard-sidebar { border-left: none; border-top: 1px solid var(--color-border); }
}
</style>
