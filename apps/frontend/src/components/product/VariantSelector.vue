<template>
  <div class="variant-selector">
    <!-- Header -->
    <div class="selector-header">
      <h3>{{ title || 'Seleccionar Variante' }}</h3>
      <p v-if="description" class="description">{{ description }}</p>
    </div>

    <!-- Error Display -->
    <div v-if="error" class="error-banner">
      <AlertTriangle :size="16" style="vertical-align: middle" /> {{ error }}
    </div>

    <!-- Search Section -->
    <div class="search-section">
      <div class="form-group">
        <label for="quick-search-input">
          SKU, código de barras o referencia
        </label>
        <div class="search-input-group">
          <input
            id="quick-search-input"
            ref="quickSearchInput"
            v-model="quickSearchQuery"
            type="text"
            class="form-control"
            placeholder="Ej: TST001, TST001-SIZE.M, código de barras..."
            @keyup.enter="performSmartSearch"
          />
          <button
            @click="performSmartSearch"
            class="btn btn-search"
            :disabled="isProcessing || !quickSearchQuery"
          >
            <Search :size="18" />
          </button>
        </div>
        <small class="form-text">
          Introduce un SKU (parcial o completo), código de barras o referencia y pulsa Enter
        </small>
      </div>

      <!-- Smart Search: product list result -->
      <div v-if="smartSearchProducts.length > 0 && !selectedProductId" class="product-list-result">
        <h4>Productos encontrados</h4>
        <div
          v-for="prod in smartSearchProducts"
          :key="prod.id"
          class="product-item"
          @click="selectProductFromSearch(prod)"
        >
          <span class="product-item-sku">{{ prod.sku }}</span>
          <span class="product-item-name">{{ prod.name }}</span>
        </div>
      </div>

      <!-- Smart Search: exact variant found (no product context) -->
      <div v-if="selectedVariant && !selectedProductId" class="selected-variant">
        <div class="variant-card">
          <div class="card-header">
            <span class="badge">{{ formatStatus(selectedVariant.status) }}</span>
            <span v-if="!selectedVariant.is_active" class="badge inactive">Inactivo</span>
          </div>
          <div class="card-body">
            <p v-if="selectedVariant.product_name"><strong>Producto:</strong> {{ selectedVariant.product_name }}</p>
            <p><strong>SKU:</strong> <code>{{ selectedVariant.sku }}</code></p>
            <p v-if="selectedVariant.barcode">
              <strong>Código de barras:</strong> {{ selectedVariant.barcode }}
            </p>
            <div class="card-footer">
              <button @click="confirmSelection" class="btn btn-success btn-add">
                <Check :size="16" style="margin-right: 4px; vertical-align: middle" /> Agregar
              </button>
              <button @click="clearSelection" class="btn btn-link">
                Cancelar
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Divider -->
    <hr class="section-divider" />

    <!-- Product + Attribute Selection -->
    <div class="attribute-selection-section">
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
            <option value="" disabled>-- Seleccionar --</option>
            <option
              v-for="value in attr.values"
              :key="value.id"
              :value="value.id"
            >
              {{ value.value }}
            </option>
          </select>
        </div>
      </div>

      <!-- SKU Preview -->
      <div v-if="previewSku" class="sku-preview">
        <strong>SKU:</strong> <code>{{ previewSku }}</code>
      </div>

      <!-- JIT + Add (single button) -->
      <div v-if="canCreateVariant && !selectedVariant" class="jit-section">
        <button
          @click="findOrCreateAndAdd"
          class="btn btn-success btn-add"
          :disabled="isProcessing"
        >
          <span v-if="isProcessing">Cargando...</span>
          <span v-else><Check :size="16" /> Agregar</span>
        </button>
      </div>

      <!-- Selected Variant Display -->
      <div v-if="selectedVariant && selectedProductId" class="selected-variant">
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
              <button @click="confirmSelection" class="btn btn-success btn-add">
                <Check :size="16" /> Agregar
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
import { ref, computed, onMounted, nextTick } from 'vue'
import { productApi } from '@/services/productApi'
import { AlertTriangle, Search, Check } from 'lucide-vue-next'

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
  initialQuery: {
    type: String,
    default: '', // If provided, auto-trigger smart search on mount
  },
})

