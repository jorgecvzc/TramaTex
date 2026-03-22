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
    <div v-else-if="invoice" ref="invoiceContentRef" class="invoice-content">
      <!-- Header -->
      <div class="page-header">
        <div>
          <button class="btn-back" @click="goBack">← Volver</button>
          <h1>Factura {{ invoice.invoiceNumber }}
            <span :class="['status-badge', 'status-' + salesApi.getStatusClass(invoice.status)]">{{ salesApi.getStatusLabel(invoice.status) }}</span>
          </h1>
          <span :class="['type-badge', `type-${invoice.type.toLowerCase()}`]">
            {{ getTypeLabel(invoice.type) }}
          </span>
        </div>
        <div class="header-actions">
          <button
            v-if="invoice.status === 'BORRADOR'"
            class="btn btn-primary"
            @click="emitInvoice"
            :disabled="isChangingStatus"
          >
            📤 Emitir
          </button>
          <button
            v-if="invoice.status === 'EMITIDA' || invoice.status === 'VENCIDA'"
            class="btn btn-success"
            @click="markAsPaid"
            :disabled="isChangingStatus"
          >
            💰 Marcar como Pagada
          </button>
          <button
            v-if="invoice.status !== 'ANULADA' && invoice.status !== 'PAGADA'"
            class="btn btn-danger"
            @click="voidInvoice"
            :disabled="isChangingStatus"
          >
            🚫 Anular
          </button>
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
          <span>Cliente: {{ partyName }}</span>
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
            <span class="value">{{ partyName }}</span>
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

      <!-- Documentos Relacionados -->
      <div v-if="(invoice.salesOrderIds && invoice.salesOrderIds.length > 0) || (invoice.deliveryNoteIds && invoice.deliveryNoteIds.length > 0)" class="related-documents-section">
        <h3>Documentos Relacionados</h3>
        <div class="documents-list">
          <div v-for="order in relatedOrders" :key="order.id" class="document-item">
            <span class="document-icon">📄</span>
            <div class="document-info">
              <span class="document-label">Pedido</span>
              <a href="#" @click.prevent="$router.push(`/sales/orders/${order.id}`)" class="document-link">
                {{ order.orderNumber || formatId(order.id) }}
              </a>
            </div>
          </div>

          <div v-for="dn in relatedDeliveryNotes" :key="dn.id" class="document-item">
            <span class="document-icon">📦</span>
            <div class="document-info">
              <span class="document-label">Albarán</span>
              <a href="#" @click.prevent="$router.push(`/sales/delivery-notes/${dn.id}`)" class="document-link">
                {{ dn.deliveryNoteNumber || formatId(dn.id) }}
              </a>
            </div>
          </div>
        </div>
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
                <th>Referencia</th>
                <th>Nombre</th>
                <th>Cantidad</th>
                <th>Precio Unitario</th>
                <th>Dto. %</th>
                <th>IVA %</th>
                <th>IVA línea</th>
                <th>Subtotal</th>
                <th>Total línea</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in invoice.lineItems" :key="item.id || item.productVariantID">
                <td class="variant-id">
                  <span v-if="item.variantSku">{{ item.variantSku }}</span>
                  <span v-else>{{ formatId(item.productVariantID) }}</span>
                </td>
                <td>{{ buildDisplayName(item) }}</td>
                <td>{{ item.quantity }}</td>
                <td>{{ salesApi.formatUnitPrice(item.unitPrice) }}</td>
                <td>{{ item.discountAmount && item.unitPrice?.amount ? ((item.discountAmount.amount / item.unitPrice.amount) * 100).toFixed(2) + '%' : '—' }}</td>
                <td>{{ typeof item.taxRate === 'number' ? `${item.taxRate}%` : '—' }}</td>
                <td>{{ item.taxAmount ? salesApi.formatMoney(item.taxAmount) : '—' }}</td>
                <td class="amount">{{ salesApi.formatMoney(item.subtotal) }}</td>
                <td class="amount">{{ salesApi.formatMoney(item.total) }}</td>
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

    <!-- Post-Issue Actions Modal -->
    <div v-if="showPostIssueModal" class="modal-overlay" @click="showPostIssueModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Factura Emitida</h3>
          <button class="btn-close" @click="showPostIssueModal = false">✕</button>
        </div>
        <div class="modal-body">
          <p class="modal-description">
            La factura <strong>{{ invoice?.invoiceNumber }}</strong> se ha emitido correctamente.
          </p>
          <p class="post-issue-prompt">¿Qué desea hacer a continuación?</p>
          <div class="post-issue-actions">
            <button class="btn btn-primary post-issue-btn" @click="postIssuePrint">
              🖨️ Imprimir
            </button>
            <button class="btn btn-primary post-issue-btn" @click="postIssueEmail">
              📧 Descargar PDF y enviar por correo
            </button>
            <button class="btn btn-primary post-issue-btn" @click="postIssueBoth">
              🖨️📧 Descargar PDF, imprimir y enviar
            </button>
          </div>
          <p class="post-issue-hint">El PDF se descargará automáticamente para adjuntarlo al correo.</p>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showPostIssueModal = false">Cerrar sin acción</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import Navbar from '@/components/layout/Navbar.vue';
