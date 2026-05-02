<template>
  <BaseDashboardPage :is-loading="isLoading" class="no-print">
    <template #header>
      <PageHeader title="Monitor de Producción (MES)">
        <template #icon><Factory :size="28" /></template>
        <template #actions>
          <button @click="loadDashboard" class="btn btn-outline btn-sm" :disabled="isLoading">
            <RefreshCw :size="16" :class="{ 'spin': isLoading }" />
            Actualizar
          </button>
        </template>
      </PageHeader>
    </template>

    <div class="module-dashboard-content">
      <!-- 1. KPIs de Resumen -->
      <section class="stats-grid">
        <div class="stat-card clickable" @click="router.push('/mes/work-orders')">
          <div class="stat-icon blue"><Factory :size="22" /></div>
          <div class="stat-info">
            <span class="stat-label">Total Trabajos</span>
            <span class="stat-value">{{ stats?.total ?? 0 }}</span>
          </div>
        </div>
        <div class="stat-card clickable" @click="router.push('/mes/work-orders?status=OVERDUE')">
          <div class="stat-icon red"><History :size="22" /></div>
          <div class="stat-info">
            <span class="stat-label">Vencidos</span>
            <span class="stat-value text-danger">{{ stats?.overdue ?? 0 }}</span>
          </div>
        </div>
        <div class="stat-card clickable" @click="router.push('/mes/work-orders?status=IN_PROGRESS')">
          <div class="stat-icon yellow"><AlertCircle :size="22" /></div>
          <div class="stat-info">
            <span class="stat-label">En Producción</span>
            <span class="stat-value text-warning">{{ inProgressOrders.length }}</span>
          </div>
        </div>
        <div class="stat-card clickable" @click="router.push('/mes/terminal')">
          <div class="stat-icon purple"><ClipboardList :size="22" /></div>
          <div class="stat-info">
            <span class="stat-label">Tareas Taller</span>
            <span class="stat-value">{{ pendingTasksCount }}</span>
          </div>
        </div>
      </section>

      <!-- 2. Accesos a Listados -->
      <section class="listings-grid">
        <RouterLink to="/mes/work-orders" class="listing-link">
          <List :size="20" />
          <span>Órdenes de Trabajo</span>
        </RouterLink>
        <RouterLink to="/mes/tasks" class="listing-link">
          <ClipboardCheck :size="20" />
          <span>Tareas de Taller</span>
        </RouterLink>
        <RouterLink to="/mes/positions" class="listing-link">
          <MapPin :size="20" />
          <span>Mapa de Posiciones</span>
        </RouterLink>
        <RouterLink to="/mes/terminal" class="listing-link highlight-subtle">
          <Tablet :size="20" />
          <span>Terminal Taller</span>
        </RouterLink>
      </section>

      <!-- 3. Actividad (Órdenes en Cola de Producción) -->
      <section class="dashboard-section">
        <div class="section-header">
          <Timer :size="20" class="text-info" />
          <h2>Cola de Producción (Sin Lanzar)</h2>
          <span class="header-tag">{{ pendingOrders.length }}</span>
        </div>
        <div class="table-wrapper">
          <table class="data-table">
            <thead>
              <tr>
                <th>Nº Orden</th>
                <th>Descripción / Trabajo</th>
                <th>Pedido Origen</th>
                <th>Prioridad</th>
                <th class="align-right">Acción</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="wo in pendingOrders" :key="wo.id" class="row-hover">
                <td><code class="code-badge">{{ wo.work_number }}</code></td>
                <td><strong>{{ wo.work_name }}</strong></td>
                <td>#{{ wo.sales_order_number || 'Interno' }}</td>
                <td><span :class="['priority-pill', `prio-${wo.priority}`]">{{ wo.priority }}</span></td>
                <td class="align-right">
                  <div class="actions-cell">
                    <button v-if="wo.work_setup_id" class="btn btn-secondary btn-sm" @click="launchToWorkshop(wo)" title="Enviar a la Tablet del taller">
                      <Rocket :size="16" />
                      Lanzar
                    </button>
                    <button v-else class="btn btn-primary btn-sm" @click="configureOrder(wo)">Configurar</button>
                    
                    <button v-if="wo.work_setup_id" class="btn btn-ghost btn-icon btn-sm" @click="configureOrder(wo)" title="Cambiar configuración técnica">
                      <Settings :size="16" />
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="pendingOrders.length === 0">
                <td colspan="5" class="p-4 text-center text-muted italic">No hay órdenes esperando en la cola de producción.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <!-- 4. Actividad en Curso (Trabajos en el Taller) -->
      <section class="dashboard-section mt-6">
        <div class="section-header">
          <PlayCircle :size="20" class="text-success" />
          <h2>Trabajos en Curso (Taller)</h2>
          <span class="header-tag success">{{ inProgressOrders.length }}</span>
        </div>
        <div class="table-wrapper">
          <table class="data-table">
            <thead>
              <tr>
                <th>Nº Orden</th>
                <th>Trabajo / Descripción</th>
                <th>Pedido</th>
                <th>Progreso Tareas</th>
                <th class="align-right">Acción</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="wo in inProgressOrders" :key="wo.id" class="row-hover">
                <td><code class="code-badge">{{ wo.work_number }}</code></td>
                <td><strong>{{ wo.work_name }}</strong></td>
                <td>#{{ wo.sales_order_number || 'Interno' }}</td>
                <td>
                  <div class="progress-info">
                    <span class="text-xs font-bold">{{ getCompletedTasksCount(wo) }} / {{ getTotalTasksCount(wo) }}</span>
                    <div class="progress-bar-mini">
                      <div class="progress-fill" :style="{ width: getProgressPercentage(wo) + '%' }"></div>
                    </div>
                  </div>
                </td>
                <td class="align-right">
                  <div class="actions-cell">
                    <button class="btn btn-outline btn-sm" @click="router.push('/mes/terminal')">
                      <Tablet :size="16" />
                      Ver Terminal
                    </button>
                    <button class="btn btn-ghost btn-icon btn-sm" @click="router.push(`/mes/work-orders/${wo.id}`)" title="Ver detalle completo">
                      <Eye :size="16" />
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="inProgressOrders.length === 0">
                <td colspan="5" class="p-4 text-center text-muted italic">No hay trabajos activos en el taller en este momento.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </div>

    <!-- Diálogo de Configuración -->
    <WorkSetupSelectorDialog
      :show="showSetupDialog"
      :work-order="selectedWorkOrder"
      @close="showSetupDialog = false"
      @assigned="handleSetupAssigned"
    />

    <template #sidebar>
      <section class="sidebar-section">
        <div class="section-header">
          <Zap :size="20" />
          <h2>Mantenimiento</h2>
        </div>
        <div class="quick-actions-list">
          <RouterLink to="/mes/work-setups" class="admin-card clickable">
            <Cpu :size="20" class="text-primary" />
            <div class="admin-card-info">
              <strong>Setups Técnicos</strong>
              <p>Configuraciones base</p>
            </div>
          </RouterLink>
          <RouterLink to="/mes/tasks" class="admin-card clickable mt-2">
            <Sliders :size="20" class="text-secondary" />
            <div class="admin-card-info">
              <strong>Maestro de Tareas</strong>
              <p>Definición de procesos</p>
            </div>
          </RouterLink>
        </div>
      </section>
    </template>
  </BaseDashboardPage>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { 
  Factory, 
  RefreshCw, 
  History, 
  AlertCircle, 
  ClipboardList, 
  List, 
  ClipboardCheck, 
  MapPin, 
  Tablet, 
  Timer, 
  Rocket, 
  Settings, 
  PlayCircle, 
  Eye, 
  Zap, 
  Cpu, 
  Sliders 
} from 'lucide-vue-next'
import PageHeader from '@/components/layout/PageHeader.vue'
import BaseDashboardPage from '@/components/shared/BaseDashboardPage.vue'
import WorkSetupSelectorDialog from '@/components/mes/WorkSetupSelectorDialog.vue'
import { mesApi } from '@/services/mesApi'
import { partyApi } from '@/services/partyApi'
import type { WorkOrder, WorkOrderDashboardStats } from '@/types/mes'

