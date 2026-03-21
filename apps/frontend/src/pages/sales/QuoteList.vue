<template>
  <Navbar />
  <div class="quotes-list-container">
    <!-- Header -->
    <div class="page-header">
      <h1>Presupuestos</h1>
      <button class="btn btn-primary" @click="navigateToCreate">
        <span class="icon">+</span>
        Nuevo Presupuesto
      </button>
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
          <label>Estado</label>
          <select v-model="filters.status" class="filter-select">
            <option value="">Todos</option>
            <option value="DRAFT">Borrador</option>
            <option value="ISSUED">Emitido</option>
            <option value="ACCEPTED">Aceptado</option>
            <option value="REJECTED">Rechazado</option>
            <option value="EXPIRED">Expirado</option>
          </select>
        </div>

        <div class="filter-group">
          <label>Desde</label>
          <input
            v-model="filters.fromDate"
            type="date"
            class="filter-input"
          />
        </div>

        <div class="filter-group">
          <label>Hasta</label>
          <input
            v-model="filters.toDate"
            type="date"
            class="filter-input"
          />
        </div>
      </div>

      <div class="filters-actions">
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
        <button class="btn btn-secondary" @click="clearFilters" v-if="hasFilters">
          Limpiar Filtros
        </button>
        <button class="btn btn-primary" @click="applyFilters" :disabled="!isDateRangeValid" :title="!isDateRangeValid ? 'Completa ambas fechas o vacía ambas para buscar' : ''">
          Buscar
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="loading-state">
      <div class="spinner"></div>
      <p>Cargando presupuestos...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-state">
      <p class="error-message">{{ error }}</p>
      <button class="btn btn-secondary" @click="fetchQuotes">Reintentar</button>
    </div>

    <!-- Empty State -->
    <div v-else-if="filteredQuotes.length === 0" class="empty-state">
      <p>No se encontraron presupuestos</p>
      <button class="btn btn-primary" @click="navigateToCreate">
        Crear Primer Presupuesto
      </button>
    </div>

    <!-- Quotes Table -->
    <div v-else class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>Número</th>
            <th>Cliente</th>
            <th>Fecha</th>
            <th>Válido Hasta</th>
            <th>Estado</th>
            <th>Total</th>
            <th>MES</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="quote in filteredQuotes" :key="quote.id" @click="navigateToDetail(quote.id)" class="clickable-row">
            <td class="quote-number">{{ quote.quoteNumber }}</td>
            <td>{{ formatPartyId(quote.partyId) }}</td>
            <td>{{ formatDate(quote.quoteDate) }}</td>
            <td>{{ formatDate(quote.validUntil) }}</td>
            <td>
              <span :class="['status-badge', `status-${getStatusClass(quote.status)}`]">
                {{ getStatusLabel(quote.status) }}
              </span>
            </td>
            <td>{{ salesApi.formatMoney(quote.total) }}</td>
            <td class="mes-cell" @click.stop>
              <router-link
                v-if="getMesWorkIds(quote.lineItems).length > 0"
                :to="`/sales/quotes/${quote.id}`"
                class="mes-link"
                :title="getMesWorkIds(quote.lineItems).length + ' definición(es) MES'"
              >
                {{ getMesSummary(quote.lineItems) }}
              </router-link>
              <span v-else>—</span>
            </td>
          </tr>
        </tbody>
      </table>
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
import { mesApi } from '@/services/mesApi';

const router = useRouter();

const quotes = ref([]);
const isLoading = ref(false);
const error = ref('');
const partiesCache = ref({});
const mesWorksCache = ref({});

const filters = ref({
  partyId: '',
  searchText: '',
  status: '',
  fromDate: '',
  toDate: '',
  limit: 50,
});

const hasFilters = computed(() => {
  return filters.value.partyId !== '' || 
         filters.value.searchText.trim() !== '' ||
         filters.value.status !== '' || 
         filters.value.fromDate !== '' || 
         filters.value.toDate !== '';
});

const isDateRangeValid = computed(() => {
  const hasFromDate = Boolean(filters.value.fromDate);
  const hasToDate = Boolean(filters.value.toDate);
  return hasFromDate === hasToDate;
});

const filteredQuotes = computed(() => {
  return quotes.value;
});

