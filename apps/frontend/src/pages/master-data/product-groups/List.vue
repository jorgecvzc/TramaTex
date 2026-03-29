<template>
  <BaseCatalog
    title="Gestión de Categorías"
    icon="category"
    :breadcrumbs="[{ label: 'Datos Maestros', to: '/master-data/product-groups' }, { label: 'Categorías' }]"
    :items="productGroups"
    :is-loading="isLoading"
    :error="error"
    create-text="Nueva Categoría"
    empty-icon="folder_off"
    empty-text="No hay categorías registradas en el sistema"
    @clear-filters="loadGroups"
    @refresh="loadGroups"
    @click-item="editGroup"
  >
    <template #filters>
      <div class="filter-group">
        <label>Búsqueda rápida</label>
        <input 
          v-model="search" 
          type="text" 
          placeholder="Filtrar por nombre de categoría..." 
          @input="onSearchInput"
        />
      </div>

      <div class="filter-group">
        <label>Estado</label>
        <select v-model="statusFilter" @change="onSearchInput">
          <option value="">Cualquier estado</option>
          <option value="true">Activas</option>
          <option value="false">Inactivas</option>
        </select>
      </div>
    </template>

    <template #header-actions>
      <button @click="openCreateModal" class="btn btn-primary">
        <span class="material-symbols-outlined">add_box</span>
        Nueva Categoría
      </button>
    </template>

    <template #table-header>
      <th>Nombre de la Categoría</th>
      <th class="text-center">Tipo</th>
      <th class="text-center">Estado</th>
      <th class="align-right">Acciones</th>
    </template>

    <template #item="{ item }">
      <td><strong>{{ item.name }}</strong></td>
      <td class="text-center">
        <span class="status-badge-sm">{{ formatType(item.type) }}</span>
      </td>
      <td class="text-center">
        <span :class="['status-pill', item.is_active ? 'status-active' : 'status-inactive']">
          {{ item.is_active ? 'Activa' : 'Inactiva' }}
        </span>
      </td>
      <td class="align-right" @click.stop>
        <div class="action-buttons">
          <button @click="editGroup(item)" class="btn-icon" title="Editar"><span class="material-symbols-outlined">edit</span></button>
          <button 
            @click="toggleActive(item)" 
            class="btn-icon" 
            :title="item.is_active ? 'Desactivar' : 'Activar'"
            :class="{ 'text-warning': item.is_active }"
          >
            <span class="material-symbols-outlined">{{ item.is_active ? 'block' : 'check_circle' }}</span>
          </button>
          <button @click="confirmDelete(item)" class="btn-icon text-danger" title="Eliminar"><span class="material-symbols-outlined">delete</span></button>
        </div>
      </td>
    </template>
  </BaseCatalog>

  <!-- Modals (Standardized) -->
  <Transition name="fade">
    <div v-if="showModal" class="modal-backdrop">
      <div class="modal card">
        <div class="modal-header">
          <span class="material-symbols-outlined">edit_square</span>
          <h2>{{ modalMode === 'create' ? 'Nueva Categoría' : 'Editar Categoría' }}</h2>
        </div>
        
        <div class="modal-body">
          <div class="form-group">
            <label>Nombre de la Categoría</label>
            <input v-model="currentGroup.name" type="text" placeholder="Ej: Hilos, Tejidos..." required @keyup.enter="saveGroup" />
          </div>
          <div class="form-row mt-4">
            <div class="form-group">
              <label>Tipo</label>
              <select v-model="currentGroup.type">
                <option value="TANGIBLE">Producto Físico</option>
                <option value="SERVICE">Servicio</option>
              </select>
            </div>
            <div class="form-group">
              <label class="checkbox-label mt-8">
                <input v-model="currentGroup.isActive" type="checkbox" />
                <span>Categoría activa</span>
              </label>
            </div>
          </div>
        </div>
        
        <div class="modal-actions">
          <button @click="showModal = false" class="btn btn-outline">Cancelar</button>
          <button @click="saveGroup" class="btn btn-primary" :disabled="isSaving">
            <span class="material-symbols-outlined">{{ isSaving ? 'sync' : 'save' }}</span>
            <span>Confirmar</span>
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import BaseCatalog from '@/components/shared/BaseCatalog.vue'
import { productApi } from '@/services/productApi'