const router = useRouter()
const isLoading = ref(true)
const stats = ref<WorkOrderDashboardStats | null>(null)
const pendingOrders = ref<WorkOrder[]>([])
const inProgressOrders = ref<WorkOrder[]>([])
const pendingSalesWork = ref<any[]>([])

// Dialog state
const showSetupDialog = ref(false)
const selectedWorkOrder = ref<any | null>(null)

const pendingTasksCount = computed(() => {
  let count = 0
  for (const wo of inProgressOrders.value) { 
    for (const line of wo.lines || []) { 
      for (const task of line.tasks || []) { 
        if (task.status === 'PENDING') count++ 
      } 
    } 
  }
  return count
})

async function loadDashboard() {
  isLoading.value = true
  try {
    const [statsResult, pendingResult, inProgressResult, pendingSetupsResult] = await Promise.allSettled([
      mesApi.getWorkOrderDashboardStats(),
      mesApi.listWorkOrders({ status: 'PENDING' }),
      mesApi.listWorkOrders({ status: 'IN_PROGRESS' }),
      mesApi.listPendingWorkSetups()
    ])

    if (statsResult.status === 'fulfilled') {
      stats.value = statsResult.value
    } else {
      console.error('[MES Dashboard] Error loading stats:', statsResult.reason)
      stats.value = null
    }

    if (pendingResult.status === 'fulfilled') {
      pendingOrders.value = pendingResult.value
    } else {
      console.error('[MES Dashboard] Error loading pending work orders:', pendingResult.reason)
      pendingOrders.value = []
    }

    if (inProgressResult.status === 'fulfilled') {
      inProgressOrders.value = inProgressResult.value
    } else {
      console.error('[MES Dashboard] Error loading in-progress work orders:', inProgressResult.reason)
      inProgressOrders.value = []
    }

    const pendingSetups = pendingSetupsResult.status === 'fulfilled' ? pendingSetupsResult.value : []
    if (pendingSetupsResult.status === 'rejected') {
      console.error('[MES Dashboard] Error loading pending work setups:', pendingSetupsResult.reason)
    }
    
    // Enriquecer solicitudes de ventas con nombres de clientes
    const enrichedSW = await Promise.all(pendingSetups.map(async (item: any) => {
      if (item.party_id) {
        try {
          const party = await partyApi.getParty(item.party_id)
          return { ...item, party_name: party?.name || party?.displayName }
        } catch (e) {
          return item
        }
      }
      return item
    }))

    pendingSalesWork.value = enrichedSW
  } catch (err) {
    console.error('Error loading dashboard:', err)
  }
  finally { isLoading.value = false }
}

