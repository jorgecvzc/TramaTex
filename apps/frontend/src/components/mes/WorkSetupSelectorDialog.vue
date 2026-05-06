<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { Info, List, PlusCircle, SearchX, Plus, Trash2, AlertCircle } from 'lucide-vue-next'
import BaseDialog from '@/components/shared/BaseDialog.vue'
import PartySelector from '@/components/party/PartySelector.vue'
import BaseSkeleton from '@/components/shared/BaseSkeleton.vue'
import { mesApi } from '@/services/mesApi'
import { productApi } from '@/services/productApi'
import { useToastStore } from '@/stores/toast'
import type { WorkOrder, WorkSetup, MESWorkType, MESPosition } from '@/types/mes'

interface Props {
  show: boolean
  workOrder: WorkOrder | null
}

const props = defineProps<Props>()
const emit = defineEmits(['close', 'assigned'])

const toastStore = useToastStore()
const mode = ref<'assign' | 'create'>('assign')
const isLoading = ref(false)
const isSaving = ref(false)
const error = ref('')
const fieldErrors = ref<Record<string, string>>({})

// ... (rest of state remains same)
const existingWorkSetups = ref<WorkSetup[]>([])
const selectedWorkSetupId = ref('')

// Create mode state
const workTypes = ref<MESWorkType[]>([])
const positions = ref<MESPosition[]>([])
const productGroups = ref<any[]>([])
const newSetupForm = ref({
  name: '',
  party_id: '',
  tangible_group_id: '',
  description: '',
  is_active: true,
  lines: [] as Array<{ work_type_id: string; position_id: string; sequence: number; design_file_path?: string }>,
})

const tangibleGroups = computed(() => productGroups.value.filter(g => g.type === 'TANGIBLE'))

async function loadExistingSetups() {
  const partyId = props.workOrder?.party_id || (props.workOrder as any)?.partyId
  if (!partyId) return
  isLoading.value = true
  try {
    existingWorkSetups.value = await mesApi.listWorkSetups({ 
      party_id: partyId, 
      is_active: true 
    })
  } catch (err: any) {
    error.value = 'Error al cargar configuraciones existentes'
  } finally {
    isLoading.value = false
  }
}

async function loadMasters() {
  try {
    const [wt, pos, groups] = await Promise.all([
      mesApi.listWorkTypes({ is_active: true }),
      mesApi.listPositions({ is_active: true }),
      productApi.listProductGroups({ isActive: true }),
    ])
    workTypes.value = wt
    positions.value = pos
    productGroups.value = groups.data || groups
  } catch (err: any) {
    error.value = 'Error al cargar maestros para nueva configuración'
  }
}

function handleModeChange(newMode: 'assign' | 'create') {
  mode.value = newMode
  error.value = ''
  fieldErrors.value = {}
  if (newMode === 'create' && workTypes.value.length === 0) {
    loadMasters()
  }
}

function addSetupLine() {
  newSetupForm.value.lines.push({ 
    work_type_id: '', 
    position_id: '', 
    sequence: newSetupForm.value.lines.length + 1 
  })
}

function removeSetupLine(index: number) {
  newSetupForm.value.lines.splice(index, 1)
  // Re-sequence
  newSetupForm.value.lines.forEach((l, i) => l.sequence = i + 1)
}

function validateCreateMode(): boolean {
  fieldErrors.value = {}
  const f = newSetupForm.value
  let isValid = true

  if (!f.name?.trim()) {
    fieldErrors.value.name = 'El nombre es obligatorio'
    isValid = false
  }
  if (!f.tangible_group_id) {
    fieldErrors.value.tangible_group_id = 'La categoría tangible es obligatoria'
    isValid = false
  }
  if (!f.lines?.length) {
    error.value = 'Debe añadir al menos un paso de producción'
    isValid = false
  } else {
    const hasInvalidLine = f.lines.some(l => !l.work_type_id || !l.position_id)
    if (hasInvalidLine) {
      error.value = 'Todos los pasos de producción deben tener tipo y posición'
      isValid = false
    }
  }

  return isValid
}

