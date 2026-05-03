<template>
  <!-- PÁGINA: DETALLE / EDICIÓN / ALTA DE ENTIDAD -->
  
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
    <!-- CAPA 1: IDENTIDAD -->
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
            <div v-else class="header-actions-group">
              <button class="btn btn-outline" @click="exitEditMode" :disabled="isSaving">
                <X :size="18" />
                <span>Cancelar</span>
              </button>
              <button class="btn btn-secondary ml-2" @click="saveParty" :disabled="isSaving">
                <component :is="isSaving ? RefreshCw : Save" :size="18" :class="{ 'spin': isSaving }" />
                <span>{{ mode === 'create' ? 'Crear Entidad' : 'Guardar Cambios' }}</span>
              </button>
            </div>
          </template>
        </BasePageHeader>

        <!-- NAVEGACIÓN POR PESTAÑAS (Solo en detalle) -->
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

    <!-- CAPA 2: CONTEXTO (Resumen) -->
    <template #summary v-if="mode !== 'create' && party">
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
    <div class="party-master-content">
      <!-- TAB: GENERAL / FORMULARIO -->
      <div v-if="activeTab === 'general'" class="tab-fade-in">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <!-- SECCIÓN: IDENTIDAD -->
          <FormSection title="Datos de Identidad" icon="person">
            <div v-if="mode === 'detail'">
              <DataRow label="Nombre Comercial / Completo" :value="party?.name" icon="badge" />
              <DataRow label="Tipo de Entidad" :value="party?.type === 'ORGANIZATION' ? 'Empresa / Organización' : 'Persona Física'" icon="users" />
              <DataRow label="Identificador Fiscal" :value="party?.tax_id" is-mono icon="fingerprint" />
              <DataRow label="Notas Internas" icon="notes">
                <p class="notes-text">{{ party?.notes || 'Sin observaciones registradas.' }}</p>
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
                  <select v-model="formData.type" class="form-input" :disabled="mode === 'edit'">
                    <option value="ORGANIZATION">Empresa / Organización</option>
                    <option value="PERSON">Persona Física</option>
                  </select>
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

          <!-- SECCIÓN: CONFIGURACIÓN -->
          <FormSection title="Configuración de Cuenta" icon="settings">
            <div v-if="mode === 'detail'">
              <DataRow label="Estado Actual" :value="party?.status" icon="shield-check" />
              <DataRow label="Rol en el Sistema" :value="party?.role" icon="git-fork" />
              <DataRow label="Descuento por Defecto" :value="(party?.default_discount_percentage || 0) + '%'" icon="tag" />
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

      <!-- TABS SECUNDARIAS (Solo modo detalle) -->
      <div v-if="activeTab === 'addresses' && mode === 'detail'" class="tab-fade-in">
        <AddressManager :party-id="party?.id" />
      </div>

      <div v-if="activeTab === 'contacts' && mode === 'detail'" class="tab-fade-in">
        <PersonManager :party-id="party?.id" />
      </div>
    </div>

    <!-- DIÁLOGO DE ELIMINACIÓN -->
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
  status: 'ACTIVE',
  role: 'CLIENT',
  notes: '',
  defaultDiscountPercentage: 0
})

const tabs = [
  { id: 'general', label: 'General', icon: Info },
  { id: 'addresses', label: 'Direcciones', icon: MapPin },
  { id: 'contacts', label: 'Contactos', icon: ContactRound }
]

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
    name: '', type: 'ORGANIZATION', taxId: '', status: 'ACTIVE',
    role: 'CLIENT', notes: '', defaultDiscountPercentage: 0
  })
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
    const payload = {
      name: formData.name,
      type: formData.type,
      role: formData.role,
      taxId: formData.taxId,
      status: formData.status,
      notes: formData.notes,
      default_discount_percentage: formData.defaultDiscountPercentage || 0
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
    toastStore.error('Error al guardar: ' + err.message)
  } finally {
    isSaving.value = false
  }
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

.overview-tags-row { display: flex; flex-wrap: wrap; gap: 1rem; }
.summary-tag { flex: 1; min-width: 240px; padding: 0.75rem 1.25rem; background: white; border: 1px solid var(--color-border); border-radius: 12px; display: flex; align-items: center; gap: 1rem; box-shadow: var(--box-shadow-sm); }
.summary-tag .icon { width: 40px; height: 40px; border-radius: 10px; display: flex; align-items: center; justify-content: center; }
.summary-tag .icon.blue { background: rgba(59, 130, 246, 0.1); color: #2563eb; }
.summary-tag .icon.yellow { background: rgba(230, 184, 0, 0.1); color: #d97706; }
.summary-tag .icon.purple { background: rgba(168, 85, 247, 0.1); color: #a855f7; }
.tag-content { display: flex; flex-direction: column; line-height: 1.2; }
.tag-content label { font-size: 0.65rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); }
.tag-content strong { font-size: 1rem; color: var(--color-text-primary); }

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
