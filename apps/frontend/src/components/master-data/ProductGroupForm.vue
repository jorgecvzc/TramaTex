<template>
  <div class="product-group-form">
    <div class="form-group">
      <label for="name">Nombre de la categoría <span class="required">*</span></label>
      <input 
        id="name"
        v-model="formData.name" 
        type="text" 
        class="form-input"
        placeholder="Ej: Calzado Deportivo, Ropa"
        @input="clearError('name')"
      />
      <span v-if="errors.name" class="error-message">{{ errors.name }}</span>
    </div>

    <div class="form-group">
      <label>Tipo de categoría <span class="required">*</span></label>
      <div class="radio-group">
        <label class="radio-label">
          <input 
            type="radio" 
            v-model="formData.type" 
            value="TANGIBLE"
            name="groupType"
          />
          <div class="radio-content">
            <span class="radio-title"><Wrench :size="18" style="vertical-align: middle; margin-right: 4px" /> Productos Tangibles</span>
            <span class="radio-description">Productos físicos: calzado, ropa, accesorios, equipamiento</span>
          </div>
        </label>
        <label class="radio-label">
          <input 
            type="radio" 
            v-model="formData.type" 
            value="SERVICE"
            name="groupType"
          />
          <div class="radio-content">
            <span class="radio-title"><Settings :size="18" style="vertical-align: middle; margin-right: 4px" /> Servicios</span>
            <span class="radio-description">Servicios profesionales: consultoría, mantenimiento, instalación</span>
          </div>
        </label>
      </div>
      <span v-if="errors.type" class="error-message">{{ errors.type }}</span>
    </div>

    <div class="form-group">
      <label for="parent">Categoría padre (opcional)</label>
      <select 
        id="parent"
        v-model="formData.parentGroupId" 
        class="form-input"
        :disabled="loadingGroups"
      >
        <option value="">Sin categoría padre (categoría raíz)</option>
        <option 
          v-for="group in availableParentGroups" 
          :key="group.id" 
          :value="group.id"
        >
          {{ group.name }}
        </option>
      </select>
      <small class="hint">Selecciona una categoría padre para crear una jerarquía</small>
    </div>

    <div class="form-group">
      <label class="checkbox-label">
        <input 
          type="checkbox" 
          v-model="formData.isActive"
        />
        <span>Categoría activa</span>
      </label>
      <small class="hint">Las categorías inactivas no se mostrarán en los formularios de productos</small>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, computed, onMounted } from 'vue'
import { Wrench, Settings } from 'lucide-vue-next'
import { productApi } from '@/services/productApi'

const props = defineProps({
  productGroup: {
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
  name: props.productGroup?.name || '',
  type: props.productGroup?.type || 'TANGIBLE',
  parentGroupId: props.productGroup?.parent_group_id || '',
  isActive: props.productGroup?.is_active ?? true
})

const errors = reactive({
  name: '',
  type: ''
})

const allGroups = ref([])
const loadingGroups = ref(false)

const availableParentGroups = computed(() => {
  if (props.mode === 'edit' && props.productGroup) {
    // Exclude self and own descendants to prevent circular references
    return allGroups.value.filter(g => g.id !== props.productGroup.id)
  }
  return allGroups.value
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
  
  if (!formData.type || !['TANGIBLE', 'SERVICE'].includes(formData.type)) {
    errors.type = 'Debe seleccionar un tipo válido'
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
    type: formData.type,
    parentGroupId: formData.parentGroupId || null,
    isActive: formData.isActive
  }
  
  if (props.mode === 'edit' && props.productGroup) {
    payload.id = props.productGroup.id
  }
  
  emit('submit', payload)
}

async function loadProductGroups() {
  loadingGroups.value = true
  try {
    const response = await productApi.listProductGroups({})
    allGroups.value = response.data || []
  } catch (error) {
    console.error('Error loading product groups:', error)
  } finally {
    loadingGroups.value = false
  }
}

onMounted(() => {
  loadProductGroups()
})

defineExpose({
  handleSubmit
})
</script>

<style scoped>
.product-group-form {
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

.mb-3 {
  margin-bottom: 1.5rem;
}

.radio-group {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.radio-label {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  padding: 1rem;
  border: 2px solid #e2e8f0;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: all 0.2s;
  background: #f8fafc;
}

.radio-label:hover {
  border-color: #94a3b8;
  background: #ffffff;
}

.radio-label input[type="radio"] {
  margin-top: 0.25rem;
  width: 1.25rem;
  height: 1.25rem;
  cursor: pointer;
  flex-shrink: 0;
}

.radio-label input[type="radio"]:checked {
  accent-color: #3b82f6;
}

.radio-label:has(input:checked) {
  border-color: #3b82f6;
  background: #eff6ff;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.radio-content {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  flex: 1;
}

.radio-title {
  font-weight: 600;
  font-size: 0.9375rem;
  color: #1e293b;
}

.radio-description {
  font-size: 0.8125rem;
  color: #64748b;
  line-height: 1.4;
}
</style>
