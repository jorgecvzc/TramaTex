<template>
  <BaseEntityPage :is-loading="isLoading" :error="error">
    <!-- 1. IDENTITY HEADER -->
    <template #header>
      <BasePageHeader 
        :title="work?.work_name || 'Detalle de Orden'" 
        :subtitle="work?.work_number"
        :breadcrumbs="[{ label: 'MES', to: '/mes/dashboard' }, { label: 'Órdenes', to: '/mes/work-orders' }, { label: work?.work_number || 'Detalle' }]"
        show-back
      >
        <template #icon>
          <span class="material-symbols-outlined">precision_manufacturing</span>
        </template>
        <template #actions>
          <button v-if="!isEditing" class="btn btn-primary" @click="startEdit">
            <span class="material-symbols-outlined">edit</span>
            <span>Editar Orden</span>
          </button>
          <template v-else>
            <button class="btn btn-outline" @click="cancelEdit" :disabled="isSaving">Cancelar</button>
            <button class="btn btn-secondary" @click="saveChanges" :disabled="isSaving">
              <span class="material-symbols-outlined">{{ isSaving ? 'sync' : 'save' }}</span>
              <span>{{ isSaving ? 'Guardando...' : 'Guardar Cambios' }}</span>
            </button>
          </template>
        </template>
      </BasePageHeader>
    </template>

    <!-- 2. TOOLBAR (Summary state) -->
    <template #toolbar v-if="work && !isEditing">
      <div class="action-toolbar card">
        <div class="toolbar-info">
          <span :class="['status-badge', `status-${getStatusClass(work.status)}`]">
            {{ mesApi.getWorkStatusLabel(work.status) }}
          </span>
          <span :class="['priority-pill', `prio-${work.priority}`]">{{ mesApi.getPriorityLabel(work.priority) }}</span>
        </div>
        <div class="toolbar-actions">
          <button v-if="work.status === 'PENDING' && work.work_setup_id" class="btn btn-success btn-sm" @click="startWorkOrder">
            <span class="material-symbols-outlined">play_arrow</span> Iniciar Producción
          </button>
          <button v-if="work.status === 'IN_PROGRESS'" class="btn btn-outline btn-sm" @click="router.push('/mes/terminal')">
            <span class="material-symbols-outlined">tablet_mac</span> Ir a Terminal
          </button>
        </div>
      </div>
    </template>

    <!-- 3. SUMMARY TAGS -->
    <template #summary v-if="work && !isEditing">
      <div class="overview-tags-row">
        <div class="summary-tag">
          <div class="icon blue"><span class="material-symbols-outlined">person</span></div>
          <div class="tag-content"><label>Cliente</label><strong>{{ partyName || '—' }}</strong></div>
        </div>
        <div class="summary-tag">
          <div class="icon yellow"><span class="material-symbols-outlined">settings_suggest</span></div>
          <div class="tag-content"><label>Configuración</label><strong>{{ workSetupName || 'No configurada' }}</strong></div>
        </div>
        <div class="summary-tag">
          <div class="icon purple"><span class="material-symbols-outlined">calendar_today</span></div>
          <div class="tag-content"><label>Vencimiento</label><strong :class="{'text-danger': isOverdue}">{{ formatDate(work.due_date) }}</strong></div>
        </div>
      </div>
    </template>

    <!-- 4. MAIN CONTENT -->
    <div v-if="work" class="work-detail-content">
      
      <!-- Edit Form -->
      <FormSection v-if="isEditing" title="Información de la Orden" icon="edit_note">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div class="form-group">
            <label class="form-label">Nombre del Trabajo</label>
            <input v-model="editForm.work_name" class="form-input" type="text" />
          </div>
          <div class="form-group">
            <label class="form-label">Prioridad</label>
            <select v-model="editForm.priority" class="form-input">
              <option value="LOW">Baja</option>
              <option value="NORMAL">Normal</option>
              <option value="HIGH">Alta</option>
              <option value="URGENT">Urgente</option>
            </select>
          </div>
          <div class="form-group">
            <label class="form-label">Estado de Producción</label>
            <select v-model="editForm.status" class="form-input">
              <option value="PENDING">Pendiente</option>
              <option value="IN_PROGRESS">En progreso</option>
              <option value="ON_HOLD">En espera</option>
              <option value="COMPLETED">Completado</option>
              <option value="CANCELLED">Cancelado</option>
            </select>
          </div>
          <div class="form-group">
            <label class="form-label">Fecha de Vencimiento</label>
            <input v-model="editForm.due_date" class="form-input" type="date" />
          </div>
          <div class="form-group md:col-span-2">
            <label class="form-label">Notas / Instrucciones</label>
            <textarea v-model="editForm.notes" class="form-textarea" rows="3" />
          </div>
        </div>
      </FormSection>

      <!-- Read Only Details -->
      <template v-else>
        <FormSection title="Información General" icon="info">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <DataRow label="Cliente" :value="partyName" icon="person" />
            <DataRow label="Configuración Técnica" :value="workSetupName || 'Sin asignar'" icon="settings_suggest" />
            <DataRow label="Referencia de Pedido" :value="work.sales_order_number ? `#${work.sales_order_number}` : 'Venta Directa'" icon="shopping_cart" />
            <DataRow label="Fecha de Vencimiento" :value="formatDate(work.due_date)" icon="calendar_today" :class="{'text-danger font-bold': isOverdue}" />
            <DataRow label="Notas Técnicas" icon="notes" class="md:col-span-2">
              <p class="notes-text">{{ work.notes || 'Sin observaciones.' }}</p>
            </DataRow>
          </div>
        </FormSection>

        <FormSection title="Ruta de Producción" icon="account_tree">
          <div v-if="work.lines.length === 0" class="empty-state p-8">
            <span class="material-symbols-outlined text-muted" style="font-size: 3rem">route</span>
            <p>No se han definido pasos para esta orden.</p>
          </div>
          
          <div v-for="line in work.lines" :key="line.id" class="production-line-card border rounded-lg p-4 mb-4">
            <div class="line-header flex justify-between items-center mb-4">
              <div class="line-identity flex items-center gap-3">
                <span class="line-seq">{{ line.sequence }}</span>
                <div class="line-title">
                  <h4 class="m-0 text-primary">{{ workTypeNames[line.work_type_id] || 'Cargando...' }}</h4>
                  <p class="text-xs text-muted uppercase font-bold">{{ positionNames[line.position_id] || 'Cargando...' }}</p>
                </div>
              </div>
              <div class="line-meta">
                <p v-if="line.notes" class="text-sm italic text-muted m-0">"{{ line.notes }}"</p>
              </div>
            </div>

            <div v-if="line.tasks.length > 0" class="table-wrapper">
              <table class="data-table">
                <thead>
                  <tr>
                    <th style="width: 50px">#</th>
                    <th>Tarea de Taller</th>
                    <th>Estado</th>
                    <th>Notas</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="task in line.tasks" :key="task.id">
                    <td class="text-center font-mono">{{ task.sequence }}</td>
                    <td><strong>{{ taskNames[task.task_id] || task.task_id }}</strong></td>
                    <td>
                      <span :class="['status-badge', `status-${getTaskStatusClass(task.status)}`]">
                        {{ taskStatusLabel(task.status) }}
                      </span>
                    </td>
                    <td class="text-xs text-muted italic">{{ task.notes || '—' }}</td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div v-if="line.design_file_path" class="design-file-row flex items-center gap-2 mt-4 p-2 bg-light rounded">
              <span class="material-symbols-outlined text-secondary">draft</span>
              <span class="text-xs font-bold text-muted">ARCHIVO:</span>
              <code class="text-xs flex-1 truncate">{{ line.design_file_path }}</code>
              <button class="btn btn-ghost btn-sm" @click="copyPath(line.design_file_path)">
                <span class="material-symbols-outlined">content_copy</span>
              </button>
            </div>
          </div>
        </FormSection>
      </template>
    </div>
  </BaseEntityPage>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import BaseEntityPage from '@/components/shared/BaseEntityPage.vue'
