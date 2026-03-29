<template>
  <BaseDashboardPage :is-loading="isLoading" class="no-print">
    <template #header>
      <PageHeader title="Monitor de Producción (MES)" :breadcrumbs="[{ label: 'MES', to: '/mes/dashboard' }, { label: 'Monitor' }]">
        <template #icon><span class="material-symbols-outlined">precision_manufacturing</span></template>
        <template #actions>
          <button @click="loadDashboard" class="btn btn-outline btn-sm" :disabled="isLoading">
            <span class="material-symbols-outlined" :class="{ 'spin': isLoading }">refresh</span>
            Actualizar
          </button>
          <button class="btn btn-primary btn-sm ml-2" @click="router.push('/mes/terminal')">
            <span class="material-symbols-outlined">tablet_mac</span>
            <span>Terminal Operario</span>
          </button>
        </template>
      </PageHeader>
    </template>

    <div class="module-dashboard-content">
      <!-- 1. KPIs de Resumen -->
      <section class="stats-grid">
        <div class="stat-card clickable" @click="router.push('/mes/work-orders')">
          <div class="stat-icon blue"><span class="material-symbols-outlined">factory</span></div>
          <div class="stat-info">
            <span class="stat-label">Total Trabajos</span>
            <span class="stat-value">{{ stats?.total ?? 0 }}</span>
          </div>
        </div>
        <div class="stat-card clickable" @click="router.push('/mes/work-orders?status=OVERDUE')">
          <div class="stat-icon red"><span class="material-symbols-outlined">history</span></div>
          <div class="stat-info">
            <span class="stat-label">Vencidos</span>
            <span class="stat-value text-danger">{{ stats?.overdue ?? 0 }}</span>
          </div>
        </div>
        <div class="stat-card clickable" @click="router.push('/mes/work-orders?status=IN_PROGRESS')">
          <div class="stat-icon yellow"><span class="material-symbols-outlined">running_with_errors</span></div>
          <div class="stat-info">
            <span class="stat-label">En Producción</span>
            <span class="stat-value text-warning">{{ inProgressOrders.length }}</span>
          </div>
        </div>
        <div class="stat-card clickable" @click="router.push('/mes/terminal')">
          <div class="stat-icon purple"><span class="material-symbols-outlined">assignment</span></div>
          <div class="stat-info">
            <span class="stat-label">Tareas Taller</span>
            <span class="stat-value">{{ pendingTasksCount }}</span>
          </div>
        </div>
      </section>

      <!-- 2. Accesos a Listados -->
      <section class="listings-grid">
        <RouterLink to="/mes/work-orders" class="listing-link">
          <span class="material-symbols-outlined">format_list_bulleted</span>
          <span>Órdenes de Trabajo</span>
        </RouterLink>
        <RouterLink to="/mes/tasks" class="listing-link">
          <span class="material-symbols-outlined">checklist</span>
          <span>Tareas de Taller</span>
        </RouterLink>
        <RouterLink to="/mes/positions" class="listing-link">
          <span class="material-symbols-outlined">location_on</span>
          <span>Mapa de Posiciones</span>
        </RouterLink>
        <RouterLink to="/mes/terminal" class="listing-link highlight-subtle">
          <span class="material-symbols-outlined">tablet_mac</span>
          <span>Terminal Taller</span>
        </RouterLink>
      </section>

      <!-- 3. Actividad (Órdenes Pendientes) -->
      <section class="dashboard-section">
        <div class="section-header">
          <span class="material-symbols-outlined text-info">pending_actions</span>
          <h2>Pendientes de Inicio</h2>
          <span class="header-tag">{{ pendingOrders.length }}</span>
        </div>
        <div class="table-wrapper">
          <table class="data-table">
            <thead>
              <tr>
                <th>Nº</th>
                <th>Descripción</th>
                <th>Pedido</th>
                <th>Prioridad</th>
                <th class="align-right">Acción</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="wo in pendingOrders" :key="wo.id" class="row-hover">
                <td><code class="code-badge">{{ wo.work_number }}</code></td>
                <td><strong>{{ wo.work_name }}</strong></td>
                <td>#{{ wo.sales_order_number }}</td>
                <td><span :class="['priority-pill', `prio-${wo.priority}`]">{{ wo.priority }}</span></td>
                <td class="align-right">
                  <button class="btn btn-primary btn-sm" @click="configureOrder(wo)">Configurar</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </div>

    <template #sidebar>
      <section class="sidebar-section">
        <div class="section-header">
          <span class="material-symbols-outlined">bolt</span>
          <h2>Mantenimiento</h2>
        </div>
        <div class="quick-actions-list">
          <RouterLink to="/mes/work-setups" class="admin-card clickable">
            <span class="material-symbols-outlined text-primary">settings_input_component</span>
            <div class="admin-card-info">
              <strong>Setups Técnicos</strong>
              <p>Configuraciones base</p>
            </div>
          </RouterLink>
          <RouterLink to="/mes/tasks" class="admin-card clickable mt-2">
            <span class="material-symbols-outlined text-secondary">tune</span>
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
import { computed, onMounted, ref, reactive } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import PageHeader from '@/components/layout/PageHeader.vue'
import BaseDashboardPage from '@/components/shared/BaseDashboardPage.vue'
import { mesApi } from '@/services/mesApi'
import type { WorkOrder, WorkOrderDashboardStats, WorkSetup } from '@/types/mes'
const router = useRouter()
const isLoading = ref(true)
const stats = ref<WorkOrderDashboardStats | null>(null)
const pendingOrders = ref<WorkOrder[]>([])
const inProgressOrders = ref<WorkOrder[]>([])
const pendingTasksCount = computed(() => {
  let count = 0
  for (const wo of inProgressOrders.value) { for (const line of wo.lines || []) { for (const task of line.tasks || []) { if (task.status === 'PENDING') count++ } } }
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
  } catch (err) {}
  finally { isLoading.value = false }
}
function configureOrder(wo: any) { router.push(`/mes/work-orders/${wo.id}`) }
onMounted(loadDashboard)
</script>

