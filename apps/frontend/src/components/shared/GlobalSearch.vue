<template>
  <BaseDialog :show="show" @close="$emit('close')" size="lg" :hide-header="true" :initial-focus="searchInputRef" content-class="search-dialog-content">
    <div class="global-search-wrapper">
      <div class="search-input-container">
        <span class="material-symbols-outlined">search</span>
        <input
          ref="searchInputRef"
          v-model="query"
          type="text"
          placeholder="Buscar pedidos, productos, clientes..."
          class="search-input"
          @keydown.down.prevent="navigateResults('down')"
          @keydown.up.prevent="navigateResults('up')"
          @keydown.enter.prevent="selectResult"
        />
      </div>

      <div v-if="isLoading" class="results-container loading">
        <div class="spinner"></div>
        <p>Buscando en TramaTex...</p>
      </div>
      
      <div v-else-if="!query" class="results-container empty">
        <span class="material-symbols-outlined large-icon">manage_search</span>
        <p>Introduzca un término para buscar en todo el sistema.</p>
      </div>

      <div v-else-if="results.length === 0 && !isLoading" class="results-container empty">
        <span class="material-symbols-outlined large-icon">search_off</span>
        <p>No se encontraron resultados para "<strong>{{ query }}</strong>".</p>
      </div>

      <div v-else class="results-container">
        <ul class="results-list">
          <template v-for="(group, index) in groupedResults" :key="index">
            <li class="group-header">{{ group.type }}</li>
            <li
              v-for="item in group.items"
              :key="item.id"
              :class="['result-item', { active: activeResultIndex === getOverallIndex(group.type, item.id) }]"
              @click="goTo(item.url)"
              @mouseenter="activeResultIndex = getOverallIndex(group.type, item.id)"
            >
              <div class="item-icon">
                <span class="material-symbols-outlined">{{ item.icon }}</span>
              </div>
              <div class="item-content">
                <strong>{{ item.title }}</strong>
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

function asArray(payload: any): any[] {
  if (Array.isArray(payload)) return payload
  if (Array.isArray(payload?.data)) return payload.data
  return []
}

async function fetchGlobalResults(searchText: string) {
  const [ordersResp, productsResp, partiesResp] = await Promise.allSettled([
    api.get('/sales/orders', { params: { search: searchText, limit: 5 } }),
    api.get('/products', { params: { search: searchText, page_size: 5 } }),
    api.get('/parties', { params: { name: searchText, page_size: 5 } }),
  ])

  const orders = ordersResp.status === 'fulfilled' ? asArray(ordersResp.value.data) : []
  const products = productsResp.status === 'fulfilled' ? asArray(productsResp.value.data) : []
  const parties = partiesResp.status === 'fulfilled' ? asArray(partiesResp.value.data) : []

  return [
    ...orders.map((o: any) => ({
      id: `order-${o.id}`,
      type: 'order',
      title: o.orderNumber ? `Pedido ${o.orderNumber}` : 'Pedido',
      subtitle: o.status || 'Sin estado',
      url: `/sales/orders/${o.id}`,
    })),
    ...products.map((p: any) => ({
      id: `product-${p.id}`,
      type: 'product',
      title: p.name || p.long_name || 'Producto',
      subtitle: p.sku ? `SKU: ${p.sku}` : 'Sin SKU',
      url: `/products/${p.id}`,
    })),
    ...parties.map((p: any) => ({
      id: `party-${p.id}`,
      type: 'party',
      title: p.name || 'Entidad',
      subtitle: p.tax_id ? `NIF/CIF: ${p.tax_id}` : 'Sin identificación fiscal',
      url: `/parties/${p.id}`,
    })),
  ]
}

watch(query, (newQuery) => {
  clearTimeout(debounceTimer)
  if (!newQuery.trim()) {
    results.value = []
    return
  }

  isLoading.value = true
  debounceTimer = setTimeout(async () => {
    try {
      results.value = await fetchGlobalResults(newQuery)
    } catch (err) {
      console.error('Error en búsqueda global:', err)
      results.value = []
    } finally {
      isLoading.value = false
      activeResultIndex.value = 0
    }
  }, 300)
})


