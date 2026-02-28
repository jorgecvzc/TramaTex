<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">MES / Definiciones de trabajo</p>
          <h1>{{ work?.work_name || 'Detalle de definición' }}</h1>
          <p class="subtitle">{{ work?.work_number || '—' }}</p>
        </div>
        <div class="header-actions" v-if="work && !isEditing">
          <button class="btn btn-primary" @click="startEdit">Editar</button>
          <RouterLink to="/mes/work-definitions" class="btn btn-secondary">Volver</RouterLink>
        </div>
        <div class="header-actions" v-else>
          <RouterLink to="/mes/work-definitions" class="btn btn-secondary">Volver</RouterLink>
        </div>
      </header>

      <section class="card" v-if="isLoading">
        <div class="empty-state">Cargando detalle...</div>
      </section>

      <section class="card" v-else-if="error">
        <div class="alert">{{ error }}</div>
      </section>

      <template v-else-if="work">
        <section class="card" v-if="isEditing">
          <h3>Editar definición de trabajo</h3>
          <div class="edit-grid">
            <div>
              <label class="label">Nombre</label>
              <input v-model="editForm.work_name" class="input" type="text" />
            </div>
            <div>
              <label class="label">Estado</label>
              <select v-model="editForm.status" class="input">
                <option value="DRAFT">DRAFT</option>
                <option value="PENDING">PENDING</option>
                <option value="IN_PROGRESS">IN_PROGRESS</option>
                <option value="ON_HOLD">ON_HOLD</option>
                <option value="COMPLETED">COMPLETED</option>
                <option value="CANCELLED">CANCELLED</option>
              </select>
            </div>
            <div>
              <label class="label">Prioridad</label>
              <select v-model="editForm.priority" class="input">
                <option value="LOW">LOW</option>
                <option value="NORMAL">NORMAL</option>
                <option value="HIGH">HIGH</option>
                <option value="URGENT">URGENT</option>
              </select>
            </div>
            <div>
              <label class="label">Fecha de vencimiento</label>
              <input v-model="editForm.due_date" class="input" type="date" />
            </div>
            <div class="full">
              <label class="label">Notas</label>
              <textarea v-model="editForm.garment_notes" class="input" rows="3" />
            </div>
          </div>
          <div class="edit-actions">
            <button class="btn btn-secondary" @click="cancelEdit" :disabled="isSaving">Cancelar</button>
            <button class="btn btn-primary" @click="saveChanges" :disabled="isSaving">
              {{ isSaving ? 'Guardando...' : 'Guardar cambios' }}
            </button>
          </div>
        </section>

        <section class="card details-grid">
          <div><strong>Estado:</strong> {{ work.status }}</div>
          <div><strong>Prioridad:</strong> {{ work.priority }}</div>
          <div><strong>Cliente:</strong> {{ partyName || work.party_id }}</div>
          <div><strong>Categoría tangible:</strong> {{ tangibleGroupName || work.tangible_group_id }}</div>
          <div class="full"><strong>Notas:</strong> {{ work.garment_notes || '—' }}</div>
        </section>

        <section class="card">
          <h3>Asignaciones de plantilla de proceso</h3>
          <div v-if="work.service_groups.length === 0" class="empty-state">Sin asignaciones</div>
          <div v-for="group in work.service_groups" :key="group.id" class="group-box">
            <p><strong>Plantilla de proceso:</strong> {{ serviceTemplateNames[group.service_group_id] || group.service_group_id }}</p>
            <p><strong>Posición:</strong> {{ positionNames[group.position_id] || group.position_id }}</p>
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
import { partyApi } from '@/services/partyApi'
import { productApi } from '@/services/productApi'
import type { MESPosition, MESServiceTemplate, MESWork } from '@/types/mes'

const route = useRoute()
const isLoading = ref(false)
const error = ref('')
const work = ref<MESWork | null>(null)
const partyName = ref('')
const tangibleGroupName = ref('')
const serviceTemplateNames = ref<Record<string, string>>({})
const positionNames = ref<Record<string, string>>({})
const isEditing = ref(false)
const isSaving = ref(false)
const editForm = ref({
  work_name: '',
  status: 'DRAFT',
  priority: 'NORMAL',
  due_date: '',
  garment_notes: '',
})