import BasePageHeader from '@/components/shared/BasePageHeader.vue'
import FormSection from '@/components/shared/FormSection.vue'
import DataRow from '@/components/shared/DataRow.vue'
import { mesApi } from '@/services/mesApi'
import { partyApi } from '@/services/partyApi'
import type { MESPosition, MESWorkType, WorkOrder } from '@/types/mes'

const route = useRoute()
const router = useRouter()
const isLoading = ref(false)
const error = ref('')
const work = ref<WorkOrder | null>(null)
const partyName = ref('')
const workSetupName = ref('')
const workTypeNames = ref<Record<string, string>>({})
const positionNames = ref<Record<string, string>>({})
const taskNames = ref<Record<string, string>>({})

const isEditing = ref(false)
const isSaving = ref(false)
const editForm = ref({
  work_name: '',
  status: 'PENDING',
  priority: 'NORMAL',
  due_date: '',
  notes: '',
})

const isOverdue = computed(() => {
  if (!work.value?.due_date || work.value.status === 'COMPLETED') return false
  return new Date(work.value.due_date) < new Date()
})

function getStatusClass(status: string) {
  const map: Record<string, string> = { PENDING: 'warning', IN_PROGRESS: 'info', COMPLETED: 'success', CANCELLED: 'danger', ON_HOLD: 'secondary' }
  return map[status] || 'secondary'
}