const emit = defineEmits(['variant-selected', 'error'])

// Refs
const quickSearchInput = ref(null)

// State
const selectedProductId = ref(props.productId || '')
const selectedProductName = ref('')
const availableProducts = ref([])
const productAttributes = ref([])
const selectedAttributes = ref({})
const selectedVariant = ref(null)
const quickSearchQuery = ref('')
const smartSearchProducts = ref([])
const isProcessing = ref(false)
const error = ref('')

// Computed
const canCreateVariant = computed(() => {
  if (!props.allowJitCreation) return false
  if (!selectedProductId.value) return false
  
  // Products with no attributes: variant creation is handled via auto-trigger
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
  loadAvailableProducts()
  if (props.productId) {
    selectedProductId.value = props.productId
    loadProductAttributes(props.productId)
  }
  // Auto-search if initialQuery is provided
  if (props.initialQuery) {
    quickSearchQuery.value = props.initialQuery
    nextTick(() => {
      performSmartSearch()
    })
  } else {
    nextTick(() => {
      quickSearchInput.value?.focus()
    })
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
      id: opt.id,
      name: opt.name,
      code: opt.code,
      values: opt.values || [],
    }))
    
    // Product with no attributes → auto-trigger default variant
    if (productAttributes.value.length === 0 && props.allowJitCreation) {
      await autoSelectDefaultVariant(productId)
      return
    }
    
    // Initialize selectedAttributes keys for reactivity
    const newSelectedAttributes = {}
    productAttributes.value.forEach(attr => {
      newSelectedAttributes[attr.id] = selectedAttributes.value[attr.id] || ''
    })
    selectedAttributes.value = newSelectedAttributes
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

/**
 * Auto-select default variant for products with no attributes.
 * Calls JIT endpoint with empty optionConfiguration.
 */
async function autoSelectDefaultVariant(productId) {
  isProcessing.value = true
  error.value = ''

  try {
    const result = await productApi.findOrCreateVariant(productId, {})

    if (result.variant) {
      const product = availableProducts.value.find(p => p.id === productId)
      selectedVariant.value = {
        ...result.variant,
        product_name: product?.name || '',
        product_description: product?.description || '',
        product_base_price: product?.base_price ?? null,
        product_tax_rate: product?.tax_rate ?? null,
      }
      // Immediately emit selection
      emit('variant-selected', {
        variantId: result.variant.id,
        variant: selectedVariant.value,
      })
    } else {
      setError('No se pudo cargar la variante por defecto')
    }
  } catch (err) {
    console.error('[VariantSelector] Default variant error:', err)
    setError(err.message || 'Error al cargar la variante por defecto')
  } finally {
    isProcessing.value = false
  }
}

async function findOrCreateSelectedVariant() {
  if (!canCreateVariant.value) return
  
  isProcessing.value = true
  error.value = ''
  selectedVariant.value = null
  
  try {
    // Build option configuration with AttributeCode: AttributeValue format
    const optionConfiguration = {}
    for (const attr of productAttributes.value) {
      const valueId = selectedAttributes.value[attr.id]
      if (valueId) {
        // Find the selected value to get its string value
        const selectedValue = attr.values.find(v => v.id === valueId)
        if (selectedValue) {
          optionConfiguration[attr.code] = selectedValue.value
        }
      }
    }
    
    // Call JIT endpoint
    const result = await productApi.findOrCreateVariant(
      selectedProductId.value,
      optionConfiguration
    )
    
    if (result.variant) {
      const product = availableProducts.value.find(p => p.id === selectedProductId.value)
      selectedVariant.value = {
        ...result.variant,
        product_base_price: product?.base_price ?? null,
        product_tax_rate: product?.tax_rate ?? null,
      }
    } else {
      setError('No se pudo cargar la variante')
    }
  } catch (err) {
    console.error('[VariantSelector] JIT error:', err)
    setError(err.message || 'Error al cargar la variante')
  } finally {
    isProcessing.value = false
  }
}

/**
 * Find or create variant AND immediately emit (single "Agregar" in attributes mode)
 */
