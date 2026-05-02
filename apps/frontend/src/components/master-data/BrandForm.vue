<template>
  <div class="brand-form">
    <div class="form-group">
      <label for="name">Nombre de la marca <span class="required">*</span></label>
      <input 
        id="name"
        v-model="formData.name" 
        type="text" 
        class="form-input"
        :class="{ 'border-error': errors.name }"
        placeholder="Ej: Nike, Adidas"
        @input="clearError('name')"
      />
      <span v-if="errors.name" class="error-message">{{ errors.name }}</span>
    </div>

    <div class="form-group">
      <label for="markupPercentage">Porcentaje de Beneficio (%)</label>
      <input 
        id="markupPercentage"
        v-model.number="formData.defaultMarkupPercentage" 
        type="number" 
        step="0.01"
        min="0"
        class="form-input"
        :class="{ 'border-error': errors.defaultMarkupPercentage }"
        placeholder="Ej: 30.00"
        @input="clearError('defaultMarkupPercentage')"
      />
      <small class="hint">Porcentaje aplicado al costo base de las variantes para calcular el precio de venta base. Ej: 30.00 = 30%</small>
      <span v-if="errors.defaultMarkupPercentage" class="error-message">{{ errors.defaultMarkupPercentage }}</span>
    </div>

    <div class="form-group">
      <label class="checkbox-label">
        <input 
          type="checkbox" 
          v-model="formData.isActive"
        />
        <span>Marca activa</span>
      </label>
      <small class="hint">Las marcas inactivas no se mostrarán en los formularios de productos</small>
    </div>
  </div>
</template>

<script setup>
import { reactive } from 'vue'

const props = defineProps({
  brand: {
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
  name: props.brand?.name || '',
  defaultMarkupPercentage: props.brand?.defaultMarkupPercentage ?? 0,
  isActive: props.brand?.is_active ?? true
})

const errors = reactive({
  name: '',
  defaultMarkupPercentage: ''
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
  
  if (formData.defaultMarkupPercentage < 0) {
    errors.defaultMarkupPercentage = 'El porcentaje de beneficio no puede ser negativo'
    isValid = false
  }
  
  return isValid
}

function handleSubmit() {
  if (!validate()) {
    return
  }
  
  const payload = {
    name: formData.name.trim(),
    defaultMarkupPercentage: formData.defaultMarkupPercentage,
    isActive: formData.isActive
  }
  
  if (props.mode === 'edit' && props.brand) {
    payload.id = props.brand.id
  }
  
  emit('submit', payload)
}

defineExpose({
  handleSubmit
})
</script>

<style scoped>
.brand-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.info-box {
  display: flex;
  gap: 1rem;
  padding: 1rem;
  border-radius: 0.5rem;
  border: 1px solid;
}

.info-box.warning {
  background: #fef3c7;
  border-color: #fde68a;
  color: #92400e;
}

.info-icon {
  font-size: 1.25rem;
  flex-shrink: 0;
}

.info-box strong {
  display: block;
  margin-bottom: 0.25rem;
}

.info-box p {
  margin: 0;
  font-size: 0.875rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #1e293b;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: normal;
  cursor: pointer;
}

.required {
  color: #ef4444;
}

.form-input {
  padding: 0.625rem;
  border: 1px solid #e2e8f0;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  transition: border-color 0.15s;
}

.form-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.error-message {
  color: #ef4444;
  font-size: 0.75rem;
}

.hint {
  color: #64748b;
  font-size: 0.75rem;
}

.mb-3 {
  margin-bottom: 1.5rem;
}
</style>
