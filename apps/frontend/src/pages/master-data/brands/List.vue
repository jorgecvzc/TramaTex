<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">Datos Maestros / Marcas</p>
          <h1>Gestión de Marcas</h1>
          <p class="subtitle">Administra las marcas disponibles para los productos</p>
        </div>
        <button @click="openCreateModal" class="btn btn-primary">
          + Nueva Marca
        </button>
      </header>

      <!-- Loading State -->
      <div v-if="isLoading" class="loading-container">
        <div class="spinner"></div>
        <p>Cargando marcas...</p>
      </div>

      <!-- Error State -->
      <div v-if="error" class="alert alert-error">
        <span class="alert-icon">✗</span>
        <div class="alert-content">
          <strong>Error al cargar marcas</strong>
          <p>{{ error }}</p>
        </div>
        <button @click="loadBrands" class="btn btn-secondary btn-sm">Reintentar</button>
      </div>

      <!-- List -->
      <div v-if="!isLoading && !error" class="card">
        <table class="data-table">
          <thead>
            <tr>
              <th>Nombre</th>
              <th>ID</th>
              <th class="text-center">Beneficio (%)</th>
              <th class="text-center">Productos</th>
              <th class="text-right">Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="brand in brands" :key="brand.id">
              <td>
                <strong>{{ brand.name }}</strong>
              </td>
              <td>
                <code class="id-badge">{{ brand.id.substring(0, 8) }}...</code>
              </td>
              <td class="text-center">
                <span class="badge badge-info">{{ brand.defaultMarkupPercentage }}%</span>
              </td>
              <td class="text-center">
                <span class="badge">-</span>
              </td>
              <td class="text-right">
                <button @click="editBrand(brand)" class="btn btn-sm btn-secondary mr-2">Editar</button>
                <button @click="deleteBrand(brand)" class="btn btn-sm btn-danger">Eliminar</button>
              </td>
            </tr>
            <tr v-if="brands.length === 0">
              <td colspan="5" class="text-center empty-state">
                No hay marcas registradas. Crea una nueva para comenzar.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="showCreateModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content">
        <div class="modal-header">
          <h2>{{ modalMode === 'create' ? 'Nueva Marca' : 'Editar Marca' }}</h2>
          <button @click="closeModal" class="modal-close">&times;</button>
        </div>
        <div class="modal-body">
          <BrandForm 
            ref="brandFormRef"
            :brand="selectedBrand" 
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import BrandForm from '@/components/master-data/BrandForm.vue'
import { productApi } from '@/services/productApi'

const brands = ref([])
const isLoading = ref(false)
const error = ref('')
const showCreateModal = ref(false)
const modalMode = ref('create')
const selectedBrand = ref(null)
const isSaving = ref(false)
const brandFormRef = ref(null)

async function loadBrands() {
  isLoading.value = true
  error.value = ''
  try {
    const response = await productApi.listBrands({})
    brands.value = response.data || []
  } catch (err) {
    error.value = err.message || 'No se pudieron cargar las marcas'
    console.error('Error loading brands:', err)
  } finally {
    isLoading.value = false
  }
}

function openCreateModal() {
  modalMode.value = 'create'
  selectedBrand.value = null
  showCreateModal.value = true
}

function editBrand(brand) {
  modalMode.value = 'edit'
  selectedBrand.value = brand
  showCreateModal.value = true
}

function viewProducts(brand) {
  console.log('View products for brand:', brand)
}

async function deleteBrand(brand) {
  if (!confirm(`¿Estás seguro de eliminar la marca "${brand.name}"? Esta acción no se puede deshacer.`)) {
    return
  }
  
  isLoading.value = true
  error.value = ''
  
  try {
    await productApi.deleteBrand(brand.id)
    await loadBrands()
  } catch (err) {
    error.value = err.message || 'No se pudo eliminar la marca'
    console.error('Error deleting brand:', err)
  } finally {
    isLoading.value = false
  }
}

function closeModal() {
  showCreateModal.value = false
  selectedBrand.value = null
  modalMode.value = 'create'
}

function submitForm() {
  if (brandFormRef.value) {
    brandFormRef.value.handleSubmit()
  }
}

async function handleSubmit(payload) {
  isSaving.value = true
  error.value = ''
  
  try {
    if (modalMode.value === 'create') {
      await productApi.createBrand(payload)
    } else {
      await productApi.updateBrand(payload.id, payload)
    }
    
    closeModal()
    await loadBrands()
  } catch (err) {
    error.value = err.message || `No se pudo ${modalMode.value === 'create' ? 'crear' : 'actualizar'} la marca`
    console.error('Error saving brand:', err)
  } finally {
    isSaving.value = false
  }
}

onMounted(() => {
  loadBrands()
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

.btn-danger {
  background: #ef4444;
  color: #ffffff;
  font-weight: 500;
}

.btn-danger:hover {
  background: #dc2626;
  box-shadow: 0 4px 6px rgba(239, 68, 68, 0.3);
  transform: translateY(-1px);
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

.id-badge {
  background: #f1f5f9;
  padding: 0.25rem 0.5rem;
  border-radius: 0.25rem;
  font-size: 0.75rem;
  font-family: 'Monaco', 'Courier New', monospace;
  color: #64748b;
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

.text-muted {
  color: #64748b;
}

.mb-4 {
  margin-bottom: 1rem;
}

.mr-2 {
  margin-right: 0.5rem;
}
</style>