function dateToInput(value?: string) {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return date.toISOString().slice(0, 10)
}

function startEdit() {
  if (!work.value) return
  editForm.value = {
    work_name: work.value.work_name || '',
    status: work.value.status || 'DRAFT',
    priority: work.value.priority || 'NORMAL',
    due_date: dateToInput(work.value.due_date),
    garment_notes: work.value.garment_notes || '',
  }
  isEditing.value = true
}

function cancelEdit() {
  isEditing.value = false
}

async function saveChanges() {
  if (!work.value) return
  isSaving.value = true
  error.value = ''
  try {
    const updated = await mesApi.updateWorkDefinition(work.value.id, {
      work_name: editForm.value.work_name.trim(),
      status: editForm.value.status,
      priority: editForm.value.priority,
      due_date: editForm.value.due_date || '',
      garment_notes: editForm.value.garment_notes.trim(),
    })
    work.value = updated
    isEditing.value = false
  } catch (err: any) {
    error.value = err.message || 'No se pudo actualizar la definición de trabajo MES'
  } finally {
    isSaving.value = false
  }
}

async function loadLookupData(currentWork: MESWork) {
  const [partyResult, groupResult, templates, positions] = await Promise.all([
    partyApi.getParty(currentWork.party_id),
    productApi.getProductGroup(currentWork.tangible_group_id),
    mesApi.listServiceTemplates({}),
    mesApi.listPositions({}),
  ])

  partyName.value = partyResult?.name || ''
  tangibleGroupName.value = groupResult?.name || ''

  const templateMap: Record<string, string> = {}
  for (const template of templates as MESServiceTemplate[]) {
    templateMap[template.id] = template.name
  }
  serviceTemplateNames.value = templateMap

  const positionMap: Record<string, string> = {}
  for (const position of positions as MESPosition[]) {
    positionMap[position.id] = `${position.name} (${position.code})`
  }
  positionNames.value = positionMap
}

async function loadDetail() {
  isLoading.value = true
  error.value = ''
  try {
    work.value = await mesApi.getWorkDefinition(String(route.params.id))
    await loadLookupData(work.value)
  } catch (err: any) {
    error.value = err.message || 'No se pudo cargar la definición de trabajo MES'
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
.header-actions { display: flex; gap: .5rem; }
.breadcrumb { font-size: .75rem; text-transform: uppercase; letter-spacing: .08em; color: #64748b; margin: 0; }
.subtitle { color: #64748b; margin: .5rem 0 0; }
.card { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 1rem; }
.details-grid { display: grid; grid-template-columns: 1fr 1fr; gap: .75rem; }
.details-grid .full { grid-column: 1 / -1; }
.edit-grid { display: grid; grid-template-columns: 1fr 1fr; gap: .75rem; }
.edit-grid .full { grid-column: 1 / -1; }
.label { display: block; font-size: .85rem; color: #334155; margin-bottom: .35rem; font-weight: 600; }
.input { border: 1px solid #cbd5e1; border-radius: 8px; padding: .6rem .75rem; font: inherit; width: 100%; }
.edit-actions { margin-top: .75rem; display: flex; justify-content: flex-end; gap: .5rem; }
.group-box { border: 1px solid #e2e8f0; border-radius: 8px; padding: .75rem; margin-top: .75rem; }
.group-box p { margin: .25rem 0; }
.btn { border: none; border-radius: 8px; padding: .6rem 1rem; cursor: pointer; text-decoration: none; display: inline-flex; align-items: center; justify-content: center; width: fit-content; }
.btn-primary { background: #f4c430; color: #1b3a6b; font-weight: 600; }
.btn-secondary { background: #fff; border: 1px solid #e2e8f0; color: #1e293b; }
.empty-state { text-align: center; color: #64748b; padding: 1rem; }
.alert { background: #fef2f2; color: #b91c1c; border: 1px solid #fecaca; border-radius: 8px; padding: .75rem; }
</style>
