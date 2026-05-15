<template>
  <div class="page-layout">
    <BaseCatalog
      title="Gestión de Entidades"
      icon="groups"
      :breadcrumbs="[{ label: 'Entidades', to: '/parties' }, { label: 'Listado' }]"
      :items="parties"
      :is-loading="isLoading"
      :error="error"
      :has-filters="hasActiveFilters"
      create-route="/parties/new"
      create-text="Nueva Entidad"
      @clear-filters="clearAllFilters"
      @refresh="loadParties"
      @click-item="goToDetail"
    >
      <template #filters>
        <div class="filter-group">
          <label>Búsqueda rápida</label>
          <input 
            v-model="filters.search" 
            type="text" 
            placeholder="Nombre, NIF o ID..." 
          />
        </div>

        <div class="filter-group">
          <label>Tipo de Entidad</label>
          <select v-model="filters.type">
            <option value="">Cualquier tipo</option>
            <option value="ORGANIZATION">Organización / Empresa</option>
            <option value="PERSON">Persona Física</option>
          </select>
        </div>

        <div class="filter-group">
          <label>Rol</label>
          <select v-model="filters.role">
            <option value="">Cualquier rol</option>
            <option value="CLIENT">Cliente</option>
            <option value="SUPPLIER">Proveedor</option>
            <option value="BOTH">Ambos (Dual)</option>
          </select>
        </div>
      </template>

      <template #table-header>
        <th>Identidad / Nombre</th>
        <th>NIF/CIF</th>
        <th>Tipo</th>
        <th>Roles</th>
        <th>Estado</th>
        <th class="align-right">Acciones</th>
      </template>

      <template #item="{ item }">
        <td>
          <div class="party-identity">
            <div class="party-avatar">
              <component :is="item.type === 'ORGANIZATION' ? Building2 : User" :size="20" />
            </div>
            <strong>{{ item.name }}</strong>
          </div>
        </td>
        <td><code class="code-badge">{{ item.tax_id || '—' }}</code></td>
        <td><span class="type-tag">{{ item.type === 'ORGANIZATION' ? 'Empresa' : 'Persona' }}</span></td>
        <td>
          <div class="roles-stack">
            <span v-if="item.role === 'CLIENT' || item.role === 'BOTH'" class="role-pill client">Cliente</span>
            <span v-if="item.role === 'SUPPLIER' || item.role === 'BOTH'" class="role-pill supplier">Proveedor</span>
          </div>
        </td>
        <td>
          <span :class="['status-badge', item.status === 'ACTIVE' ? 'status-success' : 'status-secondary']">
            {{ item.status === 'ACTIVE' ? 'Activo' : 'Inactivo' }}
          </span>
        </td>
        <td class="align-right" @click.stop>
          <div class="action-buttons">
            <button @click="goToDetail(item)" class="btn btn-ghost" title="Ver ficha completa"><Eye :size="18" /></button>
            <button @click="promptDelete(item)" class="btn btn-ghost text-danger" title="Eliminar entidad"><Trash2 :size="18" /></button>
          </div>
        </td>
      </template>
    </BaseCatalog>

    <!-- MODAL DE CONFIRMACIÓN DE ELIMINACIÓN -->
    <BaseDialog
      :show="confirmDelete.show"
      title="Eliminar Entidad"
      icon="warning"
      confirm-text="Eliminar Definitivamente"
      confirm-class="btn-danger"
      :is-confirming="isDeleting"
      @close="confirmDelete.show = false"
      @confirm="executeDelete"
    >
      <p>¿Estás seguro de que deseas eliminar permanentemente a <strong>{{ confirmDelete.party?.name }}</strong>?</p>
      <p class="mt-2 text-muted italic">Esta acción solo se completará si no existen documentos vinculados (pedidos, facturas, etc.).</p>
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Building2, User, Eye, Trash2 } from 'lucide-vue-next'
import BaseCatalog from '@/components/shared/BaseCatalog.vue'
import BaseDialog from '@/components/shared/BaseDialog.vue'
import { partyApi } from '@/services/partyApi'
import { useToastStore } from '@/stores/toast'

const router = useRouter()
const route = useRoute()
const toastStore = useToastStore()
const allParties = ref([])
const isLoading = ref(false)
const isDeleting = ref(false)
const error = ref('')

const filters = reactive({
  search: (route.query.search as string) || '',
  type: (route.query.type as string) || '',
  role: (route.query.role as string) || ''
})

const hasActiveFilters = computed(() => {
  return filters.search.trim() !== '' || filters.type !== '' || filters.role !== ''
})

function clearAllFilters() {
  filters.search = ''
  filters.type = ''
  filters.role = ''
  router.replace({ query: {} })
}

// Update filters if route changes (e.g. clicking a link in navbar while already on this page)
watch(() => route.query, (newQuery) => {
  filters.search = (newQuery.search as string) || ''
  filters.type = (newQuery.type as string) || ''
  filters.role = (newQuery.role as string) || ''
})

const parties = computed(() => {
  let result = allParties.value;
  if (filters.search) {
    const q = filters.search.toLowerCase();
    result = result.filter(p => 
      p.name.toLowerCase().includes(q) || 
      (p.tax_id && p.tax_id.toLowerCase().includes(q))
    );
  }
  if (filters.type) result = result.filter(p => p.type === filters.type);
  if (filters.role) result = result.filter(p => p.role === filters.role || p.role === 'BOTH');
  return result;
});

// --- Confirm Dialog Logic ---
const confirmDelete = reactive({
  show: false,
  party: null
})

function promptDelete(party) {
  confirmDelete.party = party
  confirmDelete.show = true
}

async function executeDelete() {
  if (!confirmDelete.party) return
  isDeleting.value = true
  try {
    await partyApi.deleteParty(confirmDelete.party.id)
    toastStore.success(`Entidad "${confirmDelete.party.name}" eliminada`)
    await loadParties()
    confirmDelete.show = false
  } catch (err) {
    toastStore.error(err.message || 'Error al eliminar entidad')
    confirmDelete.show = false // Close on error too to avoid stuck modal
  } finally {
    isDeleting.value = false
  }
}

async function loadParties() {
  isLoading.value = true; error.value = '';
  try { 
    const res = await partyApi.listParties({ limit: 1000 }); 
    allParties.value = res.data || [];
  } catch (err) { error.value = err.message; } 
  finally { isLoading.value = false }
}

function goToDetail(party) { router.push(`/parties/${party.id}`); }

onMounted(() => loadParties())
</script>

<style scoped>
.page-layout { background-color: var(--color-background); min-height: 100vh; }
.party-identity { display: flex; align-items: center; gap: 0.75rem; }
.party-avatar { width: 36px; height: 36px; border-radius: 8px; background: var(--color-background); display: flex; align-items: center; justify-content: center; color: var(--color-secondary); border: 1px solid var(--color-border); }

.code-badge { background: var(--color-background); padding: 0.2rem 0.4rem; border-radius: 4px; font-family: var(--font-family-mono); font-size: 0.8rem; }
.type-tag { font-size: 0.8rem; color: var(--color-text-secondary); }

.roles-stack { display: flex; gap: 0.4rem; }
.role-pill { font-size: 0.65rem; font-weight: 800; text-transform: uppercase; padding: 0.15rem 0.4rem; border-radius: 4px; }
.role-pill.client { background: rgba(34, 197, 94, 0.1); color: #16a34a; }
.role-pill.supplier { background: rgba(147, 51, 234, 0.1); color: #9333ea; }

.action-buttons { display: flex; justify-content: flex-end; gap: 0.25rem; }
.align-right { text-align: right; }
</style>
