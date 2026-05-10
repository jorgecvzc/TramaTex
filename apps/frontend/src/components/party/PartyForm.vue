<script setup lang="ts">
import { ref, reactive, onMounted, watch, computed } from 'vue'
import { Building2, User, Percent, Save, RefreshCw, X, RotateCcw } from 'lucide-vue-next'
import FormSection from '@/components/shared/FormSection.vue'
import { partyApi } from '@/services/partyApi'
import { useToastStore } from '@/stores/toast'

interface Props {
  partyId?: string
  initialData?: any
}

const props = defineProps<Props>()
const emit = defineEmits(['submit', 'update', 'cancel'])

const toastStore = useToastStore()
const isSaving = ref(false)
const errors = reactive<Record<string, string>>({
  name: '',
  taxId: '',
  website: ''
})

const formData = reactive({
  name: '',
  firstName: '', 
  lastName: '',  
  type: '' as 'ORGANIZATION' | 'PERSON' | '', 
  taxId: '',
  taxIdType: '',
  status: 'ACTIVE',
  role: '', 
  notes: '',
  website: '',
  defaultDiscountPercentage: 0
})

// --- Document Filtering ---
const taxIdOptions = computed(() => {
  if (formData.type === 'ORGANIZATION') {
    return [
      { value: 'CIF', label: 'CIF (España)' },
      { value: 'VAT', label: 'Número VAT (UE)' },
      { value: 'OTHER', label: 'Otros' }
    ]
  } else if (formData.type === 'PERSON') {
    return [
      { value: 'NIF', label: 'NIF / DNI' },
      { value: 'NIE', label: 'NIE' },
      { value: 'RESIDENT_CARD', label: 'Tarjeta de Residente' },
      { value: 'PASSPORT', label: 'Pasaporte' },
      { value: 'VAT', label: 'VAT Personal' },
      { value: 'OTHER', label: 'Otros' }
    ]
  }
  return []
})

function resetForm() {
  Object.assign(formData, {
    name: '', firstName: '', lastName: '', type: '', taxId: '', taxIdType: '', status: 'ACTIVE',
    role: '', notes: '', website: '', defaultDiscountPercentage: 0
  })
  Object.keys(errors).forEach(key => errors[key] = '')
}

function populateForm(data: any) {
  if (!data) return
  
  formData.type = (data.entityType || data.type || '') as any
  formData.name = data.name || ''
  
  // Handle names for PERSON
  formData.firstName = data.firstName || data.first_name || ''
  formData.lastName = data.lastName || data.last_name || ''
  
  if (formData.type === 'PERSON' && formData.name && !formData.firstName) {
    const parts = (formData.name as string).trim().split(/\s+/)
    formData.firstName = parts[0] || ''
    formData.lastName = parts.slice(1).join(' ') || ''
  }

  formData.taxId = data.taxId || data.tax_id || ''
  formData.taxIdType = data.taxIdType || data.tax_id_type || (formData.type === 'ORGANIZATION' ? 'CIF' : 'NIF')
  formData.status = data.status || 'ACTIVE'
  formData.role = data.role || ''
  formData.notes = data.notes || ''
  formData.website = data.website || ''
  formData.defaultDiscountPercentage = data.defaultDiscountPercentage || data.default_discount_percentage || 0
}

onMounted(() => {
  if (props.initialData) {
    populateForm(props.initialData)
  }
  window.addEventListener('tramatex-save', handleGlobalSave)
  window.addEventListener('tramatex-esc', handleGlobalEsc)
})

onBeforeUnmount(() => {
  window.removeEventListener('tramatex-save', handleGlobalSave)
  window.removeEventListener('tramatex-esc', handleGlobalEsc)
})

function handleGlobalSave() {
  if (!isSaving.value) handleSubmit()
}

function handleGlobalEsc() {
  emit('cancel')
}

watch(() => props.initialData, (newVal) => {
  if (newVal) populateForm(newVal)
}, { deep: true })

// --- Validation ---
function validateName() {
  if (formData.type === 'ORGANIZATION') {
    if (!formData.name.trim()) {
      errors.name = 'El nombre es obligatorio'
      return false
    }
    if (formData.name.trim().length < 3) {
      errors.name = 'Mínimo 3 caracteres'
      return false
    }
  } else if (formData.type === 'PERSON') {
    if (!formData.firstName.trim()) {
      errors.name = 'El nombre es obligatorio'
      return false
    }
    if (!formData.lastName.trim()) {
      errors.name = 'Los apellidos son obligatorios'
      return false
    }
  }
  errors.name = ''
  return true
}

function validateTaxId() {
  if (formData.taxId && formData.taxId.trim().length < 4) {
    errors.taxId = 'Formato inválido'
    return false
  }
  errors.taxId = ''
  return true
}

