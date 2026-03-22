<template>
  <Navbar />
  <div class="invoice-list-container">
    <!-- Page Header -->
    <div class="page-header">
      <h1>Facturas</h1>
      <div class="header-actions">
        <button class="btn btn-secondary" @click="navigateToCreateTicket">
          + Nuevo Ticket
        </button>
        <button class="btn btn-primary" @click="showCreateInvoiceModal = true">
          + Nueva Factura
        </button>
      </div>
    </div>

    <!-- Filters -->
    <div class="filters-card">
      <div class="filters-grid">
        <div class="filter-group">
          <PartySelector
            v-model="filters.partyId"
            label="Cliente"
            placeholder="Buscar cliente por nombre o referencia..."
            role-filter="CLIENT"
            :required="false"
          />
        </div>
        <div class="filter-group">
          <label>Búsqueda</label>
          <input
            v-model="filters.searchText"
            type="text"
            class="filter-input"
            placeholder="Buscar por referencia o nombre"
          />
        </div>
        <div class="filter-group">
          <label>Tipo</label>
          <select v-model="filters.type" class="filter-select">
            <option value="">Todos</option>
            <option value="COMPLETA">Estándar</option>
            <option value="SIMPLIFICADA">Simplificada</option>
          </select>
        </div>
        <div class="filter-group">
          <label>Desde</label>
          <input v-model="filters.fromDate" type="date" class="filter-input" />
        </div>
        <div class="filter-group">
          <label>Hasta</label>
          <input v-model="filters.toDate" type="date" class="filter-input" />
        </div>
      </div>
      <div class="filter-actions">
        <div class="limit-group">
          <label>Mostrar</label>
          <select v-model.number="filters.limit" class="filter-select limit-select" @change="applyFilters">
            <option :value="25">25</option>
            <option :value="50">50</option>
            <option :value="100">100</option>
            <option :value="0">Todos</option>
          </select>
          <span class="limit-label">registros</span>
        </div>
        <button class="btn btn-secondary" @click="clearFilters">Limpiar</button>
        <button class="btn btn-primary" @click="applyFilters" :disabled="!isDateRangeValid" :title="!isDateRangeValid ? 'Completa ambas fechas o vacía ambas para buscar' : ''">Buscar</button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="loading-state">
      <div class="spinner"></div>
      <p>Cargando facturas...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-state">
      <p class="error-message">{{ error }}</p>
      <button class="btn btn-secondary" @click="fetchInvoices">Reintentar</button>
    </div>

    <!-- Empty State -->
    <div v-else-if="!filteredInvoices || filteredInvoices.length === 0" class="empty-state">
      <p>No se encontraron facturas con los criterios seleccionados</p>
      <button class="btn btn-primary" @click="showCreateInvoiceModal = true">
        Crear Primera Factura
      </button>
    </div>

    <!-- Data Table -->
    <div v-else class="table-card">
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>Número</th>
              <th>Cliente</th>
              <th>Fecha</th>
              <th>Vencimiento</th>
              <th>Tipo</th>
              <th>Estado</th>
              <th>Total</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="invoice in filteredInvoices"
              :key="invoice.id"
              class="clickable-row"
              @click="navigateToDetail(invoice.id)"
            >
              <td class="invoice-number">{{ invoice.invoiceNumber }}</td>
              <td class="party-id">{{ formatPartyId(invoice.partyId) }}</td>
              <td>{{ formatDate(invoice.issueDate) }}</td>
              <td>{{ formatDate(invoice.dueDate) }}</td>
              <td>
                <span :class="['type-badge', `type-${invoice.type.toLowerCase()}`]">
                  {{ getTypeLabel(invoice.type) }}
                </span>
              </td>
              <td>
                <span :class="['status-badge', `status-${salesApi.getStatusClass(invoice.status)}`]">
                  {{ salesApi.getStatusLabel(invoice.status) }}
                </span>
              </td>
              <td class="amount">{{ salesApi.formatMoney(invoice.total) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="table-summary">
        Mostrando {{ filteredInvoices.length }} de {{ invoices.length }} factura(s)
      </div>
    </div>

    <!-- Create Invoice Modal -->
    <div v-if="showCreateInvoiceModal" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Nueva Factura Estándar</h3>
          <button class="btn-close" @click="closeModal">✕</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <PartySelector
              v-model="invoiceForm.partyId"
              label="Cliente"
              placeholder="Buscar cliente por nombre o referencia..."
              role-filter="CLIENT"
              :required="true"
              help-text="Cliente para el que se emitirá la factura"
            />
          </div>
          <div class="form-group">
            <label>IDs de Pedidos (separados por coma)</label>
            <input
              v-model="invoiceForm.salesOrderIds"
              type="text"
              placeholder="uuid1,uuid2,uuid3"
              class="form-input"
            />
            <span class="help-text">
              Puede dejar vacío si solo factura albaranes
            </span>
          </div>
          <div class="form-group">
            <label>IDs de Albaranes (separados por coma)</label>
            <input
              v-model="invoiceForm.deliveryNoteIds"
              type="text"
              placeholder="uuid1,uuid2,uuid3"
              class="form-input"
            />
          </div>
          <div class="form-group">
            <label>Condiciones de Pago</label>
            <textarea
              v-model="invoiceForm.paymentTerms"
              class="form-textarea"
              rows="2"
              placeholder="Ej: Pago a 30 días..."
            ></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="closeModal">Cancelar</button>
          <button
            class="btn btn-primary"
            @click="createInvoice"
            :disabled="isCreating || !isInvoiceFormValid"
          >
            {{ isCreating ? 'Creando...' : 'Crear Factura' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, watch } from 'vue';
import { useRouter } from 'vue-router';
import Navbar from '@/components/layout/Navbar.vue';
import PartySelector from '@/components/party/PartySelector.vue';
import salesApi from '@/services/salesApi';
import { partyApi } from '@/services/partyApi';

const router = useRouter();

const invoices = ref([]);
const isLoading = ref(false);
const error = ref('');
const partiesCache = ref({});

const filters = ref({
  partyId: '',
  searchText: '',
  type: '',
  fromDate: '',
  toDate: '',
  limit: 50,
});

const filteredInvoices = computed(() => {
  return invoices.value;
});

const isDateRangeValid = computed(() => {
  const hasFromDate = Boolean(filters.value.fromDate);
  const hasToDate = Boolean(filters.value.toDate);
  return hasFromDate === hasToDate;
});

let searchDebounceTimer = null;
let autoFetchEnabled = false;

function scheduleInvoicesFetch() {
  if (searchDebounceTimer) {
    clearTimeout(searchDebounceTimer);
  }

  searchDebounceTimer = setTimeout(() => {
    fetchInvoices();
  }, 300);
}

watch(
  () => filters.value.searchText,
  (newSearch, oldSearch) => {
    const normalizedNew = (newSearch || '').trim();
    const normalizedOld = (oldSearch || '').trim();
    if (normalizedNew === normalizedOld) return;

    scheduleInvoicesFetch();
  },
);

watch(
  () => [filters.value.partyId, filters.value.type],
  (newValues, oldValues) => {
    if (!oldValues) return;
    if (newValues[0] === oldValues[0] && newValues[1] === oldValues[1]) return;

    scheduleInvoicesFetch();
  },
);

watch(
  () => [filters.value.fromDate, filters.value.toDate],
  (newValues, oldValues) => {
    if (!autoFetchEnabled || !oldValues) return;

    const [newFromDate, newToDate] = newValues;
    const [oldFromDate, oldToDate] = oldValues;
    if (newFromDate === oldFromDate && newToDate === oldToDate) return;

    const hasFromDate = Boolean(newFromDate);
    const hasToDate = Boolean(newToDate);
    if (hasFromDate !== hasToDate) return;

    scheduleInvoicesFetch();
  },
);

onBeforeUnmount(() => {
  if (searchDebounceTimer) {
    clearTimeout(searchDebounceTimer);
  }
});

const showCreateInvoiceModal = ref(false);
const isCreating = ref(false);

const invoiceForm = ref({
  partyId: '',
  salesOrderIds: '',
  deliveryNoteIds: '',
  paymentTerms: '',
});

const isInvoiceFormValid = computed(() => {
  return (
    invoiceForm.value.partyId &&
    (invoiceForm.value.salesOrderIds || invoiceForm.value.deliveryNoteIds)
  );
});

onMounted(() => {
  autoFetchEnabled = true;
  fetchInvoices();
});

async function fetchInvoices() {
  isLoading.value = true;
  error.value = '';

  try {
    const params = {};

    if (filters.value.searchText) {
      params.searchText = filters.value.searchText;
    }
    if (filters.value.partyId) {
      params.partyId = filters.value.partyId;
    }
    if (filters.value.type) {
      params.invoiceType = filters.value.type;
    }
    if (filters.value.fromDate) {
      params.fromDate = filters.value.fromDate;
    }
    if (filters.value.toDate) {
      params.toDate = filters.value.toDate;
    }
    if (filters.value.limit) {
      params.limit = filters.value.limit;
    }

    const response = await salesApi.listInvoices(params);

    // Handle both array and object with data property
    invoices.value = Array.isArray(response) ? response : (response.data || []);
    
    // Load party names for display
    await loadPartyNames();
  } catch (err) {
    error.value = err?.message || 'No se pudieron cargar las facturas';
    console.error('Error loading invoices:', err);
  } finally {
    isLoading.value = false;
  }
}

async function loadPartyNames() {
  const partyIds = [...new Set(invoices.value.map(i => i.partyId).filter(Boolean))];
  
  // Filter out already cached parties
  const uncachedIds = partyIds.filter(id => !partiesCache.value[id]);
  
  if (uncachedIds.length === 0) {
    return;
  }

  try {
    const partiesMap = await partyApi.getPartiesBatch(uncachedIds);
    
    // Update cache with batch results
    for (const partyId of uncachedIds) {
      if (partiesMap[partyId]) {
        partiesCache.value[partyId] = partiesMap[partyId].name || 'Desconocido';
      } else {
        partiesCache.value[partyId] = 'No encontrado';
      }
    }
  } catch (err) {
    console.error('Error loading party names:', err);
    // Fallback: mark as error
    for (const partyId of uncachedIds) {
      if (!partiesCache.value[partyId]) {
        partiesCache.value[partyId] = 'Error al cargar';
      }
    }
  }
}

async function createInvoice() {
  if (!isInvoiceFormValid.value || isCreating.value) return;

  isCreating.value = true;

  try {
    const data = {
      partyId: invoiceForm.value.partyId,
      salesOrderIds: invoiceForm.value.salesOrderIds
        ? invoiceForm.value.salesOrderIds.split(',').map((id) => id.trim()).filter(Boolean)
        : [],
      deliveryNoteIds: invoiceForm.value.deliveryNoteIds
        ? invoiceForm.value.deliveryNoteIds.split(',').map((id) => id.trim()).filter(Boolean)
        : [],
    };

    if (invoiceForm.value.paymentTerms) {
      data.paymentTerms = invoiceForm.value.paymentTerms;
    }

    const newInvoice = await salesApi.createInvoice(data);
    router.push(`/sales/invoices/${newInvoice.id}`);
  } catch (err) {
    alert(err?.message || 'No se pudo crear la factura');
  } finally {
    isCreating.value = false;
  }
}

function clearFilters() {
  filters.value = {
    partyId: '',
    searchText: '',
    type: '',
    fromDate: '',
    toDate: '',
    limit: 50,
  };
  fetchInvoices();
}

function applyFilters() {
  if (!isDateRangeValid.value) return;
  fetchInvoices();
}

function closeModal() {
  showCreateInvoiceModal.value = false;
  invoiceForm.value = {
    partyId: '',
    salesOrderIds: '',
    deliveryNoteIds: '',
    paymentTerms: '',
  };
}

function navigateToDetail(invoiceId) {
  router.push(`/sales/invoices/${invoiceId}`);
}

function navigateToCreateTicket() {
  router.push('/sales/tickets/new');
}

function formatDate(dateString) {
  if (!dateString) return '—';
  const date = new Date(dateString);
  return date.toLocaleDateString('es-ES');
}

function formatPartyId(partyId) {
  if (!partyId) return '—';
  return partiesCache.value[partyId] || 'Cargando...';
}

function getTypeLabel(type) {
  const labels = {
    COMPLETA: 'Estándar',
    SIMPLIFICADA: 'Simplificada',
  };
  return labels[type] || type;
}
</script>

<style scoped>
.invoice-list-container {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.page-header h1 {
  font-size: 2rem;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 0.5rem;
}

.filters-card {
  background: white;
  border-radius: 8px;
  padding: 1rem 1.5rem;
  margin-bottom: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-wrap: nowrap;
  gap: 0.5rem;
  align-items: flex-end;
}

.filters-grid {
  display: contents;
}

.filter-group {
  display: flex;
  flex-direction: column;
  flex: 1 1 0;
  min-width: 0;
}

.filter-group label {
  display: block;
  font-size: 0.8rem;
  font-weight: 500;
  color: #4a5568;
  margin-bottom: 0.25rem;
}

.filter-input,
.filter-select {
  width: 100%;
  padding: 0.4rem 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-size: 0.85rem;
}

.filter-input:focus,
.filter-select:focus {
  outline: none;
  border-color: #E6B800;
  box-shadow: 0 0 0 3px rgba(230, 184, 0, 0.1);
}

.filter-actions {
  display: contents;
}

.limit-group {
  display: flex;
  align-items: center;
  gap: 0.35rem;
}

.limit-group label {
  font-size: 0.8rem;
  color: #6b7280;
  white-space: nowrap;
}

.limit-select {
  width: 70px;
}

.limit-label {
  font-size: 0.8rem;
  color: #6b7280;
  white-space: nowrap;
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

.loading-state,
.error-state,
.empty-state {
  text-align: center;
  padding: 3rem 1rem;
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
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

.table-card {
  background: white;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  overflow: hidden;
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

.clickable-row {
  cursor: pointer;
  transition: background 0.2s;
}

.clickable-row:hover {
  background: #f9fafb;
}

.invoice-number {
  font-family: 'Courier New', monospace;
  color: #002395;
  font-weight: 500;
}

.party-id {
  font-family: 'Courier New', monospace;
  color: #6b7280;
}

.amount {
  font-weight: 600;
  text-align: right;
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

.type-completa {
  background: #dbeafe;
  color: #1e40af;
}

.type-simplificada {
  background: #fef3c7;
  color: #92400e;
}

.actions-cell {
  display: flex;
  gap: 0.5rem;
}

.btn-icon {
  background: transparent;
  border: none;
  padding: 0.25rem 0.5rem;
  cursor: pointer;
  font-size: 1.25rem;
  opacity: 0.7;
  transition: opacity 0.2s;
}

.btn-icon:hover {
  opacity: 1;
}

.table-summary {
  padding: 1rem 1.5rem;
  border-top: 1px solid #f3f4f6;
  font-size: 0.875rem;
  color: #6b7280;
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

.help-text {
  display: block;
  font-size: 0.75rem;
  color: #9ca3af;
  margin-top: 0.25rem;
}

.modal-footer {
  display: flex;
  gap: 0.5rem;
  justify-content: flex-end;
  padding: 1.5rem;
  border-top: 1px solid #f3f4f6;
}
</style>
