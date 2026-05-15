<template>
  <div class="page-layout">
    <BaseCatalog
      title="Gestión de Categorías"
      icon="category"
      :breadcrumbs="[{ label: 'Catálogo', to: '/products/dashboard' }, { label: 'Categorías' }]"
      :items="productGroups"
      :is-loading="isLoading"
      :error="error"
      :has-filters="hasFilters"
      create-text="Nueva Categoría"
      empty-icon="folder_off"
      empty-text="No hay categorías registradas en el sistema"
      @clear-filters="clearFilters"
      @refresh="loadGroups"
      @click-item="openGroupDetail"
    >
      <template #header-actions>
        <button @click="openCreateModal" class="btn btn-primary">
          <PlusSquare :size="18" />
          <span>Nueva Categoría</span>
        </button>
      </template>

      <template #filters>
        <div class="filter-group">
          <label>Búsqueda rápida</label>
          <input 
            v-model="filters.search" 
            type="text" 
            placeholder="Nombre de categoría..." 
          />
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
        <th>Nombre de la Categoría</th>
        <th>Categoría Padre</th>
        <th class="text-center">Tipo</th>
        <th class="text-center">Estado</th>
        <th class="align-right">Acciones</th>
      </template>

      <template #item="{ item }">
        <td><strong>{{ item.name }}</strong></td>
        <td class="parent-cell">
          <span v-if="item.parent_group_id" class="parent-badge">
            <CornerDownRight :size="14" />
            {{ getParentName(item.parent_group_id) }}
          </span>
          <span v-else class="text-muted">—</span>
        </td>
        <td class="text-center">
          <span class="status-badge-sm">{{ formatType(item.type) }}</span>
        </td>
        <td class="text-center">
          <span :class="['status-badge', item.is_active ? 'status-success' : 'status-secondary']">
            {{ item.is_active ? 'Activa' : 'Inactiva' }}
          </span>
        </td>
        <td class="align-right" @click.stop>
          <div class="action-buttons">
            <button @click="editGroup(item)" class="btn-icon" title="Editar"><Pencil :size="18" /></button>
            <button 
              @click="toggleActive(item)" 
              class="btn-icon" 
              :title="item.is_active ? 'Desactivar' : 'Activar'"
            >
              <component :is="item.is_active ? Ban : CheckCircle" :size="18" />
            </button>
            <button @click="confirmDelete(item)" class="btn-icon text-danger" title="Eliminar"><Trash2 :size="18" /></button>
          </div>
        </td>
      </template>
    </BaseCatalog>

    <!-- MODAL: CREAR/EDITAR CATEGORÍA -->
    <BaseDialog
      :show="showModal"
      :title="modalMode === 'create' ? 'Nueva Categoría' : 'Editar Categoría'"
      icon="category"
      confirm-text="Guardar Cambios"
      :is-confirming="isSaving"
      @close="showModal = false"
      @confirm="submitForm"
    >
      <ProductGroupForm 
        v-if="showModal"
        :key="modalMode + (currentGroup?.id || 'new')"
        ref="groupFormRef" 
        :product-group="currentGroup" 
        :mode="modalMode" 
        @submit="handleSubmit" 
      />
    </BaseDialog>

    <!-- MODAL: CONFIRMAR ELIMINACIÓN -->
    <BaseDialog
      :show="showDeleteConfirm"
      title="Eliminar Categoría"
      icon="warning"
      confirm-text="Eliminar Definitivamente"
      confirm-class="btn-danger"
      @close="showDeleteConfirm = false"
      @confirm="deleteGroup"
    >
      <p>¿Está seguro de que desea eliminar la categoría <strong>{{ groupToDelete?.name }}</strong>?</p>
      <p class="mt-2 text-muted">Esta acción solo se completará si la categoría no tiene productos asociados.</p>
    </BaseDialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { 
  PlusSquare, 
  CornerDownRight, 
  Pencil, 
  Ban, 
  CheckCircle, 
  Trash2
} from 'lucide-vue-next'
import BaseCatalog from '@/components/shared/BaseCatalog.vue'
import BaseDialog from '@/components/shared/BaseDialog.vue'
import ProductGroupForm from '@/components/master-data/ProductGroupForm.vue'
import { productApi } from '@/services/productApi'
import { useToastStore } from '@/stores/toast'

const toastStore = useToastStore()
const allGroups = ref([])
const isLoading = ref(false)
const error = ref('')
const filters = reactive({
  search: '',
  isActive: ''
})

