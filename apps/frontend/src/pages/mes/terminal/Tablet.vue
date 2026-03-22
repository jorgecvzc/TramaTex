<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">MES / Terminal</p>
          <h1>Terminal de taller</h1>
          <p class="subtitle">Ejecución operativa de tareas de producción.</p>
        </div>
        <button @click="loadData" class="btn btn-secondary">Actualizar</button>
      </header>

      <section class="card filters">
        <input v-model="search" type="text" class="input" placeholder="Buscar por trabajo, número o cliente (nombre/referencia)" />
        <select v-model="taskFilter" class="input">
          <option value="">Todas las tareas</option>
          <option v-for="t in taskOptions" :key="t.id" :value="t.id">
            {{ t.name }}
          </option>
        </select>
        <select v-model="taskStatusFilter" class="input">
          <option value="">Todos los estados</option>
          <option value="ACTIVE">Activas (Pend. + En prog. + Bloq.)</option>
          <option value="PENDING">Pendiente</option>
          <option value="IN_PROGRESS">En progreso</option>
          <option value="BLOCKED">Bloqueada</option>
          <option value="COMPLETED">Completada</option>
          <option value="SKIPPED">Omitida</option>
        </select>
      </section>

      <section class="card">
        <div v-if="isLoading" class="empty-state">Cargando terminal...</div>
        <div v-else-if="error" class="alert">{{ error }}</div>
        <table v-else class="data-table">
          <thead>
            <tr>
              <th>Trabajo</th>
              <th>Tipo de trabajo</th>
              <th>Tarea</th>
              <th>Posición</th>
              <th>Archivo</th>
              <th>Estado</th>
              <th>Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in filteredRows" :key="row.taskId" class="row-clickable" @click="openDetail(row)">
              <td>
                <strong>{{ row.workNumber }}</strong><br />
                <small>{{ row.workName }}</small>
              </td>
              <td>{{ row.workTypeName }}</td>
              <td>
                <strong>{{ row.taskName }}</strong><br />
                <small>Sec. {{ row.sequence }}</small>
              </td>
              <td>{{ row.positionName || '—' }}</td>
              <td>
                <span v-if="row.designFilePath" class="row-file-path">📄 {{ row.designFilePath }}</span>
                <span v-else class="text-muted">—</span>
              </td>
              <td><span class="status-pill" :class="row.taskStatus.toLowerCase()">{{ taskStatusLabel(row.taskStatus) }}</span></td>
              <td class="actions-cell" @click.stop>
                <button
                  v-if="row.taskStatus === 'PENDING' || row.taskStatus === 'BLOCKED'"
                  class="btn-action btn-start"
                  @click="runAction(row, 'START')"
                >Iniciar</button>
                <button
                  v-if="row.taskStatus === 'IN_PROGRESS'"
                  class="btn-action btn-complete"
                  @click="runAction(row, 'COMPLETE')"
                >Completar</button>
                <button
                  v-if="row.taskStatus !== 'COMPLETED' && row.taskStatus !== 'SKIPPED'"
                  class="btn-action btn-block"
                  @click="runAction(row, 'BLOCK')"
                >Bloquear</button>
              </td>
            </tr>
            <tr v-if="!isLoading && filteredRows.length === 0">
              <td colspan="7" class="empty-state">No hay tareas para mostrar.</td>
            </tr>
          </tbody>
        </table>
      </section>
    </div>

    <!-- Task detail dialog -->
    <Teleport to="body">
      <div v-if="detailRow" class="dialog-overlay" @click.self="closeDetail">
        <div class="dialog" role="dialog" aria-modal="true">

          <!-- Header: task identity + action area -->
          <div class="dialog-header">
            <div class="dialog-header-top">
              <div class="task-heading">
                <!-- Proceso general -->
                <div class="work-badge">
                  <span class="work-number">{{ detailWork?.work_number ?? '—' }}</span>
                  <span class="work-badge-sep">·</span>
                  <span class="work-name">{{ detailWork?.work_name ?? '—' }}</span>
                  <span v-if="detailWork?.status" class="status-pill work-status-pill" :class="detailWork.status.toLowerCase()">{{ workStatusLabel(detailWork.status) }}</span>
                </div>
                <!-- Tipo + posición + secuencia -->
                <p class="task-meta">
                  <span class="meta-type">{{ detailRow.workTypeName }}</span>
                  <span v-if="positionNames[detailLine?.position_id ?? '']" class="meta-position">
                    · {{ positionNames[detailLine?.position_id ?? ''] }}
                  </span>
                  <span class="meta-seq">· Sec. {{ detailRow.sequence }}</span>
                </p>
                <button
                  v-if="detailLine?.design_file_path"
                  class="line-file-path"
                  :title="'Copiar ruta: ' + detailLine.design_file_path"
                  @click.stop="copyPath(detailLine.design_file_path!)"
                >📄 {{ detailLine.design_file_path }}</button>
                <h2 class="task-title"><span class="task-label">Tarea: </span>{{ detailRow.taskName }}<span class="status-pill task-status-inline" :class="detailRow.taskStatus.toLowerCase()">{{ taskStatusLabel(detailRow.taskStatus) }}</span></h2>
              </div>
              <button class="btn-close" @click="closeDetail" aria-label="Cerrar">✕</button>
            </div>
            <!-- Notas + botones de acción a ancho completo -->
            <div v-if="detailRow.taskStatus !== 'COMPLETED' && detailRow.taskStatus !== 'SKIPPED'" class="action-section">
              <textarea
                v-model="detailNotes"
                class="notes-input"
                placeholder="Notas (opcional)…"
                rows="2"
              ></textarea>
              <div class="action-buttons">
                <button
                  v-if="detailRow.taskStatus === 'PENDING' || detailRow.taskStatus === 'BLOCKED'"
                  class="btn btn-primary btn-action-lg"
                  @click="dialogAction('START')"
                >▶ Iniciar</button>
                <button
                  v-if="detailRow.taskStatus === 'IN_PROGRESS' || detailRow.taskStatus === 'BLOCKED'"
                  class="btn btn-success btn-action-lg"
                  @click="dialogAction('COMPLETE')"
                >✓ Completar</button>
                <button
                  class="btn btn-danger btn-action-lg"
                  @click="dialogAction('BLOCK')"
                >⊘ Bloquear</button>
              </div>
            </div>
          </div>

          <div class="dialog-body">

            <!-- Active line tasks -->
            <section v-if="detailLine" class="detail-section">
              <h3 class="section-title">Tareas de esta línea</h3>
              <div class="line-block active-line">
                <div class="line-header">
                  <strong>{{ workTypeNames[detailLine.work_type_id] || detailLine.work_type_id }}</strong>
                  <span v-if="positionNames[detailLine.position_id]" class="line-position">{{ positionNames[detailLine.position_id] }}</span>
                  <span v-if="detailLine.notes" class="line-notes">{{ detailLine.notes }}</span>
                </div>
                <table class="tasks-table">
                  <thead>
                    <tr><th>#</th><th>Tarea</th><th>Estado</th><th>Notas</th></tr>
                  </thead>
                  <tbody>
                    <tr
                      v-for="task in detailLine.tasks"
                      :key="task.id"
                      :class="{ 'task-current': task.id === detailRow.taskId }"
                    >
                      <td>{{ task.sequence }}</td>
                      <td>{{ taskNames[task.task_id] || task.task_id }}</td>
                      <td><span class="status-pill" :class="task.status.toLowerCase()">{{ taskStatusLabel(task.status) }}</span></td>
                      <td class="task-notes">{{ task.notes || '—' }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </section>

            <!-- Other lines: headers only -->
            <section v-if="otherLines.length > 0" class="detail-section">
              <h3 class="section-title">Otras líneas del trabajo</h3>
              <div v-for="line in otherLines" :key="line.id" class="other-line-row">
                <strong>{{ workTypeNames[line.work_type_id] || line.work_type_id }}</strong>
                <span v-if="positionNames[line.position_id]" class="line-position">{{ positionNames[line.position_id] }}</span>
                <span v-if="line.notes" class="line-notes">{{ line.notes }}</span>
                <button
                  v-if="line.design_file_path"
                  class="line-file-path"
                  :title="'Copiar ruta: ' + line.design_file_path"
                  @click.stop="copyPath(line.design_file_path!)"
                >📄 {{ line.design_file_path }}</button>
              </div>
            </section>

            <!-- Work info at bottom -->
            <section class="detail-section work-info-section">
              <h3 class="section-title">Trabajo</h3>
              <div class="info-grid">
                <div class="info-cell"><label>Número</label><span>{{ detailWork?.work_number }}</span></div>
                <div class="info-cell"><label>Nombre</label><span>{{ detailWork?.work_name }}</span></div>
                <div class="info-cell"><label>Cliente</label><span>{{ detailRow.partyName || '—' }}{{ detailRow.partyReference ? ` (${detailRow.partyReference})` : '' }}</span></div>
                <div class="info-cell"><label>Estado</label><span class="status-pill" :class="detailWork?.status?.toLowerCase()">{{ workStatusLabel(detailWork?.status) }}</span></div>
                <div class="info-cell"><label>Prioridad</label><span>{{ priorityLabel(detailWork?.priority) }}</span></div>
                <div v-if="detailWork?.due_date" class="info-cell"><label>Fecha límite</label><span>{{ formatDate(detailWork.due_date) }}</span></div>
              </div>
              <div v-if="detailWork?.notes" class="info-notes">
                <label>Notas del trabajo</label>
                <p>{{ detailWork.notes }}</p>
              </div>
            </section>

          </div>

          <!-- Footer: close only (actions are in body) -->
          <div class="dialog-footer">
            <div class="footer-actions">
              <button class="btn btn-secondary" @click="closeDetail">Cerrar</button>
            </div>
          </div>

        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { mesApi } from '@/services/mesApi'
import { partyApi } from '@/services/partyApi'
import type { WorkOrder, WorkOrderLine, WorkOrderTaskAction, MESWorkType } from '@/types/mes'

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

function workStatusLabel(status?: string) {
  const map: Record<string, string> = {
    PENDING: 'Pendiente',
    IN_PROGRESS: 'En progreso',
    COMPLETED: 'Completado',
    CANCELLED: 'Cancelado',
    SUSPENDED: 'Suspendido',
    OVERDUE: 'Vencido',
  }
  return status ? (map[status] ?? status) : '—'
}

function priorityLabel(priority?: string) {
  const map: Record<string, string> = {
    LOW: 'Baja',
    NORMAL: 'Normal',
    HIGH: 'Alta',
    URGENT: 'Urgente',
  }
  return priority ? (map[priority] ?? priority) : '—'
}

function formatDate(iso?: string) {
  if (!iso) return '—'
  return new Date(iso).toLocaleDateString('es-ES', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

const route = useRoute()

interface TaskRow {
  workId: string
  workNumber: string
  workName: string
  partyName: string
  partyReference: string
  workTypeId: string
  workTypeName: string
  lineId: string
  taskId: string
  taskTypeId: string
  taskName: string
  taskStatus: string
  sequence: number
  positionName: string
  designFilePath: string
}

const works = ref<WorkOrder[]>([])
const taskNames = ref<Record<string, string>>({})
const workTypeNames = ref<Record<string, string>>({})
const positionNames = ref<Record<string, string>>({})
const taskOptions = ref<{ id: string; name: string }[]>([])
const partiesCache = ref<Record<string, { name: string; reference: string }>>({})
const isLoading = ref(false)
const error = ref('')
const search = ref('')
const taskStatusFilter = ref('ACTIVE')
const taskFilter = ref('')

// Initialise filters from URL query params (e.g. /mes/terminal?status=PENDING)
const _initStatus = route.query.status as string | undefined
if (_initStatus) taskStatusFilter.value = _initStatus

const rows = computed<TaskRow[]>(() => {
  const result: TaskRow[] = []

  for (const work of works.value) {
    for (const line of work.lines || []) {
      for (const task of line.tasks || []) {
        result.push({
          workId: work.id,
          workNumber: work.work_number,
          workName: work.work_name,
          partyName: partiesCache.value[work.party_id]?.name || '',
          partyReference: partiesCache.value[work.party_id]?.reference || '',
          workTypeId: line.work_type_id,
          workTypeName: workTypeNames.value[line.work_type_id] || shortId(line.work_type_id),
          lineId: line.id,
          taskId: task.id,
          taskTypeId: task.task_id,
          taskName: taskNames.value[task.task_id] || shortId(task.task_id),
          taskStatus: task.status,
          sequence: task.sequence,
          positionName: positionNames.value[line.position_id] || '',
          designFilePath: line.design_file_path || '',
        })
      }
    }
  }

  return result
})

const filteredRows = computed(() => {
  const term = search.value.trim().toLowerCase()

  return rows.value.filter((row) => {
    const activeStatuses = ['PENDING', 'IN_PROGRESS', 'BLOCKED']
    const matchesStatus = !taskStatusFilter.value
      || (taskStatusFilter.value === 'ACTIVE' ? activeStatuses.includes(row.taskStatus) : row.taskStatus === taskStatusFilter.value)
    const matchesTask = !taskFilter.value || row.taskTypeId === taskFilter.value
    const matchesTerm =
      !term ||
      row.workNumber.toLowerCase().includes(term) ||
      row.workName.toLowerCase().includes(term) ||
      row.taskName.toLowerCase().includes(term) ||
      row.workTypeName.toLowerCase().includes(term) ||
      row.partyName.toLowerCase().includes(term) ||
      row.partyReference.toLowerCase().includes(term)

    return matchesStatus && matchesTask && matchesTerm
  })
})

function shortId(id?: string) {
  if (!id) return '—'
  return id.length > 8 ? `${id.slice(0, 8)}...` : id
}

async function loadData() {
  isLoading.value = true
  error.value = ''

  try {
    const [worksResult, tasksResult, workTypesResult, positionsResult] = await Promise.all([
      mesApi.listWorkOrders({}),
      mesApi.listTasks({}),
      mesApi.listWorkTypes({}),
      mesApi.listPositions({}),
    ])

    works.value = worksResult
    await loadPartyDetails()

    const names: Record<string, string> = {}
    const tOptions: { id: string; name: string }[] = []
    for (const task of tasksResult) {
      names[task.id] = task.name
      tOptions.push({ id: task.id, name: task.name })
    }
    taskNames.value = names
    taskOptions.value = tOptions

    const wtNames: Record<string, string> = {}
    for (const wt of workTypesResult as MESWorkType[]) {
      wtNames[wt.id] = wt.name
    }
    workTypeNames.value = wtNames

    const posNames: Record<string, string> = {}
    for (const pos of positionsResult as { id: string; name: string }[]) {
      posNames[pos.id] = pos.name
    }
    positionNames.value = posNames
  } catch (err: any) {
    error.value = err.message || 'No se pudo cargar el terminal'
  } finally {
    isLoading.value = false
  }
}

async function loadPartyDetails() {
  const partyIds = [...new Set(works.value.map((work) => work.party_id).filter(Boolean))]
  const uncachedIds = partyIds.filter((id) => !partiesCache.value[id])

  if (uncachedIds.length === 0) {
    return
  }

  try {
    const partiesMap = await partyApi.getPartiesBatch(uncachedIds)

    for (const partyId of uncachedIds) {
      const party = partiesMap[partyId]
      partiesCache.value[partyId] = {
        name: party?.name || '',
        reference: party?.reference || '',
      }
    }
  } catch (loadError) {
    console.error('Error loading MES terminal party details:', loadError)
  }
}

async function runAction(row: TaskRow, action: WorkOrderTaskAction) {
  let notes: string | undefined
  if (action === 'BLOCK' || action === 'COMPLETE') {
    const input = window.prompt(action === 'BLOCK' ? 'Motivo del bloqueo (opcional):' : 'Notas de finalizón (opcional):', '')
    if (input === null) return  // cancelado
    notes = input.trim() || undefined
  }

  try {
    await mesApi.updateWorkOrderTaskStatus(row.workId, row.taskId, { action, notes })
    await loadData()
  } catch (err: any) {
    error.value = err.message || 'No se pudo actualizar la tarea'
  }
}

// --- Dialog state ---
const detailRow = ref<TaskRow | null>(null)
const detailNotes = ref('')

const detailWork = computed<WorkOrder | undefined>(() =>
  detailRow.value ? works.value.find((w) => w.id === detailRow.value!.workId) : undefined,
)

const detailLine = computed(() =>
  detailWork.value?.lines.find((l) => l.id === detailRow.value?.lineId),
)

const otherLines = computed(() =>
  detailWork.value?.lines.filter((l) => l.id !== detailRow.value?.lineId) ?? [],
)

function openDetail(row: TaskRow) {
  detailRow.value = row
  detailNotes.value = ''
}

function closeDetail() {
  detailRow.value = null
  detailNotes.value = ''
}

async function dialogAction(action: WorkOrderTaskAction) {
  if (!detailRow.value) return
  const row = detailRow.value
  const notes = detailNotes.value.trim() || undefined

  try {
    await mesApi.updateWorkOrderTaskStatus(row.workId, row.taskId, { action, notes })
    closeDetail()
    await loadData()
  } catch (err: any) {
    error.value = err.message || 'No se pudo actualizar la tarea'
  }
}

onMounted(loadData)

function copyPath(path: string) {
  navigator.clipboard.writeText(path)
}
</script>

<style scoped>
.dashboard { min-height: 100vh; background-color: #f1f5f9; }
.dashboard-content { max-width: 1200px; margin: 0 auto; padding: 1rem; display: flex; flex-direction: column; gap: 1rem; }
.page-header { display: flex; justify-content: space-between; align-items: center; gap: 1rem; }
.breadcrumb { font-size: .75rem; text-transform: uppercase; letter-spacing: .08em; color: #64748b; margin: 0; }
.subtitle { color: #64748b; margin: .35rem 0 0; }
.card { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 1rem; }
.filters { display: grid; grid-template-columns: 1fr 260px 220px; gap: .75rem; }
.input { border: 1px solid #cbd5e1; border-radius: 8px; padding: .6rem .75rem; font-size: .9rem; }
.btn { border: none; border-radius: 8px; padding: .6rem .9rem; cursor: pointer; text-decoration: none; display: inline-flex; align-items: center; justify-content: center; font-size: .9rem; }
.btn-primary { background: #f4c430; color: #1b3a6b; font-weight: 600; }
.btn-secondary { background: #fff; border: 1px solid #cbd5e1; color: #1e293b; }
.btn-success { background: #16a34a; color: #fff; }
.btn-danger { background: #dc2626; color: #fff; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th, .data-table td { padding: .65rem; border-bottom: 1px solid #e2e8f0; text-align: left; vertical-align: top; }
.row-clickable { cursor: pointer; transition: background .15s; }
.row-clickable:hover { background: #f8fafc; }
.status-pill { display: inline-block; padding: .2rem .55rem; border-radius: 999px; font-size: .75rem; font-weight: 700; }
.status-pill.pending { background: #fef3c7; color: #92400e; }
.status-pill.in_progress { background: #dbeafe; color: #1d4ed8; }
.status-pill.blocked { background: #fee2e2; color: #991b1b; }
.status-pill.completed { background: #dcfce7; color: #166534; }
.status-pill.skipped { background: #e2e8f0; color: #334155; }
.empty-state { text-align: center; color: #64748b; padding: 1rem; }
.alert { background: #fef2f2; color: #b91c1c; border: 1px solid #fecaca; border-radius: 8px; padding: .75rem; }

/* Dialog */
.dialog-overlay { position: fixed; inset: 0; background: rgba(15,23,42,.45); display: flex; align-items: center; justify-content: center; z-index: 1000; padding: 1rem; }
.dialog { background: #fff; border-radius: 14px; width: 80vw; max-width: 900px; max-height: 90vh; display: flex; flex-direction: column; overflow: hidden; box-shadow: 0 20px 60px rgba(0,0,0,.25); }
.dialog-header { display: flex; flex-direction: column; padding: 1.25rem 1.5rem 1rem; border-bottom: 1px solid #e2e8f0; background: #f8fafc; gap: .75rem; }
.dialog-header-top { display: flex; justify-content: space-between; align-items: flex-start; gap: 1rem; }
.task-heading { display: flex; flex-direction: column; gap: .6rem; }
.work-badge { display: flex; align-items: center; gap: .5rem; background: #1b3a6b; border-radius: 8px; padding: .4rem 1rem; width: fit-content; margin-bottom: .5rem; }
.work-number { font-size: .95rem; font-weight: 800; color: #f4c430; letter-spacing: .04em; }
.work-badge-sep { color: #94a3b8; font-size: .85rem; }
.work-name { font-size: .9rem; font-weight: 600; color: #e2e8f0; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 380px; }
.work-status-pill { font-size: .75rem; padding: .15rem .55rem; margin-left: .5rem; vertical-align: middle; opacity: .92; }
.task-meta { display: flex; align-items: center; gap: .35rem; flex-wrap: wrap; font-size: 1.05rem; margin: 0; }
.meta-type { color: #64748b; text-transform: uppercase; letter-spacing: .06em; font-weight: 600; }
.meta-position { color: #1d4ed8; font-weight: 700; font-size: 1.1rem; text-transform: none; letter-spacing: 0; }
.meta-seq { color: #94a3b8; font-size: .95rem; }
.task-title { font-size: 1.25rem; font-weight: 700; color: #0f172a; margin: .5rem 0 0; }
.task-label { font-size: 1rem; font-weight: 400; color: #64748b; margin-right: .35rem; }
.task-status-inline { margin-left: 1.5rem; vertical-align: middle; }
.btn-close { background: none; border: none; font-size: 1.1rem; cursor: pointer; color: #64748b; padding: .25rem .5rem; border-radius: 6px; line-height: 1; flex-shrink: 0; }
.btn-close:hover { background: #e2e8f0; }
.dialog-body { overflow-y: auto; padding: 1.25rem 1.5rem; display: flex; flex-direction: column; gap: 1.5rem; }
.detail-section { display: flex; flex-direction: column; gap: .75rem; }
.section-title { font-size: .8rem; text-transform: uppercase; letter-spacing: .08em; color: #64748b; margin: 0; border-bottom: 1px solid #e2e8f0; padding-bottom: .4rem; }
.info-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(160px, 1fr)); gap: .65rem; }
.info-cell { display: flex; flex-direction: column; gap: .15rem; }
.info-cell label { font-size: .7rem; text-transform: uppercase; letter-spacing: .05em; color: #94a3b8; }
.info-cell span { font-size: .9rem; color: #0f172a; font-weight: 500; }
.info-notes { display: flex; flex-direction: column; gap: .2rem; }
.info-notes label { font-size: .7rem; text-transform: uppercase; letter-spacing: .05em; color: #94a3b8; }
.info-notes p { font-size: .9rem; color: #334155; background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 8px; padding: .6rem .75rem; margin: 0; }
.line-block { border: 1px solid #e2e8f0; border-radius: 8px; overflow: hidden; }
.line-header { display: flex; align-items: center; gap: .5rem; padding: .6rem .85rem; background: #f8fafc; border-bottom: 1px solid #e2e8f0; font-size: .9rem; flex-wrap: wrap; }
.line-notes { color: #64748b; font-size: .85rem; }
.line-file-ref { display: flex; align-items: center; margin-left: auto; }
.line-file-path { font-size: .9rem; font-family: monospace; color: #1d4ed8; background: #f1f5f9; border: 1px solid #e2e8f0; border-radius: 4px; padding: .1rem .35rem; max-width: 260px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; cursor: pointer; text-align: left; }
.line-file-path:hover { background: #eff6ff; border-color: #93c5fd; }
.row-file-path { font-family: monospace; color: #1d4ed8; font-size: .75rem; max-width: 180px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; display: inline-block; }
.text-muted { color: #94a3b8; }
.tasks-table { width: 100%; border-collapse: collapse; font-size: .85rem; }
.tasks-table th, .tasks-table td { padding: .5rem .85rem; border-bottom: 1px solid #f1f5f9; text-align: left; }
.tasks-table th { color: #64748b; font-weight: 600; font-size: .75rem; text-transform: uppercase; letter-spacing: .05em; background: #fff; }
.tasks-table tr:last-child td { border-bottom: none; }
.task-notes { color: #475569; font-size: .82rem; max-width: 260px; }
.task-selected { background: #eff6ff !important; }
.task-selected td { font-weight: 600; }
.dialog-footer { padding: .75rem 1.5rem; border-top: 1px solid #e2e8f0; background: #f8fafc; }
.notes-input { border: 1px solid #cbd5e1; border-radius: 8px; padding: .6rem .75rem; font-size: .9rem; resize: vertical; font-family: inherit; width: 100%; box-sizing: border-box; }
.action-section { background: #f8fafc; border-top: 1px solid #e2e8f0; padding-top: .75rem; display: flex; flex-direction: column; gap: .75rem; }
.action-buttons { display: flex; gap: .5rem; flex-wrap: wrap; margin-bottom: 1rem; }
.btn-action-lg { padding: .45rem .9rem; font-size: .82rem; font-weight: 700; }
.task-header-row { display: flex; align-items: center; gap: .5rem; flex-wrap: wrap; margin-top: .2rem; }
.active-line { border-color: #3b82f6; }
.active-line .line-header { background: #eff6ff; border-bottom-color: #bfdbfe; }
.task-current { background: #eff6ff !important; }
.task-current td { font-weight: 700; }
.line-position { font-size: .82rem; color: #1d4ed8; font-weight: 600; background: #eff6ff; border: 1px solid #bfdbfe; border-radius: 4px; padding: .1rem .4rem; }
.other-line-row { display: flex; align-items: center; gap: .5rem; flex-wrap: wrap; padding: .55rem .85rem; border: 1px solid #e2e8f0; border-radius: 8px; font-size: .9rem; background: #f8fafc; }
.work-info-section { opacity: .85; }
.footer-actions { display: flex; justify-content: flex-end; gap: .5rem; flex-wrap: wrap; }

.actions-cell { white-space: nowrap; display: flex; gap: .3rem; align-items: center; padding: .5rem .65rem; }
.btn-action { border: none; border-radius: 6px; padding: .3rem .65rem; font-size: .78rem; font-weight: 600; cursor: pointer; white-space: nowrap; }
.btn-start { background: #dbeafe; color: #1d4ed8; }
.btn-start:hover { background: #bfdbfe; }
.btn-complete { background: #dcfce7; color: #15803d; }
.btn-complete:hover { background: #bbf7d0; }
.btn-block { background: #fee2e2; color: #b91c1c; }
.btn-block:hover { background: #fecaca; }

@media (max-width: 900px) {
  .filters { grid-template-columns: 1fr; }
  .info-grid { grid-template-columns: 1fr 1fr; }
}
@media (max-width: 600px) {
  .dialog { max-height: 100vh; border-radius: 0; }
  .info-grid { grid-template-columns: 1fr; }
}
</style>