watch(() => props.show, (isShown) => {
  if (isShown) {
    query.value = ''
    results.value = []
    nextTick(() => {
      searchInputRef.value?.focus()
    })
  }
})

const groupedResults = computed(() => {
  const groups: Record<string, { type: string, items: any[] }> = {
    Pedidos: { type: 'Pedidos', items: [] },
    Productos: { type: 'Productos', items: [] },
    Clientes: { type: 'Clientes', items: [] },
  }

  for (const item of results.value) {
    if (item.type === 'order') {
      groups.Pedidos.items.push({
        id: item.id,
        icon: 'shopping_cart',
        title: item.title,
        subtitle: item.subtitle,
        url: item.url,
      })
    } else if (item.type === 'product') {
      groups.Productos.items.push({
        id: item.id,
        icon: 'inventory_2',
        title: item.title,
        subtitle: item.subtitle,
        url: item.url,
      })
    } else if (item.type === 'party') {
      groups.Clientes.items.push({
        id: item.id,
        icon: 'person',
        title: item.title,
        subtitle: item.subtitle,
        url: item.url,
      })
    }
  }

  return Object.values(groups).filter(g => g.items.length > 0)
})


function goTo(url: string) {
  router.push(url)
  emit('close')
}

const flatResults = computed(() => groupedResults.value.flatMap(g => g.items))

function navigateResults(direction: 'up' | 'down') {
  const total = flatResults.value.length
  if (total === 0) return
  if (direction === 'down') {
    activeResultIndex.value = (activeResultIndex.value + 1) % total
  } else {
    activeResultIndex.value = (activeResultIndex.value - 1 + total) % total
  }
}

function getOverallIndex(groupType: string, itemId: string) {
  return flatResults.value.findIndex(item => item.id === itemId)
}

function selectResult() {
  const selected = flatResults.value[activeResultIndex.value]
  if (selected) {
    goTo(selected.url)
  }
}
</script>

<style scoped>
.search-dialog-content {
  padding: 0 !important;
  border-radius: 12px;
  overflow: hidden;
}
.global-search-wrapper {
  display: flex;
  flex-direction: column;
}
.search-input-container {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid var(--color-border);
}
.search-input-container .material-symbols-outlined {
  font-size: 28px;
  color: var(--color-text-secondary);
}
.search-input {
  width: 100%;
  border: none;
  background: transparent;
  font-size: 1.1rem;
  padding: 0.5rem 0;
  color: var(--color-text-primary);
}
.search-input:focus {
  outline: none;
}

.results-container {
  min-height: 400px;
  max-height: 70vh;
  overflow-y: auto;
}
.results-container.empty, .results-container.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 2rem;
  color: var(--color-text-secondary);
}
.large-icon {
  font-size: 4rem;
  opacity: 0.5;
}
.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--color-border);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

.results-list {
  list-style: none;
  padding: 0.75rem;
  margin: 0;
}
.group-header {
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  color: var(--color-text-secondary);
  padding: 1rem 1rem 0.5rem;
}
.result-item {
  display: flex;
  align-items: center;
  padding: 1rem;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.15s;
}
.result-item:hover, .result-item.active {
  background: var(--color-background-soft);
}

.item-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-background);
  color: var(--color-text-secondary);
  margin-right: 1rem;
}
.item-content {
  flex: 1;
}
.item-content strong {
  color: var(--color-text-primary);
  font-size: 0.9rem;
}
.item-content p {
  color: var(--color-text-secondary);
  font-size: 0.8rem;
  margin: 0.1rem 0 0;
}
.item-action {
  color: var(--color-border-strong);
}
.result-item:hover .item-action, .result-item.active .item-action {
  color: var(--color-primary);
}
</style>
