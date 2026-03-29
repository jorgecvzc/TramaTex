<template>
  <div class="modal-overlay" @click="handleClickOutside">
    <div class="modal-content large" @click.stop>
      <!-- Header -->
      <div class="modal-header">
        <h3>Vista Previa de Combinaciones</h3>
        <button @click="$emit('close')" class="btn-close"><X :size="16" /></button>
      </div>

      <!-- Body -->
      <div class="modal-body">
        <div class="info-banner">
          <p>
            <strong>Total de combinaciones posibles:</strong> {{ allCombinations.length }}
          </p>
          <p class="hint">
            Selecciona las variantes que deseas crear y asigna códigos de barras si los tienes.
            Las variantes se crearán con estado CONFIRMADO y coste base heredado del producto.
          </p>
        </div>

        <!-- Loading State -->
        <div v-if="isLoading" class="loading">
          <div class="spinner"></div>
          <p>Cargando combinaciones...</p>
        </div>

        <!-- Error State -->
        <div v-if="error" class="error-message">
          <AlertTriangle :size="16" /> {{ error }}
        </div>

        <!-- Combinations Table -->
        <div v-if="!isLoading && allCombinations.length > 0" class="combinations-wrapper">
          <div class="table-actions">
            <button @click="selectAll" class="btn btn-secondary btn-sm">
              Seleccionar todas
            </button>
            <button @click="deselectAll" class="btn btn-secondary btn-sm">
              Deseleccionar todas
            </button>
          </div>

          <table class="combinations-table">
            <thead>
              <tr>
                <th class="col-checkbox">
                  <input 
                    type="checkbox" 
                    :checked="allSelected"
                    @change="toggleAll"
                  />
                </th>
                <th class="col-sku">SKU</th>
                <th class="col-attributes">Combinación</th>
                <th class="col-barcode">Código de Barras</th>
                <th class="col-status">Estado</th>
              </tr>
            </thead>
            <tbody>
              <tr 
                v-for="(combo, index) in allCombinations" 
                :key="index"
                :class="{ selected: combo.selected }"
              >
                <td class="col-checkbox">
                  <input 
                    type="checkbox" 
                    v-model="combo.selected"
                  />
                </td>
                <td class="col-sku">
                  <code>{{ combo.sku }}</code>
                </td>
                <td class="col-attributes">
                  <div class="attribute-tags">
                    <span 
                      v-for="(value, attr) in combo.attributeDisplay" 
                      :key="attr"
                      class="attribute-tag"
                    >
                      <span class="attr-name">{{ attr }}:</span>
                      <span class="attr-value">{{ value }}</span>
                    </span>
                  </div>
                </td>
                <td class="col-barcode">
                  <input 
                    type="text" 
                    v-model="combo.barcode"
                    placeholder="Opcional"
                    class="barcode-input"
                    :disabled="!combo.selected"
                  />
                </td>
                <td class="col-status">
                  <span v-if="combo.exists" class="pill exists">Ya existe</span>
                  <span v-else class="pill new">Nueva</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Empty State -->
        <div v-if="!isLoading && allCombinations.length === 0" class="empty-state">
          <p>No hay atributos configurados para este producto.</p>
          <p class="hint">
            Añade atributos al producto para generar variantes.
          </p>
        </div>

        <!-- Summary -->
        <div v-if="selectedCount > 0" class="summary-box">
          <p>
            <strong>{{ selectedCount }}</strong> variante(s) seleccionada(s) para crear
          </p>
        </div>
      </div>

      <!-- Footer -->
      <div class="modal-footer">
        <button 
          @click="$emit('close')" 
          class="btn btn-secondary"
          :disabled="isCreating"
        >
          Cancelar
        </button>
        <button 
          @click="handleCreate" 
          class="btn btn-primary"
          :disabled="isCreating || selectedCount === 0"
        >
          <span v-if="isCreating">Creando...</span>
          <span v-else>Crear {{ selectedCount }} variante(s)</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { X, AlertTriangle } from 'lucide-vue-next'
import { productApi } from '@/services/productApi'

const props = defineProps({
  productId: {
    type: String,
    required: true,
  },
  productSku: {
    type: String,
    required: true,
  },
})

const emit = defineEmits(['close', 'created'])