function configureOrder(wo: WorkOrder) {
  selectedWorkOrder.value = wo
  showSetupDialog.value = true
}

async function launchToWorkshop(wo: WorkOrder) {
  isLoading.value = true
  try {
    await mesApi.updateWorkOrder(wo.id, { status: 'IN_PROGRESS' })
    await loadDashboard()
  } catch (err) {
    console.error('Error launching order:', err)
  } finally {
    isLoading.value = false
  }
}

function configurePending(setup: any) {
  // Mapeamos el objeto de solicitud pendiente al formato que espera el diálogo
  selectedWorkOrder.value = {
    id: setup.id,
    work_number: setup.order_number,
    work_name: setup.description,
    party_id: setup.party_id,
    party_name: setup.party_name,
    order_work_setup_id: setup.id // Referencia crucial para vincular al crear la orden
  }
  showSetupDialog.value = true
}

async function handleSetupAssigned() {
  await loadDashboard()
}

// Helpers de progreso
function getTotalTasksCount(wo: WorkOrder) {
  let total = 0
  wo.lines?.forEach(l => { total += l.tasks?.length || 0 })
  return total
}

function getCompletedTasksCount(wo: WorkOrder) {
  let completed = 0
  wo.lines?.forEach(l => { 
    l.tasks?.forEach(t => { if (t.status === 'COMPLETED') completed++ })
  })
  return completed
}

function getProgressPercentage(wo: WorkOrder) {
  const total = getTotalTasksCount(wo)
  if (total === 0) return 0
  return Math.round((getCompletedTasksCount(wo) / total) * 100)
}

function formatDate(d: any) {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('es-ES', { day: 'numeric', month: 'short' })
}

