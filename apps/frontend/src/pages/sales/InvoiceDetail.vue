<template>
  <Navbar />
  <div class="invoice-detail-container">
    <!-- Loading State -->
    <div v-if="isLoading" class="loading-state">
      <div class="spinner"></div>
      <p>Cargando factura...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-state">
      <p class="error-message">{{ error }}</p>
      <button class="btn btn-secondary" @click="fetchInvoice">Reintentar</button>
      <button class="btn btn-secondary" @click="goBack">Volver</button>
    </div>

    <!-- Invoice Detail -->
    <div v-else-if="invoice" class="invoice-content">
      <!-- Header -->
      <div class="page-header">
        <div>
          <button class="btn-back" @click="goBack">← Volver</button>
          <h1>Factura {{ invoice.invoiceNumber }}</h1>
          <span :class="['type-badge', `type-${invoice.type.toLowerCase()}`]">
            {{ getTypeLabel(invoice.type) }}
          </span>
        </div>
        <div class="header-actions">
          <button
            class="btn btn-secondary"
            @click="printInvoice"
            title="Imprimir factura"
          >
            🖨️ Imprimir
          </button>
        </div>
      </div>

      <div class="print-doc-header">
        <p class="print-brand">{{ issuerProfile.displayName }}</p>
        <p v-if="issuerProfile.taxId" class="print-issuer-line">{{ issuerProfile.taxLabel }}: {{ issuerProfile.taxId }}</p>
        <p v-if="issuerProfile.addressLine || issuerProfile.cityLine" class="print-issuer-line">{{ issuerProfile.addressLine }}<span v-if="issuerProfile.addressLine && issuerProfile.cityLine"> · </span>{{ issuerProfile.cityLine }}</p>
        <p v-if="issuerProfile.contactLine" class="print-issuer-line">{{ issuerProfile.contactLine }}</p>
        <h2>Factura {{ invoice.invoiceNumber }}</h2>
        <div class="print-doc-meta">
          <span>Cliente: {{ formatPartyId(invoice.partyId) }}</span>
          <span>Fecha emisión: {{ formatDate(invoice.issueDate) }}</span>
          <span>Tipo: {{ getTypeLabel(invoice.type) }}</span>
        </div>
      </div>

      <!-- Info Cards -->
      <div class="info-grid">
        <div class="info-card">
          <h3>Información General</h3>
          <div class="info-row">
            <span class="label">Cliente:</span>
            <span class="value">{{ formatPartyId(invoice.partyId) }}</span>
          </div>
          <div class="info-row">
            <span class="label">Fecha de Emisión:</span>
            <span class="value">{{ formatDate(invoice.issueDate) }}</span>
          </div>
          <div class="info-row">
            <span class="label">Fecha de Vencimiento:</span>
            <span class="value">{{ formatDate(invoice.dueDate) }}</span>
          </div>
        </div>

        <div class="info-card">
          <h3>Totales</h3>
          <div class="info-row">
            <span class="label">Subtotal:</span>
            <span class="value">{{ salesApi.formatMoney(invoice.subtotal) }}</span>
          </div>
          <div class="info-row">
            <span class="label">IVA:</span>
            <span class="value">{{ salesApi.formatMoney(invoice.taxAmount) }}</span>
          </div>
          <div class="info-row total">
            <span class="label">Total:</span>
            <span class="value">{{ salesApi.formatMoney(invoice.total) }}</span>
          </div>
        </div>
      </div>

      <!-- Payment Terms -->
      <div v-if="invoice.paymentTerms" class="notes-card">
        <h3>Condiciones de Pago</h3>
        <p>{{ invoice.paymentTerms }}</p>
      </div>

      <!-- Related Orders -->
      <div v-if="invoice.salesOrderIds && invoice.salesOrderIds.length > 0" class="related-section">
        <h3>Pedidos Relacionados</h3>
        <ul class="related-list">
          <li v-for="orderId in invoice.salesOrderIds" :key="orderId">
            <a :href="`/sales/orders/${orderId}`" class="related-link">
              {{ formatId(orderId) }}
            </a>
          </li>
        </ul>
      </div>

      <!-- Related Delivery Notes -->
      <div v-if="invoice.deliveryNoteIds && invoice.deliveryNoteIds.length > 0" class="related-section">
        <h3>Albaranes Relacionados</h3>
        <ul class="related-list">
          <li v-for="noteId in invoice.deliveryNoteIds" :key="noteId">
            {{ formatId(noteId) }}
          </li>
        </ul>
      </div>

      <!-- Line Items -->
      <div class="line-items-section">
        <h2>Líneas de la Factura</h2>

        <div v-if="!invoice.lineItems || invoice.lineItems.length === 0" class="empty-state">
          <p>No hay líneas en esta factura</p>
        </div>

        <div v-else class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>Variante</th>
                <th>Cantidad</th>
                <th>Precio Unitario</th>
                <th>Descuento</th>
                <th>Subtotal</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in invoice.lineItems" :key="item.id || item.productVariantID">
                <td class="variant-id">{{ formatId(item.productVariantID) }}</td>
                <td>{{ item.quantity }}</td>
                <td>{{ salesApi.formatMoney(item.unitPrice) }}</td>
                <td>{{ item.discountAmount ? salesApi.formatMoney(item.discountAmount) : '—' }}</td>
                <td class="amount">{{ salesApi.formatMoney(item.subtotal) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="print-doc-footer">
        <span>Documento generado por {{ issuerProfile.displayName }}</span>
        <span>Vencimiento: {{ formatDate(invoice.dueDate) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import Navbar from '@/components/layout/Navbar.vue';
import salesApi from '@/services/salesApi';
import { getPrintIssuerProfile } from '@/services/printIssuerProfile';

const route = useRoute();
const router = useRouter();

const invoice = ref(null);
const isLoading = ref(false);
const error = ref('');
const issuerProfile = getPrintIssuerProfile();

onMounted(() => {
  fetchInvoice();
});

async function fetchInvoice() {
  const invoiceId = route.params.id;
  if (!invoiceId) {
    error.value = 'ID de factura no válido';
    return;
  }

  isLoading.value = true;
  error.value = '';

  try {
    invoice.value = await salesApi.getInvoice(invoiceId);
  } catch (err) {
    error.value = err?.message || 'No se pudo cargar la factura';
    console.error('Error loading invoice:', err);
  } finally {
    isLoading.value = false;
  }
}

function goBack() {
  router.push('/sales/invoices');
}

function printInvoice() {
  window.print();
}

function formatDate(dateString) {
  if (!dateString) return '—';
  const date = new Date(dateString);
  return date.toLocaleDateString('es-ES', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
}

function formatPartyId(partyId) {
  if (!partyId) return '—';
  return partyId.substring(0, 8) + '...';
}

function formatId(id) {
  if (!id) return '—';
  return id.substring(0, 8) + '...';
}

function getTypeLabel(type) {
  const labels = {
    STANDARD: 'Estándar',
    SIMPLIFIED: 'Simplificada',
  };
  return labels[type] || type;
}
</script>

<style scoped>
.invoice-detail-container {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

.loading-state,
.error-state {
  text-align: center;
  padding: 3rem 1rem;
  background: white;
  border-radius: 8px;
}

.spinner {
  width: 40px;
  height: 40px;
  margin: 0 auto 1rem;
  border: 3px solid #f3f4f6;
  border-top-color: #E6B800;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.error-message {
  color: #dc2626;
  margin-bottom: 1rem;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
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

.btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary {
  background: #f3f4f6;
  color: #4a5568;
}

.btn-secondary:hover {
  background: #e5e7eb;
}

.header-actions {
  display: flex;
  gap: 0.5rem;
}

.type-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.025em;
}

.type-standard {
  background: #dbeafe;
  color: #1e40af;
}

.type-simplified {
  background: #fef3c7;
  color: #92400e;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.info-card,
.notes-card,
.related-section {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.info-card h3,
.notes-card h3,
.related-section h3 {
  font-size: 1rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 1rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid #f3f4f6;
}

.info-row {
  display: flex;
  justify-content: space-between;
  padding: 0.5rem 0;
  font-size: 0.875rem;
}

.info-row.total {
  margin-top: 0.5rem;
  padding-top: 0.75rem;
  border-top: 2px solid #f3f4f6;
  font-weight: 600;
  font-size: 1rem;
}

.info-row .label {
  color: #6b7280;
}

.info-row .value {
  color: #1f2937;
  font-weight: 500;
}

.notes-card {
  margin-bottom: 2rem;
}

.notes-card p {
  font-size: 0.875rem;
  color: #4b5563;
  line-height: 1.6;
  margin: 0;
}

.related-section {
  margin-bottom: 2rem;
}

.related-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.related-list li {
  padding: 0.5rem 0;
  border-bottom: 1px solid #f3f4f6;
}

.related-list li:last-child {
  border-bottom: none;
}

.related-link {
  font-family: 'Courier New', monospace;
  color: #002395;
  text-decoration: none;
  transition: color 0.2s;
}

.related-link:hover {
  color: #E6B800;
  text-decoration: underline;
}

.line-items-section {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.line-items-section h2 {
  font-size: 1.25rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 1.5rem;
}

.table-container {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table thead {
  background: #f9fafb;
}

.data-table th {
  text-align: left;
  padding: 0.75rem 1rem;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  color: #6b7280;
  letter-spacing: 0.05em;
}

.data-table td {
  padding: 1rem;
  border-top: 1px solid #f3f4f6;
  font-size: 0.875rem;
  color: #1f2937;
}

.variant-id {
  font-family: 'Courier New', monospace;
  color: #6b7280;
}

.amount {
  font-weight: 600;
  text-align: right;
}

.empty-state {
  text-align: center;
  padding: 2rem;
  color: #9ca3af;
}

.print-doc-header,
.print-doc-footer {
  display: none;
}

.print-brand {
  margin: 0;
  font-size: 0.875rem;
  font-weight: 700;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #374151;
}

.print-doc-header h2 {
  margin: 0.35rem 0;
  font-size: 1.3rem;
  font-weight: 700;
  color: #111827;
}

.print-issuer-line {
  margin: 0;
  font-size: 0.75rem;
  color: #4b5563;
}

.print-doc-meta {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  font-size: 0.8rem;
  color: #4b5563;
}

@media print {
  :deep(.navbar),
  :deep(nav),
  .page-header,
  .btn-back,
  .header-actions,
  .btn,
  .related-link {
    display: none !important;
  }

  .invoice-detail-container {
    padding: 0;
    max-width: none;
  }

  .print-doc-header,
  .print-doc-footer {
    display: block;
    border: 1px solid #d1d5db;
    padding: 0.75rem 1rem;
    margin-bottom: 0.75rem;
    background: white;
  }

  .print-doc-footer {
    margin-top: 0.75rem;
    margin-bottom: 0;
    display: flex;
    justify-content: space-between;
    font-size: 0.75rem;
    color: #4b5563;
  }

  .info-card,
  .notes-card,
  .related-section,
  .line-items-section,
  .table-container {
    box-shadow: none !important;
    border: 1px solid #d1d5db;
  }

  .type-badge {
    border: 1px solid #9ca3af;
    color: #111827 !important;
    background: transparent !important;
  }

  .data-table {
    font-size: 0.75rem;
  }
}
</style>
