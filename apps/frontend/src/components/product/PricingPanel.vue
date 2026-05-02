<template>
  <div class="pricing-panel">
    <div class="pricing-header">
      <div>
        <h3>Gestión de Precios</h3>
        <p class="subtitle">
          Precios base por variante y calculadora de precios finales
        </p>
      </div>
      <button class="btn btn-primary" @click="showCalculator = !showCalculator">
        <component :is="getIcon(showCalculator ? 'bar_chart' : 'calculate')" :size="16" style="margin-right: 6px" />
        {{ showCalculator ? 'Ver Precios Base': 'Calculadora de Precios' }}
      </button>
    </div>

    <!-- Pricing Calculator -->
    <div v-if="showCalculator" class="calculator-section card">
      <h4>
        <component :is="getIcon('calculate')" style="vertical-align: middle; margin-right: 6px" :size="20" />
        Calculadora de Precios
      </h4>
      <p class="help-text">
        Simula el precio final para un cliente específico, cantidad y fecha de venta.
      </p>

      <div class="calculator-form">
        <div class="form-group">
          <label for="calc-variant">Variante del Producto</label>
          <select
            id="calc-variant"
            v-model="calculator.selectedVariantId"
            class="form-control"
            :disabled="isLoadingVariants || variants.length === 0"
          >
            <option value="">Seleccionar variante...</option>
            <option v-for="variant in variants" :key="variant.id" :value="variant.id">
              {{ productName }} - {{ variant.sku }} - {{ formatVariantAttributes(variant) }}
            </option>
          </select>
        </div>

        <div class="form-group">
          <PartySelector
            v-model="calculator.clientId"
            label="Cliente"
            placeholder="Buscar cliente por nombre..."
            role-filter="CLIENT"
            help-text="Selecciona el cliente para calcular el precio final."
          />
        </div>

        <div class="form-row">
          <div class="form-group">
            <label for="calc-quantity">Cantidad</label>
            <input
              id="calc-quantity"
              v-model.number="calculator.quantity"
              type="number"
              class="form-control"
              min="1"
              placeholder="100"
            />
          </div>

          <div class="form-group">
            <label for="calc-date">Fecha de Venta</label>
            <input
              id="calc-date"
              v-model="calculator.saleDate"
              type="date"
              class="form-control"
            />
          </div>
        </div>

        <button
          class="btn btn-success"
          :disabled="!canCalculate || isCalculating"
          @click="calculateFinalPrice"
        >
          <component :is="getIcon('attach_money')" :size="16" style="margin-right: 6px" />
          {{ isCalculating ? 'Calculando...' : 'Calcular Precio Final' }}
        </button>
      </div>

      <!-- Calculation Result -->
      <div v-if="calculationResult" class="calculation-result">
        <h5>
          <component :is="getIcon('check_circle')" style="vertical-align: middle; margin-right: 6px" :size="18" />
          Resultado del Cálculo
        </h5>
        <div class="result-grid">
          <div class="result-item">
            <span class="result-label">Precio Base Venta:</span>
            <span class="result-value">
              {{ formatMoney(calculationResult.calculatedItems[0]?.baseSalesPrice) }}
            </span>
          </div>
          <div class="result-item" :class="calculationResult.calculatedItems[0]?.discountPercent > 0 ? 'discount' : ''">
            <span class="result-label">Descuento Cliente:</span>
            <span class="result-value">
              <template v-if="calculationResult.calculatedItems[0]?.discountPercent > 0">
                {{ calculationResult.calculatedItems[0].discountPercent.toFixed(2) }}%
                ({{ formatMoney(calculationResult.calculatedItems[0].discountAmount) }})
              </template>
              <span v-else class="text-muted">Sin descuento</span>
            </span>
          </div>
          <div class="result-item total">
            <span class="result-label">Precio Final (ud.):</span>
            <span class="result-value">
              {{ formatMoney(calculationResult.calculatedItems[0]?.finalPrice) }}
            </span>
          </div>
          <div class="result-item">
            <span class="result-label">Total Venta (sin IVA):</span>
            <span class="result-value">
              {{ formatMoney(calculationResult.saleTotal) }}
            </span>
          </div>
          <div class="result-item">
            <span class="result-label">Total con IVA:</span>
            <span class="result-value">
              {{ formatMoney(calculationResult.saleTotalWithTax) }}
            </span>
          </div>
        </div>
      </div>

      <!-- Calculation Error -->
      <div v-if="calculationError" class="alert-error">
        {{ calculationError }}
      </div>
    </div>

    <!-- Base Prices Table -->
    <div v-if="!showCalculator" class="prices-section card">
      <h4>
        <component :is="getIcon('bar_chart')" style="vertical-align: middle; margin-right: 6px" :size="20" />
        Precios Base por Variante
      </h4>
      <p class="help-text">
        Estos son los precios base de venta configurados para cada variante.
        Los precios finales pueden variar según reglas de modificación y cliente.
      </p>

      <!-- Loading State -->
      <div v-if="isLoadingVariants" class="loading-inline">
        <div class="spinner-small"></div>
        <span>Cargando precios...</span>
      </div>

      <!-- Empty State -->
      <div v-if="!isLoadingVariants && variants.length === 0" class="empty-state">
        <component :is="getIcon('inventory_2')" class="empty-icon" :size="64" />
        <p>No hay variantes configuradas para este producto.</p>
        <p class="empty-hint">
          Las variantes deben crearse primero en la pestaña "Variantes".
        </p>
      </div>

      <!-- Prices Table -->
      <div v-if="!isLoadingVariants && variants.length > 0" class="table-wrapper">
        <table class="data-table">
          <thead>
            <tr>
              <th>SKU</th>
              <th>Atributos</th>
              <th>Precio Base Venta</th>
              <th>Estado</th>
              <th>Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="variant in variants" :key="variant.id">
              <td>
                <code class="sku-code">{{ variant.sku }}</code>
              </td>
              <td>
                <div class="attributes-cell">
                  {{ formatVariantAttributes(variant) }}
                </div>
              </td>
              <td>
                <div v-if="variantPrices[variant.id]">
                  <span class="price-value">
                    {{ formatMoney(variantPrices[variant.id]) }}
                  </span>
                </div>
                <div v-else class="loading-price">
                  <span v-if="loadingPrices[variant.id]">⏳ Calculando...</span>
                  <span v-else class="text-muted">No calculado</span>
                </div>
              </td>
              <td>
                <span class="pill" :class="variant.is_active ? 'active' : 'inactive'">
                  {{ variant.is_active ? 'Activo' : 'Inactivo' }}
                </span>
              </td>
              <td>
                <button
                  class="btn btn-sm btn-outline"
                  :disabled="loadingPrices[variant.id]"
                  @click="loadBasePriceForVariant(variant.id)"
                >
                  🔄 Recalcular
                </button>
                <button
                  class="btn btn-sm btn-secondary"
                  @click="viewPriceHistory(variant.id)"
                >
                  <component :is="getIcon('assignment')" :size="14" style="margin-right: 4px" />
                  Historial
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Pricing Rules Section (Future Enhancement) -->
    <div class="rules-info card">
      <h4>
        <component :is="getIcon('settings')" style="vertical-align: middle; margin-right: 6px" :size="20" />
        Reglas de Pricing
      </h4>
      <div class="info-message">
        <component :is="getIcon('info')" class="info-icon" :size="20" />
        <div>
          <p>
            <strong>Próximamente:</strong> Visualización de reglas de pricing aplicables
            (reglas de precio base, modificaciones de venta, descuentos por volumen).
          </p>
          <p class="help-text">
            Las reglas se configuran desde el módulo Pricing (Admin → Pricing → Reglas).
          </p>
        </div>
      </div>
    </div>

    <!-- Price History Modal -->
    <div v-if="showHistoryModal" class="modal-overlay" @click.self="closeHistoryModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3><component :is="getIcon('assignment')" :size="16" style="vertical-align: middle; margin-right: 6px" /> Historial de Precios</h3>
          <button class="modal-close" @click="closeHistoryModal"><component :is="getIcon('close')" :size="16" /></button>
        </div>

        <div class="modal-body">
          <!-- Loading -->
          <div v-if="isLoadingHistory" class="loading-inline">
            <div class="spinner-small"></div>
            <span>Cargando historial...</span>
          </div>

          <!-- Empty History -->
          <div v-if="!isLoadingHistory && priceHistory.length === 0" class="empty-state">
            <component :is="getIcon('bar_chart')" :size="64" class="empty-icon" />
            <p>No hay historial de cálculos para esta variante.</p>
          </div>

          <!-- History Table -->
          <div v-if="!isLoadingHistory && priceHistory.length > 0" class="history-table">
            <table class="data-table">
              <thead>
                <tr>
                  <th>Fecha</th>
                  <th>Cliente</th>
                  <th>Cantidad</th>
                  <th>Precio Base</th>
                  <th>Precio Final</th>
                  <th>Reglas Aplicadas</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="record in priceHistory" :key="record.id">
                  <td>{{ formatDate(record.calculatedAt) }}</td>
                  <td>
                    <code class="client-id">{{ record.clientId.substring(0, 8) }}...</code>
                  </td>
                  <td>{{ record.quantity }}</td>
                  <td>{{ formatMoney(record.baseCost) }}</td>
                  <td>{{ formatMoney(record.finalPrice) }}</td>
                  <td>
                    <div v-if="record.appliedRules && record.appliedRules.length > 0" class="rules-list">
                      <span v-for="(rule, idx) in record.appliedRules" :key="idx" class="rule-tag">
                        {{ rule }}
                      </span>
                    </div>
                    <span v-else class="text-muted">Sin reglas</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn btn-secondary" @click="closeHistoryModal">Cerrar</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { getIcon } from '@/utils/icons'
