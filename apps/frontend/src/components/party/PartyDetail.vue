<template>
  <BaseEntityPage :is-loading="isLoading" :error="error" @refresh="loadPartyData">
    <!-- CAPA 1: IDENTIDAD -->
    <template #header>
      <BasePageHeader 
        :title="party?.name || 'Cargando...'" 
        :breadcrumbs="[{ label: 'Entidades', to: '/parties' }, { label: party?.name || 'Detalle' }]"
        show-back
      >
        <template #icon>
          <component :is="party?.type === 'ORGANIZATION' ? Building2 : User" :size="28" />
        </template>
        <template #actions>
          <div v-if="party" class="header-actions-group">
            <template v-if="mode === 'detail'">
              <button class="btn btn-primary btn-sm" @click="enterEditMode">
                <Pencil :size="16" />
                <span>Editar Datos</span>
              </button>
              <button class="btn btn-danger btn-sm ml-2" @click="promptDelete">
                <Trash2 :size="16" />
                <span>Eliminar Entidad</span>
              </button>
            </template>
            <template v-else>
              <button class="btn btn-outline btn-sm" @click="exitEditMode" :disabled="isSaving">
                <X :size="16" />
                <span>Cancelar</span>
              </button>
              <button class="btn btn-secondary btn-sm ml-2" @click="saveParty" :disabled="isSaving">
                <component :is="isSaving ? RefreshCw : Save" :size="16" :class="{ 'spin': isSaving }" />
                <span>Guardar Cambios</span>
              </button>
            </template>
          </div>
        </template>
      </BasePageHeader>
    </template>

    <!-- CAPA 2: CONTEXTO -->
    <template #toolbar v-if="party">
      <div class="action-toolbar card">
        <div class="toolbar-info">
          <span :class="['status-badge', party.status === 'ACTIVE' ? 'status-success' : 'status-secondary']">
            {{ party.status === 'ACTIVE' ? 'ACTIVA' : 'INACTIVA' }}
          </span>
          <div class="roles-inline">
            <span v-if="party.role === 'CLIENT' || party.role === 'BOTH'" class="role-pill client">Cliente</span>
            <span v-if="party.role === 'SUPPLIER' || party.role === 'BOTH'" class="role-pill supplier">Proveedor</span>
          </div>
        </div>
      </div>
    </template>

    <template #summary v-if="party">
      <div class="overview-tags-row">
        <div class="summary-tag">
          <div class="icon blue"><Fingerprint :size="20" /></div>
          <div class="tag-content">
            <label>NIF / CIF</label>
            <strong class="text-mono">{{ party.tax_id || '—' }}</strong>
          </div>
        </div>
        <div class="summary-tag">
          <div class="icon yellow"><MapPin :size="20" /></div>
          <div class="tag-content">
            <label>Ubicación</label>
            <strong>{{ primaryAddressLine }}</strong>
          </div>
        </div>
        <div class="summary-tag">
          <div class="icon purple"><History :size="20" /></div>
          <div class="tag-content">
            <label>Alta en Sistema</label>
            <strong>{{ formatDate(party.created_at) }}</strong>
          </div>
        </div>
      </div>
    </template>

    <!-- CAPA 3: TRABAJO -->
    <div v-if="party" class="party-detail-tabs">
      <div class="tabs-header" v-if="mode === 'detail'">
        <button v-for="tab in tabs" :key="tab.id" :class="['tab-btn', { active: activeTab === tab.id }]" @click="activeTab = tab.id">
          <component :is="tab.icon" :size="18" />
          <span>{{ tab.label }}</span>
        </button>
      </div>

      <div class="tab-content mt-6">
        <!-- TAB: INFORMACIÓN GENERAL -->
        <div v-if="activeTab === 'general'" class="tab-pane animate-fade-in">
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <FormSection title="Datos de Identidad" icon="person">
              <div v-if="mode === 'detail'">
                <DataRow label="Nombre Comercial / Completo" :value="party.name" icon="badge" />
                <DataRow label="Tipo de Entidad" :value="party.type === 'ORGANIZATION' ? 'Empresa / Organización' : 'Persona Física'" icon="users" />
                <DataRow label="Identificador Fiscal" :value="party.tax_id" is-mono icon="fingerprint" />
                <DataRow label="Notas Internas" icon="notes">
                  <p class="notes-text">{{ party.notes || 'Sin observaciones registradas.' }}</p>
                </DataRow>
              </div>
              <div v-else class="form-container">
                <div class="form-group">
                  <label>Nombre de la entidad *</label>
                  <input v-model="formData.name" type="text" class="form-input" required placeholder="Nombre comercial o completo" />
                </div>
                <div class="form-row mt-4">
                  <div class="form-group">
                    <label>Tipo de Entidad</label>
                    <select v-model="formData.type" class="form-input" disabled>
                      <option value="ORGANIZATION">Empresa / Organización</option>
                      <option value="PERSON">Persona Física</option>
                    </select>
                    <span class="help-text">El tipo no se puede cambiar tras la creación</span>
                  </div>
                  <div class="form-group">
                    <label>Identificador Fiscal</label>
                    <input v-model="formData.taxId" type="text" class="form-input text-mono" placeholder="NIF / CIF" />
                  </div>
                </div>
                <div class="form-group mt-4">
                  <label>Notas Internas</label>
                  <textarea v-model="formData.notes" class="form-textarea" rows="3" placeholder="Observaciones privadas..."></textarea>
                </div>
              </div>
            </FormSection>

            <FormSection title="Configuración de Cuenta" icon="settings">
              <div v-if="mode === 'detail'">
                <DataRow label="Estado Actual" :value="party.status" icon="shield-check" />
                <DataRow label="Rol en el Sistema" :value="party.role" icon="git-fork" />
                <DataRow label="Descuento por Defecto" :value="(party.default_discount_percentage || 0) + '%'" icon="tag" />
              </div>
              <div v-else class="form-container">
                <div class="form-group">
                  <label>Estado en Sistema</label>
                  <select v-model="formData.status" class="form-input">
                    <option value="ACTIVE">ACTIVO</option>
                    <option value="INACTIVE">INACTIVO</option>
                  </select>
                </div>
                <div class="form-group mt-4">
                  <label>Rol Principal</label>
                  <select v-model="formData.role" class="form-input">
                    <option value="CLIENT">Cliente</option>
                    <option value="SUPPLIER">Proveedor</option>
                    <option value="BOTH">Cliente y Proveedor</option>
                  </select>
                </div>
                <div class="form-group mt-4">
                  <label>Bonificación Comercial (%)</label>
                  <div class="input-with-icon">
                    <Percent :size="18" class="icon-start" />
                    <input v-model.number="formData.defaultDiscountPercentage" type="number" step="0.01" min="0" max="100" class="form-input" />
                  </div>
                </div>
              </div>
            </FormSection>
          </div>
        </div>

        <!-- TAB: DIRECCIONES -->
        <div v-if="activeTab === 'addresses' && mode === 'detail'" class="tab-pane animate-fade-in">
          <AddressManager :party-id="party.id" />
        </div>

        <!-- TAB: CONTACTOS -->
        <div v-if="activeTab === 'contacts' && mode === 'detail'" class="tab-pane animate-fade-in">
          <PersonManager :party-id="party.id" />
        </div>
      </div>
    </div>

    <!-- MODAL DE CONFIRMACIÓN DE ELIMINACIÓN -->
    <BaseDialog
      :show="confirmDelete.show"
      title="Eliminar Entidad"
      icon="warning"
      confirm-text="Eliminar Definitivamente"
      confirm-class="btn-danger"
      :is-confirming="isDeleting"
      @close="confirmDelete.show = false"
      @confirm="executeDelete"
    >
      <p>¿Estás seguro de que deseas eliminar permanentemente a <strong>{{ party?.name }}</strong>?</p>
      <p class="mt-2 text-muted italic">Esta acción solo se completará si no existen documentos vinculados (pedidos, facturas, etc.).</p>
    </BaseDialog>
  </BaseEntityPage>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { 
  Building2, User, Pencil, Trash2, Fingerprint, MapPin, 
  History, Info, Settings, ShieldCheck, GitFork, Tag, 
  Contact2, List, LayoutDashboard, X, Save, RefreshCw, Percent
} from 'lucide-vue-next'
import BaseEntityPage from '@/components/shared/BaseEntityPage.vue'
import BasePageHeader from '@/components/shared/BasePageHeader.vue'
import FormSection from '@/components/shared/FormSection.vue'
import DataRow from '@/components/shared/DataRow.vue'
import BaseDialog from '@/components/shared/BaseDialog.vue'
import AddressManager from '@/components/party/AddressManager.vue'
import PersonManager from '@/components/party/PersonManager.vue'
import { partyApi } from '@/services/partyApi'
import { useToastStore } from '@/stores/toast'

