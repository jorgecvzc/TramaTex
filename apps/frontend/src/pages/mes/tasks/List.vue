<template>
  <BaseCatalog
    title="Catálogo de Tareas Básicas"
    icon="assignment"
    :breadcrumbs="[{ label: 'MES', to: '/mes/dashboard' }, { label: 'Tareas' }]"
    :items="tasks"
    :is-loading="isLoading"
    :error="error"
    :has-filters="hasFilters"
    empty-icon="assignment_late"
    empty-text="No hay tareas básicas registradas"
    @clear-filters="clearFilters"
    @refresh="loadTasks"
    @click-item="editTask"
  >
      <template #header-actions>
        <button @click="openCreateModal" class="btn btn-primary">
          <span class="material-symbols-outlined">add</span>
          <span>Nueva Tarea</span>
        </button>
      </template>

      <template #filters>
        <div class="filter-group">
          <label>Búsqueda</label>
          <input v-model="filters.search" type="text" placeholder="Nombre o referencia..." />
        </div>

        <div class="filter-group">
          <label>Estado</label>
          <select v-model="filters.isActive">
            <option value="">Cualquier estado</option>
            <option value="true">Activas</option>
            <option value="false">Inactivas</option>
          </select>
        </div>
      </template>

      <template #table-header>
        <th>Nombre de la Tarea</th>
        <th>Referencia / Código</th>
        <th>Descripción</th>
        <th class="text-center">Estado</th>
        <th class="align-right">Acciones</th>
      </template>

      <template #item="{ item }">
        <td><strong>{{ item.name }}</strong></td>
        <td><code class="code-badge">{{ item.reference || '—' }}</code></td>
        <td><span class="text-muted">{{ item.description || '—' }}</span></td>
        <td class="text-center">
          <span :class="['status-badge', item.is_active ? 'status-success' : 'status-secondary']">
            {{ item.is_active ? 'Activa' : 'Inactiva' }}
          </span>
        </td>
        <td class="align-right" @click.stop>
          <div class="action-buttons">
            <button @click="editTask(item)" class="btn-icon" title="Editar"><span class="material-symbols-outlined">edit</span></button>
            <button 
              @click="toggleActive(item)" 
              class="btn-icon" 
              :title="item.is_active ? 'Desactivar' : 'Activar'"
            >
              <span class="material-symbols-outlined">{{ item.is_active ? 'block' : 'check_circle' }}</span>
            </button>
          </div>
        </td>
      </template>
    </BaseCatalog>

    <!-- MODAL: CREAR/EDITAR TAREA -->
    <BaseDialog
      :show="showModal"
      :title="modalMode === 'create' ? 'Nueva Tarea Básica' : 'Editar Tarea'"
      icon="assignment"
      confirm-text="Guardar Tarea"
      :is-confirming="isSaving"
      @close="showModal = false"
      @confirm="saveTask"
    >
      <div class="form-group">
        <label>Nombre de la Tarea *</label>
        <input v-model="formData.name" type="text" class="form-input" placeholder="Ej: Bordado de logotipo..." required />
      </div>
      <div class="form-row mt-4">
        <div class="form-group">
          <label>Referencia Técnica</label>
          <input v-model="formData.reference" type="text" class="form-input" placeholder="Ej: BORD-01" />
        </div>
        <div class="form-group">
          <label>Estado</label>
          <label class="checkbox-label mt-2">
            <input v-model="formData.is_active" type="checkbox" class="form-checkbox" />
            <span>Tarea activa</span>
          </label>
        </div>
      </div>
      <div class="form-group mt-4">
        <label>Descripción detallada</label>
        <textarea v-model="formData.description" class="form-textarea" rows="3"></textarea>
      </div>
    </BaseDialog>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive, computed, watch, onUnmounted } from 'vue'
import BaseCatalog from '@/components/shared/BaseCatalog.vue'
import BaseDialog from '@/components/shared/BaseDialog.vue'
import { mesApi } from '@/services/mesApi'
import { useToastStore } from '@/stores/toast'
import type { MESTask } from '@/types/mes'

const toastStore = useToastStore()
const tasks = ref<MESTask[]>([])
const isLoading = ref(false)
const isSaving = ref(false)
const error = ref('')

const filters = reactive({ search: '', isActive: '' })
const hasFilters = computed(() => filters.search.trim() !== '' || filters.isActive !== '')

// Modal State
const showModal = ref(false)
const modalMode = ref<'create' | 'edit'>('create')
const selectedTaskId = ref<string | null>(null)
const formData = reactive({ name: '', reference: '', description: '', is_active: true })

let debounceTimer: any = null
watch(filters, () => {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => loadTasks(), 350)
}, { deep: true })

async function loadTasks() {
  isLoading.value = true
  error.value = ''
  try {
    const isActive = filters.isActive === '' ? undefined : filters.isActive === 'true'
    tasks.value = await mesApi.listTasks({ search: filters.search.trim() || undefined, is_active: isActive })
  } catch (err: any) { error.value = 'Error al cargar tareas básicas.' }
  finally { isLoading.value = false }
}

function openCreateModal() {
  modalMode.value = 'create'; selectedTaskId.value = null
  Object.assign(formData, { name: '', reference: '', description: '', is_active: true })
  showModal.value = true
}

function editTask(item: MESTask) {
  modalMode.value = 'edit'; selectedTaskId.value = item.id
  Object.assign(formData, { name: item.name, reference: item.reference, description: item.description, is_active: item.is_active })
  showModal.value = true
}

async function saveTask() {
  if (!formData.name) return;
  isSaving.value = true
  try {
    if (modalMode.value === 'create') await mesApi.createTask(formData)
    else await mesApi.updateTask(selectedTaskId.value!, formData)
    showModal.value = false; await loadTasks()
  } catch (err: any) { toastStore.addToast(err.message, 'error') }
  finally { isSaving.value = false }
}

async function toggleActive(item: MESTask) {
  try { await mesApi.updateTask(item.id, { is_active: !item.is_active }); await loadTasks() }
  catch (err: any) { toastStore.addToast(err.message, 'error') }
}

function clearFilters() { filters.search = ''; filters.isActive = ''; }

onMounted(loadTasks)
onUnmounted(() => { if (debounceTimer) clearTimeout(debounceTimer) })
</script>

<style scoped>
.page-layout { background-color: var(--color-background); min-height: 100vh; }
.code-badge { background: var(--color-background); padding: 0.2rem 0.4rem; border-radius: 4px; font-family: var(--font-family-mono); font-size: 0.8rem; font-weight: 700; color: var(--color-secondary); }
.align-right { text-align: right; }
.action-buttons { display: flex; justify-content: flex-end; gap: 0.25rem; }
.btn-icon { background: transparent; border: none; cursor: pointer; color: var(--color-text-secondary); padding: 0.4rem; border-radius: 6px; display: inline-flex; align-items: center; justify-content: center; transition: all 0.2s; }
.btn-icon:hover { background: rgba(0,0,0,0.05); color: var(--color-text-primary); }

.form-group { display: flex; flex-direction: column; gap: 0.5rem; }
.form-group label { font-size: var(--font-size-xs); font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); }
.form-input, .form-textarea { width: 100%; padding: 0.75rem 1rem; border-radius: 8px; border: 1px solid var(--color-border); font-family: inherit; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
.checkbox-label { display: flex; align-items: center; gap: 0.75rem; cursor: pointer; font-size: 0.9rem; }
.form-checkbox { width: 18px; height: 18px; }
</style>
