<template>
  <Navbar />
  <div class="order-create-container">
    <div class="page-header">
      <div>
        <button class="btn-back" @click="goBack">← Volver</button>
        <h1>Nuevo Pedido</h1>
      </div>
    </div>

    <div class="form-card">
      <form @submit.prevent="handleSubmit">
        <!-- Customer Selection -->
        <div class="form-section">
          <h2>Cliente</h2>
          <PartySelector
            v-model="formData.partyId"
            label="Cliente"
            placeholder="Buscar cliente por nombre o referencia..."
            role-filter="CLIENT"
            :required="true"
            help-text="Seleccione el cliente para este pedido"
            @select="onPartySelected"
          />
        </div>

        <!-- Order Details -->
        <div class="form-section">
          <h2>Detalles del Pedido</h2>
          <div class="form-row">
            <div class="form-group">
              <label for="orderDate">Fecha de Pedido *</label>
              <input
                id="orderDate"
                v-model="formData.orderDate"
                type="date"
                class="form-input"
                required
              />
            </div>
            <div class="form-group">
              <label for="deliveryDate">Fecha de Entrega *</label>
              <input
                id="deliveryDate"
                v-model="formData.deliveryDate"
                type="date"
                class="form-input"
                :min="minDeliveryDate"
                required
              />
            </div>
          </div>
          <div class="form-group">
            <label for="notes">Notas</label>
            <textarea
              id="notes"
              v-model="formData.notes"
              class="form-textarea"
              rows="3"
              placeholder="Notas adicionales sobre el pedido..."
            ></textarea>
          </div>
        </div>

        <!-- MES Work Definitions (Document-level) -->
        <div v-if="formData.partyId" class="form-section">
          <h2>Trabajos MES</h2>
          <p class="help-text">Asocie trabajos MES a este pedido. Para cada trabajo puede añadir observaciones.</p>
          <div class="mes-selector">
            <div v-if="isLoadingMesWorks" class="mes-loading">Cargando trabajos MES...</div>
            <div v-else-if="mesWorks.length === 0" class="mes-empty">No hay trabajos MES disponibles para este cliente.</div>
            <div v-else class="mes-ref-list">
              <div v-for="work in mesWorks" :key="work.id" class="mes-ref-item">
                <label class="mes-checkbox-label">
                  <input
                    type="checkbox"
                    :checked="isMesWorkSelected(work.id)"
                    @change="toggleMesWork(work.id)"
                  />
                  <span>{{ work.work_number }} - {{ work.work_name }}</span>
                </label>
                <textarea
                  v-if="isMesWorkSelected(work.id)"
                  class="mes-observations"
                  rows="2"
                  placeholder="Observaciones para este trabajo MES..."
                  :value="getMesWorkObservations(work.id)"
                  @input="setMesWorkObservations(work.id, $event.target.value)"
                ></textarea>
              </div>
            </div>
          </div>
        </div>

        <!-- Line Items -->
        <div class="form-section">
          <div class="section-header">
            <h2>Líneas del Pedido</h2>
            <button type="button" class="btn btn-secondary" @click="addLineItem">
              + Agregar Línea
            </button>
          </div>

          <p class="help-text">El precio final de venta e IVA se calculan automáticamente en Pricing al crear el pedido.</p>

          <div v-if="formData.lineItems.length === 0" class="empty-state">
            <p>No hay líneas agregadas. Agregue al menos una línea para crear el pedido.</p>
          </div>

          <div v-else class="line-items-table-wrapper">
            <table class="line-items-table">
              <thead>
                <tr>
                  <th class="col-num">#</th>
                  <th class="col-variant">Referencia</th>
                  <th class="col-product-name">Nombre</th>
                  <th class="col-qty">Cant.</th>
                  <th class="col-list-price">P. Tarifa</th>
                  <th class="col-price">Precio</th>
                  <th class="col-discount">Dto. %</th>
                  <th class="col-subtotal">Subtotal</th>
                  <th class="col-actions"></th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(item, index) in formData.lineItems" :key="index" class="line-item-row">
                  <td class="col-num">{{ index + 1 }}</td>
                  <td class="col-variant">
                    <div class="variant-inline-search">
                      <input
                        v-if="!item.productVariantId"
                        v-model="item.quickSearchQuery"
                        type="text"
                        class="form-input"
                        placeholder="SKU o código de barras..."
                        @keyup.enter="inlineSmartSearch(index)"
                      />
                      <span
                        v-else
                        class="variant-selected-label"
                        @click="clearLineVariant(index)"
                        title="Haz clic para cambiar"
                      >
                        {{ item.selectedVariantName }}
                      </span>
                      <button
                        type="button"
                        class="btn-browse-variant"
                        @click="openVariantSelector(index)"
                        title="Buscar en catálogo"
                      >
                        📋
                      </button>
                    </div>
                    <small v-if="item.inlineSearchError" class="inline-search-error">
                      {{ item.inlineSearchError }}
                    </small>
                    <input v-model="item.productVariantId" type="hidden" required />
                  </td>
                  <td class="col-product-name">
                    <span class="product-name-readonly">{{ buildDisplayName(item) }}</span>
                  </td>
                  <td class="col-qty">
                    <input
                      v-model.number="item.quantity"
                      type="number"
                      min="1"
                      class="form-input"
                      required
                      @input="calculateTotals"
                    />
                  </td>
                  <td class="col-list-price">
                    <input
                      :value="item.listPrice != null ? item.listPrice.toFixed(3) : ''"
                      type="text"
                      class="form-input input-readonly"
                      readonly
                      tabindex="-1"
                      placeholder="—"
                    />
                  </td>
                  <td class="col-price">
                    <input
                      v-model.number="item.unitPrice"
                      type="number"
                      step="0.001"
                      min="0"
                      class="form-input"
                      @input="calculateTotals"
                    />
                  </td>
                  <td class="col-discount">
                    <input
                      v-model.number="item.discountPercent"
                      type="number"
                      step="0.01"
                      min="0"
                      max="100"
                      placeholder="0.00"
                      class="form-input"
                      @input="calculateTotals"
                    />
                  </td>
                  <td class="col-subtotal">
                    <span v-if="isPreviewLoading" class="subtotal-loading">…</span>
                    <span v-else class="subtotal-value">{{ formatMoney(calculateLineSubtotal(item)) }}</span>
                  </td>
                  <td class="col-actions">
                    <button
                      type="button"
                      class="btn-remove"
                      @click="removeLineItem(index)"
                      title="Eliminar línea"
                    >
                      ✕
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Totals Summary -->
        <div v-if="formData.lineItems.length > 0" class="totals-section">
          <h3>Resumen de Totales</h3>
          <div class="totals-grid">
            <div class="total-row">
              <span class="total-label">Subtotal:</span>
              <span class="total-value">{{ formatMoney(calculatedTotals.subtotal) }}</span>
            </div>
            <div class="total-row">
              <span class="total-label">IVA:</span>
              <span class="total-value">{{ formatMoney(calculatedTotals.tax) }}</span>
            </div>
            <div class="total-row total-final">
              <span class="total-label">Total:</span>
              <span class="total-value">{{ formatMoney(calculatedTotals.total) }}</span>
            </div>
          </div>
        </div>

        <!-- Form Actions -->
        <div class="form-actions">
          <button type="button" class="btn btn-secondary" @click="goBack">
            Cancelar
          </button>
          <button
            type="submit"
            class="btn btn-primary"
            :disabled="!isFormValid || isSubmitting"
          >
            {{ isSubmitting ? "Creando..." : "Crear Pedido" }}
          </button>
        </div>
      </form>

      <!-- Error Display -->
      <div v-if="submitError" class="error-box">
        {{ submitError }}
      </div>
    </div>

    <!-- Variant Selector Modal -->
    <div v-if="showVariantSelector" class="modal-overlay" @click.self="showVariantSelector = false">
      <div class="modal-content">
        <div class="modal-header">
          <h3>Seleccionar Variante de Producto</h3>
          <button class="btn-close" @click="showVariantSelector = false">✕</button>
        </div>
        <div class="modal-body">
          <VariantSelector
            :key="variantSelectorQuery + '-' + editingLineIndex"
            :product-id="null"
            :initial-query="variantSelectorQuery"
            initial-mode="quick"
            title=""
            description="Seleccione una variante de producto"
            @variant-selected="handleVariantSelected"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';
