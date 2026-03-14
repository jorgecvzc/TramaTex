<template>
  <Navbar />
  <div class="ticket-create-container">
    <div class="page-header">
      <div>
        <button class="btn-back" @click="goBack">← Volver</button>
        <h1>Nuevo Ticket (Factura Simplificada)</h1>
      </div>
    </div>

    <div class="form-card">
      <form @submit.prevent="handleSubmit">
        <!-- Client Selector -->
        <div class="form-section">
          <div class="client-selector-row">
            <div class="client-selector-field">
              <PartySelector
                v-model="selectedPartyId"
                label="Cliente"
                placeholder="Buscar cliente... (por defecto: Consumidor Final)"
                role-filter="CLIENT"
                @select="onPartySelected"
              />
            </div>
            <button
              v-if="selectedPartyId !== CONSUMIDOR_FINAL_ID"
              type="button"
              class="btn btn-secondary btn-reset-client"
              @click="resetToConsumidorFinal"
              title="Volver a Consumidor Final"
            >
              ↺ Consumidor Final
            </button>
          </div>
        </div>

        <!-- Line Items Section -->
        <div class="form-section">
          <div class="section-header">
            <h2>Líneas del Ticket</h2>
            <button type="button" class="btn btn-secondary" @click="addLineItem">
              + Agregar Línea
            </button>
          </div>

          <div v-if="formData.lineItems.length === 0" class="empty-state">
            <p>No hay líneas agregadas. Agregue al menos una línea para crear el ticket.</p>
          </div>

          <div v-else class="line-items-table-wrapper">
            <table class="line-items-table">
              <thead>
                <tr>
                  <th class="col-num">#</th>
                  <th class="col-variant">Producto / Variante</th>
                  <th class="col-qty">Cantidad</th>
                  <th class="col-price">P. Unit.</th>
                  <th class="col-discount">Dto. %</th>
                  <th class="col-total">Total</th>
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
                  <td class="col-qty">
                    <input
                      v-model.number="item.quantity"
                      type="number"
                      min="1"
                      class="form-input"
                      required
                      @change="recalculateItemPrice(index)"
                    />
                  </td>
                  <td class="col-price">
                    <span v-if="item.unitPrice != null">{{ formatEur(item.unitPrice) }}</span>
                    <span v-else-if="item.loadingPrice" class="price-loading">…</span>
                    <span v-else class="price-pending">—</span>
                  </td>
                  <td class="col-discount">
                    <input
                      v-model.number="item.discountPercent"
                      type="number"
                      min="0"
                      max="100"
                      step="0.5"
                      class="form-input"
                      placeholder="0"
                    />
                  </td>
                  <td class="col-total">
                    <span v-if="item.unitPrice != null">{{ formatEur(lineTotal(item)) }}</span>
                    <span v-else>—</span>
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

        <!-- Total Summary -->
        <div v-if="formData.lineItems.length > 0" class="total-summary">
          <div class="total-row">
            <span class="label">Líneas:</span>
            <span class="value">{{ formData.lineItems.length }}</span>
          </div>
          <div class="total-row">
            <span class="label">Subtotal (sin IVA):</span>
            <span class="value">{{ formatEur(ticketSubtotal) }}</span>
          </div>
          <div class="total-row" v-if="ticketTaxAmount > 0">
            <span class="label">IVA:</span>
            <span class="value">{{ formatEur(ticketTaxAmount) }}</span>
          </div>
          <div class="total-row total">
            <span class="label">Total estimado:</span>
            <span class="value">{{ formatEur(ticketTotal) }}</span>
          </div>
          <p class="pricing-note">Los precios se obtienen del motor de Pricing. El total final se confirma al crear el ticket.</p>
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
            {{ isSubmitting ? 'Creando...' : 'Crear Ticket' }}
          </button>
        </div>
      </form>

      <!-- Error Display -->
      <div v-if="submitError" class="error-box">
        {{ submitError }}
      </div>
    </div>

    <!-- Ticket Receipt Modal -->
    <div v-if="createdInvoice" class="modal-overlay" @click.self="closeReceipt">
      <div class="receipt-modal">
        <div class="receipt-header">
          <p class="receipt-brand">{{ issuerProfile.displayName }}</p>
          <p v-if="issuerProfile.taxId" class="receipt-issuer">{{ issuerProfile.taxLabel }}: {{ issuerProfile.taxId }}</p>
          <p v-if="issuerProfile.addressLine" class="receipt-issuer">{{ issuerProfile.addressLine }}</p>
          <p v-if="issuerProfile.cityLine" class="receipt-issuer">{{ issuerProfile.cityLine }}</p>
          <div class="receipt-divider"></div>
          <p class="receipt-title">TICKET {{ createdInvoice.invoiceNumber }}</p>
          <p class="receipt-date">{{ formatReceiptDate(createdInvoice.issueDate) }}</p>
          <div class="receipt-divider"></div>
        </div>

        <div class="receipt-body">
          <table class="receipt-table">
            <thead>
              <tr>
                <th class="rt-name">Producto</th>
                <th class="rt-qty">Cant.</th>
                <th class="rt-price">P. Unit.</th>
                <th class="rt-total">Total</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in createdInvoice.lineItems" :key="item.id || item.productVariantID">
                <td class="rt-name">
                  {{ item.productName || '—' }}
                  <small v-if="item.discountAmount && item.discountAmount.amount > 0" class="rt-discount">
                    (Dto. {{ salesApi.formatMoney(item.discountAmount) }})
                  </small>
                </td>
                <td class="rt-qty">{{ item.quantity }}</td>
                <td class="rt-price">{{ salesApi.formatUnitPrice(item.unitPrice) }}</td>
                <td class="rt-total">{{ salesApi.formatMoney(item.total) }}</td>
              </tr>
            </tbody>
          </table>

          <div class="receipt-divider"></div>

          <div class="receipt-totals">
            <div class="receipt-total-row">
              <span>Subtotal</span>
              <span>{{ salesApi.formatMoney(createdInvoice.subtotal) }}</span>
            </div>
            <div class="receipt-total-row">
              <span>IVA</span>
              <span>{{ salesApi.formatMoney(createdInvoice.taxAmount) }}</span>
            </div>
            <div class="receipt-total-row receipt-grand-total">
              <span>TOTAL</span>
              <span>{{ salesApi.formatMoney(createdInvoice.total) }}</span>
            </div>
          </div>

          <div class="receipt-divider"></div>
          <p class="receipt-footer-text">Cliente: {{ selectedPartyName }}</p>
          <p class="receipt-footer-text">Gracias por su compra</p>
        </div>

        <div class="receipt-actions no-print">
          <button class="btn btn-secondary" @click="printReceipt">🖨️ Imprimir</button>
          <button class="btn btn-primary" @click="newTicket">+ Nuevo Ticket</button>
          <button class="btn btn-secondary" @click="closeReceipt">Cerrar</button>
        </div>
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
            :product-id="null"
            initial-mode="quick"
            title=""
            description="Seleccione una variante de producto"
            :initial-query="variantSelectorQuery"
            @variant-selected="handleVariantSelected"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { useRouter } from 'vue-router';
