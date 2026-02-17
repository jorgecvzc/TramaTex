<template>
  <div class="variant-selector">
    <!-- Header -->
    <div class="selector-header">
      <h3>{{ title || 'Seleccionar Variante' }}</h3>
      <p v-if="description" class="description">{{ description }}</p>
    </div>

    <!-- Mode Selector -->
    <div class="mode-tabs">
      <button
        :class="['tab-btn', { active: mode === 'attributes' }]"
        @click="switchMode('attributes')"
      >
        📋 Por Atributos
      </button>
      <button
        :class="['tab-btn', { active: mode === 'sku' }]"
        @click="switchMode('sku')"
      >
        🔍 Por SKU
      </button>
    </div>

    <!-- Error Display -->
    <div v-if="error" class="error-banner">
      ⚠️ {{ error }}
    </div>

    <!-- Mode A: Attribute Selection (Interactive with JIT) -->
    <div v-if="mode === 'attributes'" class="mode-content">
      <!-- Product Selection -->
      <div class="form-group">
        <label for="product-select">
          Producto <span class="required">*</span>
        </label>
        <select
          id="product-select"
          v-model="selectedProductId"
          class="form-control"
          @change="handleProductChange"
          :disabled="productId !== null"
        >
          <option value="">-- Seleccionar producto --</option>
          <option
            v-for="prod in availableProducts"
            :key="prod.id"
            :value="prod.id"
          >
            {{ prod.name }} ({{ prod.sku }})
          </option>
        </select>
      </div>

      <!-- Attribute Selectors -->
      <div v-if="productAttributes.length > 0" class="attributes-section">
        <h4>Configuración de Atributos</h4>
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
            v-model="selectedAttributes[attr.id]"
            class="form-control"
            @change="handleAttributeChange"
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

      <!-- SKU Preview -->
      <div v-if="previewSku" class="sku-preview">
        <strong>SKU:</strong> <code>{{ previewSku }}</code>
      </div>

      <!-- JIT Creation Button -->
      <div v-if="canCreateVariant" class="jit-section">
        <button
          @click="findOrCreateSelectedVariant"
          class="btn btn-primary"
          :disabled="isProcessing"
        >
          <span v-if="isProcessing">🔄 Procesando...</span>
          <span v-else>✨ Buscar o Crear Variante</span>
        </button>
        <p class="hint">
          Si la variante no existe, se creará automáticamente con estado PROVISIONAL
        </p>
      </div>

      <!-- Selected Variant Display -->
      <div v-if="selectedVariant" class="selected-variant">
        <div class="variant-card">
          <div class="card-header">
            <span class="badge">{{ formatStatus(selectedVariant.status) }}</span>
            <span v-if="!selectedVariant.is_active" class="badge inactive">Inactivo</span>
          </div>
          <div class="card-body">
            <p><strong>SKU:</strong> <code>{{ selectedVariant.sku }}</code></p>
            <p v-if="selectedVariant.barcode">
              <strong>Código de barras:</strong> {{ selectedVariant.barcode }}
            </p>
            <div class="card-footer">
              <button @click="confirmSelection" class="btn btn-success">
                ✓ Confirmar Selección
              </button>
              <button @click="clearSelection" class="btn btn-link">
                Cancelar
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Mode B: SKU Search (Direct) -->
    <div v-if="mode === 'sku'" class="mode-content">
      <div class="form-group">
        <label for="sku-input">
          SKU de Variante <span class="required">*</span>
        </label>
        <div class="search-input-group">
          <input
            id="sku-input"
            v-model="skuSearchQuery"
            type="text"
            class="form-control"
            placeholder="Ej: TST001-SIZE.M-COLOR.BLUE"
            @keyup.enter="searchBySku"
          />
          <button
            @click="searchBySku"
            class="btn btn-search"
            :disabled="isProcessing || !skuSearchQuery"
          >
            🔍 Buscar
          </button>
        </div>
        <small class="form-text">
          Ingresa el SKU completo de la variante que deseas seleccionar
        </small>
      </div>

      <!-- Search Result -->
      <div v-if="selectedVariant" class="selected-variant">
        <div class="variant-card">
          <div class="card-header">
            <span class="badge">{{ formatStatus(selectedVariant.status) }}</span>
            <span v-if="!selectedVariant.is_active" class="badge inactive">Inactivo</span>
          </div>
          <div class="card-body">
            <p><strong>Producto:</strong> {{ selectedVariant.product_name || 'N/A' }}</p>
            <p><strong>SKU:</strong> <code>{{ selectedVariant.sku }}</code></p>
            <p v-if="selectedVariant.barcode">
              <strong>Código de barras:</strong> {{ selectedVariant.barcode }}
            </p>
            <div class="attribute-values">
              <span
                v-for="(value, key) in selectedVariant.attribute_values"
                :key="key"
                class="attr-tag"
              >
                {{ key }}: {{ value }}
              </span>
            </div>
            <div class="card-footer">
              <button @click="confirmSelection" class="btn btn-success">
                ✓ Confirmar Selección
              </button>
              <button @click="clearSelection" class="btn btn-link">
                Cancelar
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { productApi } from '@/services/productApi'

