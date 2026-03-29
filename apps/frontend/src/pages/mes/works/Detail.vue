<template>
  <div class="dashboard">
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">MES / Órdenes de Trabajo</p>
          <h1>{{ work?.work_name || 'Detalle de orden' }}</h1>
          <p class="subtitle subtitle-ref">{{ work?.work_number || '—' }}</p>
          <p class="subtitle subtitle-setup">{{ workSetupName || work?.work_setup_id || '—' }}<span v-if="work" class="status-pill status-pill-header" :class="work.status.toLowerCase()">{{ mesApi.getWorkStatusLabel(work.status) }}</span></p>
        </div>
        <div class="header-actions" v-if="work && !isEditing">
          <button class="btn btn-primary" @click="startEdit">Editar</button>
          <RouterLink to="/mes/work-orders" class="btn btn-secondary">Volver</RouterLink>
        </div>
        <div class="header-actions" v-else>
          <RouterLink to="/mes/work-orders" class="btn btn-secondary">Volver</RouterLink>
        </div>
      </header>

      <section class="card" v-if="isLoading">
        <div class="empty-state">Cargando detalle...</div>
      </section>

      <section class="card" v-else-if="error">
        <div class="alert">{{ error }}</div>
      </section>

      <template v-else-if="work">
        <section class="card" v-if="isEditing">
          <h3>Editar orden de trabajo</h3>
          <div class="edit-grid">
            <div>
              <label class="label">Nombre</label>
              <input v-model="editForm.work_name" class="input" type="text" />
            </div>
            <div>
              <label class="label">Estado</label>
              <select v-model="editForm.status" class="input">
                <option value="PENDING">Pendiente</option>
                <option value="IN_PROGRESS">En progreso</option>
                <option value="ON_HOLD">En espera</option>
                <option value="COMPLETED">Completado</option>
                <option value="CANCELLED">Cancelado</option>
              </select>
            </div>
            <div>
              <label class="label">Prioridad</label>
              <select v-model="editForm.priority" class="input">
                <option value="LOW">Baja</option>
                <option value="NORMAL">Normal</option>
                <option value="HIGH">Alta</option>
                <option value="URGENT">Urgente</option>
              </select>
            </div>
            <div>
              <label class="label">Fecha de vencimiento</label>
              <input v-model="editForm.due_date" class="input" type="date" />
            </div>
            <div class="full">
              <label class="label">Notas</label>
              <textarea v-model="editForm.notes" class="input" rows="3" />
            </div>
          </div>
          <div class="edit-actions">
            <button class="btn btn-secondary" @click="cancelEdit" :disabled="isSaving">Cancelar</button>
            <button class="btn btn-primary" @click="saveChanges" :disabled="isSaving">
              {{ isSaving ? 'Guardando...' : 'Guardar cambios' }}
            </button>
          </div>
        </section>

        <section class="card details-grid">
          <div><strong>Prioridad:</strong> {{ mesApi.getPriorityLabel(work.priority) }}</div>
          <div><strong>Cliente:</strong> {{ partyName || work.party_id }}</div>
          <div><strong>Configuración:</strong> {{ workSetupName || work.work_setup_id }}</div>
          <div v-if="work.due_date"><strong>Vencimiento:</strong> {{ formatDate(work.due_date) }}</div>
          <div class="full"><strong>Notas:</strong> {{ work.notes || '—' }}</div>
        </section>

        <section class="card">
          <h3>Asignaciones de tipo de trabajo</h3>
          <div v-if="work.lines.length === 0" class="empty-state">Sin asignaciones</div>
          <div v-for="line in work.lines" :key="line.id" class="group-box">
            <p><strong>Tipo de trabajo:</strong> {{ workTypeNames[line.work_type_id] || line.work_type_id }}</p>
            <p><strong>Posición:</strong> {{ positionNames[line.position_id] || line.position_id }}</p>
            <p><strong>Secuencia:</strong> {{ line.sequence }}</p>
            <p><strong>Notas:</strong> {{ line.notes || '—' }}</p>
            <table v-if="line.tasks.length > 0" class="tasks-table">
              <thead>
                <tr><th>#</th><th>Tarea</th><th>Estado</th><th>Notas</th></tr>
              </thead>
              <tbody>
                <tr v-for="task in line.tasks" :key="task.id">
                  <td>{{ task.sequence }}</td>
                  <td>{{ taskNames[task.task_id] || task.task_id }}</td>
                  <td><span class="status-pill" :class="task.status.toLowerCase()">{{ taskStatusLabel(task.status) }}</span></td>
                  <td class="task-notes">{{ task.notes || '—' }}</td>
                </tr>
              </tbody>
            </table>
            <div v-if="line.design_file_path" class="file-ref">
              <span class="file-ref-label">Archivo:</span>
              <button class="file-ref-path" :title="'Copiar ruta: ' + line.design_file_path" @click="copyPath(line.design_file_path!)">{{ line.design_file_path }}</button>
            </div>
          </div>
        </section>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { mesApi } from '@/services/mesApi'