import Navbar from '@/components/layout/Navbar.vue';
import VariantSelector from '@/components/product/VariantSelector.vue';
import PartySelector from '@/components/party/PartySelector.vue';
import salesApi from '@/services/salesApi';
import { productApi } from '@/services/productApi';
import { calculateFinalSalePrice } from '@/services/pricingApi';
import { getPrintIssuerProfile } from '@/services/printIssuerProfile';

const router = useRouter();
const issuerProfile = getPrintIssuerProfile();

// CONSUMIDOR_FINAL UUID (from backend seed)
const CONSUMIDOR_FINAL_ID = '00000000-0000-0000-0000-000000000001';

const selectedPartyId = ref(CONSUMIDOR_FINAL_ID);
const selectedPartyName = ref('CONSUMIDOR FINAL');
const partyDefaultDiscount = ref(null);

const formData = ref({
  lineItems: [],
});

const isSubmitting = ref(false);
const submitError = ref('');
const showVariantSelector = ref(false);
const editingLineIndex = ref(null);
const variantSelectorQuery = ref('');
const createdInvoice = ref(null);

function onPartySelected(party) {
  selectedPartyName.value = party?.name || 'CONSUMIDOR FINAL';
  partyDefaultDiscount.value = party?.default_discount_percentage || null;
}

