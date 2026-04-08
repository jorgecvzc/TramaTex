<template>
  <div class="dashboard">
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">MES / Tareas</p>
          <h1>Nueva tarea</h1>
          <p class="subtitle">Crea una tarea para usar en tipos de trabajo.</p>
        </div>
        <RouterLink to="/mes/tasks" class="btn btn-secondary">Volver</RouterLink>
      </header>

      <section class="card form-card">
        <div v-if="error" class="alert">{{ error }}</div>

        <label class="label">Nombre *</label>
        <input v-model="form.name" type="text" class="input" placeholder="Ej: Diseñar" />

        <label class="label">Referencia</label>
        <input v-model="form.reference" type="text" class="input" placeholder="Código de referencia opcional" />

        <label class="label">Descripción</label>
        <textarea v-model="form.description" class="input" rows="4" placeholder="Descripción opcional" />

        <label class="check">
          <input v-model="form.is_active" type="checkbox" />
          Activa
        </label>

        <button @click="submit" :disabled="isSaving" class="btn btn-primary">
          {{ isSaving ? 'Guardando...' : 'Crear tarea' }}
        </button>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { mesApi } from '@/services/mesApi'

const router = useRouter()
const isSaving = ref(false)
const error = ref('')

const form = reactive({
  name: '',
  reference: '',
  description: '',
  is_active: true,
})

async function submit() {
  error.value = ''

  if (!form.name.trim()) {
    error.value = 'El nombre es obligatorio'
    return
  }

  isSaving.value = true
  try {
    await mesApi.createTask({
      name: form.name.trim(),
      reference: form.reference.trim() || undefined,
      description: form.description.trim() || undefined,
      is_active: form.is_active,
    })
    await router.push('/mes/tasks')
  } catch (err: any) {
    error.value = err.message || 'No se pudo crear la tarea'
  } finally {
    isSaving.value = false
  }
}
</script>

<style scoped>
.dashboard { min-height: 100vh; background-color: #f1f5f9; }
.dashboard-content { max-width: 900px; margin: 0 auto; padding: 2rem; display: flex; flex-direction: column; gap: 1.5rem; }
.page-header { display: flex; justify-content: space-between; align-items: center; gap: 1rem; }
.breadcrumb { font-size: .75rem; text-transform: uppercase; letter-spacing: .08em; color: #64748b; margin: 0; }
.subtitle { color: #64748b; margin: .5rem 0 0; }
.card { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 1rem; }
.form-card { display: flex; flex-direction: column; gap: .75rem; max-width: 720px; }
.label { font-size: .85rem; color: #334155; font-weight: 600; }
.input { border: 1px solid #cbd5e1; border-radius: 8px; padding: .6rem .75rem; font: inherit; }
.check { display: inline-flex; align-items: center; gap: .5rem; color: #334155; }
.btn { border: none; border-radius: 8px; padding: .6rem 1rem; cursor: pointer; text-decoration: none; display: inline-flex; align-items: center; justify-content: center; width: fit-content; }
.btn-primary { background: #f4c430; color: #1b3a6b; font-weight: 600; }
.btn-secondary { background: #fff; border: 1px solid #e2e8f0; color: #1e293b; }
.alert { background: #fef2f2; color: #b91c1c; border: 1px solid #fecaca; border-radius: 8px; padding: .75rem; }
</style>