import salesApi from '@/services/salesApi';
import { partyApi } from '@/services/partyApi';
import { getPrintIssuerProfile } from '@/services/printIssuerProfile';
import '@/assets/sales-print.css';

const route = useRoute();
const router = useRouter();

const invoice = ref(null);
const isLoading = ref(false);
const isChangingStatus = ref(false);
const error = ref('');
const issuerProfile = getPrintIssuerProfile();
const invoiceContentRef = ref(null);
const showPostIssueModal = ref(false);
const partyName = ref('');
const relatedOrders = ref([]);
const relatedDeliveryNotes = ref([]);

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
    await Promise.all([loadPartyName(), loadRelatedOrders(), loadRelatedDeliveryNotes()]);
  } catch (err) {
    error.value = err?.message || 'No se pudo cargar la factura';
    console.error('Error loading invoice:', err);
  } finally {
    isLoading.value = false;
  }
}

async function emitInvoice() {
  if (!confirm('¿Emitir esta factura? Una vez emitida no se puede volver a borrador.')) return;
  isChangingStatus.value = true;
  try {
    invoice.value = await salesApi.changeInvoiceStatus(invoice.value.id, 'EMITIDA');
    showPostIssueModal.value = true;
  } catch (err) {
    alert(err?.message || 'No se pudo emitir la factura');
  } finally {
    isChangingStatus.value = false;
  }
}

async function loadPartyName() {
  if (!invoice.value?.partyId) { partyName.value = 'Desconocido'; return; }
  try {
    const party = await partyApi.getParty(invoice.value.partyId);
    partyName.value = party.name || 'Sin nombre';
  } catch {
    partyName.value = 'Desconocido';
  }
}

async function loadRelatedOrders() {
  const ids = invoice.value?.salesOrderIds;
  if (!ids || ids.length === 0) { relatedOrders.value = []; return; }
  const results = await Promise.all(
    ids.map(id => salesApi.getOrder(id).catch(() => ({ id, orderNumber: null })))
  );
  relatedOrders.value = results;
}

async function loadRelatedDeliveryNotes() {
  const ids = invoice.value?.deliveryNoteIds;
  if (!ids || ids.length === 0) { relatedDeliveryNotes.value = []; return; }
  const results = await Promise.all(
    ids.map(id => salesApi.getDeliveryNote(id).catch(() => ({ id, deliveryNoteNumber: null })))
  );
  relatedDeliveryNotes.value = results;
}

function openMailClient() {
  let email = '';
  try {
    partyApi.listContacts(invoice.value.partyId).then(({ data: contacts }) => {
      const primary = contacts.find(c => c.email);
      if (primary) email = primary.email;
    }).catch(() => {});
  } catch { /* no bloquear */ }

  const inv = invoice.value;
  const subject = encodeURIComponent(`Factura ${inv.invoiceNumber || inv.id}`);
  const total = inv.total?.amount != null ? Number(inv.total.amount).toFixed(2) + ' €' : '(pendiente)';
  const body = encodeURIComponent(
    `Estimado/a ${partyName.value},\n\n` +
    `Adjunto le enviamos la factura ${inv.invoiceNumber || ''} ` +
    `por un total de ${total}.\n\n` +
    `Quedamos a su disposición para cualquier consulta.\n\n` +
    `Un saludo.`
  );
  window.open(`mailto:${email}?subject=${subject}&body=${body}`, '_self');
}

