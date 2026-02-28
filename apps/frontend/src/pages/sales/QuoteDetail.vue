<template>
  <Navbar />
  <div class="quote-detail-container">
    <!-- Loading State -->
    <div v-if="isLoading" class="loading-state">
      <div class="spinner"></div>
      <p>Cargando presupuesto...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-state">
      <p class="error-message">{{ error }}</p>
      <button class="btn btn-secondary" @click="fetchQuote">Reintentar</button>
      <button class="btn btn-secondary" @click="goBack">Volver</button>
    </div>

    <!-- Quote Detail -->
    <div v-else-if="quote" class="quote-content">
      <!-- Header -->
      <div class="page-header">
        <div>
          <button class="btn-back" @click="goBack">← Volver</button>
          <h1>Presupuesto {{ quote.quoteNumber }}</h1>
          <span :class="['status-badge', `status-${getStatusClass(quote.status)}`]">
            {{ getStatusLabel(quote.status) }}
          </span>
        </div>
        <div class="header-actions">
          <button
            class="btn btn-secondary"
            @click="printQuote"
            title="Imprimir presupuesto"
          >
            🖨️ Imprimir
          </button>
          <!-- DRAFT: Enviar, Editar (implícito), Eliminar -->
          <button 
            v-if="quote.status === 'DRAFT'" 
            class="btn btn-primary" 
            @click="sendQuote"
            title="Enviar presupuesto al cliente"
          >
            📧 Enviar
          </button>
          <button 
            v-if="quote.status === 'DRAFT'" 
            class="btn btn-danger" 
            @click="deleteQuote"
            title="Eliminar presupuesto"
          >
            🗑️ Eliminar
          </button>
          
          <!-- SENT: Aceptar, Rechazar, Convertir a Pedido -->
          <button 
            v-if="quote.status === 'SENT'" 
            class="btn btn-success" 
            @click="acceptQuote"
            title="Marcar como aceptado"
          >
            ✓ Aceptar
          </button>
          <button 
            v-if="quote.status === 'SENT'" 
            class="btn btn-danger" 
            @click="rejectQuote"
            title="Marcar como rechazado"
          >
            ✕ Rechazar
          </button>
          <button 
            v-if="quote.status === 'SENT' && !isExpired" 
            class="btn btn-primary" 
            @click="showConvertModal = true"
            title="Convertir a pedido"
          >
            🔄 Convertir a Pedido
          </button>

          <!-- ACCEPTED: Convertir a Pedido -->
          <button 
            v-if="quote.status === 'ACCEPTED'" 
            class="btn btn-primary" 
            @click="showConvertModal = true"
            title="Convertir a pedido"
          >
            🔄 Convertir a Pedido
          </button>
        </div>
      </div>

      <div class="print-doc-header">
        <p class="print-brand">{{ issuerProfile.displayName }}</p>
        <p v-if="issuerProfile.taxId" class="print-issuer-line">{{ issuerProfile.taxLabel }}: {{ issuerProfile.taxId }}</p>
        <p v-if="issuerProfile.addressLine || issuerProfile.cityLine" class="print-issuer-line">{{ issuerProfile.addressLine }}<span v-if="issuerProfile.addressLine && issuerProfile.cityLine"> · </span>{{ issuerProfile.cityLine }}</p>
        <p v-if="issuerProfile.contactLine" class="print-issuer-line">{{ issuerProfile.contactLine }}</p>
        <h2>Presupuesto {{ quote.quoteNumber }}</h2>
        <div class="print-doc-meta">
          <span>Cliente: {{ partyName }}</span>
          <span>Fecha: {{ formatDate(quote.quoteDate) }}</span>
          <span>Estado: {{ getStatusLabel(quote.status) }}</span>
        </div>
      </div>

      <!-- Expiration Warning (SENT and approaching expiration) -->
      <div v-if="quote.status === 'SENT' && daysUntilExpiration <= 7 && daysUntilExpiration > 0" class="warning-card">
        ⚠️ Este presupuesto vence en <strong>{{ daysUntilExpiration }}</strong> día{{ daysUntilExpiration !== 1 ? 's' : '' }}.
      </div>
      <div v-if="isExpired && quote.status === 'SENT'" class="error-card">
        ❌ Este presupuesto ha <strong>EXPIRADO</strong>. No se puede convertir a pedido.
      </div>

      <!-- Quote Info Cards -->
      <div class="info-grid">
        <div class="info-card">
          <h3>Información del Cliente</h3>
          <div class="info-row">
            <span class="label">Cliente:</span>
            <span class="value">{{ partyName }}</span>
          </div>
          <div class="info-row">
            <span class="label">Party ID:</span>
            <span class="value party-id">{{ formatPartyId(quote.partyId) }}</span>
          </div>
        </div>

        <div class="info-card">
          <h3>Fechas</h3>
          <div class="info-row">
            <span class="label">Fecha de Presupuesto:</span>
            <span class="value">{{ formatDate(quote.quoteDate) }}</span>
          </div>
          <div class="info-row">
            <span class="label">Válido Hasta:</span>
            <span class="value" :class="{'text-danger': isExpired, 'text-warning': daysUntilExpiration <= 7 && daysUntilExpiration > 0}">
              {{ formatDate(quote.validUntil) }}
            </span>
          </div>
        </div>

        <div class="info-card">
          <h3>Totales</h3>
          <div class="info-row">
            <span class="label">Subtotal:</span>
            <span class="value">{{ salesApi.formatMoney(quote.subtotal) }}</span>
          </div>
          <div class="info-row">
            <span class="label">IVA:</span>
            <span class="value">{{ salesApi.formatMoney(quote.taxAmount) }}</span>
          </div>
          <div class="info-row total">
            <span class="label">Total:</span>
            <span class="value">{{ salesApi.formatMoney(quote.total) }}</span>
          </div>
        </div>
      </div>

      <!-- Internal Notes -->
      <div v-if="quote.internalNotes" class="notes-card">
        <h3>Notas Internas</h3>
        <p>{{ quote.internalNotes }}</p>
      </div>

      <!-- Line Items -->
      <div class="line-items-section">
        <div class="section-header">
          <h2>Líneas del Presupuesto</h2>
        </div>

        <div v-if="!quote.lineItems || quote.lineItems.length === 0" class="empty-state">
          <p>No hay líneas en este presupuesto</p>
        </div>

        <div v-else class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>Variante de Producto</th>
                <th>Cantidad</th>
                <th>Precio Unitario</th>
                <th>Descuento</th>
                <th>IVA %</th>
                <th>IVA línea</th>
                <th>Definición MES</th>
                <th>Subtotal</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in quote.lineItems" :key="item.id">
                <td class="variant-id">{{ formatVariantId(item.productVariantID) }}</td>
                <td>{{ item.quantity }}</td>
                <td>
                  {{ salesApi.formatMoney(item.finalUnitPrice) }}
                  <span v-if="item.manualUnitPrice" class="manual-badge">Manual</span>
                </td>
                <td>{{ item.finalDiscountPerUnit ? salesApi.formatMoney(item.finalDiscountPerUnit) : '—' }}</td>
                <td>{{ typeof item.taxRate === 'number' ? `${item.taxRate}%` : '—' }}</td>
                <td>{{ salesApi.formatMoney(item.taxAmount) }}</td>
                <td>
                  <button
                    v-if="item.mesWorkId"
                    class="mes-link"
                    @click="goToMesWork(item.mesWorkId)"
                    type="button"
                    :title="getMesWorkTooltip(item.mesWorkId)"
                  >
                    {{ formatMesWorkId(item.mesWorkId) }}
                  </button>
                  <span v-else>—</span>
                </td>
                <td class="amount">{{ salesApi.formatMoney(item.subtotal) }}</td>
              </tr>
            </tbody>
            <tfoot>
              <tr class="totals-row">
                <td colspan="7" class="totals-label">Subtotal:</td>
                <td class="amount">{{ salesApi.formatMoney(quote.subtotal) }}</td>
              </tr>
              <tr class="totals-row">
                <td colspan="7" class="totals-label">IVA:</td>
                <td class="amount">{{ salesApi.formatMoney(quote.taxAmount) }}</td>
              </tr>
              <tr class="totals-row total">
                <td colspan="7" class="totals-label">Total:</td>
                <td class="amount">{{ salesApi.formatMoney(quote.total) }}</td>
              </tr>
            </tfoot>
          </table>
        </div>
      </div>

      <div class="print-doc-footer">
        <span>Documento generado por {{ issuerProfile.displayName }}</span>
        <span>Válido hasta: {{ formatDate(quote.validUntil) }}</span>
      </div>
    </div>

    <!-- Convert to Order Modal -->
    <div v-if="showConvertModal" class="modal-overlay" @click="showConvertModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Convertir Presupuesto a Pedido</h3>
          <button class="btn-close" @click="showConvertModal = false">✕</button>
        </div>
        <div class="modal-body">
          <p class="modal-description">
            Presupuesto: <strong>{{ quote?.quoteNumber }}</strong><br>
            Total: <strong>{{ quote ? salesApi.formatMoney(quote.total) : '' }}</strong>
          </p>
          <div class="form-group">
            <label for="deliveryDate">Fecha de Entrega *</label>
            <input
              id="deliveryDate"
              v-model="convertForm.deliveryDate"
              type="date"
              class="form-input"
              :min="minDeliveryDate"
              required
            />
            <small class="help-text">La fecha de entrega debe ser al menos mañana</small>
          </div>
          <div class="form-group">
            <label for="orderNotes">Notas del Pedido</label>
            <textarea
              id="orderNotes"
              v-model="convertForm.notes"
              class="form-textarea"
              rows="3"
              placeholder="Notas adicionales para el pedido..."
            ></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showConvertModal = false">Cancelar</button>
          <button 
            class="btn btn-primary" 
            @click="convertToOrder"
            :disabled="!convertForm.deliveryDate || isConverting"
          >
            {{ isConverting ? 'Convirtiendo...' : 'Crear Pedido' }}
          </button>
        </div>
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
import { mesApi } from '@/services/mesApi';
import { getPrintIssuerProfile } from '@/services/printIssuerProfile';

