<template>
  <div class="position-form">
    <div class="form-group">
      <label for="name">Nombre de la Posición <span class="required">*</span></label>
      <input 
        id="name"
        v-model="formData.name" 
        type="text" 
        class="form-input"
        :class="{ 'border-error': errors.name }"
        placeholder="Ej: Pecho Izquierdo, Espalda..."
        @input="clearError('name')"
      />
      <span v-if="errors.name" class="error-message">{{ errors.name }}</span>
    </div>

    <div class="form-row mt-4">
      <div class="form-group">
        <label for="code">Código Interno</label>
        <input 
          id="code"
          v-model="formData.code" 
          type="text" 
          class="form-input"
          placeholder="Ej: P-IZQ"
        />
      </div>
      <div class="form-group">
        <label>Estado</label>
        <label class="checkbox-label mt-2">
          <input v-model="formData.is_active" type="checkbox" class="form-checkbox" />
          <span>Posición activa</span>
        </label>
      </div>
    </div>

    <div class="form-group mt-4">
      <label for="description">Notas de la posición</label>
      <textarea 
        id="description"
        v-model="formData.description" 
        class="form-textarea" 
        rows="3"
      ></textarea>
    </div>
  </div>
</template>

<script setup>
import { reactive } from 'vue'

const props = defineProps({
  position: {
    type: Object,
    default: null
  },
  mode: {
    type: String,
    default: 'create',
    validator: (value) => ['create', 'edit'].includes(value)
  }
})

const emit = defineEmits(['submit'])

const formData = reactive({
  name: props.position?.name || '',
  code: props.position?.code || '',
  description: props.position?.description || '',
  is_active: props.position?.is_active ?? true
})

const errors = reactive({
  name: ''
})

function clearError(field) {
  errors[field] = ''
}

function validate() {
  let isValid = true
  if (!formData.name.trim()) {
    errors.name = 'El nombre es obligatorio'
    isValid = false
  }
  return isValid
}

function handleSubmit() {
  if (!validate()) return
  
  const payload = {
    name: formData.name.trim(),
    code: formData.code.trim(),
    description: formData.description.trim(),
    is_active: formData.is_active
  }
  
  if (props.mode === 'edit' && props.position) {
    payload.id = props.position.id
  }
  
  emit('submit', payload)
}

defineExpose({ handleSubmit })
</script>

<style scoped>
.position-form { display: flex; flex-direction: column; gap: 1.5rem; }
.form-group { display: flex; flex-direction: column; gap: 0.5rem; }
.form-group label { font-size: 0.875rem; font-weight: 500; color: #1e293b; }
.form-input, .form-textarea { width: 100%; padding: 0.75rem 1rem; border-radius: 8px; border: 1px solid var(--color-border); font-family: inherit; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
.checkbox-label { display: flex; align-items: center; gap: 0.75rem; cursor: pointer; font-size: 0.9rem; }
.form-checkbox { width: 18px; height: 18px; }
.required { color: #ef4444; }
.mt-4 { margin-top: 1rem; }
.mt-2 { margin-top: 0.5rem; }
</style>
