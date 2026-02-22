<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">MES / Datos Maestros</p>
          <h1>Nuevo grupo de servicio MES</h1>
          <p class="subtitle">Define tareas y secuencia para un servicio de producción.</p>
        </div>
        <RouterLink to="/mes/service-groups" class="btn btn-secondary">Volver</RouterLink>
      </header>

      <section class="card form-card">
        <div v-if="error" class="alert">{{ error }}</div>

        <label class="label">Nombre *</label>
        <input v-model="form.name" type="text" class="input" placeholder="Ej: Serigrafía básica" />

        <label class="label">Descripción</label>
        <textarea v-model="form.description" class="input" rows="3" placeholder="Descripción opcional" />

        <label class="label">Product Group ID (opcional)</label>
        <input v-model="form.product_group_id" type="text" class="input" placeholder="UUID de product group" />

        <label class="check">
          <input v-model="form.is_active" type="checkbox" />
          Activo
        </label>

        <div class="tasks-block">
          <div class="tasks-header">
            <h3>Asignación de tareas</h3>
            <button @click="addAssignment" type="button" class="btn btn-secondary btn-sm">+ Añadir</button>
          </div>

          <div v-for="(assignment, index) in form.task_assignments" :key="index" class="assignment-row">
            <select v-model="assignment.task_id" class="input">
              <option value="">Seleccionar tarea</option>
              <option v-for="task in tasks" :key="task.id" :value="task.id">{{ task.name }}</option>
            </select>

            <input v-model.number="assignment.sequence" type="number" min="1" class="input seq" placeholder="Secuencia" />

            <button @click="removeAssignment(index)" type="button" class="btn btn-danger btn-sm">Quitar</button>
          </div>

          <p v-if="form.task_assignments.length === 0" class="empty">Sin tareas asignadas</p>
        </div>

        <button @click="submit" :disabled="isSaving" class="btn btn-primary">
          {{ isSaving ? 'Guardando...' : 'Crear grupo de servicio' }}
        </button>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { mesApi } from '@/services/mesApi'
import type { MESTask } from '@/types/mes'

const router = useRouter()
const isSaving = ref(false)
const error = ref('')
const tasks = ref<MESTask[]>([])

const form = reactive({
  name: '',
  description: '',
  product_group_id: '',
  is_active: true,
  task_assignments: [] as Array<{ task_id: string; sequence: number }>,
})

function addAssignment() {
  form.task_assignments.push({ task_id: '', sequence: form.task_assignments.length + 1 })
}

function removeAssignment(index: number) {
  form.task_assignments.splice(index, 1)
}

async function loadTasks() {
  try {
    tasks.value = await mesApi.listTasks({ is_active: true })
  } catch (err: any) {
    error.value = err.message || 'No se pudieron cargar las tareas MES'
  }
}

function validate() {
  if (!form.name.trim()) {
    error.value = 'El nombre es obligatorio'
    return false
  }

  const invalid = form.task_assignments.find((item) => !item.task_id || !item.sequence || item.sequence < 1)
  if (invalid) {
    error.value = 'Todas las asignaciones deben tener tarea y secuencia válida'
    return false
  }

  return true
}

async function submit() {
  error.value = ''

  if (!validate()) {
    return
  }

  isSaving.value = true
  try {
    await mesApi.createServiceGroup({
      name: form.name.trim(),
      description: form.description.trim() || undefined,
      product_group_id: form.product_group_id.trim() || undefined,
      is_active: form.is_active,
      task_assignments: form.task_assignments,
    })

    await router.push('/mes/service-groups')
  } catch (err: any) {
    error.value = err.message || 'No se pudo crear el grupo de servicio MES'
  } finally {
    isSaving.value = false
  }
}

onMounted(loadTasks)
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
.check { display: inline-flex; align-items: center; gap: .5rem; color: #334155; }
.tasks-block { border: 1px solid #e2e8f0; border-radius: 10px; padding: .75rem; display: flex; flex-direction: column; gap: .6rem; }
.tasks-header { display: flex; justify-content: space-between; align-items: center; }
.tasks-header h3 { margin: 0; font-size: 1rem; color: #1e293b; }
.assignment-row { display: grid; grid-template-columns: 1fr 140px auto; gap: .5rem; }
.seq { width: 100%; }
.empty { margin: 0; color: #64748b; }
.btn { border: none; border-radius: 8px; padding: .6rem 1rem; cursor: pointer; text-decoration: none; display: inline-flex; align-items: center; justify-content: center; width: fit-content; }
.btn-sm { padding: .4rem .75rem; font-size: .8rem; }
.btn-primary { background: #f4c430; color: #1b3a6b; font-weight: 600; }
.btn-secondary { background: #fff; border: 1px solid #e2e8f0; color: #1e293b; }
.btn-danger { background: #ef4444; color: #fff; }
.alert { background: #fef2f2; color: #b91c1c; border: 1px solid #fecaca; border-radius: 8px; padding: .75rem; }
</style>