function resetToConsumidorFinal() {
  selectedPartyId.value = CONSUMIDOR_FINAL_ID;
  selectedPartyName.value = 'CONSUMIDOR FINAL';
  partyDefaultDiscount.value = null;
}

// Recalculate all prices when client changes (pricing may differ per client)
watch(selectedPartyId, async (newId, oldId) => {
  if (!newId || newId === oldId) return;
  const items = formData.value.lineItems;
  const promises = items
    .map((_, idx) => items[idx].productVariantId ? fetchItemPrice(idx) : null)
    .filter(Boolean);
  await Promise.all(promises);
});

const isFormValid = computed(() => {
  return (
    formData.value.lineItems.length > 0 &&
    formData.value.lineItems.every(
      (item) =>
        item.productVariantId &&
        item.quantity > 0
    )
  );
});

function addLineItem() {
  formData.value.lineItems.push({
    productVariantId: '',
    selectedVariantName: '',
    quickSearchQuery: '',
    inlineSearchError: '',
    quantity: 1,
    unitPrice: null,
    taxRate: null,
    finalPriceWithTax: null,
    loadingPrice: false,
    discountPercent: 0,
  });
}

function lineTotal(item) {
  if (item.unitPrice == null) return 0;
  const disc = item.discountPercent || 0;
  return item.unitPrice * item.quantity * (1 - disc / 100);
}

const ticketSubtotal = computed(() => {
  return formData.value.lineItems.reduce((sum, item) => {
    if (item.unitPrice != null) return sum + lineTotal(item);
    return sum;
  }, 0);
});

const ticketTaxAmount = computed(() => {
  return formData.value.lineItems.reduce((sum, item) => {
    if (item.unitPrice != null && item.finalPriceWithTax != null) {
      return sum + (item.finalPriceWithTax - item.unitPrice) * item.quantity;
    }
    return sum;
  }, 0);
});

const ticketTotal = computed(() => ticketSubtotal.value + ticketTaxAmount.value);

function formatEur(amount) {
  if (amount == null) return '—';
  return new Intl.NumberFormat('es-ES', { style: 'currency', currency: 'EUR' }).format(amount);
}

async function fetchItemPrice(index) {
  const item = formData.value.lineItems[index];
  if (!item.productVariantId) return;
  item.loadingPrice = true;
  try {
    const result = await calculateFinalSalePrice(
      [{ productVariantId: item.productVariantId, quantity: item.quantity }],
      selectedPartyId.value,
      new Date()
    );
    const calcItem = result.calculatedItems?.[0];
    if (calcItem) {
      item.unitPrice = calcItem.baseSalesPrice?.amount ?? calcItem.finalPrice?.amount ?? null;
      item.taxRate = calcItem.taxRate ?? null;
      item.finalPriceWithTax = calcItem.finalPriceWithTax?.amount ?? null;
      // Pre-fill discount from party's configured default (exact value)
      if (partyDefaultDiscount.value && partyDefaultDiscount.value > 0) {
        item.discountPercent = partyDefaultDiscount.value;
      } else if (calcItem.discountPercent > 0) {
        item.discountPercent = calcItem.discountPercent;
      }
    }
  } catch (err) {
    console.warn('[TicketCreate] Price lookup failed:', err);
    item.unitPrice = null;
    item.taxRate = null;
    item.finalPriceWithTax = null;
  } finally {
    item.loadingPrice = false;
  }
}

async function recalculateItemPrice(index) {
  const item = formData.value.lineItems[index];
  if (item.productVariantId && item.quantity > 0) {
    await fetchItemPrice(index);
  }
}

function openVariantSelector(index) {
  editingLineIndex.value = index;
  variantSelectorQuery.value = '';
  showVariantSelector.value = true;
}