import { pricingApi } from '@/services/pricingApi'
import PartySelector from '@/components/party/PartySelector.vue'

const props = defineProps({
  productId: {
    type: String,
    required: true,
  },
  productName: {
    type: String,
    default: '',
  },
  variants: {
    type: Array,
    default: () => [],
  },
  isLoadingVariants: {
    type: Boolean,
    default: false,
  },
})

// State
const showCalculator = ref(false)
const variantPrices = ref({})
const loadingPrices = ref({})
const showHistoryModal = ref(false)
const isLoadingHistory = ref(false)
const priceHistory = ref([])
const selectedHistoryVariantId = ref(null)

// Calculator state
const calculator = reactive({
  selectedVariantId: '',
  clientId: '',
  quantity: 100,
  saleDate: new Date().toISOString().split('T')[0],
})
const isCalculating = ref(false)
const calculationResult = ref(null)
const calculationError = ref('')

// Computed
const canCalculate = computed(() => {
  return (
    calculator.selectedVariantId &&
    calculator.clientId &&
    calculator.quantity > 0 &&
    calculator.saleDate
  )
})

// Lifecycle
onMounted(async () => {
  // Load base prices for all variants
  if (props.variants && props.variants.length > 0) {
    props.variants.forEach((variant) => {
      loadBasePriceForVariant(variant.id)
    })
  }
})