const props = defineProps({
  title: {
    type: String,
    default: 'Seleccionar Variante de Producto',
  },
  description: {
    type: String,
    default: '',
  },
  productId: {
    type: String,
    default: null, // If provided, locks product selection
  },
  allowJitCreation: {
    type: Boolean,
    default: true, // Enable/disable JIT creation
  },
  initialMode: {
    type: String,
    default: 'attributes', // 'attributes' or 'sku'
    validator: (value) => ['attributes', 'sku'].includes(value),
  },
})

const emit = defineEmits(['variant-selected', 'error'])

// State
const mode = ref(props.initialMode)
const selectedProductId = ref(props.productId || '')
const availableProducts = ref([])
const productAttributes = ref([])
const selectedAttributes = ref({})
const selectedVariant = ref(null)
const skuSearchQuery = ref('')
const isProcessing = ref(false)
const error = ref('')

// Computed
const canCreateVariant = computed(() => {
  if (!props.allowJitCreation) return false
  if (!selectedProductId.value) return false
  if (productAttributes.value.length === 0) return false
  
  // Check all attributes are selected
  return productAttributes.value.every(
    attr => selectedAttributes.value[attr.id]
  )
})

const previewSku = computed(() => {
  if (!selectedProductId.value) return ''
  
  const product = availableProducts.value.find(p => p.id === selectedProductId.value)
  if (!product) return ''
  
  const attrCodes = []
  for (const attr of productAttributes.value) {
    const valueId = selectedAttributes.value[attr.id]
    if (valueId) {
      const value = attr.values.find(v => v.id === valueId)
      if (value) {
        attrCodes.push(`${attr.code}.${value.code}`)
      }
    }
  }
  
  if (attrCodes.length === 0) return product.sku
  return `${product.sku}-${attrCodes.join('-')}`
})

// Lifecycle
onMounted(() => {
  if (props.productId) {
    selectedProductId.value = props.productId
    loadProductAttributes(props.productId)
  } else {
    loadAvailableProducts()
  }
})

// Methods
async function loadAvailableProducts() {
  try {
    const response = await productApi.listProducts({ isActive: true })
    availableProducts.value = response.data || []
  } catch (err) {
    console.error('[VariantSelector] Error loading products:', err)
    setError('No se pudieron cargar los productos disponibles')
  }
}

async function loadProductAttributes(productId) {
  try {
    const product = await productApi.getProduct(productId)
    
    // Store product in available products for SKU generation
    if (!availableProducts.value.find(p => p.id === product.id)) {
      availableProducts.value.push(product)
    }
    
    const optionSets = product.calculated_option_sets || []
    productAttributes.value = optionSets.map(opt => ({
      id: opt.attribute_id,
      name: opt.attribute_name,
      code: opt.attribute_code,
      values: opt.values || [],
    }))
    
    // Reset selected attributes
    selectedAttributes.value = {}
  } catch (err) {
    console.error('[VariantSelector] Error loading product attributes:', err)
    setError('No se pudieron cargar los atributos del producto')
  }
}

function handleProductChange() {
  selectedAttributes.value = {}
  selectedVariant.value = null
  error.value = ''
  
  if (selectedProductId.value) {
    loadProductAttributes(selectedProductId.value)
  } else {
    productAttributes.value = []
  }
}

function handleAttributeChange() {
  selectedVariant.value = null
  error.value = ''
}

async function findOrCreateSelectedVariant() {
  if (!canCreateVariant.value) return
  
  isProcessing.value = true
  error.value = ''
  
  try {
    // Build option configuration
    const optionConfiguration = {}
    for (const attr of productAttributes.value) {
      const valueId = selectedAttributes.value[attr.id]
      if (valueId) {
        optionConfiguration[attr.id] = valueId
      }
    }
    
    // Call JIT endpoint
    const result = await productApi.findOrCreateVariant(
      selectedProductId.value,
      optionConfiguration
    )
    
    if (result.variant) {
      selectedVariant.value = result.variant
      
      // Check if variant is inactive
      if (!result.variant.is_active) {
        setError('⚠️ Esta variante está marcada como inactiva')
      }
    } else {
      setError('No se pudo crear/obtener la variante')
    }
  } catch (err) {
    console.error('[VariantSelector] JIT error:', err)
    setError(err.message || 'Error al buscar/crear la variante')
  } finally {
    isProcessing.value = false
  }
}

async function searchBySku() {
  if (!skuSearchQuery.value) return
  
  isProcessing.value = true
  error.value = ''
  selectedVariant.value = null
  
  try {
    const variant = await productApi.getVariantBySku(skuSearchQuery.value.trim())
    
    if (variant) {
      selectedVariant.value = variant
      
      if (!variant.is_active) {
        setError('⚠️ Esta variante está marcada como inactiva')
      }
    }
  } catch (err) {
    console.error('[VariantSelector] SKU search error:', err)
    setError(err.message || 'No se encontró ninguna variante con ese SKU')
  } finally {
    isProcessing.value = false
  }
}