async function handleConfirm() {
  if (!props.workOrder || isSaving.value) return
  error.value = ''
  fieldErrors.value = {}
  
  let setupId = selectedWorkSetupId.value

  if (mode.value === 'create') {
    if (!validateCreateMode()) return
    
    isSaving.value = true
    try {
      const f = newSetupForm.value
      console.log('Creando nueva configuración MES...', f)
      const created = await mesApi.createWorkSetup({
        ...f,
        lines: f.lines.map(l => ({ ...l }))
      })
      setupId = created.id
      console.log('Configuración creada con éxito:', setupId)
    } catch (err: any) {
      console.error('Error detallado en creación MES:', err)
      error.value = err.message || 'Error al crear la configuración'
      isSaving.value = false
      return
    }
  }

  if (!setupId) {
    error.value = 'Debes seleccionar o crear una configuración válida'
    return
  }

  isSaving.value = true
  try {
    console.log(`Asignando configuración ${setupId} a la orden ${props.workOrder.id}...`)
    await mesApi.updateWorkOrder(props.workOrder.id, { 
      work_setup_id: setupId,
      status: 'PENDING' 
    })
    
    toastStore.success('Orden de trabajo configurada correctamente')
    emit('assigned', setupId)
    emit('close')
  } catch (err: any) {
    console.error('Error detallado en asignación MES:', err)
    error.value = err.message || 'Error al asignar la configuración'
  } finally {
    isSaving.value = false
  }
}

function initializeState() {
  const wo = props.workOrder
  error.value = ''
  fieldErrors.value = {}
  mode.value = 'assign'
  selectedWorkSetupId.value = ''
  
  if (wo) {
    const partyId = wo.party_id || (wo as any).partyId
    if (partyId) {
      newSetupForm.value = {
        name: `Configuración para ${wo.work_name || (wo as any).description || 'nueva orden'}`,
        party_id: partyId,
        tangible_group_id: '',
        description: '',
        is_active: true,
        lines: []
      }
      loadExistingSetups()
    }
  }
}

watch(() => props.show, (newVal) => {
  if (newVal) {
    initializeState()
  }
})

watch(() => props.workOrder, () => {
  if (props.show) {
    initializeState()
  }
})

onMounted(() => {
  if (props.show) {
    initializeState()
  }
})
</script>

