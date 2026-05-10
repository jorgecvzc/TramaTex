<script setup lang="ts">
import BaseSkeleton from '@/components/shared/BaseSkeleton.vue'

/**
 * BaseDashboardPage.vue - Master Template for Operating Panels
 * 
 * Designed for tools such as Ticket Terminal, Workshop or Analytics.
 * Supports a two-column layout (Main + Sidebar) with independent scroll.
 */
defineProps<{
  isLoading?: boolean
  sidebarPosition?: 'left' | 'right'
}>()
</script>

<template>
  <div class="dashboard-page-container">
    <!-- LAYER 1: IDENTITY -->
    <header class="dashboard-header-sticky">
      <div class="header-content-wrapper">
        <slot name="header"></slot>
      </div>
    </header>

    <!-- LOADING -->
    <div v-if="isLoading" class="dashboard-loading">
      <div class="skeleton-container">
        <BaseSkeleton type="title" width="300px" height="40px" />
        <div class="skeleton-grid">
          <BaseSkeleton v-for="i in 4" :key="i" type="row" height="100px" />
        </div>
      </div>
      <p>Iniciando panel operativo...</p>
    </div>

    <!-- LAYER 2 and 3: DASHBOARD LAYOUT -->
    <div v-else :class="['dashboard-body-layout', sidebarPosition || 'right']">
      <!-- MAIN AREA -->
      <main class="dashboard-main-content">
        <slot></slot>
      </main>

      <!-- SIDE AREA (SIDEBAR / CONTEXT) -->
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
  height: calc(100vh - 76px); /* Subtract Navbar height (76px) */
  overflow: hidden;
  background-color: var(--color-background);
}

.dashboard-header-sticky {
  background: white;
  border-bottom: 1px solid var(--color-border);
  box-shadow: var(--box-shadow-sm);
  z-index: 100;
  display: flex;
  align-items: center;
  min-height: 88px;
  position: sticky;
  top: 0;
  flex-shrink: 0;
}

.header-content-wrapper {
  display: flex;
  align-items: center;
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
  gap: var(--spacing-md);
  padding: 2rem;
}

.skeleton-container {
  width: 100%;
  max-width: 800px;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
  align-items: center;
}

.skeleton-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: var(--spacing-md);
  width: 100%;
}

/* Custom Scrollbars for Dashboard mode */
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