// Methods
async function loadBasePriceForVariant(variantId) {
  loadingPrices.value[variantId] = true

  try {
    const result = await pricingApi.calculateBaseSalesPrice(props.productId, variantId)
    variantPrices.value[variantId] = result.baseSalesPrice
  } catch (err) {
    console.error('Error loading base price for variant:', err)
    variantPrices.value[variantId] = null
  } finally {
    loadingPrices.value[variantId] = false
  }
}

async function calculateFinalPrice() {
  isCalculating.value = true
  calculationError.value = ''
  calculationResult.value = null

  try {
    const saleItems = [
      {
        productVariantId: calculator.selectedVariantId,
        quantity: calculator.quantity,
      },
    ]

    const result = await pricingApi.calculateFinalSalePrice(
      saleItems,
      calculator.clientId,
      new Date(calculator.saleDate)
    )

    calculationResult.value = result
  } catch (err) {
    calculationError.value = err?.message || 'Error al calcular el precio'
    console.error('Error calculating final price:', err)
  } finally {
    isCalculating.value = false
  }
}

async function viewPriceHistory(variantId) {
  selectedHistoryVariantId.value = variantId
  showHistoryModal.value = true
  isLoadingHistory.value = true

  try {
    const history = await pricingApi.getPricingHistory(variantId)
    priceHistory.value = history || []
  } catch (err) {
    console.error('Error loading price history:', err)
    priceHistory.value = []
  } finally {
    isLoadingHistory.value = false
  }
}

