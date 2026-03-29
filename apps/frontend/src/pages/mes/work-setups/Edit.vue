<template>
  <div class="dashboard">
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">MES / Configuraciones de Cliente</p>
          <h1>Editar configuración de cliente</h1>
          <p class="subtitle">Modifica la asignación de tipos de trabajo y posiciones.</p>
        </div>
        <RouterLink to="/mes/work-setups" class="btn btn-secondary">Volver</RouterLink>
      </header>

      <section class="card form-card">
        <div v-if="isLoading" class="empty-state">Cargando configuración...</div>
        <template v-else>
          <div v-if="error" class="alert">{{ error }}</div>

          <label class="label">Nombre *</label>
          <input v-model="form.name" type="text" class="input" placeholder="Ej: Uniformes Empresa XYZ" />

          <PartySelector
            v-model="form.party_id"
            label="Cliente *"
            placeholder="Buscar cliente por nombre..."
            role-filter="CLIENT"
            :required="true"
            help-text="Selecciona el cliente"
          />

          <label class="label">Categoría tangible *</label>
          <select v-model="form.tangible_group_id" class="input">
            <option value="">Seleccionar categoría tangible</option>
            <option v-for="group in tangibleGroups" :key="group.id" :value="group.id">
              {{ group.name }}
            </option>
          </select>

          <label class="label">Descripción</label>
          <textarea v-model="form.description" class="input" rows="3" placeholder="Descripción opcional" />

          <label class="check">
            <input v-model="form.is_active" type="checkbox" />
            Activa
          </label>

          <div class="lines-block">
            <div class="lines-header">
              <h3>Líneas de configuración</h3>
              <button @click="addLine" type="button" class="btn btn-secondary btn-sm">+ Añadir línea</button>
            </div>

            <div v-for="(line, index) in form.lines" :key="index" class="line-row">
              <select v-model="line.work_type_id" class="input">
                <option value="">Tipo de trabajo</option>
                <option v-for="wt in workTypes" :key="wt.id" :value="wt.id">
                  {{ wt.name }}
                </option>
              </select>
              <select v-model="line.position_id" class="input">
                <option value="">Posición</option>
                <option v-for="pos in positions" :key="pos.id" :value="pos.id">
                  {{ pos.name }} ({{ pos.code }})
                </option>
              </select>
              <input v-model.number="line.sequence" type="number" min="1" class="input seq" placeholder="Seq" />
              <input
                v-model="line.design_file_path"
                type="text"
                class="input file-path-input"
                placeholder="Ruta del archivo (opcional)"
                title="Ruta completa del archivo de diseño, p.ej. C:\Diseños\logo.ai"
              />
              <button type="button" class="btn btn-secondary btn-sm file-pick-btn" title="Seleccionar archivo" @click="pickFile(index)">📂</button>
              <input :ref="el => { fileInputRefs[index] = el as HTMLInputElement }" type="file" class="file-input-hidden" @change="onFileSelected($event, index)" />
              <button @click="removeLine(index)" type="button" class="btn btn-danger btn-sm">Quitar</button>
            </div>

            <p v-if="form.lines.length === 0" class="empty-state">Sin líneas.</p>
          </div>

          <button @click="submit" :disabled="isSaving" class="btn btn-primary">
            {{ isSaving ? 'Guardando...' : 'Guardar cambios' }}
          </button>
        </template>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import PartySelector from '@/components/party/PartySelector.vue'
import { mesApi } from '@/services/mesApi'
import { productApi } from '@/services/productApi'
import type { MESPosition, MESWorkType } from '@/types/mes'

type ProductGroupOption = { id: string; name: string; type: string }

const route = useRoute()
const router = useRouter()
const isLoading = ref(true)
const isSaving = ref(false)
const error = ref('')
const workTypes = ref<MESWorkType[]>([])
const positions = ref<MESPosition[]>([])
const productGroups = ref<ProductGroupOption[]>([])

const tangibleGroups = computed(() => productGroups.value.filter(g => g.type === 'TANGIBLE'))

const form = reactive({
  name: '',
  party_id: '',
  tangible_group_id: '',
  description: '',
  is_active: true,
  lines: [] as Array<{ work_type_id: string; position_id: string; sequence: number; design_file_path: string }>,
})

function addLine() {
  form.lines.push({ work_type_id: '', position_id: '', sequence: form.lines.length + 1, design_file_path: '' })
}

function removeLine(index: number) {
  form.lines.splice(index, 1)
}

const fileInputRefs = ref<HTMLInputElement[]>([])

function pickFile(index: number) {
  fileInputRefs.value[index]?.click()
}

function onFileSelected(event: Event, index: number) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (file) {
    form.lines[index].design_file_path = file.name
  }
}