let searchDebounceTimer = null;
let autoFetchEnabled = false;

function scheduleQuotesFetch() {
  if (searchDebounceTimer) {
    clearTimeout(searchDebounceTimer);
  }

  searchDebounceTimer = setTimeout(() => {
    fetchQuotes();
  }, 300);
}

watch(
  () => filters.value.searchText,
  (newSearch, oldSearch) => {
    const normalizedNew = (newSearch || '').trim();
    const normalizedOld = (oldSearch || '').trim();
    if (normalizedNew === normalizedOld) return;

    scheduleQuotesFetch();
  },
);

watch(
  () => [filters.value.partyId, filters.value.status],
  (newValues, oldValues) => {
    if (!oldValues) return;
    if (newValues[0] === oldValues[0] && newValues[1] === oldValues[1]) return;

    scheduleQuotesFetch();
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

    scheduleQuotesFetch();
  },
);

onBeforeUnmount(() => {
  if (searchDebounceTimer) {
    clearTimeout(searchDebounceTimer);
  }
});

onMounted(() => {
  autoFetchEnabled = true;
  
  fetchQuotes();
});

async function fetchQuotes() {
  isLoading.value = true;
  error.value = '';

  try {
    const apiFilters = {};
    
    if (filters.value.searchText) apiFilters.searchText = filters.value.searchText;
    if (filters.value.partyId) apiFilters.partyId = filters.value.partyId;
    if (filters.value.status) apiFilters.status = filters.value.status;
    if (filters.value.fromDate) apiFilters.fromDate = filters.value.fromDate;
    if (filters.value.toDate) apiFilters.toDate = filters.value.toDate;
    if (filters.value.limit) apiFilters.limit = filters.value.limit;

    const response = await salesApi.listQuotes(apiFilters);
    quotes.value = Array.isArray(response) ? response : (response.data || []);
    
    // Load party names for display
    await loadPartyNames();
    await loadMesWorksForQuotes();
  } catch (err) {
    error.value = err?.message || 'No se pudieron cargar los presupuestos';
    console.error('Error loading quotes:', err);
  } finally {
    isLoading.value = false;
  }
}

async function loadPartyNames() {
  const partyIds = [...new Set(quotes.value.map(q => q.partyId).filter(Boolean))];
  
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

async function loadMesWorksForQuotes() {
  const allLineItems = quotes.value.flatMap((quote) => quote?.lineItems || []);
  const mesWorkIds = [...new Set(
    allLineItems
      .map((item) => item?.mesWorkId)
      .filter((mesWorkId) => typeof mesWorkId === 'string' && mesWorkId.length > 0),
  )];

  const uncachedIds = mesWorkIds.filter((id) => !mesWorksCache.value[id]);
  if (uncachedIds.length === 0) return;

  const results = await Promise.allSettled(uncachedIds.map((id) => mesApi.getWorkOrder(id)));
  results.forEach((result, index) => {
    const mesWorkId = uncachedIds[index];
    mesWorksCache.value[mesWorkId] = result.status === 'fulfilled' ? result.value : null;
  });
}

function applyFilters() {
  if (!isDateRangeValid.value) return;
  fetchQuotes();
}

function clearFilters() {
  filters.value.partyId = '';
  filters.value.searchText = '';
  filters.value.status = '';
  filters.value.fromDate = '';
  filters.value.toDate = '';
  filters.value.limit = 50;
  fetchQuotes();
}

function navigateToCreate() {
  router.push('/sales/quotes/new');
}

function navigateToDetail(quoteId) {
  router.push(`/sales/quotes/${quoteId}`);
}

function canConvertToOrder(status) {
  return status === 'ACCEPTED';
}

async function convertToOrder(quoteId) {
  if (!confirm('¿Convertir este presupuesto en pedido?')) return;

  try {
    const deliveryDate = prompt('Fecha de entrega (YYYY-MM-DD):');
    if (!deliveryDate) return;

    const order = await salesApi.convertQuoteToOrder(quoteId, deliveryDate);
    alert(`Pedido creado: ${order.orderNumber}`);
    router.push(`/sales/orders/${order.id}`);
  } catch (err) {
    alert(err?.message || 'No se pudo convertir el presupuesto');
  }
}

function formatDate(dateString) {
  if (!dateString) return '—';
  const date = new Date(dateString);
  return date.toLocaleDateString('es-ES', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  });
}