function handleVariantSelected(payload) {
  const variant = payload?.variant || payload;
  if (editingLineIndex.value !== null && variant) {
    const idx = editingLineIndex.value;
    const item = formData.value.lineItems[idx];
    item.productVariantId = variant.id;
    item.selectedVariantName = `${variant.sku} - ${variant.product_name || 'Producto'}`;
    item.quickSearchQuery = '';
    item.inlineSearchError = '';
    fetchItemPrice(idx);
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
      item.selectedVariantName = `${result.variant.sku} - ${result.product?.name || 'Producto'}`;
      item.quickSearchQuery = '';
      fetchItemPrice(index);
    } else if (result.type === 'no_match') {
      item.inlineSearchError = 'No encontrado. Usa 📋 para buscar en catálogo.';
    } else {
      editingLineIndex.value = index;
      variantSelectorQuery.value = query;
      showVariantSelector.value = true;
    }
  } catch (err) {
    console.error('[TicketCreate] Inline search error:', err);
    item.inlineSearchError = err.message || 'Error en búsqueda';
  }
}

function clearLineVariant(index) {
  const item = formData.value.lineItems[index];
  item.productVariantId = '';
  item.selectedVariantName = '';
  item.quickSearchQuery = '';
  item.inlineSearchError = '';
  item.unitPrice = null;
  item.taxRate = null;
  item.finalPriceWithTax = null;
}

function removeLineItem(index) {
  formData.value.lineItems.splice(index, 1);
}

async function handleSubmit() {
  if (!isFormValid.value || isSubmitting.value) return;

  isSubmitting.value = true;
  submitError.value = '';

  try {
    const items = formData.value.lineItems.map((item) => ({
      productVariantId: item.productVariantId,
      quantity: item.quantity,
      discountPercent: item.discountPercent || 0,
    }));

    const newInvoice = await salesApi.createSimplifiedInvoice({
      partyId: selectedPartyId.value,
      invoiceDate: new Date().toISOString(),
      items,
    });

    // Show the ticket receipt modal
    createdInvoice.value = newInvoice;
  } catch (err) {
    submitError.value = err?.message || 'No se pudo crear el ticket';
    console.error('Error creating ticket:', err);
  } finally {
    isSubmitting.value = false;
  }
}

function goBack() {
  router.push('/sales/invoices');
}

function formatReceiptDate(dateString) {
  if (!dateString) return '';
  const d = new Date(dateString);
  return d.toLocaleDateString('es-ES', { day: '2-digit', month: '2-digit', year: 'numeric' }) +
    ' ' + d.toLocaleTimeString('es-ES', { hour: '2-digit', minute: '2-digit' });
}

function printReceipt() {
  window.print();
}

function newTicket() {
  createdInvoice.value = null;
  formData.value.lineItems = [];
  submitError.value = '';
  selectedPartyId.value = CONSUMIDOR_FINAL_ID;
  selectedPartyName.value = 'CONSUMIDOR FINAL';
  partyDefaultDiscount.value = null;
}

function closeReceipt() {
  router.push('/sales/invoices');
}
</script>

<style scoped>
.ticket-create-container {
  padding: 1.5rem 2rem;
}

.page-header {
  margin-bottom: 2rem;
}

.page-header h1 {
  font-size: 2rem;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0.5rem 0 0.25rem;
}

.subtitle {
  font-size: 0.875rem;
  color: #6b7280;
  margin: 0;
}

.client-selector-row {
  display: flex;
  align-items: flex-end;
  gap: 0.75rem;
}

.client-selector-field {
  flex: 1;
  max-width: 450px;
}