function getTaskStatusClass(status: string) {
  const map: Record<string, string> = { PENDING: 'secondary', IN_PROGRESS: 'info', COMPLETED: 'success', BLOCKED: 'danger', SKIPPED: 'muted' }
  return map[status] || 'secondary'
}

function taskStatusLabel(status: string) {
  const map: Record<string, string> = {
    PENDING: 'Pendiente',
    IN_PROGRESS: 'En progreso',
    BLOCKED: 'Bloqueada',
    COMPLETED: 'Completada',
    SKIPPED: 'Omitida',
  }
  return map[status] ?? status
}

function dateToInput(value?: string) {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return date.toISOString().slice(0, 10)
}

function startEdit() {
  if (!work.value) return
  editForm.value = {
    work_name: work.value.work_name || '',
    status: work.value.status || 'PENDING',
    priority: work.value.priority || 'NORMAL',
    due_date: dateToInput(work.value.due_date),
    notes: work.value.notes || '',
  }
  isEditing.value = true
}

function cancelEdit() {
  isEditing.value = false
}

async function saveChanges() {
  if (!work.value) return
  isSaving.value = true
  error.value = ''
  try {
    const updated = await mesApi.updateWorkOrder(work.value.id, {
      work_name: editForm.value.work_name.trim(),
      status: editForm.value.status,
      priority: editForm.value.priority,
      due_date: editForm.value.due_date || undefined,
      notes: editForm.value.notes.trim(),
    })
    work.value = updated
    isEditing.value = false
  } catch (err: any) {
    error.value = err.message || 'No se pudo actualizar la orden de trabajo'
  } finally {
    isSaving.value = false
  }
}

async function startWorkOrder() {
  if (!work.value) return
  try {
    await mesApi.updateWorkOrder(work.value.id, { status: 'IN_PROGRESS' })
    await loadDetail()
  } catch (err) {}
}