// State
const isLoading = ref(false)
const isCreating = ref(false)
const error = ref('')
const productAttributes = ref([])
const existingVariants = ref([])
const allCombinations = reactive([])

// Computed
const selectedCount = computed(() => {
  return allCombinations.filter(c => c.selected && !c.exists).length
})

const allSelected = computed(() => {
  const selectableCount = allCombinations.filter(c => !c.exists).length
  return selectableCount > 0 && selectedCount.value === selectableCount
})

// Lifecycle
onMounted(() => {
  loadData()
})

// Methods
async function loadData() {
  isLoading.value = true
  error.value = ''

  try {
    // Load product details with attributes
    const product = await productApi.getProduct(props.productId)
    productAttributes.value = product.calculated_option_sets || []

    // Load existing variants
    const variantsResponse = await productApi.listProductVariants(props.productId)
    existingVariants.value = variantsResponse.data || []

    // Generate all possible combinations
    generateCombinations()
  } catch (err) {
    console.error('[VariantBatchCreator] Error loading data:', err)
    error.value = 'No se pudieron cargar los datos del producto'
  } finally {
    isLoading.value = false
  }
}

function generateCombinations() {
  if (productAttributes.value.length === 0) {
    allCombinations.splice(0)
    return
  }

  // Build cartesian product of all attribute values
  const attributeArrays = productAttributes.value.map(attr => ({
    id: attr.id,
    name: attr.name,
    code: attr.code,
    values: attr.values || [],
  }))

  const cartesianProduct = (arrays) => {
    if (arrays.length === 0) return [[]]
    const [first, ...rest] = arrays
    const restProduct = cartesianProduct(rest)
    return first.values.flatMap(value =>
      restProduct.map(product => [
        { attrId: first.id, attrName: first.name, attrCode: first.code, ...value },
        ...product
      ])
    )
  }

  const combinations = cartesianProduct(attributeArrays)

  // Map to display format
  allCombinations.splice(0)
  combinations.forEach(combo => {
    // Generate SKU
    const skuParts = combo.map(item => `${item.attrCode}.${item.code}`)
    const sku = `${props.productSku}-${skuParts.join('-')}`

    // Build attribute display
    const attributeDisplay = {}
    const attributeValueIds = []
    const optionConfig = {} // Store for API call
    combo.forEach(item => {
      attributeDisplay[item.attrName] = item.value
      attributeValueIds.push(item.id)
      optionConfig[item.attrCode] = item.value // Backend expects: { AttributeCode: AttributeValue }
    })

    // Check if already exists
    const exists = existingVariants.value.some(v => v.sku === sku)

    allCombinations.push(reactive({
      sku,
      attributeDisplay,
      attributeValueIds,
      optionConfiguration: optionConfig,
      barcode: '',
      selected: !exists, // Auto-select new combinations
      exists,
    }))
  })
}

function selectAll() {
  allCombinations.forEach(combo => {
    if (!combo.exists) {
      combo.selected = true
    }
  })
}

function deselectAll() {
  allCombinations.forEach(combo => {
    combo.selected = false
  })
}

function toggleAll() {
  const newState = !allSelected.value
  allCombinations.forEach(combo => {
    if (!combo.exists) {
      combo.selected = newState
    }
  })
}

function handleClickOutside(event) {
  if (event.target.classList.contains('modal-overlay')) {
    emit('close')
  }
}

