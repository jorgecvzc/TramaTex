<template>
  <div class="person-manager">
    <div class="manager-header">
      <h3 class="flex items-center gap-2"><Contact2 :size="20" /> Personas de Contacto</h3>
      <button class="btn btn-primary btn-sm" @click="openCreateModal">
        <Plus :size="16" /> Añadir Persona
      </button>
    </div>

    <div v-if="persons.length === 0" class="empty-state">
      <p class="text-muted italic">No hay personas de contacto registradas para esta organización.</p>
    </div>

    <div v-else class="person-grid">
      <div v-for="person in persons" :key="person.id" class="person-card">
        <div class="person-main">
          <div class="avatar">
            <User :size="24" />
          </div>
          <div class="person-info">
            <strong>{{ person.first_name }} {{ person.last_name }}</strong>
            <span v-if="person.job_title" class="job-title">{{ person.job_title }}</span>
          </div>
        </div>
        
        <div class="person-details">
          <div v-if="person.email" class="detail-row">
            <Mail :size="14" /> <span>{{ person.email }}</span>
          </div>
          <div v-if="person.phone" class="detail-row">
            <Phone :size="14" /> <span>{{ person.phone }}</span>
          </div>
        </div>

        <div class="person-actions">
          <button class="btn-icon" @click="editPerson(person)" title="Editar"><Pencil :size="16" /></button>
          <button class="btn-icon text-danger" @click="promptDelete(person)" title="Eliminar"><Trash2 :size="16" /></button>
        </div>
      </div>
    </div>

    <!-- MODAL: CREAR/EDITAR PERSONA -->
    <BaseDialog
      :show="showModal"
      :title="editingId ? 'Editar Persona de Contacto' : 'Nueva Persona de Contacto'"
      icon="contact_2"
      confirm-text="Guardar Cambios"
      :is-confirming="isSaving"
      @close="showModal = false"
      @confirm="savePerson"
    >
      <div class="form-grid">
        <div class="form-group">
          <label>Nombre *</label>
          <input v-model="formData.first_name" type="text" class="form-input" required />
        </div>
        <div class="form-group">
          <label>Apellidos *</label>
          <input v-model="formData.last_name" type="text" class="form-input" required />
        </div>
        <div class="form-group">
          <label>Email</label>
          <input v-model="formData.email" type="email" class="form-input" placeholder="ejemplo@correo.com" />
        </div>
        <div class="form-group">
          <label>Teléfono</label>
          <input v-model="formData.phone" type="tel" class="form-input" />
        </div>
        <div class="form-group full-width">
          <label>Cargo / Departamento</label>
          <input v-model="formData.job_title" type="text" class="form-input" placeholder="Ej: Responsable de compras" />
        </div>
      </div>
    </BaseDialog>

    <!-- MODAL DE CONFIRMACIÓN -->
    <BaseDialog
      :show="confirmDelete.show"
      title="Eliminar Persona"
      icon="warning"
      confirm-text="Eliminar"
      confirm-class="btn-danger"
      @close="confirmDelete.show = false"
      @confirm="executeDelete"
    >
      <p>¿Estás seguro de que deseas eliminar a <strong>{{ confirmDelete.person?.first_name }} {{ confirmDelete.person?.last_name }}</strong> como contacto?</p>
    </BaseDialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Contact2, User, Plus, Pencil, Trash2, Mail, Phone } from 'lucide-vue-next'
import BaseDialog from '@/components/shared/BaseDialog.vue'
import { partyApi } from '@/services/partyApi'
import { useToastStore } from '@/stores/toast'

const props = defineProps({
  partyId: { type: String, required: true }
})

const toastStore = useToastStore()
const persons = ref([])
const showModal = ref(false)
const editingId = ref(null)
const isSaving = ref(false)

const formData = reactive({
  first_name: '',
  last_name: '',
  email: '',
  phone: '',
  job_title: ''
})

// --- Confirm Dialog Logic ---
const confirmDelete = reactive({
  show: false,
  person: null
})

function promptDelete(person) {
  confirmDelete.person = person
  confirmDelete.show = true
}

async function executeDelete() {
  if (!confirmDelete.person) return
  try {
    await partyApi.deletePerson(props.partyId, confirmDelete.person.id)
    toastStore.success('Persona de contacto eliminada')
    await loadPersons()
    confirmDelete.show = false
  } catch (err) {
    toastStore.error(err.message)
  }
}

async function loadPersons() {
  if (!props.partyId) return
  try {
    persons.value = await partyApi.listPersons(props.partyId)
  } catch (err) {
    console.error('Error loading persons:', err)
  }
}

function openCreateModal() {
  editingId.value = null
  Object.assign(formData, { first_name: '', last_name: '', email: '', phone: '', job_title: '' })
  showModal.value = true
}

function editPerson(person) {
  editingId.value = person.id
  Object.assign(formData, { ...person })
  showModal.value = true
}

async function savePerson() {
  if (!formData.first_name || !formData.last_name) return
  isSaving.value = true
  try {
    if (editingId.value) {
      await partyApi.updatePerson(props.partyId, editingId.value, formData)
      toastStore.success('Datos actualizados')
    } else {
      await partyApi.createPerson(props.partyId, formData)
      toastStore.success('Contacto añadido')
    }
    showModal.value = false
    await loadPersons()
  } catch (err) {
    toastStore.error(err.message)
  } finally {
    isSaving.value = false
  }
}

onMounted(() => loadPersons())
</script>

<style scoped>
.person-manager { display: flex; flex-direction: column; gap: 1.5rem; }
.manager-header { display: flex; justify-content: space-between; align-items: center; }
.manager-header h3 { font-size: 0.95rem; font-weight: 800; text-transform: uppercase; color: var(--color-text-secondary); margin: 0; }

.person-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 1rem; }
.person-card { background: white; border: 1px solid var(--color-border); border-radius: 12px; padding: 1.25rem; transition: 0.2s; position: relative; }
.person-card:hover { border-color: var(--color-primary); box-shadow: var(--box-shadow-sm); }

.person-main { display: flex; align-items: center; gap: 1rem; margin-bottom: 1rem; }
.avatar { width: 44px; height: 44px; background: var(--color-background); border-radius: 50%; display: flex; align-items: center; justify-content: center; color: var(--color-secondary); border: 1px solid var(--color-border); }

.person-info { display: flex; flex-direction: column; }
.person-info strong { font-size: 1rem; color: var(--color-text-primary); }
.job-title { font-size: 0.75rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); }

.person-details { border-top: 1px solid var(--color-background); pt: 0.75rem; display: flex; flex-direction: column; gap: 0.4rem; }
.detail-row { display: flex; align-items: center; gap: 0.6rem; font-size: 0.85rem; color: var(--color-text-secondary); }

.person-actions { position: absolute; top: 1rem; right: 1rem; display: flex; gap: 0.25rem; }
.btn-icon { background: transparent; border: none; cursor: pointer; color: var(--color-text-secondary); padding: 0.4rem; border-radius: 6px; }
.btn-icon:hover { background: rgba(0,0,0,0.05); color: var(--color-text-primary); }

.form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
.full-width { grid-column: 1 / -1; }
.form-group label { display: block; font-size: 0.75rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); margin-bottom: 0.4rem; }
.form-input { width: 100%; padding: 0.75rem; border: 1px solid var(--color-border); border-radius: 8px; font-family: inherit; }

.empty-state { padding: 2rem; text-align: center; background: var(--color-background); border-radius: 12px; border: 2px dashed var(--color-border); }
</style>