import Navbar from '@/components/layout/Navbar.vue';
import PartySelector from '@/components/party/PartySelector.vue';
import VariantSelector from '@/components/product/VariantSelector.vue';
import salesApi from '@/services/salesApi';
import { productApi } from '@/services/productApi';
import { calculateBaseSalesPrice } from '@/services/pricingApi';
import { mesApi } from '@/services/mesApi';

const router = useRouter();

const formData = ref({
  partyId: '',
  orderDate: '',
  deliveryDate: '',
  notes: '',
  mesWorkRefs: [],
  lineItems: [],
});

const isSubmitting = ref(false);
const submitError = ref('');
const mesWorks = ref([]);
const isLoadingMesWorks = ref(false);
const partyDefaultDiscount = ref(null);

const calculatedTotals = ref({
  subtotal: 0,
  tax: 0,
  total: 0,
});

const previewResult = ref(null);
const isPreviewLoading = ref(false);
let previewDebounceTimer = null;

const minDeliveryDate = computed(() => {
  const tomorrow = new Date();
  tomorrow.setDate(tomorrow.getDate() + 1);
  return tomorrow.toISOString().split('T')[0];
});

const isFormValid = computed(() => {
  return (
    formData.value.partyId &&
    formData.value.orderDate &&
    formData.value.deliveryDate &&
    formData.value.lineItems.length > 0 &&
    formData.value.lineItems.every(
      (item) => item.productVariantId && item.quantity > 0
    )
  );
});