function formatPartyId(partyId) {
  if (!partyId) return '—';
  return partiesCache.value[partyId] || 'Cargando...';
}

function getMesSummary(lineItems) {
  if (!Array.isArray(lineItems) || lineItems.length === 0) {
    return '—';
  }

  const uniqueMesWorkIds = [...new Set(
    lineItems
      .map((item) => item?.mesWorkId)
      .filter((mesWorkId) => typeof mesWorkId === 'string' && mesWorkId.length > 0),
  )];

  if (uniqueMesWorkIds.length === 0) {
    return '—';
  }

  if (uniqueMesWorkIds.length === 1) {
    const mesWork = mesWorksCache.value[uniqueMesWorkIds[0]];
    return mesWork?.work_number || '1 ref';
  }

  const firstMesWork = mesWorksCache.value[uniqueMesWorkIds[0]];
  const firstLabel = firstMesWork?.work_number || 'MES';
  return `${firstLabel} +${uniqueMesWorkIds.length - 1}`;
}

function getMesWorkIds(lineItems) {
  if (!Array.isArray(lineItems)) return [];
  return [...new Set(
    lineItems
      .map((item) => item?.mesWorkId)
      .filter((id) => typeof id === 'string' && id.length > 0),
  )];
}

function getStatusClass(status) {
  const classes = {
    DRAFT: 'warning',
    ISSUED: 'info',
    ACCEPTED: 'success',
    REJECTED: 'danger',
    EXPIRED: 'secondary',
    CONVERTED: 'primary',
  };
  return classes[status] || 'secondary';
}

function getStatusLabel(status) {
  const labels = {
    DRAFT: 'Borrador',
    ISSUED: 'Emitido',
    ACCEPTED: 'Aceptado',
    REJECTED: 'Rechazado',
    EXPIRED: 'Expirado',
    CONVERTED: 'Convertido a Pedido',
  };
  return labels[status] || status;
}
</script>

<style scoped>
.quotes-list-container {
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

.icon {
  font-size: 1.25rem;
  margin-right: 0.5rem;
}

.filters-card {
  background: white;
  border-radius: 8px;
  padding: 1rem 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  margin-bottom: 2rem;
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
  font-size: 0.8rem;
  font-weight: 500;
  color: #374151;
  margin-bottom: 0.25rem;
}

.filter-input,
.filter-select {
  padding: 0.4rem 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 4px;
  font-size: 0.85rem;
}

.filters-actions {
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

.loading-state,
.error-state,
.empty-state {
  text-align: center;
  padding: 3rem;
  background: white;
  border-radius: 8px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f4f6;
  border-top-color: #E6B800;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.table-container {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table thead {
  background: #f9fafb;
}

.data-table th {
  padding: 0.75rem 1rem;
  text-align: left;
  font-weight: 600;
  font-size: 0.875rem;
  color: #374151;
  border-bottom: 1px solid #e5e7eb;
}

.data-table td {
  padding: 1rem;
  border-bottom: 1px solid #f3f4f6;
  font-size: 0.875rem;
}

.clickable-row {
  cursor: pointer;
  transition: background 0.2s;
}

.clickable-row:hover {
  background: #f9fafb;
}

.quote-number {
  font-family: 'Courier New', monospace;
  font-weight: 600;
  color: #1e40af;
}

.status-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
}

.status-success {
  background: #d1fae5;
  color: #065f46;
}

.status-warning {
  background: #fef3c7;
  color: #92400e;
}

.status-info {
  background: #dbeafe;
  color: #1e40af;
}

.status-danger {
  background: #fee2e2;
  color: #991b1b;
}

.status-secondary {
  background: #f3f4f6;
  color: #6b7280;
}

.amount {
  font-weight: 600;
  text-align: right;
}

.mes-link {
  color: #1d4ed8;
  text-decoration: none;
  font-weight: 500;
  font-size: 0.8125rem;
}

.mes-link:hover {
  text-decoration: underline;
  color: #1e40af;
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

.btn-primary:hover {
  background: #d4a700;
}

.btn-secondary {
  background: #f3f4f6;
  color: #374151;
}

.btn-secondary:hover {
  background: #e5e7eb;
}
</style>
