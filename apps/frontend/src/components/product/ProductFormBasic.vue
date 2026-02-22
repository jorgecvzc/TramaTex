<template>
  <div class="form-step">
    <h3 class="step-title">Información básica</h3>
    <p class="step-description">
      Define los datos principales del producto
    </p>

    <form @submit.prevent="handleNext" class="step-form">
      <!-- Product Type -->
      <div class="form-group">
        <label for="productType">
          Tipo de producto *
          <span class="required">obligatorio</span>
        </label>
        <select
          id="productType"
          v-model="localData.productType"
          required
          @change="validateField('productType')"
        >
          <option value="">-- Selecciona el tipo --</option>
          <option value="TANGIBLE">Tangible (Producto físico)</option>
          <option value="SERVICE">Servicio</option>
        </select>
        <span v-if="errors.productType" class="error">{{ errors.productType }}</span>
        <span class="hint">
          Los productos tangibles son artículos físicos. Los servicios son actividades prestadas.
        </span>
      </div>

      <!-- SKU -->
      <div class="form-group">
        <label for="sku">
          SKU (Código del producto) *
          <span class="required">obligatorio</span>
        </label>
        <input
          id="sku"
          v-model="localData.sku"
          type="text"
          placeholder="ej. FYR2040"
          required
          @blur="validateField('sku')"
          @input="handleSkuInput"
        />
        <span v-if="errors.sku" class="error">{{ errors.sku }}</span>
        <span class="hint">
          Código único que identifica este producto. Usa solo letras y números (sin espacios).
        </span>
      </div>

      <!-- Name -->
      <div class="form-group">
        <label for="name">
          Nombre *
          <span class="required">obligatorio</span>
        </label>
        <input
          id="name"
          v-model="localData.name"
          type="text"
          placeholder="ej. Camiseta Clásica"
          required
          @blur="validateField('name')"
        />
        <span v-if="errors.name" class="error">{{ errors.name }}</span>
        <span class="hint">
          Nombre corto para tickets y listas.
        </span>
      </div>

      <!-- Long Name -->
      <div class="form-group">
        <label for="longName">
          Nombre completo
        </label>
        <input
          id="longName"
          v-model="localData.longName"
          type="text"
          placeholder="ej. Camiseta Clásica de Algodón Orgánico"
          @blur="validateField('longName')"
        />
        <span v-if="errors.longName" class="error">{{ errors.longName }}</span>
        <span class="hint">
          Nombre extendido para facturas y presupuestos (opcional).
        </span>
      </div>

      <!-- Description -->
      <div class="form-group">
        <label for="description">
          Descripción
        </label>
        <textarea
          id="description"
          v-model="localData.description"
          placeholder="Describe las características principales del producto..."
          rows="4"
        />
        <span class="hint">
          Información adicional que ayude a identificar y describir el producto.
        </span>
      </div>

      <!-- Base Price -->
      <div class="form-group">
        <label for="basePrice">
          Precio base (coste) *
          <span class="required">obligatorio</span>
        </label>
        <input
          id="basePrice"
          v-model.number="localData.basePrice"
          type="number"
          step="0.01"
          min="0"
          placeholder="ej. 25.50"
          required
          @blur="validateField('basePrice')"
        />
        <span v-if="errors.basePrice" class="error">{{ errors.basePrice }}</span>
        <span class="hint">
          Coste base del producto (usado para calcular precios de venta). Debe ser mayor o igual a 0.
        </span>
      </div>

      <!-- Form Actions -->
      <div class="form-actions">
        <button
          type="button"
          @click="$emit('cancel')"
          class="btn btn-secondary"
        >
          Cancelar
        </button>
        <button
          type="submit"
          class="btn btn-primary"
          :disabled="!isStepValid"
        >
          Siguiente: Clasificación →
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { reactive, computed, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['update:modelValue', 'next', 'cancel'])

// Local copy of data
const localData = reactive({
  productType: props.modelValue.productType || '',
  sku: props.modelValue.sku || '',
  name: props.modelValue.name || '',
  longName: props.modelValue.longName || '',
  description: props.modelValue.description || '',
  basePrice: props.modelValue.basePrice !== undefined ? props.modelValue.basePrice : 0,
})

// Validation errors
const errors = reactive({
  productType: '',
  sku: '',
  name: '',
  longName: '',
  basePrice: '',
})

// Watch local data changes and emit to parent
watch(localData, (newValue) => {
  emit('update:modelValue', { ...newValue })
}, { deep: true })

