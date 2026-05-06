<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { RouterLink } from 'vue-router'
import { 
  FilterX, 
  RefreshCw, 
  Plus, 
  AlertCircle
} from 'lucide-vue-next'
import { getIcon } from '@/utils/icons'
import BasePageHeader from '@/components/shared/BasePageHeader.vue'
import BaseSkeleton from '@/components/shared/BaseSkeleton.vue'

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
  skeletonRows?: number
  skeletonColumns?: number
}>()

const emit = defineEmits(['clear-filters', 'refresh', 'click-item'])

const selectedIndex = ref(-1)
const tableBodyRef = ref<HTMLElement | null>(null)

function handleRowClick(item: any, index: number) {
  selectedIndex.value = index
  emit('click-item', item)
}

function handleKeyDown(e: KeyboardEvent) {
  // Ignore if the user is typing in an input or textarea
  const target = e.target as HTMLElement
  if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable) {
    return
  }

  if (props.items.length === 0) return

  switch (e.key) {
    case 'ArrowDown':
      e.preventDefault()
      selectedIndex.value = Math.min(selectedIndex.value + 1, props.items.length - 1)
      scrollToSelected()
      break
    case 'ArrowUp':
      e.preventDefault()
      selectedIndex.value = Math.max(selectedIndex.value - 1, 0)
      scrollToSelected()
      break
    case 'Enter':
      if (selectedIndex.value >= 0 && selectedIndex.value < props.items.length) {
        e.preventDefault()
        handleRowClick(props.items[selectedIndex.value])
      }
      break
    case 'Escape':
      selectedIndex.value = -1
      break
  }
}

function scrollToSelected() {
  if (selectedIndex.value < 0 || !tableBodyRef.value) return
  
  const rows = tableBodyRef.value.querySelectorAll('tr')
  const selectedRow = rows[selectedIndex.value] as HTMLElement
  
  if (selectedRow) {
    selectedRow.scrollIntoView({
      block: 'nearest',
      behavior: 'smooth'
    })
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
})

const resolvedIcon = computed(() => getIcon(props.icon || 'list_alt'))
const resolvedEmptyIcon = computed(() => getIcon(props.emptyIcon || 'search_off'))
</script>

<template>
  <div class="catalog-page-container">
    <!-- LAYER 1: IDENTITY (STICKY) -->
    <div class="catalog-header-shell">
      <div class="catalog-header">
        <BasePageHeader :title="title" :breadcrumbs="breadcrumbs">
          <template #icon>
            <component :is="resolvedIcon" :size="28" />
          </template>
          <template #actions>
            <div class="header-actions">
              <div class="navigation-hint no-mobile">
                <kbd>↑</kbd> <kbd>↓</kbd> Navegar | <kbd>Enter</kbd> Seleccionar
              </div>
              <button v-if="hasFilters" class="btn btn-outline btn-sm" @click="$emit('clear-filters')">
                <FilterX :size="16" /> Limpiar
              </button>
              <button class="btn btn-outline btn-sm" @click="$emit('refresh')" :disabled="isLoading">
                <RefreshCw :size="16" :class="{ 'spin': isLoading }" /> Actualizar
              </button>
              <RouterLink v-if="createRoute" :to="createRoute" class="btn btn-primary btn-sm">
                <Plus :size="16" /> {{ createText || 'Nuevo' }}
              </RouterLink>
              <slot name="header-actions"></slot>
            </div>
          </template>
        </BasePageHeader>
      </div>
    </div>

    <div class="catalog-content">
      <!-- LAYER 2: CONTEXT (FILTERS) -->
      <div v-if="$slots.filters" class="card filters-card">
        <slot name="filters"></slot>
      </div>

      <!-- LAYER 3: WORK (TABLE) -->
      <div class="card table-card">
        <div v-if="isLoading && items.length === 0" class="skeleton-table-wrapper">
          <table class="data-table">
            <thead>
              <tr>
                <th v-for="c in (skeletonColumns || 4)" :key="c">
                  <BaseSkeleton type="text" width="60%" />
                </th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="r in (skeletonRows || 5)" :key="r">
                <td v-for="c in (skeletonColumns || 4)" :key="c">
                  <BaseSkeleton type="row" :width="Math.floor(Math.random() * (90 - 40 + 1) + 40) + '%'" />
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div v-else-if="error" class="error-state">
          <AlertCircle :size="48" class="state-icon" />
          <p>{{ error }}</p>
          <button class="btn btn-outline btn-sm mt-4" @click="$emit('refresh')">Reintentar</button>
        </div>

        <div v-else-if="items.length === 0" class="empty-state">
          <component :is="resolvedEmptyIcon" :size="48" class="state-icon" />
          <p>{{ emptyText || 'No se han encontrado resultados' }}</p>
        </div>

        <div v-else class="table-wrapper">
          <table class="data-table">
            <thead>
              <tr>
                <slot name="table-header"></slot>
              </tr>
            </thead>
            <tbody ref="tableBodyRef">
              <tr 
                v-for="(item, index) in items" 
                :key="item.id || index" 
                class="row-clickable" 
                :class="{ 'is-selected': selectedIndex === index }"
                @click="handleRowClick(item, index)"
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
  top: 76px; /* Synchronized with top Navbar */
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

.navigation-hint {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  padding-right: 1rem;
  border-right: 1px solid var(--color-border);
  margin-right: 0.5rem;
}

.navigation-hint kbd {
  margin: 0 0.1rem;
}

@media (max-width: 1024px) {
  .no-mobile { display: none !important; }
}

.filters-card {
  padding: 0.75rem 1rem;
  margin-bottom: 1rem;
  display: flex;
  gap: 1.5rem;
  align-items: flex-end;
  flex-wrap: wrap;
}

/* Filters now inherit from _forms.css automatically */
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

.row-clickable { cursor: pointer; transition: background-color 0.15s; position: relative; }
.row-clickable:hover { background-color: var(--color-background-soft); }

.row-clickable.is-selected {
  background-color: rgba(230, 184, 0, 0.12) !important;
  outline: 2px solid var(--color-primary);
  outline-offset: -2px;
  z-index: 1;
}

.row-clickable.is-selected::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 4px;
  background-color: var(--color-primary);
}

.empty-state, .loading-state, .error-state {
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  padding: 4rem 2rem; text-align: center; color: var(--color-text-secondary);
}

.state-icon {
  opacity: 0.3;
  margin-bottom: 1rem;
}

.spinner {
  width: 32px; height: 32px; border: 3px solid var(--color-border);
  border-top-color: var(--color-primary); border-radius: 50%;
  animation: spin 0.8s linear infinite; margin-bottom: 1rem;
}
@keyframes spin { to { transform: rotate(360deg); } }
.spin { animation: spin 1s linear infinite; }
</style>
