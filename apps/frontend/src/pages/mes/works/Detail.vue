<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">MES / Trabajos</p>
          <h1>{{ work?.work_name || 'Detalle de trabajo' }}</h1>
          <p class="subtitle">{{ work?.work_number || '—' }}</p>
        </div>
        <RouterLink to="/mes/works" class="btn btn-secondary">Volver</RouterLink>
      </header>

      <section class="card" v-if="isLoading">
        <div class="empty-state">Cargando detalle...</div>
      </section>

      <section class="card" v-else-if="error">
        <div class="alert">{{ error }}</div>
      </section>

      <template v-else-if="work">
        <section class="card details-grid">
          <div><strong>Estado:</strong> {{ work.status }}</div>
          <div><strong>Prioridad:</strong> {{ work.priority }}</div>
          <div><strong>Party ID:</strong> {{ work.party_id }}</div>
          <div><strong>Tangible Group ID:</strong> {{ work.tangible_group_id }}</div>
          <div class="full"><strong>Notas:</strong> {{ work.garment_notes || '—' }}</div>
        </section>

        <section class="card">
          <h3>Asignaciones de servicio</h3>
          <div v-if="work.service_groups.length === 0" class="empty-state">Sin asignaciones</div>
          <div v-for="group in work.service_groups" :key="group.id" class="group-box">
            <p><strong>Service Group ID:</strong> {{ group.service_group_id }}</p>
            <p><strong>Position ID:</strong> {{ group.position_id }}</p>
            <p><strong>Secuencia:</strong> {{ group.sequence }}</p>
            <p><strong>Notas:</strong> {{ group.notes || '—' }}</p>
            <p><strong>Tareas generadas:</strong> {{ group.tasks.length }}</p>
          </div>
        </section>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { mesApi } from '@/services/mesApi'
import type { MESWork } from '@/types/mes'

const route = useRoute()
const isLoading = ref(false)
const error = ref('')
const work = ref<MESWork | null>(null)

async function loadDetail() {
  isLoading.value = true
  error.value = ''
  try {
    work.value = await mesApi.getWork(String(route.params.id))
  } catch (err: any) {
    error.value = err.message || 'No se pudo cargar el trabajo MES'
  } finally {
    isLoading.value = false
  }
}

onMounted(loadDetail)
</script>

<style scoped>
.dashboard { min-height: 100vh; background-color: #f1f5f9; }
.dashboard-content { max-width: 1100px; margin: 0 auto; padding: 2rem; display: flex; flex-direction: column; gap: 1.5rem; }
.page-header { display: flex; justify-content: space-between; align-items: center; gap: 1rem; }
.breadcrumb { font-size: .75rem; text-transform: uppercase; letter-spacing: .08em; color: #64748b; margin: 0; }
.subtitle { color: #64748b; margin: .5rem 0 0; }
.card { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 1rem; }
.details-grid { display: grid; grid-template-columns: 1fr 1fr; gap: .75rem; }
.details-grid .full { grid-column: 1 / -1; }
.group-box { border: 1px solid #e2e8f0; border-radius: 8px; padding: .75rem; margin-top: .75rem; }
.group-box p { margin: .25rem 0; }
.btn { border: none; border-radius: 8px; padding: .6rem 1rem; cursor: pointer; text-decoration: none; display: inline-flex; align-items: center; justify-content: center; width: fit-content; }
.btn-secondary { background: #fff; border: 1px solid #e2e8f0; color: #1e293b; }
.empty-state { text-align: center; color: #64748b; padding: 1rem; }
.alert { background: #fef2f2; color: #b91c1c; border: 1px solid #fecaca; border-radius: 8px; padding: .75rem; }
</style>
