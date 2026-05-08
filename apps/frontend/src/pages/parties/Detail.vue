<template>
  <!-- PAGE: DETAIL / EDIT / ENTITY CREATION -->
  
  <BaseEntityPage v-if="isLoading" class="no-print">
    <template #header>
      <BasePageHeader title="Cargando..." :breadcrumbs="[{ label: 'Entidades', to: '/parties' }, { label: 'Detalle' }]" />
    </template>
    <div class="state-container loading-state">
      <RefreshCw :size="48" class="spin mb-4 opacity-50" />
      <p>Sincronizando ficha de entidad...</p>
    </div>
  </BaseEntityPage>

  <BaseEntityPage v-else-if="error" class="no-print">
    <template #header>
      <BasePageHeader title="Error" :breadcrumbs="[{ label: 'Entidades', to: '/parties' }, { label: 'Listado' }]" />
    </template>
    <div class="state-container error-state">
      <AlertCircle :size="56" class="mb-4 text-danger opacity-80" />
      <h3>Error de Carga</h3>
      <p>{{ error }}</p>
      <button class="btn btn-outline btn-sm mt-6" @click="router.push('/parties')">Volver al listado</button>
    </div>
  </BaseEntityPage>

  <BaseEntityPage v-else class="no-print">
    <!-- LAYER 1: IDENTITY -->
    <template #header>
      <div class="sticky-header-container">
        <BasePageHeader 
          :title="mode === 'create' ? 'Nueva Entidad' : (mode === 'edit' ? `Editando ${party?.name}` : party?.name)" 
          :breadcrumbs="[{ label: 'Entidades', to: '/parties' }, { label: mode === 'create' ? 'Alta' : party?.name }]"
          show-back
        >
          <template #icon>
            <component :is="getIdentityIcon()" :size="28" />
          </template>
          <template #actions>
            <div v-if="mode === 'detail'" class="header-actions-group">
              <button class="btn btn-primary" @click="enterEditMode">
                <Pencil :size="18" />
                <span>Editar Datos</span>
              </button>
              <button class="btn btn-danger ml-2" @click="promptDelete">
                <Trash2 :size="18" />
                <span>Eliminar Entidad</span>
              </button>
            </div>
          </template>
        </BasePageHeader>

        <!-- TAB NAVIGATION (Only in detail) -->
        <nav v-if="mode === 'detail'" class="entity-tabs">
          <button 
            v-for="tab in tabs" 
            :key="tab.id" 
            @click="activeTab = tab.id"
            :class="['tab-btn', { active: activeTab === tab.id }]"
          >
            <component :is="tab.icon" :size="18" />
            <span>{{ tab.label }}</span>
          </button>
        </nav>
      </div>
    </template>

    <!-- LAYER 2: CONTEXT (Summary) -->
    <template #summary v-if="mode !== 'create' && party">
      <div class="overview-details-strip">
        <div class="detail-item">
          <div class="icon blue"><Fingerprint :size="20" /></div>
          <div class="text-box">
            <label>Identificación Fiscal</label>
            <strong>{{ party.tax_id || '—' }}</strong>
          </div>
        </div>
        <div class="detail-item">
          <div class="icon yellow"><MapPin :size="20" /></div>
          <div class="text-box">
            <label>Ubicación Principal</label>
            <strong>{{ primaryAddressLine }}</strong>
          </div>
        </div>
        <div class="detail-item" v-if="party.email || party.phone">
          <div class="icon purple"><ContactRound :size="20" /></div>
          <div class="text-box">
            <label>Contacto Directo</label>
            <div class="inline-contacts">
              <span v-if="party.email">{{ party.email }}</span>
              <span v-if="party.email && party.phone" class="separator">|</span>
              <span v-if="party.phone">{{ party.phone }}</span>
            </div>
          </div>
        </div>
        <div class="detail-item">
          <div class="icon green"><History :size="20" /></div>
          <div class="text-box">
            <label>Fecha de Alta</label>
            <strong>{{ formatDate(party.created_at) }}</strong>
          </div>
        </div>
      </div>
    </template>

    <!-- LAYER 3: WORK -->
    <div class="party-master-content">
      <!-- TAB: GENERAL / FORM -->
      <div v-if="activeTab === 'general'" class="tab-fade-in">
        <template v-if="mode === 'detail'">
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <!-- SECTION: IDENTITY -->
            <FormSection title="Datos de Identidad" icon="person">
              <DataRow label="Nombre Comercial / Completo" :value="party?.name" icon="badge" />
              <DataRow label="Tipo de Entidad" :value="getEntityTypeLabel(party?.type)" icon="users" />
              <DataRow label="Identificador Fiscal" :value="`${party?.tax_id_type || 'CIF'}: ${party?.tax_id || '—'}`" is-mono icon="fingerprint" />
              <DataRow label="Notas Internas" icon="notes">
                <p class="notes-text">{{ party?.notes || 'Sin observaciones registradas.' }}</p>
              </DataRow>
            </FormSection>

            <!-- SECTION: CONFIGURATION -->
            <FormSection title="Configuración de Cuenta" icon="settings">
              <DataRow label="Estado Actual" :value="getStatusLabel(party?.status)" icon="shield-check" />
              <DataRow label="Rol en el Sistema" :value="getRoleLabel(party?.role)" icon="git-fork" />
              <DataRow v-if="party?.role !== 'SUPPLIER'" label="Descuento por Defecto" :value="(party?.default_discount_percentage || 0) + '%'" icon="tag" />
            </FormSection>
          </div>
        </template>
        <template v-else>
          <PartyForm 
            :party-id="mode === 'edit' ? party?.id : undefined" 
            :initial-data="mode === 'edit' ? party : undefined" 
            @submit="(p) => router.push(`/parties/${p.id}`)"
            @update="loadData"
            @cancel="exitEditMode"
          />
        </template>
      </div>

      <!-- SECONDARY TABS (Detail mode only) -->
      <div v-if="activeTab === 'addresses' && mode === 'detail'" class="tab-fade-in">
        <AddressManager :party-id="party?.id" />
      </div>

      <div v-if="activeTab === 'contacts' && mode === 'detail'" class="tab-fade-in">
        <PersonManager :party-id="party?.id" />
      </div>
    </div>

    <!-- DELETE DIALOG -->
    <BaseDialog
      :show="confirmDelete.show"
      title="Eliminar Entidad"
      icon="warning"
      confirm-text="Eliminar Definitivamente"
      confirm-class="btn-danger"
      :is-confirming="isSaving"
      @close="confirmDelete.show = false"
      @confirm="executeDelete"
    >
      <p>¿Estás seguro de que deseas eliminar permanentemente a <strong>{{ party?.name }}</strong>?</p>
      <p class="mt-2 text-muted italic">Esta acción solo se completará si no existen documentos vinculados.</p>
    </BaseDialog>
  </BaseEntityPage>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { 
  Building2, User, Pencil, Trash2, Fingerprint, MapPin, 
  History, Info, Settings, ShieldCheck, GitFork, Tag, 
  ContactRound, X, Save, RefreshCw, Percent, AlertCircle
} from 'lucide-vue-next'
import BaseEntityPage from '@/components/shared/BaseEntityPage.vue'
import BasePageHeader from '@/components/shared/BasePageHeader.vue'
import FormSection from '@/components/shared/FormSection.vue'
import DataRow from '@/components/shared/DataRow.vue'
import BaseDialog from '@/components/shared/BaseDialog.vue'
import AddressManager from '@/components/party/AddressManager.vue'
import PersonManager from '@/components/party/PersonManager.vue'
import PartyForm from '@/components/party/PartyForm.vue'
import { partyApi } from '@/services/partyApi'
import { useToastStore } from '@/stores/toast'

