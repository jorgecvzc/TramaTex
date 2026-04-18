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

      <!-- Filtro de Especialidad (Tipo de Tarea) -->
      <div class="specialty-filter">
        <label class="filter-label">Especialidad</label>
        <select v-model="taskFilter" class="terminal-select">
          <option value="">TODAS LAS TAREAS</option>
          <option v-for="type in availableTaskTypes" :key="type.id" :value="type.id">
            {{ type.name.toUpperCase() }}
          </option>
        </select>
      </div>

      <div class="terminal-search-wrap">
        <span class="material-symbols-outlined">search</span>
        <input v-model="search" type="text" placeholder="Escanear código o buscar trabajo..." class="terminal-input" />
      </div>
    </div>

    <!-- LISTADO TÁCTIL (TABULAR ESTRICTO) -->
    <div class="task-list-container">
      <div v-if="filteredRows.length === 0" class="terminal-empty">
        <span class="material-symbols-outlined">checklist</span>
        <p>No hay tareas pendientes asignadas a esta estación.</p>
      </div>

      <div 
        v-for="row in filteredRows" 
        :key="row.taskId" 
        class="task-row-tabular"
        :class="['status-' + row.taskStatus.toLowerCase()]"
        @click="openDetail(row)"
      >
        <!-- COL 1: CÓDIGO -->
        <div class="col-code">
          <code class="terminal-code">{{ row.workNumber }}</code>
        </div>

        <!-- COL 2: TRABAJO Y CLIENTE -->
        <div class="col-work-client">
          <strong class="work-title">{{ row.workName }}</strong>
          <div class="client-info">
            <span class="material-symbols-outlined">person</span>
            <span>{{ row.partyName }}</span>
          </div>
        </div>

        <!-- COL 3: TIPO Y POSICIÓN -->
        <div class="col-context">
          <span class="work-type-badge">{{ row.workTypeName.toUpperCase() }}</span>
          <div v-if="row.positionName" class="pos-info-inline">
            <span class="material-symbols-outlined">location_on</span>
            <strong>{{ row.positionName }}</strong>
          </div>
        </div>

        <!-- COL 4: TAREA Y ACCIÓN -->
        <div class="col-task-action">
          <div class="task-desc">
            <label>TAREA ACTUAL</label>
            <span class="task-name">{{ row.taskName }}</span>
          </div>

          <div class="action-wrap" @click.stop>
            <!-- Adjuntos -->
            <button 
              v-if="row.hasAttachments" 
              class="btn-terminal-icon attachment"
              @click="viewAttachments(row)"
            >
              <span class="material-symbols-outlined">description</span>
            </button>

            <!-- Botón de Acción -->
            <button
              v-if="['PENDING', 'BLOCKED'].includes(row.taskStatus)"
              class="btn-terminal-action start"
              @click="runAction(row, 'START')"
            >
              <span class="material-symbols-outlined">play_arrow</span>
              INICIAR
            </button>
            <button
              v-if="row.taskStatus === 'IN_PROGRESS'"
              class="btn-terminal-action complete"
              @click="runAction(row, 'COMPLETE')"
            >
              <span class="material-symbols-outlined">check_circle</span>
              FINALIZAR
            </button>
          </div>
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
import { partyApi } from '@/services/partyApi';
import type { WorkOrder, WorkOrderTaskAction, MESWorkType } from '@/types/mes';

const router = useRouter();
const isLoading = ref(true);
const search = ref('');
const taskFilter = ref('');
const works = ref<WorkOrder[]>([]);
const taskNames = ref<Record<string, string>>({});
const workTypeNames = ref<Record<string, string>>({});
const positionNames = ref<Record<string, string>>({});
const partyNames = ref<Record<string, string>>({});
const setupToPartyMap = ref<Record<string, string>>({}); // Mapeo setupId -> partyId

const goBack = () => router.push('/mes/dashboard');

// --- DATOS ---
const rows = computed(() => {
  const res: any[] = [];
  for (const w of works.value) {
    // Determinamos el party_id real (directo o via setup)
    const workSetupId = w.work_setup_id || (w as any).workSetupID || '';
    const effectivePartyId = w.party_id || (w as any).partyId || setupToPartyMap.value[workSetupId] || '';
    
    for (const l of w.lines || []) {
      for (const t of l.tasks || []) {
        res.push({
          workId: w.id, workNumber: w.work_number, workName: w.work_name,
          partyId: effectivePartyId,
          partyName: partyNames.value[effectivePartyId] || w.party_name || 'Interno',
          workTypeId: l.work_type_id,
          workTypeName: workTypeNames.value[l.work_type_id] || l.work_type_id,
          taskId: t.id, taskTypeId: t.task_id,
          taskName: taskNames.value[t.task_id] || t.task_id,
          taskStatus: t.status, sequence: t.sequence,
          positionName: positionNames.value[l.position_id] || '',
        });
      }
    }
  }
  return res;
});

