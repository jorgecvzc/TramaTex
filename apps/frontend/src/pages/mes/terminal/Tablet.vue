<template>
  <Navbar class="no-print" />
  
  <BaseDashboardPage :is-loading="isLoading" class="no-print">
    <!-- CAPA 1: IDENTIDAD -->
    <template #header>
      <PageHeader 
        title="Terminal de Taller" 
        :breadcrumbs="[{ label: 'MES', to: '/mes/dashboard' }, { label: 'Terminal Operativa' }]"
      >
        <template #icon>
          <span class="material-symbols-outlined">tablet_mac</span>
        </template>
        <template #actions>
          <button class="btn btn-outline" @click="loadData" :disabled="isLoading">
            <span class="material-symbols-outlined" :class="{ 'spin': isLoading }">refresh</span>
            <span>Actualizar</span>
          </button>
        </template>
      </PageHeader>
    </template>

    <!-- CAPA 3: TRABAJO (Listado de Tareas en Cola) -->
    <div class="terminal-main-area">
      <!-- Filtros Operativos -->
      <section class="card filters-panel mb-6">
        <div class="filter-row">
          <div class="filter-group grow">
            <label>Búsqueda Rápida</label>
            <div class="input-with-icon">
              <span class="material-symbols-outlined icon-start">search</span>
              <input v-model="search" type="text" class="form-input" placeholder="Nº Trabajo, Producto o Cliente..." />
            </div>
          </div>
          <div class="filter-group">
            <label>Estado Tarea</label>
            <select v-model="taskStatusFilter" class="form-input">
              <option value="ACTIVE">Tareas Activas</option>
              <option value="PENDING">Pendientes</option>
              <option value="IN_PROGRESS">En progreso</option>
              <option value="BLOCKED">Bloqueadas</option>
              <option value="COMPLETED">Completadas</option>
            </select>
          </div>
        </div>
      </section>

      <!-- Tabla de Ejecución -->
      <section class="card table-card overflow-hidden">
        <div class="table-wrapper">
          <table class="data-table">
            <thead>
              <tr>
                <th>Trabajo / Referencia</th>
                <th>Proceso</th>
                <th>Tarea Actual</th>
                <th>Posición</th>
                <th>Estado</th>
                <th class="align-right">Acciones Rápidas</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="row in filteredRows" :key="row.taskId" class="row-clickable" @click="openDetail(row)">
                <td>
                  <div class="work-ref-cell">
                    <code class="code-badge">{{ row.workNumber }}</code>
                    <small>{{ row.workName }}</small>
                  </div>
                </td>
                <td><span class="text-secondary fw-bold">{{ row.workTypeName }}</span></td>
                <td>
                  <div class="task-cell">
                    <strong>{{ row.taskName }}</strong>
                    <small class="text-muted">Secuencia {{ row.sequence }}</small>
                  </div>
                </td>
                <td>
                  <span v-if="row.positionName" class="position-tag">
                    <span class="material-symbols-outlined">location_on</span>
                    {{ row.positionName }}
                  </span>
                  <span v-else class="text-muted">—</span>
                </td>
                <td>
                  <span :class="['status-badge', `status-${getStatusClass(row.taskStatus)}`]">
                    {{ taskStatusLabel(row.taskStatus) }}
                  </span>
                </td>
                <td class="align-right" @click.stop>
                  <div class="action-buttons-compact">
                    <button
                      v-if="['PENDING', 'BLOCKED'].includes(row.taskStatus)"
                      class="btn btn-primary btn-sm"
                      @click="runAction(row, 'START')"
                    >
                      <span class="material-symbols-outlined">play_arrow</span>
                    </button>
                    <button
                      v-if="row.taskStatus === 'IN_PROGRESS'"
                      class="btn btn-success btn-sm"
                      @click="runAction(row, 'check_circle')"
                    >
                      <span class="material-symbols-outlined">done</span>
                    </button>
                    <button
                      v-if="row.taskStatus !== 'COMPLETED' && row.taskStatus !== 'SKIPPED'"
                      class="btn btn-outline btn-sm text-danger"
                      @click="runAction(row, 'BLOCK')"
                    >
                      <span class="material-symbols-outlined">block</span>
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="filteredRows.length === 0">
                <td colspan="6" class="empty-row-msg">No hay tareas pendientes en la cola actual.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </div>

    <!-- CAPA 2: CONTEXTO (Panel Informativo / Resumen) -->
    <template #sidebar>
      <div class="terminal-sidebar">
        <section class="sidebar-section">
          <div class="section-header">
            <span class="material-symbols-outlined">analytics</span>
            <h2>Estado de la Sesión</h2>
          </div>
          <div class="stats-mini-grid mt-4">
            <div class="stat-mini">
              <label>Pendientes</label>
              <strong>{{ counts.pending }}</strong>
            </div>
            <div class="stat-mini">
              <label>En Marcha</label>
              <strong class="text-primary">{{ counts.inProgress }}</strong>
            </div>
          </div>
        </section>

        <section class="sidebar-section mt-10">
          <div class="section-header">
            <span class="material-symbols-outlined">help_outline</span>
            <h2>Instrucciones</h2>
          </div>
          <ul class="terminal-guide mt-4">
            <li>Pinche en una fila para ver el <strong>detalle técnico</strong> y archivos adjuntos.</li>
            <li>Use los botones rápidos para cambiar el estado de la tarea.</li>
            <li>Las tareas bloqueadas requieren nota de motivo.</li>
          </ul>
        </section>

        <div class="help-notice mt-10">
          <div class="notice-header">
            <span class="material-symbols-outlined text-warning">warning</span>
            <h3>Atención</h3>
          </div>
          <p class="help-text">
            Recuerde que completar la última tarea de una orden marcará el trabajo como <strong>Terminado</strong> en el sistema central.
          </p>
        </div>
      </div>
    </template>
  </BaseDashboardPage>

  <!-- DIÁLOGO: DETALLE DE TAREA Y CONTROL -->
  <BaseDialog
    :show="!!detailRow"
    :title="detailRow ? `Tarea: ${detailRow.taskName}` : ''"
    icon="assignment"
    size="lg"
    hide-actions
    @close="closeDetail"
  >
    <div class="task-detail-dialog" v-if="detailRow">
      <!-- Header contextual del diálogo -->
      <div class="dialog-context-header mb-6">
        <div class="work-info">
          <code class="code-badge large">{{ detailRow.workNumber }}</code>
          <strong class="ml-2">{{ detailRow.workName }}</strong>
        </div>
        <div class="task-meta mt-2">
          <span class="meta-item"><span class="material-symbols-outlined">tune</span> {{ detailRow.workTypeName }}</span>
          <span class="meta-item ml-4"><span class="material-symbols-outlined">location_on</span> {{ detailRow.positionName || 'Sin posición' }}</span>
        </div>
      </div>

      <!-- Acciones de Control -->
      <div class="task-actions-panel card bg-light p-6 mb-8">
        <div class="form-group">
          <label>Notas de la Operación</label>
          <textarea v-model="detailNotes" class="form-textarea" rows="2" placeholder="Observaciones sobre la ejecución..."></textarea>
        </div>
        <div class="flex-row gap-4 mt-6">
          <button
            v-if="['PENDING', 'BLOCKED'].includes(detailRow.taskStatus)"
            class="btn btn-primary btn-lg grow"
            @click="dialogAction('START')"
          >
            <span class="material-symbols-outlined">play_arrow</span>
            <span>INICIAR TAREA</span>
          </button>
          <button
            v-if="detailRow.taskStatus === 'IN_PROGRESS'"
            class="btn btn-success btn-lg grow"
            @click="dialogAction('COMPLETE')"
          >
            <span class="material-symbols-outlined">check_circle</span>
            <span>COMPLETAR TAREA</span>
          </button>
          <button
            class="btn btn-danger btn-lg"
            @click="dialogAction('BLOCK')"
          >
            <span class="material-symbols-outlined">block</span>
            <span>BLOQUEAR</span>
          </button>
        </div>
      </div>

      <!-- Historial / Subtareas -->
      <FormSection title="Flujo de Trabajo" icon="account_tree">
        <div class="workflow-steps" v-if="detailLine">
          <div 
            v-for="task in detailLine.tasks" 
            :key="task.id" 
            :class="['step-item', { active: task.id === detailRow.taskId }]"
          >
            <div class="step-icon">
              <span class="material-symbols-outlined" v-if="task.status === 'COMPLETED'">check</span>
              <span v-else>{{ task.sequence }}</span>
            </div>
            <div class="step-content">
              <strong>{{ taskNames[task.task_id] }}</strong>
              <small :class="`text-${getStatusClass(task.status)}`">{{ taskStatusLabel(task.status) }}</small>
            </div>
            <div class="step-notes" v-if="task.notes">
              "{{ task.notes }}"
            </div>
          </div>
        </div>
      </FormSection>
    </div>
  </BaseDialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import Navbar from '@/components/layout/Navbar.vue';