const route = useRoute()
const router = useRouter()
const toastStore = useToastStore()

const mode = ref<'detail' | 'edit' | 'create'>('detail')
const activeTab = ref('general')
const isLoading = ref(true)
const isSaving = ref(false)
const error = ref('')

const party = ref<any>(null)
const addresses = ref<any[]>([])

const formData = reactive({
  name: '',
  type: 'ORGANIZATION',
  taxId: '',
  taxIdType: 'CIF',
  status: 'ACTIVE',
  role: 'CLIENT',
  notes: '',
  defaultDiscountPercentage: 0
})

const tabs = computed(() => {
  const baseTabs = [
    { id: 'general', label: 'General', icon: Info },
    { id: 'addresses', label: 'Direcciones', icon: MapPin }
  ]
  
  // Contacts tab only applies to organizations
  if (party.value?.type === 'ORGANIZATION') {
    baseTabs.push({ id: 'contacts', label: 'Contactos', icon: ContactRound })
  }
  
  return baseTabs
})

const confirmDelete = reactive({ show: false })

const primaryAddressLine = computed(() => {
  if (!addresses.value.length) return 'Sin dirección'
  const primary = addresses.value.find(a => a.is_primary) || addresses.value[0]
  return `${primary.city}, ${primary.province}`
})