onMounted(loadDashboard)
</script>

<style scoped>
.module-dashboard-content { display: flex; flex-direction: column; gap: 1.5rem; }

.stats-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 0.75rem; }
.stat-card { background: white; padding: 0.75rem 1rem; border-radius: 10px; border: 1px solid var(--color-border); display: flex; align-items: center; gap: 0.75rem; position: relative; transition: 0.2s; cursor: pointer; }
.stat-card:hover { transform: translateY(-2px); box-shadow: var(--box-shadow-md); border-color: var(--color-primary); }
.stat-icon { width: 40px; height: 40px; border-radius: 8px; display: flex; align-items: center; justify-content: center; }
.stat-icon :deep(svg) { width: 22px; height: 22px; }
.stat-icon.blue { background: rgba(59, 130, 246, 0.1); color: #3b82f6; }
.stat-icon.red { background: rgba(239, 68, 68, 0.1); color: #ef4444; }
.stat-icon.yellow { background: rgba(230, 184, 0, 0.1); color: #E6B800; }
.stat-icon.purple { background: rgba(168, 85, 247, 0.1); color: #a855f7; }
.stat-info { display: flex; flex-direction: column; gap: 0.25rem; }
.stat-label { font-size: 0.65rem; color: var(--color-text-secondary); font-weight: 600; text-transform: uppercase; }
.stat-value { font-size: 1.25rem; font-weight: 700; }

.listings-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 0.75rem; }
.listing-link { display: flex; align-items: center; gap: 0.75rem; padding: 0.75rem 1rem; background: white; border: 1px solid var(--color-border); border-radius: 8px; text-decoration: none; color: var(--color-text-primary); font-size: 0.85rem; font-weight: 600; transition: 0.2s; }
.listing-link:hover { background: var(--color-background); border-color: var(--color-secondary); color: var(--color-secondary); }
.listing-link :deep(svg) { color: var(--color-secondary); }
.listing-link.highlight-subtle { border-left: 3px solid var(--color-primary); }

.dashboard-section { background: white; padding: 0.75rem 1rem; border-radius: 10px; border: 1px solid var(--color-border); box-shadow: var(--box-shadow-sm); }
.section-header { display: flex; align-items: center; gap: 0.5rem; margin-bottom: 0.75rem; padding-bottom: 0.5rem; border-bottom: 1px solid var(--color-background); }
.section-header h2 { font-size: 0.85rem; font-weight: 700; text-transform: uppercase; margin: 0; flex: 1; }
.header-tag { font-size: 0.6rem; font-weight: 800; padding: 0.1rem 0.5rem; background: var(--color-background); color: var(--color-secondary); border-radius: 20px; }

.quick-actions-list { display: flex; flex-direction: column; gap: 0.75rem; }
.admin-card { display: flex; align-items: center; gap: 0.75rem; padding: 0.75rem; background: var(--color-background); border-radius: 8px; border: 1px solid transparent; text-decoration: none; color: var(--color-text-primary); transition: 0.2s; }
.admin-card:hover { background: white; border-color: var(--color-primary); transform: translateX(4px); box-shadow: var(--box-shadow-sm); }
.admin-card-info strong { font-size: 0.8rem; display: block; }
.admin-card-info p { font-size: 0.65rem; color: var(--color-text-secondary); margin: 0; }

.code-badge { background: var(--color-background); padding: 0.15rem 0.35rem; border-radius: 4px; font-family: var(--font-family-mono); font-size: 0.75rem; }
.priority-pill { font-size: 0.65rem; font-weight: 800; text-transform: uppercase; padding: 0.15rem 0.4rem; border-radius: 4px; }
.prio-HIGH, .prio-URGENT { background: rgba(220, 38, 38, 0.1); color: #dc2626; }
.align-right { text-align: right; }
.mt-6 { margin-top: 1.5rem; }

.actions-cell { display: flex; justify-content: flex-end; gap: 0.5rem; }
.btn-icon { padding: 0.25rem; min-width: auto; }
.btn-ghost { background: transparent; border-color: transparent; color: var(--color-text-secondary); }
.btn-ghost:hover { color: var(--color-primary); background: var(--color-background); }

.spin { animation: spin 1s linear infinite; }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
</style>