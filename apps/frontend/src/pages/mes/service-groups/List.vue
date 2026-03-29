<template>
  <BaseCatalog
    title="Tipos de Trabajo"
    icon="account_tree"
    :breadcrumbs="[{ label: 'MES', to: '/mes/dashboard' }, { label: 'Tipos de Trabajo' }]"
    :items="workTypes"
    :is-loading="isLoading"
    :error="error"
    :has-filters="!!search || statusFilter !== ''"
    create-route="/mes/work-types/new"
    create-text="Nuevo Tipo"
    empty-icon="tree_off"
    empty-text="No hay tipos de trabajo registrados"
    @clear-filters="clearFilters"
    @refresh="loadWorkTypes"
    @click-item="(item) => navigateToEdit(item.id)"
  >
    <template #filters>
      <div class="filter-group">
        <label>Búsqueda</label>
        <input 
          v-model="search" 
          type="text" 
          placeholder="Nombre o referencia..." 
          @input="debouncedSearch"
        />
      </div>

      <div class="filter-group">
        <label>Estado</label>
        <select v-model="statusFilter" @change="loadWorkTypes">
          <option value="">Todos los estados</option>
          <option value="true">Activos</option>
          <option value="false">Inactivos</option>
        </select>
      </div>
    </template>

    <template #table-header>
      <th>Nombre del Tipo</th>
      <th>Referencia</th>
      <th>Descripción</th>
      <th class="text-center">Tareas</th>
      <th class="text-center">Estado</th>
      <th class="align-right">Acciones</th>
    </template>

    <template #item="{ item }">
      <td><strong>{{ item.name }}</strong></td>
      <td><code class="code-badge">{{ item.reference || '—' }}</code></td>
      <td><span class="text-muted">{{ item.description || '—' }}</span></td>
      <td class="text-center">
        <span class="status-badge-sm">{{ item.tasks?.length || 0 }} pasos</span>
      </td>
      <td class="text-center">
        <span :class="['status-pill', item.is_active ? 'status-active' : 'status-inactive']">
          {{ item.is_active ? 'Activo' : 'Inactivo' }}
        </span>
      </td>
      <td class="align-right" @click.stop>
        <div class="action-buttons">
          <router-link :to="`/mes/work-types/${item.id}/edit`" class="btn-icon" title="Editar">
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
import BaseCatalog from '@/components/shared/BaseCatalog.vue'
import { mesApi } from '@/services/mesApi'
import type { MESWorkType } from '@/types/mes'

const router = useRouter()
const workTypes = ref<MESWorkType[]>([])
const isLoading = ref(false)
const error = ref('')
const search = ref('')
const statusFilter = ref('')

let searchTimeout: any = null

function debouncedSearch() {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => loadWorkTypes(), 350)
}

async function loadWorkTypes() {
  isLoading.value = true
  error.value = ''
  try {
    const isActive = statusFilter.value === '' ? undefined : statusFilter.value === 'true'
    workTypes.value = await mesApi.listWorkTypes({
      search: search.value.trim() || undefined,
      is_active: isActive,
    })
  } catch (err: any) {
    error.value = err.message || 'Error al cargar tipos de trabajo'
  } finally {
    isLoading.value = false
  }
}

async function toggleActive(item: MESWorkType) {
  try {
    await mesApi.updateWorkType(item.id, { is_active: !item.is_active })
    await loadWorkTypes()
  } catch (err: any) {
    alert(err.message)
  }
}

function clearFilters() {
  search.value = ''
  statusFilter.value = ''
  loadWorkTypes()
}

function navigateToEdit(id: string) {
  router.push(`/mes/work-types/${id}/edit`)
}

onMounted(loadWorkTypes)
onUnmounted(() => { if (searchTimeout) clearTimeout(searchTimeout) })
</script>

<style scoped>
.code-badge { background: var(--color-background); padding: 0.2rem 0.5rem; border-radius: 4px; font-family: var(--font-family-mono); font-size: 0.8rem; font-weight: 700; color: var(--color-secondary); }
.align-right { text-align: right; }
.text-center { text-align: center; }

.status-badge-sm { background: var(--color-background); padding: 0.1rem 0.5rem; border-radius: 4px; font-size: 0.7rem; font-weight: 600; color: var(--color-text-secondary); }

.action-buttons { display: flex; justify-content: flex-end; gap: 0.25rem; }
.btn-icon { background: transparent; border: none; cursor: pointer; color: var(--color-text-secondary); padding: 0.4rem; border-radius: 6px; }
.btn-icon:hover { background: rgba(0,0,0,0.05); color: var(--color-text-primary); }
</style>