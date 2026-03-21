<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">MES / Tipos de Trabajo</p>
          <h1>Tipos de trabajo</h1>
          <p class="subtitle">Administra las secuencias de tareas para producción.</p>
        </div>
        <RouterLink to="/mes/work-types/new" class="btn btn-primary">Nuevo tipo</RouterLink>
      </header>

      <section class="card filters">
        <input v-model="search" type="text" placeholder="Buscar por nombre" class="input" />
        <select v-model="statusFilter" class="input">
          <option value="">Todos</option>
          <option value="true">Activos</option>
          <option value="false">Inactivos</option>
        </select>
        <button @click="loadWorkTypes" class="btn btn-secondary">Filtrar</button>
      </section>

      <section class="card">
        <div v-if="isLoading" class="empty-state">Cargando tipos de trabajo...</div>
        <div v-else-if="error" class="alert">{{ error }}</div>
        <table v-else class="data-table">
          <thead>
            <tr>
              <th>Nombre</th>
              <th>Referencia</th>
              <th>Descripción</th>
              <th>Tasks</th>
              <th>Estado</th>
              <th>Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="group in workTypes" :key="group.id">
              <td><strong>{{ group.name }}</strong></td>
              <td>{{ group.reference || '—' }}</td>
              <td>{{ group.description || '—' }}</td>
              <td>{{ group.tasks?.length || 0 }}</td>
              <td>
                <span class="badge" :class="group.is_active ? 'ok' : 'off'">
                  {{ group.is_active ? 'Activo' : 'Inactivo' }}
                </span>
              </td>
              <td class="actions">
                <RouterLink :to="`/mes/work-types/${group.id}/edit`" class="btn btn-sm">Editar</RouterLink>
                <button @click="toggleActive(group)" class="btn btn-sm" :class="group.is_active ? 'btn-off' : 'btn-on'">
                  {{ group.is_active ? 'Desactivar' : 'Activar' }}
                </button>
              </td>
            </tr>
            <tr v-if="workTypes.length === 0">
              <td colspan="6" class="empty-state">No hay tipos de trabajo registrados.</td>
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
import type { MESWorkType } from '@/types/mes'

const workTypes = ref<MESWorkType[]>([])
const isLoading = ref(false)
const error = ref('')
const search = ref('')
const statusFilter = ref('')

async function loadWorkTypes() {
  isLoading.value = true
  error.value = ''

  try {
    const isActive = statusFilter.value === '' ? undefined : statusFilter.value === 'true'
    workTypes.value = await mesApi.listWorkTypes({
      search: search.value.trim() || undefined,
      is_active: isActive,
    })
  } catch (err: any) {
    error.value = err.message || 'No se pudieron cargar los tipos de trabajo'
  } finally {
    isLoading.value = false
  }
}

async function toggleActive(group: MESWorkType) {
  try {
    await mesApi.updateWorkType(group.id, { is_active: !group.is_active })
    await loadWorkTypes()
  } catch (err: any) {
    error.value = err.message || 'No se pudo cambiar el estado del tipo de trabajo'
  }
}

onMounted(loadWorkTypes)
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
