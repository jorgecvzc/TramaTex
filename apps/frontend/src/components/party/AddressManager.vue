<template>
  <div class="address-manager">
    <div class="manager-header">
      <h3 class="flex items-center gap-2"><MapPin :size="20" /> Direcciones Registradas</h3>
      <button class="btn btn-primary btn-sm" @click="openCreateModal">
        <Plus :size="16" /> Añadir Dirección
      </button>
    </div>

    <div v-if="addresses.length === 0" class="empty-state">
      <p class="text-muted italic">No hay direcciones registradas para esta entidad.</p>
    </div>

    <div v-else class="address-list">
      <div v-for="addr in addresses" :key="addr.id" class="address-row" :class="{ 'is-primary': addr.is_primary }">
        <div class="address-main-info">
          <div class="row-indicator" v-if="addr.is_primary" title="Dirección Principal">
            <CheckCircle :size="16" class="text-success" />
          </div>
          <div class="address-text">
            <span class="street">{{ addr.street }}</span>
            <span class="location">{{ addr.postal_code }} {{ addr.city }}, {{ addr.province }}</span>
            <span class="country text-muted">{{ addr.country }}</span>
          </div>
        </div>
        <div class="address-actions">
          <button class="btn-icon" @click="editAddress(addr)" title="Editar"><Pencil :size="18" /></button>
          <button class="btn-icon text-danger" @click="promptDelete(addr)" title="Eliminar"><Trash2 :size="18" /></button>
        </div>
      </div>
    </div>

    <!-- MODAL: CREAR/EDITAR DIRECCIÓN -->
    <BaseDialog
      :show="showModal"
      :title="editingId ? 'Editar Dirección' : 'Nueva Dirección'"
      icon="map_pin"
      confirm-text="Guardar Dirección"
      :is-confirming="isSaving"
      @close="showModal = false"
      @confirm="saveAddress"
    >
      <div class="form-grid">
        <div class="form-group full-width">
          <label>Calle y Número *</label>
          <input v-model="formData.street" type="text" class="form-input" placeholder="Ej: Av. Constitución 45, 2ºB" required />
        </div>
        <div class="form-group">
          <label>Ciudad *</label>
          <input v-model="formData.city" type="text" class="form-input" required />
        </div>
        <div class="form-group">
          <label>Código Postal *</label>
          <input v-model="formData.postal_code" type="text" class="form-input" required />
        </div>
        <div class="form-group">
          <label>Provincia *</label>
          <input v-model="formData.province" type="text" class="form-input" required />
        </div>
        <div class="form-group">
          <label>País *</label>
          <input v-model="formData.country" type="text" class="form-input" required />
        </div>
        <div class="form-group full-width">
          <label class="checkbox-label mt-2">
            <input v-model="formData.is_primary" type="checkbox" />
            <span>Establecer como dirección principal</span>
          </label>
        </div>
      </div>
    </BaseDialog>

    <!-- MODAL DE CONFIRMACIÓN -->
    <BaseDialog
      :show="confirmDelete.show"
      title="Eliminar Dirección"
      icon="warning"
      confirm-text="Eliminar"
      confirm-class="btn-danger"
      @close="confirmDelete.show = false"
      @confirm="executeDelete"
    >
      <p>¿Estás seguro de que deseas eliminar la dirección en <strong>{{ confirmDelete.address?.street }}</strong>?</p>
    </BaseDialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { MapPin, Plus, Pencil, Trash2 } from 'lucide-vue-next'
import BaseDialog from '@/components/shared/BaseDialog.vue'
import { partyApi } from '@/services/partyApi'
import { useToastStore } from '@/stores/toast'

const props = defineProps({
  partyId: { type: String, required: true }
})

const toastStore = useToastStore()
const addresses = ref([])
const showModal = ref(false)
const editingId = ref(null)
const isSaving = ref(false)

const formData = reactive({
  street: '',
  city: '',
  postal_code: '',
  province: '',
  country: 'España',
  is_primary: false
})

// --- Confirm Dialog Logic ---
const confirmDelete = reactive({
  show: false,
  address: null
})

function promptDelete(address) {
  confirmDelete.address = address
  confirmDelete.show = true
}

