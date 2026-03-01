<template>
  <div class="dashboard">
    <Navbar />
    <div class="dashboard-content">
      <header class="page-header">
        <div>
          <p class="breadcrumb">Operaciones / Productos / Consulta de Precios</p>
          <h1>Consulta de Precios</h1>
          <p class="subtitle">
            Consulta precios de venta para variantes. Selecciona producto, cliente y condiciones.
          </p>
        </div>
        <RouterLink to="/products" class="btn btn-secondary">
          ← Volver al catálogo
        </RouterLink>
      </header>

      <!-- Filters Card -->
      <div class="card filters-card">
        <h3>
          <Filter :size="20" class="icon-inline" />
          Condiciones de cálculo
        </h3>

        <div class="filters-grid">
          <!-- Product Selector -->
          <div class="form-group">
            <label for="filter-product">Producto</label>
            <select
              id="filter-product"
              v-model="selectedProductId"
              class="form-control"
              :disabled="isLoadingProducts"
              @change="onProductChange"
            >
              <option value="">Seleccionar producto...</option>
              <option v-for="p in products" :key="p.id" :value="p.id">
                {{ p.name }} {{ p.sku ? `(${p.sku})` : '' }}
              </option>
            </select>
          </div>

          <!-- Client Search -->
          <div class="form-group">
            <label for="filter-client">Cliente</label>
            <div class="client-search-container" @click.stop>
              <input
                v-model="clientSearchQuery"
                type="text"
                class="form-control"
                placeholder="Buscar cliente por nombre o NIF..."
                @input="searchClients"
              />
              <button
                v-if="selectedClientName"
                type="button"
                class="clear-btn"
                @click="clearClientSelection"
                title="Limpiar selección"
              >
                ✕
              </button>
            </div>
            <div v-if="showClientDropdown && filteredClients.length > 0" class="dropdown-list">
              <div
                v-for="client in filteredClients"
                :key="client.id"
                class="dropdown-option"
                @click="selectClient(client)"
              >
                <div class="option-name">{{ client.legalName || client.id }}</div>
                <div class="option-meta">
                  <span class="meta-role">{{ formatRole(client.role) }}</span>
                  <span v-if="client.taxId" class="meta-taxid">{{ client.taxId }}</span>
                </div>
              </div>
            </div>
            <div v-if="selectedClientName" class="selected-info">
              <CheckCircle :size="14" />
              <span>{{ selectedClientName }}</span>
            </div>
          </div>

          <!-- Quantity & Date -->
          <div class="form-group">
            <label for="filter-quantity">Cantidad</label>
            <input
              id="filter-quantity"
              v-model.number="quantity"
              type="number"
              class="form-control"
              min="1"
              placeholder="100"
            />
          </div>

          <div class="form-group">
            <label for="filter-date">Fecha de venta</label>
            <input
              id="filter-date"
              v-model="saleDate"
              type="date"
              class="form-control"
            />
          </div>
        </div>
      </div>

      <!-- Error Alert -->
      <div v-if="calculationError" class="alert-error">
        <strong>Error:</strong> {{ calculationError }}
        <button class="alert-close" @click="calculationError = ''">&times;</button>
      </div>

      <!-- Variants Price Table -->
      <div class="card">
        <div class="section-header">
          <h3>
            <BarChart3 :size="20" class="icon-inline" />
            Precios de variantes
          </h3>
          <div class="section-actions">
            <button
              v-if="selectedProductId && variants.length > 0"
              class="btn btn-secondary"
              :disabled="isCalculatingBase"
              @click="calculateBasePrices"
            >
              {{ isCalculatingBase ? 'Calculando...' : 'Calcular Base' }}
            </button>
            <button
              v-if="selectedProductId && variants.length > 0 && selectedClientId"
              class="btn btn-primary"
              :disabled="isCalculatingFinal"
              @click="calculateFinalPrices"
            >
              {{ isCalculatingFinal ? 'Calculando...' : 'Calcular Precio Final' }}
            </button>
          </div>
        </div>

        <!-- Info: no client -->
        <div v-if="selectedProductId && variants.length > 0 && !selectedClientId" class="info-banner">
          Selecciona un <strong>cliente</strong> para poder calcular el precio final de venta.
        </div>

        <!-- No product selected -->
        <div v-if="!selectedProductId" class="empty-state">
          <Package :size="48" class="empty-icon" />
          <p>Selecciona un producto para ver sus variantes y precios.</p>
        </div>

        <!-- Loading variants -->
        <div v-if="selectedProductId && isLoadingVariants" class="loading-inline">
          <div class="spinner-small"></div>
          <span>Cargando variantes...</span>
        </div>

        <!-- No variants -->
        <div v-if="selectedProductId && !isLoadingVariants && variants.length === 0" class="empty-state">
          <Package :size="48" class="empty-icon" />
          <p>Este producto no tiene variantes configuradas.</p>
        </div>

        <!-- Price Table -->
        <div v-if="!isLoadingVariants && variants.length > 0" class="table-wrapper">
          <table class="data-table">
            <thead>
              <tr>
                <th>SKU</th>
                <th>Atributos</th>
                <th class="align-right">Coste Base</th>
                <th class="align-right">Precio Base Venta</th>
                <th class="align-right">Precio Final</th>
                <th>Estado</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="variant in variants" :key="variant.id">
                <td>
                  <code class="sku-code">{{ variant.sku }}</code>
                </td>
                <td>
                  <div class="attribute-tags">
                    <span
                      v-for="(value, attr) in (variant.option_configuration || variant.attribute_values || {})"
                      :key="attr"
                      class="attribute-tag"
                    >
                      <span class="attr-name">{{ attr }}:</span>
                      <span class="attr-value">{{ value }}</span>
                    </span>
                    <span v-if="!variant.option_configuration && !variant.attribute_values" class="text-muted">
                      Sin atributos
                    </span>
                  </div>
                </td>
                <td class="align-right">
                  <span v-if="variant.base_cost != null" class="price-value muted">
                    {{ formatEUR(variant.base_cost) }}
                  </span>
                  <span v-else class="text-muted">—</span>
                </td>
                <td class="align-right">
                  <span v-if="variantPrices[variant.id]?.baseSalesPrice" class="price-value">
                    {{ formatMoney(variantPrices[variant.id].baseSalesPrice) }}
                  </span>
                  <span v-else-if="variantPrices[variant.id]?.baseError" class="text-error" :title="variantPrices[variant.id].baseError">
                    Error
                  </span>
                  <span v-else-if="loadingPrices[variant.id]" class="text-muted">Calculando...</span>
                  <span v-else class="text-muted">—</span>
                </td>
                <td class="align-right">
                  <span v-if="variantPrices[variant.id]?.finalPrice" class="price-value highlight">
                    {{ formatMoney(variantPrices[variant.id].finalPrice) }}
                  </span>
                  <span v-else-if="variantPrices[variant.id]?.finalError" class="text-error" :title="variantPrices[variant.id].finalError">
                    Error
                  </span>
                  <span v-else-if="loadingPrices[variant.id]" class="text-muted">Calculando...</span>
                  <span v-else-if="!selectedClientId" class="text-muted">Sin cliente</span>
                  <span v-else class="text-muted">—</span>
                </td>
                <td>
                  <span class="pill" :class="variant.is_active ? 'active' : 'inactive'">
                    {{ variant.is_active ? 'Activo' : 'Inactivo' }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { RouterLink } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import { productApi } from '@/services/productApi'
import { pricingApi } from '@/services/pricingApi'
import { partyApi } from '@/services/partyApi'
import { Package, Filter, BarChart3, CheckCircle } from 'lucide-vue-next'

// State - Products
const products = ref([])
const isLoadingProducts = ref(false)
const selectedProductId = ref('')

// State - Variants
const variants = ref([])
const isLoadingVariants = ref(false)

// State - Prices
const variantPrices = ref({})
const loadingPrices = ref({})
const isCalculatingBase = ref(false)
const isCalculatingFinal = ref(false)
const calculationError = ref('')

// State - Filters
const quantity = ref(100)
const saleDate = ref(new Date().toISOString().split('T')[0])

// State - Client search
const clientSearchQuery = ref('')
const filteredClients = ref([])
const showClientDropdown = ref(false)
const selectedClientId = ref('')
const selectedClientName = ref('')
let searchTimeout = null

// Lifecycle
onMounted(async () => {
  await loadProducts()
  document.addEventListener('click', handleOutsideClick)
})

onUnmounted(() => {
  document.removeEventListener('click', handleOutsideClick)
  if (searchTimeout) clearTimeout(searchTimeout)
})

function handleOutsideClick() {
  showClientDropdown.value = false
}

// Methods
async function loadProducts() {
  isLoadingProducts.value = true
  try {
    const result = await productApi.listProducts()
    products.value = result.data || result || []
  } catch (err) {
    console.error('Error loading products:', err)
    products.value = []
  } finally {
    isLoadingProducts.value = false
  }
}

async function onProductChange() {
  variants.value = []
  variantPrices.value = {}
  calculationError.value = ''
  if (!selectedProductId.value) return

  isLoadingVariants.value = true
  try {
    const data = await productApi.listProductVariants(selectedProductId.value)
    variants.value = data.variants || data || []
    // Auto-calculate base prices when product is selected
    if (variants.value.length > 0) {
      await calculateBasePrices()
    }
  } catch (err) {
    console.error('Error loading variants:', err)
    variants.value = []
  } finally {
    isLoadingVariants.value = false
  }
}

async function calculateBasePrices() {
  if (!selectedProductId.value || variants.value.length === 0) return

  isCalculatingBase.value = true
  calculationError.value = ''
  let errorCount = 0

  for (const variant of variants.value) {
    loadingPrices.value[variant.id] = true
    try {
      const baseResult = await pricingApi.calculateBaseSalesPrice(
        selectedProductId.value,
        variant.id
      )
      variantPrices.value[variant.id] = {
        ...variantPrices.value[variant.id],
        baseSalesPrice: baseResult?.baseSalesPrice || null,
        baseError: null,
      }
    } catch (err) {
      errorCount++
      console.error(`Error calculating base price for variant ${variant.id}:`, err)
      variantPrices.value[variant.id] = {
        ...variantPrices.value[variant.id],
        baseSalesPrice: null,
        baseError: err?.message || 'Error de cálculo',
      }
    } finally {
      loadingPrices.value[variant.id] = false
    }
  }

  if (errorCount > 0 && errorCount === variants.value.length) {
    calculationError.value = 'No se pudieron calcular los precios base. Verifica que existan reglas de precio configuradas.'
  }

  isCalculatingBase.value = false
}

async function calculateFinalPrices() {
  if (!selectedProductId.value || variants.value.length === 0 || !selectedClientId.value) return

  isCalculatingFinal.value = true
  calculationError.value = ''

  try {
    // Build all sale items at once for a single API call
    const saleItems = variants.value.map(v => ({
      productVariantId: v.id,
      quantity: quantity.value || 1,
    }))

    const result = await pricingApi.calculateFinalSalePrice(
      saleItems,
      selectedClientId.value,
      new Date(saleDate.value)
    )

    // Map results back to variants
    if (result?.calculatedItems) {
      for (const item of result.calculatedItems) {
        const variantId = item.productVariantId
        variantPrices.value[variantId] = {
          ...variantPrices.value[variantId],
          baseSalesPrice: item.baseSalesPrice || variantPrices.value[variantId]?.baseSalesPrice || null,
          finalPrice: item.finalPrice || null,
          finalError: null,
        }
      }
    }
  } catch (err) {
    console.error('Error calculating final prices:', err)
    calculationError.value = err?.message || 'Error al calcular el precio final de venta. Verifica la configuración de reglas.'
    // Mark all variants with final price error
    for (const variant of variants.value) {
      variantPrices.value[variant.id] = {
        ...variantPrices.value[variant.id],
        finalPrice: null,
        finalError: err?.message || 'Error de cálculo',
      }
    }
  } finally {
    isCalculatingFinal.value = false
  }
}

// Client search
async function searchClients() {
  if (searchTimeout) clearTimeout(searchTimeout)
  const query = clientSearchQuery.value.trim()
  if (query.length < 2) {
    filteredClients.value = []
    showClientDropdown.value = false
    return
  }
  searchTimeout = setTimeout(async () => {
    try {
      const response = await partyApi.listParties({ name: query, role: 'client', pageSize: 10 })
      filteredClients.value = response.data || []
      showClientDropdown.value = true
    } catch {
      filteredClients.value = []
    }
  }, 300)
}

function selectClient(client) {
  selectedClientId.value = client.id
  selectedClientName.value = client.legalName || client.id
  clientSearchQuery.value = ''
  showClientDropdown.value = false
  filteredClients.value = []
}

function clearClientSelection() {
  selectedClientId.value = ''
  selectedClientName.value = ''
  clientSearchQuery.value = ''
  filteredClients.value = []
  showClientDropdown.value = false
}

function formatRole(role) {
  const names = { client: 'Cliente', supplier: 'Proveedor', employee: 'Empleado', partner: 'Socio' }
  return names[role] || role
}

function formatEUR(value) {
  if (value === null || value === undefined) return '—'
  return new Intl.NumberFormat('es-ES', {
    style: 'currency',
    currency: 'EUR',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(value)
}

function formatMoney(moneyDTO) {
  if (!moneyDTO || moneyDTO.amount === undefined) return '—'
  const currency = moneyDTO.currency || 'EUR'
  return new Intl.NumberFormat('es-ES', {
    style: 'currency',
    currency: currency,
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(parseFloat(moneyDTO.amount))
}
</script>

<style scoped>
.dashboard {
  min-height: 100vh;
  background-color: #f1f5f9;
  font-family: 'Inter', sans-serif;
}

.dashboard-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.page-header {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

.page-header h1 {
  color: #1b3a6b;
  margin: 0.25rem 0 0;
}

.breadcrumb {
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #64748b;
  margin: 0;
}

.subtitle {
  color: #64748b;
  margin: 0.5rem 0 0;
  font-size: 0.95rem;
}

.card {
  background-color: #ffffff;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
  border: 1px solid #e2e8f0;
}

.filters-card h3 {
  margin: 0 0 1rem;
  color: #1b3a6b;
  font-size: 1.1rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.icon-inline {
  vertical-align: middle;
}

.filters-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.form-group label {
  font-size: 0.8rem;
  font-weight: 600;
  color: #475569;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.form-control {
  width: 100%;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  padding: 0.6rem 0.8rem;
  font-size: 0.9rem;
  color: #1e293b;
  font-family: inherit;
  background: #ffffff;
}

.form-control:focus {
  outline: none;
  border-color: #1b3a6b;
  box-shadow: 0 0 0 3px rgba(27, 58, 107, 0.1);
}

/* Client search */
.client-search-container {
  position: relative;
}

.clear-btn {
  position: absolute;
  right: 8px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  cursor: pointer;
  color: #94a3b8;
  font-size: 1rem;
  padding: 0.2rem 0.4rem;
}

.clear-btn:hover {
  color: #475569;
}

.dropdown-list {
  position: absolute;
  z-index: 50;
  background: #fff;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  margin-top: 0.25rem;
  max-height: 200px;
  overflow-y: auto;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
  width: 100%;
}

.dropdown-option {
  padding: 0.6rem 0.8rem;
  cursor: pointer;
  transition: background 0.15s;
}

.dropdown-option:hover {
  background: #f1f5f9;
}

.option-name {
  font-weight: 500;
  color: #1e293b;
}

.option-meta {
  display: flex;
  gap: 0.5rem;
  font-size: 0.75rem;
  color: #64748b;
  margin-top: 0.15rem;
}

.selected-info {
  display: flex;
  align-items: center;
  gap: 0.3rem;
  color: #059669;
  font-size: 0.8rem;
  margin-top: 0.25rem;
}

/* Section header */
.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.section-header h3 {
  margin: 0;
  color: #1b3a6b;
  font-size: 1.1rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.section-actions {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

/* Alert & Info banner */
.alert-error {
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 8px;
  padding: 0.75rem 1rem;
  color: #991b1b;
  font-size: 0.9rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.alert-close {
  margin-left: auto;
  background: none;
  border: none;
  color: #991b1b;
  font-size: 1.2rem;
  cursor: pointer;
  padding: 0 0.3rem;
}

.info-banner {
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  border-radius: 8px;
  padding: 0.6rem 1rem;
  color: #1e40af;
  font-size: 0.85rem;
  margin-bottom: 1rem;
}

.text-error {
  color: #dc2626;
  font-size: 0.8rem;
  font-weight: 500;
  cursor: help;
}

/* Table */
.table-wrapper {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th {
  text-align: left;
  padding: 0.75rem 0.8rem;
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #64748b;
  border-bottom: 2px solid #e2e8f0;
  white-space: nowrap;
}

.data-table td {
  padding: 0.75rem 0.8rem;
  border-bottom: 1px solid #f1f5f9;
  font-size: 0.9rem;
}

.data-table tr:hover {
  background: #f8fafc;
}

.align-right {
  text-align: right;
}

.sku-code {
  background: #f1f5f9;
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  font-size: 0.85rem;
  color: #1b3a6b;
}

.attribute-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.3rem;
}

.attribute-tag {
  background: #eff6ff;
  padding: 0.15rem 0.5rem;
  border-radius: 4px;
  font-size: 0.8rem;
}

.attr-name {
  color: #64748b;
  font-weight: 500;
}

.attr-value {
  color: #1e293b;
}

.price-value {
  font-weight: 600;
  color: #059669;
  font-size: 0.95rem;
}

.price-value.muted {
  color: #64748b;
  font-weight: 500;
}

.price-value.highlight {
  color: #1b3a6b;
  font-size: 1rem;
}

.text-muted {
  color: #94a3b8;
}

.pill {
  display: inline-block;
  padding: 0.2rem 0.6rem;
  border-radius: 20px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
}

.pill.active {
  background: #dcfce7;
  color: #166534;
}

.pill.inactive {
  background: #fee2e2;
  color: #991b1b;
}

/* Empty state */
.empty-state {
  text-align: center;
  padding: 3rem;
  color: #64748b;
}

.empty-icon {
  color: #cbd5e1;
  margin-bottom: 1rem;
}

/* Loading */
.loading-inline {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 1.5rem;
  color: #64748b;
}

.spinner-small {
  width: 20px;
  height: 20px;
  border: 3px solid rgba(27, 58, 107, 0.12);
  border-top-color: #1b3a6b;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Buttons */
.btn {
  border: none;
  border-radius: 8px;
  padding: 0.6rem 1rem;
  font-size: 0.85rem;
  cursor: pointer;
  transition: background 0.2s ease, box-shadow 0.2s ease;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  gap: 0.4rem;
}

.btn-primary {
  background: #1b3a6b;
  color: #ffffff;
}

.btn-primary:hover {
  background: #15325d;
}

.btn-primary:disabled {
  background: #94a3b8;
  cursor: not-allowed;
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
  .filters-grid {
    grid-template-columns: 1fr;
  }

  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.75rem;
  }
}
</style>
