<template>
  <div class="page-layout">
    
    <BaseCatalog
      title="Atributos de Producto"
      icon="tune"
      :breadcrumbs="[{ label: 'Catálogo', to: '/products/dashboard' }, { label: 'Atributos' }]"
      :items="attributes"
      :is-loading="isLoading"
      :error="error"
      :has-filters="hasFilters"
      create-text="Nuevo Atributo"
      empty-icon="tune"
      empty-text="No hay atributos registrados"
      @clear-filters="resetFilters"
      @refresh="loadAttributes"
      @click-item="openAttributeDetail"
    >
      <template #header-actions>
        <button @click="openCreateModal" class="btn btn-primary">
          <Plus :size="18" />
          <span>Nuevo Atributo</span>
        </button>
      </template>

      <template #filters>
        <div class="filter-group">
          <label>Búsqueda</label>
          <input v-model="search" type="text" placeholder="Nombre o código..." @input="onSearchInput" />
        </div>
      </template>

      <template #table-header>
        <th>Nombre del Atributo</th>
        <th>Código Técnico</th>
        <th class="text-center">Valores Registrados</th>
        <th class="align-right">Acciones</th>
      </template>

      <template #item="{ item }">
        <td><strong>{{ item.name }}</strong></td>
        <td><code class="code-badge">{{ item.code }}</code></td>
        <td class="text-center">
          <span class="status-badge status-secondary">
            {{ item.values ? item.values.length : 0 }} valores
          </span>
        </td>
        <td class="align-right" @click.stop>
          <div class="action-buttons">
            <button @click="editAttribute(item)" class="btn-icon" title="Editar"><Pencil :size="18" /></button>
            <button @click="confirmDelete(item)" class="btn-icon text-danger" title="Eliminar"><Trash2 :size="18" /></button>
          </div>
        </td>
      </template>
    </BaseCatalog>

    <!-- MODAL: CREAR/EDITAR ATRIBUTO -->
    <BaseDialog
      :show="showCreateModal"
      :title="modalMode === 'create' ? 'Nuevo Atributo' : 'Editar Atributo'"
      icon="tune"
      confirm-text="Guardar Atributo"
      :is-confirming="isSaving"
      @close="closeModal"
      @confirm="submitForm"
    >
      <AttributeForm ref="attributeFormRef" :attribute="selectedAttribute" :mode="modalMode" @submit="handleSubmit" />
    </BaseDialog>

    <!-- MODAL: CONFIRMAR ELIMINACIÓN -->
    <BaseDialog
      :show="showDeleteModal"
      title="Eliminar Atributo"
      icon="warning"
      confirm-text="Eliminar Definitivamente"
      confirm-class="btn-danger"
      :is-confirming="isDeleting"
      @close="showDeleteModal = false"
      @confirm="deleteAttribute"
    >
      <p>¿Está seguro de que desea eliminar el atributo <strong>{{ attributeToDelete?.name }}</strong>?</p>
      <p class="mt-2 text-muted">Esta acción no se puede deshacer y puede afectar a los productos vinculados.</p>
    </BaseDialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Plus, Pencil, Trash2 } from 'lucide-vue-next'
import BaseCatalog from '@/components/shared/BaseCatalog.vue'
import BaseDialog from '@/components/shared/BaseDialog.vue'
import AttributeForm from '@/components/master-data/AttributeForm.vue'
import { productApi } from '@/services/productApi'
import { useToastStore } from '@/stores/toast'

const toastStore = useToastStore()
const attributes = ref([])
const allAttributes = ref([])
const isLoading = ref(false)
const error = ref('')
const search = ref('')
const showCreateModal = ref(false)
const showDeleteModal = ref(false)
const modalMode = ref('create')
const selectedAttribute = ref(null)
const attributeToDelete = ref(null)
const isSaving = ref(false)
const isDeleting = ref(false)
const attributeFormRef = ref(null)

const hasFilters = computed(() => search.value.trim() !== '')

async function loadAttributes() {
  isLoading.value = true; error.value = '';
  try { 
    const res = await productApi.listAttributes({}); 
    allAttributes.value = res.data || (Array.isArray(res) ? res : []);
    onSearchInput();
  }
  catch (err) { 
    error.value = 'Error al cargar atributos.'; 
    console.error(err);
  } finally { isLoading.value = false; }
}

function onSearchInput() {
  const query = search.value.toLowerCase();
  attributes.value = allAttributes.value.filter(a => 
    a.name.toLowerCase().includes(query) || a.code.toLowerCase().includes(query)
  );
}

function resetFilters() { search.value = ''; attributes.value = allAttributes.value; }
function openCreateModal() { modalMode.value = 'create'; selectedAttribute.value = null; showCreateModal.value = true; }
function editAttribute(attr) { modalMode.value = 'edit'; selectedAttribute.value = attr; showCreateModal.value = true; }
function openAttributeDetail(attr) { editAttribute(attr); }
function confirmDelete(attr) { attributeToDelete.value = attr; showDeleteModal.value = true; }
function closeModal() { showCreateModal.value = false; selectedAttribute.value = null; }
function submitForm() { if (attributeFormRef.value) attributeFormRef.value.handleSubmit() }

async function handleSubmit(payload) {
  isSaving.value = true;
  try {
    if (modalMode.value === 'create') await productApi.createAttribute(payload);
    else await productApi.updateAttribute(payload.id, payload);
    closeModal(); await loadAttributes();
  } catch (err) { toastStore.addToast(err.message, 'error'); } finally { isSaving.value = false; }
}

async function deleteAttribute() {
  if (!attributeToDelete.value) return;
  isDeleting.value = true;
  try { await productApi.deleteAttribute(attributeToDelete.value.id); await loadAttributes(); showDeleteModal.value = false; }
  catch (err) { toastStore.addToast(err.message, 'error'); } finally { isDeleting.value = false; }
}

onMounted(() => loadAttributes())
</script>

<style scoped>
.page-layout { background-color: var(--color-background); min-height: 100vh; }
.code-badge { background: var(--color-background); padding: 0.25rem 0.5rem; border-radius: 4px; font-family: var(--font-family-mono); font-size: 0.75rem; color: var(--color-secondary); }
.action-buttons { display: flex; justify-content: flex-end; gap: 0.25rem; }
.btn-icon { background: transparent; border: none; cursor: pointer; color: var(--color-text-secondary); padding: 0.4rem; border-radius: 6px; display: inline-flex; align-items: center; justify-content: center; }
.btn-icon:hover { background: rgba(0,0,0,0.05); color: var(--color-text-primary); }
</style>