async function findOrCreateAndAdd() {
  if (!canCreateVariant.value) return
  
  isProcessing.value = true
  error.value = ''
  
  try {
    const optionConfiguration = {}
    for (const attr of productAttributes.value) {
      const valueId = selectedAttributes.value[attr.id]
      if (valueId) {
        const selectedValue = attr.values.find(v => v.id === valueId)
        if (selectedValue) {
          optionConfiguration[attr.code] = selectedValue.value
        }
      }
    }
    
    const result = await productApi.findOrCreateVariant(
      selectedProductId.value,
      optionConfiguration
    )
    
    if (result.variant) {
      const product = availableProducts.value.find(p => p.id === selectedProductId.value)
      selectedVariant.value = {
        ...result.variant,
        product_name: product?.name || '',
        product_description: product?.description || '',
        product_base_price: product?.base_price ?? null,
        product_tax_rate: product?.tax_rate ?? null,
      }
      // Immediately emit selection
      emit('variant-selected', {
        variantId: result.variant.id,
        variant: selectedVariant.value,
      })
    } else {
      setError('No se pudo cargar la variante')
    }
  } catch (err) {
    console.error('[VariantSelector] JIT + Add error:', err)
    setError(err.message || 'Error al cargar la variante')
  } finally {
    isProcessing.value = false
  }
}

/**
 * Perform smart search via backend endpoint
 */
async function performSmartSearch() {
  const query = quickSearchQuery.value?.trim()
  if (!query) return
  
  isProcessing.value = true
  error.value = ''
  selectedVariant.value = null
  selectedProductId.value = ''
  selectedProductName.value = ''
  productAttributes.value = []
  selectedAttributes.value = {}
  smartSearchProducts.value = []
  
  try {
    const result = await productApi.smartSearch(query)
    
    switch (result.type) {
      case 'exact_variant':
        // Variant found directly → Auto-select and confirm immediately
        if (result.variant) {
          selectedVariant.value = {
            ...result.variant,
            product_name: result.product?.name || '',
            product_description: result.product?.description || '',
            product_base_price: result.product?.base_price ?? null,
            product_tax_rate: result.product?.tax_rate ?? null,
          }
          // Emit the event right away to close modal and add line
          confirmSelection()
        }
        break

        
      case 'exact_product':
        // Product found → load its attribute selectors
        if (result.product) {
          selectedProductId.value = result.product.id
          selectedProductName.value = `${result.product.name} (${result.product.sku})`
          if (!availableProducts.value.find(p => p.id === result.product.id)) {
            availableProducts.value.push(result.product)
          }
          if (result.optionSets) {
            loadOptionSetsIntoAttributes(result.optionSets)
          } else {
            await loadProductAttributes(result.product.id)
          }
        }
        break
        
      case 'partial_match':
        // Partial SKU match → load product with pre-selected attributes
        if (result.product) {
          selectedProductId.value = result.product.id
          selectedProductName.value = `${result.product.name} (${result.product.sku})`
          if (!availableProducts.value.find(p => p.id === result.product.id)) {
            availableProducts.value.push(result.product)
          }
          if (result.optionSets) {
            loadOptionSetsIntoAttributes(result.optionSets, result.selectedAttributes)
          } else {
            await loadProductAttributes(result.product.id)
          }
        }
        break
        
      case 'product_list':
        // Multiple products found → show list to pick from
        if (result.products && result.products.length > 0) {
          smartSearchProducts.value = result.products
        } else {
          setError('No se encontraron productos')
        }
        break
        
      case 'no_match':
        setError('No se encontró ningún producto o variante con esa referencia')
        break
    }
  } catch (err) {
    console.error('[VariantSelector] Smart search error:', err)
    setError(err.message || 'Error en la búsqueda')
  } finally {
    isProcessing.value = false
  }
}

/**
 * Load option sets from smart search response into attribute selectors
 */