function validateWebsite() {
  if (formData.website && !formData.website.startsWith('http')) {
    errors.website = 'URL inválido'
    return false
  }
  errors.website = ''
  return true
}

function validateForm() {
  const isNameValid = validateName()
  const isTaxValid = validateTaxId()
  const isWebValid = validateWebsite()
  
  if (!formData.role) {
    toastStore.warning('Selecciona un rol para la entidad')
    return false
  }
  
  if (!formData.type) {
    toastStore.warning('Selecciona el tipo de entidad')
    return false
  }

  return isNameValid && isTaxValid && isWebValid
}

async function handleSubmit() {
  if (!validateForm()) {
    toastStore.error('Corrige los errores antes de continuar')
    return
  }

  isSaving.value = true
  try {
    const finalName = formData.type === 'PERSON' 
      ? `${formData.firstName} ${formData.lastName}`.trim()
      : formData.name

    const payload = {
      name: finalName,
      firstName: formData.firstName,
      lastName: formData.lastName,
      type: formData.type as any,
      entityType: formData.type as any,
      role: formData.role as any,
      taxId: formData.taxId,
      taxIdType: formData.taxIdType as any,
      status: formData.status as any,
      notes: formData.notes,
      website: formData.website,
      defaultDiscountPercentage: formData.defaultDiscountPercentage,
      default_discount_percentage: formData.defaultDiscountPercentage, // Sync names
      hasPerson: formData.type === 'PERSON'
    }

    if (props.partyId) {
      await partyApi.updateParty(props.partyId, payload)
      toastStore.success('Cambios guardados con éxito')
      emit('update', payload)
    } else {
      const result = await partyApi.createParty(payload)
      toastStore.success('Entidad creada correctamente')
      emit('submit', result)
    }
  } catch (err: any) {
    toastStore.error(err.message || 'Error al guardar la entidad')
  } finally {
    isSaving.value = false
  }
}

// Watch for type changes to adjust taxIdType defaults
watch(() => formData.type, (newType) => {
  if (newType === 'ORGANIZATION') formData.taxIdType = 'CIF'
  else if (newType === 'PERSON') formData.taxIdType = 'NIF'
})
</script>

