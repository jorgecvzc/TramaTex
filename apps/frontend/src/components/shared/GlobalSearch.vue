<template>
  <BaseDialog :show="show" @close="$emit('close')" size="lg" :hide-header="true" :initial-focus="searchInputRef" content-class="search-dialog-content">
    <div class="global-search-wrapper">
      <div class="search-input-container">
        <span class="material-symbols-outlined">search</span>
        <input
          ref="searchInputRef"
          v-model="query"
          type="text"
          placeholder="Buscar pedidos, facturas, clientes, productos..."
          class="search-input"
          @keydown.down.prevent="navigateResults('down')"
          @keydown.up.prevent="navigateResults('up')"
          @keydown.enter.prevent="selectResult"
        />
        <div class="search-hint">
          <kbd>ESC</kbd> para cerrar
        </div>
      </div>

      <div v-if="isLoading" class="results-container loading">
        <div class="spinner"></div>
        <p>Buscando en TramaTex...</p>
      </div>
      
      <div v-else-if="!query" class="results-container empty">
        <span class="material-symbols-outlined large-icon">manage_search</span>
        <p>Introduce un término para buscar en todo el sistema.</p>
        <small class="text-muted">Busca por Nº de documento, SKU o Nombre de Cliente</small>
      </div>

      <div v-else-if="results.length === 0 && !isLoading" class="results-container empty">
        <span class="material-symbols-outlined large-icon">search_off</span>
        <p>No se encontraron resultados para "<strong>{{ query }}</strong>".</p>
      </div>

      <div v-else class="results-container">
        <ul class="results-list">
          <template v-for="(group, gName) in groupedResults" :key="gName">
            <li class="group-header">{{ gName }}</li>
            <li
              v-for="item in group"
              :key="item.id"
              :class="['result-item', { active: activeResultIndex === getOverallIndex(item.id) }]"
              @click="goTo(item.url)"
              @mouseenter="activeResultIndex = getOverallIndex(item.id)"
            >
              <div class="item-icon" :class="item.type">
                <span class="material-symbols-outlined">{{ getIcon(item.type) }}</span>
              </div>
              <div class="item-content">
                <div class="title-row">
                  <strong>{{ item.title }}</strong>
                  <span v-if="item.type === 'invoice'" class="badge-type">Factura</span>
                </div>
                <p>{{ item.subtitle }}</p>
              </div>
              <div class="item-action">
                <span class="material-symbols-outlined">arrow_forward</span>
              </div>
            </li>
          </template>
        </ul>
      </div>
    </div>
  </BaseDialog>
</template>