const route = useRoute()
const router = useRouter()
const toastStore = useToastStore()

const party = ref(null)
const addresses = ref([])
const isLoading = ref(true)
const isSaving = ref(false)
const error = ref('')
const activeTab = ref('general')
const mode = ref('detail')

const formData = reactive({
  name: '',
  type: '',
  taxId: '',
  status: '',
  role: '',
  notes: '',
  defaultDiscountPercentage: 0
})

const tabs = [
  { id: 'general', label: 'General', icon: Info },
  { id: 'addresses', label: 'Direcciones', icon: MapPin },
  { id: 'contacts', label: 'Contactos', icon: ContactRound }
]

const confirmDelete = reactive({
  show: false
})

const primaryAddressLine = computed(() => {
  const primary = addresses.value.find(a => a.is_primary) || addresses.value[0]
  if (!primary) return 'Sin dirección'
  return `${primary.city}, ${primary.province}`
})

async function loadPartyData() {
  const id = route.params.id
  isLoading.value = true
  error.value = ''
  try {
    const [partyData, addrData] = await Promise.all([
      partyApi.getParty(id),
      partyApi.listAddresses(id)
    ])
    party.value = partyData
    addresses.value = addrData
  } catch (err) {
    error.value = 'No se pudo cargar la información de la entidad.'
    console.error(err)
  } finally {
    isLoading.value = false
  }
}

