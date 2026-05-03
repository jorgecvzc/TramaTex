<template>
  <div class="person-manager">
    <div class="manager-header">
      <h3 class="flex items-center gap-2"><Contact2 :size="20" /> Personas de Contacto</h3>
      <div class="header-actions">
        <button class="btn btn-outline btn-sm" @click="openLinkModal">
          <Search :size="16" /> Vincular Existente
        </button>
        <button class="btn btn-primary btn-sm ml-2" @click="openCreateModal">
          <Plus :size="16" /> Añadir Nueva
        </button>
      </div>
    </div>

    <div v-if="persons.length === 0" class="empty-state">
      <p class="text-muted italic">No hay personas de contacto registradas para esta organización.</p>
    </div>

    <div v-else class="person-list">
      <div v-for="person in persons" :key="person.id" class="person-row">
        <div class="person-main-info">
          <div class="avatar">
            <User :size="20" />
          </div>
          <div class="name-box">
            <strong>{{ person.first_name }} {{ person.last_name }}</strong>
            <span v-if="person.job_title" class="job-title">{{ person.job_title }}</span>
          </div>
        </div>
        
        <div class="person-contact-info">
          <div v-if="person.email" class="contact-item">
            <Mail :size="14" /> <span>{{ person.email }}</span>
          </div>
          <div v-if="person.phone" class="contact-item">
            <Phone :size="14" /> <span>{{ person.phone }}</span>
          </div>
        </div>

        <div class="person-actions">
          <button class="btn-icon" @click="editPerson(person)" title="Editar"><Pencil :size="18" /></button>
          <button class="btn-icon text-danger" @click="promptDelete(person)" title="Desvincular"><Trash2 :size="18" /></button>
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

    <!-- MODAL: VINCULAR EXISTENTE -->
    <BaseDialog
      :show="showLinkModal"
      title="Vincular Contacto Existente"
      icon="search"
      confirm-text="Vincular Seleccionado"
      :is-confirming="isSaving"
      @close="showLinkModal = false"
      @confirm="linkExistingPerson"
    >
      <div class="p-2">
        <p class="mb-4 text-sm text-muted">Busca una persona ya registrada en el sistema para vincularla a esta organización.</p>
        <PartySelector
          v-model="selectedExistingId"
          label="Persona de Contacto"
          placeholder="Escribe nombre para buscar..."
          role-filter="CONTACT"
          required
          @select="onExistingPersonSelect"
        />

        <div v-if="selectedExistingId" class="mt-6 p-4 border rounded bg-light animate-fade-in">
          <p class="section-label mb-3">Datos de contacto para esta empresa:</p>
          <div class="form-grid">
            <div class="form-group full-width">
              <label>Cargo / Departamento</label>
              <input v-model="linkFormData.job_title" type="text" class="form-input" placeholder="Ej: Jefe de Compras" />
            </div>
            <div class="form-group">
              <label>Email específico</label>
              <input v-model="linkFormData.email" type="email" class="form-input" />
            </div>
            <div class="form-group">
              <label>Teléfono directo</label>
              <input v-model="linkFormData.phone" type="tel" class="form-input" />
            </div>
          </div>
        </div>
      </div>
    </BaseDialog>

    <!-- MODAL DE CONFIRMACIÓN -->
    <BaseDialog
      :show="confirmDelete.show"
      title="Desvincular Persona"
      icon="warning"
      confirm-text="Desvincular"
      confirm-class="btn-danger"
      @close="confirmDelete.show = false"
      @confirm="executeDelete"
    >
      <p>¿Estás seguro de que deseas desvincular a <strong>{{ confirmDelete.person?.first_name }} {{ confirmDelete.person?.last_name }}</strong>?</p>
      <p class="mt-2 text-xs text-muted">Si esta persona no está vinculada a ninguna otra empresa, sus datos permanecerán en el sistema como contacto independiente.</p>
    </BaseDialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Contact2, User, Plus, Pencil, Trash2, Mail, Phone, Search } from 'lucide-vue-next'
import BaseDialog from '@/components/shared/BaseDialog.vue'
import PartySelector from '@/components/party/PartySelector.vue'
import { partyApi } from '@/services/partyApi'
import { useToastStore } from '@/stores/toast'

const props = defineProps({
  partyId: { type: String, required: true }
})

const toastStore = useToastStore()
const persons = ref([])
const showModal = ref(false)
const showLinkModal = ref(false)
const editingId = ref(null)
const isSaving = ref(false)
const selectedExistingId = ref('')
const selectedExistingContact = ref(null)

const formData = reactive({
  first_name: '',
  last_name: '',
  email: '',
  phone: '',
  job_title: ''
})

const linkFormData = reactive({
  job_title: 'Contacto Comercial',
  email: '',
  phone: ''
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
    // Para desvincular usamos el ID de la persona
    await partyApi.removeContact(props.partyId, confirmDelete.person.id, false)
    toastStore.success('Vínculo eliminado correctamente')
    await loadPersons()
    confirmDelete.show = false
  } catch (err) {
    toastStore.error('No se pudo desvincular el contacto')
  }
}

async function loadPersons() {
  if (!props.partyId) return
  try {
    const res = await partyApi.listPersons(props.partyId)
    persons.value = res || []
  } catch (err) {
    console.error('Error al cargar contactos:', err)
  }
}

