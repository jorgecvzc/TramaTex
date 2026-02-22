<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">MES / Trabajos</p>
          <h1>Nuevo trabajo MES</h1>
          <p class="subtitle">Crea una orden de manufactura y genera tareas automáticamente.</p>
        </div>
        <RouterLink to="/mes/works" class="btn btn-secondary">Volver</RouterLink>
      </header>

      <section class="card form-card">
        <div v-if="error" class="alert">{{ error }}</div>

        <label class="label">Nombre *</label>
        <input v-model="form.work_name" type="text" class="input" placeholder="Ej: Uniformes Cliente A" />

        <label class="label">Party ID *</label>
        <input v-model="form.party_id" type="text" class="input" placeholder="party-id" />

        <label class="label">Tangible Group ID *</label>
        <input v-model="form.tangible_group_id" type="text" class="input" placeholder="UUID de product group tangible" />

        <label class="label">Estado</label>
        <select v-model="form.status" class="input">
          <option value="DRAFT">DRAFT</option>
          <option value="PENDING">PENDING</option>
          <option value="IN_PROGRESS">IN_PROGRESS</option>
          <option value="ON_HOLD">ON_HOLD</option>
        </select>

        <label class="label">Prioridad</label>
        <select v-model="form.priority" class="input">
          <option value="LOW">LOW</option>
          <option value="NORMAL">NORMAL</option>
          <option value="HIGH">HIGH</option>
          <option value="URGENT">URGENT</option>
        </select>

        <label class="label">Observaciones</label>
        <textarea v-model="form.garment_notes" class="input" rows="3" placeholder="Observaciones de prenda" />

        <div class="tasks-block">
          <div class="tasks-header">
            <h3>Asignaciones de servicio *</h3>
            <button @click="addAssignment" type="button" class="btn btn-secondary btn-sm">+ Añadir</button>
          </div>

          <div v-for="(assignment, index) in form.service_group_assignments" :key="index" class="assignment-row">
            <input v-model="assignment.service_group_id" class="input" placeholder="Service Group ID" />
            <input v-model="assignment.position_id" class="input" placeholder="Position ID" />
            <input v-model.number="assignment.sequence" type="number" min="1" class="input seq" placeholder="Seq" />
            <button @click="removeAssignment(index)" type="button" class="btn btn-danger btn-sm">Quitar</button>
          </div>
        </div>

        <button @click="submit" :disabled="isSaving" class="btn btn-primary">
          {{ isSaving ? 'Guardando...' : 'Crear trabajo' }}
        </button>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { mesApi } from '@/services/mesApi'

const router = useRouter()
const isSaving = ref(false)
const error = ref('')

const form = reactive({
  work_name: '',
  party_id: '',
  tangible_group_id: '',
  status: 'DRAFT',
  priority: 'NORMAL',
  garment_notes: '',
  service_group_assignments: [
    { service_group_id: '', position_id: '', sequence: 1 },
  ] as Array<{ service_group_id: string; position_id: string; sequence: number }>,
})

function addAssignment() {
  form.service_group_assignments.push({
    service_group_id: '',
    position_id: '',
    sequence: form.service_group_assignments.length + 1,
  })
}

function removeAssignment(index: number) {
  form.service_group_assignments.splice(index, 1)
}

function validate() {
  if (!form.work_name.trim()) {
    error.value = 'El nombre es obligatorio'
    return false
  }
  if (!form.party_id.trim()) {
    error.value = 'El party_id es obligatorio'
    return false
  }
  if (!form.tangible_group_id.trim()) {
    error.value = 'El tangible_group_id es obligatorio'
    return false
  }
  if (form.service_group_assignments.length === 0) {
    error.value = 'Debes agregar al menos una asignación de servicio'
    return false
  }

  const invalid = form.service_group_assignments.find((item) => !item.service_group_id || !item.position_id || item.sequence < 1)
  if (invalid) {
    error.value = 'Todas las asignaciones deben tener service_group_id, position_id y secuencia válida'
    return false
  }

  return true
}

async function submit() {
  error.value = ''
  if (!validate()) return

  isSaving.value = true
  try {
    const created = await mesApi.createWork({
      work_name: form.work_name.trim(),
      party_id: form.party_id.trim(),
      tangible_group_id: form.tangible_group_id.trim(),
      garment_notes: form.garment_notes.trim() || undefined,
      status: form.status,
      priority: form.priority,
      service_group_assignments: form.service_group_assignments.map((item) => ({
        service_group_id: item.service_group_id.trim(),
        position_id: item.position_id.trim(),
        sequence: item.sequence,
      })),
    })

    await router.push(`/mes/works/${created.id}`)
  } catch (err: any) {
    error.value = err.message || 'No se pudo crear el trabajo MES'
  } finally {
    isSaving.value = false
  }
}
</script>

<style scoped>
.dashboard { min-height: 100vh; background-color: #f1f5f9; }
.dashboard-content { max-width: 980px; margin: 0 auto; padding: 2rem; display: flex; flex-direction: column; gap: 1.5rem; }
.page-header { display: flex; justify-content: space-between; align-items: center; gap: 1rem; }
.breadcrumb { font-size: .75rem; text-transform: uppercase; letter-spacing: .08em; color: #64748b; margin: 0; }
.subtitle { color: #64748b; margin: .5rem 0 0; }
.card { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 1rem; }
.form-card { display: flex; flex-direction: column; gap: .75rem; max-width: 760px; }
.label { font-size: .85rem; color: #334155; font-weight: 600; }
.input { border: 1px solid #cbd5e1; border-radius: 8px; padding: .6rem .75rem; font: inherit; }
.tasks-block { border: 1px solid #e2e8f0; border-radius: 10px; padding: .75rem; display: flex; flex-direction: column; gap: .6rem; }
.tasks-header { display: flex; justify-content: space-between; align-items: center; }
.tasks-header h3 { margin: 0; font-size: 1rem; color: #1e293b; }
.assignment-row { display: grid; grid-template-columns: 1fr 1fr 120px auto; gap: .5rem; }
.seq { width: 100%; }
.btn { border: none; border-radius: 8px; padding: .6rem 1rem; cursor: pointer; text-decoration: none; display: inline-flex; align-items: center; justify-content: center; width: fit-content; }
.btn-sm { padding: .4rem .75rem; font-size: .8rem; }
.btn-primary { background: #f4c430; color: #1b3a6b; font-weight: 600; }
.btn-secondary { background: #fff; border: 1px solid #e2e8f0; color: #1e293b; }
.btn-danger { background: #ef4444; color: #fff; }
.alert { background: #fef2f2; color: #b91c1c; border: 1px solid #fecaca; border-radius: 8px; padding: .75rem; }
</style>