const route = useRoute();
const router = useRouter();

const quote = ref(null);
const isLoading = ref(false);
const error = ref('');
const partyName = ref('Cargando...');
const mesWorksCache = ref({});
const issuerProfile = getPrintIssuerProfile();

const showConvertModal = ref(false);
const isConverting = ref(false);
const convertForm = ref({
  deliveryDate: '',
  notes: '',
});

const minDeliveryDate = computed(() => {
  const tomorrow = new Date();
  tomorrow.setDate(tomorrow.getDate() + 1);
  return tomorrow.toISOString().split('T')[0];
});

const isExpired = computed(() => {
  if (!quote.value?.validUntil) return false;
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  const validUntil = new Date(quote.value.validUntil);
  validUntil.setHours(0, 0, 0, 0);
  return validUntil < today;
});

const daysUntilExpiration = computed(() => {
  if (!quote.value?.validUntil) return null;
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  const validUntil = new Date(quote.value.validUntil);
  validUntil.setHours(0, 0, 0, 0);
  const diffTime = validUntil - today;
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
  return diffDays;
});

onMounted(() => {
  fetchQuote();
});

async function fetchQuote() {
  const quoteId = route.params.id;
  if (!quoteId) {
    error.value = 'ID de presupuesto no válido';
    return;
  }

  isLoading.value = true;
  error.value = '';

  try {
    quote.value = await salesApi.getQuote(quoteId);
    await loadPartyName();
    await loadMesWorksForQuote();
  } catch (err) {
    error.value = err?.message || 'No se pudo cargar el presupuesto';
    console.error('Error loading quote:', err);
  } finally {
    isLoading.value = false;
  }
}