onMounted(() => {
  // Set default order date to today
  const today = new Date().toISOString().split('T')[0];
  formData.value.orderDate = today;
  formData.value.deliveryDate = today;

  // Pre-fill from quote if navigated from QuoteDetail
  const fromQuoteRaw = sessionStorage.getItem('orderFromQuote');
  if (fromQuoteRaw) {
    sessionStorage.removeItem('orderFromQuote');
    try {
      const fromQuote = JSON.parse(fromQuoteRaw);
      if (fromQuote.partyId) formData.value.partyId = fromQuote.partyId;
      if (fromQuote.notes) formData.value.notes = fromQuote.notes;
      if (fromQuote.mesWorkRefs?.length) formData.value.mesWorkRefs = fromQuote.mesWorkRefs;
      if (fromQuote.lineItems?.length) {
        formData.value.lineItems = fromQuote.lineItems.map(item => ({
          productVariantId: item.productVariantId || '',
          productId: '',
          selectedVariantName: item.selectedVariantName || '',
          productName: item.productName || '',
          productDescription: '',
          optionConfiguration: item.optionConfiguration || {},
          quickSearchQuery: '',
          inlineSearchError: '',
          quantity: item.quantity || 1,
          listPrice: item.listPrice ?? null,
          unitPrice: item.unitPrice ?? null,
          discountPercent: item.discountPercent ?? null,
          productBasePrice: null,
        }));
        calculateTotals();
      }
    } catch (e) {
      console.error('Error parsing quote pre-fill data:', e);
    }
  }
});

watch(() => formData.value.partyId, (partyId) => {
  loadMesWorksForParty(partyId);
});

function onPartySelected(party) {
  partyDefaultDiscount.value = party?.default_discount_percentage || null;
}

const showVariantSelector = ref(false);
const editingLineIndex = ref(null);
const variantSelectorQuery = ref('');

function addLineItem() {
  formData.value.lineItems.push({
    productVariantId: '',
    productId: '',
    selectedVariantName: '',
    productName: '',
    productDescription: '',
    optionConfiguration: {},
    quickSearchQuery: '',
    inlineSearchError: '',
    quantity: 1,
    listPrice: null,
    unitPrice: null,
    discountPercent: partyDefaultDiscount.value || null,
    productBasePrice: null,
  });
}

async function loadMesWorksForParty(partyId) {
  if (!partyId) {
    mesWorks.value = [];
    formData.value.mesWorkRefs = [];
    return;
  }

  isLoadingMesWorks.value = true;
  try {
    mesWorks.value = await mesApi.listWorkDefinitions({ party_id: partyId });
  } catch (error) {
    console.error('Error loading MES work definitions for selected party:', error);
    mesWorks.value = [];
  } finally {
    isLoadingMesWorks.value = false;
  }
}

function isMesWorkSelected(workId) {
  return formData.value.mesWorkRefs.some(r => r.mesWorkId === workId);
}

