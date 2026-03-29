<template>
  <BaseTerminalPage
    title="Terminal de Taller"
    station-id="TALLER CENTRAL"
    icon="tablet_mac"
    :is-loading="isLoading"
    @refresh="loadData"
    @close="goBack"
  >
    <!-- CAPA DE ESTADO RÁPIDO -->
    <div class="status-ribbon">
      <div class="status-item">
        <label>Tareas en Cola</label>
        <strong>{{ counts.pending }}</strong>
      </div>
      <div class="status-item active">
        <label>En Marcha</label>
        <strong>{{ counts.inProgress }}</strong>
      </div>
      <div class="terminal-search-wrap">
        <span class="material-symbols-outlined">search</span>
        <input v-model="search" type="text" placeholder="Escanear código o buscar trabajo..." class="terminal-input" />
      </div>
    </div>

    <!-- LISTADO TÁCTIL -->
    <div class="task-list-container">
      <div v-if="filteredRows.length === 0" class="terminal-empty">
        <span class="material-symbols-outlined">checklist</span>
        <p>No hay tareas pendientes asignadas a esta estación.</p>
      </div>

      <div 
        v-for="row in filteredRows" 
        :key="row.taskId" 
        class="task-row-card"
        :class="['status-' + row.taskStatus.toLowerCase()]"
        @click="openDetail(row)"
      >
        <div class="card-main-info">
          <div class="work-ref">
            <code class="terminal-code">{{ row.workNumber }}</code>
            <strong>{{ row.workName }}</strong>
          </div>
          <div class="task-info">
            <span class="task-name">{{ row.taskName }}</span>
            <span class="work-type">{{ row.workTypeName }}</span>
          </div>
        </div>

        <div class="card-meta">
          <div v-if="row.positionName" class="terminal-pos">
            <span class="material-symbols-outlined">location_on</span>
            {{ row.positionName }}
          </div>
          <span :class="['terminal-status-pill', row.taskStatus.toLowerCase()]">
            {{ taskStatusLabel(row.taskStatus) }}
          </span>
        </div>

        <div class="card-quick-actions" @click.stop>
          <button
            v-if="['PENDING', 'BLOCKED'].includes(row.taskStatus)"
            class="btn-action start"
            @click="runAction(row, 'START')"
          >
            <span class="material-symbols-outlined">play_arrow</span>
            INICIAR
          </button>
          <button
            v-if="row.taskStatus === 'IN_PROGRESS'"
            class="btn-action complete"
            @click="runAction(row, 'COMPLETE')"
          >
            <span class="material-symbols-outlined">check_circle</span>
            FINALIZAR
          </button>
        </div>
      </div>
    </div>

    <!-- DIÁLOGO TÁCTIL -->
    <BaseDialog
      :show="!!detailRow"
      :title="detailRow ? `Tarea: ${detailRow.taskName}` : ''"
      icon="assignment"
      size="lg"
      hide-actions
      @close="closeDetail"
    >
      <div class="terminal-dialog-body" v-if="detailRow">
        <div class="dialog-info-grid">
          <div class="info-block">
            <label>Orden de Trabajo</label>
            <div class="val">{{ detailRow.workNumber }} - {{ detailRow.workName }}</div>
          </div>
          <div class="info-block">
            <label>Ubicación</label>
            <div class="val">{{ detailRow.positionName || 'No asignada' }}</div>
          </div>
        </div>

        <div class="terminal-textarea-wrap">
          <label>Notas / Observaciones del Operario</label>
          <textarea v-model="detailNotes" placeholder="Escriba aquí cualquier incidencia..."></textarea>
        </div>

        <div class="terminal-giant-actions">
          <button
            v-if="['PENDING', 'BLOCKED'].includes(detailRow.taskStatus)"
            class="btn-giant start"
            @click="dialogAction('START')"
          >
            <span class="material-symbols-outlined">play_arrow</span>
            INICIAR TRABAJO
          </button>
          <button
            v-if="detailRow.taskStatus === 'IN_PROGRESS'"
            class="btn-giant complete"
            @click="dialogAction('COMPLETE')"
          >
            <span class="material-symbols-outlined">check_circle</span>
            COMPLETAR TAREA
          </button>
          <button
            class="btn-giant block"
            @click="dialogAction('BLOCK')"
          >
            <span class="material-symbols-outlined">block</span>
            BLOQUEAR
          </button>
        </div>
      </div>
    </BaseDialog>
  </BaseTerminalPage>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import BaseTerminalPage from '@/components/shared/BaseTerminalPage.vue';