<script setup lang="ts">
import { ref, watch, computed, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import BaseDialog from './BaseDialog.vue'
import { api } from '@/services/api'

const props = defineProps({ show: Boolean })
const emit = defineEmits(['close'])

const router = useRouter()
const query = ref('')
const results = ref<any[]>([])
const isLoading = ref(false)
const searchInputRef = ref<HTMLInputElement | null>(null)
const activeResultIndex = ref(0)

let debounceTimer: any

async function fetchGlobalResults(searchText: string) {
  try {
    const resp = await api.get('/search', { params: { q: searchText } })
    return resp.data || []
  } catch (err) {
    console.error('Error backend search:', err)
    return []
  }
}

watch(query, (newQuery) => {
  clearTimeout(debounceTimer)
  if (!newQuery.trim()) {
    results.value = []
    return
  }

  isLoading.value = true
  debounceTimer = setTimeout(async () => {
    results.value = await fetchGlobalResults(newQuery)
    isLoading.value = false
    activeResultIndex.value = 0
  }, 300)
})

watch(() => props.show, (isShown) => {
  if (isShown) {
    query.value = ''
    results.value = []
    nextTick(() => searchInputRef.value?.focus())
  }
})

const groupedResults = computed(() => {
  const groups: Record<string, any[]> = {}
  results.value.forEach(item => {
    const module = item.module || 'Otros'
    if (!groups[module]) groups[module] = []
    groups[module].push(item)
  })
  return groups
})

const flatResults = computed(() => results.value)

function getOverallIndex(itemId: string) {
  return flatResults.value.findIndex(item => item.id === itemId)
}

function getIcon(type: string) {
  const icons: Record<string, string> = {
    order: 'shopping_cart',
    quote: 'request_quote',
    invoice: 'receipt_long',
    delivery_note: 'local_shipping',
    product: 'inventory_2',
    party: 'person',
    mes_work: 'precision_manufacturing'
  }
  return icons[type] || 'description'
}

function goTo(url: string) {
  router.push(url)
  emit('close')
}

function navigateResults(direction: 'up' | 'down') {
  const total = flatResults.value.length
  if (total === 0) return
  if (direction === 'down') {
    activeResultIndex.value = (activeResultIndex.value + 1) % total
  } else {
    activeResultIndex.value = (activeResultIndex.value - 1 + total) % total
  }
}

function selectResult() {
  const selected = flatResults.value[activeResultIndex.value]
  if (selected) goTo(selected.url)
}
</script>

<style scoped>
.search-dialog-content { padding: 0 !important; border-radius: 12px; overflow: hidden; }
.global-search-wrapper { display: flex; flex-direction: column; background: white; }
.search-input-container { display: flex; align-items: center; gap: 1rem; padding: 1.25rem 1.5rem; border-bottom: 1px solid var(--color-border); }
.search-input-container .material-symbols-outlined { font-size: 28px; color: var(--color-secondary); }
.search-input { width: 100%; border: none; background: transparent; font-size: 1.2rem; font-weight: 500; color: var(--color-text-primary); }
.search-input:focus { outline: none; }
.search-hint { font-size: 0.7rem; color: var(--color-text-secondary); white-space: nowrap; }
.search-hint kbd { background: var(--color-background); padding: 0.2rem 0.4rem; border-radius: 4px; border: 1px solid var(--color-border); }

.results-container { min-height: 350px; max-height: 65vh; overflow-y: auto; background: var(--color-background-soft); }
.results-container.empty, .results-container.loading { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 4rem 2rem; color: var(--color-text-secondary); }
.large-icon { font-size: 4rem; opacity: 0.2; margin-bottom: 1rem; }

.results-list { list-style: none; padding: 0.5rem; margin: 0; }
.group-header { font-size: 0.7rem; font-weight: 800; text-transform: uppercase; color: var(--color-text-secondary); padding: 1.25rem 1rem 0.5rem; letter-spacing: 0.05em; }
.result-item { display: flex; align-items: center; padding: 0.75rem 1rem; border-radius: 10px; cursor: pointer; transition: 0.2s; margin-bottom: 2px; }
.result-item:hover, .result-item.active { background: white; box-shadow: var(--box-shadow-md); transform: translateX(4px); }

.item-icon { width: 42px; height: 42px; border-radius: 10px; display: flex; align-items: center; justify-content: center; background: white; color: var(--color-text-secondary); margin-right: 1.25rem; border: 1px solid var(--color-border); }
.result-item.active .item-icon { border-color: var(--color-primary); color: var(--color-primary); }

.item-content { flex: 1; min-width: 0; }
.title-row { display: flex; align-items: center; gap: 0.75rem; }
.item-content strong { color: var(--color-text-primary); font-size: 0.95rem; }
.item-content p { color: var(--color-text-secondary); font-size: 0.8rem; margin: 0.15rem 0 0; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

.badge-type { font-size: 0.6rem; font-weight: 800; text-transform: uppercase; padding: 0.15rem 0.4rem; background: var(--color-background); border-radius: 4px; color: var(--color-secondary); }

.item-action { color: var(--color-border); opacity: 0; transition: 0.2s; }
.result-item:hover .item-action, .result-item.active .item-action { opacity: 1; color: var(--color-primary); }

.spinner { width: 32px; height: 32px; border: 3px solid var(--color-border); border-top-color: var(--color-primary); border-radius: 50%; animation: spin 0.8s linear infinite; margin-bottom: 1rem; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