async function handleCreate() {
  isCreating.value = true
  error.value = ''

  try {
    const selectedCombinations = allCombinations.filter(c => c.selected && !c.exists)
    
    if (selectedCombinations.length === 0) {
      error.value = 'No hay variantes seleccionadas para crear'
      return
    }

    // Create variants one by one
    const createdVariants = []
    const errors = []

    for (const combo of selectedCombinations) {
      try {
        // Use pre-built option_configuration with AttributeCode: AttributeValue format
        const optionConfiguration = combo.optionConfiguration

        // Call findOrCreateVariant
        const result = await productApi.findOrCreateVariant(
          props.productId,
          optionConfiguration
        )

        // Update barcode if provided
        if (combo.barcode && result.variant) {
          await productApi.updateVariant(result.variant.id, {
            barcode: combo.barcode,
          })
        }

        createdVariants.push(result.variant)
      } catch (err) {
        console.error(`Error creating variant ${combo.sku}:`, err)
        errors.push({ sku: combo.sku, error: err.message })
      }
    }

    if (errors.length > 0) {
      error.value = `Se crearon ${createdVariants.length} variantes, pero ${errors.length} fallaron`
    }

    emit('created', { created: createdVariants, errors })
    
    if (errors.length === 0) {
      emit('close')
    }
  } catch (err) {
    console.error('[VariantBatchCreator] Error creating variants:', err)
    error.value = 'Ocurrió un error al crear las variantes'
  } finally {
    isCreating.value = false
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
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 2rem;
}

.modal-content {
  background: white;
  border-radius: 8px;
  width: 100%;
  max-width: 1200px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
}

.modal-content.large {
  max-width: 1300px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #e2e8f0;
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
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.2s;
}

.btn-close:hover {
  background-color: #f1f5f9;
  color: #1e293b;
}

.modal-body {
  padding: 1.5rem;
  overflow-y: auto;
  flex: 1;
}

.info-banner {
  background-color: #eff6ff;
  border-left: 4px solid #3b82f6;
  padding: 1rem;
  margin-bottom: 1.5rem;
  border-radius: 4px;
}

.info-banner p {
  margin: 0.5rem 0;
  color: #1e40af;
}

.info-banner .hint {
  font-size: 0.875rem;
  color: #3730a3;
}

.loading,
.empty-state {
  text-align: center;
  padding: 3rem;
  color: #64748b;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e2e8f0;
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-message {
  background-color: #fee;
  border-left: 4px solid #ef4444;
  padding: 1rem;
  margin-bottom: 1rem;
  color: #991b1b;
  border-radius: 4px;
}

.combinations-wrapper {
  margin-top: 1rem;
}

.table-actions {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.btn-sm {
  font-size: 0.875rem;
  padding: 0.5rem 1rem;
}

.combinations-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;
}

.combinations-table th {
  background-color: #f8fafc;
  padding: 0.75rem;
  text-align: left;
  font-weight: 600;
  color: #475569;
  border-bottom: 2px solid #e2e8f0;
}

.combinations-table td {
  padding: 0.75rem;
  border-bottom: 1px solid #e2e8f0;
}

.combinations-table tr:hover {
  background-color: #f8fafc;
}

.combinations-table tr.selected {
  background-color: #eff6ff;
}

.col-checkbox {
  width: 40px;
  text-align: center;
}

.col-sku {
  width: 200px;
}

.col-attributes {
  min-width: 300px;
}

.col-barcode {
  width: 180px;
}

.col-status {
  width: 100px;
}

code {
  background-color: #f1f5f9;
  padding: 0.25rem 0.5rem;
  border-radius: 3px;
  font-family: 'Courier New', monospace;
  font-size: 0.875rem;
}

.attribute-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.attribute-tag {
  background-color: #f1f5f9;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.875rem;
}

.attr-name {
  color: #64748b;
  font-weight: 500;
}

.attr-value {
  color: #1e293b;
  margin-left: 0.25rem;
}

.barcode-input {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #cbd5e1;
  border-radius: 4px;
  font-size: 0.875rem;
}

.barcode-input:disabled {
  background-color: #f1f5f9;
  cursor: not-allowed;
}

.pill {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
  text-transform: uppercase;
}

.pill.exists {
  background-color: #fef3c7;
  color: #92400e;
}

.pill.new {
  background-color: #d1fae5;
  color: #065f46;
}

.summary-box {
  background-color: #f0fdf4;
  border: 1px solid #86efac;
  padding: 1rem;
  margin-top: 1rem;
  border-radius: 4px;
}

.summary-box p {
  margin: 0;
  color: #166534;
  font-size: 0.875rem;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  padding: 1.5rem;
  border-top: 1px solid #e2e8f0;
}

.btn {
  padding: 0.625rem 1.25rem;
  border-radius: 6px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
  font-size: 0.875rem;
}

.btn-secondary {
  background-color: #f1f5f9;
  color: #475569;
}

.btn-secondary:hover:not(:disabled) {
  background-color: #e2e8f0;
}

.btn-primary {
  background-color: #3b82f6;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background-color: #2563eb;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
