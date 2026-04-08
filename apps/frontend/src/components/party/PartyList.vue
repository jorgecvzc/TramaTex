<script setup lang="ts">
/**
 * PartyList.vue - Listado Maestro de Entidades (Clientes/Proveedores)
 * 
 * Implementa el estándar BaseCatalog con Arquitectura de 3 Capas.
 */
import { ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { partyApi } from '@/services/partyApi'
import BaseCatalog from '@/components/shared/BaseCatalog.vue'

const router = useRouter()
const route = useRoute()
const parties = ref<any[]>([])
const isLoading = ref(false)
const error = ref('')

const filters = reactive({ 
  name: (route.query.name as string) || '',
  role: (route.query.role as string) || '', 
  status: (route.query.status as string) || '' 
})

const hasFilters = computed(() => 
  filters.name.trim() !== '' || filters.role !== '' || filters.status !== ''
)

// Sync filters when navigating to this page from another route (e.g. dashboard KPI cards)
watch(() => route.query, (newQuery) => {
  filters.name = (newQuery.name as string) || ''
  filters.role = (newQuery.role as string) || ''
  filters.status = (newQuery.status as string) || ''
}, { deep: true })

// Lógica de filtrado con debounce
let debounceTimer: any = null
watch(filters, () => {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => fetchParties(), 350)
}, { deep: true })

async function fetchParties() {
  isLoading.value = true
  error.value = ''
  try {
    const res = await partyApi.listParties({ 
      name: filters.name,
      role: filters.role,
      status: filters.status,
      pageSize: 100
    })
    parties.value = res.data || (Array.isArray(res) ? res : [])
  } catch (err: any) { 
    error.value = 'No se han podido cargar las entidades.'
    console.error(err)
  } finally { 
    isLoading.value = false 
  }
}

function clearFilters() { 
  filters.name = ''
  filters.role = ''
  filters.status = ''
}

function navigateToDetail(party: any) { 
  router.push(`/parties/${party.id}`) 
}

function formatRole(r: string) { 
  const map: Record<string, string> = { 
    CLIENT: 'Cliente', 
    SUPPLIER: 'Proveedor', 
    BOTH: 'Cliente/Prov.', 
    CONTACT: 'Contacto' 
  }
  return map[r] || r
}

function getStatusLabel(s: string) { return s === 'ACTIVE' ? 'Activo' : 'Inactivo'; }
function getStatusClass(s: string) { return s === 'ACTIVE' ? 'success' : 'secondary'; }
function formatDate(d: string) { return d ? new Date(d).toLocaleDateString('es-ES', { year: 'numeric', month: 'short', day: 'numeric' }) : '—' }

async function toggleStatus(party: any) {
  const newStatus = party.status === 'ACTIVE' ? 'INACTIVE' : 'ACTIVE'
  try { 
    await partyApi.changePartyStatus(party.id, newStatus)
    party.status = newStatus 
  } catch (err: any) {
    alert('Error al cambiar estado: ' + err.message)
  }
}

async function deleteParty(party: any) {
  if (!confirm(`¿Eliminar "${party.name}"? Esta acción no se puede deshacer.`)) return
  try {
    await partyApi.deleteParty(party.id)
    await fetchParties()
  } catch (err: any) {
    alert('No se pudo eliminar la entidad: ' + (err.message || 'Error desconocido'))
  }
}

onMounted(fetchParties)
onUnmounted(() => { if (debounceTimer) clearTimeout(debounceTimer) })
</script>

<template>
  <BaseCatalog
    title="Base de Datos de Entidades"
    :breadcrumbs="[{ label: 'Entidades', to: '/parties/dashboard' }, { label: 'Clientes y Proveedores' }]"
    :items="parties"
    :is-loading="isLoading"
    :error="error"
    :has-filters="hasFilters"
    create-route="/parties/new"
    create-text="Nueva Entidad"
    empty-icon="groups"
    empty-text="No hay entidades registradas"
    @clear-filters="clearFilters"
    @refresh="fetchParties"
    @click-item="navigateToDetail"
  >
    <!-- CAPA 2: CONTEXTO (Filtros) -->
    <template #filters>
      <div class="filter-group">
        <label>Búsqueda</label>
        <input v-model="filters.name" type="text" placeholder="Nombre, NIF, Email..." />
      </div>

      <div class="filter-group">
        <label>Rol de Negocio</label>
        <select v-model="filters.role">
          <option value="">Todos los roles</option>
          <option value="CLIENT">Clientes</option>
          <option value="SUPPLIER">Proveedores</option>
          <option value="BOTH">Ambos</option>
          <option value="CONTACT">Contactos</option>
        </select>
      </div>

      <div class="filter-group">
        <label>Estado</label>
        <select v-model="filters.status">
          <option value="">Cualquier estado</option>
          <option value="ACTIVE">Activos</option>
          <option value="INACTIVE">Inactivos</option>
        </select>
      </div>
    </template>

    <!-- CAPA 3: TRABAJO (Tabla) -->
    <template #table-header>
      <th>Nombre / Razón Social</th>
      <th>Tipo</th>
      <th>Rol</th>
      <th>NIF / CIF</th>
      <th>Estado</th>
      <th>Fecha Alta</th>
      <th class="align-right">Acciones</th>
    </template>

    <template #item="{ item }">
      <td><strong class="font-bold">{{ item.name }}</strong></td>
      <td>
        <div class="type-info">
          <span class="material-symbols-outlined icon-secondary">{{ item.has_person ? 'person' : 'domain' }}</span>
          <span>{{ item.has_person ? 'Persona' : 'Organización' }}</span>
        </div>
      </td>
      <td><span class="role-badge">{{ formatRole(item.role) }}</span></td>
      <td><code class="text-mono">{{ item.tax_id || '—' }}</code></td>
      <td>
        <span :class="['status-badge', `status-${getStatusClass(item.status)}`]">
          {{ getStatusLabel(item.status) }}
        </span>
      </td>
      <td class="text-muted">{{ formatDate(item.created_at) }}</td>
      <td class="align-right" @click.stop>
        <div class="action-buttons">
          <button 
            class="btn btn-ghost" 
            @click="toggleStatus(item)" 
            :title="item.status === 'ACTIVE' ? 'Desactivar' : 'Activar'"
          >
            <span class="material-symbols-outlined">{{ item.status === 'ACTIVE' ? 'block' : 'check_circle' }}</span>
          </button>
          <button 
            v-if="item.can_delete"
            class="btn btn-ghost btn-danger"
            @click.stop="deleteParty(item)"
            title="Eliminar entidad"
          >
            <span class="material-symbols-outlined">delete</span>
          </button>
        </div>
      </td>
    </template>
  </BaseCatalog>
</template>

<style scoped>
.font-bold { font-weight: 700; color: var(--color-text-primary); }
.text-mono { font-family: var(--font-family-mono); font-size: 0.8rem; background: var(--color-background); padding: 0.2rem 0.4rem; border-radius: 4px; }
.type-info { display: flex; align-items: center; gap: 0.5rem; font-size: 0.85rem; color: var(--color-text-secondary); }
.role-badge { font-weight: 600; font-size: 0.85rem; color: var(--color-text-primary); }
.action-buttons { display: flex; justify-content: flex-end; gap: 0.25rem; }
.align-right { text-align: right; }
</style>
