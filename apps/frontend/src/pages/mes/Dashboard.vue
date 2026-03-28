<template>
  <Navbar class="no-print" />
  
  <BaseDashboardPage :is-loading="isLoading" class="no-print">
    <!-- CAPA 1: IDENTIDAD -->
    <template #header>
      <div class="dashboard-header">
        <div class="header-title">
          <h1>Monitoreo de Producción (MES)</h1>
          <p class="subtitle">Supervisión en tiempo real del flujo de taller</p>
        </div>
        <div class="header-actions">
          <button @click="loadDashboard" class="btn btn-outline btn-sm" :disabled="isLoading">
            <span class="material-symbols-outlined" :class="{ 'spin': isLoading }">refresh</span>
            Actualizar datos
          </button>
          <button class="btn btn-primary btn-sm ml-2" @click="router.push('/mes/terminal')">
            <span class="material-symbols-outlined">tablet_mac</span>
            <span>Terminal de Taller</span>
          </button>
        </div>
      </div>
    </template>

    <!-- CAPA 3: TRABAJO (Área Principal) -->
    <div class="mes-main-area">
      <!-- Dashboard Stats (Siguiendo patrón del Dashboard Principal) -->
      <section class="stats-grid mb-8">
        <div class="stat-card clickable" @click="router.push('/mes/work-orders')">
          <div class="stat-icon blue"><span class="material-symbols-outlined">factory</span></div>
          <div class="stat-info">
            <span class="stat-label">Total Trabajos</span>
            <span class="stat-value">{{ stats?.total ?? 0 }}</span>
          </div>
          <div class="stat-link-arrow"><span class="material-symbols-outlined">arrow_forward</span></div>
        </div>

        <div class="stat-card clickable" @click="router.push('/mes/work-orders?status=OVERDUE')">
          <div class="stat-icon red"><span class="material-symbols-outlined">history</span></div>
          <div class="stat-info">
            <span class="stat-label">Vencidos</span>
            <span class="stat-value text-danger">{{ stats?.overdue ?? 0 }}</span>
          </div>
          <div class="stat-link-arrow"><span class="material-symbols-outlined">arrow_forward</span></div>
        </div>

        <div class="stat-card clickable" @click="router.push('/mes/work-orders?status=IN_PROGRESS')">
          <div class="stat-icon yellow"><span class="material-symbols-outlined">running_with_errors</span></div>
          <div class="stat-info">
            <span class="stat-label">En Producción</span>
            <span class="stat-value text-warning">{{ inProgressOrders.length }}</span>
          </div>
          <div class="stat-link-arrow"><span class="material-symbols-outlined">arrow_forward</span></div>
        </div>

        <div class="stat-card clickable" @click="router.push('/mes/terminal')">
          <div class="stat-icon purple"><span class="material-symbols-outlined">assignment</span></div>
          <div class="stat-info">
            <span class="stat-label">Tareas en Taller</span>
            <span class="stat-value">{{ pendingTasksCount }}</span>
          </div>
          <div class="stat-link-arrow"><span class="material-symbols-outlined">arrow_forward</span></div>
        </div>
      </section>

      <!-- PENDING Section -->
      <section class="dashboard-section mb-8">
        <div class="section-header">
          <span class="material-symbols-outlined text-info">pending_actions</span>
          <h2>Órdenes Pendientes de Configuración</h2>
          <span class="header-tag">{{ pendingOrders.length }}</span>
        </div>
        <div class="table-wrapper">
          <table class="data-table">
            <thead>
              <tr>
                <th>Número</th>
                <th>Descripción / Trabajo</th>
                <th>Pedido Origen</th>
                <th>Prioridad</th>
                <th class="align-right">Acciones</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="wo in pendingOrders" :key="wo.id" class="row-hover">
                <td><code class="code-badge">{{ wo.work_number }}</code></td>
                <td><strong>{{ wo.work_name }}</strong></td>
                <td>
                  <router-link v-if="wo.sales_order_id" :to="`/sales/orders/${wo.sales_order_id}`" class="link-primary">#{{ wo.sales_order_number }}</router-link>
                  <span v-else class="text-muted">—</span>
                </td>
                <td>
                  <span :class="['priority-pill', `prio-${wo.priority}`]">
                    {{ mesApi.getPriorityLabel(wo.priority) }}
                  </span>
                </td>
                <td class="align-right">
                  <div class="action-buttons">
                    <button v-if="!wo.work_setup_id" class="btn btn-primary btn-sm" @click="configureOrder(wo)">
                      <span class="material-symbols-outlined">settings</span> Configurar
                    </button>
                    <button v-else class="btn btn-success btn-sm" @click="startWorkOrder(wo)">
                      <span class="material-symbols-outlined">play_arrow</span> Iniciar
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="pendingOrders.length === 0">
                <td colspan="5" class="empty-row-msg">No hay órdenes pendientes de inicio.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <!-- IN_PROGRESS Section -->
      <section class="dashboard-section">
        <div class="section-header">
          <span class="material-symbols-outlined text-warning">running_with_errors</span>
          <h2>Producción Activa en Taller</h2>
          <span class="header-tag">{{ inProgressOrders.length }}</span>
        </div>
        <div class="table-wrapper">
          <table class="data-table">
            <thead>
              <tr>
                <th>Número</th>
                <th>Descripción</th>
                <th>Pedido</th>
                <th>Vencimiento</th>
                <th class="align-right">Gestión</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="wo in inProgressOrders" :key="wo.id" class="row-hover">
                <td><code class="code-badge">{{ wo.work_number }}</code></td>
                <td><strong>{{ wo.work_name }}</strong></td>
                <td><router-link v-if="wo.sales_order_id" :to="`/sales/orders/${wo.sales_order_id}`" class="link-primary">#{{ wo.sales_order_number }}</router-link></td>
                <td :class="{ 'text-danger fw-bold': isDeliveryUrgent(wo.due_date) }">{{ formatDate(wo.due_date) }}</td>
                <td class="align-right">
                  <button class="btn btn-outline btn-sm" @click="suspendWorkOrder(wo)">Pausar</button>
                </td>
              </tr>
              <tr v-if="inProgressOrders.length === 0">
                <td colspan="5" class="empty-row-msg">No hay trabajos activos en este momento.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </div>

    <!-- CAPA 2: CONTEXTO (Panel Lateral) -->
    <template #sidebar>
      <div class="mes-sidebar-content">
        <section class="sidebar-section mb-8">
          <div class="section-header">
            <span class="material-symbols-outlined">analytics</span>
            <h2>Resumen Histórico</h2>
          </div>
          <div class="history-cards-grid">
            <div class="admin-card clickable" @click="router.push('/mes/work-orders?status=COMPLETED')">
              <span class="material-symbols-outlined text-success">task_alt</span>
              <div class="admin-card-info">
                <strong>Terminadas</strong>
                <p>{{ stats?.by_status?.COMPLETED ?? 0 }} órdenes cerradas</p>
              </div>
              <span class="material-symbols-outlined arrow">chevron_right</span>
            </div>

            <div class="admin-card clickable mt-3" @click="router.push('/mes/work-orders?status=CANCELLED')">
              <span class="material-symbols-outlined text-danger">cancel</span>
              <div class="admin-card-info">
                <strong>Canceladas</strong>
                <p>{{ stats?.by_status?.CANCELLED ?? 0 }} órdenes anuladas</p>
              </div>
              <span class="material-symbols-outlined arrow">chevron_right</span>
            </div>
          </div>
        </section>

        <!-- Accesos rápidos técnicos -->
        <section class="sidebar-section mb-8">
          <div class="section-header">
            <span class="material-symbols-outlined">bolt</span>
            <h2>Datos Maestros</h2>
          </div>
          <div class="actions-grid-mini">
            <RouterLink to="/mes/work-setups" class="action-card-sm">
              <span class="material-symbols-outlined">settings_input_component</span>
              <span>Setups</span>
            </RouterLink>
            <RouterLink to="/mes/works" class="action-card-sm">
              <span class="material-symbols-outlined">tune</span>
              <span>Procesos</span>
            </RouterLink>
            <RouterLink to="/mes/positions" class="action-card-sm">
              <span class="material-symbols-outlined">location_on</span>
              <span>Posiciones</span>
            </RouterLink>
          </div>
        </section>

        <section class="help-notice">
          <div class="notice-header">
            <span class="material-symbols-outlined">info</span>
            <h3>Monitor de Taller</h3>
          </div>
          <p class="help-text">
            Este panel sincroniza la planificación de ventas con la ejecución de taller. Las órdenes configuradas aparecen automáticamente en el <strong>Terminal de Taller</strong> para su operario.
          </p>
        </section>
      </div>
    </template>
  </BaseDashboardPage>

  <!-- MODAL: CONFIGURACIÓN RÁPIDA -->
  <BaseDialog
    :show="showSetupDialog"
    title="Configurar Orden de Trabajo"
    icon="settings"
    confirm-text="Guardar Configuración"
    :is-confirming="isSaving"
    @close="showSetupDialog = false"
    @confirm="saveSetup"
  >
    <div class="setup-dialog-body" v-if="selectedOrder">
      <p class="mb-4">Asigne una configuración técnica predefinida a la orden <strong>{{ selectedOrder.work_number }}</strong>.</p>
      
      <div class="form-group">
        <label>Configuración de Trabajo (Setup) *</label>
        <select v-model="setupFormData.workSetupId" class="form-input">
          <option :value="null">-- Seleccione configuración --</option>
          <option v-for="setup in availableSetups" :key="setup.id" :value="setup.id">{{ setup.name }}</option>
        </select>
      </div>
      
      <div class="form-group mt-4">
        <label>Instrucciones Adicionales</label>
        <textarea v-model="setupFormData.notes" class="form-textarea" rows="3" placeholder="Observaciones para el taller..."></textarea>
      </div>
    </div>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, reactive } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import PageHeader from '@/components/layout/PageHeader.vue'