function postIssuePrint() {
  showPostIssueModal.value = false;
  window.print();
}

async function postIssueEmail() {
  showPostIssueModal.value = false;
  await generateInvoicePdf();
  openMailClient();
}

async function postIssueBoth() {
  showPostIssueModal.value = false;
  await generateInvoicePdf();
  window.print();
  setTimeout(() => openMailClient(), 500);
}

async function generateInvoicePdf() {
  const el = invoiceContentRef.value;
  if (!el) return;
  const { default: html2pdf } = await import('html2pdf.js');
  const filename = `Factura_${invoice.value.invoiceNumber || invoice.value.id}.pdf`;

  el.classList.add('pdf-rendering');
  await new Promise(r => setTimeout(r, 100));

  try {
    await html2pdf()
      .set({
        margin: [10, 10, 10, 10],
        filename,
        image: { type: 'jpeg', quality: 0.95 },
        html2canvas: { scale: 2, useCORS: true },
        jsPDF: { unit: 'mm', format: 'a4', orientation: 'portrait' },
        pagebreak: { mode: ['avoid-all', 'css', 'legacy'] },
      })
      .from(el)
      .save();
  } finally {
    el.classList.remove('pdf-rendering');
  }
}

async function markAsPaid() {
  if (!confirm('¿Marcar esta factura como pagada?')) return;
  isChangingStatus.value = true;
  try {
    invoice.value = await salesApi.changeInvoiceStatus(invoice.value.id, 'PAGADA');
  } catch (err) {
    alert(err?.message || 'No se pudo marcar como pagada');
  } finally {
    isChangingStatus.value = false;
  }
}

async function voidInvoice() {
  if (!confirm('¿Anular esta factura? Esta acción no se puede deshacer.')) return;
  isChangingStatus.value = true;
  try {
    invoice.value = await salesApi.changeInvoiceStatus(invoice.value.id, 'ANULADA');
  } catch (err) {
    alert(err?.message || 'No se pudo anular la factura');
  } finally {
    isChangingStatus.value = false;
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

function buildDisplayName(item) {
  const name = item.productName || '';
  const config = item.optionConfiguration;
  if (!name) return '—';
  if (!config || Object.keys(config).length === 0) return name;
  return name + ' - ' + Object.values(config).join(', ');
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
.notes-card {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
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

.related-documents-section {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  margin-bottom: 2rem;
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
  color: #374151;
}

.variant-sku {
  color: #6b7280;
  font-family: 'Courier New', monospace;
  font-size: 0.85em;
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

/* Print styles: see @/assets/sales-print.css */

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
  background: white;
  border-radius: 8px;
  max-width: 500px;
  width: 90%;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #f3f4f6;
}

.modal-header h3 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #1f2937;
}

.btn-close {
  background: transparent;
  border: none;
  font-size: 1.5rem;
  color: #9ca3af;
  cursor: pointer;
  padding: 0;
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.2s;
}

.btn-close:hover {
  background: #f3f4f6;
  color: #1f2937;
}

.modal-body {
  padding: 1.5rem;
}

.modal-description {
  font-size: 0.875rem;
  color: #4b5563;
  line-height: 1.6;
  margin-bottom: 1.5rem;
  padding: 1rem;
  background: #f9fafb;
  border-radius: 4px;
}

.post-issue-prompt {
  font-size: 0.95rem;
  color: #374151;
  margin-bottom: 1rem;
}

.post-issue-actions {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.post-issue-btn {
  width: 100%;
  justify-content: center;
  font-size: 0.95rem;
  padding: 0.75rem 1rem;
}

.post-issue-hint {
  margin-top: 0.5rem;
  font-size: 0.8rem;
  color: var(--tt-text-secondary, #64748b);
  text-align: center;
}

.modal-footer {
  display: flex;
  gap: 0.5rem;
  justify-content: flex-end;
  padding: 1.5rem;
  border-top: 1px solid #f3f4f6;
}
</style>