<style scoped>
.module-dashboard-content { display: flex; flex-direction: column; gap: 1.5rem; }

.stats-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 0.75rem; }
.stat-card { background: white; padding: 0.75rem 1rem; border-radius: 10px; border: 1px solid var(--color-border); display: flex; align-items: center; gap: 0.75rem; position: relative; transition: 0.2s; cursor: pointer; }
.stat-card:hover { transform: translateY(-2px); box-shadow: var(--box-shadow-md); border-color: var(--color-primary); }
.stat-icon { width: 40px; height: 40px; border-radius: 8px; display: flex; align-items: center; justify-content: center; }
.stat-icon .material-symbols-outlined { font-size: 22px; }
.stat-icon.blue { background: rgba(59, 130, 246, 0.1); color: #3b82f6; }
.stat-icon.red { background: rgba(239, 68, 68, 0.1); color: #ef4444; }
.stat-icon.yellow { background: rgba(230, 184, 0, 0.1); color: #E6B800; }
.stat-icon.purple { background: rgba(168, 85, 247, 0.1); color: #a855f7; }
.stat-label { font-size: 0.65rem; color: var(--color-text-secondary); font-weight: 600; text-transform: uppercase; }
.stat-value { font-size: 1.25rem; font-weight: 700; }

.listings-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 0.75rem; }
.listing-link { display: flex; align-items: center; gap: 0.75rem; padding: 0.75rem 1rem; background: white; border: 1px solid var(--color-border); border-radius: 8px; text-decoration: none; color: var(--color-text-primary); font-size: 0.85rem; font-weight: 600; transition: 0.2s; }
.listing-link:hover { background: var(--color-background); border-color: var(--color-secondary); color: var(--color-secondary); }
.listing-link .material-symbols-outlined { color: var(--color-secondary); font-size: 1.25rem; }
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
</style>
