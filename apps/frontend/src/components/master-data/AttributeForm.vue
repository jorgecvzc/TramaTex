<template>
  <div class="attribute-form">
    <div class="form-group">
      <label for="name">Nombre <span class="required">*</span></label>
      <input 
        id="name"
        v-model="formData.name" 
        type="text" 
        class="form-input"
        placeholder="Ej: Talla, Color, Material"
        @input="clearError('name')"
      />
      <span v-if="errors.name" class="error-message">{{ errors.name }}</span>
    </div>

    <div class="form-group">
      <label for="code">Código <span class="required">*</span></label>
      <input 
        id="code"
        v-model="formData.code" 
        type="text" 
        class="form-input"
        placeholder="Ej: SIZE, COLOR (mayúsculas, sin espacios)"
        @input="handleCodeInput"
      />
      <span v-if="errors.code" class="error-message">{{ errors.code }}</span>
      <small class="hint">Solo letras mayúsculas, números y guiones bajos</small>
    </div>

    <div class="form-group">
      <label for="order">Orden de visualización</label>
      <input 
        id="order"
        v-model.number="formData.order" 
        type="number" 
        class="form-input"
        min="0"
        placeholder="0"
      />
      <small class="hint">Menor número = mayor prioridad</small>
    </div>

    <div class="form-section">
      <h3>Valores del Atributo</h3>
      <p class="section-description">Define los valores posibles para este atributo</p>
      
      <div class="values-list">
        <div v-for="(val, index) in formData.values" :key="index" class="value-item">
          <input 
            v-model="val.value" 
            type="text" 
            class="form-input"
            placeholder="Nombre del valor"
          />
          <input 
            v-model="val.code" 
            type="text" 
            class="form-input"
            placeholder="Código"
            @input="val.code = val.code.toUpperCase()"
          />
          <button 
            type="button" 
            @click="removeValue(index)" 
            class="btn-icon btn-danger"
            :disabled="formData.values.length === 1"
          >
            ✕
          </button>
        </div>
      </div>
      
      <button type="button" @click="addValue" class="btn btn-secondary btn-sm">
        + Agregar valor
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'

const props = defineProps({
  attribute: {
    type: Object,
    default: null
  },
  mode: {
    type: String,
    default: 'create', // 'create' or 'edit'
    validator: (value) => ['create', 'edit'].includes(value)
  }
})

const emit = defineEmits(['submit', 'cancel'])

const formData = reactive({
  name: '',
  code: '',
  order: 0,
  values: [
    { value: '', code: '' }
  ]
})

const errors = reactive({
  name: '',
  code: ''
})

function handleCodeInput(event) {
  formData.code = event.target.value.toUpperCase().replace(/[^A-Z0-9_]/g, '')
  clearError('code')
}

function clearError(field) {
  errors[field] = ''
}

function addValue() {
  formData.values.push({ value: '', code: '' })
}

function removeValue(index) {
  if (formData.values.length > 1) {
    formData.values.splice(index, 1)
  }
}

function validate() {
  let isValid = true
  
  if (!formData.name.trim()) {
    errors.name = 'El nombre es obligatorio'
    isValid = false
  }
  
  if (!formData.code.trim()) {
    errors.code = 'El código es obligatorio'
    isValid = false
  } else if (!/^[A-Z0-9_]+$/.test(formData.code)) {
    errors.code = 'El código solo puede contener mayúsculas, números y guiones bajos'
    isValid = false
  }
  
  return isValid
}

function handleSubmit() {
  if (!validate()) {
    return
  }
  
  // Filter empty values
  const cleanedValues = formData.values.filter(v => v.value.trim() && v.code.trim())
  
  const payload = {
    name: formData.name.trim(),
    code: formData.code.trim(),
    order: formData.order,
    values: cleanedValues
  }
  
  if (props.mode === 'edit' && props.attribute) {
    payload.id = props.attribute.id
  }
  
  emit('submit', payload)
}

onMounted(() => {
  if (props.attribute && props.mode === 'edit') {
    formData.name = props.attribute.name || ''
    formData.code = props.attribute.code || ''
    formData.order = props.attribute.sort_order || 0
    
    if (props.attribute.values && props.attribute.values.length > 0) {
      formData.values = props.attribute.values.map(v => ({
        value: v.value || '',
        code: v.code || '',
        id: v.id // Preserve ID for existing values
      }))
    }
  }
})

defineExpose({
  handleSubmit
})
</script>

<style scoped>
.attribute-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-section {
  border-top: 1px solid #e2e8f0;
  padding-top: 1.5rem;
}

.form-section h3 {
  margin: 0 0 0.5rem 0;
  color: #1e293b;
  font-size: 1rem;
}

.section-description {
  margin: 0 0 1rem 0;
  color: #64748b;
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
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.radio-label {
  margin-left: 0.5rem;
  font-weight: normal;
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

.form-input:disabled {
  background: #f8fafc;
  cursor: not-allowed;
}

.error-message {
  color: #ef4444;
  font-size: 0.75rem;
}

.hint {
  color: #64748b;
  font-size: 0.75rem;
}

.mt-2 {
  margin-top: 0.5rem;
}

.values-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-bottom: 1rem;
}

.value-item {
  display: grid;
  grid-template-columns: 1fr 1fr auto;
  gap: 0.5rem;
  align-items: center;
}

.btn-icon {
  padding: 0.5rem;
  border: none;
  background: none;
  cursor: pointer;
  border-radius: 0.25rem;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  transition: all 0.15s;
}

.btn-icon:hover:not(:disabled) {
  background: #f1f5f9;
}

.btn-icon:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-danger:hover:not(:disabled) {
  background: #fee2e2;
  color: #dc2626;
}

.btn {
  border: none;
  border-radius: 8px;
  padding: 0.6rem 1rem;
  font-size: 0.85rem;
  cursor: pointer;
  transition: background 0.2s ease;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.btn-secondary {
  background: #ffffff;
  color: #1b3a6b;
  border: 1px solid #dde3ed;
  font-weight: 500;
}

.btn-secondary:hover {
  background: #f8fafc;
  border-color: #cbd5e1;
}

.btn-sm {
  padding: 0.4rem 0.75rem;
  font-size: 0.8rem;
}

.checkbox-label {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 0.375rem;
  transition: background 0.15s;
}

.checkbox-label:hover {
  background: #f8fafc;
}

.checkbox-label input[type="checkbox"] {
  margin-top: 0.25rem;
  width: 1.125rem;
  height: 1.125rem;
  cursor: pointer;
}

.scope-selectors {
  margin-top: 1rem;
  padding: 1rem;
  background: #f8fafc;
  border-radius: 0.5rem;
  border: 1px solid #e2e8f0;
}

.scope-hint {
  margin: 0 0 1rem 0;
  color: #475569;
  font-size: 0.875rem;
  font-weight: 500;
}

.info-box {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  padding: 0.875rem;
  border-radius: 0.5rem;
  margin-top: 1rem;
  font-size: 0.875rem;
}

.info-box.warning {
  background: #fef3cd;
  border: 1px solid #f4c430;
  color: #856404;
}

.info-box.info {
  background: #dbeafe;
  border: 1px solid #3b82f6;
  color: #1e40af;
}

.info-icon {
  font-size: 1.125rem;
  flex-shrink: 0;
  margin-top: 0.125rem;
}
</style>
