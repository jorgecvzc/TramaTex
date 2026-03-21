<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">MES / Órdenes de Trabajo</p>
          <h1>Órdenes de trabajo</h1>
          <p class="subtitle">Seguimiento y gestión de órdenes de producción.</p>
        </div>
        <RouterLink to="/mes/work-orders/new" class="btn btn-primary">Nueva orden</RouterLink>
      </header>

      <section class="card filters">
        <input v-model="search" type="text" placeholder="Buscar por nombre, número o cliente (nombre/referencia)" class="input" />
        <select v-model="statusFilter" class="input">
          <option value="">Todos</option>
          <option value="PENDING">Pendiente</option>
          <option value="IN_PROGRESS">En progreso</option>
          <option value="ON_HOLD">En espera</option>
          <option value="COMPLETED">Completado</option>
          <option value="CANCELLED">Cancelado</option>
        </select>
        <button @click="loadWorks" class="btn btn-secondary">Filtrar</button>
      </section>

      <section class="card">
        <div v-if="isLoading" class="empty-state">Cargando trabajos...</div>
        <div v-else-if="error" class="alert">{{ error }}</div>
        <table v-else class="data-table">
          <thead>
            <tr>
              <th>Número</th>
              <th>Nombre</th>
              <th>Estado</th>
              <th>Prioridad</th>
              <th>Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="work in works" :key="work.id">
              <td><strong>{{ work.work_number }}</strong></td>
              <td>{{ work.work_name }}</td>
              <td>{{ mesApi.getWorkStatusLabel(work.status) }}</td>
              <td>{{ mesApi.getPriorityLabel(work.priority) }}</td>
              <td>
                <RouterLink :to="`/mes/work-orders/${work.id}`" class="btn-link">Ver</RouterLink>
              </td>
            </tr>
            <tr v-if="works.length === 0">
              <td colspan="5" class="empty-state">No hay órdenes de trabajo registradas.</td>
            </tr>
          </tbody>
        </table>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { mesApi } from '@/services/mesApi'
import type { WorkOrder } from '@/types/mes'

const works = ref<WorkOrder[]>([])
const isLoading = ref(false)
const error = ref('')
const route = useRoute()
const search = ref('')
const statusFilter = ref((route.query.status as string) || '')

async function loadWorks() {
  isLoading.value = true
  error.value = ''

  try {
    works.value = await mesApi.listWorkOrders({
      search: search.value.trim() || undefined,
      status: statusFilter.value || undefined,
    })
  } catch (err: any) {
    error.value = err.message || 'No se pudieron cargar las órdenes de trabajo'
  } finally {
    isLoading.value = false
  }
}

onMounted(loadWorks)
</script>

<style scoped>
.dashboard { min-height: 100vh; background-color: #f1f5f9; }
.dashboard-content { max-width: 1200px; margin: 0 auto; padding: 2rem; display: flex; flex-direction: column; gap: 1.5rem; }
.page-header { display: flex; justify-content: space-between; align-items: center; gap: 1rem; }
.breadcrumb { font-size: .75rem; text-transform: uppercase; letter-spacing: .08em; color: #64748b; margin: 0; }
.subtitle { color: #64748b; margin: .5rem 0 0; }
.card { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 1rem; }
.filters { display: grid; grid-template-columns: 1fr 220px auto; gap: .75rem; }
.input { border: 1px solid #cbd5e1; border-radius: 8px; padding: .6rem .75rem; }
.btn { border: none; border-radius: 8px; padding: .6rem 1rem; cursor: pointer; text-decoration: none; display: inline-flex; align-items: center; justify-content: center; }
.btn-primary { background: #f4c430; color: #1b3a6b; font-weight: 600; }
.btn-secondary { background: #fff; border: 1px solid #e2e8f0; color: #1e293b; }
.btn-link { color: #1d4ed8; text-decoration: none; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th, .data-table td { padding: .75rem; border-bottom: 1px solid #e2e8f0; text-align: left; }
.empty-state { text-align: center; color: #64748b; padding: 1rem; }
.alert { background: #fef2f2; color: #b91c1c; border: 1px solid #fecaca; border-radius: 8px; padding: .75rem; }
</style>
