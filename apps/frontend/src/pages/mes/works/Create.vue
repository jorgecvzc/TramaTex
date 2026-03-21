<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">MES / Órdenes de Trabajo</p>
          <h1>Nueva orden de trabajo</h1>
          <p class="subtitle">Selecciona una configuración de trabajo y las tareas se generarán automáticamente.</p>
        </div>
        <RouterLink to="/mes/work-orders" class="btn btn-secondary">Volver</RouterLink>
      </header>

      <section class="card form-card">
        <div v-if="error" class="alert">{{ error }}</div>

        <label class="label">Nombre *</label>
        <input v-model="form.work_name" type="text" class="input" placeholder="Ej: Uniformes Cliente A" />

        <PartySelector
          v-model="form.party_id"
          label="Cliente *"
          placeholder="Buscar cliente por nombre..."
          role-filter="CLIENT"
          :required="true"
          help-text="Selecciona el cliente para filtrar las configuraciones disponibles"
          :disabled="isPrefilledFromOrder"
          @update:modelValue="onPartyChanged"
        />

        <label class="label">Configuración de trabajo *</label>
        <div v-if="isPrefilledFromOrder && prefilledSetupName" class="prefilled-config">
          <span class="prefilled-value">{{ prefilledSetupName }}</span>
          <span class="prefilled-badge">Pre-configurado</span>
        </div>
        <template v-else>
          <select v-model="form.work_setup_id" class="input" :disabled="!form.party_id || isLoadingSetups">
            <option value="">{{ form.party_id ? (isLoadingSetups ? 'Cargando...' : 'Seleccionar configuración') : 'Selecciona primero un cliente' }}</option>
            <option v-for="ws in filteredWorkSetups" :key="ws.id" :value="ws.id">
              {{ ws.name }}
            </option>
          </select>
          <p v-if="form.party_id && filteredWorkSetups.length === 0 && !isLoadingSetups" class="help-text">
            No hay configuraciones activas para este cliente.
          </p>
        </template>

        <label class="label">Prioridad</label>
        <select v-model="form.priority" class="input">
          <option value="LOW">Baja</option>
          <option value="NORMAL">Normal</option>
          <option value="HIGH">Alta</option>
          <option value="URGENT">Urgente</option>
        </select>

        <label class="label">Fecha de vencimiento</label>
        <input v-model="form.due_date" type="date" class="input" />

        <label class="label">Notas</label>
        <textarea v-model="form.notes" class="input" rows="3" placeholder="Observaciones adicionales" />

        <button @click="submit" :disabled="isSaving" class="btn btn-primary">
          {{ isSaving ? 'Guardando...' : 'Crear orden' }}
        </button>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import PartySelector from '@/components/party/PartySelector.vue'
import { mesApi } from '@/services/mesApi'
import type { WorkSetup } from '@/types/mes'

const route = useRoute()
const router = useRouter()
const isSaving = ref(false)
const isLoadingSetups = ref(false)
const error = ref('')
const workSetups = ref<WorkSetup[]>([])
const filteredWorkSetups = ref<WorkSetup[]>([])
const prefilledSetupName = ref('')

const form = reactive({
  work_name: (route.query.name as string) || '',
  party_id: (route.query.partyId as string) || '',
  work_setup_id: (route.query.workSetupId as string) || '',
  priority: 'NORMAL',
  due_date: '',
  notes: (route.query.notes as string) || '',
  order_work_setup_id: (route.query.orderWorkSetupId as string) || '',
})

const isPrefilledFromOrder = computed(() => !!route.query.workSetupId && !!route.query.orderWorkSetupId)

function onPartyChanged(partyId: string) {
  form.work_setup_id = ''
  filteredWorkSetups.value = workSetups.value.filter((ws) => ws.party_id === partyId)
}

function validate() {
  if (!form.work_name.trim()) {
    error.value = 'El nombre es obligatorio'
    return false
  }
  if (!form.party_id.trim()) {
    error.value = 'Debes seleccionar un cliente'
    return false
  }
  if (!form.work_setup_id.trim()) {
    error.value = 'Debes seleccionar una configuración de trabajo'
    return false
  }
  return true
}

async function loadWorkSetups() {
  isLoadingSetups.value = true
  try {
    // If arriving with a pre-filled workSetupId, fetch it directly to guarantee it's available
    if (form.work_setup_id) {
      try {
        const setup = await mesApi.getWorkSetup(form.work_setup_id)
        if (setup) {
          prefilledSetupName.value = setup.name
          workSetups.value = [setup]
          filteredWorkSetups.value = [setup]
        }
      } catch {
        // Setup not found — fall through to load all
      }
    }

    // Load all active setups for the dropdown (in case user changes party)
    const allSetups = await mesApi.listWorkSetups({ is_active: true })

    // Merge: keep the pre-fetched setup even if not in the active list
    const existingIds = new Set(allSetups.map(s => s.id))
    const extras = workSetups.value.filter(s => !existingIds.has(s.id))
    workSetups.value = [...allSetups, ...extras]

    if (form.party_id) {
      filteredWorkSetups.value = workSetups.value.filter((ws) => ws.party_id === form.party_id)
    }
  } catch (err: any) {
    error.value = err.message || 'No se pudieron cargar las configuraciones de trabajo'
  } finally {
    isLoadingSetups.value = false
  }
}

async function submit() {
  error.value = ''
  if (!validate()) return

  isSaving.value = true
  try {
    const created = await mesApi.createWorkOrder({
      work_name: form.work_name.trim(),
      party_id: form.party_id.trim(),
      work_setup_id: form.work_setup_id.trim(),
      notes: form.notes.trim() || undefined,
      priority: form.priority,
      due_date: form.due_date || undefined,
      order_work_setup_id: form.order_work_setup_id.trim() || undefined,
    })

    await router.push(`/mes/work-orders/${created.id}`)
  } catch (err: any) {
    error.value = err.message || 'No se pudo crear la orden de trabajo'
  } finally {
    isSaving.value = false
  }
}

onMounted(loadWorkSetups)
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
.help-text { font-size: .8rem; color: #94a3b8; margin: 0; }
.btn { border: none; border-radius: 8px; padding: .6rem 1rem; cursor: pointer; text-decoration: none; display: inline-flex; align-items: center; justify-content: center; width: fit-content; }
.btn-primary { background: #f4c430; color: #1b3a6b; font-weight: 600; }
.btn-secondary { background: #fff; border: 1px solid #e2e8f0; color: #1e293b; }
.alert { background: #fef2f2; color: #b91c1c; border: 1px solid #fecaca; border-radius: 8px; padding: .75rem; }
.prefilled-config { display: flex; align-items: center; gap: .5rem; padding: .6rem .75rem; border: 1px solid #e2e8f0; border-radius: 8px; background: #f8fafc; }
.prefilled-value { font-size: .875rem; font-weight: 500; color: #1e293b; }
.prefilled-badge { font-size: .7rem; font-weight: 500; padding: .15rem .5rem; border-radius: 9999px; background: #dbeafe; color: #1e40af; }
</style>