async function loadLookupData(currentWork: WorkOrder) {
  const [partyResult, wsResult, templates, positions, tasks] = await Promise.all([
    partyApi.getParty(currentWork.party_id),
    currentWork.work_setup_id ? mesApi.getWorkSetup(currentWork.work_setup_id) : Promise.resolve(null),
    mesApi.listWorkTypes({}),
    mesApi.listPositions({}),
    mesApi.listTasks({}),
  ])

  partyName.value = partyResult?.name || ''
  workSetupName.value = wsResult?.name || ''

  const templateMap: Record<string, string> = {}
  for (const wt of templates as MESWorkType[]) {
    templateMap[wt.id] = wt.name
  }
  workTypeNames.value = templateMap

  const positionMap: Record<string, string> = {}
  for (const position of positions as MESPosition[]) {
    positionMap[position.id] = `${position.name} (${position.code})`
  }
  positionNames.value = positionMap

  const taskMap: Record<string, string> = {}
  for (const task of tasks as { id: string; name: string }[]) {
    taskMap[task.id] = task.name
  }
  taskNames.value = taskMap
}

function formatDate(value?: string) {
  if (!value) return '—'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleDateString('es-ES', { year: 'numeric', month: 'short', day: 'numeric' })
}

async function loadDetail() {
  isLoading.value = true
  error.value = ''
  try {
    work.value = await mesApi.getWorkOrder(String(route.params.id))
    await loadLookupData(work.value)
  } catch (err: any) {
    error.value = err.message || 'No se pudo cargar la orden de trabajo'
  } finally {
    isLoading.value = false
  }
}

function copyPath(path: string) {
  navigator.clipboard.writeText(path)
}

onMounted(loadDetail)
</script>

<style scoped>
@import "@/design-system/_sections.css";

.overview-tags-row { display: flex; flex-wrap: wrap; gap: 1rem; }
.summary-tag { flex: 1; min-width: 240px; padding: 0.6rem 1rem; background: white; border: 1px solid var(--color-border); border-radius: 12px; display: flex; align-items: center; gap: 0.75rem; box-shadow: var(--box-shadow-sm); }

.icon { width: 36px; height: 36px; border-radius: 8px; display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
.icon .material-symbols-outlined { font-size: 22px; }
.icon.blue { background: rgba(59, 130, 246, 0.1); color: #2563eb; }
.icon.yellow { background: rgba(230, 184, 0, 0.1); color: #d97706; }
.icon.purple { background: rgba(168, 85, 247, 0.1); color: #9333ea; }

.tag-content { display: flex; flex-direction: column; gap: 0.15rem; }
.tag-content label { font-size: 0.65rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); }
.tag-content strong { font-size: 0.95rem; color: var(--color-text-primary); }

.action-toolbar { display: flex; justify-content: space-between; align-items: center; padding: 0.75rem 1.5rem; background: white; border: 1px solid var(--color-border); border-radius: 8px; box-shadow: var(--box-shadow-sm); }
.toolbar-info { display: flex; align-items: center; gap: 1rem; }
.status-badge { padding: 0.4rem 1rem; font-size: 0.85rem; font-weight: 800; }
.priority-pill { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; padding: 0.2rem 0.6rem; border-radius: 4px; }
.prio-HIGH, .prio-URGENT { background: rgba(220, 38, 38, 0.1); color: #dc2626; }
.prio-NORMAL { background: rgba(217, 119, 6, 0.1); color: #d97706; }
.prio-LOW { background: rgba(37, 99, 235, 0.1); color: #2563eb; }

.work-detail-content { display: flex; flex-direction: column; gap: 1.5rem; }

.form-group { display: flex; flex-direction: column; gap: 0.4rem; }
.form-label { font-size: var(--font-size-xs); font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); }
.form-input, .form-textarea { width: 100%; padding: 0.75rem 1rem; border-radius: 8px; border: 1px solid var(--color-border); font-family: inherit; }

.line-seq { width: 28px; height: 28px; background: var(--color-primary); color: var(--color-text-on-primary); border-radius: 50%; display: flex; align-items: center; justify-content: center; font-weight: 800; font-size: 0.85rem; }
.production-line-card { background: white; transition: 0.2s; }
.production-line-card:hover { border-color: var(--color-primary); box-shadow: var(--box-shadow-sm); }

.bg-light { background-color: var(--color-background); }
.notes-text { color: var(--color-text-primary); line-height: 1.5; margin: 0; }
</style>