<template>
  <BaseDialog
    :show="show"
    title="Configurar Orden de Trabajo"
    :confirm-text="mode === 'assign' ? 'Asignar Configuración' : 'Crear y Asignar'"
    :is-confirming="isSaving"
    size="xl"
    @close="$emit('close')"
    @confirm="handleConfirm"
  >
    <div class="setup-selector-content">
      <!-- Banner informativo -->
      <div class="info-banner mb-4">
        <Info :size="18" />
        <p>Configurando orden <strong>{{ workOrder?.work_number }}</strong> para el cliente <strong>{{ workOrder?.party_name || workOrder?.party_id }}</strong></p>
      </div>

      <!-- Selector de modo -->
      <div class="mode-tabs mb-6">
        <button 
          class="tab-btn" 
          :class="{ active: mode === 'assign' }" 
          @click="handleModeChange('assign')"
        >
          <List :size="18" />
          Asignar Existente
        </button>
        <button 
          class="tab-btn" 
          :class="{ active: mode === 'create' }" 
          @click="handleModeChange('create')"
        >
          <PlusCircle :size="18" />
          Crear Nueva Configuración
        </button>
      </div>

      <div v-if="error" class="alert-box alert-danger mb-4">
        <AlertCircle :size="20" />
        <span>{{ error }}</span>
      </div>

      <!-- MODO: ASIGNAR EXISTENTE -->
      <div v-if="mode === 'assign'" class="assign-mode">
        <div v-if="isLoading" class="setup-list">
          <p class="section-label mb-3">Buscando configuraciones...</p>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
            <div v-for="i in 2" :key="i" class="setup-card" style="opacity: 0.7;">
              <div class="card-radio"><BaseSkeleton type="circle" width="20px" /></div>
              <div class="card-info w-full">
                <BaseSkeleton type="title" width="70%" />
                <BaseSkeleton type="text" width="40%" />
              </div>
            </div>
          </div>
        </div>
        <div v-else-if="existingWorkSetups.length === 0" class="empty-state p-8 text-center card bg-light">
          <SearchX class="text-muted mb-2 mx-auto" :size="48" />
          <p>No se han encontrado configuraciones previas para este cliente.</p>
          <button class="btn btn-outline-primary btn-sm mt-4" @click="handleModeChange('create')">
            Crear la primera configuración
          </button>
        </div>
        <div v-else class="setup-list">
          <p class="section-label mb-3">Selecciona una configuración base:</p>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
            <div 
              v-for="setup in existingWorkSetups" 
              :key="setup.id" 
              class="setup-card"
              :class="{ selected: selectedWorkSetupId === setup.id }"
              @click="selectedWorkSetupId = setup.id"
            >
              <div class="card-radio">
                <div class="radio-circle"></div>
              </div>
              <div class="card-info">
                <span class="setup-name">{{ setup.name }}</span>
                <span class="setup-meta">{{ setup.lines?.length || 0 }} pasos de producción</span>
                <p v-if="setup.description" class="setup-desc mt-1">{{ setup.description }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- MODO: CREAR NUEVA -->
      <div v-if="mode === 'create'" class="create-mode">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
          <div class="form-group">
            <label class="form-label">Nombre del Setup *</label>
            <input v-model="newSetupForm.name" type="text" class="form-input" :class="{ 'is-invalid': fieldErrors.name }" placeholder="Ej: Estampado Camiseta Algodón" />
            <span v-if="fieldErrors.name" class="field-error">{{ fieldErrors.name }}</span>
          </div>
          <div class="form-group">
            <label class="form-label">Categoría Tangible *</label>
            <select v-model="newSetupForm.tangible_group_id" class="form-input" :class="{ 'is-invalid': fieldErrors.tangible_group_id }">
              <option value="">-- Seleccionar categoría --</option>
              <option v-for="group in tangibleGroups" :key="group.id" :value="group.id">{{ group.name }}</option>
            </select>
            <span v-if="fieldErrors.tangible_group_id" class="field-error">{{ fieldErrors.tangible_group_id }}</span>
          </div>
          <div class="form-group md:col-span-2">
            <label class="form-label">Descripción / Observaciones</label>
            <textarea v-model="newSetupForm.description" class="form-textarea" rows="2" placeholder="Detalles técnicos generales..."></textarea>
          </div>
        </div>

        <div class="setup-lines-section">
          <div class="section-header-row mb-3">
            <h4 class="m-0">Pasos de Producción *</h4>
            <button type="button" class="btn btn-outline btn-sm" @click="addSetupLine">
              <Plus :size="16" /> Añadir Paso
            </button>
          </div>

          <div class="table-wrapper border rounded">
            <table class="data-table">
              <thead>
                <tr>
                  <th style="width: 60px">Orden</th>
                  <th>Tipo de Trabajo</th>
                  <th>Posición / Máquina</th>
                  <th style="width: 50px"></th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(line, idx) in newSetupForm.lines" :key="idx">
                  <td class="text-center font-bold">{{ line.sequence }}</td>
                  <td>
                    <select v-model="line.work_type_id" class="form-input-sm w-full">
                      <option value="">-- Tipo --</option>
                      <option v-for="wt in workTypes" :key="wt.id" :value="wt.id">{{ wt.name }}</option>
                    </select>
                  </td>
                  <td>
                    <select v-model="line.position_id" class="form-input-sm w-full">
                      <option value="">-- Posición --</option>
                      <option v-for="pos in positions" :key="pos.id" :value="pos.id">{{ pos.name }} ({{ pos.code }})</option>
                    </select>
                  </td>
                  <td>
                    <button type="button" class="btn-icon text-danger" @click="removeSetupLine(idx)">
                      <Trash2 :size="18" />
                    </button>
                  </td>
                </tr>
                <tr v-if="newSetupForm.lines.length === 0">
                  <td colspan="4" class="p-4 text-center text-muted italic">
                    No se han añadido pasos. Define el flujo de trabajo para esta orden.
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </BaseDialog>
</template>

<style scoped>
.setup-selector-content { min-height: 400px; }

.info-banner { display: flex; align-items: center; gap: 0.75rem; padding: 0.75rem 1rem; background: var(--color-background); border: 1px solid var(--color-border); border-radius: 8px; color: var(--color-text-secondary); font-size: 0.9rem; }
.info-banner :deep(svg) { color: var(--color-primary); }
.info-banner p { margin: 0; }

.mode-tabs { display: flex; border-bottom: 1px solid var(--color-border); gap: 1rem; }
.tab-btn { display: flex; align-items: center; gap: 0.5rem; padding: 0.75rem 1.5rem; border: none; background: none; cursor: pointer; color: var(--color-text-secondary); font-weight: 600; border-bottom: 3px solid transparent; transition: 0.2s; }
.tab-btn:hover { color: var(--color-primary); }
.tab-btn.active { color: var(--color-primary); border-bottom-color: var(--color-primary); }

.setup-card { display: flex; align-items: flex-start; gap: 1rem; padding: 1rem; border: 1px solid var(--color-border); border-radius: 10px; cursor: pointer; transition: 0.2s; background: white; }
.setup-card:hover { border-color: var(--color-primary); background: var(--color-background); }
.setup-card.selected { border-color: var(--color-primary); border-width: 2px; background: rgba(230, 184, 0, 0.05); }

.card-radio { width: 20px; height: 20px; border: 2px solid var(--color-border); border-radius: 50%; margin-top: 2px; position: relative; flex-shrink: 0; }
.setup-card.selected .card-radio { border-color: var(--color-primary); }
.setup-card.selected .radio-circle { width: 10px; height: 10px; background: var(--color-primary); border-radius: 50%; position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); }