import BaseDialog from '@/components/shared/BaseDialog.vue';
import { mesApi } from '@/services/mesApi';
import type { WorkOrder, WorkOrderTaskAction, MESWorkType } from '@/types/mes';

const router = useRouter();
const isLoading = ref(true);
const search = ref('');
const works = ref<WorkOrder[]>([]);
const taskNames = ref<Record<string, string>>({});
const workTypeNames = ref<Record<string, string>>({});
const positionNames = ref<Record<string, string>>({});

const goBack = () => router.push('/mes/dashboard');

// --- DATOS ---
const rows = computed(() => {
  const res: any[] = [];
  for (const w of works.value) {
    for (const l of w.lines || []) {
      for (const t of l.tasks || []) {
        res.push({
          workId: w.id, workNumber: w.work_number, workName: w.work_name,
          workTypeName: workTypeNames.value[l.work_type_id] || l.work_type_id,
          taskId: t.id, taskName: taskNames.value[t.task_id] || t.task_id,
          taskStatus: t.status, sequence: t.sequence,
          positionName: positionNames.value[l.position_id] || '',
        });
      }
    }
  }
  return res;
});

const filteredRows = computed(() => {
  const term = search.value.toLowerCase();
  return rows.value.filter(r => {
    const matchesStatus = ['PENDING', 'IN_PROGRESS', 'BLOCKED'].includes(r.taskStatus);
    const matchesTerm = !term || r.workNumber.toLowerCase().includes(term) || r.workName.toLowerCase().includes(term);
    return matchesStatus && matchesTerm;
  });
});

const counts = computed(() => ({
  pending: rows.value.filter(r => r.taskStatus === 'PENDING').length,
  inProgress: rows.value.filter(r => r.taskStatus === 'IN_PROGRESS').length
}));

async function loadData() {
  isLoading.value = true;
  try {
    const [worksR, tasksR, typesR, posR] = await Promise.all([
      mesApi.listWorkOrders({}), mesApi.listTasks({}), mesApi.listWorkTypes({}), mesApi.listPositions({})
    ]);
    works.value = worksR;
    tasksR.forEach(t => taskNames.value[t.id] = t.name);
    (typesR as MESWorkType[]).forEach(wt => workTypeNames.value[wt.id] = wt.name);
    (posR as any[]).forEach(p => positionNames.value[p.id] = p.name);
  } finally { isLoading.value = false; }
}

async function runAction(row: any, action: string) {
  try {
    await mesApi.updateWorkOrderTaskStatus(row.workId, row.taskId, { action: action as WorkOrderTaskAction });
    await loadData();
  } catch (err: any) { alert(err.message); }
}

// --- DETALLE ---
const detailRow = ref<any | null>(null);
const detailNotes = ref('');
function openDetail(row: any) { detailRow.value = row; detailNotes.value = ''; }
function closeDetail() { detailRow.value = null; }
async function dialogAction(action: WorkOrderTaskAction) {
  if (!detailRow.value) return;
  try {
    await mesApi.updateWorkOrderTaskStatus(detailRow.value.workId, detailRow.value.taskId, { action, notes: detailNotes.value });
    closeDetail(); await loadData();
  } catch (err: any) { alert(err.message); }
}

