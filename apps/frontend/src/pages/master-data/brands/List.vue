<template>
  <div class="page-layout">
    <Navbar style="z-index: 2000;" />
    
    <BaseCatalog
      title="Gestión de Marcas"
      :breadcrumbs="[{ label: 'Operaciones', to: '/products' }, { label: 'Marcas' }]"
      :items="brands"
      :is-loading="isLoading"
      :error="error"
      create-text="Nueva Marca"
      empty-icon="branding_watermark"
      empty-text="No hay marcas registradas"
      @clear-filters="resetFilters"
      @refresh="loadBrands"
    >
      <template #header-actions>
        <button @click="openCreateModal" class="btn btn-primary">
          <span class="material-symbols-outlined">add</span>
          <span>Nueva Marca</span>
        </button>
      </template>

      <template #filters>
        <div class="filter-group">
          <label>Búsqueda</label>
          <input v-model="search" type="text" placeholder="Nombre de marca..." @input="onSearchInput" />
        </div>
      </template>

      <template #table-header>
        <th>Nombre de la Marca</th>
        <th class="text-center">Estado</th>
        <th class="align-right">Acciones</th>
      </template>

      <template #item="{ item }">
        <td><strong>{{ item.name }}</strong></td>
        <td class="text-center">
          <span :class="['status-badge', item.is_active ? 'status-success' : 'status-secondary']">
            {{ item.is_active ? 'Activa' : 'Inactiva' }}
          </span>
        </td>
        <td class="align-right" @click.stop>
          <div class="action-buttons">
            <button @click="editBrand(item)" class="btn-icon" title="Editar"><span class="material-symbols-outlined">edit</span></button>
            <button 
              @click="toggleActive(item)" 
              class="btn-icon" 
              :title="item.is_active ? 'Desactivar' : 'Activar'"
            >
              <span class="material-symbols-outlined">{{ item.is_active ? 'block' : 'check_circle' }}</span>
            </button>
            <button @click="confirmDelete(item)" class="btn-icon text-danger" title="Eliminar"><span class="material-symbols-outlined">delete</span></button>
          </div>
        </td>
      </template>
    </BaseCatalog>

    <!-- MODAL: CREAR/EDITAR MARCA -->
    <BaseDialog
      :show="showModal"
      :title="modalMode === 'create' ? 'Nueva Marca' : 'Editar Marca'"
      icon="branding_watermark"
      confirm-text="Guardar Marca"
      :is-confirming="isSaving"
      @close="showModal = false"
      @confirm="saveBrand"
    >
      <div class="form-group">
        <label>Nombre de la Marca</label>
        <input v-model="currentBrand.name" type="text" class="form-input" placeholder="Ej: Adidas, Nike..." required @keyup.enter="saveBrand" />
      </div>
      <div class="form-group mt-6">
        <label class="checkbox-label">
          <input v-model="currentBrand.isActive" type="checkbox" class="form-checkbox" />
          <span>Marca activa y disponible para productos</span>
        </label>
      </div>
    </BaseDialog>

    <!-- MODAL: CONFIRMAR ELIMINACIÓN -->
    <BaseDialog
      :show="showDeleteConfirm"
      title="Eliminar Marca"
      icon="warning"
      confirm-text="Eliminar Definitivamente"
      confirm-class="btn-danger"
      @close="showDeleteConfirm = false"
      @confirm="deleteBrand"
    >
      <p>¿Está seguro de que desea eliminar la marca <strong>{{ brandToDelete?.name }}</strong>?</p>
      <p class="mt-2 text-muted">Esta acción puede afectar a la visualización de los productos asociados.</p>
    </BaseDialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import BaseCatalog from '@/components/shared/BaseCatalog.vue'
import BaseDialog from '@/components/shared/BaseDialog.vue'
import { productApi } from '@/services/productApi'

const brands = ref([])
const allBrands = ref([])
const isLoading = ref(false)
const error = ref('')
const search = ref('')
const showModal = ref(false)
const showDeleteConfirm = ref(false)
const modalMode = ref('create')
const isSaving = ref(false)
const currentBrand = ref({ name: '', isActive: true })
const brandToDelete = ref(null)

async function loadBrands() {
  isLoading.value = true; error.value = '';
  try { 
    const res = await productApi.listBrands({}); 
    allBrands.value = res.data || (Array.isArray(res) ? res : []);
    brands.value = allBrands.value;
  }
  catch (err) { error.value = 'Error al cargar marcas.'; } 
  finally { isLoading.value = false }
}

function onSearchInput() {
  const query = search.value.toLowerCase();
  brands.value = allBrands.value.filter(b => b.name.toLowerCase().includes(query));
}

function resetFilters() { search.value = ''; brands.value = allBrands.value; }
function openCreateModal() { modalMode.value = 'create'; currentBrand.value = { name: '', isActive: true }; showModal.value = true; }
function editBrand(brand) { 
  modalMode.value = 'edit'; 
  currentBrand.value = { id: brand.id, name: brand.name, isActive: brand.is_active };
  showModal.value = true; 
}

async function saveBrand() {
  if (!currentBrand.value.name) { alert('El nombre es obligatorio'); return; }
  isSaving.value = true;
  try {
    const payload = { name: currentBrand.value.name, isActive: currentBrand.value.isActive };
    if (modalMode.value === 'create') await productApi.createBrand(payload);
    else await productApi.updateBrand(currentBrand.value.id, payload);
    showModal.value = false; await loadBrands();
  } catch (err) { alert(err.message); } finally { isSaving.value = false }
}

async function toggleActive(brand) {
  try {
    const newStatus = !brand.is_active;
    await productApi.updateBrand(brand.id, { isActive: newStatus });
    await loadBrands();
  } catch (err) { alert('Error al cambiar estado'); }
}

function confirmDelete(brand) { brandToDelete.value = brand; showDeleteConfirm.value = true; }
async function deleteBrand() {
  if (!brandToDelete.value) return;
  try { await productApi.deleteBrand(brandToDelete.value.id); await loadBrands(); showDeleteConfirm.value = false; }
  catch (err) { alert(err.message); }
}

onMounted(() => loadBrands())
</script>

<style scoped>
.page-layout { background-color: var(--color-background); min-height: 100vh; }
.action-buttons { display: flex; justify-content: flex-end; gap: 0.25rem; }
.btn-icon { background: transparent; border: none; cursor: pointer; color: var(--color-text-secondary); padding: 0.4rem; border-radius: 6px; display: inline-flex; align-items: center; justify-content: center; }
.btn-icon:hover { background: rgba(0,0,0,0.05); color: var(--color-text-primary); }

.form-group { display: flex; flex-direction: column; gap: 0.5rem; }
.form-group label { font-size: var(--font-size-xs); font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); }
.form-input { width: 100%; padding: 0.75rem 1rem; border-radius: 8px; border: 1px solid var(--color-border); font-family: inherit; }

.checkbox-label { display: flex; align-items: center; gap: 0.75rem; cursor: pointer; font-size: 0.9rem; }
.form-checkbox { width: 18px; height: 18px; }
</style>