.btn-reset-client {
  white-space: nowrap;
  margin-bottom: 0.125rem;
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

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.section-header h2 {
  font-size: 1.25rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
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
  min-width: 200px;
}

.line-items-table .col-qty {
  width: 80px;
}

.line-items-table .col-price,
.line-items-table .col-total {
  width: 100px;
  text-align: right;
  white-space: nowrap;
  font-variant-numeric: tabular-nums;
}

.line-items-table .col-discount {
  width: 80px;
}

.line-items-table thead .col-price,
.line-items-table thead .col-total {
  text-align: right;
}

.price-loading {
  color: #9ca3af;
}

.price-pending {
  color: #d1d5db;
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

.form-group {
  margin-bottom: 0;
}

.form-group label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: #4a5568;
  margin-bottom: 0.5rem;
}

.form-input {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-size: 0.875rem;
  font-family: inherit;
}

.form-input:focus {
  outline: none;
  border-color: #E6B800;
  box-shadow: 0 0 0 3px rgba(230, 184, 0, 0.1);
}

.line-item-summary {
  display: flex;
  justify-content: flex-end;
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid #e5e7eb;
  font-size: 0.875rem;
}

.line-item-summary .label {
  color: #6b7280;
  margin-right: 0.5rem;
}

.line-item-summary .value {
  font-weight: 600;
  color: #1f2937;
}

.total-summary {
  background: #f9fafb;
  border-radius: 6px;
  padding: 1.5rem;
  margin-bottom: 2rem;
}

.total-row {
  display: flex;
  justify-content: space-between;
  padding: 0.5rem 0;
  font-size: 0.875rem;
}

.total-row.total {
  margin-top: 0.5rem;
  padding-top: 0.75rem;
  border-top: 2px solid #e5e7eb;
  font-weight: 600;
  font-size: 1.125rem;
}

.total-row .label {
  color: #6b7280;
}

.total-row .value {
  color: #1f2937;
  font-weight: 500;
}

.total-row.total .label,
.total-row.total .value {
  color: #1f2937;
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

.error-box {
  margin-top: 1rem;
  padding: 1rem;
  background: #fee2e2;
  border: 1px solid #fecaca;
  border-radius: 4px;
  color: #991b1b;
  font-size: 0.875rem;
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

/* ── Receipt Modal ── */
.receipt-modal {
  background: white;
  border-radius: 8px;
  width: 380px;
  max-width: 95vw;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
  font-family: 'Courier New', Courier, monospace;
}

.receipt-header {
  padding: 1.5rem 1.5rem 0;
  text-align: center;
}

.receipt-brand {
  font-size: 1.1rem;
  font-weight: 900;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: #002395;
  margin: 0 0 0.25rem;
}

.receipt-issuer {
  font-size: 0.72rem;
  color: #64748b;
  margin: 0;
  line-height: 1.4;
}

.receipt-title {
  font-size: 1rem;
  font-weight: 700;
  color: #1e293b;
  margin: 0.5rem 0 0.15rem;
}

.receipt-date {
  font-size: 0.8rem;
  color: #64748b;
  margin: 0;
}

.receipt-divider {
  border: none;
  border-top: 1px dashed #cbd5e1;
  margin: 0.75rem 0;
}

.receipt-body {
  padding: 0 1.5rem 1rem;
}

.receipt-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.8rem;
}

.receipt-table thead th {
  font-weight: 600;
  font-size: 0.72rem;
  text-transform: uppercase;
  color: #4b5563;
  padding: 0.35rem 0.25rem;
  border-bottom: 1px solid #e5e7eb;
  text-align: left;
}

.receipt-table tbody td {
  padding: 0.35rem 0.25rem;
  vertical-align: top;
  color: #1e293b;
}

.rt-name {
  width: 45%;
}

.rt-qty {
  width: 12%;
  text-align: center !important;
}

.rt-price {
  width: 22%;
  text-align: right !important;
}

.rt-total {
  width: 21%;
  text-align: right !important;
}

.rt-discount {
  display: block;
  color: #dc2626;
  font-size: 0.7rem;
}

.receipt-totals {
  font-size: 0.85rem;
}

.receipt-total-row {
  display: flex;
  justify-content: space-between;
  padding: 0.25rem 0;
  color: #4b5563;
}

.receipt-grand-total {
  font-weight: 800;
  font-size: 1.1rem;
  color: #1e293b;
  padding-top: 0.5rem;
  margin-top: 0.25rem;
  border-top: 2px solid #1e293b;
}

.receipt-footer-text {
  text-align: center;
  font-size: 0.75rem;
  color: #64748b;
  margin: 0.25rem 0;
}

.receipt-actions {
  display: flex;
  gap: 0.75rem;
  justify-content: center;
  padding: 1rem 1.5rem 1.5rem;
  border-top: 1px solid #e5e7eb;
}

/* ── Print mode for receipt ── */
@media print {
  .ticket-create-container > .page-header,
  .ticket-create-container > .form-card,
  .navbar,
  nav,
  .receipt-actions {
    display: none !important;
  }

  .ticket-create-container {
    padding: 0 !important;
  }

  .modal-overlay {
    display: block !important;
    position: static;
    background: none;
  }

  .receipt-modal {
    display: block !important;
    box-shadow: none;
    border-radius: 0;
    width: 80mm;
    max-width: 80mm;
    margin: 0 auto;
  }
}
</style>
