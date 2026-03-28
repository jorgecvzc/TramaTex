<template>
  <Navbar />
  
  <BaseCatalog
    title="Posiciones"
    icon="factory"
    :breadcrumbs="[{ label: 'MES', to: '/mes/dashboard' }, { label: 'Posiciones' }]"
    :items="positions"
    :is-loading="isLoading"
    :error="error"
    :has-filters="!!search || statusFilter !== ''"
    create-route="/mes/positions/new"
    create-text="Nueva Posición"
    empty-icon="location_off"
    empty-text="No hay posiciones registradas"
    @clear-filters="clearFilters"
    @refresh="loadPositions"
    @click-item="(item) => navigateToEdit(item.id)"
  >
    <template #filters>
      <div class="filter-group">
        <label>Búsqueda</label>
        <input 
          v-model="search" 
          type="text" 
          placeholder="Nombre o código de posición..." 
          @input="debouncedSearch"
        />
      </div>

      <div class="filter-group">
        <label>Estado</label>
        <select v-model="statusFilter" @change="loadPositions">
          <option value="">Todos los estados</option>
          <option value="true">Activas</option>
          <option value="false">Inactivas</option>
        </select>
      </div>
    </template>

    <template #table-header>
      <th>Nombre</th>
      <th>Código</th>
      <th>Descripción</th>
      <th class="text-center">Estado</th>
      <th class="align-right">Acciones</th>
    </template>

    <template #item="{ item }">
      <td><strong>{{ item.name }}</strong></td>
      <td><code class="code-badge">{{ item.code }}</code></td>
      <td><span class="text-muted">{{ item.description || '—' }}</span></td>
      <td class="text-center">
        <span :class="['status-pill', item.is_active ? 'status-active' : 'status-inactive']">
          {{ item.is_active ? 'Activa' : 'Inactiva' }}
        </span>
      </td>
      <td class="align-right" @click.stop>
        <div class="action-buttons">
          <router-link :to="`/mes/positions/${item.id}/edit`" class="btn-icon" title="Editar">
            <span class="material-symbols-outlined">edit</span>
          </router-link>
          <button 
            class="btn-icon" 
            @click="toggleActive(item)" 
            :title="item.is_active ? 'Desactivar' : 'Activar'"
            :class="{ 'text-warning': item.is_active }"
          >
            <span class="material-symbols-outlined">{{ item.is_active ? 'block' : 'check_circle' }}</span>
          </button>
        </div>
      </td>
    </template>
  </BaseCatalog>
</template>

<script setup lang="ts">
import { onMounted, ref, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import BaseCatalog from '@/components/shared/BaseCatalog.vue'
import { mesApi } from '@/services/mesApi'
import type { MESPosition } from '@/types/mes'

const router = useRouter()
const positions = ref<MESPosition[]>([])
const isLoading = ref(false)
const error = ref('')
const search = ref('')
const statusFilter = ref('')

let searchTimeout: any = null

function debouncedSearch() {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => loadPositions(), 350)
}

async function loadPositions() {
  isLoading.value = true
  error.value = ''
  try {
    const isActive = statusFilter.value === '' ? undefined : statusFilter.value === 'true'
    positions.value = await mesApi.listPositions({
      search: search.value.trim() || undefined,
      is_active: isActive,
    })
  } catch (err: any) {
    error.value = err.message || 'Error al cargar posiciones'
  } finally {
    isLoading.value = false
  }
}

async function toggleActive(position: MESPosition) {
  try {
    await mesApi.updatePosition(position.id, { is_active: !position.is_active })
    await loadPositions()
  } catch (err: any) {
    alert(err.message)
  }
}

function clearFilters() {
  search.value = ''
  statusFilter.value = ''
  loadPositions()
}

function navigateToEdit(id: string) {
  router.push(`/mes/positions/${id}/edit`)
}

onMounted(loadPositions)
onUnmounted(() => { if (searchTimeout) clearTimeout(searchTimeout) })
</script>

<style scoped>
.code-badge { background: var(--color-background); padding: 0.2rem 0.5rem; border-radius: 4px; font-family: var(--font-family-mono); font-size: 0.8rem; font-weight: 700; color: var(--color-secondary); }
.align-right { text-align: right; }
.text-center { text-align: center; }

.action-buttons { display: flex; justify-content: flex-end; }
.btn-icon { color: var(--color-text-secondary); transition: 0.2s; padding: 0.4rem; border-radius: 6px; border: none; background: transparent; cursor: pointer; }
.btn-icon:hover { color: var(--color-text-primary); background: rgba(0,0,0,0.05); }
</style>