function closeHistoryModal() {
  showHistoryModal.value = false
  selectedHistoryVariantId.value = null
  priceHistory.value = []
}

function formatVariantAttributes(variant) {
  if (!variant.attribute_values || variant.attribute_values.length === 0) {
    return 'Sin atributos'
  }

  return variant.attribute_values
    .map((av) => `${av.attribute_name}: ${av.value}`)
    .join(', ')
}

function formatMoney(moneyDTO) {
  if (!moneyDTO || moneyDTO.amount === undefined) {
    return '—'
  }

  const currency = moneyDTO.currency || 'EUR'
  const amount = parseFloat(moneyDTO.amount).toFixed(2)

  return `${amount} ${currency}`
}

function formatDate(dateString) {
  if (!dateString) return '—'

  const date = new Date(dateString)
  return date.toLocaleString('es-ES', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}
</script>

<style scoped>
.pricing-panel {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.pricing-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 1rem;
}

.pricing-header h3 {
  margin: 0;
  color: #1b3a6b;
  font-size: 1.5rem;
}

.subtitle {
  margin: 0.5rem 0 0;
  color: #64748b;
  font-size: 0.95rem;
}

.help-text {
  color: #64748b;
  font-size: 0.875rem;
  margin-top: 0.25rem;
  margin-bottom: 0;
}

/* Card */
.card {
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 1.5rem;
}

.card h4 {
  margin: 0 0 0.75rem;
  color: #1b3a6b;
  font-size: 1.125rem;
}

/* Calculator Section */
.calculator-section {
  background: #f8fafc;
  border-color: #cbd5e1;
}

.calculator-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-top: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  font-weight: 600;
  color: #1b3a6b;
  font-size: 0.875rem;
}

.form-control {
  padding: 0.5rem 0.75rem;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-size: 0.875rem;
  transition: border-color 0.2s;
}

.form-control:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-control:disabled {
  background-color: #f1f5f9;
  cursor: not-allowed;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

/* Calculation Result */
.calculation-result {
  margin-top: 1.5rem;
  padding: 1rem;
  background: white;
  border: 2px solid #10b981;
  border-radius: 8px;
}

.calculation-result h5 {
  margin: 0 0 1rem;
  color: #059669;
  font-size: 1rem;
}

.result-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.result-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  padding: 0.75rem;
  background: #f8fafc;
  border-radius: 6px;
}

.result-item.total {
  background: #ecfdf5;
  border: 1px solid #10b981;
}

.result-item.discount {
  background: #fef9c3;
  border: 1px solid #eab308;
}

.result-item.discount .result-value {
  color: #a16207;
}

.result-label {
  font-size: 0.875rem;
  color: #64748b;
  font-weight: 500;
}

.result-value {
  font-size: 1.25rem;
  color: #1b3a6b;
  font-weight: 700;
}

.result-item.total .result-value {
  color: #059669;
}

/* Prices Table */
.prices-section {
  background: white;
}

.table-wrapper {
  overflow-x: auto;
  margin-top: 1rem;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;
}

.data-table thead {
  background: #f8fafc;
}

.data-table th {
  padding: 0.75rem;
  text-align: left;
  font-weight: 600;
  color: #1b3a6b;
  border-bottom: 2px solid #e2e8f0;
  white-space: nowrap;
}