function getIdentityIcon() {
  const type = mode.value === 'create' ? formData.type : party.value?.type
  return type === 'ORGANIZATION' ? Building2 : User
}

function getEntityTypeLabel(type) {
  return type === 'ORGANIZATION' ? 'Empresa / Jurídica' : 'Persona Física'
}

function getStatusLabel(status) {
  const map = { ACTIVE: 'ACTIVO', INACTIVE: 'INACTIVO' }
  return map[status] || status
}

function getRoleLabel(role) {
  const map = { CLIENT: 'Cliente', SUPPLIER: 'Proveedor', BOTH: 'Cliente y Proveedor' }
  return map[role] || role
}

async function loadData() {
  const id = route.params.id as string
  if (!id || id === 'new') {
    mode.value = 'create'
    resetForm()
    isLoading.value = false
    return
  }

  mode.value = 'detail'
  isLoading.value = true
  error.value = ''
  try {
    const [partyData, addrData] = await Promise.all([
      partyApi.getParty(id),
      partyApi.listAddresses(id)
    ])
    party.value = partyData
    addresses.value = addrData
  } catch (err: any) {
    error.value = 'No se pudo cargar la información de la entidad.'
    console.error(err)
  } finally {
    isLoading.value = false
  }
}

function resetForm() {
  Object.assign(formData, {
    name: '', type: 'ORGANIZATION', taxId: '', taxIdType: 'CIF', status: 'ACTIVE',
    role: 'CLIENT', notes: '', defaultDiscountPercentage: 0
  })
}

function enterEditMode() {
  Object.assign(formData, {
    name: party.value.name,
    type: party.value.type,
    taxId: party.value.tax_id,
    taxIdType: party.value.tax_id_type || (party.value.type === 'ORGANIZATION' ? 'CIF' : 'NIF'),
    status: party.value.status,
    role: party.value.role,
    notes: party.value.notes || '',
    defaultDiscountPercentage: party.value.default_discount_percentage || 0
  })
  mode.value = 'edit'
}

function exitEditMode() {
  if (mode.value === 'edit') mode.value = 'detail'
  else router.push('/parties')
}