const availableTaskTypes = computed(() => {
  const typesMap: Record<string, string> = {};
  rows.value.forEach(r => {
    if (r.taskTypeId) typesMap[r.taskTypeId] = r.taskName;
  });
  return Object.entries(typesMap).map(([id, name]) => ({ id, name }));
});

const filteredRows = computed(() => {
  const term = search.value.toLowerCase();
  const filterId = taskFilter.value;
  
  return rows.value.filter(r => {
    const matchesStatus = ['PENDING', 'IN_PROGRESS', 'BLOCKED'].includes(r.taskStatus);
    const matchesTerm = !term || r.workNumber.toLowerCase().includes(term) || r.workName.toLowerCase().includes(term);
    const matchesType = !filterId || r.taskTypeId === filterId;
    return matchesStatus && matchesTerm && matchesType;
  });
});

const counts = computed(() => ({
  pending: filteredRows.value.filter(r => r.taskStatus === 'PENDING').length,
  inProgress: filteredRows.value.filter(r => r.taskStatus === 'IN_PROGRESS').length
}));

async function loadData() {
  isLoading.value = true;
  console.log('[Terminal] Iniciando carga de datos...');
  try {
    const [worksR, tasksR, typesR, posR, setupsR] = await Promise.all([
      mesApi.listWorkOrders({}), 
      mesApi.listTasks({}), 
      mesApi.listWorkTypes({}), 
      mesApi.listPositions({}),
      mesApi.listWorkSetups({})
    ]);

    console.log(`[Terminal] Cargadas ${worksR.length} órdenes y ${setupsR.length} setups.`);
    
    works.value = worksR;
    tasksR.forEach(t => taskNames.value[t.id] = t.name);
    (typesR as MESWorkType[]).forEach(wt => workTypeNames.value[wt.id] = wt.name);
    (posR as any[]).forEach(p => positionNames.value[p.id] = p.name);

    // Mapeamos setups a clientes para el fallback
    setupsR.forEach(s => {
      const pId = (s as any).party_id || (s as any).partyId;
      if (s.id && pId) setupToPartyMap.value[s.id] = pId;
    });

    // Hidratar nombres de clientes de forma robusta
    const partyIdsToFetch = new Set<string>();
    
    worksR.forEach(w => {
      const workSetupId = w.work_setup_id || (w as any).workSetupID || '';
      const pId = w.party_id || (w as any).partyId || setupToPartyMap.value[workSetupId];
      if (pId) {
        partyIdsToFetch.add(pId);
      } else {
        console.warn(`[Terminal] La orden ${w.work_number} no tiene party_id ni setup vinculado.`);
      }
    });

    const uniquePartyIds = Array.from(partyIdsToFetch);
    console.log(`[Terminal] Hidratando ${uniquePartyIds.length} clientes únicos...`);

    if (uniquePartyIds.length > 0) {
      await Promise.all(uniquePartyIds.map(async (id) => {
        try {
          const p = await partyApi.getParty(id);
          if (p) {
            partyNames.value[id] = p.name || p.displayName || p.businessName || 'Cliente sin nombre';
            console.log(`[Terminal] Cliente hidratado: ${id} -> ${partyNames.value[id]}`);
          }
        } catch (e) { 
          console.error(`[Terminal] Error cargando cliente ${id}:`, e); 
        }
      }));
    }

  } catch (err) {
    console.error('[Terminal] Error crítico en loadData:', err);
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
  display: flex; gap: 2rem; align-items: flex-end; background: #1e293b; 
  padding: 1.25rem 2rem; border-radius: 16px; margin-bottom: 1.5rem;
}
.status-item { display: flex; flex-direction: column; gap: 0.25rem; min-width: 100px; }
.status-item label { font-size: 0.7rem; color: #94a3b8; text-transform: uppercase; font-weight: 700; margin-left: 0.25rem; }
.status-item strong { font-size: 2.25rem; line-height: 1; color: white; }
.status-item.active strong { color: var(--color-primary); }

.specialty-filter { display: flex; flex-direction: column; gap: 0.5rem; min-width: 250px; }
.filter-label { font-size: 0.7rem; color: #94a3b8; text-transform: uppercase; font-weight: 700; margin-left: 0.5rem; }
.terminal-select {
  background: #0f172a; border: 2px solid #334155; color: white;
  padding: 0 1rem; border-radius: 12px; font-weight: 700; font-size: 1rem;
  cursor: pointer; outline: none; transition: 0.2s;
  height: 54px; 
}
.terminal-select:focus { border-color: var(--color-primary); }

.terminal-search-wrap { flex: 1; position: relative; }
.terminal-search-wrap .material-symbols-outlined { position: absolute; left: 1rem; top: 50%; transform: translateY(-50%); font-size: 2rem; color: #64748b; }
.terminal-input { 
  width: 100%; background: #0f172a; border: 2px solid #334155; border-radius: 12px;
  padding: 0 1rem 0 4rem; font-size: 1.25rem; color: white;
  height: 54px;
}

.task-list-container { flex: 1; overflow-y: auto; display: flex; flex-direction: column; gap: 0.75rem; }

.task-row-tabular { 
  background: #1e293b; border-radius: 12px; padding: 1rem 1.5rem; 
  display: grid; grid-template-columns: 120px 2fr 1fr 2.5fr; 
  align-items: center; gap: 3rem; border: 2px solid transparent;
  transition: 0.2s; cursor: pointer;
}
.task-row-tabular.status-in_progress { border-color: var(--color-primary); box-shadow: 0 0 15px rgba(230, 184, 0, 0.1); }

/* COL 1: CÓDIGO */
.col-code { display: flex; justify-content: center; }
.terminal-code { background: var(--color-primary); color: #000; padding: 0.5rem 0.75rem; border-radius: 8px; font-weight: 900; font-size: 1.1rem; }

/* COL 2: TRABAJO Y CLIENTE */
.col-work-client { display: flex; flex-direction: column; gap: 0.25rem; min-width: 0; }
.work-title { font-size: 1.25rem; color: white; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.client-info { display: flex; align-items: center; gap: 0.4rem; color: #94a3b8; font-size: 0.9rem; font-weight: 600; }
.client-info .material-symbols-outlined { font-size: 1.1rem; opacity: 0.7; }

/* COL 3: TIPO Y POSICIÓN */
.col-context { display: flex; flex-direction: column; gap: 0.5rem; }
.work-type-badge {
  align-self: flex-start; background: #334155; color: white;
  padding: 0.2rem 0.6rem; border-radius: 6px; font-size: 0.7rem; font-weight: 900;
  border: 1px solid #475569; letter-spacing: 0.05em;
}
.pos-info-inline { display: flex; align-items: center; gap: 0.4rem; color: #3b82f6; }
.pos-info-inline .material-symbols-outlined { font-size: 1.25rem; }
.pos-info-inline strong { font-size: 1.1rem; font-weight: 800; text-transform: uppercase; }

/* COL 4: TAREA Y ACCIÓN */
.col-task-action { display: flex; align-items: center; gap: 1.5rem; justify-content: space-between; }
.task-desc { display: flex; flex-direction: column; gap: 0.1rem; flex: 1; min-width: 0; }
.task-desc label { font-size: 0.6rem; color: #64748b; font-weight: 800; text-transform: uppercase; }
.task-name { font-size: 1.1rem; font-weight: 700; color: #e2e8f0; line-height: 1.2; }

.action-wrap { display: flex; align-items: center; gap: 1rem; }
.btn-terminal-icon {
  background: rgba(255, 255, 255, 0.05); border: 1px solid #334155; color: #94a3b8;
  width: 48px; height: 48px; border-radius: 10px; display: flex; align-items: center; justify-content: center;
  cursor: pointer; transition: 0.2s;
}
.btn-terminal-icon:hover { background: #334155; color: white; }

.btn-terminal-action { 
  padding: 0 1.5rem; height: 50px; border-radius: 10px; border: none; font-weight: 900; font-size: 1rem;
  display: flex; align-items: center; gap: 0.5rem; cursor: pointer; min-width: 150px; justify-content: center;
}
.btn-terminal-action.start { background: var(--color-primary); color: black; }
.btn-terminal-action.complete { background: #16a34a; color: white; }

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