function openCreateModal() {
  editingId.value = null
  Object.assign(formData, { first_name: '', last_name: '', email: '', phone: '', job_title: '' })
  showModal.value = true
}

function openLinkModal() {
  selectedExistingId.value = ''
  selectedExistingContact.value = null
  Object.assign(linkFormData, { job_title: 'Contacto Comercial', email: '', phone: '' })
  showLinkModal.value = true
}

function onExistingPersonSelect(party) {
  selectedExistingContact.value = party
  if (party) {
    linkFormData.email = party.email || ''
    linkFormData.phone = party.phone || ''
  }
}

const editingContactDetailsId = ref(null)

function editPerson(person) {
  editingId.value = person.id // Party ID de la persona
  editingContactDetailsId.value = person.contact_details_id // ID del vínculo
  Object.assign(formData, { ...person })
  showModal.value = true
}

async function savePerson() {
  if (!formData.first_name || !formData.last_name) {
    toastStore.warning('El nombre y los apellidos son obligatorios')
    return
  }
  isSaving.value = true
  try {
    if (editingId.value) {
      // 1. Actualizar los detalles del vínculo (email, tel, cargo)
      await partyApi.updatePerson(props.partyId, editingContactDetailsId.value, {
        email: formData.email,
        phone: formData.phone,
        job_title: formData.job_title
      })
      
      // 2. Opcional: Actualizar el nombre en la entidad persona (si ha cambiado)
      // Por ahora el backend de updatePerson solo toca los detalles del vínculo
      
      toastStore.success('Datos actualizados correctamente')
    } else {
      await partyApi.createPerson(props.partyId, {
        firstName: formData.first_name,
        lastName: formData.last_name,
        email: formData.email,
        phone: formData.phone,
        jobTitle: formData.job_title
      })
      toastStore.success('Contacto creado y vinculado')
    }
    showModal.value = false
    await loadPersons()
  } catch (err) {
    toastStore.error('Error al guardar el contacto')
  } finally {
    isSaving.value = false
  }
}

async function linkExistingPerson() {
  if (!selectedExistingId.value) {
    toastStore.warning('Selecciona una persona de la lista')
    return
  }
  isSaving.value = true
  try {
    await partyApi.linkExistingContact(props.partyId, selectedExistingId.value, {
      jobTitle: linkFormData.job_title,
      email: linkFormData.email,
      phone: linkFormData.phone
    })
    toastStore.success('Contacto vinculado correctamente')
    showLinkModal.value = false
    await loadPersons()
  } catch (err) {
    toastStore.error('Error al vincular el contacto')
  } finally {
    isSaving.value = false
  }
}

onMounted(() => loadPersons())
</script>

<style scoped>
.person-manager { display: flex; flex-direction: column; gap: 1rem; }
.manager-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem; }
.manager-header h3 { font-size: 0.95rem; font-weight: 800; text-transform: uppercase; color: var(--color-text-secondary); margin: 0; }

.person-list { display: flex; flex-direction: column; gap: 0.5rem; }
.person-row { 
  background: white; border: 1px solid var(--color-border); border-radius: 8px; padding: 0.75rem 1.25rem; 
  display: flex; justify-content: space-between; align-items: center; transition: 0.2s; 
}
.person-row:hover { border-color: var(--color-primary); box-shadow: var(--box-shadow-sm); }

.person-main-info { display: flex; align-items: center; gap: 1rem; flex: 1; }
.avatar { width: 36px; height: 36px; background: var(--color-background); border-radius: 50%; display: flex; align-items: center; justify-content: center; color: var(--color-secondary); border: 1px solid var(--color-border); flex-shrink: 0; }

.name-box { display: flex; flex-direction: column; line-height: 1.2; min-width: 200px; }
.name-box strong { font-size: 0.95rem; color: var(--color-text-primary); }
.job-title { font-size: 0.65rem; font-weight: 800; text-transform: uppercase; color: var(--color-text-secondary); letter-spacing: 0.025em; }

.person-contact-info { display: flex; flex: 2; gap: 2rem; align-items: center; }
.contact-item { display: flex; align-items: center; gap: 0.5rem; font-size: 0.85rem; color: var(--color-text-primary); font-weight: 500; }
.contact-item :deep(svg) { color: var(--color-text-secondary); }

.person-actions { display: flex; gap: 0.25rem; }
.btn-icon { background: transparent; border: none; cursor: pointer; color: var(--color-text-secondary); padding: 0.4rem; border-radius: 6px; }
.btn-icon:hover { background: rgba(0,0,0,0.05); color: var(--color-text-primary); }

.form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
.full-width { grid-column: 1 / -1; }
.form-group label { display: block; font-size: 0.75rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); margin-bottom: 0.4rem; }
.form-input { width: 100%; padding: 0.75rem; border: 1px solid var(--color-border); border-radius: 8px; font-family: inherit; }

.empty-state { padding: 2rem; text-align: center; background: var(--color-background); border-radius: 12px; border: 2px dashed var(--color-border); }

@media (max-width: 900px) {
  .person-row { flex-direction: column; align-items: flex-start; gap: 1rem; }
  .person-contact-info { flex-direction: column; align-items: flex-start; gap: 0.5rem; }
  .person-actions { position: absolute; top: 0.75rem; right: 0.75rem; }
}
</style>
