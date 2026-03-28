<template>
  <div class="page-layout">
    <Navbar style="z-index: 2000;" />
    
    <BaseCatalog
      title="Órdenes de Trabajo (Producción)"
      icon="precision_manufacturing"
      :breadcrumbs="[{ label: 'MES', to: '/mes/dashboard' }, { label: 'Órdenes de Trabajo' }]"
      :items="works"
      :is-loading="isLoading"
      :error="error"
      :has-filters="hasFilters"
      create-route="/mes/work-orders/new"
      create-text="Nueva Orden"
      empty-icon="pending_actions"
      empty-text="No hay órdenes de trabajo registradas"
      @clear-filters="clearFilters"
      @refresh="loadWorks"
      @click-item="(item) => navigateToDetail(item.id)"
    >
      <template #filters>
        <div class="filter-group">
          <label>Búsqueda</label>
          <input v-model="search" type="text" placeholder="Número, nombre o cliente..." />
        </div>

        <div class="filter-group">
          <label>Estado de Producción</label>
          <select v-model="statusFilter">
            <option value="">Todos los estados</option>
            <option value="PENDING">Pendiente</option>
            <option value="IN_PROGRESS">En progreso</option>
            <option value="ON_HOLD">En espera</option>
            <option value="COMPLETED">Completado</option>
            <option value="CANCELLED">Cancelado</option>
          </select>
        </div>
      </template>

      <template #table-header>
        <th>Número</th>
        <th>Descripción del Trabajo</th>
        <th>Pedido Origen</th>
        <th>Estado</th>
        <th>Prioridad</th>
        <th class="align-right">Acciones</th>
      </template>

      <template #item="{ item }">
        <td><code class="code-badge">{{ item.work_number }}</code></td>
        <td><strong>{{ item.work_name }}</strong></td>
        <td>
          <span v-if="item.sales_order_number" class="text-muted">Pedido #{{ item.sales_order_number }}</span>
          <span v-else class="text-muted">—</span>
        </td>
        <td>
          <span :class="['status-badge', `status-${getStatusClass(item.status)}`]">
            {{ mesApi.getWorkStatusLabel(item.status) }}
          </span>
        </td>
        <td>
          <span :class="['priority-pill', `prio-${item.priority}`]">
            {{ mesApi.getPriorityLabel(item.priority) }}
          </span>
        </td>
        <td class="align-right" @click.stop>
          <div class="action-buttons">
            <button class="btn-icon" @click="navigateToDetail(item.id)" title="Ver detalle">
              <span class="material-symbols-outlined">visibility</span>
            </button>
          </div>
        </td>
      </template>
    </BaseCatalog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, computed, watch, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import BaseCatalog from '@/components/shared/BaseCatalog.vue'
import { mesApi } from '@/services/mesApi'
import type { WorkOrder } from '@/types/mes'

const router = useRouter()
const route = useRoute()
const works = ref<WorkOrder[]>([])
const isLoading = ref(false)
const error = ref('')
const search = ref('')
const statusFilter = ref((route.query.status as string) || '')

const hasFilters = computed(() => search.value.trim() !== '' || statusFilter.value !== '')

let searchTimeout: any = null
watch([search, statusFilter], () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => loadWorks(), 350)
})

async function loadWorks() {
  isLoading.value = true
  error.value = ''
  try {
    works.value = await mesApi.listWorkOrders({
      search: search.value.trim() || undefined,
      status: statusFilter.value || undefined,
    })
  } catch (err: any) {
    error.value = 'No se han podido cargar las órdenes de trabajo.'
  } finally {
    isLoading.value = false
  }
}

function clearFilters() { search.value = ''; statusFilter.value = ''; }
function navigateToDetail(id: string) { router.push(`/mes/work-orders/${id}`); }

function getStatusClass(status: string) {
  const map: Record<string, string> = { PENDING: 'warning', IN_PROGRESS: 'info', COMPLETED: 'success', CANCELLED: 'danger', ON_HOLD: 'secondary' }
  return map[status] || 'secondary'
}

onMounted(loadWorks)
onUnmounted(() => { if (searchTimeout) clearTimeout(searchTimeout) })
</script>

<style scoped>
.page-layout { background-color: var(--color-background); min-height: 100vh; }
.code-badge { background: var(--color-background); padding: 0.2rem 0.4rem; border-radius: 4px; font-family: var(--font-family-mono); font-size: 0.8rem; font-weight: 700; color: var(--color-secondary); }
.priority-pill { font-size: 0.7rem; font-weight: 800; text-transform: uppercase; padding: 0.2rem 0.5rem; border-radius: 4px; }
.prio-HIGH, .prio-URGENT { background: rgba(220, 38, 38, 0.1); color: #dc2626; }
.prio-MEDIUM { background: rgba(217, 119, 6, 0.1); color: #d97706; }
.prio-LOW { background: rgba(37, 99, 235, 0.1); color: #2563eb; }
.align-right { text-align: right; }
.action-buttons { display: flex; justify-content: flex-end; }
.btn-icon { color: var(--color-text-secondary); transition: 0.2s; padding: 0.4rem; border-radius: 6px; border: none; background: transparent; cursor: pointer; }
.btn-icon:hover { color: var(--color-text-primary); background: rgba(0,0,0,0.05); }
</style>