import PageHeader from '@/components/layout/PageHeader.vue';
import BaseDashboardPage from '@/components/shared/BaseDashboardPage.vue';
import BaseDialog from '@/components/shared/BaseDialog.vue';
import FormSection from '@/components/shared/FormSection.vue';
import { mesApi } from '@/services/mesApi';
import { partyApi } from '@/services/partyApi';
import type { WorkOrder, WorkOrderLine, WorkOrderTaskAction, MESWorkType } from '@/types/mes';
import salesApi from '@/services/salesApi';

const route = useRoute();
const isLoading = ref(true);
const error = ref('');
const search = ref('');
const taskStatusFilter = ref('ACTIVE');
const taskFilter = ref('');

const works = ref<WorkOrder[]>([]);
const taskNames = ref<Record<string, string>>({});
const workTypeNames = ref<Record<string, string>>({});
const positionNames = ref<Record<string, string>>({});
const partiesCache = ref<Record<string, any>>({});

interface TaskRow {
  workId: string; workNumber: string; workName: string;
  workTypeName: string; lineId: string; taskId: string;
  taskTypeId: string; taskName: string; taskStatus: string;
  sequence: number; positionName: string; designFilePath: string;
}

const rows = computed<TaskRow[]>(() => {
  const result: TaskRow[] = [];
  for (const work of works.value) {
    for (const line of work.lines || []) {
      for (const task of line.tasks || []) {
        result.push({
          workId: work.id, workNumber: work.work_number, workName: work.work_name,
          workTypeName: workTypeNames.value[line.work_type_id] || line.work_type_id,
          lineId: line.id, taskId: task.id, taskTypeId: task.task_id,
          taskName: taskNames.value[task.task_id] || task.task_id,
          taskStatus: task.status, sequence: task.sequence,
          positionName: positionNames.value[line.position_id] || '',
          designFilePath: line.design_file_path || '',
        });
      }
    }
  }
  return result;
});

