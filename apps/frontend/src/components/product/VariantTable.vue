<template>
  <div class="variant-table">
    <div class="table-header">
      <div>
        <h3>Variantes del Producto</h3>
        <p class="table-description">
          Listado de todas las variantes configuradas con sus precios base de venta.
        </p>
      </div>
      <div class="header-actions">
        <button @click="openAddVariantModal" class="btn btn-primary" :disabled="isLoading">
          <Plus :size="16" />
          Añadir Variante
        </button>
        <button @click="$emit('refresh')" class="btn btn-secondary" :disabled="isLoading">
          <RefreshCw :size="16" />
          Actualizar
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="loading">
      <div class="spinner"></div>
      <p>Cargando variantes...</p>
    </div>

    <!-- Empty State -->
    <div v-if="!isLoading && variants.length === 0" class="empty-state">
      <Package :size="64" class="empty-icon" />
      <p>No hay variantes creadas para este producto.</p>
      <p class="empty-hint">
        Las variantes se crean automáticamente (Just-in-Time) cuando se añaden a una orden,
        o puedes generarlas manualmente según las combinaciones de atributos disponibles.
      </p>
      <button class="btn btn-primary" @click="openBatchCreator">
        <Zap :size="16" />
        Generar variantes
      </button>
    </div>

    <!-- Variants Table -->
    <div v-if="!isLoading && variants.length > 0" class="table-wrapper">
      <table>
        <thead>
          <tr>
            <th>SKU Variante</th>
            <th>Configuración</th>
            <th>Código de barras</th>
            <th>Estado</th>
            <th class="align-right">Costo Base</th>
            <th class="align-right">Precio Base Venta</th>
            <th class="align-center">Acciones</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="variant in variants"
            :key="variant.id"
            :class="{ inactive: !variant.is_active }"
          >
            <td>
              <code class="sku-code">{{ variant.sku }}</code>
            </td>
            <td>
              <div class="attribute-values">
                <span
                  v-for="(value, attr) in variant.attribute_values"
                  :key="attr"
                  class="attribute-tag"
                >
                  <span class="attr-name">{{ attr }}:</span>
                  <span class="attr-value">{{ value }}</span>
                </span>
                <span v-if="!variant.attribute_values || Object.keys(variant.attribute_values).length === 0" class="text-muted">
                  Sin atributos
                </span>
              </div>
            </td>
            <td>
              <code v-if="variant.barcode" class="barcode">{{ variant.barcode }}</code>
              <span v-else class="text-muted">—</span>
            </td>
            <td>
              <span class="pill" :class="getStatusClass(variant.status)">
                {{ formatStatus(variant.status) }}
              </span>
              <span v-if="variant.is_active" class="pill active">Activo</span>
              <span v-else class="pill inactive">Inactivo</span>
            </td>
            <td class="align-right">
              <span v-if="variant.base_cost !== undefined && variant.base_cost !== null" class="price base-cost">
                {{ formatPrice(variant.base_cost) }}
              </span>
              <span v-else class="text-muted">—</span>
            </td>
            <td class="align-right">
              <span v-if="variantPrices[variant.id]" class="price sale-price">
                {{ formatPrice(variantPrices[variant.id]) }}
              </span>
              <span v-else-if="loadingPrices[variant.id]" class="loading-price">
                <div class="spinner-small"></div>
              </span>
              <span v-else class="text-muted">—</span>
            </td>
            <td class="align-center">
              <button
                @click="viewVariantDetail(variant)"
                class="btn-icon"
                title="Ver detalles"
              >
                <Eye :size="16" />
              </button>
              <button
                @click="openEditVariantModal(variant)"
                class="btn-icon"
                title="Editar variante"
              >
                <Edit2 :size="16" />
              </button>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Summary Footer -->
      <div class="table-footer">
        <p class="footer-info">
          Mostrando {{ variants.length }} variante{{ variants.length !== 1 ? 's' : '' }}
        </p>
      </div>
    </div>

    <!-- Variant Form Modal -->
    <VariantFormModal
      v-if="showVariantForm"
      :product-id="productId"
      :product-sku="productSku"
      :variant="editingVariant"
      @close="closeVariantForm"
      @saved="handleVariantSaved"
    />

    <!-- Batch Creator Modal -->
    <VariantBatchCreator
      v-if="showBatchCreator"
      :product-id="productId"
      :product-sku="productSku"
      @close="closeBatchCreator"
      @created="handleBatchCreated"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { Plus, RefreshCw, Package, Zap, Eye, Edit2 } from 'lucide-vue-next'