async function executeDelete() {
  if (!confirmDelete.address) return
  try {
    await partyApi.deleteAddress(props.partyId, confirmDelete.address.id)
    toastStore.success('Dirección eliminada correctamente')
    await loadAddresses()
    confirmDelete.show = false
  } catch (err) {
    toastStore.error('No se pudo eliminar la dirección')
  }
}

async function loadAddresses() {
  if (!props.partyId) return
  try {
    addresses.value = await partyApi.listAddresses(props.partyId)
  } catch (err) {
    console.error('Error al cargar direcciones:', err)
  }
}

function openCreateModal() {
  editingId.value = null
  Object.assign(formData, { 
    street: '', 
    city: '', 
    postal_code: '', 
    province: '', 
    country: 'España', 
    is_primary: addresses.value.length === 0 
  })
  showModal.value = true
}

function editAddress(addr) {
  editingId.value = addr.id
  Object.assign(formData, { ...addr })
  showModal.value = true
}

async function saveAddress() {
  if (!formData.street || !formData.city) {
    toastStore.warning('La calle y la ciudad son obligatorias')
    return
  }
  isSaving.value = true
  try {
    const payload = {
      id: editingId.value || generateUUID(),
      street: formData.street,
      city: formData.city,
      province: formData.province,
      postalCode: formData.postal_code,
      country: formData.country,
      is_primary: formData.is_primary
    }

    if (editingId.value) {
      await partyApi.updateAddress(props.partyId, editingId.value, payload)
      toastStore.success('Dirección actualizada correctamente')
    } else {
      await partyApi.createAddress(props.partyId, payload)
      toastStore.success('Dirección añadida correctamente')
    }
    showModal.value = false
    await loadAddresses()
  } catch (err) {
    toastStore.error('Error al guardar la dirección')
  } finally {
    isSaving.value = false
  }
}

function generateUUID() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0
    return (c === 'x' ? r : (r & 0x3) | 0x8).toString(16)
  })
}

onMounted(() => loadAddresses())
</script>

<style scoped>
.address-manager { display: flex; flex-direction: column; gap: 1rem; }
.manager-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem; }
.manager-header h3 { font-size: 0.95rem; font-weight: 800; text-transform: uppercase; color: var(--color-text-secondary); margin: 0; }

.address-list { display: flex; flex-direction: column; gap: 0.5rem; }
.address-row { 
  background: white; border: 1px solid var(--color-border); border-radius: 8px; padding: 0.75rem 1.25rem; 
  display: flex; justify-content: space-between; align-items: center; transition: 0.2s; 
}
.address-row:hover { border-color: var(--color-primary); box-shadow: var(--box-shadow-sm); }
.address-row.is-primary { border-left: 4px solid var(--color-success); background: #f0fdf4; }

.address-main-info { display: flex; align-items: center; gap: 1rem; flex: 1; }
.row-indicator { display: flex; align-items: center; justify-content: center; }

.address-text { display: flex; flex-wrap: wrap; align-items: center; gap: 0.75rem; font-size: 0.95rem; }
.address-text .street { font-weight: 700; color: var(--color-text-primary); min-width: 250px; }
.address-text .location { color: var(--color-text-primary); }
.address-text .country { font-size: 0.8rem; font-weight: 600; text-transform: uppercase; }

.address-actions { display: flex; gap: 0.25rem; }
.btn-icon { background: transparent; border: none; cursor: pointer; color: var(--color-text-secondary); padding: 0.4rem; border-radius: 6px; }
.btn-icon:hover { background: rgba(0,0,0,0.05); color: var(--color-text-primary); }

.form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
.full-width { grid-column: 1 / -1; }
.form-group label { display: block; font-size: 0.75rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); margin-bottom: 0.4rem; }
.form-input { width: 100%; padding: 0.75rem; border: 1px solid var(--color-border); border-radius: 8px; font-family: inherit; }

.checkbox-label { display: flex; align-items: center; gap: 0.75rem; font-size: 0.85rem; cursor: pointer; }
.empty-state { padding: 2rem; text-align: center; background: var(--color-background); border-radius: 12px; border: 2px dashed var(--color-border); }
</style>
