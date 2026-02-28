<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">MES / Dashboard</p>
          <h1>Monitoreo de producción</h1>
          <p class="subtitle">Resumen de trabajos, estados y vencimientos.</p>
        </div>
        <div class="header-actions">
          <RouterLink to="/mes/terminal" class="btn btn-primary">Terminal Tablet</RouterLink>
          <RouterLink to="/mes/work-definitions" class="btn btn-secondary">Ver definiciones</RouterLink>
        </div>
      </header>

      <section class="stats-grid">
        <article class="card stat-card">
          <p class="stat-label">Total trabajos</p>
          <p class="stat-value">{{ stats?.total ?? 0 }}</p>
        </article>
        <article class="card stat-card">
          <p class="stat-label">Vencidos</p>
          <p class="stat-value danger">{{ stats?.overdue ?? 0 }}</p>
        </article>
        <article class="card stat-card">
          <p class="stat-label">Vencen hoy</p>
          <p class="stat-value warning">{{ stats?.due_today ?? 0 }}</p>
        </article>
        <article class="card stat-card">
          <p class="stat-label">En progreso</p>
          <p class="stat-value">{{ stats?.by_status?.IN_PROGRESS ?? 0 }}</p>
        </article>
      </section>

      <section class="card">
        <h2>Distribución por estado</h2>
        <div v-if="statusRows.length === 0" class="empty-state">Sin datos de estado.</div>
        <table v-else class="data-table">
          <thead>
            <tr>
              <th>Estado</th>
              <th>Cantidad</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in statusRows" :key="row.status">
              <td>{{ row.status }}</td>
              <td>{{ row.count }}</td>
            </tr>
          </tbody>
        </table>
      </section>

      <section class="card">
        <h2>Definiciones vencidas</h2>
        <div v-if="isLoading" class="empty-state">Cargando dashboard MES...</div>
        <div v-else-if="error" class="alert">{{ error }}</div>
        <div v-else-if="overdueWorks.length === 0" class="empty-state">No hay definiciones vencidas.</div>
        <table v-else class="data-table">
          <thead>
            <tr>
              <th>Número</th>
              <th>Nombre</th>
              <th>Estado</th>
              <th>Vencimiento</th>
              <th>Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="work in overdueWorks" :key="work.id">
              <td><strong>{{ work.work_number }}</strong></td>
              <td>{{ work.work_name }}</td>
              <td>{{ work.status }}</td>
              <td>{{ formatDate(work.due_date) }}</td>
              <td>
                <RouterLink :to="`/mes/work-definitions/${work.id}`" class="btn-link">Ver</RouterLink>
              </td>
            </tr>
          </tbody>
        </table>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { mesApi } from '@/services/mesApi'
import type { MESWork, MESWorkDashboardStats } from '@/types/mes'

const isLoading = ref(false)
const error = ref('')
const stats = ref<MESWorkDashboardStats | null>(null)
const overdueWorks = ref<MESWork[]>([])

const statusRows = computed(() => {
  if (!stats.value?.by_status) {
    return []
  }

  return Object.entries(stats.value.by_status)
    .map(([status, count]) => ({ status, count }))
    .sort((a, b) => b.count - a.count)
})

function formatDate(value?: string) {
  if (!value) {
    return '—'
  }
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  return date.toLocaleDateString('es-ES')
}

async function loadDashboard() {
  isLoading.value = true
  error.value = ''

  try {
    const [statsResult, overdueResult] = await Promise.all([
      mesApi.getWorkDashboardStats(),
      mesApi.listOverdueWorks(20),
    ])

    stats.value = statsResult
    overdueWorks.value = overdueResult
  } catch (err: any) {
    error.value = err.message || 'No se pudo cargar el dashboard MES'
  } finally {
    isLoading.value = false
  }
}

onMounted(loadDashboard)
</script>

<style scoped>
.dashboard { min-height: 100vh; background-color: #f1f5f9; }
.dashboard-content { max-width: 1200px; margin: 0 auto; padding: 2rem; display: flex; flex-direction: column; gap: 1.5rem; }
.page-header { display: flex; justify-content: space-between; align-items: center; gap: 1rem; }
.header-actions { display: flex; gap: .5rem; }
.breadcrumb { font-size: .75rem; text-transform: uppercase; letter-spacing: .08em; color: #64748b; margin: 0; }
.subtitle { color: #64748b; margin: .5rem 0 0; }
.stats-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: .75rem; }
.stat-card { display: flex; flex-direction: column; gap: .25rem; }
.stat-label { margin: 0; color: #64748b; font-size: .85rem; }
.stat-value { margin: 0; font-size: 1.8rem; font-weight: 700; color: #1e293b; }
.stat-value.danger { color: #b91c1c; }
.stat-value.warning { color: #b45309; }
.card { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 1rem; }
.btn { border: none; border-radius: 8px; padding: .6rem 1rem; cursor: pointer; text-decoration: none; display: inline-flex; align-items: center; justify-content: center; }
.btn-primary { background: #f4c430; color: #1b3a6b; font-weight: 600; }
.btn-secondary { background: #fff; border: 1px solid #e2e8f0; color: #1e293b; }
.btn-link { color: #1d4ed8; text-decoration: none; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th, .data-table td { padding: .75rem; border-bottom: 1px solid #e2e8f0; text-align: left; }
.empty-state { text-align: center; color: #64748b; padding: 1rem; }
.alert { background: #fef2f2; color: #b91c1c; border: 1px solid #fecaca; border-radius: 8px; padding: .75rem; }
@media (max-width: 900px) { .stats-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); } }
</style>
