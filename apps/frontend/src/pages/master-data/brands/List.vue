<template>
  <div class="page-layout">
    
    <BaseCatalog
      title="Gestión de Marcas"
      icon="branding_watermark"
      :breadcrumbs="[{ label: 'Catálogo', to: '/products/dashboard' }, { label: 'Marcas' }]"
      :items="brands"
      :is-loading="isLoading"
      :error="error"
      create-text="Nueva Marca"
      empty-icon="branding_watermark"
      empty-text="No hay marcas registradas"
      @clear-filters="resetFilters"
      @refresh="loadBrands"
      @click-item="editBrand"
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
          <input v-model="filters.search" type="text" placeholder="Nombre de marca..." />
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
        <th>Nombre de la Marca</th>
        <th class="text-right">Margen de Beneficio</th>
        <th class="text-center">Estado</th>
        <th class="align-right">Acciones</th>
      </template>

      <template #item="{ item }">
        <td><strong>{{ item.name }}</strong></td>
        <td class="text-right">
          <span class="markup-badge">{{ item.defaultMarkupPercentage }}%</span>
        </td>
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
      
      <div class="form-group mt-4">
        <label>Margen de Beneficio por defecto (%)</label>
        <div class="input-with-suffix">
          <input v-model.number="currentBrand.defaultMarkupPercentage" type="number" step="0.01" class="form-input" placeholder="0.00" />
          <span class="suffix">%</span>
        </div>
        <p class="form-hint">Este margen se aplicará a los productos de esta marca para calcular su precio de venta sugerido.</p>
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
      <p class="mt-2 text-muted">Esta acción solo se completará si la marca no tiene productos asociados.</p>
    </BaseDialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import BaseCatalog from '@/components/shared/BaseCatalog.vue'
import BaseDialog from '@/components/shared/BaseDialog.vue'
import { productApi } from '@/services/productApi'

const allBrands = ref([])
const isLoading = ref(false)
const error = ref('')
const filters = reactive({
  search: '',
  isActive: ''
})

const brands = computed(() => {
  let result = allBrands.value;
  if (filters.search) {
    const q = filters.search.toLowerCase();
    result = result.filter(b => b.name.toLowerCase().includes(q));
  }
  if (filters.isActive !== '') {
    const active = filters.isActive === 'true';
    result = result.filter(b => b.is_active === active);
  }
  return result;
});

const showModal = ref(false)
const showDeleteConfirm = ref(false)
const modalMode = ref('create')
const isSaving = ref(false)
const currentBrand = ref({ name: '', isActive: true, defaultMarkupPercentage: 0 })
const brandToDelete = ref(null)

async function loadBrands() {
  isLoading.value = true; error.value = '';
  try { 
    const res = await productApi.listBrands({}); 
    allBrands.value = res.data || (Array.isArray(res) ? res : []);
  }
  catch (err) { error.value = 'Error al cargar marcas.'; } 
  finally { isLoading.value = false }
}

function resetFilters() { filters.search = ''; filters.isActive = ''; }
function openCreateModal() { modalMode.value = 'create'; currentBrand.value = { name: '', isActive: true, defaultMarkupPercentage: 0 }; showModal.value = true; }
function editBrand(brand) { 
  modalMode.value = 'edit'; 
  currentBrand.value = { 
    id: brand.id, 
    name: brand.name, 
    isActive: brand.is_active,
    defaultMarkupPercentage: brand.defaultMarkupPercentage
  };
  showModal.value = true; 
}

async function saveBrand() {
  if (!currentBrand.value.name) { alert('El nombre es obligatorio'); return; }
  isSaving.value = true;
  try {
    const payload = { 
      name: currentBrand.value.name, 
      is_active: currentBrand.value.isActive,
      default_markup_percentage: currentBrand.value.defaultMarkupPercentage
    };
    if (modalMode.value === 'create') await productApi.createBrand(payload);
    else await productApi.updateBrand(currentBrand.value.id, payload);
    showModal.value = false; await loadBrands();
  } catch (err) { alert(err.message); } finally { isSaving.value = false }
}

async function toggleActive(brand) {
  try {
    const newStatus = !brand.is_active;
    await productApi.updateBrand(brand.id, { is_active: newStatus });
    await loadBrands();
  } catch (err) { alert('Error al cambiar estado'); }
}

function confirmDelete(brand) { brandToDelete.value = brand; showDeleteConfirm.value = true; }
async function deleteBrand() {
  if (!brandToDelete.value) return;
  try { 
    await productApi.deleteBrand(brandToDelete.value.id); 
    await loadBrands(); 
    showDeleteConfirm.value = false; 
  }
  catch (err) { 
    alert('No se puede eliminar: ' + (err.message || 'La marca tiene productos asociados.')); 
  }
}

onMounted(() => loadBrands())
</script>

<style scoped>
.page-layout { background-color: var(--color-background); min-height: 100vh; }
.action-buttons { display: flex; justify-content: flex-end; gap: 0.25rem; }
.btn-icon { background: transparent; border: none; cursor: pointer; color: var(--color-text-secondary); padding: 0.4rem; border-radius: 6px; display: inline-flex; align-items: center; justify-content: center; }
.btn-icon:hover { background: rgba(0,0,0,0.05); color: var(--color-text-primary); }

.markup-badge { background: #f1f5f9; color: var(--color-secondary); font-weight: 700; padding: 0.2rem 0.5rem; border-radius: 4px; font-size: 0.85rem; }

.form-group { display: flex; flex-direction: column; gap: 0.5rem; }
.form-group label { font-size: var(--font-size-xs); font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); }
.form-input { width: 100%; padding: 0.75rem 1rem; border-radius: 8px; border: 1px solid var(--color-border); font-family: inherit; }

.input-with-suffix { position: relative; display: flex; align-items: center; }
.input-with-suffix .form-input { padding-right: 3rem; }
.input-with-suffix .suffix { position: absolute; right: 1rem; color: var(--color-text-secondary); font-weight: 700; }

.form-hint { font-size: 0.75rem; color: var(--color-text-secondary); margin-top: 0.25rem; }

.checkbox-label { display: flex; align-items: center; gap: 0.75rem; cursor: pointer; font-size: 0.9rem; }
.form-checkbox { width: 18px; height: 18px; }

.text-right { text-align: right; }
.text-center { text-align: center; }
</style>
