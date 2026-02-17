<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">Datos Maestros / Atributos</p>
          <h1>Gestión de Atributos</h1>
          <p class="subtitle">Administra los atributos configurables de productos</p>
        </div>
        <button @click="openCreateModal" class="btn btn-primary">
          + Nuevo Atributo
        </button>
      </header>

      <!-- Loading State -->
      <div v-if="isLoading" class="loading-container">
        <div class="spinner"></div>
        <p>Cargando atributos...</p>
      </div>

      <!-- Error State -->
      <div v-if="error" class="alert alert-error">
        <span class="alert-icon">✗</span>
        <div class="alert-content">
          <strong>Error al cargar atributos</strong>
          <p>{{ error }}</p>
        </div>
        <button @click="loadAttributes" class="btn btn-secondary btn-sm">Reintentar</button>
      </div>

      <!-- List -->
      <div v-if="!isLoading && !error" class="card">
        <table class="data-table">
          <thead>
            <tr>
              <th>Nombre</th>
              <th>Código</th>
              <th class="text-center">Valores</th>
              <th class="text-right">Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="attr in attributes" :key="attr.id">
              <td>
                <strong>{{ attr.name }}</strong>
              </td>
              <td>
                <code class="code-badge">{{ attr.code }}</code>
              </td>
              <td class="text-center">
                <span class="badge">
                  {{ attr.values ? attr.values.length : 0 }} valores
                </span>
              </td>
              <td class="text-right">
                <button @click="editAttribute(attr)" class="btn btn-sm btn-secondary mr-2">Editar</button>
                <button @click="confirmDelete(attr)" class="btn btn-sm btn-danger">Eliminar</button>
              </td>
            </tr>
            <tr v-if="attributes.length === 0">
              <td colspan="4" class="text-center empty-state">
                No hay atributos registrados.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="showCreateModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content modal-large">
        <div class="modal-header">
          <h2>{{ modalMode === 'create' ? 'Nuevo Atributo' : 'Editar Atributo' }}</h2>
          <button @click="closeModal" class="modal-close">&times;</button>
        </div>
        <div class="modal-body">
          <!-- Error message in modal -->
          <div v-if="error" class="alert alert-error mb-3">
            <span class="alert-icon">✗</span>
            <div class="alert-content">
              <strong>Error</strong>
              <p>{{ error }}</p>
            </div>
          </div>
          
          <AttributeForm 
            ref="attributeFormRef"
            :attribute="selectedAttribute" 
            :mode="modalMode"
            @submit="handleSubmit"
          />
        </div>
        <div class="modal-footer">
          <button @click="closeModal" class="btn btn-secondary">Cancelar</button>
          <button @click="submitForm" class="btn btn-primary" :disabled="isSaving">
            {{ isSaving ? 'Guardando...' : (modalMode === 'create' ? 'Crear' : 'Actualizar') }}
          </button>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="modal-overlay" @click.self="showDeleteModal = false">
      <div class="modal-content">
        <div class="modal-header">
          <h2>Confirmar eliminación</h2>
          <button @click="showDeleteModal = false" class="modal-close">&times;</button>
        </div>
        <div class="modal-body">
          <p>¿Estás seguro de que deseas eliminar el atributo <strong>{{ attributeToDelete?.name }}</strong>?</p>
          <p class="text-muted mt-2">Esta acción no se puede deshacer.</p>
        </div>
        <div class="modal-footer">
          <button @click="showDeleteModal = false" class="btn btn-secondary">Cancelar</button>
          <button @click="deleteAttribute" class="btn btn-danger" :disabled="isDeleting">
            {{ isDeleting ? 'Eliminando...' : 'Eliminar' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import AttributeForm from '@/components/master-data/AttributeForm.vue'
import { productApi } from '@/services/productApi'

const attributes = ref([])
const isLoading = ref(false)
const error = ref('')
const showCreateModal = ref(false)
const showDeleteModal = ref(false)
const modalMode = ref('create')
const selectedAttribute = ref(null)
const attributeToDelete = ref(null)
const isSaving = ref(false)
const isDeleting = ref(false)
const attributeFormRef = ref(null)

async function loadAttributes() {
  isLoading.value = true
  error.value = ''
  try {
    const response = await productApi.listAttributes({})
    attributes.value = response.data || []
  } catch (err) {
    error.value = err.message || 'No se pudieron cargar los atributos'
    console.error('Error loading attributes:', err)
  } finally {
    isLoading.value = false
  }
}

function openCreateModal() {
  modalMode.value = 'create'
  selectedAttribute.value = null
  showCreateModal.value = true
}

function editAttribute(attribute) {
  modalMode.value = 'edit'
  selectedAttribute.value = attribute
  showCreateModal.value = true
}

function confirmDelete(attribute) {
  attributeToDelete.value = attribute
  showDeleteModal.value = true
}

function closeModal() {
  showCreateModal.value = false
  selectedAttribute.value = null
  modalMode.value = 'create'
}

function submitForm() {
  if (attributeFormRef.value) {
    attributeFormRef.value.handleSubmit()
  }
}

async function handleSubmit(payload) {
  isSaving.value = true
  error.value = ''
  
  try {
    if (modalMode.value === 'create') {
      await productApi.createAttribute(payload)
    } else {
      await productApi.updateAttribute(payload.id, payload)
    }
    
    closeModal()
    await loadAttributes()
  } catch (err) {
    error.value = err.message || `No se pudo ${modalMode.value === 'create' ? 'crear' : 'actualizar'} el atributo`
    console.error('Error saving attribute:', err)
  } finally {
    isSaving.value = false
  }
}

async function deleteAttribute() {
  if (!attributeToDelete.value) return
  
  isDeleting.value = true
  error.value = ''
  
  try {
    // Note: Delete endpoint not yet implemented in backend
    // await productApi.deleteAttribute(attributeToDelete.value.id)
    console.warn('Delete attribute not implemented in backend yet')
    error.value = 'La funcionalidad de eliminación aún no está implementada en el backend'
    showDeleteModal.value = false
  } catch (err) {
    error.value = err.message || 'No se pudo eliminar el atributo'
    console.error('Error deleting attribute:', err)
  } finally {
    isDeleting.value = false
    attributeToDelete.value = null
  }
}

onMounted(() => {
  loadAttributes()
})
</script>

<style scoped>
.dashboard {
  min-height: 100vh;
  background-color: #f1f5f9;
  font-family: 'Inter', sans-serif;
}

.dashboard-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.page-header {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

.page-header h1 {
  color: #1b3a6b;
  margin: 0.25rem 0 0;
}

.breadcrumb {
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #64748b;
  margin: 0;
}

.subtitle {
  color: #64748b;
  margin: 0.5rem 0 0;
  font-size: 0.95rem;
}

.card {
  background-color: #ffffff;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
  border: 1px solid #e2e8f0;
}

.btn {
  border: none;
  border-radius: 8px;
  padding: 0.6rem 1rem;
  font-size: 0.85rem;
  cursor: pointer;
  transition: background 0.2s ease, box-shadow 0.2s ease, transform 0.2s ease;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.btn-primary {
  background: #f4c430;
  color: #1b3a6b;
  font-weight: 600;
}

.btn-primary:hover {
  background: #e6b82a;
  box-shadow: 0 4px 6px rgba(244, 196, 48, 0.3);
  transform: translateY(-1px);
}

.btn-secondary {
  background: #ffffff;
  color: #1b3a6b;
  border: 1px solid #dde3ed;
  font-weight: 500;
}

.btn-secondary:hover {
  background: #f8fafc;
  border-color: #cbd5e1;
}

.btn-sm {
  padding: 0.4rem 0.75rem;
  font-size: 0.8rem;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  color: #64748b;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e2e8f0;
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.alert {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  border-radius: 8px;
  margin-bottom: 1rem;
}

.alert-error {
  background: #fef2f2;
  border: 1px solid #fecaca;
  color: #991b1b;
}

.alert-icon {
  font-size: 1.25rem;
  font-weight: bold;
}

.alert-content {
  flex: 1;
}

.alert-content strong {
  display: block;
  margin-bottom: 0.25rem;
}

.alert-content p {
  margin: 0;
  font-size: 0.875rem;
}

.mb-3 {
  margin-bottom: 1.5rem;
}

.filters {
  display: flex;
  gap: 1rem;
  padding: 0.5rem 0;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.filter-group label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #64748b;
}

.filter-select {
  padding: 0.5rem 0.75rem;
  border: 1px solid #e2e8f0;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  color: #1e293b;
  background: white;
  cursor: pointer;
}

.filter-select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table thead {
  background: #f8fafc;
  border-bottom: 2px solid #e2e8f0;
}

.data-table th,
.data-table td {
  padding: 1rem;
  text-align: left;
}

.data-table th {
  font-weight: 600;
  color: #1e293b;
  font-size: 0.875rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.data-table tbody tr {
  border-bottom: 1px solid #e2e8f0;
  transition: background-color 0.15s;
}

.data-table tbody tr:hover {
  background: #f8fafc;
}

.code-badge {
  background: #f1f5f9;
  padding: 0.25rem 0.5rem;
  border-radius: 0.25rem;
  font-size: 0.75rem;
  font-family: 'Monaco', 'Courier New', monospace;
  color: #64748b;
}

.scope-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 1rem;
  font-size: 0.75rem;
  font-weight: 600;
}

.scope-badge.generic {
  background: #e0e7ff;
  color: #4338ca;
}

.scope-badge.brand {
  background: #fef3c7;
  color: #b45309;
}

.scope-badge.group {
  background: #d1fae5;
  color: #047857;
}

.scope-badge.both {
  background: #e9d5ff;
  color: #7c3aed;
}

.badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 1rem;
  font-size: 0.75rem;
  font-weight: 600;
  background: #e2e8f0;
  color: #64748b;
}

.text-center {
  text-align: center;
}

.text-right {
  text-align: right;
}

.text-muted {
  color: #64748b;
}

.empty-state {
  padding: 3rem 1rem;
  color: #94a3b8;
  font-style: italic;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 0.5rem;
  max-width: 500px;
  width: 90%;
  max-height: 90vh;
  overflow: auto;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
}

.modal-large {
  max-width: 700px;
}

.modal-header {
  padding: 1.5rem;
  border-bottom: 1px solid #e2e8f0;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.modal-header h2 {
  margin: 0;
  font-size: 1.25rem;
  color: #1e293b;
}

.modal-close {
  background: none;
  border: none;
  font-size: 1.5rem;
  color: #94a3b8;
  cursor: pointer;
  padding: 0;
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 0.25rem;
  transition: all 0.15s;
}

.modal-close:hover {
  background: #f1f5f9;
  color: #1e293b;
}

.modal-body {
  padding: 1.5rem;
}

.modal-footer {
  padding: 1.5rem;
  border-top: 1px solid #e2e8f0;
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
}

.mb-4 {
  margin-bottom: 1rem;
}

.mt-2 {
  margin-top: 0.5rem;
}

.mr-2 {
  margin-right: 0.5rem;
}

.btn-danger {
  background: #fee2e2;
  color: #dc2626;
  border: 1px solid #fecaca;
  font-weight: 500;
}

.btn-danger:hover:not(:disabled) {
  background: #fecaca;
  border-color: #fca5a5;
}

.btn-danger:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