.card-info { display: flex; flex-direction: column; }
.setup-name { font-weight: 700; color: var(--color-text-primary); }
.setup-meta { font-size: 0.75rem; color: var(--color-text-secondary); font-weight: 600; text-transform: uppercase; }
.setup-desc { font-size: 0.85rem; color: var(--color-text-secondary); }

.section-label { font-size: 0.8rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); letter-spacing: 0.05em; }
.section-header-row { display: flex; justify-content: space-between; align-items: center; }

.form-group { display: flex; flex-direction: column; gap: 0.4rem; }
.form-label { font-size: 0.8rem; font-weight: 700; color: var(--color-text-secondary); }
.form-input, .form-textarea { width: 100%; padding: 0.6rem 0.75rem; border: 1px solid var(--color-border); border-radius: 8px; font-family: inherit; }
.form-input-sm { padding: 0.4rem; border: 1px solid var(--color-border); border-radius: 4px; font-size: 0.85rem; }

.is-invalid { border-color: var(--color-error) !important; }
.field-error { font-size: 0.75rem; color: var(--color-error); font-weight: 600; margin-top: 0.1rem; }

.alert-box { display: flex; align-items: center; gap: 0.75rem; padding: 0.75rem 1rem; border-radius: 8px; font-size: 0.9rem; font-weight: 500; }
.alert-danger { background: #fef2f2; border: 1px solid #fee2e2; color: #b91c1c; }

.btn-icon { background: none; border: none; cursor: pointer; padding: 0.25rem; border-radius: 4px; display: flex; align-items: center; justify-content: center; }
.btn-icon:hover { background: rgba(0,0,0,0.05); }

.spinner { width: 30px; height: 30px; border: 3px solid var(--color-border); border-top-color: var(--color-primary); border-radius: 50%; animation: spin 1s linear infinite; margin: 0 auto; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>