.data-table td {
  padding: 0.75rem;
  border-bottom: 1px solid #e2e8f0;
}

.data-table tbody tr:hover {
  background: #f8fafc;
}

.sku-code {
  display: inline-block;
  padding: 0.25rem 0.5rem;
  background: #f1f5f9;
  border-radius: 4px;
  font-size: 0.8125rem;
  font-family: 'Courier New', monospace;
  color: #1b3a6b;
}

.attributes-cell {
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.price-value {
  font-weight: 600;
  color: #059669;
  font-size: 1rem;
}

.loading-price {
  color: #94a3b8;
  font-size: 0.875rem;
}

/* Buttons */
.btn {
  padding: 0.625rem 1.25rem;
  border-radius: 6px;
  font-weight: 600;
  font-size: 0.875rem;
  border: none;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #2563eb;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.btn-success {
  background: #10b981;
  color: white;
}

.btn-success:hover:not(:disabled) {
  background: #059669;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);
}

.btn-secondary {
  background: #64748b;
  color: white;
}

.btn-secondary:hover:not(:disabled) {
  background: #475569;
}

.btn-outline {
  background: white;
  color: #3b82f6;
  border: 1px solid #3b82f6;
}

.btn-outline:hover:not(:disabled) {
  background: #eff6ff;
}

.btn-sm {
  padding: 0.375rem 0.75rem;
  font-size: 0.8125rem;
}

/* Pills */
.pill {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.pill.active {
  background: #d1fae5;
  color: #065f46;
}

.pill.inactive {
  background: #fee2e2;
  color: #991b1b;
}

/* Loading States */
.loading-inline {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem;
  color: #64748b;
}

.spinner-small {
  width: 20px;
  height: 20px;
  border: 3px solid #e2e8f0;
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 3rem 1rem;
  color: #64748b;
}

.empty-icon {
  font-size: 3rem;
  display: block;
  margin-bottom: 1rem;
}

.empty-state p {
  margin: 0.5rem 0;
}

.empty-hint {
  font-size: 0.875rem;
  color: #94a3b8;
}

/* Alerts */
.alert-error {
  padding: 1rem;
  background: #fee2e2;
  border: 1px solid #f87171;
  border-radius: 6px;
  color: #991b1b;
  margin-top: 1rem;
}

.info-message {
  display: flex;
  gap: 1rem;
  padding: 1rem;
  background: #eff6ff;
  border: 1px solid #93c5fd;
  border-radius: 6px;
  margin-top: 1rem;
}

.info-icon {
  font-size: 1.5rem;
}

/* Modal */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  padding: 1rem;
}

.modal-content {
  background: white;
  border-radius: 8px;
  max-width: 900px;
  width: 100%;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
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
}

.modal-close {
  background: none;
  border: none;
  font-size: 1.5rem;
  color: #64748b;
  cursor: pointer;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
}

.modal-close:hover {
  background: #f1f5f9;
  color: #1b3a6b;
}

.modal-body {
  padding: 1.5rem;
  overflow-y: auto;
  flex: 1;
}

.modal-footer {
  padding: 1.5rem;
  border-top: 1px solid #e2e8f0;
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
}

/* History Table Specific */
.history-table {
  margin-top: 1rem;
}

.client-id {
  font-family: 'Courier New', monospace;
  font-size: 0.8125rem;
  color: #64748b;
}

.rules-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
}

.rule-tag {
  display: inline-block;
  padding: 0.125rem 0.5rem;
  background: #dbeafe;
  color: #1e40af;
  border-radius: 4px;
  font-size: 0.75rem;
}

.text-muted {
  color: #94a3b8;
  font-style: italic;
}

/* Responsive */
@media (max-width: 768px) {
  .pricing-header {
    flex-direction: column;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .result-grid {
    grid-template-columns: 1fr;
  }

  .data-table {
    font-size: 0.8125rem;
  }

  .data-table th,
  .data-table td {
    padding: 0.5rem;
  }
}
</style>