const filteredRows = computed(() => {
  const term = search.value.trim().toLowerCase();
  return rows.value.filter((row) => {
    const activeStatuses = ['PENDING', 'IN_PROGRESS', 'BLOCKED'];
    const matchesStatus = taskStatusFilter.value === 'ACTIVE' ? activeStatuses.includes(row.taskStatus) : row.taskStatus === taskStatusFilter.value;
    const matchesTerm = !term || row.workNumber.toLowerCase().includes(term) || row.workName.toLowerCase().includes(term) || row.taskName.toLowerCase().includes(term);
    return matchesStatus && matchesTerm;
  });
});

const counts = computed(() => ({
  pending: rows.value.filter(r => r.taskStatus === 'PENDING').length,
  inProgress: rows.value.filter(r => r.taskStatus === 'IN_PROGRESS').length
}));

async function loadData() {
  isLoading.value = true; error.value = '';
  try {
    const [worksResult, tasksResult, workTypesResult, positionsResult] = await Promise.all([
      mesApi.listWorkOrders({}), mesApi.listTasks({}), mesApi.listWorkTypes({}), mesApi.listPositions({})
    ]);
    works.value = worksResult;
    tasksResult.forEach(t => taskNames.value[t.id] = t.name);
    (workTypesResult as MESWorkType[]).forEach(wt => workTypeNames.value[wt.id] = wt.name);
    (positionsResult as any[]).forEach(pos => positionNames.value[pos.id] = pos.name);
  } catch (err: any) { error.value = 'Error al conectar con la terminal de taller.'; }
  finally { isLoading.value = false; }
}

async function runAction(row: TaskRow, action: WorkOrderTaskAction | string) {
  let notes: string | undefined;
  if (action === 'BLOCK') {
    const input = window.prompt('Motivo del bloqueo:', '');
    if (input === null) return;
    notes = input.trim();
  }
  
  // Mapeo interno para iconos rápidos
  const finalAction = action === 'check_circle' ? 'COMPLETE' : action as WorkOrderTaskAction;

  try {
    await mesApi.updateWorkOrderTaskStatus(row.workId, row.taskId, { action: finalAction, notes });
    await loadData();
  } catch (err: any) { alert(err.message); }
}

// --- Detalle ---
const detailRow = ref<TaskRow | null>(null);
const detailNotes = ref('');
const detailWork = computed(() => detailRow.value ? works.value.find(w => w.id === detailRow.value!.workId) : null);
const detailLine = computed(() => detailWork.value?.lines.find(l => l.id === detailRow.value?.lineId));

function openDetail(row: TaskRow) { detailRow.value = row; detailNotes.value = ''; }
function closeDetail() { detailRow.value = null; }