<template>
  <form @submit.prevent="handleSubmit" class="party-form">
    <header class="form-header-box mb-6">
      <h2 class="text-xl font-bold flex items-center gap-2">
        <Building2 v-if="formData.type === 'ORGANIZATION'" :size="24" />
        <User v-else :size="24" />
        {{ props.partyId ? 'Editar entidad' : 'Crear entidad' }}
      </h2>
    </header>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- SECTION: IDENTITY -->
      <FormSection title="Identidad" icon="person">
        <div class="form-group mb-4">
          <label for="entity-type">Tipo de entidad *</label>
          <select id="entity-type" v-model="formData.type" class="form-input" required :disabled="!!props.partyId">
            <option value="">-- Selecciona tipo --</option>
            <option value="ORGANIZATION">Empresa / Jurídica</option>
            <option value="PERSON">Persona Física</option>
          </select>
        </div>

        <div class="form-group mb-4">
          <label for="party-role">Rol de la entidad *</label>
          <select id="party-role" v-model="formData.role" class="form-input" required>
            <option value="">-- Selecciona rol --</option>
            <option value="CLIENT">Cliente</option>
            <option value="SUPPLIER">Proveedor</option>
            <option value="BOTH">Cliente y proveedor</option>
            <option value="CONTACT">Contacto</option>
          </select>
        </div>

        <div class="animate-fade-in" v-if="formData.type">
          <!-- Organization Name -->
          <div v-if="formData.type === 'ORGANIZATION'" class="form-group mb-4">
            <label for="party-name">Nombre de la organización *</label>
            <input 
              id="party-name"
              v-model="formData.name" 
              type="text" 
              class="form-input" 
              :class="{ 'is-invalid': errors.name && formData.type === 'ORGANIZATION' }"
              required 
              @blur="validateName"
              placeholder="p. ej. Acme Corp"
            />
            <span v-if="errors.name && formData.type === 'ORGANIZATION'" class="error-msg">{{ errors.name }}</span>
          </div>

          <!-- Person Name & Last Name -->
          <div v-else-if="formData.type === 'PERSON'" class="form-row-names mb-4">
            <div class="form-group">
              <label for="first-name">Nombre *</label>
              <input 
                id="first-name"
                v-model="formData.firstName" 
                type="text" 
                class="form-input" 
                :class="{ 'is-invalid': errors.name && !formData.firstName }"
                required 
                @blur="validateName"
                placeholder="p. ej. Juan"
              />
            </div>
            <div class="form-group">
              <label for="last-name">Apellidos *</label>
              <input 
                id="last-name"
                v-model="formData.lastName" 
                type="text" 
                class="form-input" 
                :class="{ 'is-invalid': errors.name && !formData.lastName }"
                required 
                @blur="validateName"
                placeholder="p. ej. Pérez"
              />
            </div>
            <span v-if="errors.name && formData.type === 'PERSON'" class="error-msg full-width-msg">{{ errors.name }}</span>
          </div>

          <div class="form-row-tax">
            <div class="form-group">
              <label for="tax-id-type">Tipo identificación</label>
              <select id="tax-id-type" v-model="formData.taxIdType" class="form-input">
                <option v-for="opt in taxIdOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
              </select>
            </div>
            <div class="form-group">
              <label for="tax-id">Número</label>
              <input 
                id="tax-id"
                v-model="formData.taxId" 
                type="text" 
                class="form-input" 
                :class="{ 'is-invalid': errors.taxId }"
                placeholder="p. ej., 12345678A"
                @blur="validateTaxId"
              />
              <span v-if="errors.taxId" class="error-msg">{{ errors.taxId }}</span>
            </div>
          </div>
        </div>
      </FormSection>

      <!-- SECTION: CONFIG -->
      <FormSection title="Configuración" icon="settings">
        <div class="form-group mb-4">
          <label for="website">Sitio web</label>
          <input 
            id="website"
            v-model="formData.website" 
            type="url" 
            class="form-input" 
            :class="{ 'is-invalid': errors.website }"
            placeholder="https://..."
            @blur="validateWebsite"
          />
          <span v-if="errors.website" class="error-msg">{{ errors.website }}</span>
        </div>

        <div class="form-group mb-4">
          <label for="notes">Notas</label>
          <textarea id="notes" v-model="formData.notes" class="form-textarea" rows="3"></textarea>
        </div>

        <div class="form-group" v-if="formData.role !== 'SUPPLIER'">
          <label>Bonificación Comercial (%)</label>
          <div class="input-with-icon">
            <Percent :size="18" class="icon-start" />
            <input v-model.number="formData.defaultDiscountPercentage" type="number" step="0.01" class="form-input" />
          </div>
        </div>
      </FormSection>
    </div>

    <div class="form-actions mt-8">
      <button type="button" class="btn btn-outline" @click="resetForm">
        <RotateCcw :size="18" />
        <span>Reiniciar</span>
      </button>
      <div class="flex-1"></div>
      <button type="button" class="btn btn-ghost mr-2" @click="$emit('cancel')">Cancelar</button>
      <button type="submit" class="btn btn-secondary" :disabled="isSaving">
        <component :is="isSaving ? RefreshCw : Save" :size="18" :class="{ 'spin': isSaving }" />
        <span>{{ props.partyId ? 'Actualizar entidad' : 'Crear entidad' }}</span>
      </button>
    </div>
  </form>
</template>

<style scoped>
.party-form { background: white; padding: 2rem; border-radius: 12px; box-shadow: var(--box-shadow-md); border: 1px solid var(--color-border); }
.form-header-box { border-bottom: 2px solid var(--color-background); padding-bottom: 1rem; color: var(--color-secondary); }

.form-row-names { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
.form-row-tax { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }

.form-group label { display: block; font-size: 0.75rem; font-weight: 700; text-transform: uppercase; color: var(--color-text-secondary); margin-bottom: 0.5rem; }
.form-input, .form-textarea { width: 100%; padding: 0.75rem 1rem; border-radius: 8px; border: 1px solid var(--color-border); font-family: inherit; transition: 0.2s; }
.form-input:focus, .form-textarea:focus { outline: none; border-color: var(--color-primary); box-shadow: 0 0 0 3px rgba(230, 184, 0, 0.1); }
.form-input.is-invalid { border-color: var(--color-error); }
.error-msg { font-size: 0.7rem; color: var(--color-error); font-weight: 600; margin-top: 0.25rem; display: block; }
.full-width-msg { grid-column: 1 / span 2; }

.input-with-icon { position: relative; display: flex; align-items: center; }
.icon-start { position: absolute; left: 0.75rem; color: var(--color-text-secondary); }
.input-with-icon input { padding-left: 2.5rem; }

.form-actions { display: flex; align-items: center; gap: 1rem; border-top: 1px solid var(--color-border); padding-top: 2rem; }

.animate-fade-in { animation: fadeIn 0.4s ease-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
.spin { animation: spin 1s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

@media (max-width: 600px) {
  .form-row-names, .form-row-tax { grid-template-columns: 1fr; }
}
</style>
