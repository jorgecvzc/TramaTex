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

    <!-- MODAL: CREATE/EDIT PERSON -->
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

    <!-- MODAL: LINK EXISTING -->
    <BaseDialog
      :show="showLinkModal"
      title="Vincular Contacto Existente"
      icon="search"
      confirm-text="Vincular Seleccionado"
      :is-confirming="isSaving"
      @close="showLinkModal = false"
      @confirm="linkExistingPerson"
    >
      <div class="p-1">
        <p class="mb-4 text-sm text-muted">Selecciona una persona física ya registrada en el sistema para vincularla a esta organización.</p>
        
        <!-- Search Box -->
        <div class="search-box mb-3">
          <div class="input-with-icon">
            <Search :size="18" class="icon-start" />
            <input 
              v-model="linkSearchTerm" 
              type="text" 
              class="form-input" 
              placeholder="Filtrar por nombre o identificación..." 
            />
          </div>
        </div>

        <!-- Scrollable List -->
        <div class="selectable-list-container">
          <div v-if="isLoadingCandidates" class="loading-mini">
            <RefreshCw :size="20" class="spin" />
            <span>Cargando personas...</span>
          </div>
          <div v-else-if="filteredCandidates.length === 0" class="empty-mini">
            <p>No se encontraron personas físicas disponibles.</p>
          </div>
          <div v-else class="selectable-list scrollable">
            <div 
              v-for="candidate in filteredCandidates" 
              :key="candidate.id"
              :class="['list-item', { selected: selectedExistingId === candidate.id }]"
              @click="selectCandidate(candidate)"
            >
              <div class="candidate-info">
                <span class="candidate-name">{{ candidate.name }}</span>
                <span v-if="candidate.tax_id" class="candidate-tax">{{ candidate.tax_id }}</span>
              </div>
              <div v-if="selectedExistingId === candidate.id" class="selection-mark">
                <Check :size="18" />
              </div>
            </div>
          </div>
        </div>

        <!-- Specific Link Data -->
        <div v-if="selectedExistingId" class="mt-4 p-4 border rounded bg-light animate-fade-in">
          <p class="section-label mb-3">Datos específicos para esta organización:</p>
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

    <!-- CONFIRMATION MODAL -->
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
import { ref, reactive, onMounted, computed } from 'vue'
import { Contact2, User, Plus, Pencil, Trash2, Mail, Phone, Search, Check, RefreshCw } from 'lucide-vue-next'
import BaseDialog from '@/components/shared/BaseDialog.vue'
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

const linkSearchTerm = ref('')
const candidates = ref([])
const isLoadingCandidates = ref(false)

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

const filteredCandidates = computed(() => {
  if (!linkSearchTerm.value) return candidates.value
  const term = linkSearchTerm.value.toLowerCase()
  return candidates.value.filter(c => 
    c.name.toLowerCase().includes(term) || 
    (c.tax_id && c.tax_id.toLowerCase().includes(term))
  )
})

// --- Validation ---
function isValidEmail(email) {
  if (!email) return true
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
}

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
    console.error('Error loading contacts:', err)
  }
}

function openCreateModal() {
  editingId.value = null
  Object.assign(formData, { first_name: '', last_name: '', email: '', phone: '', job_title: '' })
  showModal.value = true
}

async function openLinkModal() {
  selectedExistingId.value = ''
  linkSearchTerm.value = ''
  Object.assign(linkFormData, { job_title: 'Contacto Comercial', email: '', phone: '' })
  showLinkModal.value = true
  
  isLoadingCandidates.value = true
  try {
    // Only physical persons
    const res = await partyApi.listParties({ type: 'person', pageSize: 500 })
    const linkedIds = persons.value.map(p => p.id)
    candidates.value = (res.data || []).filter(c => !linkedIds.includes(c.id))
  } catch (err) {
    toastStore.error('Error al cargar candidatos')
  } finally {
    isLoadingCandidates.value = false
  }
}

function selectCandidate(candidate) {
  selectedExistingId.value = candidate.id
  linkFormData.email = candidate.email || ''
  linkFormData.phone = candidate.phone || ''
}

const editingContactDetailsId = ref(null)

function editPerson(person) {
  editingId.value = person.id
  editingContactDetailsId.value = person.contact_details_id
  Object.assign(formData, { ...person })
  showModal.value = true
}

async function savePerson() {
  if (!formData.first_name || !formData.last_name) {
    toastStore.warning('El nombre y los apellidos son obligatorios')
    return
  }
  
  if (!isValidEmail(formData.email)) {
    toastStore.warning('El formato del email no es válido')
    return
  }

  isSaving.value = true
  try {
    if (editingId.value) {
      await partyApi.updatePerson(props.partyId, editingContactDetailsId.value, {
        email: formData.email,
        phone: formData.phone,
        job_title: formData.job_title
      })
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
    toastStore.error('Error al guardar el contacto: ' + (err.message || ''))
  } finally {
    isSaving.value = false
  }
}

async function linkExistingPerson() {
  if (!selectedExistingId.value) {
    toastStore.warning('Selecciona una persona de la lista')
    return
  }
  
  if (!isValidEmail(linkFormData.email)) {
    toastStore.warning('El formato del email específico no es válido')
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
    toastStore.error('Error al vincular el contacto: ' + (err.message || ''))
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
.form-input { width: 100%; padding: 0.75rem; border: 1px solid var(--color-border); border-radius: 8px; font-family: inherit; transition: all 0.2s; }
.form-input:focus { outline: none; border-color: var(--color-primary); box-shadow: 0 0 0 3px rgba(230, 184, 0, 0.1); }

.empty-state { padding: 2rem; text-align: center; background: var(--color-background); border-radius: 12px; border: 2px dashed var(--color-border); }

/* Selectable List Styles */
.selectable-list-container { border: 1px solid var(--color-border); border-radius: 8px; overflow: hidden; background: white; }
.selectable-list.scrollable { max-height: 250px; overflow-y: auto; }
.list-item { padding: 0.75rem 1rem; border-bottom: 1px solid var(--color-background); cursor: pointer; display: flex; justify-content: space-between; align-items: center; transition: 0.2s; }
.list-item:last-child { border-bottom: none; }
.list-item:hover { background-color: var(--color-background-soft); }
.list-item.selected { background-color: rgba(230, 184, 0, 0.1); border-left: 4px solid var(--color-primary); }

.candidate-info { display: flex; flex-direction: column; gap: 0.1rem; }
.candidate-name { font-weight: 600; color: var(--color-text-primary); font-size: 0.9rem; }
.candidate-tax { font-size: 0.7rem; color: var(--color-text-secondary); font-family: var(--font-family-mono); }
.selection-mark { color: var(--color-success); }

.loading-mini, .empty-mini { padding: 2rem; text-align: center; color: var(--color-text-secondary); font-size: 0.85rem; display: flex; flex-direction: column; align-items: center; gap: 0.5rem; }

.section-label { font-size: 0.7rem; font-weight: 800; text-transform: uppercase; color: var(--color-text-secondary); letter-spacing: 0.05em; border-bottom: 1px solid var(--color-border); padding-bottom: 0.5rem; }

.animate-fade-in { animation: fadeIn 0.3s ease-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(5px); } to { opacity: 1; transform: translateY(0); } }

@media (max-width: 900px) {
  .person-row { flex-direction: column; align-items: flex-start; gap: 1rem; }
  .person-contact-info { flex-direction: column; align-items: flex-start; gap: 0.5rem; }
  .person-actions { position: absolute; top: 0.75rem; right: 0.75rem; }
}
</style>