const productGroups = computed(() => {
  let result = allGroups.value;
  if (filters.search) {
    const q = filters.search.toLowerCase();
    result = result.filter(g => g.name.toLowerCase().includes(q));
  }
  if (filters.isActive !== '') {
    const active = filters.isActive === 'true';
    result = result.filter(g => g.is_active === active);
  }
  return result;
});

const showModal = ref(false)
const modalMode = ref('create')
const isSaving = ref(false)
const showDeleteConfirm = ref(false)
const groupToDelete = ref(null)
const currentGroup = ref(null)
const groupFormRef = ref(null)

const hasFilters = computed(() => filters.search.trim() !== '' || filters.isActive !== '')

function getParentName(parentId) {
  const parent = allGroups.value.find(g => g.id === parentId)
  return parent ? parent.name : '—'
}

async function loadGroups() {
  isLoading.value = true; error.value = '';
  try { 
    const res = await productApi.listProductGroups({}); 
    allGroups.value = res.data || [];
  }
  catch (err) { error.value = err.message; } 
  finally { isLoading.value = false }
}

function clearFilters() {
  filters.search = '';
  filters.isActive = '';
}

function openCreateModal() { 
  modalMode.value = 'create'; 
  currentGroup.value = null; 
  showModal.value = true; 
}

function editGroup(group) { 
  modalMode.value = 'edit'; 
  currentGroup.value = group; 
  showModal.value = true; 
}

function openGroupDetail(group) { editGroup(group); }

function submitForm() {
  if (groupFormRef.value) {
    groupFormRef.value.handleSubmit()
  }
}

async function handleSubmit(payload) {
  isSaving.value = true;
  try {
    if (modalMode.value === 'create') {
      await productApi.createProductGroup(payload);
      toastStore.success('Categoría creada correctamente');
    } else {
      await productApi.updateProductGroup(payload.id, payload);
      toastStore.success('Categoría actualizada correctamente');
    }
    showModal.value = false;
    await loadGroups();
  } catch (err) {
    toastStore.addToast(err.message, 'error');
  } finally {
    isSaving.value = false;
  }
}

async function toggleActive(group) {
  try {
    await productApi.updateProductGroup(group.id, { isActive: !group.is_active });
    await loadGroups();
  } catch (err) { toastStore.addToast(err.message, 'error'); }
}

function confirmDelete(group) { 
  groupToDelete.value = group; 
  showDeleteConfirm.value = true; 
}

async function deleteGroup() {
  if (!groupToDelete.value) return;
  try { 
    await productApi.deleteProductGroup(groupToDelete.value.id); 
    await loadGroups(); 
    showDeleteConfirm.value = false;
  } catch (err) { 
    toastStore.addToast('No se puede eliminar: ' + (err.message || 'La categoría tiene productos asociados.'), 'error'); 
  }
}

function formatType(t) { return t === 'TANGIBLE' ? 'Producto' : 'Servicio'; }

onMounted(() => loadGroups())
</script>

<style scoped>
.page-layout { background-color: var(--color-background); min-height: 100vh; }
.action-buttons { display: flex; justify-content: flex-end; gap: 0.25rem; }

.parent-cell { color: var(--color-text-secondary); font-size: 0.875rem; }
.parent-badge { display: inline-flex; align-items: center; gap: 0.25rem; color: var(--color-text-secondary); }
.text-muted { color: var(--color-text-secondary); opacity: 0.5; }
.btn-icon { background: transparent; border: none; cursor: pointer; color: var(--color-text-secondary); padding: 0.4rem; border-radius: 6px; display: inline-flex; align-items: center; justify-content: center; }
.btn-icon:hover { background: rgba(0,0,0,0.05); color: var(--color-text-primary); }

.status-badge-sm { background: var(--color-background); padding: 0.1rem 0.5rem; border-radius: 4px; font-size: 0.7rem; font-weight: 600; color: var(--color-text-secondary); }

.form-group { display: flex; flex-direction: column; gap: 0.5rem; }
.form-group label { font-size: var(--font-size-xs); font-weight: 600; text-transform: uppercase; color: var(--color-text-secondary); }
.form-input { width: 100%; padding: 0.75rem 1rem; border-radius: 8px; border: 1px solid var(--color-border); font-family: inherit; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }

.checkbox-label { display: flex; align-items: center; gap: 0.75rem; cursor: pointer; font-size: 0.9rem; }
.mt-8 { margin-top: 2rem; }
.text-center { text-align: center; }
.align-right { text-align: right; }
</style>