import { productApi } from '@/services/productApi'
import { pricingApi } from '@/services/pricingApi'
import VariantFormModal from './VariantFormModal.vue'
import VariantBatchCreator from './VariantBatchCreator.vue'

const props = defineProps({
  productId: {
    type: String,
    required: true,
  },
  productSku: {
    type: String,
    required: true,
  },
  variants: {
    type: Array,
    default: () => [],
  },
  isLoading: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['refresh'])

// State
const variantPrices = ref({})
const loadingPrices = ref({})
const showVariantForm = ref(false)
const editingVariant = ref(null)
const showBatchCreator = ref(false)

// Lifecycle
onMounted(() => {
  fetchPricesForVariants()
})

// Watchers
watch(() => props.variants, () => {
  fetchPricesForVariants()
}, { deep: true })

// Methods
async function fetchPricesForVariants() {
  if (!props.variants || props.variants.length === 0) return

  // Mark all as loading
  props.variants.forEach(variant => {
    loadingPrices.value[variant.id] = true
  })

  // Fetch prices using real Pricing API
  for (const variant of props.variants) {
    try {
      const result = await pricingApi.calculateBaseSalesPrice(props.productId, variant.id)
      variantPrices.value[variant.id] = {
        amount: result.baseSalesPrice.amount,
        currency: result.baseSalesPrice.currency
      }
    } catch (err) {
      console.error(`Error fetching price for variant ${variant.id}:`, err)
      variantPrices.value[variant.id] = null
    } finally {
      loadingPrices.value[variant.id] = false
    }
  }
}

function formatPrice(price) {
  if (!price || !price.amount) return '—'
  return new Intl.NumberFormat('es-ES', {
    style: 'currency',
    currency: price.currency || 'EUR'
  }).format(price.amount)
}

function formatStatus(status) {
  const statusMap = {
    PROVISIONAL: 'Provisional',
    CONFIRMED: 'Confirmado',
  }
  return statusMap[status] || status
}

function getStatusClass(status) {
  return status === 'CONFIRMED' ? 'status-confirmed' : 'status-provisional'
}

function viewVariantDetail(variant) {
  // TODO: Implement variant detail view
  console.log('View variant:', variant)
  alert(`Vista de detalle de variante:\n\nSKU: ${variant.sku}\n\nEsta funcionalidad se implementará en futuras iteraciones.`)
}

function openAddVariantModal() {
  editingVariant.value = null
  showVariantForm.value = true
}

function openEditVariantModal(variant) {
  editingVariant.value = variant
  showVariantForm.value = true
}

function closeVariantForm() {
  showVariantForm.value = false
  editingVariant.value = null
}

function handleVariantSaved() {
  // Refresh the variant list
  emit('refresh')
}

function openBatchCreator() {
  showBatchCreator.value = true
}

function closeBatchCreator() {
  showBatchCreator.value = false
}

function handleBatchCreated(result) {
  console.log('Batch created:', result)
  if (result.created && result.created.length > 0) {
    alert(`✅ Se crearon ${result.created.length} variante(s) exitosamente`)
  }
  if (result.errors && result.errors.length > 0) {
    console.error('Errors creating variants:', result.errors)
  }
  emit('refresh')
}
</script>

<style scoped>
.variant-table {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 1rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #e2e8f0;
}

.header-actions {
  display: flex;
  gap: 0.75rem;
}

.table-header h3 {
  color: #1b3a6b;
  margin: 0 0 0.5rem 0;
  font-size: 1.2rem;
}

.table-description {
  color: #64748b;
  margin: 0;
  font-size: 0.9rem;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  gap: 1rem;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid rgba(27, 58, 107, 0.12);
  border-top-color: #1b3a6b;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.spinner-small {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(27, 58, 107, 0.12);
  border-top-color: #1b3a6b;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
  display: inline-block;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  text-align: center;
  color: #64748b;
  gap: 1rem;
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 0.5rem;
  opacity: 0.5;
}

.empty-state p {
  margin: 0.25rem 0;
  font-size: 0.95rem;
}

.empty-hint {
  font-size: 0.85rem;
  color: #94a3b8;
  max-width: 600px;
  line-height: 1.5;
}

.table-wrapper {
  overflow-x: auto;
  border-radius: 10px;
  border: 1px solid #e2e8f0;
}

table {
  width: 100%;
  border-collapse: collapse;
  background: #ffffff;
}

thead {
  background: #f8fafc;
  border-bottom: 2px solid #e2e8f0;
}

th {
  text-align: left;
  padding: 0.75rem 1rem;
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: #475569;
}

th.align-right {
  text-align: right;
}

th.align-center {
  text-align: center;
}

tbody tr {
  border-bottom: 1px solid #e2e8f0;
  transition: background 0.15s ease;
}

tbody tr:hover {
  background: #f8fafc;
}

tbody tr.inactive {
  background: rgba(220, 38, 38, 0.04);
  border-left: 3px solid rgba(220, 38, 38, 0.4);
}

tbody tr.inactive:hover {
  background: rgba(220, 38, 38, 0.08);
}

tbody tr.inactive td {
  color: #94a3b8;
}

tbody tr.inactive .sku-code {
  background: rgba(220, 38, 38, 0.1);
  color: #dc2626;
}

tbody tr.inactive .pill.inactive {
  background-color: rgba(220, 38, 38, 0.1);
  color: #dc2626;
  font-weight: 700;
}

td {
  padding: 0.875rem 1rem;
  font-size: 0.9rem;
  color: #1e293b;
}

td.align-right {
  text-align: right;
}

td.align-center {
  text-align: center;
}

.sku-code {
  background: #f1f5f9;
  color: #475569;
  padding: 0.25rem 0.6rem;
  border-radius: 4px;
  font-size: 0.8rem;
  font-family: 'Monaco', 'Menlo', monospace;
  font-weight: 600;
}

.barcode {
  background: #fef3c7;
  color: #92400e;
  padding: 0.25rem 0.6rem;
  border-radius: 4px;
  font-size: 0.8rem;
  font-family: 'Monaco', 'Menlo', monospace;
  font-weight: 600;
}

.attribute-values {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.attribute-tag {
  display: inline-flex;
  align-items: center;
  background: #f1f5f9;
  padding: 0.2rem 0.6rem;
  border-radius: 4px;
  font-size: 0.8rem;
}

.attr-name {
  color: #64748b;
  font-weight: 600;
  margin-right: 0.25rem;
}

.attr-value {
  color: #1e293b;
  font-weight: 500;
}

.text-muted {
  color: #94a3b8;
  font-style: italic;
  font-size: 0.85rem;
}

.pill {
  display: inline-block;
  padding: 0.25rem 0.6rem;
  border-radius: 999px;
  font-weight: 600;
  font-size: 0.7rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-right: 0.5rem;
}

.pill.status-confirmed {
  background-color: rgba(34, 197, 94, 0.1);
  color: #22c55e;
}

.pill.status-provisional {
  background-color: rgba(251, 191, 36, 0.1);
  color: #f59e0b;
}

.pill.active {
  background-color: rgba(34, 197, 94, 0.1);
  color: #22c55e;
}

.pill.inactive {
  background-color: rgba(148, 163, 184, 0.1);
  color: #94a3b8;
}

.price {
  font-weight: 700;
  color: #1b3a6b;
  font-size: 1rem;
}

.loading-price {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 60px;
}

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
  gap: 0.5rem;
}

.btn-primary {
  background: #f4d03f;
  color: #1e293b;
  font-weight: 700;
}

.btn-primary:hover:not(:disabled) {
  background: #e6c530;
  box-shadow: 0 2px 8px rgba(244, 208, 63, 0.3);
}

.btn-secondary {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  color: #1e293b;
}

.btn-secondary:hover:not(:disabled) {
  background: #f8fafc;
}

.btn-secondary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-icon {
  background: none;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
  padding: 0.25rem;
  opacity: 0.6;
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.btn-icon:hover {
  opacity: 1;
  transform: scale(1.15);
}

.table-footer {
  padding: 1rem;
  background: #f8fafc;
  border-top: 1px solid #e2e8f0;
}

.footer-info {
  margin: 0;
  font-size: 0.85rem;
  color: #64748b;
}

/* Modal Styles */
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
}

.modal-content {
  background: #ffffff;
  border-radius: 12px;
  max-width: 500px;
  width: 90%;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
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

.modal-body p {
  margin: 0.5rem 0;
  color: #475569;
}

.modal-footer {
  padding: 1rem 1.5rem;
  border-top: 1px solid #e2e8f0;
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
}

@media (max-width: 768px) {
  .table-wrapper {
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
  }

  table {
    min-width: 800px;
  }

  .table-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .attribute-values {
    flex-direction: column;
  }
}

/* Price styling */
.price.base-cost {
  color: #64748b;
  font-weight: 500;
}

.price.sale-price {
  color: #16a34a;
  font-weight: 700;
}
</style>