async function loadPartyName() {
  if (!quote.value?.partyId) {
    partyName.value = 'Desconocido';
    return;
  }
  
  try {
    const party = await partyApi.getParty(quote.value.partyId);
    partyName.value = party.name || 'Sin nombre';
  } catch (err) {
    console.error('Error loading party name:', err);
    partyName.value = 'Error al cargar';
  }
}

async function loadMesWorksForQuote() {
  const lineItems = quote.value?.lineItems;
  if (!Array.isArray(lineItems) || lineItems.length === 0) return;

  const mesWorkIds = [...new Set(
    lineItems
      .map((item) => item?.mesWorkId)
      .filter((mesWorkId) => typeof mesWorkId === 'string' && mesWorkId.length > 0),
  )];

  const uncachedIds = mesWorkIds.filter((id) => !mesWorksCache.value[id]);
  if (uncachedIds.length === 0) return;

  const results = await Promise.allSettled(uncachedIds.map((id) => mesApi.getWorkDefinition(id)));
  results.forEach((result, index) => {
    const mesWorkId = uncachedIds[index];
    mesWorksCache.value[mesWorkId] = result.status === 'fulfilled' ? result.value : null;
  });
}

async function sendQuote() {
  if (!confirm('¿Enviar este presupuesto al cliente? El estado cambiará a "Enviado".')) return;

  try {
    await salesApi.changeQuoteStatus(quote.value.id, 'SENT');
    await fetchQuote();
  } catch (err) {
    alert(err?.message || 'No se pudo enviar el presupuesto');
  }
}