function validate(): boolean {
  if (!form.name.trim()) { error.value = 'El nombre es obligatorio'; return false }
  if (!form.party_id.trim()) { error.value = 'El cliente es obligatorio'; return false }
  if (!form.tangible_group_id.trim()) { error.value = 'La categoría tangible es obligatoria'; return false }
  const invalid = form.lines.find(l => !l.work_type_id || !l.position_id || l.sequence < 1)
  if (invalid) { error.value = 'Todas las líneas deben tener tipo de trabajo, posición y secuencia válida'; return false }
  return true
}

async function loadData() {
  try {
    const [setup, wt, pos, groups] = await Promise.all([
      mesApi.getWorkSetup(route.params.id as string),
      mesApi.listWorkTypes({ is_active: true }),
      mesApi.listPositions({ is_active: true }),
      productApi.listProductGroups({ isActive: true }),
    ])

    workTypes.value = wt
    positions.value = pos
    productGroups.value = groups.data

    form.name = setup.name
    form.party_id = setup.party_id
    form.tangible_group_id = setup.tangible_group_id
    form.description = setup.description || ''
    form.is_active = setup.is_active
    form.lines = (setup.lines || []).map(l => ({
      work_type_id: l.work_type_id,
      position_id: l.position_id,
      sequence: l.sequence,
      design_file_path: l.design_file_path || '',
    }))
  } catch (err: any) {
    error.value = err.message || 'No se pudo cargar la configuración'
  } finally {
    isLoading.value = false
  }
}

async function submit() {
  error.value = ''
  if (!validate()) return

  isSaving.value = true
  try {
    await mesApi.updateWorkSetup(route.params.id as string, {
      name: form.name.trim(),
      party_id: form.party_id.trim(),
      tangible_group_id: form.tangible_group_id.trim(),
      description: form.description.trim() || undefined,
      is_active: form.is_active,
      lines: form.lines.map(l => ({
        work_type_id: l.work_type_id,
        position_id: l.position_id,
        sequence: l.sequence,
        design_file_path: l.design_file_path || undefined,
      })),
    })
    await router.push('/mes/work-setups')
  } catch (err: any) {
    error.value = err.message || 'No se pudo actualizar la configuración'
  } finally {
    isSaving.value = false
  }
}

onMounted(loadData)
</script>

<style scoped>
.dashboard { min-height: 100vh; background-color: #f1f5f9; }
.dashboard-content { max-width: 1200px; margin: 0 auto; padding: 2rem; display: flex; flex-direction: column; gap: 1.5rem; }
.page-header { display: flex; justify-content: space-between; align-items: center; gap: 1rem; }
.breadcrumb { font-size: .75rem; text-transform: uppercase; letter-spacing: .08em; color: #64748b; margin: 0; }
.subtitle { color: #64748b; margin: .5rem 0 0; }
.card { background: #fff; border: 1px solid #e2e8f0; border-radius: 12px; padding: 1rem; }
.form-card { display: flex; flex-direction: column; gap: .75rem; }
.label { font-size: .85rem; color: #334155; font-weight: 600; }
.input { border: 1px solid #cbd5e1; border-radius: 8px; padding: .6rem .75rem; font: inherit; }
.check { display: inline-flex; align-items: center; gap: .5rem; color: #334155; }
.btn { border: none; border-radius: 8px; padding: .6rem 1rem; cursor: pointer; text-decoration: none; display: inline-flex; align-items: center; justify-content: center; width: fit-content; }
.btn-primary { background: #f4c430; color: #1b3a6b; font-weight: 600; }
.btn-secondary { background: #fff; border: 1px solid #e2e8f0; color: #1e293b; }
.btn-sm { font-size: .8rem; padding: .35rem .65rem; border-radius: 6px; }
.btn-danger { background: #fef2f2; color: #b91c1c; border: 1px solid #fecaca; }
.lines-block { border: 1px solid #e2e8f0; border-radius: 8px; padding: 1rem; display: flex; flex-direction: column; gap: .75rem; }
.lines-header { display: flex; justify-content: space-between; align-items: center; }
.lines-header h3 { margin: 0; font-size: .95rem; color: #1e293b; }
.line-row { display: grid; grid-template-columns: 1fr 1fr 70px 1fr auto auto auto; gap: .5rem; align-items: center; }
.file-pick-btn { padding: .4rem .6rem; font-size: .85rem; }
.file-input-hidden { display: none; }
.file-path-input { font-size: .8rem; }
.seq { text-align: center; }
.empty-state { text-align: center; color: #64748b; padding: .5rem; font-size: .85rem; }
.alert { background: #fef2f2; color: #b91c1c; border: 1px solid #fecaca; border-radius: 8px; padding: .75rem; }
</style>