// Validation rules
const validationRules = {
  productType: (value) => {
    if (!value) {
      return 'El tipo de producto es obligatorio'
    }
    if (!['TANGIBLE', 'SERVICE'].includes(value)) {
      return 'Tipo de producto inválido'
    }
    return ''
  },
  sku: (value) => {
    if (!value || value.trim().length === 0) {
      return 'El SKU es obligatorio'
    }
    if (value.length < 3) {
      return 'El SKU debe tener al menos 3 caracteres'
    }
    if (value.length > 20) {
      return 'El SKU no debe superar 20 caracteres'
    }
    if (!/^[A-Z0-9]+$/.test(value)) {
      return 'El SKU solo puede contener letras mayúsculas y números'
    }
    return ''
  },
  name: (value) => {
    if (!value || value.trim().length === 0) {
      return 'El nombre es obligatorio'
    }
    if (value.length < 3) {
      return 'El nombre debe tener al menos 3 caracteres'
    }
    if (value.length > 100) {
      return 'El nombre no debe superar 100 caracteres'
    }
    return ''
  },
  longName: (value) => {
    if (value && value.length > 200) {
      return 'El nombre completo no debe superar 200 caracteres'
    }
    return ''
  },
  basePrice: (value) => {
    if (value === undefined || value === null || value === '') {
      return 'El precio base es obligatorio'
    }
    const numValue = parseFloat(value)
    if (isNaN(numValue)) {
      return 'El precio base debe ser un número válido'
    }
    if (numValue < 0) {
      return 'El precio base no puede ser negativo'
    }
    return ''
  },
}

// Validate a single field
function validateField(fieldName) {
  if (validationRules[fieldName]) {
    errors[fieldName] = validationRules[fieldName](localData[fieldName])
  }
}

// Validate all fields
function validateAll() {
  let isValid = true
  for (const field in validationRules) {
    validateField(field)
    if (errors[field]) {
      isValid = false
    }
  }
  return isValid
}

// Auto-uppercase SKU
function handleSkuInput(event) {
  localData.sku = event.target.value.toUpperCase()
}

// Check if step is valid
const isStepValid = computed(() => {
  return (
    localData.productType !== '' &&
    localData.sku !== '' &&
    localData.name !== '' &&
    localData.basePrice !== undefined &&
    localData.basePrice !== null &&
    localData.basePrice !== '' &&
    !errors.productType &&
    !errors.sku &&
    !errors.name &&
    !errors.longName &&
    !errors.basePrice
  )
})

// Handle next step
function handleNext() {
  if (validateAll() && isStepValid.value) {
    emit('next')
  }
}
</script>

<style scoped>
.form-step {
  max-width: 720px;
  margin: 0 auto;
}

.step-title {
  color: #1b3a6b;
  font-size: 1.5rem;
  margin: 0 0 0.5rem;
}

.step-description {
  color: #64748b;
  margin: 0 0 2rem;
  font-size: 0.95rem;
}

.step-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  font-weight: 600;
  color: #1e293b;
  font-size: 0.9rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.required {
  font-size: 0.75rem;
  font-weight: 400;
  color: #94a3b8;
  text-transform: lowercase;
}

.form-group input,
.form-group select,
.form-group textarea {
  padding: 0.65rem 0.85rem;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-size: 0.95rem;
  font-family: inherit;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #1b3a6b;
  box-shadow: 0 0 0 3px rgba(27, 58, 107, 0.1);
}

.form-group textarea {
  resize: vertical;
  min-height: 100px;
}

.hint {
  font-size: 0.8rem;
  color: #64748b;
  margin-top: -0.25rem;
}

.error {
  color: #ef4444;
  font-size: 0.8rem;
  margin-top: -0.25rem;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  padding-top: 1rem;
  border-top: 1px solid #e2e8f0;
  margin-top: 1rem;
}

.btn {
  border: none;
  border-radius: 8px;
  padding: 0.65rem 1.25rem;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background: #f4d03f;
  color: #1e293b;
  box-shadow: 0 2px 4px rgba(244, 208, 63, 0.3);
}

.btn-primary:not(:disabled):hover {
  background: #f0c929;
  box-shadow: 0 4px 8px rgba(244, 208, 63, 0.4);
  transform: translateY(-1px);
}

.btn-primary:not(:disabled):active {
  transform: translateY(0);
}

.btn-secondary {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  color: #1e293b;
}

.btn-secondary:hover {
  background: #f8fafc;
}

@media (max-width: 768px) {
  .form-step {
    max-width: 100%;
  }

  .form-actions {
    flex-direction: column-reverse;
  }

  .btn {
    width: 100%;
  }
}
</style>
