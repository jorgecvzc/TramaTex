<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">MES / Posiciones</p>
          <h1>Posiciones</h1>
          <p class="subtitle">Administra posiciones físicas de aplicación en la prenda.</p>
        </div>
        <RouterLink to="/mes/positions/new" class="btn btn-primary">Nueva posición</RouterLink>
      </header>

      <section class="card filters">
        <input v-model="search" type="text" placeholder="Buscar por nombre o código" class="input" />
        <select v-model="statusFilter" class="input">
          <option value="">Todas</option>
          <option value="true">Activas</option>
          <option value="false">Inactivas</option>
        </select>
        <button @click="loadPositions" class="btn btn-secondary">Filtrar</button>
      </section>

      <section class="card">
        <div v-if="isLoading" class="empty-state">Cargando posiciones...</div>
        <div v-else-if="error" class="alert">{{ error }}</div>
        <table v-else class="data-table">
          <thead>
            <tr>
              <th>Nombre</th>
              <th>Código</th>
              <th>Descripción</th>
              <th>Estado</th>
              <th>Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="position in positions" :key="position.id">
              <td><strong>{{ position.name }}</strong></td>
              <td><code class="code">{{ position.code }}</code></td>
              <td>{{ position.description || '—' }}</td>
              <td>
                <span class="badge" :class="position.is_active ? 'ok' : 'off'">
                  {{ position.is_active ? 'Activa' : 'Inactiva' }}
                </span>
              </td>
              <td class="actions">
                <RouterLink :to="`/mes/positions/${position.id}/edit`" class="btn btn-sm">Editar</RouterLink>
                <button @click="toggleActive(position)" class="btn btn-sm" :class="position.is_active ? 'btn-off' : 'btn-on'">
                  {{ position.is_active ? 'Desactivar' : 'Activar' }}
                </button>
              </td>
            </tr>
            <tr v-if="positions.length === 0">
              <td colspan="5" class="empty-state">No hay posiciones registradas.</td>
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
import type { MESPosition } from '@/types/mes'

const positions = ref<MESPosition[]>([])
const isLoading = ref(false)
const error = ref('')
const search = ref('')
const statusFilter = ref('')

async function loadPositions() {
  isLoading.value = true
  error.value = ''

  try {
    const isActive = statusFilter.value === '' ? undefined : statusFilter.value === 'true'
    positions.value = await mesApi.listPositions({
      search: search.value.trim() || undefined,
      is_active: isActive,
    })
  } catch (err: any) {
    error.value = err.message || 'No se pudieron cargar las posiciones'
  } finally {
    isLoading.value = false
  }
}

async function toggleActive(position: MESPosition) {
  try {
    await mesApi.updatePosition(position.id, { is_active: !position.is_active })
    await loadPositions()
  } catch (err: any) {
    error.value = err.message || 'No se pudo cambiar el estado de la posición'
  }
}

onMounted(loadPositions)
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
.code { background: #f1f5f9; padding: .2rem .45rem; border-radius: 6px; }
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