function toggleMesWork(workId) {
  const idx = formData.value.mesWorkRefs.findIndex(r => r.mesWorkId === workId);
  if (idx >= 0) {
    formData.value.mesWorkRefs.splice(idx, 1);
  } else {
    formData.value.mesWorkRefs.push({ mesWorkId: workId, observations: '' });
  }
}

function getMesWorkObservations(workId) {
  return formData.value.mesWorkRefs.find(r => r.mesWorkId === workId)?.observations || '';
}

function setMesWorkObservations(workId, value) {
  const ref = formData.value.mesWorkRefs.find(r => r.mesWorkId === workId);
  if (ref) ref.observations = value;
}

function openVariantSelector(index) {
  editingLineIndex.value = index;
  const item = formData.value.lineItems[index];
  variantSelectorQuery.value = item?.quickSearchQuery?.trim() || '';
  showVariantSelector.value = true;
}

function handleVariantSelected(payload) {
  const variant = payload?.variant || payload

  if (editingLineIndex.value !== null && variant) {
    const item = formData.value.lineItems[editingLineIndex.value];
    item.productVariantId = variant.id;
    item.productId = variant.product_id || '';
    item.selectedVariantName = variant.sku;
    item.productName = variant.product_name || '';
    item.optionConfiguration = variant.option_configuration || {};
    item.productDescription = variant.product_description || '';
    item.quickSearchQuery = '';
    item.inlineSearchError = '';
    item.productBasePrice = variant.product_base_price ?? variant.base_cost ?? null;
    // Immediately show base_cost so price is never blank
    if (item.productBasePrice != null && item.unitPrice == null) {
      item.listPrice = item.productBasePrice;
      item.unitPrice = item.productBasePrice;
    }
    fetchPriceForLineItem(item);
  }
  showVariantSelector.value = false;
  editingLineIndex.value = null;
}

async function inlineSmartSearch(index) {
  const item = formData.value.lineItems[index];
  const query = item.quickSearchQuery?.trim();
  if (!query) return;
  
  item.inlineSearchError = '';
  
  try {
    const result = await productApi.smartSearch(query);
    
    if (result.type === 'exact_variant' && result.variant) {
      item.productVariantId = result.variant.id;
      item.productId = result.variant.product_id || result.product?.id || '';
      item.selectedVariantName = result.variant.sku;
      item.productName = result.product?.name || '';
      item.optionConfiguration = result.variant.option_configuration || {};
      item.productDescription = result.product?.description || '';
      item.quickSearchQuery = '';
      item.productBasePrice = result.product?.base_price ?? result.variant.base_cost ?? null;
      if (item.productBasePrice != null && item.unitPrice == null) {
        item.listPrice = item.productBasePrice;
        item.unitPrice = item.productBasePrice;
      }
      fetchPriceForLineItem(item);
      calculateTotals();
    } else if (result.type === 'no_match') {
      item.inlineSearchError = 'No encontrado. Usa 📋 para buscar en catálogo.';
    } else {
      // Ambiguous or partial result → open modal with query pre-filled
      editingLineIndex.value = index;
      variantSelectorQuery.value = query;
      showVariantSelector.value = true;
    }
  } catch (err) {
    console.error('[OrderCreate] Inline search error:', err);
    item.inlineSearchError = err.message || 'Error en búsqueda';
  }
}

function clearLineVariant(index) {
  const item = formData.value.lineItems[index];
  item.productVariantId = '';
  item.productId = '';
  item.selectedVariantName = '';
  item.productName = '';
  item.optionConfiguration = {};
  item.productDescription = '';
  item.quickSearchQuery = '';
  item.inlineSearchError = '';
  item.listPrice = null;
  item.unitPrice = null;
  item.discountPercent = null;
  item.productBasePrice = null;
  calculateTotals();
}

/**
 * Fetch sale price from Pricing Engine (ADR-015) and auto-fill listPrice + unitPrice.
 * Sale price = base_price + attribute modifiers + brand markup.
 * Falls back to product base_price when pricing engine is unavailable.
 */
