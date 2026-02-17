<template>
  <div class="modal-overlay" @click="handleClickOutside">
    <div class="modal-content" @click.stop>
      <!-- Header -->
      <div class="modal-header">
        <h3>{{ isEditMode ? 'Editar Variante' : 'Crear Variante' }}</h3>
        <button @click="$emit('close')" class="btn-close">✕</button>
      </div>

      <!-- Body -->
      <div class="modal-body">
        <form @submit.prevent="handleSubmit">
          <!-- Attribute Selection (only for create mode) -->
          <div v-if="!isEditMode" class="form-section">
            <h4>Seleccionar Combinación de Atributos</h4>
            <p class="hint">
              Elige los valores específicos para cada atributo del producto.
              La combinación debe ser única.
            </p>
            
            <div
              v-for="attr in productAttributes"
              :key="attr.id"
              class="form-group"
            >
              <label :for="`attr-${attr.id}`">
                {{ attr.name }} <span class="required">*</span>
              </label>
              <select
                :id="`attr-${attr.id}`"
                v-model="form.attributeValues[attr.id]"
                class="form-control"
                required
              >
                <option value="">-- Seleccionar --</option>
                <option
                  v-for="value in attr.values"
                  :key="value.id"
                  :value="value.id"
                >
                  {{ value.display_value }}
                </option>
              </select>
            </div>
          </div>

          <!-- SKU Preview (create mode) -->
          <div v-if="!isEditMode && computedSku" class="sku-preview">
            <strong>SKU Generado:</strong> 
            <code>{{ computedSku }}</code>
          </div>

          <!-- Existing SKU (edit mode) -->
          <div v-if="isEditMode" class="existing-sku">
            <strong>SKU:</strong> 
            <code>{{ variant?.sku }}</code>
            <p class="hint-small">La combinación de atributos no puede modificarse</p>
          </div>

          <!-- Metadata Fields -->
          <div class="form-section">
            <h4>Metadatos de la Variante</h4>

            <!-- Barcode -->
            <div class="form-group">
              <label for="barcode">Código de barras</label>
              <input
                id="barcode"
                v-model="form.barcode"
                type="text"
                class="form-control"
                placeholder="Ej: 7501234567890"
                maxlength="50"
              />
              <small class="form-text">Opcional. Código EAN/UPC para escáner</small>
            </div>

            <!-- Base Cost -->
            <div class="form-group">
              <label for="baseCost">Costo base</label>
              <input
                id="baseCost"
                v-model.number="form.baseCost"
                type="number"
                step="0.01"
                min="0"
                class="form-control"
                placeholder="0.00"
              />
              <small class="form-text">Costo unitario (opcional)</small>
            </div>

            <!-- Status -->
            <div class="form-group">
              <label for="status">
                Estado <span class="required">*</span>
              </label>
              <select
                id="status"
                v-model="form.status"
                class="form-control"
                required
              >
                <option value="PROVISIONAL">Provisional</option>
                <option value="CONFIRMED">Confirmado</option>
              </select>
              <small class="form-text">
                Las variantes nuevas se crean como PROVISIONAL hasta su confirmación
              </small>
            </div>

            <!-- Is Active -->
            <div class="form-group form-checkbox">
              <input
                id="isActive"
                v-model="form.isActive"
                type="checkbox"
                class="form-check-input"
              />
              <label for="isActive" class="form-check-label">
                Variante activa
              </label>
              <small class="form-text">
                Las variantes inactivas no aparecen en ventas por defecto
              </small>
            </div>
          </div>

          <!-- Error Display -->
          <div v-if="error" class="error-message">
            ⚠️ {{ error }}
          </div>

          <!-- Actions -->
          <div class="modal-footer">
            <button
              type="button"
              @click="$emit('close')"
              class="btn btn-secondary"
              :disabled="isSubmitting"
            >
              Cancelar
            </button>
            <button
              type="submit"
              class="btn btn-primary"
              :disabled="isSubmitting || !isFormValid"
            >
              <span v-if="isSubmitting">Guardando...</span>
              <span v-else>{{ isEditMode ? 'Actualizar' : 'Crear Variante' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { productApi } from '@/services/productApi'

const props = defineProps({
  productId: {
    type: String,
    required: true,
  },
  variant: {
    type: Object,
    default: null, // If provided, it's edit mode
  },
  productSku: {
    type: String,
    required: true,
  },
})

const emit = defineEmits(['close', 'saved'])

// State
const isEditMode = computed(() => props.variant !== null)
const productAttributes = ref([])
const isSubmitting = ref(false)
const error = ref('')

const form = ref({
  attributeValues: {}, // { attributeId: valueId }
  barcode: '',
  baseCost: null,
  status: 'PROVISIONAL',
  isActive: true,
})

// Computed
const isFormValid = computed(() => {
  if (!isEditMode.value) {
    // Create mode: all attributes must be selected
    const allAttributesSelected = productAttributes.value.every(
      attr => form.value.attributeValues[attr.id]
    )
    return allAttributesSelected && form.value.status
  }
  // Edit mode: always valid (only metadata)
  return true
})

const computedSku = computed(() => {
  if (isEditMode.value || !props.productSku) return ''
  
  const attrCodes = []
  for (const attr of productAttributes.value) {
    const valueId = form.value.attributeValues[attr.id]
    if (valueId) {
      const value = attr.values.find(v => v.id === valueId)
      if (value) {
        attrCodes.push(`${attr.code}.${value.code}`)
      }
    }
  }
  
  if (attrCodes.length === 0) return props.productSku
  return `${props.productSku}-${attrCodes.join('-')}`
})

// Lifecycle
onMounted(() => {
  loadProductAttributes()
  if (isEditMode.value) {
    populateFormFromVariant()
  }
})

// Methods
async function loadProductAttributes() {
  try {
    const product = await productApi.getProduct(props.productId)
    
    // Get calculated option sets (attributes with values)
    const optionSets = product.calculated_option_sets || []
    productAttributes.value = optionSets.map(opt => ({
      id: opt.attribute_id,
      name: opt.attribute_name,
      code: opt.attribute_code,
      values: opt.values || [],
    }))
  } catch (err) {
    console.error('[VariantFormModal] Error loading attributes:', err)
    error.value = 'No se pudieron cargar los atributos del producto'
  }
}

function populateFormFromVariant() {
  if (!props.variant) return
  
  form.value.barcode = props.variant.barcode || ''
  form.value.baseCost = props.variant.base_cost || null
  form.value.status = props.variant.status || 'PROVISIONAL'
  form.value.isActive = props.variant.is_active !== false
}

async function handleSubmit() {
  error.value = ''
  isSubmitting.value = true

  try {
    if (isEditMode.value) {
      await updateExistingVariant()
    } else {
      await createNewVariant()
    }
    
    emit('saved')
    emit('close')
  } catch (err) {
    console.error('[VariantFormModal] Submit error:', err)
    error.value = err.message || 'Ocurrió un error al guardar la variante'
  } finally {
    isSubmitting.value = false
  }
}

async function createNewVariant() {
  // Build option_configuration for JIT creation
  const optionConfiguration = {}
  
  for (const attr of productAttributes.value) {
    const valueId = form.value.attributeValues[attr.id]
    if (valueId) {
      optionConfiguration[attr.id] = valueId
    }
  }

  // Call findOrCreateVariant
  const result = await productApi.findOrCreateVariant(
    props.productId,
    optionConfiguration
  )

  // If it was just created (not found), update its metadata
  if (result.variant) {
    const variantId = result.variant.id
    
    // Update metadata fields if provided
    const hasMetadata = form.value.barcode || 
                        form.value.baseCost !== null || 
                        form.value.status !== 'PROVISIONAL' ||
                        !form.value.isActive

    if (hasMetadata) {
      await productApi.updateVariant(variantId, {
        barcode: form.value.barcode || undefined,
        base_cost: form.value.baseCost || undefined,
        status: form.value.status,
        is_active: form.value.isActive,
      })
    }
  }
}

async function updateExistingVariant() {
  await productApi.updateVariant(props.variant.id, {
    barcode: form.value.barcode || undefined,
    base_cost: form.value.baseCost || undefined,
    status: form.value.status,
    is_active: form.value.isActive,
  })
}

function handleClickOutside() {
  if (!isSubmitting.value) {
    emit('close')
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.modal-content {
  background: #ffffff;
  border-radius: 12px;
  max-width: 600px;
  width: 90%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  animation: slideUp 0.3s ease;
}

@keyframes slideUp {
  from {
    transform: translateY(30px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #e2e8f0;
  position: sticky;
  top: 0;
  background: #ffffff;
  z-index: 10;
}

.modal-header h3 {
  margin: 0;
  color: #1b3a6b;
  font-size: 1.25rem;
}

.btn-close {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #64748b;
  padding: 0;
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: background 0.2s ease;
}

.btn-close:hover {
  background: #f1f5f9;
}

.modal-body {
  padding: 1.5rem;
}

.form-section {
  margin-bottom: 2rem;
}

.form-section h4 {
  color: #1b3a6b;
  font-size: 1rem;
  margin: 0 0 0.5rem 0;
  font-weight: 700;
}

.hint {
  color: #64748b;
  font-size: 0.85rem;
  margin: 0 0 1rem 0;
  line-height: 1.5;
}

.hint-small {
  color: #94a3b8;
  font-size: 0.8rem;
  margin: 0.25rem 0 0 0;
  font-style: italic;
}

.form-group {
  margin-bottom: 1.25rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.4rem;
  color: #475569;
  font-weight: 600;
  font-size: 0.9rem;
}

.required {
  color: #dc2626;
  font-weight: 700;
}

.form-control {
  width: 100%;
  padding: 0.6rem 0.75rem;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-size: 0.9rem;
  color: #1e293b;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.form-control:focus {
  outline: none;
  border-color: #1b3a6b;
  box-shadow: 0 0 0 3px rgba(27, 58, 107, 0.1);
}

.form-control::placeholder {
  color: #94a3b8;
}

.form-text {
  display: block;
  margin-top: 0.3rem;
  color: #94a3b8;
  font-size: 0.8rem;
}

.form-checkbox {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
}

.form-check-input {
  width: 18px;
  height: 18px;
  margin-top: 0.2rem;
  cursor: pointer;
}

.form-check-label {
  margin: 0;
  cursor: pointer;
  flex: 1;
}

.sku-preview,
.existing-sku {
  padding: 1rem;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  margin-bottom: 1.5rem;
}

.sku-preview strong,
.existing-sku strong {
  color: #475569;
  margin-right: 0.5rem;
}

.sku-preview code,
.existing-sku code {
  background: #1e293b;
  color: #f4d03f;
  padding: 0.3rem 0.6rem;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 0.85rem;
  font-weight: 600;
}

.error-message {
  padding: 0.75rem 1rem;
  background: rgba(220, 38, 38, 0.1);
  border: 1px solid rgba(220, 38, 38, 0.3);
  border-radius: 6px;
  color: #dc2626;
  font-size: 0.9rem;
  margin-bottom: 1rem;
}

.modal-footer {
  padding: 1rem 1.5rem;
  border-top: 1px solid #e2e8f0;
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  position: sticky;
  bottom: 0;
  background: #ffffff;
}

.btn {
  border: none;
  border-radius: 8px;
  padding: 0.65rem 1.25rem;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.2s ease;
  font-weight: 600;
}

.btn-primary {
  background: #f4d03f;
  color: #1e293b;
}

.btn-primary:hover:not(:disabled) {
  background: #e6c530;
  box-shadow: 0 2px 8px rgba(244, 208, 63, 0.3);
}

.btn-primary:disabled {
  background: #e2e8f0;
  color: #94a3b8;
  cursor: not-allowed;
}

.btn-secondary {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  color: #475569;
}

.btn-secondary:hover:not(:disabled) {
  background: #f8fafc;
}

.btn-secondary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

@media (max-width: 640px) {
  .modal-content {
    width: 95%;
    max-height: 95vh;
  }

  .modal-header,
  .modal-body,
  .modal-footer {
    padding: 1rem;
  }
}
</style>