import BaseDashboardPage from '@/components/shared/BaseDashboardPage.vue'
import BaseDialog from '@/components/shared/BaseDialog.vue'
import { mesApi } from '@/services/mesApi'
import type { WorkOrder, WorkOrderDashboardStats, WorkSetup } from '@/types/mes'

const router = useRouter()
const isLoading = ref(true)
const isSaving = ref(false)
const error = ref('')
const stats = ref<WorkOrderDashboardStats | null>(null)

const pendingOrders = ref<WorkOrder[]>([])
const inProgressOrders = ref<WorkOrder[]>([])

// Estado para el modal de configuración
const showSetupDialog = ref(false)
const selectedOrder = ref<WorkOrder | null>(null)
const availableSetups = ref<WorkSetup[]>([])
const setupFormData = reactive({ workSetupId: null as string | null, notes: '' })

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
    const [s, p, i] = await Promise.all([
      mesApi.getWorkOrderDashboardStats(),
      mesApi.listWorkOrders({ status: 'PENDING' }),
      mesApi.listWorkOrders({ status: 'IN_PROGRESS' })
    ])
    stats.value = s
    pendingOrders.value = p
    inProgressOrders.value = i
  } catch (err: any) { 
    error.value = 'No se ha podido cargar el dashboard de producción.'
    console.error(err)
  } finally { 
    isLoading.value = false 
  }
}