function taskStatusLabel(s: string) { const map = { PENDING: 'Pte.', IN_PROGRESS: 'En curso', BLOCKED: 'Bloq.', COMPLETED: 'Fin' }; return map[s] || s; }

onMounted(loadData);
</script>

<style scoped>
.status-ribbon { 
  display: flex; gap: 1.5rem; align-items: center; background: #1e293b; 
  padding: 1rem 2rem; border-radius: 16px; margin-bottom: 1.5rem;
}
.status-item { display: flex; flex-direction: column; }
.status-item label { font-size: 0.7rem; color: #94a3b8; text-transform: uppercase; font-weight: 700; }
.status-item strong { font-size: 2rem; line-height: 1; color: white; }
.status-item.active strong { color: var(--color-primary); }

.terminal-search-wrap { flex: 1; position: relative; margin-left: 2rem; }
.terminal-search-wrap .material-symbols-outlined { position: absolute; left: 1rem; top: 50%; transform: translateY(-50%); font-size: 2rem; color: #64748b; }
.terminal-input { 
  width: 100%; background: #0f172a; border: 2px solid #334155; border-radius: 12px;
  padding: 1rem 1rem 1.25rem 4rem; font-size: 1.25rem; color: white;
}

.task-list-container { flex: 1; overflow-y: auto; display: flex; flex-direction: column; gap: 1rem; }

.task-row-card { 
  background: #1e293b; border-radius: 16px; padding: 1.5rem; display: flex;
  align-items: center; justify-content: space-between; border: 2px solid transparent;
  transition: 0.2s; cursor: pointer;
}
.task-row-card.status-in_progress { border-color: var(--color-primary); box-shadow: 0 0 20px rgba(230, 184, 0, 0.1); }

.work-ref { display: flex; align-items: center; gap: 1rem; margin-bottom: 0.5rem; }
.terminal-code { background: var(--color-primary); color: #000; padding: 0.4rem 0.8rem; border-radius: 8px; font-weight: 900; font-size: 1.1rem; }
.work-ref strong { font-size: 1.4rem; color: white; }

.task-info { display: flex; flex-direction: column; }
.task-name { font-size: 1.2rem; font-weight: 700; color: #e2e8f0; }
.work-type { font-size: 0.85rem; color: #94a3b8; font-weight: 600; }

.card-meta { display: flex; flex-direction: column; align-items: flex-end; gap: 0.5rem; }
.terminal-pos { display: flex; align-items: center; gap: 0.4rem; color: #3b82f6; font-weight: 700; font-size: 1.1rem; }

.btn-action { 
  padding: 1rem 2rem; border-radius: 12px; border: none; font-weight: 800; font-size: 1rem;
  display: flex; align-items: center; gap: 0.75rem; cursor: pointer;
}
.btn-action.start { background: var(--color-primary); color: black; }
.btn-action.complete { background: #16a34a; color: white; }

.terminal-status-pill { font-size: 0.75rem; font-weight: 800; padding: 0.25rem 0.75rem; border-radius: 20px; text-transform: uppercase; }
.terminal-status-pill.pending { background: rgba(230, 184, 0, 0.1); color: var(--color-primary); }
.terminal-status-pill.in_progress { background: rgba(59, 130, 246, 0.1); color: #3b82f6; }

.terminal-giant-actions { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; margin-top: 2rem; }
.btn-giant { height: 120px; border-radius: 20px; border: none; font-size: 1.5rem; font-weight: 900; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 0.5rem; cursor: pointer; }
.btn-giant.start { background: var(--color-primary); color: black; }
.btn-giant.complete { background: #16a34a; color: white; }
.btn-giant.block { background: #dc2626; color: white; grid-column: span 2; height: 80px; font-size: 1.1rem; }

.terminal-textarea-wrap textarea { width: 100%; height: 150px; background: #f1f5f9; border: 2px solid #cbd5e1; border-radius: 12px; padding: 1rem; font-size: 1.25rem; }
</style>