async function fetchPriceForLineItem(item) {
  if (!item.productVariantId || !item.productId) {
    // Without productId we cannot call the pricing engine
    if (item.productBasePrice != null) {
      item.listPrice = item.productBasePrice;
      item.unitPrice = item.productBasePrice;
      calculateTotals();
    }
    return;
  }
  try {
    const result = await calculateBaseSalesPrice(item.productId, item.productVariantId);
    const rawPrice = result.baseSalesPrice?.amount ?? item.productBasePrice ?? null;
    const salePrice = rawPrice != null ? Math.round(rawPrice * 1000) / 1000 : null;
    item.listPrice = salePrice;
    item.unitPrice = salePrice;
    calculateTotals();
  } catch (err) {
    console.warn('[OrderCreate] Error fetching sale price:', err.message);
    // Pricing engine unavailable — fall back to product base price
    if (item.productBasePrice != null) {
      item.listPrice = item.productBasePrice;
      item.unitPrice = item.productBasePrice;
      calculateTotals();
    }
  }
}

function removeLineItem(index) {
  formData.value.lineItems.splice(index, 1);
  calculateTotals();
}

function calculateLineSubtotal(item) {
  if (!previewResult.value) return 0;
  const match = previewResult.value.lineItems.find(
    li => li.productVariantId === item.productVariantId
  );
  return match?.subtotal?.amount ?? 0;
}

function calculateTotals() {
  clearTimeout(previewDebounceTimer);
  previewDebounceTimer = setTimeout(fetchPreviewCalculation, 400);
}

function buildPreviewItems() {
  return formData.value.lineItems
    .filter(item => item.productVariantId)
    .map(item => ({
      productVariantId: item.productVariantId,
      quantity: item.quantity || 1,
      unitPrice: item.unitPrice ? { amount: item.unitPrice, currency: 'EUR' } : undefined,
      discountPercent: item.discountPercent != null ? item.discountPercent : undefined,
    }));
}

async function fetchPreviewCalculation() {
  const partyId = formData.value.partyId;
  const items = buildPreviewItems();
  if (!partyId || items.length === 0) {
    previewResult.value = null;
    calculatedTotals.value = { subtotal: 0, tax: 0, total: 0 };
    return;
  }
  isPreviewLoading.value = true;
  try {
    previewResult.value = await salesApi.previewOrderCalculation(partyId, items);
    calculatedTotals.value = {
      subtotal: previewResult.value.subtotal?.amount ?? 0,
      tax: previewResult.value.taxAmount?.amount ?? 0,
      total: previewResult.value.total?.amount ?? 0,
    };
  } catch (err) {
    console.warn('[OrderCreate] Preview calculation error:', err.message);
  } finally {
    isPreviewLoading.value = false;
  }
}

function buildDisplayName(item) {
  if (!item.productName) return '';
  const config = item.optionConfiguration;
  if (!config || Object.keys(config).length === 0) return item.productName;
  return item.productName + ' - ' + Object.values(config).join(', ');
}