function configureOrder(wo: WorkOrder) {
  selectedOrder.value = wo
  setupFormData.workSetupId = wo.work_setup_id || null
  setupFormData.notes = wo.description || ''
  loadSetups(wo.party_id)
  showSetupDialog.value = true
}

async function loadSetups(partyId: string) {
  try {
    const res = await mesApi.listWorkSetups({ party_id: partyId })
    availableSetups.value = res.data || (Array.isArray(res) ? res : [])
  } catch (err) {}
}

async function saveSetup() {
  if (!selectedOrder.value || !setupFormData.workSetupId) return
  isSaving.value = true
  try {
    await mesApi.updateWorkOrder(selectedOrder.value.id, {
      work_setup_id: setupFormData.workSetupId,
      description: setupFormData.notes
    })
    showSetupDialog.value = false
    await loadDashboard()
  } catch (err: any) {
    alert('Error al guardar: ' + err.message)
  } finally {
    isSaving.value = false
  }
}

async function startWorkOrder(wo: WorkOrder) {
  try { await mesApi.updateWorkOrder(wo.id, { status: 'IN_PROGRESS' }); await loadDashboard() }
  catch (err: any) {}
}

async function suspendWorkOrder(wo: WorkOrder) {
  try { await mesApi.updateWorkOrder(wo.id, { status: 'SUSPENDED' }); await loadDashboard() }
  catch (err: any) {}
}

function formatDate(v?: string) { return v ? new Date(v).toLocaleDateString('es-ES', { year: 'numeric', month: 'short', day: 'numeric' }) : '—' }
function isDeliveryUrgent(d?: string) { if(!d) return false; const diff = (new Date(d).getTime() - new Date().getTime()) / 86400000; return diff <= 3 }

onMounted(loadDashboard)
</script>

<style scoped>
/* Header Styles */
.dashboard-header { display: flex; justify-content: space-between; align-items: flex-start; }
.dashboard-header h1 { font-size: 1.75rem; color: var(--color-text-primary); margin: 0; font-family: var(--font-family-brand); }
.subtitle { color: var(--color-text-secondary); font-size: 1rem; margin: 0.25rem 0 0; }

/* Stats Grid (Patrón Dashboard Principal) */
.stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(240px, 1fr)); gap: 1.25rem; }
.stat-card { background: white; padding: 1.25rem; border-radius: var(--border-radius-lg); display: flex; align-items: center; gap: 1rem; position: relative; box-shadow: var(--box-shadow-sm); border: 1px solid var(--color-border); transition: all 0.2s ease; cursor: pointer; }
.stat-card:hover { transform: translateY(-3px); box-shadow: var(--box-shadow-md); border-color: var(--color-primary); }