async function dialogAction(action: WorkOrderTaskAction) {
  if (!detailRow.value) return;
  try {
    await mesApi.updateWorkOrderTaskStatus(detailRow.value.workId, detailRow.value.taskId, { action, notes: detailNotes.value });
    closeDetail(); await loadData();
  } catch (err: any) { alert(err.message); }
}

function taskStatusLabel(s: string) { const map = { PENDING: 'Pendiente', IN_PROGRESS: 'En curso', BLOCKED: 'Bloqueada', COMPLETED: 'Terminada', SKIPPED: 'Omitida' }; return map[s] || s; }
function getStatusClass(s: string) { const map = { PENDING: 'warning', IN_PROGRESS: 'primary', BLOCKED: 'danger', COMPLETED: 'success' }; return map[s] || 'secondary'; }

onMounted(loadData);
</script>

<style scoped>
@import "@/design-system/_sections.css";

.filter-row { display: flex; gap: 1.5rem; align-items: flex-end; }
.grow { flex: 1; }

.work-ref-cell { display: flex; flex-direction: column; gap: 0.25rem; }
.code-badge { background: var(--color-background); padding: 0.2rem 0.5rem; border-radius: 4px; font-family: var(--font-family-mono); font-size: 0.8rem; font-weight: 700; color: var(--color-secondary); }
.code-badge.large { font-size: 1rem; padding: 0.4rem 0.75rem; }

.task-cell { display: flex; flex-direction: column; }
.position-tag { display: inline-flex; align-items: center; gap: 0.25rem; font-size: 0.85rem; font-weight: 600; color: #2563eb; background: rgba(37, 99, 235, 0.05); padding: 0.2rem 0.5rem; border-radius: 4px; }
.position-tag .material-symbols-outlined { font-size: 16px; }

.action-buttons-compact { display: flex; gap: 0.5rem; justify-content: flex-end; }

.stat-mini { display: flex; flex-direction: column; gap: 0.25rem; padding: 1rem; background: var(--color-background); border-radius: 8px; border: 1px solid var(--color-border); }
.stat-mini label { font-size: 0.65rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); }
.stat-mini strong { font-size: 1.5rem; }

.terminal-guide { padding-left: 1.25rem; font-size: 0.85rem; color: var(--color-text-secondary); line-height: 1.6; }
.terminal-guide li { margin-bottom: 0.75rem; }

/* Dialog Styles */
.meta-item { display: inline-flex; align-items: center; gap: 0.35rem; color: var(--color-text-secondary); font-size: 0.9rem; font-weight: 600; }
.meta-item .material-symbols-outlined { font-size: 18px; }

.task-actions-panel { border-left: 6px solid var(--color-primary); }
.flex-row { display: flex; }
.gap-4 { gap: 1rem; }

.workflow-steps { display: flex; flex-direction: column; gap: 1rem; }
.step-item { display: flex; align-items: center; gap: 1rem; padding: 0.75rem; border-radius: 8px; border: 1px solid var(--color-border); opacity: 0.6; }
.step-item.active { opacity: 1; border-color: var(--color-primary); background: rgba(230, 184, 0, 0.05); box-shadow: var(--box-shadow-sm); }
.step-icon { width: 28px; height: 28px; border-radius: 50%; background: var(--color-background); border: 2px solid var(--color-border); display: flex; align-items: center; justify-content: center; font-size: 0.75rem; font-weight: 800; }
.step-item.active .step-icon { border-color: var(--color-primary); color: var(--color-primary); }
.step-content { flex: 1; display: flex; flex-direction: column; line-height: 1.2; }
.step-notes { font-size: 0.75rem; font-style: italic; color: var(--color-text-secondary); max-width: 200px; text-align: right; }

.help-notice { padding: 1.25rem; background: rgba(59, 130, 246, 0.05); border-radius: 12px; border: 1px dashed rgba(59, 130, 246, 0.3); }
.notice-header { display: flex; align-items: center; gap: 0.5rem; color: #2563eb; font-size: 0.85rem; font-weight: 700; text-transform: uppercase; }
.help-text { font-size: 0.8rem; color: var(--color-text-secondary); margin-top: 0.5rem; line-height: 1.5; }

.input-with-icon { position: relative; display: flex; align-items: center; width: 100%; }
.icon-start { position: absolute; left: 0.85rem; font-size: 20px; color: var(--color-text-secondary); }
.input-with-icon .form-input { padding-left: 2.75rem; }

.form-input, .form-textarea { width: 100%; padding: 0.75rem 1rem; border-radius: 8px; border: 1px solid var(--color-border); font-family: inherit; }
</style>