async function saveParty() {
  if (!formData.name) {
    toastStore.warning('El nombre es obligatorio')
    return
  }

  isSaving.value = true
  try {
    const partyId = mode.value === 'create' ? generateUUID() : party.value.id
    const payload = {
      id: partyId,
      name: formData.name,
      type: formData.type,
      role: formData.role,
      taxId: formData.taxId,
      taxIdType: formData.taxIdType,
      status: formData.status,
      notes: formData.notes,
      default_discount_percentage: formData.role === 'SUPPLIER' ? 0 : (formData.defaultDiscountPercentage || 0),
      hasPerson: formData.type === 'PERSON'
    }

    if (mode.value === 'create') {
      const newParty = await partyApi.createParty({
        ...payload,
        entityType: formData.type // API naming alignment
      })
      toastStore.success('Entidad creada con éxito')
      router.push(`/parties/${newParty.id}`)
    } else {
      await partyApi.updateParty(party.value.id, payload)
      toastStore.success('Ficha técnica actualizada')
      await loadData()
      mode.value = 'detail'
    }
  } catch (err: any) {
    toastStore.error('Error al guardar: ' + (err.message || 'Error desconocido'))
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

function promptDelete() { confirmDelete.show = true }

async function executeDelete() {
  isSaving.value = true
  try {
    await partyApi.deleteParty(party.value.id)
    toastStore.success('Entidad eliminada')
    router.push('/parties')
  } catch (err: any) {
    toastStore.error('Error al eliminar: ' + err.message)
  } finally {
    isSaving.value = false
    confirmDelete.show = false
  }
}

function formatDate(d: string) { return d ? new Date(d).toLocaleDateString('es-ES') : '—' }

watch(() => route.params.id, () => loadData(), { immediate: true })
</script>

<style scoped>
@import "@/design-system/_sections.css";

.sticky-header-container { background: white; margin-top: -1.5rem; padding-top: 1.5rem; border-bottom: 1px solid var(--color-border); box-shadow: var(--box-shadow-sm); }
.entity-tabs { display: flex; max-width: 1300px; margin: 0 auto; padding: 0 2rem; gap: 0.25rem; }
.tab-btn { display: flex; align-items: center; gap: 0.6rem; padding: 1rem 1.25rem; background: transparent; border: none; border-bottom: 3px solid transparent; color: var(--color-text-secondary); font-weight: 700; cursor: pointer; transition: all 0.2s; font-size: 0.85rem; text-transform: uppercase; letter-spacing: 0.025em; margin-bottom: -1px; }
.tab-btn:hover { color: var(--color-primary); background: rgba(0,0,0,0.02); }
.tab-btn.active { border-bottom-color: var(--color-secondary); color: var(--color-secondary); background: rgba(0, 35, 149, 0.03); }

.overview-details-strip { 
  display: flex; flex-wrap: wrap; gap: 2rem; background: white; border: 1px solid var(--color-border); 
  border-radius: 12px; padding: 1.25rem 2rem; box-shadow: var(--box-shadow-sm); 
}
.detail-item { display: flex; align-items: center; gap: 1rem; }
.detail-item .icon { 
  width: 36px; height: 36px; border-radius: 8px; display: flex; 
  align-items: center; justify-content: center; flex-shrink: 0; 
}
.detail-item .icon.blue { background: rgba(59, 130, 246, 0.1); color: #2563eb; }
.detail-item .icon.yellow { background: rgba(230, 184, 0, 0.1); color: #d97706; }
.detail-item .icon.purple { background: rgba(168, 85, 247, 0.1); color: #a855f7; }
.detail-item .icon.green { background: rgba(34, 197, 94, 0.1); color: #16a34a; }

.text-box { display: flex; flex-direction: column; line-height: 1.2; }
.text-box label { font-size: 0.65rem; font-weight: 800; text-transform: uppercase; color: var(--color-text-secondary); letter-spacing: 0.05em; }
.text-box strong { font-size: 0.95rem; color: var(--color-text-primary); font-weight: 700; }

.inline-contacts { display: flex; align-items: center; gap: 0.75rem; font-size: 0.9rem; font-weight: 600; color: var(--color-text-primary); }
.inline-contacts .separator { color: var(--color-border-strong); font-weight: 300; }

.party-master-content { padding-top: 1rem; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
.form-group label { display: block; font-size: 0.75rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); margin-bottom: 0.5rem; }
.form-input, .form-textarea { width: 100%; padding: 0.75rem 1rem; border-radius: 8px; border: 1px solid var(--color-border); font-family: inherit; }
.form-input:focus, .form-textarea:focus { outline: none; border-color: var(--color-primary); box-shadow: 0 0 0 3px rgba(230, 184, 0, 0.1); }
.form-input:disabled { background: var(--color-background); cursor: not-allowed; }

.input-with-icon { position: relative; display: flex; align-items: center; }
.icon-start { position: absolute; left: 0.75rem; color: var(--color-text-secondary); }
.input-with-icon input { padding-left: 2.5rem; }

.action-toolbar { display: flex; align-items: center; padding: 0.75rem 1.5rem; background: white; border: 1px solid var(--color-border); border-radius: 8px; margin-bottom: 1.5rem; }
.role-pill { font-size: 0.65rem; font-weight: 800; text-transform: uppercase; padding: 0.2rem 0.6rem; border-radius: 4px; margin-left: 0.5rem; }
.role-pill.client { background: rgba(34, 197, 94, 0.1); color: #16a34a; }
.role-pill.supplier { background: rgba(147, 51, 234, 0.1); color: #9333ea; }

.notes-text { font-style: italic; color: var(--color-text-secondary); line-height: 1.5; margin: 0; }
.tab-fade-in { animation: fadeIn 0.3s ease-in-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(5px); } to { opacity: 1; transform: translateY(0); } }

.state-container { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 8rem 2rem; text-align: center; }
.spin { animation: spin 1s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