function enterEditMode() {
  Object.assign(formData, {
    name: party.value.name,
    type: party.value.type,
    taxId: party.value.tax_id,
    status: party.value.status,
    role: party.value.role,
    notes: party.value.notes || '',
    defaultDiscountPercentage: party.value.default_discount_percentage || 0
  })
  mode.value = 'edit'
  activeTab.value = 'general'
}

function exitEditMode() {
  mode.value = 'detail'
}

async function saveParty() {
  if (!formData.name) {
    toastStore.warning('El nombre de la entidad es obligatorio')
    return
  }
  
  isSaving.value = true
  try {
    const payload = {
      name: formData.name,
      role: formData.role,
      taxId: formData.taxId,
      status: formData.status,
      notes: formData.notes,
      default_discount_percentage: formData.defaultDiscountPercentage || 0
    }

    await partyApi.updateParty(party.value.id, payload)
    toastStore.success('Datos de entidad actualizados')
    await loadPartyData()
    mode.value = 'detail'
  } catch (err) {
    toastStore.error('Error al guardar: ' + err.message)
  } finally {
    isSaving.value = false
  }
}

function promptDelete() {
  confirmDelete.show = true
}

async function executeDelete() {
  isDeleting.value = true
  try {
    await partyApi.deleteParty(party.value.id)
    toastStore.success(`Entidad "${party.value.name}" eliminada`)
    router.push('/parties')
  } catch (err) {
    toastStore.error(err.message || 'Error al eliminar entidad')
    confirmDelete.show = false
  } finally {
    isDeleting.value = false
  }
}

function formatDate(d) { return d ? new Date(d).toLocaleDateString('es-ES') : '—' }

onMounted(() => loadPartyData())
</script>

<style scoped>
@import "@/design-system/_sections.css";

.action-toolbar { display: flex; align-items: center; padding: 0.75rem 1.5rem; background: white; border: 1px solid var(--color-border); border-radius: 8px; box-shadow: var(--box-shadow-sm); }
.roles-inline { display: flex; gap: 0.5rem; margin-left: 1.5rem; border-left: 1px solid var(--color-border); padding-left: 1.5rem; }
.role-pill { font-size: 0.65rem; font-weight: 800; text-transform: uppercase; padding: 0.2rem 0.6rem; border-radius: 4px; }
.role-pill.client { background: rgba(34, 197, 94, 0.1); color: #16a34a; }
.role-pill.supplier { background: rgba(147, 51, 234, 0.1); color: #9333ea; }

.overview-tags-row { display: flex; flex-wrap: wrap; gap: 1rem; }
.summary-tag { flex: 1; min-width: 240px; padding: 0.75rem 1.25rem; background: white; border: 1px solid var(--color-border); border-radius: 12px; display: flex; align-items: center; gap: 1rem; box-shadow: var(--box-shadow-sm); }
.summary-tag .icon { width: 40px; height: 40px; border-radius: 10px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.summary-tag .icon.blue { background: rgba(59, 130, 246, 0.1); color: #2563eb; }
.summary-tag .icon.yellow { background: rgba(230, 184, 0, 0.1); color: #d97706; }
.summary-tag .icon.purple { background: rgba(168, 85, 247, 0.1); color: #a855f7; }

.tag-content { display: flex; flex-direction: column; gap: 0.2rem; }
.tag-content label { font-size: 0.65rem; font-weight: 800; text-transform: uppercase; color: var(--color-text-secondary); letter-spacing: 0.05em; }
.tag-content strong { font-size: 1rem; color: var(--color-text-primary); font-weight: 700; }

.tabs-header { display: flex; gap: 0.5rem; border-bottom: 2px solid var(--color-border); padding-bottom: 0; }
.tab-btn { display: flex; align-items: center; gap: 0.6rem; padding: 0.75rem 1.5rem; border: none; background: transparent; cursor: pointer; color: var(--color-text-secondary); font-weight: 600; border-bottom: 3px solid transparent; transition: 0.2s; margin-bottom: -2px; }
.tab-btn:hover { color: var(--color-primary); }
.tab-btn.active { color: var(--color-primary); border-bottom-color: var(--color-primary); }

.notes-text { font-style: italic; color: var(--color-text-secondary); line-height: 1.5; margin: 0; }

/* Form Inline Styles */
.form-group label { display: block; font-size: 0.75rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); margin-bottom: 0.5rem; }
.form-input, .form-textarea { width: 100%; padding: 0.75rem 1rem; border-radius: 8px; border: 1px solid var(--color-border); font-family: inherit; }
.form-input:focus, .form-textarea:focus { outline: none; border-color: var(--color-primary); box-shadow: 0 0 0 3px rgba(230, 184, 0, 0.1); }
.form-input:disabled { background: var(--color-background); cursor: not-allowed; }
.help-text { font-size: 0.7rem; color: var(--color-text-secondary); margin-top: 0.25rem; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
.input-with-icon { position: relative; display: flex; align-items: center; }
.icon-start { position: absolute; left: 0.75rem; color: var(--color-text-secondary); }
.input-with-icon input { padding-left: 2.5rem; }
</style>