function confirmSelection() {
  if (!selectedVariant.value) return
  
  emit('variant-selected', {
    variantId: selectedVariant.value.id,
    variant: selectedVariant.value,
  })
}

function clearSelection() {
  selectedVariant.value = null
  selectedAttributes.value = {}
  skuSearchQuery.value = ''
  error.value = ''
}

function switchMode(newMode) {
  mode.value = newMode
  clearSelection()
}

function formatStatus(status) {
  const map = {
    PROVISIONAL: 'Provisional',
    CONFIRMED: 'Confirmado',
  }
  return map[status] || status
}

function setError(message) {
  error.value = message
  emit('error', message)
}
</script>

<style scoped>
.variant-selector {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  padding: 1.5rem;
}

.selector-header h3 {
  color: #1b3a6b;
  margin: 0 0 0.5rem 0;
  font-size: 1.2rem;
}

.description {
  color: #64748b;
  font-size: 0.9rem;
  margin: 0 0 1rem 0;
}

.mode-tabs {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
  border-bottom: 2px solid #e2e8f0;
}

.tab-btn {
  background: none;
  border: none;
  padding: 0.75rem 1.25rem;
  font-size: 0.9rem;
  font-weight: 600;
  color: #64748b;
  cursor: pointer;
  border-bottom: 3px solid transparent;
  transition: all 0.2s ease;
  position: relative;
  bottom: -2px;
}

.tab-btn:hover {
  color: #1b3a6b;
}

.tab-btn.active {
  color: #1b3a6b;
  border-bottom-color: #f4d03f;
}

.error-banner {
  background: rgba(220, 38, 38, 0.1);
  border: 1px solid rgba(220, 38, 38, 0.3);
  border-radius: 8px;
  padding: 0.75rem 1rem;
  color: #dc2626;
  margin-bottom: 1rem;
  font-size: 0.9rem;
}

.mode-content {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.attributes-section h4 {
  color: #475569;
  font-size: 1rem;
  margin: 0 0 1rem 0;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.form-group label {
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
  padding: 0.65rem 0.75rem;
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

.form-control:disabled {
  background: #f1f5f9;
  color: #94a3b8;
  cursor: not-allowed;
}

.form-text {
  color: #94a3b8;
  font-size: 0.8rem;
}

.search-input-group {
  display: flex;
  gap: 0.5rem;
}

.search-input-group .form-control {
  flex: 1;
}

.sku-preview {
  padding: 0.75rem 1rem;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
}

.sku-preview strong {
  color: #475569;
  margin-right: 0.5rem;
}

.sku-preview code {
  background: #1e293b;
  color: #f4d03f;
  padding: 0.3rem 0.6rem;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 0.85rem;
  font-weight: 600;
}

.jit-section {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0.5rem;
}

.hint {
  color: #94a3b8;
  font-size: 0.85rem;
  margin: 0;
}

.selected-variant {
  margin-top: 1rem;
}

.variant-card {
  border: 2px solid #22c55e;
  border-radius: 8px;
  overflow: hidden;
  background: #f0fdf4;
}

.card-header {
  padding: 0.75rem 1rem;
  background: rgba(34, 197, 94, 0.1);
  display: flex;
  gap: 0.5rem;
}

.badge {
  display: inline-block;
  padding: 0.25rem 0.6rem;
  border-radius: 999px;
  font-weight: 600;
  font-size: 0.7rem;
  text-transform: uppercase;
  background: rgba(34, 197, 94, 0.2);
  color: #16a34a;
}

.badge.inactive {
  background: rgba(220, 38, 38, 0.2);
  color: #dc2626;
}

.card-body {
  padding: 1rem;
}

.card-body p {
  margin: 0.5rem 0;
  color: #475569;
  font-size: 0.9rem;
}

.card-body code {
  background: #1e293b;
  color: #f4d03f;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 0.85rem;
}

.attribute-values {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

.attr-tag {
  background: #e0f2fe;
  color: #0c4a6e;
  padding: 0.25rem 0.6rem;
  border-radius: 4px;
  font-size: 0.8rem;
  font-weight: 600;
}

.card-footer {
  display: flex;
  gap: 0.75rem;
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid rgba(34, 197, 94, 0.2);
}

.btn {
  border: none;
  border-radius: 6px;
  padding: 0.6rem 1rem;
  font-size: 0.85rem;
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

.btn-search {
  background: #1b3a6b;
  color: #ffffff;
}

.btn-search:hover:not(:disabled) {
  background: #152d52;
}

.btn-search:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-success {
  background: #22c55e;
  color: #ffffff;
}

.btn-success:hover {
  background: #16a34a;
  box-shadow: 0 2px 8px rgba(34, 197, 94, 0.3);
}

.btn-link {
  background: none;
  color: #64748b;
  text-decoration: underline;
}

.btn-link:hover {
  color: #475569;
}

@media (max-width: 640px) {
  .search-input-group {
    flex-direction: column;
  }
  
  .card-footer {
    flex-direction: column;
  }
  
  .btn {
    width: 100%;
  }
}
</style>
