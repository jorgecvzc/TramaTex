<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { RouterLink } from 'vue-router'
import PageHeader from '@/components/layout/PageHeader.vue'

const props = defineProps<{
  title: string
  icon?: string
  breadcrumbs: { label: string; to?: string }[]
  items: any[]
  isLoading?: boolean
  error?: string
  hasFilters?: boolean
  createText?: string
  createRoute?: string
  emptyIcon?: string
  emptyText?: string
}>()

const emit = defineEmits(['clear-filters', 'refresh', 'click-item'])

function handleRowClick(item: any) {
  emit('click-item', item)
}
</script>

<template>
  <div class="catalog-page-container">
    <!-- CAPA 1: IDENTIDAD (STICKY) -->
    <div class="catalog-header-shell">
      <div class="catalog-header">
        <PageHeader :title="title" :breadcrumbs="breadcrumbs">
          <template #icon><span class="material-symbols-outlined">{{ icon || 'list_alt' }}</span></template>
          <template #actions>
            <div class="header-actions">
              <button v-if="hasFilters" class="btn btn-outline btn-sm" @click="$emit('clear-filters')">
                <span class="material-symbols-outlined">filter_alt_off</span> Limpiar
              </button>
              <button class="btn btn-outline btn-sm" @click="$emit('refresh')" :disabled="isLoading">
                <span class="material-symbols-outlined" :class="{ 'spin': isLoading }">refresh</span> Actualizar
              </button>
              <RouterLink v-if="createRoute" :to="createRoute" class="btn btn-primary btn-sm">
                <span class="material-symbols-outlined">add</span> {{ createText || 'Nuevo' }}
              </RouterLink>
              <slot name="header-actions"></slot>
            </div>
          </template>
        </PageHeader>
      </div>
    </div>

    <div class="catalog-content">
      <!-- CAPA 2: CONTEXTO (FILTROS) -->
      <div v-if="$slots.filters" class="card filters-card">
        <slot name="filters"></slot>
      </div>

      <!-- CAPA 3: TRABAJO (TABLA) -->
      <div class="card table-card">
        <div v-if="isLoading && items.length === 0" class="loading-state">
          <div class="spinner"></div>
          <p>Cargando datos...</p>
        </div>

        <div v-else-if="error" class="error-state">
          <span class="material-symbols-outlined">error</span>
          <p>{{ error }}</p>
          <button class="btn btn-outline btn-sm mt-4" @click="$emit('refresh')">Reintentar</button>
        </div>

        <div v-else-if="items.length === 0" class="empty-state">
          <span class="material-symbols-outlined">{{ emptyIcon || 'search_off' }}</span>
          <p>{{ emptyText || 'No se han encontrado resultados' }}</p>
        </div>

        <div v-else class="table-wrapper">
          <table class="data-table">
            <thead>
              <tr>
                <slot name="table-header"></slot>
              </tr>
            </thead>
            <tbody>
              <tr 
                v-for="(item, index) in items" 
                :key="item.id || index" 
                class="row-clickable" 
                @click="handleRowClick(item)"
              >
                <slot name="item" :item="item"></slot>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.catalog-page-container {
  width: 100%;
  min-height: 100vh;
  background-color: var(--color-background);
  display: flex;
  flex-direction: column;
}

.catalog-header-shell {
  background: white;
  border-bottom: 1px solid var(--color-border);
  position: sticky;
  top: 76px; /* Sincronizado con Navbar superior */
  z-index: 900;
}

.catalog-header {
  max-width: 1300px;
  margin: 0 auto;
  padding: 0.5rem 1rem;
}

.catalog-content {
  flex: 1;
  padding: 1rem;
  max-width: 1300px;
  width: 100%;
  margin: 0 auto;
}

.header-actions { display: flex; gap: 0.5rem; align-items: center; }

.filters-card {
  padding: 0.75rem 1rem;
  margin-bottom: 1rem;
  display: flex;
  gap: 1.5rem;
  align-items: flex-end;
  flex-wrap: wrap;
}

/* Los filtros ahora heredan de _forms.css automáticamente */
:deep(.filter-group) {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.table-card { overflow: hidden; }
.table-wrapper { overflow-x: auto; }

.data-table { width: 100%; border-collapse: collapse; }
.data-table th, .data-table td {
  padding: 0.75rem 1rem;
  text-align: left;
  border-bottom: 1px solid var(--color-border-soft);
  vertical-align: middle;
}

.data-table th {
  background: var(--color-background-soft);
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.row-clickable { cursor: pointer; transition: background-color 0.15s; }
.row-clickable:hover { background-color: var(--color-background-soft); }

.empty-state, .loading-state, .error-state {
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  padding: 4rem 2rem; text-align: center; color: var(--color-text-secondary);
}
.empty-state .material-symbols-outlined, .error-state .material-symbols-outlined {
  font-size: 3rem; opacity: 0.3; margin-bottom: 1rem;
}

.spinner {
  width: 32px; height: 32px; border: 3px solid var(--color-border);
  border-top-color: var(--color-primary); border-radius: 50%;
  animation: spin 0.8s linear infinite; margin-bottom: 1rem;
}
@keyframes spin { to { transform: rotate(360deg); } }
.spin { animation: spin 1s linear infinite; }
</style>
