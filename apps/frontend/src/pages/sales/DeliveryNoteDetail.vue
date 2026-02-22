<template>
  <Navbar />
  <div class="delivery-note-detail-container">
    <!-- Loading State -->
    <div v-if="isLoading" class="loading-state">
      <div class="spinner"></div>
      <p>Cargando albarán...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-state">
      <p class="error-message">{{ error }}</p>
      <button class="btn btn-secondary" @click="fetchDeliveryNote">Reintentar</button>
      <button class="btn btn-secondary" @click="goBack">Volver</button>
    </div>

    <!-- Delivery Note Detail -->
    <div v-else-if="deliveryNote" class="delivery-note-content">
      <!-- Header -->
      <div class="page-header">
        <div>
          <button class="btn-back" @click="goBack">← Volver</button>
          <h1>Albarán {{ deliveryNote.deliveryNoteNumber }}</h1>
          <span class="subtitle">
            Pedido: 
            <a 
              href="#" 
              @click.prevent="navigateToOrder(deliveryNote.salesOrderId)" 
              class="order-link"
            >
              {{ formatOrderId(deliveryNote.salesOrderId) }}
            </a>
          </span>
        </div>
        <div class="header-actions">
          <button 
            class="btn btn-secondary" 
            @click="printDeliveryNote"
            title="Imprimir albarán (PDF)"
          >
            🖨️ Imprimir Albarán
          </button>
        </div>
      </div>

      <div class="print-doc-header">
        <p class="print-brand">{{ issuerProfile.displayName }}</p>
        <p v-if="issuerProfile.taxId" class="print-issuer-line">{{ issuerProfile.taxLabel }}: {{ issuerProfile.taxId }}</p>
        <p v-if="issuerProfile.addressLine || issuerProfile.cityLine" class="print-issuer-line">{{ issuerProfile.addressLine }}<span v-if="issuerProfile.addressLine && issuerProfile.cityLine"> · </span>{{ issuerProfile.cityLine }}</p>
        <p v-if="issuerProfile.contactLine" class="print-issuer-line">{{ issuerProfile.contactLine }}</p>
        <h2>Albarán {{ deliveryNote.deliveryNoteNumber }}</h2>
        <div class="print-doc-meta">
          <span>Cliente: {{ partyName }}</span>
          <span>Fecha entrega: {{ formatDate(deliveryNote.deliveryDate) }}</span>
          <span>Pedido: {{ orderNumber || formatOrderId(deliveryNote.salesOrderId) }}</span>
        </div>
      </div>

      <!-- Delivery Note Info Cards -->
      <div class="info-grid">
        <div class="info-card">
          <h3>Información del Pedido</h3>
          <div class="info-row">
            <span class="label">Número de Pedido:</span>
            <a 
              href="#" 
              @click.prevent="navigateToOrder(deliveryNote.salesOrderId)" 
              class="value-link"
            >
              {{ orderNumber || 'Cargando...' }}
            </a>
          </div>
          <div class="info-row">
            <span class="label">Order ID:</span>
            <span class="value order-id">{{ formatOrderId(deliveryNote.salesOrderId) }}</span>
          </div>
        </div>

        <div class="info-card">
          <h3>Información del Cliente</h3>
          <div class="info-row">
            <span class="label">Cliente:</span>
            <span class="value">{{ partyName }}</span>
          </div>
          <div class="info-row">
            <span class="label">Party ID:</span>
            <span class="value party-id">{{ formatPartyId(deliveryNote.partyId) }}</span>
          </div>
        </div>

        <div class="info-card">
          <h3>Fecha de Entrega</h3>
          <div class="info-row">
            <span class="label">Fecha:</span>
            <span class="value">{{ formatDate(deliveryNote.deliveryDate) }}</span>
          </div>
          <div class="info-row">
            <span class="label">Creado:</span>
            <span class="value">{{ formatDateTime(deliveryNote.createdAt) }}</span>
          </div>
        </div>
      </div>

      <!-- Delivery Address -->
      <div v-if="deliveryNote.deliveryAddress" class="info-card address-card">
        <h3>Dirección de Entrega</h3>
        <div class="address-content">
          <p v-if="deliveryNote.deliveryAddress.street">{{ deliveryNote.deliveryAddress.street }}</p>
          <p v-if="deliveryNote.deliveryAddress.city || deliveryNote.deliveryAddress.postalCode">
            {{ deliveryNote.deliveryAddress.postalCode }} {{ deliveryNote.deliveryAddress.city }}
          </p>
          <p v-if="deliveryNote.deliveryAddress.state">{{ deliveryNote.deliveryAddress.state }}</p>
          <p v-if="deliveryNote.deliveryAddress.country">{{ deliveryNote.deliveryAddress.country }}</p>
        </div>
      </div>

      <!-- Delivery Notes / Comments -->
      <div v-if="deliveryNote.notes" class="notes-card">
        <h3>Observaciones</h3>
        <p>{{ deliveryNote.notes }}</p>
      </div>

      <!-- Line Items -->
      <div class="line-items-section">
        <div class="section-header">
          <h2>Líneas Entregadas</h2>
        </div>

        <div v-if="!deliveryNote.lineItems || deliveryNote.lineItems.length === 0" class="empty-state">
          <p>No hay líneas en este albarán</p>
        </div>

        <div v-else class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>Variante de Producto</th>
                <th>Cantidad Entregada</th>
                <th>Observaciones</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in deliveryNote.lineItems" :key="item.id">
                <td class="variant-id">{{ formatVariantId(item.productVariantID) }}</td>
                <td class="quantity">{{ item.quantityDelivered || item.quantity }}</td>
                <td class="observations">{{ item.observations || '—' }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Signatures Section -->
      <div v-if="deliveryNote.signatures || showSignaturesSection" class="signatures-section">
        <h3>Firmas</h3>
        <div class="signatures-grid">
          <div class="signature-box">
            <div class="signature-label">Firma del Cliente</div>
            <div class="signature-area">
              <span v-if="deliveryNote.signatures?.customer" class="signature-present">✓ Firmado</span>
              <span v-else class="signature-pending">Pendiente</span>
            </div>
            <div v-if="deliveryNote.signatures?.customer" class="signature-info">
              <small>{{ deliveryNote.signatures.customer.name }}</small><br>
              <small>{{ formatDateTime(deliveryNote.signatures.customer.timestamp) }}</small>
            </div>
          </div>
          
          <div class="signature-box">
            <div class="signature-label">Firma del Repartidor</div>
            <div class="signature-area">
              <span v-if="deliveryNote.signatures?.driver" class="signature-present">✓ Firmado</span>
              <span v-else class="signature-pending">Pendiente</span>
            </div>
            <div v-if="deliveryNote.signatures?.driver" class="signature-info">
              <small>{{ deliveryNote.signatures.driver.name }}</small><br>
              <small>{{ formatDateTime(deliveryNote.signatures.driver.timestamp) }}</small>
            </div>
          </div>
        </div>
      </div>

      <!-- Related Documents -->
      <div class="related-documents-section">
        <h3>Documentos Relacionados</h3>
        <div class="documents-list">
          <div class="document-item">
            <span class="document-icon">📄</span>
            <div class="document-info">
              <span class="document-label">Pedido origen</span>
              <a 
                href="#" 
                @click.prevent="navigateToOrder(deliveryNote.salesOrderId)" 
                class="document-link"
              >
                {{ orderNumber || formatOrderId(deliveryNote.salesOrderId) }}
              </a>
            </div>
          </div>
          
          <div v-if="relatedInvoice" class="document-item">
            <span class="document-icon">💰</span>
            <div class="document-info">
              <span class="document-label">Factura asociada</span>
              <a 
                href="#" 
                @click.prevent="navigateToInvoice(relatedInvoice.id)" 
                class="document-link"
              >
                {{ relatedInvoice.invoiceNumber }}
              </a>
            </div>
          </div>

          <div v-else class="document-item disabled">
            <span class="document-icon">💰</span>
            <div class="document-info">
              <span class="document-label">Factura asociada</span>
              <span class="document-link-disabled">No generada aún</span>
            </div>
          </div>
        </div>
      </div>

      <div class="print-doc-footer">
        <span>Documento generado por {{ issuerProfile.displayName }}</span>
        <span>Creado: {{ formatDateTime(deliveryNote.createdAt) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import Navbar from '@/components/layout/Navbar.vue';
import salesApi from '@/services/salesApi';
import partyApi from '@/services/partyApi';
import { getPrintIssuerProfile } from '@/services/printIssuerProfile';

const route = useRoute();
const router = useRouter();

const deliveryNote = ref(null);
const isLoading = ref(false);
const error = ref('');
const partyName = ref('Cargando...');
const orderNumber = ref(null);
const relatedInvoice = ref(null);
const issuerProfile = getPrintIssuerProfile();

const showSignaturesSection = computed(() => {
  // Show signatures section even if empty for visual consistency
  return true;
});

onMounted(() => {
  fetchDeliveryNote();
});

async function fetchDeliveryNote() {
  const noteId = route.params.id;
  if (!noteId) {
    error.value = 'ID de albarán no válido';
    return;
  }

  isLoading.value = true;
  error.value = '';

  try {
    deliveryNote.value = await salesApi.getDeliveryNote(noteId);
    
    // Load related data
    await Promise.all([
      loadPartyName(),
      loadOrderNumber(),
      loadRelatedInvoice(),
    ]);
  } catch (err) {
    error.value = err?.message || 'No se pudo cargar el albarán';
    console.error('Error loading delivery note:', err);
  } finally {
    isLoading.value = false;
  }
}

async function loadPartyName() {
  if (!deliveryNote.value?.partyId) {
    partyName.value = 'Desconocido';
    return;
  }
  
  try {
    const party = await partyApi.getParty(deliveryNote.value.partyId);
    partyName.value = party.name || 'Sin nombre';
  } catch (err) {
    console.error('Error loading party name:', err);
    partyName.value = 'Error al cargar';
  }
}

async function loadOrderNumber() {
  if (!deliveryNote.value?.salesOrderId) {
    return;
  }
  
  try {
    const order = await salesApi.getOrder(deliveryNote.value.salesOrderId);
    orderNumber.value = order.orderNumber || formatOrderId(order.id);
  } catch (err) {
    console.error('Error loading order number:', err);
  }
}

async function loadRelatedInvoice() {
  if (!deliveryNote.value?.salesOrderId) {
    return;
  }
  
  try {
    // Try to find invoice related to this order
    const invoices = await salesApi.listInvoices({ 
      orderId: deliveryNote.value.salesOrderId 
    });
    
    if (invoices && invoices.length > 0) {
      relatedInvoice.value = invoices[0];
    }
  } catch (err) {
    console.error('Error loading related invoice:', err);
    // Non-critical, don't show error to user
  }
}

function printDeliveryNote() {
  window.print();
}

function navigateToOrder(orderId) {
  if (!orderId) return;
  router.push(`/sales/orders/${orderId}`);
}

function navigateToInvoice(invoiceId) {
  if (!invoiceId) return;
  router.push(`/sales/invoices/${invoiceId}`);
}

function goBack() {
  router.push('/sales/delivery-notes');
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

function formatDateTime(dateString) {
  if (!dateString) return '—';
  const date = new Date(dateString);
  return date.toLocaleString('es-ES', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  });
}

function formatOrderId(orderId) {
  if (!orderId) return '—';
  return orderId.substring(0, 8) + '...';
}

function formatPartyId(partyId) {
  if (!partyId) return '—';
  return partyId.substring(0, 8) + '...';
}

function formatVariantId(variantId) {
  if (!variantId) return '—';
  return variantId.substring(0, 8) + '...';
}
</script>

<style scoped>
.delivery-note-detail-container {
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
  to { transform: rotate(360deg); }
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
  margin: 0.5rem 0 0.25rem;
}

.subtitle {
  font-size: 0.875rem;
  color: #6b7280;
}

.order-link {
  color: #3b82f6;
  text-decoration: none;
  font-weight: 500;
}

.order-link:hover {
  text-decoration: underline;
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

.header-actions {
  display: flex;
  gap: 0.5rem;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.info-card,
.notes-card {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.address-card {
  grid-column: 1 / -1;
}

.info-card h3,
.notes-card h3 {
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

.info-row .label {
  color: #6b7280;
}

.info-row .value {
  color: #1f2937;
  font-weight: 500;
}

.value-link {
  color: #3b82f6;
  text-decoration: none;
  font-weight: 500;
}

.value-link:hover {
  text-decoration: underline;
}

.order-id,
.party-id {
  font-family: 'Courier New', monospace;
  font-size: 0.8rem;
}

.address-content {
  font-size: 0.875rem;
  color: #4b5563;
  line-height: 1.6;
}

.address-content p {
  margin: 0.25rem 0;
}

.notes-card {
  margin-bottom: 2rem;
}

.notes-card p {
  font-size: 0.875rem;
  color: #4b5563;
  line-height: 1.6;
}

.line-items-section {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  margin-bottom: 2rem;
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

.quantity {
  font-weight: 600;
  color: #059669;
}

.observations {
  color: #6b7280;
  font-style: italic;
}

.empty-state {
  text-align: center;
  padding: 2rem;
  color: #9ca3af;
}

.signatures-section {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  margin-bottom: 2rem;
}

.signatures-section h3 {
  font-size: 1rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 1.5rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid #f3f4f6;
}

.signatures-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
}

.signature-box {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 1rem;
}

.signature-label {
  font-size: 0.875rem;
  font-weight: 600;
  color: #4b5563;
  margin-bottom: 0.75rem;
}

.signature-area {
  min-height: 80px;
  background: #f9fafb;
  border: 2px dashed #d1d5db;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 0.5rem;
}

.signature-present {
  color: #059669;
  font-weight: 600;
  font-size: 1.25rem;
}

.signature-pending {
  color: #9ca3af;
  font-style: italic;
}

.signature-info {
  font-size: 0.75rem;
  color: #6b7280;
  text-align: center;
}

.related-documents-section {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.related-documents-section h3 {
  font-size: 1rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 1rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid #f3f4f6;
}

.documents-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.document-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.75rem;
  background: #f9fafb;
  border-radius: 6px;
}

.document-item.disabled {
  opacity: 0.6;
}

.document-icon {
  font-size: 1.5rem;
}

.document-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.document-label {
  font-size: 0.75rem;
  color: #6b7280;
  text-transform: uppercase;
  font-weight: 600;
  letter-spacing: 0.025em;
}

.document-link {
  color: #3b82f6;
  text-decoration: none;
  font-weight: 500;
  font-size: 0.875rem;
}

.document-link:hover {
  text-decoration: underline;
}

.document-link-disabled {
  color: #9ca3af;
  font-style: italic;
  font-size: 0.875rem;
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
  .document-link,
  .order-link,
  .value-link {
    display: none !important;
  }

  .delivery-note-detail-container {
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
  .line-items-section,
  .signatures-section,
  .related-documents-section,
  .table-container,
  .signature-box,
  .document-item {
    box-shadow: none !important;
    border: 1px solid #d1d5db;
  }

  .data-table {
    font-size: 0.75rem;
  }

  .signatures-section,
  .related-documents-section {
    break-inside: avoid;
  }
}
</style>