const productGroups = ref([])
const allGroups = ref([])
const isLoading = ref(false)
const error = ref('')
const search = ref('')
const statusFilter = ref('')
const showModal = ref(false)
const modalMode = ref('create')
const isSaving = ref(false)
const currentGroup = ref({ name: '', type: 'TANGIBLE', isActive: true })

async function loadGroups() {
  isLoading.value = true; error.value = '';
  try { 
    const res = await productApi.listProductGroups({}); 
    allGroups.value = res.data || [];
    productGroups.value = allGroups.value;
  }
  catch (err) { error.value = err.message; } 
  finally { isLoading.value = false }
}

function onSearchInput() {
  const query = search.value.toLowerCase();
  productGroups.value = allGroups.value.filter(g => {
    const matchesSearch = g.name.toLowerCase().includes(query);
    const matchesStatus = statusFilter.value === '' || String(g.is_active) === statusFilter.value;
    return matchesSearch && matchesStatus;
  });
}

function openCreateModal() { modalMode.value = 'create'; currentGroup.value = { name: '', type: 'TANGIBLE', isActive: true }; showModal.value = true; }
function editGroup(group) { modalMode.value = 'edit'; currentGroup.value = { ...group, isActive: group.is_active }; showModal.value = true; }

async function saveGroup() {
  if (!currentGroup.value.name) return;
  isSaving.value = true;
  try {
    if (modalMode.value === 'create') await productApi.createProductGroup(currentGroup.value);
    else await productApi.updateProductGroup(currentGroup.value.id, currentGroup.value);
    showModal.value = false; await loadGroups();
  } catch (err) { alert(err.message); } finally { isSaving.value = false }
}

async function toggleActive(group) {
  try {
    await productApi.updateProductGroup(group.id, { isActive: !group.is_active });
    await loadGroups();
  } catch (err) { alert(err.message); }
}

async function confirmDelete(group) {
  if (!confirm(`¿Eliminar la categoría ${group.name}?`)) return;
  try { await productApi.deleteProductGroup(group.id); await loadGroups(); } catch (err) { alert(err.message); }
}

function formatType(t) { return t === 'TANGIBLE' ? 'Producto' : 'Servicio'; }

onMounted(() => loadGroups())
</script>

<style scoped>
.action-buttons { display: flex; justify-content: flex-end; gap: 0.25rem; }
.btn-icon { background: transparent; border: none; cursor: pointer; color: var(--color-text-secondary); padding: 0.4rem; border-radius: 6px; }
.btn-icon:hover { background: rgba(0,0,0,0.05); color: var(--color-text-primary); }
.text-danger:hover { color: var(--color-error); }
.text-warning:hover { color: #d97706; }

.status-badge-sm { background: var(--color-background); padding: 0.1rem 0.5rem; border-radius: 4px; font-size: 0.7rem; font-weight: 600; color: var(--color-text-secondary); }

.form-group { display: flex; flex-direction: column; gap: 0.5rem; }
.form-group label { font-size: var(--font-size-xs); font-weight: 600; text-transform: uppercase; color: var(--color-text-secondary); }
input[type="text"], select { width: 100%; padding: 0.75rem 1rem; border-radius: var(--border-radius-md); border: 1px solid var(--color-border); font-size: var(--font-size-sm); }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }

.checkbox-label { flex-direction: row; align-items: center; gap: 0.75rem; cursor: pointer; }
.mt-8 { margin-top: 2rem; }
.text-center { text-align: center; }
.align-right { text-align: right; }
</style>