import { partyApi } from '@/services/partyApi'
import type { MESPosition, MESWorkType, WorkOrder } from '@/types/mes'

const route = useRoute()
const isLoading = ref(false)
const error = ref('')
const work = ref<WorkOrder | null>(null)
const partyName = ref('')
const workSetupName = ref('')
const workTypeNames = ref<Record<string, string>>({})
const positionNames = ref<Record<string, string>>({})
const taskNames = ref<Record<string, string>>({})

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
const isEditing = ref(false)
const isSaving = ref(false)
const editForm = ref({
  work_name: '',
  status: 'PENDING',
  priority: 'NORMAL',
  due_date: '',
  notes: '',
})

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
  return date.toLocaleDateString('es-ES')
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

onMounted(loadDetail)

function copyPath(path: string) {
  navigator.clipboard.writeText(path)
}
</script>

<style scoped>
.dashboard { min-height: 100vh; background-color: #f1f5f9; }
.dashboard-content { max-width: 1100px; margin: 0 auto; padding: 2rem; display: flex; flex-direction: column; gap: 1.5rem; }
.page-header { display: flex; justify-content: space-between; align-items: center; gap: 1rem; }
.header-actions { display: flex; gap: .5rem; }
.breadcrumb { font-size: .75rem; text-transform: uppercase; letter-spacing: .08em; color: #64748b; margin: 0 0 .6rem; }
.subtitle { color: #64748b; margin: 0; }
.page-header h1 { margin: 0; line-height: 1.2; }
.subtitle-ref { font-size: 1rem; font-weight: 600; color: #334155; margin-top: 0; }
.subtitle-setup { font-size: .875rem; margin-top: 1rem; }
.card { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 1rem; }
.details-grid { display: grid; grid-template-columns: 1fr 1fr; gap: .75rem; }
.details-grid .full { grid-column: 1 / -1; }
.edit-grid { display: grid; grid-template-columns: 1fr 1fr; gap: .75rem; }
.edit-grid .full { grid-column: 1 / -1; }
.label { display: block; font-size: .85rem; color: #334155; margin-bottom: .35rem; font-weight: 600; }
.input { border: 1px solid #cbd5e1; border-radius: 8px; padding: .6rem .75rem; font: inherit; width: 100%; }
.edit-actions { margin-top: .75rem; display: flex; justify-content: flex-end; gap: .5rem; }
.group-box { border: 1px solid #e2e8f0; border-radius: 8px; padding: .75rem; margin-top: .75rem; }
.group-box p { margin: .25rem 0; }
.tasks-table { width: 100%; border-collapse: collapse; margin-top: .6rem; font-size: .875rem; }
.tasks-table th { text-align: left; padding: .35rem .5rem; border-bottom: 2px solid #e2e8f0; color: #64748b; font-weight: 600; font-size: .8rem; }
.tasks-table td { padding: .35rem .5rem; border-bottom: 1px solid #f1f5f9; vertical-align: top; }
.tasks-table tr:last-child td { border-bottom: none; }
.task-notes { color: #475569; font-size: .82rem; max-width: 320px; }
.file-ref { display: flex; align-items: center; gap: .4rem; margin-top: .4rem; flex-wrap: wrap; }
.file-ref-label { font-size: .85rem; font-weight: 600; color: #334155; white-space: nowrap; }
.file-ref-path { font-size: .8rem; font-family: monospace; color: #1d4ed8; background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 4px; padding: .15rem .4rem; max-width: 400px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; cursor: pointer; text-align: left; }
.file-ref-path:hover { background: #eff6ff; border-color: #93c5fd; }
.btn { border: none; border-radius: 8px; padding: .6rem 1rem; cursor: pointer; text-decoration: none; display: inline-flex; align-items: center; justify-content: center; width: fit-content; }
.btn-primary { background: #f4c430; color: #1b3a6b; font-weight: 600; }
.btn-secondary { background: #fff; border: 1px solid #e2e8f0; color: #1e293b; }
.empty-state { text-align: center; color: #64748b; padding: 1rem; }
.alert { background: #fef2f2; color: #b91c1c; border: 1px solid #fecaca; border-radius: 8px; padding: .75rem; }
.status-pill { display: inline-flex; align-items: center; padding: .2rem .65rem; border-radius: 999px; font-size: .78rem; font-weight: 700; letter-spacing: .03em; white-space: nowrap; }
.status-pill.pending { background: #fef9c3; color: #854d0e; }
.status-pill.in_progress { background: #dbeafe; color: #1d4ed8; }
.status-pill.completed { background: #dcfce7; color: #166534; }
.status-pill.cancelled { background: #fee2e2; color: #b91c1c; }
.status-pill.on_hold { background: #f3f4f6; color: #374151; }
.status-pill.overdue { background: #fee2e2; color: #b91c1c; }
.status-pill.suspended { background: #fef3c7; color: #92400e; }
.status-pill-header { margin-left: .75rem; vertical-align: middle; font-size: .82rem; }
</style>