.stat-icon { width: 48px; height: 48px; border-radius: 10px; display: flex; align-items: center; justify-content: center; }
.stat-icon .material-symbols-outlined { font-size: 28px; }
.stat-icon.blue { background-color: rgba(59, 130, 246, 0.1); color: #3b82f6; }
.stat-icon.red { background-color: rgba(239, 68, 68, 0.1); color: #ef4444; }
.stat-icon.yellow { background-color: rgba(230, 184, 0, 0.1); color: #E6B800; }
.stat-icon.purple { background-color: rgba(168, 85, 247, 0.1); color: #a855f7; }

.stat-info { display: flex; flex-direction: column; }
.stat-label { font-size: 0.7rem; color: var(--color-text-secondary); font-weight: 600; text-transform: uppercase; letter-spacing: 0.05em; }
.stat-value { font-size: 1.5rem; font-weight: 700; color: var(--color-text-primary); }
.stat-link-arrow { position: absolute; right: 1rem; color: var(--color-border); font-size: 18px; transition: 0.2s; }
.stat-card:hover .stat-link-arrow { color: var(--color-primary); transform: translateX(3px); }

/* Section Headers */
.dashboard-section { background: white; padding: 1.5rem; border-radius: var(--border-radius-lg); border: 1px solid var(--color-border); box-shadow: var(--box-shadow-sm); }
.section-header { display: flex; align-items: center; gap: 0.75rem; margin-bottom: 1.25rem; padding-bottom: 0.75rem; border-bottom: 1px solid var(--color-background); }
.section-header h2 { font-size: 0.9rem; font-weight: 700; margin: 0; color: var(--color-text-primary); text-transform: uppercase; letter-spacing: 0.05em; flex: 1; }
.section-header .material-symbols-outlined { color: var(--color-text-secondary); font-size: 20px; }
.header-tag { font-size: 0.65rem; font-weight: 800; padding: 0.2rem 0.6rem; background: var(--color-background); color: var(--color-secondary); border-radius: 20px; }

/* Table Elements */
.code-badge { background: var(--color-background); padding: 0.2rem 0.4rem; border-radius: 4px; font-family: var(--font-family-mono); font-size: 0.8rem; }
.priority-pill { font-size: 0.7rem; font-weight: 800; text-transform: uppercase; padding: 0.2rem 0.5rem; border-radius: 4px; }
.prio-HIGH, .prio-URGENT { background: rgba(220, 38, 38, 0.1); color: #dc2626; }
.prio-MEDIUM { background: rgba(217, 119, 6, 0.1); color: #d97706; }
.prio-LOW { background: rgba(37, 99, 235, 0.1); color: #2563eb; }

/* Sidebar Elements (Patrón Admin del Dashboard Principal) */
.admin-card { display: flex; align-items: center; gap: 1rem; padding: 1rem; background-color: var(--color-background); border-radius: 10px; border: 1px solid transparent; transition: 0.2s; cursor: pointer; }
.admin-card:hover { background-color: white; border-color: var(--color-primary); transform: translateX(4px); box-shadow: var(--box-shadow-sm); }
.admin-card-info { flex: 1; }
.admin-card-info strong { display: block; font-size: 0.85rem; color: var(--color-text-primary); }
.admin-card-info p { font-size: 0.7rem; color: var(--color-text-secondary); margin: 0; }
.admin-card .arrow { color: var(--color-border); font-size: 16px; transition: 0.2s; }
.admin-card:hover .arrow { color: var(--color-primary); }

.actions-grid-mini { display: grid; grid-template-columns: repeat(3, 1fr); gap: 0.75rem; }
.action-card-sm { display: flex; flex-direction: column; align-items: center; gap: 0.5rem; padding: 1rem 0.5rem; background: var(--color-background); border-radius: 8px; text-decoration: none; color: var(--color-text-secondary); transition: 0.2s; }
.action-card-sm:hover { background: white; color: var(--color-primary); box-shadow: var(--box-shadow-sm); transform: translateY(-2px); }
.action-card-sm .material-symbols-outlined { font-size: 24px; }
.action-card-sm span:not(.material-symbols-outlined) { font-size: 0.65rem; font-weight: 700; text-transform: uppercase; }

.help-notice { padding: 1.25rem; background: rgba(59, 130, 246, 0.05); border-radius: 12px; border: 1px dashed rgba(59, 130, 246, 0.3); }
.notice-header { display: flex; align-items: center; gap: 0.5rem; margin-bottom: 0.75rem; color: #2563eb; }
.notice-header h3 { margin: 0; font-size: 0.85rem; text-transform: uppercase; }
.help-text { font-size: 0.8rem; color: var(--color-text-secondary); line-height: 1.5; margin: 0; }

.spin { animation: spin 1s linear infinite; }
@keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }

.align-right { text-align: right; }
.link-primary { color: var(--color-secondary); font-weight: 600; text-decoration: none; }
.form-input, .form-textarea { width: 100%; padding: 0.75rem; border: 1px solid var(--color-border); border-radius: 8px; font-family: inherit; }
</style>
