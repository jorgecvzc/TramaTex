<template>
  <div class="form-step">
    <h3 class="step-title">Clasificación</h3>
    <p class="step-description">
      Asigna la marca y las categorías del producto
    </p>

    <form @submit.prevent="handleNext" class="step-form">
      <!-- Brand Selection -->
      <div class="form-group">
        <label for="brandId">
          Marca *
          <span class="required">obligatorio</span>
        </label>
        <select
          id="brandId"
          v-model="localData.brandId"
          required
          :disabled="isLoadingBrands"
          @change="validateField('brandId')"
        >
          <option value="">
            {{ isLoadingBrands ? 'Cargando marcas...' : '-- Selecciona una marca --' }}
          </option>
          <option
            v-for="brand in brands"
            :key="brand.id"
            :value="brand.id"
          >
            {{ brand.name }}
          </option>
        </select>
        <span v-if="errors.brandId" class="error">{{ errors.brandId }}</span>
        <span v-if="loadError.brands" class="error">{{ loadError.brands }}</span>
        <span class="hint">
          La marca ayuda a categorizar el producto y define el alcance de atributos y precios.
        </span>
      </div>

      <!-- Product Groups Selection -->
      <div class="form-group">
        <label>
          Categorías
          <span class="optional">(opcional)</span>
        </label>
        <div v-if="isLoadingGroups" class="loading-state">
          Cargando categorías...
        </div>
        <div v-else-if="loadError.groups" class="error">
          {{ loadError.groups }}
        </div>
        <div v-else-if="groups.length === 0" class="empty-state">
          No hay categorías disponibles. Puedes crearlas más tarde.
        </div>
        <div v-else class="checkbox-group">
          <label
            v-for="group in groups"
            :key="group.id"
            class="checkbox-label"
          >
            <input
              type="checkbox"
              :value="group.id"
              v-model="localData.groupIds"
            />
            <span class="checkbox-text">
              {{ group.parent_group_id ? '  └ ' : '' }}{{ group.name }}
            </span>
          </label>
        </div>
        <span class="hint">
          Selecciona una o más categorías para organizar el producto. Las categorías también influyen en los atributos heredados.
        </span>
      </div>

      <!-- Selected Summary -->
      <div v-if="selectedSummary" class="summary-box">
        <h4>Resumen de clasificación</h4>
        <div class="summary-content">
          <div class="summary-item">
            <span class="summary-label">Marca:</span>
            <span class="summary-value">{{ selectedBrandName }}</span>
          </div>
          <div v-if="localData.groupIds.length > 0" class="summary-item">
            <span class="summary-label">Categorías:</span>
            <span class="summary-value">{{ selectedGroupNames }}</span>
          </div>
          <div v-else class="summary-item">
            <span class="summary-label">Categorías:</span>
            <span class="summary-value empty">(ninguna seleccionada)</span>
          </div>
        </div>
      </div>

      <!-- Form Actions -->
      <div class="form-actions">
        <button
          type="button"
          @click="$emit('prev')"
          class="btn btn-secondary"
        >
          ← Anterior
        </button>
        <button
          type="submit"
          class="btn btn-primary"
          :disabled="!isStepValid"
        >
          Siguiente: Atributos →
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { reactive, computed, onMounted, watch, ref } from 'vue'
import { productApi } from '@/services/productApi'

const props = defineProps({
  modelValue: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['update:modelValue', 'next', 'prev'])

// Local copy of data
const localData = reactive({
  brandId: props.modelValue.brandId || '',
  groupIds: props.modelValue.groupIds || [],
})

// Catalog data
const brands = reactive([])
const groups = reactive([])
const isLoadingBrands = ref(false)
const isLoadingGroups = ref(false)
const loadError = reactive({
  brands: '',
  groups: '',
})

// Validation errors
const errors = reactive({
  brandId: '',
})

// Watch local data changes and emit to parent
watch(localData, (newValue) => {
  console.log('[ProductFormClassification] Data changed:', newValue)
  emit('update:modelValue', { ...newValue })
}, { deep: true })

// Computed properties
const selectedBrandName = computed(() => {
  const brand = brands.find(b => b.id === localData.brandId)
  return brand ? brand.name : ''
})

const selectedGroupNames = computed(() => {
  return localData.groupIds
    .map(id => {
      const group = groups.find(g => g.id === id)
      return group ? group.name : ''
    })
    .filter(Boolean)
    .join(', ')
})

const selectedSummary = computed(() => {
  return localData.brandId !== ''
})

const isStepValid = computed(() => {
  return localData.brandId !== '' && !errors.brandId
})

// Validation rules
const validationRules = {
  brandId: (value) => {
    if (!value) {
      return 'La marca es obligatoria'
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

// Handle next step
function handleNext() {
  if (validateAll() && isStepValid.value) {
    emit('next')
  }
}

// Load brands from API
async function loadBrands() {
  isLoadingBrands.value = true
  loadError.brands = ''
  try {
    const response = await productApi.listBrands({ isActive: true })
    brands.length = 0
    brands.push(...response.data)
  } catch (error) {
    loadError.brands = 'No se pudieron cargar las marcas. Intenta recargar la página.'
    console.error('Error loading brands:', error)
  } finally {
    isLoadingBrands.value = false
  }
}

// Load product groups from API
async function loadGroups() {
  isLoadingGroups.value = true
  loadError.groups = ''
  try {
    const response = await productApi.listProductGroups({ isActive: true })
    groups.length = 0
    groups.push(...response.data)
  } catch (error) {
    loadError.groups = 'No se pudieron cargar las categorías.'
    console.error('Error loading groups:', error)
  } finally {
    isLoadingGroups.value = false
  }
}

// Load data on mount
onMounted(() => {
  loadBrands()
  loadGroups()
})
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

.optional {
  font-size: 0.75rem;
  font-weight: 400;
  color: #94a3b8;
}

.form-group select {
  padding: 0.65rem 0.85rem;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-size: 0.95rem;
  font-family: inherit;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.form-group select:focus {
  outline: none;
  border-color: #1b3a6b;
  box-shadow: 0 0 0 3px rgba(27, 58, 107, 0.1);
}

.form-group select:disabled {
  background-color: #f1f5f9;
  cursor: not-allowed;
}

.checkbox-group {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  padding: 0.75rem;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  max-height: 300px;
  overflow-y: auto;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 400;
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 4px;
  transition: background-color 0.2s ease;
}

.checkbox-label:hover {
  background-color: #f8fafc;
}

.checkbox-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.checkbox-text {
  font-size: 0.9rem;
  color: #1e293b;
}

.hint {
  font-size: 0.8rem;
  color: #64748b;
}

.error {
  color: #ef4444;
  font-size: 0.8rem;
}

.loading-state,
.empty-state {
  padding: 1rem;
  text-align: center;
  color: #64748b;
  font-size: 0.9rem;
  border: 1px dashed #cbd5e1;
  border-radius: 6px;
}

.summary-box {
  background-color: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 1rem;
  margin-top: 0.5rem;
}

.summary-box h4 {
  margin: 0 0 0.75rem;
  font-size: 0.9rem;
  color: #1b3a6b;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.summary-content {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.summary-item {
  display: flex;
  gap: 0.5rem;
  font-size: 0.9rem;
}

.summary-label {
  font-weight: 600;
  color: #475569;
  min-width: 100px;
}

.summary-value {
  color: #1e293b;
}

.summary-value.empty {
  color: #94a3b8;
  font-style: italic;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: space-between;
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