async function acceptQuote() {
  if (!confirm('¿Marcar este presupuesto como aceptado?')) return;

  try {
    await salesApi.changeQuoteStatus(quote.value.id, 'ACCEPTED');
    await fetchQuote();
  } catch (err) {
    alert(err?.message || 'No se pudo aceptar el presupuesto');
  }
}

async function rejectQuote() {
  if (!confirm('¿Marcar este presupuesto como rechazado? Esta acción no se puede deshacer.')) return;

  try {
    await salesApi.changeQuoteStatus(quote.value.id, 'REJECTED');
    await fetchQuote();
  } catch (err) {
    alert(err?.message || 'No se pudo rechazar el presupuesto');
  }
}

async function deleteQuote() {
  if (!confirm('¿Eliminar este presupuesto? Esta acción no se puede deshacer.')) return;

  try {
    // Note: salesApi might not have deleteQuote yet, this is a placeholder
    // await salesApi.deleteQuote(quote.value.id);
    alert('Funcionalidad de eliminación pendiente de implementación en el backend');
    // router.push('/sales/quotes');
  } catch (err) {
    alert(err?.message || 'No se pudo eliminar el presupuesto');
  }
}

async function convertToOrder() {
  if (!convertForm.value.deliveryDate) {
    alert('Por favor, especifica una fecha de entrega');
    return;
  }

  if (isExpired.value) {
    alert('No se puede convertir un presupuesto expirado a pedido');
    return;
  }

  isConverting.value = true;

  try {
    const order = await salesApi.convertQuoteToOrder(quote.value.id, convertForm.value.deliveryDate);
    
    alert(`✓ Presupuesto convertido exitosamente.\nPedido creado: ${order.orderNumber || order.id}`);
    
    // Navigate to the new order
    if (order.id) {
      router.push(`/sales/orders/${order.id}`);
    } else {
      router.push('/sales/orders');
    }
  } catch (err) {
    alert(err?.message || 'No se pudo convertir el presupuesto a pedido');
    isConverting.value = false;
  }
}

function goBack() {
  router.push('/sales/quotes');
}