function loadOptionSetsIntoAttributes(optionSets, preselected = null) {
  productAttributes.value = optionSets.map(os => ({
    id: os.attributeId,
    name: os.attributeName,
    code: os.attributeCode,
    values: os.values || [],
  }))
  
  // Product with no attributes → auto-trigger default variant
  if (productAttributes.value.length === 0 && props.allowJitCreation && selectedProductId.value) {
    autoSelectDefaultVariant(selectedProductId.value)
    return
  }
  
  // Initialize selectedAttributes
  const newSelectedAttributes = {}
  for (const attr of productAttributes.value) {
    newSelectedAttributes[attr.id] = ''
    
    // Pre-select attribute if we have a match from partial reference (case-insensitive)
    if (preselected) {
      const attrCodeUpper = (attr.code || '').toUpperCase()
      const preselectedCode = Object.keys(preselected).find(k => k.toUpperCase() === attrCodeUpper)
      if (preselectedCode) {
        const valueCode = preselected[preselectedCode]
        const matchingValue = attr.values.find(v => (v.code || '').toUpperCase() === (valueCode || '').toUpperCase())
        if (matchingValue) {
          newSelectedAttributes[attr.id] = matchingValue.id
        }
      }
    }
  }
  selectedAttributes.value = newSelectedAttributes
  
  // If all attributes are pre-selected, auto-trigger JIT
  if (preselected && canCreateVariant.value) {
    findOrCreateAndAdd()
  }
}

/**
 * Select a product from smart search product list
 */
async function selectProductFromSearch(product) {
  selectedProductId.value = product.id
  selectedProductName.value = `${product.name} (${product.sku})`
  smartSearchProducts.value = []
  await loadProductAttributes(product.id)
}

function confirmSelection() {
  if (!selectedVariant.value) return
  
  // Enrich with product info if not already present
  if (selectedProductId.value) {
    const product = availableProducts.value.find(p => p.id === selectedProductId.value)
    if (product) {
      if (!selectedVariant.value.product_name) {
        selectedVariant.value.product_name = product.name || ''
      }
      if (selectedVariant.value.product_base_price == null) {
        selectedVariant.value.product_base_price = product.base_price ?? null
      }
      if (selectedVariant.value.product_tax_rate == null) {
        selectedVariant.value.product_tax_rate = product.tax_rate ?? null
      }
      if (!selectedVariant.value.product_description) {
        selectedVariant.value.product_description = product.description || ''
      }
    }
  }
  
  emit('variant-selected', {
    variantId: selectedVariant.value.id,
    variant: selectedVariant.value,
  })
}

function clearSelection() {
  selectedVariant.value = null
  selectedAttributes.value = {}
  quickSearchQuery.value = ''
  smartSearchProducts.value = []
  selectedProductName.value = ''
  selectedProductId.value = props.productId || ''
  productAttributes.value = []
  error.value = ''
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

.section-divider {
  border: none;
  border-top: 1px solid #e2e8f0;
  margin: 1.5rem 0;
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

.search-section,
.attribute-selection-section {
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
  background: #ffffff;
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

/* Better styling for select dropdowns */
.form-control option {
  color: #1e293b;
  background: #ffffff;
  padding: 0.5rem;
}

.form-control option:disabled {
  color: #94a3b8;
  font-style: italic;
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

.btn-add {
  min-width: 120px;
  padding: 0.7rem 1.5rem;
  font-size: 0.95rem;
}

.btn-sm {
  padding: 0.3rem 0.6rem;
  font-size: 0.8rem;
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

/* Quick search product list */
.product-list-result h4 {
  color: #475569;
  font-size: 0.95rem;
  margin: 0 0 0.75rem 0;
}

.product-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.65rem 0.75rem;
  border: 1px solid #e2e8f0;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s ease;
  margin-bottom: 0.4rem;
}

.product-item:hover {
  background: #f0f9ff;
  border-color: #1b3a6b;
}

.product-item-sku {
  background: #1e293b;
  color: #f4d03f;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 0.8rem;
  font-weight: 600;
  white-space: nowrap;
}

.product-item-name {
  color: #475569;
  font-size: 0.9rem;
}

.selected-product-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.5rem 0.75rem;
  background: #f0f9ff;
  border: 1px solid #bae6fd;
  border-radius: 6px;
  margin-bottom: 0.75rem;
}

.selected-product-header strong {
  color: #0c4a6e;
  font-size: 0.9rem;
}
</style>
