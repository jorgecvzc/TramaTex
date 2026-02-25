<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">MES / Terminal</p>
          <h1>Terminal de Taller (Tablet)</h1>
          <p class="subtitle">Ejecución operativa de tareas de producción.</p>
        </div>
        <button @click="loadData" class="btn btn-secondary">Actualizar</button>
      </header>

      <section class="card filters">
        <input v-model="search" type="text" class="input" placeholder="Buscar por trabajo, número o cliente (nombre/referencia)" />
        <select v-model="taskStatusFilter" class="input">
          <option value="">Todos los estados</option>
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
              <th>Tarea</th>
              <th>Estado</th>
              <th>Asignado</th>
              <th>Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in filteredRows" :key="row.taskId">
              <td>
                <strong>{{ row.workNumber }}</strong><br />
                <small>{{ row.workName }}</small>
              </td>
              <td>
                <strong>{{ row.taskName }}</strong><br />
                <small>Secuencia {{ row.sequence }}</small>
              </td>
              <td><span class="status-pill" :class="row.taskStatus.toLowerCase()">{{ row.taskStatus }}</span></td>
              <td>{{ row.assignedTo || '—' }}</td>
              <td>
                <div class="actions">
                  <button
                    v-if="row.taskStatus === 'PENDING' || row.taskStatus === 'BLOCKED'"
                    @click="runAction(row, 'START')"
                    class="btn btn-primary btn-sm"
                  >Iniciar</button>
                  <button
                    v-if="row.taskStatus === 'IN_PROGRESS'"
                    @click="runAction(row, 'PAUSE')"
                    class="btn btn-secondary btn-sm"
                  >Pausar</button>
                  <button
                    v-if="row.taskStatus === 'IN_PROGRESS' || row.taskStatus === 'BLOCKED'"
                    @click="runAction(row, 'COMPLETE')"
                    class="btn btn-success btn-sm"
                  >Completar</button>
                  <button
                    v-if="row.taskStatus !== 'COMPLETED' && row.taskStatus !== 'SKIPPED'"
                    @click="runAction(row, 'BLOCK')"
                    class="btn btn-danger btn-sm"
                  >Bloquear</button>
                </div>
              </td>
            </tr>
            <tr v-if="!isLoading && filteredRows.length === 0">
              <td colspan="5" class="empty-state">No hay tareas para mostrar.</td>
            </tr>
          </tbody>
        </table>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import Navbar from '@/components/layout/Navbar.vue'
import { mesApi } from '@/services/mesApi'
import { partyApi } from '@/services/partyApi'
import type { MESWork, MESWorkTaskAction } from '@/types/mes'

interface TaskRow {
  workId: string
  workNumber: string
  workName: string
  partyName: string
  partyReference: string
  taskId: string
  taskName: string
  taskStatus: string
  sequence: number
  assignedTo: string
}

const works = ref<MESWork[]>([])
const taskNames = ref<Record<string, string>>({})
const partiesCache = ref<Record<string, { name: string; reference: string }>>({})
const isLoading = ref(false)
const error = ref('')
const search = ref('')
const taskStatusFilter = ref('')

const rows = computed<TaskRow[]>(() => {
  const result: TaskRow[] = []

  for (const work of works.value) {
    for (const group of work.service_groups || []) {
      for (const task of group.tasks || []) {
        result.push({
          workId: work.id,
          workNumber: work.work_number,
          workName: work.work_name,
          partyName: partiesCache.value[work.party_id]?.name || '',
          partyReference: partiesCache.value[work.party_id]?.reference || '',
          taskId: task.id,
          taskName: taskNames.value[task.task_id] || shortId(task.task_id),
          taskStatus: task.status,
          sequence: task.sequence,
          assignedTo: task.assigned_to ? shortId(task.assigned_to) : '',
        })
      }
    }
  }

  return result
})

const filteredRows = computed(() => {
  const term = search.value.trim().toLowerCase()

  return rows.value.filter((row) => {
    const matchesStatus = !taskStatusFilter.value || row.taskStatus === taskStatusFilter.value
    const matchesTerm =
      !term ||
      row.workNumber.toLowerCase().includes(term) ||
      row.workName.toLowerCase().includes(term) ||
      row.taskName.toLowerCase().includes(term) ||
      row.partyName.toLowerCase().includes(term) ||
      row.partyReference.toLowerCase().includes(term)

    return matchesStatus && matchesTerm
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
    const [worksResult, tasksResult] = await Promise.all([
      mesApi.listWorkDefinitions({}),
      mesApi.listTasks({}),
    ])

    works.value = worksResult
    await loadPartyDetails()

    const names: Record<string, string> = {}
    for (const task of tasksResult) {
      names[task.id] = task.name
    }
    taskNames.value = names
  } catch (err: any) {
    error.value = err.message || 'No se pudo cargar el terminal MES'
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

async function runAction(row: TaskRow, action: MESWorkTaskAction) {
  let notes: string | undefined
  if (action === 'BLOCK' || action === 'COMPLETE') {
    const input = window.prompt('Notas (opcional):', '')
    if (input !== null) {
      notes = input.trim() || undefined
    }
  }

  try {
    await mesApi.updateWorkTaskStatus(row.workId, row.taskId, { action, notes })
    await loadData()
  } catch (err: any) {
    error.value = err.message || 'No se pudo actualizar la tarea'
  }
}

onMounted(loadData)
</script>

<style scoped>
.dashboard { min-height: 100vh; background-color: #f1f5f9; }
.dashboard-content { max-width: 1200px; margin: 0 auto; padding: 1rem; display: flex; flex-direction: column; gap: 1rem; }
.page-header { display: flex; justify-content: space-between; align-items: center; gap: 1rem; }
.breadcrumb { font-size: .75rem; text-transform: uppercase; letter-spacing: .08em; color: #64748b; margin: 0; }
.subtitle { color: #64748b; margin: .35rem 0 0; }
.card { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 1rem; }
.filters { display: grid; grid-template-columns: 1fr 260px; gap: .75rem; }
.input { border: 1px solid #cbd5e1; border-radius: 8px; padding: .6rem .75rem; }
.btn { border: none; border-radius: 8px; padding: .6rem .9rem; cursor: pointer; text-decoration: none; display: inline-flex; align-items: center; justify-content: center; }
.btn-sm { padding: .4rem .65rem; font-size: .8rem; }
.btn-primary { background: #f4c430; color: #1b3a6b; font-weight: 600; }
.btn-secondary { background: #fff; border: 1px solid #cbd5e1; color: #1e293b; }
.btn-success { background: #16a34a; color: #fff; }
.btn-danger { background: #dc2626; color: #fff; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th, .data-table td { padding: .65rem; border-bottom: 1px solid #e2e8f0; text-align: left; vertical-align: top; }
.actions { display: flex; flex-wrap: wrap; gap: .35rem; }
.status-pill { display: inline-block; padding: .2rem .55rem; border-radius: 999px; font-size: .75rem; font-weight: 700; }
.status-pill.pending { background: #fef3c7; color: #92400e; }
.status-pill.in_progress { background: #dbeafe; color: #1d4ed8; }
.status-pill.blocked { background: #fee2e2; color: #991b1b; }
.status-pill.completed { background: #dcfce7; color: #166534; }
.status-pill.skipped { background: #e2e8f0; color: #334155; }
.empty-state { text-align: center; color: #64748b; padding: 1rem; }
.alert { background: #fef2f2; color: #b91c1c; border: 1px solid #fecaca; border-radius: 8px; padding: .75rem; }
@media (max-width: 900px) {
  .filters { grid-template-columns: 1fr; }
}
</style>
