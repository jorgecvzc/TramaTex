<template>
  <div class="page-layout">
    
    <BaseCatalog
      title="Configuraciones Técnicas por Cliente"
      icon="settings_input_component"
      :breadcrumbs="[{ label: 'MES', to: '/mes/dashboard' }, { label: 'Configuraciones' }]"
      :items="setups"
      :is-loading="isLoading"
      :error="error"
      :has-filters="hasFilters"
      create-route="/mes/work-setups/new"
      create-text="Nueva Configuración"
      empty-icon="settings_off"
      empty-text="No hay configuraciones técnicas registradas"
      @clear-filters="clearFilters"
      @refresh="loadSetups"
      @click-item="(item) => navigateToEdit(item.id)"
    >
      <template #filters>
        <div class="filter-group">
          <label>Búsqueda</label>
          <input v-model="search" type="text" placeholder="Nombre de configuración..." />
        </div>

        <div class="filter-group">
          <label>Estado</label>
          <select v-model="statusFilter">
            <option value="">Cualquier estado</option>
            <option value="true">Activas</option>
            <option value="false">Inactivas</option>
          </select>
        </div>
      </template>

      <template #table-header>
        <th>Nombre de la Configuración</th>
        <th>Descripción / Notas Técnicas</th>
        <th class="text-center">Operaciones</th>
        <th class="text-center">Estado</th>
        <th class="align-right">Acciones</th>
      </template>

      <template #item="{ item }">
        <td><strong>{{ item.name }}</strong></td>
        <td><span class="text-muted">{{ item.description || '—' }}</span></td>
        <td class="text-center">
          <span class="count-badge">{{ item.lines?.length || 0 }} etapas</span>
        </td>
        <td class="text-center">
          <span :class="['status-badge', item.is_active ? 'status-success' : 'status-secondary']">
            {{ item.is_active ? 'Activa' : 'Inactiva' }}
          </span>
        </td>
        <td class="align-right" @click.stop>
          <div class="action-buttons">
            <router-link :to="`/mes/work-setups/${item.id}/edit`" class="btn-icon" title="Editar">
              <Pencil :size="18" />
            </router-link>
            <button 
              class="btn-icon" 
              @click="toggleActive(item)" 
              :title="item.is_active ? 'Desactivar' : 'Activar'"
            >
              <component :is="item.is_active ? Ban : CheckCircle" :size="18" />
            </button>
          </div>
        </td>
      </template>
    </BaseCatalog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, computed, watch, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { Pencil, Ban, CheckCircle } from 'lucide-vue-next'
import BaseCatalog from '@/components/shared/BaseCatalog.vue'
import { mesApi } from '@/services/mesApi'
import { useToastStore } from '@/stores/toast'
import type { WorkSetup } from '@/types/mes'

const router = useRouter()
const toastStore = useToastStore()
const setups = ref<WorkSetup[]>([])
const isLoading = ref(false)
const error = ref('')
const search = ref('')
const statusFilter = ref('')

const hasFilters = computed(() => search.value.trim() !== '' || statusFilter.value !== '')

let searchTimeout: any = null
watch([search, statusFilter], () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => loadSetups(), 350)
})

async function loadSetups() {
  isLoading.value = true
  error.value = ''
  try {
    const isActive = statusFilter.value === '' ? undefined : statusFilter.value === 'true'
    const res = await mesApi.listWorkSetups({
      search: search.value.trim() || undefined,
      is_active: isActive,
    })
    setups.value = res.data || res || []
  } catch (err: any) {
    error.value = 'Error al cargar las configuraciones técnicas.'
  } finally {
    isLoading.value = false
  }
}

async function toggleActive(setup: WorkSetup) {
  try {
    await mesApi.updateWorkSetup(setup.id, { is_active: !setup.is_active })
    toastStore.addToast(`Configuración ${setup.is_active ? 'desactivada' : 'activada'} correctamente`, 'info')
    await loadSetups()
  } catch (err: any) {
    toastStore.addToast(err.message, 'error')
  }
}

function clearFilters() { search.value = ''; statusFilter.value = ''; }
function navigateToEdit(id: string) { router.push(`/mes/work-setups/${id}/edit`); }

onMounted(loadSetups)
onUnmounted(() => { if (searchTimeout) clearTimeout(searchTimeout) })
</script>

<style scoped>
.page-layout { background-color: var(--color-background); min-height: 100vh; }
.count-badge { background: var(--color-background); padding: 0.2rem 0.6rem; border-radius: 20px; font-size: 0.75rem; font-weight: 600; color: var(--color-text-secondary); }
.align-right { text-align: right; }
.text-center { text-align: center; }
.action-buttons { display: flex; justify-content: flex-end; gap: 0.25rem; }
.btn-icon { background: transparent; border: none; cursor: pointer; color: var(--color-text-secondary); padding: 0.4rem; border-radius: 6px; display: inline-flex; align-items: center; justify-content: center; }
.btn-icon:hover { background: rgba(0,0,0,0.05); color: var(--color-text-primary); }
</style>