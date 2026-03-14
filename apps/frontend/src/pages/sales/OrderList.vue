<template>
  <Navbar />
  <div class="orders-list-container">
    <!-- Header -->
    <div class="page-header">
      <h1>Pedidos</h1>
      <button class="btn btn-primary" @click="navigateToCreate">
        <span class="icon">+</span>
        Nuevo Pedido
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
            <option value="PENDING">Pendiente</option>
            <option value="CONFIRMED">En Preparación</option>
            <option value="PARTIALLY_DELIVERED">Entregado Parcialmente</option>
            <option value="DELIVERED">Entregado</option>
            <option value="CANCELLED">Cancelado</option>
            <option value="PARTIALLY_INVOICED">Facturado Parcialmente</option>
            <option value="INVOICED">Facturado</option>
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
      <p>Cargando pedidos...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-state">
      <p class="error-message">{{ error }}</p>
      <button class="btn btn-secondary" @click="fetchOrders">Reintentar</button>
    </div>

    <!-- Empty State -->
    <div v-else-if="filteredOrders.length === 0" class="empty-state">
      <p>No se encontraron pedidos</p>
      <button class="btn btn-primary" @click="navigateToCreate">
        Crear Primer Pedido
      </button>
    </div>

    <!-- Orders Table -->
    <div v-else class="table-container">
      <table class="data-table">
        <thead>
          <tr>
            <th>Número</th>
            <th>Cliente</th>
            <th>Fecha Pedido</th>
            <th>Fecha Entrega</th>
            <th>Estado</th>
            <th>Total</th>
            <th>MES</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="order in filteredOrders" :key="order.id" @click="navigateToDetail(order.id)" class="clickable-row">
            <td class="order-number">{{ order.orderNumber }}</td>
            <td>{{ formatPartyId(order.partyId) }}</td>
            <td>{{ formatDate(order.orderDate) }}</td>
            <td>{{ formatDate(order.deliveryDate) }}</td>
            <td>
              <span :class="['status-badge', `status-${salesApi.getStatusClass(order.status)}`]">
                {{ salesApi.getStatusLabel(order.status) }}
              </span>
            </td>
            <td class="amount">{{ salesApi.formatMoney(order.total) }}</td>
            <td class="mes-cell" @click.stop>
              <router-link
                v-if="getMesWorkIds(order.lineItems).length > 0"
                :to="`/sales/orders/${order.id}`"
                class="mes-link"
                :title="getMesWorkIds(order.lineItems).length + ' definición(es) MES'"
              >
                {{ getMesSummary(order.lineItems) }}
              </router-link>
              <span v-else>—</span>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Summary -->
      <div class="table-summary">
        <p>Mostrando {{ filteredOrders.length }} de {{ orders.length }} pedido(s)</p>
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
import { mesApi } from '@/services/mesApi';

const router = useRouter();

const orders = ref([]);
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

const filteredOrders = computed(() => {
  return orders.value;
});

let searchDebounceTimer = null;
let autoFetchEnabled = false;

function scheduleOrdersFetch() {
  if (searchDebounceTimer) {
    clearTimeout(searchDebounceTimer);
  }

  searchDebounceTimer = setTimeout(() => {
    fetchOrders();
  }, 300);
}

watch(
  () => filters.value.searchText,
  (newSearch, oldSearch) => {
    const normalizedNew = (newSearch || '').trim();
    const normalizedOld = (oldSearch || '').trim();
    if (normalizedNew === normalizedOld) return;

    scheduleOrdersFetch();
  },
);

watch(
  () => [filters.value.partyId, filters.value.status],
  (newValues, oldValues) => {
    if (!oldValues) return;
    if (newValues[0] === oldValues[0] && newValues[1] === oldValues[1]) return;

    scheduleOrdersFetch();
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

    scheduleOrdersFetch();
  },
);

onBeforeUnmount(() => {
  if (searchDebounceTimer) {
    clearTimeout(searchDebounceTimer);
  }
});

onMounted(() => {
  // Set default date range (last 30 days)
  const today = new Date();
  const thirtyDaysAgo = new Date(today);
  thirtyDaysAgo.setDate(today.getDate() - 30);
  
  filters.value.fromDate = thirtyDaysAgo.toISOString().split('T')[0];
  filters.value.toDate = today.toISOString().split('T')[0];

  autoFetchEnabled = true;
  
  fetchOrders();
});

async function fetchOrders() {
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

    const response = await salesApi.listOrders(apiFilters);
    orders.value = Array.isArray(response) ? response : (response.data || []);
    
    // Load party names for display
    await loadPartyNames();
    await loadMesWorksForOrders();
  } catch (err) {
    error.value = err?.message || 'No se pudieron cargar los pedidos';
    console.error('Error loading orders:', err);
  } finally {
    isLoading.value = false;
  }
}

async function loadPartyNames() {
  const partyIds = [...new Set(orders.value.map(o => o.partyId).filter(Boolean))];
  
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

async function loadMesWorksForOrders() {
  const allLineItems = orders.value.flatMap((order) => order?.lineItems || []);
  const mesWorkIds = [...new Set(
    allLineItems
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

function applyFilters() {
  if (!isDateRangeValid.value) return;
  fetchOrders();
}

function clearFilters() {
  filters.value.partyId = '';
  filters.value.searchText = '';
  filters.value.status = '';
  filters.value.fromDate = '';
  filters.value.toDate = '';
  filters.value.limit = 50;
  fetchOrders();
}

function navigateToCreate() {
  router.push('/sales/orders/new');
}

function navigateToDetail(orderId) {
  router.push(`/sales/orders/${orderId}`);
}

async function confirmOrder(orderId) {
  if (!confirm('¿Confirmar este pedido?')) return;

  try {
    await salesApi.changeOrderStatus(orderId, 'CONFIRMED');
    await fetchOrders();
  } catch (err) {
    alert(err?.message || 'No se pudo confirmar el pedido');
  }
}

async function cancelOrder(orderId) {
  if (!confirm('¿Cancelar este pedido? Esta acción no se puede deshacer.')) return;

  try {
    await salesApi.changeOrderStatus(orderId, 'CANCELLED');
    await fetchOrders();
  } catch (err) {
    alert(err?.message || 'No se pudo cancelar el pedido');
  }
}

function canCancel(status) {
  return ['PENDING', 'CONFIRMED'].includes(status);
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
</script>

<style scoped>
.orders-list-container {
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

.filters-card {
  background: white;
  border-radius: 8px;
  padding: 1rem 1.5rem;
  margin-bottom: 2rem;
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

.btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
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
  color: #4a5568;
}

.btn-secondary:hover {
  background: #e5e7eb;
}

.icon {
  font-size: 1.25rem;
  line-height: 1;
}

.loading-state,
.error-state,
.empty-state {
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
  transition: background-color 0.2s;
}

.clickable-row:hover {
  background-color: #f9fafb;
}

.order-number {
  font-family: 'Courier New', monospace;
  font-weight: 600;
  color: #002395;
}

.amount {
  font-weight: 600;
  text-align: right;
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

.status-primary {
  background: #e0e7ff;
  color: #3730a3;
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

.table-summary {
  padding: 1rem;
  border-top: 1px solid #f3f4f6;
  background: #f9fafb;
  font-size: 0.875rem;
  color: #6b7280;
}
</style>
