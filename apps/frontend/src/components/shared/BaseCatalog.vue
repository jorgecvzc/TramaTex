<script setup lang="ts">
/**
 * BaseCatalog.vue - Maestro de Listados/Catálogos TramaTex
 * 
 * Sigue la Arquitectura de Tres Capas:
 * 1. Identidad: PageHeader + Acciones Globales
 * 2. Contexto: Filtros y Summary (opcional)
 * 3. Trabajo: Tabla de datos y paginación
 */
import { useRouter } from 'vue-router'
import PageHeader from '@/components/layout/PageHeader.vue'

interface Breadcrumb {
  label: string
  to?: string
}

const props = defineProps<{
  title: string
  breadcrumbs: Breadcrumb[]
  items: any[]
  isLoading?: boolean
  error?: string
  hasFilters?: boolean
  createRoute?: string
  createText?: string
  emptyIcon?: string
  emptyText?: string
  // Control de visibilidad
  hideHeader?: boolean
  hideToolbar?: boolean
}>()

const emit = defineEmits(['refresh', 'clear-filters', 'click-item'])
const router = useRouter()

function handleRowClick(item: any) {
  emit('click-item', item)
}
</script>

<template>
  <div class="catalog-master-layout">
    <!-- CAPA 1: IDENTIDAD (Sticky Header) -->
    <header class="catalog-identity-layer" v-if="!hideHeader">
      <PageHeader :title="title" :breadcrumbs="breadcrumbs">
        <template #actions>
          <slot name="header-actions">
            <button v-if="createRoute" class="btn btn-primary" @click="router.push(createRoute)">
              <span class="material-symbols-outlined">add</span>
              <span>{{ createText || 'Nuevo Registro' }}</span>
            </button>
          </slot>
        </template>
      </PageHeader>
    </header>

    <!-- CAPA 2: CONTEXTO (Filtros y Dashboard de Listado) -->
    <section class="catalog-context-layer" v-if="!hideToolbar">
      <div class="filters-card card compact">
        <div class="catalog-toolbar">
          <div class="filters-scroll-area">
            <div class="filters-grid-inline">
              <slot name="filters"></slot>
            </div>
          </div>
          
          <div class="filter-actions-inline">
            <button class="btn-icon" title="Limpiar filtros" @click="emit('clear-filters')" :disabled="!hasFilters">
              <span class="material-symbols-outlined">filter_alt_off</span>
            </button>
            <button class="btn-icon" title="Refrescar datos" @click="emit('refresh')" :disabled="isLoading">
              <span class="material-symbols-outlined" :class="{ 'spin': isLoading }">refresh</span>
            </button>
          </div>
        </div>
      </div>
      
      <!-- Slot para KPIs o Resúmenes del listado (opcional) -->
      <div v-if="$slots.summary" class="catalog-summary-area mb-6">
        <slot name="summary"></slot>
      </div>
    </section>

    <!-- CAPA 3: TRABAJO (Área de Datos) -->
    <main class="catalog-work-layer">
      <div v-if="error" class="alert-card card error">
        <span class="material-symbols-outlined">error</span>
        <p>{{ error }}</p>
        <button class="btn btn-outline btn-sm" @click="emit('refresh')">Reintentar</button>
      </div>

      <div class="card table-card">
        <div class="table-wrapper" :class="{ 'is-loading': isLoading }">
          <table class="data-table">
            <thead>
              <tr>
                <slot name="table-header"></slot>
              </tr>
            </thead>
            <tbody v-if="items.length > 0">
              <tr 
                v-for="(item, idx) in items" 
                :key="item.id || idx" 
                class="row-clickable"
                @click="handleRowClick(item)"
              >
                <slot name="item" :item="item"></slot>
              </tr>
            </tbody>
          </table>

          <!-- Estados Vacíos y Carga -->
          <div v-if="isLoading && items.length === 0" class="loading-state">
            <div class="spinner"></div>
            <p>Consultando catálogo...</p>
          </div>

          <div v-else-if="items.length === 0" class="empty-state">
            <div class="empty-icon-wrapper">
              <span class="material-symbols-outlined">{{ emptyIcon || 'inventory_2' }}</span>
            </div>
            <h3>{{ emptyText || 'No se han encontrado registros' }}</h3>
            <p v-if="hasFilters">Pruebe a cambiar los criterios de búsqueda o limpie los filtros.</p>
            <button v-if="hasFilters" class="btn btn-nav btn-sm mt-4" @click="emit('clear-filters')">
              Limpiar Filtros
            </button>
          </div>
        </div>
      </div>
      
      <!-- Footer de tabla (Paginación, etc.) -->
      <footer v-if="$slots.pagination" class="catalog-pagination mt-4">
        <slot name="pagination"></slot>
      </footer>
    </main>
  </div>
</template>

<style scoped>
.catalog-master-layout {
  display: flex;
  flex-direction: column;
  gap: 0;
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 1.5rem 3rem 1.5rem;
}

/* Layer 1: Identity */
.catalog-identity-layer {
  position: sticky;
  top: 60px;
  background: white;
  z-index: 500;
  padding-top: 0.75rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid var(--color-border);
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  margin-left: -1.5rem;
  margin-right: -1.5rem;
  margin-bottom: 2rem;
}

.catalog-identity-layer :deep(.page-header) {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 1.5rem;
  margin-bottom: 0;
}

/* Layer 2: Context */
.catalog-context-layer {
  margin-bottom: 1rem;
}

.catalog-toolbar {
  display: flex;
  align-items: flex-end; /* Alinear a la base */
  gap: 1rem;
  width: 100%;
}

.filters-scroll-area {
  flex: 1;
  overflow-x: auto;
  padding-bottom: 4px; /* Espacio para el focus ring de los inputs */
}

.filters-grid-inline {
  display: flex;
  align-items: flex-end; /* Alinear campos a la base */
  gap: 1rem;
  min-width: min-content;
}

.filter-actions-inline {
  display: flex;
  gap: 0.25rem;
  padding-bottom: 6px; /* Ajuste fino para alinear iconos con texto de inputs */
}

.filter-actions-inline .btn-icon {
  width: 34px;
  height: 34px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* Layer 3: Work */
.catalog-work-layer {
  flex: 1;
}

.table-card {
  overflow: hidden;
  border-radius: 12px;
  box-shadow: var(--box-shadow-md);
}

.table-wrapper {
  position: relative;
  min-height: 200px;
}

.table-wrapper.is-loading {
  opacity: 0.6;
  pointer-events: none;
}

/* Utility classes */
.spin { animation: rotate 1s linear infinite; }
@keyframes rotate { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

.alert-card.error {
  display: flex;
  align-items: center;
  gap: 1rem;
  border-left: 4px solid var(--color-danger);
  color: var(--color-danger);
  padding: 1rem 1.5rem;
  margin-bottom: 1.5rem;
}

.empty-state {
  padding: 5rem 2rem;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.empty-icon-wrapper {
  width: 64px;
  height: 64px;
  background: var(--color-background);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 1.5rem;
  color: var(--color-text-secondary);
}

.empty-icon-wrapper .material-symbols-outlined {
  font-size: 32px;
}

.loading-state {
  position: absolute;
  inset: 0;
  background: rgba(255,255,255,0.5);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 10;
}
</style>
