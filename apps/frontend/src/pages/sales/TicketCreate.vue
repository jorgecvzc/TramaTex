<template>
  <Navbar />
  <div class="ticket-create-container">
    <div class="page-header">
      <div>
        <button class="btn-back" @click="goBack">← Volver</button>
        <h1>Nuevo Ticket (Factura Simplificada)</h1>
        <p class="subtitle">Cliente: CONSUMIDOR FINAL</p>
      </div>
    </div>

    <div class="form-card">
      <form @submit.prevent="handleSubmit">
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

          <div v-else class="line-items-list">
            <div
              v-for="(item, index) in formData.lineItems"
              :key="index"
              class="line-item-card"
            >
              <div class="line-item-header">
                <span class="line-item-number">Línea {{ index + 1 }}</span>
                <button
                  type="button"
                  class="btn-remove"
                  @click="removeLineItem(index)"
                  title="Eliminar línea"
                >
                  ✕
                </button>
              </div>
              <div class="line-item-fields">
                <div class="form-group field-variant">
                  <label>Variante de Producto *</label>
                  <input
                    v-model="item.productVariantId"
                    type="text"
                    placeholder="UUID de la variante"
                    class="form-input"
                    required
                  />
                </div>
                <div class="form-group">
                  <label>Cantidad *</label>
                  <input
                    v-model.number="item.quantity"
                    type="number"
                    min="1"
                    class="form-input"
                    required
                  />
                </div>
                <div class="form-group">
                  <label>Precio Unit. (EUR) *</label>
                  <input
                    v-model.number="item.unitPrice"
                    type="number"
                    step="0.01"
                    min="0"
                    class="form-input"
                    required
                  />
                </div>
                <div class="form-group">
                  <label>Descuento (EUR)</label>
                  <input
                    v-model.number="item.discountAmount"
                    type="number"
                    step="0.01"
                    min="0"
                    placeholder="Opcional"
                    class="form-input"
                  />
                </div>
              </div>
              <div class="line-item-summary">
                <span class="label">Subtotal línea:</span>
                <span class="value">{{ formatMoney(calculateLineTotal(item)) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Total Summary -->
        <div v-if="formData.lineItems.length > 0" class="total-summary">
          <div class="total-row">
            <span class="label">Subtotal:</span>
            <span class="value">{{ formatMoney(calculateSubtotal()) }}</span>
          </div>
          <div class="total-row">
            <span class="label">IVA (21%):</span>
            <span class="value">{{ formatMoney(calculateTax()) }}</span>
          </div>
          <div class="total-row total">
            <span class="label">TOTAL:</span>
            <span class="value">{{ formatMoney(calculateTotal()) }}</span>
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
            {{ isSubmitting ? 'Creando...' : 'Crear Ticket' }}
          </button>
        </div>
      </form>

      <!-- Error Display -->
      <div v-if="submitError" class="error-box">
        {{ submitError }}
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import Navbar from '@/components/layout/Navbar.vue';
import salesApi from '@/services/salesApi.js';

const router = useRouter();

// CONSUMIDOR_FINAL UUID (from backend seed)
const CONSUMIDOR_FINAL_ID = '00000000-0000-0000-0000-000000000001';

const formData = ref({
  lineItems: [],
});

const isSubmitting = ref(false);
const submitError = ref('');

const isFormValid = computed(() => {
  return (
    formData.value.lineItems.length > 0 &&
    formData.value.lineItems.every(
      (item) =>
        item.productVariantId &&
        item.quantity > 0 &&
        item.unitPrice >= 0
    )
  );
});

function addLineItem() {
  formData.value.lineItems.push({
    productVariantId: '',
    quantity: 1,
    unitPrice: 0,
    discountAmount: null,
  });
}

function removeLineItem(index) {
  formData.value.lineItems.splice(index, 1);
}

function calculateLineTotal(item) {
  const subtotal = item.quantity * item.unitPrice;
  const discount = item.discountAmount || 0;
  return Math.max(0, subtotal - discount);
}

function calculateSubtotal() {
  return formData.value.lineItems.reduce((total, item) => {
    return total + calculateLineTotal(item);
  }, 0);
}

function calculateTax() {
  const subtotal = calculateSubtotal();
  return subtotal * 0.21; // 21% IVA
}

function calculateTotal() {
  return calculateSubtotal() + calculateTax();
}

function formatMoney(amount) {
  return new Intl.NumberFormat('es-ES', {
    style: 'currency',
    currency: 'EUR',
  }).format(amount);
}

async function handleSubmit() {
  if (!isFormValid.value || isSubmitting.value) return;

  isSubmitting.value = true;
  submitError.value = '';

  try {
    // Prepare line items for API
    const lineItems = formData.value.lineItems.map((item) => {
      const lineItem = {
        productVariantId: item.productVariantId,
        quantity: item.quantity,
        unitPrice: {
          amount: item.unitPrice,
          currency: 'EUR',
        },
      };

      if (item.discountAmount && item.discountAmount > 0) {
        lineItem.discountAmount = {
          amount: item.discountAmount,
          currency: 'EUR',
        };
      }

      return lineItem;
    });

    // Create simplified invoice (ticket)
    const newInvoice = await salesApi.createSimplifiedInvoice({ lineItems });

    // Navigate to invoice detail (if route exists) or back to invoice list
    router.push(`/sales/invoices/${newInvoice.id}`);
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
</script>

<style scoped>
.ticket-create-container {
  padding: 2rem;
  max-width: 900px;
  margin: 0 auto;
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

.line-items-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.line-item-card {
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  padding: 1rem;
  background: #f9fafb;
}

.line-item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.line-item-number {
  font-size: 0.875rem;
  font-weight: 600;
  color: #4b5563;
}

.btn-remove {
  background: transparent;
  border: none;
  color: #dc2626;
  cursor: pointer;
  font-size: 1.25rem;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  transition: background 0.2s;
}

.btn-remove:hover {
  background: rgba(220, 38, 38, 0.1);
}

.line-item-fields {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr;
  gap: 1rem;
}

.field-variant {
  grid-column: span 2;
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

@media (max-width: 768px) {
  .line-item-fields {
    grid-template-columns: 1fr;
  }

  .field-variant {
    grid-column: span 1;
  }
}
</style>
