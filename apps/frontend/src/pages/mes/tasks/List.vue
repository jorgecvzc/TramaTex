<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">MES / Tareas</p>
          <h1>Tareas</h1>
          <p class="subtitle">Administra las tareas base de producción.</p>
        </div>
        <RouterLink to="/mes/tasks/new" class="btn btn-primary">Nueva tarea</RouterLink>
      </header>

      <section class="card filters">
        <input v-model="search" type="text" placeholder="Buscar por nombre" class="input" />
        <select v-model="statusFilter" class="input">
          <option value="">Todas</option>
          <option value="true">Activas</option>
          <option value="false">Inactivas</option>
        </select>
        <button @click="loadTasks" class="btn btn-secondary">Filtrar</button>
      </section>

      <section class="card">
        <div v-if="isLoading" class="empty-state">Cargando tareas...</div>
        <div v-else-if="error" class="alert">{{ error }}</div>
        <table v-else class="data-table">
          <thead>
            <tr>
              <th>Nombre</th>
              <th>Referencia</th>
              <th>Descripción</th>
              <th>Estado</th>
              <th>Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="task in tasks" :key="task.id">
              <td><strong>{{ task.name }}</strong></td>
              <td>{{ task.reference || '—' }}</td>
              <td>{{ task.description || '—' }}</td>
              <td>
                <span class="badge" :class="task.is_active ? 'ok' : 'off'">
                  {{ task.is_active ? 'Activa' : 'Inactiva' }}
                </span>
              </td>
              <td class="actions">
                <RouterLink :to="`/mes/tasks/${task.id}/edit`" class="btn btn-sm">Editar</RouterLink>
                <button @click="toggleActive(task)" class="btn btn-sm" :class="task.is_active ? 'btn-off' : 'btn-on'">
                  {{ task.is_active ? 'Desactivar' : 'Activar' }}
                </button>
              </td>
            </tr>
            <tr v-if="tasks.length === 0">
              <td colspan="5" class="empty-state">No hay tareas registradas.</td>
            </tr>
          </tbody>
        </table>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { mesApi } from '@/services/mesApi'
import type { MESTask } from '@/types/mes'

const tasks = ref<MESTask[]>([])
const isLoading = ref(false)
const error = ref('')
const search = ref('')
const statusFilter = ref('')

async function loadTasks() {
  isLoading.value = true
  error.value = ''

  try {
    const isActive = statusFilter.value === '' ? undefined : statusFilter.value === 'true'
    tasks.value = await mesApi.listTasks({
      search: search.value.trim() || undefined,
      is_active: isActive,
    })
  } catch (err: any) {
    error.value = err.message || 'No se pudieron cargar las tareas'
  } finally {
    isLoading.value = false
  }
}

async function toggleActive(task: MESTask) {
  try {
    await mesApi.updateTask(task.id, { is_active: !task.is_active })
    await loadTasks()
  } catch (err: any) {
    error.value = err.message || 'No se pudo cambiar el estado de la tarea'
  }
}

onMounted(loadTasks)
</script>

<style scoped>
.dashboard { min-height: 100vh; background-color: #f1f5f9; }
.dashboard-content { max-width: 1200px; margin: 0 auto; padding: 2rem; display: flex; flex-direction: column; gap: 1.5rem; }
.page-header { display: flex; justify-content: space-between; align-items: center; gap: 1rem; }
.breadcrumb { font-size: .75rem; text-transform: uppercase; letter-spacing: .08em; color: #64748b; margin: 0; }
.subtitle { color: #64748b; margin: .5rem 0 0; }
.card { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 1rem; }
.filters { display: grid; grid-template-columns: 1fr 180px auto; gap: .75rem; }
.input { border: 1px solid #cbd5e1; border-radius: 8px; padding: .6rem .75rem; }
.btn { border: none; border-radius: 8px; padding: .6rem 1rem; cursor: pointer; text-decoration: none; display: inline-flex; align-items: center; justify-content: center; }
.btn-primary { background: #f4c430; color: #1b3a6b; font-weight: 600; }
.btn-secondary { background: #fff; border: 1px solid #e2e8f0; color: #1e293b; }
.data-table { width: 100%; border-collapse: collapse; }
.data-table th, .data-table td { padding: .75rem; border-bottom: 1px solid #e2e8f0; text-align: left; }
.badge { padding: .2rem .55rem; border-radius: 999px; font-size: .75rem; font-weight: 600; }
.badge.ok { background: #dcfce7; color: #166534; }
.badge.off { background: #e2e8f0; color: #475569; }
.empty-state { text-align: center; color: #64748b; padding: 1rem; }
.alert { background: #fef2f2; color: #b91c1c; border: 1px solid #fecaca; border-radius: 8px; padding: .75rem; }
.actions { display: flex; gap: .5rem; align-items: center; }
.btn-sm { font-size: .8rem; padding: .35rem .65rem; border-radius: 6px; }
.btn-off { background: #fef2f2; color: #b91c1c; border: 1px solid #fecaca; }
.btn-on { background: #dcfce7; color: #166534; border: 1px solid #bbf7d0; }
</style>