function printQuote() {
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

function formatVariantId(variantId) {
  if (!variantId) return '—';
  return variantId.substring(0, 8) + '...';
}

function formatMesWorkId(mesWorkId) {
  if (!mesWorkId) return '—';
  const mesWork = mesWorksCache.value[mesWorkId];
  if (mesWork?.work_number) return mesWork.work_number;
  return mesWorkId.substring(0, 8) + '...';
}

function getMesWorkTooltip(mesWorkId) {
  if (!mesWorkId) return 'Definición de trabajo MES';
  const mesWork = mesWorksCache.value[mesWorkId];
  if (mesWork?.work_number && mesWork?.work_name) {
    return `${mesWork.work_number} · ${mesWork.work_name}`;
  }
  if (mesWork?.work_number) return mesWork.work_number;
  return `Definición MES ${mesWorkId}`;
}

function goToMesWork(mesWorkId) {
  if (!mesWorkId) return;
  router.push(`/mes/work-definitions/${mesWorkId}`);
}

function getStatusLabel(status) {
  const labels = {
    DRAFT: 'Borrador',
    SENT: 'Enviado',
    ACCEPTED: 'Aceptado',
    REJECTED: 'Rechazado',
    EXPIRED: 'Expirado',
  };
  return labels[status] || status;
}

function getStatusClass(status) {
  const classes = {
    DRAFT: 'warning',
    SENT: 'info',
    ACCEPTED: 'success',
    REJECTED: 'danger',
    EXPIRED: 'secondary',
  };
  return classes[status] || 'secondary';
}
</script>

<style scoped>
.quote-detail-container {
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

.header-actions {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.warning-card,
.error-card {
  padding: 1rem 1.5rem;
  border-radius: 8px;
  margin-bottom: 1.5rem;
  font-size: 0.875rem;
}

.warning-card {
  background: #fef3c7;
  border-left: 4px solid #f59e0b;
  color: #92400e;
}

.error-card {
  background: #fee2e2;
  border-left: 4px solid #ef4444;
  color: #991b1b;
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

.party-id {
  font-family: 'Courier New', monospace;
  font-size: 0.8rem;
}

.text-danger {
  color: #dc2626 !important;
  font-weight: 600;
}

.text-warning {
  color: #f59e0b !important;
  font-weight: 600;
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

.btn-success {
  background: #10b981;
  color: white;
}

.btn-success:hover {
  background: #059669;
}

.btn-danger {
  background: #ef4444;
  color: white;
}

.btn-danger:hover {
  background: #dc2626;
}

.status-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.025em;
}

.status-warning {
  background: #fef3c7;
  color: #92400e;
}

.status-info {
  background: #dbeafe;
  color: #1e40af;
}

.status-success {
  background: #d1fae5;
  color: #065f46;
}

.status-danger {
  background: #fee2e2;
  color: #991b1b;
}

.status-secondary {
  background: #f3f4f6;
  color: #4b5563;
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

.data-table tfoot {
  background: #f9fafb;
  font-weight: 600;
}

.totals-row td {
  padding: 0.75rem 1rem;
  border-top: 1px solid #e5e7eb;
}

.totals-row.total td {
  border-top: 2px solid #d1d5db;
  font-size: 1rem;
  padding-top: 1rem;
}

.totals-label {
  text-align: right;
  color: #6b7280;
}

.variant-id {
  font-family: 'Courier New', monospace;
  color: #6b7280;
}

.amount {
  font-weight: 600;
  text-align: right;
}

.manual-badge {
  display: inline-block;
  margin-left: 0.5rem;
  padding: 0.125rem 0.375rem;
  background: #e0e7ff;
  color: #3730a3;
  border-radius: 4px;
  font-size: 0.65rem;
  font-weight: 600;
  text-transform: uppercase;
}

.mes-link {
  background: transparent;
  border: none;
  color: #2563eb;
  cursor: pointer;
  padding: 0;
  font-size: 0.875rem;
  text-decoration: underline;
}

.mes-link:hover {
  color: #1d4ed8;
}

.empty-state {
  text-align: center;
  padding: 2rem;
  color: #9ca3af;
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

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: #4a5568;
  margin-bottom: 0.25rem;
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

.form-textarea {
  resize: vertical;
}

.help-text {
  display: block;
  margin-top: 0.25rem;
  font-size: 0.75rem;
  color: #6b7280;
}

.modal-footer {
  display: flex;
  gap: 0.5rem;
  justify-content: flex-end;
  padding: 1.5rem;
  border-top: 1px solid #f3f4f6;
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
  .modal-overlay,
  .btn,
  .btn-icon {
    display: none !important;
  }

  .quote-detail-container {
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
  .table-container {
    box-shadow: none !important;
    border: 1px solid #d1d5db;
  }

  .info-grid {
    grid-template-columns: 1fr 1fr;
    gap: 0.75rem;
  }

  .status-badge {
    border: 1px solid #9ca3af;
    color: #111827 !important;
    background: transparent !important;
  }

  .data-table {
    font-size: 0.75rem;
  }
}
</style>
