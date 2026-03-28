<script setup lang="ts">
/**
 * ProductList.vue - Listado Maestro de Productos
 * 
 * Implementa el estándar BaseCatalog con Arquitectura de 3 Capas.
 */
import { ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { productApi } from '@/services/productApi'
import BaseCatalog from '@/components/shared/BaseCatalog.vue'

const router = useRouter()
const products = ref<any[]>([])
const brands = ref<any[]>([])
const productGroups = ref<any[]>([])
const isLoading = ref(false)
const error = ref('')

const filters = reactive({ 
  search: '', 
  brandId: '', 
  groupId: '', 
  isActive: '' 
})

const hasFilters = computed(() => 
  filters.search.trim() !== '' || filters.brandId !== '' || filters.groupId !== '' || filters.isActive !== ''
)

// Lógica de filtrado con debounce
let debounceTimer: any = null
watch(filters, () => {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => fetchProducts(), 350)
}, { deep: true })

async function fetchProducts() {
  isLoading.value = true
  error.value = ''
  try {
    const res = await productApi.listProducts({ 
      searchText: filters.search,
      brandId: filters.brandId,
      groupId: filters.groupId,
      isActive: filters.isActive === '' ? undefined : filters.isActive === 'true',
      pageSize: 100
    })
    products.value = res.data || (Array.isArray(res) ? res : [])
  } catch (err: any) { 
    error.value = 'No se han podido cargar los productos.'
    console.error(err)
  } finally { 
    isLoading.value = false 
  }
}

async function loadMasters() {
  try {
    const [bRes, gRes] = await Promise.all([
      productApi.listBrands({ isActive: true }),
      productApi.listProductGroups({ isActive: true })
    ])
    brands.value = bRes.data || []
    productGroups.value = gRes.data || []
  } catch (err) {}
}

function clearFilters() { 
  filters.search = ''
  filters.brandId = ''
  filters.groupId = ''
  filters.isActive = ''
}

function navigateToDetail(product: any) { 
  router.push(`/products/${product.id}`) 
}

async function toggleStatus(product: any) {
  const newStatus = !product.is_active
  try { 
    await productApi.changeProductStatus(product.id, newStatus)
    product.is_active = newStatus 
  } catch (err: any) {
    alert('Error al cambiar estado: ' + err.message)
  }
}

function getBrandName(id: string) { 
  return brands.value.find(b => b.id === id)?.name || '—' 
}

function formatPrice(v: number) { 
  return new Intl.NumberFormat('es-ES', { style: 'currency', currency: 'EUR' }).format(v || 0) 
}

onMounted(() => {
  fetchProducts()
  loadMasters()
})

onUnmounted(() => { if (debounceTimer) clearTimeout(debounceTimer) })
</script>

<template>
  <BaseCatalog
    title="Catálogo de Productos"
    :breadcrumbs="[{ label: 'Ventas', to: '/sales/orders' }, { label: 'Productos' }]"
    :items="products"
    :is-loading="isLoading"
    :error="error"
    :has-filters="hasFilters"
    create-route="/products/new"
    create-text="Nuevo Producto"
    empty-icon="inventory_2"
    empty-text="No hay productos registrados"
    @clear-filters="clearFilters"
    @refresh="fetchProducts"
    @click-item="navigateToDetail"
  >
    <!-- CAPA 2: CONTEXTO (Filtros) -->
    <template #filters>
      <div class="filter-group">
        <label>Búsqueda</label>
        <input v-model="filters.search" type="text" placeholder="Nombre o SKU..." />
      </div>

      <div class="filter-group">
        <label>Marca</label>
        <select v-model="filters.brandId">
          <option value="">Todas las marcas</option>
          <option v-for="brand in brands" :key="brand.id" :value="brand.id">{{ brand.name }}</option>
        </select>
      </div>

      <div class="filter-group">
        <label>Familia / Categoría</label>
        <select v-model="filters.groupId">
          <option value="">Todas las categorías</option>
          <option v-for="group in productGroups" :key="group.id" :value="group.id">{{ group.name }}</option>
        </select>
      </div>

      <div class="filter-group">
        <label>Estado</label>
        <select v-model="filters.isActive">
          <option value="">Cualquier estado</option>
          <option value="true">Activos</option>
          <option value="false">Inactivos</option>
        </select>
      </div>
    </template>

    <!-- CAPA 3: TRABAJO (Tabla) -->
    <template #table-header>
      <th>SKU / Referencia</th>
      <th>Nombre del Producto</th>
      <th>Marca</th>
      <th class="align-right">Precio Base</th>
      <th class="text-center">Estado</th>
      <th class="align-right">Acciones</th>
    </template>

    <template #item="{ item }">
      <td><code class="sku-code">{{ item.sku || '—' }}</code></td>
      <td>
        <div class="product-info-cell">
          <strong>{{ item.name }}</strong>
          <div v-if="item.long_name" class="product-subtitle">{{ item.long_name }}</div>
        </div>
      </td>
      <td><span class="text-muted">{{ getBrandName(item.brand_id) }}</span></td>
      <td class="align-right"><strong class="price-text">{{ formatPrice(item.base_price) }}</strong></td>
      <td class="text-center">
        <span :class="['status-badge', item.is_active ? 'status-success' : 'status-secondary']">
          {{ item.is_active ? 'Activo' : 'Inactivo' }}
        </span>
      </td>
      <td class="align-right" @click.stop>
        <div class="action-buttons">
          <button 
            class="btn-icon" 
            @click="toggleStatus(item)" 
            :title="item.is_active ? 'Desactivar' : 'Activar'"
          >
            <span class="material-symbols-outlined">{{ item.is_active ? 'block' : 'check_circle' }}</span>
          </button>
          <button class="btn-icon" @click="navigateToDetail(item)" title="Ver detalle">
            <span class="material-symbols-outlined">visibility</span>
          </button>
        </div>
      </td>
    </template>
  </BaseCatalog>
</template>

<style scoped>
.sku-code { background: var(--color-background); padding: 0.2rem 0.5rem; border-radius: 4px; font-family: var(--font-family-mono); font-size: 0.8rem; font-weight: 700; color: var(--color-secondary); }
.product-subtitle { font-size: 0.75rem; color: var(--color-text-secondary); margin-top: 0.1rem; }
.price-text { color: #16a34a; font-weight: 700; }
.action-buttons { display: flex; justify-content: flex-end; gap: 0.25rem; }
.btn-icon { background: transparent; border: none; cursor: pointer; color: var(--color-text-secondary); padding: 0.4rem; border-radius: 6px; display: inline-flex; align-items: center; justify-content: center; transition: all 0.2s; }
.btn-icon:hover { background: rgba(0,0,0,0.05); color: var(--color-text-primary); }
.align-right { text-align: right; }
</style>
