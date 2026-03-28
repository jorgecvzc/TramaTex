<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import PageHeader from '@/components/layout/PageHeader.vue'

const props = defineProps<{
  title: string
  icon?: string
  breadcrumbs: { label: string; to: string }[]
  items: any[]
  isLoading: boolean
  error?: string
  createText?: string
  createRoute?: string
  emptyText?: string
  emptyIcon?: string
  hasFilters?: boolean
}>()

const emit = defineEmits(['refresh', 'clear-filters', 'click-item'])

const activeIndex = ref(-1)

function handleKeydown(e: KeyboardEvent) {
  if (props.items && props.items.length > 0) {
    if (e.key === 'ArrowDown') {
      e.preventDefault()
      activeIndex.value = (activeIndex.value + 1) % props.items.length
    } else if (e.key === 'ArrowUp') {
      e.preventDefault()
      activeIndex.value = (activeIndex.value - 1 + props.items.length) % props.items.length
    } else if (e.key === 'Enter') {
      e.preventDefault()
      if (activeIndex.value >= 0 && activeIndex.value < props.items.length) {
        emit('click-item', props.items[activeIndex.value])
      }
    }
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<template>
  <div class="base-catalog-layout">
    <PageHeader 
      :title="title" 
      :icon="icon"
      :breadcrumbs="breadcrumbs"
      class="catalog-header"
    >
      <template #actions>
        <div class="header-actions">
          <button v-if="hasFilters" class="btn btn-outline btn-sm" @click="$emit('clear-filters')">
            Limpiar Filtros
          </button>
          <slot name="header-actions">
            <RouterLink v-if="createRoute" :to="createRoute" class="btn btn-primary">
              <span class="material-symbols-outlined">add</span>
              <span>{{ createText || 'Nuevo' }}</span>
            </RouterLink>
          </slot>
        </div>
      </template>
    </PageHeader>
    
    <div class="catalog-content">
      <div class="card filters-card" v-if="$slots.filters">
        <slot name="filters"></slot>
      </div>

      <div class="card table-card">
        <div v-if="isLoading" class="loading-state">
          <div class="spinner"></div>
          <p>Cargando datos...</p>
        </div>
        <div v-else-if="error" class="error-state">
          <span class="material-symbols-outlined">error</span>
          <h3>Error al cargar</h3>
          <p>{{ error }}</p>
          <button class="btn btn-sm btn-outline" @click="$emit('refresh')">Reintentar</button>
        </div>
        <div v-else-if="items.length === 0" class="empty-state">
          <span class="material-symbols-outlined">{{ emptyIcon || 'folder_off' }}</span>
          <p>{{ emptyText || 'No hay elementos que mostrar.' }}</p>
          <RouterLink v-if="createRoute" :to="createRoute" class="btn btn-primary btn-sm mt-4">
            <span class="material-symbols-outlined">add</span>
            <span>{{ createText || 'Crear el primero' }}</span>
          </RouterLink>
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
                :key="item.id" 
                class="row-clickable" 
                :class="{ 'row-active': index === activeIndex }"
                @click="$emit('click-item', item)"
                @mouseenter="activeIndex = index"
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
.base-catalog-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: var(--color-background);
}
.catalog-header {
  background: white;
  border-bottom: 1px solid var(--color-border);
  position: sticky;
  top: 65px; /* Altura de la Navbar */
  z-index: 900;
  padding-bottom: 1.5rem !important;
  margin-bottom: 0 !important;
}
.catalog-content {
  flex: 1;
  padding: 2rem;
  max-width: 1400px;
  width: 100%;
  margin: 0 auto;
  overflow-y: auto;
}
.header-actions {
  display: flex;
  gap: 1rem;
}
.card {
  background: white;
  border-radius: 12px;
  border: 1px solid var(--color-border);
  box-shadow: var(--box-shadow-sm);
}
.filters-card {
  padding: 1.5rem;
  margin-bottom: 1.5rem;
  display: flex;
  gap: 1.5rem;
  flex-wrap: wrap;
}
.filter-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}
.filter-group label {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--color-text-secondary);
}
.filter-group input, .filter-group select {
  min-width: 220px;
}
.table-card {
  overflow: hidden;
}
.table-wrapper {
  overflow-x: auto;
}
.data-table {
  width: 100%;
  border-collapse: collapse;
}
.data-table th, .data-table td {
  padding: 1rem 1.25rem;
  text-align: left;
  border-bottom: 1px solid var(--color-border-soft);
  vertical-align: middle;
}
.data-table th {
  background: var(--color-background);
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--color-text-secondary);
  text-transform: uppercase;
}
.row-clickable {
  cursor: pointer;
  transition: background-color 0.15s ease-in-out;
}
.row-clickable:hover {
  background-color: var(--color-background-soft);
}
.row-active {
  background-color: var(--color-primary-soft) !important;
  outline: 2px solid var(--color-primary);
  outline-offset: -2px;
}
.empty-state, .loading-state, .error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  text-align: center;
  color: var(--color-text-secondary);
}
.empty-state .material-symbols-outlined, .error-state .material-symbols-outlined {
  font-size: 3rem;
  opacity: 0.5;
  margin-bottom: 1rem;
}
.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--color-border);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 1rem;
}
@keyframes spin { to { transform: rotate(360deg); } }
</style>