function formatMoney(amount) {
  return new Intl.NumberFormat('es-ES', {
    style: 'currency',
    currency: 'EUR',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(amount);
}

async function handleSubmit() {
  if (!isFormValid.value || isSubmitting.value) return;

  isSubmitting.value = true;
  submitError.value = '';

  try {
    // Prepare line items
    const lineItems = formData.value.lineItems.map((item) => {
      const lineItem = {
        productVariantId: item.productVariantId,
        quantity: item.quantity,
      };

      if (item.unitPrice != null && item.unitPrice > 0) {
        lineItem.unitPrice = {
          amount: item.unitPrice,
          currency: 'EUR',
        };
      }

      if (item.discountPercent != null) {
        lineItem.discountPercent = item.discountPercent;
      }

      return lineItem;
    });

    // Prepare order data — field names must match backend CreateOrderCommand json tags
    const orderData = {
      partyId: formData.value.partyId,
      deliveryDate: salesApi.formatDateForAPI(new Date(formData.value.deliveryDate)),
      items: lineItems,
    };

    if (formData.value.mesWorkRefs.length > 0) {
      orderData.mesWorkRefs = formData.value.mesWorkRefs;
    }

    if (formData.value.notes) {
      orderData.notes = formData.value.notes;
    }

    // Create order
    const newOrder = await salesApi.createOrder(orderData);

    // Navigate to order detail
    router.push(`/sales/orders/${newOrder.id}`);
  } catch (err) {
    submitError.value = err?.message || 'No se pudo crear el pedido';
    console.error('Error creating order:', err);
  } finally {
    isSubmitting.value = false;
  }
}

function goBack() {
  router.push('/sales/orders');
}
</script>

<style scoped>
.input-readonly {
  background-color: #f3f4f6;
  cursor: not-allowed;
  color: #6b7280;
}

.order-create-container {
  padding: 1.5rem 2rem;
}

.page-header {
  margin-bottom: 2rem;
}

.page-header h1 {
  font-size: 2rem;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0.5rem 0;
}

.btn-back {
  background: transparent;
  border: none;
  color: #6b7280;
  cursor: pointer;
  padding: 0.25rem 0.5rem;
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
  transition: color 0.2s;
}

.btn-back:hover {
  color: #1f2937;
}

.form-card {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.form-section {
  margin-bottom: 2rem;
  padding-bottom: 2rem;
  border-bottom: 1px solid #f3f4f6;
}

.form-section:last-of-type {
  border-bottom: none;
}

.form-section h2 {
  font-size: 1.25rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 1.5rem;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.section-header h2 {
  margin: 0;
}

.form-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: #4a5568;
  margin-bottom: 0.5rem;
}

.form-input,
.form-textarea {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-size: 0.875rem;
  font-family: inherit;
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: #E6B800;
  box-shadow: 0 0 0 3px rgba(230, 184, 0, 0.1);
}

.help-text {
  display: block;
  font-size: 0.75rem;
  color: #9ca3af;
  margin-top: 0.25rem;
}

.empty-state {
  text-align: center;
  padding: 2rem;
  color: #9ca3af;
  background: #f9fafb;
  border-radius: 4px;
}

/* ── Line Items Table ── */
.line-items-table-wrapper {
  overflow-x: auto;
  margin-top: 0.5rem;
}

.line-items-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;
}

.line-items-table thead th {
  background: #f3f4f6;
  color: #4b5563;
  font-weight: 600;
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.03em;
  padding: 0.625rem 0.5rem;
  border-bottom: 2px solid #e5e7eb;
  text-align: left;
  white-space: nowrap;
}

.line-items-table tbody tr {
  border-bottom: 1px solid #f3f4f6;
  transition: background 0.15s;
}

.line-items-table tbody tr:hover {
  background: #f9fafb;
}

.line-items-table td {
  padding: 0.5rem;
  vertical-align: middle;
}

.line-items-table .col-num {
  width: 2.5rem;
  text-align: center;
  color: #9ca3af;
  font-weight: 600;
}

.line-items-table .col-variant {
  width: 280px;
  min-width: 240px;
}

.line-items-table .col-product-name {
  min-width: 200px;
}

.product-name-readonly {
  font-size: 0.8125rem;
  color: #374151;
}

.line-items-table .col-mes {
  min-width: 160px;
}

.line-items-table .col-qty,
.line-items-table .col-list-price,
.line-items-table .col-price,
.line-items-table .col-discount {
  width: 90px;
}

.line-items-table .col-subtotal {
  width: 100px;
  text-align: right;
}

.line-items-table .col-actions {
  width: 40px;
  text-align: center;
}

.line-items-table .form-input {
  width: 100%;
  padding: 0.375rem 0.5rem;
  font-size: 0.8125rem;
}

.btn-remove {
  background: transparent;
  border: none;
  color: #dc2626;
  cursor: pointer;
  font-size: 1.125rem;
  padding: 0.25rem 0.4rem;
  border-radius: 4px;
  transition: background 0.2s;
  line-height: 1;
}

.btn-remove:hover {
  background: rgba(220, 38, 38, 0.1);
}

.subtotal-value {
  font-size: 0.875rem;
  color: #1f2937;
  font-weight: 600;
  white-space: nowrap;
}

.totals-section {
  background: #f9fafb;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 2rem;
}

.totals-section h3 {
  font-size: 1rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 1rem;
}

.totals-grid {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.total-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.875rem;
}

.total-row.total-final {
  margin-top: 0.5rem;
  padding-top: 0.75rem;
  border-top: 2px solid #d1d5db;
  font-size: 1.125rem;
  font-weight: 700;
}

.total-label {
  color: #6b7280;
}

.total-row.total-final .total-label {
  color: #1f2937;
}

.total-value {
  color: #1f2937;
  font-weight: 600;
}

.total-row.total-final .total-value {
  color: #E6B800;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 2rem;
}

.btn {
  padding: 0.625rem 1.25rem;
  border: none;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: #E6B800;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #d4a700;
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background: #f3f4f6;
  color: #4a5568;
}

.btn-secondary:hover {
  background: #e5e7eb;
}
.btn-select-variant {
  width: 100%;
  padding: 0.625rem 0.875rem;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  background: white;
  text-align: left;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 0.875rem;
}

.btn-select-variant:hover {
  border-color: #3b82f6;
  background: #f9fafb;
}

.variant-inline-search {
  display: flex;
  gap: 0.35rem;
  align-items: center;
}

.variant-inline-search .form-input {
  flex: 1;
}

.variant-selected-label {
  flex: 1;
  padding: 0.5rem 0.75rem;
  background: #f0fdf4;
  border: 1px solid #86efac;
  border-radius: 4px;
  font-size: 0.85rem;
  color: #166534;
  cursor: pointer;
  overflow: hidden;
  text-overflow: ellipsis;
  display: flex;
  flex-direction: column;
  gap: 0.15rem;
}

.variant-description {
  display: block;
  color: #6b7280;
  font-size: 0.75rem;
  font-weight: 400;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.variant-selected-label:hover {
  background: #dcfce7;
  border-color: #22c55e;
}

.btn-browse-variant {
  background: #f8fafc;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  padding: 0.5rem 0.65rem;
  cursor: pointer;
  font-size: 1rem;
  transition: all 0.15s;
}

.btn-browse-variant:hover {
  background: #e2e8f0;
  border-color: #94a3b8;
}

.inline-search-error {
  color: #dc2626;
  font-size: 0.75rem;
  margin-top: 0.2rem;
}

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
  background: white;
  border-radius: 8px;
  width: 90%;
  max-width: 900px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
}

.modal-header {
  padding: 1.5rem;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  color: #1b3a6b;
}

.btn-close {
  background: none;
  border: none;
  font-size: 1.5rem;
  color: #6b7280;
  cursor: pointer;
  padding: 0;
  line-height: 1;
  transition: color 0.2s;
}

.btn-close:hover {
  color: #dc2626;
}

.modal-body {
  padding: 1.5rem;
  overflow-y: auto;
}

.error-box {
  margin-top: 1rem;
  padding: 1rem;
  background: #fee2e2;
  border: 1px solid #fecaca;
  border-radius: 4px;
  color: #991b1b;
  font-size: 0.875rem;
}

@media (max-width: 1024px) {
  .line-items-table .col-mes {
    display: none;
  }
  .line-items-table thead th.col-mes {
    display: none;
  }
}

@media (max-width: 768px) {
  .line-items-table thead {
    display: none;
  }
  .line-items-table,
  .line-items-table tbody,
  .line-items-table tr,
  .line-items-table td {
    display: block;
    width: 100%;
  }
  .line-items-table tr {
    border: 1px solid #e5e7eb;
    border-radius: 6px;
    padding: 0.75rem;
    margin-bottom: 0.75rem;
    background: #f9fafb;
  }
  .line-items-table td {
    padding: 0.25rem 0;
  }
  .line-items-table td::before {
    content: attr(data-label);
    font-weight: 600;
    font-size: 0.75rem;
    color: #6b7280;
    display: block;
    margin-bottom: 0.15rem;
  }
  .line-items-table .col-num {
    display: none;
  }
}

/* MES selector styles */
.mes-selector {
  margin-top: 0.5rem;
}

.mes-loading,
.mes-empty {
  color: #6b7280;
  font-size: 0.875rem;
  font-style: italic;
  padding: 0.5rem 0;
}

.mes-checkboxes {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.mes-ref-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.mes-ref-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.mes-observations {
  margin-left: 1.75rem;
  padding: 0.5rem;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 0.875rem;
  resize: vertical;
}

.mes-checkbox-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 0.75rem;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.875rem;
  transition: background 0.15s, border-color 0.15s;
}

.mes-checkbox-label:hover {
  background: #f0f7ff;
  border-color: #93c5fd;
}

.mes-checkbox-label input[type="checkbox"] {
  accent-color: #2563eb;
}
</style>
